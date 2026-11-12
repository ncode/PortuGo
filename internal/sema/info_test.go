package sema

import (
	"testing"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/runtime"
	"github.com/ncode/PortuGo/internal/token"
)

func TestAnalyzeFacts(t *testing.T) {
	p := parseInfoProgram(t, `algoritmo "facts"
var v: vetor[0..2, 2..3] de inteiro
procedimento p(v: inteiro)
inicio
  escreval(v)
fimprocedimento
inicio
  v[0,2] <- 7
  p(v[0,2])
fimalgoritmo`)
	info, ds := Analyze(p)
	if diag.HasErrors(ds) || !info.ValidFor(p) || info.ValidFor(&ast.Program{}) {
		t.Fatalf("invalid analysis contract: %v", ds)
	}
	decl := p.Globals[0].Names[0]
	b, ok := info.Binding(decl)
	if !ok || b.Slots != 6 || b.Type.Kind != runtime.VectorType {
		t.Fatalf("missing vector layout: %+v", b)
	}
	b.Type.Ranges[0].Low = 99
	b.Type.Elem.Kind = runtime.StringType
	copy, _ := info.Binding(decl)
	if copy.Type.Ranges[0].Low != 0 || copy.Type.Elem.Kind != runtime.IntegerType {
		t.Fatal("caller mutated semantic facts")
	}
	assign := p.Body[0].(*ast.AssignStmt)
	typ, ok := info.TypeOf(assign.Target)
	if !ok || typ.Kind != runtime.IntegerType {
		t.Fatal("missing designator type")
	}
	base := assign.Target.(*ast.IndexExpr).X.(*ast.IdentExpr)
	use, _ := info.Binding(base.Name)
	if use.ID != decl.Pos {
		t.Fatal("global use bound to a different declaration")
	}
	proc := p.Subs[0].(*ast.ProcedureDecl)
	local := proc.Body[0].(*ast.WriteStmt).Args[0].Expr.(*ast.IdentExpr)
	use, _ = info.Binding(local.Name)
	if use.ID != proc.Params[0].Name.Pos {
		t.Fatal("shadowed use did not bind to its parameter")
	}
}

func TestInfoScalarQueries(t *testing.T) {
	p := parseInfoProgram(t, "algoritmo \"facts\"\nvar V: vetor[1..2] de inteiro\ninicio\nescreval(v[1])\nfimalgoritmo\n")
	info, ds := Analyze(p)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	name := p.Globals[0].Names[0]
	ref, ok := info.Reference(name)
	if !ok || ref.ID != name.Pos || ref.Name != "v" || ref.Builtin {
		t.Fatalf("reference = %+v, present=%t", ref, ok)
	}
	ref.Name = "changed"
	ref, _ = info.Reference(name)
	if ref.Name != "v" {
		t.Fatal("reference mutation affected semantic information")
	}
	expr := p.Body[0].(*ast.WriteStmt).Args[0].Expr.(*ast.IndexExpr).X
	if kind, ok := info.KindOf(expr); !ok || kind != runtime.VectorType {
		t.Fatalf("kind = %v, present=%t", kind, ok)
	}
	binding, _ := info.Binding(name)
	if !info.Compatible(expr, binding.Type) {
		t.Fatal("matching vector layout rejected")
	}
	binding.Type.Ranges[0].Low = 0
	if info.Compatible(expr, binding.Type) {
		t.Fatal("incompatible vector layout accepted")
	}
	if _, ok := info.Reference(token.Token{Pos: name.Pos, Text: "other"}); ok {
		t.Fatal("reference ignored token identity")
	}
	var missing *Info
	if _, ok := missing.Reference(name); ok {
		t.Fatal("nil info has a reference")
	}
	if _, ok := missing.KindOf(expr); ok || missing.Compatible(expr, runtime.Type{}) {
		t.Fatal("nil info has type facts")
	}
}

func TestAnalyzeInvalidInfo(t *testing.T) {
	for _, src := range []string{
		"algoritmo \"x\"\ninicio\nescreval(missing)\nfimalgoritmo",
		"algoritmo \"x\"\nvar v: vetor[0..9223372036854775807] de inteiro\ninicio\nfimalgoritmo",
	} {
		p := parseInfoProgram(t, src)
		info, ds := Analyze(p)
		if !diag.HasErrors(ds) || info.ValidFor(p) {
			t.Fatalf("invalid program has executable info: %v", ds)
		}
	}
	info, ds := Analyze(nil)
	if !diag.HasErrors(ds) || info.ValidFor(nil) {
		t.Fatal("nil program has executable info")
	}
}

func TestUncalledSubprogramDefersUndeclaredIdentifier(t *testing.T) {
	p := parseInfoProgram(t, `algoritmo "lazy body"
procedimento P
inicio
  leia(missing)
fimprocedimento
inicio
  escreval("ok")
fimalgoritmo`)
	info, ds := Analyze(p)
	if len(ds) != 0 || !info.ValidFor(p) {
		t.Fatalf("uncalled body rejected: %v", ds)
	}
	missing := p.Subs[0].(*ast.ProcedureDecl).Body[0].(*ast.ReadStmt).Targets[0]
	if d, ok := info.DeferredDiagnostic(missing.Start()); !ok || d.Code != diag.EUndeclared {
		t.Fatalf("missing deferred diagnostic: %+v, present=%t", d, ok)
	}
}

func parseInfoProgram(t *testing.T, src string) *ast.Program {
	t.Helper()
	_, toks, ds := lexer.Scan("sample.alg", src)
	if diag.HasErrors(ds) {
		t.Fatal(ds)
	}
	p, ds := parser.Parse(toks)
	if diag.HasErrors(ds) {
		t.Fatal(ds)
	}
	return p
}
