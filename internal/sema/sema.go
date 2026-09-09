package sema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

// Analyze resolves names, types, layouts, calls, loop control, and returns.
func Analyze(prog *ast.Program) (*Info, []diag.Diagnostic) {
	info := &Info{program: prog, types: make(map[ast.Expr]runtime.Type), bindings: make(map[token.Pos]Binding), names: make(map[token.Pos]string)}
	c := &checker{scope: newScope(nil), subs: make(map[string]symbol), info: info}
	if prog == nil {
		c.error(token.NoPos, diag.ETypeMismatch, "missing program")
		return info, c.diags
	}
	if ds := ast.CheckLimits(prog); len(ds) != 0 {
		return info, ds
	}
	c.declareBuiltins()
	c.checkProgram(prog)
	info.valid = !diag.HasErrors(c.diags)
	return info, diag.Ordered(c.diags)
}

type symbolKind int

const (
	varSym symbolKind = iota
	constSym
	procSym
	funcSym
	builtinSym
)

type paramSig struct {
	name  string
	typ   runtime.Type
	byRef bool
}

type symbol struct {
	pos    token.Pos
	name   string
	kind   symbolKind
	typ    runtime.Type
	params []paramSig
}

type scope struct {
	parent *scope
	syms   map[string]symbol
}

type checker struct {
	info       *Info
	scope      *scope
	subs       map[string]symbol
	diags      []diag.Diagnostic
	loopDepth  int
	returnType runtime.Type
	inFunction bool
}

func newScope(parent *scope) *scope {
	return &scope{parent: parent, syms: make(map[string]symbol)}
}

func (s *scope) lookup(name string) (symbol, bool) {
	for cur := s; cur != nil; cur = cur.parent {
		if sym, ok := cur.syms[name]; ok {
			return sym, true
		}
	}
	return symbol{}, false
}

func (s *scope) declare(sym symbol) bool {
	if _, ok := s.syms[sym.name]; ok {
		return false
	}
	s.syms[sym.name] = sym
	return true
}

func (c *checker) checkProgram(prog *ast.Program) {
	if !c.declareConsts(prog.Consts) {
		return
	}
	for _, decl := range prog.Globals {
		c.declareVars(decl)
	}
	for _, sub := range prog.Subs {
		c.declareSub(sub)
	}
	if diag.HasErrors(c.diags) {
		return
	}
	for _, sub := range prog.Subs {
		c.checkSub(sub)
	}
	c.checkStmts(prog.Body)
}

func (c *checker) declareVars(decl ast.VarDecl) {
	typ := runtime.TypeFromSpec(decl.Type)
	c.validateType(decl.Type, typ)
	for _, name := range decl.Names {
		key := canon(name.Text)
		sym := symbol{pos: name.Pos, name: key, kind: varSym, typ: typ}
		c.recordBinding(name, sym)
		if !c.scope.declare(sym) {
			c.error(name.Pos, diag.ERedeclared, "redeclared identifier %q", name.Text)
		}
	}
}

func (c *checker) declareSub(sub ast.Subprogram) {
	nameTok := sub.NameToken()
	key := canon(nameTok.Text)
	sym := symbol{pos: nameTok.Pos, name: key}
	switch d := sub.(type) {
	case *ast.ProcedureDecl:
		sym.kind = procSym
		sym.typ = runtime.Type{Kind: runtime.VoidType}
		sym.params = paramsFromAST(d.Params)
	case *ast.FunctionDecl:
		sym.kind = funcSym
		sym.typ = runtime.TypeFromSpec(d.Return)
		sym.params = paramsFromAST(d.Params)
		c.validateType(d.Return, sym.typ)
	}
	c.recordBinding(nameTok, sym)
	previous, reserved := c.scope.syms[key]
	if _, exists := c.subs[key]; exists || reserved && previous.kind == builtinSym {
		c.error(nameTok.Pos, diag.ERedeclared, "redeclared identifier %q", nameTok.Text)
		return
	}
	c.subs[key] = sym
}

func paramsFromAST(params []ast.Param) []paramSig {
	out := make([]paramSig, len(params))
	for i, p := range params {
		out[i] = paramSig{name: canon(p.Name.Text), typ: runtime.TypeFromSpec(p.Type), byRef: p.ByRef}
	}
	return out
}

