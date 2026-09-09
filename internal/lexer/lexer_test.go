package lexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/testprocess"
	"github.com/ncode/portugol-go/internal/token"
)

func TestScanGolden(t *testing.T) {
	src := "algoritmo \"x\"\ninicio\nx <- 1 + 2\ny := 3\nfimalgoritmo\n"
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
		"IDENT \"y\" @4:1",
		"<- \":=\" @4:3",
		"NUMBER \"3\" @4:6",
		"fimalgoritmo \"fimalgoritmo\" @5:1",
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
		"/* line only\nescreval(8 / 2 * 3)\n*/\n",
		"escreval(\"BEFORE//AFTER\")\n",
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

func TestCommentLinePositions(t *testing.T) {
	for _, ending := range []string{"\n", "\r\n"} {
		for _, comment := range []string{"//", "{", "}", "/", "*"} {
			src := "inicio" + ending + "  " + comment + " comentário" + ending + "  escreval(3)"
			file, tokens, ds := Scan("comments.alg", src)
			if len(ds) != 0 {
				t.Fatalf("ending %q, prefix %q: %v", ending, comment, ds)
			}
			if len(tokens) != 6 || tokens[1].Kind != token.ESCREVAL {
				t.Fatalf("ending %q, prefix %q: tokens = %v", ending, comment, tokens)
			}
			if pos := file.Position(tokens[1].Pos); pos.Line != 3 || pos.Column != 3 {
				t.Fatalf("ending %q, prefix %q: position = %v, want 3:3", ending, comment, pos)
			}
		}
	}
}

func TestRecordedKeywordSpellings(t *testing.T) {
	for _, tt := range []struct {
		text string
		kind token.Kind
	}{
		{"caracter", token.CARACTERE},
		{"caracter_extra", token.IDENT},
		{"função", token.FUNCAO},
		{"então", token.ENTAO},
		{"senão", token.SENAO},
		{"faça", token.FACA},
		{"até", token.ATE},
		{"não", token.NAO},
		{"lógico", token.IDENT},
		{"início", token.IDENT},
		{"até_que", token.IDENT},
	} {
		for _, spelling := range []string{tt.text, strings.ToUpper(tt.text)} {
			file, tokens, ds := Scan("keywords.alg", "  "+spelling)
			if len(ds) != 0 || len(tokens) != 2 || tokens[0].Kind != tt.kind || tokens[0].Text != spelling {
				t.Fatalf("%q: tokens = %v, diagnostics = %v", spelling, tokens, ds)
			}
			if pos := file.Position(tokens[0].Pos); pos.Line != 1 || pos.Column != 3 {
				t.Fatalf("%q: position = %v, want 1:3", spelling, pos)
			}
		}
	}
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
