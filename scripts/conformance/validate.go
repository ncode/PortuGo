package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func validate(root string, m manifest, mode string, previous *manifest) error {
	var problems []error
	add := func(err error) {
		if err != nil {
			problems = append(problems, err)
		}
	}
	if mode != "evidence" && mode != "incremental" && mode != "implementation-acceptance" {
		return fmt.Errorf("unknown validation mode %q", mode)
	}
	r := m.Reference
	u, urlErr := url.Parse(r.SourceURL)
	_, dateErr := time.Parse("2006-01-02", r.AcquiredDate)
	if r.Product != "VisuAlg" || r.Version != "3.0.7.0" || dateErr != nil || urlErr != nil || u.Host == "" || u.Scheme != "https" || !validHash(r.ArchiveSHA256) || !validHash(r.ExecutableSHA256) || r.OS == "" || r.Culture == "" {
		add(fmt.Errorf("incomplete reference provenance"))
	}
	if m.Version != 1 || m.RecorderVersion != "manual-v1" || m.NormalizerVersion != "panel-v1" {
		add(fmt.Errorf("unsupported manifest or recorder/normalizer version"))
	}
	add(checkProhibited(root, r))
	tasks, completed, err := taskStates(root, m.TasksPath)
	add(err)
	probes := make(map[string]probe)
	for _, p := range m.Probes {
		if p.ID == "" {
			add(fmt.Errorf("empty probe ID"))
		}
		if _, exists := probes[p.ID]; exists {
			add(fmt.Errorf("duplicate probe %s", p.ID))
		}
		probes[p.ID] = p
		if err := validateProbe(root, p, mode, tasks, completed); err != nil {
			add(fmt.Errorf("%s: %w", p.ID, err))
		}
	}
	if len(m.Probes) == 0 {
		add(fmt.Errorf("empty probe inventory"))
	}
	add(validateInventory(root, m, probes))
	retired := make(map[string]bool)
	for _, d := range m.Retired {
		if d.ID == "" || retired[d.ID] {
			add(fmt.Errorf("invalid retired ID %q", d.ID))
		}
		add(checkReview(root, &d.Review))
		retired[d.ID] = true
		if d.Replacement != "" {
			if _, ok := probes[d.Replacement]; !ok {
				add(fmt.Errorf("missing replacement %s", d.Replacement))
			}
		}
	}
	if previous != nil {
		for _, old := range previous.Probes {
			if old.Implementation.State != "verified" {
				continue
			}
			p, exists := probes[old.ID]
			if !exists && !retired[old.ID] {
				add(fmt.Errorf("removed verified probe %s without reviewed disposition", old.ID))
			}
			if exists && p.Implementation.State != "verified" {
				if err := checkReview(root, p.Implementation.Review); err != nil {
					add(fmt.Errorf("unreviewed downgrade of %s: %w", old.ID, err))
				}
			}
		}
	}
	return errors.Join(problems...)
}

func validateProbe(root string, p probe, mode string, tasks map[string]bool, completed map[int]bool) error {
	var problems []error
	add := func(err error) {
		if err != nil {
			problems = append(problems, err)
		}
	}
	if p.OwnerGroup < 1 || len(p.Tasks) == 0 {
		add(fmt.Errorf("missing owner group or task links"))
	}
	owned := false
	for _, task := range p.Tasks {
		if _, ok := tasks[task]; !ok {
			add(fmt.Errorf("stale task link %s", task))
		}
		if strings.HasPrefix(task, fmt.Sprintf("%d.", p.OwnerGroup)) {
			owned = true
		}
	}
	if !owned {
		add(fmt.Errorf("owner group has no linked task"))
	}
	if p.TimeoutMS < 1 || p.TimeoutMS > 30000 {
		add(fmt.Errorf("replay budget must be 1..30000 milliseconds"))
	}
	if p.Evidence.State != "not-applicable" {
		for _, a := range []artifact{p.Source, p.Input} {
			_, err := readArtifact(root, a)
			add(err)
		}
	}
	e := p.Evidence
	switch e.State {
	case "recorded":
		if e.Accepted == nil {
			add(fmt.Errorf("missing reference acceptance disposition"))
		}
		if _, err := time.Parse(time.RFC3339Nano, e.CapturedAt); err != nil {
			add(fmt.Errorf("invalid capture date"))
		}
		raw, err := readArtifact(root, e.Raw)
		add(err)
		normalized, readErr := readArtifact(root, e.Normalized)
		add(readErr)
		if err == nil && readErr == nil {
			want, err := normalize(e.Normalizer, raw)
			add(err)
			if err == nil && !bytes.Equal(want, normalized) {
				add(fmt.Errorf("normalized evidence differs from raw observation"))
			}
		}
		if e.GUIOnly && (e.Screenshot == nil || e.Transcription == nil) {
			add(fmt.Errorf("GUI evidence requires screenshot and transcription"))
		}
		for _, a := range []*artifact{e.Screenshot, e.Transcription} {
			if a != nil {
				_, err := readArtifact(root, *a)
				add(err)
			}
		}
	case "not-applicable":
		add(checkReview(root, e.Review))
	default:
		add(fmt.Errorf("unrecorded evidence"))
	}
	i := p.Implementation
	switch i.State {
	case "verified":
		add(checkTests(root, i.Tests))
	case "pending":
		if mode == "implementation-acceptance" {
			add(fmt.Errorf("pending implementation"))
		}
		if mode != "evidence" && completed[p.OwnerGroup] {
			add(fmt.Errorf("pending behavior owned by completed group %d", p.OwnerGroup))
		}
	case "not-applicable":
		add(checkReview(root, i.Review))
		if e.State != "not-applicable" {
			add(fmt.Errorf("reference acceptance or rejection requires implementation coverage"))
		}
	default:
		add(fmt.Errorf("invalid implementation state %q", i.State))
	}
	for _, link := range i.Tests {
		add(checkLink(root, link, true))
	}
	if e.State == "recorded" {
		out, err := readArtifact(root, i.Expected.Stdout)
		add(err)
		if e.Accepted != nil && *e.Accepted {
			want, err := readArtifact(root, e.Normalized)
			add(err)
			if err == nil && !bytes.Equal(out, want) {
				add(fmt.Errorf("expected stdout differs from reference evidence"))
			}
		}
	}
	for _, file := range append(append([]generatedFile(nil), p.Files...), i.Expected.Generated...) {
		_, err := safePath(root, file.Path)
		add(err)
		_, err = readArtifact(root, file.Content)
		add(err)
	}
	if e.State == "recorded" {
		generated := make(map[string]string)
		for _, file := range e.Generated {
			if _, exists := generated[file.Path]; exists {
				add(fmt.Errorf("duplicate generated reference path %s", file.Path))
			}
			_, err := safePath(root, file.Path)
			add(err)
			_, err = readArtifact(root, file.Content)
			add(err)
			generated[file.Path] = file.Content.SHA256
		}
		if len(e.Generated) != len(i.Expected.Generated) {
			add(fmt.Errorf("generated reference inventory differs from replay expectation"))
		}
		for _, file := range i.Expected.Generated {
			if generated[file.Path] != file.Content.SHA256 {
				add(fmt.Errorf("generated reference differs from replay expectation: %s", file.Path))
			}
		}
	}
	for _, a := range []*artifact{i.Expected.State, i.Expected.HostTrace} {
		if a != nil {
			_, err := readArtifact(root, *a)
			add(err)
		}
	}
	return errors.Join(problems...)
}
