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

func TestRecordedEnvironmentCommands(t *testing.T) {
	for _, id := range []string{
		"echo-command-on", "echo-command-off", "echo-command-bare", "echo-command-number",
		"echo-command-unknown", "echo-command-quoted", "echo-command-tail",
		"random-input-echo-console", "random-input-echo-random",
		"chronometer-bare", "chronometer-off", "chronometer-expression",
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
