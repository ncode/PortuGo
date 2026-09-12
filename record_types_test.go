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

func TestRecordedRecords(t *testing.T) {
	for _, id := range []string{
		"record-field", "record-copy", "record-layout-primitives",
		"record-layout-field-case", "record-layout-local",
		"record-layout-scalar-field", "record-layout-duplicate-field",
		"record-boundary-empty", "record-boundary-ending-semicolon",
		"record-boundary-field-group", "record-boundary-field-parameter",
		"record-boundary-header-semicolon", "record-boundary-read-field",
		"record-boundary-self-copy", "record-boundary-vector-elements",
		"record-field-type-alias-unused", "record-field-type-nested-unused",
		"record-field-type-local-shadow", "record-field-type-vector-cell-copy",
		"record-alias-value-alias-value", "record-alias-value-base-value",
		"record-name-duplicate-real-first", "record-name-record-first", "record-name-value-name",
		"record-reference-replace-record", "record-reference-replace-element",
		"record-reference-self-copy", "record-reference-argument-replacement",
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
			input, err := os.ReadFile(filepath.Join(dir, "input.txt"))
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			checkFormattingPreservesExecutionWithInput(t, src, input, want)
		})
	}
}

func TestRecordedRecordRejections(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"record-layout-missing-field", diag.EUndeclared, 9},
		{"record-layout-nested-inline", diag.EParse, 4},
		{"record-layout-nested-named", diag.EUndeclared, 12},
		{"record-layout-vector-field", diag.EParse, 4},
		{"record-layout-parameter", diag.EParse, 8},
		{"record-boundary-alias-name", diag.EUndeclared, 10},
		{"record-boundary-alias-chain", diag.EUndeclared, 11},
		{"record-boundary-field-semicolon", diag.EParse, 4},
		{"record-field-type-nested-value", diag.EUndeclared, 12},
		{"record-field-type-nested-scalar-assignment", diag.EUndeclared, 12},
		{"record-field-type-unknown-field-type", diag.EParse, 4},
		{"record-field-type-result-header", diag.EParse, 7},
		{"record-name-alias-first", diag.EUndeclared, 10},
		{"record-name-keyword-field", diag.EParse, 4},
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
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics = %v, want %s on line %d", ds, tt.code, tt.line)
			}
		})
	}
}
