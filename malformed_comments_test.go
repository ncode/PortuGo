package portugol_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/interp"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/sema"
	"github.com/ncode/PortuGo/internal/source"
)

func TestRecordedMalformedComments(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
	}{
		{"brace-inline-math", ""}, {"c-comment-in-expression", diag.EUndeclared},
		{"brace-comment-in-expression", diag.EParse},
	} {
		t.Run(tt.id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			checkMalformedCommentExecution(t, src, "", tt.code, 3)
		})
	}
}

func TestMalformedCommentExecutionControls(t *testing.T) {
	const tick = "funcao tick: inteiro\ninicio\nescreval(\"TICK\")\nretorne 3\nfimfuncao\n"
	for _, tt := range []struct {
		name, declarations, body, output string
		code                             diag.Code
		line                             int
	}{
		{"adjacent brace", "", `escreval(1{ note })`, "BEFORE\nAFTER\n", "", 0},
		{"adjacent brace tail", "", `escreval(1{ note } + 2)`, "BEFORE\nAFTER\n", "", 0},
		{"spaced brace", "", `escreval(1 { note })`, "BEFORE\n", diag.EParse, 4},
		{"spaced brace tail", "", `escreval(1 { note } + 2)`, "BEFORE\n", diag.EParse, 4},
		{"parenthesized brace", "", `escreval((1){ note })`, "BEFORE\n", diag.EParse, 4},
		{"undeclared operand", "", `escreval(1 /* note */ + 2)`, "BEFORE\n", diag.EUndeclared, 4},
		{"declared operand", "var\nnote: inteiro\n", "note <- 3\nescreval(1 /* note */ + 2)", "BEFORE\nAFTER\n", "", 0},
		{"literal operand", "", `escreval(1 /* 3 */ + 2)`, "BEFORE\nAFTER\n", "", 0},
		{"call operand", tick, `escreval(1 /* tick() */ + 2)`, "BEFORE\nTICK\nAFTER\n", "", 0},
		{"brace after call", tick, `escreval(tick() + 1{ note })`, "BEFORE\nTICK\nAFTER\n", "", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := "algoritmo \"comment control\"\n" + tt.declarations + "inicio\nescreval(\"BEFORE\")\n" + tt.body + "\nescreval(\"AFTER\")\nfimalgoritmo\n"
			for _, ending := range []string{"\n", "\r\n", "\r\r\n"} {
				checkMalformedCommentExecution(t, strings.ReplaceAll(src, "\n", ending), tt.output, tt.code, tt.line)
			}
		})
	}
}

func checkMalformedCommentExecution(t *testing.T, src, want string, code diag.Code, line int) {
	t.Helper()
	for pass := range 2 {
		file, tokens, ds := lexer.Scan("source.alg", src)
		if len(ds) != 0 {
			t.Fatalf("lexical diagnostics: %v", ds)
		}
		prog, ds := parser.Parse(tokens)
		if len(ds) != 0 {
			t.Fatalf("parse diagnostics: %v", ds)
		}
		info, ds := sema.Analyze(prog)
		if len(ds) != 0 {
			t.Fatalf("semantic diagnostics: %v", ds)
		}
		var out bytes.Buffer
		ds = interp.New(interp.Options{Output: &out}).Run(prog, info)
		if code == "" && len(ds) != 0 || code != "" && (len(ds) != 1 || ds[0].Code != code || file.Position(ds[0].Pos).Line != line) {
			t.Fatalf("execution diagnostics = %v, want %s on line %d", ds, code, line)
		}
		if out.String() != want {
			t.Fatalf("output = %q, want %q", out.String(), want)
		}
		var printed bytes.Buffer
		if err := ast.Fprint(&printed, prog); err != nil {
			t.Fatal(err)
		}
		if pass != 0 && printed.String() != src {
			t.Fatal("formatting is not idempotent")
		}
		src = printed.String()
	}
}

func TestMalformedCommentOriginalBytes(t *testing.T) {
	original := []byte("\ufeffalgoritmo \"\xe9\"\r\ninicio\r\nescreval(1{ \xe9 })\r\nescreval(7)\r\nfimalgoritmo\r\n")
	decoded, err := source.DecodeFile("source.alg", original)
	if err != nil {
		t.Fatal(err)
	}
	file, tokens, ds := lexer.ScanFile(decoded)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	prog, ds := parser.Parse(tokens)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	stmt := prog.Body[0].(*ast.WriteStmt)
	if stmt.Unclosed == 0 || int(stmt.Unclosed) != bytes.IndexByte(original, '{') || file.Position(stmt.Unclosed).Line != 3 {
		t.Fatal("recovered brace lost its original-byte position")
	}
	var printed bytes.Buffer
	if err := ast.Fprint(&printed, prog); err != nil {
		t.Fatal(err)
	}
	if strings.Count(printed.String(), "{ é })") != 1 {
		t.Fatal("comment contents or adjacency changed")
	}
	checkMalformedCommentExecution(t, printed.String(), " 7\n", "", 0)
}

func TestRecoveredCallDefersTextDiagnostic(t *testing.T) {
	const src = `algoritmo "recovered call"

funcao tick: inteiro
inicio
escreval("TICK")
retorne 3
fimfuncao
inicio
escreval("BEFORE")
escreval(copia(1 /* tick() */, falso, 2))
escreval("AFTER")
fimalgoritmo
`
	checkMalformedCommentExecution(t, src, "BEFORE\nTICK\n", diag.ETypeMismatch, 10)
}
