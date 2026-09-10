package parser

import (
	"strconv"

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
	start := p.peek()
	prog := &ast.Program{At: start.Pos}
	if !p.match(token.ALGORITMO) {
		p.error(start, "expected algoritmo")
		return prog
	}
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Kind != token.STRING {
		p.error(start, "expected algorithm name string on the same line")
		return prog
	}
	prog.Name = p.advance().Text
	prog.Consts = p.parseConstBlock()
	prog.Types = p.parseTypeBlock()
	if len(p.diags) != 0 {
		return prog
	}
	if p.peek().Kind == token.VAR {
		prog.Globals = p.parseVarBlock()
	}
	for p.peek().Kind == token.PROCEDIMENTO || p.peek().Kind == token.FUNCAO {
		prog.Subs = append(prog.Subs, p.parseSubprogram())
		if len(p.diags) != 0 {
			return prog
		}
	}
	if !p.match(token.INICIO) {
		p.error(p.peek(), "expected inicio")
		return prog
	}
	prog.Body = p.parseStmtList(stopSet(token.FIMALGORITMO))
	p.expect(token.FIMALGORITMO, "expected fimalgoritmo")
	return prog
}

func (p *parser) parseConstBlock() []ast.ConstDecl {
	if !p.match(token.CONST) {
		return nil
	}
	var decls []ast.ConstDecl
	for p.peek().Kind == token.IDENT {
		name := p.advance()
		p.expect(token.EQL, "expected '=' after constant name")
		value := p.parseExpr(0)
		p.parseDeclarationSemicolon()
		decls = append(decls, ast.ConstDecl{Name: name, Value: value})
	}
	if p.peek().Kind != token.VAR && p.peek().Kind != token.TIPO {
		p.error(p.peek(), "expected var after constants")
	}
	return decls
}

func (p *parser) parseVarBlock() []ast.VarDecl {
	p.expect(token.VAR, "expected var")
	p.parseDeclarationSemicolon()
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
	p.parseDeclarationSemicolon()
	return ast.VarDecl{At: first.Pos, Names: names, Type: typ}
}

func (p *parser) parseDeclarationSemicolon() {
	if p.atLineEnd() || !p.match(token.SEMI) {
		return
	}
	if !p.atLineEnd() {
		p.error(p.peek(), "expected end of line after declaration semicolon")
		p.skipLine()
	}
}

