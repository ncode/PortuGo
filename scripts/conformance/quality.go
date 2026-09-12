package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// qualityReport records trusted local results, not independent attestations.
type qualityReport struct {
	Version int             `json:"version"`
	Commit  string          `json:"commit"`
	Results []qualityResult `json:"results"`
}

type qualityResult struct {
	Check    string   `json:"check"`
	Status   string   `json:"status"`
	Evidence artifact `json:"evidence"`
}

func validateQuality(root, name, candidate string) error {
	if name == "" {
		return fmt.Errorf("missing quality results: supply --quality")
	}
	if candidate != "" {
		return fmt.Errorf("acceptance requires the default current-checkout candidate")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	git := func(args ...string) ([]byte, error) {
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = root
		return cmd.Output()
	}
	commit, err := git("rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return fmt.Errorf("resolve quality candidate: %w", err)
	}
	status, err := git("status", "--porcelain", "--untracked-files=all", "--ignored")
	if err != nil {
		return fmt.Errorf("inspect quality candidate: %w", err)
	}
	if len(status) != 0 {
		return fmt.Errorf("quality candidate requires a clean checkout including untracked files")
	}
	data, err := readFile(filepath.Dir(name), filepath.Base(name))
	if err != nil {
		return fmt.Errorf("read quality results: %w", err)
	}
	var report qualityReport
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&report); err != nil {
		return fmt.Errorf("decode quality results: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("trailing quality results JSON")
	}
	return checkQualityReport(filepath.Dir(name), report, strings.TrimSpace(string(commit)))
}

func checkQualityReport(root string, report qualityReport, commit string) error {
	var problems []error
	if report.Version != 1 {
		problems = append(problems, fmt.Errorf("unsupported quality results version"))
	}
	if report.Commit != commit || len(commit) != 40 && len(commit) != 64 {
		problems = append(problems, fmt.Errorf("quality results belong to a different candidate commit"))
	}
	required := map[string]bool{
		"build": false, "gofmt": false, "vet": false, "staticcheck": false,
		"golangci-lint": false, "tests": false, "race": false,
		"windows": false, "macos": false, "linux": false,
		"fuzz-lexer": false, "fuzz-parser": false, "openspec": false,
	}
	for _, result := range report.Results {
		seen, known := required[result.Check]
		if !known || seen {
			problems = append(problems, fmt.Errorf("unknown or duplicate quality check %q", result.Check))
			continue
		}
		required[result.Check] = true
		if result.Status != "pass" {
			problems = append(problems, fmt.Errorf("quality check %s did not pass", result.Check))
		}
		data, err := readArtifact(root, result.Evidence)
		if err != nil {
			problems = append(problems, fmt.Errorf("quality check %s evidence: %w", result.Check, err))
		} else if len(bytes.TrimSpace(data)) == 0 {
			problems = append(problems, fmt.Errorf("quality check %s has empty evidence", result.Check))
		}
	}
	var missing []string
	for check, present := range required {
		if !present {
			missing = append(missing, check)
		}
	}
	slices.Sort(missing)
	if len(missing) != 0 {
		problems = append(problems, fmt.Errorf("missing quality results: %s", strings.Join(missing, ", ")))
	}
	return errors.Join(problems...)
}

func incompleteTasks(root, name string) ([]string, error) {
	tasks, _, err := taskStates(root, name)
	if err != nil {
		return nil, err
	}
	var incomplete []string
	for id, complete := range tasks {
		if !complete {
			incomplete = append(incomplete, id)
		}
	}
	slices.Sort(incomplete)
	return incomplete, nil
}
