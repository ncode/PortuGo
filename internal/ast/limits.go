package ast

import (
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

// MaxDepth bounds syntax nesting and subsequent AST traversals.
const MaxDepth = 256

// CheckLimits validates traversal depth before analysis or printing begins.
// Program/declaration lists are containers; their root statements/types start
// at depth one. Each nested statement, expression, or element type adds a level.
func CheckLimits(prog *Program) []diag.Diagnostic {
	if prog == nil {
		return []diag.Diagnostic{{Code: diag.EParse, Message: "missing program"}}
	}
	c := &limitChecker{}
	c.consts(prog.Consts)
	c.types(prog.Types)
	c.decls(prog.Globals)
	for _, sub := range prog.Subs {
		switch s := sub.(type) {
		case *ProcedureDecl:
			c.consts(s.Consts)
			c.types(s.Types)
			c.params(s.Params)
			c.decls(s.Locals)
			c.stmts(s.Body, 1)
		case *FunctionDecl:
			c.consts(s.Consts)
			c.types(s.Types)
			c.params(s.Params)
			c.typ(s.Return, 1)
			c.decls(s.Locals)
			c.stmts(s.Body, 1)
		}
	}
	c.stmts(prog.Body, 1)
	if c.failure != nil {
		return []diag.Diagnostic{*c.failure}
	}
	return nil
}

type limitChecker struct{ failure *diag.Diagnostic }

func (c *limitChecker) consts(decls []ConstDecl) {
	for _, decl := range decls {
		c.expr(decl.Value, 1)
	}
}

func (c *limitChecker) types(decls []TypeDecl) {
	for _, decl := range decls {
		c.typ(decl.Type, 1)
	}
}

func (c *limitChecker) enter(pos token.Pos, depth int) bool {
	if c.failure != nil {
		return false
	}
	if depth > MaxDepth {
		c.failure = &diag.Diagnostic{Code: diag.EResource, Pos: pos, End: pos + 1, Message: "AST traversal depth limit exceeded"}
		return false
	}
	return true
}

func (c *limitChecker) decls(decls []VarDecl) {
	for _, decl := range decls {
		c.typ(decl.Type, 1)
	}
}

func (c *limitChecker) params(params []Param) {
	for _, param := range params {
		c.typ(param.Type, 1)
	}
}

func (c *limitChecker) typ(t TypeSpec, depth int) {
	if !c.enter(t.At, depth) {
		return
	}
	for _, r := range t.Ranges {
		c.expr(r.Low, depth+1)
		c.expr(r.High, depth+1)
	}
	if t.Elem != nil {
		c.typ(*t.Elem, depth+1)
	}
	for _, field := range t.Fields {
		c.typ(field.Type, depth+1)
	}
}

func (c *limitChecker) stmts(stmts []Stmt, depth int) {
	for _, stmt := range stmts {
		if !c.enter(stmt.Start(), depth) {
			return
		}
		next := depth + 1
		switch s := stmt.(type) {
		case *AssignStmt:
			c.expr(s.Target, next)
			c.expr(s.Value, next)
		case *CallStmt:
			c.expr(s.Call, next)
		case *ColorStmt:
			c.expr(s.Color, next)
			c.expr(s.Target, next)
		case *IfStmt:
			c.expr(s.Cond, next)
			c.stmts(s.Then, next)
			c.stmts(s.Else, next)
		case *SwitchStmt:
			c.expr(s.X, next)
			for _, cc := range s.Cases {
				for _, label := range cc.Labels {
					c.expr(label.Low, next)
					c.expr(label.High, next)
				}
				c.stmts(cc.Body, next)
			}
			c.stmts(s.Default, next)
		case *WhileStmt:
			c.expr(s.Cond, next)
			c.stmts(s.Body, next)
		case *RepeatStmt:
			c.stmts(s.Body, next)
			c.expr(s.Cond, next)
		case *ForStmt:
			c.expr(s.From, next)
			c.expr(s.To, next)
			c.expr(s.Step, next)
			c.stmts(s.Body, next)
		case *ReturnStmt:
			c.expr(s.Value, next)
		case *ReadStmt:
			c.exprs(s.Targets, next)
		case *WriteStmt:
			for _, arg := range s.Args {
				c.expr(arg.Expr, next)
				c.expr(arg.Width, next)
				c.expr(arg.Decimals, next)
			}
		}
	}
}

func (c *limitChecker) exprs(exprs []Expr, depth int) {
	for _, expr := range exprs {
		c.expr(expr, depth)
	}
}

func (c *limitChecker) expr(expr Expr, depth int) {
	if expr == nil || !c.enter(expr.Start(), depth) {
		return
	}
	switch e := expr.(type) {
	case *UnaryExpr:
		c.expr(e.X, depth+1)
	case *BinaryExpr:
		c.expr(e.Left, depth+1)
		c.expr(e.Right, depth+1)
	case *IndexExpr:
		c.expr(e.X, depth+1)
		c.exprs(e.Indices, depth+1)
	case *FieldExpr:
		c.expr(e.X, depth+1)
	case *CallExpr:
		c.exprs(e.Args, depth+1)
	}
}
