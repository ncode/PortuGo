package portugol_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/ncode/PortuGo/internal/interp"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/sema"
)

func TestRecordedRecursionDepth(t *testing.T) {
	source, err := os.ReadFile("testdata/conformance/visualg-3.0.7/probes/recursive-depth-32/source.alg")
	if err != nil {
		t.Fatal(err)
	}
	_, tokens, lexDiagnostics := lexer.Scan("recursive-depth-32.alg", string(source))
	program, parseDiagnostics := parser.Parse(tokens)
	if len(lexDiagnostics) != 0 || len(parseDiagnostics) != 0 {
		t.Fatalf("frontend diagnostics = lexer %v, parser %v", lexDiagnostics, parseDiagnostics)
	}
	info, diagnostics := sema.Analyze(program)
	if len(diagnostics) != 0 {
		t.Fatalf("semantic diagnostics = %v", diagnostics)
	}
	var output bytes.Buffer
	if diagnostics := interp.New(interp.Options{Output: &output}).Run(program, info); len(diagnostics) != 0 {
		t.Fatalf("runtime diagnostics = %v", diagnostics)
	}
	if got := output.String(); got != "DONE\n" {
		t.Fatalf("stdout = %q, want %q", got, "DONE\n")
	}
}
