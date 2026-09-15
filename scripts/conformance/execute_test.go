package main

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ncode/PortuGo/internal/interp"
	"github.com/ncode/PortuGo/internal/testprocess"
)

func TestObservationAdapter(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.alg")
	if err := os.WriteFile(src, []byte("algoritmo \"state\"\ndos\nvar x: inteiro\ninicio\neco off\neco on\nx <- 7\nescreval(x)\nfimalgoritmo"), 0600); err != nil {
		t.Fatal(err)
	}
	state, host := filepath.Join(dir, "state.json"), filepath.Join(dir, "host.json")
	var out, stderr bytes.Buffer
	status := executeProbe([]string{"--state", state, "--host-trace", host, src}, strings.NewReader(""), &out, &stderr)
	if status != 0 || out.String() != " 7\n" || stderr.Len() != 0 {
		t.Fatalf("status %d, streams %q %q", status, &out, &stderr)
	}
	for path, want := range map[string]string{state: "{\"x\":7}\n", host: "[{\"operation\":\"console\"},{\"operation\":\"echo\",\"echo\":false},{\"operation\":\"echo\",\"echo\":true}]\n"} {
		data, err := os.ReadFile(path)
		if err != nil || string(data) != want {
			t.Fatalf("observation %s: %q (%v)", filepath.Base(path), data, err)
		}
	}
}

func TestReplayBudgetExhaustion(t *testing.T) {
	testprocess.Run(t, func() {
		source := filepath.Join(t.TempDir(), "source.alg")
		program := "algoritmo \"budget\"\ninicio\nenquanto verdadeiro faca\nfimenquanto\nfimalgoritmo\n"
		if err := os.WriteFile(source, []byte(program), 0600); err != nil {
			t.Fatal(err)
		}
		var out, stderr bytes.Buffer
		status := executeProbe([]string{"--max-steps", "8", source}, strings.NewReader(""), &out, &stderr)
		diagnostics := strings.ReplaceAll(stderr.String(), source, filepath.Base(source))
		if status != 1 || out.Len() != 0 || strings.Count(diagnostics, ": R006:") != 1 {
			t.Fatalf("status=%d output=%q diagnostics=%q; want one bounded R006", status, &out, diagnostics)
		}
	})
}

func TestObservationAdapterUsesClockFixture(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.alg")
	if err := os.WriteFile(src, []byte("algoritmo \"clock\"\ninicio\ncronometro on\ncronometro off\nfimalgoritmo"), 0600); err != nil {
		t.Fatal(err)
	}
	clock := filepath.Join(dir, "clock.json")
	if err := os.WriteFile(clock, []byte("{\"nowMS\":[0,16]}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	host := filepath.Join(dir, "host.json")
	var out, stderr bytes.Buffer
	status := executeProbe([]string{"--clock", clock, "--host-trace", host, src}, strings.NewReader(""), &out, &stderr)
	want := "\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 16 ms.\n"
	if status != 0 || out.String() != want || stderr.Len() != 0 {
		t.Fatalf("status %d, streams %q %q", status, &out, &stderr)
	}
	data, err := os.ReadFile(host)
	if err != nil || string(data) != "[{\"operation\":\"now\"},{\"operation\":\"now\"}]\n" {
		t.Fatalf("host trace: %q (%v)", data, err)
	}
}

func TestObservationAdapterRejectsExhaustedClockFixture(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.alg")
	if err := os.WriteFile(src, []byte("algoritmo \"clock\"\ninicio\ncronometro on\ncronometro off\nfimalgoritmo"), 0600); err != nil {
		t.Fatal(err)
	}
	clock := filepath.Join(dir, "clock.json")
	if err := os.WriteFile(clock, []byte("{\"nowMS\":[0]}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	status := executeProbe([]string{"--clock", clock, src}, strings.NewReader(""), &out, &stderr)
	if status != 1 || out.Len() != 0 || !strings.Contains(stderr.String(), "clock fixture exhausted") {
		t.Fatalf("status %d, streams %q %q", status, &out, &stderr)
	}
}

func TestObservationAdapterRejectsOversizedClockFixture(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.alg")
	if err := os.WriteFile(src, []byte("algoritmo \"clock\"\ninicio\nfimalgoritmo"), 0600); err != nil {
		t.Fatal(err)
	}
	clock := filepath.Join(dir, "clock.json")
	data := append([]byte(`{"nowMS":[0]}`), []byte(strings.Repeat(" ", maxArtifactBytes))...)
	if err := os.WriteFile(clock, data, 0600); err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	status := executeProbe([]string{"--clock", clock, src}, strings.NewReader(""), &out, &stderr)
	if status != 1 || !strings.Contains(stderr.String(), "clock fixture exceeds artifact size limit") {
		t.Fatalf("status %d, stderr %q", status, stderr.String())
	}
}

func TestObservationAdapterRejectsNonRegularSource(t *testing.T) {
	source := filepath.Join(t.TempDir(), "source.alg")
	if err := os.Mkdir(source, 0o700); err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	status := executeProbe([]string{source}, strings.NewReader(""), &out, &stderr)
	if status != 1 || out.Len() != 0 || !strings.Contains(stderr.String(), "invalid file size or type") {
		t.Fatalf("status %d, streams %q %q", status, &out, &stderr)
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

func TestRecordingHostBoundsEvents(t *testing.T) {
	h := &recordingHost{}
	var err error
	for range maxObservationBytes {
		err = h.ClearScreen()
		if err != nil {
			break
		}
	}
	if err == nil || !strings.Contains(err.Error(), "observation size limit") {
		t.Fatalf("error = %v, want observation size limit", err)
	}
	if !h.overflow || h.eventSize >= maxObservationBytes {
		t.Fatalf("event trace exceeded limit: overflow=%v size=%d", h.overflow, h.eventSize)
	}
	if err := h.ClearScreen(); err == nil {
		t.Fatal("accepted event after overflow")
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

func TestReplayClockFixture(t *testing.T) {
	root, m := testManifest(t)
	p := m.Probes[0]
	p.TimeoutMS = 5000
	p.Source = writeArtifact(t, root, p.Source.Path, "algoritmo \"clock\"\ninicio\ncronometro on\ncronometro off\nfimalgoritmo")
	p.Implementation.Expected.Stdout = writeArtifact(t, root, "clock-output.txt", "\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 16 ms.\n")
	clock := writeArtifact(t, root, "clock.json", "{\"nowMS\":[0,16]}\n")
	host := writeArtifact(t, root, "host.json", "[{\"operation\":\"now\"},{\"operation\":\"now\"}]\n")
	p.Implementation.Expected.Clock = &clock
	p.Implementation.Expected.HostTrace = &host
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	prefix := []string{"-test.run=^TestReplayObserverChild$", "--"}
	if err := replayProbe(root, p, "unused", prefix, executable); err != nil {
		t.Fatal(err)
	}
}
