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

func TestRecordedTextArgumentOrderGaps(t *testing.T) {
	for _, tt := range []struct {
		id     string
		stdout string
		code   diag.Code
		line   int
	}{
		{id: "text-order-no-value-bound", stdout: " 1\n 3\n"},
		{id: "text-order-bad-bound", stdout: " 1\n", code: diag.ETypeMismatch, line: 13},
		{id: "text-order-extra", stdout: " 1\n 2\n 3\n", code: diag.EParse, line: 13},
		{id: "text-order-no-value-code", stdout: " 9\n"},
		{id: "text-state-code-absent-stored", stdout: " 9\n"},
		{id: "text-state-code-absent-tail-error", stdout: " 9\n"},
	} {
		t.Run(tt.id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			for pass := range 2 {
				file, tokens, ds := lexer.Scan("source.alg", src)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				program, ds := parser.Parse(tokens)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				info, ds := sema.Analyze(program)
				if len(ds) != 0 {
					t.Fatalf("semantic diagnostics = %v", ds)
				}
				var out bytes.Buffer
				if len(ds) == 0 {
					ds = interp.New(interp.Options{Output: &out}).Run(program, info)
				}
				if out.String() != tt.stdout {
					t.Fatalf("pass %d output = %q, want %q", pass, out.String(), tt.stdout)
				}
				if tt.code == "" {
					if len(ds) != 0 {
						t.Fatalf("diagnostics = %v", ds)
					}
				} else if len(ds) != 1 || ds[0].Code != tt.code || pass == 0 && file.Position(ds[0].Pos).Line != tt.line {
					t.Fatalf("diagnostics = %v, want %s on line %d", ds, tt.code, tt.line)
				}
				var formatted bytes.Buffer
				if err := ast.Fprint(&formatted, program); err != nil {
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
