package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func (c *checker) declareBuiltins() {
	for _, name := range []string{
		"abs", "raizq", "exp", "log", "logn", "pi", "sen", "cos", "tan", "int",
		"aleatorio", "copia", "maiusc", "minusc", "asc", "carac", "compr", "pos", "numpcarac", "caracpnum", "randi", "rand",
		"arccos", "arcsen", "arctan", "cotan", "grauprad", "radpgrau", "quad",
	} {
		c.scope.declare(symbol{name: name, kind: builtinSym})
	}
}

func (c *checker) builtinCallType(name string, call *ast.CallExpr) (runtime.Type, bool) {
	switch name {
	case "pi", "rand":
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
		case runtime.DynamicType:
			return t, true
		case runtime.IntegerType, runtime.RealType, runtime.NumericType:
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
		possiblyAbsent := false
		for index, arg := range call.Args[:min(len(call.Args), 2)] {
			typ := c.expr(arg)
			switch typ.Kind {
			case runtime.VoidType, runtime.StringType, runtime.BoolType:
				return runtime.Type{Kind: runtime.VoidType}, true
			case runtime.VectorType:
				c.error(arg.Start(), diag.EParse, "expected '[' after vector")
			}
			absent, numericAbsent := c.possibleAbsence(arg)
			if index == 1 && numericAbsent && len(call.Args) > 2 {
				c.expr(call.Args[2])
			}
			possiblyAbsent = possiblyAbsent || absent
		}
		if len(call.Args) != 2 && !possiblyAbsent {
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
	case "caracpnum":
		c.textArgs(call, 1, 1)
		return runtime.Type{Kind: runtime.NumericType}, true
	case "randi":
		if len(call.Args) > 1 {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after random bound")
		} else if len(call.Args) == 1 {
			switch c.expr(call.Args[0]).Kind {
			case runtime.VectorType:
				c.error(call.Args[0].Start(), diag.EParse, "expected '[' after vector")
			case runtime.IntegerType, runtime.NumericType, runtime.DynamicType, runtime.InvalidType:
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
		c.textArgs(call, 3, 1)
		return runtime.Type{Kind: runtime.StringType}, true
	case "maiusc", "minusc":
		c.textArgs(call, 1, 1)
		return runtime.Type{Kind: runtime.StringType}, true
	case "asc":
		c.textArgs(call, 1, 1)
		return runtime.Type{Kind: runtime.IntegerType}, true
	case "carac":
		if len(call.Args) > 1 {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after character code")
		} else if len(call.Args) == 1 {
			t := c.expr(call.Args[0])
			if t.Kind != runtime.IntegerType && t.Kind != runtime.NumericType && t.Kind != runtime.VoidType && t.Kind != runtime.InvalidType {
				c.error(call.Args[0].Start(), diag.ETypeMismatch, "expected inteiro argument")
			}
		}
		return runtime.Type{Kind: runtime.StringType}, true
	case "compr":
		c.textArgs(call, 1, 1)
		return runtime.Type{Kind: runtime.IntegerType}, true
	case "pos":
		c.textArgs(call, 2, 2)
		return runtime.Type{Kind: runtime.IntegerType}, true
	default:
		return runtime.Type{}, false
	}
}

// possibleAbsence distinguishes generic absence from a numeric domain result.
func (c *checker) possibleAbsence(expr ast.Expr) (possible, numeric bool) {
	if c.info.types[expr].Kind == runtime.VoidType {
		return true, false
	}
	switch e := expr.(type) {
	case *ast.UnaryExpr:
		if e.Op.Kind == token.ADD {
			return c.possibleAbsence(e.X)
		}
	case *ast.BinaryExpr:
		switch e.Op.Kind {
		case token.QUO, token.IDIV:
			return c.possibleAbsence(e.Right)
		case token.REM, token.MOD:
			if c.info.types[e.Left].Kind == runtime.BoolType {
				return c.possibleAbsence(e.Right)
			}
		case token.POW:
			left, _ := c.possibleAbsence(e.Left)
			right, _ := c.possibleAbsence(e.Right)
			return left || right, false
		}
	}
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false, false
	}
	binding, ok := c.info.Binding(call.Name)
	if !ok || !binding.Builtin {
		return false, false
	}
	switch binding.Name {
	case "arccos", "arcsen", "cotan", "radpgrau":
		return true, true
	case "asc", "carac":
		return true, false
	}
	for _, arg := range call.Args {
		absent, numericAbsent := c.possibleAbsence(arg)
		possible = possible || absent
		numeric = numeric || numericAbsent
	}
	if binding.Name == "exp" || binding.Name == "numpcarac" {
		numeric = false
	}
	return possible, numeric
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

func (c *checker) textArgs(call *ast.CallExpr, arity, strings int) {
	for index, arg := range call.Args[:min(len(call.Args), arity)] {
		t := c.expr(arg)
		if t.Kind == runtime.InvalidType {
			return
		}
		if t.Kind == runtime.DynamicType {
			continue
		}
		if index < strings {
			if t.Kind != runtime.StringType {
				c.error(arg.Start(), diag.ETypeMismatch, "expected caractere argument")
				return
			}
		} else if !isNumeric(t) && t.Kind != runtime.VoidType {
			c.error(arg.Start(), diag.ETypeMismatch, "expected numeric argument")
			return
		}
	}
	if len(call.Args) < strings {
		c.error(call.Name.Pos, diag.ETypeMismatch, "expected caractere argument")
	} else if len(call.Args) != arity {
		c.error(call.Name.Pos, diag.EParse, "expected ')' after text arguments")
	}
}
