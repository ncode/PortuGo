package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckProhibitedSkipsRepositoryMetadata(t *testing.T) {
	for _, kind := range []string{"file", "directory"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			metadata := filepath.Join(root, ".git")
			if kind == "file" {
				if err := os.WriteFile(metadata, []byte("metadata\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			} else if err := os.Mkdir(metadata, 0o700); err != nil {
				t.Fatal(err)
			}
			if err := checkProhibited(root, reference{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestManifestPathComponents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, path string
		valid      bool
	}{
		{name: "file period", path: "input.dat."},
		{name: "file space", path: "input.dat "},
		{name: "directory period", path: "data./input.dat"},
		{name: "directory space", path: "data /input.dat"},
		{name: "parent alias", path: ".. /input.dat"},
		{name: "only periods", path: "data/.../input.dat"},
		{name: "repository metadata", path: ".git/HEAD"},
		{name: "repository metadata alias", path: ".GIT/config"},
		{name: "plain", path: "data/input.dat", valid: true},
		{name: "leading periods", path: ".data/.input.dat", valid: true},
		{name: "interior periods and spaces", path: "data.. files/input value.dat", valid: true},
		{name: "nonbreaking space", path: "input.dat\u00a0", valid: true},
	}
	for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
		for _, tt := range tests {
			t.Run(mode+"/"+tt.name, func(t *testing.T) {
				t.Parallel()
				root, m := testManifest(t)
				content := writeArtifact(t, root, "fixture.txt", "Input bytes.\n")
				p := &m.Probes[0]
				p.Files = []generatedFile{{Path: tt.path, Content: content}}
				if mode == "implementation-acceptance" {
					p.Implementation.State = "verified"
					p.Implementation.Tests = []string{"output_test.go#TestOutput"}
				}
				err := validate(root, m, mode, nil)
				if tt.valid {
					if err != nil {
						t.Fatal(err)
					}
				} else if err == nil || !strings.Contains(err.Error(), "unsafe path") {
					t.Fatalf("error = %v, want unsafe path", err)
				}
			})
		}
	}
}
