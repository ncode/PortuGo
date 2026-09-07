package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/token"
)

// Run starts a small complete-program REPL. Submit an empty line to run.
func Run(options interp.Options, errout io.Writer) (bool, error) {
	if options.Input == nil {
		options.Input = strings.NewReader("")
	}
	if options.Output == nil {
		options.Output = io.Discard
	}
	reader := bufio.NewReader(options.Input)
	options.Input = reader
	i := interp.New(options)
	out := options.Output
	ok := true
	var lines []string
	size := 0
	if _, err := fmt.Fprintln(out, "Portugol REPL. Enter a complete program, blank line runs it, :sair exits."); err != nil {
		return false, err
	}
	for {
		if len(lines) == 0 {
			if _, err := fmt.Fprint(out, "portugol> "); err != nil {
				return false, err
			}
		} else {
			if _, err := fmt.Fprint(out, "... "); err != nil {
				return false, err
			}
		}
		// A separator or exit command must still fit after an exact-boundary
		// submission. Other text is rejected before joining the source buffer.
		line, err := readLine(reader, max(source.MaxBytes-size, len(":sair\r\n")))
		if err == io.EOF {
			return ok, nil
		}
		if err != nil {
			var d diag.Diagnostic
			if errors.As(err, &d) && d.Code == diag.EResource {
				return false, submissionLimit(len(lines)+1, source.MaxBytes-size+1)
			}
			return false, err
		}
		if strings.TrimSpace(line) == ":sair" {
			return ok, nil
		}
		if strings.TrimSpace(line) != "" {
			if len(line) > source.MaxBytes-size {
				return false, submissionLimit(len(lines)+1, source.MaxBytes-size+1)
			}
			lines = append(lines, line)
			size += len(line)
			continue
		}
		if len(lines) == 0 {
			continue
		}
		if !runSource(strings.Join(lines, ""), i, errout) {
			ok = false
		}
		clear(lines)
		lines = lines[:0]
		size = 0
	}
}

func runSource(src string, i *interp.Interpreter, errout io.Writer) bool {
	file, toks, lexDiags := lexer.Scan("<repl>", src)
	if len(lexDiags) > 0 {
		diag.Render(errout, file, lexDiags)
		return false
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) > 0 {
		diag.Render(errout, file, parseDiags)
		return false
	}
	info, semaDiags := sema.Analyze(prog)
	if diag.HasErrors(semaDiags) {
		diag.Render(errout, file, semaDiags)
		return false
	}
	ds := i.Run(prog, info)
	diag.Render(errout, file, ds)
	return !diag.HasErrors(ds)
}

func readLine(reader *bufio.Reader, limit int) (string, error) {
	var line strings.Builder
	for {
		part, err := reader.ReadSlice('\n')
		if err != nil && err != bufio.ErrBufferFull && err != io.EOF {
			return "", err
		}
		if len(part) > limit-line.Len() {
			return "", diag.Diagnostic{Code: diag.EResource, Message: "source size limit exceeded"}
		}
		line.Write(part)
		if err == io.EOF && line.Len() == 0 {
			return "", io.EOF
		}
		if err != bufio.ErrBufferFull {
			return line.String(), nil
		}
	}
}

func submissionLimit(line, column int) error {
	d := diag.Diagnostic{Code: diag.EResource, Pos: token.Pos(source.MaxBytes), End: token.Pos(source.MaxBytes + 1), Message: "source size limit exceeded"}
	return fmt.Errorf("<repl>:%d:%d: %w", line, column, d)
}
