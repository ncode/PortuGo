package parser

import (
	"testing"

	"github.com/ncode/portugol-go/internal/lexer"
)

func FuzzParser(f *testing.F) {
	for _, seed := range []string{
		"",
		"algoritmo \"x\"\ninicio\nfimalgoritmo",
		"algoritmo \"x\"\nvar\nx: inteiro\ninicio\nx <- 1\nfimalgoritmo",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, src string) {
		_, toks, _ := lexer.Scan("fuzz.alg", src)
		Parse(toks)
	})
}
