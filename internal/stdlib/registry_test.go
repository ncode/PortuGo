package stdlib

import (
	"math"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/runtime"
)

func TestRegistryValidation(t *testing.T) {
	base := Descriptor{
		name: "sample", aliases: []string{"alias"},
		signature: Signature{Form: Parenthesized, Rule: UnaryNumericCall, Arity: 1,
			Parameters: [3]Parameter{{Type: runtime.NumericType, Mode: ByValue}},
			Result:     runtime.RealType, Empty: runtime.RealType},
		domain:   Domain{Description: "real numeric result"},
		evaluate: func(_ *Library, args []runtime.Value) (runtime.Value, bool, error) { return abs(args) },
	}
	for _, tt := range []struct {
		name string
		edit func(*Descriptor)
	}{
		{"empty name", func(d *Descriptor) { d.name = "" }},
		{"noncanonical name", func(d *Descriptor) { d.name = "SAMPLE" }},
		{"empty alias", func(d *Descriptor) { d.aliases = []string{""} }},
		{"noncanonical alias", func(d *Descriptor) { d.aliases = []string{"Alias"} }},
		{"invalid identifier", func(d *Descriptor) { d.aliases = []string{"bad-name"} }},
		{"duplicate alias", func(d *Descriptor) { d.aliases = []string{"sample"} }},
		{"missing form", func(d *Descriptor) { d.signature.Form = 0 }},
		{"missing rule", func(d *Descriptor) { d.signature.Rule = 0 }},
		{"contradictory form", func(d *Descriptor) { d.signature.Form = Bare }},
		{"contradictory arity", func(d *Descriptor) { d.signature.Arity = 0; d.signature.Parameters = [3]Parameter{} }},
		{"contradictory argument", func(d *Descriptor) { d.signature.Parameters[0].Type = runtime.StringType }},
		{"unsupported mode", func(d *Descriptor) { d.signature.Parameters[0].Mode = 2 }},
		{"contradictory result", func(d *Descriptor) { d.signature.Result = runtime.StringType }},
		{"contradictory empty result", func(d *Descriptor) { d.signature.Empty = runtime.StringType }},
		{"contradictory absence", func(d *Descriptor) { d.signature.Parameters[0].AbsenceAsZero = true }},
		{"unsupported arity", func(d *Descriptor) { d.signature.Arity = 4 }},
		{"negative arity", func(d *Descriptor) { d.signature.Arity = -1 }},
		{"missing parameter", func(d *Descriptor) { d.signature.Parameters[0].Type = runtime.InvalidType }},
		{"missing mode", func(d *Descriptor) { d.signature.Parameters[0].Mode = 0 }},
		{"extra parameter", func(d *Descriptor) { d.signature.Parameters[1] = d.signature.Parameters[0] }},
		{"missing result", func(d *Descriptor) { d.signature.Result = runtime.InvalidType }},
		{"missing empty result", func(d *Descriptor) { d.signature.Empty = runtime.InvalidType }},
		{"missing domain", func(d *Descriptor) { d.domain.Description = "" }},
		{"invalid absence", func(d *Descriptor) { d.domain.Absence = 99 }},
		{"missing evaluator", func(d *Descriptor) { d.evaluate = nil }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := base
			tt.edit(&d)
			if _, err := newRegistry([]Descriptor{d}); err == nil {
				t.Fatal("invalid descriptor accepted")
			}
		})
	}
	if _, err := newRegistry([]Descriptor{base, base}); err == nil {
		t.Fatal("duplicate canonical names accepted")
	}
	other := base
	other.name = "second"
	if _, err := newRegistry([]Descriptor{base, other}); err == nil {
		t.Fatal("shared alias accepted")
	}
	registry, err := newRegistry([]Descriptor{base})
	if err != nil {
		t.Fatal(err)
	}
	base.aliases[0] = "changed"
	if d := registry.byName["alias"]; d.Name() != "sample" || d.Names()[1] != "alias" {
		t.Fatal("registry retains mutable input aliases")
	}
}

