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
		"numeric-boundary-abs", "numeric-boundary-arccos", "numeric-boundary-arcsen",
		"numeric-boundary-arctan", "numeric-boundary-cos", "numeric-boundary-cotan",
		"numeric-boundary-grauprad", "numeric-boundary-int", "numeric-boundary-log",
		"numeric-boundary-logn", "numeric-boundary-quad", "numeric-boundary-radpgrau",
		"numeric-boundary-raizq", "numeric-boundary-sen", "numeric-boundary-tan",
		"numeric-boundary-exp-positive", "numeric-boundary-exp-negative-odd",
		"numeric-boundary-exp-negative-even", "numeric-boundary-exp-zero-positive",
		"numeric-boundary-exp-zero-zero", "numeric-boundary-exp-underflow",
		"numeric-boundary-grauprad-large", "numeric-boundary-radpgrau-large", "numeric-boundary-abs-large",
		"numeric-boundary-exp-absent-short", "numeric-boundary-exp-absent-effect",
		"numeric-domain-radians-large", "numeric-domain-degrees-negative-large",
		"numeric-domain-exp-domain-short-absent", "numeric-domain-exp-domain-extra-absent",
		"numeric-domain-exp-second-extra-absent", "numeric-domain-exp-nested-abs-absent",
		"numeric-domain-exp-nested-arctan-absent", "numeric-domain-exp-cotangent-absent",
		"numeric-domain-exp-angle-absent",
		"exp-absent-tail-conversion", "exp-absent-tail-logical", "exp-absent-tail-text",
		"exp-absent-tail-twice", "exp-absent-tail-nested-absence",
		"numeric-absence-type-integer", "numeric-absence-type-conversion",
		"numeric-absence-type-exponentiation", "numeric-absence-type-cotangent",
		"numeric-absence-type-angle",
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
		line int
	}{
		{"legacy-log-empty", diag.RBuiltin, 3},
		{"legacy-logn-empty", diag.RBuiltin, 3},
		{"legacy-frac-value", diag.EUndeclared, 3},
		{"legacy-frac-empty", diag.EUndeclared, 3},
		{"abs-extra", diag.EParse, 3},
		{"int-extra", diag.EParse, 3},
		{"log-negative", diag.RBuiltin, 3},
		{"logn-extra", diag.EParse, 3},
		{"exp-one-argument", diag.EParse, 3},
		{"exp-negative-domain", diag.RBuiltin, 3},
		{"exp-overflow", diag.RBuiltin, 3},
		{"cos-string", diag.EParse, 3},
		{"cos-logical", diag.EParse, 3},
		{"arctan-bare", diag.EParse, 3},
		{"quad-bare", diag.EParse, 3},
		{"abs-string-extra", diag.EParse, 3},
		{"quad-string-extra", diag.EParse, 3},
		{"numpcarac-string-extra", diag.EParse, 3},
		{"numeric-boundary-arctan-absent-extra", diag.EParse, 3},
		{"numeric-boundary-log-zero", diag.RBuiltin, 3},
		{"numeric-boundary-logn-zero", diag.RBuiltin, 3},
		{"numeric-domain-exp-domain-short-numeric", diag.EParse, 6},
		{"numeric-domain-exp-domain-extra-numeric", diag.EParse, 11},
		{"numeric-domain-exp-second-extra-numeric", diag.EParse, 11},
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
			if out.Len() != 0 {
				t.Fatalf("output = %q, want empty", &out)
			}
		})
	}
}
