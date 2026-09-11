package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/stdlib"
	"github.com/ncode/portugol-go/internal/token"
)

func (c *checker) declareBuiltins() {
	for _, descriptor := range stdlib.Catalog() {
		for _, name := range descriptor.Names() {
			c.scope.declare(symbol{name: name, kind: builtinSym})
		}
	}
}

func (c *checker) builtinCallType(name string, call *ast.CallExpr) (runtime.Type, bool) {
	descriptor, ok := stdlib.Lookup(name)
	if !ok {
		return runtime.Type{}, false
	}
	signature := descriptor.Signature()
	result := runtime.Type{Kind: signature.Result}
	switch signature.Rule {
	case stdlib.ConstantCall:
		c.requireArity(call, signature.Arity, signature.Arity)
	case stdlib.UnaryNumericCall:
		if len(call.Args) == 0 {
			return runtime.Type{Kind: signature.Empty}, true
		}
		if len(call.Args) > signature.Arity && signature.ArityFirst {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after numeric argument")
			return result, true
		}
		t := c.expr(call.Args[0])
		switch t.Kind {
		case runtime.DynamicType:
			return t, true
		case runtime.IntegerType, runtime.RealType, runtime.NumericType:
			if signature.PreserveNumeric {
				return t, true
			}
		case runtime.VoidType:
			return t, true
		case runtime.StringType, runtime.BoolType:
			if signature.NonNumericAbsent {
				return runtime.Type{Kind: runtime.VoidType}, true
			}
			c.error(call.Args[0].Start(), diag.EParse, "expected numeric expression")
		case runtime.VectorType:
			c.error(call.Args[0].Start(), diag.EParse, "expected '[' after vector")
		}
		if len(call.Args) > signature.Arity {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after numeric argument")
		}
	case stdlib.PowerCall:
		if len(call.Args) == 0 {
			return runtime.Type{Kind: signature.Empty}, true
		}
		possiblyAbsent := false
		for index, arg := range call.Args[:min(len(call.Args), signature.Arity)] {
			typ := c.expr(arg)
			switch typ.Kind {
			case runtime.VoidType, runtime.StringType, runtime.BoolType:
				return runtime.Type{Kind: runtime.VoidType}, true
			case runtime.VectorType:
				c.error(arg.Start(), diag.EParse, "expected '[' after vector")
			}
			absent, numericAbsent := c.possibleAbsence(arg)
			if index == signature.Arity-1 && numericAbsent && len(call.Args) > signature.Arity {
				c.expr(call.Args[signature.Arity])
			}
			possiblyAbsent = possiblyAbsent || absent
		}
		if len(call.Args) != signature.Arity && !possiblyAbsent {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after numeric argument")
		}
	case stdlib.NumberTextCall:
		if len(call.Args) > signature.Arity {
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
	case stdlib.RandomIntegerCall:
		if len(call.Args) > signature.Arity {
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
	case stdlib.CharacterCall:
		if len(call.Args) > signature.Arity {
			c.error(call.Name.Pos, diag.EParse, "expected ')' after character code")
		} else if len(call.Args) == 1 {
			t := c.expr(call.Args[0])
			if t.Kind != runtime.IntegerType && t.Kind != runtime.NumericType && t.Kind != runtime.VoidType && t.Kind != runtime.InvalidType {
				c.error(call.Args[0].Start(), diag.ETypeMismatch, "expected inteiro argument")
			}
		}
	case stdlib.TextCall:
		c.textArgs(call, signature)
	}
	return result, true
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
	descriptor, ok := stdlib.Lookup(binding.Name)
	if !ok {
		return false, false
	}
	switch descriptor.Domain().Absence {
	case stdlib.NumericDomainAbsence:
		return true, true
	case stdlib.GenericDomainAbsence:
		return true, false
	}
	for _, arg := range call.Args {
		absent, numericAbsent := c.possibleAbsence(arg)
		possible = possible || absent
		numeric = numeric || numericAbsent
	}
	if descriptor.Signature().ClearsNumericAbsence() {
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

func (c *checker) textArgs(call *ast.CallExpr, signature stdlib.Signature) {
	strings := 0
	for _, parameter := range signature.Parameters[:signature.Arity] {
		if parameter.Type == runtime.StringType {
			strings++
		}
	}
	for index, arg := range call.Args[:min(len(call.Args), signature.Arity)] {
		t := c.expr(arg)
		if t.Kind == runtime.InvalidType {
			return
		}
		if t.Kind == runtime.DynamicType {
			continue
		}
		parameter := signature.Parameters[index]
		if parameter.Type == runtime.StringType {
			if t.Kind != runtime.StringType {
				c.error(arg.Start(), diag.ETypeMismatch, "expected caractere argument")
				return
			}
		} else if !isNumeric(t) && (!parameter.AbsenceAsZero || t.Kind != runtime.VoidType) {
			c.error(arg.Start(), diag.ETypeMismatch, "expected numeric argument")
			return
		}
	}
	if len(call.Args) < strings {
		c.error(call.Name.Pos, diag.ETypeMismatch, "expected caractere argument")
	} else if len(call.Args) != signature.Arity {
		c.error(call.Name.Pos, diag.EParse, "expected ')' after text arguments")
	}
}
