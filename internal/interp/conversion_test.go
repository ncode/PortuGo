package interp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func TestNumericConversionUsesRuntimeType(t *testing.T) {
	src := "algoritmo \"dynamic conversion\"\nvar\ns: caractere\ninicio\nleia(s)\nescreval(randi(caracpnum(s)))\nfimalgoritmo"
	p, info := analyzed(t, src)
	call := p.Body[1].(*ast.WriteStmt).Args[0].Expr.(*ast.CallExpr)
	if typ, ok := info.TypeOf(call.Args[0]); !ok || typ.Kind != runtime.NumericType {
		t.Fatalf("conversion type = %v, present = %v; want dynamic numeric type", typ, ok)
	}
	for _, tt := range []struct {
		input   string
		output  string
		invalid bool
	}{
		{"1\n", "1\n 0\n", false},
		{"1.0\n", "1.0\n", true},
	} {
		t.Run(tt.input, func(t *testing.T) {
			var out bytes.Buffer
			i := New(Options{Input: strings.NewReader(tt.input), Output: &out})
			ds := i.Run(p, info)
			if tt.invalid {
				if len(ds) != 1 || ds[0].Code != diag.ETypeMismatch || ds[0].Pos != token.Pos(strings.Index(src, "caracpnum")) {
					t.Fatalf("real conversion result: %v", ds)
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
