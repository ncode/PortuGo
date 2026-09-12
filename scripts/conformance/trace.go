package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func taskStates(root, name string) (map[string]bool, map[int]bool, error) {
	data, err := readFile(root, name)
	if err != nil {
		return nil, nil, err
	}
	tasks, groups := make(map[string]bool), make(map[int]bool)
	pattern := regexp.MustCompile(`^- \[([ x])\] (([0-9]+)\.[0-9]+) `)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "- [") {
			continue
		}
		match := pattern.FindStringSubmatch(line)
		if match == nil {
			return nil, nil, fmt.Errorf("invalid task checkbox in %s", name)
		}
		if _, exists := tasks[match[2]]; exists {
			return nil, nil, fmt.Errorf("duplicate task %s", match[2])
		}
		group, err := strconv.Atoi(match[3])
		if err != nil {
			return nil, nil, err
		}
		if _, exists := groups[group]; !exists {
			groups[group] = true
		}
		done := match[1] == "x"
		tasks[match[2]] = done
		groups[group] = groups[group] && done
	}
	if len(tasks) == 0 {
		return nil, nil, fmt.Errorf("no task links in %s", name)
	}
	return tasks, groups, nil
}

func checkLink(root, link string, test bool) error {
	name, anchor, hasAnchor := strings.Cut(link, "#")
	data, err := readFile(root, name)
	if err != nil {
		return err
	}
	if !hasAnchor || anchor == "" {
		if test {
			return fmt.Errorf("test link lacks a test function: %s", link)
		}
		return nil
	}
	if test {
		if !strings.HasSuffix(name, "_test.go") || !strings.HasPrefix(anchor, "Test") {
			return fmt.Errorf("invalid test link %s", link)
		}
		file, err := parser.ParseFile(token.NewFileSet(), name, data, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == anchor {
				return nil
			}
		}
		return fmt.Errorf("stale test link %s", link)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(strings.TrimLeft(line, "#")) == anchor {
			return nil
		}
	}
	return fmt.Errorf("stale trace link %s", link)
}

func checkTests(root string, tests []string) error {
	if len(tests) == 0 {
		return fmt.Errorf("missing test link")
	}
	for _, test := range tests {
		if err := checkLink(root, test, true); err != nil {
			return err
		}
	}
	return nil
}

func checkReview(root string, r *review) error {
	if r == nil || strings.TrimSpace(r.Reason) == "" || strings.TrimSpace(r.Link) == "" {
		return fmt.Errorf("missing reviewed rationale and review link")
	}
	if err := checkLink(root, r.Link, false); err != nil {
		return fmt.Errorf("review: %w", err)
	}
	return nil
}

func validateInventory(root string, m manifest, probes map[string]probe) error {
	var problems []error
	ids, links, used := make(map[string]bool), make(map[string]bool), make(map[string]bool)
	for _, item := range m.Inventory {
		if item.ID == "" || ids[item.ID] {
			problems = append(problems, fmt.Errorf("duplicate or empty inventory ID %q", item.ID))
		}
		ids[item.ID] = true
		if item.Kind == "requirement" {
			links[item.Link] = true
		}
		switch item.Kind {
		case "requirement", "checklist", "defect", "assumption", "feature", "example":
		default:
			problems = append(problems, fmt.Errorf("invalid inventory kind %q", item.Kind))
		}
		if err := checkLink(root, item.Link, false); err != nil {
			problems = append(problems, fmt.Errorf("%s trace link: %w", item.ID, err))
		}
		if len(item.Probes) == 0 {
			problems = append(problems, fmt.Errorf("missing probe link for %s", item.ID))
		}
		for _, id := range item.Probes {
			if _, ok := probes[id]; !ok {
				problems = append(problems, fmt.Errorf("stale probe link %s in %s", id, item.ID))
			}
			used[id] = true
		}
	}
	for id := range probes {
		if !used[id] {
			problems = append(problems, fmt.Errorf("untraced probe %s", id))
		}
	}
	if len(m.InventorySources) == 0 {
		problems = append(problems, fmt.Errorf("missing inventory sources"))
	}
	// Discover the change's complete specification tree independently of the
	// declared sources, so removing a source cannot remove its obligations.
	sources := append([]string(nil), m.InventorySources...)
	seen := make(map[string]bool)
	for _, name := range sources {
		seen[name] = true
	}
	err := fs.WalkDir(os.DirFS(root), path.Join(path.Dir(m.TasksPath), "specs"), func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(name, ".md") && !seen[name] {
			sources = append(sources, name)
			seen[name] = true
		}
		return nil
	})
	if err != nil {
		problems = append(problems, fmt.Errorf("discover specifications: %w", err))
	}
	for _, name := range sources {
		data, err := readFile(root, name)
		if err != nil {
			problems = append(problems, err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "### Requirement: ") {
				link := name + "#" + strings.TrimPrefix(line, "### ")
				if !links[link] {
					problems = append(problems, fmt.Errorf("untraced requirement %s", link))
				}
			}
		}
	}
	return errors.Join(problems...)
}
