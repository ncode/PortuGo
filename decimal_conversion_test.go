package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedDecimalConversions(t *testing.T) {
	for _, id := range []string{
		"format-precision-small-rounding", "format-identity-tiny-decimals",
		"decimal-conversion-tiny-lengths", "decimal-conversion-leading-zeroes",
		"decimal-conversion-significant-tails", "decimal-conversion-exponent-scaling",
		"decimal-zero-fractions", "decimal-zero-scaled-fractions", "decimal-zero-signs-zero",
		"decimal-type-zero", "decimal-type-fraction", "decimal-type-negative",
		"decimal-type-leading-zeroes", "decimal-type-scaled",
		"decimal-invalid-suffix", "decimal-invalid-zeroes-suffix", "decimal-invalid-exponent-overflow",
		"decimal-invalid-fraction-overflow", "decimal-invalid-repeated-point",
		"decimal-invalid-underscore", "decimal-invalid-alphabetic",
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
