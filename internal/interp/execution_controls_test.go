package interp

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/lexer"
	"github.com/ncode/PortuGo/internal/parser"
	"github.com/ncode/PortuGo/internal/sema"
	"github.com/ncode/PortuGo/internal/source"
	"github.com/ncode/PortuGo/internal/token"
)

type executionHost struct {
	HeadlessHost
	out      *bytes.Buffer
	delays   []time.Duration
	prefixes []string
	breaks   []Breakpoint
	elapsed  time.Duration
	failAt   int
	err      error
}

func (h *executionHost) Delay(d time.Duration) error {
	h.delays = append(h.delays, d)
	h.prefixes = append(h.prefixes, h.out.String())
	if len(h.delays) == h.failAt {
		return h.err
	}
	h.elapsed += d
	return nil
}

func (h *executionHost) Breakpoint(b Breakpoint) error {
	h.breaks = append(h.breaks, b)
	h.prefixes = append(h.prefixes, h.out.String())
	return h.err
}

func (h *executionHost) Now() time.Time { return time.Unix(0, 0).Add(h.elapsed) }

func TestTimerHost(t *testing.T) {
	for _, tt := range []struct {
		name, body string
		prefixes   []string
	}{
		{"ordinary", "timer 2\nescreva(\"A\")\ntimer 0\nescreva(\"B\")", []string{"", "A"}},
		{"ignored types", "timer 2\nescreva(\"A\")\ntimer \"off\"\ntimer verdadeiro\nescreva(\"B\")\ntimer 0", []string{"", "A", "A", "A", "AB"}},
		{"negative", "timer 2\nescreva(\"A\")\ntimer -1\nescreva(\"B\")", []string{"", "A"}},
		{"fraction", "timer 2\nescreva(\"A\")\ntimer 0.5\nescreva(\"B\")", []string{"", "A"}},
		{"fraction below one", "timer 2\nescreva(\"A\")\ntimer 0.9\nescreva(\"B\")", []string{"", "A"}},
		{"conditional", "timer 2\nescreva(\"A\")\nse verdadeiro entao\nescreva(\"B\")\nfimse\nescreva(\"C\")\ntimer 0", []string{"", "A", "A", "AB", "ABC"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p, info := analyzed(t, "algoritmo \"timer\"\ninicio\n"+tt.body+"\nfimalgoritmo")
			var out bytes.Buffer
			h := &executionHost{out: &out}
			if ds := New(Options{Host: h, Output: &out}).Run(p, info); len(ds) != 0 {
				t.Fatal(ds)
			}
			if !reflect.DeepEqual(h.prefixes, tt.prefixes) {
				t.Fatalf("delay ordering: %q; want %q", h.prefixes, tt.prefixes)
			}
			for _, d := range h.delays {
				if d != 2*time.Millisecond {
					t.Fatalf("delay=%s; want 2ms", d)
				}
			}
		})
	}
}

func TestTimerClockOrder(t *testing.T) {
	for _, body := range []string{
		"timer 1000\ncronometro on\ntimer 0\ncronometro off",
		"cronometro on\ntimer 1000\ncronometro off\ntimer 0",
	} {
		p, info := analyzed(t, "algoritmo \"clock order\"\ninicio\n"+body+"\nfimalgoritmo")
		var out bytes.Buffer
		h := &executionHost{out: &out}
		if ds := New(Options{Host: h, Output: &out}).Run(p, info); len(ds) != 0 || out.String() != "\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 1 segundo(s).\n" || h.elapsed != 2*time.Second {
			t.Fatalf("clock ordering: %v %q elapsed=%s", ds, out.String(), h.elapsed)
		}
	}
}

