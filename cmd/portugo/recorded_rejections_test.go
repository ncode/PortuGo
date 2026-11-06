package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRecordedRejectedExamples(t *testing.T) {
	if encoded := os.Getenv("PORTUGO_REJECTION_ARGS"); encoded != "" {
		if err := json.Unmarshal([]byte(encoded), &os.Args); err != nil {
			t.Fatal(err)
		}
		main()
		return
	}
	root, probes := recordedRejectionProbes(t)
	names := map[string]string{
		"bundled-24919c4a2e4c": "assignment",
		"bundled-49eb7f47065e": "character table",
		"bundled-f251f2cd0c22": "school",
		"bundled-f4fa0c514cf7": "repeat",
	}
	seen := 0
	for _, probe := range probes {
		name, ok := names[probe.ID]
		if !ok {
			continue
		}
		seen++
		t.Run(name, func(t *testing.T) {
			if err := rejectedExampleReady(probe); err != nil {
				t.Fatal(err)
			}
			testRecordedRejection(t, root, probe)
		})
	}
	if seen != len(names) {
		t.Fatalf("found %d probes, want %d", seen, len(names))
	}
}

func rejectedExampleReady(probe rejectionProbe) error {
	if probe.Implementation.State != "verified" || probe.Implementation.Expected.ExitCode != 1 || len(probe.Implementation.Expected.Diagnostics) != 1 || !slices.Contains(probe.Implementation.Tests, "cmd/portugo/recorded_rejections_test.go#TestRecordedRejectedExamples") {
		return fmt.Errorf("%s: rejected example is incomplete or has a mismatching rejection contract", probe.ID)
	}
	return nil
}

func TestRejectedExampleCompletionGate(t *testing.T) {
	_, probes := recordedRejectionProbes(t)
	for _, p := range probes {
		if p.ID != "bundled-24919c4a2e4c" {
			continue
		}
		p.Implementation.State = "verified"
		if err := rejectedExampleReady(p); err != nil {
			t.Fatal(err)
		}
		for _, change := range []func(*rejectionProbe){
			func(p *rejectionProbe) { p.Implementation.State = "pending" },
			func(p *rejectionProbe) { p.Implementation.Expected.ExitCode = 0 },
			func(p *rejectionProbe) { p.Implementation.Expected.Diagnostics = nil },
			func(p *rejectionProbe) { p.Implementation.Tests = nil },
		} {
			broken := p
			change(&broken)
			if err := rejectedExampleReady(broken); err == nil {
				t.Fatal("accepted-example success concealed an incomplete rejection")
			}
		}
		return
	}
	t.Fatal("missing rejection probe")
}

type rejectionProbe struct {
	ID     string `json:"id"`
	Source struct {
		Path string `json:"path"`
	} `json:"source"`
	Input struct {
		Path string `json:"path"`
	} `json:"input"`
	MaxSteps       uint64 `json:"maxSteps"`
	Implementation struct {
		State    string   `json:"state"`
		Tests    []string `json:"tests"`
		Expected struct {
			ExitCode int `json:"exitCode"`
			Stdout   struct {
				Path string `json:"path"`
			} `json:"stdout"`
			Diagnostics []struct {
				Code string `json:"code"`
				Line int    `json:"line"`
			} `json:"diagnostics"`
		} `json:"expected"`
	} `json:"implementation"`
}

func recordedRejectionProbes(t *testing.T) (string, []rejectionProbe) {
	t.Helper()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "testdata/conformance/visualg-3.0.7/manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m struct {
		Probes []rejectionProbe `json:"probes"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	return root, m.Probes
}

func testRecordedRejection(t *testing.T, root string, probe rejectionProbe) {
	t.Helper()
	budget := probe.MaxSteps
	if budget == 0 {
		budget = 1000000
	}

	source, err := os.ReadFile(filepath.Join(root, probe.Source.Path))
	if err != nil {
		t.Fatal(err)
	}
	input, err := os.ReadFile(filepath.Join(root, probe.Input.Path))
	if err != nil {
		t.Fatal(err)
	}
	want := probe.Implementation.Expected
	stdout, err := os.ReadFile(filepath.Join(root, want.Stdout.Path))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "source.alg")
	if err := os.WriteFile(path, source, 0o600); err != nil {
		t.Fatal(err)
	}
	exit, out, stderr := rejectionCommand(t, []string{"run", "--max-steps", fmt.Sprint(budget), path}, string(input))
	if exit != want.ExitCode {
		t.Errorf("exit = %d, want %d", exit, want.ExitCode)
	}
	if !bytes.Equal(out, stdout) {
		t.Errorf("stdout differs: got %d bytes, want %d", len(out), len(stdout))
	}
	matches := regexp.MustCompile(`:(\d+):(\d+): ([A-Z]\d{3}):`).FindAllStringSubmatch(string(stderr), -1)
	if len(matches) != len(want.Diagnostics) {
		t.Fatalf("diagnostic count = %d, want %d", len(matches), len(want.Diagnostics))
	}
	for index, diagnostic := range want.Diagnostics {
		line, err := strconv.Atoi(matches[index][1])
		if err != nil {
			t.Fatal(err)
		}
		column, err := strconv.Atoi(matches[index][2])
		if err != nil {
			t.Fatal(err)
		}
		if line != diagnostic.Line || matches[index][3] != diagnostic.Code || column < 1 {
			t.Errorf("diagnostic = %s line %d column %d, want %s line %d with a positive column", matches[index][3], line, column, diagnostic.Code, diagnostic.Line)
		}
	}
	if len(want.Diagnostics) == 0 && len(stderr) != 0 {
		t.Error("successful control wrote diagnostics")
	}
}

func TestRecordedRejectionControls(t *testing.T) {
	root, probes := recordedRejectionProbes(t)
	seen := 0
	for _, probe := range probes {
		if !strings.HasPrefix(probe.ID, "rejection-") {
			continue
		}
		seen++
		t.Run(strings.TrimPrefix(probe.ID, "rejection-"), func(t *testing.T) { testRecordedRejection(t, root, probe) })
	}
	if seen == 0 {
		t.Fatal("no recorded rejection controls")
	}
}

func rejectionCommand(t *testing.T, args []string, input string) (int, []byte, []byte) {
	t.Helper()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(append([]string{"portugo"}, args...))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, "-test.run=^TestRecordedRejectedExamples$")
	cmd.Env = append(os.Environ(), "PORTUGO_REJECTION_ARGS="+string(encoded))
	cmd.Dir = t.TempDir()
	cmd.Stdin = strings.NewReader(input)
	cmd.WaitDelay = time.Second
	out, stderr := &rejectionOutput{}, &rejectionOutput{}
	cmd.Stdout, cmd.Stderr = out, stderr
	err = cmd.Run()
	if ctx.Err() != nil || cmd.ProcessState == nil || out.overflow || stderr.overflow {
		t.Fatalf("CLI exceeded watchdog/output bounds: %v", err)
	}
	return cmd.ProcessState.ExitCode(), out.Bytes(), stderr.Bytes()
}

type rejectionOutput struct {
	bytes.Buffer
	overflow bool
}

func (b *rejectionOutput) Write(p []byte) (int, error) {
	if len(p) > (1<<20)-b.Len() {
		b.overflow = true
		return 0, fmt.Errorf("CLI output exceeds test limit")
	}
	return b.Buffer.Write(p)
}
