package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCommandContracts(t *testing.T) {
	if args := os.Getenv("PORTUGOL_COMMAND_ARGS"); args != "" {
		if err := json.Unmarshal([]byte(args), &os.Args); err != nil {
			panic(err)
		}
		main()
		os.Exit(0)
	}
	path := filepath.Join(t.TempDir(), "program.alg")
	src := "algoritmo \"partial\"\ninicio\nescreva(\"before\")\nescreval(1 / 0)\nfimalgoritmo\n"
	if err := os.WriteFile(path, []byte(src), 0600); err != nil {
		t.Fatal(err)
	}
	oversized := filepath.Join(t.TempDir(), "oversized.alg")
	if err := os.WriteFile(oversized, []byte(strings.Repeat(" ", (4<<20)+1)), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name                  string
		args                  []string
		input, stdout, stderr string
		exit                  int
	}{
		{"runtime", []string{"run", path}, "", "before", ":4:12: R002:", 1},
		{"budget", []string{"run", "--max-steps", "2", path}, "", "before", ":4:1: R006:", 1},
		{"check", []string{"check", path}, "", "", "", 0},
		{"source limit", []string{"check", oversized}, "", "", "oversized.alg:1:4194305: E900:", 1},
		{"fmt", []string{"fmt", path}, "", "algoritmo \"partial\"\ninicio\n  escreva(\"before\")\n  escreval((1 / 0))\nfimalgoritmo\n", "", 0},
		{"missing", []string{"run"}, "", "", "usage:", 2},
		{"unknown flag", []string{"check", "--bogus"}, "", "", "flag provided", 2},
		{"negative budget", []string{"run", "--max-steps", "-1", path}, "", "", "invalid value", 2},
		{"repl args", []string{"repl", path}, "", "", "usage:", 2},
		{"unknown command", []string{"unknown"}, "", "", "usage:", 2},
		{"repl read", []string{"repl", "--max-steps", "20"}, "algoritmo \"read\"\nvar x: inteiro\ninicio\nleia(x)\nescreval(x)\nfimalgoritmo\n\n42\n:sair\n", "", "", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			executable, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, executable, "-test.run=^TestCommandContracts$")
			args, err := json.Marshal(append([]string{"portugol"}, tt.args...))
			if err != nil {
				t.Fatal(err)
			}
			cmd.Env = append(os.Environ(), "PORTUGOL_COMMAND_ARGS="+string(args))
			cmd.Stdin = strings.NewReader(tt.input)
			var out, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &out, &stderr
			err = cmd.Run()
			if ctx.Err() != nil || cmd.ProcessState == nil {
				t.Fatalf("CLI did not complete: %v", err)
			}
			if got := cmd.ProcessState.ExitCode(); got != tt.exit {
				t.Errorf("exit %d want %d: %s", got, tt.exit, &stderr)
			}
			if tt.name == "repl read" {
				if !strings.Contains(out.String(), " 42\nportugol> ") {
					t.Errorf("lost shared input: %q", &out)
				}
			} else if out.String() != tt.stdout {
				t.Errorf("stdout %q want %q", &out, tt.stdout)
			}
			if (tt.stderr == "" && stderr.Len() != 0) || !strings.Contains(stderr.String(), tt.stderr) {
				t.Errorf("stderr %q want %q", &stderr, tt.stderr)
			}
		})
	}
}
