package parser

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/lexer"
)

func TestFormatterASTForms(t *testing.T) {
	data, err := os.ReadFile("../../testdata/cli/format_forms.alg")
	if err != nil {
		t.Fatal(err)
	}
	parse := func(text string) *ast.Program {
		t.Helper()
		_, tokens, ds := lexer.Scan("forms.alg", text)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		prog, ds := Parse(tokens)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		return prog
	}
	original := parse(string(data))
	var first bytes.Buffer
	if err := ast.Fprint(&first, original); err != nil {
		t.Fatal(err)
	}
	reparsed := parse(first.String())
	var second bytes.Buffer
	if err := ast.Fprint(&second, reparsed); err != nil {
		t.Fatal(err)
	}
	if first.String() != second.String() {
		t.Fatal("formatting is not idempotent")
	}
	for _, program := range []*ast.Program{original, reparsed} {
		program.Fragments = nil // Formatting trivia is checked by idempotence.
		clearSyntaxPositions(reflect.ValueOf(program))
	}
	if !reflect.DeepEqual(original, reparsed) {
		t.Fatal("formatting changed the syntax tree")
	}
}
