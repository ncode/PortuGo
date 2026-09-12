package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func (c *checker) checkBound(expr ast.Expr, at token.Pos) bool {
	if lit, ok := expr.(*ast.LiteralExpr); ok && lit.Kind == ast.IntLiteral && lit.Int >= 0 {
		c.info.types[expr] = runtime.Type{Kind: runtime.IntegerType}
		return true
	}
	if name, ok := expr.(*ast.IdentExpr); ok {
		if sym, found := c.lookup(name.Name); found && sym.kind == constSym && (sym.typ.Kind == runtime.IntegerType || sym.typ.Kind == runtime.NumericType) {
			c.info.types[expr] = sym.typ
			return true
		}
	}
	if expr != nil {
		at = expr.Start()
	}
	c.error(at, diag.EParse, "expected unsigned integer literal or integer constant bound")
	return false
}
