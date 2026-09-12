package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func testQualityReport(t *testing.T, commit string) (string, qualityReport) {
	t.Helper()
	root := t.TempDir()
	report := qualityReport{Version: 1, Commit: commit}
	for _, name := range []string{"build", "gofmt", "vet", "staticcheck", "golangci-lint", "tests", "race", "windows", "macos", "linux", "fuzz-lexer", "fuzz-parser", "openspec"} {
		evidence := writeArtifact(t, root, name+".txt", "check="+name+"\nstatus=pass\n")
		report.Results = append(report.Results, qualityResult{Check: name, Status: "pass", Evidence: evidence})
	}
	return root, report
}

func TestQualityResults(t *testing.T) {
	t.Parallel()
	commit := strings.Repeat("a", 40)
	for _, tt := range []struct {
		name, want string
		mutate     func(*testing.T, string, *qualityReport)
	}{
		{name: "complete"},
		{name: "missing", want: "missing quality results", mutate: func(_ *testing.T, _ string, r *qualityReport) { r.Results = nil }},
		{name: "stale candidate", want: "different candidate", mutate: func(_ *testing.T, _ string, r *qualityReport) { r.Commit = strings.Repeat("b", 40) }},
		{name: "unknown version", want: "version", mutate: func(_ *testing.T, _ string, r *qualityReport) { r.Version++ }},
		{name: "failed", want: "did not pass", mutate: func(_ *testing.T, _ string, r *qualityReport) { r.Results[0].Status = "fail" }},
		{name: "skipped", want: "did not pass", mutate: func(_ *testing.T, _ string, r *qualityReport) { r.Results[0].Status = "skip" }},
		{name: "duplicate", want: "duplicate quality check", mutate: func(_ *testing.T, _ string, r *qualityReport) { r.Results = append(r.Results, r.Results[0]) }},
		{name: "unknown check", want: "unknown or duplicate", mutate: func(_ *testing.T, _ string, r *qualityReport) { r.Results[0].Check = "other" }},
		{name: "altered evidence", want: "hash mismatch", mutate: func(t *testing.T, root string, r *qualityReport) {
			writeArtifact(t, root, r.Results[0].Evidence.Path, "changed")
		}},
		{name: "missing evidence", want: "evidence", mutate: func(_ *testing.T, _ string, r *qualityReport) { r.Results[0].Evidence.Path = "missing.txt" }},
		{name: "empty evidence", want: "empty evidence", mutate: func(t *testing.T, root string, r *qualityReport) {
			r.Results[0].Evidence = writeArtifact(t, root, "empty.txt", " \n")
		}},
		{name: "unsafe evidence", want: "unsafe path", mutate: func(_ *testing.T, _ string, r *qualityReport) { r.Results[0].Evidence.Path = "../outside.txt" }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, report := testQualityReport(t, commit)
			if tt.mutate != nil {
				tt.mutate(t, root, &report)
			}
			err := checkQualityReport(root, report, commit)
			if tt.want == "" && err != nil || tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
	root, report := testQualityReport(t, commit)
	for i, result := range report.Results {
		t.Run("missing_"+result.Check, func(t *testing.T) {
			partial := report
			partial.Results = slices.Delete(slices.Clone(report.Results), i, i+1)
			if err := checkQualityReport(root, partial, commit); err == nil || !strings.Contains(err.Error(), "missing quality results: "+result.Check) {
				t.Fatalf("error = %v, want missing %s", err, result.Check)
			}
		})
	}
}

func testAcceptanceCandidate(t *testing.T, tasks, output string) (string, string) {
	t.Helper()
	root, m := testManifest(t)
	m.Probes[0].Implementation.State = "verified"
	m.Probes[0].Implementation.Tests = []string{"output_test.go#TestOutput"}
	// Newly built executable startup can be slow under concurrent race builds.
	// Keep the existing maximum replay watchdog; expiration still fails the test.
	m.Probes[0].TimeoutMS = 30000
	writeArtifact(t, root, "tasks.md", tasks)
	writeArtifact(t, root, "go.mod", "module example.invalid/fixture\n\ngo 1.22\n")
	writeArtifact(t, root, "cmd/portugol/main.go", fmt.Sprintf("package main\nimport \"fmt\"\nfunc main() { fmt.Print(%q) }\n", output))
	writeArtifact(t, root, ".gitignore", "ignored.go\n")
	writeJSON(t, root, "manifest.json", m)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, args := range [][]string{{"init"}, {"add", "."}, {"-c", "user.name=Test", "-c", "user.email=test@example.invalid", "-c", "commit.gpgsign=false", "-c", "core.hooksPath=.git/no-hooks", "commit", "-m", "Synthetic acceptance candidate"}} {
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("initialize candidate: %v: %s", err, out)
		}
	}
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	cmd.Dir = root
	commit, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	qualityRoot, report := testQualityReport(t, strings.TrimSpace(string(commit)))
	writeJSON(t, qualityRoot, "quality.json", report)
	return root, filepath.Join(qualityRoot, "quality.json")
}

func TestQualityCandidate(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, want string
		mutate     func(*testing.T, string, string)
	}{
		{name: "matching clean candidate"},
		{name: "tracked change", want: "clean checkout", mutate: func(t *testing.T, root, _ string) { writeArtifact(t, root, "cmd/portugol/main.go", "changed") }},
		{name: "untracked source", want: "clean checkout", mutate: func(t *testing.T, root, _ string) { writeArtifact(t, root, "new.go", "package example") }},
		{name: "ignored source", want: "clean checkout", mutate: func(t *testing.T, root, _ string) { writeArtifact(t, root, "ignored.go", "package example") }},
		{name: "unknown report field", want: "unknown field", mutate: func(t *testing.T, _ string, name string) {
			writeArtifact(t, filepath.Dir(name), filepath.Base(name), `{"versoin":1}`)
		}},
		{name: "trailing report JSON", want: "trailing quality results", mutate: func(t *testing.T, _ string, name string) {
			writeArtifact(t, filepath.Dir(name), filepath.Base(name), `{} {}`)
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, name := testAcceptanceCandidate(t, "- [x] 10.1 Implement\n", " 1\n")
			if tt.mutate != nil {
				tt.mutate(t, root, name)
			}
			err := validateQuality(root, name, "")
			if tt.want == "" && err != nil || tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
			if err := validateQuality(root, name, "other-candidate"); err == nil || !strings.Contains(err.Error(), "current-checkout candidate") {
				t.Fatalf("custom candidate: %v", err)
			}
		})
	}
}

