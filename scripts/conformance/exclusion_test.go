package main

import (
	"strings"
	"testing"
)

func TestReviewedExclusionArtifactHashes(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"source.alg", "input.txt", "raw.txt", "output.txt"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			writeArtifact(t, root, "review.md", "Only finite completion is excluded; original evidence remains required.\n")
			m.Probes[0].Evidence.State = "not-applicable"
			m.Probes[0].Evidence.Review = &review{Reason: "Reviewed external-stop exclusion", Link: "review.md"}
			if err := validate(root, m, "evidence", nil); err != nil {
				t.Fatal(err)
			}
			writeArtifact(t, root, "probes/output/"+name, "altered")
			if err := validate(root, m, "evidence", nil); err == nil || !strings.Contains(err.Error(), "hash mismatch") {
				t.Fatalf("error = %v, want retained artifact hash mismatch", err)
			}
		})
	}
}
