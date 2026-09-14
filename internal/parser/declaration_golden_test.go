package parser

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/source"
)

func TestAggregateDeclarationGolden(t *testing.T) {
	sourcePath := "../../testdata/format/aggregate_declarations.alg"
	wantPath := "../../testdata/format/aggregate_declarations.formatted.alg"
	original, err := source.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatal(err)
	}
	_, tokens, lexDiags := lexer.Scan(sourcePath, original)
	if len(lexDiags) != 0 {
		t.Fatalf("lexer diagnostics: %v", lexDiags)
	}
	program, parseDiags := Parse(tokens)
	if len(parseDiags) != 0 {
		t.Fatalf("parser diagnostics: %v", parseDiags)
	}
	var formatted bytes.Buffer
	if err := ast.Fprint(&formatted, program); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(formatted.Bytes(), want) {
		t.Fatalf("formatted declaration golden mismatch:\ngot:\n%s\nwant:\n%s", formatted.Bytes(), want)
	}
	_, reparsedTokens, lexDiags := lexer.Scan("aggregate_declarations.alg", formatted.String())
	if len(lexDiags) != 0 {
		t.Fatalf("reformatted lexer diagnostics: %v", lexDiags)
	}
	reparsed, parseDiags := Parse(reparsedTokens)
	if len(parseDiags) != 0 {
		t.Fatalf("reformatted parser diagnostics: %v", parseDiags)
	}
	program.Fragments = nil
	reparsed.Fragments = nil
	clearSyntaxPositions(reflect.ValueOf(program))
	clearSyntaxPositions(reflect.ValueOf(reparsed))
	if !reflect.DeepEqual(program, reparsed) {
		t.Fatal("aggregate declaration golden changed the syntax tree")
	}
}

func TestGlobalVarAfterTypedSubprograms(t *testing.T) {
	src := `algoritmo "global var after type subprogram"
tipo
  tdado = registro
    codigo: inteiro
  fimregistro
procedimento incrementar
inicio
  dados.codigo <- dados.codigo + 1
fimprocedimento
var
  dados: tdado
inicio
  dados.codigo <- 4
  incrementar()
  escreval(dados.codigo)
fimalgoritmo
`
	_, tokens, lexDiags := lexer.Scan("trailing-var.alg", src)
	if len(lexDiags) != 0 {
		t.Fatalf("lexer diagnostics: %v", lexDiags)
	}
	program, parseDiags := Parse(tokens)
	if len(parseDiags) != 0 {
		t.Fatalf("parser diagnostics: %v", parseDiags)
	}
	if len(program.Types) != 1 || len(program.Subs) != 1 || len(program.Globals) != 1 {
		t.Fatalf("unexpected declaration counts: types=%d subs=%d globals=%d", len(program.Types), len(program.Subs), len(program.Globals))
	}
	var formatted bytes.Buffer
	if err := ast.Fprint(&formatted, program); err != nil {
		t.Fatal(err)
	}
	_, reparsedTokens, lexDiags := lexer.Scan("trailing-var.alg", formatted.String())
	if len(lexDiags) != 0 {
		t.Fatalf("reformatted lexer diagnostics: %v", lexDiags)
	}
	if _, parseDiags := Parse(reparsedTokens); len(parseDiags) != 0 {
		t.Fatalf("reformatted parser diagnostics: %v", parseDiags)
	}
}
