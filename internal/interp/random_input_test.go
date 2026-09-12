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
