package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type replayResult struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Error string `json:"error,omitempty"`
}

func replayFailed(results []replayResult) bool {
	for _, r := range results {
		if r.State == "verified" && r.Error != "" {
			return true
		}
	}
	return false
}

// Separate stdout/stderr buffers are each written by a single os/exec copier.
type boundedOutput struct {
	buffer   bytes.Buffer
	overflow bool
	cancel   context.CancelFunc
}

func (b *boundedOutput) Write(p []byte) (int, error) {
	if len(p) > (1<<20)-b.buffer.Len() {
		b.overflow = true
		b.cancel()
		return 0, fmt.Errorf("output limit exceeded")
	}
	return b.buffer.Write(p)
}

func replayProbe(root string, p probe, executable string, prefix []string, observer string) (runErr error) {
	want := p.Implementation.Expected
	observe := want.State != nil || want.HostTrace != nil
	if observe && observer == "" {
		return fmt.Errorf("missing state/host observation adapter")
	}
	if p.TimeoutMS < 1 || p.TimeoutMS > 30000 {
		return fmt.Errorf("invalid replay budget")
	}
	source, err := readArtifact(root, p.Source)
	if err != nil {
		return err
	}
	input, err := readArtifact(root, p.Input)
	if err != nil {
		return err
	}
	if len(source) > 64<<10 {
		return fmt.Errorf("source exceeds replay profile")
	}
	dir, err := os.MkdirTemp("", "portugol-probe-")
	if err != nil {
		return err
	}
	defer func() { runErr = errors.Join(runErr, os.RemoveAll(dir)) }()
	if err := os.WriteFile(filepath.Join(dir, "source.alg"), source, 0o600); err != nil {
		return err
	}
	for _, file := range p.Files {
		if file.Path == "source.alg" || observe && (file.Path == "state.json" || file.Path == "host.json") {
			return fmt.Errorf("input file conflicts with probe source")
		}
		data, err := readArtifact(root, file.Content)
		if err != nil {
			return err
		}
		name, err := safePath(dir, file.Path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(name), 0o700); err != nil {
			return err
		}
		if err := os.WriteFile(name, data, 0o600); err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.TimeoutMS)*time.Millisecond)
	defer cancel()
	stdout, stderr := &boundedOutput{cancel: cancel}, &boundedOutput{cancel: cancel}
	steps := p.MaxSteps
	if steps == 0 {
		steps = 10000
	}
	command := "run"
	if observe {
		command, executable = "execute", observer
	}
	args := append(append([]string(nil), prefix...), command, "--max-steps", strconv.FormatUint(steps, 10))
	if observe {
		args = append(args, "--state", "state.json", "--host-trace", "host.json")
	}
	args = append(args, "source.alg")
	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Dir, cmd.Stdin, cmd.Stdout, cmd.Stderr = dir, bytes.NewReader(input), stdout, stderr
	cmd.WaitDelay = time.Second
	err = cmd.Run()
	if stdout.overflow || stderr.overflow {
		return fmt.Errorf("output limit exceeded")
	}
	if ctx.Err() != nil {
		return fmt.Errorf("replay deadline exceeded")
	}
	exitCode := 0
	if err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			return fmt.Errorf("start or collect replay process: %w", err)
		}
		exitCode = exit.ExitCode()
	}
	if exitCode != want.ExitCode {
		return fmt.Errorf("exit status %d, expected %d", exitCode, want.ExitCode)
	}
	wantOut, err := readArtifact(root, want.Stdout)
	if err != nil {
		return err
	}
	if !bytes.Equal(stdout.buffer.Bytes(), wantOut) {
		return fmt.Errorf("stdout mismatch: got %q, expected %q", stdout.buffer.Bytes(), wantOut)
	}
	if err := compareDiagnostics(stderr.buffer.String(), want.Diagnostics); err != nil {
		return err
	}
	for name, artifact := range map[string]*artifact{"state.json": want.State, "host.json": want.HostTrace} {
		if artifact == nil {
			continue
		}
		actual, err := readFile(dir, name)
		if err != nil {
			return fmt.Errorf("missing %s observation", name)
		}
		expected, err := readArtifact(root, *artifact)
		if err != nil {
			return err
		}
		if !bytes.Equal(actual, expected) {
			return fmt.Errorf("%s observation mismatch", name)
		}
	}
	for _, file := range want.Generated {
		actual, err := readFile(dir, file.Path)
		if err != nil {
			return fmt.Errorf("generated file %s unavailable", file.Path)
		}
		expected, err := readArtifact(root, file.Content)
		if err != nil {
			return err
		}
		if !bytes.Equal(actual, expected) {
			return fmt.Errorf("generated file mismatch: %s", file.Path)
		}
	}
	for _, name := range want.Absent {
		if err := checkAbsent(dir, name); err != nil {
			return err
		}
	}
	return nil
}

func compareDiagnostics(stderr string, want []diagnostic) error {
	if len(want) == 0 && stderr != "" {
		return fmt.Errorf("unexpected diagnostic output")
	}
	pattern := regexp.MustCompile(`^source\.alg:([0-9]+):([0-9]+): ([LPSER][0-9]{3}): .+$`)
	var got []diagnostic
	lines := strings.Split(strings.TrimSuffix(stderr, "\n"), "\n")
	for index, line := range lines {
		if line == "" {
			continue
		}
		// The current CLI appends this fixed summary after static diagnostics.
		if line == "check failed" && len(got) > 0 && index == len(lines)-1 {
			continue
		}
		match := pattern.FindStringSubmatch(line)
		if match == nil {
			return fmt.Errorf("unmapped diagnostic output")
		}
		lineNo, err := strconv.Atoi(match[1])
		if err != nil {
			return err
		}
		column, err := strconv.Atoi(match[2])
		if err != nil {
			return err
		}
		got = append(got, diagnostic{Code: match[3], Line: lineNo, Column: column})
	}
	if len(got) != len(want) {
		return fmt.Errorf("diagnostic count %d, expected %d", len(got), len(want))
	}
	for i, d := range got {
		w := want[i]
		if d.Code != w.Code || d.Line != w.Line || w.Column != 0 && d.Column != w.Column {
			return fmt.Errorf("diagnostic %d = %+v, expected %+v", i+1, d, w)
		}
	}
	return nil
}
