package stdlib

import (
	"testing"

	"github.com/ncode/portugol-go/internal/runtime"
)

func TestExp(t *testing.T) {
	integer := runtime.Value{Kind: runtime.IntegerValue, Int: 2}
	real := runtime.Value{Kind: runtime.RealValue, Real: 3}
	text := runtime.Value{Kind: runtime.StringValue, Str: "2"}
	for _, tt := range []struct {
		name    string
		args    []runtime.Value
		want    float64
		wantErr bool
	}{
		{"integer", []runtime.Value{integer, integer}, 4, false},
		{"mixed", []runtime.Value{integer, real}, 8, false},
		{"real", []runtime.Value{real, real}, 27, false},
		{"none", nil, 0, true},
		{"one", []runtime.Value{integer}, 0, true},
		{"three", []runtime.Value{integer, real, real}, 0, true},
		{"bad base", []runtime.Value{text, integer}, 0, true},
		{"bad exponent", []runtime.Value{integer, text}, 0, true},
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
