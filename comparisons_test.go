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

func TestRecordedComparisons(t *testing.T) {
	for _, id := range []string{
		"comparison-numeric-mixed", "comparison-numeric-order", "comparison-numeric-boundary",
		"comparison-logical-order", "comparison-logical-same",
		"comparison-string-case", "comparison-string-accent", "comparison-string-order",
		"comparison-integer-string", "comparison-string-integer",
		"comparison-integer-logical", "comparison-logical-integer",
		"comparison-logical-string", "comparison-string-logical", "comparison-operand-order",
		"comparison-real-string", "comparison-string-real", "comparison-real-logical", "comparison-logical-real",
		"comparison-logical-reverse", "comparison-empty-left", "comparison-empty-right",
		"comparison-byte-order", "comparison-mixed-order",
		"comparison-record-values", "comparison-record-scalar",
		"comparison-domain-empty-left", "comparison-code-empty-left", "comparison-code-empty-right",
		"comparison-logical-assignment", "comparison-logical-condition",
		"comparison-logical-number", "comparison-logical-text", "comparison-mixed-numeric-condition",
		"comparison-numeric-call", "comparison-text-call",
		"comparison-formatted-value", "comparison-chained", "comparison-arithmetic",
		"comparison-stored-numeric-call", "comparison-stored-text-call",
		"comparison-domain-left-order", "comparison-generic-left-order",
		"comparison-reassign-integer", "comparison-stored-global-call", "comparison-stored-branch",
		"comparison-stored-field", "comparison-stored-element", "comparison-stored-record",
		"comparison-negation", "comparison-integer-parameter", "comparison-logical-return",
		"comparison-e-left", "comparison-e-right", "comparison-e-text", "comparison-constant-value",
		"comparison-reference-scalar", "comparison-reference-chain", "comparison-reference-element", "comparison-reference-field",
		"comparison-and-assignment", "comparison-and-condition", "comparison-division-value", "comparison-negative-category",
		"comparison-sign-divided", "comparison-sign-plus-numeric", "comparison-sign-plus-sum",
		"comparison-division-sign-values",
		"logical-precedence-relation-and", "logical-precedence-relation-or", "logical-precedence-or-xor",
		"logical-precedence-and-or", "logical-precedence-relation-xor",
		"logical-binding-add-first", "logical-binding-multiply", "logical-binding-multiply-first",
		"rand-suffix-number-control", "rand-suffix-assignment-line",
		"logical-binding-values-multiply",
		"logical-binding-and-before-multiply", "logical-binding-values-add",
		"logical-reduction-grouped-add", "logical-reduction-different-add", "logical-reduction-subtract",
		"logical-reduction-multiply", "logical-reduction-divide", "logical-reduction-outer-add",
		"logical-reduction-after-add", "logical-reduction-right-add",
		"logical-context-ordinary-add", "logical-context-function", "logical-context-builtin",
		"logical-context-two-inner", "logical-context-computed-left", "logical-context-comparison-left",
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
			checkFormattingPreservesExecutionWithInput(t, src, nil, want)
		})
	}
}

func TestRecordedComparisonRejection(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"comparison-vector-values", diag.EParse, 5},
		{"logical-precedence-relation-chain", diag.EParse, 3},
		{"logical-binding-add", diag.EParse, 3},
		{"logical-binding-and-first", diag.EParse, 3},
		{"logical-context-comparison-right", diag.EParse, 3},
		{"logical-context-both-comparisons", diag.EParse, 3},
		{"comparison-domain-empty-right", diag.EParse, 3},
		{"comparison-record-result", diag.RType, 11},
		{"comparison-record-types", diag.RType, 15},
		{"comparison-scalar-record", diag.RType, 10},
		{"comparison-integer-assignment", diag.RType, 5},
		{"comparison-string-assignment", diag.RType, 5},
		{"comparison-integer-format", diag.EParse, 3},
		{"comparison-stored-condition", diag.ETypeMismatch, 6},
		{"comparison-reassign-logical", diag.RType, 6},
		{"comparison-stored-condition-prefix", diag.ETypeMismatch, 7},
		{"comparison-logical-operators", diag.EParse, 4},
		{"comparison-logical-parameter", diag.ETypeMismatch, 2},
		{"comparison-stored-record-copy", diag.RType, 12},
		{"comparison-stored-record-overwrite", diag.RType, 12},
		{"comparison-stored-record-field", diag.EUndeclared, 12},
		{"comparison-nested-record-slot", diag.EUndeclared, 12},
		{"string-size-format-values", diag.EParse, 10},
		{"comparison-ou-left", diag.EParse, 8},
		{"comparison-ou-right", diag.EParse, 3},
		{"comparison-ou-text", diag.EParse, 3},
		{"comparison-xou-left", diag.EParse, 8},
		{"comparison-xou-right", diag.EParse, 3},
		{"comparison-xou-text", diag.EParse, 3},
		{"comparison-plus-category", diag.EParse, 4},
		{"comparison-division-category", diag.RType, 4},
		{"comparison-integer-return", diag.ETypeMismatch, 4},
		{"comparison-consumed-category", diag.EParse, 3},
	} {
		t.Run(tt.id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			var out bytes.Buffer
			if len(ds) == 0 {
				prog, parsed := parser.Parse(tokens)
				ds = parsed
				if len(ds) == 0 {
					var info *sema.Info
					info, ds = sema.Analyze(prog)
					if len(ds) == 0 {
						ds = interp.New(interp.Options{Output: &out}).Run(prog, info)
					}
				}
			}
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line || !bytes.Equal(out.Bytes(), want) {
				t.Fatalf("diagnostics = %v, output = %q; want %s on line %d and output %q", ds, out.String(), tt.code, tt.line, want)
			}
		})
	}
}