func (c *checker) checkSub(sub ast.Subprogram) {
	outer := c.scope
	c.scope = newScope(outer)
	prevReturn, prevInFunction := c.returnType, c.inFunction
	defer func() {
		c.scope = outer
		c.returnType, c.inFunction = prevReturn, prevInFunction
	}()

	switch d := sub.(type) {
	case *ast.ProcedureDecl:
		c.inFunction = false
		c.declareParams(d.Params)
		if !c.declareConsts(d.Consts) {
			return
		}
		for _, decl := range d.Locals {
			c.declareVars(decl)
		}
		c.checkStmts(d.Body)
	case *ast.FunctionDecl:
		c.inFunction = true
		c.returnType = runtime.TypeFromSpec(d.Return)
		c.declareParams(d.Params)
		if !c.declareConsts(d.Consts) {
			return
		}
		for _, decl := range d.Locals {
			c.declareVars(decl)
		}
		c.checkStmts(d.Body)
	}
}

func (c *checker) declareParams(params []ast.Param) {
	for _, p := range params {
		typ := runtime.TypeFromSpec(p.Type)
		c.validateType(p.Type, typ)
		sym := symbol{pos: p.Name.Pos, name: canon(p.Name.Text), kind: varSym, typ: typ}
		c.recordBinding(p.Name, sym)
		if !c.scope.declare(sym) {
			c.error(p.Name.Pos, diag.ERedeclared, "redeclared parameter %q", p.Name.Text)
		}
	}
}

func (c *checker) validateType(spec ast.TypeSpec, typ runtime.Type) {
	if typ.Kind == runtime.InvalidType {
		c.error(spec.At, diag.ETypeMismatch, "invalid type")
	}
	if typ.Kind == runtime.VectorType {
		for _, r := range spec.Ranges {
			if !c.checkBound(r.Low, r.At) || !c.checkBound(r.High, r.At) {
				return
			}
		}
		if spec.Elem != nil && typ.Elem != nil {
			before := len(c.diags)
			c.validateType(*spec.Elem, *typ.Elem)
			if len(c.diags) != before {
				return
			}
		}
		if typ.DynamicBounds() {
			return
		}
		if _, err := typ.Slots(); err != nil {
			code := diag.ETypeMismatch
			if errors.Is(err, runtime.ErrVectorSize) {
				code = diag.EResource
			}
			c.error(spec.At, code, "%s", err)
			return
		}
	}
}

func (c *checker) checkStmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		c.checkStmt(stmt)
	}
}

func (c *checker) checkStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if sym, ok := c.assignmentProcedure(s.Target); ok {
			c.error(sym.pos, diag.ECall, "%q is a procedure in statement context", sym.name)
			return
		}
		dst, ok := c.writable(s.Target)
		if ok && dst.Kind == runtime.VectorType {
			c.error(s.Target.Start(), diag.ETypeMismatch, "vector assignment requires an indexed target")
			return
		}
		src := c.expr(s.Value)
		if ok && !runtime.Assignable(dst, src) {
			c.error(s.Value.Start(), diag.ETypeMismatch, "cannot assign %s to %s", src, dst)
		}
	case *ast.CallStmt:
		c.checkCall(s.Call, true)
	case *ast.IfStmt:
		c.requireBool(s.Cond)
		c.checkStmts(s.Then)
		c.checkStmts(s.Else)
	case *ast.SwitchStmt:
		x := c.expr(s.X)
		for _, cc := range s.Cases {
			for _, label := range cc.Labels {
				mismatch := false
				for _, expr := range []ast.Expr{label.Low, label.High} {
					if expr == nil {
						continue
					}
					t := c.expr(expr)
					if !mismatch && x.Kind != runtime.VoidType && t.Kind != runtime.VoidType && !runtime.Assignable(x, t) && !runtime.Assignable(t, x) {
						c.error(expr.Start(), diag.ETypeMismatch, "cannot compare %s with %s", x, t)
						mismatch = true
					}
				}
			}
			c.checkStmts(cc.Body)
		}
		c.checkStmts(s.Default)
	case *ast.WhileStmt:
		c.requireBool(s.Cond)
		c.withLoop(func() { c.checkStmts(s.Body) })
	case *ast.RepeatStmt:
		c.withLoop(func() { c.checkStmts(s.Body) })
		c.requireBool(s.Cond)
	case *ast.ForStmt:
		sym, ok := c.lookup(s.Name)
		if !ok || sym.kind != varSym {
			c.error(s.Name.Pos, diag.EUndeclared, "undeclared loop variable %q", s.Name.Text)
		} else if sym.typ.Kind != runtime.IntegerType {
			c.error(s.Name.Pos, diag.ETypeMismatch, "loop variable must be inteiro")
		}
		c.requireInt(s.From)
		c.requireInt(s.To)
		if s.Step != nil {
			c.requireInt(s.Step)
		}
		c.withLoop(func() { c.checkStmts(s.Body) })
	case *ast.BreakStmt:
		if c.loopDepth == 0 {
			c.error(s.At, diag.EBreak, "interrompa outside loop")
		}
	case *ast.ReturnStmt:
		if !c.inFunction {
			c.error(s.At, diag.EReturn, "retorne outside function")
			return
		}
		if s.Value == nil {
			c.error(s.At, diag.ETypeMismatch, "retorne requires a value")
			return
		}
		t := c.expr(s.Value)
		if !runtime.Assignable(c.returnType, t) {
			c.error(s.Value.Start(), diag.ETypeMismatch, "cannot return %s from %s function", t, c.returnType)
		}
	case *ast.ReadStmt:
		for _, target := range s.Targets {
			if t, ok := c.writable(target); ok && t.Kind == runtime.VectorType {
				c.error(target.Start(), diag.ETypeMismatch, "cannot read a whole vector")
			}
		}
	case *ast.WriteStmt:
		for _, arg := range s.Args {
			typ := c.expr(arg.Expr)
			if arg.Width != nil {
				if typ.Kind == runtime.BoolType {
					c.error(arg.Expr.Start(), diag.ETypeMismatch, "cannot format logico with a field width")
				}
				c.requireInt(arg.Width)
			}
			if arg.Decimals != nil {
				c.requireInt(arg.Decimals)
			}
		}
	}
}

