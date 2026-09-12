// Command conformance validates, replays, and stages VisuAlg reference evidence.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func loadManifest(root, name string) (manifest, error) {
	var m manifest
	data, err := readFile(root, name)
	if err != nil {
		return m, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&m); err != nil {
		return m, err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return m, fmt.Errorf("trailing manifest JSON")
	}
	return m, nil
}

func run(args []string, out, stderr io.Writer) (status int) {
	if len(args) > 0 && args[0] == "execute" {
		return executeProbe(args[1:], os.Stdin, out, stderr)
	}
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "usage: conformance validate|stage|capture [options]")
		return 2
	}
	if args[0] != "validate" && args[0] != "stage" && args[0] != "capture" {
		_, _ = fmt.Fprintln(stderr, "unknown conformance command")
		return 2
	}
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("root", ".", "repository root")
	name := flags.String("manifest", "testdata/conformance/visualg-3.0.7/manifest.json", "manifest path relative to root")
	mode := flags.String("mode", "evidence", "evidence, incremental, or implementation-acceptance")
	candidate := flags.String("candidate", "", "CLI executable; default builds the current checkout")
	previous := flags.String("previous", "", "previous manifest path relative to root, for downgrade validation")
	base := flags.String("base", "", "Git ref containing the previous manifest, for downgrade validation")
	id := flags.String("probe", "", "probe ID to stage")
	stage := flags.String("staging", "", "private recording directory outside repository")
	accepted := flags.String("accepted", "", "recorded reference acceptance: true or false")
	capturedAt := flags.String("captured-at", "", "observed UTC capture time in RFC3339 format")
	normalizer := flags.String("normalizer", "text-v1", "panel-v1, text-v1, or bytes-v1")
	guiOnly := flags.Bool("gui-only", false, "require screenshot.png and transcription.txt")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		_, _ = fmt.Fprintln(stderr, "unexpected positional arguments")
		return 2
	}
	// A failed error stream cannot change the already failing exit status.
	fail := func(err error) int { _, _ = fmt.Fprintln(stderr, err); return 1 }
	if args[0] == "capture" {
		accepts, err := strconv.ParseBool(*accepted)
		if *stage == "" || *capturedAt == "" || err != nil {
			_, _ = fmt.Fprintln(stderr, "capture requires --staging, --accepted and --captured-at")
			return 2
		}
		e, err := captureRecording(*stage, accepts, *capturedAt, *normalizer, *guiOnly)
		if err != nil {
			return fail(err)
		}
		if err := json.NewEncoder(out).Encode(e); err != nil {
			return fail(err)
		}
		return 0
	}
	if args[0] == "stage" && (*id == "" || *stage == "") {
		_, _ = fmt.Fprintln(stderr, "stage requires --probe and --staging")
		return 2
	}
	rootPath, err := filepath.Abs(*root)
	if err != nil {
		return fail(err)
	}
	m, err := loadManifest(rootPath, *name)
	if err != nil {
		return fail(err)
	}
	if args[0] == "stage" {
		for _, p := range m.Probes {
			if p.ID == *id {
				if err := prepareRecording(rootPath, p, *stage); err != nil {
					return fail(err)
				}
				if _, err := fmt.Fprintln(out, "Private recording staged; follow instructions.txt."); err != nil {
					return fail(err)
				}
				return 0
			}
		}
		return fail(fmt.Errorf("unknown probe %q", *id))
	}
	if *base != "" && *previous != "" {
		_, _ = fmt.Fprintln(stderr, "choose --base or --previous")
		return 2
	}
	if *base == "" && *previous == "" {
		_, _ = fmt.Fprintln(stderr, "validate requires --base or --previous for history validation")
		return 2
	}
	var old *manifest
	if *previous != "" {
		value, err := loadManifest(rootPath, *previous)
		if err != nil {
			return fail(err)
		}
		old = &value
	}
	if *base != "" {
		value, err := previousManifest(rootPath, *base, *name)
		if err != nil {
			return fail(err)
		}
		old = value
	}
	if err := validate(rootPath, m, *mode, old); err != nil {
		return fail(err)
	}
	executable := *candidate
	needsObserver := false
	for _, p := range m.Probes {
		needsObserver = needsObserver || p.Implementation.Expected.State != nil || p.Implementation.Expected.HostTrace != nil
	}
	if executable != "" && needsObserver {
		return fail(fmt.Errorf("state/host observations require the default current-checkout candidate"))
	}
	observer := ""
	if executable == "" {
		dir, err := os.MkdirTemp("", "portugol-candidate-")
		if err != nil {
			return fail(err)
		}
		defer func() {
			if err := os.RemoveAll(dir); err != nil {
				status = fail(fmt.Errorf("clean candidate directory: %w", err))
			}
		}()
		executable = filepath.Join(dir, "portugol.exe")
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, "go", "build", "-o", executable, "./cmd/portugol")
		cmd.Dir, cmd.Stdout, cmd.Stderr = rootPath, stderr, stderr
		if err := cmd.Run(); err != nil {
			return fail(fmt.Errorf("build candidate: %w", err))
		}
		if needsObserver {
			observer = filepath.Join(dir, "observations.exe")
			cmd := exec.CommandContext(ctx, "go", "build", "-o", observer, "./scripts/conformance")
			cmd.Dir, cmd.Stdout, cmd.Stderr = rootPath, stderr, stderr
			if err := cmd.Run(); err != nil {
				return fail(fmt.Errorf("build observation adapter: %w", err))
			}
		}
	} else {
		executable, err = filepath.Abs(executable)
		if err != nil {
			return fail(err)
		}
	}
	results := make([]replayResult, 0, len(m.Probes))
	for _, p := range m.Probes {
		r := replayResult{ID: p.ID, State: p.Implementation.State}
		if p.Evidence.State == "recorded" && p.Implementation.State != "not-applicable" {
			if err := replayProbe(rootPath, p, executable, nil, observer); err != nil {
				r.Error = err.Error()
			}
		}
		results = append(results, r)
	}
	if err := json.NewEncoder(out).Encode(results); err != nil {
		return fail(err)
	}
	if replayFailed(results) {
		return fail(fmt.Errorf("verified reference regression"))
	}
	return 0
}

func previousManifest(root, base, name string) (*manifest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resolve := exec.CommandContext(ctx, "git", "rev-parse", "--verify", "--end-of-options", base+"^{commit}")
	resolve.Dir = root
	sha, err := resolve.Output()
	if err != nil {
		return nil, fmt.Errorf("resolve previous manifest base: %w", err)
	}
	object := string(bytes.TrimSpace(sha)) + ":" + filepath.ToSlash(name)
	cmd := exec.CommandContext(ctx, "git", "show", object)
	cmd.Dir = root
	data, err := cmd.Output()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			// A new corpus has no previous manifest. Confirm absence in the tree;
			// other Git errors must not silently disable downgrade checking.
			list := exec.CommandContext(ctx, "git", "ls-tree", "--name-only", string(bytes.TrimSpace(sha)), "--", name)
			list.Dir = root
			paths, listErr := list.Output()
			if listErr == nil && len(bytes.TrimSpace(paths)) == 0 {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("read previous manifest: %w", err)
	}
	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
