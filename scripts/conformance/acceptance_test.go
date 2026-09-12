package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestImplementationAcceptance(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, want string
		mutate     func(*testing.T, string, *manifest)
	}{
		{name: "verified behavior with unfinished handoff"},
		{name: "pending behavior", want: "pending implementation", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Implementation.State = "pending" }},
		{name: "mismatched expectation", want: "stdout differs", mutate: func(t *testing.T, root string, m *manifest) {
			m.Probes[0].Implementation.Expected.Stdout = writeArtifact(t, root, "different.txt", "1\n")
		}},
		{name: "stale requirement", want: "stale trace link", mutate: func(_ *testing.T, _ string, m *manifest) { m.Inventory[0].Link += " missing" }},
		{name: "stale test", want: "stale test link", mutate: func(_ *testing.T, _ string, m *manifest) {
			m.Probes[0].Implementation.Tests[0] = "output_test.go#TestRemoved"
		}},
		{name: "stale task", want: "stale task link", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Tasks[0] = "10.99" }},
		{name: "missing tests", want: "missing test link", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Implementation.Tests = nil }},
		{name: "unsupported accepted example", want: "pending implementation", mutate: func(t *testing.T, root string, m *manifest) {
			name := "examples/output.alg"
			id := "example." + hashBytes([]byte(name))[:12]
			src, err := readArtifact(root, m.Probes[0].Source)
			if err != nil {
				t.Fatal(err)
			}
			catalog := []bundledExample{{ID: id, Name: name, SHA256: m.Probes[0].Source.SHA256, Bytes: len(src), Disposition: "accepted"}}
			writeJSON(t, root, "official-examples.json", catalog)
			m.Inventory = append(m.Inventory, inventoryItem{ID: id, Kind: "example", Link: "official-examples.json", Probes: []string{"output"}})
			m.Probes[0].Implementation.State = "pending"
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			m.Probes[0].Implementation.State = "verified"
			m.Probes[0].Implementation.Tests = []string{"output_test.go#TestOutput"}
			writeArtifact(t, root, "tasks.md", "- [x] 10.1 Implement output\n- [ ] 17.9 Report\n- [ ] 17.10 Handoff\n")
			if tt.mutate != nil {
				tt.mutate(t, root, &m)
			}
			err := validate(root, m, "implementation-acceptance", nil)
			if tt.want == "" && err != nil || tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}

func writeJSON(t *testing.T, root, path string, value any) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, root, path, string(data))
}

func TestAcceptanceRequiresQualityResults(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	p := &m.Probes[0]
	p.Implementation.State = "verified"
	p.Implementation.Tests = []string{"output_test.go#TestOutput"}
	p.Evidence.State = "not-applicable"
	p.Evidence.Review = &review{Reason: "Project-specific contract", Link: "review.md"}
	writeArtifact(t, root, "review.md", "Project-specific contract.\n")
	writeJSON(t, root, "manifest.json", m)
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		mode   string
		status int
	}{
		{"evidence", 0}, {"incremental", 0}, {"implementation-acceptance", 1}, {"release", 1},
	} {
		t.Run(tt.mode, func(t *testing.T) {
			var out, stderr bytes.Buffer
			status := run([]string{"validate", "--root", root, "--manifest", "manifest.json", "--previous", "manifest.json", "--mode", tt.mode, "--candidate", executable}, &out, &stderr)
			if status != tt.status || tt.status == 1 && !strings.Contains(stderr.String(), "quality results") {
				t.Fatalf("status = %d, stderr = %s; want status %d and quality results only for final modes", status, &stderr, tt.status)
			}
		})
	}
}
