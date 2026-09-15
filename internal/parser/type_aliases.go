package parser

import (
	"github.com/ncode/PortuGo/internal/ast"
	"github.com/ncode/PortuGo/internal/token"
)

func (p *parser) parseTypeBlock() []ast.TypeDecl {
	if !p.match(token.TIPO) {
		return nil
	}
	p.parseDeclarationSemicolon()
	var decls []ast.TypeDecl
	for isTypeDeclarationName(p.peek(), p.peekN(1)) {
		name := p.advance()
		p.expect(token.EQL, "expected '=' after type name")
		var typ ast.TypeSpec
		if p.peek().Kind == token.REGISTRO {
			typ = p.parseRecordType(name.Pos)
		} else {
			typ = p.parseType()
			p.rememberFragment(name.Pos)
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

func isTypeDeclarationName(name, next token.Token) bool {
	if name.Kind == token.IDENT {
		return true
	}
	// The recorded alias spelling "E" lexes as the logical conjunction
	// keyword. In a type declaration, the following '=' disambiguates it
	// from an expression and preserves the reference's declaration boundary.
	return name.Kind == token.E && next.Kind == token.EQL
}
