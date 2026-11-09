package portugol_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/interp"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/sema"
	"github.com/ncode/PortuGo/internal/token"
)

func benchmarkPrograms() []string {
	return []string{"small", "arithmetic", "calls", "vectors", "strings", "declarations"}
}

type benchmarkProgram struct {
	source  string
	tokens  []token.Token
	program *ast.Program
	info    *sema.Info
}

func loadBenchmarkProgram(tb testing.TB, name string) benchmarkProgram {
	tb.Helper()
	path := filepath.Join("testdata", "bench", name)
	source, err := os.ReadFile(path + ".alg")
	if err != nil {
		tb.Fatal(err)
	}
	want, err := os.ReadFile(path + ".out")
	if err != nil {
		tb.Fatal(err)
	}
	_, tokens, ds := lexer.Scan(name+".alg", string(source))
	if len(ds) != 0 {
		tb.Fatal(ds)
	}
	program, ds := parser.Parse(tokens)
	if len(ds) != 0 {
		tb.Fatal(ds)
	}
	info, ds := sema.Analyze(program)
	if len(ds) != 0 {
		tb.Fatal(ds)
	}
	var out bytes.Buffer
	i := interp.New(interp.Options{Output: &out, WorkingDir: ".", MaxSteps: 1_000_000})
	if ds := i.Run(program, info); len(ds) != 0 {
		tb.Fatal(ds)
	}
	if !bytes.Equal(out.Bytes(), want) {
		tb.Fatalf("workload output = %q, want %q", out.Bytes(), want)
	}
	return benchmarkProgram{source: string(source), tokens: tokens, program: program, info: info}
}

func TestBenchmarkPrograms(t *testing.T) {
	for _, name := range benchmarkPrograms() {
		t.Run(name, func(t *testing.T) { loadBenchmarkProgram(t, name) })
	}
}

// BenchmarkPortuGo measures in-memory stages and the complete decoded-source pipeline.
func BenchmarkPortuGo(b *testing.B) {
	for _, name := range benchmarkPrograms() {
		b.Run(name, func(b *testing.B) {
			data := loadBenchmarkProgram(b, name)
			filename := name + ".alg"
			for _, stage := range []string{"scan", "parse", "analyze", "run", "format", "pipeline"} {
				b.Run(stage, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						switch stage {
						case "scan":
							_, _, ds := lexer.Scan(filename, data.source)
							if len(ds) != 0 {
								b.Fatal(ds)
							}
						case "parse":
							_, ds := parser.Parse(data.tokens)
							if len(ds) != 0 {
								b.Fatal(ds)
							}
						case "analyze":
							_, ds := sema.Analyze(data.program)
							if len(ds) != 0 {
								b.Fatal(ds)
							}
						case "run":
							i := interp.New(interp.Options{Output: io.Discard, WorkingDir: "."})
							if ds := i.Run(data.program, data.info); len(ds) != 0 {
								b.Fatal(ds)
							}
						case "format":
							if err := ast.Fprint(io.Discard, data.program); err != nil {
								b.Fatal(err)
							}
						case "pipeline":
							_, tokens, ds := lexer.Scan(filename, data.source)
							if len(ds) != 0 {
								b.Fatal(ds)
							}
							program, ds := parser.Parse(tokens)
							if len(ds) != 0 {
								b.Fatal(ds)
							}
							info, ds := sema.Analyze(program)
							if len(ds) != 0 {
								b.Fatal(ds)
							}
							i := interp.New(interp.Options{Output: io.Discard, WorkingDir: "."})
							if ds := i.Run(program, info); len(ds) != 0 {
								b.Fatal(ds)
							}
						}
					}
				})
			}
		})
	}
}
