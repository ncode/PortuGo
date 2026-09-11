package interp

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
)

func TestFileInputCleanupAndReset(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "failure"}[fail], func(t *testing.T) {
			dir := t.TempDir()
			body := "leia(x)"
			if fail {
				body += "\nescreval(1/0)"
			}
			p, info := analyzed(t, "algoritmo \"cleanup\"\narquivo \"data.txt\"\nvar\nx: inteiro\ninicio\n"+body+"\nfimalgoritmo")
			var out bytes.Buffer
			i := New(Options{WorkingDir: dir, Input: strings.NewReader("7\n8\n"), Output: &out})
			ds := i.Run(p, info)
			if fail && (len(ds) != 1 || ds[0].Code != diag.RArithmetic) || !fail && len(ds) != 0 {
				t.Fatal(ds)
			}
			if i.fileInput.file != nil {
				t.Fatal("file remains open after execution")
			}
			name := filepath.Join(dir, "data.txt")
			data, err := os.ReadFile(name)
			want := "7\r\n"
			if fail {
				want = ""
			}
			if err != nil || string(data) != want {
				t.Fatalf("recording=%q error=%v", data, err)
			}
			if err := os.Rename(name, name+".saved"); err != nil {
				t.Fatal(err) // Windows also verifies that the file handle was closed.
			}
			next, nextInfo := analyzed(t, "algoritmo \"console\"\nvar x: inteiro\ninicio\nleia(x)\nfimalgoritmo")
			if ds := i.Run(next, nextInfo); len(ds) != 0 {
				t.Fatal(ds)
			}
			if out.String() != "7\n8\n" {
				t.Fatalf("console buffer or file state leaked between runs: %q", &out)
			}
		})
	}
}

func TestFileInputFailures(t *testing.T) {
	for _, path := range []string{"", "bad\x00name", "missing/data.txt", "."} {
		t.Run(path, func(t *testing.T) {
			i := New(Options{WorkingDir: t.TempDir()})
			err := i.configureFile(&ast.FileInputStmt{At: 19, Path: path})
			var d diag.Diagnostic
			if !errors.As(err, &d) || d.Code != diag.RHost || d.Pos != 19 || i.fileInput.file != nil {
				t.Fatalf("file failure=%v, want positioned R008 without a handle", err)
			}
		})
	}
	for _, mode := range []string{"read", "write", "flush", "close", "encode", "decode"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			name := filepath.Join(dir, "data.txt")
			if mode == "read" || mode == "write" || mode == "decode" {
				data := []byte("7\r\n")
				if mode == "decode" {
					data = []byte{'A', 0x81, '\n'}
				}
				if err := os.WriteFile(name, data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			i := New(Options{WorkingDir: dir})
			if err := i.configureFile(&ast.FileInputStmt{At: 23, Path: "data.txt"}); err != nil {
				t.Fatal(err)
			}
			if mode == "write" {
				i.fileInput.reader = nil // The read-only handle rejects a recording write.
				i.fileInput.writer = bufio.NewWriterSize(i.fileInput.file, 128)
			} else if mode == "flush" {
				if err := i.recordInput("7"); err != nil {
					t.Fatal(err)
				}
				if err := i.fileInput.file.Close(); err != nil {
					t.Fatal(err)
				}
			} else if mode != "encode" && mode != "decode" {
				if err := i.fileInput.file.Close(); err != nil {
					t.Fatal(err)
				}
			}
			var err error
			switch mode {
			case "read", "decode":
				_, err = i.readFileLine()
			case "write":
				err = i.recordInput(strings.Repeat("7", 129))
			case "flush":
				err = i.closeInputFile(true)
			case "close":
				err = i.closeInputFile(false)
			case "encode":
				err = i.recordInput("A😀")
			}
			var d diag.Diagnostic
			if !errors.As(err, &d) || d.Code != diag.RHost || d.Pos != 23 || strings.Contains(d.Message, dir) {
				t.Fatalf("failure=%v, want positioned R008 without host details", err)
			}
			if (mode == "read" || mode == "flush" || mode == "close") && !errors.Is(err, os.ErrClosed) {
				t.Fatal("file diagnostic lost the underlying error")
			}
			_ = i.closeInputFile(false)
			if mode == "encode" {
				data, err := os.ReadFile(name)
				if err != nil || len(data) != 0 {
					t.Fatalf("failed encoding wrote partial bytes: %q, %v", data, err)
				}
			}
			if mode == "write" {
				data, err := os.ReadFile(name)
				if err != nil || string(data) != "7\r\n" {
					t.Fatalf("failed write changed existing bytes: %q, %v", data, err)
				}
			}
		})
	}
}

func TestFileInputLineLimit(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "data.txt")
	if err := os.WriteFile(name, bytes.Repeat([]byte{0x80}, maxTextBytes/3+1), 0o600); err != nil {
		t.Fatal(err)
	}
	i := New(Options{WorkingDir: dir})
	if err := i.configureFile(&ast.FileInputStmt{At: 29, Path: "data.txt"}); err != nil {
		t.Fatal(err)
	}
	_, err := i.readFileLine()
	var d diag.Diagnostic
	if !errors.As(err, &d) || d.Code != diag.RStorage || d.Pos != 29 {
		t.Fatalf("expanded file input error=%v, want positioned R003", err)
	}
	if err := i.closeInputFile(false); err != nil {
		t.Fatal(err)
	}
}
