package interp

import (
	"fmt"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func (i *Interpreter) resolveType(spec ast.TypeSpec) (runtime.Type, error) {
	typ := runtime.TypeFromSpec(spec)
	if typ.Kind == runtime.InvalidType {
		if binding, ok := i.info.Binding(token.Token{Text: spec.Name, Pos: spec.At}); ok {
			return binding.Type, nil
		}
		return runtime.Type{}, failure(spec.At, diag.RType, fmt.Errorf("missing named type binding"))
	}
	if typ.Kind != runtime.VectorType {
		return typ, nil
	}
	for n, r := range spec.Ranges {
		bounds := []*int64{&typ.Ranges[n].Low, &typ.Ranges[n].High}
		for j, expr := range []ast.Expr{r.Low, r.High} {
			if _, literal := expr.(*ast.LiteralExpr); literal {
				continue
			}
			v, err := i.eval(expr)
			if err != nil {
				return runtime.Type{}, err
			}
			if v.Kind != runtime.IntegerValue {
				return runtime.Type{}, failure(expr.Start(), diag.EParse, fmt.Errorf("vector bound is not an integer"))
			}
			*bounds[j] = v.Int
		}
		typ.Ranges[n].LowDynamic, typ.Ranges[n].HighDynamic = false, false
		if typ.Ranges[n].High < typ.Ranges[n].Low {
			return runtime.Type{}, failure(r.At, diag.EParse, fmt.Errorf("vector upper bound is smaller than lower bound"))
		}
	}
	if spec.Elem != nil {
		elem, err := i.resolveType(*spec.Elem)
		if err != nil {
			return runtime.Type{}, err
		}
		typ.Elem = &elem
	}
	if _, err := typ.Slots(); err != nil {
		return runtime.Type{}, failure(spec.At, diag.RStorage, err)
	}
	return typ, nil
}
