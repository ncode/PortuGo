package repl

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/source"
)

func TestSourceLimitRecovery(t *testing.T) {
	prefix, end := "algoritmo \"large\"\ninicio\n//", "\nfimalgoritmo\n"
	good := "algoritmo \"next\"\nvar n: inteiro\ninicio\nleia(n)\nescreval(n)\nfimalgoritmo\n7\n:sair\n"
	oversizedLine := strings.Repeat("x", source.MaxBytes+1) + "\n"
	for _, tt := range []struct{ name, rejected string }{
		{"first line", oversizedLine},
		{"continuation", prefix + oversizedLine + "fimalgoritmo\n"},
		{"terminator crosses cap", prefix + strings.Repeat("x", source.MaxBytes-len(prefix)-len(end)+1) + end},
		{"short line crosses cap", prefix + strings.Repeat("x", source.MaxBytes-len(prefix)-3) + "\n   \nfimalgoritmo\n"},
		{"discarded long lines", oversizedLine + oversizedLine + "fimalgoritmo\n"},
		{"discarded lookalikes", oversizedLine + "escreval(\"algoritmo\")\n// algoritmo \"fake\"\nalgoritmo_extra\nfimalgoritmo\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out, stderr bytes.Buffer
			ok, err := Run(interp.Options{Input: strings.NewReader(tt.rejected + good), Output: &out, MaxSteps: 100}, &stderr)
			if err != nil || ok || strings.Count(stderr.String(), ": E900:") != 1 ||
				strings.Count(stderr.String(), "<repl>:") != 1 ||
				strings.Count(out.String(), " 7\n") != 1 || !strings.HasSuffix(out.String(), " 7\nportugol> ") {
				t.Fatalf("ok=%t, error=%v, output=%q, diagnostics=%q", ok, err, &out, &stderr)
			}
		})
	}
	for _, tail := range []string{"", "\n:sair\n"} {
		var stderr bytes.Buffer
		ok, err := Run(interp.Options{Input: strings.NewReader(strings.TrimSuffix(oversizedLine, "\n") + tail)}, &stderr)
		if err != nil || ok || strings.Count(stderr.String(), ": E900:") != 1 {
			t.Fatalf("recovery exit: ok=%t, error=%v, diagnostics=%q", ok, err, &stderr)
		}
	}
}

func TestHostDiagnosticRecovery(t *testing.T) {
	input := "algoritmo \"host\"\ninicio\nescreva(\"BEFORE\")\nescreval(\"FAIL\")\nfimalgoritmo\n" +
		"algoritmo \"next\"\ninicio\nescreval(7)\nfimalgoritmo\n:sair\n"
	var out failOnceWriter
	var stderr bytes.Buffer
	ok, err := Run(interp.Options{Input: strings.NewReader(input), Output: &out}, &stderr)
	if err != nil || ok || strings.Count(stderr.String(), ": R008:") != 1 || !strings.Contains(stderr.String(), "<repl>:4:") ||
		!strings.Contains(out.output.String(), "BEFOREportugol> ") || !strings.HasSuffix(out.output.String(), " 7\nportugol> ") {
		t.Fatalf("ok=%t, error=%v, output=%q, diagnostics=%q", ok, err, &out.output, &stderr)
	}
}

func TestNilDiagnosticWriter(t *testing.T) {
	for _, rejected := range []string{
		"algoritmo \"bad\"\ninicio\n@\nfimalgoritmo\n",
		strings.Repeat("x", source.MaxBytes+1) + "\n",
	} {
		input := rejected + "algoritmo \"next\"\ninicio\nescreval(7)\nfimalgoritmo\n:sair\n"
		var out bytes.Buffer
		ok, err := Run(interp.Options{Input: strings.NewReader(input), Output: &out, MaxSteps: 100}, nil)
		if err != nil || ok || !strings.HasSuffix(out.String(), " 7\nportugol> ") {
			t.Fatalf("nil diagnostic writer: ok=%t, error=%v, output=%q", ok, err, &out)
		}
	}
}

type failOnceWriter struct {
	output bytes.Buffer
	failed bool
}

func (w *failOnceWriter) Write(p []byte) (int, error) {
	if !w.failed && bytes.Contains(p, []byte("FAIL")) {
		w.failed = true
		return 0, errors.New("synthetic output failure")
	}
	return w.output.Write(p)
}
