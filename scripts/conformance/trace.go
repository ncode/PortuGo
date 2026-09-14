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
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
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
				if !isTestFunction(file, fn) {
					return fmt.Errorf("invalid test link %s: not a Go test function", link)
				}
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

// isTestFunction applies Go test discovery's name and declaration checks.
func isTestFunction(file *ast.File, fn *ast.FuncDecl) bool {
	suffix, ok := strings.CutPrefix(fn.Name.Name, "Test")
	if !ok || fn.Recv != nil {
		return false
	}
	if suffix != "" {
		r, _ := utf8.DecodeRuneInString(suffix)
		if unicode.IsLower(r) {
			return false
		}
	}
	if fn.Type.TypeParams.NumFields() != 0 || fn.Type.Results.NumFields() != 0 || fn.Type.Params.NumFields() != 1 {
		return false
	}
	ptr, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	if typ, ok := ptr.X.(*ast.Ident); ok {
		if typ.Name != "T" {
			return false
		}
		for _, spec := range file.Imports {
			if spec.Name != nil && spec.Name.Name == "." && importPath(spec) == "testing" {
				return true
			}
		}
		return false
	}
	typ, ok := ptr.X.(*ast.SelectorExpr)
	if !ok || typ.Sel.Name != "T" {
		return false
	}
	pkg, ok := typ.X.(*ast.Ident)
	if !ok {
		return false
	}
	for _, spec := range file.Imports {
		if importPath(spec) != "testing" || spec.Name != nil && (spec.Name.Name == "." || spec.Name.Name == "_") {
			continue
		}
		name := "testing"
		if spec.Name != nil {
			name = spec.Name.Name
		}
		if name == pkg.Name {
			return true
		}
	}
	return false
}

func importPath(spec *ast.ImportSpec) string {
	path, err := strconv.Unquote(spec.Path.Value)
	if err != nil {
		return ""
	}
	return path
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

func isNumberedSourceItem(line string) bool {
	line = strings.TrimSpace(line)
	for i := 0; i < len(line) && line[i] >= '0' && line[i] <= '9'; i++ {
		if i+1 < len(line) && (line[i+1] == '.' || line[i+1] == ')') && i+2 < len(line) && line[i+2] == ' ' {
			return true
		}
	}
	return false
}

// Verification markers may carry a free-form qualifier inside the brackets,
// such as "[VERIFICAR posição exata]". Treat every bracketed spelling as a
// marker, then let the source-line shape and explicit contextual list decide
// whether it creates an inventory obligation.
func hasVerificationMarker(line string) bool {
	return strings.Contains(line, "[VERIFICAR")
}

func isSourceObligationLine(line string) bool {
	line = strings.TrimSpace(line)
	if !hasVerificationMarker(line) {
		return false
	}
	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "+ ") || strings.HasPrefix(line, "|") || strings.HasPrefix(line, "`") || strings.HasPrefix(line, "[VERIFICAR]") || isNumberedSourceItem(line) {
		return true
	}
	return false
}

func checklistSourceItems(data []byte) []string {
	const header = "## 12. Checklist de conformidade"
	inChecklist := false
	var items []string
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, header) {
			inChecklist = true
			continue
		}
		if inChecklist && strings.HasPrefix(line, "## ") {
			break
		}
		if inChecklist && isNumberedSourceItem(line) {
			items = append(items, line)
		}
	}
	return items
}

func sourceObligationLink(name, line string) string {
	return name + "#" + strings.TrimSpace(strings.TrimLeft(line, "#"))
}

type pendingSourceObligation struct {
	line   int
	prefix string
}

type contextualSourceMarker struct {
	line    int
	prefix  string
	aliasID string
}

// These markers explain the marker/checklist notation or repeat a question
// already represented by a nearby source item. They are anchored explicitly
// so future edits cannot silently turn contextual prose or grammar comments
// into untracked obligations.
var contextualSourceMarkers = map[string][]contextualSourceMarker{
	"especificacao-visualg-3.md": {
		{line: 18, prefix: "Itens marcados"},
		{line: 171, prefix: "Pontos a validar empiricamente"},
		{line: 528, prefix: "grupo_param", aliasID: "assumption.vector-procedure-parameter"},
		{line: 549, prefix: "valor_caso", aliasID: "assumption.choice-range-ate"},
		{line: 562, prefix: "(* expressões:"},
	},
}

