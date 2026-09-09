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
		"bundled-183b40a8f657", // randomicos.alg
		"bundled-0d33ed8315bc", // combin.alg
		"bundled-c1a2be69b9ce", // fatorial.alg
		"bundled-2ecb6b470914", // fatorial2.alg
		"bundled-9ddfb3dd3ea1", // menorde3.alg
		"bundled-8c99b59ea498", // CALCULOMEDIA2.ALG
		"bundled-b233b8edad6d", // decomp.alg
		"bundled-a942295b3cc1", // dectobin.alg
		"bundled-fb5b43336226", // dectohex.alg
		"bundled-eedf876719f1", // EXEMPLO1.alg.ALG
		"bundled-82d1c4dba544", // media_aluno.alg
		"bundled-7afef248aa17", // MEDIA_SIMPLES.ALG
		"bundled-8a30556ea4fc", // MEDIA_VETOR.ALG
		"bundled-1dfa641f1d7a", // MODULO.ALG
		"bundled-dc0f54ded747", // Nome_inverso.alg
		"bundled-4fdb89658892", // PERFEITOS.ALG
		"bundled-f8eb1bdcb2cd", // rqpaprox.alg
		"bundled-3a5ff4924627", // times.alg
		"bundled-61a67cb1f566", // Troca de Valores.alg
		"bundled-116164dd09bb", // ENCRYPT.ALG
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
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			checkFormattingPreservesExecutionWithInput(t, src, input, want)
		})
	}
}
