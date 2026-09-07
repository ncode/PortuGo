package parser

import (
	"strconv"
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

// Parse builds an AST from a token stream.
func Parse(tokens []token.Token) (*ast.Program, []diag.Diagnostic) {
	p := &parser{tokens: tokens}
	prog := p.parseProgram()
	if p.limited {
		return nil, p.diags
	}
	if ds := ast.CheckLimits(prog); len(ds) != 0 {
		return nil, diag.Ordered(append(p.diags, ds...))
	}
	return prog, p.diags
}

type parser struct {
	tokens  []token.Token
	pos     int
	diags   []diag.Diagnostic
	depth   int
	limited bool
}

func (p *parser) parseProgram() *ast.Program {
	start := p.expect(token.ALGORITMO, "expected algoritmo")
	name := ""
	if p.peek().Kind == token.STRING {
		name = p.advance().Text
	} else {
		p.error(p.peek(), "expected algorithm name string")
	}
	prog := &ast.Program{At: start.Pos, Name: name}
	if p.peek().Kind == token.VAR {
		prog.Globals = p.parseVarBlock()
	}
	for p.peek().Kind == token.PROCEDIMENTO || p.peek().Kind == token.FUNCAO {
		prog.Subs = append(prog.Subs, p.parseSubprogram())
	}
	p.expect(token.INICIO, "expected inicio")
	prog.Body = p.parseStmtList(stopSet(token.FIMALGORITMO))
	p.expect(token.FIMALGORITMO, "expected fimalgoritmo")
	return prog
}

func (p *parser) parseVarBlock() []ast.VarDecl {
	p.expect(token.VAR, "expected var")
	var decls []ast.VarDecl
	for p.peek().Kind == token.IDENT {
		decls = append(decls, p.parseVarDecl())
	}
	return decls
}

func (p *parser) parseVarDecl() ast.VarDecl {
	first := p.expect(token.IDENT, "expected identifier")
	names := []token.Token{first}
	for p.match(token.COMMA) {
		names = append(names, p.expect(token.IDENT, "expected identifier"))
	}
	p.expect(token.COLON, "expected ':' after variable name")
	typ := p.parseType()
	return ast.VarDecl{At: first.Pos, Names: names, Type: typ}
}

func (p *parser) parseType() ast.TypeSpec {
	if !p.enter() {
		return ast.TypeSpec{Name: "inteiro"}
	}
	defer func() { p.depth-- }()
	tok := p.peek()
	switch tok.Kind {
	case token.INTEIRO, token.REAL, token.CARACTERE, token.LOGICO:
		p.advance()
		return ast.TypeSpec{At: tok.Pos, Name: strings.ToLower(tok.Text)}
	case token.VETOR:
		p.advance()
		p.expect(token.LBRACK, "expected '[' after vetor")
		var ranges []ast.Range
		for {
			at := p.peek().Pos
			low := p.parseBoundInt()
			p.expect(token.DOTDOT, "expected '..' in vector bound")
			high := p.parseBoundInt()
			ranges = append(ranges, ast.Range{At: at, Low: low, High: high})
			if !p.match(token.COMMA) {
				break
			}
		}
		p.expect(token.RBRACK, "expected ']' after vector bounds")
		p.expect(token.DE, "expected de after vector bounds")
		elem := p.parseType()
		return ast.TypeSpec{At: tok.Pos, Name: "vetor", Ranges: ranges, Elem: &elem}
	default:
		p.error(tok, "expected type")
		p.advance()
		return ast.TypeSpec{At: tok.Pos, Name: "inteiro"}
	}
}

func (p *parser) parseBoundInt() int64 {
	neg := p.match(token.SUB)
	tok := p.expect(token.NUMBER, "expected integer bound")
	v, err := strconv.ParseInt(tok.Text, 10, 64)
	if err != nil {
		p.error(tok, "expected integer bound")
		return 0
	}
	if neg {
		return -v
	}
	return v
}

func (p *parser) parseSubprogram() ast.Subprogram {
	if p.peek().Kind == token.PROCEDIMENTO {
		return p.parseProcedure()
	}
	return p.parseFunction()
}

func (p *parser) parseProcedure() *ast.ProcedureDecl {
	start := p.expect(token.PROCEDIMENTO, "expected procedimento")
	name := p.expect(token.IDENT, "expected procedure name")
	params := p.parseParamList()
	decl := &ast.ProcedureDecl{At: start.Pos, Name: name, Params: params}
	if p.peek().Kind == token.VAR {
		decl.Locals = p.parseVarBlock()
	}
	p.expect(token.INICIO, "expected inicio in procedure")
	decl.Body = p.parseStmtList(stopSet(token.FIMPROCEDIMENTO))
	p.expect(token.FIMPROCEDIMENTO, "expected fimprocedimento")
	return decl
}

func (p *parser) parseFunction() *ast.FunctionDecl {
	start := p.expect(token.FUNCAO, "expected funcao")
	name := p.expect(token.IDENT, "expected function name")
	params := p.parseParamList()
	p.expect(token.COLON, "expected ':' before function return type")
	ret := p.parseType()
	decl := &ast.FunctionDecl{At: start.Pos, Name: name, Params: params, Return: ret}
	if p.peek().Kind == token.VAR {
		decl.Locals = p.parseVarBlock()
	}
	p.expect(token.INICIO, "expected inicio in function")
	decl.Body = p.parseStmtList(stopSet(token.FIMFUNCAO))
	p.expect(token.FIMFUNCAO, "expected fimfuncao")
	return decl
}

func (p *parser) parseParamList() []ast.Param {
	p.expect(token.LPAREN, "expected '('")
	if p.match(token.RPAREN) {
		return nil
	}
	var params []ast.Param
	for {
		byRef := p.match(token.VAR)
		first := p.expect(token.IDENT, "expected parameter name")
		names := []token.Token{first}
		for p.match(token.COMMA) {
			names = append(names, p.expect(token.IDENT, "expected parameter name"))
		}
		p.expect(token.COLON, "expected ':' after parameter name")
		typ := p.parseType()
		for _, name := range names {
			params = append(params, ast.Param{At: name.Pos, Name: name, Type: typ, ByRef: byRef})
		}
		if !p.match(token.SEMI) {
			break
		}
	}
	p.expect(token.RPAREN, "expected ')'")
	return params
}

func (p *parser) parseStmtList(stops map[token.Kind]bool) []ast.Stmt {
	var stmts []ast.Stmt
	for !stops[p.peek().Kind] && p.peek().Kind != token.EOF {
		if p.match(token.SEMI) {
			continue
		}
		stmt := p.parseStmt()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}
	return stmts
}

func (p *parser) parseStmt() ast.Stmt {
	if !p.enter() {
		return nil
	}
	defer func() { p.depth-- }()
	switch p.peek().Kind {
	case token.IDENT:
		return p.parseIdentStmt()
	case token.SE:
		return p.parseIf()
	case token.ESCOLHA:
		return p.parseSwitch()
	case token.ENQUANTO:
		return p.parseWhile()
	case token.REPITA:
		return p.parseRepeat()
	case token.PARA:
		return p.parseFor()
	case token.INTERROMPA:
		return &ast.BreakStmt{At: p.advance().Pos}
	case token.RETORNE:
		tok := p.advance()
		return &ast.ReturnStmt{At: tok.Pos, Value: p.parseExpr(0)}
	case token.LEIA:
		return p.parseRead()
	case token.ESCREVA, token.ESCREVAL:
		return p.parseWrite()
	default:
		p.error(p.peek(), "expected statement")
		p.advance()
		return nil
	}
}

func (p *parser) parseIdentStmt() ast.Stmt {
	if p.peekN(1).Kind == token.LPAREN {
		call := p.parseCall()
		return &ast.CallStmt{Call: call}
	}
	target := p.parseDesignator()
	at := target.Start()
	p.expect(token.ASSIGN, "expected '<-' in assignment")
	value := p.parseExpr(0)
	return &ast.AssignStmt{At: at, Target: target, Value: value}
}

func (p *parser) parseIf() ast.Stmt {
	start := p.expect(token.SE, "expected se")
	cond := p.parseExpr(0)
	p.expect(token.ENTAO, "expected entao")
	thenBody := p.parseStmtList(stopSet(token.SENAO, token.FIMSE))
	var elseBody []ast.Stmt
	if p.match(token.SENAO) {
		elseBody = p.parseStmtList(stopSet(token.FIMSE))
	}
	p.expect(token.FIMSE, "expected fimse")
	return &ast.IfStmt{At: start.Pos, Cond: cond, Then: thenBody, Else: elseBody}
}

func (p *parser) parseSwitch() ast.Stmt {
	start := p.expect(token.ESCOLHA, "expected escolha")
	x := p.parseExpr(0)
	sw := &ast.SwitchStmt{At: start.Pos, X: x}
	for p.peek().Kind != token.FIMESCOLHA && p.peek().Kind != token.EOF {
		switch p.peek().Kind {
		case token.CASO:
			cstart := p.advance()
			values := []ast.Expr{p.parseExpr(0)}
			for p.match(token.COMMA) {
				values = append(values, p.parseExpr(0))
			}
			p.expect(token.COLON, "expected ':' after caso")
			body := p.parseStmtList(stopSet(token.CASO, token.OUTROCASO, token.FIMESCOLHA))
			sw.Cases = append(sw.Cases, ast.CaseClause{At: cstart.Pos, Values: values, Body: body})
		case token.OUTROCASO:
			p.advance()
			p.match(token.COLON)
			sw.Default = p.parseStmtList(stopSet(token.FIMESCOLHA))
		default:
			p.error(p.peek(), "expected caso or fimescolha")
			p.advance()
		}
	}
	p.expect(token.FIMESCOLHA, "expected fimescolha")
	return sw
}

func (p *parser) parseWhile() ast.Stmt {
	start := p.expect(token.ENQUANTO, "expected enquanto")
	cond := p.parseExpr(0)
	p.expect(token.FACA, "expected faca")
	body := p.parseStmtList(stopSet(token.FIMENQUANTO))
	p.expect(token.FIMENQUANTO, "expected fimenquanto")
	return &ast.WhileStmt{At: start.Pos, Cond: cond, Body: body}
}

func (p *parser) parseRepeat() ast.Stmt {
	start := p.expect(token.REPITA, "expected repita")
	body := p.parseStmtList(stopSet(token.ATE))
	p.expect(token.ATE, "expected ate")
	cond := p.parseExpr(0)
	return &ast.RepeatStmt{At: start.Pos, Body: body, Cond: cond}
}

func (p *parser) parseFor() ast.Stmt {
	start := p.expect(token.PARA, "expected para")
	name := p.expect(token.IDENT, "expected loop variable")
	p.expect(token.DE, "expected de")
	from := p.parseExpr(0)
	p.expect(token.ATE, "expected ate")
	to := p.parseExpr(0)
	var step ast.Expr
	if p.match(token.PASSO) {
		step = p.parseExpr(0)
	}
	p.expect(token.FACA, "expected faca")
	body := p.parseStmtList(stopSet(token.FIMPARA))
	p.expect(token.FIMPARA, "expected fimpara")
	return &ast.ForStmt{At: start.Pos, Name: name, From: from, To: to, Step: step, Body: body}
}

func (p *parser) parseRead() ast.Stmt {
	start := p.expect(token.LEIA, "expected leia")
	p.expect(token.LPAREN, "expected '(' after leia")
	var targets []ast.Expr
	if !p.match(token.RPAREN) {
		for {
			targets = append(targets, p.parseDesignator())
			if !p.match(token.COMMA) {
				break
			}
		}
		p.expect(token.RPAREN, "expected ')'")
	}
	return &ast.ReadStmt{At: start.Pos, Targets: targets}
}

func (p *parser) parseWrite() ast.Stmt {
	start := p.advance()
	stmt := &ast.WriteStmt{At: start.Pos, Newline: start.Kind == token.ESCREVAL}
	p.expect(token.LPAREN, "expected '(' after write")
	if !p.match(token.RPAREN) {
		for {
			arg := ast.WriteArg{Expr: p.parseExpr(0)}
			if p.match(token.COLON) {
				arg.Width = p.parseExpr(0)
				if p.match(token.COLON) {
					arg.Decimals = p.parseExpr(0)
				}
			}
			stmt.Args = append(stmt.Args, arg)
			if !p.match(token.COMMA) {
				break
			}
		}
		p.expect(token.RPAREN, "expected ')'")
	}
	return stmt
}

func (p *parser) parseDesignator() ast.Expr {
	name := p.expect(token.IDENT, "expected identifier")
	var expr ast.Expr = &ast.IdentExpr{Name: name}
	for p.match(token.LBRACK) {
		at := expr.Start()
		var indices []ast.Expr
		if !p.match(token.RBRACK) {
			for {
				indices = append(indices, p.parseExpr(0))
				if !p.match(token.COMMA) {
					break
				}
			}
			p.expect(token.RBRACK, "expected ']'")
		}
		expr = &ast.IndexExpr{At: at, X: expr, Indices: indices}
	}
	return expr
}

func (p *parser) expect(kind token.Kind, msg string) token.Token {
	if p.peek().Kind == kind {
		return p.advance()
	}
	tok := p.peek()
	p.error(tok, msg)
	return token.Token{Kind: kind, Pos: tok.Pos}
}

func (p *parser) match(kind token.Kind) bool {
	if p.peek().Kind != kind {
		return false
	}
	p.advance()
	return true
}

func (p *parser) advance() token.Token {
	tok := p.peek()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

func (p *parser) peek() token.Token {
	return p.peekN(0)
}

func (p *parser) peekN(n int) token.Token {
	i := p.pos + n
	if i >= 0 && i < len(p.tokens) {
		return p.tokens[i]
	}
	return token.Token{Kind: token.EOF}
}

func (p *parser) error(tok token.Token, msg string) {
	if p.limited {
		return
	}
	p.diags = append(p.diags, diag.Diagnostic{Code: diag.EParse, Pos: tok.Pos, Message: msg})
}

func (p *parser) enter() bool {
	if p.limited {
		return false
	}
	if p.depth == ast.MaxDepth {
		pos := p.peek().Pos
		p.diags = append(p.diags, diag.Diagnostic{Code: diag.EResource, Pos: pos, End: pos + 1, Message: "syntax nesting limit exceeded"})
		p.limited = true
		p.pos = len(p.tokens)
		return false
	}
	p.depth++
	return true
}

func stopSet(kinds ...token.Kind) map[token.Kind]bool {
	m := make(map[token.Kind]bool, len(kinds))
	for _, kind := range kinds {
		m[kind] = true
	}
	return m
}
