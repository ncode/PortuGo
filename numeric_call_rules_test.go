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

func TestRecordedNumericCallRules(t *testing.T) {
	for _, id := range []string{
		"legacy-abs-value", "legacy-abs-empty", "legacy-raizq-value", "legacy-raizq-empty",
		"legacy-exp-value", "legacy-exp-empty", "legacy-log-value", "legacy-logn-value",
		"legacy-sen-value", "legacy-sen-empty", "legacy-cos-value", "legacy-cos-empty",
		"legacy-tan-value", "legacy-tan-empty", "legacy-int-value", "legacy-int-empty",
		"abs-string", "abs-logical", "abs-no-value", "int-string", "int-logical",
		"log-unit", "logn-unit", "exp-string", "exp-no-value", "sqrt-string",
		"exp-no-value-right-effect", "exp-no-value-right-failure", "exp-no-value-left-effect",
		"exp-logical-right", "int-no-value", "sqrt-logical",
		"exp-string-right-effect", "exp-logical-right-effect", "exp-string-one",
		"exp-logical-one", "exp-no-value-one", "sqrt-no-value", "log-no-value", "logn-no-value",
		"exp-string-extra", "int-string-extra", "arctan-no-value-extra",
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

func TestNumericCallRuleDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
	}{
		{"legacy-log-empty", diag.RBuiltin},
		{"legacy-logn-empty", diag.RBuiltin},
		{"legacy-frac-value", diag.EUndeclared},
		{"legacy-frac-empty", diag.EUndeclared},
		{"abs-extra", diag.EParse},
		{"int-extra", diag.EParse},
		{"log-negative", diag.RBuiltin},
		{"logn-extra", diag.EParse},
		{"exp-one-argument", diag.EParse},
		{"exp-negative-domain", diag.RBuiltin},
		{"exp-overflow", diag.RBuiltin},
		{"cos-string", diag.EParse},
		{"cos-logical", diag.EParse},
		{"arctan-bare", diag.EParse},
		{"quad-bare", diag.EParse},
		{"abs-string-extra", diag.EParse},
		{"quad-string-extra", diag.EParse},
		{"numpcarac-string-extra", diag.EParse},
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
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != 3 {
				t.Fatalf("diagnostics = %v, want %s on line 3", ds, tt.code)
			}
			if out.Len() != 0 {
				t.Fatalf("output = %q, want empty", &out)
			}
		})
	}
}
