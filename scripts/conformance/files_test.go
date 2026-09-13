package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadBoundedRejectsOversizedReader(t *testing.T) {
	_, err := readBounded(strings.NewReader(strings.Repeat("x", maxArtifactBytes+1)), maxArtifactBytes)
	if !errors.Is(err, errArtifactTooLarge) {
		t.Fatalf("error = %v, want oversized artifact", err)
	}
}

func TestReadFileLimitRejectsOversizedFile(t *testing.T) {
	root := t.TempDir()
	name := filepath.Join(root, "output.txt")
	if err := os.WriteFile(name, []byte(strings.Repeat("x", maxObservationBytes+1)), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := readFileLimit(root, "output.txt", maxObservationBytes)
	if !errors.Is(err, errArtifactTooLarge) {
		t.Fatalf("error = %v, want oversized artifact", err)
	}
}

func TestReadReplaySourceArtifactRejectsOversizedFile(t *testing.T) {
	root := t.TempDir()
	data := []byte(strings.Repeat("x", 64<<10+1))
	if err := os.WriteFile(filepath.Join(root, "source.alg"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := readReplaySourceArtifact(root, artifact{Path: "source.alg", SHA256: hashBytes(data)})
	if err == nil || err.Error() != "source exceeds replay profile" {
		t.Fatalf("error = %v, want source profile limit", err)
	}
}
