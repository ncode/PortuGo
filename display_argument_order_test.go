package portugol_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedDisplayArgumentOrder(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		out  string
	}{
		{"display-color-bad-second", diag.ETypeMismatch, "amarelo\n"},
		{"display-color-no-value-first", diag.EParse, ""},
		{"display-color-no-value-second", diag.EParse, "amarelo\n"},
	} {
		t.Run(tt.id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			for pass := range 2 {
				file, tokens, lexDiags := lexer.Scan("source.alg", src)
				if len(lexDiags) != 0 {
					t.Fatal(lexDiags)
				}
				prog, parseDiags := parser.Parse(tokens)
				if len(parseDiags) != 0 {
					t.Fatal(parseDiags)
				}
				info, diags := sema.Analyze(prog)
				var out bytes.Buffer
				if len(diags) == 0 {
					diags = interp.New(interp.Options{Output: &out}).Run(prog, info)
				}
				if len(diags) != 1 || diags[0].Code != tt.code || out.String() != tt.out {
					t.Fatalf("pass %d diagnostics=%v output=%q", pass, diags, out.String())
				}
				if pass == 0 && file.Position(diags[0].Pos).Line != 8 {
					t.Fatalf("diagnostic line=%d, want 8", file.Position(diags[0].Pos).Line)
				}
				var formatted bytes.Buffer
				if err := ast.Fprint(&formatted, prog); err != nil {
					t.Fatal(err)
				}
				if pass == 1 && formatted.String() != src {
					t.Fatal("formatting is not idempotent")
				}
				src = formatted.String()
			}
		})
	}
}
