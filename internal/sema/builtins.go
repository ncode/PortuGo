package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func (c *checker) declareBuiltins() {
	for _, name := range []string{
		"abs", "raizq", "exp", "log", "logn", "pi", "sen", "cos", "tan", "int",
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
	case "abs", "raizq", "log", "logn", "sen", "cos", "tan", "int",
		"arccos", "arcsen", "arctan", "cotan", "grauprad", "radpgrau", "quad":
		result := runtime.Type{Kind: runtime.RealType}
		if name == "int" {
			result.Kind = runtime.IntegerType
		}
		if len(call.Args) == 0 {
			if name == "abs" || name == "quad" {
				result.Kind = runtime.VoidType
			}
			return result, true
		}
		if len(call.Args) > 1 && (name == "abs" || name == "quad" || name == "raizq") {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after numeric argument")
			return result, true
		}
		t := c.expr(call.Args[0])
		switch t.Kind {
		case runtime.IntegerType, runtime.RealType:
			if name == "abs" || name == "quad" {
				return t, true
			}
		case runtime.VoidType:
			return t, true
		case runtime.StringType, runtime.BoolType:
			if name == "abs" || name == "quad" || name == "int" || name == "raizq" {
				return runtime.Type{Kind: runtime.VoidType}, true
			}
			c.error(call.Args[0].Start(), diag.EParse, "expected numeric expression")
		case runtime.VectorType:
			c.error(call.Args[0].Start(), diag.EParse, "expected '[' after vector")
		}
		if len(call.Args) > 1 {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after numeric argument")
		}
		return result, true
	case "exp":
		if len(call.Args) == 0 {
			return runtime.Type{Kind: runtime.VoidType}, true
		}
		for _, arg := range call.Args[:min(len(call.Args), 2)] {
			switch c.expr(arg).Kind {
			case runtime.VoidType, runtime.StringType, runtime.BoolType:
				return runtime.Type{Kind: runtime.VoidType}, true
			case runtime.VectorType:
				c.error(arg.Start(), diag.EParse, "expected '[' after vector")
			}
		}
		if len(call.Args) != 2 {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after numeric argument")
		}
		return runtime.Type{Kind: runtime.RealType}, true
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
