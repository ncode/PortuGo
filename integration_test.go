package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRunFixtures(t *testing.T) {
	files, err := filepath.Glob("testdata/run/*.alg")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no run fixtures")
	}
	for _, path := range files {
		t.Run(strings.TrimSuffix(filepath.Base(path), ".alg"), func(t *testing.T) {
			src, err := source.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			_, toks, lexDiags := lexer.Scan(path, src)
			if len(lexDiags) > 0 {
				t.Fatalf("lexer diagnostics: %v", lexDiags)
			}
			prog, parseDiags := parser.Parse(toks)
			if len(parseDiags) > 0 {
				t.Fatalf("parser diagnostics: %v", parseDiags)
			}
			if semaDiags := sema.Check(prog); len(semaDiags) > 0 {
				t.Fatalf("sema diagnostics: %v", semaDiags)
			}
			in := readOptional(t, strings.TrimSuffix(path, ".alg")+".in")
			want := readOptional(t, strings.TrimSuffix(path, ".alg")+".out")
			var out bytes.Buffer
			if err := interp.New(strings.NewReader(in), &out).Run(prog); err != nil {
				t.Fatal(err)
			}
			if out.String() != want {
				t.Fatalf("stdout mismatch\nwant: %q\n got: %q", want, out.String())
			}
		})
	}
}

func readOptional(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
