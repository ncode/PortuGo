package interp

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func TestArgumentSetupFailure(t *testing.T) {
	src := `algoritmo "arguments"
var
  n: real
  v: vetor[1..1] de inteiro
procedimento P(var x: inteiro; y: inteiro)
inicio
  escreval("BODY")
  x <- 9
fimprocedimento
inicio
  n <- 3.5
  P(n, v[2])
fimalgoritmo`
	p, info := analyzed(t, src)
	var out bytes.Buffer
	i := New(Options{Output: &out})
	ds := i.Run(p, info)
	if len(ds) != 1 || ds[0].Code != diag.RStorage || out.Len() != 0 {
		t.Fatalf("failed setup entered body: diagnostics %v, output %q", ds, &out)
	}
	if got := i.State()["n"]; got.Kind != runtime.RealValue || got.Real != 3.5 || i.calls != 0 || i.env != i.global {
		t.Fatalf("failed setup changed caller: value %+v, calls %d", got, i.calls)
	}
	p, info = analyzed(t, "algoritmo \"reuse\"\ninicio\nescreval(1)\nfimalgoritmo")
	if ds := i.Run(p, info); len(ds) != 0 || out.String() != " 1\n" {
		t.Fatalf("failed setup polluted next run: %v %q", ds, &out)
	}
}

func TestIntegerArgumentRange(t *testing.T) {
	for _, tt := range []struct {
		argument string
		output   string
		code     diag.Code
	}{
		{"-9223372036854775808.0", " -9223372036854775808\n", ""},
		{"9223372036854774784.0", " 9223372036854774784\n", ""},
		{"9223372036854775808.0", "", diag.RCall},
		{"-9223372036854777856.0", "", diag.RCall},
		{"exp(10, 1000)", "", diag.RBuiltin},
		{"-exp(10, 1000)", "", diag.RBuiltin},
		{"raizq(-1)", "", diag.RBuiltin},
	} {
		t.Run(tt.argument, func(t *testing.T) {
			src := fmt.Sprintf("algoritmo \"arguments\"\nprocedimento P(x: inteiro)\ninicio\nescreval(x)\nfimprocedimento\ninicio\nP(%s)\nfimalgoritmo", tt.argument)
			p, info := analyzed(t, src)
			var out bytes.Buffer
			ds := New(Options{Output: &out}).Run(p, info)
			if tt.output == "" {
				if len(ds) != 1 || ds[0].Code != tt.code {
					t.Fatalf("out-of-range argument accepted: %v", ds)
				}
			} else if len(ds) != 0 {
				t.Fatal(ds)
			}
			if out.String() != tt.output {
				t.Fatalf("output = %q, want %q", &out, tt.output)
			}
		})
	}
}

func TestNumericAbsenceExpressionFacts(t *testing.T) {
	for _, tt := range []struct {
		name, call, want string
	}{
		{"unary plus", "exp(arccos(1), +arccos(2), marker())", "MARKER\n"},
		{"division", "exp(arccos(1), 2 / arccos(2), marker())", "MARKER\n"},
		{"integer division", "exp(arccos(1), 2 \\ arccos(2), marker())", "MARKER\n"},
		{"logical remainder", "exp(arccos(1), verdadeiro mod arccos(2), marker())", "MARKER\n"},
		{"power", "exp(arccos(1), 2 ^ arccos(2), marker())", ""},
		{"short power", "exp(2 ^ arccos(2))", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := fmt.Sprintf("algoritmo \"numeric absence\"\nfuncao marker: inteiro\ninicio\nescreval(\"MARKER\")\nretorne 2\nfimfuncao\ninicio\nescreval(%s)\nfimalgoritmo", tt.call)
			p, info := analyzed(t, src)
			var out bytes.Buffer
			if ds := New(Options{Output: &out}).Run(p, info); len(ds) != 0 || out.String() != tt.want {
				t.Fatalf("diagnostics %v, output %q, want %q", ds, &out, tt.want)
			}
		})
	}
}

func TestFunctionTailFailureAndReuse(t *testing.T) {
	p, info := analyzed(t, `algoritmo "returns"
var
  n: inteiro
funcao F(var x: inteiro): inteiro
inicio
  retorne x
  x <- 9
  escreval(1 \ 0)
fimfuncao
inicio
  n <- 7
  escreval(F(n))
fimalgoritmo`)
	var out bytes.Buffer
	i := New(Options{Output: &out})
	ds := i.Run(p, info)
	if len(ds) != 1 || ds[0].Code != diag.RArithmetic || out.Len() != 0 {
		t.Fatalf("tail failure: diagnostics %v, output %q", ds, &out)
	}
	if got := i.State()["n"]; got.Int != 7 || i.calls != 0 || i.env != i.global {
		t.Fatalf("failed body changed caller: value %+v, calls %d", got, i.calls)
	}
	p, info = analyzed(t, "algoritmo \"reuse\"\nfuncao F: inteiro\ninicio\nfimfuncao\ninicio\nescreval(F())\nfimalgoritmo")
	if ds := i.Run(p, info); len(ds) != 0 || out.String() != " 0\n" {
		t.Fatalf("failed call polluted next run: %v %q", ds, &out)
	}
}

func TestCrossTypeFallthroughUsesTypedZero(t *testing.T) {
	p, info := analyzed(t, `algoritmo "returns"
var
  a: inteiro
  b: real
funcao F: inteiro
inicio
  retorne 9
fimfuncao
funcao G: real
inicio
fimfuncao
inicio
  a <- F()
  b <- G()
  escreval(a, b)
fimalgoritmo`)
	var out bytes.Buffer
	if ds := New(Options{Output: &out}).Run(p, info); len(ds) != 0 || out.String() != " 9 0\n" {
		t.Fatalf("cross-type fallthrough: %v, %q", ds, &out)
	}
}
