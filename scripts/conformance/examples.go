package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
)

type bundledExample struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	SHA256      string  `json:"sha256"`
	Bytes       int     `json:"bytes"`
	Disposition string  `json:"disposition"`
	Review      *review `json:"review,omitempty"`
}

func validateExampleCatalogs(root string, inventory []inventoryItem, probes map[string]probe) error {
	catalogs := make(map[string]bool)
	items := make(map[string]inventoryItem)
	for _, item := range inventory {
		if item.Kind == "example" || path.Base(item.Link) == "official-examples.json" {
			catalogs[item.Link] = true
		}
		if item.Kind == "example" {
			items[item.ID] = item
		}
	}
	var problems []error
	for name := range catalogs {
		data, err := readFile(root, name)
		if err != nil {
			problems = append(problems, err)
			continue
		}
		var examples []bundledExample
		if err := json.Unmarshal(data, &examples); err != nil {
			problems = append(problems, fmt.Errorf("example catalog %s: %w", name, err))
			continue
		}
		seen := make(map[string]bool)
		for _, e := range examples {
			if e.ID != "example."+hashBytes([]byte(e.Name))[:12] || seen[e.ID] {
				problems = append(problems, fmt.Errorf("invalid or duplicate example ID %s", e.ID))
			}
			seen[e.ID] = true
			if _, err := safePath(root, e.Name); err != nil || !strings.EqualFold(path.Ext(e.Name), ".alg") || !validHash(e.SHA256) || e.Bytes < 1 {
				problems = append(problems, fmt.Errorf("invalid example source metadata %s", e.ID))
			}
			item, ok := items[e.ID]
			if !ok || item.Link != name {
				problems = append(problems, fmt.Errorf("untraced example %s", e.ID))
				continue
			}
			switch e.Disposition {
			case "pending":
				if len(item.Probes) != 0 {
					problems = append(problems, fmt.Errorf("pending example disposition has recorded links: %s", e.ID))
				}
			case "accepted":
			case "unusable", "non-goal":
				if err := checkReview(root, e.Review); err != nil {
					problems = append(problems, fmt.Errorf("example %s: %w", e.ID, err))
				}
			default:
				problems = append(problems, fmt.Errorf("invalid example disposition %s", e.ID))
			}
			for _, id := range item.Probes {
				p := probes[id]
				if e.Disposition == "non-goal" {
					if p.Evidence.State != "not-applicable" {
						problems = append(problems, fmt.Errorf("example non-goal requires reviewed non-applicability: %s", e.ID))
					}
					continue
				}
				src, err := readArtifact(root, p.Source)
				if err != nil || len(src) != e.Bytes || p.Source.SHA256 != e.SHA256 {
					problems = append(problems, fmt.Errorf("example source differs from catalog: %s", e.ID))
				}
				if p.Evidence.State != "recorded" || p.Evidence.Accepted == nil || *p.Evidence.Accepted != (e.Disposition == "accepted") {
					problems = append(problems, fmt.Errorf("example acceptance differs from recording: %s", e.ID))
				}
			}
		}
		for id, item := range items {
			if item.Link == name && !seen[id] {
				problems = append(problems, fmt.Errorf("untraced example inventory entry %s", id))
			}
		}
	}
	return errors.Join(problems...)
}
