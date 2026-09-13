package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestCommandErrors(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{nil, {"unknown"}, {"validate", "--unknown"}, {"stage"}, {"capture"}} {
		var out, stderr bytes.Buffer
		if status := run(args, &out, &stderr); status != 2 {
			t.Errorf("%v: status=%d, stderr=%s", args, status, &stderr)
		}
		if stderr.Len() == 0 {
			t.Errorf("%v: missing error", args)
		}
	}
}

func TestCommandRejectsInvalidManifest(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	m.Reference.ExecutableSHA256 = ""
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, root, "manifest.json", string(data))
	var out, stderr bytes.Buffer
	if status := run([]string{"validate", "--root", root, "--manifest", "manifest.json", "--previous", "manifest.json"}, &out, &stderr); status != 1 {
		t.Fatalf("status=%d, stderr=%s", status, &stderr)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("reference provenance")) {
		t.Fatal(stderr.String())
	}
}

func TestCommandHistoryValidation(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	m.Probes[0].Implementation.State = "verified"
	m.Probes[0].Implementation.Tests = []string{"output_test.go#TestOutput"}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, root, "manifest.json", string(data))
	writeArtifact(t, root, "previous.json", string(data))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, args := range [][]string{{"init"}, {"add", "."}, {"-c", "user.name=Test", "-c", "user.email=test@example.invalid", "-c", "commit.gpgsign=false", "-c", "core.hooksPath=.git/no-hooks", "commit", "-m", "Initial corpus"}} {
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = root
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	m.Probes[0].Implementation.State = "pending"
	data, err = json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, root, "manifest.json", string(data))
	for _, tt := range []struct {
		name, mode, message string
		history             []string
		status              int
	}{
		{name: "evidence requires history", mode: "evidence", status: 2, message: "requires --base or --previous"},
		{name: "incremental requires history", mode: "incremental", status: 2, message: "requires --base or --previous"},
		{name: "acceptance requires history", mode: "implementation-acceptance", status: 2, message: "requires --base or --previous"},
		{name: "base catches downgrade", mode: "evidence", history: []string{"--base", "HEAD"}, status: 1, message: "unreviewed downgrade"},
		{name: "file catches downgrade", mode: "evidence", history: []string{"--previous", "previous.json"}, status: 1, message: "unreviewed downgrade"},
		{name: "invalid base fails closed", mode: "evidence", history: []string{"--base", "missing-ref"}, status: 1, message: "resolve previous manifest base"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out, stderr bytes.Buffer
			args := []string{"validate", "--root", root, "--manifest", "manifest.json", "--mode", tt.mode, "--candidate", "must-not-run"}
			status := run(append(args, tt.history...), &out, &stderr)
			if status != tt.status || !bytes.Contains(stderr.Bytes(), []byte(tt.message)) {
				t.Fatalf("status=%d, stderr=%s; want status=%d, message=%q", status, &stderr, tt.status, tt.message)
			}
		})
	}
	if old, err := previousManifest(root, "HEAD", "new-corpus.json"); err != nil || old != nil {
		t.Fatalf("initial corpus: previous=%+v, error=%v", old, err)
	}
}

func TestLoadRejectsUnknownFieldsAndTrailingJSON(t *testing.T) {
	t.Parallel()
	for _, content := range []string{`{"versoin":1}`, `{"version":1} {"version":2}`} {
		root := t.TempDir()
		writeArtifact(t, root, "manifest.json", content)
		if _, err := loadManifest(root, "manifest.json"); err == nil {
			t.Errorf("accepted %s", content)
		}
	}
}

func TestPreviousManifestRejectsLooseJSON(t *testing.T) {
	for _, tt := range []struct {
		name, content, want string
	}{
		{name: "unknown field", content: `{"versoin":1}`, want: "unknown field"},
		{name: "trailing JSON", content: `{"version":1} {"version":2}`, want: "trailing manifest JSON"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeArtifact(t, root, "manifest.json", tt.content)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			for _, args := range [][]string{{"init"}, {"add", "manifest.json"}, {"-c", "user.name=Test", "-c", "user.email=test@example.invalid", "-c", "commit.gpgsign=false", "commit", "-m", "Initial manifest"}} {
				cmd := exec.CommandContext(ctx, "git", args...)
				cmd.Dir = root
				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("git %v: %v: %s", args, err, output)
				}
			}
			if _, err := previousManifest(root, "HEAD", "manifest.json"); err == nil || !bytes.Contains([]byte(err.Error()), []byte(tt.want)) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestPreviousManifestRejectsOversizedJSON(t *testing.T) {
	root := t.TempDir()
	writeArtifact(t, root, "manifest.json", strings.Repeat("x", maxArtifactBytes+1))
	if _, err := loadManifest(root, "manifest.json"); err == nil || !strings.Contains(err.Error(), "invalid file size or type") {
		t.Fatalf("current manifest error = %v, want size limit", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, args := range [][]string{{"init"}, {"add", "manifest.json"}, {"-c", "user.name=Test", "-c", "user.email=test@example.invalid", "-c", "commit.gpgsign=false", "commit", "-m", "Initial manifest"}} {
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = root
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	if _, err := previousManifest(root, "HEAD", "manifest.json"); err == nil || !strings.Contains(err.Error(), "previous manifest exceeds artifact size limit") {
		t.Fatalf("error = %v, want size limit", err)
	}
}

func TestLoadManifestAcceptsFixtureAccess(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeArtifact(t, root, "manifest.json", `{"probes":[{"fixtureAccess":{"path":"input.dat","mode":"unavailable"}}]}`)
	if _, err := loadManifest(root, "manifest.json"); err != nil {
		t.Fatalf("load fixture access manifest: %v", err)
	}
}

func TestCommandStagesSelectedProbe(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, root, "manifest.json", string(data))
	stage := t.TempDir() + string(os.PathSeparator) + "recording"
	var out, stderr bytes.Buffer
	if status := run([]string{"stage", "--root", root, "--manifest", "manifest.json", "--probe", "output", "--staging", stage}, &out, &stderr); status != 0 {
		t.Fatalf("status=%d, stderr=%s", status, &stderr)
	}
}
