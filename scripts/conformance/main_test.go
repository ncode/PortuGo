package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestCommandErrors(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{nil, {"unknown"}, {"validate", "--unknown"}, {"stage"}, {"capture"}} {
		var out, stderr bytes.Buffer
		if status := run(args, &out, &stderr); status != 2 {
			t.Errorf("%v: status=%d, stderr=%s", args, status, &stderr)
		}
		if stderr.Len() == 0 {
			t.Errorf("%v: missing error", args)
		}
	}
}

func TestCommandRejectsInvalidManifest(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	m.Reference.ExecutableSHA256 = ""
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, root, "manifest.json", string(data))
	var out, stderr bytes.Buffer
	if status := run([]string{"validate", "--root", root, "--manifest", "manifest.json"}, &out, &stderr); status != 1 {
		t.Fatalf("status=%d, stderr=%s", status, &stderr)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("reference provenance")) {
		t.Fatal(stderr.String())
	}
}

func TestLoadRejectsUnknownFieldsAndTrailingJSON(t *testing.T) {
	t.Parallel()
	for _, content := range []string{`{"versoin":1}`, `{"version":1} {"version":2}`} {
		root := t.TempDir()
		writeArtifact(t, root, "manifest.json", content)
		if _, err := loadManifest(root, "manifest.json"); err == nil {
			t.Errorf("accepted %s", content)
		}
	}
}

func TestCommandStagesSelectedProbe(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, root, "manifest.json", string(data))
	stage := t.TempDir() + string(os.PathSeparator) + "recording"
	var out, stderr bytes.Buffer
	if status := run([]string{"stage", "--root", root, "--manifest", "manifest.json", "--probe", "output", "--staging", stage}, &out, &stderr); status != 0 {
		t.Fatalf("status=%d, stderr=%s", status, &stderr)
	}
}
