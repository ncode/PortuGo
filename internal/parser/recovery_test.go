package parser

import (
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
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
