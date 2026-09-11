package portugol_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/stdlib"
)

// This table is independent of Catalog: changing both implementation consumers
// cannot change the required signatures, bindings, or expected values here.
func TestBuiltinRegistryAgreement(t *testing.T) {
	const (
		i = runtime.IntegerType
		r = runtime.RealType
		s = runtime.StringType
		n = runtime.NumericType
	)
	rows := []struct {
		name, expression, parameters string
		result, binding              runtime.TypeKind
		absence                      stdlib.AbsenceRule
		output                       string
	}{
		{"abs", "abs(-7)", "n", n, i, 0, " 7\n"},
		{"arccos", "arccos(1)", "n", r, r, stdlib.NumericDomainAbsence, " 0\n"},
		{"arcsen", "arcsen(0)", "n", r, r, stdlib.NumericDomainAbsence, " 0\n"},
		{"arctan", "arctan(0)", "n", r, r, 0, " 0\n"},
		{"asc", `asc("A")`, "s", i, i, stdlib.GenericDomainAbsence, " 65\n"},
		{"carac", "carac(65)", "i", s, s, stdlib.GenericDomainAbsence, "A\n"},
		{"caracpnum", `caracpnum("7")`, "s", n, n, 0, " 7\n"},
		{"compr", `compr("ABC")`, "s", i, i, 0, " 3\n"},
		{"copia", `copia("ABC",2,1)`, "snn", s, s, 0, "B\n"},
		{"cos", "cos(0)", "n", r, r, 0, " 1\n"},
		{"cotan", "cotan(1)", "n", r, r, stdlib.NumericDomainAbsence, " 0.642092615934331\n"},
		{"exp", "exp(2,3)", "nn", r, r, 0, " 8\n"},
		{"grauprad", "grauprad(0)", "n", r, r, 0, " 0\n"},
		{"int", "int(2.9)", "n", i, i, 0, " 2\n"},
		{"log", "log(1)", "n", r, r, 0, " 0\n"},
		{"logn", "logn(1)", "n", r, r, 0, " 0\n"},
		{"maiusc", `maiusc("abc")`, "s", s, s, 0, "ABC\n"},
		{"minusc", `minusc("ABC")`, "s", s, s, 0, "abc\n"},
		{"numpcarac", "numpcarac(7)", "n", s, s, 0, "7\n"},
		{"pi", "pi", "", r, r, 0, " 3.14159265358979\n"},
		{"pos", `pos("B","ABC")`, "ss", i, i, 0, " 2\n"},
		{"quad", "quad(2)", "n", n, i, 0, " 4\n"},
		{"radpgrau", "radpgrau(0)", "n", r, r, stdlib.NumericDomainAbsence, " 0\n"},
		{"raizq", "raizq(4)", "n", r, r, 0, " 2\n"},
		{"rand", "rand", "", r, r, 0, " 0.25\n"},
		{"randi", "randi(1)", "i", i, i, 0, " 0\n"},
		{"sen", "sen(0)", "n", r, r, 0, " 0\n"},
		{"tan", "tan(0)", "n", r, r, 0, " 0\n"},
	}
	covered := make(map[string]bool)
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			d, ok := stdlib.Lookup(row.name)
			if !ok || covered[row.name] {
				t.Fatal("missing or duplicate catalog coverage")
			}
			covered[row.name] = true
			sig := d.Signature()
			form := stdlib.Parenthesized
			if row.parameters == "" {
				form = stdlib.Bare
			}
			if sig.Arity != len(row.parameters) || sig.Result != row.result || sig.Form != form || d.Domain().Absence != row.absence {
				t.Fatalf("signature/domain drift: %+v %+v", sig, d.Domain())
			}
			for index, parameter := range row.parameters {
				want := map[rune]runtime.TypeKind{'i': i, 'r': r, 's': s, 'n': n}[parameter]
				if sig.Parameters[index].Type != want || sig.Parameters[index].Mode != stdlib.ByValue {
					t.Fatalf("parameter %d disagrees with independent signature", index)
				}
			}
			for _, spelling := range []string{row.name, strings.ToUpper(row.name)} {
				expression := strings.Replace(row.expression, row.name, spelling, 1)
				src := "algoritmo \"catalog\"\ninicio\nescreval(" + expression + ")\nfimalgoritmo\n"
				_, tokens, ds := lexer.Scan("catalog.alg", src)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				program, ds := parser.Parse(tokens)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				info, ds := sema.Analyze(program)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				found := false
				for _, tok := range tokens {
					if tok.Text == spelling {
						binding, ok := info.Binding(tok)
						if !ok || !binding.Builtin || binding.Type.Kind != row.binding {
							t.Fatalf("semantic binding disagrees: %+v", binding)
						}
						found = true
					}
				}
				if !found {
					t.Fatal("call was not bound")
				}
				var output bytes.Buffer
				if ds := interp.New(interp.Options{Output: &output, Random: catalogRandom{}}).Run(program, info); len(ds) != 0 {
					t.Fatal(ds)
				}
				if output.String() != row.output {
					t.Fatalf("runtime output = %q, want %q", output.String(), row.output)
				}
			}
		})
	}
	for _, d := range stdlib.Catalog() {
		if !covered[d.Name()] {
			t.Errorf("no signature/runtime row for %s", d.Name())
		}
	}
	t.Logf("catalog agreement: %d descriptors, semantic bindings and runtime outputs in both letter cases", len(covered))
}

type catalogRandom struct{}

func (catalogRandom) Float64() float64        { return 0.25 }
func (catalogRandom) Uint64N(n uint64) uint64 { return n - 1 }
