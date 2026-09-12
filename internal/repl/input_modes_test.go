package repl

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/interp"
)

func TestInputModesPreserveNextSubmission(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, header, body, input, output, diagnostic string
		file, recording                               string
	}{
		{
			name: "console", body: "leia(x)\nescreval(x)",
			input: "42\n", output: "42\n 42\n",
		},
		{
			name: "file with unread values", header: "arquivo \"data.txt\"\n",
			body: "leia(x)\nescreval(x)", file: "41\r\n42\r\n", output: "41\n 41\n",
		},
		{
			name: "file exhausted", header: "arquivo \"data.txt\"\n",
			body: "leia(x)\nescreval(x)\nleia(x)\nescreval(x)",
			file: "41\r\n", input: "42\n", output: "41\n 41\n42\n 42\n",
		},
		{
			name: "new file recording", header: "arquivo \"data.txt\"\n",
			body: "leia(x)\nescreval(x)", input: "42\n", output: "42\n 42\n", recording: "42\r\n",
		},
		{
			name: "random mode resets", body: "aleatorio 7, 7\nleia(x)\nescreval(x)",
			output: "7\n 7\n",
		},
		{
			name: "random mode off", body: "aleatorio 7, 7\nleia(x)\naleatorio off\nleia(x)\nescreval(x)",
			input: "42\n", output: "7\n42\n 42\n",
		},
		{
			name: "echo and environment", body: "eco off\npausa\ndebug verdadeiro\nlimpatela\nmudacor(\"amarelo\", \"frente\")\ntimer 0\ncronometro off\ncronometro on\neco on\nleia(x)\nescreval(x)",
			input: "42\n", output: "\nO cronômetro não foi iniciado.\n\nCronômetro iniciado.\n42\n 42\n",
		},
		{
			name: "file failure", header: "arquivo \"missing/data.txt\"\n",
			body: "leia(x)", diagnostic: "R008",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "data.txt")
			if tt.file != "" {
				if err := os.WriteFile(path, []byte(tt.file), 0600); err != nil {
					t.Fatal(err)
				}
			}
			first := "algoritmo \"mode\"\n" + tt.header + "var x: inteiro\ninicio\n" + tt.body + "\nfimalgoritmo\n"
			next := "algoritmo \"next\"\nvar x: inteiro\ninicio\nleia(x)\nescreval(\"next=\",x)\nfimalgoritmo\n"
			input := first + tt.input + next + "99\n:sair\n"
			var out, stderr bytes.Buffer
			ok, err := Run(interp.Options{WorkingDir: dir, Input: strings.NewReader(input), Output: &out, MaxSteps: 1000}, &stderr)
			if err != nil || ok != (tt.diagnostic == "") {
				t.Fatalf("ok=%t error=%v diagnostics=%q", ok, err, &stderr)
			}
			want := "Portugol REPL. Enter a complete program, fimalgoritmo runs it, :sair exits.\nportugol> " +
				strings.Repeat("... ", strings.Count(first, "\n")-1) + tt.output + "portugol> " +
				strings.Repeat("... ", strings.Count(next, "\n")-1) + "99\nnext= 99\nportugol> "
			if out.String() != want {
				t.Errorf("shared input transcript: got %q, want %q", &out, want)
			}
			if tt.diagnostic == "" && stderr.Len() != 0 || tt.diagnostic != "" && strings.Count(stderr.String(), ": "+tt.diagnostic+":") != 1 {
				t.Errorf("diagnostics=%q, want %q", &stderr, tt.diagnostic)
			}
			wantFile := tt.file + tt.recording
			if wantFile != "" {
				data, err := os.ReadFile(path)
				if err != nil || string(data) != wantFile {
					t.Errorf("file content=%q, want %q (error=%v)", data, wantFile, err)
				}
				if err := os.Rename(path, path+".saved"); err != nil {
					t.Fatal(err) // Also checks closed handles on Windows.
				}
			}
		})
	}
}
