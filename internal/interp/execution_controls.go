package interp

import (
	"fmt"
	"math"
	"time"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func (i *Interpreter) execTimer(s *ast.TimerStmt) error {
	v, err := i.eval(s.Value)
	if err != nil {
		return err
	}
	if v.Kind == runtime.VoidValue {
		return failure(s.Value.Start(), diag.EParse, fmt.Errorf("timer requires a value"))
	}
	n, ok := asFloat(v)
	if !ok || v.Comparison {
		return nil
	}
	n = math.Trunc(n)
	const maxMilliseconds = math.MaxInt64 / int64(time.Millisecond)
	if math.IsNaN(n) || math.IsInf(n, 0) || n > float64(maxMilliseconds) {
		return failure(s.Value.Start(), diag.RHost, fmt.Errorf("timer duration out of range"))
	}
	i.timerDelay = time.Duration(int64(max(0, n))) * time.Millisecond
	return nil
}

func (i *Interpreter) delay(at token.Pos) error {
	if i.timerDelay <= 0 {
		return nil
	}
	if err := i.options.Host.Delay(i.timerDelay); err != nil {
		return diag.Diagnostic{Code: diag.RHost, Pos: at, Message: "cannot delay execution", Cause: err}
	}
	return nil
}

func (i *Interpreter) breakpoint(at token.Pos) error {
	if err := i.options.Host.Breakpoint(Breakpoint{Pos: at}); err != nil {
		return diag.Diagnostic{Code: diag.RHost, Pos: at, Message: "cannot pause execution", Cause: err}
	}
	return nil
}
