package portugol_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
)

// TestRecordedWindowsProbes replays the selected desktop observations. This is
// deliberately narrower than the planned complete conformance corpus.
func TestRecordedWindowsProbes(t *testing.T) {
	for _, name := range []string{"observations.json", "followup-observations.json"} {
		data, err := os.ReadFile("docs/validation/visualg-2026-09-07/" + name)
		if err != nil {
			t.Fatal(err)
		}
		var record struct {
			Probes []struct {
				ID, Source, SourceSHA256, ReferencePanelText string
			}
		}
		if err := json.Unmarshal(data, &record); err != nil {
			t.Fatal(err)
		}
		if len(record.Probes) == 0 {
			t.Fatal("empty reference observations")
		}
		for _, probe := range record.Probes {
			t.Run(probe.ID, func(t *testing.T) {
				// These reduced sources are ASCII, so UTF-8 and CP1252 bytes agree.
				if hash := fmt.Sprintf("%x", sha256.Sum256([]byte(probe.Source))); !strings.EqualFold(hash, probe.SourceSHA256) {
					t.Fatal("recorded source hash mismatch")
				}
				file, toks, lexDiags := lexer.Scan(probe.ID+".alg", probe.Source)
				prog, parseDiags := parser.Parse(toks)
				if len(lexDiags)+len(parseDiags) != 0 {
					t.Fatalf("unexpected diagnostics: %v %v", lexDiags, parseDiags)
				}
				info, diags := sema.Analyze(prog)
				var code diag.Code
				var line int
				count := 1
				switch probe.ID {
				case "exp-one":
					code, line = diag.ECall, 4
				case "exact-division-integer":
					code, line = diag.ETypeMismatch, 6
				case "output-rounding":
					code, line = diag.ETypeMismatch, 5
					count = 3
				}
				if code != "" {
					if len(diags) != count {
						t.Fatalf("diagnostics = %v, want %d", diags, count)
					}
					for _, d := range diags {
						if d.Code != code || file.Position(d.Pos).Line != line {
							t.Fatalf("diagnostics = %v, want %s on line %d", diags, code, line)
						}
					}
					return // The CLI rejects these before execution; VisuAlg reports them while running.
				}
				if len(diags) != 0 {
					t.Fatalf("unexpected diagnostics: %v", diags)
				}
				want, ok := strings.CutPrefix(probe.ReferencePanelText, "Início da execução\r\n")
				if !ok {
					t.Fatal("missing reference start notice")
				}
				want, ok = strings.CutSuffix(want, "\r\nFim da execução.\r\n")
				if !ok {
					t.Fatal("missing reference completion notice")
				}
				// The portable CLI writes LF; no other program-output bytes change.
				want = strings.ReplaceAll(want, "\r\n", "\n")
				var out bytes.Buffer
				if err := interp.New(interp.Options{Output: &out}).Run(prog, info); err != nil {
					t.Fatal(err)
				}
				if out.String() != want {
					t.Fatalf("stdout = %q, reference program output = %q", out.String(), want)
				}
			})
		}
	}
}