func TestTimerStructures(t *testing.T) {
	for _, tt := range []struct {
		name, declarations, body string
		calls                    int
	}{
		{"procedure", "procedimento P\ninicio\nescreval(\"B\")\nfimprocedimento\n", "P", 5},
		{"one local", "procedimento P\nvar\ni: inteiro\ninicio\nescreval(\"B\")\nfimprocedimento\n", "P", 7},
		{"grouped locals", "procedimento P\nvar\ni,j: inteiro\ninicio\nescreval(\"B\")\nfimprocedimento\n", "P", 7},
		{"separate locals", "procedimento P\nvar\ni: inteiro\nj: inteiro\ninicio\nescreval(\"B\")\nfimprocedimento\n", "P", 8},
		{"function", "funcao F: inteiro\ninicio\nescreval(\"B\")\nretorne 1\nfimfuncao\n", "escreval(F())", 6},
		{"for break", "var\ni: inteiro\n", "para i de 1 ate 2 faca\ninterrompa\nfimpara", 3},
		{"while break", "", "enquanto verdadeiro faca\ninterrompa\nfimenquanto", 3},
		{"switch", "", "escolha 2\ncaso 1\nescreval(\"A\")\ncaso 2\nescreval(\"B\")\nfimescolha", 3},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p, info := analyzed(t, "algoritmo \"timer structures\"\n"+tt.declarations+"inicio\ntimer 500\n"+tt.body+"\ntimer 0\nfimalgoritmo")
			var out bytes.Buffer
			h := &executionHost{out: &out}
			if ds := New(Options{Host: h, Output: &out}).Run(p, info); len(ds) != 0 || len(h.delays) != tt.calls || h.elapsed != time.Duration(tt.calls)*500*time.Millisecond {
				t.Fatalf("structure timing: %v delays=%v elapsed=%s", ds, h.delays, h.elapsed)
			}
		})
	}
}

func TestBreakOutsideLoopIsIgnored(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"outside break\"\ninicio\ninterrompa\nescreval(\"AFTER\")\nfimalgoritmo")
	var out bytes.Buffer
	if ds := New(Options{Output: &out}).Run(p, info); len(ds) != 0 || out.String() != "AFTER\n" {
		t.Fatalf("outside-loop break: diagnostics=%v output=%q", ds, out.String())
	}
}

func TestBreakDoesNotEscapeCallee(t *testing.T) {
	src := "algoritmo \"callee break\"\nvar\ni: inteiro\nprocedimento P\ninicio\ninterrompa\nfimprocedimento\ninicio\npara i de 1 ate 2 faca\nP\nfimpara\nescreval(i)\nfimalgoritmo"
	p, info := analyzed(t, src)
	var out bytes.Buffer
	if ds := New(Options{Output: &out}).Run(p, info); len(ds) != 0 || out.String() != " 2\n" {
		t.Fatalf("callee break: diagnostics=%v output=%q", ds, out.String())
	}
}

func TestBreakpointHost(t *testing.T) {
	for _, command := range []string{"pausa", "pausa()", "pausa ignored", "debug verdadeiro", "debug 1=1"} {
		t.Run(command, func(t *testing.T) {
			src := "algoritmo \"breakpoint\"\ninicio\nescreva(\"A\")\n" + command + "\nescreva(\"B\")\nfimalgoritmo"
			p, info := analyzed(t, src)
			for _, fail := range []bool{false, true} {
				var out bytes.Buffer
				h := &executionHost{out: &out}
				if fail {
					h.err = errors.New("simulated host failure")
				}
				ds := New(Options{Host: h, Output: &out}).Run(p, info)
				at := token.Pos(strings.Index(src, command))
				if len(h.breaks) != 1 || h.breaks[0].Pos != at || !reflect.DeepEqual(h.prefixes, []string{"A"}) {
					t.Fatalf("breakpoint trace: %+v %q", h.breaks, h.prefixes)
				}
				if fail {
					if len(ds) != 1 || ds[0].Code != diag.RHost || ds[0].Pos != at || !errors.Is(ds[0], h.err) || out.String() != "A" {
						t.Fatalf("breakpoint failure: %v %q", ds, out.String())
					}
				} else if len(ds) != 0 || out.String() != "AB" {
					t.Fatalf("breakpoint success: %v %q", ds, out.String())
				}
			}
		})
	}
}

func TestBreakpointSuppression(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"no breakpoint\"\ninicio\ndebug falso\ndebug 1=2\nescreval(timer(),debug(),pausa())\nescreval(\"DONE\")\nfimalgoritmo")
	var out bytes.Buffer
	h := &executionHost{out: &out}
	if ds := New(Options{Host: h, Output: &out}).Run(p, info); len(ds) != 0 || len(h.breaks) != 0 || len(h.delays) != 0 || out.String() != "DONE\n" {
		t.Fatalf("inactive commands: %v breaks=%v delays=%v output=%q", ds, h.breaks, h.delays, out.String())
	}
}

