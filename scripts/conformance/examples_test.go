package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestManifestExampleCatalog(t *testing.T) {
	for _, tt := range []struct {
		name, want string
		mutate     func(map[string]any, *manifest)
	}{
		{name: "accepted source"},
		{name: "recorded rejection", mutate: func(e map[string]any, m *manifest) {
			e["disposition"] = "unusable"
			e["review"] = map[string]string{"reason": "Recorded rejection", "link": "review.md"}
			accepted := false
			m.Probes[0].Evidence.Accepted = &accepted
			m.Probes[0].Implementation.Expected.ExitCode = 1
			m.Probes[0].Implementation.Expected.Diagnostics = []diagnostic{{Code: "P001", Line: 3}}
		}},
		{name: "reviewed non-goal", mutate: func(e map[string]any, m *manifest) {
			e["disposition"] = "non-goal"
			e["review"] = map[string]string{"reason": "Explicit non-goal", "link": "review.md"}
			m.Probes[0].Evidence.State = "not-applicable"
			m.Probes[0].Evidence.Review = &review{Reason: "Explicit non-goal", Link: "review.md"}
		}},
		{name: "source hash", want: "example source", mutate: func(e map[string]any, _ *manifest) { e["sha256"] = strings.Repeat("a", 64) }},
		{name: "source length", want: "example source", mutate: func(e map[string]any, _ *manifest) { e["bytes"] = 1 }},
		{name: "unstable ID", want: "example ID", mutate: func(e map[string]any, _ *manifest) { e["name"] = "Exemplos/renamed.alg" }},
		{name: "untraced source", want: "untraced example", mutate: func(_ map[string]any, m *manifest) { m.Inventory[1].ID = "example.missing" }},
		{name: "unclassified recording", want: "example disposition", mutate: func(e map[string]any, _ *manifest) { e["disposition"] = "pending" }},
		{name: "invented disposition", want: "example disposition", mutate: func(e map[string]any, _ *manifest) { e["disposition"] = "complete" }},
		{name: "false rejection", want: "example acceptance", mutate: func(e map[string]any, _ *manifest) {
			e["disposition"] = "unusable"
			e["review"] = map[string]string{"reason": "Recorded rejection", "link": "review.md"}
		}},
		{name: "unreviewed exclusion", want: "review", mutate: func(e map[string]any, _ *manifest) { e["disposition"] = "unusable" }},
		{name: "wrong inventory kind", want: "untraced example", mutate: func(_ map[string]any, m *manifest) {
			m.Inventory[1].Kind = "assumption"
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root, m := testManifest(t)
			name := "Exemplos/output.alg"
			id := fmt.Sprintf("example.%x", sha256.Sum256([]byte(name)))[:20]
			src, err := readArtifact(root, m.Probes[0].Source)
			if err != nil {
				t.Fatal(err)
			}
			e := map[string]any{"id": id, "name": name, "sha256": m.Probes[0].Source.SHA256, "bytes": len(src), "disposition": "accepted"}
			m.Inventory = append(m.Inventory, inventoryItem{ID: id, Kind: "example", Link: "official-examples.json", Probes: []string{"output"}})
			writeArtifact(t, root, "review.md", "Recorded reference rejection.\n")
			if tt.mutate != nil {
				tt.mutate(e, &m)
			}
			data, err := json.Marshal([]any{e})
			if err != nil {
				t.Fatal(err)
			}
			writeArtifact(t, root, "official-examples.json", string(data))
			err = validate(root, m, "evidence", nil)
			if tt.want == "" && err != nil {
				t.Fatal(err)
			}
			if tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}
