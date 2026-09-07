package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"time"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/token"
)

// executeProbe is a deterministic subprocess adapter for state/host fixtures.
// Ordinary output/diagnostic fixtures continue to use the actual CLI.
func executeProbe(args []string, in io.Reader, out, stderr io.Writer) int {
	flags := flag.NewFlagSet("execute", flag.ContinueOnError)
	flags.SetOutput(stderr)
	steps := flags.Uint64("max-steps", 10000, "finite probe work budget")
	state := flags.String("state", "", "global-state observation file")
	host := flags.String("host-trace", "", "host-call observation file")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 1 || *steps == 0 {
		return 2
	}
	fail := func(err error) int { _, _ = fmt.Fprintln(stderr, err); return 1 }
	f, err := os.Open(flags.Arg(0))
	if err != nil {
		return fail(err)
	}
	data, err := io.ReadAll(io.LimitReader(f, (64<<10)+1))
	closeErr := f.Close()
	if err != nil {
		return fail(err)
	}
	if closeErr != nil {
		return fail(closeErr)
	}
	if len(data) > 64<<10 {
		return fail(fmt.Errorf("source exceeds replay profile"))
	}
	decoded, err := source.Decode(data)
	if err != nil {
		return fail(err)
	}
	file, toks, ds := lexer.Scan(flags.Arg(0), decoded)
	if diag.HasErrors(ds) {
		diag.Render(stderr, file, ds)
		return 1
	}
	p, ds := parser.Parse(toks)
	if diag.HasErrors(ds) {
		diag.Render(stderr, file, ds)
		return 1
	}
	info, ds := sema.Analyze(p)
	if diag.HasErrors(ds) {
		diag.Render(stderr, file, ds)
		return 1
	}
	capture := &boundedOutput{cancel: func() {}}
	h := &recordingHost{events: []hostEvent{}}
	i := interp.New(interp.Options{Input: in, Output: capture, Host: h, Random: rand.New(rand.NewPCG(1, 2)), MaxSteps: *steps})
	ds = i.Run(p, info)
	if _, err := out.Write(capture.buffer.Bytes()); err != nil {
		return fail(err)
	}
	diag.Render(stderr, file, ds)
	if *state != "" {
		values := make(map[string]any)
		for name, value := range i.State() {
			values[name] = observedValue(value)
		}
		if err := writeObservation(*state, values); err != nil {
			return fail(err)
		}
	}
	if *host != "" {
		if err := writeObservation(*host, h.events); err != nil {
			return fail(err)
		}
	}
	if diag.HasErrors(ds) || capture.overflow {
		return 1
	}
	return 0
}

func writeObservation(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if len(data) >= 1<<20 {
		return fmt.Errorf("observation size limit exceeded")
	}
	return os.WriteFile(path, append(data, '\n'), 0600)
}

func observedValue(v runtime.Value) any {
	switch v.Kind {
	case runtime.IntegerValue:
		return v.Int
	case runtime.RealValue:
		return v.Real
	case runtime.StringValue:
		return v.Str
	case runtime.BoolValue:
		return v.Bool
	case runtime.VectorValue:
		if v.Vec == nil {
			return nil
		}
		values := make([]any, len(v.Vec.Elements))
		for n, cell := range v.Vec.Elements {
			values[n] = observedValue(cell.Value)
		}
		return struct {
			Bounds []runtime.Range `json:"bounds"`
			Values []any           `json:"values"`
		}{v.Vec.Type.Ranges, values}
	default:
		return nil
	}
}

type hostEvent struct {
	Operation string               `json:"operation"`
	Duration  time.Duration        `json:"duration,omitempty"`
	Pos       token.Pos            `json:"pos,omitempty"`
	Display   *interp.DisplayState `json:"display,omitempty"`
}

type recordingHost struct {
	elapsed time.Duration
	events  []hostEvent
}

func (h *recordingHost) Delay(d time.Duration) error {
	h.events = append(h.events, hostEvent{Operation: "delay", Duration: d})
	h.elapsed += d
	return nil
}
func (h *recordingHost) Breakpoint(e interp.Breakpoint) error {
	h.events = append(h.events, hostEvent{Operation: "breakpoint", Pos: e.Pos})
	return nil
}
func (h *recordingHost) ClearScreen() error {
	h.events = append(h.events, hostEvent{Operation: "clearScreen"})
	return nil
}
func (h *recordingHost) SetDisplay(s interp.DisplayState) error {
	h.events = append(h.events, hostEvent{Operation: "display", Display: &s})
	return nil
}
func (h *recordingHost) Now() time.Time {
	h.events = append(h.events, hostEvent{Operation: "now"})
	return time.Unix(0, 0).UTC().Add(h.elapsed)
}
