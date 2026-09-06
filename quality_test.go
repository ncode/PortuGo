package portugol_test

import (
	"os"
	"strings"
	"testing"
)

// TestQualityPolicy keeps required checks from disappearing from the workflow.
// GitHub validates the YAML and executes the jobs; this checks their contract.
func TestQualityPolicy(t *testing.T) {
	t.Parallel()
	checks := []struct {
		path string
		want []string
	}{
		{".github/workflows/quality.yml", []string{
			"pull_request:", "push:", "go build ./...", "go run ./scripts/checkfmt",
			"go vet ./...", "staticcheck ./...", "golangci-lint run",
			"go test -count=1 ./...", "go test -race -count=1 ./...",
			"ubuntu-latest", "macos-latest", "windows-latest",
			"-fuzz=FuzzLexer -fuzztime=30s", "-fuzz=FuzzParser -fuzztime=30s",
			"openspec validate --all --strict --no-interactive", "bash scripts/check-oracle.sh",
		}},
		{".golangci.yml", []string{"version: \"2\"", "errcheck", "govet", "ineffassign", "staticcheck", "unused"}},
		{"docs/development.md", []string{"honnef.co/go/tools/cmd/staticcheck@v0.6.1", "golangci-lint/v2/cmd/golangci-lint@v2.4.0"}},
	}
	for _, check := range checks {
		t.Run(check.path, func(t *testing.T) {
			data, err := os.ReadFile(check.path)
			if err != nil {
				t.Fatal(err)
			}
			var active []string
			for _, line := range strings.Split(string(data), "\n") {
				if !strings.HasPrefix(strings.TrimSpace(line), "#") {
					active = append(active, line)
				}
			}
			text := strings.Join(active, "\n")
			for _, want := range check.want {
				if !strings.Contains(text, want) {
					t.Errorf("missing required quality gate %q", want)
				}
			}
		})
	}
}
