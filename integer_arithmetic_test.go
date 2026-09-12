package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedIntegerArithmetic(t *testing.T) {
	for _, id := range []string{
		"integer-overflow-add", "integer-overflow-subtract",
		"integer-multiply-wrap", "integer-multiply-zero", "integer-negate-minimum",
		"integer-divide-signs", "integer-real-sum", "integer-assignment-wrap",
		"integer-function-wrap", "integer-abs-minimum", "integer-int-large", "integer-int-fraction",
		"int-random-bound", "int-truncation", "int-wrap-zero", "abs-random-bound",
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
