package sema

import (
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
)

func TestAbsentComparisonGuard(t *testing.T) {
	for _, op := range []string{"=", "<>", "<", ">", "<=", ">="} {
		t.Run(op, func(t *testing.T) {
			src := "algoritmo \"absent comparison\"\ninicio\nescreval(abs() " + op + " abs())\nfimalgoritmo\n"
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			prog, ds := parser.Parse(tokens)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			_, ds = Analyze(prog)
			if len(ds) != 1 || ds[0].Code != diag.ETypeMismatch || file.Position(ds[0].Pos).Line != 3 || file.Position(ds[0].Pos).Column != 16 {
				t.Fatalf("diagnostics = %v, want E001 at 3:16", ds)
			}
		})
	}
}
