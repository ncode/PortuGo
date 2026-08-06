package lexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/token"
)

func TestScanGolden(t *testing.T) {
	src := "algoritmo \"x\"\ninicio\nx <- 1 + 2\nfimalgoritmo\n"
	file, toks, diags := Scan("test.alg", src)
	if len(diags) > 0 {
		t.Fatalf("diagnostics: %v", diags)
	}
	var lines []string
	for _, tok := range toks {
		if tok.Kind == token.EOF {
			break
		}
		pos := file.Position(tok.Pos)
		lines = append(lines, fmt.Sprintf("%s %q @%d:%d", tok.Kind, tok.Text, pos.Line, pos.Column))
	}
	got := strings.Join(lines, "\n")
	want := strings.Join([]string{
		"algoritmo \"algoritmo\" @1:1",
		"STRING \"x\" @1:11",
		"inicio \"inicio\" @2:1",
		"IDENT \"x\" @3:1",
		"<- \"<-\" @3:3",
		"NUMBER \"1\" @3:6",
		"+ \"+\" @3:8",
		"NUMBER \"2\" @3:10",
		"fimalgoritmo \"fimalgoritmo\" @4:1",
	}, "\n")
	if got != want {
		t.Fatalf("tokens mismatch\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}

func FuzzLexer(f *testing.F) {
	for _, seed := range []string{
		"",
		"algoritmo \"x\"\ninicio\nfimalgoritmo",
		"vetor[1..10] de inteiro",
		"se verdadeiro e falso entao fimse",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, src string) {
		Scan("fuzz.alg", src)
	})
}
