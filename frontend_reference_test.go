package portugol_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/sema"
	"github.com/ncode/PortuGo/internal/source"
	"github.com/ncode/PortuGo/internal/token"
)

func TestRecordedFrontendPrograms(t *testing.T) {
	for _, id := range []string{
		"empty-var-section", "post-terminator-words", "post-terminator-statement",
		"leading-underscore", "internal-underscore", "mixed-case-identifier",
		"numeric-exponent", "trailing-decimal-point",
		"write-bare-newline", "write-bare-no-newline", "choice-faca", "choice-faca-accented",
		"write-bare-comment", "choice-faca-next-line",
		"physical-newline-lf", "physical-newline-crlf",
		"output-tail-comment-next-write",
	} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			checkFormattingPreservesExecution(t, src, want)
		})
	}
}

func TestRecordedFrontendRejections(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"leading-decimal-point", diag.ELexer, 3},
		{"unterminated-string", diag.ELexer, 3},
		{"single-quoted-string", diag.ELexer, 3},
		{"expression-next-line", diag.EParse, 3},
		{"single-slash-inline", diag.EParse, 3},
		{"single-star-inline", diag.EParse, 3},
		{"logical-type-accented", diag.EParse, 3},
		{"accented-identifier", diag.ELexer, 3},
		{"accented-keywords", diag.ELexer, 2},
		{"write-no-parenthesis-string", diag.EParse, 3},
		{"write-no-parenthesis-number", diag.EParse, 3},
		{"write-bare-same-line", diag.EParse, 3},
		{"same-line-statements", diag.EParse, 3},
		{"semicolon-statements", diag.EParse, 3},
		{"write-parentheses-next-line", diag.EParse, 4},
		{"header-next-line-lf", diag.EParse, 1},
		{"header-next-line-crlf", diag.EParse, 1},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			src, err := source.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) == 0 {
				var prog *ast.Program
				prog, ds = parser.Parse(tokens)
				if len(ds) == 0 {
					_, ds = sema.Analyze(prog)
				}
			}
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics = %v, want %s on line %d", ds, tt.code, tt.line)
			}
		})
	}
}

func TestRecordedLineEndings(t *testing.T) {
	var want string
	for _, id := range []string{"physical-newline-lf", "physical-newline-crlf"} {
		src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", id, "source.alg"))
		if err != nil {
			t.Fatal(err)
		}
		file, tokens, ds := lexer.Scan("source.alg", src)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		var got strings.Builder
		for _, tok := range tokens {
			text := tok.Text
			if tok.Kind == token.NEWLINE {
				text = "\n"
			}
			if tok.Kind == token.SUFFIX {
				text = strings.ReplaceAll(text, "\r\n", "\n")
			}
			pos := file.Position(tok.Pos)
			fmt.Fprintf(&got, "%s %q @%d:%d\n", tok.Kind, text, pos.Line, pos.Column)
		}
		prog, ds := parser.Parse(tokens)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		if err := ast.Fprint(&got, prog); err != nil {
			t.Fatal(err)
		}
		if want != "" && got.String() != want {
			t.Fatalf("%s changed positioned token kinds or syntax", id)
		}
		want = got.String()
	}
}

func TestAccentedFormPositions(t *testing.T) {
	for _, tt := range []struct {
		id, spelling string
		line         int
	}{
		{"accented-identifier", "ação", 3},
		{"accented-keywords", "início", 2},
	} {
		t.Run(tt.id, func(t *testing.T) {
			original, err := os.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := source.DecodeFile("source.alg", original)
			if err != nil {
				t.Fatal(err)
			}
			for _, variant := range []struct {
				data     []byte
				spelling string
			}{
				{original, tt.spelling},
				{[]byte(decoded.Text), tt.spelling},
				{[]byte("\ufeff" + strings.ToUpper(decoded.Text)), strings.ToUpper(tt.spelling)},
			} {
				file, err := source.DecodeFile("source.alg", variant.data)
				if err != nil {
					t.Fatal(err)
				}
				positions, tokens, ds := lexer.ScanFile(file)
				if len(ds) != 0 {
					t.Fatalf("tokenization changed: %v", ds)
				}
				_, ds = parser.Parse(tokens)
				if len(ds) != 1 || ds[0].Code != diag.ELexer || positions.Position(ds[0].Pos).Line != tt.line {
					t.Fatalf("diagnostics = %v, want L001 on line %d", ds, tt.line)
				}
				var found bool
				for _, tok := range tokens {
					if tok.Pos == ds[0].Pos && tok.Text == variant.spelling {
						found = true
					}
				}
				wantPos := 0
				for _, line := range bytes.SplitAfter(variant.data, []byte("\n"))[:tt.line-1] {
					wantPos += len(line)
				}
				if !found || int(ds[0].Pos) != wantPos {
					t.Fatal("diagnostic lost the original token spelling or byte position")
				}
			}
		})
	}
}
