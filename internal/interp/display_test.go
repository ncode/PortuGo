package interp

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

type displayHost struct {
	HeadlessHost
	out    *bytes.Buffer
	events []string
	colors []DisplayState
	err    error
}

func (h *displayHost) ClearScreen() error {
	h.events = append(h.events, "clear:"+h.out.String())
	return h.err
}

func (h *displayHost) SetDisplay(state DisplayState) error {
	h.events = append(h.events, "color:"+h.out.String())
	h.colors = append(h.colors, state)
	return h.err
}

func TestDisplayHostOrder(t *testing.T) {
	src := `algoritmo "display"
funcao mark(s: caractere): caractere
inicio
escreval(s)
retorne s
fimfuncao
inicio
escreval("BEFORE")
limpatela
mudacor(mark("Amarelo"),mark("Frente"))
mudacor("azul","fundos")
escreval("AFTER")
fimalgoritmo`
	p, info := analyzed(t, src)
	var out bytes.Buffer
	h := &displayHost{out: &out}
	i := New(Options{Output: &out, Host: h})
	for range 2 {
		out.Reset()
		h.events, h.colors = nil, nil
		ds := i.Run(p, info)
		wantEvents := []string{"clear:BEFORE\n", "color:BEFORE\nAmarelo\nFrente\n", "color:BEFORE\nAmarelo\nFrente\n"}
		wantColors := []DisplayState{{Color: Yellow}, {Color: Blue, Background: true}}
		if len(ds) != 0 || !reflect.DeepEqual(h.events, wantEvents) || !reflect.DeepEqual(h.colors, wantColors) || out.String() != "BEFORE\nAmarelo\nFrente\nAFTER\n" {
			t.Fatalf("display execution: diagnostics %v, events %q, colors %v, output %q", ds, h.events, h.colors, &out)
		}
	}
}

func TestDisplayHostFailures(t *testing.T) {
	for _, command := range []string{"limpatela", `mudacor("amarelo","frente")`} {
		t.Run(command, func(t *testing.T) {
			src := "algoritmo \"failure\"\ninicio\nescreval(\"BEFORE\")\n" + command + "\nescreval(\"AFTER\")\nfimalgoritmo"
			p, info := analyzed(t, src)
			var out bytes.Buffer
			cause := errors.New("private host details")
			h := &displayHost{out: &out, err: cause}
			i := New(Options{Output: &out, Host: h})
			ds := i.Run(p, info)
			if len(ds) != 1 || ds[0].Code != diag.RHost || ds[0].Pos != token.Pos(strings.Index(src, command)) || !errors.Is(ds[0], cause) || strings.Contains(ds[0].Error(), cause.Error()) || out.String() != "BEFORE\n" || len(h.events) != 1 {
				t.Fatalf("host failure: diagnostics %v, events %q, output %q", ds, h.events, &out)
			}
			h.err = nil
			out.Reset()
			if ds = i.Run(p, info); len(ds) != 0 || out.String() != "BEFORE\nAFTER\n" {
				t.Fatalf("failed host polluted next run: %v %q", ds, &out)
			}
		})
	}
}

func TestDisplayKeywordsHaveNoValue(t *testing.T) {
	p, info := analyzed(t, `algoritmo "no value"
inicio
escreval("discarded",limpatela)
escreval(mudacor("amarelo","frente"))
escreval("DONE")
fimalgoritmo`)
	var out bytes.Buffer
	h := &displayHost{out: &out}
	if ds := New(Options{Output: &out, Host: h}).Run(p, info); len(ds) != 0 || len(h.events) != 0 || out.String() != "DONE\n" {
		t.Fatalf("expression executed display command: %v, events %q, output %q", ds, h.events, &out)
	}
}

func TestDisplayPalette(t *testing.T) {
	for _, tt := range []struct {
		name  string
		color Color
	}{
		{"Preto", Black}, {"Azul", Blue}, {"Verde", Green}, {"Vermelho", Red},
		{"Roxo", Purple}, {"AmArElO", Yellow}, {"Branco", White},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p, info := analyzed(t, fmt.Sprintf("algoritmo \"palette\"\ninicio\nmudacor(%q,\"FrEnTe\")\nmudacor(%q,\"FuNdOs\")\nfimalgoritmo", tt.name, tt.name))
			var out bytes.Buffer
			h := &displayHost{out: &out}
			want := []DisplayState{{Color: tt.color}, {Color: tt.color, Background: true}}
			if ds := New(Options{Host: h}).Run(p, info); len(ds) != 0 || !reflect.DeepEqual(h.colors, want) {
				t.Fatalf("palette mapping: diagnostics %v, colors %v, want %v", ds, h.colors, want)
			}
		})
	}
	for _, args := range []string{
		`"inexistente","frente"`, `" amarelo ","frente"`, `"cinza","frente"`,
		`"azulclaro","frente"`, `"ciano","frente"`, `"magenta","frente"`,
		`"amarelo","fundo"`, `"amarelo"," fundos "`, `"amarelo","inexistente"`,
	} {
		p, info := analyzed(t, "algoritmo \"unknown color\"\ninicio\nmudacor("+args+")\nfimalgoritmo")
		var out bytes.Buffer
		h := &displayHost{out: &out}
		if ds := New(Options{Host: h}).Run(p, info); len(ds) != 0 || len(h.events) != 0 {
			t.Fatalf("unknown display option %s: diagnostics %v, events %q", args, ds, h.events)
		}
	}
}
