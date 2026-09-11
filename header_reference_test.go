package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedProgramHeaders(t *testing.T) {
	for _, id := range []string{
		"header-extra-word", "header-extra-string", "header-extra-broken-string",
		"declaration-semicolon-scalar", "declaration-semicolon-real",
		"declaration-semicolon-vector", "declaration-semicolon-local",
		"declaration-semicolon-var", "declaration-without-separator",
		"declaration-semicolon-empty-var", "declaration-semicolon-comment",
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

func TestRecordedHeaderRejections(t *testing.T) {
	for _, tt := range []struct {
		id   string
		line int
	}{
		{"header-same-line", 2},
		{"header-name-next-line", 1},
		{"header-unquoted-name", 1},
		{"header-missing-name", 1},
		{"header-omitted", 1},
		{"header-leading-unrelated-word", 1},
		{"alias-caracter-variable", 3},
		{"legacy-type-variable", 3},
		{"division-alias-identifier", 3},
		{"function-ending-variable", 3},
		{"rand-context-declaration", 3},
		{"rand-context-expression-call", 7},
		{"rand-suffix-direct-comparisons", 5},
		{"rand-logical-left-only", 6},
		{"rand-logical-right-only", 6},
		{"rand-logical-or", 6},
		{"rand-logical-xor", 6},
		{"rand-logical-right-number", 6},
		{"rand-logical-right-randi", 6},
		{"declaration-semicolon-repeated", 3},
		{"declaration-semicolon-same-line", 3},
		{"declaration-semicolon-next-line", 4},
		{"declaration-semicolon-var-same-line", 2},
		{"declaration-semicolon-var-repeated", 2},
		{"bundled-ce6fa8f9a2a3", 13}, // estcivil.alg
		{"bundled-b3549ab1faa5", 1},  // Calculo_media2.alg
		{"bundled-50f837d55875", 1},  // Calculo_media2.alg.ALG
	} {
		t.Run(tt.id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) == 0 {
				var prog *ast.Program
				prog, ds = parser.Parse(tokens)
				if len(ds) == 0 {
					_, ds = sema.Analyze(prog)
				}
			}
			if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics = %v, want P001 on line %d", ds, tt.line)
			}
		})
	}
}
