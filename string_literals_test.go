package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedStringLiterals(t *testing.T) {
	for _, id := range []string{
		"backslash-string-newline", "backslash-string-tab", "backslash-string-pair",
		"backslash-string-unknown", "backslash-string-final", "backslash-string-quote",
	} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			if id == "backslash-string-quote" {
				file, _, ds := lexer.Scan("source.alg", src)
				if len(ds) != 1 || ds[0].Code != diag.ELexer || file.Position(ds[0].Pos).Line != 3 {
					t.Fatalf("diagnostics = %v, want L001 on line 3", ds)
				}
				return
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			checkStringFormatting(t, src, want)
		})
	}
}

func TestStringFormattingPreservesProgramName(t *testing.T) {
	src, err := source.ReadFile("testdata/run/literal_backslashes.alg")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/run/literal_backslashes.out")
	if err != nil {
		t.Fatal(err)
	}
	checkStringFormatting(t, src, want)
}

func checkStringFormatting(t *testing.T, src string, want []byte) {
	t.Helper()
	var name string
	for pass := range 2 {
		_, toks, lexDiags := lexer.Scan("source.alg", src)
		prog, parseDiags := parser.Parse(toks)
		if len(lexDiags)+len(parseDiags) != 0 {
			t.Fatalf("diagnostics: %v %v", lexDiags, parseDiags)
		}
		if pass == 0 {
			name = prog.Name
		} else if prog.Name != name {
			t.Fatalf("formatted name = %q, want %q", prog.Name, name)
		}
		info, ds := sema.Analyze(prog)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		var out bytes.Buffer
		if ds := interp.New(interp.Options{Output: &out}).Run(prog, info); len(ds) != 0 {
			t.Fatal(ds)
		}
		if !bytes.Equal(out.Bytes(), want) {
			t.Fatalf("pass %d output = %q, want %q", pass, out.Bytes(), want)
		}
		var formatted bytes.Buffer
		if err := ast.Fprint(&formatted, prog); err != nil {
			t.Fatal(err)
		}
		if pass != 0 && formatted.String() != src {
			t.Fatal("formatting is not idempotent")
		}
		src = formatted.String()
	}
}
