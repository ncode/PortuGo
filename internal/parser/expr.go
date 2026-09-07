package parser

import (
	"strconv"
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/token"
)

func (p *parser) parseExpr(minPrec int) ast.Expr {
	if !p.enter() {
		return &ast.LiteralExpr{Kind: ast.IntLiteral}
	}
	defer func() { p.depth-- }()
	left := p.parseUnary()
	for {
		op := p.peek()
		prec := precedence(op.Kind)
		if prec < minPrec {
			return left
		}
		p.advance()
		right := p.parseExpr(prec + 1)
		left = &ast.BinaryExpr{Op: op, Left: left, Right: right}
	}
}

func (p *parser) parseUnary() ast.Expr {
	if p.peek().Kind == token.SUB || p.peek().Kind == token.NAO {
		op := p.advance()
		prec := 7
		if op.Kind == token.SUB {
			prec = 9 // VisuAlg binds unary minus more tightly than power.
		}
		return &ast.UnaryExpr{Op: op, X: p.parseExpr(prec)}
	}
	return p.parsePrimary()
}

func (p *parser) parsePrimary() ast.Expr {
	tok := p.peek()
	switch tok.Kind {
	case token.NUMBER:
		p.advance()
		if strings.ContainsAny(tok.Text, ".eE") {
			v, err := strconv.ParseFloat(tok.Text, 64)
			if err != nil {
				p.error(tok, "invalid real literal")
			}
			return &ast.LiteralExpr{At: tok.Pos, Kind: ast.RealLiteral, Real: v}
		}
		v, err := strconv.ParseInt(tok.Text, 10, 64)
		if err != nil {
			p.error(tok, "invalid integer literal")
		}
		return &ast.LiteralExpr{At: tok.Pos, Kind: ast.IntLiteral, Int: v}
	case token.STRING:
		p.advance()
		return &ast.LiteralExpr{At: tok.Pos, Kind: ast.StringLiteral, Str: tok.Text}
	case token.VERDADEIRO, token.FALSO:
		p.advance()
		return &ast.LiteralExpr{At: tok.Pos, Kind: ast.BoolLiteral, Bool: tok.Kind == token.VERDADEIRO}
	case token.IDENT:
		if p.peekN(1).Kind == token.LPAREN {
			return p.parseCall()
		}
		return p.parseDesignator()
	case token.LPAREN:
		p.advance()
		expr := p.parseExpr(0)
		p.expect(token.RPAREN, "expected ')'")
		return expr
	default:
		p.error(tok, "expected expression")
		p.advance()
		return &ast.LiteralExpr{At: tok.Pos, Kind: ast.IntLiteral}
	}
}

func (p *parser) parseCall() *ast.CallExpr {
	name := p.expect(token.IDENT, "expected call name")
	call := &ast.CallExpr{Name: name}
	p.expect(token.LPAREN, "expected '('")
	if !p.match(token.RPAREN) {
		for {
			call.Args = append(call.Args, p.parseExpr(0))
			if !p.match(token.COMMA) {
				break
			}
		}
		p.expect(token.RPAREN, "expected ')'")
	}
	return call
}

func precedence(kind token.Kind) int {
	switch kind {
	case token.OU:
		return 1
	case token.XOU:
		return 2
	case token.E:
		return 3
	case token.EQL, token.NEQ, token.LSS, token.GTR, token.LEQ, token.GEQ:
		return 4
	case token.ADD, token.SUB:
		return 5
	case token.MUL, token.QUO, token.IDIV, token.REM, token.MOD:
		return 6
	case token.POW:
		return 8
	default:
		return -1
	}
}
