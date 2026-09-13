package main

import (
	"strings"
	"testing"
)

func TestManifestRejectsOversizedObservation(t *testing.T) {
	for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
		t.Run(mode, func(t *testing.T) {
			root, m := testManifest(t)
			p := &m.Probes[0]
			p.Implementation.State = "verified"
			p.Implementation.Tests = []string{"output_test.go#TestOutput"}
			observation := writeArtifact(t, root, "state.json", strings.Repeat("x", maxObservationBytes+1))
			p.Implementation.Expected.State = &observation
			if err := validate(root, m, mode, nil); err == nil || !strings.Contains(err.Error(), "replay output exceeds size limit") {
				t.Fatalf("error = %v, want oversized observation failure", err)
			}
		})
	}
}

func TestReplayRejectsOversizedObservation(t *testing.T) {
	root, m := testManifest(t)
	observation := writeArtifact(t, root, "state.json", strings.Repeat("x", maxObservationBytes+1))
	m.Probes[0].Implementation.Expected.State = &observation
	if err := replayProbe(root, m.Probes[0], "unused", nil, "unused"); err == nil || !strings.Contains(err.Error(), "replay output exceeds size limit") {
		t.Fatalf("error = %v, want oversized observation failure", err)
	}
}

func TestManifestRejectsOversizedReplayOutput(t *testing.T) {
	for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
		t.Run(mode, func(t *testing.T) {
			root, m := testManifest(t)
			p := &m.Probes[0]
			p.Implementation.State = "verified"
			p.Implementation.Tests = []string{"output_test.go#TestOutput"}
			p.Implementation.Expected.Stdout = writeArtifact(t, root, "stdout.txt", strings.Repeat("x", maxObservationBytes+1))
			if err := validate(root, m, mode, nil); err == nil || !strings.Contains(err.Error(), "replay output exceeds size limit") {
				t.Fatalf("error = %v, want oversized output failure", err)
			}
		})
	}
}

func TestReplayRejectsOversizedReplayOutput(t *testing.T) {
	root, m := testManifest(t)
	m.Probes[0].Implementation.Expected.Stdout = writeArtifact(t, root, "stdout.txt", strings.Repeat("x", maxObservationBytes+1))
	if err := replayProbe(root, m.Probes[0], "unused", nil, ""); err == nil || !strings.Contains(err.Error(), "replay output exceeds size limit") {
		t.Fatalf("error = %v, want oversized output failure", err)
	}
}
