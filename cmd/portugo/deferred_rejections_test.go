package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestDeferredScalarAssignment(t *testing.T) {
	types := []struct{ name, typ, value string }{
		{"integer", "inteiro", "1"}, {"real", "real", "1.5"}, {"text", "caractere", `"x"`}, {"logical", "logico", "verdadeiro"},
	}
	for _, dst := range types {
		for _, src := range types {
			if dst.name == src.name || dst.name == "real" && src.name == "integer" {
				continue
			}
			t.Run(src.name+" to "+dst.name, func(t *testing.T) {
				source := fmt.Sprintf("algoritmo \"assignment control\"\nvar\nx:%s\ninicio\nescreval(\"BEFORE\")\nx<-%s\nescreval(\"AFTER\")\nfimalgoritmo\n", dst.typ, src.value)
				path := filepath.Join(t.TempDir(), "source.alg")
				if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
					t.Fatal(err)
				}
				code, checkCode := "E001", "E001"
				if dst.name == "integer" && src.name == "real" {
					code, checkCode = "R001", ""
				}
				exit, out, stderr := rejectionCommand(t, []string{"run", "--max-steps", "10000", path}, "")
				assertRejection(t, exit, out, stderr, "BEFORE\n", code, 6)
				exit, out, stderr = rejectionCommand(t, []string{"check", path}, "")
				if checkCode != "" {
					assertRejection(t, exit, out, stderr, "", checkCode, 6)
				} else if exit != 0 || len(out) != 0 || len(stderr) != 0 {
					t.Fatal("static checking executed or rejected the recorded narrowing form")
				}
			})
		}
	}
}

func TestDeferredRejectionPaths(t *testing.T) {
	function := `algoritmo "assignment order"
var
v:vetor[1..1] de caractere
funcao Valor():inteiro
inicio
escreval("RHS")
retorne 1
fimfuncao
inicio
escreval("BEFORE")
v[%d]<-Valor()
fimalgoritmo
`
	for _, tt := range []struct {
		name, source, input, out, code string
		line                           int
	}{
		{"RHS before assignment failure", fmt.Sprintf(function, 1), "", "BEFORE\nRHS\n", "E001", 11},
		{"target failure before RHS", fmt.Sprintf(function, 2), "", "BEFORE\n", "R003", 11},
		{"malformed condition", "algoritmo \"condition control\"\nvar\nx:inteiro\ninicio\nescreval(\"BEFORE\")\nse x>=1 e x<=3 entao\nescreval(\"BODY\")\nfimse\nfimalgoritmo\n", "", "BEFORE\n", "P001", 6},
		{"missing fimse at EOF", "algoritmo \"eof control\"\nvar\ns:caractere\ninicio\nescreva(\"PROMPT?\")\nleia(s)\nse falso entao\nescreval(\"NO\")\nfimalgoritmo\n", "SAIR\n", "PROMPT?SAIR\n", "P001", 10},
		{"called malformed procedure", "algoritmo \"body control\"\nprocedimento P()\ninicio\nse verdadeiro entao\nescreval(\"BODY\")\nfimprocedimento\ninicio\nP()\nescreval(\"MAIN\")\nfimalgoritmo\n", "", "BODY\nMAIN\n", "P001", 11},
		{"unused malformed procedure", "algoritmo \"body control\"\nprocedimento P()\ninicio\nse verdadeiro entao\nescreval(\"BODY\")\nfimprocedimento\ninicio\nescreval(\"MAIN\")\nfimalgoritmo\n", "", "MAIN\n", "", 0},
		{"unselected extra repeat marker", "algoritmo \"repeat boundary\"\ninicio\nescreval(\"BEFORE\")\nse falso entao\nrepita\nescreval(\"BODY\")\nate verdadeiro\nfimrepita\nfimse\nescreval(\"AFTER\")\nfimalgoritmo\n", "", "BEFORE\nAFTER\n", "", 0},
		{"reached extra repeat marker", "algoritmo \"repeat boundary\"\ninicio\nescreval(\"BEFORE\")\nse verdadeiro entao\nrepita\nescreval(\"BODY\")\nate verdadeiro\nfimrepita\nfimse\nescreval(\"AFTER\")\nfimalgoritmo\n", "", "BEFORE\nBODY\n", "P001", 8},
		{"reached incompatible label", "algoritmo \"case control\"\nvar\ns:caractere\ninicio\nescreval(\"BEFORE\")\nse verdadeiro entao\nescolha s\ncaso 5\nescreval(\"BODY\")\nfimescolha\nfimse\nescreval(\"AFTER\")\nfimalgoritmo\n", "", "BEFORE\n", "E001", 8},
		{"unreached incompatible label", "algoritmo \"case control\"\nvar\ns:caractere\ninicio\nescreval(\"BEFORE\")\nse falso entao\nescolha s\ncaso 5\nescreval(\"BODY\")\nfimescolha\nfimse\nescreval(\"AFTER\")\nfimalgoritmo\n", "", "BEFORE\nAFTER\n", "", 0},
		{"blocking lookup", "algoritmo \"lookup\"\ninicio\nescreval(\"BEFORE\")\nescreval(missing)\nfimalgoritmo\n", "", "", "E002", 4},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "source.alg")
			if err := os.WriteFile(path, []byte(tt.source), 0o600); err != nil {
				t.Fatal(err)
			}
			exit, out, stderr := rejectionCommand(t, []string{"run", "--max-steps", "10000", path}, tt.input)
			if tt.code == "" {
				if exit != 0 || string(out) != tt.out || len(stderr) != 0 {
					t.Errorf("successful control: exit %d, stdout %q, diagnostics %d bytes", exit, out, len(stderr))
				}
				return
			}
			assertRejection(t, exit, out, stderr, tt.out, tt.code, tt.line)
			if tt.code == "P001" {
				exit, out, _ = rejectionCommand(t, []string{"fmt", "-w", path}, tt.input)
				after, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if exit != 1 || len(out) != 0 || !bytes.Equal(after, []byte(tt.source)) {
					t.Fatal("formatter modified malformed input or emitted source")
				}
			}
		})
	}
}

