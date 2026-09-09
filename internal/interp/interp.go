package interp

import (
	"bufio"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/stdlib"
	"github.com/ncode/portugol-go/internal/token"
)

// Interpreter is an instantiable tree-walking Portugol evaluator.
type Interpreter struct {
	in           *bufio.Reader
	out          io.Writer
	writeNewline bool
	writeBytes   int
	lib          *stdlib.Library
	env          *env
	subs         map[token.Pos]ast.Subprogram
	global       *env
	info         *sema.Info
	program      *ast.Program
	options      Options
	initErr      error
	steps        uint64
	depth, calls int
	result       *runtime.Cell
	results      []runtime.Value
}

// New resolves defaults once and owns the buffered input for subsequent runs.
// Passing an existing buffered reader shares its state with a REPL safely.
func New(options Options) *Interpreter {
	if options.Input == nil {
		options.Input = strings.NewReader("")
	}
	if options.Output == nil {
		options.Output = io.Discard
	}
	if options.Host == nil {
		options.Host = HeadlessHost{}
	}
	if options.Random == nil {
		options.Random = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}
	var err error
	if options.WorkingDir == "" {
		options.WorkingDir, err = os.Getwd()
	}
	return &Interpreter{
		in:      bufio.NewReader(options.Input),
		out:     options.Output,
		lib:     stdlib.New(options.Random),
		options: options,
		initErr: err,
	}
}

// Run executes a checked program.
func (i *Interpreter) Run(prog *ast.Program, info *sema.Info) []diag.Diagnostic {
	i.writeNewline = false
	i.writeBytes = 0
	i.program = nil
	i.global = nil
	i.result = nil
	i.results = nil
	pos := token.NoPos
	if prog != nil {
		pos = prog.At
	}
	if !info.ValidFor(prog) {
		return []diag.Diagnostic{{Code: diag.RType, Pos: pos, Message: "missing or inconsistent semantic information"}}
	}
	if i.initErr != nil {
		return []diag.Diagnostic{{Code: diag.RHost, Pos: pos, Message: "cannot resolve working directory", Cause: i.initErr}}
	}
	i.info = info
	i.program = prog
	i.steps, i.depth, i.calls = 0, 0, 0
	i.env = newEnv(nil)
	i.global = i.env
	i.subs = make(map[token.Pos]ast.Subprogram)
	if err := i.defineConsts(prog.Consts); err != nil {
		return diagnostics(err, pos, diag.RType)
	}
	for _, decl := range prog.Globals {
		if err := i.defineVars(decl); err != nil {
			return diagnostics(err, decl.At, diag.RStorage)
		}
	}
	for _, sub := range prog.Subs {
		b, ok := info.Binding(sub.NameToken())
		if !ok {
			return diagnostics(fmt.Errorf("missing subprogram binding"), sub.Start(), diag.RType)
		}
		i.subs[b.ID] = sub
	}
	ctrl, err := i.execStmts(prog.Body)
	if err != nil {
		return diagnostics(err, pos, diag.RType)
	}
	switch ctrl.kind {
	case noControl:
		return nil
	case breakControl:
		return diagnostics(fmt.Errorf("interrompa outside loop"), ctrl.at, diag.RLoop)
	default:
		return nil
	}
}

// State returns independent copies of global variables after the most recent run.
func (i *Interpreter) State() map[string]runtime.Value {
	values := make(map[string]runtime.Value)
	if i.program == nil || i.global == nil {
		return values
	}
	for _, decl := range i.program.Globals {
		for _, name := range decl.Names {
			b, ok := i.info.Binding(name)
			if !ok {
				continue
			}
			if cell, ok := i.global.lookup(b.ID); ok {
				values[b.Name] = runtime.Clone(cell.Value)
			}
		}
	}
	return values
}

func (i *Interpreter) defineVars(decl ast.VarDecl) error {
	for _, name := range decl.Names {
		b, ok := i.info.Binding(name)
		if !ok {
			return failure(name.Pos, diag.RType, fmt.Errorf("missing declaration layout"))
		}
		if slots, err := b.Type.Slots(); err != nil || slots != b.Slots {
			return failure(name.Pos, diag.RStorage, fmt.Errorf("inconsistent storage layout"))
		}
		i.env.define(b.ID, b.Type)
	}
	return nil
}

func (i *Interpreter) defineConsts(decls []ast.ConstDecl) error {
	for _, decl := range decls {
		b, ok := i.info.Binding(decl.Name)
		if !ok {
			return failure(decl.Name.Pos, diag.RType, fmt.Errorf("missing constant binding"))
		}
		value, err := i.eval(decl.Value)
		if err != nil {
			return err
		}
		i.env.cells[b.ID] = &runtime.Cell{Type: value.Type(), Value: runtime.Clone(value)}
	}
	return nil
}

type controlKind int

const (
	noControl controlKind = iota
	breakControl
)

type control struct {
	kind controlKind
	at   token.Pos
}
