package interp

import (
	"cmp"
	"fmt"
	"math"
	"slices"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/cp1252"
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
	// Call frames reset the depth guard but share the enclosing expression's
	// retained operands. Release those values when the outer evaluation ends.
	i.evalFrames++
	defer func() {
		i.evalFrames--
		if i.evalFrames == 0 {
			clear(i.operands)
			i.operands = i.operands[:0]
		}
	}()
	i.depth++
	defer func() {
		i.depth--
		err = failure(expr.Start(), diag.RType, err)
		if err == nil && value.Kind == runtime.StringValue {
			value.Str = runtime.LimitText(value.Str)
		}
		// Numeric reference parameters can change a caller's runtime type.
		actual := value.Type()
		// Numeric built-ins can produce no value for a domain failure.
		if err == nil && !value.Comparison && actual.Kind != runtime.VoidType && !runtime.Assignable(typ, actual) && !runtime.Assignable(actual, typ) {
			err = failure(expr.Start(), diag.RType, fmt.Errorf("inconsistent expression value"))
		}
	}()
	switch e := expr.(type) {
	case *ast.NoValueExpr:
		return runtime.Value{Kind: runtime.VoidValue}, nil
	case *ast.LiteralExpr:
		if len(e.Str) > maxTextBytes {
			return value, failure(e.At, diag.RStorage, fmt.Errorf("text size limit exceeded"))
		}
		return literalValue(e), nil
	case *ast.IdentExpr:
		if binding, ok := i.info.Binding(e.Name); ok && (binding.Builtin || i.subs[binding.ID] != nil) {
			return i.callFunction(&ast.CallExpr{Name: e.Name})
		}
		cell, err := i.lookupCell(e.Name)
		if err != nil {
			return runtime.Value{}, err
		}
		return cell.Value, nil
	case *ast.IndexExpr, *ast.FieldExpr:
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
	case token.ADD:
		if v.Kind == runtime.VoidValue {
			return v, nil
		}
		if v.Comparison || v.RealFallback && v.Kind != runtime.RealValue {
			return runtime.Value{}, failure(e.Op.Pos, diag.EParse, fmt.Errorf("expected numeric operand"))
		}
		if v.Kind != runtime.IntegerValue && v.Kind != runtime.RealValue {
			return runtime.Value{}, failure(e.Op.Pos, diag.ETypeMismatch, fmt.Errorf("expected numeric operand"))
		}
		return v, nil
	case token.SUB:
		if v.Comparison {
			return runtime.Value{Kind: runtime.VoidValue}, nil
		}
		if v.RealFallback && v.Kind == runtime.IntegerValue {
			return runtime.Value{Kind: runtime.RealValue, Real: -math.Float64frombits(uint64(uint32(v.Int)))}, nil
		}
		switch v.Kind {
		case runtime.IntegerValue:
			return runtime.Value{Kind: runtime.IntegerValue, Int: int64(-int32(v.Int))}, nil
		case runtime.RealValue:
			return runtime.Value{Kind: runtime.RealValue, Real: -v.Real}, nil
		case runtime.StringValue, runtime.BoolValue:
			return runtime.Value{Kind: runtime.VoidValue}, nil
		}
	case token.NAO:
		if v.Kind == runtime.BoolValue || v.Comparison {
			return runtime.Value{Kind: runtime.BoolValue, Bool: !v.Bool}, nil
		}
		return runtime.Value{}, failure(e.Op.Pos, diag.ETypeMismatch, fmt.Errorf("expected logico"))
	}
	return runtime.Value{}, fmt.Errorf("invalid unary operator %s", e.Op.Text)
}

