package stdlib

import (
	"fmt"
	"math"

	"github.com/ncode/portugol-go/internal/runtime"
)

func (l *Library) randi(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) > 1 {
		return runtime.Value{}, true, fmt.Errorf("randi expects zero or one argument")
	}
	var n int64
	if len(args) == 1 {
		var err error
		n, err = asInt(args[0])
		if err != nil {
			return runtime.Value{}, true, err
		}
	}
	if n < math.MinInt32 || n > math.MaxInt32 {
		return runtime.Value{}, true, fmt.Errorf("randi bound must fit a signed 32-bit integer")
	}
	// The reference interprets a negative bound as an unsigned 32-bit width.
	width := uint64(uint32(n))
	if width == 0 {
		return runtime.Value{Kind: runtime.IntegerValue}, true, nil
	}
	draw, err := l.randomUint64(width)
	if err != nil {
		return runtime.Value{}, true, err
	}
	return runtime.Value{Kind: runtime.IntegerValue, Int: int64(int32(draw))}, true, nil
}

func (l *Library) randomUint64(width uint64) (uint64, error) {
	if width == 0 {
		return 0, fmt.Errorf("random bound must be positive")
	}
	if l.rng == nil {
		return 0, fmt.Errorf("random source is unavailable")
	}
	draw := l.rng.Uint64N(width)
	if draw >= width {
		return 0, fmt.Errorf("random source returned an out-of-range value")
	}
	return draw, nil
}
