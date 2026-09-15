package portugol_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/ncode/PortuGo/internal/interp"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/sema"
	"github.com/ncode/PortuGo/internal/source"
)

func TestRecordedZeroValue(t *testing.T) {
	const dir = "testdata/conformance/visualg-3.0.7/probes/zero-value"
	src, err := source.ReadFile(dir + "/source.alg")
	if err != nil {
		t.Fatal(err)
	}
	_, tokens, diagnostics := lexer.Scan("zero-value.alg", src)
	if len(diagnostics) != 0 {
		t.Fatalf("lexer diagnostics = %v", diagnostics)
	}
	program, diagnostics := parser.Parse(tokens)
	if len(diagnostics) != 0 {
		t.Fatalf("parser diagnostics = %v", diagnostics)
	}
	info, diagnostics := sema.Analyze(program)
	if len(diagnostics) != 0 {
		t.Fatalf("semantic diagnostics = %v", diagnostics)
	}
	var output bytes.Buffer
	if diagnostics := interp.New(interp.Options{Output: &output}).Run(program, info); len(diagnostics) != 0 {
		t.Fatalf("runtime diagnostics = %v", diagnostics)
	}
	want, err := os.ReadFile(dir + "/stdout.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(output.Bytes(), want) {
		t.Fatalf("stdout = %q, want %q", output.Bytes(), want)
	}
}
