package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedArithmeticDomains(t *testing.T) {
	for _, id := range []string{
		"power-zero-zero", "power-negative-exponent", "power-underflow", "power-real-exponent",
		"power-string-left", "power-string-right", "power-logical-left", "power-logical-right",
		"slash-string-left", "slash-string-right", "slash-logical-left", "slash-logical-right",
		"unary-minus-string", "unary-minus-logical", "unary-plus-integer",
	} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			checkFormattingPreservesExecution(t, src, want)
		})
	}
}

func TestRecordedPowerDiagnostics(t *testing.T) {
	for _, id := range []string{"power-negative-fraction", "power-zero-negative", "power-overflow"} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			prog, ds := parser.Parse(tokens)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			info, ds := sema.Analyze(prog)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			var out bytes.Buffer
			ds = interp.New(interp.Options{Output: &out}).Run(prog, info)
			if len(ds) != 1 || ds[0].Code != diag.RArithmetic || file.Position(ds[0].Pos).Line != 3 {
				t.Fatalf("diagnostics = %v, want R002 on line 3", ds)
			}
			if out.Len() != 0 {
				t.Fatalf("output = %q, want empty", &out)
			}
		})
	}
}
