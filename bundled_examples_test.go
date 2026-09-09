package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedBundledExamples(t *testing.T) {
	for _, id := range []string{
		"bundled-6c81c979714d", // ajustes.alg
		"bundled-6fca58cfe4ac", // passo.alg
		"bundled-a2fd5d8d6435", // taxaspop.alg
		"bundled-069318ac942f", // Exemplos/TESTE.ALG
		"bundled-6f7f325a3d14", // troca.alg
		"bundled-a3ab1c0f85e6", // vetr2dim.alg
		"bundled-08dbc8f0cc9e", // PRIMOS.ALG
		"bundled-b8d8516450ea", // caracfun.alg
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
