package ast

import "github.com/ncode/portugol-go/internal/token"

// Expr is a Portugol expression.
type Expr interface {
	exprNode()
	Start() token.Pos
}

// NoValueExpr is a display keyword used as an expression, without a host effect.
type NoValueExpr struct{ Keyword token.Token }

func (*NoValueExpr) exprNode()          {}
func (e *NoValueExpr) Start() token.Pos { return e.Keyword.Pos }

// LiteralKind identifies the concrete literal value field.
type LiteralKind int

const (
	IntLiteral LiteralKind = iota
	RealLiteral
	StringLiteral
	BoolLiteral
)

// LiteralExpr is a scalar literal.
type LiteralExpr struct {
	At   token.Pos
	Kind LiteralKind
	Int  int64
	Real float64
	Str  string
	Bool bool
}

func (*LiteralExpr) exprNode()          {}
func (e *LiteralExpr) Start() token.Pos { return e.At }

// IdentExpr refers to a variable, function, or built-in by name.
type IdentExpr struct {
	Name token.Token
}

func (*IdentExpr) exprNode()          {}
func (e *IdentExpr) Start() token.Pos { return e.Name.Pos }

// IndexExpr indexes a vector value.
type IndexExpr struct {
	At      token.Pos
	X       Expr
	Indices []Expr
}

func (*IndexExpr) exprNode()          {}
func (e *IndexExpr) Start() token.Pos { return e.At }

// FieldExpr selects a named record field.
type FieldExpr struct {
	At   token.Pos
	X    Expr
	Name token.Token
}

func (*FieldExpr) exprNode()          {}
func (e *FieldExpr) Start() token.Pos { return e.At }

// UnaryExpr applies one unary operator.
type UnaryExpr struct {
	Op token.Token
	X  Expr
}

func (*UnaryExpr) exprNode()          {}
func (e *UnaryExpr) Start() token.Pos { return e.Op.Pos }

// BinaryExpr applies one binary operator.
type BinaryExpr struct {
	Op    token.Token
	Left  Expr
	Right Expr
}

func (*BinaryExpr) exprNode()          {}
func (e *BinaryExpr) Start() token.Pos { return e.Op.Pos }

// IsComparison reports whether the operator has the logical result category.
func (e *BinaryExpr) IsComparison() bool {
	switch e.Op.Kind {
	case token.EQL, token.NEQ, token.LSS, token.GTR, token.LEQ, token.GEQ:
		return true
	}
	return false
}

// CallExpr calls a function or built-in.
type CallExpr struct {
	Name token.Token
	Args []Expr
	Bare bool // Statement call written without parentheses.
}

func (*CallExpr) exprNode()          {}
func (e *CallExpr) Start() token.Pos { return e.Name.Pos }
