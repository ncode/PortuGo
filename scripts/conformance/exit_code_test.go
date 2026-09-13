package main

import (
	"os"
	"strings"
	"testing"
)

func TestExpectedExitCodeDomain(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name     string
		exitCode int
		valid    bool
	}{
		{"success", 0, true},
		{"failure", 1, true},
		{"usage", 2, false},
		{"negative", -1, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			p := &m.Probes[0]
			writeArtifact(t, root, "review.md", "Reviewed exclusion.\n")
			review := &review{Reason: "Reviewed exclusion", Link: "review.md"}
			p.Evidence.State, p.Evidence.Review = "not-applicable", review
			p.Implementation.State, p.Implementation.Review = "not-applicable", review
			p.Implementation.Expected.ExitCode = tt.exitCode
			for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
				t.Run(mode, func(t *testing.T) {
					err := validate(root, m, mode, nil)
					if tt.valid {
						if err != nil {
							t.Fatal(err)
						}
					} else if err == nil || !strings.Contains(err.Error(), "invalid expected exit status") {
						t.Fatalf("validation error = %v", err)
					}
				})
			}
		})
	}
}

func TestReplayRejectsInvalidExpectedExitCode(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	p := m.Probes[0]
	p.Source = writeArtifact(t, root, p.Source.Path, "unknown")
	p.Implementation.Expected.Stdout = writeArtifact(t, root, "empty.txt", "")
	p.Implementation.Expected.ExitCode = 2
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	err = replayProbe(root, p, executable, []string{"-test.run=^TestReplayChild$", "--"}, "")
	if err == nil || !strings.Contains(err.Error(), "invalid expected exit status") {
		t.Fatalf("replay error = %v", err)
	}
}
