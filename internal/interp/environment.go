package interp

import (
	"fmt"
	"io"
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

func (i *Interpreter) execEcho(s *ast.EchoStmt) error {
	if s.Mode.Kind == token.IDENT {
		if err := i.options.Host.SetEcho(strings.EqualFold(s.Mode.Text, "on")); err != nil {
			return diag.Diagnostic{Code: diag.RHost, Pos: s.At, Message: "cannot configure input echo", Cause: err}
		}
	}
	return nil
}

func (i *Interpreter) execChronometer(s *ast.ChronometerStmt) error {
	text := "\nO cronômetro não foi iniciado.\n"
	if !s.Off {
		i.chronometerStart = i.options.Host.Now()
		i.chronometerRunning = true
		text = "\nCronômetro iniciado.\n"
	} else if i.chronometerRunning {
		elapsed := i.options.Host.Now().Sub(i.chronometerStart)
		i.chronometerRunning = false
		if elapsed < 0 {
			return failure(s.At, diag.RHost, fmt.Errorf("chronometer clock moved backwards"))
		}
		milliseconds := elapsed.Milliseconds()
		seconds, remainder := milliseconds/1000, milliseconds%1000
		amount := fmt.Sprintf("%d segundo(s).", seconds)
		if remainder != 0 {
			amount = fmt.Sprintf("%d ms.", remainder)
			if seconds != 0 {
				amount = fmt.Sprintf("%d segundo(s) e %s", seconds, amount)
			}
		}
		text = "\nCronômetro terminado. Tempo decorrido: " + amount + "\n"
	}
	if _, err := io.WriteString(i.out, text); err != nil {
		return diag.Diagnostic{Code: diag.RHost, Pos: s.At, Message: "cannot write chronometer output", Cause: err}
	}
	return nil
}
