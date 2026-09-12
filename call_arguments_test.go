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

func TestRecordedCallArguments(t *testing.T) {
	for _, id := range []string{
		"reference-integer-to-real",
		"reference-integer-as-real-read", "reference-integer-as-real-write",
		"reference-real-as-integer-read", "reference-real-as-integer-write",
		"value-argument-evaluation-order", "reference-designator-capture-order",
		"callee-global-binding", "ordinary-local-shadowing", "mutual-function-recursion",
		"reference-global-write-visibility", "reference-caller-write-visibility",
		"reference-widening-visibility", "reference-narrowing-visibility",
		"reference-alias-read-after-write", "reference-alias-conversion-order",
		"reference-narrowing-positive-half", "reference-narrowing-negative-half",
		"value-narrowing-positive-half", "value-narrowing-negative-half",
		"value-narrowing-positive-fraction", "value-narrowing-negative-fraction",
		"value-widening-format", "reference-nested-write-visibility",
		"reference-recursive-alias",
		"call-supplied-value-integer",
		"call-supplied-value-real",
		"call-supplied-value-logical",
		"call-supplied-value-character",
		"call-supplied-reference-integer",
		"call-supplied-reference-real",
		"call-empty-zero-parameters",
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

func TestRecordedReferenceTypeChangeFailure(t *testing.T) {
	dir := "testdata/conformance/visualg-3.0.7/probes/reference-narrowing-followup"
	src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file, tokens, ds := lexer.Scan("source.alg", src)
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
	var out bytes.Buffer
	ds = interp.New(interp.Options{Output: &out}).Run(prog, info)
	if len(ds) != 1 || ds[0].Code != diag.RType || file.Position(ds[0].Pos).Line != 12 || !bytes.Equal(out.Bytes(), want) {
		t.Fatalf("diagnostics = %v, output = %q; want R001 on line 12, output %q", ds, &out, want)
	}
}

func TestRecordedArgumentRejections(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"procedure-extra-arguments", diag.ECall, 2},
		{"bare-procedure-missing-argument", diag.ECall, 2},
		{"function-extra-arguments", diag.ECall, 7},
		{"function-missing-tail-argument", diag.ECall, 9},
		{"reference-string-as-integer", diag.ETypeMismatch, 4},
		{"procedure-missing-tail-argument", diag.ECall, 2},
		{"procedure-empty-two-arguments", diag.ECall, 2},
	} {
		t.Run(tt.id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg")
			src, err := source.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			prog, ds := parser.Parse(tokens)
			if len(ds) != 0 {
				t.Fatal(ds)
			}
			_, ds = sema.Analyze(prog)
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics = %v, want %s on line %d", ds, tt.code, tt.line)
			}
		})
	}
}
