package portugol_test

import (
	"bytes"
	"errors"
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

func TestRecordedEnvironmentCommands(t *testing.T) {
	for _, id := range []string{
		"echo-command-on", "echo-command-off", "echo-command-bare", "echo-command-number",
		"echo-command-unknown", "echo-command-quoted", "echo-command-tail",
		"random-input-echo-console", "random-input-echo-random",
		"chronometer-bare", "chronometer-off", "chronometer-expression",
		"execution-pause-bare", "execution-pause-call", "execution-pause-tail",
		"execution-debug-false", "execution-debug-true",
		"environment-timer-straight", "environment-timer-for", "environment-timer-while",
		"environment-timer-real", "environment-timer-text", "environment-timer-logical",
		"environment-timer-number-expression", "environment-timer-text-reset", "environment-timer-logical-reset",
		"environment-timer-negative-reset", "environment-timer-fraction-reset", "environment-timer-conditional",
		"environment-timer-procedure", "environment-timer-repeat", "environment-timer-call-expression", "environment-timer-parenthesized",
		"environment-debug-expression-true", "environment-debug-expression-false", "environment-debug-tail",
		"environment-debug-call-expression", "environment-pause-call-expression",
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

func TestRecordedExecutionFailures(t *testing.T) {
	for _, tt := range []struct {
		id     string
		code   diag.Code
		line   int
		output string
	}{
		{"environment-debug-text", diag.ETypeMismatch, 3, ""},
		{"environment-timer-empty-value", diag.EParse, 3, ""},
		{"environment-timer-empty-reset", diag.EParse, 5, "A\n"},
	} {
		t.Run(tt.id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			var out bytes.Buffer
			if len(ds) == 0 {
				p, parseDiags := parser.Parse(tokens)
				ds = parseDiags
				if len(ds) == 0 {
					info, semaDiags := sema.Analyze(p)
					ds = semaDiags
					if len(ds) == 0 {
						ds = interp.New(interp.Options{Output: &out}).Run(p, info)
					}
				}
			}
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line || out.String() != tt.output {
				t.Fatalf("diagnostics=%v stdout=%q; want %s line %d and %q", ds, out.String(), tt.code, tt.line, tt.output)
			}
		})
	}
}

func TestEnvironmentCommandDiagnostics(t *testing.T) {
	for _, id := range []string{
		"echo-command-expression", "echo-command-argument", "echo-command-declaration",
		"chronometer-number", "chronometer-quoted", "chronometer-unknown",
	} {
		t.Run(id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", id, "source.alg"))
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
			if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != 3 {
				t.Fatalf("diagnostics=%v; want P001 on line 3", ds)
			}
		})
	}
}
