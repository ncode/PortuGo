package ast

import "github.com/ncode/portugol-go/internal/token"

// Program is a complete Portugol source file.
type Program struct {
	Sections   DeclSections
	Begin, End token.Pos
	Comments   []CommentGroup
	Fragments  map[token.Pos]CommentFragment
	At         token.Pos
	Name       string
	Config     []Stmt
	Consts     []ConstDecl
	Types      []TypeDecl
	Globals    []VarDecl
	Subs       []Subprogram
	Body       []Stmt
	Suffix     token.Token // Opaque decoded source immediately after fimalgoritmo.
}

// DeclSections retains the introducing tokens, including empty sections.
type DeclSections struct {
	Const, Type, Var token.Token
}

// CommentGroup is an ordered physical-line comment with its source span.
// Inline groups follow code on the same line; other groups precede the next
// construct or closing delimiter. Text excludes indentation and the newline.
type CommentGroup struct {
	Pos, End  token.Pos
	Text      string
	Inline    bool
	Semicolon bool // The comment follows an otherwise empty separator line.
}

// CommentFragment retains a construct whose internal comments require physical
// line breaks, such as a multiline expression or parameter list.
type CommentFragment struct {
	End  token.Pos
	Text string
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
	End    token.Pos
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
	Sections   DeclSections
	Begin, End token.Pos
	At         token.Pos
	Name       token.Token
	Params     []Param
	Config     []Stmt
	Consts     []ConstDecl
	Types      []TypeDecl
	Locals     []VarDecl
	Body       []Stmt
}

func (*ProcedureDecl) subprogramNode()          {}
func (d *ProcedureDecl) Start() token.Pos       { return d.At }
func (d *ProcedureDecl) NameToken() token.Token { return d.Name }

// FunctionDecl declares a function.
type FunctionDecl struct {
	Sections   DeclSections
	Begin, End token.Pos
	At         token.Pos
	Name       token.Token
	Params     []Param
	Return     TypeSpec
	Config     []Stmt
	Consts     []ConstDecl
	Types      []TypeDecl
	Locals     []VarDecl
	Body       []Stmt
}

func (*FunctionDecl) subprogramNode()          {}
func (d *FunctionDecl) Start() token.Pos       { return d.At }
func (d *FunctionDecl) NameToken() token.Token { return d.Name }
