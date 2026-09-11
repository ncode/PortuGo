package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedRandomBuiltins(t *testing.T) {
	for _, id := range []string{
		"randi-unit", "randi-zero", "randi-empty", "randi-domain",
		"randi-negative-domain", "randi-evaluation",
		"rand-context-domain-bare", "rand-context-domain-call", "rand-context-direct-comparison",
		"rand-context-assigned-sum", "rand-context-grouped-sum", "rand-context-argument",
		"rand-context-command", "rand-context-argument-effect", "rand-context-prefix-error", "rand-context-extra",
		"rand-suffix-ignored-addition", "rand-suffix-grouped-comparisons",
		"rand-suffix-command-call", "rand-suffix-binary-call",
	} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			checkFormattingPreservesExecution(t, src, want)
		})
	}
}

func TestRandomBuiltinDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"randi-fraction", diag.ETypeMismatch, 3},
		{"randi-string", diag.ETypeMismatch, 3},
		{"randi-logical", diag.ETypeMismatch, 3},
		{"randi-extra", diag.EParse, 3},
		{"randi-bare", diag.EParse, 3},
		{"randi-vector", diag.EParse, 5},
		{"randi-no-value", diag.ETypeMismatch, 3},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			checkSemanticDiagnostic(t, path, tt.code, tt.line)
		})
	}
}
