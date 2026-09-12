package stdlib

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/ncode/portugol-go/internal/runtime"
)

func TestRandomInputDraws(t *testing.T) {
	for _, tt := range []struct {
		name      string
		kind      runtime.TypeKind
		low, high float64
		decimals  int
		bounds    []uint64
		want      runtime.Value
	}{
		{"integer endpoints", runtime.IntegerType, -3, 5, 3, []uint64{9}, runtime.Value{Kind: runtime.IntegerValue, Int: 5}},
		{"fixed integer", runtime.IntegerType, 7, 7, 0, []uint64{1}, runtime.Value{Kind: runtime.IntegerValue, Int: 7}},
		{"fixed real", runtime.RealType, 2.75, 2.75, 0, []uint64{1}, runtime.Value{Kind: runtime.RealValue, Real: 2.75}},
		{"fractional real", runtime.RealType, 2, 2, 3, []uint64{1, 1000}, runtime.Value{Kind: runtime.RealValue, Real: 2.999}},
		{"text", runtime.StringType, 0, 0, 0, []uint64{26, 26, 26, 26, 26}, runtime.Value{Kind: runtime.StringValue, Str: "ZZZZZ"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := &recordingRandom{}
			v, err := New(r).RandomInput(tt.kind, tt.low, tt.high, tt.decimals)
			if err != nil || v != tt.want || !slices.Equal(r.bounds, tt.bounds) || r.draws != 0 {
				t.Fatalf("value=%v error=%v bounds=%v fractions=%d", v, err, r.bounds, r.draws)
			}
		})
	}
}

func TestRandomInputDomains(t *testing.T) {
	for seed := uint64(0); seed < 64; seed++ {
		lib := New(rand.New(rand.NewPCG(seed, seed+1)))
		for range 64 {
			integer, err := lib.RandomInput(runtime.IntegerType, -3, 5, 3)
			if err != nil || integer.Kind != runtime.IntegerValue || integer.Int < -3 || integer.Int > 5 {
				t.Fatalf("integer: seed=%d value=%v error=%v", seed, integer, err)
			}
			real, err := lib.RandomInput(runtime.RealType, 2, 2, 3)
			if err != nil || real.Kind != runtime.RealValue || real.Real < 2 || real.Real >= 3 || math.Abs(real.Real*1000-math.Round(real.Real*1000)) > 1e-9 {
				t.Fatalf("real: seed=%d value=%v error=%v", seed, real, err)
			}
			text, err := lib.RandomInput(runtime.StringType, 0, 0, 0)
			if err != nil || text.Kind != runtime.StringValue || len(text.Str) != 5 {
				t.Fatalf("text: seed=%d value=%v error=%v", seed, text, err)
			}
			for _, ch := range text.Str {
				if ch < 'A' || ch > 'Z' {
					t.Fatalf("random letter outside uppercase ASCII: %q", text.Str)
				}
			}
		}
	}
}

func TestRandomInputInvalidRange(t *testing.T) {
	for _, tt := range []struct {
		low, high float64
		decimals  int
	}{
		{math.NaN(), 1, 0}, {0, math.Inf(1), 0}, {math.Inf(-1), 0, 0},
		{0, 0x1p64, 0}, {3, 0, 0}, {0, 1, -1}, {0, 1, 6},
	} {
		r := &recordingRandom{}
		if _, err := New(r).RandomInput(runtime.RealType, tt.low, tt.high, tt.decimals); err == nil || len(r.bounds) != 0 || r.draws != 0 {
			t.Fatalf("invalid range accepted or consumed source: range=%+v error=%v bounds=%v", tt, err, r.bounds)
		}
	}
}
