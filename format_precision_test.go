package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedFormatPrecision(t *testing.T) {
	for _, id := range []string{
		"string-field-real-precision", "string-field-integer-precision", "string-field-wide-precision",
		"format-precision-real-digits", "format-precision-integer-digits", "format-precision-mantissa-length",
		"format-precision-tenths", "format-precision-thirds", "format-precision-finite-decimal",
		"format-precision-negative-count", "format-precision-zero", "format-precision-significant-values",
		"format-precision-large-magnitude", "format-precision-large-count", "format-precision-rounding",
		"format-precision-tiny-values", "format-precision-exponent-threshold", "format-precision-scientific-width",
		"format-precision-scientific-mantissa",
		"format-precision-tiny-cutoff", "format-precision-large-cutoff", "format-precision-rounded-zero", "format-precision-exact-ties",
		"format-identity-tiny-arithmetic", "format-identity-large-threshold", "format-identity-leading-digits",
		"format-scaling-leading-range", "format-scaling-binary-threshold",
		"format-transition-binary-digits", "format-transition-exponent-grid",
		"format-transition-scientific-transition", "format-transition-scientific-rounding",
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
