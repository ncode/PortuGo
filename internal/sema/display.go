package sema

import (
	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/runtime"
	"github.com/ncode/PortuGo/internal/token"
)

func (c *checker) checkColorArg(expr ast.Expr, deferInvalid bool) bool {
	before := len(c.diags)
	typ := c.expr(expr)
	switch typ.Kind {
	case runtime.StringType, runtime.InvalidType:
		if typ.Kind == runtime.InvalidType && c.mixedStringAddition(expr) {
			if deferInvalid {
				c.diags = c.diags[:before]
			} else {
				for i := before; i < len(c.diags); i++ {
					if c.diags[i].Code == diag.ETypeMismatch {
						c.diags[i].Code = diag.EParse
					}
				}
			}
		}
	case runtime.VoidType:
		if !deferInvalid {
			c.error(expr.Start(), diag.EParse, "missing color argument value")
		}
	default:
		if !deferInvalid {
			c.error(expr.Start(), diag.ETypeMismatch, "expected caractere, got %s", typ)
		}
	}
	return len(c.diags) == before
}

func (c *checker) mixedStringAddition(expr ast.Expr) bool {
	e, ok := expr.(*ast.BinaryExpr)
	if !ok || e.Op.Kind != token.ADD {
		return false
	}
	left, right := c.info.types[e.Left], c.info.types[e.Right]
	return (isNumeric(left) && right.Kind == runtime.StringType) ||
		(isNumeric(right) && left.Kind == runtime.StringType)
}
