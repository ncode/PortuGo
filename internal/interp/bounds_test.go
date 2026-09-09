package interp

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func TestDynamicBoundFactsStayImmutable(t *testing.T) {
	p, info := analyzed(t, `algoritmo "layouts"
funcao value(x: inteiro): inteiro
const
  n = x
var
  v: vetor[1..n] de inteiro
inicio
  v[n] <- x
  retorne v[n]
fimfuncao
inicio
  escreval(value(2), "|", value(4))
fimalgoritmo`)
	fn := p.Subs[0].(*ast.FunctionDecl)
	name := fn.Locals[0].Names[0]
	before, ok := info.Binding(name)
	if !ok || before.Slots != 0 || !before.Type.DynamicBounds() {
		t.Fatalf("missing dynamic layout facts: %+v", before)
	}
	for range 2 {
		var out bytes.Buffer
		if ds := New(Options{Output: &out}).Run(p, info); len(ds) != 0 || out.String() != " 2| 4\n" {
			t.Fatalf("per-call bounds: %v, output %q", ds, &out)
		}
		after, _ := info.Binding(name)
		if !reflect.DeepEqual(before, after) {
			t.Fatal("execution mutated shared layout facts")
		}
	}
}

func TestConstantBoundAllocationGuard(t *testing.T) {
	for _, n := range []string{"1048577", "2147483647"} {
		src := "algoritmo \"allocation\"\nconst\nn=" + n + "\nvar\nv: vetor[1..n] de inteiro\ninicio\nfimalgoritmo"
		p, info := analyzed(t, src)
		i := New(Options{})
		ds := i.Run(p, info)
		if len(ds) != 1 || ds[0].Code != diag.RStorage || ds[0].Pos != token.Pos(strings.Index(src, "vetor")) {
			t.Fatalf("large bound was not rejected before allocation: %v", ds)
		}
		if _, allocated := i.State()["v"]; allocated {
			t.Fatal("failed layout allocated a vector")
		}
		p, info = analyzed(t, "algoritmo \"reuse\"\nconst\nn=2\nvar\nv: vetor[1..n] de inteiro\ninicio\nfimalgoritmo")
		if ds := i.Run(p, info); len(ds) != 0 || i.State()["v"].Kind != runtime.VectorValue {
			t.Fatalf("failed allocation polluted next run: %v", ds)
		}
	}
}
