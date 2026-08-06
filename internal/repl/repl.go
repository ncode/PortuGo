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
func Run(in io.Reader, out, errout io.Writer) error {
	scanner := bufio.NewScanner(in)
	var lines []string
	if _, err := fmt.Fprintln(out, "Portugol REPL. Enter a complete program, blank line runs it, :sair exits."); err != nil {
		return err
	}
	for {
		if len(lines) == 0 {
			if _, err := fmt.Fprint(out, "portugol> "); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprint(out, "... "); err != nil {
				return err
			}
		}
		if !scanner.Scan() {
			return scanner.Err()
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == ":sair" {
			return nil
		}
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
			continue
		}
		if len(lines) == 0 {
			continue
		}
		runSource(strings.Join(lines, "\n"), out, errout)
		lines = lines[:0]
	}
}

func runSource(src string, out, errout io.Writer) {
	file, toks, lexDiags := lexer.Scan("<repl>", src)
	if len(lexDiags) > 0 {
		diag.Render(errout, file, lexDiags)
		return
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) > 0 {
		diag.Render(errout, file, parseDiags)
		return
	}
	if semaDiags := sema.Check(prog); len(semaDiags) > 0 {
		diag.Render(errout, file, semaDiags)
		return
	}
	if err := interp.New(strings.NewReader(""), out).Run(prog); err != nil {
		_, _ = fmt.Fprintln(errout, err)
	}
}
