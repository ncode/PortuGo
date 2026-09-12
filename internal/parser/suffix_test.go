package parser

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/lexer"
)

func TestIgnoredSuffixRoundTrip(t *testing.T) {
	for _, suffix := range []string{
		"\n\"unterminated\n",
		"\n@#$?\n",
		"\nse verdadeiro entao\nescreval(\"IGNORED\")\n",
		" // trailing note\nPlain notes: ação\n  { preserved\n",
		"\n", "", "   ",
	} {
		for _, ending := range []string{"\n", "\r\n"} {
			src := strings.ReplaceAll("algoritmo \"Suffix\"\ninicio\nescreval(\"BODY\")\nfimalgoritmo"+suffix, "\n", ending)
			want := "algoritmo \"Suffix\"\ninicio\n  escreval(\"BODY\")\nfimalgoritmo" + suffix
			if suffix == "" {
				want += "\n"
			}
			for pass := range 2 {
				file, tokens, ds := lexer.Scan("suffix.alg", src)
				if len(ds) != 0 {
					t.Fatalf("suffix %q, pass %d: %v", suffix, pass, ds)
				}
				if got := file.Position(tokens[len(tokens)-1].Pos).Line; got != strings.Count(src, "\n")+1 {
					t.Fatalf("EOF line = %d; suffix line positions were lost", got)
				}
				prog, ds := Parse(tokens)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				if prog.Suffix.Text != "" && prog.Suffix.Text != src[int(prog.Suffix.Pos):] {
					t.Fatal("suffix text or position no longer matches the decoded source")
				}
				var out bytes.Buffer
				if err := ast.Fprint(&out, prog); err != nil {
					t.Fatal(err)
				}
				if out.String() != want {
					t.Fatalf("suffix %q, pass %d:\ngot %q\nwant %q", suffix, pass, out.String(), want)
				}
				src = out.String()
			}
		}
	}
}

func TestSuffixDoesNotHideEarlierErrors(t *testing.T) {
	for _, body := range []string{"escreval(\"unterminated\n", "escreval(1 + )\n"} {
		_, tokens, lex := lexer.Scan("suffix.alg", "algoritmo \"Invalid\"\ninicio\n"+body+"fimalgoritmo\n\"ignored")
		_, parse := Parse(tokens)
		if len(lex)+len(parse) == 0 {
			t.Fatalf("accepted malformed body %q", body)
		}
	}
	_, _, ds := lexer.Scan("suffix.alg", "algoritmo \"Invalid\"\ninicio\nfimalgoritmo \"unterminated\n")
	if len(ds) != 1 {
		t.Fatalf("terminator-line lexical errors must remain visible: %v", ds)
	}
}
