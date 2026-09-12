package lexer

import (
	"bytes"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/token"
)

func TestScanFileOriginalPositions(t *testing.T) {
	for _, newline := range []string{"\n", "\r\n"} {
		t.Run(newline, func(t *testing.T) {
			original := []byte("\ufeffalgoritmo \"\xe9\"" + newline + "// \x80" + newline + "inicio" + newline + "escreval(\"\xe9\x80\")" + newline + "fimalgoritmo // \xe9" + newline + "\x80")
			decoded, err := source.DecodeFile("positions.alg", original)
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := ScanFile(decoded)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			if file != decoded.Positions {
				t.Fatal("scanner discarded original position table")
			}
			seen := make(map[token.Kind]bool)
			for _, tok := range tokens {
				var want int
				switch tok.Kind {
				case token.ALGORITMO, token.INICIO, token.ESCREVAL, token.FIMALGORITMO:
					want = bytes.Index(original, []byte(tok.Text))
				case token.SUFFIX:
					want = bytes.Index(original, []byte("fimalgoritmo")) + len("fimalgoritmo")
				case token.NEWLINE:
					start := int(tok.Pos)
					if start < 0 || start+len(tok.Text) > len(original) || string(original[start:start+len(tok.Text)]) != tok.Text {
						t.Fatalf("newline does not match original bytes: %+v", tok)
					}
					continue
				case token.EOF:
					want = len(original)
				default:
					continue
				}
				seen[tok.Kind] = true
				if tok.Pos != token.Pos(want) {
					t.Errorf("%s offset = %d, want %d", tok.Kind, tok.Pos, want)
				}
			}
			for _, kind := range []token.Kind{token.ALGORITMO, token.INICIO, token.ESCREVAL, token.FIMALGORITMO, token.SUFFIX, token.EOF} {
				if !seen[kind] {
					t.Errorf("missing %s token", kind)
				}
			}
		})
	}
}
