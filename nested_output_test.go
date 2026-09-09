package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedNestedOutput(t *testing.T) {
	for _, id := range []string{
		"nested-escreva-inner-newline", "nested-escreva-inner-write",
		"nested-escreval-inner-newline", "nested-escreval-inner-write",
		"nested-newline-bare-newline", "nested-newline-bare-write",
		"nested-newline-no-inner-write", "nested-newline-two-calls", "nested-newline-width-call",
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
