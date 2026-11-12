package interp

import (
	"fmt"

	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/runtime"
	"github.com/ncode/PortuGo/internal/token"
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

func (e *env) lookup(key token.Pos) (*runtime.Cell, bool) {
	for cur := e; cur != nil; cur = cur.parent {
		if cell, ok := cur.cells[key]; ok {
			return cell, true
		}
	}
	return nil, false
}

func (i *Interpreter) assign(cell *runtime.Cell, v runtime.Value) error {
	if v.Comparison && cell.Type.Kind == runtime.BoolType && v.Kind == runtime.RecordValue {
		if err := i.reserveStorage(1); err != nil {
			return err
		}
		defer func() { i.storageUsed-- }()
	}
	converted, err := runtime.ConvertForAssign(cell.Type, v)
	if err != nil {
		return err
	}
	units := valueStorage(converted)
	if err := i.reserveStorage(units); err != nil {
		return err
	}
	old := cell.Value
	if v.Comparison && cell.Type.Kind == runtime.BoolType {
		cell.Type = converted.Type().Clone()
		cell.Value = runtime.Clone(converted)
	} else if cell.Type.Kind == runtime.RecordType {
		var replaced []runtime.Value
		if cell.Value.Rec != nil {
			for _, field := range cell.Value.Rec.Fields {
				if field.Value.Kind == runtime.RecordValue || field.Value.Kind == runtime.VectorValue {
					replaced = append(replaced, field.Value)
				}
			}
		}
		if err := cell.Value.Rec.AssignFrom(converted.Rec); err != nil {
			i.storageUsed -= units
			return err
		}
		// The cloned container/backing cells are temporary; copied field values remain.
		i.storageUsed -= 1 + len(converted.Rec.Fields)
		if len(replaced) != 0 {
			roots := i.retentionRoots()
			for _, value := range replaced {
				i.releaseValue(value, roots)
			}
		}
		return nil
	} else {
		cell.Value = runtime.Clone(converted)
	}
	if old.Kind == runtime.RecordValue || old.Kind == runtime.VectorValue {
		i.releaseValue(old, i.retentionRoots())
	}
	return nil
}

func (i *Interpreter) lookupCell(name token.Token) (*runtime.Cell, error) {
	b, ok := i.info.Reference(name)
	if !ok {
		return nil, failure(name.Pos, diag.RType, fmt.Errorf("missing variable binding"))
	}
	cell, ok := i.env.lookup(b.ID)
	if !ok {
		return nil, failure(name.Pos, diag.RStorage, fmt.Errorf("missing storage for %q", name.Text))
	}
	return cell, nil
}
