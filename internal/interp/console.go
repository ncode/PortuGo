package interp

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
)

func (i *Interpreter) configure(settings []ast.Stmt) error {
	for _, setting := range settings {
		if err := i.charge(setting.Start()); err != nil {
			return err
		}
		switch s := setting.(type) {
		case *ast.ConsoleStmt:
			if err := i.options.Host.UseConsole(); err != nil {
				return diag.Diagnostic{Code: diag.RHost, Pos: s.At, Message: "cannot configure console display", Cause: err}
			}
		case *ast.FileInputStmt:
			if err := i.configureFile(s); err != nil {
				return err
			}
		}
	}
	return nil
}
