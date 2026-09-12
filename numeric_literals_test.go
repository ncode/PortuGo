package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedNumericLiterals(t *testing.T) {
	for _, id := range []string{
		"literal-large-output", "literal-large-precision", "literal-int64-output",
		"literal-uint64-output", "literal-leading-zeros", "integer-computed-minimum",
		"real-output-fractions", "real-output-thresholds", "real-output-zero",
		"real-output-default-width", "real-output-addition", "real-output-exponents",
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

func TestNumericLiteralDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		line int
	}{
		{"literal-large-return", 4},
		{"literal-negative-minimum", 3},
		{"literal-negative-outside", 3},
		{"randi-wrap-zero", 3},
		{"randi-wrap-unit", 3},
		{"randi-signed-boundary", 8},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			checkSemanticDiagnostic(t, path, diag.ETypeMismatch, tt.line)
		})
	}
}
