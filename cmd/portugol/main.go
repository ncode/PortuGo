package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/repl"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func main() { os.Exit(dispatch(os.Args[1:])) }

func dispatch(args []string) int {
	if len(args) == 0 {
		usage()
		return 2
	}
	command := args[0]
	switch command {
	case "run", "check", "fmt", "repl":
	default:
		usage()
		return 2
	}
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var steps uint64
	var formatCheck, formatWrite bool
	if command == "run" || command == "repl" {
		fs.Uint64Var(&steps, "max-steps", 0, "maximum execution steps (0 is unlimited)")
	}
	if command == "fmt" {
		fs.BoolVar(&formatCheck, "check", false, "report whether the file is already formatted")
		fs.BoolVar(&formatWrite, "w", false, "write formatted source back to the file")
	}
	if err := fs.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if formatCheck && formatWrite {
		fmt.Fprintln(os.Stderr, "fmt: --check and -w cannot be used together")
		return 2
	}
	wantArgs := 1
	if command == "repl" {
		wantArgs = 0
	}
	if fs.NArg() != wantArgs {
		usage()
		return 2
	}
	options := interp.Options{Input: os.Stdin, Output: os.Stdout, MaxSteps: steps}
	if command == "repl" {
		ok, err := repl.Run(options, os.Stderr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if !ok {
			return 1
		}
		return 0
	}
	file, prog, ok, err := parsedProgram(fs.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if !ok {
		return 1
	}
	if command == "fmt" {
		var output formatBuffer
		if err := ast.Fprint(&output, prog); err != nil {
			var d diag.Diagnostic
			if errors.As(err, &d) {
				d.Pos = prog.At
				diag.Render(os.Stderr, file.Positions, []diag.Diagnostic{d})
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			return 1
		}
		// Formatting can add syntax nesting; validate before emitting any bytes.
		formattedFile, toks, ds := lexer.Scan(file.Positions.Name, output.data.String())
		_, parseDiags := parser.Parse(toks)
		ds = append(ds, parseDiags...)
		if diag.HasErrors(ds) {
			diag.Render(os.Stderr, formattedFile, ds)
			return 1
		}
		if formatCheck || formatWrite {
			if bytes.Equal(file.Original, output.data.Bytes()) {
				return 0
			}
			if formatCheck {
				fmt.Fprintf(os.Stderr, "%s: not formatted\n", fs.Arg(0))
				return 1
			}
			if err := writeFormatted(fs.Arg(0), output.data.Bytes()); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
			return 0
		}
		if _, err := output.data.WriteTo(os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}
	info, ds := sema.Analyze(prog)
	diag.Render(os.Stderr, file.Positions, ds)
	if diag.HasErrors(ds) {
		return 1
	}
	if command == "check" {
		return 0
	}
	ds = interp.New(options).Run(prog, info)
	diag.Render(os.Stderr, file.Positions, ds)
	if diag.HasErrors(ds) {
		return 1
	}
	return 0
}

type formatBuffer struct{ data bytes.Buffer }

func (b *formatBuffer) Write(p []byte) (int, error) {
	if len(p) > source.MaxBytes-b.data.Len() {
		return 0, diag.Diagnostic{Code: diag.EResource, Message: "formatted source size limit exceeded"}
	}
	return b.data.Write(p)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: portugol <run|check|fmt> [options] file.alg | portugol repl [--max-steps N]")
}

func parsedProgram(path string) (*source.File, *ast.Program, bool, error) {
	src, err := source.LoadFile(path)
	if err != nil {
		return nil, nil, false, err
	}
	file, toks, lexDiags := lexer.ScanFile(src)
	if diag.HasErrors(lexDiags) {
		diag.Render(os.Stderr, file, lexDiags)
		return src, nil, false, nil
	}
	prog, parseDiags := parser.Parse(toks)
	if diag.HasErrors(parseDiags) {
		diag.Render(os.Stderr, file, parseDiags)
		return src, prog, false, nil
	}
	return src, prog, true, nil
}
