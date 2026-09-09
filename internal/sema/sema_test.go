package sema

import (
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
)

func TestBreakOutsideLoop(t *testing.T) {
	diags := checkSource(t, `algoritmo "x"
inicio
  interrompa
fimalgoritmo`)
	if len(diags) != 1 || diags[0].Code != diag.EBreak {
		t.Fatalf("got diagnostics %#v, want one %s", diags, diag.EBreak)
	}
}

func TestFunctionMayFallThrough(t *testing.T) {
	diags := checkSource(t, `algoritmo "x"

funcao f(n: inteiro): inteiro
inicio
  se n > 0 entao
    retorne n
  fimse
fimfuncao

inicio
  escreval(f(1))
fimalgoritmo`)
	if len(diags) != 0 {
		t.Fatalf("fallthrough function rejected: %v", diags)
	}
}

func checkSource(t *testing.T, src string) []diag.Diagnostic {
	t.Helper()
	_, toks, lexDiags := lexer.Scan("test.alg", src)
	if len(lexDiags) > 0 {
		t.Fatalf("lexer diagnostics: %v", lexDiags)
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) > 0 {
		t.Fatalf("parser diagnostics: %v", parseDiags)
	}
	_, diags := Analyze(prog)
	return diags
}

func TestExpArguments(t *testing.T) {
	for _, tt := range []struct {
		expr string
		code diag.Code
	}{
		{"exp(2, 3)", ""},
		{"exp(2.0, 0.5)", ""},
		{"exp()", ""},
		{"exp(2)", diag.EParse},
		{"exp(2, 3, 4)", diag.EParse},
		{`exp("2", 3)`, ""},
		{"exp(2, verdadeiro)", ""},
		{"exp(arccos(2))", ""},
	} {
		t.Run(tt.expr, func(t *testing.T) {
			diags := checkSource(t, "algoritmo \"exp\"\ninicio\nescreval("+tt.expr+")\nfimalgoritmo")
			if tt.code == "" {
				if len(diags) != 0 {
					t.Fatalf("unexpected diagnostics: %v", diags)
				}
			} else if len(diags) != 1 || diags[0].Code != tt.code {
				t.Fatalf("got diagnostics %v, want one %s", diags, tt.code)
			}
		})
	}
}
