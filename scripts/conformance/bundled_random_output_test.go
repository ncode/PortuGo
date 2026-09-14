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
	for _, p := range doc.Probes {
		if p.ID == "bundled-991ec2bd1566" {
			if len(p.Implementation.Expected.RandomOutput) == 0 || string(p.Implementation.Expected.RandomOutput) == "null" {
				t.Fatal("bundled randi example has no output contract")
			}
			return
		}
	}
	t.Fatal("bundled randi example not found")
}

func TestRecordedBundledRandiOutput(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	m, err := loadManifest(root, "testdata/conformance/visualg-3.0.7/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	var probe *probe
	for index := range m.Probes {
		if m.Probes[index].ID == "bundled-991ec2bd1566" {
			probe = &m.Probes[index]
			break
		}
	}
	if probe == nil || probe.Implementation.Expected.RandomOutput == nil {
		t.Fatal("bundled randi output contract is missing")
	}
	contract := probe.Implementation.Expected.RandomOutput
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
		for seed := uint64(0); seed < 16; seed++ {
			generators = append(generators, rand.New(rand.NewPCG(seed, seed+1)))
		}
		for index, rng := range generators {
			var out bytes.Buffer
			i := interp.New(interp.Options{Output: &out, Random: rng, MaxSteps: 10000})
			if ds := i.Run(program, info); len(ds) != 0 {
				t.Fatalf("pass=%d generator=%d: %v", pass, index, ds)
			}
			if err := contract.compare(out.Bytes()); err != nil {
				t.Fatalf("pass=%d generator=%d: %v", pass, index, err)
			}
		}
		var formatted bytes.Buffer
		if err := ast.Fprint(&formatted, program); err != nil {
			t.Fatal(err)
		}
		src = formatted.String()
	}
}
