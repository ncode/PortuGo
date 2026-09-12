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

func TestRecordedVectorPrograms(t *testing.T) {
	for _, id := range []string{
		"vector-500-slots", "vector-501-slots", "vector-zero-bound",
		"vector-positive-bound", "vector-5000-slots", "vector-5001-slots",
		"vector-single-zero",
		"vector-two-dimensional", "vector-zero-elements",
		"vector-missing-index", "vector-omitted-column-write-first",
		"vector-omitted-column-write-second", "vector-omitted-column-read",
		"vector-omitted-column-positive-bound", "vector-omitted-column-zero-bound",
		"vector-omitted-column-after-explicit",
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

func TestRecordedVectorRejections(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"vector-copy", diag.ETypeMismatch, 6},
		{"vector-non-one-bound", diag.EParse, 3},
		{"vector-negative-zero-lower", diag.EParse, 3},
		{"vector-plus-lower", diag.EParse, 3},
		{"vector-negative-upper", diag.EParse, 3},
		{"vector-plus-upper", diag.EParse, 3},
		{"vector-reversed-positive", diag.EParse, 3},
		{"vector-real-lower", diag.EParse, 3},
		{"vector-real-upper", diag.EParse, 3},
		{"vector-expression-lower", diag.EParse, 3},
		{"vector-parenthesized-lower", diag.EParse, 3},
		{"vector-expression-upper", diag.EParse, 3},
		{"vector-self-assignment", diag.ETypeMismatch, 6},
		{"vector-to-scalar-assignment", diag.ETypeMismatch, 6},
		{"vector-three-dimensional", diag.EParse, 3},
		{"vector-extra-index", diag.ETypeMismatch, 6},
		{"vector-second-index-oob", diag.RStorage, 6},
		{"vector-negative-index", diag.RStorage, 6},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			src, err := source.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			var out bytes.Buffer
			if len(ds) == 0 {
				prog, parseDiags := parser.Parse(tokens)
				ds = parseDiags
				if len(ds) == 0 {
					var info *sema.Info
					info, ds = sema.Analyze(prog)
					if len(ds) == 0 {
						ds = interp.New(interp.Options{Output: &out}).Run(prog, info)
					}
				}
			}
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics = %v, want %s on line %d", ds, tt.code, tt.line)
			}
			var want []byte
			if tt.code == diag.RStorage {
				want, err = os.ReadFile(filepath.Join(filepath.Dir(path), "stdout.txt"))
				if err != nil {
					t.Fatal(err)
				}
			}
			if !bytes.Equal(out.Bytes(), want) {
				t.Fatalf("output = %q, want %q", &out, want)
			}
		})
	}
}
