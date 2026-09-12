package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestReplayChild(t *testing.T) {
	if len(os.Args) < 5 || os.Args[len(os.Args)-4] != "run" {
		return
	}
	src, err := os.ReadFile(os.Args[len(os.Args)-1])
	if err != nil {
		os.Exit(2)
	}
	switch string(src) {
	case "pass":
		fmt.Print(" 1\n")
	case "mismatch":
		fmt.Print("1\n")
	case "random-input":
		fmt.Print("2.8470000000\n 2.847\n")
	case "reject":
		fmt.Fprintln(os.Stderr, "source.alg:4:2: E004: invalid call")
		os.Exit(1)
	case "unpositioned":
		fmt.Fprintln(os.Stderr, "source.alg:4:0: E004: invalid call")
		os.Exit(1)
	case "hang":
		time.Sleep(time.Hour)
	case "flood":
		for {
			fmt.Print(strings.Repeat("x", 4096))
		}
	case "file":
		b, err := os.ReadFile("input.dat")
		if err != nil {
			os.Exit(2)
		}
		if err := os.WriteFile("result.dat", b, 0o600); err != nil {
			os.Exit(2)
		}
		fmt.Print(" 1\n")
	default:
		os.Exit(2)
	}
	os.Exit(0)
}

func TestReplay(t *testing.T) {
	t.Parallel()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ source, want string }{
		{"pass", ""}, {"mismatch", "stdout mismatch"}, {"reject", ""}, {"unpositioned", "unpositioned diagnostic"}, {"hang", "deadline"}, {"flood", "output limit"}, {"file", ""}, {"unexpected-file", "expected absent file"},
	} {
		t.Run(tt.source, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			p := m.Probes[0]
			p.Source = writeArtifact(t, root, p.Source.Path, tt.source)
			p.TimeoutMS = 5000
			if tt.source == "hang" {
				p.TimeoutMS = 100
			}
			if tt.source == "reject" || tt.source == "unpositioned" {
				p.Implementation.Expected.Stdout = writeArtifact(t, root, "empty.txt", "")
				p.Implementation.Expected.ExitCode = 1
				p.Implementation.Expected.Diagnostics = []diagnostic{{Code: "E004", Line: 4, Column: 2}}
				if tt.source == "unpositioned" {
					p.Implementation.Expected.Diagnostics[0].Column = 0
				}
			}
			if tt.source == "file" {
				a := writeArtifact(t, root, "fixture.dat", "\xe9\r\n ")
				p.Files = []generatedFile{{Path: "input.dat", Content: a}}
				p.Implementation.Expected.Generated = []generatedFile{{Path: "result.dat", Content: a}}
			}
			p.Implementation.Expected.Absent = []string{"absent.dat"}
			if tt.source == "unexpected-file" {
				p.Source = writeArtifact(t, root, p.Source.Path, "pass")
				a := writeArtifact(t, root, "present.dat", "unexpected")
				p.Files = []generatedFile{{Path: "absent.dat", Content: a}}
			}
			err := replayProbe(root, p, executable, []string{"-test.run=^TestReplayChild$", "--"}, "")
			if tt.want == "" && err != nil {
				t.Fatal(err)
			}
			if tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestReplayDoesNotHideUnsupportedObservation(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	p := m.Probes[0]
	p.Implementation.Expected.State = &p.Evidence.Normalized
	if err := replayProbe(root, p, "unused", nil, ""); err == nil || !strings.Contains(err.Error(), "observation adapter") {
		t.Fatalf("error = %v", err)
	}
}

func TestReplaySummary(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		state          string
		mismatch, fail bool
	}{
		{"pending", false, false}, {"pending", true, false}, {"verified", false, false}, {"verified", true, true},
	} {
		r := replayResult{ID: "sample", State: tt.state}
		if tt.mismatch {
			r.Error = "stdout mismatch"
		}
		if got := replayFailed([]replayResult{r}); got != tt.fail {
			t.Errorf("%+v: failed=%v, want %v", r, got, tt.fail)
		}
	}
}

func TestDiagnosticMapping(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, stderr string
		want         []diagnostic
		valid        bool
	}{
		{"lexer", "source.alg:2:1: L001: invalid character\ncheck failed\n", []diagnostic{{Code: "L001", Line: 2}}, true},
		{"parser", "source.alg:4:2: P001: expected declaration\ncheck failed\n", []diagnostic{{Code: "P001", Line: 4, Column: 2}}, true},
		{"encoding", "source.alg:1:1: S001: bad encoding\n", []diagnostic{{Code: "S001", Line: 1}}, true},
		{"bare summary", "check failed\n", nil, false},
		{"unknown trailer", "source.alg:2:1: E001: type mismatch\nother failure\n", []diagnostic{{Code: "E001", Line: 2}}, false},
		{"whitespace stderr", "\n", nil, false},
		{"wrong line", "source.alg:3:1: E001: type mismatch\n", []diagnostic{{Code: "E001", Line: 2}}, false},
		{"runtime", "source.alg:4:2: R001: invalid operation\n", []diagnostic{{Code: "R001", Line: 4}}, true},
		{"missing position", "R001: invalid operation\n", []diagnostic{{Code: "R001", Line: 4}}, false},
		{"zero line", "source.alg:0:2: R001: invalid operation\n", []diagnostic{{Code: "R001"}}, false},
		{"zero column", "source.alg:4:0: R001: invalid operation\n", []diagnostic{{Code: "R001", Line: 4}}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := compareDiagnostics(tt.stderr, tt.want); (err == nil) != tt.valid {
				t.Fatalf("error=%v, valid=%v", err, tt.valid)
			}
		})
	}
}
