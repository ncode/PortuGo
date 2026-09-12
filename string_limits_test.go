package portugol_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedStringLimits(t *testing.T) {
	for _, id := range []string{
		"conversion-enormous", "conversion-edge-literal-length", "conversion-edge-joined-length",
		"conversion-edge-joined-number", "conversion-edge-suffix-after-255",
		"string-size-literal-254", "string-size-literal-255", "string-size-literal-256", "string-size-literal-511",
		"string-size-accent-tail", "string-size-joined-tail", "string-size-accent-count", "string-size-uppercase",
		"string-size-search-tail", "string-size-copy-tail", "string-size-comparison", "string-size-assignment",
		"string-size-parameter-result", "string-size-vector", "string-size-output-width",
		"string-size-constant-value", "string-size-constant-joined", "string-size-input-255", "string-size-input-256",
		"string-field-widths",
		"string-field-large-width",
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
			input, err := os.ReadFile(filepath.Join(dir, "input.txt"))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
			checkFormattingPreservesExecutionWithInput(t, src, input, want)
		})
	}
}
