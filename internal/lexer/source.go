package lexer

import (
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/token"
)

// ScanFile tokenizes decoded source, preserving original-byte source positions.
func ScanFile(src *source.File) (*token.File, []token.Token, []diag.Diagnostic) {
	_, tokens, diags := Scan(src.Positions.Name, src.Text)
	for i := range tokens {
		tokens[i].Pos = src.OriginalPos(tokens[i].Pos)
	}
	for i := range diags {
		diags[i].Pos = src.OriginalPos(diags[i].Pos)
		if diags[i].End != token.NoPos {
			diags[i].End = src.OriginalPos(diags[i].End)
		}
	}
	return src.Positions, tokens, diags
}
