package lexer

import (
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

// Scan tokenizes decoded Portugol source text.
func Scan(filename, src string) (*token.File, []token.Token, []diag.Diagnostic) {
	l := &scanner{
		src:  src,
		file: token.NewFile(filename, len(src)),
	}
	l.scan()
	return l.file, l.tokens, l.diags
}

type scanner struct {
	src    string
	file   *token.File
	offset int
	tokens []token.Token
	diags  []diag.Diagnostic
}

func (s *scanner) scan() {
	for {
		s.skipSpaceAndComments()
		if s.offset >= len(s.src) {
			s.emit(token.EOF, "", token.Pos(s.offset))
			return
		}
		start := s.offset
		r := s.advance()
		switch {
		case isIdentStart(r):
			s.scanIdent(start)
		case unicode.IsDigit(r):
			s.scanNumber(start)
		case r == '"':
			s.scanString(start)
		default:
			s.scanSymbol(start, r)
		}
	}
}

func (s *scanner) skipSpaceAndComments() {
	for s.offset < len(s.src) {
		r := s.peek()
		switch r {
		case ' ', '\t', '\r':
			s.advance()
		case '\n':
			s.advance()
			s.file.AddLine(s.offset)
		case '/':
			if s.peekNext() != '/' {
				return
			}
			for s.offset < len(s.src) && s.peek() != '\n' {
				s.advance()
			}
		case '{':
			start := s.offset
			s.advance()
			for s.offset < len(s.src) && s.peek() != '}' {
				if s.peek() == '\n' {
					s.advance()
					s.file.AddLine(s.offset)
					continue
				}
				s.advance()
			}
			if s.offset >= len(s.src) {
				s.error(token.Pos(start), "unterminated block comment")
				return
			}
			s.advance()
		default:
			return
		}
	}
}

func (s *scanner) scanIdent(start int) {
	for s.offset < len(s.src) && isIdentPart(s.peek()) {
		s.advance()
	}
	text := s.src[start:s.offset]
	s.emit(token.Lookup(text), text, token.Pos(start))
}

func (s *scanner) scanNumber(start int) {
	for s.offset < len(s.src) && unicode.IsDigit(s.peek()) {
		s.advance()
	}
	if s.offset < len(s.src) && s.peek() == '.' && s.peekNext() != '.' {
		s.advance()
		for s.offset < len(s.src) && unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}
	if s.offset < len(s.src) && (s.peek() == 'e' || s.peek() == 'E') {
		save := s.offset
		s.advance()
		if s.peek() == '+' || s.peek() == '-' {
			s.advance()
		}
		if s.offset >= len(s.src) || !unicode.IsDigit(s.peek()) {
			s.offset = save
		} else {
			for s.offset < len(s.src) && unicode.IsDigit(s.peek()) {
				s.advance()
			}
		}
	}
	s.emit(token.NUMBER, s.src[start:s.offset], token.Pos(start))
}

func (s *scanner) scanString(start int) {
	var text []rune
	for s.offset < len(s.src) {
		r := s.advance()
		switch r {
		case '"':
			s.emit(token.STRING, string(text), token.Pos(start))
			return
		case '\n':
			s.file.AddLine(s.offset)
			s.error(token.Pos(start), "unterminated string literal")
			return
		case '\\':
			if s.offset >= len(s.src) {
				break
			}
			esc := s.advance()
			switch esc {
			case 'n':
				text = append(text, '\n')
			case 't':
				text = append(text, '\t')
			case '"', '\\':
				text = append(text, esc)
			default:
				text = append(text, esc)
			}
		default:
			text = append(text, r)
		}
	}
	s.error(token.Pos(start), "unterminated string literal")
}

func (s *scanner) scanSymbol(start int, r rune) {
	switch r {
	case '+':
		s.emit(token.ADD, "+", token.Pos(start))
	case '-':
		s.emit(token.SUB, "-", token.Pos(start))
	case '*':
		s.emit(token.MUL, "*", token.Pos(start))
	case '/':
		s.emit(token.QUO, "/", token.Pos(start))
	case '\\':
		s.emit(token.IDIV, "\\", token.Pos(start))
	case '%':
		s.emit(token.REM, "%", token.Pos(start))
	case '^':
		s.emit(token.POW, "^", token.Pos(start))
	case '(':
		s.emit(token.LPAREN, "(", token.Pos(start))
	case ')':
		s.emit(token.RPAREN, ")", token.Pos(start))
	case '[':
		s.emit(token.LBRACK, "[", token.Pos(start))
	case ']':
		s.emit(token.RBRACK, "]", token.Pos(start))
	case ',':
		s.emit(token.COMMA, ",", token.Pos(start))
	case ':':
		s.emit(token.COLON, ":", token.Pos(start))
	case ';':
		s.emit(token.SEMI, ";", token.Pos(start))
	case '.':
		if s.match('.') {
			s.emit(token.DOTDOT, "..", token.Pos(start))
			return
		}
		s.error(token.Pos(start), "unexpected '.'")
	case '<':
		switch {
		case s.match('-'):
			s.emit(token.ASSIGN, "<-", token.Pos(start))
		case s.match('>'):
			s.emit(token.NEQ, "<>", token.Pos(start))
		case s.match('='):
			s.emit(token.LEQ, "<=", token.Pos(start))
		default:
			s.emit(token.LSS, "<", token.Pos(start))
		}
	case '>':
		if s.match('=') {
			s.emit(token.GEQ, ">=", token.Pos(start))
			return
		}
		s.emit(token.GTR, ">", token.Pos(start))
	case '=':
		s.emit(token.EQL, "=", token.Pos(start))
	default:
		s.error(token.Pos(start), "unexpected character "+strconv.QuoteRune(r))
	}
}

func (s *scanner) emit(kind token.Kind, text string, pos token.Pos) {
	s.tokens = append(s.tokens, token.Token{Kind: kind, Text: text, Pos: pos})
}

func (s *scanner) error(pos token.Pos, msg string) {
	s.diags = append(s.diags, diag.Diagnostic{Code: diag.ELexer, Pos: pos, Message: msg})
}

func (s *scanner) advance() rune {
	r, size := utf8.DecodeRuneInString(s.src[s.offset:])
	if r == utf8.RuneError && size == 0 {
		return 0
	}
	s.offset += size
	return r
}

func (s *scanner) peek() rune {
	if s.offset >= len(s.src) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(s.src[s.offset:])
	return r
}

func (s *scanner) peekNext() rune {
	if s.offset >= len(s.src) {
		return 0
	}
	_, size := utf8.DecodeRuneInString(s.src[s.offset:])
	if s.offset+size >= len(s.src) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(s.src[s.offset+size:])
	return r
}

func (s *scanner) match(want rune) bool {
	if s.offset >= len(s.src) || s.peek() != want {
		return false
	}
	s.advance()
	return true
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentPart(r rune) bool {
	return isIdentStart(r) || unicode.IsDigit(r)
}
