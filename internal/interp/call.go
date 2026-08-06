package interp

import (
	"fmt"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/stdlib"
)

func (i *Interpreter) callFunction(call *ast.CallExpr) (runtime.Value, error) {
	name := canon(call.Name.Text)
	if stdlib.IsBuiltin(name) {
		args, err := i.evalArgs(call.Args)
		if err != nil {
			return runtime.Value{}, err
		}
		if v, ok, err := i.lib.Call(name, args); ok {
			return v, err
		}
	}
	sub, ok := i.subs[name]
	if !ok {
		return runtime.Value{}, fmt.Errorf("undefined function %q", call.Name.Text)
	}
	fn, ok := sub.(*ast.FunctionDecl)
	if !ok {
		return runtime.Value{}, fmt.Errorf("%q is not a function", call.Name.Text)
	}
	return i.callUserFunction(fn, call.Args)
}

func (i *Interpreter) callProcedure(call *ast.CallExpr) error {
	name := canon(call.Name.Text)
	sub, ok := i.subs[name]
	if !ok {
		return fmt.Errorf("undefined procedure %q", call.Name.Text)
	}
	proc, ok := sub.(*ast.ProcedureDecl)
	if !ok {
		return fmt.Errorf("%q is not a procedure", call.Name.Text)
	}
	_, err := i.callSub(proc.Params, proc.Locals, proc.Body, runtime.Type{Kind: runtime.VoidType}, call.Args)
	return err
}

func (i *Interpreter) callUserFunction(fn *ast.FunctionDecl, args []ast.Expr) (runtime.Value, error) {
	retType := runtime.TypeFromSpec(fn.Return)
	return i.callSub(fn.Params, fn.Locals, fn.Body, retType, args)
}

func (i *Interpreter) callSub(params []ast.Param, locals []ast.VarDecl, body []ast.Stmt, retType runtime.Type, args []ast.Expr) (runtime.Value, error) {
	if len(args) != len(params) {
		return runtime.Value{}, fmt.Errorf("expected %d arguments, got %d", len(params), len(args))
	}
	outer := i.env
	callEnv := newEnv(outer)
	i.env = callEnv
	defer func() { i.env = outer }()
	for idx, param := range params {
		typ := runtime.TypeFromSpec(param.Type)
		if param.ByRef {
			cell, err := i.lvalueIn(outer, args[idx])
			if err != nil {
				return runtime.Value{}, err
			}
			callEnv.bind(param.Name.Text, cell)
			continue
		}
		v, err := i.evalIn(outer, args[idx])
		if err != nil {
			return runtime.Value{}, err
		}
		cell := callEnv.define(param.Name.Text, typ)
		if err := assign(cell, v); err != nil {
			return runtime.Value{}, err
		}
	}
	for _, decl := range locals {
		i.defineVars(decl)
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

func (i *Interpreter) evalIn(e *env, expr ast.Expr) (runtime.Value, error) {
	outer := i.env
	i.env = e
	defer func() { i.env = outer }()
	return i.eval(expr)
}

func (i *Interpreter) lvalueIn(e *env, expr ast.Expr) (*runtime.Cell, error) {
	outer := i.env
	i.env = e
	defer func() { i.env = outer }()
	return i.lvalue(expr)
}
