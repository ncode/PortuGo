package interp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func TestInputFailurePreservesValue(t *testing.T) {
	for _, tt := range []struct{ name, kind, input string }{
		{"exhaustion", "inteiro", ""},
		{"integer range", "inteiro", "9223372036854775808\n"},
		{"real range", "real", "1e9999\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := "algoritmo \"failure\"\nvar x: " + tt.kind + "\ninicio\nx <- 9\nescreva(\"BEFORE\")\nleia(x)\nfimalgoritmo"
			p, info := analyzed(t, src)
			var out bytes.Buffer
			i := New(Options{Input: strings.NewReader(tt.input), Output: &out})
			ds := i.Run(p, info)
			if len(ds) != 1 || ds[0].Code != diag.RInput || int(ds[0].Pos) != strings.LastIndex(src, "x)") || out.String() != "BEFORE" {
				t.Fatalf("diagnostics=%v output=%q, want positioned R004 and preceding output", ds, &out)
			}
			binding, ok := info.Binding(p.Globals[0].Names[0])
			if !ok {
				t.Fatal("missing binding")
			}
			cell, ok := i.global.lookup(binding.ID)
			want := runtime.Value{Kind: runtime.IntegerValue, Int: 9}
			if tt.kind == "real" {
				want = runtime.Value{Kind: runtime.RealValue, Real: 9}
			}
			if !ok || cell.Value != want {
				t.Fatal("failed input changed the destination")
			}
		})
	}
}

func TestWriteStateAfterFailure(t *testing.T) {
	failed, failedInfo := analyzed(t, "algoritmo \"failed\"\ninicio\nescreva(\"before\")\nescreval(1 / 0)\nfimalgoritmo")
	next, nextInfo := analyzed(t, "algoritmo \"next\"\ninicio\nescreva(\"after\")\nfimalgoritmo")
	var out bytes.Buffer
	i := New(Options{Output: &out})
	if ds := i.Run(failed, failedInfo); len(ds) != 1 || ds[0].Code != diag.RArithmetic {
		t.Fatalf("failed write diagnostics = %v", ds)
	}
	if ds := i.Run(next, nextInfo); len(ds) != 0 {
		t.Fatal(ds)
	}
	if out.String() != "beforeafter" {
		t.Fatalf("write state survived a failed run: %q", &out)
	}
}

func TestWriteBufferLimit(t *testing.T) {
	value := strings.Repeat("x", maxTextBytes/2+1)
	for _, tt := range []struct{ name, declarations, body string }{
		{"flat", "", "escreval(s, s)"},
		{"nested", "funcao F: inteiro\ninicio\nescreva(s)\nretorne 1\nfimfuncao\n", "escreval(s, F())"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p, info := analyzed(t, "algoritmo \"buffer\"\nvar\ns: caractere\n"+tt.declarations+"inicio\ns <- \""+value+"\"\n"+tt.body+"\nfimalgoritmo")
			var out bytes.Buffer
			ds := New(Options{Output: &out}).Run(p, info)
			if len(ds) != 1 || ds[0].Code != diag.RStorage || out.Len() != 0 {
				t.Fatalf("unbounded or partially emitted write: diagnostics=%v bytes=%d", ds, out.Len())
			}
		})
	}
	p, info := analyzed(t, "algoritmo \"reuse\"\nvar\ns: caractere\ninicio\ns <- \""+value+"\"\nescreva(s)\nescreva(s)\nfimalgoritmo")
	if ds := New(Options{}).Run(p, info); len(ds) != 0 {
		t.Fatalf("completed statement did not release its buffer: %v", ds)
	}
}

func TestWriteEvaluationFailure(t *testing.T) {
	for _, tt := range []struct {
		name, declarations, body, output, operator string
	}{
		{"value", "", "escreval(\"PREFIX\", 1 / 0)", "", "/"},
		{"width", "", "escreval(\"PREFIX\", 5:(1 \\ 0))", "", "\\"},
		{"nested", "funcao F: inteiro\ninicio\nescreva(\"INNER\")\nretorne 7\nfimfuncao\n", "escreval(\"PREFIX\", F(), 1 / 0)", "INNER\n", "/"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := "algoritmo \"failure\"\n" + tt.declarations + "inicio\n" + tt.body + "\nfimalgoritmo"
			p, info := analyzed(t, src)
			var out bytes.Buffer
			ds := New(Options{Output: &out}).Run(p, info)
			if len(ds) != 1 || ds[0].Code != diag.RArithmetic || int(ds[0].Pos) != strings.LastIndex(src, tt.operator) || out.String() != tt.output {
				t.Fatalf("diagnostics=%v output=%q, want positioned R002 and %q", ds, &out, tt.output)
			}
		})
	}
}
