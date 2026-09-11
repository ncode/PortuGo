package ast

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ncode/portugol-go/internal/token"
)

// Fprint writes a deterministic source-like representation of prog.
func Fprint(w io.Writer, prog *Program) error {
	if ds := CheckLimits(prog); len(ds) != 0 {
		return ds[0]
	}
	p := &printer{w: w}
	p.line(`algoritmo "%s"`, prog.Name)
	p.printConsole(prog.Console)
	p.printDecls(prog.Consts, prog.Types, prog.Globals)
	for _, sub := range prog.Subs {
		p.line("")
		p.printSub(sub)
	}
	p.line("inicio")
	p.indent++
	p.printStmts(prog.Body)
	p.indent--
	p.line("fimalgoritmo")
	return p.err
}

type printer struct {
	w      io.Writer
	indent int
	err    error
}

func (p *printer) line(format string, args ...any) {
	if p.err != nil {
		return
	}
	if format != "" {
		_, p.err = fmt.Fprint(p.w, strings.Repeat("  ", p.indent))
		if p.err != nil {
			return
		}
		_, p.err = fmt.Fprintf(p.w, format, args...)
		if p.err != nil {
			return
		}
	}
	_, p.err = fmt.Fprintln(p.w)
}

func (p *printer) printConsole(settings []ConsoleStmt) {
	for range settings {
		p.line("dos")
	}
}

func (p *printer) printDecls(consts []ConstDecl, types []TypeDecl, decls []VarDecl) {
	if len(consts) == 0 && len(types) == 0 && len(decls) == 0 {
		return
	}
	if len(consts) != 0 {
		p.line("const")
		p.indent++
		for _, d := range consts {
			p.line("%s = %s", d.Name.Text, exprString(d.Value))
		}
		p.indent--
	}
	if len(types) != 0 {
		p.line("tipo")
		p.indent++
		for _, d := range types {
			p.line("%s = %s", d.Name.Text, typeString(d.Type))
			if d.Type.Name == "registro" {
				p.indent++
				p.printVars(d.Type.Fields)
				p.indent--
				p.line("fimregistro")
			}
		}
		p.indent--
	}
	p.line("var")
	p.indent++
	p.printVars(decls)
	p.indent--
}

func (p *printer) printVars(decls []VarDecl) {
	for _, d := range decls {
		names := make([]string, len(d.Names))
		for i, name := range d.Names {
			names[i] = name.Text
		}
		p.line("%s: %s", strings.Join(names, ", "), typeString(d.Type))
	}
}

func (p *printer) printSub(sub Subprogram) {
	switch s := sub.(type) {
	case *ProcedureDecl:
		p.line("procedimento %s(%s)", s.Name.Text, paramsString(s.Params))
		p.printConsole(s.Console)
		p.printDecls(s.Consts, s.Types, s.Locals)
		p.line("inicio")
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.line("fimprocedimento")
	case *FunctionDecl:
		p.line("funcao %s(%s): %s", s.Name.Text, paramsString(s.Params), typeString(s.Return))
		p.printConsole(s.Console)
		p.printDecls(s.Consts, s.Types, s.Locals)
		p.line("inicio")
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.line("fimfuncao")
	}
}

func (p *printer) printStmts(stmts []Stmt) {
	for _, stmt := range stmts {
		p.printStmt(stmt)
	}
}

