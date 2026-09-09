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
	if v.Type.Kind != VectorType || v.Type.Elem == nil || len(v.Type.Ranges) == 0 {
		return nil, fmt.Errorf("invalid vector layout")
	}
	if kind := v.Type.Elem.Kind; kind < IntegerType || kind > VectorType {
		return nil, fmt.Errorf("invalid vector element type")
	}
	size := uint64(1)
	for _, r := range v.Type.Ranges {
		width := uint64(r.High) - uint64(r.Low) + 1
		if r.High < r.Low || width == 0 || width > uint64(len(v.Elements))/size {
			return nil, fmt.Errorf("vector layout exceeds backing storage")
		}
		size *= width
	}
	if size != uint64(len(v.Elements)) {
		return nil, fmt.Errorf("vector layout does not match backing storage")
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
