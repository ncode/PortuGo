package main

import (
	"errors"
	"strings"
	"testing"
)

func TestReadBoundedRejectsOversizedReader(t *testing.T) {
	_, err := readBounded(strings.NewReader(strings.Repeat("x", maxArtifactBytes+1)), maxArtifactBytes)
	if !errors.Is(err, errArtifactTooLarge) {
		t.Fatalf("error = %v, want oversized artifact", err)
	}
}
