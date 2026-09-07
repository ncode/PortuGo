package main

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ncode/portugol-go/internal/interp"
)

func TestObservationAdapter(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.alg")
	if err := os.WriteFile(src, []byte("algoritmo \"state\"\nvar x: inteiro\ninicio\nx <- 7\nescreval(x)\nfimalgoritmo"), 0600); err != nil {
		t.Fatal(err)
	}
	state, host := filepath.Join(dir, "state.json"), filepath.Join(dir, "host.json")
	var out, stderr bytes.Buffer
	status := executeProbe([]string{"--state", state, "--host-trace", host, src}, strings.NewReader(""), &out, &stderr)
	if status != 0 || out.String() != " 7\n" || stderr.Len() != 0 {
		t.Fatalf("status %d, streams %q %q", status, &out, &stderr)
	}
	for path, want := range map[string]string{state: "{\"x\":7}\n", host: "[]\n"} {
		data, err := os.ReadFile(path)
		if err != nil || string(data) != want {
			t.Fatalf("observation %s: %q (%v)", filepath.Base(path), data, err)
		}
	}
}

func TestRecordingHost(t *testing.T) {
	h := &recordingHost{}
	before := h.Now()
	for _, err := range []error{h.Delay(2 * time.Second), h.Breakpoint(interp.Breakpoint{Pos: 3}), h.ClearScreen(), h.SetDisplay(interp.DisplayState{})} {
		if err != nil {
			t.Fatal(err)
		}
	}
	if !h.Now().Equal(before.Add(2*time.Second)) || len(h.events) != 6 {
		t.Fatalf("non-deterministic host: %+v", h.events)
	}
}

func TestReplayObserverChild(t *testing.T) {
	marker := slices.Index(os.Args, "--")
	if marker < 0 {
		return
	}
	os.Exit(run(os.Args[marker+1:], os.Stdout, os.Stderr))
}

func TestReplayObservations(t *testing.T) {
	root, m := testManifest(t)
	p := m.Probes[0]
	p.Source = writeArtifact(t, root, p.Source.Path, "algoritmo \"observe\"\nvar x: inteiro\ninicio\nx <- 1\nescreval(x)\nfimalgoritmo")
	p.TimeoutMS = 5000
	state := writeArtifact(t, root, "state.json", "{\"x\":1}\n")
	host := writeArtifact(t, root, "host.json", "[]\n")
	p.Implementation.Expected.State, p.Implementation.Expected.HostTrace = &state, &host
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	prefix := []string{"-test.run=^TestReplayObserverChild$", "--"}
	if err := replayProbe(root, p, "unused", prefix, executable); err != nil {
		t.Fatal(err)
	}
	state = writeArtifact(t, root, "state.json", "{\"x\":2}\n")
	if err := replayProbe(root, p, "unused", prefix, executable); err == nil || !strings.Contains(err.Error(), "observation mismatch") {
		t.Fatalf("false state verification: %v", err)
	}
}
