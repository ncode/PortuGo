package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatterTrailingVarComments(t *testing.T) {
	source, err := os.ReadFile("../../testdata/cli/trailing_var_comments.alg")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("../../testdata/cli/trailing_var_comments.formatted.alg")
	if err != nil {
		t.Fatal(err)
	}
	for _, empty := range []bool{false, true} {
		name := "populated"
		input, expected := string(source), string(want)
		output := "P\n 7\n"
		if empty {
			name = "empty"
			for _, text := range []string{"  x: inteiro // variable comment\n", "  x <- F()\n", "  escreval(x)\n"} {
				input = strings.ReplaceAll(input, text, "")
				expected = strings.ReplaceAll(expected, text, "")
			}
			output = "P\n"
		}
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "comments.alg")
			if err := os.WriteFile(path, []byte(input), 0600); err != nil {
				t.Fatal(err)
			}
			out, stderr, exit := commandOutput(t, "fmt", path)
			if exit != 0 || stderr != "" || out != expected {
				t.Fatalf("format: exit=%d stderr=%q\ngot:\n%s\nwant:\n%s", exit, stderr, out, expected)
			}
			out, stderr, exit = commandOutput(t, "run", path)
			if exit != 0 || stderr != "" || out != output {
				t.Fatalf("original execution: exit=%d stderr=%q output=%q", exit, stderr, out)
			}
			for range 2 {
				out, stderr, exit = commandOutput(t, "fmt", "-w", path)
				if exit != 0 || stderr != "" || out != "" {
					t.Fatalf("write: exit=%d stderr=%q output=%q", exit, stderr, out)
				}
				assertFileBytes(t, path, []byte(expected))
				out, stderr, exit = commandOutput(t, "fmt", "--check", path)
				if exit != 0 || stderr != "" || out != "" {
					t.Fatalf("check: exit=%d stderr=%q output=%q", exit, stderr, out)
				}
			}
			out, stderr, exit = commandOutput(t, "run", path)
			if exit != 0 || stderr != "" || out != output {
				t.Fatalf("formatted execution: exit=%d stderr=%q output=%q", exit, stderr, out)
			}
		})
	}
}
