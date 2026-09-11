package parser

import (
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/token"
)

func (p *parser) parseRandomInput() ast.Stmt {
	s := &ast.RandomInputStmt{At: p.advance().Pos}
	defer p.skipLine()
	if p.atLineEnd() {
		p.error(token.Token{Pos: s.At}, "expected random-input mode or bounds")
		return s
	}
	if p.peek().Kind == token.IDENT && (strings.EqualFold(p.peek().Text, "on") || strings.EqualFold(p.peek().Text, "off")) {
		s.Off = strings.EqualFold(p.advance().Text, "off")
		return s
	}
	for {
		if p.atLineEnd() {
			p.error(token.Token{Pos: s.At}, "expected random-input bound")
			return s
		}
		s.Args = append(s.Args, p.parseExpr(0))
		if p.atLineEnd() || len(s.Args) == 3 || !p.match(token.COMMA) {
			return s
		}
	}
}
