package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

var errObservationSize = errors.New("observation size limit exceeded")
var errClockExhausted = errors.New("clock fixture exhausted")

// executeProbe is a deterministic subprocess adapter for state/host fixtures.
// Ordinary output/diagnostic fixtures continue to use the actual CLI.
func executeProbe(args []string, in io.Reader, out, stderr io.Writer) int {
	flags := flag.NewFlagSet("execute", flag.ContinueOnError)
	flags.SetOutput(stderr)
	steps := flags.Uint64("max-steps", 10000, "finite probe work budget")
	state := flags.String("state", "", "global-state observation file")
	host := flags.String("host-trace", "", "host-call observation file")
	clock := flags.String("clock", "", "deterministic clock fixture")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 1 || *steps == 0 {
		return 2
	}
	fail := func(err error) int { _, _ = fmt.Fprintln(stderr, publicError(err)); return 1 }
	data, err := readBoundedFile(flags.Arg(0), 64<<10, fmt.Errorf("source exceeds replay profile"))
	if err != nil {
		return fail(err)
	}
	decoded, err := source.DecodeFile(flags.Arg(0), data)
	if err != nil {
		return fail(err)
	}
	file, toks, ds := lexer.ScanFile(decoded)
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
	var nowMS []int64
	if *clock != "" {
		data, err := readBoundedFile(*clock, maxArtifactBytes, fmt.Errorf("clock fixture exceeds artifact size limit"))
		if err != nil {
			return fail(err)
		}
		nowMS, err = decodeClockFixture(data)
		if err != nil {
			return fail(err)
		}
	}
	h := &recordingHost{events: []hostEvent{}, nowMS: nowMS}
	i := interp.New(interp.Options{Input: in, Output: capture, Host: h, Random: rand.New(rand.NewPCG(1, 2)), MaxSteps: *steps})
	ds = i.Run(p, info)
	if h.clockExhausted {
		return fail(errClockExhausted)
	}
	if _, err := out.Write(capture.buffer.Bytes()); err != nil {
		return fail(err)
	}
	diag.Render(stderr, file, ds)
	if h.overflow {
		return fail(errObservationSize)
	}
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

func readBoundedFile(name string, limit int, tooLarge error) ([]byte, error) {
	info, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("invalid file size or type: %s", name)
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	data, readErr := io.ReadAll(io.LimitReader(f, int64(limit)+1))
	closeErr := f.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if len(data) > limit {
		return nil, tooLarge
	}
	return data, nil
}

func decodeClockFixture(data []byte) ([]int64, error) {
	var fixture struct {
		NowMS []int64 `json:"nowMS"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&fixture); err != nil {
		return nil, fmt.Errorf("decode clock fixture: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return nil, fmt.Errorf("trailing clock fixture JSON")
	}
	if len(fixture.NowMS) == 0 {
		return nil, fmt.Errorf("clock fixture has no reads")
	}
	return fixture.NowMS, nil
}

func writeObservation(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if len(data) >= maxObservationBytes {
		return errObservationSize
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
	Echo      *bool                `json:"echo,omitempty"`
}

type recordingHost struct {
	elapsed        time.Duration
	events         []hostEvent
	eventSize      int
	overflow       bool
	nowMS          []int64
	nowRead        int
	clockExhausted bool
}

func (h *recordingHost) Delay(d time.Duration) error {
	if err := h.appendEvent(hostEvent{Operation: "delay", Duration: d}); err != nil {
		return err
	}
	h.elapsed += d
	return nil
}
func (h *recordingHost) Breakpoint(e interp.Breakpoint) error {
	return h.appendEvent(hostEvent{Operation: "breakpoint", Pos: e.Pos})
}
func (h *recordingHost) ClearScreen() error {
	return h.appendEvent(hostEvent{Operation: "clearScreen"})
}
func (h *recordingHost) UseConsole() error {
	return h.appendEvent(hostEvent{Operation: "console"})
}
func (h *recordingHost) SetDisplay(s interp.DisplayState) error {
	return h.appendEvent(hostEvent{Operation: "display", Display: &s})
}
func (h *recordingHost) SetEcho(enabled bool) error {
	return h.appendEvent(hostEvent{Operation: "echo", Echo: &enabled})
}
func (h *recordingHost) Now() time.Time {
	_ = h.appendEvent(hostEvent{Operation: "now"})
	if h.nowRead < len(h.nowMS) {
		ms := h.nowMS[h.nowRead]
		h.nowRead++
		return time.Unix(0, 0).UTC().Add(time.Duration(ms) * time.Millisecond)
	}
	h.nowRead++
	if len(h.nowMS) != 0 {
		h.clockExhausted = true
	}
	return time.Unix(0, 0).UTC().Add(h.elapsed)
}

func (h *recordingHost) appendEvent(event hostEvent) error {
	if h.overflow {
		return errObservationSize
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		return err
	}
	size := h.eventSize
	if size == 0 {
		size = 2
	}
	if len(h.events) != 0 {
		size++
	}
	size += len(encoded)
	if size >= maxObservationBytes {
		h.overflow = true
		return errObservationSize
	}
	h.events = append(h.events, event)
	h.eventSize = size
	return nil
}
