package portugol_test

import (
	"errors"
	"os"
	"syscall"
	"testing"
)

func restrictFileRead(t *testing.T, name, mode string) func() {
	t.Helper()
	if mode == "denied" {
		t.Skip("explicit ACL denial has separate native reference evidence")
	}
	path, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		t.Fatal(err)
	}
	// Match the reference fixture's read/write handle. A read-only lock can
	// permit Go's backup-semantics open to read the file despite share mode zero.
	h, err := syscall.CreateFile(path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		t.Fatal(err)
	}
	restore := func() {
		if h != syscall.InvalidHandle {
			if err := syscall.CloseHandle(h); err != nil {
				t.Error(err)
			}
			h = syscall.InvalidHandle
		}
	}
	t.Cleanup(restore)
	f, err := os.Open(name)
	if err == nil {
		_ = f.Close()
		t.Fatal("exclusive lock allowed a second open")
	}
	if !errors.Is(err, syscall.Errno(32)) {
		t.Fatalf("read restriction returned %v, want sharing violation", err)
	}
	return restore
}
