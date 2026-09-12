package golden

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name            string
		want, got       []byte
		missing, update bool
		wantError       string
	}{
		{name: "equal", want: []byte("á\r\n"), got: []byte("á\r\n")},
		{name: "whitespace", want: []byte("x \n"), got: []byte("x\n"), wantError: "byte 1: want 0x20, got 0x0a"},
		{name: "encoding", want: []byte{0xe1}, got: []byte("á"), wantError: "byte 0: want 0xe1, got 0xc3"},
		{name: "truncated", want: []byte("abc"), got: []byte("ab"), wantError: "byte 2: want 0x63, got EOF"},
		{name: "extra", want: []byte("ab"), got: []byte("abc"), wantError: "byte 2: want EOF, got 0x63"},
		{name: "missing", missing: true, got: []byte("new"), wantError: "read testdata/result.out"},
		{name: "explicit create", missing: true, update: true, got: []byte{0x00, 0xe1}},
		{name: "explicit replace", want: []byte("old"), got: []byte("new"), update: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			name := filepath.Join("testdata", "result.out")
			path := filepath.Join(root, name)
			if !tt.missing {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, tt.want, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			err := Compare(root, name, tt.got, tt.update)
			if tt.wantError == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("error = %v, want %q", err, tt.wantError)
			} else if strings.Contains(err.Error(), root) {
				t.Errorf("error contains unstable root %q: %v", root, err)
			}
			data, readErr := os.ReadFile(path)
			if tt.missing && !tt.update {
				if !os.IsNotExist(readErr) {
					t.Fatalf("comparison created missing golden: %v", readErr)
				}
				return
			}
			if readErr != nil {
				t.Fatal(readErr)
			}
			want := tt.want
			if tt.update {
				want = tt.got
			}
			if string(data) != string(want) {
				t.Fatalf("golden bytes changed unexpectedly: %q", data)
			}
		})
	}
}
