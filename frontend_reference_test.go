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

func TestRecordedFrontendPrograms(t *testing.T) {
	for _, id := range []string{
		"empty-var-section", "post-terminator-words", "post-terminator-statement",
		"leading-underscore", "internal-underscore", "mixed-case-identifier",
		"numeric-exponent", "trailing-decimal-point",
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

func TestRecordedFrontendRejections(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
	}{
		{"leading-decimal-point", diag.ELexer},
		{"unterminated-string", diag.ELexer},
		{"single-slash-inline", diag.EParse},
		{"single-star-inline", diag.EParse},
		{"logical-type-accented", diag.EParse},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			src, err := source.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) == 0 {
				_, ds = parser.Parse(tokens)
			}
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != 3 {
				t.Fatalf("diagnostics = %v, want %s on line 3", ds, tt.code)
			}
		})
	}
}
