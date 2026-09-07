package interp

import (
	"errors"
	"fmt"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

const (
	maxDepth     = 256
	maxCalls     = 256
	maxTextBytes = runtime.MaxTextBytes
)

func failure(pos token.Pos, code diag.Code, err error) error {
	if err == nil {
		return nil
	}
	var d diag.Diagnostic
	if errors.As(err, &d) {
		return d
	}
	return diag.Diagnostic{Code: code, Pos: pos, Message: err.Error(), Cause: err}
}

func diagnostics(err error, pos token.Pos, code diag.Code) []diag.Diagnostic {
	if err == nil {
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
