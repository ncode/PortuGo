package stdlib

import (
	"fmt"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ncode/portugol-go/internal/cp1252"
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
		return squareRoot(args)
	case "exp":
		return exp(args)
	case "log":
		return checkedNumeric1(args, math.Log10)
	case "logn":
		return checkedNumeric1(args, math.Log)
	case "pi":
		return runtime.Value{Kind: runtime.RealValue, Real: math.Pi}, true, nil
	case "sen":
		return numeric1(args, math.Sin)
	case "cos":
		return numeric1(args, math.Cos)
	case "tan":
		return numeric1(args, math.Tan)
	case "arccos":
		return numeric1(args, math.Acos)
	case "arcsen":
		return numeric1(args, math.Asin)
	case "arctan":
		return numeric1(args, math.Atan)
	case "cotan":
		return numeric1(args, func(x float64) float64 { return 1 / math.Tan(x) })
	case "grauprad":
		return numeric1(args, func(x float64) float64 { return x * (math.Pi / 180) })
	case "radpgrau":
		return numeric1(args, func(x float64) float64 { return x * (180 / math.Pi) })
	case "quad":
		return square(args)
	case "int":
		return intval(args)
	case "aleatorio":
		return l.random(args)
	case "randi":
		return l.randi(args)
	case "copia":
		return copia(args)
	case "maiusc":
		return string1(args, func(r rune) rune {
			// These Unicode uppercase counterparts are outside Windows-1252.
			if r == 'µ' || r == 'ƒ' {
				return r
			}
			return unicode.ToUpper(r)
		})
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
	case "caracpnum":
		return caracpnum(args)
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
	return runtime.Value{Kind: runtime.StringValue, Str: runtime.FormatReal(x)}, true, nil
}

func abs(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) > 1 {
		return runtime.Value{}, true, fmt.Errorf("abs expects zero or one argument")
	}
	if len(args) == 0 {
		return runtime.Value{Kind: runtime.VoidValue}, true, nil
	}
	switch args[0].Kind {
	case runtime.IntegerValue:
		v := args[0].Int
		if v < 0 {
			v = -v
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: int64(int32(v))}, true, nil
	case runtime.RealValue:
		return runtime.Value{Kind: runtime.RealValue, Real: math.Abs(args[0].Real)}, true, nil
	case runtime.VoidValue:
		return args[0], true, nil
	case runtime.StringValue, runtime.BoolValue:
		return runtime.Value{Kind: runtime.VoidValue}, true, nil
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
	if len(args) == 0 {
		return runtime.Value{Kind: runtime.VoidValue}, true, nil
	}
	for _, arg := range args[:min(len(args), 2)] {
		if arg.Kind == runtime.StringValue || arg.Kind == runtime.BoolValue || arg.Kind == runtime.VoidValue {
			return runtime.Value{Kind: runtime.VoidValue}, true, nil
		}
	}
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
	result := math.Pow(base, exponent)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return runtime.Value{}, true, fmt.Errorf("invalid exp domain")
	}
	return runtime.Value{Kind: runtime.RealValue, Real: result}, true, nil
}

func intval(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) == 0 {
		return runtime.Value{Kind: runtime.IntegerValue}, true, nil
	}
	if args[0].Kind == runtime.VoidValue {
		return args[0], true, nil
	}
	if args[0].Kind == runtime.StringValue || args[0].Kind == runtime.BoolValue {
		return runtime.Value{Kind: runtime.VoidValue}, true, nil
	}
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("int expects 1 argument")
	}
	x, err := asFloat(args[0])
	if err != nil {
		return runtime.Value{}, true, err
	}
	if math.IsNaN(x) || math.IsInf(x, 0) || x < -0x1p63 || x >= 0x1p63 {
		return runtime.Value{}, true, fmt.Errorf("int argument exceeds the conversion range")
	}
	return runtime.Value{Kind: runtime.IntegerValue, Int: int64(int32(int64(x)))}, true, nil
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
	position, err := copyIndex(args[1])
	if err != nil {
		return runtime.Value{}, true, err
	}
	n, err := copyIndex(args[2])
	if err != nil {
		return runtime.Value{}, true, err
	}
	s := []rune(args[0].Str)
	start := max(int64(1), position) - 1
	if n < 0 || start >= int64(len(s)) {
		return runtime.Value{Kind: runtime.StringValue}, true, nil
	}
	end := start + min(n, int64(len(s))-start)
	return runtime.Value{Kind: runtime.StringValue, Str: string(s[int(start):int(end)])}, true, nil
}

func copyIndex(v runtime.Value) (int64, error) {
	if v.Kind == runtime.VoidValue {
		return 0, nil
	}
	if v.Kind == runtime.IntegerValue {
		return v.Int, nil
	}
	x, _, err := intval([]runtime.Value{v})
	if err == nil && x.Kind != runtime.IntegerValue {
		err = fmt.Errorf("copia expects numeric bounds")
	}
	return x.Int, err
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
	if args[0].Str == "" {
		return runtime.Value{Kind: runtime.VoidValue}, true, nil
	}
	r, _ := utf8.DecodeRuneInString(args[0].Str)
	code, ok := cp1252.EncodeRune(r)
	if !ok {
		return runtime.Value{}, true, fmt.Errorf("character is outside Windows-1252")
	}
	return runtime.Value{Kind: runtime.IntegerValue, Int: int64(code)}, true, nil
}

func carac(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) > 1 {
		return runtime.Value{}, true, fmt.Errorf("carac expects zero or one argument")
	}
	var code int64
	if len(args) == 1 && args[0].Kind != runtime.VoidValue {
		var err error
		code, err = asInt(args[0])
		if err != nil {
			return runtime.Value{}, true, err
		}
	}
	if code < 0 || code > 255 {
		return runtime.Value{Kind: runtime.VoidValue}, true, nil
	}
	if code < 32 {
		code = 32
	} else if code >= 127 {
		code = int64(characterBytes[code-127])
	}
	return runtime.Value{Kind: runtime.StringValue, Str: string(cp1252.DecodeByte(byte(code)))}, true, nil
}

func compr(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("compr expects 1 argument")
	}
	return runtime.Value{Kind: runtime.IntegerValue, Int: int64(utf8.RuneCountInString(args[0].Str))}, true, nil
}

func pos(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 2 {
		return runtime.Value{}, true, fmt.Errorf("pos expects 2 arguments")
	}
	haystack, needle := args[1].Str, args[0].Str
	if needle != "" {
		if offset := strings.Index(haystack, needle); offset >= 0 {
			return runtime.Value{Kind: runtime.IntegerValue, Int: int64(utf8.RuneCountInString(haystack[:offset]) + 1)}, true, nil
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
