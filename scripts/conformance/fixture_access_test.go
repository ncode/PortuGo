package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareFixtureAccessMakesExistingFileUnavailable(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := filepath.Join(root, "input.dat")
	want := []byte("unchanged\r\n")
	if err := os.WriteFile(path, want, 0o600); err != nil {
		t.Fatal(err)
	}
	restore, err := prepareFixtureAccess(root, &fixtureAccess{Path: "input.dat", Mode: "unavailable"})
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err == nil {
		_ = f.Close()
		t.Fatal("fixture remained readable")
	}
	if err := restore(); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("fixture bytes = %q, want %q", got, want)
	}
}
