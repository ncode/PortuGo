package interp

import (
	"testing"

	"github.com/ncode/portugol-go/internal/runtime"
)

func TestUnicodeComparisonExtension(t *testing.T) {
	for _, tt := range []struct {
		left, right string
		order       int
	}{
		{"€", "λ", -1},
		{"λ", "λ", 0},
		{"\u0081", "\u0181", -1},
		{"\U0010ffff", "λ", 1},
		{"λ", "λa", -1},
	} {
		t.Run(tt.left+"/"+tt.right, func(t *testing.T) {
			got, err := ordering(runtime.Value{Kind: runtime.StringValue, Str: tt.left}, runtime.Value{Kind: runtime.StringValue, Str: tt.right})
			if err != nil || got != tt.order {
				t.Fatalf("ordering = %d, %v; want %d", got, err, tt.order)
			}
		})
	}
}
