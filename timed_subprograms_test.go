package portugol_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/interp"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/sema"
	"github.com/ncode/PortuGo/internal/source"
)

func TestRecordedTimedSubprograms(t *testing.T) {
	for _, tt := range []struct {
		id         string
		elapsedMS  int64
		delayCount int
	}{
		{"environment-timed-procedure", 2547, 5},
		{"environment-timed-procedure-local", 3563, 7},
		{"environment-timed-function", 3094, 6},
		{"environment-timed-for-break", 1531, 3},
		{"environment-timed-while-break", 1578, 3},
		{"environment-timed-switch", 1532, 3},
		{"environment-timed-one-local", 3578, 7},
		{"environment-timed-separate-locals", 4078, 8},
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
				host := &recordedClockHost{nowMS: []int64{0, tt.elapsedMS}}
				var out bytes.Buffer
				ds = interp.New(interp.Options{Output: &out, Host: host}).Run(prog, info)
				wantDelays := make([]time.Duration, tt.delayCount)
				for i := range wantDelays {
					wantDelays[i] = 500 * time.Millisecond
				}
				if len(ds) != 0 || !bytes.Equal(out.Bytes(), want) || host.clockCalls != 2 || !reflect.DeepEqual(host.delays, wantDelays) {
					t.Fatalf("pass %d: diagnostics %v, output %q, clock calls %d, delays %v", pass, ds, &out, host.clockCalls, host.delays)
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
