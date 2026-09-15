package sema

import (
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
)

func TestControlFlowDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		name   string
		source string
		code   diag.Code
		line   int
		column int
	}{
		{
			name: "repeat condition",
			source: `algoritmo "invalid repeat"
inicio
  repita
    escreva("x")
  ate 1
fimalgoritmo`,
			code: diag.ETypeMismatch, line: 5, column: 7,
		},
		{
			name: "loop variable",
			source: `algoritmo "invalid loop variable"
var
  x: real
inicio
  para x de 1 ate 2 faca
  fimpara
fimalgoritmo`,
			code: diag.ETypeMismatch, line: 5, column: 8,
		},
		{
			name: "loop lower bound",
			source: `algoritmo "invalid loop bound"
var
  i: inteiro
inicio
  para i de 1.5 ate 2 faca
  fimpara
fimalgoritmo`,
			code: diag.ETypeMismatch, line: 5, column: 13,
		},
		{
			name: "loop step",
			source: `algoritmo "invalid loop step"
var
  i: inteiro
inicio
  para i de 1 ate 2 passo 1.5 faca
  fimpara
fimalgoritmo`,
			code: diag.ETypeMismatch, line: 5, column: 27,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			file, tokens, lexDiags := lexer.Scan("control-flow.alg", tt.source)
			if len(lexDiags) != 0 {
				t.Fatal(lexDiags)
			}
			program, parseDiags := parser.Parse(tokens)
			if len(parseDiags) != 0 {
				t.Fatal(parseDiags)
			}
			_, diagnostics := Analyze(program)
			if len(diagnostics) != 1 || diagnostics[0].Code != tt.code {
				t.Fatalf("diagnostics = %v, want one %s diagnostic", diagnostics, tt.code)
			}
			pos := file.Position(diagnostics[0].Pos)
			if pos.Line != tt.line || pos.Column != tt.column {
				t.Fatalf("diagnostic position = %d:%d, want %d:%d", pos.Line, pos.Column, tt.line, tt.column)
			}
		})
	}
}
