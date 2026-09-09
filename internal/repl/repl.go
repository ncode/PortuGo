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

// Run starts a complete-program REPL. A program terminator or EOF submits input.
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
	if _, err := fmt.Fprintln(out, "Portugol REPL. Enter a complete program, fimalgoritmo runs it, :sair exits."); err != nil {
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
			if len(lines) != 0 {
				valid, err := runSource(strings.Join(lines, ""), i, errout)
				return ok && valid, err
			}
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
		if len(lines) == 0 && strings.TrimSpace(line) == "" {
			continue
		}
		if len(line) > source.MaxBytes-size {
			return false, submissionLimit(len(lines)+1, source.MaxBytes-size+1)
		}
		lines = append(lines, line)
		size += len(line)
		if !submissionComplete(line) {
			continue
		}
		valid, err := runSource(strings.Join(lines, ""), i, errout)
		if err != nil {
			return false, err
		}
		ok = ok && valid
		clear(lines)
		lines = lines[:0]
		size = 0
	}
}

// Strings and comments end on each physical line. Scan each new line once for
// framing; runSource still diagnoses the full decoded program, including errors.
func submissionComplete(line string) bool {
	decoded, err := source.Decode([]byte(line))
	if err != nil {
		return false
	}
	_, toks, _ := lexer.Scan("<repl>", decoded)
	for _, tok := range toks {
		if tok.Kind == token.FIMALGORITMO {
			return true
		}
	}
	return false
}

func runSource(src string, i *interp.Interpreter, errout io.Writer) (bool, error) {
	decoded, err := source.Decode([]byte(src))
	if err != nil {
		return false, err
	}
	file, toks, lexDiags := lexer.Scan("<repl>", decoded)
	if len(lexDiags) > 0 {
		diag.Render(errout, file, lexDiags)
		return false, nil
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) > 0 {
		diag.Render(errout, file, parseDiags)
		return false, nil
	}
	info, semaDiags := sema.Analyze(prog)
	if diag.HasErrors(semaDiags) {
		diag.Render(errout, file, semaDiags)
		return false, nil
	}
	ds := i.Run(prog, info)
	diag.Render(errout, file, ds)
	return !diag.HasErrors(ds), nil
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
