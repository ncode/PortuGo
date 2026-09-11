package stdlib

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/ncode/portugol-go/internal/runtime"
)

type recordingRandom struct {
	bounds   []uint64
	bad      bool
	fraction float64
	draws    int
}

func (r *recordingRandom) Float64() float64 { r.draws++; return r.fraction }
func (r *recordingRandom) Uint64N(n uint64) uint64 {
	r.bounds = append(r.bounds, n)
	if r.bad {
		return n
	}
	return n - 1
}

func TestRandSource(t *testing.T) {
	for _, value := range []float64{0, .25, math.Nextafter(1, 0), -1, 1, math.NaN(), math.Inf(1)} {
		r := &recordingRandom{fraction: value}
		v, found, err := New(r).Call("rand", nil)
		valid := value >= 0 && value < 1
		if !found || (err == nil) != valid || r.draws != 1 || len(r.bounds) != 0 || valid && (v.Kind != runtime.RealValue || v.Real != value) {
			t.Fatalf("draw %g: result=%v found=%t error=%v calls=%d bounds=%v", value, v, found, err, r.draws, r.bounds)
		}
	}
	if _, found, err := New(nil).Call("rand", nil); !found || err == nil {
		t.Fatalf("missing source: found=%t error=%v", found, err)
	}
	r := &recordingRandom{}
	if _, found, err := New(r).Call("rand", []runtime.Value{{Kind: runtime.IntegerValue}}); !found || err == nil || r.draws != 0 {
		t.Fatalf("unexpected argument: found=%t error=%v draws=%d", found, err, r.draws)
	}
}

func TestRandDomains(t *testing.T) {
	for seed := uint64(0); seed < 64; seed++ {
		lib := New(rand.New(rand.NewPCG(seed, seed+1)))
		for range 64 {
			v, found, err := lib.Call("rand", nil)
			if !found || err != nil || v.Kind != runtime.RealValue || !(v.Real >= 0 && v.Real < 1) {
				t.Fatalf("seed=%d result=%v found=%t error=%v", seed, v, found, err)
			}
		}
	}
}

func TestLegacyRandomSourceFailures(t *testing.T) {
	for _, args := range [][]runtime.Value{
		{{Kind: runtime.IntegerValue, Int: 7}},
		{{Kind: runtime.IntegerValue, Int: -3}, {Kind: runtime.IntegerValue, Int: 5}},
	} {
		for _, source := range []RandomSource{&recordingRandom{bad: true}, nil} {
			if _, found, err := New(source).Call("aleatorio", args); !found || err == nil {
				t.Fatalf("invalid source: found=%t error=%v", found, err)
			}
		}
	}
}

func TestRandiSource(t *testing.T) {
	for _, tt := range []struct {
		name  string
		args  []runtime.Value
		bound uint64
		want  int64
	}{
		{"empty", nil, 0, 0},
		{"zero", []runtime.Value{{Kind: runtime.IntegerValue}}, 0, 0},
		{"unit", []runtime.Value{{Kind: runtime.IntegerValue, Int: 1}}, 1, 0},
		{"exclusive upper bound", []runtime.Value{{Kind: runtime.IntegerValue, Int: 7}}, 7, 6},
		{"negative bound", []runtime.Value{{Kind: runtime.IntegerValue, Int: -7}}, 4294967289, -8},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := &recordingRandom{}
			got, found, err := New(r).Call("randi", tt.args)
			if err != nil || !found || got.Kind != runtime.IntegerValue || got.Int != tt.want {
				t.Fatalf("result=%v found=%t error=%v", got, found, err)
			}
			if tt.bound == 0 && len(r.bounds) != 0 || tt.bound != 0 && (len(r.bounds) != 1 || r.bounds[0] != tt.bound) {
				t.Fatalf("draws=%v, expected bound %d", r.bounds, tt.bound)
			}
		})
	}
}

func TestRandiFailures(t *testing.T) {
	for _, args := range [][]runtime.Value{
		{{Kind: runtime.RealValue, Real: 1}},
		{{Kind: runtime.StringValue, Str: "7"}},
		{{Kind: runtime.BoolValue}},
		{{Kind: runtime.VoidValue}},
		{{Kind: runtime.VectorValue}},
		{{Kind: runtime.IntegerValue, Int: 2147483648}},
		{{Kind: runtime.IntegerValue, Int: -2147483649}},
		{{Kind: runtime.IntegerValue}, {Kind: runtime.IntegerValue}},
	} {
		r := &recordingRandom{}
		if _, found, err := New(r).Call("randi", args); !found || err == nil || len(r.bounds) != 0 {
			t.Fatalf("arguments=%v found=%t error=%v draws=%v", args, found, err, r.bounds)
		}
	}
	for _, r := range []RandomSource{nil, &recordingRandom{bad: true}} {
		if _, found, err := New(r).Call("randi", []runtime.Value{{Kind: runtime.IntegerValue, Int: 7}}); !found || err == nil {
			t.Fatalf("invalid source: found=%t error=%v", found, err)
		}
	}
}

func TestRandiDomains(t *testing.T) {
	for seed := uint64(0); seed < 64; seed++ {
		lib := New(rand.New(rand.NewPCG(seed, seed+1)))
		for _, bound := range []int64{1, 2, 7, 2147483647, -7} {
			for draw := 0; draw < 64; draw++ {
				got, found, err := lib.Call("randi", []runtime.Value{{Kind: runtime.IntegerValue, Int: bound}})
				outside := got.Int < -2147483648 || got.Int > 2147483647
				if bound > 0 {
					outside = outside || got.Int < 0 || got.Int >= bound
				} else {
					outside = outside || got.Int < 0 && got.Int >= bound
				}
				if !found || err != nil || got.Kind != runtime.IntegerValue || outside {
					t.Fatalf("seed=%d bound=%d result=%v found=%t error=%v", seed, bound, got, found, err)
				}
			}
		}
	}
}
