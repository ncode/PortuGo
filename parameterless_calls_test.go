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

func TestRecordedParameterlessCalls(t *testing.T) {
	for _, id := range []string{
		"procedure-no-parentheses", "procedure-bare-call",
		"function-no-parentheses", "function-bare-call",
		"bare-function-side-effects", "local-variable-shadows-function",
		"function-name-priority", "function-name-priority-over-parameter",
		"procedure-local-name-value", "procedure-local-name-parenthesized",
		"procedure-local-name-bare", "procedure-parameter-name-value",
		"global-variable-function-name",
		"global-procedure-variable-read", "global-procedure-variable-call",
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

func TestRecordedCallRejections(t *testing.T) {
	for _, tt := range []struct {
		id   string
		line int
	}{
		{"function-used-as-statement", 8},
		{"bare-function-missing-argument", 7},
		{"bare-variable-statement", 5},
		{"variable-call-statement", 5},
		{"procedure-value-expression", 7},
		{"function-name-priority-with-arguments", 11},
		{"procedure-name-priority", 2},
		{"procedure-local-name-assignment", 2},
		{"procedure-parameter-name-assignment", 2},
		{"global-variable-procedure-name", 4},
		{"procedure-vector-name-assignment", 2},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			checkSemanticDiagnostic(t, path, diag.ECall, tt.line)
		})
	}
}

func TestBareFunctionIsNotReferenceStorage(t *testing.T) {
	checkSemanticDiagnostic(t, "testdata/check/bare_function_reference.alg", diag.ECall, 11)
}

func TestRecordedDuplicateCallables(t *testing.T) {
	for _, id := range []string{
		"function-procedure-same-name", "procedure-function-same-name",
		"duplicate-function-name", "duplicate-procedure-name",
	} {
		t.Run(id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id, "source.alg")
			checkSemanticDiagnostic(t, path, diag.ERedeclared, 6)
		})
	}
}

func checkSemanticDiagnostic(t *testing.T, path string, code diag.Code, line int) {
	t.Helper()
	src, err := source.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	file, tokens, ds := lexer.Scan("source.alg", src)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	prog, ds := parser.Parse(tokens)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	_, ds = sema.Analyze(prog)
	if len(ds) != 1 || ds[0].Code != code || file.Position(ds[0].Pos).Line != line {
		t.Fatalf("diagnostics = %v, want %s on line %d", ds, code, line)
	}
}
