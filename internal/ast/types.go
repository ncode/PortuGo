package ast

import "github.com/ncode/portugol-go/internal/token"

// Program is a complete Portugol source file.
type Program struct {
	At      token.Pos
	Name    string
	Console []ConsoleStmt
	Consts  []ConstDecl
	Types   []TypeDecl
	Globals []VarDecl
	Subs    []Subprogram
	Body    []Stmt
}

// ConstDecl binds an expression once when its declaration section is entered.
type ConstDecl struct {
	Name  token.Token
	Value Expr
}

// TypeDecl defines a record or gives a name to an earlier type.
type TypeDecl struct {
	Name token.Token
	Type TypeSpec
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
	Fields []VarDecl
}

// Range is one vector dimension bound.
type Range struct {
	At   token.Pos
	Low  Expr
	High Expr
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
	At      token.Pos
	Name    token.Token
	Params  []Param
	Console []ConsoleStmt
	Consts  []ConstDecl
	Types   []TypeDecl
	Locals  []VarDecl
	Body    []Stmt
}

func (*ProcedureDecl) subprogramNode()          {}
func (d *ProcedureDecl) Start() token.Pos       { return d.At }
func (d *ProcedureDecl) NameToken() token.Token { return d.Name }

// FunctionDecl declares a function.
type FunctionDecl struct {
	At      token.Pos
	Name    token.Token
	Params  []Param
	Return  TypeSpec
	Console []ConsoleStmt
	Consts  []ConstDecl
	Types   []TypeDecl
	Locals  []VarDecl
	Body    []Stmt
}

func (*FunctionDecl) subprogramNode()          {}
func (d *FunctionDecl) Start() token.Pos       { return d.At }
func (d *FunctionDecl) NameToken() token.Token { return d.Name }
