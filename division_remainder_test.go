package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedDivisionRemainder(t *testing.T) {
	for _, id := range []string{
		"integer-divide-real", "integer-modulo-overflow", "integer-modulo-real", "integer-modulo-signs",
		"modulo-negative-right", "modulo-negative-both", "modulo-negative-eight", "modulo-negative-variable",
		"modulo-small-left", "modulo-percent-negative", "divide-real-left", "divide-real-right", "divide-real-both",
		"modulo-real-right", "modulo-real-both", "divide-integral-real-left", "divide-integral-real-right",
		"divide-left-result-type", "divide-real-zero", "modulo-integral-real-right", "modulo-real-result-type",
		"modulo-negative-real-left", "modulo-zero-divisor", "modulo-large-real-left",
		"divide-string-left", "divide-string-right", "divide-logical-left", "divide-logical-right",
		"modulo-logical-left", "modulo-logical-right", "divide-real-negative-right",
		"modulo-unit-divisor", "modulo-real-wrap-zero", "operator-evaluation-order",
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

func TestDivisionResultDiagnostic(t *testing.T) {
	path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", "divide-right-result-type", "source.alg")
	checkSemanticDiagnostic(t, path, diag.ETypeMismatch, 3)
}
