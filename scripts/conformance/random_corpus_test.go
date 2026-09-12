package main

import (
	"bytes"
	"math/rand/v2"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/source"
)

// endpointRandom forces every draw to an endpoint, independently of a seed.
type endpointRandom struct{ upper bool }

func (r endpointRandom) Float64() float64 { return 0 }
func (r endpointRandom) Uint64N(n uint64) uint64 {
	if r.upper {
		return n - 1
	}
	return 0
}

func TestRecordedRandomInputDomains(t *testing.T) {
	t.Parallel()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	m, err := loadManifest(root, "testdata/conformance/visualg-3.0.7/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	qualified := 0
	for _, p := range m.Probes {
		contract := p.Implementation.Expected.RandomInput
		if contract == nil {
			continue
		}
		qualified++
		t.Run(p.ID, func(t *testing.T) {
			t.Parallel()
			recorded, err := readArtifact(root, p.Evidence.Normalized)
			if err != nil {
				t.Fatal(err)
			}
			if err := contract.compare(recorded); err != nil {
				t.Fatal(err)
			}
			data, err := readArtifact(root, p.Source)
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
						t.Fatal(ds)
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
		})
	}
	if qualified != 17 {
		t.Fatalf("qualified %d random-input recordings, expected 17", qualified)
	}
}
