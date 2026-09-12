package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedConstants(t *testing.T) {
	for _, id := range []string{
		"constant-var", "constant-empty-var",
		"constant-form-scalar-types", "constant-form-negative", "constant-form-arithmetic",
		"constant-form-dependency", "constant-form-builtin", "constant-form-call",
		"constant-form-local", "constant-form-function-use", "constant-form-loop-bound",
		"constant-form-overflow", "constant-form-semicolon", "constant-form-empty",
		"constant-context-parameter", "constant-context-global-variable", "constant-context-shadow",
		"constant-context-random-domain", "constant-context-random-order",
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

func TestConstantDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"constant-declaration", diag.EParse, 4},
		{"constant-form-forward", diag.EParse, 3},
		{"constant-form-cycle", diag.EParse, 3},
		{"constant-form-unknown", diag.EParse, 3},
		{"constant-form-assignment", diag.EUndeclared, 6},
		{"constant-form-duplicate", diag.ERedeclared, 4},
		{"constant-form-case-duplicate", diag.ERedeclared, 4},
		{"constant-form-variable-collision", diag.ERedeclared, 5},
		{"constant-form-after-variable", diag.EParse, 4},
		{"constant-context-local-no-var", diag.EParse, 5},
		{"constant-context-bound-arithmetic", diag.EParse, 5},
		{"constant-context-real-bound", diag.EParse, 5},
		{"constant-context-text-bound", diag.EParse, 5},
	} {
		t.Run(tt.id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) == 0 {
				prog, parseDiags := parser.Parse(tokens)
				ds = parseDiags
				if len(ds) == 0 {
					_, ds = sema.Analyze(prog)
				}
			}
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics = %v, want %s on line %d", ds, tt.code, tt.line)
			}
		})
	}
}
