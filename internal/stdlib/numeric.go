package stdlib

import (
	"fmt"
	"math"

	"github.com/ncode/portugol-go/internal/runtime"
)

func numeric1(args []runtime.Value, fn func(float64) float64) (runtime.Value, bool, error) {
	if len(args) == 0 {
		args = []runtime.Value{{Kind: runtime.IntegerValue}}
	} else if args[0].Kind == runtime.VoidValue {
		return args[0], true, nil
	}
	v, ok, err := real1(args, fn)
	if err == nil && (math.IsNaN(v.Real) || math.IsInf(v.Real, 0)) {
		v = runtime.Value{Kind: runtime.VoidValue, NumericAbsence: true}
	}
	return v, ok, err
}

func checkedNumeric1(args []runtime.Value, fn func(float64) float64) (runtime.Value, bool, error) {
	v, ok, err := numeric1(args, fn)
	if err == nil && v.Kind == runtime.VoidValue && (len(args) == 0 || args[0].Kind != runtime.VoidValue) {
		err = fmt.Errorf("invalid numeric function domain")
	}
	return v, ok, err
}

func squareRoot(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) == 1 && (args[0].Kind == runtime.StringValue || args[0].Kind == runtime.BoolValue) {
		return runtime.Value{Kind: runtime.VoidValue}, true, nil
	}
	return checkedNumeric1(args, math.Sqrt)
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
			return checkedNumeric1(args, func(x float64) float64 { return x * x })
		case runtime.VoidValue:
			return v, true, nil
		}
	}
	return runtime.Value{Kind: runtime.VoidValue}, true, nil
}
