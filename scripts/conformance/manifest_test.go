package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeArtifact(t *testing.T, root, path, content string) artifact {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, path)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return artifact{Path: path, SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(content)))}
}

func testManifest(t *testing.T) (string, manifest) {
	t.Helper()
	root := t.TempDir()
	writeArtifact(t, root, "specs/output/spec.md", "### Requirement: Output\n")
	writeArtifact(t, root, "tasks.md", "## 2. Evidence\n- [ ] 2.1 Record output\n## 10. Output\n- [ ] 10.1 Implement output\n")
	writeArtifact(t, root, "output_test.go", "package example\nimport \"testing\"\nfunc TestOutput(t *testing.T) {}\n")
	source := writeArtifact(t, root, "probes/output/source.alg", "algoritmo \"output\"\ninicio\nescreval(1)\nfimalgoritmo\n")
	input := writeArtifact(t, root, "probes/output/input.txt", "")
	raw := writeArtifact(t, root, "probes/output/raw.txt", "Início da execução\r\n 1\r\n\r\nFim da execução.\r\n")
	normalized := writeArtifact(t, root, "probes/output/output.txt", " 1\n")
	accepted := true
	m := manifest{
		Version: 1, RecorderVersion: "manual-v1", NormalizerVersion: "panel-v1",
		Reference: reference{Product: "VisuAlg", Version: "3.0.7.0", SourceURL: "https://example.invalid/reference", AcquiredDate: "2026-09-07", ArchiveSHA256: strings.Repeat("a", 64), ExecutableSHA256: strings.Repeat("b", 64), OS: "Windows amd64", Culture: "en-US"},
		TasksPath: "tasks.md", InventorySources: []string{"specs/output/spec.md"},
		Inventory: []inventoryItem{{ID: "requirement.output", Kind: "requirement", Link: "specs/output/spec.md#Requirement: Output", Probes: []string{"output"}}},
		Probes: []probe{{ID: "output", OwnerGroup: 10, Tasks: []string{"10.1"}, Source: source, Input: input, TimeoutMS: 1000,
			Evidence:       evidence{State: "recorded", Accepted: &accepted, CapturedAt: "2026-09-07T12:00:00Z", Raw: raw, Normalized: normalized, Normalizer: "panel-v1"},
			Implementation: implementation{State: "pending", Expected: observation{Stdout: normalized}},
		}},
	}
	return root, m
}

