package portugol_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedRandomInput(t *testing.T) {
	for _, id := range []string{
		"random-input-fixed-integer", "random-input-fixed-negative",
		"random-input-console-transition", "random-input-logical-console", "random-input-legacy-expression",
		"random-range-fraction-fixed", "random-range-integer-precision", "random-range-negative-precision", "random-range-variable-bound",
		"random-context-call-zero", "random-context-call-bound", "random-context-call-range",
		"random-context-call-argument-effect", "random-context-bare-expression",
		"random-context-enable-tail", "random-context-disable-tail", "random-context-extra-argument",
		"random-value-real-zero", "random-value-real-default", "random-value-fraction-origin",
		"random-value-fraction-integer", "random-value-expression-fixed",
		"random-state-precision-reset", "random-state-subprogram", "random-state-argument-order", "random-state-grouped-bound",
	} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			input, err := os.ReadFile(filepath.Join(dir, "input.txt"))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			checkFormattingPreservesExecutionWithInput(t, src, input, want)
		})
	}
}

func TestRandomInputDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		line int
	}{
		{"random-range-default-real", 5}, {"random-range-missing-bound", 5}, {"random-range-parentheses", 5},
		{"random-context-declaration", 3}, {"random-context-unknown-mode", 3},
		{"random-state-logical-bound", 5}, {"random-state-no-value-bound", 5},
		{"random-state-text-bound", 5}, {"random-state-text-precision", 5},
		{"random-quoted-on", 3}, {"random-quoted-off", 3},
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
			if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics=%v; want P001 on line %d", ds, tt.line)
			}
		})
	}
}