func (i *Interpreter) evalBinary(e *ast.BinaryExpr) (value runtime.Value, err error) {
	defer func() {
		if err == nil && value.Kind == runtime.RealValue && (math.IsNaN(value.Real) || math.IsInf(value.Real, 0)) {
			err = fmt.Errorf("invalid result for operator %s", e.Op.Text)
		}
		err = failure(e.Op.Pos, diag.RArithmetic, err)
	}()
	left, err := i.eval(e.Left)
	if err != nil {
		return runtime.Value{}, err
	}
	if e.IsComparison() && left.NumericAbsence {
		return left, nil
	}
	if len(i.operands) == maxOperands {
		return runtime.Value{}, failure(e.Op.Pos, diag.RStorage, fmt.Errorf("retained operand limit exceeded"))
	}
	i.operands = append(i.operands, left)
	right, err := i.eval(e.Right)
	if err != nil {
		return runtime.Value{}, err
	}
	// Comparisons keep their own left value. Other operators reduce the most
	// recent operand, including one retained by a nested mixed expression.
	if !e.IsComparison() {
		left = i.operands[len(i.operands)-1]
	}
	keepLeft := false
	defer func() {
		if !keepLeft {
			last := len(i.operands) - 1
			i.operands[last] = runtime.Value{}
			i.operands = i.operands[:last]
		}
	}()
	switch e.Op.Kind {
	case token.ADD:
		if left.Kind == runtime.StringValue && right.Kind == runtime.StringValue {
			if len(left.Str) > maxTextBytes-len(right.Str) {
				return value, failure(e.Op.Pos, diag.RStorage, fmt.Errorf("text size limit exceeded"))
			}
			return runtime.Value{Kind: runtime.StringValue, Str: left.Str + right.Str}, nil
		}
		if _, _, err := floats(left, right); err != nil {
			return runtime.Value{}, failure(e.Op.Pos, diag.EParse, err)
		}
		return numeric(left, right, func(a, b int64) int64 { return a + b }, func(a, b float64) float64 { return a + b })
	case token.SUB:
		return numeric(left, right, func(a, b int64) int64 { return a - b }, func(a, b float64) float64 { return a - b })
	case token.MUL:
		if _, _, err := floats(left, right); err != nil {
			keepLeft = true
			right.Comparison = false
			return right, nil
		}
		return numeric(left, right, func(a, b int64) int64 { return a * b }, func(a, b float64) float64 { return a * b })
	case token.QUO:
		a, b, err := floats(left, right)
		if err != nil {
			right.Comparison = false
			right.RealFallback = true
			return right, nil
		}
		if b == 0 {
			return runtime.Value{}, fmt.Errorf("division by zero")
		}
		return runtime.Value{Kind: runtime.RealValue, Real: a / b}, nil
	case token.IDIV:
		if left.Kind != runtime.IntegerValue || right.Kind != runtime.IntegerValue {
			right.Comparison = false
			return right, nil
		}
		a, b := left.Int, right.Int
		if b == 0 {
			return runtime.Value{}, fmt.Errorf("division by zero")
		}
		if a == math.MinInt32 && b == -1 {
			return runtime.Value{}, fmt.Errorf("integer division overflow")
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: a / b}, nil
	case token.REM, token.MOD:
		if left.Kind == runtime.BoolValue || right.Kind == runtime.BoolValue {
			right.Comparison = false
			return right, nil
		}
		if _, _, err := floats(left, right); err != nil {
			return runtime.Value{}, err
		}
		if right.Kind != runtime.IntegerValue || right.Int <= 0 {
			return runtime.Value{Kind: runtime.IntegerValue, Int: -1}, nil
		}
		a := int64(int32(left.Int))
		if left.Kind == runtime.RealValue {
			converted, _, err := i.lib.Call("int", []runtime.Value{left})
			if err != nil {
				return runtime.Value{}, err
			}
			a = converted.Int
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: a % right.Int}, nil
	case token.POW:
		a, b, err := floats(left, right)
		if err != nil {
			return runtime.Value{Kind: runtime.VoidValue}, nil
		}
		return runtime.Value{Kind: runtime.RealValue, Real: math.Pow(a, b)}, nil
	case token.EQL, token.NEQ, token.LSS, token.GTR, token.LEQ, token.GEQ:
		if right.NumericAbsence {
			return runtime.Value{}, failure(e.Op.Pos, diag.EParse, fmt.Errorf("expected comparison operand"))
		}
		value, consumed, err := compare(e.Op.Kind, left, right)
		keepLeft = !consumed
		value.Comparison = true
		value.RealFallback = false
		return value, err
	case token.E, token.OU, token.XOU:
		if left.Kind != runtime.BoolValue || right.Kind != runtime.BoolValue {
			if e.Op.Kind == token.E {
				keepLeft = true
				return right, nil
			}
			return runtime.Value{}, failure(e.Op.Pos, diag.EParse, fmt.Errorf("expected logico operands"))
		}
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
	case *ast.FieldExpr:
		base, err := i.eval(e.X)
		if err != nil {
			return nil, err
		}
		baseType, _ := i.info.TypeOf(e.X)
		if base.Kind != runtime.RecordValue || base.Rec == nil {
			if baseType.Kind == runtime.DynamicType {
				return nil, failure(e.Name.Pos, diag.EUndeclared, fmt.Errorf("unknown field %q", e.Name.Text))
			}
			return nil, fmt.Errorf("cannot select a field of a non-record")
		}
		cell, err := base.Rec.Cell(e.Name.Text)
		if err != nil && baseType.Kind == runtime.DynamicType {
			err = failure(e.Name.Pos, diag.EUndeclared, err)
		}
		return cell, err
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
			indices[idx], err = i.evalInt(expr, diag.ETypeMismatch)
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
	truth, err := runtime.Truth(v)
	return truth, failure(expr.Start(), diag.ETypeMismatch, err)
}

func (i *Interpreter) evalInt(expr ast.Expr, code diag.Code) (int64, error) {
	v, err := i.eval(expr)
	if err != nil {
		return 0, err
	}
	if v.Kind != runtime.IntegerValue || v.Comparison {
		return 0, failure(expr.Start(), code, fmt.Errorf("expected inteiro"))
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

// compare also reports whether the reduction consumes a retained operand.
func compare(op token.Kind, left, right runtime.Value) (runtime.Value, bool, error) {
	if right.Kind == runtime.VoidValue && left.Kind != runtime.VoidValue {
		return left, true, nil
	}
	if left.Comparison {
		if !right.Comparison && right.Kind != runtime.BoolValue {
			return right, false, nil
		}
		left = runtime.Value{Kind: runtime.BoolValue, Bool: left.Bool}
	}
	if right.Comparison {
		if left.Kind != runtime.BoolValue {
			return right, false, nil
		}
		right = runtime.Value{Kind: runtime.BoolValue, Bool: right.Bool}
	}
	if left.Kind == runtime.RecordValue || right.Kind == runtime.RecordValue {
		return right, false, nil
	}
	if left.Kind != right.Kind {
		if _, _, err := floats(left, right); err != nil {
			return right, false, nil
		}
	}
	if op == token.EQL || op == token.NEQ {
		eq, err := equalValues(left, right)
		if err != nil {
			return runtime.Value{}, true, err
		}
		if op == token.NEQ {
			eq = !eq
		}
		return runtime.Value{Kind: runtime.BoolValue, Bool: eq}, true, nil
	}
	cmp, err := ordering(left, right)
	if err != nil {
		return runtime.Value{}, true, err
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
	return runtime.Value{Kind: runtime.BoolValue, Bool: ok}, true, nil
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
	if left.Kind == runtime.BoolValue && right.Kind == runtime.BoolValue {
		if left.Bool == right.Bool {
			return 0, nil
		}
		if left.Bool {
			return 1, nil
		}
		return -1, nil
	}
	if left.Kind == runtime.StringValue && right.Kind == runtime.StringValue {
		return slices.CompareFunc([]rune(left.Str), []rune(right.Str), func(a, b rune) int {
			x, y := int(a)+256, int(b)+256
			if code, ok := cp1252.EncodeRune(a); ok {
				x = int(code)
			}
			if code, ok := cp1252.EncodeRune(b); ok {
				y = int(code)
			}
			return cmp.Compare(x, y)
		}), nil
	}
	return 0, fmt.Errorf("values are not ordered")
}
