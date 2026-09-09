package repl

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/interp"
)

func TestAutomaticSubmission(t *testing.T) {
	for _, newline := range []string{"\n", "\r\n"} {
		input := `algoritmo "read"
var n: inteiro
inicio
  leia(n)
  escreval(n)
fimalgoritmo
42
algoritmo "next"
inicio
  escreval(7)
FIMALGORITMO // completed
:sair
`
		input = strings.ReplaceAll(input, "\n", newline)
		var out, stderr bytes.Buffer
		ok, err := Run(interp.Options{Input: strings.NewReader(input), Output: &out, MaxSteps: 100}, &stderr)
		if err != nil || !ok || stderr.Len() != 0 {
			t.Fatalf("ok=%t, error=%v, diagnostics=%q", ok, err, &stderr)
		}
		if strings.Count(out.String(), " 42\n") != 1 || strings.Count(out.String(), " 7\n") != 1 ||
			strings.Index(out.String(), " 42\n") >= strings.Index(out.String(), " 7\n") ||
			!strings.HasSuffix(out.String(), " 7\nportugol> ") {
			t.Fatalf("program or input submission order: %q", &out)
		}
	}
}

func TestIncompleteBlankLinesAndTerminatorLookalikes(t *testing.T) {
	input := `algoritmo "lookalikes"

var fimalgoritmo_extra: inteiro

inicio

  // fimalgoritmo
  { fimalgoritmo
  * fimalgoritmo
  escreval("fimalgoritmo")
  fimalgoritmo_extra <- 7

  escreval(fimalgoritmo_extra)
fimalgoritmo
:sair
`
	var out, stderr bytes.Buffer
	ok, err := Run(interp.Options{Input: strings.NewReader(input), Output: &out, MaxSteps: 100}, &stderr)
	if err != nil || !ok || stderr.Len() != 0 || strings.Count(out.String(), "fimalgoritmo\n 7\nportugol> ") != 1 {
		t.Fatalf("ok=%t, error=%v, output=%q, diagnostics=%q", ok, err, &out, &stderr)
	}
}

func TestAutomaticSubmissionDiagnostics(t *testing.T) {
	for _, tt := range []struct{ name, body, diagnostic string }{
		{"lexical", "@", "L001"},
		{"syntax", "+", "P001"},
		{"semantic", "escreval(missing)", "E002"},
		{"runtime", "escreval(1 / 0)", "R002"},
		{"budget", "enquanto verdadeiro faca\nfimenquanto", "R006"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := "algoritmo \"bad\"\ninicio\n\n" + tt.body + "\nfimalgoritmo\n" +
				"algoritmo \"next\"\ninicio\nescreval(7)\nfimalgoritmo\n:sair\n"
			var out, stderr bytes.Buffer
			ok, err := Run(interp.Options{Input: strings.NewReader(input), Output: &out, MaxSteps: 20}, &stderr)
			if err != nil || ok || strings.Count(stderr.String(), ": "+tt.diagnostic+":") != 1 ||
				!strings.Contains(stderr.String(), "<repl>:4:") ||
				!strings.HasSuffix(out.String(), " 7\nportugol> ") {
				t.Fatalf("ok=%t, error=%v, output=%q, diagnostics=%q", ok, err, &out, &stderr)
			}
		})
	}
}

func TestCP1252TerminatorLookalike(t *testing.T) {
	input := "algoritmo \"cancel\"\nvar \xe9fimalgoritmo: inteiro\n:sair\n"
	var stderr bytes.Buffer
	ok, err := Run(interp.Options{Input: strings.NewReader(input)}, &stderr)
	if err != nil || !ok || stderr.Len() != 0 {
		t.Fatalf("lookalike submitted an unfinished program: ok=%t, error=%v, diagnostics=%q", ok, err, &stderr)
	}
}
