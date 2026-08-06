package interp

import (
	"fmt"
	"strings"

	"github.com/ncode/portugol-go/internal/runtime"
)

type env struct {
	parent *env
	cells  map[string]*runtime.Cell
}

func newEnv(parent *env) *env {
	return &env{parent: parent, cells: make(map[string]*runtime.Cell)}
}

func (e *env) define(name string, typ runtime.Type) *runtime.Cell {
	cell := &runtime.Cell{Type: typ, Value: runtime.Zero(typ)}
	e.cells[canon(name)] = cell
	return cell
}

func (e *env) bind(name string, cell *runtime.Cell) {
	e.cells[canon(name)] = cell
}

func (e *env) lookup(name string) (*runtime.Cell, bool) {
	key := canon(name)
	for cur := e; cur != nil; cur = cur.parent {
		if cell, ok := cur.cells[key]; ok {
			return cell, true
		}
	}
	return nil, false
}

func canon(name string) string {
	return strings.ToLower(name)
}

func assign(cell *runtime.Cell, v runtime.Value) error {
	converted, err := runtime.ConvertForAssign(cell.Type, v)
	if err != nil {
		return err
	}
	cell.Value = runtime.Clone(converted)
	return nil
}

func lookupCell(e *env, name string) (*runtime.Cell, error) {
	cell, ok := e.lookup(name)
	if !ok {
		return nil, fmt.Errorf("undeclared identifier %q", name)
	}
	return cell, nil
}
