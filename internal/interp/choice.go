package interp

import (
	"fmt"
	"strings"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/runtime"
)

func (i *Interpreter) matchesCase(selector runtime.Value, label ast.CaseLabel) (bool, error) {
	low, err := i.eval(label.Low)
	if err != nil || low.Kind == runtime.VoidValue {
		return false, err
	}
	if label.High == nil {
		if selector.Kind == runtime.VoidValue {
			return false, nil
		}
		if selector.Kind == runtime.StringValue && low.Kind == runtime.StringValue {
			// Only the label is uppercased by the reference.
			return selector.Str == strings.ToUpper(low.Str), nil
		}
		return equalValues(selector, low)
	}
	high, err := i.eval(label.High)
	if err != nil {
		return false, err
	}
	if high.Kind == runtime.VoidValue {
		return false, failure(label.High.Start(), diag.EParse, fmt.Errorf("expected range upper value"))
	}
	if selector.NumericAbsence && i.dynamicChoiceBound(label.High) {
		i.halted = true
		return false, nil
	}
	x, xNumeric := asFloat(selector)
	lo, loNumeric := asFloat(low)
	hi, hiNumeric := asFloat(high)
	return xNumeric && loNumeric && hiNumeric && lo <= x && x <= hi, nil
}

func (i *Interpreter) dynamicChoiceBound(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.IdentExpr:
		binding, ok := i.info.Binding(e.Name)
		return ok && !binding.Builtin
	case *ast.CallExpr:
		binding, ok := i.info.Binding(e.Name)
		return ok && !binding.Builtin
	default:
		return false
	}
}
