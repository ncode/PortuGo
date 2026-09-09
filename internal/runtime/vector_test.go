package runtime

import (
	"math"
	"reflect"
	"slices"
	"testing"
)

func TestVectorCellRejectsInvalidStorage(t *testing.T) {
	integer := Type{Kind: IntegerType}
	invalid := Type{Kind: InvalidType}
	void := Type{Kind: VoidType}
	for _, tt := range []struct {
		name    string
		vector  *Vector
		indices []int64
	}{
		{"short backing", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{1, 2}}}, make([]Cell, 1)}, []int64{2}},
		{"long backing", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{1, 2}}}, make([]Cell, 3)}, []int64{1}},
		{"nil vector", nil, []int64{1}},
		{"missing ranges", &Vector{Type{Kind: VectorType, Elem: &integer}, nil}, nil},
		{"missing element", &Vector{Type{Kind: VectorType, Ranges: []Range{{1, 1}}}, make([]Cell, 1)}, []int64{1}},
		{"invalid element", &Vector{Type{Kind: VectorType, Elem: &invalid, Ranges: []Range{{1, 1}}}, make([]Cell, 1)}, []int64{1}},
		{"void element", &Vector{Type{Kind: VectorType, Elem: &void, Ranges: []Range{{1, 1}}}, make([]Cell, 1)}, []int64{1}},
		{"wrong kind", &Vector{Type{Kind: IntegerType, Elem: &integer, Ranges: []Range{{1, 1}}}, make([]Cell, 1)}, []int64{1}},
		{"reversed bound", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{2, 1}}}, nil}, []int64{1}},
		{"wrapped width", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{math.MinInt64, math.MaxInt64}}}, make([]Cell, 1)}, []int64{0}},
		{"wrapped product", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{0, 1 << 32}, {0, 1 << 32}}}, make([]Cell, 1)}, []int64{0, 0}},
		{"missing index", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{1, 2}}}), nil},
		{"extra index", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{1, 2}}}), []int64{1, 1}},
		{"below bound", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{0, 2}}}), []int64{-1}},
		{"above bound", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{0, 2}}}), []int64{3}},
		{"second below bound", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{0, 1}, {4, 5}}}), []int64{0, 3}},
		{"second above bound", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{0, 1}, {4, 5}}}), []int64{0, 6}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var before []Cell
			if tt.vector != nil {
				before = slices.Clone(tt.vector.Elements)
			}
			cell, err := tt.vector.Cell(tt.indices)
			if err == nil || cell != nil {
				t.Fatalf("Cell(%v) = %v, %v; want nil and an error", tt.indices, cell, err)
			}
			if tt.vector != nil && !reflect.DeepEqual(before, tt.vector.Elements) {
				t.Fatal("failed lookup changed backing storage")
			}
		})
	}
}

func TestVectorCellOffsets(t *testing.T) {
	integer := Type{Kind: IntegerType}
	for _, ranges := range [][]Range{
		{{0, 2}}, {{2, 4}}, {{math.MaxInt64 - 1, math.MaxInt64}},
		{{math.MinInt64, math.MinInt64 + 1}}, {{2, 4}, {5, 7}},
	} {
		v := NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: ranges})
		width := int64(1)
		if len(ranges) == 2 {
			width = ranges[1].High - ranges[1].Low + 1
		}
		for offset := range v.Elements {
			indices := []int64{ranges[0].Low + int64(offset)/width}
			if len(ranges) == 2 {
				indices = append(indices, ranges[1].Low+int64(offset)%width)
			}
			cell, err := v.Cell(indices)
			if err != nil || cell != &v.Elements[offset] || cell.Value.Kind != IntegerValue || cell.Value.Int != 0 {
				t.Fatalf("Cell(%v) = %v, %v; want zero element %d", indices, cell, err, offset)
			}
			if int64(offset)%width == 0 {
				omitted, err := v.Cell(indices[:1])
				if err != nil || omitted != cell {
					t.Fatalf("omitted column = %v, %v; want %v", omitted, err, cell)
				}
			}
		}
	}
}
