package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedKeywordPrograms(t *testing.T) {
	for _, id := range []string{
		"string-type-caracter", "string-type-parameter-caracter",
		"keyword-entao-accented", "keyword-senao-accented",
		"keyword-nao-accented", "keyword-ate-accented",
		"keyword-faca-accented", "keyword-funcao-accented",
		"alias-caracter-lookalike",
		"division-alias-values", "division-alias-case", "division-alias-order",
		"function-ending-accented", "function-ending-uppercase",
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
