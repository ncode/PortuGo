package interp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/token"
)

func TestDynamicBuiltinArguments(t *testing.T) {
	for _, target := range []string{"b", "v[1]", "r.value"} {
		for _, operation := range []struct {
			name, body, output string
			color              bool
		}{
			{"character", target + ` <- "x" = 65
escreval(carac(` + target + `), carac(abs(` + target + `)))`, "AA\n", false},
			{"absent character code", target + ` <- 7 = "x"
escreval(carac(abs(` + target + `)))`, " \n", false},
			{"color", target + ` <- 7 = "vermelho"
c <- 7 = "frente"
mudacor(` + target + `, c)`, "", true},
		} {
			t.Run(target+"/"+operation.name, func(t *testing.T) {
				p, info := analyzed(t, `algoritmo "dynamic arguments"
tipo Holder = registro
value: logico
fimregistro
var b, c: logico
v: vetor[1..1] de logico
r: Holder
inicio
`+operation.body+"\nfimalgoritmo\n")
				var out bytes.Buffer
				host := &displayHost{out: &out}
				ds := New(Options{Output: &out, Host: host}).Run(p, info)
				if len(ds) != 0 || out.String() != operation.output {
					t.Fatalf("diagnostics=%v output=%q, want %q", ds, &out, operation.output)
				}
				if operation.color && (len(host.colors) != 1 || host.colors[0] != (DisplayState{Color: Red})) {
					t.Fatalf("display changes=%v, want one red foreground", host.colors)
				}
			})
		}
	}
}

func TestDynamicBuiltinArgumentRejections(t *testing.T) {
	for _, tt := range []struct{ name, value, call string }{
		{"character real", `"x" = 65.5`, "escreval(carac(b))"},
		{"character text", `7 = "A"`, "escreval(carac(b))"},
		{"character logical", `verdadeiro = falso`, "escreval(carac(b))"},
		{"color integer", `"x" = 65`, `mudacor(b,"frente")`},
		{"color target integer", `"x" = 65`, `mudacor("vermelho",b)`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := "algoritmo \"dynamic rejection\"\nvar b: logico\ninicio\nb <- " + tt.value + "\nescreval(\"BEFORE\")\n" + tt.call + "\nfimalgoritmo\n"
			p, info := analyzed(t, src)
			var out bytes.Buffer
			host := &displayHost{out: &out}
			ds := New(Options{Output: &out, Host: host}).Run(p, info)
			at := token.Pos(strings.Index(src, tt.call) + strings.Index(tt.call, "b"))
			if len(ds) != 1 || ds[0].Code != diag.ETypeMismatch || ds[0].Pos != at || out.String() != "BEFORE\n" || len(host.colors) != 0 {
				t.Fatalf("diagnostics=%v output=%q display=%v", ds, &out, host.colors)
			}
		})
	}
}
