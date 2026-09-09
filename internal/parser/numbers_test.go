package parser

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
)

func TestNumericLiteralKinds(t *testing.T) {
	for _, tt := range []struct {
		source  string
		kind    ast.LiteralKind
		integer int64
		real    float64
	}{
		{"0", ast.IntLiteral, 0, 0},
		{"2147483647", ast.IntLiteral, 2147483647, 0},
		{"2147483648", ast.RealLiteral, 0, 2147483648},
		{"9223372036854775807", ast.RealLiteral, 0, 9223372036854775808},
		{"18446744073709551615", ast.RealLiteral, 0, 18446744073709551616},
		{"00000000000000000000000000000001", ast.IntLiteral, 1, 0},
		{"1.0", ast.RealLiteral, 0, 1},
		{"1e0", ast.RealLiteral, 0, 1},
	} {
		t.Run(tt.source, func(t *testing.T) {
			source := "algoritmo \"literal\"\ninicio\nescreval(" + tt.source + ")\nfimalgoritmo"
			for round := 0; round < 2; round++ {
				_, tokens, ds := lexer.Scan("literal.alg", source)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				prog, ds := Parse(tokens)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				got := prog.Body[0].(*ast.WriteStmt).Args[0].Expr.(*ast.LiteralExpr)
				if got.Kind != tt.kind || got.Int != tt.integer || got.Real != tt.real {
					t.Fatalf("round %d: literal=%+v, want kind=%d integer=%d real=%g", round, got, tt.kind, tt.integer, tt.real)
				}
				var printed bytes.Buffer
				if err := ast.Fprint(&printed, prog); err != nil {
					t.Fatal(err)
				}
				source = printed.String()
			}
		})
	}
}

func TestUnrepresentableNumericLiteral(t *testing.T) {
	for _, text := range []string{"1e400", strings.Repeat("9", 400)} {
		file, tokens, ds := lexer.Scan("literal.alg", "algoritmo \"literal\"\ninicio\nescreval("+text+")\nfimalgoritmo")
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		_, ds = Parse(tokens)
		if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != 3 {
			t.Fatalf("diagnostics=%v, want positioned P001", ds)
		}
	}
}
