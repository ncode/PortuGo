package diag

import (
	"fmt"
	"io"

	"github.com/ncode/portugol-go/internal/token"
)

// Code is a stable diagnostic identifier.
type Code string

const (
	// Syntax and lexical diagnostics.
	ELexer  Code = "L001"
	EParse  Code = "P001"
	EEncode Code = "S001"

	// Semantic diagnostics.
	ETypeMismatch Code = "E001"
	EUndeclared   Code = "E002"
	ERedeclared   Code = "E003"
	ECall         Code = "E004"
	EReturn       Code = "E005"
	EBreak        Code = "E006"
)

// Diagnostic describes a user-facing problem.
type Diagnostic struct {
	Code    Code
	Pos     token.Pos
	Message string
}

// Error implements error for convenient test output.
func (d Diagnostic) Error() string {
	if d.Code == "" {
		return d.Message
	}
	return string(d.Code) + ": " + d.Message
}

// Render writes diagnostics with source positions.
func Render(w io.Writer, file *token.File, diags []Diagnostic) {
	for _, d := range diags {
		pos := token.Position{}
		if file != nil {
			pos = file.Position(d.Pos)
		}
		if pos.Filename != "" {
			_, _ = fmt.Fprintf(w, "%s:%d:%d: %s: %s\n", pos.Filename, pos.Line, pos.Column, d.Code, d.Message)
			continue
		}
		_, _ = fmt.Fprintf(w, "%s: %s\n", d.Code, d.Message)
	}
}

// HasErrors reports whether diagnostics contains any entries.
func HasErrors(diags []Diagnostic) bool {
	return len(diags) > 0
}
