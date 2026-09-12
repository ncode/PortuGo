package parser

import (
	"sort"
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/token"
)

// rememberFragment keeps only syntax with a comment before its last code token.
// Raw decoded token slices avoid confusing UTF-8 text with original byte offsets.
func (p *parser) rememberFragment(start token.Pos) {
	if p.pos == 0 || p.pos > len(p.tokens) {
		return
	}
	last := p.tokens[p.pos-1]
	begin := sort.Search(len(p.original), func(i int) bool { return p.original[i].Pos >= start })
	end := sort.Search(len(p.original), func(i int) bool { return p.original[i].Pos > last.Pos })
	if begin >= end {
		return
	}
	hasComment := false
	for _, tok := range p.original[begin:end] {
		if tok.Kind == token.COMMENT && tok.Pos < last.Pos {
			hasComment = true
		}
	}
	if !hasComment {
		return
	}
	var text strings.Builder
	for _, tok := range p.original[begin:end] {
		text.WriteString(tok.Raw)
	}
	if p.fragments == nil {
		p.fragments = make(map[token.Pos]ast.CommentFragment)
	}
	p.fragments[start] = ast.CommentFragment{End: last.Pos, Text: strings.TrimLeft(text.String(), " \t")}
}

// collectComments removes trivia without disturbing physical newline tokens.
func collectComments(tokens []token.Token) ([]token.Token, []ast.CommentGroup) {
	var comments []ast.CommentGroup
	var syntax []token.Token
	inline := false
	onlySemicolons, semicolon := true, false
	for i, tok := range tokens {
		if tok.Kind == token.COMMENT {
			end := tok.Pos
			if i+1 < len(tokens) {
				end = tokens[i+1].Pos
			}
			comments = append(comments, ast.CommentGroup{Pos: tok.Pos, End: end, Text: tok.Text, Inline: inline, Semicolon: onlySemicolons && semicolon})
			continue
		}
		syntax = append(syntax, tok)
		inline = tok.Kind != token.NEWLINE
		switch tok.Kind {
		case token.NEWLINE:
			onlySemicolons, semicolon = true, false
		case token.SEMI:
			semicolon = true
		default:
			onlySemicolons = false
		}
	}
	return syntax, comments
}

func (p *parser) section(kind token.Kind) token.Token {
	if p.peek().Kind == kind {
		return p.peek()
	}
	return token.Token{}
}
