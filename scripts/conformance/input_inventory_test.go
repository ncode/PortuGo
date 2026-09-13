package main

import (
	"strings"
	"testing"
)

func TestManifestInputInventory(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
		for _, name := range []string{"distinct paths", "shared content", "duplicate content", "conflicting content", "input output reuse", "parent before child", "child before parent", "deep parent", "shared directory", "similar prefix", "absolute path", "repeated slash", "input parent of output", "output parent of input", "output parent child", "shared input output directory", "similar input output prefix", "removed input below output", "removed input above output", "removed conflicting input"} {
			t.Run(mode+"/"+name, func(t *testing.T) {
				t.Parallel()
				root, m := testManifest(t)
				first := writeArtifact(t, root, "first.txt", "First bytes.\n")
				second := writeArtifact(t, root, "second.txt", "Second bytes.\n")
				p := &m.Probes[0]
				p.Files = []generatedFile{{Path: "first.dat", Content: first}, {Path: "second.dat", Content: second}}
				want := ""
				switch name {
				case "shared content":
					p.Files[1].Content = first
				case "duplicate content":
					p.Files[1] = p.Files[0]
					want = "duplicate input path"
				case "conflicting content":
					p.Files[1].Path = p.Files[0].Path
					want = "duplicate input path"
				case "parent before child":
					p.Files[1].Path = "first.dat/child.dat"
					want = "file path is both file and directory"
				case "child before parent":
					p.Files[0].Path = "second.dat/child.dat"
					want = "file path is both file and directory"
				case "deep parent":
					p.Files[1].Path = "first.dat/nested/child.dat"
					want = "file path is both file and directory"
				case "shared directory":
					p.Files[0].Path, p.Files[1].Path = "data/first.dat", "data/second.dat"
				case "similar prefix":
					p.Files[0].Path, p.Files[1].Path = "data", "database/second.dat"
				case "absolute path":
					p.Files[1].Path = "/"
					want = "unsafe path"
				case "repeated slash":
					p.Files[1].Path = "data//child.dat"
					want = "unsafe path"
				case "input parent of output":
					p.Evidence.Generated = []generatedFile{{Path: "first.dat/result.dat", Content: second}}
					p.Implementation.Expected.Generated = p.Evidence.Generated
					want = "file path is both file and directory"
				case "output parent of input":
					p.Files[0].Path = "data/input.dat"
					p.Evidence.Generated = []generatedFile{{Path: "data", Content: second}}
					p.Implementation.Expected.Generated = p.Evidence.Generated
					want = "file path is both file and directory"
				case "output parent child":
					p.Evidence.Generated = []generatedFile{{Path: "result.dat", Content: first}, {Path: "result.dat/child.dat", Content: second}}
					p.Implementation.Expected.Generated = p.Evidence.Generated
					want = "file path is both file and directory"
				case "shared input output directory":
					p.Files[0].Path = "data/input.dat"
					p.Evidence.Generated = []generatedFile{{Path: "data/result.dat", Content: second}}
					p.Implementation.Expected.Generated = p.Evidence.Generated
				case "similar input output prefix":
					p.Files[0].Path = "data"
					p.Evidence.Generated = []generatedFile{{Path: "database/result.dat", Content: second}}
					p.Implementation.Expected.Generated = p.Evidence.Generated
				case "removed input below output":
					p.Files[0].Path = "data/input.dat"
					p.Evidence.Generated = []generatedFile{{Path: "data", Content: second}}
					p.Implementation.Expected.Generated = p.Evidence.Generated
					p.Evidence.Absent = []string{"data/input.dat"}
					p.Implementation.Expected.Absent = p.Evidence.Absent
				case "removed input above output":
					p.Files[0].Path = "data"
					p.Evidence.Generated = []generatedFile{{Path: "data/result.dat", Content: second}}
					p.Implementation.Expected.Generated = p.Evidence.Generated
					p.Evidence.Absent = []string{"data"}
					p.Implementation.Expected.Absent = p.Evidence.Absent
					want = "file path is both file and directory"
				case "removed conflicting input":
					p.Files[0].Path, p.Files[1].Path = "data", "data/child.dat"
					p.Evidence.Absent = []string{"data/child.dat"}
					p.Implementation.Expected.Absent = p.Evidence.Absent
					want = "file path is both file and directory"
				case "input output reuse":
					p.Evidence.Generated = []generatedFile{{Path: "first.dat", Content: second}}
					p.Implementation.Expected.Generated = p.Evidence.Generated
				}
				if mode == "implementation-acceptance" {
					p.Implementation.State = "verified"
					p.Implementation.Tests = []string{"output_test.go#TestOutput"}
				}
				err := validate(root, m, mode, nil)
				if want != "" {
					if err == nil || !strings.Contains(err.Error(), want) {
						t.Fatalf("error = %v, want %q", err, want)
					}
				} else if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
