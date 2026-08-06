package main

import (
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
	"github.com/ncode/portugol-go/internal/token"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "run":
		err = runCmd(os.Args[2:])
	case "check":
		err = checkCmd(os.Args[2:])
	case "fmt":
		err = fmtCmd(os.Args[2:])
	case "repl":
		err = repl.Run(os.Stdin, os.Stdout, os.Stderr)
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: portugol <run|check|fmt|repl> [file.alg]")
}

func runCmd(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: portugol run file.alg")
	}
	file, prog, ok, err := checkedProgram(fs.Arg(0))
	if err != nil || !ok {
		return err
	}
	_ = file
	return interp.New(os.Stdin, os.Stdout).Run(prog)
}

func checkCmd(args []string) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: portugol check file.alg")
	}
	_, _, ok, err := checkedProgram(fs.Arg(0))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("check failed")
	}
	return nil
}

func fmtCmd(args []string) error {
	fs := flag.NewFlagSet("fmt", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: portugol fmt file.alg")
	}
	file, prog, ok, err := parsedProgram(fs.Arg(0))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("parse failed for %s", file.Name)
	}
	return ast.Fprint(os.Stdout, prog)
}

func checkedProgram(path string) (*token.File, *ast.Program, bool, error) {
	file, prog, ok, err := parsedProgram(path)
	if err != nil || !ok {
		return file, prog, ok, err
	}
	diags := sema.Check(prog)
	if len(diags) > 0 {
		diag.Render(os.Stderr, file, diags)
		return file, prog, false, nil
	}
	return file, prog, true, nil
}

func parsedProgram(path string) (*token.File, *ast.Program, bool, error) {
	src, err := source.ReadFile(path)
	if err != nil {
		return nil, nil, false, err
	}
	file, toks, lexDiags := lexer.Scan(path, src)
	if len(lexDiags) > 0 {
		diag.Render(os.Stderr, file, lexDiags)
		return file, nil, false, nil
	}
	prog, parseDiags := parser.Parse(toks)
	if len(parseDiags) > 0 {
		diag.Render(os.Stderr, file, parseDiags)
		return file, prog, false, nil
	}
	return file, prog, true, nil
}