func (p *printer) printStmt(stmt Stmt) {
	switch s := stmt.(type) {
	case *ConsoleStmt:
		p.line("dos")
	case *EchoStmt:
		if s.Mode.Kind == token.IDENT {
			p.line("eco %s", strings.ToLower(s.Mode.Text))
		} else {
			p.line("eco")
		}
	case *ChronometerStmt:
		if s.Off {
			p.line("cronometro off")
		} else {
			p.line("cronometro on")
		}
	case *TimerStmt:
		p.line("timer %s", exprString(s.Value))
	case *PauseStmt:
		p.line("pausa")
	case *DebugStmt:
		p.line("debug %s", exprString(s.Cond))
	case *RandomInputStmt:
		if s.Off {
			p.line("aleatorio off")
		} else if len(s.Args) == 0 {
			p.line("aleatorio on")
		} else {
			args := make([]string, len(s.Args))
			for n, arg := range s.Args {
				args[n] = exprString(arg)
			}
			p.line("aleatorio %s", strings.Join(args, ", "))
		}
	case *AssignStmt:
		p.line("%s <- %s", exprString(s.Target), exprString(s.Value))
	case *CallStmt:
		p.line("%s", exprString(s.Call))
	case *IfStmt:
		p.line("se %s entao", exprString(s.Cond))
		p.indent++
		p.printStmts(s.Then)
		p.indent--
		if len(s.Else) > 0 {
			p.line("senao")
			p.indent++
			p.printStmts(s.Else)
			p.indent--
		}
		p.line("fimse")
	case *SwitchStmt:
		p.line("escolha %s", exprString(s.X))
		p.indent++
		for _, cc := range s.Cases {
			values := make([]string, len(cc.Labels))
			for i, label := range cc.Labels {
				values[i] = exprString(label.Low)
				if label.High != nil {
					values[i] += " ate " + exprString(label.High)
				}
			}
			p.line("caso %s:", strings.Join(values, ", "))
			p.indent++
			p.printStmts(cc.Body)
			p.indent--
		}
		if len(s.Default) > 0 {
			p.line("outrocaso:")
			p.indent++
			p.printStmts(s.Default)
			p.indent--
		}
		p.indent--
		p.line("fimescolha")
	case *WhileStmt:
		p.line("enquanto %s faca", exprString(s.Cond))
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.line("fimenquanto")
	case *RepeatStmt:
		p.line("repita")
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.line("ate %s", exprString(s.Cond))
	case *ForStmt:
		step := ""
		if s.Step != nil {
			step = " passo " + exprString(s.Step)
		}
		p.line("para %s de %s ate %s%s faca", s.Name.Text, exprString(s.From), exprString(s.To), step)
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.line("fimpara")
	case *BreakStmt:
		p.line("interrompa")
	case *ClearStmt:
		p.line("limpatela")
	case *ColorStmt:
		p.line("mudacor(%s, %s)", exprString(s.Color), exprString(s.Target))
	case *ReturnStmt:
		p.line("retorne %s", exprString(s.Value))
	case *ReadStmt:
		targets := make([]string, len(s.Targets))
		for i, t := range s.Targets {
			targets[i] = exprString(t)
		}
		p.line("leia(%s)", strings.Join(targets, ", "))
	case *WriteStmt:
		name := "escreva"
		if s.Newline {
			name = "escreval"
		}
		args := make([]string, len(s.Args))
		for i, arg := range s.Args {
			args[i] = writeArgString(arg)
		}
		p.line("%s(%s)", name, strings.Join(args, ", "))
	}
}

func typeString(t TypeSpec) string {
	if t.Name != "vetor" {
		return t.Name
	}
	ranges := make([]string, len(t.Ranges))
	for i, r := range t.Ranges {
		ranges[i] = exprString(r.Low) + ".." + exprString(r.High)
	}
	elem := "invalido"
	if t.Elem != nil {
		elem = typeString(*t.Elem)
	}
	return "vetor[" + strings.Join(ranges, ", ") + "] de " + elem
}

func paramsString(params []Param) string {
	parts := make([]string, len(params))
	for i, p := range params {
		prefix := ""
		if p.ByRef {
			prefix = "var "
		}
		parts[i] = fmt.Sprintf("%s%s: %s", prefix, p.Name.Text, typeString(p.Type))
	}
	return strings.Join(parts, "; ")
}

func writeArgString(arg WriteArg) string {
	out := exprString(arg.Expr)
	if arg.Width != nil {
		out += ":" + exprString(arg.Width)
		if arg.Decimals != nil {
			out += ":" + exprString(arg.Decimals)
		}
	}
	return out
}

func exprString(expr Expr) string {
	switch e := expr.(type) {
	case *NoValueExpr:
		return e.Keyword.Kind.String()
	case *LiteralExpr:
		switch e.Kind {
		case IntLiteral:
			return strconv.FormatInt(e.Int, 10)
		case RealLiteral:
			text := strconv.FormatFloat(e.Real, 'f', -1, 64)
			if !strings.ContainsRune(text, '.') {
				text += ".0"
			}
			return text
		case StringLiteral:
			return `"` + e.Str + `"`
		case BoolLiteral:
			if e.Bool {
				return "verdadeiro"
			}
			return "falso"
		}
	case *IdentExpr:
		return e.Name.Text
	case *FieldExpr:
		return exprString(e.X) + "." + e.Name.Text
	case *IndexExpr:
		indices := make([]string, len(e.Indices))
		for i, idx := range e.Indices {
			indices[i] = exprString(idx)
		}
		return exprString(e.X) + "[" + strings.Join(indices, ", ") + "]"
	case *UnaryExpr:
		return e.Op.Text + " " + operandString(e.X, e.Op.Kind.UnaryPrecedence())
	case *BinaryExpr:
		prec := e.Op.Kind.BinaryPrecedence()
		leftPrec := prec
		if e.IsComparison() {
			leftPrec++
		}
		return operandString(e.Left, leftPrec) + " " + e.Op.Text + " " + operandString(e.Right, prec+1)
	case *CallExpr:
		args := make([]string, len(e.Args))
		for i, arg := range e.Args {
			args[i] = exprString(arg)
		}
		return e.Name.Text + "(" + strings.Join(args, ", ") + ")"
	}
	return "<expr>"
}

func operandString(expr Expr, minPrec int) string {
	out := exprString(expr)
	prec := minPrec
	switch e := expr.(type) {
	case *UnaryExpr:
		prec = e.Op.Kind.UnaryPrecedence()
	case *BinaryExpr:
		prec = e.Op.Kind.BinaryPrecedence()
	}
	if prec < minPrec {
		return "(" + out + ")"
	}
	return out
}
