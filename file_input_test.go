package portugol_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedFileInput(t *testing.T) {
	type artifact struct{ Path string }
	type file struct {
		Path    string
		Content artifact
	}
	var manifest struct {
		Probes []struct {
			ID             string
			Source, Input  artifact
			Files          []file
			Evidence       struct{ Accepted bool }
			Implementation struct {
				Expected struct {
					Stdout    artifact
					Generated []file
					Absent    []string
				}
			}
		}
	}
	read := func(name string) []byte {
		t.Helper()
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	if err := json.Unmarshal(read("testdata/conformance/visualg-3.0.7/manifest.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, p := range manifest.Probes {
		if !strings.HasPrefix(p.ID, "file-") || !p.Evidence.Accepted {
			continue
		}
		if strings.Contains(p.ID, "missing-parent") {
			continue // The documented file-error guard differs from this reference case.
		}
		count++
		t.Run(p.ID, func(t *testing.T) {
			src, err := source.Decode(read(p.Source.Path))
			if err != nil {
				t.Fatal(err)
			}
			for _, formatted := range []bool{false, true} {
				_, tokens, ds := lexer.Scan("source.alg", src)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				prog, ds := parser.Parse(tokens)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				info, ds := sema.Analyze(prog)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				dir := t.TempDir()
				for _, f := range p.Files {
					name := filepath.Join(dir, f.Path)
					if err := os.MkdirAll(filepath.Dir(name), 0o700); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(name, read(f.Content.Path), 0o600); err != nil {
						t.Fatal(err)
					}
				}
				var out bytes.Buffer
				i := interp.New(interp.Options{WorkingDir: dir, Input: bytes.NewReader(read(p.Input.Path)), Output: &out})
				if ds := i.Run(prog, info); len(ds) != 0 {
					t.Fatal(ds)
				}
				want := p.Implementation.Expected
				if !bytes.Equal(out.Bytes(), read(want.Stdout.Path)) {
					t.Fatalf("formatted=%v: stdout=%q, want %q", formatted, &out, read(want.Stdout.Path))
				}
				for _, f := range want.Generated {
					if !bytes.Equal(read(filepath.Join(dir, f.Path)), read(f.Content.Path)) {
						t.Fatalf("formatted=%v: file %s differs from recorded bytes", formatted, f.Path)
					}
				}
				for _, name := range want.Absent {
					if _, err := os.Lstat(filepath.Join(dir, name)); !errors.Is(err, os.ErrNotExist) {
						t.Fatalf("formatted=%v: expected absent file %s", formatted, name)
					}
				}
				var printed bytes.Buffer
				if err := ast.Fprint(&printed, prog); err != nil {
					t.Fatal(err)
				}
				src = printed.String()
			}
		})
	}
	if count < 25 {
		t.Fatal("missing file-input recordings")
	}
}

func TestRecordedFileInputDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		name string
		line int
	}{
		{"file-detail-bare", 2}, {"file-detail-numeric", 2}, {"file-detail-unquoted", 2}, {"file-detail-parenthesized", 2},
		{"file-detail-after-var", 4}, {"file-detail-body", 5},
		{"file-atomic-positioned-failure", 8}, {"file-atomic-buffer-127", 7},
		{"file-atomic-buffer-128", 7}, {"file-atomic-buffer-129", 7}, {"file-atomic-buffer-256", 7},
		{"file-atomic-buffer-126", 7}, {"file-atomic-buffer-254", 7},
	} {
		t.Run(tt.name, func(t *testing.T) {
			probeDir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", tt.name)
			src, err := source.ReadFile(filepath.Join(probeDir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			input, err := os.ReadFile(filepath.Join(probeDir, "input.txt"))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(probeDir, "stdout.txt"))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
			dir := t.TempDir()
			var out bytes.Buffer
			file, tokens, ds := lexer.Scan("source.alg", src)
			if len(ds) == 0 {
				p, parseDiags := parser.Parse(tokens)
				ds = parseDiags
				if len(ds) == 0 {
					info, semaDiags := sema.Analyze(p)
					ds = semaDiags
					if len(ds) == 0 {
						ds = interp.New(interp.Options{WorkingDir: dir, Input: bytes.NewReader(input), Output: &out}).Run(p, info)
					}
				}
			}
			if len(ds) != 1 || ds[0].Code != diag.EParse || file.Position(ds[0].Pos).Line != tt.line {
				t.Fatalf("diagnostics=%v, want P001 on line %d", ds, tt.line)
			}
			if !bytes.Equal(out.Bytes(), want) {
				t.Fatalf("stdout=%q, want %q", &out, want)
			}
			if strings.HasPrefix(tt.name, "file-atomic-") {
				wantFile, err := os.ReadFile(filepath.Join(probeDir, "final-0.dat"))
				if err != nil {
					t.Fatal(err)
				}
				data, err := os.ReadFile(filepath.Join(dir, "record.txt"))
				if err != nil || !bytes.Equal(data, wantFile) {
					t.Fatalf("retained file bytes=%q error=%v, want %q", data, err, wantFile)
				}
			}
		})
	}
}
