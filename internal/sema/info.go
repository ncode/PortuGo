package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
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

// Info owns semantic facts for one unchanged AST. Reanalyze after editing syntax.
// Query results are independent copies; Info can be shared by interpreters.
type Info struct {
	program  *ast.Program
	valid    bool
	types    map[ast.Expr]runtime.Type
	bindings map[token.Pos]Binding
	names    map[token.Pos]string
}

// ValidFor reports whether this result successfully analyzed program.
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

// Binding returns the declaration resolved at a declaration or use token.
func (i *Info) Binding(name token.Token) (Binding, bool) {
	if i == nil {
		return Binding{}, false
	}
	b, ok := i.bindings[name.Pos]
	if !ok || i.names[name.Pos] != canon(name.Text) {
		return Binding{}, false
	}
	b.Type = b.Type.Clone()
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
	default:
		return symbol{}, false
	}
}
