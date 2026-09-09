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
		{"short backing", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 1, High: 2}}}, make([]Cell, 1)}, []int64{2}},
		{"long backing", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 1, High: 2}}}, make([]Cell, 3)}, []int64{1}},
		{"nil vector", nil, []int64{1}},
		{"missing ranges", &Vector{Type{Kind: VectorType, Elem: &integer}, nil}, nil},
		{"missing element", &Vector{Type{Kind: VectorType, Ranges: []Range{{Low: 1, High: 1}}}, make([]Cell, 1)}, []int64{1}},
		{"invalid element", &Vector{Type{Kind: VectorType, Elem: &invalid, Ranges: []Range{{Low: 1, High: 1}}}, make([]Cell, 1)}, []int64{1}},
		{"void element", &Vector{Type{Kind: VectorType, Elem: &void, Ranges: []Range{{Low: 1, High: 1}}}, make([]Cell, 1)}, []int64{1}},
		{"wrong kind", &Vector{Type{Kind: IntegerType, Elem: &integer, Ranges: []Range{{Low: 1, High: 1}}}, make([]Cell, 1)}, []int64{1}},
		{"reversed bound", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 2, High: 1}}}, nil}, []int64{1}},
		{"unresolved bound", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 1, High: 2, HighDynamic: true}}}, make([]Cell, 2)}, []int64{1}},
		{"wrapped width", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: math.MinInt64, High: math.MaxInt64}}}, make([]Cell, 1)}, []int64{0}},
		{"wrapped product", &Vector{Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 0, High: 1 << 32}, {Low: 0, High: 1 << 32}}}, make([]Cell, 1)}, []int64{0, 0}},
		{"missing index", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 1, High: 2}}}), nil},
		{"extra index", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 1, High: 2}}}), []int64{1, 1}},
		{"below bound", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 0, High: 2}}}), []int64{-1}},
		{"above bound", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 0, High: 2}}}), []int64{3}},
		{"second below bound", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 0, High: 1}, {Low: 4, High: 5}}}), []int64{0, 3}},
		{"second above bound", NewVector(Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 0, High: 1}, {Low: 4, High: 5}}}), []int64{0, 6}},
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
		{{Low: 0, High: 2}}, {{Low: 2, High: 4}}, {{Low: math.MaxInt64 - 1, High: math.MaxInt64}},
		{{Low: math.MinInt64, High: math.MinInt64 + 1}}, {{Low: 2, High: 4}, {Low: 5, High: 7}},
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

func TestVectorSlotLimit(t *testing.T) {
	integer := Type{Kind: IntegerType}
	for _, tt := range []struct {
		rows, cols int64
		allowed    bool
	}{
		{1, 1 << 20, true},
		{1, 1<<20 + 1, false},
		{2048, 512, true},
		{2048, 513, false},
	} {
		typ := Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 1, High: tt.rows}, {Low: 1, High: tt.cols}}}
		n, err := typ.Slots()
		if (err == nil) != tt.allowed || tt.allowed && int64(n) != tt.rows*tt.cols {
			t.Fatalf("Slots(%d,%d) = %d, %v; allowed=%t", tt.rows, tt.cols, n, err, tt.allowed)
		}
	}
}
