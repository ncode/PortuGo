package runtime

import "errors"

// MaxStorageUnits bounds live cells and aggregate containers in one interpreter.
const MaxStorageUnits = 2 * MaxVectorSlots

// ErrStorageSize reports a rejected live-storage reservation before allocation.
var ErrStorageSize = errors.New("live storage limit exceeded")

// StorageUnits counts a cell and its complete zero-value backing storage without
// expanding elements. The existing scalar-slot limit remains independent.
func (t Type) StorageUnits(limit int) (int, error) {
	if _, err := t.Slots(); err != nil {
		return 0, err
	}
	return t.storageUnits(limit)
}

func (t Type) storageUnits(limit int) (int, error) {
	if limit < 1 {
		return 0, ErrStorageSize
	}
	switch t.Kind {
	case RecordType:
		if len(t.Fields) > limit-2 {
			return 0, ErrStorageSize
		}
		return 2 + len(t.Fields), nil
	case VectorType:
		if limit < 2 {
			return 0, ErrStorageSize
		}
		elem, err := t.Elem.storageUnits(limit - 2)
		if err != nil {
			return 0, err
		}
		n := elem
		for _, r := range t.Ranges {
			width := uint64(r.High) - uint64(r.Low) + 1
			if width > uint64((limit-2)/n) {
				return 0, ErrStorageSize
			}
			n *= int(width)
		}
		return 2 + n, nil
	default:
		return 1, nil
	}
}
