package runtime

import "fmt"

// Vector stores a flat checked representation of a Portugol vector.
type Vector struct {
	Type     Type
	Elements []Cell
}

// NewVector creates a zero-filled vector for typ.
func NewVector(typ Type) *Vector {
	size := int64(1)
	for _, r := range typ.Ranges {
		if r.High < r.Low {
			size = 0
			break
		}
		size *= r.High - r.Low + 1
	}
	elem := Type{Kind: InvalidType}
	if typ.Elem != nil {
		elem = *typ.Elem
	}
	elements := make([]Cell, int(size))
	for i := range elements {
		elements[i] = Cell{Type: elem, Value: Zero(elem)}
	}
	return &Vector{Type: typ, Elements: elements}
}

// Cell returns the storage location for indices.
func (v *Vector) Cell(indices []int64) (*Cell, error) {
	if v == nil {
		return nil, fmt.Errorf("nil vector")
	}
	if len(indices) == 1 && len(v.Type.Ranges) == 2 {
		indices = []int64{indices[0], v.Type.Ranges[1].Low}
	}
	if len(indices) != len(v.Type.Ranges) {
		return nil, fmt.Errorf("expected %d indices, got %d", len(v.Type.Ranges), len(indices))
	}
	var offset int64
	stride := int64(1)
	for i := len(indices) - 1; i >= 0; i-- {
		r := v.Type.Ranges[i]
		idx := indices[i]
		if idx < r.Low || idx > r.High {
			return nil, fmt.Errorf("index %d out of bounds %d..%d", idx, r.Low, r.High)
		}
		offset += (idx - r.Low) * stride
		stride *= r.High - r.Low + 1
	}
	return &v.Elements[offset], nil
}
