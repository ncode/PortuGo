package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedNumericText(t *testing.T) {
	for _, id := range []string{
		"numpcarac-integers", "numpcarac-reals", "numpcarac-rounding", "numpcarac-extremes",
		"numpcarac-signed-zero", "numpcarac-expression", "numpcarac-missing-argument",
		"numpcarac-string-argument", "numpcarac-logical-argument", "numpcarac-empty-followup",
		"numpcarac-large-integer", "numpcarac-logical-stop", "numpcarac-string-stop", "numpcarac-thresholds",
		"numpcarac-empty-width", "numpcarac-nested-newline", "numpcarac-skips-later-items",
		"numpcarac-inner-pending", "numpcarac-nested-empty", "numpcarac-pending-newline",
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

func TestNumericTextDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"numpcarac-extra-argument", diag.EParse, 3},
		{"numpcarac-bare", diag.EParse, 3},
		{"numpcarac-vector-argument", diag.EParse, 5},
		{"numpcarac-void-return", diag.ETypeMismatch, 4},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			checkSemanticDiagnostic(t, path, tt.code, tt.line)
		})
	}
}
