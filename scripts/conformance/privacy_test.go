package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublicErrorRedactsFilesystemDetails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "private-name")
	for _, input := range []error{
		fmt.Errorf("start: %w", &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}),
		fmt.Errorf("start: %w", &exec.Error{Name: path, Err: os.ErrNotExist}),
	} {
		err := publicError(input)
		if err == nil || err.Error() != "filesystem operation failed" {
			t.Fatalf("error = %v, want generic filesystem failure", err)
		}
		if strings.Contains(err.Error(), path) {
			t.Fatalf("error leaked filesystem path: %v", err)
		}
	}
}

func TestRunRedactsFilesystemFailure(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing-root")
	var out, stderr strings.Builder
	status := run([]string{"validate", "--root", root, "--manifest", "manifest.json", "--previous", "manifest.json"}, &out, &stderr)
	if status != 1 || stderr.String() != "filesystem operation failed\n" {
		t.Fatalf("status = %d, stderr = %q; want generic filesystem failure", status, stderr.String())
	}
	if strings.Contains(stderr.String(), root) {
		t.Fatalf("stderr leaked filesystem path: %s", stderr.String())
	}
}

func TestRunRedactsReplayFilesystemFailure(t *testing.T) {
	root, m := testManifest(t)
	writeJSON(t, root, "manifest.json", m)
	candidate := filepath.Join(t.TempDir(), "missing-candidate")
	var out, stderr strings.Builder
	status := run([]string{"validate", "--root", root, "--manifest", "manifest.json", "--previous", "manifest.json", "--candidate", candidate}, &out, &stderr)
	if status != 0 || !strings.Contains(out.String(), `"error":"filesystem operation failed"`) {
		t.Fatalf("status = %d, stdout = %q, stderr = %q; want sanitized replay failure", status, out.String(), stderr.String())
	}
	if strings.Contains(out.String(), candidate) {
		t.Fatalf("stdout leaked executable path: %s", out.String())
	}
}

func TestExecuteProbeRedactsFilesystemFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source.alg")
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatal(err)
	}
	var out, stderr strings.Builder
	status := executeProbe([]string{path}, strings.NewReader(""), &out, &stderr)
	if status != 1 || !strings.Contains(stderr.String(), "invalid file size or type") || strings.Contains(stderr.String(), path) {
		t.Fatalf("status = %d, stderr = %q; want sanitized filesystem failure", status, stderr.String())
	}
}
