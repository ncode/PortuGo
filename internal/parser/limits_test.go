package parser

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/testprocess"
)

func TestFormatDepthRoundTrip(t *testing.T) {
	parse := func(src string) *ast.Program {
		t.Helper()
		_, toks, ds := lexer.Scan("format.alg", src)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		// Source positions change with formatting; the syntax tree must not.
		for i := range toks {
			toks[i].Pos = 0
		}
		prog, ds := Parse(toks)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		return prog
	}
	for _, expr := range []string{
		strings.Repeat("1+", 254) + "1",
		strings.Repeat("-", 254) + "1",
		strings.Repeat("nao ", 254) + "verdadeiro",
		"1 - (2 - 3) + 4 * (5 + 6)",
		"(2 ^ 3) ^ 2 + 2 ^ (3 ^ 2)",
		"-2 ^ 2 + -(2 ^ 2)",
		"nao (verdadeiro e falso) ou (falso xou verdadeiro)",
		"(nao 1) ^ 2 + - nao 3",
	} {
		prog := parse("algoritmo \"format\"\ninicio\nescreval(" + expr + ")\nfimalgoritmo\n")
		var first, second bytes.Buffer
		if err := ast.Fprint(&first, prog); err != nil {
			t.Fatal(err)
		}
		reparsed := parse(first.String())
		if !reflect.DeepEqual(prog, reparsed) {
			t.Fatalf("formatting changed expression %q", expr)
		}
		if err := ast.Fprint(&second, reparsed); err != nil {
			t.Fatal(err)
		}
		if first.String() != second.String() {
			t.Fatalf("formatting is not idempotent for %q", expr)
		}
	}
}

func TestStructuralLimits(t *testing.T) {
	for _, tt := range []struct {
		name, body string
		limit      bool
	}{
		{"flat boundary", "escreva(" + strings.Repeat("1+", 254) + "1)", false},
		{"flat excessive", "escreva(" + strings.Repeat("1+", 255) + "1)", true},
		{"case lower boundary", "escolha 1\ncaso " + strings.Repeat("1+", 254) + "1 ate 5\nfimescolha", false},
		{"case lower excessive", "escolha 1\ncaso " + strings.Repeat("1+", 255) + "1 ate 5\nfimescolha", true},
		{"case upper boundary", "escolha 1\ncaso 1 ate " + strings.Repeat("1+", 254) + "1\nfimescolha", false},
		{"case upper excessive", "escolha 1\ncaso 1 ate " + strings.Repeat("1+", 255) + "1\nfimescolha", true},
		{"parentheses boundary", "escreva(" + strings.Repeat("(", 254) + "1" + strings.Repeat(")", 254) + ")", false},
		{"parentheses excessive", "escreva(" + strings.Repeat("(", 255) + "1" + strings.Repeat(")", 255) + ")", true},
		{"nested statements", strings.Repeat("se verdadeiro entao\n", 300) + strings.Repeat("fimse\n", 300), true},
		{"nested calls", "escreva(" + strings.Repeat("abs(", 300) + "1" + strings.Repeat(")", 300) + ")", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testprocess.Run(t, func() {
				_, tokens, ds := lexer.Scan("limit.alg", "algoritmo \"limit\"\ninicio\n"+tt.body+"\nfimalgoritmo")
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				p, ds := Parse(tokens)
				if tt.limit {
					if p != nil || len(ds) != 1 || ds[0].Code != diag.EResource || ds[0].Pos == 0 {
						t.Fatalf("missing single limit diagnostic: %v", ds)
					}
				} else if len(ds) != 0 {
					t.Fatalf("rejected boundary: %v", ds)
				}
			})
		})
	}
}

func TestTypeDepthLimit(t *testing.T) {
	for _, vectors := range []int{255, 256} {
		src := "algoritmo \"type\"\nvar v: " + strings.Repeat("vetor[1..1] de ", vectors) + "inteiro\ninicio\nfimalgoritmo"
		_, toks, ds := lexer.Scan("type.alg", src)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		p, ds := Parse(toks)
		if vectors == 255 {
			if len(ds) != 0 {
				t.Fatalf("type boundary rejected: %v", ds)
			}
		} else if p != nil || len(ds) != 1 || ds[0].Code != diag.EResource {
			t.Fatalf("unbounded type: %v", ds)
		}
	}
}

func TestLimitKeepsEarlierDiagnostics(t *testing.T) {
	for _, expr := range []string{strings.Repeat("(", 300) + "1" + strings.Repeat(")", 300), strings.Repeat("1+", 300) + "1"} {
		_, toks, ds := lexer.Scan("invalid.alg", "algoritmo \"invalid\"\nvar\nv: vetor[2..1] de inteiro\ninicio\nescreva("+expr+")\nfimalgoritmo")
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		p, ds := Parse(toks)
		if p != nil || len(ds) != 2 || ds[0].Code != diag.EParse || ds[1].Code != diag.EResource {
			t.Fatalf("lost error or skipped guard: %v", ds)
		}
	}
}
