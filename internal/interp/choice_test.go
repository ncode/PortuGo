package interp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/token"
)

func TestDynamicChoiceLabelDiagnostic(t *testing.T) {
	for _, tt := range []struct{ name, value, label string }{
		{"integer selector", `"x" = 7`, `"x"`},
		{"text selector", `7 = "x"`, "7"},
		{"logical selector", "verdadeiro = falso", "7"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := "algoritmo \"choice\"\nvar b: logico\ninicio\nb <- " + tt.value + "\nescreval(\"BEFORE\")\nescolha b\ncaso " + tt.label + ":\nescreval(\"BAD\")\nfimescolha\nfimalgoritmo\n"
			p, info := analyzed(t, src)
			var out bytes.Buffer
			ds := New(Options{Output: &out}).Run(p, info)
			at := token.Pos(strings.Index(src, "caso ") + len("caso "))
			if len(ds) != 1 || ds[0].Code != diag.ETypeMismatch || ds[0].Pos != at || out.String() != "BEFORE\n" {
				t.Fatalf("diagnostics=%v output=%q, want E001 at %d after BEFORE", ds, &out, at)
			}
		})
	}
}
