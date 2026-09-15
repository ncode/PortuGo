package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func TestPrepareRecordingRejectsRecorderOwnedPaths(t *testing.T) {
	for _, name := range recorderOwnedPaths {
		t.Run(name, func(t *testing.T) {
			root, m := testManifest(t)
			p := m.Probes[0]
			p.Implementation.Expected.Generated = []generatedFile{{Path: name}}
			if err := prepareRecording(root, p, filepath.Join(t.TempDir(), "recording")); err == nil || !strings.Contains(err.Error(), "recorder-owned path") {
				t.Fatalf("error = %v, want recorder-owned path", err)
			}
		})
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
	got, err := captureRecording(root, stage, true, "2026-09-07T12:00:00Z", "panel-v1", false)
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
	if _, err := captureRecording(root, stage, false, "2026-09-07T12:00:00Z", "text-v1", true); err == nil {
		t.Fatal("GUI-only capture requires screenshot and transcription")
	}
	if _, err := captureRecording(root, stage, true, "invalid", "panel-v1", false); err == nil {
		t.Fatal("invalid capture date")
	}
}

func TestCaptureRejectsLooseStagedJSON(t *testing.T) {
	for _, tt := range []struct {
		name, want string
	}{
		{name: "unknown field", want: "unknown field"},
		{name: "trailing JSON", want: "trailing staged JSON"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root, m := testManifest(t)
			stage := filepath.Join(t.TempDir(), "recording")
			if err := prepareRecording(root, m.Probes[0], stage); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(stage, "staged.json"))
			if err != nil {
				t.Fatal(err)
			}
			if tt.name == "unknown field" {
				text := strings.TrimSuffix(strings.TrimSpace(string(data)), "}") + `,"extra":true}`
				data = []byte(text)
			} else {
				data = append(bytes.TrimSpace(data), []byte("\n{}")...)
			}
			if err := os.WriteFile(filepath.Join(stage, "staged.json"), data, 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := captureRecording(root, stage, true, "2026-09-07T12:00:00Z", "panel-v1", false); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestCaptureRejectsDuplicateStagedPaths(t *testing.T) {
	for _, tt := range []struct {
		name string
		set  func(*stagedRecording)
		want string
	}{
		{name: "generated", set: func(staged *stagedRecording) { staged.Generated = []string{"result.dat", "result.dat"} }, want: "duplicate staged generated path"},
		{name: "absent", set: func(staged *stagedRecording) { staged.Absent = []string{"missing.dat", "missing.dat"} }, want: "duplicate staged absent path"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root, m := testManifest(t)
			p := m.Probes[0]
			if tt.name == "generated" {
				p.Implementation.Expected.Generated = []generatedFile{{Path: "result.dat"}}
			} else {
				p.Implementation.Expected.Absent = []string{"missing.dat"}
			}
			stage := filepath.Join(t.TempDir(), "recording")
			if err := prepareRecording(root, p, stage); err != nil {
				t.Fatal(err)
			}
			writeArtifact(t, stage, "raw.txt", "Início da execução\r\n 1\r\n\r\nFim da execução.\r\n")
			if tt.name == "generated" {
				writeArtifact(t, stage, "result.dat", "generated")
			}
			data, err := os.ReadFile(filepath.Join(stage, "staged.json"))
			if err != nil {
				t.Fatal(err)
			}
			var staged stagedRecording
			if err := json.Unmarshal(data, &staged); err != nil {
				t.Fatal(err)
			}
			tt.set(&staged)
			data, err = json.Marshal(staged)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(stage, "staged.json"), append(data, '\n'), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := captureRecording(root, stage, true, "2026-09-07T12:00:00Z", "panel-v1", false); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
			for _, name := range []string{"normalized.txt", "evidence.json"} {
				if _, err := os.Stat(filepath.Join(stage, name)); !os.IsNotExist(err) {
					t.Fatalf("capture created %s for duplicate staged path", name)
				}
			}
		})
	}
}

func TestPrepareRecordingRejectsDuplicateStagedPaths(t *testing.T) {
	for _, tt := range []struct {
		name string
		set  func(*probe)
		want string
	}{
		{name: "generated", set: func(p *probe) {
			p.Implementation.Expected.Generated = []generatedFile{{Path: "result.dat"}, {Path: "result.dat"}}
		}, want: "duplicate staged generated path"},
		{name: "absent", set: func(p *probe) { p.Implementation.Expected.Absent = []string{"missing.dat", "missing.dat"} }, want: "duplicate staged absent path"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root, m := testManifest(t)
			p := m.Probes[0]
			tt.set(&p)
			if err := prepareRecording(root, p, filepath.Join(t.TempDir(), "recording")); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestCaptureRejectsSymlinkedStage(t *testing.T) {
	root, _ := testManifest(t)
	target := filepath.Join(root, "capture-target")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	stage := filepath.Join(t.TempDir(), "recording")
	if err := os.Symlink(target, stage); err != nil {
		t.Skip(err)
	}
	var out, stderr bytes.Buffer
	if status := run([]string{"capture", "--root", root, "--staging", stage, "--accepted", "true", "--captured-at", "2026-09-07T12:00:00Z"}, &out, &stderr); status != 1 || !strings.Contains(stderr.String(), "must not be a symlink") {
		t.Fatalf("status = %d, stderr = %q, want symlink rejection", status, stderr.String())
	}
	for _, name := range []string{"normalized.txt", "evidence.json"} {
		if _, err := os.Stat(filepath.Join(target, name)); !os.IsNotExist(err) {
			t.Fatalf("capture created %s through symlink", name)
		}
	}
}

func TestCaptureRejectsNonPrivateStage(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows staging privacy is enforced by ACLs")
	}
	root, m := testManifest(t)
	stage := filepath.Join(t.TempDir(), "recording")
	if err := prepareRecording(root, m.Probes[0], stage); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(stage, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := captureRecording(root, stage, true, "2026-09-07T12:00:00Z", "panel-v1", false); err == nil || !strings.Contains(err.Error(), "stage must be private") {
		t.Fatalf("error = %v, want private-stage rejection", err)
	}
	for _, name := range []string{"normalized.txt", "evidence.json"} {
		if _, err := os.Stat(filepath.Join(stage, name)); !os.IsNotExist(err) {
			t.Fatalf("capture created %s in non-private stage", name)
		}
	}
}

func TestCaptureRejectsRecorderOwnedGeneratedPath(t *testing.T) {
	root, m := testManifest(t)
	stage := filepath.Join(t.TempDir(), "recording")
	if err := prepareRecording(root, m.Probes[0], stage); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(stage, "staged.json"))
	if err != nil {
		t.Fatal(err)
	}
	var staged stagedRecording
	if err := json.Unmarshal(data, &staged); err != nil {
		t.Fatal(err)
	}
	staged.Generated = []string{"instructions.txt"}
	data, err = json.Marshal(staged)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(stage, "staged.json"), append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := captureRecording(root, stage, true, "2026-09-07T12:00:00Z", "panel-v1", false); err == nil || !strings.Contains(err.Error(), "recorder-owned path") {
		t.Fatalf("error = %v, want recorder-owned path", err)
	}
}

func TestCaptureRejectsUndeclaredStageFile(t *testing.T) {
	root, m := testManifest(t)
	stage := filepath.Join(t.TempDir(), "recording")
	if err := prepareRecording(root, m.Probes[0], stage); err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, stage, "raw.txt", "Início da execução\r\n 1\r\n\r\nFim da execução.\r\n")
	if err := os.WriteFile(filepath.Join(stage, "unexpected.txt"), []byte("private\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := captureRecording(root, stage, true, "2026-09-07T12:00:00Z", "panel-v1", false); err == nil || !strings.Contains(err.Error(), "undeclared recording file") {
		t.Fatalf("error = %v, want undeclared file rejection", err)
	}
	for _, name := range []string{"normalized.txt", "evidence.json"} {
		if _, err := os.Stat(filepath.Join(stage, name)); !os.IsNotExist(err) {
			t.Fatalf("capture created %s after undeclared file", name)
		}
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
	e, err := captureRecording(root, stage, true, "2026-09-07T12:00:00Z", "panel-v1", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(e.Generated) != 1 || e.Generated[0].Content != want {
		t.Fatalf("generated evidence = %+v, want %+v", e.Generated, want)
	}
}

func TestCaptureAbsentFiles(t *testing.T) {
	for _, present := range []bool{false, true} {
		t.Run(fmt.Sprint(present), func(t *testing.T) {
			root, m := testManifest(t)
			p := m.Probes[0]
			p.Implementation.Expected.Absent = []string{"absent.dat"}
			stage := filepath.Join(t.TempDir(), "recording")
			if err := prepareRecording(root, p, stage); err != nil {
				t.Fatal(err)
			}
			writeArtifact(t, stage, "raw.txt", "Início da execução\r\n 1\r\n\r\nFim da execução.\r\n")
			if present {
				writeArtifact(t, stage, "absent.dat", "unexpected")
			}
			e, err := captureRecording(root, stage, true, "2026-09-07T12:00:00Z", "panel-v1", false)
			if present && (err == nil || !strings.Contains(err.Error(), "expected absent file")) {
				t.Fatalf("unexpected capture error=%v", err)
			}
			if !present && (err != nil || len(e.Absent) != 1 || e.Absent[0] != "absent.dat") {
				t.Fatalf("absence evidence=%v error=%v", e.Absent, err)
			}
		})
	}
}

func TestCaptureInitialFiles(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, content string
		generated     bool
		valid         bool
	}{
		{name: "unchanged input", content: "original", valid: true},
		{name: "changed input", content: "changed"},
		{name: "declared output", content: "changed", generated: true, valid: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			p := m.Probes[0]
			a := writeArtifact(t, root, "fixture.dat", "original")
			p.Files = []generatedFile{{Path: "data/input.dat", Content: a}}
			if tt.generated {
				p.Implementation.Expected.Generated = p.Files
			}
			stage := filepath.Join(t.TempDir(), "recording")
			if err := prepareRecording(root, p, stage); err != nil {
				t.Fatal(err)
			}
			writeArtifact(t, stage, "data/input.dat", tt.content)
			writeArtifact(t, stage, "raw.txt", "Início da execução\r\n 1\r\n\r\nFim da execução.\r\n")
			e, err := captureRecording(root, stage, true, "2026-09-07T12:00:00Z", "panel-v1", false)
			if tt.valid && err != nil {
				t.Fatal(err)
			}
			if !tt.valid && (err == nil || !strings.Contains(err.Error(), "hash mismatch")) {
				t.Fatalf("capture = %+v, error = %v, want changed input rejection", e, err)
			}
			if tt.generated && (len(e.Generated) != 1 || e.Generated[0].Content.SHA256 != hashBytes([]byte(tt.content))) {
				t.Fatalf("generated evidence = %+v", e.Generated)
			}
		})
	}
}
