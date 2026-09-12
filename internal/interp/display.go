package interp

import (
	"fmt"
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func (i *Interpreter) execColor(s *ast.ColorStmt) error {
	color, err := i.displayArgument(s.Color)
	if err != nil {
		return err
	}
	target, err := i.displayArgument(s.Target)
	if err != nil {
		return err
	}
	background := false
	switch strings.ToLower(target) {
	case "fundos":
		background = true
	case "frente":
	default:
		return nil
	}
	selected, ok := displayColor(color)
	if !ok {
		return nil
	}
	if err := i.options.Host.SetDisplay(DisplayState{Color: selected, Background: background}); err != nil {
		return diag.Diagnostic{Code: diag.RHost, Pos: s.Start(), Message: "cannot change display color", Cause: err}
	}
	return nil
}

func (i *Interpreter) displayArgument(expr ast.Expr) (string, error) {
	v, err := i.eval(expr)
	if err != nil {
		return "", err
	}
	switch v.Kind {
	case runtime.StringValue:
		return v.Str, nil
	case runtime.VoidValue:
		return "", failure(expr.Start(), diag.EParse, fmt.Errorf("missing color argument value"))
	default:
		return "", failure(expr.Start(), diag.ETypeMismatch, fmt.Errorf("expected caractere display argument"))
	}
}

func displayColor(name string) (Color, bool) {
	switch strings.ToLower(name) {
	case "preto":
		return Black, true
	case "azul":
		return Blue, true
	case "verde":
		return Green, true
	case "vermelho":
		return Red, true
	case "roxo":
		return Purple, true
	case "amarelo":
		return Yellow, true
	case "branco":
		return White, true
	}
	return Black, false
}
