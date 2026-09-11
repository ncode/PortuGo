package interp

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

type environmentHost struct {
	HeadlessHost
	times      []time.Time
	clockCalls int
	echo       []bool
	err        error
}

func (h *environmentHost) SetEcho(enabled bool) error {
	h.echo = append(h.echo, enabled)
	return h.err
}

func (h *environmentHost) Now() time.Time {
	n := h.clockCalls
	h.clockCalls++
	if n >= len(h.times) {
		return time.Time{}
	}
	return h.times[n]
}

func TestEchoHost(t *testing.T) {
	src := "algoritmo \"echo\"\ninicio\nescreva(\"before\")\neco off\neco on\nfimalgoritmo"
	p, info := analyzed(t, src)
	for _, fail := range []bool{false, true} {
		var out bytes.Buffer
		h := &environmentHost{}
		if fail {
			h.err = errors.New("simulated host failure")
		}
		ds := New(Options{Host: h, Output: &out}).Run(p, info)
		if fail {
			if len(ds) != 1 || ds[0].Code != diag.RHost || ds[0].Pos != token.Pos(strings.Index(src, "eco off")) || !errors.Is(ds[0], h.err) || !reflect.DeepEqual(h.echo, []bool{false}) {
				t.Fatalf("echo failure: diagnostics=%v events=%v", ds, h.echo)
			}
		} else if len(ds) != 0 || !reflect.DeepEqual(h.echo, []bool{false, true}) {
			t.Fatalf("echo events: diagnostics=%v events=%v", ds, h.echo)
		}
		if out.String() != "before" {
			t.Fatalf("unexpected terminal output: %q", out.String())
		}
	}
}

func TestChronometerHost(t *testing.T) {
	start := time.Unix(0, 0)
	for _, tt := range []struct {
		name, body, want string
		times            []time.Time
	}{
		{"inactive", "cronometro off", "\nO cronômetro não foi iniciado.\n", nil},
		{"zero", "cronometro on\ncronometro off", "\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 0 segundo(s).\n", []time.Time{start, start}},
		{"milliseconds", "cronometro on\ncronometro off\ncronometro off", "\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 16 ms.\n\nO cronômetro não foi iniciado.\n", []time.Time{start, start.Add(16 * time.Millisecond)}},
		{"seconds", "cronometro on\ncronometro off", "\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 2 segundo(s) e 47 ms.\n", []time.Time{start, start.Add(2047 * time.Millisecond)}},
		{"fractional seconds", "cronometro on\ncronometro off", "\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 2 segundo(s) e 422 ms.\n", []time.Time{start, start.Add(2422 * time.Millisecond)}},
		{"over one minute", "cronometro on\ncronometro off", "\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 72 segundo(s) e 141 ms.\n", []time.Time{start, start.Add(72141 * time.Millisecond)}},
		{"restart", "cronometro on\ncronometro on\ncronometro off", "\nCronômetro iniciado.\n\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 16 ms.\n", []time.Time{start, start.Add(10 * time.Millisecond), start.Add(26 * time.Millisecond)}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p, info := analyzed(t, "algoritmo \"clock\"\ninicio\n"+tt.body+"\nfimalgoritmo")
			h := &environmentHost{times: tt.times}
			var out bytes.Buffer
			i := New(Options{Host: h, Output: &out})
			if ds := i.Run(p, info); len(ds) != 0 || out.String() != tt.want || h.clockCalls != len(tt.times) || i.chronometerRunning {
				t.Fatalf("clock: diagnostics=%v output=%q calls=%d running=%t", ds, out.String(), h.clockCalls, i.chronometerRunning)
			}
		})
	}
}

func TestChronometerResetAndFailures(t *testing.T) {
	on, onInfo := analyzed(t, "algoritmo \"start\"\ninicio\ncronometro on\nfimalgoritmo")
	off, offInfo := analyzed(t, "algoritmo \"stop\"\ninicio\ncronometro off\nfimalgoritmo")
	h := &environmentHost{}
	var out bytes.Buffer
	i := New(Options{Host: h, Output: &out})
	if ds := i.Run(on, onInfo); len(ds) != 0 {
		t.Fatal(ds)
	}
	out.Reset()
	if ds := i.Run(off, offInfo); len(ds) != 0 || out.String() != "\nO cronômetro não foi iniciado.\n" || h.clockCalls != 1 {
		t.Fatalf("clock state crossed runs: %v %q calls=%d", ds, out.String(), h.clockCalls)
	}
	src := "algoritmo \"clock regression\"\ninicio\ncronometro on\ncronometro off\nfimalgoritmo"
	p, info := analyzed(t, src)
	h = &environmentHost{times: []time.Time{time.Unix(1, 0), time.Unix(0, 0)}}
	out.Reset()
	i = New(Options{Host: h, Output: &out})
	ds := i.Run(p, info)
	if len(ds) != 1 || ds[0].Code != diag.RHost || ds[0].Pos != token.Pos(strings.Index(src, "cronometro off")) || out.String() != "\nCronômetro iniciado.\n" || i.chronometerRunning {
		t.Fatalf("clock regression: diagnostics=%v output=%q", ds, out.String())
	}
	ds = New(Options{Output: failingWriter{}}).Run(on, onInfo)
	if len(ds) != 1 || ds[0].Code != diag.RHost || ds[0].Pos != on.Body[0].Start() {
		t.Fatalf("clock output failure: %v", ds)
	}
}
