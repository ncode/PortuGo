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
		prec := op.Kind.BinaryPrecedence()
		if prec < minPrec {
			return left
		}
		p.advance()
		right := p.parseExpr(prec + 1)
		binary := &ast.BinaryExpr{Op: op, Left: left, Right: right}
		left = binary
		if binary.IsComparison() {
			return left
		}
	}
}

func (p *parser) parseUnary() ast.Expr {
	if p.peek().Kind == token.ADD || p.peek().Kind == token.SUB || p.peek().Kind == token.NAO {
		op := p.advance()
		return &ast.UnaryExpr{Op: op, X: p.parseExpr(op.Kind.UnaryPrecedence())}
	}
	return p.parsePrimary()
}

func (p *parser) parsePrimary() ast.Expr {
	tok := p.peek()
	switch tok.Kind {
	case token.NUMBER:
		p.advance()
		if v, err := strconv.ParseInt(tok.Text, 10, 32); err == nil {
			return &ast.LiteralExpr{At: tok.Pos, Kind: ast.IntLiteral, Int: v}
		}
		// A bare exponent marker contributes zero and retains the real type.
		v, err := strconv.ParseFloat(strings.TrimRight(tok.Text, "eE"), 64)
		if err != nil {
			p.error(tok, "invalid real literal")
		}
		return &ast.LiteralExpr{At: tok.Pos, Kind: ast.RealLiteral, Real: v}
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
	case token.RAND:
		return &ast.IdentExpr{Name: p.advance()}
	case token.LIMPATELA, token.MUDACOR, token.DOS:
		if p.peekN(1).Kind == token.LPAREN {
			p.parseCall() // The reference consumes this syntax without evaluating it.
		} else {
			p.advance()
		}
		return &ast.NoValueExpr{Keyword: tok}
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
	name := p.advance()
	if strings.EqualFold(name.Text, "pi") {
		p.error(name, "pi does not accept parentheses")
	}
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
