package interp

import (
	"fmt"
	"math"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
)

type randomInputState struct {
	active    bool
	low, high float64
	decimals  int
}

func (i *Interpreter) execRandomInput(s *ast.RandomInputStmt) error {
	if s.Off {
		i.randomInput.active = false
		return nil
	}
	state := randomInputState{active: true, high: 100}
	for index, arg := range s.Args {
		value, err := i.eval(arg)
		if err != nil {
			return err
		}
		n, ok := asFloat(value)
		if !ok {
			return failure(arg.Start(), diag.EParse, fmt.Errorf("expected numeric random-input bound"))
		}
		if math.IsNaN(n) || math.IsInf(n, 0) {
			return failure(arg.Start(), diag.RInput, fmt.Errorf("random input requires finite numeric bounds"))
		}
		switch index {
		case 0:
			state.low = n
		case 1:
			state.high = n
		case 2:
			state.decimals = int(max(0, min(5, n)))
		}
	}
	if state.high < state.low {
		state.low, state.high = state.high, state.low
	}
	i.randomInput = state
	return nil
}
