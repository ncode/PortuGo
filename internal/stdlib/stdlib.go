package stdlib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ncode/portugol-go/internal/runtime"
)

// Library stores stateful random built-ins.
type Library struct {
	rng RandomSource
}

// RandomSource supplies draws for random built-ins.
type RandomSource interface {
	Float64() float64
	Uint64N(uint64) uint64
}

// New creates a standard library instance.
func New(random RandomSource) *Library {
	return &Library{rng: random}
}

// Call invokes a built-in function by canonical lowercase name.
func (l *Library) Call(name string, args []runtime.Value) (runtime.Value, bool, error) {
	switch name {
	case "abs":
		return abs(args)
	case "raizq":
		return real1(args, math.Sqrt)
	case "exp":
		return exp(args)
	case "log":
		return real1(args, math.Log)
	case "logn":
		return logn(args)
	case "pi":
		return runtime.Value{Kind: runtime.RealValue, Real: math.Pi}, true, nil
	case "sen":
		return real1(args, math.Sin)
	case "cos":
		return real1(args, math.Cos)
	case "tan":
		return real1(args, math.Tan)
	case "int":
		return intval(args)
	case "frac":
		return frac(args)
	case "aleatorio":
		return l.random(args)
	case "randi":
		return l.randi(args)
	case "copia":
		return copia(args)
	case "maiusc":
		return string1(args, unicode.ToUpper)
	case "minusc":
		return string1(args, unicode.ToLower)
	case "asc":
		return asc(args)
	case "carac":
		return carac(args)
	case "compr":
		return compr(args)
	case "pos":
		return pos(args)
	case "numpcarac":
		return numpcarac(args)
	default:
		return runtime.Value{}, false, nil
	}
}

func numpcarac(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) > 1 {
		return runtime.Value{}, true, fmt.Errorf("numpcarac expects zero or one argument")
	}
	var x float64
	if len(args) == 1 {
		if args[0].Kind == runtime.StringValue || args[0].Kind == runtime.BoolValue || args[0].Kind == runtime.VoidValue {
			return runtime.Value{Kind: runtime.VoidValue}, true, nil
		}
		var err error
		x, err = asFloat(args[0])
		if err != nil {
			return runtime.Value{}, true, err
		}
	}
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return runtime.Value{}, true, fmt.Errorf("numpcarac requires a finite value")
	}
	if x == 0 {
		return runtime.Value{Kind: runtime.StringValue, Str: "0"}, true, nil
	}
	text := strconv.FormatFloat(x, 'G', 15, 64)
	if mantissa, exponent, ok := strings.Cut(text, "E"); ok {
		n, _ := strconv.Atoi(exponent) // FormatFloat always emits a valid exponent.
		text = mantissa + "E" + strconv.Itoa(n)
	}
	return runtime.Value{Kind: runtime.StringValue, Str: text}, true, nil
}

func abs(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("abs expects 1 argument")
	}
	switch args[0].Kind {
	case runtime.IntegerValue:
		v := args[0].Int
		if v < 0 {
			v = -v
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: v}, true, nil
	case runtime.RealValue:
		return runtime.Value{Kind: runtime.RealValue, Real: math.Abs(args[0].Real)}, true, nil
	default:
		return runtime.Value{}, true, fmt.Errorf("abs expects numeric argument")
	}
}

func real1(args []runtime.Value, fn func(float64) float64) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("function expects 1 argument")
	}
	x, err := asFloat(args[0])
	if err != nil {
		return runtime.Value{}, true, err
	}
	return runtime.Value{Kind: runtime.RealValue, Real: fn(x)}, true, nil
}

func exp(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 2 {
		return runtime.Value{}, true, fmt.Errorf("exp expects 2 arguments")
	}
	base, err := asFloat(args[0])
	if err != nil {
		return runtime.Value{}, true, err
	}
	exponent, err := asFloat(args[1])
	if err != nil {
		return runtime.Value{}, true, err
	}
	return runtime.Value{Kind: runtime.RealValue, Real: math.Pow(base, exponent)}, true, nil
}

func logn(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 2 {
		return runtime.Value{}, true, fmt.Errorf("logn expects 2 arguments")
	}
	x, err := asFloat(args[0])
	if err != nil {
		return runtime.Value{}, true, err
	}
	base, err := asFloat(args[1])
	if err != nil {
		return runtime.Value{}, true, err
	}
	return runtime.Value{Kind: runtime.RealValue, Real: math.Log(x) / math.Log(base)}, true, nil
}

