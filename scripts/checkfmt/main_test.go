package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheck(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, source string
		wantError    bool
	}{
		{"formatted", "package example\n\nvar x = 1\n", false},
		{"unformatted", "package example;var x=1\n", true},
		{"invalid", "package {", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			path := filepath.Join(root, "example.go")
			if err := os.WriteFile(path, []byte(tt.source), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := check(root); (err != nil) != tt.wantError {
				t.Errorf("check = %v, want error %v", err, tt.wantError)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.source {
				t.Fatal("check rewrote the file")
			}
		})
	}
}
