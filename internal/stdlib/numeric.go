package stdlib

import (
	"fmt"
	"math"

	"github.com/ncode/portugol-go/internal/runtime"
)

func numeric1(args []runtime.Value, fn func(float64) float64) (runtime.Value, bool, error) {
	if len(args) == 0 {
		args = []runtime.Value{{Kind: runtime.IntegerValue}}
	} else if len(args) == 1 && args[0].Kind == runtime.VoidValue {
		return args[0], true, nil
	}
	v, ok, err := real1(args, fn)
	if err == nil && (math.IsNaN(v.Real) || math.IsInf(v.Real, 0)) {
		v = runtime.Value{Kind: runtime.VoidValue}
	}
	return v, ok, err
}

func square(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) > 1 {
		return runtime.Value{}, true, fmt.Errorf("quad expects zero or one argument")
	}
	if len(args) == 1 {
		switch v := args[0]; v.Kind {
		case runtime.IntegerValue:
			return runtime.Value{Kind: runtime.IntegerValue, Int: int64(int32(v.Int * v.Int))}, true, nil
		case runtime.RealValue:
			return numeric1(args, func(x float64) float64 { return x * x })
		}
	}
	return runtime.Value{Kind: runtime.VoidValue}, true, nil
}
