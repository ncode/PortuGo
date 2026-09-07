package stdlib

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/ncode/portugol-go/internal/runtime"
)

// Library stores stateful built-ins such as aleatorio.
type Library struct {
	rng *rand.Rand
}

var builtinNames = map[string]bool{
	"abs": true, "raizq": true, "exp": true, "log": true, "logn": true, "pi": true,
	"sen": true, "cos": true, "tan": true, "int": true, "frac": true, "aleatorio": true,
	"copia": true, "maiusc": true, "minusc": true, "asc": true, "carac": true,
	"compr": true, "pos": true,
}

// IsBuiltin reports whether name is a standard function.
func IsBuiltin(name string) bool {
	return builtinNames[name]
}

// New creates a standard library instance.
func New() *Library {
	return &Library{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
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
	case "copia":
		return copia(args)
	case "maiusc":
		return string1(args, strings.ToUpper)
	case "minusc":
		return string1(args, strings.ToLower)
	case "asc":
		return asc(args)
	case "carac":
		return carac(args)
	case "compr":
		return compr(args)
	case "pos":
		return pos(args)
	default:
		return runtime.Value{}, false, nil
	}
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
		return runtime.Value{Kind: runtime.IntegerValue, Int: l.rng.Int63n(n)}, true, nil
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
		return runtime.Value{Kind: runtime.IntegerValue, Int: lo + l.rng.Int63n(hi-lo+1)}, true, nil
	default:
		return runtime.Value{}, true, fmt.Errorf("aleatorio expects 0 to 2 arguments")
	}
}

func copia(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 3 {
		return runtime.Value{}, true, fmt.Errorf("copia expects 3 arguments")
	}
	s := []rune(args[0].Str)
	start := args[1].Int - 1
	n := args[2].Int
	if start < 0 {
		start = 0
	}
	if n < 0 || start >= int64(len(s)) {
		return runtime.Value{Kind: runtime.StringValue}, true, nil
	}
	end := start + n
	if end > int64(len(s)) {
		end = int64(len(s))
	}
	return runtime.Value{Kind: runtime.StringValue, Str: string(s[int(start):int(end)])}, true, nil
}

func string1(args []runtime.Value, fn func(string) string) (runtime.Value, bool, error) {
	if len(args) != 1 {
		return runtime.Value{}, true, fmt.Errorf("function expects 1 argument")
	}
	return runtime.Value{Kind: runtime.StringValue, Str: fn(args[0].Str)}, true, nil
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
