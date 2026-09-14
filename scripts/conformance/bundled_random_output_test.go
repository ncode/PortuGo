package main

import (
	"bytes"
	"encoding/json"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

type nonZeroEndpointRandom struct{ upper bool }

func (r nonZeroEndpointRandom) Float64() float64 { return 0 }
func (r nonZeroEndpointRandom) Uint64N(n uint64) uint64 {
	if n <= 1 {
		return 0
	}
	if r.upper {
		return n - 1
	}
	return 1
}

func TestBundledRandiOutputContractMetadata(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "testdata/conformance/visualg-3.0.7/manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Probes []struct {
			ID             string `json:"id"`
			Implementation struct {
				Expected struct {
					RandomOutput json.RawMessage `json:"randomOutput"`
				} `json:"expected"`
			} `json:"implementation"`
		} `json:"probes"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	required := map[string]string{
		"randi-lines":                    "bundled randi example",
		"random-int-text-lines":          "bundled mixed random example",
		"random-int-sort-lines":          "bundled random sort example",
		"random-real-sort-lines":         "bundled real sort examples",
		"random-record-sort-lines":       "bundled record sort examples",
		"random-record-search-lines":     "bundled record search example",
		"random-int-search-table-lines":  "bundled random search table example",
		"random-int-search-lines":        "bundled random search example",
		"random-int-search-sorted-lines": "bundled sorted random search example",
		"random-int-countsort-lines":     "bundled counting sort example",
		"random-randi-repeat-lines":      "bundled repeated random example",
	}
	found := make(map[string]bool)
	for _, p := range doc.Probes {
		var contract struct {
			Kind string `json:"kind"`
		}
		if len(p.Implementation.Expected.RandomOutput) == 0 || string(p.Implementation.Expected.RandomOutput) == "null" || json.Unmarshal(p.Implementation.Expected.RandomOutput, &contract) != nil {
			continue
		}
		if _, ok := required[contract.Kind]; !ok {
			continue
		}
		found[contract.Kind] = true
	}
	for kind, label := range required {
		if !found[kind] {
			t.Fatalf("%s not found", label)
		}
	}
}

func TestRecordedBundledRandiOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "randi-lines")
}

func TestRecordedBundledMixedRandomOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-int-text-lines")
}

func TestRecordedBundledRandomSortOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-int-sort-lines")
}

func TestRecordedBundledRandomRealSortOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-real-sort-lines")
}

func TestRecordedBundledRandomRecordSortOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-record-sort-lines")
}

func TestRecordedBundledRandomRecordSearchOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-record-search-lines")
}

func TestRecordedBundledRandomSearchTableOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-int-search-table-lines")
}

func TestRecordedBundledRandomSearchOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-int-search-lines")
}

func TestRecordedBundledRandomSearchSortedOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-int-search-sorted-lines")
}

func TestRecordedBundledRandomCountSortOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-int-countsort-lines")
}

func TestRecordedBundledRandomRandiRepeatOutput(t *testing.T) {
	testRecordedBundledRandomOutput(t, "random-randi-repeat-lines")
}

func testRecordedBundledRandomOutput(t *testing.T, kind string) {
	t.Helper()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	m, err := loadManifest(root, "testdata/conformance/visualg-3.0.7/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	var probes []*probe
	for index := range m.Probes {
		contract := m.Probes[index].Implementation.Expected.RandomOutput
		if contract != nil && contract.Kind == kind {
			probes = append(probes, &m.Probes[index])
		}
	}
	if len(probes) == 0 {
		t.Fatal("bundled random output contract is missing")
	}
	for probeIndex, probe := range probes {
		contract := probe.Implementation.Expected.RandomOutput
		var input []byte
		if probe.Input.Path != "" {
			input, err = readArtifact(root, probe.Input)
			if err != nil {
				t.Fatal(err)
			}
		}
		recorded, err := readArtifact(root, probe.Evidence.Normalized)
		if err != nil {
			t.Fatal(err)
		}
		if err := contract.compare(recorded); err != nil {
			t.Fatal(err)
		}
		data, err := readArtifact(root, probe.Source)
		if err != nil {
			t.Fatal(err)
		}
		src, err := source.Decode(data)
		if err != nil {
			t.Fatal(err)
		}
		for pass := range 2 {
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
			generators := []interp.RandomSource{endpointRandom{}, endpointRandom{upper: true}}
			if contract.Kind == "random-randi-repeat-lines" {
				generators = []interp.RandomSource{nonZeroEndpointRandom{}, nonZeroEndpointRandom{upper: true}}
			}
			for seed := uint64(0); seed < 16; seed++ {
				generators = append(generators, rand.New(rand.NewPCG(seed, seed+1)))
			}
			for index, rng := range generators {
				var out bytes.Buffer
				i := interp.New(interp.Options{Input: bytes.NewReader(input), Output: &out, Random: rng, MaxSteps: 10000})
				if ds := i.Run(program, info); len(ds) != 0 {
					t.Fatalf("case=%d pass=%d generator=%d: %v", probeIndex, pass, index, ds)
				}
				if err := contract.compare(out.Bytes()); err != nil {
					t.Fatalf("case=%d pass=%d generator=%d: %v", probeIndex, pass, index, err)
				}
			}
			var formatted bytes.Buffer
			if err := ast.Fprint(&formatted, program); err != nil {
				t.Fatal(err)
			}
			src = formatted.String()
		}
	}
}
