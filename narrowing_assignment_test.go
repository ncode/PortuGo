package portugol_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedRealToIntegerAssignmentDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		line int
	}{
		{"literal-large-assignment", 5},
		{"numeric-source-empty-exponent-type", 5},
		{"numeric-source-signed-exponent-type", 5},
		{"type-alias-duplicate", 8},
	} {
		t.Run(tt.id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			program, ds := parser.Parse(tokens)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			info, ds := sema.Analyze(program)
			if len(ds) != 0 {
				t.Fatalf("semantic diagnostics = %v", ds)
			}
			var out bytes.Buffer
			ds = interp.New(interp.Options{Output: &out}).Run(program, info)
			if len(ds) != 1 || ds[0].Code != diag.RType || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics = %v, want R001 on line %d", ds, tt.line)
			}
			if out.Len() != 0 {
				t.Fatalf("output = %q, want empty", out.String())
			}
		})
	}
}
