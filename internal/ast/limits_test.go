package ast_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/token"
)

func TestTraversalLimitsBeforeOutput(t *testing.T) {
	for _, depth := range []int{255, 256} {
		var expr ast.Expr = &ast.LiteralExpr{At: 2, Kind: ast.IntLiteral, Int: 1}
		for n := 1; n < depth; n++ {
			expr = &ast.UnaryExpr{Op: token.Token{Kind: token.SUB, Text: "-", Pos: token.Pos(n + 2)}, X: expr}
		}
		p := &ast.Program{Name: "depth", Body: []ast.Stmt{&ast.WriteStmt{At: 1, Args: []ast.WriteArg{{Expr: expr}}}}}
		var out bytes.Buffer
		err := ast.Fprint(&out, p)
		info, ds := sema.Analyze(p)
		if depth == 255 {
			if err != nil || len(ds) != 0 || !info.ValidFor(p) {
				t.Fatalf("boundary rejected: %v %v", err, ds)
			}
		} else {
			var d diag.Diagnostic
			if !errors.As(err, &d) || d.Code != diag.EResource || out.Len() != 0 || len(ds) != 1 || ds[0].Code != diag.EResource || info.ValidFor(p) {
				t.Fatalf("unprotected traversal: %v %v, output %d bytes", err, ds, out.Len())
			}
		}
	}
}
