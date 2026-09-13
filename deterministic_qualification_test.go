package portugol_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
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

type recordedClockHost struct {
	interp.HeadlessHost
	nowMS      []int64
	clockCalls int
	delays     []time.Duration
	operations []string
}

func (h *recordedClockHost) Now() time.Time {
	h.operations = append(h.operations, "now")
	if h.clockCalls >= len(h.nowMS) {
		return time.Time{}
	}
	ms := h.nowMS[h.clockCalls]
	h.clockCalls++
	return time.Unix(0, 0).Add(time.Duration(ms) * time.Millisecond)
}

func (h *recordedClockHost) Delay(d time.Duration) error {
	h.operations = append(h.operations, "delay")
	h.delays = append(h.delays, d)
	return nil
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

func TestRecordedTimedChronometers(t *testing.T) {
	for _, tt := range []struct {
		id     string
		nowMS  []int64
		delays []time.Duration
	}{
		{"chronometer-repeat-stop", []int64{0, 16}, nil},
		{"environment-chronometer-milliseconds", []int64{0, 359}, []time.Duration{150 * time.Millisecond, 150 * time.Millisecond}},
		{"environment-chronometer-seconds", []int64{0, 2047}, []time.Duration{time.Second, time.Second}},
		{"environment-chronometer-fractional-seconds", []int64{0, 2422}, []time.Duration{1200 * time.Millisecond, 1200 * time.Millisecond}},
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
				host := &recordedClockHost{nowMS: tt.nowMS}
				var out bytes.Buffer
				ds = interp.New(interp.Options{Output: &out, Host: host}).Run(prog, info)
				if len(ds) != 0 || !bytes.Equal(out.Bytes(), want) || host.clockCalls != len(tt.nowMS) || !reflect.DeepEqual(host.delays, tt.delays) {
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

func TestRecordedTimerTimingBoundaries(t *testing.T) {
	for _, id := range []string{
		"environment-timer-clock-start", "environment-timer-clock-stop",
		"environment-timer-upper-bound", "environment-timer-real-delay",
		"environment-timer-large-delay", "environment-timer-enormous-delay",
		"environment-timer-fraction-loop", "environment-timer-rounding-zero",
		"environment-timer-rounding-one", "environment-timer-rounding-one-half",
		"environment-timer-rounding-two-half", "environment-chronometer-minute",
	} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
			clockData, err := os.ReadFile(filepath.Join(dir, "expected-clock.json"))
			if err != nil {
				t.Fatal(err)
			}
			var clock struct {
				NowMS []int64 `json:"nowMS"`
			}
			if err := json.Unmarshal(clockData, &clock); err != nil {
				t.Fatal(err)
			}
			if len(clock.NowMS) == 0 {
				t.Fatal("clock fixture has no reads")
			}
			traceData, err := os.ReadFile(filepath.Join(dir, "expected-host.json"))
			if err != nil {
				t.Fatal(err)
			}
			var trace []struct {
				Operation string        `json:"operation"`
				Duration  time.Duration `json:"duration"`
			}
			if err := json.Unmarshal(traceData, &trace); err != nil {
				t.Fatal(err)
			}
			wantOperations := make([]string, len(trace))
			var wantDelays []time.Duration
			for n, event := range trace {
				wantOperations[n] = event.Operation
				if event.Operation == "delay" {
					wantDelays = append(wantDelays, event.Duration)
				}
			}
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
				host := &recordedClockHost{nowMS: clock.NowMS}
				var out bytes.Buffer
				ds = interp.New(interp.Options{Output: &out, Host: host}).Run(prog, info)
				if len(ds) != 0 || !bytes.Equal(out.Bytes(), want) || host.clockCalls != len(clock.NowMS) ||
					!reflect.DeepEqual(host.operations, wantOperations) || !reflect.DeepEqual(host.delays, wantDelays) {
					t.Fatalf("pass %d: diagnostics %v, output %q, operations %v, delays %v", pass, ds, &out, host.operations, host.delays)
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
