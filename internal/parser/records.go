package parser

import (
	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/token"
)

func (p *parser) parseRecordType(header token.Pos) ast.TypeSpec {
	start := p.expect(token.REGISTRO, "expected registro")
	p.rememberFragment(header)
	typ := ast.TypeSpec{At: start.Pos, Name: "registro"}
	p.parseDeclarationSemicolon()
	for p.peek().Kind == token.IDENT {
		first := p.advance()
		names := []token.Token{first}
		for p.match(token.COMMA) {
			names = append(names, p.expect(token.IDENT, "expected field name"))
		}
		p.expect(token.COLON, "expected ':' after field name")
		field := p.parseType()
		p.rememberFragment(first.Pos)
		if field.Name == "vetor" && len(p.diags) == 0 {
			p.error(token.Token{Pos: field.At}, "vector fields are unsupported")
		}
		typ.Fields = append(typ.Fields, ast.VarDecl{At: first.Pos, Names: names, Type: field})
		if len(p.diags) != 0 {
			return typ
		}
	}
	typ.End = p.expect(token.FIMREGISTRO, "expected fimregistro").Pos
	p.parseDeclarationSemicolon()
	return typ
}
