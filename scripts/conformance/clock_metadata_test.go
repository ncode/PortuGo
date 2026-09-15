package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestManifestClockFixture(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct{ name, content, want string }{
		{"zero", `{"nowMS":[0]}`, ""},
		{"multiple reads", `{"nowMS":[0,16,16]}`, ""},
		{"clock regression", `{"nowMS":[5,-1]}`, ""},
		{"malformed", `{`, "decode clock fixture"},
		{"unknown field", `{"nowMS":[0],"extra":1}`, "decode clock fixture"},
		{"trailing JSON", `{"nowMS":[0]} {}`, "trailing clock fixture JSON"},
		{"missing reads", `{}`, "clock fixture has no reads"},
		{"empty reads", `{"nowMS":[]}`, "clock fixture has no reads"},
		{"null reads", `{"nowMS":null}`, "clock fixture has no reads"},
		{"null fixture", `null`, "clock fixture has no reads"},
		{"string read", `{"nowMS":["0"]}`, "decode clock fixture"},
		{"fractional read", `{"nowMS":[0.5]}`, "decode clock fixture"},
		{"overflowing read", `{"nowMS":[9223372036854775808]}`, "decode clock fixture"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			p := &m.Probes[0]
			p.Implementation.State = "verified"
			p.Implementation.Tests = []string{"output_test.go#TestOutput"}
			clock := writeArtifact(t, root, "clock.json", tt.content)
			p.Implementation.Expected.Clock = &clock
			var out, stderr bytes.Buffer
			status := executeProbe([]string{"--clock", filepath.Join(root, clock.Path), filepath.Join(root, p.Source.Path)}, strings.NewReader(""), &out, &stderr)
			if tt.want == "" {
				if status != 0 || out.String() != " 1\n" || stderr.Len() != 0 {
					t.Fatalf("adapter status = %d, stdout = %q, stderr = %q", status, &out, &stderr)
				}
			} else if status != 1 || out.Len() != 0 || !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("adapter status = %d, stdout = %q, stderr = %q; want %q", status, &out, &stderr, tt.want)
			}
			for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
				t.Run(mode, func(t *testing.T) {
					err := validate(root, m, mode, nil)
					if tt.want == "" {
						if err != nil {
							t.Fatal(err)
						}
					} else if err == nil || !strings.Contains(err.Error(), tt.want) {
						t.Fatalf("error = %v, want %q", err, tt.want)
					}
				})
			}
		})
	}
}
