package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedConstantBounds(t *testing.T) {
	for _, id := range []string{
		"constant-vector-bound", "constant-context-negative-bound",
		"constant-context-expression-bound", "constant-context-local-vector-bound",
		"constant-bound-per-call", "constant-bound-case-insensitive",
		"constant-bound-dimensions", "constant-bound-size-5001",
	} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			checkFormattingPreservesExecution(t, src, want)
		})
	}
}

func TestVectorDeclarationDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		line int
		out  string
	}{
		{"constant-bound-global-variable", 4, ""},
		{"constant-bound-local-variable", 5, ""},
		{"constant-bound-parameter-direct", 4, ""},
		{"constant-bound-logical", 5, ""},
		{"constant-bound-reversed", 5, ""},
		{"constant-bound-local-reversed", 6, " 7\n"},
		{"constant-bound-vector-parameter", 6, ""},
		{"vector-context-procedure-parameter", 2, ""},
		{"vector-context-function-parameter", 2, ""},
		{"vector-context-function-result", 4, ""},
	} {
		t.Run(tt.id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			var out bytes.Buffer
			if len(ds) == 0 {
				prog, parseDiags := parser.Parse(tokens)
				ds = parseDiags
				if len(ds) == 0 {
					var info *sema.Info
					info, ds = sema.Analyze(prog)
					if len(ds) == 0 {
						ds = interp.New(interp.Options{Output: &out}).Run(prog, info)
					}
				}
			}
			if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != tt.line || out.String() != tt.out {
				t.Fatalf("diagnostics = %v, output = %q; want P001 on line %d, output %q", ds, &out, tt.line, tt.out)
			}
		})
	}
}