func intval(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("int expects 1 argument")
	}
	x, err := asFloat(args[0])
	if err != nil {
		return runtime.Value{}, true, err
	}
	return runtime.Value{Kind: runtime.IntegerValue, Int: int64(x)}, true, nil
}

func frac(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("frac expects 1 argument")
	}
	x, err := asFloat(args[0])
	if err != nil {
		return runtime.Value{}, true, err
	}
	return runtime.Value{Kind: runtime.RealValue, Real: x - float64(int64(x))}, true, nil
}

func (l *Library) random(args []runtime.Value) (runtime.Value, bool, error) {
	switch len(args) {
	case 0:
		return runtime.Value{Kind: runtime.RealValue, Real: l.rng.Float64()}, true, nil
	case 1:
		n, err := asInt(args[0])
		if err != nil {
			return runtime.Value{}, true, err
		}
		if n <= 0 {
			return runtime.Value{}, true, fmt.Errorf("aleatorio upper bound must be positive")
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: int64(l.rng.Uint64N(uint64(n)))}, true, nil
	case 2:
		lo, err := asInt(args[0])
		if err != nil {
			return runtime.Value{}, true, err
		}
		hi, err := asInt(args[1])
		if err != nil {
			return runtime.Value{}, true, err
		}
		if hi < lo {
			lo, hi = hi, lo
		}
		width := uint64(hi) - uint64(lo) + 1
		if width == 0 {
			return runtime.Value{}, true, fmt.Errorf("aleatorio range exceeds representable draw size")
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: int64(uint64(lo) + l.rng.Uint64N(width))}, true, nil
	default:
		return runtime.Value{}, true, fmt.Errorf("aleatorio expects 0 to 2 arguments")
	}
}

func copia(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 3 {
		return runtime.Value{}, true, fmt.Errorf("copia expects 3 arguments")
	}
	s := []rune(args[0].Str)
	start := max(int64(1), args[1].Int) - 1
	n := args[2].Int
	if n < 0 || start >= int64(len(s)) {
		return runtime.Value{Kind: runtime.StringValue}, true, nil
	}
	end := start + min(n, int64(len(s))-start)
	return runtime.Value{Kind: runtime.StringValue, Str: string(s[int(start):int(end)])}, true, nil
}

func string1(args []runtime.Value, fn func(rune) rune) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("function expects 1 argument")
	}
	size := 0
	for _, r := range args[0].Str {
		n := utf8.RuneLen(fn(r))
		if size > runtime.MaxTextBytes-n {
			return runtime.Value{}, true, runtime.ErrTextSize
		}
		size += n
	}
	return runtime.Value{Kind: runtime.StringValue, Str: strings.Map(fn, args[0].Str)}, true, nil
}

func asc(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("asc expects 1 argument")
	}
	for _, r := range args[0].Str {
		return runtime.Value{Kind: runtime.IntegerValue, Int: int64(r)}, true, nil
	}
	return runtime.Value{Kind: runtime.IntegerValue}, true, nil
}

func carac(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("carac expects 1 argument")
	}
	return runtime.Value{Kind: runtime.StringValue, Str: string(rune(args[0].Int))}, true, nil
}

func compr(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("compr expects 1 argument")
	}
	return runtime.Value{Kind: runtime.IntegerValue, Int: int64(len([]rune(args[0].Str)))}, true, nil
}

func pos(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 2 {
		return runtime.Value{}, true, fmt.Errorf("pos expects 2 arguments")
	}
	haystack := []rune(args[1].Str)
	needle := []rune(args[0].Str)
	if len(needle) == 0 {
		return runtime.Value{Kind: runtime.IntegerValue, Int: 1}, true, nil
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if string(haystack[i:i+len(needle)]) == string(needle) {
			return runtime.Value{Kind: runtime.IntegerValue, Int: int64(i + 1)}, true, nil
		}
	}
	return runtime.Value{Kind: runtime.IntegerValue}, true, nil
}

func asFloat(v runtime.Value) (float64, error) {
	switch v.Kind {
	case runtime.IntegerValue:
		return float64(v.Int), nil
	case runtime.RealValue:
		return v.Real, nil
	default:
		return 0, fmt.Errorf("expected numeric argument")
	}
}

func asInt(v runtime.Value) (int64, error) {
	if v.Kind != runtime.IntegerValue {
		return 0, fmt.Errorf("expected inteiro argument")
	}
	return v.Int, nil
}
