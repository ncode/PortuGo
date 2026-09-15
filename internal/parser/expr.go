package parser

import (
	"strconv"
	"strings"

	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/stdlib"
	"github.com/ncode/PortuGo/internal/token"
)

func (p *parser) parseExpr(minPrec int) ast.Expr {
	if !p.enter() {
		return &ast.LiteralExpr{Kind: ast.IntLiteral}
	}
	defer func() { p.depth-- }()
	start := p.peek().Pos
	left := p.parseUnary()
	for {
		op := p.peek()
		if p.recoveryPair && p.adjacentPair(token.MUL, token.QUO) {
			return left
		}
		prec := op.Kind.BinaryPrecedence()
		if prec < minPrec {
			return left
		}
		if p.writeExpr && p.adjacentPair(token.QUO, token.MUL) {
			p.advance()
			p.advance()
			previous := p.recoveryPair
			p.recoveryPair = true
			inner := p.parseExpr(0)
			p.recoveryPair = previous
			if !p.adjacentPair(token.MUL, token.QUO) {
				p.error(p.peek(), "expected closing operator pair")
				return left
			}
			p.advance()
			p.advance()
			left = &ast.RecoveryExpr{At: start, Operands: []ast.Expr{left, inner}}
			p.writeRecovery = true
			continue
		}
		p.advance()
		continued := p.writeExpr && p.pos < len(p.tokens) && p.tokens[p.pos].Kind == token.NEWLINE
		before := len(p.diags)
		right := p.parseExpr(prec + 1)
		if continued && len(p.diags) == before {
			// Keep parsing the operand for recovery, but do not let formatting
			// turn an incomplete physical output line into accepted syntax.
			p.error(op, "expected output operand on the same line")
		}
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
	case token.NUMBER, token.ABSENT_NUMBER:
		p.advance()
		if tok.Kind == token.ABSENT_NUMBER && p.writeExpr {
			p.writeRecovery = true
			return &ast.RecoveryExpr{At: tok.Pos, Text: tok.Text + "{"}
		}
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
	case token.IDIV:
		if p.peekN(1).Kind == token.LPAREN {
			p.parseCall() // The reference accepts this candidate as a no-value expression.
			return &ast.NoValueExpr{Keyword: tok}
		}
		p.error(tok, "expected expression")
		p.advance()
		return &ast.LiteralExpr{At: tok.Pos, Kind: ast.IntLiteral}
	case token.RAND:
		return &ast.IdentExpr{Name: p.advance()}
	case token.ECO:
		p.error(tok, "eco is a statement")
		fallthrough
	case token.LIMPATELA, token.MUDACOR, token.DOS, token.ALEATORIO, token.CRONOMETRO, token.TIMER, token.PAUSA, token.DEBUG:
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
	if descriptor, ok := stdlib.Lookup(strings.ToLower(name.Text)); ok && descriptor.Signature().Form == stdlib.Bare {
		p.error(name, descriptor.Name()+" does not accept parentheses")
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

func (p *parser) adjacentPair(first, second token.Kind) bool {
	a, b := p.peek(), p.peekN(1)
	return a.Kind == first && b.Kind == second && b.Pos == a.Pos+1
}
