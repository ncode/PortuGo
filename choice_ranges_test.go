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

func TestRecordedChoiceRanges(t *testing.T) {
	for _, id := range []string{
		"choice-range-inside", "choice-range-lower", "choice-range-upper",
		"choice-range-below", "choice-range-above", "choice-range-descending",
		"choice-range-accented", "choice-range-comma", "choice-range-overlap",
		"choice-range-duplicate", "choice-range-real", "choice-range-mixed",
		"choice-range-text", "choice-range-logical", "choice-range-bound-expression",
		"choice-range-bound-variable", "choice-range-evaluation-order", "choice-range-unmatched-order",
		"choice-range-equal-endpoints", "choice-range-negative-endpoints",
		"choice-range-text-lower", "choice-range-text-upper", "choice-range-text-point",
		"choice-range-text-ordinary", "choice-range-text-case",
		"choice-range-logical-ordinary-false", "choice-range-logical-ordinary-true",
		"choice-range-logical-upper", "choice-range-real-upper", "choice-range-no-value-lower",
		"choice-range-no-value-selector", "choice-range-stop-after-first",
		"choice-range-comma-after-miss", "choice-range-variable-values",
		"choice-match-uppercase-literal", "choice-match-uppercase-lower-label",
		"choice-match-empty-literal", "choice-match-two-characters", "choice-match-quoted-label",
		"choice-match-concatenated-selector", "choice-match-real-outside-lower",
		"choice-match-real-outside-upper", "choice-match-variable-lowercase",
		"choice-match-variable-uppercase", "choice-match-variable-lower-upper-label",
		"choice-match-variable-upper-lower-label", "choice-match-variable-range-lower",
		"choice-match-variable-range-middle", "choice-match-variable-range-upper",
		"choice-match-variable-uppercase-range", "choice-match-two-variables",
		"choice-match-function-selector", "choice-match-no-value-lower-order",
		"choice-numeric-point-integer-low", "choice-numeric-point-integer-high",
		"choice-numeric-point-real-low", "choice-numeric-point-real-high",
		"choice-numeric-point-real-exact", "choice-numeric-point-integer-selector",
		"choice-numeric-point-integer-selector-high", "choice-numeric-range-real-lower",
		"choice-numeric-range-both-real-below", "choice-numeric-range-both-real-above",
		"choice-numeric-range-integer-lower-fraction", "choice-numeric-range-integer-lower-next",
		"choice-numeric-range-integer-lower-equal", "choice-numeric-range-integer-selector",
		"choice-numeric-range-integer-selector-below", "choice-numeric-range-negative-real",
		"choice-numeric-range-negative-integer-low", "choice-numeric-range-zero-upper",
		"choice-numeric-range-integer-upper", "choice-numeric-range-real-equal-endpoints",
		"choice-selector-signed-boundary", "choice-selector-unsigned-boundary",
		"choice-selector-negative-boundary", "choice-selector-large-equal",
		"choice-selector-negative-fraction", "choice-selector-negative-fraction-range",
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

func TestChoiceRangeDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"case-range", diag.EParse, 4},
		{"choice-range-missing-upper", diag.EParse, 4},
		{"choice-range-numeric-text-bounds", diag.ETypeMismatch, 4},
		{"choice-range-text-numeric-bounds", diag.ETypeMismatch, 4},
		{"choice-range-logical-numeric-bounds", diag.ETypeMismatch, 4},
		{"choice-range-numeric-logical-bounds", diag.ETypeMismatch, 4},
		{"choice-range-no-value-upper", diag.EParse, 4},
		{"choice-range-next-line-selected", diag.EParse, 5},
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
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics = %v, want %s on line %d", ds, tt.code, tt.line)
			}
			if out.Len() != 0 {
				t.Fatalf("output = %q, want empty", &out)
			}
		})
	}
}
