package parser

import (
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/testprocess"
)

func FuzzParser(f *testing.F) {
	for _, seed := range []string{
		"",
		"algoritmo \"x\"\ninicio\nfimalgoritmo",
		"algoritmo \"x\"\nvar\nx: inteiro\ninicio\nx <- 1\nfimalgoritmo",
		"algoritmo \"x\"\r\ninicio // end\r\nfimalgoritmo\r\n",
		"algoritmo \"x\"\nvar\nv: vetor[-2..2,1..3] de real\ninicio\nfimalgoritmo",
		"algoritmo \"x\"\ninicio\nescreval((1+2)^3)\nfimalgoritmo",
		"algoritmo \"x\"\ninicio\nse entao senao fimse\nfimalgoritmo",
		"\xff\x00\"",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, src string) {
		if len(src) > testprocess.MaxSourceBytes {
			t.Skip("outside 64 KiB fuzz profile")
		}
		_, toks, _ := lexer.Scan("fuzz.alg", src)
		Parse(toks)
	})
}

func TestFuzzAdversarial(t *testing.T) {
	for _, tt := range []struct{ name, src string }{
		{"truncated nesting", "algoritmo \"x\"\ninicio\nescreval(" + strings.Repeat("(", 4096)},
		{"unary nesting", "algoritmo \"x\"\ninicio\nescreval(" + strings.Repeat("-", 4096) + "1)\nfimalgoritmo"},
		{"flat expression", "algoritmo \"x\"\ninicio\nescreval(" + strings.Repeat("1+", 8192) + "1)\nfimalgoritmo"},
		{"recovery delimiters", "algoritmo \"x\"\ninicio\n" + strings.Repeat("fimse\n", 8192)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testprocess.Run(t, func() {
				_, toks, _ := lexer.Scan("adversarial.alg", tt.src)
				Parse(toks)
			})
		})
	}
}
