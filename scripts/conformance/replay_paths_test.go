package main

import (
	"os"
	"strings"
	"testing"
)

func TestReplayInputPaths(t *testing.T) {
	t.Parallel()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ name, path, observer, want string }{
		{"source", "source.alg", "", "input file conflicts"},
		{"source case", "SOURCE.ALG", "", "input file conflicts"},
		{"source descendant", "source.alg/input.dat", "", "input file conflicts"},
		{"source case descendant", "SOURCE.ALG/input.dat", "", "input file conflicts"},
		{"state", "state.json", "state", "input file conflicts"},
		{"state case", "STATE.JSON", "state", "input file conflicts"},
		{"state descendant", "state.json/input.dat", "state", "input file conflicts"},
		{"host", "host.json", "host", "input file conflicts"},
		{"host case", "HOST.JSON", "host", "input file conflicts"},
		{"host descendant", "host.json/input.dat", "host", "input file conflicts"},
		{"clock", "clock.json", "clock", "input file conflicts"},
		{"clock case", "CLOCK.JSON", "clock", "input file conflicts"},
		{"clock descendant", "clock.json/input.dat", "clock", "input file conflicts"},
		{"clock with state adapter", "clock.json", "state", "input file conflicts"},
		{"host with clock adapter", "host.json", "clock", "input file conflicts"},
		{"source period", "source.alg.", "", "unsafe path"},
		{"clock directory space", "clock.json /input.dat", "clock", "unsafe path"},
		{"plain state", "state.json", "", ""},
		{"plain host", "host.json", "", ""},
		{"plain clock", "clock.json", "", ""},
		{"plain state directory", "state.json/input.dat", "", ""},
		{"nested source", "data/source.alg", "state", ""},
		{"nested state", "data/state.json", "state", ""},
		{"source prefix", "source.alg.backup", "state", ""},
		{"state prefix", "state.json.backup", "state", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			p := &m.Probes[0]
			p.TimeoutMS = 5000
			p.Implementation.State = "verified"
			p.Implementation.Tests = []string{"output_test.go#TestOutput"}
			p.Files = []generatedFile{{Path: tt.path, Content: writeArtifact(t, root, "fixture.txt", "Input bytes.\n")}}
			p.Evidence.Generated = p.Files
			p.Implementation.Expected.Generated = p.Files
			prefix, observer := []string{"-test.run=^TestReplayChild$", "--"}, ""
			switch tt.observer {
			case "state":
				a := writeArtifact(t, root, "expected-state.txt", "{}\n")
				p.Implementation.Expected.State = &a
			case "host":
				a := writeArtifact(t, root, "expected-host.txt", "[]\n")
				p.Implementation.Expected.HostTrace = &a
			case "clock":
				a := writeArtifact(t, root, "expected-clock.txt", `{"nowMS":[0]}`)
				p.Implementation.Expected.Clock = &a
			}
			if tt.observer == "" {
				p.Source = writeArtifact(t, root, p.Source.Path, "pass")
			} else {
				prefix, observer = []string{"-test.run=^TestReplayObserverChild$", "--"}, executable
			}
			for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
				t.Run(mode, func(t *testing.T) {
					err := validate(root, m, mode, nil)
					if tt.want == "" {
						if err != nil {
							t.Fatal(err)
						}
					} else if err == nil || !strings.Contains(err.Error(), tt.want) {
						t.Fatalf("validation error = %v, want %q", err, tt.want)
					}
				})
			}
			err := replayProbe(root, *p, executable, prefix, observer)
			if tt.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("replay error = %v, want %q", err, tt.want)
			}
		})
	}
}