func TestTimerResetAndFailures(t *testing.T) {
	var out bytes.Buffer
	h := &executionHost{out: &out}
	i := New(Options{Host: h, Output: &out})
	for _, body := range []string{"timer 2", "escreva(\"B\")"} {
		p, info := analyzed(t, "algoritmo \"reset\"\ninicio\n"+body+"\nfimalgoritmo")
		if ds := i.Run(p, info); len(ds) != 0 {
			t.Fatal(ds)
		}
	}
	if len(h.delays) != 1 || i.timerDelay != 0 {
		t.Fatalf("timer crossed runs: %v %s", h.delays, i.timerDelay)
	}
	for _, tt := range []struct {
		name, body, at, output string
		code                   diag.Code
		failAt, calls          int
	}{
		{"host failure", "timer 2\nescreva(\"A\")\nescreva(\"B\")", "escreva(\"A\")", "A", diag.RHost, 2, 2},
		{"no value", "escreva(\"A\")\ntimer abs()\nescreva(\"B\")", "abs()", "A", diag.EParse, 0, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := "algoritmo \"failure\"\ninicio\n" + tt.body + "\nfimalgoritmo"
			p, info := analyzed(t, src)
			out.Reset()
			h = &executionHost{out: &out, failAt: tt.failAt, err: errors.New("simulated host failure")}
			ds := New(Options{Host: h, Output: &out}).Run(p, info)
			if len(ds) != 1 || ds[0].Code != tt.code || ds[0].Pos != token.Pos(strings.Index(src, tt.at)) || out.String() != tt.output || len(h.delays) != tt.calls {
				t.Fatalf("failure: %v output=%q calls=%d", ds, out.String(), len(h.delays))
			}
			if tt.failAt != 0 && !errors.Is(ds[0], h.err) {
				t.Fatal("host error cause lost")
			}
		})
	}
}

func TestTimerArgumentEvaluation(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"timer argument\"\nfuncao Marker: inteiro\ninicio\nescreva(\"A\")\nretorne 2\nfimfuncao\ninicio\ntimer Marker()\ntimer 0\nfimalgoritmo")
	var out bytes.Buffer
	h := &executionHost{out: &out}
	if ds := New(Options{Host: h, Output: &out}).Run(p, info); len(ds) != 0 || !reflect.DeepEqual(h.prefixes, []string{"A"}) || !reflect.DeepEqual(h.delays, []time.Duration{2 * time.Millisecond}) {
		t.Fatalf("argument evaluation: %v prefixes=%q delays=%v", ds, h.prefixes, h.delays)
	}
}

func TestTimerCallFailureRestoresFrame(t *testing.T) {
	src := "algoritmo \"frame failure\"\nprocedimento P\ninicio\nescreva(\"unreached\")\nfimprocedimento\ninicio\ntimer 2\nP\nfimalgoritmo"
	p, info := analyzed(t, src)
	var out bytes.Buffer
	h := &executionHost{out: &out, failAt: 2, err: errors.New("simulated host failure")}
	i := New(Options{Host: h, Output: &out})
	ds := i.Run(p, info)
	if len(ds) != 1 || ds[0].Code != diag.RHost || ds[0].Pos != token.Pos(strings.LastIndex(src, "\nP\n")+1) || !errors.Is(ds[0], h.err) || out.Len() != 0 || i.calls != 0 || i.env != i.global || i.result != nil {
		t.Fatalf("frame failure: %v output=%q calls=%d", ds, out.String(), i.calls)
	}
}

func TestTimerFractionalBoundaries(t *testing.T) {
	for _, tt := range []struct {
		literal string
		want    time.Duration
	}{
		{"2.9", 2 * time.Millisecond},
		{"30000", 10 * time.Second},
		{"2147483648", 10 * time.Second},
		{"9223372036855", 10 * time.Second},
		{"9223372036854.5", 10 * time.Second},
	} {
		p, info := analyzed(t, "algoritmo \"timer boundary\"\ninicio\ntimer "+tt.literal+"\nfimalgoritmo")
		var out bytes.Buffer
		h := &executionHost{out: &out}
		if ds := New(Options{Host: h}).Run(p, info); len(ds) != 0 || !reflect.DeepEqual(h.delays, []time.Duration{tt.want}) {
			t.Fatalf("timer %s: %v delays=%v", tt.literal, ds, h.delays)
		}
	}
}

