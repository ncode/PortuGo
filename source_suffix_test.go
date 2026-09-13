package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedIgnoredSuffixes(t *testing.T) {
	for _, name := range []string{"unterminated-next-line", "invalid-symbols", "unclosed-block", "notes"} {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", "source-suffix-"+name)
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

func TestRecordedMalformedTerminatorSuffix(t *testing.T) {
	dir := "testdata/conformance/visualg-3.0.7/probes/source-suffix-unterminated-same-line"
	original, err := source.ReadFile(filepath.Join(dir, "source.alg"))
	if err != nil {
		t.Fatal(err)
	}
	original = strings.ReplaceAll(original, "\r\n", "\n")
	want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, ending := range []string{"\n", "\r\n", "\r\r\n"} {
		for _, finalNewline := range []bool{false, true} {
			src := strings.TrimRight(original, "\r\n")
			if finalNewline {
				src += "\n"
			}
			src = strings.ReplaceAll(src, "\n", ending)
			for pass := range 2 {
				file, tokens, ds := lexer.Scan("source.alg", src)
				if len(ds) != 0 {
					t.Fatalf("lexical diagnostics: %v", ds)
				}
				prog, ds := parser.Parse(tokens)
				if len(ds) != 0 {
					t.Fatalf("parse diagnostics: %v", ds)
				}
				info, ds := sema.Analyze(prog)
				if len(ds) != 0 {
					t.Fatalf("semantic diagnostics: %v", ds)
				}
				var out bytes.Buffer
				ds = interp.New(interp.Options{Output: &out}).Run(prog, info)
				if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != 4 {
					t.Fatalf("execution diagnostics: %v, want P001 on line 4", ds)
				}
				if !bytes.Equal(out.Bytes(), want) {
					t.Fatalf("output = %q, want %q", out.Bytes(), want)
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
	}
}
