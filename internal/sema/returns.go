package sema

import "github.com/ncode/portugol-go/internal/ast"

func allPathsReturn(stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.ReturnStmt:
			return true
		case *ast.IfStmt:
			if len(s.Else) > 0 && allPathsReturn(s.Then) && allPathsReturn(s.Else) {
				return true
			}
		case *ast.SwitchStmt:
			if len(s.Default) == 0 || !allPathsReturn(s.Default) {
				continue
			}
			allCases := true
			for _, cc := range s.Cases {
				if !allPathsReturn(cc.Body) {
					allCases = false
					break
				}
			}
			if allCases {
				return true
			}
		}
	}
	return false
}