func (c *checker) expr(expr ast.Expr) (typ runtime.Type) {
	defer func() { c.info.types[expr] = typ }()
	switch e := expr.(type) {
	case *ast.LiteralExpr:
		switch e.Kind {
		case ast.IntLiteral:
			return runtime.Type{Kind: runtime.IntegerType}
		case ast.RealLiteral:
			return runtime.Type{Kind: runtime.RealType}
		case ast.StringLiteral:
			return runtime.Type{Kind: runtime.StringType}
		case ast.BoolLiteral:
			return runtime.Type{Kind: runtime.BoolType}
		}
	case *ast.IdentExpr:
		sym, ok := c.lookupCallable(e.Name, funcSym)
		if !ok {
			c.error(e.Name.Pos, diag.EUndeclared, "undeclared identifier %q", e.Name.Text)
			return runtime.Type{Kind: runtime.InvalidType}
		}
		if sym.kind == funcSym {
			c.checkArgs(&ast.CallExpr{Name: e.Name}, sym)
			return sym.typ
		}
		if sym.kind == builtinSym && sym.name == "pi" {
			return c.checkCall(&ast.CallExpr{Name: e.Name}, false)
		}
		if sym.kind == builtinSym {
			c.error(e.Name.Pos, diag.EParse, "expected '(' after %s", sym.name)
			return runtime.Type{Kind: runtime.InvalidType}
		}
		if sym.kind != varSym && sym.kind != constSym {
			c.error(e.Name.Pos, diag.ETypeMismatch, "%q is not a variable", e.Name.Text)
			return runtime.Type{Kind: runtime.InvalidType}
		}
		return sym.typ
	case *ast.IndexExpr:
		base := c.expr(e.X)
		for _, idx := range e.Indices {
			c.requireInt(idx)
		}
		if base.Kind != runtime.VectorType || base.Elem == nil {
			c.error(e.Start(), diag.ETypeMismatch, "cannot index %s", base)
			return runtime.Type{Kind: runtime.InvalidType}
		}
		omittedColumn := len(e.Indices) == 1 && len(base.Ranges) == 2
		if len(e.Indices) != len(base.Ranges) && !omittedColumn {
			c.error(e.Start(), diag.ETypeMismatch, "expected %d indices, got %d", len(base.Ranges), len(e.Indices))
		}
		return *base.Elem
	case *ast.UnaryExpr:
		return c.unary(e)
	case *ast.BinaryExpr:
		return c.binary(e)
	case *ast.CallExpr:
		return c.checkCall(e, false)
	}
	return runtime.Type{Kind: runtime.InvalidType}
}

