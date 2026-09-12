package parser

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/token"
)

func (p *parser) parseTypeBlock() []ast.TypeDecl {
	if !p.match(token.TIPO) {
		return nil
	}
	p.parseDeclarationSemicolon()
	var decls []ast.TypeDecl
	for p.peek().Kind == token.IDENT {
		name := p.advance()
		p.expect(token.EQL, "expected '=' after type name")
		var typ ast.TypeSpec
		if p.peek().Kind == token.REGISTRO {
			typ = p.parseRecordType()
		} else {
			typ = p.parseType()
		}
		if typ.Name == "vetor" {
			p.error(token.Token{Pos: typ.At}, "vector type aliases are unsupported")
		}
		decls = append(decls, ast.TypeDecl{Name: name, Type: typ})
		if len(p.diags) != 0 {
			return decls
		}
	}
	if p.peek().Kind != token.VAR && len(p.diags) == 0 {
		p.error(p.peek(), "expected var after types")
	}
	return decls
}
