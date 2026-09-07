package diag

import (
	"fmt"
	"io"
	"slices"

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
	EResource     Code = "E900"

	// Runtime diagnostic categories.
	RType       Code = "R001"
	RArithmetic Code = "R002"
	RStorage    Code = "R003"
	RInput      Code = "R004"
	RCall       Code = "R005"
	RLoop       Code = "R006"
	RBuiltin    Code = "R007"
	RHost       Code = "R008"
)

// Severity describes whether a diagnostic prevents execution.
type Severity uint8

const (
	Error Severity = iota
	Warning
	Note
)

// Diagnostic describes a user-facing problem.
type Diagnostic struct {
	Code     Code
	Pos      token.Pos
	Message  string
	Severity Severity
	End      token.Pos // Exclusive end; zero means no recorded span.
	Cause    error     // Retained for inspection, never appended to rendered output.
}

// Unwrap exposes the optional underlying cause.
func (d Diagnostic) Unwrap() error { return d.Cause }

// Ordered returns a stable source-ordered copy without modifying its input.
func Ordered(diags []Diagnostic) []Diagnostic {
	out := slices.Clone(diags)
	slices.SortStableFunc(out, func(a, b Diagnostic) int {
		if a.Pos < b.Pos {
			return -1
		}
		if a.Pos > b.Pos {
			return 1
		}
		return 0
	})
	return out
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
	for _, d := range Ordered(diags) {
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

// HasErrors reports whether diagnostics contains an error.
func HasErrors(diags []Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == Error {
			return true
		}
	}
	return false
}
