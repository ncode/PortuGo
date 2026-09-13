package sema

import (
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func TestExpressionTypeMatrix(t *testing.T) {
	for _, tt := range []struct {
		name string
		expr string
		kind runtime.TypeKind
	}{
		{name: "integer addition", expr: "1 + 2", kind: runtime.IntegerType},
		{name: "numeric promotion", expr: "1 + 2.0", kind: runtime.RealType},
		{name: "numeric promotion reversed", expr: "1.0 + 2", kind: runtime.RealType},
		{name: "string concatenation", expr: `"a" + "b"`, kind: runtime.StringType},
		{name: "integer subtraction", expr: "3 - 1", kind: runtime.IntegerType},
		{name: "real subtraction", expr: "3.0 - 1", kind: runtime.RealType},
		{name: "real multiplication", expr: "3 * 1.5", kind: runtime.RealType},
		{name: "real multiplication reversed", expr: "3.0 * 1", kind: runtime.RealType},
		{name: "real division", expr: "3 / 2", kind: runtime.RealType},
		{name: "real division reversed", expr: "3.0 / 2", kind: runtime.RealType},
		{name: "integer division", expr: `3 \ 2`, kind: runtime.IntegerType},
		{name: "integer division with real", expr: `3.0 \ 2`, kind: runtime.IntegerType},
		{name: "integer division reversed", expr: `3 \ 2.0`, kind: runtime.RealType},
		{name: "integer modulo", expr: "3 mod 2", kind: runtime.IntegerType},
		{name: "integer modulo with real", expr: "3.0 mod 2", kind: runtime.IntegerType},
		{name: "power", expr: "3 ^ 2", kind: runtime.RealType},
		{name: "power with real", expr: "3.0 ^ 2", kind: runtime.RealType},
		{name: "unary integer", expr: "-3", kind: runtime.IntegerType},
		{name: "unary real", expr: "+3.0", kind: runtime.RealType},
		{name: "logical negation", expr: "nao falso", kind: runtime.BoolType},
		{name: "numeric equality", expr: "1 = 1.0", kind: runtime.BoolType},
		{name: "string equality", expr: `"a" = "a"`, kind: runtime.BoolType},
		{name: "logical equality", expr: "verdadeiro = falso", kind: runtime.BoolType},
		{name: "logical and", expr: "verdadeiro e falso", kind: runtime.BoolType},
		{name: "logical xor", expr: "verdadeiro xou falso", kind: runtime.BoolType},
		{name: "logical or", expr: "verdadeiro ou falso", kind: runtime.BoolType},
	} {
		t.Run(tt.name, func(t *testing.T) {
			prog := parseInfoProgram(t, "algoritmo \"types\"\ninicio\nescreval("+tt.expr+")\nfimalgoritmo")
			info, ds := Analyze(prog)
			if diag.HasErrors(ds) {
				t.Fatalf("analysis diagnostics: %v", ds)
			}
			stmt, ok := prog.Body[0].(*ast.WriteStmt)
			if !ok || len(stmt.Args) != 1 {
				t.Fatalf("unexpected body: %#v", prog.Body)
			}
			got, ok := info.TypeOf(stmt.Args[0].Expr)
			if !ok || got.Kind != tt.kind {
				t.Fatalf("expression type = %v (present=%t), want kind %v", got, ok, tt.kind)
			}
		})
	}
}

func TestExpressionTypeDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		name string
		expr string
	}{
		{name: "string subtraction", expr: `"a" - "b"`},
		{name: "logical multiplication", expr: "verdadeiro * falso"},
		{name: "logical negation of text", expr: `nao "x"`},
		{name: "mixed logical or", expr: "verdadeiro ou 1"},
		{name: "mixed numeric addition", expr: "1 + verdadeiro"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			prog := parseInfoProgram(t, "algoritmo \"invalid types\"\ninicio\nescreval("+tt.expr+")\nfimalgoritmo")
			_, ds := Analyze(prog)
			if len(ds) != 1 || ds[0].Code != diag.ETypeMismatch {
				t.Fatalf("diagnostics = %v, want one E001 diagnostic", ds)
			}
		})
	}
}