func (c *checker) unary(e *ast.UnaryExpr) runtime.Type {
	t := c.expr(e.X)
	switch e.Op.Kind {
	case token.ADD, token.SUB:
		if isNumeric(t) {
			return t
		}
		if e.Op.Kind == token.SUB && isScalar(t) {
			return runtime.Type{Kind: runtime.VoidType}
		}
		c.error(e.Op.Pos, diag.ETypeMismatch, "operator %s requires numeric operand", e.Op.Text)
	case token.NAO:
		if t.Kind == runtime.BoolType {
			return t
		}
		c.error(e.Op.Pos, diag.ETypeMismatch, "operator nao requires logico operand")
	}
	return runtime.Type{Kind: runtime.InvalidType}
}

func (c *checker) binary(e *ast.BinaryExpr) runtime.Type {
	left := c.expr(e.Left)
	right := c.expr(e.Right)
	switch e.Op.Kind {
	case token.ADD:
		if left.Kind == runtime.StringType && right.Kind == runtime.StringType {
			return runtime.Type{Kind: runtime.StringType}
		}
		return c.numericBinary(e.Op.Pos, left, right)
	case token.SUB, token.MUL:
		return c.numericBinary(e.Op.Pos, left, right)
	case token.QUO, token.POW:
		if isNumeric(left) && isNumeric(right) {
			return runtime.Type{Kind: runtime.RealType}
		}
		if isScalar(left) && isScalar(right) {
			if e.Op.Kind == token.QUO {
				return right
			}
			return runtime.Type{Kind: runtime.VoidType}
		}
		c.error(e.Op.Pos, diag.ETypeMismatch, "operator %s requires numeric operands", e.Op.Text)
	case token.IDIV:
		if isScalar(left) && isScalar(right) {
			return right
		}
		c.error(e.Op.Pos, diag.ETypeMismatch, "operator %s requires scalar operands", e.Op.Text)
	case token.REM, token.MOD:
		if isScalar(left) && isScalar(right) && left.Kind != runtime.StringType && right.Kind != runtime.StringType {
			if isNumeric(left) && isNumeric(right) {
				return runtime.Type{Kind: runtime.IntegerType}
			}
			return right
		}
		c.error(e.Op.Pos, diag.ETypeMismatch, "operator %s requires numeric or logico operands", e.Op.Text)
	case token.EQL, token.NEQ:
		if runtime.Assignable(left, right) || runtime.Assignable(right, left) {
			return runtime.Type{Kind: runtime.BoolType}
		}
		c.error(e.Op.Pos, diag.ETypeMismatch, "cannot compare %s with %s", left, right)
	case token.LSS, token.GTR, token.LEQ, token.GEQ:
		if (isNumeric(left) && isNumeric(right)) || (left.Kind == runtime.StringType && right.Kind == runtime.StringType) {
			return runtime.Type{Kind: runtime.BoolType}
		}
		c.error(e.Op.Pos, diag.ETypeMismatch, "operator %s requires comparable operands", e.Op.Text)
	case token.E, token.OU, token.XOU:
		if left.Kind == runtime.BoolType && right.Kind == runtime.BoolType {
			return runtime.Type{Kind: runtime.BoolType}
		}
		c.error(e.Op.Pos, diag.ETypeMismatch, "logical operator requires logico operands")
	}
	return runtime.Type{Kind: runtime.InvalidType}
}

func (c *checker) numericBinary(pos token.Pos, left, right runtime.Type) runtime.Type {
	if isNumeric(left) && isNumeric(right) {
		if left.Kind == runtime.RealType || right.Kind == runtime.RealType {
			return runtime.Type{Kind: runtime.RealType}
		}
		return runtime.Type{Kind: runtime.IntegerType}
	}
	c.error(pos, diag.ETypeMismatch, "operator requires numeric operands")
	return runtime.Type{Kind: runtime.InvalidType}
}

