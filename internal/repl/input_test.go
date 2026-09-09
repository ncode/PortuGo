package repl

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/interp"
)

func TestBufferedSourceAtEOF(t *testing.T) {
	t.Parallel()
	program := "algoritmo \"eof\"\ninicio\nescreval(42)\nfimalgoritmo"
	for _, tt := range []struct {
		name, input, output, diagnostic string
		ok                              bool
	}{
		{name: "empty", ok: true},
		{name: "complete without newline", input: program, output: " 42\n", ok: true},
		{name: "complete with newline", input: program + "\n", output: " 42\n", ok: true},
		{name: "incomplete", input: "algoritmo \"unfinished\"\ninicio\n", diagnostic: "P001"},
		{name: "lexical error", input: "algoritmo \"bad\"\ninicio\n@\nfimalgoritmo", diagnostic: "L001"},
		{name: "semantic error", input: "algoritmo \"bad\"\ninicio\nescreval(missing)\nfimalgoritmo", diagnostic: "E002"},
		{name: "runtime error", input: "algoritmo \"bad\"\ninicio\nescreval(42)\nescreval(1 / 0)\nfimalgoritmo", output: " 42\n", diagnostic: "R002"},
		{name: "prior error still runs final program", input: "algoritmo \"bad\"\ninicio\n@\nfimalgoritmo\n\n" + program, output: " 42\n", diagnostic: "L001"},
		{name: "cancel incomplete", input: "algoritmo \"cancel\"\n:sair\n", ok: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var out, stderr bytes.Buffer
			ok, err := Run(interp.Options{Input: strings.NewReader(tt.input), Output: &out, MaxSteps: 100}, &stderr)
			if err != nil || ok != tt.ok {
				t.Fatalf("ok=%t, error=%v, diagnostics=%q", ok, err, &stderr)
			}
			if tt.output != "" && !strings.HasSuffix(out.String(), tt.output+"portugol> ") {
				t.Errorf("output %q lacks final result %q", &out, tt.output)
			}
			if tt.diagnostic == "" {
				if stderr.Len() != 0 {
					t.Errorf("unexpected diagnostics: %s", &stderr)
				}
			} else if strings.Count(stderr.String(), ": "+tt.diagnostic+":") != 1 || !strings.Contains(stderr.String(), "<repl>:") {
				t.Errorf("missing positioned %s: %q", tt.diagnostic, &stderr)
			}
		})
	}
}

func TestSharedInputPreservesFinalProgram(t *testing.T) {
	t.Parallel()
	input := "algoritmo \"read\"\nvar x: inteiro\ninicio\nleia(x)\nescreval(x)\nfimalgoritmo\n42\n" +
		"algoritmo \"next\"\ninicio\nescreval(7)\nfimalgoritmo"
	var out, stderr bytes.Buffer
	ok, err := Run(interp.Options{Input: strings.NewReader(input), Output: &out, MaxSteps: 100}, &stderr)
	if err != nil || !ok || stderr.Len() != 0 {
		t.Fatalf("ok=%t, error=%v, diagnostics=%q", ok, err, &stderr)
	}
	if strings.Count(out.String(), "42\n 42\n") != 1 || !strings.HasSuffix(out.String(), " 7\nportugol> ") {
		t.Fatalf("shared input lost or replayed: %q", &out)
	}
}

func TestSourceEncoding(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct{ name, source string }{
		{"UTF-8", "algoritmo \"text\"\ninicio\nescreval(\"café\")\nfimalgoritmo\n\n"},
		{"UTF-8 BOM", "\xef\xbb\xbfalgoritmo \"text\"\ninicio\nescreval(\"café\")\nfimalgoritmo\n\n"},
		{"Windows-1252", "algoritmo \"text\"\ninicio\nescreval(\"caf\xe9\")\nfimalgoritmo\n\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var out, stderr bytes.Buffer
			ok, err := Run(interp.Options{Input: strings.NewReader(tt.source), Output: &out, MaxSteps: 100}, &stderr)
			if err != nil || !ok || stderr.Len() != 0 {
				t.Fatalf("ok=%t, error=%v, diagnostics=%q", ok, err, &stderr)
			}
			if !strings.Contains(out.String(), "café\n") {
				t.Errorf("source was not decoded: %q", &out)
			}
		})
	}
}
