package lexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/testprocess"
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
		"Algoritmo \"á\"\r\ninicio // comment\r\nfimalgoritmo\r\n",
		"{ unterminated\n",
		"\"unterminated\nnext",
		"1..10 1.5 1e+2 1e-",
		"\x00\xff\xc3",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, src string) {
		if len(src) > testprocess.MaxSourceBytes {
			t.Skip("outside 64 KiB fuzz profile")
		}
		_, toks, _ := Scan("fuzz.alg", src)
		if len(toks) == 0 || toks[len(toks)-1].Kind != token.EOF {
			t.Fatal("missing final EOF token")
		}
		for i, tok := range toks {
			if tok.Pos < 0 || int(tok.Pos) > len(src) {
				t.Fatalf("token position %d outside source", tok.Pos)
			}
			if i > 0 && tok.Pos < toks[i-1].Pos {
				t.Fatal("token positions are not monotonic")
			}
		}
	})
}

func TestFuzzAdversarial(t *testing.T) {
	for _, tt := range []struct{ name, src string }{
		{"long comment", "//" + strings.Repeat("x", testprocess.MaxSourceBytes-2)},
		{"unterminated string", "\"" + strings.Repeat("x", testprocess.MaxSourceBytes-1)},
		{"invalid bytes", strings.Repeat("\xff", testprocess.MaxSourceBytes)},
		{"punctuation", strings.Repeat("[]():+-", 8192)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testprocess.Run(t, func() { Scan("adversarial.alg", tt.src) })
		})
	}
}
