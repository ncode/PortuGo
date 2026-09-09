package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedFunctionReturns(t *testing.T) {
	for _, id := range []string{
		"function-partial-return", "function-widening-return",
		"return-default-integer", "return-default-real", "return-default-string", "return-default-logical",
		"return-result-persists", "return-continues-body", "return-last-value-wins",
		"return-finite-while", "return-finite-repeat", "return-finite-for",
		"return-switch-continues", "return-loop-break",
		"return-recursive-shared-result", "return-recursive-fallthrough",
		"return-independent-functions", "return-nested-expression",
		"return-reference-tail-copyback", "return-logical-value", "return-string-value",
		"recursive-local-accumulation",
		"return-storage-procedure-local", "return-storage-procedure-parameter",
		"return-storage-nested-calls", "return-storage-frame-size",
		"return-storage-after-recursion", "return-storage-procedure-depth",
		"return-retains-real", "return-retains-caractere", "return-retains-logico",
		"return-nested-argument-default", "return-nested-argument-retained",
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

func TestRecordedReturnRejections(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"procedure-valued-return", diag.EReturn, 5},
		{"main-valued-return", diag.EReturn, 3},
		{"function-narrowing-return", diag.ETypeMismatch, 4},
		{"function-negative-narrowing-return", diag.ETypeMismatch, 4},
		{"function-incompatible-return", diag.ETypeMismatch, 4},
		{"function-name-assignment-result", diag.EUndeclared, 4},
		{"return-integer-valued-real", diag.ETypeMismatch, 6},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			checkSemanticDiagnostic(t, path, tt.code, tt.line)
		})
	}
}
