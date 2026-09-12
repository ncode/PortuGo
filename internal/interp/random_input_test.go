package interp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

func TestRandomInputFailure(t *testing.T) {
	for _, kind := range []string{"inteiro", "real", "caractere"} {
		t.Run(kind, func(t *testing.T) {
			src := "algoritmo \"input failure\"\nvar value: " + kind + "\ninicio\nescreva(\"before\")\naleatorio 7,7\nleia(value)\nfimalgoritmo"
			p, info := analyzed(t, src)
			r := &scriptedRandom{bad: true, badFraction: true}
			var out bytes.Buffer
			i := New(Options{Random: r, Output: &out})
			ds := i.Run(p, info)
			at := token.Pos(strings.Index(src, "value)"))
			if len(ds) != 1 || ds[0].Code != diag.RInput || ds[0].Pos != at || out.String() != "before" {
				t.Fatalf("random-input failure: diagnostics=%v output=%q", ds, out.String())
			}
		})
	}
}

func TestRandomInputRunReset(t *testing.T) {
	random, randomInfo := analyzed(t, "algoritmo \"random mode\"\nvar value: inteiro\ninicio\naleatorio 7,7\nleia(value)\nfimalgoritmo")
	console, consoleInfo := analyzed(t, "algoritmo \"console mode\"\nvar value: inteiro\ninicio\nleia(value)\nfimalgoritmo")
	var out bytes.Buffer
	i := New(Options{Random: &scriptedRandom{}, Input: strings.NewReader("41\n"), Output: &out})
	if ds := i.Run(random, randomInfo); len(ds) != 0 {
		t.Fatal(ds)
	}
	if ds := i.Run(console, consoleInfo); len(ds) != 0 || out.String() != "7\n41\n" {
		t.Fatalf("input mode leaked between runs: diagnostics=%v output=%q", ds, out.String())
	}
}

func TestRandomInputRejectsUnrepresentableIntegerRangeBeforeDraw(t *testing.T) {
	src := "algoritmo \"integer range\"\nvar\nlow, high: real\nvalue: inteiro\ninicio\nlow <- 9223372036854775808\nhigh <- low\nescreva(\"before\")\naleatorio low, high\nleia(value)\nfimalgoritmo"
	p, info := analyzed(t, src)
	r := &scriptedRandom{}
	var out bytes.Buffer
	ds := New(Options{Random: r, Output: &out}).Run(p, info)
	at := token.Pos(strings.Index(src, "value)"))
	if len(ds) != 1 || ds[0].Code != diag.RInput || ds[0].Pos != at || out.String() != "before" || len(r.bounds) != 0 {
		t.Fatalf("unrepresentable integer range: diagnostics=%v output=%q bounds=%v", ds, out.String(), r.bounds)
	}
}

func TestRandomInputSwapsReversedFullSigned32Bounds(t *testing.T) {
	src := "algoritmo \"reversed range\"\nvar\nlow, high, value: inteiro\ninicio\nlow <- 2147483647\nhigh <- -2147483647 - 1\naleatorio low, high\nleia(value)\nescreval(value)\nfimalgoritmo"
	p, info := analyzed(t, src)
	r := &scriptedRandom{}
	var out bytes.Buffer
	if ds := New(Options{Random: r, Output: &out}).Run(p, info); len(ds) != 0 {
		t.Fatalf("reversed full signed-32 range: %v", ds)
	}
	if out.String() != "2147483647\n 2147483647\n" {
		t.Fatalf("reversed full signed-32 output=%q", out.String())
	}
	if len(r.bounds) != 1 || r.bounds[0] != 1<<32 {
		t.Fatalf("reversed full signed-32 bounds=%v", r.bounds)
	}
}
