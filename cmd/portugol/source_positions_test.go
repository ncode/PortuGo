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

func TestOriginalByteDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		name, source, command, position string
	}{
		{"BOM", "\ufeff@", "check", ":1:4: L001:"},
		{"CP1252 lexer", "algoritmo \"a\"\r\ninicio\r\nescreval(\"\xe9\x80\") @\r\nfimalgoritmo", "check", ":3:16: L001:"},
		{"CP1252 parser", "algoritmo \"a\"\ninicio\nescreval(\"\xe9\x80\", )\nfimalgoritmo", "check", ":3:16: P001:"},
		{"CP1252 sema", "algoritmo \"a\"\ninicio\nescreval(\"\xe9\x80\", missing)\nfimalgoritmo", "check", ":3:16: E002:"},
		{"CP1252 runtime", "algoritmo \"a\"\ninicio\nescreval(\"\xe9\x80\", 1 / 0)\nfimalgoritmo", "run", ":3:18: R002:"},
		{"CP1252 REPL", "algoritmo \"a\"\ninicio\nescreval(\"\xe9\x80\", missing)\nfimalgoritmo\n:sair\n", "repl", ":3:16: E002:"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "positions.alg")
			if err := os.WriteFile(path, []byte(tt.source), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{"portugol", tt.command, path}
			if tt.command == "repl" {
				args = args[:2]
			}
			encoded, err := json.Marshal(args)
			if err != nil {
				t.Fatal(err)
			}
			executable, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, executable, "-test.run=^TestCommandContracts$")
			cmd.Env = append(os.Environ(), "PORTUGOL_COMMAND_ARGS="+string(encoded))
			cmd.Stdin = strings.NewReader(tt.source)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err = cmd.Run()
			if ctx.Err() != nil || cmd.ProcessState == nil {
				t.Fatalf("CLI did not complete: %v", err)
			}
			if cmd.ProcessState.ExitCode() != 1 || !strings.Contains(stderr.String(), tt.position) {
				t.Fatalf("exit %d, diagnostics %q; want %q", cmd.ProcessState.ExitCode(), &stderr, tt.position)
			}
		})
	}
}
