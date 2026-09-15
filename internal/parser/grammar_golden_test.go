package parser

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/source"
	"github.com/ncode/PortuGo/internal/token"
)

func TestFrontendGrammarGolden(t *testing.T) {
	const sourcePath = "../../testdata/format/frontend_grammar.alg"
	original, err := source.LoadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	wantTokens, err := os.ReadFile(strings.TrimSuffix(sourcePath, ".alg") + ".tokens")
	if err != nil {
		t.Fatal(err)
	}
	wantCP1252Tokens, err := os.ReadFile(strings.TrimSuffix(sourcePath, ".alg") + ".cp1252.tokens")
	if err != nil {
		t.Fatal(err)
	}
	wantFormatted, err := os.ReadFile(strings.TrimSuffix(sourcePath, ".alg") + ".formatted.alg")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name       string
		ending     string
		file       *source.File
		wantTokens []byte
	}{
		{name: "LF", ending: "\n", file: original, wantTokens: wantTokens},
		{name: "CRLF", ending: "\r\n", file: decodeGrammarSource(t, []byte(strings.ReplaceAll(original.Text, "\n", "\r\n"))), wantTokens: wantTokens},
		{name: "CP1252", ending: "\n", file: decodeGrammarSource(t, cp1252GrammarBytes(original.Text)), wantTokens: wantCP1252Tokens},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			file, tokens, lexDiags := lexer.ScanFile(tc.file)
			if len(lexDiags) != 0 {
				t.Fatalf("lexer diagnostics: %v", lexDiags)
			}
			for _, tok := range tokens {
				if tok.Kind == token.NEWLINE && tok.Text != tc.ending {
					t.Fatalf("newline token = %q, want %q", tok.Text, tc.ending)
				}
			}
			if got := grammarTokenGolden(file, tokens); !bytes.Equal([]byte(got), tc.wantTokens) {
				t.Fatalf("token golden mismatch:\ngot:\n%s\nwant:\n%s", got, tc.wantTokens)
			}
			program, parseDiags := Parse(tokens)
			if len(parseDiags) != 0 {
				t.Fatalf("parser diagnostics: %v", parseDiags)
			}
			var formatted bytes.Buffer
			if err := ast.Fprint(&formatted, program); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(formatted.Bytes(), wantFormatted) {
				t.Fatalf("formatted grammar golden mismatch:\ngot:\n%s\nwant:\n%s", formatted.Bytes(), wantFormatted)
			}

			_, reparsedTokens, lexDiags := lexer.Scan("frontend_grammar.formatted.alg", formatted.String())
			if len(lexDiags) != 0 {
				t.Fatalf("reformatted lexer diagnostics: %v", lexDiags)
			}
			reparsed, parseDiags := Parse(reparsedTokens)
			if len(parseDiags) != 0 {
				t.Fatalf("reformatted parser diagnostics: %v", parseDiags)
			}
			program.Fragments = nil
			reparsed.Fragments = nil
			program.Suffix.Text = strings.ReplaceAll(program.Suffix.Text, "\r", "")
			clearSyntaxPositions(reflect.ValueOf(program))
			clearSyntaxPositions(reflect.ValueOf(reparsed))
			if !reflect.DeepEqual(program, reparsed) {
				t.Fatal("grammar golden changed the syntax tree")
			}
		})
	}
}

func TestFrontendGrammarRecoveryFixture(t *testing.T) {
	const sourcePath = "../../testdata/check/frontend_recovery.alg"
	src, err := source.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		ending string
	}{
		{name: "LF", ending: "\n"},
		{name: "CRLF", ending: "\r\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			text := strings.ReplaceAll(src, "\n", tc.ending)
			file, tokens, lexDiags := lexer.Scan(sourcePath, text)
			if len(lexDiags) != 0 {
				t.Fatalf("lexer diagnostics: %v", lexDiags)
			}
			program, parseDiags := Parse(tokens)
			want := []struct {
				code diag.Code
				line int
			}{
				{code: diag.EParse, line: 3},
				{code: diag.EParse, line: 6},
				{code: diag.EParse, line: 8},
			}
			if len(parseDiags) != len(want) {
				t.Fatalf("diagnostics = %v, want %d independent diagnostics", parseDiags, len(want))
			}
			for i, expected := range want {
				if parseDiags[i].Code != expected.code || file.Position(parseDiags[i].Pos).Line != expected.line {
					t.Fatalf("diagnostic %d = %v, want %s on line %d", i, parseDiags[i], expected.code, expected.line)
				}
			}
			if program == nil || len(program.Body) != 3 {
				t.Fatalf("recovery lost statements: %+v", program)
			}
			for i, line := range []int{6, 7, 9} {
				if got := file.Position(program.Body[i].Start()).Line; got != line {
					t.Fatalf("body statement %d starts on line %d, want %d", i, got, line)
				}
			}
		})
	}
}

func grammarTokenGolden(file *token.File, tokens []token.Token) string {
	var out strings.Builder
	for _, tok := range tokens {
		switch tok.Kind {
		case token.EOF, token.SUFFIX, token.INVALID_SUFFIX:
			return out.String()
		}
		pos := file.Position(tok.Pos)
		text := tok.Text
		if tok.Kind == token.NEWLINE {
			text = "\n"
		}
		fmt.Fprintf(&out, "%s %q @%d:%d\n", tok.Kind, text, pos.Line, pos.Column)
	}
	return out.String()
}

func decodeGrammarSource(t *testing.T, data []byte) *source.File {
	t.Helper()
	file, err := source.DecodeFile("frontend_grammar.alg", data)
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func cp1252GrammarBytes(text string) []byte {
	return bytes.ReplaceAll([]byte(text), []byte("a\xc3\xa7\xc3\xa3o"), []byte{'a', 0xe7, 0xe3, 'o'})
}
