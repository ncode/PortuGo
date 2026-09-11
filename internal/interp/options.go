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

// Color is a reference display color encoded as 0xRRGGBB.
type Color uint32

// Recorded display colors.
const (
	Black  Color = 0x000000
	Blue   Color = 0x0000ff
	Green  Color = 0x008000
	Red    Color = 0xff0000
	Purple Color = 0x800080
	Yellow Color = 0xffff00
	White  Color = 0xffffff
)

// DisplayState changes one color without replacing the other display settings.
type DisplayState struct {
	Color      Color
	Background bool
}

// Host implements the interpreter's typed environment operations.
type Host interface {
	UseConsole() error
	Delay(time.Duration) error
	Breakpoint(Breakpoint) error
	ClearScreen() error
	SetDisplay(DisplayState) error
	SetEcho(bool) error
	Now() time.Time
}

// HeadlessHost uses real time and nonblocking, silent UI operations.
type HeadlessHost struct{}

// UseConsole is a no-op because headless output already uses the supplied writer.
func (HeadlessHost) UseConsole() error { return nil }

// Delay waits for duration.
func (HeadlessHost) Delay(duration time.Duration) error { time.Sleep(duration); return nil }

// Breakpoint is a nonblocking no-op in headless execution.
func (HeadlessHost) Breakpoint(Breakpoint) error { return nil }

// ClearScreen does not emit terminal escape sequences.
func (HeadlessHost) ClearScreen() error { return nil }

// SetDisplay is a no-op in headless execution.
func (HeadlessHost) SetDisplay(DisplayState) error { return nil }

// SetEcho does not alter the deterministic console transcript.
func (HeadlessHost) SetEcho(bool) error { return nil }

// Now returns the system clock.
func (HeadlessHost) Now() time.Time { return time.Now() }
