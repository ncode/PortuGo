//go:build !windows

package portugol_test

import (
	"errors"
	"os"
	"testing"
)

func restrictFileRead(t *testing.T, name, mode string) func() {
	t.Helper()
	if mode == "locked" {
		t.Skip("exclusive file-sharing restriction requires Windows")
	}
	if err := os.Chmod(name, 0); err != nil {
		t.Fatal(err)
	}
	restore := func() {
		if err := os.Chmod(name, 0o600); err != nil {
			t.Error(err)
		}
	}
	t.Cleanup(restore)
	f, err := os.Open(name)
	if err == nil {
		_ = f.Close()
		t.Skip("current account can read despite mode denial")
	}
	if !errors.Is(err, os.ErrPermission) {
		t.Fatalf("read restriction returned %v, want permission denial", err)
	}
	return restore
}
