package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
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
		line, err := readLine(reader)
		if err == io.EOF {
			return ok, nil
		}
		if err != nil {
			return false, err
		}
		if strings.TrimSpace(line) == ":sair" {
			return ok, nil
		}
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
			continue
		}
		if len(lines) == 0 {
			continue
		}
		if !runSource(strings.Join(lines, "\n"), i, errout) {
			ok = false
		}
		lines = lines[:0]
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

func readLine(reader *bufio.Reader) (string, error) {
	var line strings.Builder
	for {
		part, more, err := reader.ReadLine()
		if err != nil {
			return "", err
		}
		if len(part) > 16<<20-line.Len() {
			return "", diag.Diagnostic{Code: diag.RStorage, Message: "input line size limit exceeded"}
		}
		line.Write(part)
		if !more {
			return line.String(), nil
		}
	}
}
