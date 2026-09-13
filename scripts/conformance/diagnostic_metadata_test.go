package main

import (
	"strings"
	"testing"
)

func TestManifestDiagnosticMetadata(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
		for _, tt := range []struct {
			name  string
			diag  diagnostic
			valid bool
		}{
			{"line only", diagnostic{Code: "E004", Line: 4}, true},
			{"exact column", diagnostic{Code: "E004", Line: 4, Column: 2}, true},
			{"lexical code", diagnostic{Code: "L001", Line: 1, Column: 1}, true},
			{"parser code", diagnostic{Code: "P001", Line: 1, Column: 1}, true},
			{"semantic code", diagnostic{Code: "S001", Line: 1, Column: 1}, true},
			{"runtime code", diagnostic{Code: "R008", Line: 1, Column: 1}, true},
			{"empty code", diagnostic{Line: 1}, false},
			{"unsupported code family", diagnostic{Code: "X001", Line: 1}, false},
			{"short code", diagnostic{Code: "E01", Line: 1}, false},
			{"lowercase code", diagnostic{Code: "e001", Line: 1}, false},
			{"nondigit code", diagnostic{Code: "E00x", Line: 1}, false},
			{"code suffix", diagnostic{Code: "E001\n", Line: 1}, false},
			{"zero line", diagnostic{Code: "E001"}, false},
			{"negative line", diagnostic{Code: "E001", Line: -1}, false},
			{"negative column", diagnostic{Code: "E001", Line: 1, Column: -1}, false},
		} {
			t.Run(mode+"/"+tt.name, func(t *testing.T) {
				t.Parallel()
				root, m := testManifest(t)
				p := &m.Probes[0]
				accepted := false
				p.Evidence.Accepted = &accepted
				a := writeArtifact(t, root, "error.txt", "Rejected\n")
				p.Evidence.Raw, p.Evidence.Normalized, p.Evidence.Normalizer = a, a, "bytes-v1"
				p.Implementation.Expected.Stdout = writeArtifact(t, root, "stdout.txt", "")
				p.Implementation.Expected.ExitCode = 1
				p.Implementation.Expected.Diagnostics = []diagnostic{tt.diag}
				if mode == "implementation-acceptance" {
					p.Implementation.State = "verified"
					p.Implementation.Tests = []string{"output_test.go#TestOutput"}
				}
				err := validate(root, m, mode, nil)
				if tt.valid && err != nil {
					t.Fatal(err)
				}
				if !tt.valid && (err == nil || !strings.Contains(err.Error(), "invalid diagnostic expectation")) {
					t.Fatalf("error = %v, want invalid diagnostic expectation", err)
				}
			})
		}
	}
}
