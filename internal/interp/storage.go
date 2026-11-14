package interp

import (
	"github.com/ncode/PortuGo/internal/runtime"
	"github.com/ncode/PortuGo/internal/token"
)

func (i *Interpreter) reserveStorage(units int) error {
	if units < 0 || units > i.storageLimit-i.storageUsed {
		return runtime.ErrStorageSize
	}
	i.storageUsed += units
	return nil
}

func (i *Interpreter) define(e *env, id token.Pos, typ runtime.Type) (*runtime.Cell, error) {
	units, err := typ.StorageUnits(i.storageLimit - i.storageUsed)
	if err != nil {
		return nil, err
	}
	if err := i.reserveStorage(units); err != nil {
		return nil, err
	}
	return e.define(id, typ), nil
}

func (i *Interpreter) releaseEnv(e *env) {
	var roots map[any]bool
	for _, cell := range e.cells {
		i.storageUsed--
		if cell.Value.Kind == runtime.VectorValue || cell.Value.Kind == runtime.RecordValue {
			if roots == nil {
				roots = i.retentionRoots()
			}
			i.releaseValue(cell.Value, roots)
		}
	}
}

func (i *Interpreter) retentionRoots() map[any]bool {
	var roots map[any]bool
	for _, values := range [2][]runtime.Value{i.operands, i.argumentRoots} {
		for _, value := range values {
			var key any
			switch {
			case value.Kind == runtime.RecordValue && value.Rec != nil:
				key = value.Rec
			case value.Kind == runtime.VectorValue && value.Vec != nil:
				key = value.Vec
			default:
				continue
			}
			if roots == nil {
				roots = make(map[any]bool)
			}
			roots[key] = true
		}
	}
	return roots
}

func (i *Interpreter) releaseValue(v runtime.Value, roots map[any]bool) {
	var key any
	var cells []runtime.Cell
	switch {
	case v.Kind == runtime.RecordValue && v.Rec != nil:
		key, cells = v.Rec, v.Rec.Fields
	case v.Kind == runtime.VectorValue && v.Vec != nil:
		key, cells = v.Vec, v.Vec.Elements
	default:
		return
	}
	if roots[key] {
		i.retainedStorage = append(i.retainedStorage, v)
		return
	}
	i.storageUsed -= 1 + len(cells)
	for _, cell := range cells {
		i.releaseValue(cell.Value, roots)
	}
}

func (i *Interpreter) releaseRetainedStorage() {
	if len(i.retainedStorage) == 0 {
		return
	}
	previous := i.retainedStorage
	i.retainedStorage = nil
	roots := i.retentionRoots()
	for _, v := range previous {
		i.releaseValue(v, roots)
	}
}

// Values own independent aggregates; captured references point to existing cells.
func valueStorage(v runtime.Value) int {
	var cells []runtime.Cell
	switch {
	case v.Kind == runtime.VectorValue && v.Vec != nil:
		cells = v.Vec.Elements
	case v.Kind == runtime.RecordValue && v.Rec != nil:
		cells = v.Rec.Fields
	default:
		return 0
	}
	n := 1 + len(cells)
	for _, cell := range cells {
		n += valueStorage(cell.Value)
	}
	return n
}

func (i *Interpreter) copyCell(destination, source *runtime.Cell) error {
	units := valueStorage(source.Value)
	if err := i.reserveStorage(units); err != nil {
		return err
	}
	old := destination.Value
	*destination = runtime.Cell{Type: source.Type.Clone(), Value: runtime.Clone(source.Value)}
	if old.Kind == runtime.RecordValue || old.Kind == runtime.VectorValue {
		i.releaseValue(old, i.retentionRoots())
	}
	return nil
}
