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
		src:       src,
		file:      token.NewFile(filename, len(src)),
		lineStart: true,
	}
	l.scan()
	return l.file, l.tokens, l.diags
}

type scanner struct {
	src        string
	file       *token.File
	offset     int
	tokens     []token.Token
	diags      []diag.Diagnostic
	lineStart  bool
	endProgram int
	endTokens  int
	rawOffset  int
}

func (s *scanner) scan() {
	for {
		s.skipSpaceAndComments()
		if s.endProgram != 0 && (s.lineStart || s.offset == len(s.src)) {
			s.scanSuffix()
		}
		if s.offset >= len(s.src) {
			s.emit(token.EOF, "", token.Pos(s.offset))
			return
		}
		start := s.offset
		s.lineStart = false
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
		case ' ', '\t':
			s.advance()
		case '\r', '\n':
			start := s.offset
			s.advance()
			if r == '\r' && !s.match('\n') {
				continue
			}
			s.file.AddLine(s.offset)
			s.lineStart = true
			s.emit(token.NEWLINE, s.src[start:s.offset], token.Pos(start))
			if s.endProgram != 0 {
				return
			}
		case '/', '*':
			if !s.lineStart && (r != '/' || s.peekNext() != '/') {
				return
			}
			s.scanComment()
		case '{', '}':
			s.scanComment()
		default:
			return
		}
	}
}

func (s *scanner) scanComment() {
	start := s.offset
	s.skipLine()
	s.emit(token.COMMENT, s.src[start:s.offset], token.Pos(start))
}

func (s *scanner) scanIgnoredLine() {
	for s.offset < len(s.src) && (s.peek() == ' ' || s.peek() == '\t') {
		s.advance()
	}
	if s.offset < len(s.src) && s.peek() != '\n' && (s.peek() != '\r' || s.peekNext() != '\n') {
		s.scanComment()
	}
}

func (s *scanner) skipLine() {
	for s.offset < len(s.src) && s.peek() != '\n' {
		if s.peek() == '\r' && s.peekNext() == '\n' {
			return
		}
		s.advance()
	}
}

func (s *scanner) scanIdent(start int) {
	for s.offset < len(s.src) && isIdentPart(s.peek()) {
		s.advance()
	}
	text := s.src[start:s.offset]
	kind := token.Lookup(text)
	leading := len(s.tokens) == 0 || s.tokens[len(s.tokens)-1].Kind == token.NEWLINE
	s.emit(kind, text, token.Pos(start))
	if kind == token.FIMALGORITMO && s.endProgram == 0 {
		// The terminator's own line still receives lexical validation.
		s.endProgram, s.endTokens = s.offset, len(s.tokens)
	}
	if kind == token.DOS && leading {
		// Configuration directives ignore the rest of their physical line.
		s.scanIgnoredLine()
	}
}

func (s *scanner) scanSuffix() {
	// Retain every decoded byte after the terminator, but do not tokenize
	// later physical lines: even malformed literals there are ignored.
	s.tokens = s.tokens[:s.endTokens]
	s.rawOffset = s.endProgram
	if s.endProgram < len(s.src) {
		s.emit(token.SUFFIX, s.src[s.endProgram:], token.Pos(s.endProgram))
	}
	for i := s.offset; i < len(s.src); i++ {
		if s.src[i] == '\n' {
			s.file.AddLine(i + 1)
		}
	}
	s.offset, s.endProgram = len(s.src), 0
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
		s.advance()
		// Exponent digits are unsigned; a following sign starts an operator.
		for s.offset < len(s.src) && unicode.IsDigit(s.peek()) {
			s.advance()
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
			header := len(s.tokens) > 0 && s.tokens[len(s.tokens)-1].Kind == token.ALGORITMO
			s.emit(token.STRING, string(text), token.Pos(start))
			if header {
				s.scanIgnoredLine()
			}
			return
		case '\r', '\n':
			newline := s.offset - 1
			if r == '\r' {
				if !s.match('\n') {
					text = append(text, r)
					continue
				}
			}
			s.file.AddLine(s.offset)
			s.lineStart = true
			s.error(token.Pos(start), "unterminated string literal")
			s.emit(token.NEWLINE, s.src[newline:s.offset], token.Pos(newline))
			return
		case '/':
			if s.peek() == '/' {
				// The reference strips // even inside a quoted value.
				s.skipLine()
				s.error(token.Pos(start), "unterminated string literal")
				return
			}
			text = append(text, r)
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
		if s.match('=') {
			s.emit(token.ASSIGN, ":=", token.Pos(start))
			return
		}
		s.emit(token.COLON, ":", token.Pos(start))
	case ';':
		s.emit(token.SEMI, ";", token.Pos(start))
	case '.':
		if s.match('.') {
			s.emit(token.DOTDOT, "..", token.Pos(start))
			return
		}
		if next := s.peek(); next >= '0' && next <= '9' {
			s.error(token.Pos(start), "unexpected '.'")
		} else {
			s.emit(token.DOT, ".", token.Pos(start))
		}
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
	s.tokens = append(s.tokens, token.Token{Kind: kind, Text: text, Pos: pos, Raw: s.src[s.rawOffset:s.offset]})
	s.rawOffset = s.offset
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
