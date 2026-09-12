package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedDeterministicFollowups(t *testing.T) {
	for _, id := range []string{"post-terminator-broken-string", "reference-widening-followup"} {
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

func TestRecordedSyntaxRejectionsBeforeOutput(t *testing.T) {
	for _, tt := range []struct {
		id   string
		line int
	}{
		{"pascal-comment", 3}, {"doubled-string-quote", 3}, {"assignment-equals", 5},
		{"c-comment-after-statement", 3}, {"c-inline-without-statement", 3},
	} {
		t.Run(tt.id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id)
			partial, err := os.ReadFile(filepath.Join(dir, "partial-output.txt"))
			if err != nil {
				t.Fatal(err)
			}
			// The recorded panel contains only the reference's execution wrapper.
			if string(partial) != "Início da execução\r\n\r\n" {
				t.Fatal("reference partial output needs renewed qualification")
			}
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			_, ds = parser.Parse(tokens)
			if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics = %v, want P001 on line %d before execution", ds, tt.line)
			}
		})
	}
}

type zeroElapsedHost struct {
	interp.HeadlessHost
	clockCalls int
}

func (h *zeroElapsedHost) Now() time.Time {
	h.clockCalls++
	return time.Unix(0, 0)
}

func TestRecordedZeroElapsedChronometers(t *testing.T) {
	for _, tt := range []struct {
		id    string
		calls int
	}{
		{"chronometer-on-off", 2}, {"chronometer-repeat-start", 3}, {"chronometer-tail", 2},
	} {
		t.Run(tt.id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			for pass := range 2 {
				_, tokens, ds := lexer.Scan("source.alg", src)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				prog, ds := parser.Parse(tokens)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				info, ds := sema.Analyze(prog)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				host := &zeroElapsedHost{}
				var out bytes.Buffer
				ds = interp.New(interp.Options{Output: &out, Host: host}).Run(prog, info)
				if len(ds) != 0 || !bytes.Equal(out.Bytes(), want) || host.clockCalls != tt.calls {
					t.Fatalf("pass %d: diagnostics %v, output %q, clock calls %d", pass, ds, &out, host.clockCalls)
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
		})
	}
}
