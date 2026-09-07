package parser

import (
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/testprocess"
)

func TestStructuralLimits(t *testing.T) {
	for _, tt := range []struct {
		name, body string
		limit      bool
	}{
		{"flat boundary", "escreva(" + strings.Repeat("1+", 254) + "1)", false},
		{"flat excessive", "escreva(" + strings.Repeat("1+", 255) + "1)", true},
		{"parentheses boundary", "escreva(" + strings.Repeat("(", 254) + "1" + strings.Repeat(")", 254) + ")", false},
		{"parentheses excessive", "escreva(" + strings.Repeat("(", 255) + "1" + strings.Repeat(")", 255) + ")", true},
		{"nested statements", strings.Repeat("se verdadeiro entao\n", 300) + strings.Repeat("fimse\n", 300), true},
		{"nested calls", "escreva(" + strings.Repeat("abs(", 300) + "1" + strings.Repeat(")", 300) + ")", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testprocess.Run(t, func() {
				_, tokens, ds := lexer.Scan("limit.alg", "algoritmo \"limit\"\ninicio\n"+tt.body+"\nfimalgoritmo")
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				p, ds := Parse(tokens)
				if tt.limit {
					if p != nil || len(ds) != 1 || ds[0].Code != diag.EResource || ds[0].Pos == 0 {
						t.Fatalf("missing single limit diagnostic: %v", ds)
					}
				} else if len(ds) != 0 {
					t.Fatalf("rejected boundary: %v", ds)
				}
			})
		})
	}
}
