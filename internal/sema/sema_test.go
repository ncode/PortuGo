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

func TestFunctionMustReturn(t *testing.T) {
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
	if len(diags) == 0 || diags[0].Code != diag.EReturn {
		t.Fatalf("got diagnostics %#v, want %s", diags, diag.EReturn)
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
	return Check(prog)
}