func TestManifestValidation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, mode, want string
		mutate           func(*testing.T, string, *manifest)
	}{
		{name: "recorded pending evidence", mode: "evidence"},
		{name: "pending project safeguard", mode: "evidence", mutate: func(t *testing.T, root string, m *manifest) {
			writeArtifact(t, root, "review.md", "Project safeguard: no reference behavior. Implementation belongs to group 10.\n")
			m.Probes[0].Evidence.State = "not-applicable"
			m.Probes[0].Evidence.Review = &review{Reason: "Project-specific safeguard with future implementation", Link: "review.md"}
			m.Probes[0].Source, m.Probes[0].Input = artifact{}, artifact{}
		}},
		{name: "pending acceptance", mode: "implementation-acceptance", want: "pending implementation"},
		{name: "missing provenance", mode: "evidence", want: "reference provenance", mutate: func(_ *testing.T, _ string, m *manifest) { m.Reference.ExecutableSHA256 = "" }},
		{name: "invalid date", mode: "evidence", want: "reference provenance", mutate: func(_ *testing.T, _ string, m *manifest) { m.Reference.AcquiredDate = "yesterday" }},
		{name: "altered bytes", mode: "evidence", want: "hash mismatch", mutate: func(t *testing.T, root string, _ *manifest) {
			writeArtifact(t, root, "probes/output/source.alg", "changed")
		}},
		{name: "missing evidence", mode: "evidence", want: "unrecorded evidence", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Evidence.State = "unrecorded" }},
		{name: "generated reference corruption", mode: "evidence", want: "hash mismatch", mutate: func(t *testing.T, root string, m *manifest) {
			a := writeArtifact(t, root, "generated.dat", "reference bytes")
			m.Probes[0].Evidence.Generated = []generatedFile{{Path: "result.dat", Content: a}}
			m.Probes[0].Implementation.Expected.Generated = []generatedFile{{Path: "result.dat", Content: a}}
			writeArtifact(t, root, "generated.dat", "changed bytes")
		}},
		{name: "candidate replaces generated reference", mode: "evidence", want: "generated reference", mutate: func(t *testing.T, root string, m *manifest) {
			a := writeArtifact(t, root, "generated.dat", "reference bytes")
			b := writeArtifact(t, root, "candidate.dat", "candidate bytes")
			m.Probes[0].Evidence.Generated = []generatedFile{{Path: "result.dat", Content: a}}
			m.Probes[0].Implementation.Expected.Generated = []generatedFile{{Path: "result.dat", Content: b}}
		}},
		{name: "missing disposition", mode: "evidence", want: "acceptance", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Evidence.Accepted = nil }},
		{name: "bad normalization", mode: "evidence", want: "normalized evidence", mutate: func(t *testing.T, root string, m *manifest) {
			m.Probes[0].Evidence.Normalized = writeArtifact(t, root, "probes/output/output.txt", "1\n")
		}},
		{name: "parent path", mode: "evidence", want: "unsafe path", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Source.Path = "../source.alg" }},
		{name: "Windows path", mode: "evidence", want: "unsafe path", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Source.Path = `C:\source.alg` }},
		{name: "duplicate probe", mode: "evidence", want: "duplicate probe", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes = append(m.Probes, m.Probes[0]) }},
		{name: "missing task", mode: "evidence", want: "task link", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Tasks = []string{"10.99"} }},
		{name: "wrong owner", mode: "evidence", want: "owner group", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].OwnerGroup = 9 }},
		{name: "stale requirement", mode: "evidence", want: "trace link", mutate: func(_ *testing.T, _ string, m *manifest) {
			m.Inventory[0].Link = "specs/output/spec.md#Requirement: Missing"
		}},
		{name: "untraced requirement", mode: "evidence", want: "untraced requirement", mutate: func(t *testing.T, root string, _ *manifest) {
			writeArtifact(t, root, "specs/output/spec.md", "### Requirement: Output\n### Requirement: Input\n")
		}},
		{name: "omitted specification", mode: "evidence", want: "untraced requirement", mutate: func(t *testing.T, root string, _ *manifest) {
			writeArtifact(t, root, "specs/input/spec.md", "### Requirement: Input\n")
		}},
		{name: "wrong requirement kind", mode: "evidence", want: "untraced requirement", mutate: func(_ *testing.T, _ string, m *manifest) {
			m.Inventory[0].Kind = "assumption"
		}},
		{name: "missing probe link", mode: "evidence", want: "probe link", mutate: func(_ *testing.T, _ string, m *manifest) { m.Inventory[0].Probes = []string{"missing"} }},
		{name: "unreviewed exemption", mode: "evidence", want: "review", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Evidence.State = "not-applicable" }},
		{name: "GUI without screenshot", mode: "evidence", want: "GUI evidence", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Evidence.GUIOnly = true }},
		{name: "verified without tests", mode: "evidence", want: "test link", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].Implementation.State = "verified" }},
		{name: "stale test symbol", mode: "evidence", want: "test link", mutate: func(_ *testing.T, _ string, m *manifest) {
			m.Probes[0].Implementation.State = "verified"
			m.Probes[0].Implementation.Tests = []string{"output_test.go#TestMissing"}
		}},
		{name: "completed owner pending", mode: "incremental", want: "completed group", mutate: func(t *testing.T, root string, _ *manifest) {
			writeArtifact(t, root, "tasks.md", "- [ ] 2.1 Record\n- [x] 10.1 Implement\n")
		}},
		{name: "unbounded replay", mode: "evidence", want: "budget", mutate: func(_ *testing.T, _ string, m *manifest) { m.Probes[0].TimeoutMS = 0 }},
		{name: "prohibited executable", mode: "evidence", want: "prohibited artifact", mutate: func(t *testing.T, root string, _ *manifest) { writeArtifact(t, root, "reference.exe", "binary") }},
		{name: "renamed executable", mode: "evidence", want: "prohibited artifact", mutate: func(t *testing.T, root string, m *manifest) {
			a := writeArtifact(t, root, "renamed.dat", "reference binary")
			m.Reference.ExecutableSHA256 = a.SHA256
		}},
		{name: "unknown mode", mode: "anything", want: "validation mode"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			if tt.mutate != nil {
				tt.mutate(t, root, &m)
			}
			err := validate(root, m, tt.mode, nil)
			if tt.want == "" && err != nil {
				t.Fatal(err)
			}
			if tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestManifestPreservesVerifiedHistory(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	previous := m
	previous.Probes = append([]probe(nil), m.Probes...)
	previous.Probes[0].Implementation.State = "verified"
	previous.Probes[0].Implementation.Tests = []string{"output_test.go#TestOutput"}
	if err := validate(root, m, "evidence", &previous); err == nil || !strings.Contains(err.Error(), "downgrade") {
		t.Fatalf("error = %v, want unreviewed downgrade", err)
	}
	m.Probes = nil
	if err := validate(root, m, "evidence", &previous); err == nil || !strings.Contains(err.Error(), "removed verified") {
		t.Fatalf("error = %v, want removed verified probe", err)
	}
}

func TestManifestReferenceDisposition(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name        string
		accepted    bool
		exitCode    int
		diagnostics []diagnostic
		valid       bool
	}{
		{name: "accepted success", accepted: true, valid: true},
		{name: "accepted failure", accepted: true, exitCode: 1},
		{name: "accepted diagnostics", accepted: true, diagnostics: []diagnostic{{Code: "E004", Line: 4}}},
		{name: "rejected success"},
		{name: "rejected without diagnostics", exitCode: 1},
		{name: "rejected wrong exit", exitCode: 2, diagnostics: []diagnostic{{Code: "E004", Line: 4}}},
		{name: "rejected with coverage", exitCode: 1, diagnostics: []diagnostic{{Code: "E004", Line: 4}}, valid: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			p := &m.Probes[0]
			p.Evidence.Accepted = &tt.accepted
			if !tt.accepted {
				a := writeArtifact(t, root, "error.txt", "Rejected on line 4\n")
				p.Evidence.Raw, p.Evidence.Normalized, p.Evidence.Normalizer = a, a, "bytes-v1"
			}
			p.Implementation.State = "verified"
			p.Implementation.Tests = []string{"output_test.go#TestOutput"}
			p.Implementation.Expected.ExitCode = tt.exitCode
			p.Implementation.Expected.Diagnostics = tt.diagnostics
			err := validate(root, m, "implementation-acceptance", nil)
			if tt.valid && err != nil {
				t.Fatal(err)
			}
			if !tt.valid && (err == nil || !strings.Contains(err.Error(), "reference disposition")) {
				t.Fatalf("error = %v, want reference disposition mismatch", err)
			}
		})
	}
}

