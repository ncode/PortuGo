package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
)

func (c *checker) declareConsts(decls []ast.ConstDecl) bool {
	for _, decl := range decls {
		before := len(c.diags)
		typ := c.expr(decl.Value)
		if len(c.diags) != before {
			// The reference rejects unresolved declaration expressions as syntax.
			if c.diags[before].Code == diag.EUndeclared {
				c.diags[before].Code = diag.EParse
				c.diags = c.diags[:before+1]
			}
			return false
		}
		if !isScalar(typ) {
			c.error(decl.Name.Pos, diag.ETypeMismatch, "constant requires a scalar expression")
			return false
		}
		sym := symbol{pos: decl.Name.Pos, name: canon(decl.Name.Text), kind: constSym, typ: typ}
		if !c.scope.declare(sym) {
			c.error(decl.Name.Pos, diag.ERedeclared, "redeclared identifier %q", decl.Name.Text)
			return false
		}
		c.recordBinding(decl.Name, sym)
	}
	return true
}
