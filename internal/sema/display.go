package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func (c *checker) checkColorArg(expr ast.Expr) bool {
	before := len(c.diags)
	switch typ := c.expr(expr); typ.Kind {
	case runtime.StringType, runtime.InvalidType:
	case runtime.VoidType:
		c.error(expr.Start(), diag.EParse, "missing color argument value")
	default:
		c.error(expr.Start(), diag.ETypeMismatch, "expected caractere, got %s", typ)
	}
	return len(c.diags) == before
}
