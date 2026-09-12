package interp

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
)

func TestConstantInitialization(t *testing.T) {
	p, info := analyzed(t, `algoritmo "constants"
const
  first = randi(3)
  second = randi(5)
var
procedimento show
const
  local = randi(7)
var
inicio
  escreval(first, "|", second, "|", local, "|", local)
fimprocedimento
inicio
  show()
  show()
fimalgoritmo`)
	r := &scriptedRandom{}
	var out bytes.Buffer
	i := New(Options{Random: r, Output: &out})
	for run := 0; run < 2; run++ {
		out.Reset()
		r.bounds = nil
		if ds := i.Run(p, info); len(ds) != 0 {
			t.Fatal(ds)
		}
		if !reflect.DeepEqual(r.bounds, []uint64{3, 5, 7, 7}) || out.String() != " 2| 4| 6| 6\n 2| 4| 6| 6\n" {
			t.Fatalf("constant initialization: draws %v, output %q", r.bounds, &out)
		}
	}
}

func TestConstantInitializationFailure(t *testing.T) {
	p, info := analyzed(t, `algoritmo "constants"
var
  x: inteiro
procedimento fail(var y: inteiro)
const
  n = 1 / 0
var
inicio
  y <- 9
fimprocedimento
inicio
  x <- 3
  fail(x)
fimalgoritmo`)
	i := New(Options{})
	ds := i.Run(p, info)
	if len(ds) != 1 || ds[0].Code != diag.RArithmetic || i.State()["x"].Int != 3 || i.calls != 0 || i.env != i.global {
		t.Fatalf("failed constant initialization changed caller: %v, calls %d", ds, i.calls)
	}
	p, info = analyzed(t, "algoritmo \"reuse\"\nconst\nn=2\nvar\ninicio\nfimalgoritmo")
	if ds := i.Run(p, info); len(ds) != 0 {
		t.Fatalf("failed initialization polluted next run: %v", ds)
	}
}
