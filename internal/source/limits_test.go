package source

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
)

func TestSourceLimit(t *testing.T) {
	for _, extra := range []int{0, 1} {
		for _, b := range []byte{'x', 0x80} {
			data := bytes.Repeat([]byte{b}, MaxBytes+extra)
			text, err := Decode(data)
			if extra == 0 {
				want := MaxBytes
				if b == 0x80 {
					want *= 3
				}
				if err != nil || len(text) != want {
					t.Fatalf("boundary decode: %d bytes, %v", len(text), err)
				}
			} else {
				var d diag.Diagnostic
				if text != "" || !errors.As(err, &d) || d.Code != diag.EResource || int(d.Pos) != MaxBytes {
					t.Fatalf("oversized decode: %d bytes, %v", len(text), err)
				}
			}
		}
	}
	path := filepath.Join(t.TempDir(), "oversized.alg")
	data := bytes.Repeat([]byte{'\n'}, MaxBytes+1)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	_, err := ReadFile(path)
	if err == nil || !strings.Contains(err.Error(), "oversized.alg:4194305:1: E900:") {
		t.Fatalf("missing source position: %v", err)
	}
}
