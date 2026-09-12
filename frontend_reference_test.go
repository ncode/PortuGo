package portugol_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/token"
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
		{"single-slash-inline", diag.EParse, 3},
		{"single-star-inline", diag.EParse, 3},
		{"logical-type-accented", diag.EParse, 3},
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
