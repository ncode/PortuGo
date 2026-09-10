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

func TestRecordedCharacterConversions(t *testing.T) {
	for _, id := range []string{
		"text-conversion-integer", "text-conversion-real", "text-conversion-comma",
		"text-conversion-signed", "text-conversion-invalid", "text-conversion-empty",
		"conversion-plus", "conversion-leading-space", "conversion-trailing-space",
		"conversion-exponent", "conversion-exponent-upper", "conversion-leading-dot",
		"conversion-leading-comma", "conversion-trailing-dot", "conversion-trailing-comma",
		"conversion-leading-zeroes", "conversion-integer-max", "conversion-integer-overflow",
		"conversion-unsigned-overflow", "conversion-large-integer",
		"conversion-negative-zero", "conversion-real-zero", "conversion-integer-kind",
		"conversion-invalid-kind", "conversion-mixed-case", "conversion-integer-assignment",
		"conversion-edge-signed-integer", "conversion-edge-signed-real", "conversion-edge-signed-exponent",
		"conversion-edge-plus-real", "conversion-edge-plus-space", "conversion-edge-only-plus", "conversion-edge-only-minus",
		"conversion-edge-overflow-middle", "conversion-edge-unsigned-max", "conversion-edge-unsigned-next",
		"conversion-edge-exponent-integer-max", "conversion-edge-real-integer-max", "conversion-edge-alphabetic-prefix",
		"conversion-edge-nan", "conversion-edge-infinity", "conversion-edge-exponent-underflow",
		"conversion-edge-incomplete-exponent", "conversion-edge-negative-kind", "conversion-edge-empty-kind",
		"conversion-edge-hex-prefix", "conversion-edge-currency-prefix",
		"conversion-context-integer-sum", "conversion-context-abs-kind", "conversion-context-character-code",
		"conversion-context-width-integer", "conversion-context-tab-space", "conversion-context-nonbreaking-space",
		"conversion-context-index-integer", "conversion-context-loop-integer", "conversion-context-integer-return",
		"conversion-context-integer-parameter", "conversion-context-integer-bound", "conversion-context-argument-effect",
		"conversion-class-double-zero", "conversion-class-empty-hex", "conversion-class-hex-invalid",
		"conversion-class-hex-kind", "conversion-class-hex-overflow", "conversion-class-hex-signed",
		"conversion-class-hex-upper", "conversion-class-incomplete-e", "conversion-class-incomplete-e-minus",
		"conversion-class-incomplete-e-plus", "conversion-class-incomplete-sign", "conversion-class-leading-zeros-real",
		"conversion-class-octal-like", "conversion-class-integer-window",
		"conversion-integer-boundary", "conversion-integer-threshold",
		"conversion-spelling-negative-window", "conversion-spelling-plus-currency",
		"conversion-spelling-negative-hex", "conversion-spelling-plus-hex",
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

func TestCharacterConversionDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"conversion-suffix", diag.RBuiltin, 3},
		{"conversion-two-dots", diag.RBuiltin, 3},
		{"conversion-two-commas", diag.RBuiltin, 3},
		{"conversion-internal-space", diag.RBuiltin, 3},
		{"conversion-bare", diag.EParse, 3},
		{"conversion-empty-call", diag.ETypeMismatch, 3},
		{"conversion-extra", diag.EParse, 3},
		{"conversion-numeric", diag.ETypeMismatch, 3},
		{"conversion-real", diag.ETypeMismatch, 3},
		{"conversion-logical", diag.ETypeMismatch, 3},
		{"conversion-no-value", diag.ETypeMismatch, 3},
		{"conversion-decimal-kind", diag.ETypeMismatch, 3},
		{"conversion-real-assignment", diag.RType, 7},
		{"conversion-edge-comma-kind", diag.ETypeMismatch, 3},
		{"conversion-edge-exponent-kind", diag.ETypeMismatch, 3},
		{"conversion-edge-trailing-dot-kind", diag.ETypeMismatch, 3},
		{"conversion-edge-exponent-overflow", diag.RBuiltin, 3},
		{"conversion-edge-mixed-separators", diag.RBuiltin, 3},
		{"conversion-edge-suffix-at-255", diag.RBuiltin, 3},
		{"conversion-context-argument-failure", diag.RBuiltin, 8},
		{"conversion-context-index-real", diag.ETypeMismatch, 6},
		{"conversion-context-logical-assignment", diag.RType, 5},
		{"conversion-context-real-bound", diag.EParse, 5},
		{"conversion-context-real-return", diag.ETypeMismatch, 4},
		{"conversion-context-real-character-code", diag.ETypeMismatch, 3},
		{"conversion-context-real-sum", diag.ETypeMismatch, 3},
		{"conversion-context-decimals-real", diag.EParse, 3},
		{"conversion-context-width-real", diag.EParse, 3},
		{"conversion-context-loop-real-start", diag.EParse, 5},
		{"conversion-context-loop-real-end", diag.EParse, 5},
		{"conversion-context-loop-real-step", diag.EParse, 5},
		{"conversion-class-hex-larger", diag.RBuiltin, 3},
		{"conversion-class-incomplete-kind", diag.ETypeMismatch, 3},
		{"conversion-spelling-underscore-integer", diag.RBuiltin, 3},
		{"conversion-spelling-underscore-real", diag.RBuiltin, 3},
		{"conversion-spelling-hex-suffix", diag.RBuiltin, 3},
	} {
		t.Run(tt.id, func(t *testing.T) {
			src, err := source.ReadFile(filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.id, "source.alg"))
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
			want := ""
			if tt.id == "conversion-context-argument-failure" {
				want = "EFFECT\n"
			}
			if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line || out.String() != want {
				t.Fatalf("diagnostics = %v, output = %q; want %s on line %d, output %q", ds, &out, tt.code, tt.line, want)
			}
		})
	}
}
