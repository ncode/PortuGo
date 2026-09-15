package portugol_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/golden"
	"github.com/ncode/PortuGo/internal/interp"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/sema"
	"github.com/ncode/PortuGo/internal/source"
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
			if _, err := os.Stat(strings.TrimSuffix(path, ".alg") + ".err"); err == nil {
				t.Skip("diagnostic fixture is checked by TestErrorFixtures")
			} else if !os.IsNotExist(err) {
				t.Fatal(err)
			}
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

func TestErrorFixtures(t *testing.T) {
	files, err := filepath.Glob("testdata/run/*.err")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no error fixtures")
	}
	for _, errPath := range files {
		algPath := strings.TrimSuffix(errPath, ".err") + ".alg"
		name := strings.TrimSuffix(filepath.Base(errPath), ".err")
		t.Run(name, func(t *testing.T) {
			src, err := source.ReadFile(algPath)
			if err != nil {
				t.Fatal(err)
			}
			file, toks, lexDiags := lexer.Scan(algPath, src)
			ds := lexDiags
			var prog *ast.Program
			var info *sema.Info
			if len(ds) == 0 {
				prog, ds = parser.Parse(toks)
			}
			if len(ds) == 0 {
				info, ds = sema.Analyze(prog)
			}
			var out bytes.Buffer
			if len(ds) == 0 {
				ds = interp.New(interp.Options{Output: &out}).Run(prog, info)
			}
			var got strings.Builder
			for _, d := range diag.Ordered(ds) {
				pos := file.Position(d.Pos)
				fmt.Fprintf(&got, "%s@%d:%d\n", d.Code, pos.Line, pos.Column)
			}
			want, err := os.ReadFile(errPath)
			if err != nil {
				t.Fatal(err)
			}
			if got.String() != string(want) {
				t.Fatalf("diagnostics = %q, want %q", got.String(), want)
			}
			outPath := strings.TrimSuffix(errPath, ".err") + ".out"
			wantOut, err := os.ReadFile(outPath)
			if err != nil {
				if os.IsNotExist(err) {
					if out.Len() != 0 {
						t.Fatalf("unexpected partial output %q without .out fixture", out.String())
					}
					return
				}
				t.Fatal(err)
			}
			if out.String() != string(wantOut) {
				t.Fatalf("stdout = %q, want %q", out.String(), wantOut)
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
