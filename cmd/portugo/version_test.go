package main

import (
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		exit int
		out  string
	}{
		{"development version", []string{"--version"}, 0, "dev\n"},
		{"extra argument", []string{"--version", "unexpected"}, 2, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			out, stderr, exit := commandOutput(t, tt.args...)
			if exit != tt.exit || out != tt.out {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want exit=%d stdout=%q", exit, out, stderr, tt.exit, tt.out)
			}
			if tt.exit == 0 && stderr != "" || tt.exit == 2 && !strings.Contains(stderr, "usage:") {
				t.Errorf("unexpected stderr: %q", stderr)
			}
		})
	}
}
