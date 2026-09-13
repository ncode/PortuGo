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

func TestNonGoalCatalogSourceMetadata(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"evidence", "incremental", "implementation-acceptance"} {
		for _, field := range []string{"hash", "bytes"} {
			t.Run(mode+"/"+field, func(t *testing.T) {
				t.Parallel()
				root, m := testManifest(t)
				writeArtifact(t, root, "review.md", "Recorded non-goal.\n")
				p := &m.Probes[0]
				p.Evidence.State = "not-applicable"
				p.Evidence.Review = &review{Reason: "Recorded non-goal", Link: "review.md"}
				p.Implementation.State = "not-applicable"
				p.Implementation.Review = &review{Reason: "Recorded non-goal", Link: "review.md"}
				name := "examples/output.alg"
				id := "example." + hashBytes([]byte(name))[:12]
				src, err := readArtifact(root, p.Source)
				if err != nil {
					t.Fatal(err)
				}
				e := bundledExample{ID: id, Name: name, SHA256: p.Source.SHA256, Bytes: len(src), Disposition: "non-goal", Review: &review{Reason: "Recorded non-goal", Link: "review.md"}}
				switch field {
				case "hash":
					e.SHA256 = strings.Repeat("0", 64)
				case "bytes":
					e.Bytes++
				}
				m.Inventory = append(m.Inventory, inventoryItem{ID: id, Kind: "example", Link: "official-examples.json", Probes: []string{"output"}})
				writeJSON(t, root, "official-examples.json", []bundledExample{e})
				if err := validate(root, m, mode, nil); err == nil || !strings.Contains(err.Error(), "example source differs from catalog") {
					t.Fatalf("error = %v, want catalog source metadata failure", err)
				}
			})
		}
	}
}
