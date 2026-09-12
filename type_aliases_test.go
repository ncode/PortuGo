package portugol_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedTypeAliases(t *testing.T) {
	for _, id := range []string{
		"named-type", "type-alias-primitives", "type-alias-chain",
		"type-alias-local", "type-alias-vector-element",
		"type-context-case-fold", "type-context-variable-name",
		"type-context-scalar-compatibility", "type-context-local-shadow",
		"type-context-constant-first", "type-boundary-empty",
		"type-boundary-duplicate-first-real", "type-boundary-header-semicolon",
		"type-callable-function-local", "type-callable-scalar-arguments",
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

func TestLongTypeAliasChain(t *testing.T) {
	var src strings.Builder
	src.WriteString("algoritmo \"alias chain\"\ntipo\nT0 = inteiro\n")
	for n := 1; n <= 1024; n++ {
		fmt.Fprintf(&src, "T%d = T%d\n", n, n-1)
	}
	src.WriteString("var\nx: T1024\ninicio\nx <- 7\nescreval(x)\nfimalgoritmo\n")
	checkFormattingPreservesExecution(t, src.String(), []byte(" 7\n"))
}

func TestRecordedTypeAliasRejections(t *testing.T) {
	for _, tt := range []struct {
		id   string
		line int
	}{
		{"type-alias-forward", 3}, {"type-alias-cycle", 3},
		{"type-alias-missing", 3}, {"type-alias-vector", 3},
		{"type-alias-vector-chain", 3}, {"type-alias-parameter", 5},
		{"type-context-keyword-name", 3},
		{"type-context-procedure-parameter", 5},
		{"type-context-function-parameter", 5},
		{"type-context-reference-parameter", 6},
		{"type-context-type-first", 4}, {"type-context-after-vars", 4},
		{"type-boundary-no-var", 4}, {"type-boundary-duplicate-invalid", 4},
		{"type-boundary-declaration-semicolon", 3},
		{"type-boundary-unknown-use", 3}, {"type-boundary-local-sibling", 12},
		{"type-callable-result", 5}, {"type-callable-result-bare", 5},
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
