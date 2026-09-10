package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func (c *checker) recordType(spec ast.TypeSpec) runtime.Type {
	typ := runtime.Type{Kind: runtime.RecordType}
	for _, decl := range spec.Fields {
		field := c.resolveType(decl.Type)
		if field.Kind == runtime.InvalidType {
			return field
		}
		if field.Kind == runtime.RecordType {
			continue // Nested declarations create no addressable fields.
		}
		if !isScalar(field) {
			c.error(decl.Type.At, diag.EParse, "unsupported record field type")
			return runtime.Type{Kind: runtime.InvalidType}
		}
		for _, name := range decl.Names {
			typ.Fields = append(typ.Fields, runtime.Field{Name: canon(name.Text), Type: field})
		}
	}
	return typ
}

func (c *checker) fieldType(expr *ast.FieldExpr) runtime.Type {
	if c.stopped {
		return runtime.Type{Kind: runtime.InvalidType}
	}
	base := c.expr(expr.X)
	if base.Kind == runtime.InvalidType {
		return base
	}
	if base.Kind == runtime.DynamicType {
		return base
	}
	if base.Kind == runtime.RecordType {
		for _, field := range base.Fields {
			if field.Name == canon(expr.Name.Text) {
				return c.valueType(expr, field.Type)
			}
		}
	}
	c.error(expr.Name.Pos, diag.EUndeclared, "unknown field %q", expr.Name.Text)
	c.stopped = true
	return runtime.Type{Kind: runtime.InvalidType}
}
