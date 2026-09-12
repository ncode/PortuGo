package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedConsoleInput(t *testing.T) {
	for _, id := range []string{
		"input-single-integer", "input-two-integers", "input-integer-edge-spaces",
		"input-string-spaces", "input-string-edge-spaces", "input-string-empty",
		"input-real-dot", "input-real-comma", "input-logical-word",
		"input-string-lines", "input-integer-empty", "input-integer-fraction", "input-integer-packed",
		"input-logical-english", "input-logical-numeric", "input-logical-false",
		"input-integer-negative-fraction", "input-integer-exponent", "input-integer-invalid",
		"input-integer-hex", "input-integer-leading-zero", "input-integer-large",
		"input-real-trailing-text", "input-real-empty", "input-real-exponent", "input-real-invalid",
		"input-logical-uppercase", "input-logical-letter", "input-logical-edge-spaces", "input-logical-empty",
		"input-real-mantissa-tail", "input-real-negative-tail", "input-real-trailing-space", "input-real-packed",
		"input-real-exponent-tail", "input-real-decimal-exponent-tail", "input-real-missing-exponent",
		"input-real-integer-missing-exponent", "input-real-repeated-dot",
		"input-integer-wrap-positive", "input-integer-wrap-negative", "input-integer-maximum",
		"input-logical-prefix", "input-logical-trailing-space", "input-logical-lowercase-letter", "input-logical-short-word",
		"input-integer-invalid-replaces", "input-real-invalid-replaces", "input-real-leading-space", "input-real-negative-zero",
	} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			input, err := os.ReadFile(filepath.Join(dir, "input.txt"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			checkFormattingPreservesExecutionWithInput(t, src, input, want)
			checkFormattingPreservesExecutionWithInput(t, src, bytes.ReplaceAll(input, []byte("\n"), []byte("\r\n")), want)
		})
	}
}
