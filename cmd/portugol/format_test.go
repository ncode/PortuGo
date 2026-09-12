package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func commandOutput(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(append([]string{"portugol"}, args...))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, "-test.run=^TestCommandContracts$")
	cmd.Env = append(os.Environ(), "PORTUGOL_COMMAND_ARGS="+string(encoded))
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err = cmd.Run()
	if ctx.Err() != nil || cmd.ProcessState == nil {
		t.Fatalf("CLI did not complete: %v (%v)", err, ctx.Err())
	}
	return stdout.String(), stderr.String(), cmd.ProcessState.ExitCode()
}

func TestFormatterModes(t *testing.T) {
	src, err := os.ReadFile("../../testdata/cli/format_modes.alg")
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := os.ReadFile("../../testdata/cli/format_modes.formatted.alg")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name string
		data []byte
	}{
		{"spacing", src},
		{"CRLF", bytes.ReplaceAll(canonical, []byte("\n"), []byte("\r\n"))},
		{"BOM", append([]byte{0xef, 0xbb, 0xbf}, canonical...)},
		{"CP1252", bytes.ReplaceAll(canonical, []byte("á"), []byte{0xe1})},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "source file.alg")
			if err := os.WriteFile(path, tt.data, 0640); err != nil {
				t.Fatal(err)
			}
			before, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			out, stderr, exit := commandOutput(t, "fmt", path)
			if exit != 0 || stderr != "" || out != string(canonical) {
				t.Fatalf("stdout format: exit=%d stderr=%q output=%q", exit, stderr, out)
			}
			out, stderr, exit = commandOutput(t, "fmt", "--check", path)
			if exit != 1 || out != "" || !strings.Contains(stderr, "not formatted") {
				t.Fatalf("dirty check: exit=%d stdout=%q stderr=%q", exit, out, stderr)
			}
			assertFileBytes(t, path, tt.data)
			originalOutput, originalError, originalExit := commandOutput(t, "run", path)
			for pass := range 2 {
				out, stderr, exit = commandOutput(t, "fmt", "-w", path)
				if exit != 0 || out != "" || stderr != "" {
					t.Fatalf("write pass %d: exit=%d stdout=%q stderr=%q", pass, exit, out, stderr)
				}
				assertFileBytes(t, path, canonical)
				out, stderr, exit = commandOutput(t, "fmt", "--check", path)
				if exit != 0 || out != "" || stderr != "" {
					t.Fatalf("clean check: exit=%d stdout=%q stderr=%q", exit, out, stderr)
				}
			}
			after, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if before.Mode().Perm() != after.Mode().Perm() {
				t.Errorf("permissions changed: %v to %v", before.Mode(), after.Mode())
			}
			out, stderr, exit = commandOutput(t, "run", path)
			if originalExit != 0 || originalError != "" || exit != 0 || stderr != "" || out != originalOutput {
				t.Errorf("format changed execution: before (%d, %q, %q), after (%d, %q, %q)", originalExit, originalOutput, originalError, exit, out, stderr)
			}
		})
	}
}

func TestFormatterPreservesMalformedFiles(t *testing.T) {
	for _, name := range []string{"lexer", "parser", "depth_limit"} {
		t.Run(name, func(t *testing.T) {
			src, err := os.ReadFile("../../testdata/cli/" + name + ".alg")
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(t.TempDir(), "malformed.alg")
			if err := os.WriteFile(path, src, 0600); err != nil {
				t.Fatal(err)
			}
			for _, flags := range [][]string{{}, {"--check"}, {"-w"}} {
				args := append([]string{"fmt"}, flags...)
				out, stderr, exit := commandOutput(t, append(args, path)...)
				if exit != 1 || out != "" || stderr == "" {
					t.Errorf("format %v: exit=%d stdout=%q stderr=%q", flags, exit, out, stderr)
				}
				assertFileBytes(t, path, src)
			}
		})
	}
}

func TestFormatterPreservesResourceLimitedFiles(t *testing.T) {
	prefix, suffix := "algoritmo \"limit\"\ninicio\nescreva(\"", "\")\nfimalgoritmo\n"
	for _, tt := range []struct {
		name, source string
	}{
		{"source size", strings.Repeat(" ", (4<<20)+1)},
		{"formatted size", prefix + strings.Repeat("a", (4<<20)-len(prefix)-len(suffix)) + suffix},
		{"formatted depth", "algoritmo \"nesting\"\ninicio\nescreval(" + strings.Repeat("- nao ", 85) + "1)\nfimalgoritmo\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "limited.alg")
			if err := os.WriteFile(path, []byte(tt.source), 0600); err != nil {
				t.Fatal(err)
			}
			for _, flag := range []string{"--check", "-w"} {
				out, stderr, exit := commandOutput(t, "fmt", flag, path)
				if exit != 1 || out != "" || !strings.Contains(stderr, "E900:") {
					t.Errorf("%s: exit=%d stdout bytes=%d stderr=%q", flag, exit, len(out), stderr)
				}
				assertFileBytes(t, path, []byte(tt.source))
			}
		})
	}
}

func TestFormatterWriteSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires additional Windows privileges")
	}
	path := filepath.Join(t.TempDir(), "target.alg")
	src := []byte("algoritmo \"x\"\ninicio\nescreval(1)\nfimalgoritmo\n")
	if err := os.WriteFile(path, src, 0600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(filepath.Dir(path), "link.alg")
	if err := os.Symlink(path, link); err != nil {
		t.Fatal(err)
	}
	out, stderr, exit := commandOutput(t, "fmt", "-w", link)
	if exit != 0 || out != "" || stderr != "" {
		t.Fatalf("write symlink: exit=%d stdout=%q stderr=%q", exit, out, stderr)
	}
	if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("symlink replaced: %v (%v)", info, err)
	}
	assertFileBytes(t, path, bytes.ReplaceAll(src, []byte("\nescreval"), []byte("\n  escreval")))
}

func TestFormatterFlagUsage(t *testing.T) {
	for _, args := range [][]string{
		{"fmt", "--check", "-w", "missing.alg"},
		{"fmt", "--check"},
		{"fmt", "-w"},
		{"run", "-w", "missing.alg"},
		{"check", "--check", "missing.alg"},
	} {
		out, stderr, exit := commandOutput(t, args...)
		if exit != 2 || out != "" || stderr == "" {
			t.Errorf("%v: exit=%d stdout=%q stderr=%q", args, exit, out, stderr)
		}
	}
}

func assertFileBytes(t *testing.T, path string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("file bytes changed: got %d bytes, want %d", len(got), len(want))
	}
}
