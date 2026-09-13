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

func replayOwnedPath(name string, observe bool) bool {
	first, _, _ := strings.Cut(name, "/")
	first = strings.ToUpper(first)
	return first == "SOURCE.ALG" || observe && (first == "STATE.JSON" || first == "HOST.JSON" || first == "CLOCK.JSON")
}

func checkExpectedExitCode(code int) error {
	if code != 0 && code != 1 {
		return fmt.Errorf("invalid expected exit status %d", code)
	}
	return nil
}

func stdoutMismatch(got, want []byte) error {
	offset := 0
	for offset < len(got) && offset < len(want) && got[offset] == want[offset] {
		offset++
	}
	return fmt.Errorf("stdout mismatch: got %d bytes, expected %d bytes, first difference at byte %d", len(got), len(want), offset)
}

func checkReplayInputPath(name string, observe bool) error {
	if replayOwnedPath(name, observe) {
		return fmt.Errorf("input file conflicts with replay-owned path: %s", name)
	}
	return nil
}

func replayProbe(root string, p probe, executable string, prefix []string, observer string) (runErr error) {
	want := p.Implementation.Expected
	observe := want.State != nil || want.HostTrace != nil || want.Clock != nil
	if observe && observer == "" {
		return fmt.Errorf("missing state/host observation adapter")
	}
	if p.TimeoutMS < 1 || p.TimeoutMS > 30000 {
		return fmt.Errorf("invalid replay budget")
	}
	if err := checkExpectedExitCode(want.ExitCode); err != nil {
		return err
	}
	for _, a := range []*artifact{want.State, want.HostTrace} {
		if a != nil {
			if _, err := readReplayOutputArtifact(root, *a); err != nil {
				return err
			}
		}
	}
	if _, err := readReplayOutputArtifact(root, want.Stdout); err != nil {
		return err
	}
	source, err := readReplaySourceArtifact(root, p.Source)
	if err != nil {
		return err
	}
	input, err := readArtifact(root, p.Input)
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", "portugol-probe-")
	if err != nil {
		return err
	}
	defer func() { runErr = errors.Join(runErr, os.RemoveAll(dir)) }()
	if err := os.WriteFile(filepath.Join(dir, "source.alg"), source, 0o600); err != nil {
		return err
	}
	if want.Clock != nil {
		data, err := readArtifact(root, *want.Clock)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, "clock.json"), data, 0o600); err != nil {
			return err
		}
	}
	for _, file := range p.Files {
		if err := checkReplayInputPath(file.Path, observe); err != nil {
			return err
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
	for _, file := range want.Generated {
		if replayOwnedPath(file.Path, observe) {
			return fmt.Errorf("generated expectation conflicts with replay-owned path: %s", file.Path)
		}
	}
	for _, name := range want.Absent {
		if replayOwnedPath(name, observe) {
			return fmt.Errorf("absent expectation conflicts with replay-owned path: %s", name)
		}
	}
	restoreAccess, err := prepareFixtureAccess(dir, p.FixtureAccess)
	if err != nil {
		return err
	}
	restoredAccess := false
	restore := func() error {
		if restoredAccess {
			return nil
		}
		restoredAccess = true
		return restoreAccess()
	}
	defer func() {
		if err := restore(); err != nil {
			runErr = errors.Join(runErr, fmt.Errorf("restore fixture access: %w", err))
		}
	}()
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
		if want.State != nil {
			args = append(args, "--state", "state.json")
		}
		if want.HostTrace != nil {
			args = append(args, "--host-trace", "host.json")
		}
		if want.Clock != nil {
			args = append(args, "--clock", "clock.json")
		}
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
	if err := restore(); err != nil {
		return fmt.Errorf("restore fixture access: %w", err)
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
	wantOut, err := readReplayOutputArtifact(root, want.Stdout)
	if err != nil {
		return err
	}
	if want.RandomInput != nil {
		if err := want.RandomInput.compare(wantOut); err != nil {
			return fmt.Errorf("recorded random input: %w", err)
		}
		if err := want.RandomInput.compare(stdout.buffer.Bytes()); err != nil {
			return fmt.Errorf("replayed random input: %w", err)
		}
	} else if !bytes.Equal(stdout.buffer.Bytes(), wantOut) {
		return stdoutMismatch(stdout.buffer.Bytes(), wantOut)
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
		expected, err := readReplayOutputArtifact(root, *artifact)
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
	for _, file := range p.Files {
		declared := false
		for _, generated := range want.Generated {
			if generated.Path == file.Path {
				declared = true
				break
			}
		}
		if !declared {
			for _, absent := range want.Absent {
				if absent == file.Path {
					declared = true
					break
				}
			}
		}
		if declared {
			continue
		}
		actual, err := readFile(dir, file.Path)
		if err != nil {
			return fmt.Errorf("input file mismatch: %s: %w", file.Path, err)
		}
		expected, err := readArtifact(root, file.Content)
		if err != nil {
			return err
		}
		if !bytes.Equal(actual, expected) {
			return fmt.Errorf("input file mismatch: %s", file.Path)
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
		if lineNo < 1 || column < 1 {
			return fmt.Errorf("unpositioned diagnostic output")
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
