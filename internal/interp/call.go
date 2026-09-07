package interp

import (
	"errors"
	"fmt"

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
		args, err := i.evalArgs(call.Args)
		if err != nil {
			return runtime.Value{}, err
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
	_, err = i.callSub(proc.Params, proc.Locals, proc.Body, runtime.Type{Kind: runtime.VoidType}, call.Args)
	return err
}

func (i *Interpreter) callUserFunction(fn *ast.FunctionDecl, args []ast.Expr) (runtime.Value, error) {
	b, ok := i.info.Binding(fn.Name)
	if !ok {
		return runtime.Value{}, failure(fn.Start(), diag.RType, fmt.Errorf("missing return type"))
	}
	return i.callSub(fn.Params, fn.Locals, fn.Body, b.Type, args)
}

func (i *Interpreter) callSub(params []ast.Param, locals []ast.VarDecl, body []ast.Stmt, retType runtime.Type, args []ast.Expr) (runtime.Value, error) {
	if i.calls == maxCalls {
		return runtime.Value{}, fmt.Errorf("active call limit exceeded")
	}
	if len(args) != len(params) {
		return runtime.Value{}, fmt.Errorf("expected %d arguments, got %d", len(params), len(args))
	}
	outer := i.env
	callEnv := newEnv(i.global)
	for idx, param := range params {
		b, ok := i.info.Binding(param.Name)
		if !ok {
			return runtime.Value{}, failure(param.Name.Pos, diag.RType, fmt.Errorf("missing parameter layout"))
		}
		typ := b.Type
		if param.ByRef {
			cell, err := i.lvalue(args[idx])
			if err != nil {
				return runtime.Value{}, err
			}
			callEnv.bind(b.ID, cell)
			continue
		}
		v, err := i.eval(args[idx])
		if err != nil {
			return runtime.Value{}, err
		}
		cell := callEnv.define(b.ID, typ)
		if err := assign(cell, v); err != nil {
			return runtime.Value{}, err
		}
	}
	i.env = callEnv
	depth := i.depth
	i.depth = 0
	i.calls++
	defer func() { i.env = outer; i.depth = depth; i.calls-- }()
	for _, decl := range locals {
		if err := i.defineVars(decl); err != nil {
			return runtime.Value{}, err
		}
	}
	ctrl, err := i.execStmts(body)
	if err != nil {
		return runtime.Value{}, err
	}
	if retType.Kind == runtime.VoidType {
		if ctrl.kind == returnControl {
			return runtime.Value{}, fmt.Errorf("procedure returned a value")
		}
		return runtime.Value{Kind: runtime.VoidValue}, nil
	}
	if ctrl.kind != returnControl {
		return runtime.Value{}, fmt.Errorf("function did not return")
	}
	return runtime.ConvertForAssign(retType, ctrl.value)
}

func (i *Interpreter) evalArgs(args []ast.Expr) ([]runtime.Value, error) {
	values := make([]runtime.Value, len(args))
	for idx, arg := range args {
		v, err := i.eval(arg)
		if err != nil {
			return nil, err
		}
		values[idx] = v
	}
	return values, nil
}
