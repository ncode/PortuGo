package parser

import (
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/token"
)

func TestRecoveryKeepsEOFPosition(t *testing.T) {
	src, err := source.ReadFile("../../testdata/check/truncated_call.alg")
	if err != nil {
		t.Fatal(err)
	}
	file, tokens, ds := lexer.Scan("truncated_call.alg", src)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	_, ds = Parse(tokens)
	if len(ds) < 2 {
		t.Fatalf("expected expression and delimiter diagnostics, got %v", ds)
	}
	for _, d := range ds {
		if d.Pos != token.Pos(len(src)) {
			t.Fatalf("recovery diagnostic moved to line %d; want EOF line %d", file.Position(d.Pos).Line, file.Position(token.Pos(len(src))).Line)
		}
	}
}

func TestStatementLineRecovery(t *testing.T) {
	file, tokens, ds := lexer.Scan("recovery.alg", "algoritmo \"recovery\"\ninicio\n(1 + 2)\nescreval(7)\n(3 + 4)\nfimalgoritmo\n")
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	prog, ds := Parse(tokens)
	if len(ds) != 2 {
		t.Fatalf("diagnostics = %v, want two independent line errors", ds)
	}
	for i, line := range []int{3, 5} {
		if ds[i].Code != diag.EParse || file.Position(ds[i].Pos).Line != line {
			t.Fatalf("diagnostic %d = %v, want P001 on line %d", i, ds[i], line)
		}
	}
	if prog == nil || len(prog.Body) != 1 {
		t.Fatal("recovery lost the valid statement between malformed lines")
	}
}

func TestReturnValueDoesNotCrossLine(t *testing.T) {
	src, err := source.ReadFile("../../testdata/conformance/visualg-3.0.7/probes/return-value-next-line/source.alg")
	if err != nil {
		t.Fatal(err)
	}
	_, tokens, ds := lexer.Scan("source.alg", src)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	// The following line is also invalid; inspect recovery without fixing the
	// precedence between syntax diagnostics and the earlier missing value here.
	prog, parseDiags := Parse(tokens)
	if len(parseDiags) != 0 {
		t.Fatalf("parse diagnostics = %v, want semantic missing-return diagnostic", parseDiags)
	}
	if prog == nil || len(prog.Subs) != 1 {
		t.Fatal("missing function after recovery")
	}
	fn, ok := prog.Subs[0].(*ast.FunctionDecl)
	if !ok || len(fn.Body) != 1 {
		t.Fatal("missing return after recovery")
	}
	ret, ok := fn.Body[0].(*ast.ReturnStmt)
	if !ok || ret.Value != nil {
		t.Fatal("return consumed an expression from the next line")
	}
}

func TestWriteTailRecovery(t *testing.T) {
	for _, command := range []string{"escreva", "escreval"} {
		for _, newline := range []struct{ name, text string }{{"LF", "\n"}, {"CRLF", "\r\n"}} {
			t.Run(command+"/"+newline.name, func(t *testing.T) {
				src := strings.Join([]string{
					`algoritmo "write recovery"`, "inicio",
					command + `(1) escreval(2)`, `escreval(3)`, `(4 + 5)`, "fimalgoritmo", "",
				}, newline.text)
				file, tokens, ds := lexer.Scan("source.alg", src)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				prog, ds := Parse(tokens)
				if len(ds) != 2 {
					t.Fatalf("diagnostics = %v, want two independent line errors", ds)
				}
				for i, line := range []int{3, 5} {
					if ds[i].Code != diag.EParse || file.Position(ds[i].Pos).Line != line {
						t.Fatalf("diagnostic %d = %v, want P001 on line %d", i, ds[i], line)
					}
				}
				if prog == nil || len(prog.Body) != 2 || file.Position(prog.Body[1].Start()).Line != 4 {
					t.Fatal("recovery retained trailing code or lost the valid next-line statement")
				}
			})
		}
	}
}

func TestWriteTailRejections(t *testing.T) {
	for _, name := range []string{"write_trailing_semicolon", "write_no_newline_trailing_semicolon", "write_same_line_terminator"} {
		t.Run(name, func(t *testing.T) {
			src, err := source.ReadFile("../../testdata/check/" + name + ".alg")
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			_, ds = Parse(tokens)
			if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != 3 {
				t.Fatalf("diagnostics = %v, want P001 on line 3", ds)
			}
		})
	}
}

func TestWriteOperandLineRecovery(t *testing.T) {
	for _, ending := range []string{"\n", "\r\n", "\r\r\n"} {
		for _, command := range []string{"escreva", "escreval"} {
			src := strings.Join([]string{`algoritmo "write recovery"`, "inicio", command + "(1 +", "2)", "escreval(3)", "(4)", "fimalgoritmo", ""}, ending)
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			prog, ds := Parse(tokens)
			if len(ds) != 2 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != 3 || ds[1].Code != diag.EParse || file.Position(ds[1].Pos).Line != 6 {
				t.Fatalf("diagnostics = %v, want P001 on lines 3 and 6", ds)
			}
			if prog == nil || len(prog.Body) != 2 || file.Position(prog.Body[1].Start()).Line != 5 {
				t.Fatal("recovery lost the following write")
			}
		}
	}
}

func TestCommentTruncatedWriteRecovery(t *testing.T) {
	for _, ending := range []string{"\n", "\r\n", "\r\r\n"} {
		src := strings.Join([]string{`algoritmo "comment recovery"`, "inicio",
			`escreval(1 { note } + 2)`, `escreval(3)`, `(4 + 5)`, "fimalgoritmo", ""}, ending)
		file, tokens, ds := lexer.Scan("source.alg", src)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		prog, ds := Parse(tokens)
		if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != 5 {
			t.Fatalf("diagnostics = %v, want the independent error on line 5", ds)
		}
		if len(prog.Body) != 2 {
			t.Fatal("recovery lost the following valid write")
		}
		if first, ok := prog.Body[0].(*ast.WriteStmt); !ok || file.Position(first.Unclosed).Line != 3 {
			t.Fatal("missing deferred closing-delimiter diagnostic on line 3")
		}
		if _, ok := prog.Body[1].(*ast.WriteStmt); !ok {
			t.Fatal("recovery did not preserve the write")
		}
	}
}

func TestRecoveredExpressionChainLimit(t *testing.T) {
	src := "algoritmo \"recovery limit\"\ninicio\nescreval(\"" + strings.Repeat("x", 1<<16) + "\"" + strings.Repeat(" /* 1 */", 512) + ")\nfimalgoritmo\n"
	_, tokens, ds := lexer.Scan("source.alg", src)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	prog, ds := Parse(tokens)
	if prog != nil || len(ds) != 1 || ds[0].Code != diag.EResource {
		t.Fatalf("uncontrolled recovery chain: %v", ds)
	}
}
