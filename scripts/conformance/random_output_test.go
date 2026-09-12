package main

import (
	"os"
	"strings"
	"testing"
)

func TestRandomInputReplay(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	p := m.Probes[0]
	p.TimeoutMS = 5000
	p.Source = writeArtifact(t, root, p.Source.Path, "random-input")
	p.Implementation.Expected.Stdout = writeArtifact(t, root, "random.txt", "2.0000000000\n 2\n")
	p.Implementation.Expected.RandomInput = &randomInputExpectation{Kind: "real", Minimum: 2000, Maximum: 2999, Decimals: 3, Samples: 1}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err := replayProbe(root, p, executable, []string{"-test.run=^TestReplayChild$", "--"}, ""); err != nil {
		t.Fatal(err)
	}
	p.Implementation.Expected.RandomInput.Maximum = 2500
	if err := replayProbe(root, p, executable, []string{"-test.run=^TestReplayChild$", "--"}, ""); err == nil || !strings.Contains(err.Error(), "replayed random input") {
		t.Fatalf("out-of-domain replay was not rejected: %v", err)
	}
}

func TestRandomInputOutputContract(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, kind, output string
		min, max           int64
		decimals, samples  int
		valid              bool
	}{
		{"lower endpoint", "real", "2.0000000000\n 2\n", 2000, 2999, 3, 1, true},
		{"upper endpoint", "real", "2.9990000000\n 2.999\n", 2000, 2999, 3, 1, true},
		{"upper excluded", "real", "3.0000000000\n 3\n", 2000, 2999, 3, 1, false},
		{"below range", "real", "1.9990000000\n 1.999\n", 2000, 2999, 3, 1, false},
		{"off grid", "real", "2.0001000000\n 2.0001\n", 2000, 2999, 3, 1, false},
		{"fractional origin", "real", "3.7490000000\n 3.749\n", 2750, 3749, 3, 1, true},
		{"default integer", "integer", "100\n 100\n", 0, 100, 0, 1, true},
		{"negative integer", "integer", "-3\n -3\n", -3, 100, 0, 1, true},
		{"wrong stored value", "integer", "7\n 8\n", 0, 100, 0, 1, false},
		{"missing numeric space", "integer", "7\n7\n", 0, 100, 0, 1, false},
		{"wrong real echo", "real", "2.847\n 2.847\n", 2000, 2999, 3, 1, false},
		{"wrong real output", "real", "2.8470000000\n 2.848\n", 2000, 2999, 3, 1, false},
		{"leading zero", "integer", "07\n 7\n", 0, 100, 0, 1, false},
		{"rational spelling", "real", "5/2\n 2.5\n", 2000, 2999, 3, 1, false},
		{"extreme exponent", "real", "1e999999999\n 1\n", 2000, 2999, 3, 1, false},
		{"not finite", "real", "NaN\n NaN\n", 2000, 2999, 3, 1, false},
		{"negative zero real", "real", "-0.0000000000\n -0\n", 0, 100, 0, 1, false},
		{"negative zero integer", "integer", "-0\n -0\n", 0, 100, 0, 1, false},
		{"text", "text", "AZABC\nAZABC\n", 0, 0, 0, 1, true},
		{"text mismatch", "text", "AZABC\nAZABD\n", 0, 0, 0, 1, false},
		{"lowercase text", "text", "azabc\nazabc\n", 0, 0, 0, 1, false},
		{"wrong text length", "text", "ABCD\nABCD\n", 0, 0, 0, 1, false},
		{"missing pair", "integer", "7\n 7\n", 0, 100, 0, 2, false},
		{"extra output", "integer", "7\n 7\n\n", 0, 100, 0, 1, false},
		{"crlf", "integer", "7\r\n 7\r\n", 0, 100, 0, 1, false},
		{"no final newline", "integer", "7\n 7", 0, 100, 0, 1, false},
		{"invalid kind", "unknown", "7\n 7\n", 0, 100, 0, 1, false},
		{"invalid integer grid", "integer", "7\n 7\n", 0, 100, 1, 1, false},
		{"invalid text bounds", "text", "ABCDE\nABCDE\n", 0, 100, 0, 1, false},
		{"reversed domain", "integer", "7\n 7\n", 100, 0, 0, 1, false},
		{"negative precision", "real", "2.0000000000\n 2\n", 2, 3, -1, 1, false},
		{"unsupported precision", "real", "2.0000000000\n 2\n", 2, 3, 6, 1, false},
		{"empty sample count", "integer", "", 0, 100, 0, 0, false},
		{"excessive sample count", "integer", "", 0, 100, 0, 4097, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			contract := randomInputExpectation{Kind: tt.kind, Minimum: tt.min, Maximum: tt.max, Decimals: tt.decimals, Samples: tt.samples}
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
}

func TestRandomInputContractValidation(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"accepted", "missing review", "wrong owner", "rejected", "invalid recorded sample", "replaced evidence"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, m := testManifest(t)
			p := m.Probes[0]
			p.OwnerGroup, p.Tasks = 13, []string{"13.6"}
			out := writeArtifact(t, root, "random.txt", "7\n 7\n")
			p.Evidence.Raw, p.Evidence.Normalized, p.Evidence.Normalizer = out, out, "bytes-v1"
			p.Implementation.Expected.Stdout = out
			p.Implementation.Expected.RandomInput = &randomInputExpectation{Kind: "integer", Minimum: 0, Maximum: 100, Samples: 1, Review: review{Reason: "Recorded generated input uses a domain contract.", Link: "tasks.md"}}
			switch name {
			case "missing review":
				p.Implementation.Expected.RandomInput.Review = review{}
			case "wrong owner":
				p.OwnerGroup = 2
			case "rejected":
				accepted := false
				p.Evidence.Accepted = &accepted
			case "invalid recorded sample":
				p.Implementation.Expected.RandomInput.Maximum = 6
			case "replaced evidence":
				p.Implementation.Expected.Stdout = writeArtifact(t, root, "other.txt", "8\n 8\n")
			}
			err := validateProbe(root, p, "evidence", map[string]bool{"13.6": false}, nil)
			if (err == nil) != (name == "accepted") {
				t.Fatalf("error=%v", err)
			}
			if name == "replaced evidence" && !strings.Contains(err.Error(), "expected stdout differs") {
				t.Fatalf("exact evidence check lost: %v", err)
			}
		})
	}
}
