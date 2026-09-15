package sema

import (
	"fmt"
	"testing"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/runtime"
)

// TestDocumentedOperandTypeMatrix keeps the static result rules for every
// scalar category exercised by the language reference in one table. Numeric
// is the analysis-only category returned by caracpnum before runtime input
// chooses its concrete integer or real representation.
func TestDocumentedOperandTypeMatrix(t *testing.T) {
	operands := []struct {
		name string
		expr string
		kind runtime.TypeKind
	}{
		{name: "integer", expr: "1", kind: runtime.IntegerType},
		{name: "real", expr: "1.0", kind: runtime.RealType},
		{name: "numeric", expr: `caracpnum("1")`, kind: runtime.NumericType},
		{name: "string", expr: `"x"`, kind: runtime.StringType},
		{name: "logical", expr: "verdadeiro", kind: runtime.BoolType},
	}

	numeric := operands[:3]
	scalar := operands

	for _, op := range []string{"+", "-", "*"} {
		for _, left := range numeric {
			for _, right := range numeric {
				name := fmt.Sprintf("%s/%s/%s", op, left.name, right.name)
				want := promotedNumericKind(left.kind, right.kind)
				assertDocumentedExpressionType(t, name, left.expr+" "+op+" "+right.expr, want)
			}
		}
	}
	assertDocumentedExpressionType(t, "add/string/string", `"x" + "y"`, runtime.StringType)

	for _, op := range []string{"/", `\`} {
		for _, left := range scalar {
			for _, right := range scalar {
				name := fmt.Sprintf("%s/%s/%s", op, left.name, right.name)
				want := right.kind
				if op == "/" && isDocumentedNumericKind(left.kind) && isDocumentedNumericKind(right.kind) {
					want = runtime.RealType
				}
				assertDocumentedExpressionType(t, name, left.expr+" "+op+" "+right.expr, want)
			}
		}
	}

	for _, op := range []string{"%", "mod"} {
		for _, left := range scalar {
			for _, right := range scalar {
				if left.kind == runtime.StringType || right.kind == runtime.StringType {
					continue
				}
				name := fmt.Sprintf("%s/%s/%s", op, left.name, right.name)
				want := right.kind
				if isDocumentedNumericKind(left.kind) && isDocumentedNumericKind(right.kind) {
					want = runtime.IntegerType
				}
				assertDocumentedExpressionType(t, name, left.expr+" "+op+" "+right.expr, want)
			}
		}
	}

	for _, left := range scalar {
		for _, right := range scalar {
			name := fmt.Sprintf("power/%s/%s", left.name, right.name)
			want := runtime.VoidType
			if isDocumentedNumericKind(left.kind) && isDocumentedNumericKind(right.kind) {
				want = runtime.RealType
			}
			assertDocumentedExpressionType(t, name, left.expr+" ^ "+right.expr, want)
		}
	}

	for _, op := range []string{"=", "<>", "<", ">", "<=", ">="} {
		for _, left := range scalar {
			for _, right := range scalar {
				name := fmt.Sprintf("%s/%s/%s", op, left.name, right.name)
				want := right.kind
				if (isDocumentedNumericKind(left.kind) && isDocumentedNumericKind(right.kind)) || left.kind == right.kind && (left.kind == runtime.StringType || left.kind == runtime.BoolType) {
					want = runtime.BoolType
				}
				assertDocumentedExpressionType(t, name, left.expr+" "+op+" "+right.expr, want)
			}
		}
	}

	for _, left := range numeric {
		for _, right := range numeric {
			name := fmt.Sprintf("e/%s/%s", left.name, right.name)
			assertDocumentedExpressionType(t, name, left.expr+" e "+right.expr, right.kind)
		}
	}
	assertDocumentedExpressionType(t, "e/logical/logical", "verdadeiro e falso", runtime.BoolType)
	for _, op := range []string{"ou", "xou"} {
		assertDocumentedExpressionType(t, op+"/logical/logical", "verdadeiro "+op+" falso", runtime.BoolType)
	}
	assertDocumentedExpressionType(t, "e/logical/comparison", `verdadeiro e ("x" = 7)`, runtime.DynamicType)

	for _, operand := range numeric {
		assertDocumentedExpressionType(t, "plus/"+operand.name, "+"+operand.expr, operand.kind)
		assertDocumentedExpressionType(t, "minus/"+operand.name, "-"+operand.expr, operand.kind)
	}
	for _, operand := range []struct {
		name string
		expr string
	}{
		{name: "string", expr: `"x"`},
		{name: "logical", expr: "verdadeiro"},
	} {
		assertDocumentedExpressionType(t, "minus/"+operand.name, "-"+operand.expr, runtime.VoidType)
	}
	assertDocumentedExpressionType(t, "not/logical", "nao falso", runtime.BoolType)
	assertDocumentedExpressionType(t, "not/comparison", `nao ("x" = 7)`, runtime.BoolType)
}

func assertDocumentedExpressionType(t *testing.T, name, expr string, want runtime.TypeKind) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Helper()
		program := parseInfoProgram(t, "algoritmo \"expression types\"\ninicio\nescreval("+expr+")\nfimalgoritmo")
		info, diagnostics := Analyze(program)
		if diag.HasErrors(diagnostics) {
			t.Fatalf("analysis diagnostics: %v", diagnostics)
		}
		statement, ok := program.Body[0].(*ast.WriteStmt)
		if !ok || len(statement.Args) != 1 {
			t.Fatalf("unexpected body: %#v", program.Body)
		}
		got, ok := info.TypeOf(statement.Args[0].Expr)
		if !ok || got.Kind != want {
			t.Fatalf("expression type = %v (present=%t), want kind %v", got, ok, want)
		}
	})
}

func isDocumentedNumericKind(kind runtime.TypeKind) bool {
	return kind == runtime.IntegerType || kind == runtime.RealType || kind == runtime.NumericType
}

func promotedNumericKind(left, right runtime.TypeKind) runtime.TypeKind {
	if left == runtime.RealType || right == runtime.RealType {
		return runtime.RealType
	}
	if left == runtime.NumericType || right == runtime.NumericType {
		return runtime.NumericType
	}
	return runtime.IntegerType
}
