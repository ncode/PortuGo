package main

import (
	"strings"
	"testing"
)

func TestManifestDowngradeRequiresNewReview(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
		for _, tt := range []struct {
			name, state string
			prior, next *review
			unreviewed  bool
		}{
			{"verified unchanged", "verified", &review{"Qualified output", "review.md"}, &review{"Qualified output", "review.md"}, false},
			{"missing review", "pending", &review{"Qualified output", "review.md"}, nil, true},
			{"reused review", "pending", &review{"Qualified output", "review.md"}, &review{"Qualified output", "review.md"}, true},
			{"whitespace change", "pending", &review{"Qualified output", "review.md"}, &review{" Qualified output \n", "review.md"}, true},
			{"new reason", "pending", &review{"Qualified output", "review.md"}, &review{"Reopened after a discovered mismatch", "review.md"}, false},
			{"new review link", "pending", &review{"Reviewed scope", "review.md"}, &review{"Reviewed scope", "correction.md"}, false},
			{"first review", "pending", nil, &review{"Reopened after a discovered mismatch", "correction.md"}, false},
		} {
			t.Run(mode+"/"+tt.name, func(t *testing.T) {
				t.Parallel()
				root, m := testManifest(t)
				writeArtifact(t, root, "review.md", "Reviewed output qualification.\n")
				writeArtifact(t, root, "correction.md", "Reviewed scope correction after a discovered mismatch.\n")
				m.Probes[0].Implementation.State = "verified"
				m.Probes[0].Implementation.Tests = []string{"output_test.go#TestOutput"}
				m.Probes[0].Implementation.Review = tt.prior
				previous := m
				previous.Probes = append([]probe(nil), m.Probes...)
				m.Probes[0].Implementation.State = tt.state
				m.Probes[0].Implementation.Review = tt.next
				err := validate(root, m, mode, &previous)
				unreviewed := err != nil && strings.Contains(err.Error(), "unreviewed downgrade")
				if unreviewed != tt.unreviewed {
					t.Fatalf("error = %v, want unreviewed downgrade = %t", err, tt.unreviewed)
				}
				if !tt.unreviewed {
					pendingAcceptance := mode == "implementation-acceptance" && tt.state == "pending"
					if pendingAcceptance {
						if err == nil || !strings.Contains(err.Error(), "pending implementation") {
							t.Fatalf("error = %v, want pending implementation", err)
						}
					} else if err != nil {
						t.Fatal(err)
					}
				}
			})
		}
	}
}
