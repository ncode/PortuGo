// Package golden compares byte-exact fixtures with an explicit update mode.
package golden

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// Compare checks a root-relative fixture against got, or updates it when requested.
func Compare(root, name string, got []byte, update bool) error {
	path := filepath.Join(root, name)
	label := filepath.ToSlash(filepath.Clean(name))
	if update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fileError("create directory for", label, err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			return fileError("update", label, err)
		}
		return nil
	}
	want, err := os.ReadFile(path)
	if err != nil {
		return fileError("read", label, err)
	}
	if !bytes.Equal(want, got) {
		offset := 0
		for offset < min(len(want), len(got)) && want[offset] == got[offset] {
			offset++
		}
		return fmt.Errorf("%s: byte %d: want %s, got %s (lengths %d and %d)\n-want %q\n+got  %q",
			label, offset, byteAt(want, offset), byteAt(got, offset), len(want), len(got),
			want[max(0, offset-16):min(len(want), offset+16)],
			got[max(0, offset-16):min(len(got), offset+16)])
	}
	return nil
}

func byteAt(data []byte, offset int) string {
	if offset == len(data) {
		return "EOF"
	}
	return fmt.Sprintf("0x%02x", data[offset])
}

func fileError(operation, label string, err error) error {
	if pathErr, ok := err.(*os.PathError); ok {
		err = pathErr.Err // Keep the cause while reporting only the stable fixture path.
	}
	return fmt.Errorf("%s %s: %w", operation, label, err)
}
