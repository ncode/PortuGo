package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatterEdgeCases(t *testing.T) {
	forms, err := os.ReadFile("../../testdata/cli/format_forms.alg")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ name, source, want string }{
		{"bare forms", string(forms), ""},
		{
			"suffix line endings",
			"algoritmo \"suffix\"\ninicio\nfimalgoritmo\r\r\nnotes\rtext\r\r\r\nlast\r\r",
			"algoritmo \"suffix\"\ninicio\nfimalgoritmo\nnotes\rtext\nlast\r\r",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "source.alg")
			if err := os.WriteFile(path, []byte(tt.source), 0600); err != nil {
				t.Fatal(err)
			}
			out, stderr, exit := commandOutput(t, "fmt", path)
			if exit != 0 || stderr != "" {
				t.Fatalf("format: exit=%d stderr=%q", exit, stderr)
			}
			if tt.want != "" && out != tt.want {
				t.Fatalf("formatted source=%q, want %q", out, tt.want)
			}
			if tt.name == "bare forms" && (!strings.Contains(out, "\n  retorne\n") || !strings.Contains(out, "\n  pi\n")) {
				t.Fatalf("bare forms were lost: %q", out)
			}
			for range 2 {
				stdout, stderr, exit := commandOutput(t, "fmt", "-w", path)
				if exit != 0 || stdout != "" || stderr != "" {
					t.Fatalf("write: exit=%d stdout=%q stderr=%q", exit, stdout, stderr)
				}
				assertFileBytes(t, path, []byte(out))
				stdout, stderr, exit = commandOutput(t, "fmt", "--check", path)
				if exit != 0 || stdout != "" || stderr != "" {
					t.Fatalf("check: exit=%d stdout=%q stderr=%q", exit, stdout, stderr)
				}
			}
		})
	}
}

func TestFormatterOperatorSpellings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "operators.alg")
	source := "algoritmo \"operators\"\ninicio\nescreval(NÃO falso E verdadeiro OU falso XOU verdadeiro)\nescreval(7 DIV 2 MOD 2)\nfimalgoritmo\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	out, stderr, exit := commandOutput(t, "fmt", path)
	want := "algoritmo \"operators\"\ninicio\n  escreval(nao falso e verdadeiro ou falso xou verdadeiro)\n  escreval(7 \\ 2 mod 2)\nfimalgoritmo\n"
	if exit != 0 || stderr != "" || out != want {
		t.Fatalf("exit=%d stderr=%q source=%q, want %q", exit, stderr, out, want)
	}
}
