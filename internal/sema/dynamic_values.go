package sema

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

// A field key covers every instance of its declared layout. Array keys cover
// every element, so branches and indices cannot hide a possible kind change.
type dynamicCell struct {
	id    token.Pos
	field string
}

func (c *checker) cellKey(expr ast.Expr) (dynamicCell, bool) {
	switch e := expr.(type) {
	case *ast.IdentExpr:
		binding, ok := c.info.Binding(e.Name)
		return dynamicCell{id: binding.ID}, ok
	case *ast.IndexExpr:
		key, ok := c.cellKey(e.X)
		key.field += "[]"
		return key, ok
	case *ast.FieldExpr:
		base := c.info.types[e.X]
		return dynamicCell{id: base.RecordID, field: canon(e.Name.Text)}, base.Kind == runtime.RecordType
	}
	return dynamicCell{}, false
}

func (c *checker) valueType(expr ast.Expr, declared runtime.Type) runtime.Type {
	if declared.Kind == runtime.BoolType {
		if key, ok := c.cellKey(expr); ok && c.dynamicCells[key] {
			return runtime.Type{Kind: runtime.DynamicType}
		}
	}
	return declared
}

// comparisonResult identifies expressions that can retain a comparison's
// temporary category. Arithmetic consumes that logical category.
func comparisonResult(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.IsComparison() {
			return true
		}
		switch e.Op.Kind {
		case token.E:
			return comparisonResult(e.Right)
		}
	}
	return false
}

func (c *checker) propagateDynamicCells() {
	queue := make([]dynamicCell, 0, len(c.dynamicCells))
	for key := range c.dynamicCells {
		queue = append(queue, key)
	}
	for head := 0; head < len(queue); head++ {
		for _, target := range c.copyBack[queue[head]] {
			if !c.dynamicCells[target] {
				c.dynamicCells[target] = true
				queue = append(queue, target)
			}
		}
	}
}
