package main

import (
	"strings"
	"testing"
)

func TestManifestInputInventory(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
		for _, name := range []string{"distinct paths", "shared content", "duplicate content", "conflicting content", "input output reuse"} {
			t.Run(mode+"/"+name, func(t *testing.T) {
				t.Parallel()
				root, m := testManifest(t)
				first := writeArtifact(t, root, "first.txt", "First bytes.\n")
				second := writeArtifact(t, root, "second.txt", "Second bytes.\n")
				p := &m.Probes[0]
				p.Files = []generatedFile{{Path: "first.dat", Content: first}, {Path: "second.dat", Content: second}}
				duplicate := false
				switch name {
				case "shared content":
					p.Files[1].Content = first
				case "duplicate content":
					p.Files[1] = p.Files[0]
					duplicate = true
				case "conflicting content":
					p.Files[1].Path = p.Files[0].Path
					duplicate = true
				case "input output reuse":
					p.Evidence.Generated = []generatedFile{{Path: "first.dat", Content: second}}
					p.Implementation.Expected.Generated = p.Evidence.Generated
				}
				if mode == "implementation-acceptance" {
					p.Implementation.State = "verified"
					p.Implementation.Tests = []string{"output_test.go#TestOutput"}
				}
				err := validate(root, m, mode, nil)
				if duplicate {
					if err == nil || !strings.Contains(err.Error(), "duplicate input path") {
						t.Fatalf("error = %v, want duplicate input path", err)
					}
				} else if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
