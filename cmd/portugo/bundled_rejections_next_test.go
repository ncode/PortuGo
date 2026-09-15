package main

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestUnsupportedCallableDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		name   string
		code   string
		source string
	}{
		{
			name: "unsupported-callable-name",
			code: "L001",
			source: `algoritmo "x"
funcao é(): inteiro
inicio
  retorne 1
fimfuncao
inicio
fimalgoritmo
`,
		},
		{
			name: "headerless-unsupported-callable-name",
			code: "L001",
			source: `funcao ok(): inteiro
inicio
  retorne 1
fimfuncao
funcao é(): inteiro
inicio
  retorne 1
fimfuncao
`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			status, stdout, stderr := runSourceCommand(t, "run", tt.source)
			if status != 1 {
				t.Fatalf("status = %d, want 1; stdout=%q stderr=%q", status, stdout, stderr)
			}
			diagnosticPattern := regexp.MustCompile(`:[0-9]+:[0-9]+: ([A-Z][0-9]{3}):`)
			matches := diagnosticPattern.FindAllStringSubmatch(stderr, -1)
			if len(matches) != 1 {
				t.Fatalf("stderr = %q, want one diagnostic, got %d", stderr, len(matches))
			}
			if matches[0][1] != tt.code {
				t.Fatalf("stderr = %q, want %s diagnostic", stderr, tt.code)
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty output", stdout)
			}
		})
	}
}

func runSourceCommand(t *testing.T, command, contents string) (status int, stdout, stderr string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "source.alg")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	oldOut, oldErr := os.Stdout, os.Stderr
	outReader, outWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errReader, errWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout, os.Stderr = outWriter, errWriter
	status = dispatch([]string{command, path})
	_ = outWriter.Close()
	_ = errWriter.Close()
	os.Stdout, os.Stderr = oldOut, oldErr
	out, outErr := io.ReadAll(outReader)
	errout, errRead := io.ReadAll(errReader)
	_ = outReader.Close()
	_ = errReader.Close()
	if outErr != nil || errRead != nil {
		t.Fatalf("read command output: %v %v", outErr, errRead)
	}
	return status, string(out), string(errout)
}
