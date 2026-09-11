package interp

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
)

func (i *Interpreter) configureConsole(settings []ast.ConsoleStmt) error {
	for _, setting := range settings {
		if err := i.charge(setting.At); err != nil {
			return err
		}
		if err := i.options.Host.UseConsole(); err != nil {
			return diag.Diagnostic{Code: diag.RHost, Pos: setting.At, Message: "cannot configure console display", Cause: err}
		}
	}
	return nil
}
