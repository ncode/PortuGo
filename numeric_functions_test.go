package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedNumericFunctions(t *testing.T) {
	for _, id := range []string{
		"numeric-arccos-value", "numeric-arccos-empty", "numeric-arcsen-value", "numeric-arcsen-empty",
		"numeric-arctan-value", "numeric-arctan-empty", "numeric-cotan-value", "numeric-cotan-empty",
		"numeric-grauprad-value", "numeric-grauprad-empty", "numeric-radpgrau-value", "numeric-radpgrau-empty",
		"numeric-quad-value", "numeric-quad-type", "numeric-quad-empty",
		"quad-kinds", "quad-fraction", "quad-string", "quad-logical", "cotan-zero",
		"arccos-outside-domain", "arcsen-outside-domain", "arctan-unit",
		"numeric-function-no-value", "numeric-angle-real",
		"legacy-pi-value",
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

func TestNumericFunctionDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
	}{
		{"numeric-arccos-type", diag.ETypeMismatch},
		{"numeric-arcsen-type", diag.ETypeMismatch},
		{"numeric-arctan-type", diag.ETypeMismatch},
		{"numeric-cotan-type", diag.ETypeMismatch},
		{"numeric-grauprad-type", diag.ETypeMismatch},
		{"numeric-radpgrau-type", diag.ETypeMismatch},
		{"quad-real-type", diag.ETypeMismatch},
		{"cotan-empty-type", diag.ETypeMismatch},
		{"cotan-string", diag.EParse},
		{"numeric-function-string", diag.EParse},
		{"numeric-function-logical", diag.EParse},
		{"numeric-function-extra-argument", diag.EParse},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			checkSemanticDiagnostic(t, path, tt.code, 3)
		})
	}
}

func TestPiParenthesesDiagnostic(t *testing.T) {
	path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", "legacy-pi-empty", "source.alg")
	src, err := source.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	file, tokens, ds := lexer.Scan("source.alg", src)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	_, ds = parser.Parse(tokens)
	if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != 3 {
		t.Fatalf("diagnostics = %v, want P001 on line 3", ds)
	}
}