func TestReleaseTaskCompletion(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, mode, tasks, want string
		status                  int
	}{
		{name: "acceptance before handoff", mode: "implementation-acceptance", tasks: "- [x] 10.1 Implement\n- [ ] 17.9 Report\n- [ ] 17.10 Handoff\n", want: "Remaining implementation/handoff tasks: [17.10 17.9]"},
		{name: "release before handoff", mode: "release", tasks: "- [x] 10.1 Implement\n- [ ] 17.9 Report\n- [ ] 17.10 Handoff\n", want: "incomplete release tasks", status: 1},
		{name: "release before implementation", mode: "release", tasks: "- [ ] 10.1 Implement\n- [x] 17.9 Report\n- [x] 17.10 Handoff\n", want: "incomplete release tasks", status: 1},
		{name: "release before archive and tag", mode: "release", tasks: "- [x] 10.1 Implement\n- [x] 17.9 Report\n- [x] 17.10 Handoff\n"},
		{name: "duplicate completion", mode: "release", tasks: "- [ ] 10.1 Implement\n- [x] 10.1 Implement\n", want: "duplicate task", status: 1},
		{name: "invalid completion", mode: "release", tasks: "- [x] 10.1 Implement\n- [?] 17.9 Report\n", want: "invalid task checkbox", status: 1},
		{name: "indented incomplete task", mode: "release", tasks: "- [x] 10.1 Implement\n  - [ ] 17.9 Report\n", want: "incomplete release tasks", status: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root, quality := testAcceptanceCandidate(t, tt.tasks, " 1\n")
			var out, stderr bytes.Buffer
			status := run([]string{"validate", "--root", root, "--manifest", "manifest.json", "--previous", "manifest.json", "--mode", tt.mode, "--quality", quality}, &out, &stderr)
			if status != tt.status || !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("status = %d, stdout = %s, stderr = %s; want %d, %q", status, &out, &stderr, tt.status, tt.want)
			}
			if status == 0 && !strings.Contains(out.String(), `"state":"verified"`) {
				t.Fatalf("missing successful replay: %s", &out)
			}
		})
	}
}

func TestAcceptanceRejectsReplayMismatch(t *testing.T) {
	t.Parallel()
	root, quality := testAcceptanceCandidate(t, "- [x] 10.1 Implement\n", "1\n")
	var out, stderr bytes.Buffer
	status := run([]string{"validate", "--root", root, "--manifest", "manifest.json", "--previous", "manifest.json", "--mode", "implementation-acceptance", "--quality", quality}, &out, &stderr)
	if status != 1 || !strings.Contains(stderr.String(), "verified reference regression") || !strings.Contains(out.String(), "stdout mismatch") {
		t.Fatalf("status = %d, stdout = %s, stderr = %s; want replay mismatch", status, &out, &stderr)
	}
}
