package ast

import (
	"fmt"
	"strings"

	"github.com/ncode/portugol-go/internal/token"
)

// before emits comments between the preceding printed line and this position.
// Keeping one pending line lets its trailing comment stay on that line even
// when the printer has already changed indentation for the following block.
func (p *printer) before(pos token.Pos, indent int) {
	for len(p.comments) > 0 && p.comments[0].Pos < pos {
		comment := p.comments[0]
		p.comments = p.comments[1:]
		if comment.Inline && !comment.Semicolon && p.hasPending {
			p.pending += " " + comment.Text
		} else {
			p.flush()
			prefix := ""
			if comment.Semicolon {
				prefix = "; "
			}
			p.pending = strings.Repeat("  ", indent) + prefix + comment.Text
			p.hasPending = true
			p.flush()
		}
	}
	p.flush()
}

func (p *printer) flush() {
	if p.hasPending && p.err == nil {
		_, p.err = fmt.Fprintln(p.w, p.pending)
	}
	p.pending, p.hasPending = "", false
}

func (p *printer) section(tok token.Token, name string) {
	if tok.Kind != token.ILLEGAL {
		p.before(tok.Pos, p.indent)
	}
	p.line("%s", name)
}

func (p *printer) begin(pos token.Pos, sections DeclSections) {
	indent := p.indent
	if sections.Var.Kind == token.VAR {
		indent++
	}
	p.before(pos, indent)
	p.line("inicio")
}

func (p *printer) end(pos token.Pos, format string, args ...any) {
	p.before(pos, p.indent+1)
	p.lineFrom(pos, format, args...)
}

func (p *printer) lineFrom(pos token.Pos, format string, args ...any) {
	if !p.fragment(pos) {
		p.line(format, args...)
	}
}

func (p *printer) fragment(pos token.Pos) bool {
	fragment, ok := p.fragments[pos]
	if !ok {
		return false
	}
	for len(p.comments) > 0 && p.comments[0].Pos < fragment.End {
		p.comments = p.comments[1:]
	}
	p.flush()
	lines := strings.Split(normalizeLineEndings(fragment.Text), "\n")
	for i, line := range lines {
		indent := p.indent
		if i > 0 {
			indent++
		}
		if strings.TrimLeft(line, " \t") != "" {
			lines[i] = strings.Repeat("  ", indent) + strings.TrimLeft(line, " \t")
		} else {
			lines[i] = ""
		}
	}
	p.pending, p.hasPending = strings.Join(lines, "\n"), true
	return true
}
