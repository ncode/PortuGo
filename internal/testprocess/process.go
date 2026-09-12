// Package testprocess isolates adversarial tests behind a failing watchdog.
package testprocess

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

// MaxSourceBytes is the generated-source fuzz profile, not a language limit.
const MaxSourceBytes = 64 << 10

// Run executes the selected test body in a child process with a five-second limit.
func Run(t *testing.T, body func()) {
	t.Helper()
	RunWithTimeout(t, 5*time.Second, body)
}

// RunWithTimeout executes the selected test body in a child process with the given limit.
func RunWithTimeout(t *testing.T, timeout time.Duration, body func()) {
	t.Helper()
	if os.Getenv("PORTUGOL_TEST_CHILD") == t.Name() {
		body()
		return
	}
	if err := run(t.Name(), timeout, ""); err != nil {
		t.Fatal(err)
	}
}

func run(name string, timeout time.Duration, mode string) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	parts := strings.Split(name, "/")
	for i, part := range parts {
		parts[i] = "^" + regexp.QuoteMeta(part) + "$"
	}
	cmd := exec.CommandContext(ctx, executable, "-test.run="+strings.Join(parts, "/"), "-test.count=1", "-test.timeout=0")
	cmd.Env = append(os.Environ(), "PORTUGOL_TEST_CHILD="+name, "PORTUGOL_TEST_MODE="+mode)
	output, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return fmt.Errorf("%s: watchdog expired after %s: %w\n%s", name, timeout, ctx.Err(), output)
	}
	if err != nil {
		return fmt.Errorf("%s: child test: %w\n%s", name, err, output)
	}
	return nil
}