func isContextualSourceMarker(name string, lineNumber int, line string) bool {
	line = strings.TrimSpace(line)
	if !hasVerificationMarker(line) {
		return false
	}
	for _, marker := range contextualSourceMarkers[name] {
		if marker.line == lineNumber && strings.HasPrefix(line, marker.prefix) {
			return true
		}
	}
	return false
}

// These source markers still contain an unresolved implementation or
// diagnostic question. Keep them explicitly enumerated so the legacy source
// cannot hide them behind its checklist section, and do not synthesize probe
// links for them until the question is resolved.
var pendingSourceObligations = map[string][]pendingSourceObligation{
	"especificacao-visualg-3.md": {
		{line: 334, prefix: "- A gravação de compatibilidade mostra uma falha interna"},
		{line: 440, prefix: "- A pilha de ativação é visível no IDE"},
		{line: 441, prefix: "- Valores iniciais de variáveis:"},
	},
}

func isPendingSourceObligation(name string, lineNumber int, line string) bool {
	line = strings.TrimSpace(line)
	for _, pending := range pendingSourceObligations[name] {
		if pending.line == lineNumber && strings.HasPrefix(line, pending.prefix) {
			return true
		}
	}
	return false
}
func validateInventory(root string, m manifest, probes map[string]probe) error {
	var problems []error
	ids, links, requirementLinks, used := make(map[string]bool), make(map[string]bool), make(map[string]bool), make(map[string]bool)
	for _, item := range m.Inventory {
		if item.ID == "" || ids[item.ID] {
			problems = append(problems, fmt.Errorf("duplicate or empty inventory ID %q", item.ID))
		}
		ids[item.ID] = true
		links[item.Link] = true
		if item.Kind == "requirement" {
			requirementLinks[item.Link] = true
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
		itemProbes := make(map[string]bool)
		for _, id := range item.Probes {
			if itemProbes[id] {
				problems = append(problems, fmt.Errorf("duplicate probe link %s in %s", id, item.ID))
			}
			itemProbes[id] = true
			if _, ok := probes[id]; !ok {
				problems = append(problems, fmt.Errorf("stale probe link %s in %s", id, item.ID))
			}
			used[id] = true
		}
	}
	var untraced []string
	for id := range probes {
		if !used[id] {
			untraced = append(untraced, id)
		}
	}
	sort.Strings(untraced)
	for _, id := range untraced {
		problems = append(problems, fmt.Errorf("untraced probe %s", id))
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
		checklist := checklistSourceItems(data)
		pendingSeen := make(map[int]bool)
		contextualSeen := make(map[int]bool)
		for lineNumber, line := range strings.Split(string(data), "\n") {
			lineNumber++
			if strings.HasPrefix(line, "### Requirement: ") {
				link := name + "#" + strings.TrimPrefix(line, "### ")
				if !requirementLinks[link] {
					problems = append(problems, fmt.Errorf("untraced requirement %s", link))
				}
			}
			if isContextualSourceMarker(name, lineNumber, line) {
				contextualSeen[lineNumber] = true
				continue
			}
			if !isSourceObligationLine(line) {
				continue
			}
			if isPendingSourceObligation(name, lineNumber, line) {
				pendingSeen[lineNumber] = true
				continue
			}
			link := sourceObligationLink(name, line)
			if !links[link] {
				problems = append(problems, fmt.Errorf("untraced verification item %s", link))
			}
		}
		for _, pending := range pendingSourceObligations[name] {
			if !pendingSeen[pending.line] {
				problems = append(problems, fmt.Errorf("pending verification item missing %s#line %d", name, pending.line))
			}
		}
		for _, marker := range contextualSourceMarkers[name] {
			if !contextualSeen[marker.line] {
				problems = append(problems, fmt.Errorf("contextual verification marker missing %s#line %d", name, marker.line))
			} else if marker.aliasID != "" && !ids[marker.aliasID] {
				problems = append(problems, fmt.Errorf("contextual verification marker %s aliases missing inventory %s", name, marker.aliasID))
			}
		}
		for _, line := range checklist {
			link := sourceObligationLink(name, line)
			if !links[link] {
				problems = append(problems, fmt.Errorf("untraced checklist item %s", link))
			}
		}
	}
	return errors.Join(problems...)
}
