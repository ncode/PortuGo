package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunExitStatus(t *testing.T) {
	if path := os.Getenv("PORTUGOL_RUN_FIXTURE"); path != "" {
		os.Args = []string{"portugol", "run", path}
		main()
		os.Exit(0)
	}
	for _, tt := range []struct {
		name, code, stdout string
		exit               int
	}{
		{"valid", "", "ok\n", 0},
		{"runtime_error", "R002", "before", 1},
		{"depth_limit", "E900", "", 1},
		{"lexer", "L001", "", 1},
		{"parser", "P001", "", 1},
		{"undeclared", "E002", "", 1},
		{"exp_one", "E004", "", 1},
		{"division_integer", "E001", "", 1},
		{"formatted_boolean", "E001", "", 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path, err := filepath.Abs(filepath.Join("..", "..", "testdata", "cli", tt.name+".alg"))
			if err != nil {
				t.Fatal(err)
			}
			executable, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, executable, "-test.run=^TestRunExitStatus$")
			cmd.Env = append(os.Environ(), "PORTUGOL_RUN_FIXTURE="+path)
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err = cmd.Run()
			if ctx.Err() != nil || cmd.ProcessState == nil {
				t.Fatalf("CLI did not complete: %v (%v)", err, ctx.Err())
			}
			if got := cmd.ProcessState.ExitCode(); got != tt.exit {
				t.Errorf("exit = %d, want %d; stderr: %s", got, tt.exit, &stderr)
			}
			if stdout.String() != tt.stdout {
				t.Errorf("stdout = %q, want %q", &stdout, tt.stdout)
			}
			if tt.code == "" {
				if stderr.Len() != 0 {
					t.Errorf("unexpected stderr: %s", &stderr)
				}
			} else if !strings.Contains(stderr.String(), tt.code+":") || !strings.Contains(stderr.String(), path+":") {
				t.Errorf("missing positioned %s diagnostic: %s", tt.code, &stderr)
			}
		})
	}
}