func (c *checker) checkCall(call *ast.CallExpr, asStmt bool) (typ runtime.Type) {
	defer func() { c.info.types[call] = typ }()
	name := canon(call.Name.Text)
	if ret, ok := c.builtinCallType(name, call); ok {
		c.recordBinding(call.Name, symbol{name: name, kind: builtinSym, typ: ret})
		if asStmt && ret.Kind != runtime.VoidType {
			return ret
		}
		return ret
	}
	kind := funcSym
	if asStmt {
		kind = procSym
	}
	sym, ok := c.lookupCallable(call.Name, kind)
	if !ok {
		c.error(call.Name.Pos, diag.EUndeclared, "undeclared callable %q", call.Name.Text)
		return runtime.Type{Kind: runtime.InvalidType}
	}
	if asStmt && sym.kind != procSym {
		c.error(call.Name.Pos, diag.ECall, "%q is not a procedure", call.Name.Text)
		return runtime.Type{Kind: runtime.InvalidType}
	}
	if !asStmt && sym.kind != funcSym {
		c.error(call.Name.Pos, diag.ECall, "%q is not a function", call.Name.Text)
		return runtime.Type{Kind: runtime.InvalidType}
	}
	c.checkArgs(call, sym)
	return sym.typ
}

func (c *checker) checkArgs(call *ast.CallExpr, sym symbol) {
	params := sym.params
	pos := call.Name.Pos
	if sym.kind == procSym {
		pos = sym.pos
	}
	if len(call.Args) != len(params) {
		c.error(pos, diag.ECall, "%q expects %d arguments, got %d", call.Name.Text, len(params), len(call.Args))
		return
	}
	for i, arg := range call.Args {
		t := c.expr(arg)
		if !runtime.Assignable(params[i].typ, t) && (!isNumeric(params[i].typ) || !isNumeric(t)) {
			argPos := arg.Start()
			if sym.kind == procSym {
				argPos = sym.pos
			}
			c.error(argPos, diag.ETypeMismatch, "argument %d: cannot use %s as %s", i+1, t, params[i].typ)
		}
		if params[i].byRef && !c.isWritableExpr(arg) {
			c.error(arg.Start(), diag.ECall, "argument %d must be assignable for var parameter", i+1)
		}
	}
}

func (c *checker) writable(expr ast.Expr) (runtime.Type, bool) {
	switch e := expr.(type) {
	case *ast.IdentExpr:
		sym, ok := c.lookup(e.Name)
		if !ok || sym.kind == funcSym || sym.kind == constSym {
			c.error(e.Name.Pos, diag.EUndeclared, "undeclared identifier %q", e.Name.Text)
			return runtime.Type{Kind: runtime.InvalidType}, false
		}
		if sym.kind != varSym {
			c.error(e.Name.Pos, diag.ETypeMismatch, "%q is not assignable", e.Name.Text)
			return runtime.Type{Kind: runtime.InvalidType}, false
		}
		c.info.types[expr] = sym.typ
		return sym.typ, true
	case *ast.IndexExpr:
		return c.expr(e), true
	default:
		c.error(expr.Start(), diag.ETypeMismatch, "expression is not assignable")
		return runtime.Type{Kind: runtime.InvalidType}, false
	}
}

func (c *checker) requireBool(expr ast.Expr) {
	if t := c.expr(expr); t.Kind != runtime.BoolType && t.Kind != runtime.InvalidType {
		c.error(expr.Start(), diag.ETypeMismatch, "expected logico, got %s", t)
	}
}

func (c *checker) requireInt(expr ast.Expr) {
	if t := c.expr(expr); t.Kind != runtime.IntegerType && t.Kind != runtime.InvalidType {
		c.error(expr.Start(), diag.ETypeMismatch, "expected inteiro, got %s", t)
	}
}

func (c *checker) withLoop(fn func()) {
	c.loopDepth++
	defer func() { c.loopDepth-- }()
	fn()
}

func (c *checker) error(pos token.Pos, code diag.Code, format string, args ...any) {
	c.diags = append(c.diags, diag.Diagnostic{Code: code, Pos: pos, Message: fmt.Sprintf(format, args...)})
}

func canon(name string) string {
	return strings.ToLower(name)
}

func isNumeric(t runtime.Type) bool {
	return t.Kind == runtime.IntegerType || t.Kind == runtime.RealType
}

func isScalar(t runtime.Type) bool {
	return isNumeric(t) || t.Kind == runtime.StringType || t.Kind == runtime.BoolType
}

func (c *checker) isWritableExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.IdentExpr:
		sym, ok := c.lookupCallable(e.Name, funcSym)
		return ok && sym.kind == varSym
	case *ast.IndexExpr:
		return c.isWritableExpr(e.X)
	default:
		return false
	}
}
