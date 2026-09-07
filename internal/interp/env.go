package interp

import (
	"fmt"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

type env struct {
	parent *env
	cells  map[token.Pos]*runtime.Cell
}

func newEnv(parent *env) *env {
	return &env{parent: parent, cells: make(map[token.Pos]*runtime.Cell)}
}

func (e *env) define(id token.Pos, typ runtime.Type) *runtime.Cell {
	cell := &runtime.Cell{Type: typ, Value: runtime.Zero(typ)}
	e.cells[id] = cell
	return cell
}

func (e *env) bind(id token.Pos, cell *runtime.Cell) {
	e.cells[id] = cell
}

func (e *env) lookup(key token.Pos) (*runtime.Cell, bool) {
	for cur := e; cur != nil; cur = cur.parent {
		if cell, ok := cur.cells[key]; ok {
			return cell, true
		}
	}
	return nil, false
}

func assign(cell *runtime.Cell, v runtime.Value) error {
	converted, err := runtime.ConvertForAssign(cell.Type, v)
	if err != nil {
		return err
	}
	cell.Value = runtime.Clone(converted)
	return nil
}

func (i *Interpreter) lookupCell(name token.Token) (*runtime.Cell, error) {
	b, ok := i.info.Binding(name)
	if !ok {
		return nil, failure(name.Pos, diag.RType, fmt.Errorf("missing variable binding"))
	}
	cell, ok := i.env.lookup(b.ID)
	if !ok {
		return nil, failure(name.Pos, diag.RStorage, fmt.Errorf("missing storage for %q", name.Text))
	}
	return cell, nil
}
