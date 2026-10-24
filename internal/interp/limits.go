package interp

import (
	"errors"
	"fmt"

	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/runtime"
	"github.com/ncode/PortuGo/internal/token"
)

const (
	maxDepth     = 256
	maxCalls     = 256
	maxOperands  = maxDepth * maxCalls
	maxTextBytes = runtime.MaxTextBytes
)

// haltControl unwinds active evaluations without producing a diagnostic.
type haltControl struct{}

func (haltControl) Error() string { return "execution halted" }

func failure(pos token.Pos, code diag.Code, err error) error {
	if err == nil || errors.Is(err, haltControl{}) {
		return err
	}
	var d diag.Diagnostic
	if errors.As(err, &d) {
		return d
	}
	return diag.Diagnostic{Code: code, Pos: pos, Message: err.Error(), Cause: err}
}

func diagnostics(err error, pos token.Pos, code diag.Code) []diag.Diagnostic {
	if err == nil || errors.Is(err, haltControl{}) {
		return nil
	}
	var d diag.Diagnostic
	if errors.As(failure(pos, code, err), &d) {
		return []diag.Diagnostic{d}
	}
	return nil
}

func (i *Interpreter) charge(pos token.Pos) error {
	if i.options.MaxSteps == 0 {
		return nil
	}
	if i.steps == i.options.MaxSteps {
		return failure(pos, diag.RLoop, fmt.Errorf("execution step budget exhausted"))
	}
	i.steps++
	return nil
}
