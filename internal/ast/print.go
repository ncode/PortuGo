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
	p := &printer{w: w, comments: prog.Comments, fragments: prog.Fragments}
	p.before(prog.At, 0)
	p.line(`algoritmo "%s"`, prog.Name)
	p.printStmts(prog.Config)
	p.printDecls(prog.Sections, prog.Consts, prog.Types, prog.Globals)
	for _, sub := range prog.Subs {
		for len(p.comments) > 0 && p.comments[0].Inline && !p.comments[0].Semicolon && p.comments[0].Pos < sub.Start() && p.hasPending {
			p.pending += " " + p.comments[0].Text
			p.comments = p.comments[1:]
		}
		p.line("")
		p.printSub(sub)
	}
	sections := prog.Sections
	if len(prog.Subs) != 0 {
		sections = DeclSections{}
	}
	p.begin(prog.Begin, sections)
	p.indent++
	p.printStmts(prog.Body)
	p.indent--
	p.before(prog.End, 1)
	if p.err == nil {
		suffix := strings.ReplaceAll(prog.Suffix.Text, "\r\n", "\n")
		if suffix == "" {
			suffix = "\n"
		}
		_, p.err = fmt.Fprintf(w, "fimalgoritmo%s", suffix)
	}
	return p.err
}

type printer struct {
	w          io.Writer
	indent     int
	err        error
	comments   []CommentGroup
	pending    string
	hasPending bool
	fragments  map[token.Pos]CommentFragment
}

func (p *printer) line(format string, args ...any) {
	p.flush()
	if p.err != nil {
		return
	}
	if format != "" {
		p.pending = strings.Repeat("  ", p.indent) + fmt.Sprintf(format, args...)
	}
	p.hasPending = true
}

func (p *printer) printDecls(sections DeclSections, consts []ConstDecl, types []TypeDecl, decls []VarDecl) {
	if len(consts) == 0 && len(types) == 0 && len(decls) == 0 && sections.Var.Kind != token.VAR && sections.Type.Kind != token.TIPO && sections.Const.Kind != token.CONST {
		return
	}
	if len(consts) != 0 || sections.Const.Kind == token.CONST {
		p.section(sections.Const, "const")
		p.indent++
		for _, d := range consts {
			p.before(d.Name.Pos, p.indent)
			p.lineFrom(d.Name.Pos, "%s = %s", d.Name.Text, exprString(d.Value))
		}
		p.indent--
	}
	if len(types) != 0 || sections.Type.Kind == token.TIPO {
		p.section(sections.Type, "tipo")
		p.indent++
		for _, d := range types {
			p.before(d.Name.Pos, p.indent)
			p.lineFrom(d.Name.Pos, "%s = %s", d.Name.Text, typeString(d.Type))
			if d.Type.Name == "registro" {
				p.indent++
				p.printVars(d.Type.Fields)
				p.indent--
				p.end(d.Type.End, "fimregistro")
			}
		}
		p.indent--
	}
	p.section(sections.Var, "var")
	p.indent++
	p.printVars(decls)
	p.indent--
}

func (p *printer) printVars(decls []VarDecl) {
	for _, d := range decls {
		p.before(d.At, p.indent)
		names := make([]string, len(d.Names))
		for i, name := range d.Names {
			names[i] = name.Text
		}
		p.lineFrom(d.At, "%s: %s", strings.Join(names, ", "), typeString(d.Type))
	}
}

func (p *printer) printSub(sub Subprogram) {
	p.before(sub.Start(), p.indent)
	switch s := sub.(type) {
	case *ProcedureDecl:
		p.lineFrom(s.At, "procedimento %s(%s)", s.Name.Text, paramsString(s.Params))
		p.printStmts(s.Config)
		p.printDecls(s.Sections, s.Consts, s.Types, s.Locals)
		p.begin(s.Begin, s.Sections)
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.end(s.End, "fimprocedimento")
	case *FunctionDecl:
		p.lineFrom(s.At, "funcao %s(%s): %s", s.Name.Text, paramsString(s.Params), typeString(s.Return))
		p.printStmts(s.Config)
		p.printDecls(s.Sections, s.Consts, s.Types, s.Locals)
		p.begin(s.Begin, s.Sections)
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.end(s.End, "fimfuncao")
	}
}

func (p *printer) printStmts(stmts []Stmt) {
	for _, stmt := range stmts {
		p.printStmt(stmt)
	}
}

func (p *printer) printStmt(stmt Stmt) {
	p.before(stmt.Start(), p.indent)
	switch stmt.(type) {
	case *IfStmt, *SwitchStmt, *WhileStmt, *RepeatStmt, *ForStmt:
	default:
		if p.fragment(stmt.Start()) {
			return
		}
	}
	switch s := stmt.(type) {
	case *FileInputStmt:
		p.line("arquivo \"%s\"", s.Path)
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
		p.lineFrom(s.At, "se %s entao", exprString(s.Cond))
		p.indent++
		p.printStmts(s.Then)
		p.indent--
		if len(s.Else) > 0 || s.ElseAt != 0 {
			p.end(s.ElseAt, "senao")
			p.indent++
			p.printStmts(s.Else)
			p.indent--
		}
		p.end(s.End, "fimse")
	case *SwitchStmt:
		p.lineFrom(s.At, "escolha %s", exprString(s.X))
		p.indent++
		for _, cc := range s.Cases {
			p.before(cc.At, p.indent+1)
			values := make([]string, len(cc.Labels))
			for i, label := range cc.Labels {
				values[i] = exprString(label.Low)
				if label.High != nil {
					values[i] += " ate " + exprString(label.High)
				}
			}
			p.lineFrom(cc.At, "caso %s:", strings.Join(values, ", "))
			p.indent++
			p.printStmts(cc.Body)
			p.indent--
		}
		if len(s.Default) > 0 || s.DefaultAt != 0 {
			p.end(s.DefaultAt, "outrocaso:")
			p.indent++
			p.printStmts(s.Default)
			p.indent--
		}
		p.indent--
		p.before(s.End, p.indent+2)
		p.line("fimescolha")
	case *WhileStmt:
		p.lineFrom(s.At, "enquanto %s faca", exprString(s.Cond))
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.end(s.End, "fimenquanto")
	case *RepeatStmt:
		p.line("repita")
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.end(s.End, "ate %s", exprString(s.Cond))
	case *ForStmt:
		step := ""
		if s.Step != nil {
			step = " passo " + exprString(s.Step)
		}
		p.lineFrom(s.At, "para %s de %s ate %s%s faca", s.Name.Text, exprString(s.From), exprString(s.To), step)
		p.indent++
		p.printStmts(s.Body)
		p.indent--
		p.end(s.End, "fimpara")
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
