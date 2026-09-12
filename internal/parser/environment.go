package parser

import (
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/token"
)

func (p *parser) parseEcho() ast.Stmt {
	s := &ast.EchoStmt{At: p.advance().Pos}
	if !p.atLineEnd() {
		mode := p.peek()
		if mode.Kind == token.IDENT && (strings.EqualFold(mode.Text, "on") || strings.EqualFold(mode.Text, "off")) {
			s.Mode = p.advance()
		}
	}
	p.skipLine()
	return s
}

func (p *parser) parseChronometer() ast.Stmt {
	s := &ast.ChronometerStmt{At: p.advance().Pos}
	if !p.atLineEnd() {
		mode := p.peek()
		if mode.Kind == token.IDENT && (strings.EqualFold(mode.Text, "on") || strings.EqualFold(mode.Text, "off")) {
			s.Off = strings.EqualFold(p.advance().Text, "off")
		} else {
			p.error(mode, "expected chronometer on or off")
		}
	}
	p.skipLine()
	return s
}
