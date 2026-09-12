package portugol_test

import (
	"bytes"
	"fmt"
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

func TestRecordedDisplayCommands(t *testing.T) {
	ids := []string{
		"display-clear-bare", "display-clear-call", "display-clear-argument", "display-clear-value",
		"display-color-foreground", "display-color-background", "display-color-variable",
		"display-color-unknown", "display-color-direction", "display-color-extra",
		"bundled-32d43fcfece0", "bundled-5a4e512c0d2c", "bundled-19f3ab9f0764",
		"display-clear-ignored-tail", "display-clear-same-line", "display-clear-side-effect",
		"display-clear-value-items", "display-color-expression", "display-color-extra-ignored",
		"display-color-extra-unknown", "display-color-no-close", "display-color-order",
		"display-context-purple", "display-context-background", "display-context-background-case",
		"display-context-background-spaces", "display-context-unknown-color-after-yellow",
		"display-context-unknown-color-background", "display-context-clear-value-effects",
		"display-context-color-bare-value", "display-context-color-value-unknown", "display-context-clear-value-unknown",
		"console-control", "console-header", "console-repeated", "console-switch",
		"console-off", "console-unknown", "console-number", "console-text", "console-tail", "console-unclosed",
		"console-case", "console-expression", "console-local",
	}
	for n := range 19 {
		ids = append(ids, fmt.Sprintf("display-palette-%d", n))
	}
	for n := range 4 {
		ids = append(ids, fmt.Sprintf("display-palette-background-%d", n))
	}
	for _, id := range ids {
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

func TestDisplayCommandDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
		out  string
	}{
		{"display-clear-variable", diag.EParse, 3, ""},
		{"display-color-bare", diag.EParse, 4, ""},
		{"display-color-missing", diag.EParse, 4, ""},
		{"display-color-numeric", diag.ETypeMismatch, 4, ""},
		{"display-color-bad-first", diag.ETypeMismatch, 8, ""},
		{"display-color-empty", diag.ETypeMismatch, 8, ""},
		{"display-color-no-open", diag.ETypeMismatch, 8, ""},
		{"display-color-variable-name", diag.EParse, 3, ""},
		{"display-context-color-args-next-line", diag.EParse, 3, ""},
		{"display-context-color-second-next-line", diag.EParse, 3, ""},
		{"display-context-color-missing-target", diag.ETypeMismatch, 3, ""},
		{"display-context-color-unknown-invalid", diag.ETypeMismatch, 3, ""},
		{"console-after-var", diag.EParse, 3, ""},
		{"console-body", diag.EParse, 3, ""},
		{"console-body-prefix", diag.EParse, 4, "BEFORE\n"},
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
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line || out.String() != tt.out {
				t.Fatalf("diagnostics = %v, output = %q; want %s on line %d, output %q", ds, &out, tt.code, tt.line, tt.out)
			}
		})
	}
}