func TestManifestGeneratedInventory(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	a := writeArtifact(t, root, "first.dat", "first")
	b := writeArtifact(t, root, "second.dat", "second")
	p := &m.Probes[0]
	p.Evidence.Generated = []generatedFile{{Path: "first.dat", Content: a}, {Path: "second.dat", Content: b}}
	p.Implementation.Expected.Generated = append([]generatedFile(nil), p.Evidence.Generated...)
	if err := validate(root, m, "evidence", nil); err != nil {
		t.Fatal(err)
	}
	p.Implementation.Expected.Generated[1] = p.Implementation.Expected.Generated[0]
	if err := validate(root, m, "evidence", nil); err == nil || !strings.Contains(err.Error(), "generated") {
		t.Fatalf("error = %v, want duplicate generated path rejection", err)
	}
}

func TestManifestRejectsSymlink(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	outside := t.TempDir()
	target := writeArtifact(t, outside, "source.alg", "outside")
	// Windows does not normally grant unprivileged symlink creation.
	if err := os.Symlink(filepath.Join(outside, target.Path), filepath.Join(root, "link")); err != nil {
		t.Skip(err)
	}
	m.Probes[0].Source.Path = "link"
	if err := validate(root, m, "evidence", nil); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("error = %v, want symlink rejection", err)
	}
}
