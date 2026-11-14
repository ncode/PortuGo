package interp

import (
	"bytes"
	"strconv"
	"testing"
)

func TestSharedSemanticInfo(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"shared\"\nvar v: vetor[1..2] de inteiro\ninicio\nv[1] <- 7\nescreval(v[1])\nfimalgoritmo\n")
	for n := range 4 {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			t.Parallel()
			for range 5 {
				binding, _ := info.Binding(p.Globals[0].Names[0])
				binding.Type.Ranges[0].Low = 99
				var out bytes.Buffer
				i := New(Options{Output: &out})
				if ds := i.Run(p, info); len(ds) != 0 || out.String() != " 7\n" {
					t.Fatalf("shared info changed: %v output=%q", ds, &out)
				}
			}
		})
	}
}
