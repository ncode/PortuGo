package ast

import "github.com/ncode/portugol-go/internal/token"

// Program is a complete Portugol source file.
type Program struct {
	At      token.Pos
	Name    string
	Consts  []ConstDecl
	Globals []VarDecl
	Subs    []Subprogram
	Body    []Stmt
}

// ConstDecl binds an expression once when its declaration section is entered.
type ConstDecl struct {
	Name  token.Token
	Value Expr
}

// VarDecl declares one or more variables with the same type.
type VarDecl struct {
	At    token.Pos
	Names []token.Token
	Type  TypeSpec
}

// TypeSpec describes a declared Portugol type.
type TypeSpec struct {
	At     token.Pos
	Name   string
	Ranges []Range
	Elem   *TypeSpec
}

// Range is one vector dimension bound.
type Range struct {
	At   token.Pos
	Low  int64
	High int64
}

// Param is a procedure or function parameter.
type Param struct {
	At    token.Pos
	Name  token.Token
	Type  TypeSpec
	ByRef bool
}

// Subprogram is a top-level procedure or function declaration.
type Subprogram interface {
	subprogramNode()
	Start() token.Pos
	NameToken() token.Token
}

// ProcedureDecl declares a procedure.
type ProcedureDecl struct {
	At     token.Pos
	Name   token.Token
	Params []Param
	Consts []ConstDecl
	Locals []VarDecl
	Body   []Stmt
}

func (*ProcedureDecl) subprogramNode()          {}
func (d *ProcedureDecl) Start() token.Pos       { return d.At }
func (d *ProcedureDecl) NameToken() token.Token { return d.Name }

// FunctionDecl declares a function.
type FunctionDecl struct {
	At     token.Pos
	Name   token.Token
	Params []Param
	Return TypeSpec
	Consts []ConstDecl
	Locals []VarDecl
	Body   []Stmt
}

func (*FunctionDecl) subprogramNode()          {}
func (d *FunctionDecl) Start() token.Pos       { return d.At }
func (d *FunctionDecl) NameToken() token.Token { return d.Name }
