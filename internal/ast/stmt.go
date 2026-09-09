package ast

import "github.com/ncode/portugol-go/internal/token"

// Stmt is a Portugol statement.
type Stmt interface {
	stmtNode()
	Start() token.Pos
}

// AssignStmt assigns to a variable or indexed vector element.
type AssignStmt struct {
	At     token.Pos
	Target Expr
	Value  Expr
}

func (*AssignStmt) stmtNode()          {}
func (s *AssignStmt) Start() token.Pos { return s.At }

// CallStmt invokes a procedure.
type CallStmt struct {
	Call *CallExpr
}

func (*CallStmt) stmtNode()          {}
func (s *CallStmt) Start() token.Pos { return s.Call.Start() }

// IfStmt is a conditional branch.
type IfStmt struct {
	At   token.Pos
	Cond Expr
	Then []Stmt
	Else []Stmt
}

func (*IfStmt) stmtNode()          {}
func (s *IfStmt) Start() token.Pos { return s.At }

// CaseLabel is a single escolha value or an inclusive range when High is non-nil.
type CaseLabel struct {
	Low  Expr
	High Expr
}

// Start returns the label's source position.
func (l CaseLabel) Start() token.Pos { return l.Low.Start() }

// CaseClause is one escolha branch.
type CaseClause struct {
	At     token.Pos
	Labels []CaseLabel
	Body   []Stmt
}

// SwitchStmt is an escolha statement.
type SwitchStmt struct {
	At      token.Pos
	X       Expr
	Cases   []CaseClause
	Default []Stmt
}

func (*SwitchStmt) stmtNode()          {}
func (s *SwitchStmt) Start() token.Pos { return s.At }

// WhileStmt is an enquanto loop.
type WhileStmt struct {
	At   token.Pos
	Cond Expr
	Body []Stmt
}

func (*WhileStmt) stmtNode()          {}
func (s *WhileStmt) Start() token.Pos { return s.At }

// RepeatStmt is a repita loop.
type RepeatStmt struct {
	At   token.Pos
	Body []Stmt
	Cond Expr
}

func (*RepeatStmt) stmtNode()          {}
func (s *RepeatStmt) Start() token.Pos { return s.At }

// ForStmt is a para loop.
type ForStmt struct {
	At   token.Pos
	Name token.Token
	From Expr
	To   Expr
	Step Expr
	Body []Stmt
}

func (*ForStmt) stmtNode()          {}
func (s *ForStmt) Start() token.Pos { return s.At }

// BreakStmt exits the innermost loop.
type BreakStmt struct {
	At token.Pos
}

func (*BreakStmt) stmtNode()          {}
func (s *BreakStmt) Start() token.Pos { return s.At }

// ReturnStmt sets the function result. Value is nil when the expression is missing.
type ReturnStmt struct {
	At    token.Pos
	Value Expr
}

func (*ReturnStmt) stmtNode()          {}
func (s *ReturnStmt) Start() token.Pos { return s.At }

// ReadStmt reads values into destinations.
type ReadStmt struct {
	At      token.Pos
	Targets []Expr
}

func (*ReadStmt) stmtNode()          {}
func (s *ReadStmt) Start() token.Pos { return s.At }

// WriteArg is one escreva/escreval argument.
type WriteArg struct {
	Expr     Expr
	Width    Expr
	Decimals Expr
}

// WriteStmt writes expressions to stdout.
type WriteStmt struct {
	At      token.Pos
	Newline bool
	Args    []WriteArg
}

func (*WriteStmt) stmtNode()          {}
func (s *WriteStmt) Start() token.Pos { return s.At }
