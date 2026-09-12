package parser

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

func (p *parser) parseColor() ast.Stmt {
	start := p.advance()
	begin, before := p.pos, len(p.diags)
	stmt := &ast.ColorStmt{At: start.Pos}
	if p.atLineEnd() && p.peek().Kind != token.LPAREN {
		p.error(start, "expected two color arguments")
		return stmt
	}
	if !p.match(token.LPAREN) || p.peek().Kind == token.RPAREN {
		p.diags = append(p.diags, diag.Diagnostic{Code: diag.ETypeMismatch, Pos: start.Pos, Message: "expected caractere color"})
		p.skipLine()
		return stmt
	}
	stmt.Color = p.parseExpr(0)
	if !p.match(token.COMMA) {
		p.error(start, "expected display target")
		p.skipLine()
		return stmt
	}
	if p.peek().Kind == token.RPAREN {
		p.diags = append(p.diags, diag.Diagnostic{Code: diag.ETypeMismatch, Pos: start.Pos, Message: "expected caractere display target"})
		p.skipLine()
		return stmt
	}
	stmt.Target = p.parseExpr(0)
	// Consume a continued argument for recovery, but reject its physical newline.
	if len(p.diags) == before {
		for _, tok := range p.tokens[begin:p.pos] {
			if tok.Kind == token.NEWLINE {
				p.error(start, "color arguments must be on the command line")
				break
			}
		}
	}
	p.skipLine()
	return stmt
}
