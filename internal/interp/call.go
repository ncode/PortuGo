package interp

import (
	"errors"
	"fmt"
	"math"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func (i *Interpreter) callFunction(call *ast.CallExpr) (value runtime.Value, err error) {
	defer func() { err = failure(call.Start(), diag.RCall, err) }()
	b, ok := i.info.Binding(call.Name)
	if !ok {
		return value, failure(call.Start(), diag.RType, fmt.Errorf("missing call binding"))
	}
	if b.Builtin {
		exprs := call.Args
		if b.Name == "exp" {
			exprs = exprs[:min(len(exprs), 2)]
		}
		args := make([]runtime.Value, len(exprs))
		for idx, arg := range exprs {
			v, err := i.eval(arg)
			if err != nil {
				return runtime.Value{}, err
			}
			if v.Kind == runtime.VoidValue && (b.Name == "carac" || b.Name == "copia" && idx > 0) {
				v = runtime.Value{Kind: runtime.IntegerValue}
			}
			textArg := b.Name == "maiusc" || b.Name == "minusc" || b.Name == "asc" || b.Name == "compr" || b.Name == "pos" || b.Name == "caracpnum" || b.Name == "copia" && idx == 0
			if textArg && v.Kind != runtime.StringValue {
				return runtime.Value{}, failure(arg.Start(), diag.ETypeMismatch, fmt.Errorf("expected caractere argument"))
			}
			if b.Name == "copia" && idx > 0 && v.Kind != runtime.IntegerValue && v.Kind != runtime.RealValue {
				return runtime.Value{}, failure(arg.Start(), diag.ETypeMismatch, fmt.Errorf("expected numeric argument"))
			}
			if (b.Name == "carac" || b.Name == "randi") && v.Kind != runtime.IntegerValue {
				return runtime.Value{}, failure(arg.Start(), diag.ETypeMismatch, fmt.Errorf("expected inteiro argument"))
			}
			if v.Kind == runtime.VoidValue || (b.Name == "exp" || b.Name == "int") && (v.Kind == runtime.StringValue || v.Kind == runtime.BoolValue) {
				// An absent numeric exponent consumes one trailing expression.
				if b.Name == "exp" && idx == 1 && v.NumericAbsence && len(call.Args) > 2 {
					if _, err := i.eval(call.Args[2]); err != nil {
						return runtime.Value{}, err
					}
				}
				if v.Kind == runtime.VoidValue && b.Name != "exp" && b.Name != "numpcarac" {
					return v, nil
				}
				return runtime.Value{Kind: runtime.VoidValue}, nil
			}
			v.Comparison = false
			v.RealFallback = false
			args[idx] = v
		}
		if b.Name == "exp" && len(call.Args) != 0 && len(call.Args) != 2 {
			return runtime.Value{}, failure(call.Start(), diag.EParse, fmt.Errorf("expected ')' after numeric argument"))
		}
		if v, ok, err := i.lib.Call(b.Name, args); ok {
			if errors.Is(err, runtime.ErrTextSize) {
				return v, failure(call.Start(), diag.RStorage, err)
			}
			return v, failure(call.Start(), diag.RBuiltin, err)
		}
	}
	sub, ok := i.subs[b.ID]
	if !ok {
		return runtime.Value{}, fmt.Errorf("undefined function %q", call.Name.Text)
	}
	fn, ok := sub.(*ast.FunctionDecl)
	if !ok {
		return runtime.Value{}, fmt.Errorf("%q is not a function", call.Name.Text)
	}
	return i.callUserFunction(fn, call.Args)
}

func (i *Interpreter) callProcedure(call *ast.CallExpr) (err error) {
	if err := i.charge(call.Start()); err != nil {
		return err
	}
	defer func() { err = failure(call.Start(), diag.RCall, err) }()
	b, ok := i.info.Binding(call.Name)
	if !ok {
		return failure(call.Start(), diag.RType, fmt.Errorf("missing call binding"))
	}
	if b.Builtin {
		_, err := i.callFunction(call)
		return err
	}
	sub, ok := i.subs[b.ID]
	if !ok {
		return fmt.Errorf("undefined procedure %q", call.Name.Text)
	}
	proc, ok := sub.(*ast.ProcedureDecl)
	if !ok {
		return fmt.Errorf("%q is not a procedure", call.Name.Text)
	}
	_, err = i.callSub(proc.Params, proc.Console, proc.Consts, proc.Locals, proc.Body, runtime.Type{Kind: runtime.VoidType}, call.Args)
	return err
}

func (i *Interpreter) callUserFunction(fn *ast.FunctionDecl, args []ast.Expr) (runtime.Value, error) {
	b, ok := i.info.Binding(fn.Name)
	if !ok {
		return runtime.Value{}, failure(fn.Start(), diag.RType, fmt.Errorf("missing return type"))
	}
	return i.callSub(fn.Params, fn.Console, fn.Consts, fn.Locals, fn.Body, b.Type, args)
}

func (i *Interpreter) callSub(params []ast.Param, console []ast.ConsoleStmt, consts []ast.ConstDecl, locals []ast.VarDecl, body []ast.Stmt, retType runtime.Type, args []ast.Expr) (runtime.Value, error) {
	if i.calls == maxCalls {
		return runtime.Value{}, fmt.Errorf("active call limit exceeded")
	}
	if len(args) != len(params) {
		return runtime.Value{}, fmt.Errorf("expected %d arguments, got %d", len(params), len(args))
	}
	outer := i.env
	callEnv := newEnv(i.global)
	var references []struct{ caller, parameter *runtime.Cell }
	for idx, param := range params {
		b, ok := i.info.Binding(param.Name)
		if !ok {
			return runtime.Value{}, failure(param.Name.Pos, diag.RType, fmt.Errorf("missing parameter layout"))
		}
		typ := b.Type
		var v runtime.Value
		var caller *runtime.Cell
		var err error
		if param.ByRef {
			caller, err = i.lvalue(args[idx])
			if err == nil {
				v = caller.Value
			}
		} else {
			v, err = i.eval(args[idx])
		}
		if err != nil {
			return runtime.Value{}, err
		}
		v.Comparison = false
		v.RealFallback = false
		if typ.Kind == runtime.IntegerType && v.Kind == runtime.RealValue {
			if math.IsNaN(v.Real) || v.Real < -0x1p63 || v.Real >= 0x1p63 {
				return runtime.Value{}, fmt.Errorf("integer argument out of range")
			}
			v = runtime.Value{Kind: runtime.IntegerValue, Int: int64(v.Real)}
		}
		cell := callEnv.define(b.ID, typ)
		if err := assign(cell, v); err != nil {
			return runtime.Value{}, err
		}
		if caller != nil {
			references = append(references, struct{ caller, parameter *runtime.Cell }{caller, cell})
		}
	}
	i.env = callEnv
	depth := i.depth
	callDepth := i.calls
	outerResult := i.result
	i.result = nil
	result := runtime.Cell{Type: retType, Value: runtime.Zero(retType)}
	if retType.Kind != runtime.VoidType {
		for len(i.results) <= callDepth {
			i.results = append(i.results, runtime.Value{})
		}
		// Completed calls reuse the result at this depth, even across function
		// names. Cross-type fallthrough gets a fresh typed zero: the reference's
		// uninitialized cross-type storage is not a stable language value.
		if previous := i.results[callDepth]; retType.Equal(previous.Type()) {
			result.Value = runtime.Clone(previous)
		}
		i.result = &result
	}
	i.depth = 0
	i.calls++
	defer func() { i.env = outer; i.depth = depth; i.result = outerResult; i.calls-- }()
	if err := i.configureConsole(console); err != nil {
		return runtime.Value{}, err
	}
	if err := i.defineConsts(consts); err != nil {
		return runtime.Value{}, err
	}
	for _, decl := range locals {
		if err := i.defineVars(decl); err != nil {
			return runtime.Value{}, err
		}
	}
	ctrl, err := i.execStmts(body)
	if err != nil {
		return runtime.Value{}, err
	}
	if ctrl.kind == breakControl {
		return runtime.Value{}, failure(ctrl.at, diag.RLoop, fmt.Errorf("interrompa outside loop"))
	}
	if retType.Kind != runtime.VoidType {
		i.results[callDepth] = runtime.Clone(result.Value)
	}
	// The reference copies parameters back in declaration order, including their
	// numeric types. Repeated destinations therefore receive the last value.
	for _, ref := range references {
		*ref.caller = runtime.Cell{Type: ref.parameter.Type.Clone(), Value: runtime.Clone(ref.parameter.Value)}
	}
	return result.Value, nil
}
