package stdlib

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/runtime"
)

func TestCaseConversionAllocationGuard(t *testing.T) {
	// Direct callers can supply values larger than language strings allow.
	text := strings.Repeat("\u023f", runtime.MaxTextBytes/2)
	_, found, err := New(nil).Call("maiusc", []runtime.Value{{Kind: runtime.StringValue, Str: text}})
	if !found || !errors.Is(err, runtime.ErrTextSize) {
		t.Fatalf("case expansion found=%t error=%v, want the allocation guard", found, err)
	}
}

func TestNumpcaracRejectsUnsupportedValues(t *testing.T) {
	for _, args := range [][]runtime.Value{
		{{Kind: runtime.RealValue, Real: math.NaN()}},
		{{Kind: runtime.RealValue, Real: math.Inf(1)}},
		{{Kind: runtime.RealValue, Real: math.Inf(-1)}},
		{{Kind: runtime.IntegerValue}, {Kind: runtime.IntegerValue}},
		{{Kind: runtime.VectorValue}},
		{{}},
	} {
		if _, found, err := New(nil).Call("numpcarac", args); !found || err == nil {
			t.Fatalf("arguments=%v found=%t error=%v, want a conversion error", args, found, err)
		}
	}
}

func TestCopiaExtremeLength(t *testing.T) {
	for _, length := range []int64{math.MaxInt32, math.MaxInt64} {
		args := []runtime.Value{
			{Kind: runtime.StringValue, Str: "abc"},
			{Kind: runtime.IntegerValue, Int: 2},
			{Kind: runtime.IntegerValue, Int: length},
		}
		got, found, err := New(nil).Call("copia", args)
		if !found || err != nil || got.Kind != runtime.StringValue || got.Str != "bc" {
			t.Fatalf("length=%d result=%v found=%t error=%v", length, got, found, err)
		}
	}
}

func TestIntRejectsUnsupportedRange(t *testing.T) {
	for _, value := range []float64{
		math.NaN(), math.Inf(1), math.Inf(-1), math.MaxFloat64,
		math.Ldexp(1, 63), math.Nextafter(-math.Ldexp(1, 63), math.Inf(-1)),
	} {
		_, found, err := New(nil).Call("int", []runtime.Value{{Kind: runtime.RealValue, Real: value}})
		if !found || err == nil {
			t.Fatalf("value=%g found=%t error=%v, want a conversion error", value, found, err)
		}
	}
}

func TestExp(t *testing.T) {
	integer := runtime.Value{Kind: runtime.IntegerValue, Int: 2}
	real := runtime.Value{Kind: runtime.RealValue, Real: 3}
	for _, tt := range []struct {
		name    string
		args    []runtime.Value
		want    float64
		wantErr bool
	}{
		{"integer", []runtime.Value{integer, integer}, 4, false},
		{"mixed", []runtime.Value{integer, real}, 8, false},
		{"real", []runtime.Value{real, real}, 27, false},
		{"one", []runtime.Value{integer}, 0, true},
		{"three", []runtime.Value{integer, real, real}, 0, true},
		{"domain", []runtime.Value{{Kind: runtime.IntegerValue, Int: -1}, {Kind: runtime.RealValue, Real: 0.5}}, 0, true},
		{"overflow", []runtime.Value{{Kind: runtime.IntegerValue, Int: 10}, {Kind: runtime.IntegerValue, Int: 400}}, 0, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, found, err := New(nil).Call("exp", tt.args)
			if !found || (err != nil) != tt.wantErr {
				t.Fatalf("found = %v, error = %v, want error %v", found, err, tt.wantErr)
			}
			if !tt.wantErr && (got.Kind != runtime.RealValue || got.Real != tt.want) {
				t.Fatalf("got %#v, want real %v", got, tt.want)
			}
		})
	}
}