func TestCatalogImmutable(t *testing.T) {
	entries := Catalog()
	if len(entries) != 28 {
		t.Fatalf("catalog has %d entries, want 28", len(entries))
	}
	for _, entry := range entries {
		d, ok := Lookup(entry.Name())
		if !ok || d.Name() != entry.Name() || d.Signature() != entry.Signature() {
			t.Fatalf("lookup disagrees for %q", entry.Name())
		}
		names := d.Names()
		names[0] = "changed"
		sig := d.Signature()
		sig.Parameters[0].Type = runtime.InvalidType
		if d.Names()[0] != entry.Name() || d.Signature() != entry.Signature() {
			t.Fatal("descriptor exposes mutable metadata")
		}
		if strings.TrimSpace(d.Domain().Description) == "" {
			t.Fatalf("%q has no domain", entry.Name())
		}
	}
	entries[0] = Descriptor{}
	if Catalog()[0].Name() == "" {
		t.Fatal("catalog exposes its storage")
	}
	if _, ok := Lookup("frac"); ok {
		t.Fatal("rejected legacy name remains in catalog")
	}
}

func TestCatalogDomainResults(t *testing.T) {
	integer := func(x int64) runtime.Value { return runtime.Value{Kind: runtime.IntegerValue, Int: x} }
	real := func(x float64) runtime.Value { return runtime.Value{Kind: runtime.RealValue, Real: x} }
	for _, row := range []struct {
		name    string
		args    []runtime.Value
		want    runtime.Value
		wantErr bool
	}{
		{"abs", []runtime.Value{integer(math.MinInt32)}, integer(math.MinInt32), false},
		{"quad", []runtime.Value{integer(65536)}, integer(0), false},
		{"int", []runtime.Value{real(2147483648)}, integer(math.MinInt32), false},
		{"exp", []runtime.Value{integer(0), integer(0)}, real(1), false},
		{"exp", []runtime.Value{integer(10), integer(-400)}, real(0), false},
		{"arccos", []runtime.Value{real(2)}, runtime.Value{Kind: runtime.VoidValue, NumericAbsence: true}, false},
		{"arcsen", []runtime.Value{real(-2)}, runtime.Value{Kind: runtime.VoidValue, NumericAbsence: true}, false},
		{"cotan", []runtime.Value{real(0)}, runtime.Value{Kind: runtime.VoidValue, NumericAbsence: true}, false},
		{"radpgrau", []runtime.Value{real(1e308)}, runtime.Value{Kind: runtime.VoidValue, NumericAbsence: true}, false},
		{"asc", []runtime.Value{{Kind: runtime.StringValue}}, runtime.Value{Kind: runtime.VoidValue}, false},
		{"carac", []runtime.Value{integer(-1)}, runtime.Value{Kind: runtime.VoidValue}, false},
		{"log", []runtime.Value{real(0)}, runtime.Value{}, true},
		{"logn", []runtime.Value{real(-1)}, runtime.Value{}, true},
		{"raizq", []runtime.Value{real(-1)}, runtime.Value{}, true},
		{"quad", []runtime.Value{real(1e200)}, runtime.Value{}, true},
		{"exp", []runtime.Value{integer(-1), real(.5)}, runtime.Value{}, true},
		{"exp", []runtime.Value{integer(10), integer(400)}, runtime.Value{}, true},
		{"int", []runtime.Value{real(math.Inf(1))}, runtime.Value{}, true},
		{"numpcarac", []runtime.Value{real(math.Inf(1))}, runtime.Value{}, true},
		{"caracpnum", []runtime.Value{{Kind: runtime.StringValue, Str: "01.5x"}}, runtime.Value{}, true},
		{"rand", nil, runtime.Value{}, true},
		{"randi", []runtime.Value{integer(1)}, runtime.Value{}, true},
	} {
		t.Run(row.name, func(t *testing.T) {
			got, found, err := New(nil).Call(row.name, row.args)
			if !found || (err != nil) != row.wantErr {
				t.Fatalf("domain result found=%t error=%v, want error=%t", found, err, row.wantErr)
			}
			if !row.wantErr && (got.Kind != row.want.Kind || got.Int != row.want.Int || got.Real != row.want.Real || got.NumericAbsence != row.want.NumericAbsence) {
				t.Fatalf("domain result=%+v, want %+v", got, row.want)
			}
		})
	}
}
