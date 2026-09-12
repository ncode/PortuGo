package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedTextFunctions(t *testing.T) {
	for _, id := range []string{
		"text-upper-ascii", "text-lower-ascii", "text-upper-accents", "text-lower-accents",
		"text-length-accents", "text-length-empty", "text-position-first", "text-position-missing",
		"text-position-empty", "text-position-empty-pair", "text-position-accents",
		"text-copy-zero", "text-copy-negative-start", "text-copy-past-end",
		"text-copy-zero-count", "text-copy-negative-count", "text-copy-large-count", "text-copy-real-start",
		"code-ascii-first", "code-ascii-empty", "code-ascii-extended",
		"code-character-extended", "code-character-zero",
		"code-detail-ascii-empty-items", "code-detail-character-256", "code-detail-character-negative",
		"code-detail-character-no-args", "code-detail-character-table", "text-detail-case-table",
		"text-detail-copy-real-count", "text-detail-copy-real-negative", "text-detail-copy-real-small",
		"text-detail-lower-empty", "text-detail-position-empty-haystack", "text-detail-position-too-long",
		"text-detail-upper-empty",
		"code-context-ascii-zero-character", "code-context-character-no-value", "text-context-copy-accents",
		"text-context-copy-count-int32", "text-context-copy-count-wrap", "text-context-copy-second-no-value",
		"text-context-copy-start-int32", "text-context-copy-start-wrap", "text-context-copy-third-no-value",
		"text-order-arguments", "text-order-ascii-empty",
		"text-state-code-zero-call", "text-state-copy-absent-earlier-call",
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

func TestTextFunctionDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
	}{
		{"code-ascii-numeric", diag.ETypeMismatch},
		{"code-character-real", diag.ETypeMismatch},
		{"code-character-wrap", diag.ETypeMismatch},
		{"text-length-empty-call", diag.ETypeMismatch},
		{"text-position-empty-call", diag.ETypeMismatch},
		{"text-upper-numeric", diag.ETypeMismatch},
		{"code-detail-ascii-extra", diag.EParse},
		{"code-detail-ascii-no-args", diag.ETypeMismatch},
		{"code-detail-character-extra", diag.EParse},
		{"code-detail-character-fraction", diag.ETypeMismatch},
		{"code-detail-character-logical", diag.ETypeMismatch},
		{"code-detail-character-real-whole", diag.ETypeMismatch},
		{"code-detail-character-string", diag.ETypeMismatch},
		{"text-detail-copy-extra", diag.EParse},
		{"text-detail-copy-logical-count", diag.ETypeMismatch},
		{"text-detail-copy-no-args", diag.ETypeMismatch},
		{"text-detail-copy-short", diag.EParse},
		{"text-detail-copy-string-start", diag.ETypeMismatch},
		{"text-detail-length-extra", diag.EParse},
		{"text-detail-lower-no-args", diag.ETypeMismatch},
		{"text-detail-position-bad-second", diag.ETypeMismatch},
		{"text-detail-position-extra", diag.EParse},
		{"text-detail-position-short", diag.ETypeMismatch},
		{"text-detail-upper-extra", diag.EParse},
		{"text-detail-upper-no-args", diag.ETypeMismatch},
		{"code-context-ascii-no-value", diag.ETypeMismatch},
		{"text-context-copy-first-no-value", diag.ETypeMismatch},
		{"text-context-length-no-value", diag.ETypeMismatch},
		{"text-context-position-no-value", diag.ETypeMismatch},
		{"text-context-upper-no-value", diag.ETypeMismatch},
	} {
		t.Run(tt.id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			line := 3
			if strings.Contains(tt.id, "-context-") {
				line = 4
			}
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
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != line || out.Len() != 0 {
				t.Fatalf("diagnostics = %v, output = %q; want %s on line %d, empty output", ds, &out, tt.code, line)
			}
		})
	}
}
