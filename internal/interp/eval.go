package interp

import (
	"fmt"
	"math"
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func (i *Interpreter) eval(expr ast.Expr) (value runtime.Value, err error) {
	if expr == nil {
		return value, failure(token.NoPos, diag.RType, fmt.Errorf("missing expression"))
	}
	if err := i.charge(expr.Start()); err != nil {
		return value, err
	}
	typ, ok := i.info.TypeOf(expr)
	if !ok {
		return value, failure(expr.Start(), diag.RType, fmt.Errorf("missing expression type"))
	}
	if i.depth == maxDepth {
		return value, failure(expr.Start(), diag.RStorage, fmt.Errorf("expression depth limit exceeded"))
	}
	i.depth++
	defer func() {
		i.depth--
		err = failure(expr.Start(), diag.RType, err)
		// Numeric reference parameters can change a caller's runtime type.
		actual := value.Type()
		if err == nil && !runtime.Assignable(typ, actual) && !runtime.Assignable(actual, typ) {
			err = failure(expr.Start(), diag.RType, fmt.Errorf("inconsistent expression value"))
		}
	}()
	switch e := expr.(type) {
	case *ast.LiteralExpr:
		if len(e.Str) > maxTextBytes {
			return value, failure(e.At, diag.RStorage, fmt.Errorf("text size limit exceeded"))
		}
		return literalValue(e), nil
	case *ast.IdentExpr:
		if binding, ok := i.info.Binding(e.Name); ok && i.subs[binding.ID] != nil {
			return i.callFunction(&ast.CallExpr{Name: e.Name})
		}
		cell, err := i.lookupCell(e.Name)
		if err != nil {
			return runtime.Value{}, err
		}
		return cell.Value, nil
	case *ast.IndexExpr:
		cell, err := i.location(e)
		if err != nil {
			return runtime.Value{}, err
		}
		return cell.Value, nil
	case *ast.UnaryExpr:
		return i.evalUnary(e)
	case *ast.BinaryExpr:
		return i.evalBinary(e)
	case *ast.CallExpr:
		return i.callFunction(e)
	default:
		return runtime.Value{}, fmt.Errorf("unsupported expression %T", expr)
	}
}

func literalValue(e *ast.LiteralExpr) runtime.Value {
	switch e.Kind {
	case ast.IntLiteral:
		return runtime.Value{Kind: runtime.IntegerValue, Int: e.Int}
	case ast.RealLiteral:
		return runtime.Value{Kind: runtime.RealValue, Real: e.Real}
	case ast.StringLiteral:
		return runtime.Value{Kind: runtime.StringValue, Str: e.Str}
	case ast.BoolLiteral:
		return runtime.Value{Kind: runtime.BoolValue, Bool: e.Bool}
	default:
		return runtime.Value{Kind: runtime.InvalidValue}
	}
}

func (i *Interpreter) evalUnary(e *ast.UnaryExpr) (runtime.Value, error) {
	v, err := i.eval(e.X)
	if err != nil {
		return runtime.Value{}, err
	}
	switch e.Op.Kind {
	case token.SUB:
		switch v.Kind {
		case runtime.IntegerValue:
			return runtime.Value{Kind: runtime.IntegerValue, Int: int64(-int32(v.Int))}, nil
		case runtime.RealValue:
			return runtime.Value{Kind: runtime.RealValue, Real: -v.Real}, nil
		}
	case token.NAO:
		if v.Kind == runtime.BoolValue {
			return runtime.Value{Kind: runtime.BoolValue, Bool: !v.Bool}, nil
		}
	}
	return runtime.Value{}, fmt.Errorf("invalid unary operator %s", e.Op.Text)
}

func (i *Interpreter) evalBinary(e *ast.BinaryExpr) (value runtime.Value, err error) {
	defer func() { err = failure(e.Op.Pos, diag.RArithmetic, err) }()
	left, err := i.eval(e.Left)
	if err != nil {
		return runtime.Value{}, err
	}
	right, err := i.eval(e.Right)
	if err != nil {
		return runtime.Value{}, err
	}
	switch e.Op.Kind {
	case token.ADD:
		if left.Kind == runtime.StringValue && right.Kind == runtime.StringValue {
			if len(left.Str) > maxTextBytes-len(right.Str) {
				return value, failure(e.Op.Pos, diag.RStorage, fmt.Errorf("text size limit exceeded"))
			}
			return runtime.Value{Kind: runtime.StringValue, Str: left.Str + right.Str}, nil
		}
		return numeric(left, right, func(a, b int64) int64 { return a + b }, func(a, b float64) float64 { return a + b })
	case token.SUB:
		return numeric(left, right, func(a, b int64) int64 { return a - b }, func(a, b float64) float64 { return a - b })
	case token.MUL:
		return numeric(left, right, func(a, b int64) int64 { return a * b }, func(a, b float64) float64 { return a * b })
	case token.QUO:
		a, b, err := floats(left, right)
		if err != nil {
			return runtime.Value{}, err
		}
		if b == 0 {
			return runtime.Value{}, fmt.Errorf("division by zero")
		}
		return runtime.Value{Kind: runtime.RealValue, Real: a / b}, nil
	case token.IDIV:
		a, b, err := ints(left, right)
		if err != nil {
			return runtime.Value{}, err
		}
		if b == 0 {
			return runtime.Value{}, fmt.Errorf("division by zero")
		}
		if a == math.MinInt32 && b == -1 {
			return runtime.Value{}, fmt.Errorf("integer division overflow")
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: a / b}, nil
	case token.REM, token.MOD:
		a, b, err := ints(left, right)
		if err != nil {
			return runtime.Value{}, err
		}
		if b == 0 {
			return runtime.Value{}, fmt.Errorf("modulo by zero")
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: a % b}, nil
	case token.POW:
		a, b, err := floats(left, right)
		if err != nil {
			return runtime.Value{}, err
		}
		return runtime.Value{Kind: runtime.RealValue, Real: math.Pow(a, b)}, nil
	case token.EQL, token.NEQ, token.LSS, token.GTR, token.LEQ, token.GEQ:
		return compare(e.Op.Kind, left, right)
	case token.E, token.OU, token.XOU:
		a, err := runtime.Truth(left)
		if err != nil {
			return runtime.Value{}, err
		}
		b, err := runtime.Truth(right)
		if err != nil {
			return runtime.Value{}, err
		}
		switch e.Op.Kind {
		case token.E:
			return runtime.Value{Kind: runtime.BoolValue, Bool: a && b}, nil
		case token.OU:
			return runtime.Value{Kind: runtime.BoolValue, Bool: a || b}, nil
		default:
			return runtime.Value{Kind: runtime.BoolValue, Bool: a != b}, nil
		}
	}
	return runtime.Value{}, fmt.Errorf("unsupported operator %s", e.Op.Text)
}

func (i *Interpreter) lvalue(expr ast.Expr) (cell *runtime.Cell, err error) {
	if err := i.charge(expr.Start()); err != nil {
		return nil, err
	}
	return i.location(expr)
}

func (i *Interpreter) location(expr ast.Expr) (cell *runtime.Cell, err error) {
	if _, ok := i.info.TypeOf(expr); !ok {
		return nil, failure(expr.Start(), diag.RType, fmt.Errorf("missing designator type"))
	}
	defer func() { err = failure(expr.Start(), diag.RStorage, err) }()
	switch e := expr.(type) {
	case *ast.IdentExpr:
		return i.lookupCell(e.Name)
	case *ast.IndexExpr:
		base, err := i.eval(e.X)
		if err != nil {
			return nil, err
		}
		if base.Kind != runtime.VectorValue || base.Vec == nil {
			return nil, fmt.Errorf("cannot index non-vector")
		}
		indices := make([]int64, len(e.Indices))
		for idx, expr := range e.Indices {
			indices[idx], err = i.evalInt(expr)
			if err != nil {
				return nil, err
			}
		}
		return base.Vec.Cell(indices)
	default:
		return nil, fmt.Errorf("expression is not assignable")
	}
}

func (i *Interpreter) evalBool(expr ast.Expr) (bool, error) {
	v, err := i.eval(expr)
	if err != nil {
		return false, err
	}
	return runtime.Truth(v)
}

func (i *Interpreter) evalInt(expr ast.Expr) (int64, error) {
	v, err := i.eval(expr)
	if err != nil {
		return 0, err
	}
	if v.Kind != runtime.IntegerValue {
		return 0, fmt.Errorf("expected inteiro")
	}
	return v.Int, nil
}

func numeric(left, right runtime.Value, intFn func(int64, int64) int64, realFn func(float64, float64) float64) (runtime.Value, error) {
	if left.Kind == runtime.IntegerValue && right.Kind == runtime.IntegerValue {
		return runtime.Value{Kind: runtime.IntegerValue, Int: int64(int32(intFn(left.Int, right.Int)))}, nil
	}
	a, b, err := floats(left, right)
	if err != nil {
		return runtime.Value{}, err
	}
	return runtime.Value{Kind: runtime.RealValue, Real: realFn(a, b)}, nil
}

func ints(left, right runtime.Value) (int64, int64, error) {
	if left.Kind != runtime.IntegerValue || right.Kind != runtime.IntegerValue {
		return 0, 0, fmt.Errorf("expected inteiro operands")
	}
	return left.Int, right.Int, nil
}

func floats(left, right runtime.Value) (float64, float64, error) {
	a, ok := asFloat(left)
	if !ok {
		return 0, 0, fmt.Errorf("expected numeric operand")
	}
	b, ok := asFloat(right)
	if !ok {
		return 0, 0, fmt.Errorf("expected numeric operand")
	}
	return a, b, nil
}

func asFloat(v runtime.Value) (float64, bool) {
	switch v.Kind {
	case runtime.IntegerValue:
		return float64(v.Int), true
	case runtime.RealValue:
		return v.Real, true
	default:
		return 0, false
	}
}

func compare(op token.Kind, left, right runtime.Value) (runtime.Value, error) {
	if op == token.EQL || op == token.NEQ {
		eq, err := equalValues(left, right)
		if err != nil {
			return runtime.Value{}, err
		}
		if op == token.NEQ {
			eq = !eq
		}
		return runtime.Value{Kind: runtime.BoolValue, Bool: eq}, nil
	}
	cmp, err := ordering(left, right)
	if err != nil {
		return runtime.Value{}, err
	}
	var ok bool
	switch op {
	case token.LSS:
		ok = cmp < 0
	case token.GTR:
		ok = cmp > 0
	case token.LEQ:
		ok = cmp <= 0
	case token.GEQ:
		ok = cmp >= 0
	}
	return runtime.Value{Kind: runtime.BoolValue, Bool: ok}, nil
}

func equalValues(left, right runtime.Value) (bool, error) {
	if left.Kind == runtime.IntegerValue && right.Kind == runtime.RealValue {
		return float64(left.Int) == right.Real, nil
	}
	if left.Kind == runtime.RealValue && right.Kind == runtime.IntegerValue {
		return left.Real == float64(right.Int), nil
	}
	if left.Kind != right.Kind {
		return false, fmt.Errorf("cannot compare %s with %s", left.Type(), right.Type())
	}
	switch left.Kind {
	case runtime.IntegerValue:
		return left.Int == right.Int, nil
	case runtime.RealValue:
		return left.Real == right.Real, nil
	case runtime.StringValue:
		return left.Str == right.Str, nil
	case runtime.BoolValue:
		return left.Bool == right.Bool, nil
	default:
		return false, fmt.Errorf("cannot compare %s", left.Type())
	}
}

func ordering(left, right runtime.Value) (int, error) {
	if a, b, err := floats(left, right); err == nil {
		switch {
		case a < b:
			return -1, nil
		case a > b:
			return 1, nil
		default:
			return 0, nil
		}
	}
	if left.Kind == runtime.StringValue && right.Kind == runtime.StringValue {
		return strings.Compare(left.Str, right.Str), nil
	}
	return 0, fmt.Errorf("values are not ordered")
}
