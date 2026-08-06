package interp

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/stdlib"
)

// Interpreter is an instantiable tree-walking Portugol evaluator.
type Interpreter struct {
	in   *bufio.Scanner
	out  io.Writer
	lib  *stdlib.Library
	env  *env
	subs map[string]ast.Subprogram
}

// New creates an interpreter using in for leia and out for escreva/escreval.
func New(in io.Reader, out io.Writer) *Interpreter {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)
	return &Interpreter{
		in:   scanner,
		out:  out,
		lib:  stdlib.New(),
		subs: make(map[string]ast.Subprogram),
	}
}

// Run executes a checked program.
func (i *Interpreter) Run(prog *ast.Program) error {
	i.env = newEnv(nil)
	i.subs = make(map[string]ast.Subprogram)
	for _, decl := range prog.Globals {
		i.defineVars(decl)
	}
	for _, sub := range prog.Subs {
		i.subs[canon(sub.NameToken().Text)] = sub
	}
	ctrl, err := i.execStmts(prog.Body)
	if err != nil {
		return err
	}
	switch ctrl.kind {
	case noControl:
		return nil
	case breakControl:
		return fmt.Errorf("interrompa outside loop")
	case returnControl:
		return fmt.Errorf("retorne outside function")
	default:
		return nil
	}
}

func (i *Interpreter) defineVars(decl ast.VarDecl) {
	typ := runtime.TypeFromSpec(decl.Type)
	for _, name := range decl.Names {
		i.env.define(name.Text, typ)
	}
}

type controlKind int

const (
	noControl controlKind = iota
	breakControl
	returnControl
)

type control struct {
	kind  controlKind
	value runtime.Value
}
