package ast

import "github.com/ncode/portugol-go/internal/token"

// Expr is a Portugol expression.
type Expr interface {
	exprNode()
	Start() token.Pos
}

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

// CallExpr calls a function or built-in.
type CallExpr struct {
	Name token.Token
	Args []Expr
}

func (*CallExpr) exprNode()          {}
func (e *CallExpr) Start() token.Pos { return e.Name.Pos }
