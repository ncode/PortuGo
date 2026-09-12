package parser

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/testprocess"
	"github.com/ncode/portugol-go/internal/token"
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

func TestChoiceRangePrinting(t *testing.T) {
	_, tokens, ds := lexer.Scan("choice.alg", "algoritmo \"ranges\"\ninicio\nescolha 2.5\ncaso 1+1 até 3, 5\nescreval(1)\nfimescolha\nfimalgoritmo\n")
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	prog, ds := Parse(tokens)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	var out bytes.Buffer
	if err := ast.Fprint(&out, prog); err != nil {
		t.Fatal(err)
	}
	want := "algoritmo \"ranges\"\ninicio\n  escolha 2.5\n    caso 1 + 1 ate 3, 5:\n      escreval(1)\n  fimescolha\nfimalgoritmo\n"
	if out.String() != want {
		t.Fatalf("formatted choice = %q, want %q", out.String(), want)
	}
}

func FuzzParser(f *testing.F) {
	for _, seed := range []string{
		"",
		"algoritmo \"x\"\ninicio\nfimalgoritmo",
		"algoritmo \"x\"\nvar\nx: inteiro\ninicio\nx <- 1\nfimalgoritmo",
		"algoritmo \"x\"\nconst\nn = 1+2\nvar\ninicio\nescreval(n)\nfimalgoritmo",
		"algoritmo \"x\"\nconst\nn = 1+2\nvar\nv: vetor[1..n] de inteiro\ninicio\nv[n] <- 7\nfimalgoritmo",
		"algoritmo \"x\"\r\ninicio // end\r\nfimalgoritmo\r\n",
		"algoritmo \"x\"\nvar\nv: vetor[-2..2,1..3] de real\ninicio\nfimalgoritmo",
		"algoritmo \"x\"\ninicio\nescreval((1+2)^3)\nfimalgoritmo",
		"algoritmo \"x\"\ninicio\nse entao senao fimse\nfimalgoritmo",
		"algoritmo \"x\"\ninicio\nescolha 2\ncaso 1 ate 3, 5\nescreval(1)\nfimescolha\nfimalgoritmo",
		"algoritmo \"x\"\ninicio\nlimpatela()\nmudacor(\"amarelo\",\"frente\", ignored)\nfimalgoritmo",
		"\xff\x00\"",
		"\xef\xbb\xbfalgoritmo \"a\xe7\xe3o\"\r\r\ninicio\r\r\nfimalgoritmo\n",
		"algoritmo \"x\"\ninicio\nescreval(1 + // first\r\r\n2) // last\nfimalgoritmo\n\"ignored",
		"algoritmo \"x\"\nvar\nprocedimento P(a: inteiro; // first\nb: inteiro)\ninicio\nfimprocedimento\ninicio\nfimalgoritmo",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, src string) {
		if len(src) > testprocess.MaxSourceBytes {
			t.Skip("outside 64 KiB fuzz profile")
		}
		decoded, err := source.DecodeFile("fuzz.alg", []byte(src))
		if err != nil {
			t.Fatal(err)
		}
		_, toks, lexDiags := lexer.ScanFile(decoded)
		prog, parseDiags := Parse(toks)
		checkDiagnosticPositions(t, len(src), lexDiags, parseDiags)
		if len(lexDiags)+len(parseDiags) != 0 {
			return
		}
		_, ds := sema.Analyze(prog)
		checkDiagnosticPositions(t, len(src), ds)
		var first, second bytes.Buffer
		if err := ast.Fprint(&first, prog); err != nil {
			t.Fatal(err)
		}
		_, toks, lexDiags = lexer.Scan("formatted.alg", first.String())
		reparsed, parseDiags := Parse(toks)
		checkDiagnosticPositions(t, first.Len(), lexDiags, parseDiags)
		if len(lexDiags) == 0 && len(parseDiags) == 1 && parseDiags[0].Code == diag.EResource {
			return // The CLI rejects output that exceeds a structural limit.
		}
		if len(lexDiags)+len(parseDiags) != 0 {
			t.Fatal("formatting produced invalid syntax")
		}
		if err := ast.Fprint(&second, reparsed); err != nil {
			t.Fatal(err)
		}
		if first.String() != second.String() {
			t.Fatal("formatting is not idempotent")
		}
	})
}

func checkDiagnosticPositions(t *testing.T, size int, groups ...[]diag.Diagnostic) {
	t.Helper()
	for _, ds := range groups {
		for _, d := range ds {
			if d.Code == "" || d.Pos < 0 || int(d.Pos) > size || d.End != token.NoPos && (d.End < d.Pos || int(d.End) > size) {
				t.Fatalf("invalid diagnostic span: %s [%d,%d)", d.Code, d.Pos, d.End)
			}
		}
	}
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
				_, toks, lexDiags := lexer.Scan("adversarial.alg", tt.src)
				_, parseDiags := Parse(toks)
				checkDiagnosticPositions(t, len(tt.src), lexDiags, parseDiags)
				want := diag.EResource
				if tt.name == "recovery delimiters" {
					want = diag.EParse
				}
				if len(lexDiags) != 0 || len(parseDiags) == 0 || parseDiags[0].Code != want {
					t.Fatalf("missing controlled %s rejection", want)
				}
			})
		})
	}
}
