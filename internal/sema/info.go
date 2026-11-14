package sema

import (
	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/runtime"
	"github.com/ncode/PortuGo/internal/token"
)

// Binding identifies a resolved declaration and its storage layout.
// ID is the declaration's source position, shared by all of its uses.
// Slots is zero when Type.DynamicBounds requires per-initialization resolution.
type Binding struct {
	ID      token.Pos
	Name    string
	Type    runtime.Type
	Slots   int
	Builtin bool
}

// Reference identifies a binding without exposing its mutable type layout.
type Reference struct {
	ID      token.Pos
	Name    string
	Builtin bool
}

// Info owns semantic facts for one unchanged AST. Reanalyze after editing syntax.
// Query results are independent copies; Info can be shared by interpreters.
type Info struct {
	program     *ast.Program
	valid       bool
	types       map[ast.Expr]runtime.Type
	bindings    map[token.Pos]Binding
	names       map[token.Pos]string
	deferred    map[token.Pos]diag.Diagnostic
	assignments map[*ast.AssignStmt]diag.Diagnostic
}

// ValidFor reports whether this result permits execution of the unchanged program.
func (i *Info) ValidFor(program *ast.Program) bool {
	return i != nil && program != nil && i.valid && i.program == program
}

// TypeOf returns the resolved type of an expression or designator.
func (i *Info) TypeOf(expr ast.Expr) (runtime.Type, bool) {
	if i == nil {
		return runtime.Type{}, false
	}
	t, ok := i.types[expr]
	return t.Clone(), ok
}

// KindOf returns only an expression's resolved type category.
func (i *Info) KindOf(expr ast.Expr) (runtime.TypeKind, bool) {
	if i == nil {
		return runtime.InvalidType, false
	}
	t, ok := i.types[expr]
	return t.Kind, ok
}

// Compatible checks the evaluator's bidirectional type compatibility without
// returning or copying the stored layout.
func (i *Info) Compatible(expr ast.Expr, actual runtime.Type) bool {
	if i == nil {
		return false
	}
	t, ok := i.types[expr]
	return ok && (runtime.Assignable(t, actual) || runtime.Assignable(actual, t))
}

// Binding returns the declaration resolved at a declaration or use token.
func (i *Info) Binding(name token.Token) (Binding, bool) {
	b, ok := i.binding(name)
	if !ok {
		return Binding{}, false
	}
	b.Type = b.Type.Clone()
	return b, true
}

// Reference returns immutable identity facts for a declaration or use token.
func (i *Info) Reference(name token.Token) (Reference, bool) {
	b, ok := i.binding(name)
	return Reference{ID: b.ID, Name: b.Name, Builtin: b.Builtin}, ok
}

func (i *Info) binding(name token.Token) (Binding, bool) {
	if i == nil {
		return Binding{}, false
	}
	b, ok := i.bindings[name.Pos]
	if !ok || i.names[name.Pos] != canon(name.Text) {
		return Binding{}, false
	}
	return b, true
}

func (c *checker) recordBinding(name token.Token, sym symbol) {
	slots, _ := sym.typ.Slots()
	c.info.bindings[name.Pos] = Binding{ID: sym.pos, Name: sym.name, Type: sym.typ, Slots: slots, Builtin: sym.kind == builtinSym}
	c.info.names[name.Pos] = canon(name.Text)
}

func (c *checker) lookup(name token.Token) (symbol, bool) {
	key := canon(name.Text)
	sym, ok := c.scope.lookup(key)
	if !ok {
		sym, ok = c.subs[key]
	}
	if ok {
		c.recordBinding(name, sym)
	}
	return sym, ok
}

// Callables take priority in the corresponding statement or value context.
func (c *checker) lookupCallable(name token.Token, kind symbolKind) (symbol, bool) {
	if sym, ok := c.subs[canon(name.Text)]; ok && sym.kind == kind {
		c.recordBinding(name, sym)
		return sym, true
	}
	return c.lookup(name)
}

// Procedure names win at the start of a statement even when a variable with
// that name exists. Other designator contexts still resolve variables normally.
func (c *checker) assignmentProcedure(expr ast.Expr) (symbol, bool) {
	switch e := expr.(type) {
	case *ast.IdentExpr:
		sym, ok := c.subs[canon(e.Name.Text)]
		return sym, ok && sym.kind == procSym
	case *ast.IndexExpr:
		return c.assignmentProcedure(e.X)
	case *ast.FieldExpr:
		return c.assignmentProcedure(e.X)
	default:
		return symbol{}, false
	}
}

// DeferredDiagnostic returns a name error retained for expression execution.
func (i *Info) DeferredDiagnostic(pos token.Pos) (diag.Diagnostic, bool) {
	if i == nil {
		return diag.Diagnostic{}, false
	}
	d, ok := i.deferred[pos]
	return d, ok
}

// AssignmentDiagnostic returns a checked scalar failure retained until assignment.
func (i *Info) AssignmentDiagnostic(stmt *ast.AssignStmt) (diag.Diagnostic, bool) {
	if i == nil {
		return diag.Diagnostic{}, false
	}
	d, ok := i.assignments[stmt]
	return d, ok
}
