package stdlib

import (
	"fmt"
	"math"

	"github.com/ncode/portugol-go/internal/runtime"
)

// RandomInput generates a scalar leia value using the library's shared source.
func (l *Library) RandomInput(kind runtime.TypeKind, low, high float64, decimals int) (runtime.Value, error) {
	if kind == runtime.StringType {
		var letters [5]byte
		for n := range letters {
			draw, err := l.randomUint64(26)
			if err != nil {
				return runtime.Value{}, err
			}
			letters[n] = 'A' + byte(draw)
		}
		return runtime.Value{Kind: runtime.StringValue, Str: string(letters[:])}, nil
	}
	if kind != runtime.IntegerType && kind != runtime.RealType {
		return runtime.Value{}, fmt.Errorf("unsupported random-input destination")
	}
	width := math.Floor(high-low) + 1
	if !(width >= 1 && width < 0x1p64) || math.IsInf(low, 0) || math.IsInf(high, 0) || decimals < 0 || decimals > 5 {
		return runtime.Value{}, fmt.Errorf("random-input range is not representable")
	}
	draw, err := l.randomUint64(uint64(width))
	if err != nil {
		return runtime.Value{}, err
	}
	n := low + float64(draw)
	if kind == runtime.IntegerType {
		if !(n >= -0x1p63 && n < 0x1p63) {
			return runtime.Value{}, fmt.Errorf("random-input integer is not representable")
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: int64(int32(int64(n)))}, nil
	}
	if decimals != 0 {
		scale := math.Pow10(decimals)
		fraction, err := l.randomUint64(uint64(scale))
		if err != nil {
			return runtime.Value{}, err
		}
		n += float64(fraction) / scale
	}
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return runtime.Value{}, fmt.Errorf("random-input value is not finite")
	}
	return runtime.Value{Kind: runtime.RealValue, Real: n}, nil
}
