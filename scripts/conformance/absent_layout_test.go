package main

import (
	"strings"
	"testing"
)

func TestManifestAbsentLayout(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                      string
		inputs, generated, absent []string
		valid                     bool
	}{
		{name: "generated child", generated: []string{"data/output.dat"}, absent: []string{"data"}},
		{name: "generated descendant", generated: []string{"data/nested/output.dat"}, absent: []string{"data"}},
		{name: "nested absent parent", generated: []string{"data/nested/output.dat"}, absent: []string{"data/nested"}},
		{name: "retained input child", inputs: []string{"data/input.dat"}, absent: []string{"data"}},
		{name: "one removed input", inputs: []string{"data/first.dat", "data/second.dat"}, absent: []string{"data/first.dat", "data"}},
		{name: "removed input and generated child", inputs: []string{"data/input.dat"}, generated: []string{"data/output.dat"}, absent: []string{"data/input.dat", "data"}},
		{name: "generated parent", generated: []string{"data"}, absent: []string{"data/child.dat"}},
		{name: "generated ancestor", generated: []string{"data"}, absent: []string{"data/nested/child.dat"}},
		{name: "retained input parent", inputs: []string{"data"}, absent: []string{"data/child.dat"}},
		{name: "generated sibling", generated: []string{"data/output.dat"}, absent: []string{"data/missing.dat"}, valid: true},
		{name: "similar prefix", generated: []string{"database/output.dat"}, absent: []string{"data"}, valid: true},
		{name: "different parents", inputs: []string{"first/input.dat"}, absent: []string{"second"}, valid: true},
		{name: "removed input and parent", inputs: []string{"data/input.dat"}, absent: []string{"data/input.dat", "data"}, valid: true},
		{name: "removed inputs and parent", inputs: []string{"data/first.dat", "data/second.dat"}, absent: []string{"data", "data/first.dat", "data/second.dat"}, valid: true},
		{name: "exact input removal", inputs: []string{"input.dat"}, absent: []string{"input.dat"}, valid: true},
		{name: "removed input parent", inputs: []string{"data"}, absent: []string{"data", "data/child.dat"}, valid: true},
		{name: "similar file prefix", generated: []string{"data"}, absent: []string{"database/child.dat"}, valid: true},
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
				if mode == "implementation-acceptance" {
					p.Implementation.State = "verified"
					p.Implementation.Tests = []string{"output_test.go#TestOutput"}
				}
				err := validate(root, m, mode, nil)
				if tt.valid {
					if err != nil {
						t.Fatal(err)
					}
				} else if err == nil || !strings.Contains(err.Error(), "absent path") {
					t.Fatalf("error = %v, want required file conflicting with an absent path", err)
				}
			})
		}
	}
}
