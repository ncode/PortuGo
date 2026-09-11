package portugol_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/stdlib"
)

func TestIndependentBuiltinInventory(t *testing.T) {
	root := "testdata/conformance/visualg-3.0.7"
	data, err := os.ReadFile(filepath.Join(root, "builtin-inventory.json"))
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Names   []string `json:"names"`
		ProbeID string   `json:"probeId"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	required := make(map[string]bool)
	for _, name := range inventory.Names {
		if _, exists := required[name]; exists || name == "" {
			t.Fatalf("invalid independent inventory name %q", name)
		}
		required[name] = false
	}
	if len(required) != 28 {
		t.Fatalf("documented inventory has %d names, want 28", len(required))
	}
	for _, descriptor := range stdlib.Catalog() {
		if _, ok := required[descriptor.Name()]; !ok {
			t.Errorf("catalog includes undocumented name %q", descriptor.Name())
		}
	}
	for name := range required {
		if _, ok := stdlib.Lookup(name); !ok {
			t.Errorf("catalog omits required reference function %q", name)
		}
	}
	dir := filepath.Join(root, "probes", inventory.ProbeID)
	src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
	if err != nil {
		t.Fatal(err)
	}
	_, tokens, ds := lexer.Scan("source.alg", src)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	program, ds := parser.Parse(tokens)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	info, ds := sema.Analyze(program)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	for _, tok := range tokens {
		if binding, ok := info.Binding(tok); ok && binding.Builtin {
			if _, documented := required[binding.Name]; !documented {
				t.Fatalf("probe calls undocumented builtin %q", binding.Name)
			}
			required[binding.Name] = true
		}
	}
	for name, exercised := range required {
		if !exercised {
			t.Errorf("independent inventory probe does not exercise %q", name)
		}
	}
	want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
	if err != nil {
		t.Fatal(err)
	}
	checkFormattingPreservesExecution(t, src, want)
}

func TestRejectedBuiltinInventoryCandidate(t *testing.T) {
	checkSemanticDiagnostic(t, "testdata/conformance/visualg-3.0.7/probes/catalog-pot-candidate/source.alg", diag.EUndeclared, 3)
}