func (p *parser) parseType() ast.TypeSpec {
	if !p.enter() {
		return ast.TypeSpec{Name: "inteiro"}
	}
	defer func() { p.depth-- }()
	tok := p.peek()
	switch tok.Kind {
	case token.IDENT:
		p.advance()
		return ast.TypeSpec{At: tok.Pos, Name: tok.Text}
	case token.INTEIRO, token.REAL, token.CARACTERE, token.LOGICO:
		p.advance()
		return ast.TypeSpec{At: tok.Pos, Name: tok.Kind.String()}
	case token.VETOR:
		p.advance()
		p.expect(token.LBRACK, "expected '[' after vetor")
		var ranges []ast.Range
		for {
			at := p.peek()
			errors := len(p.diags)
			low := p.parseBound()
			p.expect(token.DOTDOT, "expected '..' in vector bound")
			high := p.parseBound()
			lo, loLiteral := low.(*ast.LiteralExpr)
			hi, hiLiteral := high.(*ast.LiteralExpr)
			if len(p.diags) == errors && loLiteral && hiLiteral && hi.Int < lo.Int {
				p.error(at, "vector upper bound is smaller than lower bound")
			}
			ranges = append(ranges, ast.Range{At: at.Pos, Low: low, High: high})
			if !p.match(token.COMMA) {
				break
			}
		}
		if len(ranges) > 2 {
			p.error(tok, "vector declarations support at most two dimensions")
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

func (p *parser) parseBound() ast.Expr {
	tok := p.peek()
	if tok.Kind == token.NUMBER || tok.Kind == token.IDENT {
		p.advance()
	}
	v, err := strconv.ParseUint(tok.Text, 10, 63)
	if tok.Kind == token.IDENT || tok.Kind == token.NUMBER && err == nil {
		switch p.peek().Kind {
		case token.DOTDOT, token.COMMA, token.RBRACK:
			if tok.Kind == token.IDENT {
				return &ast.IdentExpr{Name: tok}
			}
			return &ast.LiteralExpr{At: tok.Pos, Kind: ast.IntLiteral, Int: int64(v)}
		}
	}
	p.error(tok, "expected unsigned integer literal or constant bound")
	for {
		switch p.peek().Kind {
		case token.DOTDOT, token.COMMA, token.RBRACK, token.DE, token.INICIO, token.EOF:
			return &ast.LiteralExpr{At: tok.Pos, Kind: ast.IntLiteral}
		}
		p.advance()
	}
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
	decl.Consts = p.parseConstBlock()
	decl.Types = p.parseTypeBlock()
	if len(p.diags) != 0 {
		return decl
	}
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
	ret := p.parseCallableType()
	decl := &ast.FunctionDecl{At: start.Pos, Name: name, Params: params, Return: ret}
	decl.Consts = p.parseConstBlock()
	decl.Types = p.parseTypeBlock()
	if len(p.diags) != 0 {
		return decl
	}
	if p.peek().Kind == token.VAR {
		decl.Locals = p.parseVarBlock()
	}
	p.expect(token.INICIO, "expected inicio in function")
	decl.Body = p.parseStmtList(stopSet(token.FIMFUNCAO))
	p.expect(token.FIMFUNCAO, "expected fimfuncao")
	return decl
}

func (p *parser) parseParamList() []ast.Param {
	if !p.match(token.LPAREN) {
		return nil
	}
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
		typ := p.parseCallableType()
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

func (p *parser) parseCallableType() ast.TypeSpec {
	typ := p.parseType()
	if len(p.diags) != 0 {
		return typ
	}
	if typ.Name == "vetor" {
		p.error(token.Token{Pos: typ.At}, "inline vector parameter and result types are unsupported")
	} else if token.Lookup(typ.Name) == token.IDENT {
		p.error(token.Token{Pos: typ.At}, "named parameter and result types are unsupported")
	}
	return typ
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
		stmt := &ast.ReturnStmt{At: tok.Pos}
		if !p.atLineEnd() {
			stmt.Value = p.parseExpr(0)
		}
		return stmt
	case token.LEIA:
		return p.parseRead()
	case token.ESCREVA, token.ESCREVAL:
		return p.parseWrite()
	case token.LIMPATELA:
		s := &ast.ClearStmt{At: p.advance().Pos}
		p.skipLine()
		return s
	case token.MUDACOR:
		return p.parseColor()
	default:
		p.error(p.advance(), "expected statement")
		p.skipLine()
		return nil
	}
}

func (p *parser) parseIdentStmt() ast.Stmt {
	if p.peekN(1).Kind == token.LPAREN {
		call := p.parseCall()
		return &ast.CallStmt{Call: call}
	}
	if p.peekN(1).Kind != token.ASSIGN && p.peekN(1).Kind != token.LBRACK {
		return &ast.CallStmt{Call: &ast.CallExpr{Name: p.advance()}}
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
	p.match(token.FACA)
	sw := &ast.SwitchStmt{At: start.Pos, X: x}
	for p.peek().Kind != token.FIMESCOLHA && p.peek().Kind != token.EOF {
		switch p.peek().Kind {
		case token.CASO:
			cstart := p.advance()
			labels := []ast.CaseLabel{p.parseCaseLabel()}
			for p.match(token.COMMA) {
				labels = append(labels, p.parseCaseLabel())
			}
			p.match(token.COLON)
			body := p.parseStmtList(stopSet(token.CASO, token.OUTROCASO, token.FIMESCOLHA))
			sw.Cases = append(sw.Cases, ast.CaseClause{At: cstart.Pos, Labels: labels, Body: body})
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

func (p *parser) parseCaseLabel() ast.CaseLabel {
	label := ast.CaseLabel{Low: p.parseExpr(0)}
	if !p.atLineEnd() && p.peek().Kind == token.ATE {
		at := p.advance()
		if p.atLineEnd() {
			p.error(at, "expected range upper bound")
		} else {
			label.High = p.parseExpr(0)
		}
	}
	return label
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
	if p.atLineEnd() {
		return stmt
	}
	if !p.match(token.LPAREN) {
		p.error(start, "expected '(' or end of line after write")
		p.skipLine()
		return stmt
	}
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

func (p *parser) skipLine() {
	for !p.atLineEnd() {
		p.pos++
	}
}

func (p *parser) atLineEnd() bool {
	return p.pos >= len(p.tokens) || p.tokens[p.pos].Kind == token.NEWLINE || p.tokens[p.pos].Kind == token.EOF
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
	i := p.peekIndex(0)
	if i < len(p.tokens) {
		p.pos = i + 1
	}
	return tok
}

func (p *parser) peek() token.Token {
	return p.peekN(0)
}

func (p *parser) peekN(n int) token.Token {
	i := p.peekIndex(n)
	if i >= 0 && i < len(p.tokens) {
		return p.tokens[i]
	}
	return token.Token{Kind: token.EOF}
}

func (p *parser) peekIndex(n int) int {
	i := p.pos
	for i < len(p.tokens) {
		if p.tokens[i].Kind != token.NEWLINE {
			if n == 0 {
				return i
			}
			n--
		}
		i++
	}
	return i
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
