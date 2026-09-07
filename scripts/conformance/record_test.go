package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareRecording(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	stage := filepath.Join(t.TempDir(), "recording")
	p := m.Probes[0]
	if err := prepareRecording(root, p, stage); err != nil {
		t.Fatal(err)
	}
	for name, a := range map[string]artifact{"source.alg": p.Source, "input.txt": p.Input} {
		got, err := os.ReadFile(filepath.Join(stage, name))
		if err != nil {
			t.Fatal(err)
		}
		want, err := readArtifact(root, a)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("staged %s differs", name)
		}
	}
	if _, err := os.Stat(filepath.Join(stage, "instructions.txt")); err != nil {
		t.Fatal(err)
	}
	if err := prepareRecording(root, p, stage); err == nil {
		t.Fatal("must not overwrite prior recording")
	}
	if err := prepareRecording(root, p, filepath.Join(root, "recording")); err == nil {
		t.Fatal("raw recording must remain outside repository")
	}
}

func TestCaptureRecording(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	stage := filepath.Join(t.TempDir(), "recording")
	if err := prepareRecording(root, m.Probes[0], stage); err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, stage, "raw.txt", "Início da execução\r\n 1\r\n\r\nFim da execução.\r\n")
	got, err := captureRecording(stage, true, "2026-09-07T12:00:00Z", "panel-v1", false)
	if err != nil {
		t.Fatal(err)
	}
	if got.State != "recorded" || got.Accepted == nil || !*got.Accepted {
		t.Fatalf("evidence = %+v", got)
	}
	output, err := readArtifact(stage, got.Normalized)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != " 1\n" {
		t.Fatalf("output = %q", output)
	}
	if _, err := captureRecording(stage, false, "2026-09-07T12:00:00Z", "text-v1", true); err == nil {
		t.Fatal("GUI-only capture requires screenshot and transcription")
	}
	if _, err := captureRecording(stage, true, "invalid", "panel-v1", false); err == nil {
		t.Fatal("invalid capture date")
	}
}

func TestCaptureGeneratedBytes(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	p := m.Probes[0]
	p.Implementation.Expected.Generated = []generatedFile{{Path: "generated.dat"}}
	stage := filepath.Join(t.TempDir(), "recording")
	if err := prepareRecording(root, p, stage); err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, stage, "raw.txt", "Início da execução\r\n 1\r\n\r\nFim da execução.\r\n")
	want := writeArtifact(t, stage, "generated.dat", "\xe9\r\n ")
	e, err := captureRecording(stage, true, "2026-09-07T12:00:00Z", "panel-v1", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(e.Generated) != 1 || e.Generated[0].Content != want {
		t.Fatalf("generated evidence = %+v, want %+v", e.Generated, want)
	}
}
