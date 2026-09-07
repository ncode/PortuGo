package interp

import (
	"io"
	"time"

	"github.com/ncode/portugol-go/internal/token"
)

// Options configures one interpreter. Zero MaxSteps leaves work unlimited.
type Options struct {
	Input      io.Reader
	Output     io.Writer
	WorkingDir string
	Host       Host
	Random     RandomSource
	MaxSteps   uint64
}

// RandomSource supplies independent draws without global interpreter state.
type RandomSource interface {
	Float64() float64
	Uint64N(uint64) uint64
}

// Breakpoint identifies the source location of a host pause request.
type Breakpoint struct{ Pos token.Pos }

// DisplayState carries display options once their syntax is reference-verified.
// No display options are exposed by the currently supported language.
type DisplayState struct{}

// Host implements the interpreter's typed environment operations.
type Host interface {
	Delay(time.Duration) error
	Breakpoint(Breakpoint) error
	ClearScreen() error
	SetDisplay(DisplayState) error
	Now() time.Time
}

// HeadlessHost uses real time and nonblocking, silent UI operations.
type HeadlessHost struct{}

// Delay waits for duration.
func (HeadlessHost) Delay(duration time.Duration) error { time.Sleep(duration); return nil }

// Breakpoint is a nonblocking no-op in headless execution.
func (HeadlessHost) Breakpoint(Breakpoint) error { return nil }

// ClearScreen does not emit terminal escape sequences.
func (HeadlessHost) ClearScreen() error { return nil }

// SetDisplay is a no-op in headless execution.
func (HeadlessHost) SetDisplay(DisplayState) error { return nil }

// Now returns the system clock.
func (HeadlessHost) Now() time.Time { return time.Now() }
