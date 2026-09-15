package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/interp"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/sema"
	"github.com/ncode/PortuGo/internal/source"
)

func TestRecordedRepeatFimrepitaPrefix(t *testing.T) {
	src := readRepeatProbe(t, "repeat-fimrepita")
	_, tokens, lexDiags := lexer.Scan("source.alg", src)
	if len(lexDiags) != 0 {
		t.Fatal(lexDiags)
	}
	prog, parseDiags := parser.Parse(tokens)
	if len(parseDiags) != 0 {
		t.Fatalf("parse diagnostics = %v", parseDiags)
	}
	info, semaDiags := sema.Analyze(prog)
	if len(semaDiags) != 0 {
		t.Fatalf("sema diagnostics = %v", semaDiags)
	}
	var out bytes.Buffer
	if diags := interp.New(interp.Options{Output: &out}).Run(prog, info); len(diags) != 0 {
		t.Fatalf("runtime diagnostics = %v", diags)
	}
	want, err := os.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", "repeat-fimrepita", "stdout.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out.Bytes(), want) {
		t.Fatalf("stdout = %q, want %q", out.Bytes(), want)
	}
	var formatted bytes.Buffer
	if err := ast.Fprint(&formatted, prog); err != nil {
		t.Fatal(err)
	}
	_, formattedTokens, lexDiags := lexer.Scan("formatted.alg", formatted.String())
	if len(lexDiags) != 0 {
		t.Fatal(lexDiags)
	}
	formattedProg, parseDiags := parser.Parse(formattedTokens)
	if len(parseDiags) != 0 {
		t.Fatalf("formatted parse diagnostics = %v", parseDiags)
	}
	formattedInfo, semaDiags := sema.Analyze(formattedProg)
	if len(semaDiags) != 0 {
		t.Fatalf("formatted sema diagnostics = %v", semaDiags)
	}
	out.Reset()
	if diags := interp.New(interp.Options{Output: &out}).Run(formattedProg, formattedInfo); len(diags) != 0 || !bytes.Equal(out.Bytes(), want) {
		t.Fatalf("formatted runtime: diagnostics = %v, stdout = %q, want %q", diags, out.Bytes(), want)
	}
	var formattedAgain bytes.Buffer
	if err := ast.Fprint(&formattedAgain, formattedProg); err != nil {
		t.Fatal(err)
	}
	if formattedAgain.String() != formatted.String() {
		t.Fatal("repeat formatting is not idempotent")
	}
}

func TestRecordedRepeatRejectedForms(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{id: "repeat-ate-que", code: diag.EUndeclared, line: 5},
		{id: "repeat-ate-que-accented", code: diag.EUndeclared, line: 5},
		{id: "repeat-fimrepita-reached", code: diag.EParse, line: 10},
	} {
		t.Run(tt.id, func(t *testing.T) {
			src := readRepeatProbe(t, tt.id)
			file, tokens, lexDiags := lexer.Scan("source.alg", src)
			if len(lexDiags) != 0 {
				t.Fatal(lexDiags)
			}
			prog, parseDiags := parser.Parse(tokens)
			if len(parseDiags) != 0 {
				t.Fatalf("parse diagnostics = %v", parseDiags)
			}
			info, semaDiags := sema.Analyze(prog)
			if tt.code == diag.EUndeclared {
				if len(semaDiags) != 1 || semaDiags[0].Code != tt.code || file.Position(semaDiags[0].Pos).Line != tt.line {
					t.Fatalf("sema diagnostics = %v, want %s on line %d", semaDiags, tt.code, tt.line)
				}
				return
			}
			if len(semaDiags) != 0 {
				t.Fatalf("sema diagnostics = %v", semaDiags)
			}
			diags := interp.New(interp.Options{}).Run(prog, info)
			if len(diags) != 1 || diags[0].Code != tt.code || file.Position(diags[0].Pos).Line != tt.line {
				t.Fatalf("runtime diagnostics = %v, want %s on line %d", diags, tt.code, tt.line)
			}
		})
	}
}

func readRepeatProbe(t *testing.T, id string) string {
	t.Helper()
	src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", id, "source.alg"))
	if err != nil {
		t.Fatal(err)
	}
	return src
}
