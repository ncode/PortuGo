package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func (c *checker) declareBuiltins() {
	for _, name := range []string{
		"abs", "raizq", "exp", "log", "logn", "pi", "sen", "cos", "tan", "int", "frac",
		"aleatorio", "copia", "maiusc", "minusc", "asc", "carac", "compr", "pos", "numpcarac", "randi",
		"arccos", "arcsen", "arctan", "cotan", "grauprad", "radpgrau", "quad",
	} {
		c.scope.declare(symbol{name: name, kind: builtinSym})
	}
}

func (c *checker) builtinCallType(name string, call *ast.CallExpr) (runtime.Type, bool) {
	switch name {
	case "pi":
		c.requireArity(call, 0, 0)
		return runtime.Type{Kind: runtime.RealType}, true
	case "abs":
		if c.requireArity(call, 1, 1) {
			t := c.expr(call.Args[0])
			if isNumeric(t) {
				return t, true
			}
			c.error(call.Args[0].Start(), diag.ETypeMismatch, "abs requires numeric argument")
		}
		return runtime.Type{Kind: runtime.InvalidType}, true
	case "raizq", "log", "sen", "cos", "tan", "frac":
		c.checkBuiltinArgs(call, runtime.Type{Kind: runtime.RealType})
		return runtime.Type{Kind: runtime.RealType}, true
	case "arccos", "arcsen", "arctan", "cotan", "grauprad", "radpgrau", "quad":
		if len(call.Args) > 1 {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after numeric argument")
		} else if len(call.Args) == 0 {
			if name == "quad" {
				return runtime.Type{Kind: runtime.VoidType}, true
			}
		} else {
			t := c.expr(call.Args[0])
			switch t.Kind {
			case runtime.IntegerType, runtime.RealType:
				if name == "quad" {
					return t, true
				}
			case runtime.VoidType:
				return t, true
			case runtime.StringType, runtime.BoolType:
				if name == "quad" {
					return runtime.Type{Kind: runtime.VoidType}, true
				}
				c.error(call.Args[0].Start(), diag.EParse, "expected numeric expression")
			case runtime.VectorType:
				c.error(call.Args[0].Start(), diag.EParse, "expected '[' after vector")
			}
		}
		return runtime.Type{Kind: runtime.RealType}, true
	case "exp", "logn":
		if c.requireArity(call, 2, 2) {
			for _, arg := range call.Args {
				t := c.expr(arg)
				if !isNumeric(t) {
					c.error(arg.Start(), diag.ETypeMismatch, "%q requires numeric arguments", call.Name.Text)
				}
			}
		}
		return runtime.Type{Kind: runtime.RealType}, true
	case "int":
		if c.requireArity(call, 1, 1) {
			t := c.expr(call.Args[0])
			if !isNumeric(t) {
				c.error(call.Args[0].Start(), diag.ETypeMismatch, "int requires numeric argument")
			}
		}
		return runtime.Type{Kind: runtime.IntegerType}, true
	case "numpcarac":
		if len(call.Args) > 1 {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after conversion argument")
		} else if len(call.Args) == 1 {
			t := c.expr(call.Args[0])
			switch t.Kind {
			case runtime.StringType, runtime.BoolType, runtime.VoidType:
				return runtime.Type{Kind: runtime.VoidType}, true
			case runtime.VectorType:
				c.error(call.Args[0].Start(), diag.EParse, "expected '[' after vector")
			}
		}
		return runtime.Type{Kind: runtime.StringType}, true
	case "randi":
		if len(call.Args) > 1 {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after random bound")
		} else if len(call.Args) == 1 {
			switch c.expr(call.Args[0]).Kind {
			case runtime.VectorType:
				c.error(call.Args[0].Start(), diag.EParse, "expected '[' after vector")
			case runtime.IntegerType, runtime.InvalidType:
			default:
				c.error(call.Args[0].Start(), diag.ETypeMismatch, "expected inteiro argument")
			}
		}
		return runtime.Type{Kind: runtime.IntegerType}, true
	case "aleatorio":
		if !c.requireArity(call, 0, 2) {
			return runtime.Type{Kind: runtime.InvalidType}, true
		}
		for _, arg := range call.Args {
			c.requireInt(arg)
		}
		if len(call.Args) == 0 {
			return runtime.Type{Kind: runtime.RealType}, true
		}
		return runtime.Type{Kind: runtime.IntegerType}, true
	case "copia":
		if c.requireArity(call, 3, 3) {
			c.requireString(call.Args[0])
			c.requireInt(call.Args[1])
			c.requireInt(call.Args[2])
		}
		return runtime.Type{Kind: runtime.StringType}, true
	case "maiusc", "minusc":
		if c.requireArity(call, 1, 1) {
			c.requireString(call.Args[0])
		}
		return runtime.Type{Kind: runtime.StringType}, true
	case "asc":
		if c.requireArity(call, 1, 1) {
			c.requireString(call.Args[0])
		}
		return runtime.Type{Kind: runtime.IntegerType}, true
	case "carac":
		if c.requireArity(call, 1, 1) {
			c.requireInt(call.Args[0])
		}
		return runtime.Type{Kind: runtime.StringType}, true
	case "compr":
		if c.requireArity(call, 1, 1) {
			c.requireString(call.Args[0])
		}
		return runtime.Type{Kind: runtime.IntegerType}, true
	case "pos":
		if c.requireArity(call, 2, 2) {
			c.requireString(call.Args[0])
			c.requireString(call.Args[1])
		}
		return runtime.Type{Kind: runtime.IntegerType}, true
	default:
		return runtime.Type{}, false
	}
}

func (c *checker) checkBuiltinArgs(call *ast.CallExpr, ret runtime.Type) {
	if c.requireArity(call, 1, 1) {
		t := c.expr(call.Args[0])
		if !isNumeric(t) {
			c.error(call.Args[0].Start(), diag.ETypeMismatch, "%q requires numeric argument", call.Name.Text)
		}
	}
	_ = ret
}

func (c *checker) requireArity(call *ast.CallExpr, min, max int) bool {
	n := len(call.Args)
	if n < min || n > max {
		if min == max {
			c.error(call.Name.Pos, diag.ECall, "%q expects %d arguments, got %d", call.Name.Text, min, n)
		} else {
			c.error(call.Name.Pos, diag.ECall, "%q expects %d to %d arguments, got %d", call.Name.Text, min, max, n)
		}
		return false
	}
	return true
}

func (c *checker) requireString(expr ast.Expr) {
	if t := c.expr(expr); t.Kind != runtime.StringType && t.Kind != runtime.InvalidType {
		c.error(expr.Start(), diag.ETypeMismatch, "expected caractere, got %s", t)
	}
}
