package main

import (
	"strings"
	"testing"
)

func TestManifestPathCase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                      string
		inputs, generated, absent []string
		access                    string
		valid                     bool
	}{
		{name: "input aliases", inputs: []string{"input.dat", "INPUT.DAT"}},
		{name: "nested aliases", inputs: []string{"data/input.dat", "DATA/INPUT.DAT"}},
		{name: "directory aliases", inputs: []string{"data/first.dat", "DATA/second.dat"}},
		{name: "parent before child", inputs: []string{"data", "DATA/input.dat"}},
		{name: "child before parent", inputs: []string{"DATA/input.dat", "data"}},
		{name: "input output aliases", inputs: []string{"input.dat"}, generated: []string{"INPUT.DAT"}},
		{name: "input parent of output", inputs: []string{"data"}, generated: []string{"DATA/output.dat"}},
		{name: "output parent of input", inputs: []string{"DATA/input.dat"}, generated: []string{"data"}},
		{name: "generated aliases", generated: []string{"output.dat", "OUTPUT.DAT"}},
		{name: "absent aliases", absent: []string{"output.dat", "OUTPUT.DAT"}},
		{name: "generated absent aliases", generated: []string{"output.dat"}, absent: []string{"OUTPUT.DAT"}},
		{name: "input absent aliases", inputs: []string{"input.dat"}, absent: []string{"INPUT.DAT"}},
		{name: "access aliases", inputs: []string{"input.dat"}, access: "INPUT.DAT"},
		{name: "input output reuse", inputs: []string{"Data/Input.dat"}, generated: []string{"Data/Input.dat"}, valid: true},
		{name: "consistent directory", inputs: []string{"Data/Input.dat"}, generated: []string{"Data/Output.dat"}, valid: true},
		{name: "removed input", inputs: []string{"Data/Input.dat"}, absent: []string{"Data/Input.dat"}, valid: true},
		{name: "different parents", inputs: []string{"First/input.dat", "Second/INPUT.DAT"}, valid: true},
		{name: "similar prefix", inputs: []string{"Data/input.dat", "Database/INPUT.DAT"}, valid: true},
		{name: "exact access", inputs: []string{"Data/Input.dat"}, access: "Data/Input.dat", valid: true},
	}
	for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
		for _, tt := range tests {
			t.Run(mode+"/"+tt.name, func(t *testing.T) {
				t.Parallel()
				root, m := testManifest(t)
				content := writeArtifact(t, root, "fixture.txt", "Fixture bytes.\n")
				p := &m.Probes[0]
				for _, name := range tt.inputs {
					p.Files = append(p.Files, generatedFile{Path: name, Content: content})
				}
				for _, name := range tt.generated {
					p.Evidence.Generated = append(p.Evidence.Generated, generatedFile{Path: name, Content: content})
				}
				p.Implementation.Expected.Generated = p.Evidence.Generated
				p.Evidence.Absent, p.Implementation.Expected.Absent = tt.absent, tt.absent
				if tt.access != "" {
					p.FixtureAccess = &fixtureAccess{Path: tt.access, Mode: "unavailable"}
				}
				if mode == "implementation-acceptance" {
					p.Implementation.State = "verified"
					p.Implementation.Tests = []string{"output_test.go#TestOutput"}
				}
				err := validate(root, m, mode, nil)
				if tt.valid {
					if err != nil {
						t.Fatal(err)
					}
				} else if err == nil || !strings.Contains(err.Error(), "case-insensitive path collision") {
					t.Fatalf("error = %v, want case-insensitive path collision", err)
				}
			})
		}
	}
}
