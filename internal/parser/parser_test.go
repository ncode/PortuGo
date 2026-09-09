package parser

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/testprocess"
)

func TestPowerPrecedence(t *testing.T) {
	for _, tt := range []struct{ source, printed string }{
		{"2^3^2", "2 ^ 3 ^ 2"},
		{"-2^2", "- 2 ^ 2"},
		{"2^(3^2)", "2 ^ (3 ^ 2)"},
		{"-(2^2)", "- (2 ^ 2)"},
	} {
		t.Run(tt.source, func(t *testing.T) {
			_, toks, lexDiags := lexer.Scan("power.alg", "algoritmo \"power\"\ninicio\nescreval("+tt.source+")\nfimalgoritmo")
			prog, parseDiags := Parse(toks)
			if len(lexDiags)+len(parseDiags) != 0 {
				t.Fatalf("diagnostics: %v %v", lexDiags, parseDiags)
			}
			var out bytes.Buffer
			if err := ast.Fprint(&out, prog); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), "escreval("+tt.printed+")") {
				t.Fatalf("want expression %s in:\n%s", tt.printed, &out)
			}
		})
	}
}

func TestVectorBoundRecovery(t *testing.T) {
	file, tokens, ds := lexer.Scan("bounds.alg", `algoritmo "bounds"
var
  v: vetor[1+1..3] de inteiro
  w: vetor[0..+2] de inteiro
  n: inteiro
inicio
  escreval(n)
fimalgoritmo`)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	prog, ds := Parse(tokens)
	if len(ds) != 2 {
		t.Fatalf("diagnostics = %v, want two independent bound errors", ds)
	}
	for n, d := range ds {
		if d.Code != diag.EParse || file.Position(d.Pos).Line != n+3 {
			t.Fatalf("diagnostic %v, want P001 on line %d", d, n+3)
		}
	}
	if prog == nil || len(prog.Globals) != 3 || len(prog.Body) != 1 {
		t.Fatalf("bound recovery lost following declarations or body: %+v", prog)
	}
}

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
		decoded, err := source.Decode([]byte(src))
		if err != nil {
			t.Fatal(err)
		}
		_, toks, lexDiags := lexer.Scan("fuzz.alg", decoded)
		prog, parseDiags := Parse(toks)
		if len(lexDiags)+len(parseDiags) == 0 {
			sema.Analyze(prog)
			if err := ast.Fprint(io.Discard, prog); err != nil {
				t.Fatal(err)
			}
		}
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