func TestDeferredRejectionREPLRecovery(t *testing.T) {
	root, probes := recordedRejectionProbes(t)
	for _, name := range []string{"assignment-integer-to-text", "condition-later-error", "eof-control", "repeat-no-newline", "body-called", "case-type-reached", "call-modulo"} {
		t.Run(name, func(t *testing.T) {
			var probe rejectionProbe
			for _, p := range probes {
				if p.ID == "rejection-"+name {
					probe = p
				}
			}
			if probe.ID == "" {
				t.Fatal("missing control")
			}
			src, err := os.ReadFile(filepath.Join(root, probe.Source.Path))
			if err != nil {
				t.Fatal(err)
			}
			stdin, err := os.ReadFile(filepath.Join(root, probe.Input.Path))
			if err != nil {
				t.Fatal(err)
			}
			input := string(src) + string(stdin) + "algoritmo \"second\"\ninicio\nescreval(\"CLEAN\")\nfimalgoritmo\n:sair\n"
			exit, out, stderr := rejectionCommand(t, []string{"repl", "--max-steps", "10000"}, input)
			code := probe.Implementation.Expected.Diagnostics[0].Code
			if exit != 1 || strings.Count(string(out), "CLEAN\n") != 1 || !strings.HasSuffix(string(out), "CLEAN\nportugo> ") || strings.Count(string(stderr), code+":") != 1 {
				t.Fatalf("REPL lost rejection or next submission: exit %d, stdout %q, stderr %q", exit, out, stderr)
			}
		})
	}
}

func assertRejection(t *testing.T, exit int, out, stderr []byte, want, code string, line int) {
	t.Helper()
	if exit != 1 || string(out) != want {
		t.Errorf("exit %d, stdout %q; want exit 1 and %q", exit, out, want)
	}
	matches := regexp.MustCompile(`:(\d+):(\d+): ([A-Z]\d{3}):`).FindAllStringSubmatch(string(stderr), -1)
	if len(matches) != 1 || matches[0][1] != fmt.Sprint(line) || matches[0][3] != code {
		t.Errorf("expected one %s on line %d, got diagnostic fields %v", code, line, matches)
	}
}

func TestAnalysisOnlyDeferredRejections(t *testing.T) {
	root, probes := recordedRejectionProbes(t)
	for _, probe := range probes {
		if !strings.HasPrefix(probe.ID, "rejection-") || probe.Implementation.Expected.ExitCode != 1 {
			continue
		}
		t.Run(probe.ID, func(t *testing.T) {
			src, err := os.ReadFile(filepath.Join(root, probe.Source.Path))
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(t.TempDir(), "source.alg")
			if err := os.WriteFile(path, src, 0o600); err != nil {
				t.Fatal(err)
			}
			_, out, _ := rejectionCommand(t, []string{"check", path}, "")
			if len(out) != 0 {
				t.Fatal("check executed program writes")
			}
			exit, out, _ := rejectionCommand(t, []string{"fmt", "-w", path}, "")
			if len(out) != 0 {
				t.Fatal("formatter executed program writes")
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if exit == 1 {
				if !bytes.Equal(src, after) {
					t.Fatal("rejected formatting changed source bytes")
				}
				return
			}
			if exit != 0 {
				t.Fatalf("formatter exit = %d", exit)
			}
			input, err := os.ReadFile(filepath.Join(root, probe.Input.Path))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(root, probe.Implementation.Expected.Stdout.Path))
			if err != nil {
				t.Fatal(err)
			}
			exit, out, stderr := rejectionCommand(t, []string{"run", "--max-steps", "10000", path}, string(input))
			code := probe.Implementation.Expected.Diagnostics[0].Code
			if exit != 1 || !bytes.Equal(out, want) || strings.Count(string(stderr), code+":") != 1 {
				t.Fatalf("formatting changed rejection behavior: exit %d, output %q, diagnostics %q", exit, out, stderr)
			}
		})
	}
}
