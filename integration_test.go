package portugol_test

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/golden"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

var updateGolden = flag.Bool("update", false, "update expected fixture bytes")

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
			info, semaDiags := sema.Analyze(prog)
			if len(semaDiags) > 0 {
				t.Fatalf("sema diagnostics: %v", semaDiags)
			}
			in := readOptional(t, strings.TrimSuffix(path, ".alg")+".in")
			var out bytes.Buffer
			if err := interp.New(interp.Options{Input: strings.NewReader(in), Output: &out}).Run(prog, info); err != nil {
				t.Fatal(err)
			}
			if err := golden.Compare(".", strings.TrimSuffix(path, ".alg")+".out", out.Bytes(), *updateGolden); err != nil {
				t.Fatal(err)
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
