package main

import (
	"os"
	"strings"
	"testing"
)

func TestReplayObserverOwnedExpectations(t *testing.T) {
	t.Parallel()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name, path, observer, kind, want string
	}{
		{"generated source", "source.alg", "state", "generated", "generated expectation conflicts"},
		{"generated state case", "STATE.JSON", "state", "generated", "generated expectation conflicts"},
		{"generated host child", "host.json/output.dat", "host", "generated", "generated expectation conflicts"},
		{"absent clock", "clock.json", "clock", "absent", "absent expectation conflicts"},
		{"absent source case", "SOURCE.ALG", "state", "absent", "absent expectation conflicts"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			p := &m.Probes[0]
			p.Implementation.State = "verified"
			p.Implementation.Tests = []string{"output_test.go#TestOutput"}
			var observation artifact
			switch tt.observer {
			case "state":
				observation = writeArtifact(t, root, "expected-state.txt", "{\"x\":1}\n")
				p.Implementation.Expected.State = &observation
			case "host":
				observation = writeArtifact(t, root, "expected-host.txt", "[]\n")
				p.Implementation.Expected.HostTrace = &observation
			case "clock":
				observation = writeArtifact(t, root, "expected-clock.txt", "{\"nowMS\":[0]}\n")
				p.Implementation.Expected.Clock = &observation
			}
			if tt.kind == "generated" {
				generated := generatedFile{Path: tt.path, Content: observation}
				p.Evidence.Generated = []generatedFile{generated}
				p.Implementation.Expected.Generated = []generatedFile{generated}
			} else {
				p.Evidence.Absent = []string{tt.path}
				p.Implementation.Expected.Absent = []string{tt.path}
			}
			for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
				t.Run(mode, func(t *testing.T) {
					err := validate(root, m, mode, nil)
					if err == nil || !strings.Contains(err.Error(), tt.want) {
						t.Fatalf("validation error = %v, want %q", err, tt.want)
					}
				})
			}
			prefix := []string{"-test.run=^TestReplayObserverChild$", "--"}
			err := replayProbe(root, *p, executable, prefix, executable)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("replay error = %v, want %q", err, tt.want)
			}
		})
	}
}
