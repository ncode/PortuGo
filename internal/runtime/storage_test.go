package runtime

import (
	"errors"
	"math"
	"testing"
)

func TestStorageUnits(t *testing.T) {
	integer := Type{Kind: IntegerType}
	record := Type{Kind: RecordType, Fields: []Field{{Name: "x", Type: integer}}}
	vector := Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 2, High: 4}}}
	for _, tt := range []struct {
		name string
		typ  Type
		want int
	}{
		{"scalar", integer, 1},
		{"record", record, 3},
		{"empty record", Type{Kind: RecordType}, 2},
		{"vector", vector, 5},
		{"nested", Type{Kind: VectorType, Elem: &vector, Ranges: []Range{{Low: 1, High: 2}}}, 12},
		{"record vector", Type{Kind: VectorType, Elem: &record, Ranges: []Range{{Low: 1, High: 2}}}, 8},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, extra := range []int{-1, 0, 1} {
				got, err := tt.typ.StorageUnits(tt.want + extra)
				if extra < 0 {
					if !errors.Is(err, ErrStorageSize) {
						t.Fatalf("limit accepted: %d, %v", got, err)
					}
				} else if err != nil || got != tt.want {
					t.Fatalf("StorageUnits = %d, %v; want %d", got, err, tt.want)
				}
			}
		})
	}
	oversize := Type{Kind: VectorType, Elem: &integer, Ranges: []Range{{Low: 0, High: math.MaxInt64}}}
	if _, err := oversize.StorageUnits(math.MaxInt); err == nil {
		t.Fatal("overflowing layout accepted")
	}
	for range 255 {
		previous := vector
		vector = Type{Kind: VectorType, Elem: &previous, Ranges: []Range{{Low: 1, High: 1}}}
	}
	if _, err := vector.StorageUnits(100); !errors.Is(err, ErrStorageSize) {
		t.Fatalf("nested containers unbounded: %v", err)
	}
}
