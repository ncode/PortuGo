package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func (c *checker) declareTypes(decls []ast.TypeDecl) bool {
	for _, decl := range decls {
		typ := c.resolveType(decl.Type)
		if typ.Kind == runtime.InvalidType {
			return false
		}
		key := canon(decl.Name.Text)
		if _, exists := c.scope.types[key]; !exists {
			c.scope.types[key] = symbol{pos: decl.Name.Pos, name: key, kind: typeSym, typ: typ}
		}
		c.recordBinding(decl.Name, c.scope.types[key])
	}
	return true
}

func (c *checker) resolveType(spec ast.TypeSpec) runtime.Type {
	typ := runtime.TypeFromSpec(spec)
	if typ.Kind == runtime.VectorType && spec.Elem != nil {
		elem := c.resolveType(*spec.Elem)
		if elem.Kind == runtime.InvalidType {
			return elem
		}
		typ.Elem = &elem
	}
	if typ.Kind != runtime.InvalidType {
		return typ
	}
	for scope := c.scope; scope != nil; scope = scope.parent {
		if sym, ok := scope.types[canon(spec.Name)]; ok {
			c.recordBinding(token.Token{Text: spec.Name, Pos: spec.At}, sym)
			return sym.typ
		}
	}
	c.error(spec.At, diag.EParse, "unknown type %q", spec.Name)
	return typ
}
