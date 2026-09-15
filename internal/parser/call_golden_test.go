package parser

import (
	"bytes"
	"os"
	"testing"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/source"
)

func TestCallFormsGolden(t *testing.T) {
	sourcePath := "../../testdata/format/call_forms.alg"
	wantPath := "../../testdata/format/call_forms.formatted.alg"
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
		t.Fatalf("formatted call golden mismatch:\ngot:\n%s\nwant:\n%s", formatted.Bytes(), want)
	}
	_, reparsedTokens, lexDiags := lexer.Scan("call_forms.alg", formatted.String())
	if len(lexDiags) != 0 {
		t.Fatalf("reformatted lexer diagnostics: %v", lexDiags)
	}
	reparsed, parseDiags := Parse(reparsedTokens)
	if len(parseDiags) != 0 {
		t.Fatalf("reformatted parser diagnostics: %v", parseDiags)
	}
	var reformatted bytes.Buffer
	if err := ast.Fprint(&reformatted, reparsed); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(reformatted.Bytes(), formatted.Bytes()) {
		t.Fatalf("call formatting is not idempotent:\nfirst:\n%s\nsecond:\n%s", formatted.Bytes(), reformatted.Bytes())
	}
}
