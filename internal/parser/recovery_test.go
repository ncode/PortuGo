package parser

import (
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/source"
)

func TestStatementLineRecovery(t *testing.T) {
	file, tokens, ds := lexer.Scan("recovery.alg", "algoritmo \"recovery\"\ninicio\n(1 + 2)\nescreval(7)\n(3 + 4)\nfimalgoritmo\n")
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	prog, ds := Parse(tokens)
	if len(ds) != 2 {
		t.Fatalf("diagnostics = %v, want two independent line errors", ds)
	}
	for i, line := range []int{3, 5} {
		if ds[i].Code != diag.EParse || file.Position(ds[i].Pos).Line != line {
			t.Fatalf("diagnostic %d = %v, want P001 on line %d", i, ds[i], line)
		}
	}
	if prog == nil || len(prog.Body) != 1 {
		t.Fatal("recovery lost the valid statement between malformed lines")
	}
}

func TestReturnValueDoesNotCrossLine(t *testing.T) {
	src, err := source.ReadFile("../../testdata/conformance/visualg-3.0.7/probes/return-value-next-line/source.alg")
	if err != nil {
		t.Fatal(err)
	}
	_, tokens, ds := lexer.Scan("source.alg", src)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	// The following line is also invalid; inspect recovery without fixing the
	// precedence between syntax diagnostics and the earlier missing value here.
	prog, _ := Parse(tokens)
	if prog == nil || len(prog.Subs) != 1 {
		t.Fatal("missing function after recovery")
	}
	fn, ok := prog.Subs[0].(*ast.FunctionDecl)
	if !ok || len(fn.Body) != 1 {
		t.Fatal("missing return after recovery")
	}
	ret, ok := fn.Body[0].(*ast.ReturnStmt)
	if !ok || ret.Value != nil {
		t.Fatal("return consumed an expression from the next line")
	}
}
