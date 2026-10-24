package interp

import (
	"bytes"
	"testing"

	"github.com/ncode/PortuGo/internal/diag"
)

func TestNestedHaltStopsCallerEvaluation(t *testing.T) {
	for _, trigger := range []struct{ name, body string }{
		{"empty numeric call", "stop()"},
		{"absent choice", "escolha arccos(2)\ncaso 1 ate bound:\nescreval(\"BAD ARM\")\nfimescolha"},
	} {
		for _, consumer := range []struct{ name, body string }{
			{"output", `escreval("prefix", f(), "tail")`},
			{"later output error", `escreval(f(), 1 \ 0)`},
			{"binary operand", `escreval(f() + (1 \ 0))`},
			{"width", `escreval(1:f(), 1 \ 0)`},
			{"decimals", `escreval(1:4:f(), 1 \ 0)`},
			{"builtin argument", `escreval(copia("abc", f(), 1 \ 0))`},
			{"user argument", `take(f(), 1 \ 0)`},
			{"assignment", `count <- f() + (1 \ 0)`},
			{"reference copy back", `escreval(changing(count))`},
		} {
			t.Run(trigger.name+"/"+consumer.name, func(t *testing.T) {
				src := `algoritmo "nested halt"
var count, bound: inteiro
procedimento stop(x: inteiro)
inicio
escreval("missing:", x)
fimprocedimento
procedimento take(a, b: inteiro)
inicio
escreval("BAD CALL")
fimprocedimento
funcao f(): inteiro
inicio
` + trigger.body + `
retorne 7
fimfuncao
funcao changing(var x: inteiro): inteiro
inicio
x <- 9
retorne f()
fimfuncao
inicio
escreval("BEFORE")
` + consumer.body + `
escreval("AFTER")
fimalgoritmo`
				p, info := analyzed(t, src)
				var out bytes.Buffer
				i := New(Options{Output: &out, MaxSteps: 100})
				ds := i.Run(p, info)
				if len(ds) != 0 || out.String() != "BEFORE\n" || i.State()["count"].Int != 0 {
					t.Fatalf("diagnostics=%v output=%q count=%+v", ds, &out, i.State()["count"])
				}
				p, info = analyzed(t, "algoritmo \"reuse\"\ninicio\nescreval(\"OK\")\nfimalgoritmo\n")
				out.Reset()
				if ds := i.Run(p, info); len(ds) != 0 || out.String() != "OK\n" {
					t.Fatalf("halt affected next run: diagnostics=%v output=%q", ds, &out)
				}
			})
		}
	}
}

func TestHaltPreservesEarlierFailure(t *testing.T) {
	// A real error before the empty call must still be reported.
	p, info := analyzed(t, `algoritmo "earlier failure"
procedimento stop(x: inteiro)
inicio
escreval("missing:", x)
fimprocedimento
inicio
escreval(1 \ 0)
stop()
fimalgoritmo`)
	ds := New(Options{}).Run(p, info)
	if len(ds) != 1 || ds[0].Code != diag.RArithmetic {
		t.Fatalf("earlier diagnostic lost: %v", ds)
	}
}
