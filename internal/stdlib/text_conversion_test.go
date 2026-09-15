package stdlib

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/ncode/PortuGo/internal/runtime"
)

func TestRecordedTextAndConversionContracts(t *testing.T) {
	integer := func(value int64) runtime.Value {
		return runtime.Value{Kind: runtime.IntegerValue, Int: value}
	}
	real := func(value float64) runtime.Value {
		return runtime.Value{Kind: runtime.RealValue, Real: value}
	}
	text := func(value string) runtime.Value {
		return runtime.Value{Kind: runtime.StringValue, Str: value}
	}
	for _, tt := range []struct {
		name    string
		builtin string
		args    []runtime.Value
		want    runtime.Value
	}{
		{"maiusc accents", "maiusc", []runtime.Value{text("ação µƒ")}, text("AÇÃO µƒ")},
		{"minusc accents", "minusc", []runtime.Value{text("AÇÃO")}, text("ação")},
		{"compr accents", "compr", []runtime.Value{text("ação€")}, integer(5)},
		{"pos first", "pos", []runtime.Value{text("ç"), text("ação")}, integer(2)},
		{"pos missing", "pos", []runtime.Value{text("x"), text("ação")}, integer(0)},
		{"pos empty", "pos", []runtime.Value{text(""), text("ação")}, integer(0)},
		{"copia one based", "copia", []runtime.Value{text("ação"), integer(2), integer(2)}, text("çã")},
		{"copia fractional bounds", "copia", []runtime.Value{text("ação"), real(2.9), real(2.5)}, text("çã")},
		{"copia clipped", "copia", []runtime.Value{text("ação"), integer(4), integer(9)}, text("o")},
		{"asc euro", "asc", []runtime.Value{text("€")}, integer(128)},
		{"asc empty", "asc", []runtime.Value{text("")}, runtime.Value{Kind: runtime.VoidValue}},
		{"carac extended", "carac", []runtime.Value{integer(128)}, text("Ç")},
		{"carac control", "carac", []runtime.Value{integer(0)}, text(" ")},
		{"carac omitted", "carac", nil, text(" ")},
		{"int truncates", "int", []runtime.Value{real(2.9)}, integer(2)},
		{"caracpnum integer", "caracpnum", []runtime.Value{text("123")}, integer(123)},
		{"caracpnum real", "caracpnum", []runtime.Value{text("12,5")}, real(12.5)},
		{"caracpnum zero fraction", "caracpnum", []runtime.Value{text("0.5")}, integer(0)},
		{"caracpnum real threshold", "caracpnum", []runtime.Value{text("01.5")}, real(1.5)},
		{"numpcarac integer", "numpcarac", []runtime.Value{integer(7)}, text("7")},
		{"numpcarac empty", "numpcarac", nil, text("0")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, found, err := New(nil).Call(tt.builtin, tt.args)
			if !found || err != nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("%s(%v) = %#v, found=%t, error=%v; want %#v", tt.builtin, tt.args, got, found, err, tt.want)
			}
		})
	}
}

func TestCopiaChecksResultSizeBeforeMaterializing(t *testing.T) {
	input := strings.Repeat("a", runtime.MaxTextBytes+1)
	got, found, err := New(nil).Call("copia", []runtime.Value{
		{Kind: runtime.StringValue, Str: input},
		{Kind: runtime.IntegerValue, Int: 1},
		{Kind: runtime.IntegerValue, Int: runtime.MaxTextBytes + 1},
	})
	if !found || !errors.Is(err, runtime.ErrTextSize) {
		t.Fatalf("copia result kind=%v length=%d found=%t error=%v, want ErrTextSize", got.Kind, len(got.Str), found, err)
	}
	got, found, err = New(nil).Call("copia", []runtime.Value{
		{Kind: runtime.StringValue, Str: input},
		{Kind: runtime.IntegerValue, Int: 1},
		{Kind: runtime.IntegerValue, Int: runtime.MaxTextBytes},
	})
	if !found || err != nil || got.Kind != runtime.StringValue || len(got.Str) != runtime.MaxTextBytes {
		t.Fatalf("copia boundary kind=%v length=%d found=%t error=%v, want a %d-byte result", got.Kind, len(got.Str), found, err, runtime.MaxTextBytes)
	}
}

func TestTextConversionContractHelpersRemainTyped(t *testing.T) {
	got, found, err := New(nil).Call("caracpnum", []runtime.Value{{Kind: runtime.StringValue, Str: "01.5"}})
	if !found || err != nil || got.Kind != runtime.RealValue || got.Real != 1.5 {
		t.Fatalf("caracpnum result=%#v found=%t error=%v, want real 1.5", got, found, err)
	}

	got, found, err = New(nil).Call("caracpnum", []runtime.Value{{Kind: runtime.StringValue, Str: "01.5x"}})
	if !found || err == nil || got != (runtime.Value{}) {
		t.Fatalf("malformed caracpnum result=%#v found=%t error=%v, want an error without a value", got, found, err)
	}
}