func TestRecordedDeferredExecutionCommands(t *testing.T) {
	for _, tt := range []struct {
		id   string
		code diag.Code
		line int
	}{
		{"execution-debug-missing", diag.EParse, 4},
		{"execution-debug-number", diag.ETypeMismatch, 4},
		{"execution-timer-zero", diag.EParse, 6},
		{"execution-timer-negative", diag.EParse, 6},
		{"execution-timer-fraction", diag.EParse, 6},
		{"execution-timer-on-off", diag.EParse, 4},
		{"execution-timer-bare", diag.EParse, 4},
		{"execution-timer-unknown", diag.EParse, 4},
		{"execution-timer-quoted", diag.EParse, 6},
		{"execution-timer-expression", diag.EParse, 6},
	} {
		t.Run(tt.id, func(t *testing.T) {
			dir := filepath.Join("../../testdata/conformance/visualg-3.0.7/probes", tt.id)
			data, err := os.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			for pass := range 2 {
				src, err := source.DecodeFile("source.alg", data)
				if err != nil {
					t.Fatal(err)
				}
				file, tokens, ds := lexer.ScanFile(src)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				prog, ds := parser.Parse(tokens)
				if len(ds) != 0 {
					t.Fatalf("premature parse diagnostics: %v", ds)
				}
				info, ds := sema.Analyze(prog)
				if len(ds) != 0 {
					t.Fatalf("premature semantic diagnostics: %v", ds)
				}
				var out bytes.Buffer
				host := &executionHost{out: &out}
				ds = New(Options{Output: &out, Host: host}).Run(prog, info)
				if len(ds) != 1 || ds[0].Code != tt.code || file.Position(ds[0].Pos).Line != tt.line || !bytes.Equal(out.Bytes(), want) {
					t.Fatalf("diagnostics=%v stdout=%q; want %s line %d and %q", ds, out.String(), tt.code, tt.line, want)
				}
				if len(host.delays) != 0 || len(host.breaks) != 0 {
					t.Fatalf("unexpected host effects: delays=%v breakpoints=%v", host.delays, host.breaks)
				}
				var formatted bytes.Buffer
				if err := ast.Fprint(&formatted, prog); err != nil {
					t.Fatal(err)
				}
				if pass != 0 && !bytes.Equal(formatted.Bytes(), data) {
					t.Fatal("formatting is not idempotent")
				}
				data = formatted.Bytes()
			}
		})
	}
}

func TestDeferredExecutionCommandEffects(t *testing.T) {
	for _, tt := range []struct {
		command string
		code    diag.Code
	}{
		{"timer", diag.EParse}, {"timer unknown", diag.EParse},
		{"debug", diag.EParse}, {"debug 0", diag.ETypeMismatch},
	} {
		t.Run(tt.command, func(t *testing.T) {
			prog, info := analyzed(t, "algoritmo \"command failure\"\ninicio\ntimer 2\nescreva(\"A\")\n"+tt.command+"\nescreva(\"B\")\nfimalgoritmo")
			var out bytes.Buffer
			host := &executionHost{out: &out}
			ds := New(Options{Output: &out, Host: host}).Run(prog, info)
			if len(ds) != 1 || ds[0].Code != tt.code || out.String() != "A" {
				t.Fatalf("diagnostics=%v stdout=%q", ds, out.String())
			}
			if !reflect.DeepEqual(host.delays, []time.Duration{2 * time.Millisecond, 2 * time.Millisecond}) ||
				!reflect.DeepEqual(host.prefixes, []string{"", "A"}) || len(host.breaks) != 0 {
				t.Fatalf("failed command changed host effects: delays=%v prefixes=%q breakpoints=%v", host.delays, host.prefixes, host.breaks)
			}
		})
	}

	const src = `algoritmo "timer operand order"
funcao marker: inteiro
inicio
escreva("A")
retorne 2
fimfuncao
inicio
timer marker() + unknown
escreva("B")
fimalgoritmo
`
	prog, info := analyzed(t, src)
	var out bytes.Buffer
	host := &executionHost{out: &out}
	ds := New(Options{Output: &out, Host: host}).Run(prog, info)
	if len(ds) != 1 || ds[0].Code != diag.EParse || ds[0].Pos != token.Pos(strings.Index(src, "unknown")) || out.String() != "A" || len(host.delays) != 0 || len(host.breaks) != 0 {
		t.Fatalf("operand error changed effects: diagnostics=%v stdout=%q delays=%v breakpoints=%v", ds, out.String(), host.delays, host.breaks)
	}
}
