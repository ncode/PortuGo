package main

import (
	"math"
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

func TestRandomOutputReplay(t *testing.T) {
	t.Parallel()
	root, m := testManifest(t)
	p := m.Probes[0]
	p.TimeoutMS = 5000
	p.Source = writeArtifact(t, root, p.Source.Path, "random-output")
	p.Implementation.Expected.Stdout = writeArtifact(t, root, "random-output.txt", " 4\n 6\n 8\n 2\n 1\n 7\n 4\n 6\n 0\n 0\n")
	p.Implementation.Expected.RandomOutput = &randomOutputExpectation{Kind: "randi-lines", Lines: 10, Bound: 10, Review: review{Reason: "Generated lines are qualified by domain and framing.", Link: "tasks.md"}}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err := replayProbe(root, p, executable, []string{"-test.run=^TestReplayChild$", "--"}, ""); err != nil {
		t.Fatal(err)
	}
	p.Source = writeArtifact(t, root, p.Source.Path, "random-output-invalid")
	if err := replayProbe(root, p, executable, []string{"-test.run=^TestReplayChild$", "--"}, ""); err == nil || !strings.Contains(err.Error(), "replayed random output") {
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

func TestRandomOutputContract(t *testing.T) {
	t.Parallel()
	valid := " 4\n 6\n 8\n 2\n 1\n 7\n 4\n 6\n 0\n 0\n"
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain and framing", valid, true},
		{"upper bound", " 4\n 6\n 8\n 2\n 1\n 7\n 4\n 6\n 0\n 10\n", false},
		{"fixed tail", " 4\n 6\n 8\n 2\n 1\n 7\n 4\n 6\n 0\n 1\n", false},
		{"missing leading space", "4\n 6\n 8\n 2\n 1\n 7\n 4\n 6\n 0\n 0\n", false},
		{"extra leading space", "  4\n 6\n 8\n 2\n 1\n 7\n 4\n 6\n 0\n 0\n", false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			contract := randomOutputExpectation{Kind: "randi-lines", Lines: 10, Bound: 10, FixedTail: 0}
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "unknown", Lines: 10, Bound: 10}).compare([]byte(valid)); err == nil {
		t.Fatal("invalid contract accepted")
	}
	if err := (randomOutputExpectation{Kind: "randi-lines", Lines: 1, Bound: int64(math.MaxInt32) + 1}).compare([]byte(" 0\n")); err == nil {
		t.Fatal("bound beyond randi signed-32-bit range accepted")
	}
}

func TestRandomIntegerTextOutputContract(t *testing.T) {
	t.Parallel()
	valid := "70\nNHEEX\n99\nMWFJX\n36\nTCWSY\n97\nJVXUJ\n80\nUWQIS\n"
	contract := randomOutputExpectation{Kind: "random-int-text-lines", Lines: 10, Bound: 101}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain and framing", valid, true},
		{"integer upper bound", "70\nNHEEX\n101\nMWFJX\n36\nTCWSY\n97\nJVXUJ\n80\nUWQIS\n", false},
		{"negative integer", "70\nNHEEX\n-1\nMWFJX\n36\nTCWSY\n97\nJVXUJ\n80\nUWQIS\n", false},
		{"leading zero", "070\nNHEEX\n99\nMWFJX\n36\nTCWSY\n97\nJVXUJ\n80\nUWQIS\n", false},
		{"lowercase text", "70\nNHEEX\n99\nmwfjx\n36\nTCWSY\n97\nJVXUJ\n80\nUWQIS\n", false},
		{"wrong text length", "70\nNHEEX\n99\nMWFJ\n36\nTCWSY\n97\nJVXUJ\n80\nUWQIS\n", false},
		{"wrong line count", "70\nNHEEX\n99\nMWFJX\n36\nTCWSY\n97\nJVXUJ\n", false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-int-text-lines", Lines: 10, Bound: 101, FixedTail: 1}).compare([]byte(valid)); err == nil {
		t.Fatal("fixed tail accepted for mixed random contract")
	}
}

func TestRandomIntegerSortOutputContract(t *testing.T) {
	t.Parallel()
	valid := "79\n69\n34\n76\n63\n67\n44\n51\n91\n31\n18\n51\n4\n89\n29\n64\n92\n25\n78\n42\n 4\n 18\n 25\n 29\n 31\n 34\n 42\n 44\n 51\n 51\n 63\n 64\n 67\n 69\n 76\n 78\n 79\n 89\n 91\n 92\n"
	contract := randomOutputExpectation{Kind: "random-int-sort-lines", Lines: 40, Minimum: 1, Bound: 101}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain and sorted permutation", valid, true},
		{"input below domain", strings.Replace(valid, "79\n", "0\n", 1), false},
		{"input above domain", strings.Replace(valid, "79\n", "101\n", 1), false},
		{"unsorted output", strings.Replace(valid, " 18\n 25\n", " 25\n 18\n", 1), false},
		{"wrong permutation", strings.Replace(valid, " 92\n", " 93\n", 1), false},
		{"missing output space", strings.Replace(valid, "\n 4\n", "\n4\n", 1), false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-int-sort-lines", Lines: 3, Minimum: 1, Bound: 101}).compare([]byte("1\n 1\n")); err == nil {
		t.Fatal("odd sort line count accepted")
	}
}

func TestRandomRealSortOutputContract(t *testing.T) {
	t.Parallel()
	valid := "99.2680000000\n7.8610000000\n44.1160000000\n21.3800000000\n3.5940000000\n40.1850000000\n15.1450000000\n50.2560000000\n85.0640000000\n20.0420000000\n  1 -      3.594\n  2 -      7.861\n  3 -     15.145\n  4 -     20.042\n  5 -     21.380\n  6 -     40.185\n  7 -     44.116\n  8 -     50.256\n  9 -     85.064\n 10 -     99.268\n"
	contract := randomOutputExpectation{Kind: "random-real-sort-lines", Lines: 20, Minimum: 0, Bound: 101000, Decimals: 3}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain and sorted permutation", valid, true},
		{"input below domain", strings.Replace(valid, "99.2680000000", "-0.0010000000", 1), false},
		{"input above domain", strings.Replace(valid, "99.2680000000", "101.0000000000", 1), false},
		{"input off precision grid", strings.Replace(valid, "7.8610000000", "7.8615000000", 1), false},
		{"unsorted output", strings.Replace(valid, "  2 -      7.861\n  3 -     15.145\n", "  2 -     15.145\n  3 -      7.861\n", 1), false},
		{"wrong row number", strings.Replace(valid, "  1 -      3.594", "  2 -      3.594", 1), false},
		{"wrong output precision", strings.Replace(valid, "  1 -      3.594", "  1 -      3.5940", 1), false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-real-sort-lines", Lines: 20, Minimum: 0, Bound: 101000, Decimals: 0}).compare([]byte(valid)); err == nil {
		t.Fatal("zero-decimal real sort contract accepted")
	}
}

func TestRandomRecordSortOutputContract(t *testing.T) {
	t.Parallel()
	valid := "Digite o codigo do  1o registro:60\nDigite o nome do  1o registro:HEGLU\nDigite o codigo do  2o registro:3\nDigite o nome do  2o registro:YUDCK\nDigite o codigo do  3o registro:49\nDigite o nome do  3o registro:QNYWK\nDigite o codigo do  4o registro:54\nDigite o nome do  4o registro:XWVRH\nDigite o codigo do  5o registro:10\nDigite o nome do  5o registro:FJVCP\nDigite o codigo do  6o registro:18\nDigite o nome do  6o registro:IPYCX\nDigite o codigo do  7o registro:42\nDigite o nome do  7o registro:KGLFO\nDigite o codigo do  8o registro:46\nDigite o nome do  8o registro:GWOMD\nDigite o codigo do  9o registro:7\nDigite o nome do  9o registro:HJYLV\nDigite o codigo do  10o registro:75\nDigite o nome do  10o registro:WDBXW\nItem - Codigo Nome\n   1 -     10 FJVCP\n   2 -     46 GWOMD\n   3 -     60 HEGLU\n   4 -      7 HJYLV\n   5 -     18 IPYCX\n   6 -     42 KGLFO\n   7 -     49 QNYWK\n   8 -     75 WDBXW\n   9 -     54 XWVRH\n  10 -      3 YUDCK\nItem - Codigo Nome\n   1 -      3 YUDCK\n   2 -      7 HJYLV\n   3 -     10 FJVCP\n   4 -     18 IPYCX\n   5 -     42 KGLFO\n   6 -     46 GWOMD\n   7 -     49 QNYWK\n   8 -     54 XWVRH\n   9 -     60 HEGLU\n  10 -     75 WDBXW\n"
	contract := randomOutputExpectation{Kind: "random-record-sort-lines", Lines: 42, Minimum: 0, Bound: 101}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain and sorted permutations", valid, true},
		{"code below domain", strings.Replace(valid, ":60\n", ":-1\n", 1), false},
		{"code above domain", strings.Replace(valid, ":60\n", ":101\n", 1), false},
		{"name alphabet", strings.Replace(valid, ":HEGLU\n", ":HEGLU1\n", 1), false},
		{"name sort", strings.Replace(valid, "   1 -     10 FJVCP\n   2 -     46 GWOMD\n", "   1 -     46 GWOMD\n   2 -     10 FJVCP\n", 1), false},
		{"code sort", strings.Replace(valid, "   1 -      3 YUDCK\n   2 -      7 HJYLV\n", "   1 -      7 HJYLV\n   2 -      3 YUDCK\n", 1), false},
		{"wrong header", strings.Replace(valid, "Item - Codigo Nome", "Item - Codigo", 1), false},
		{"wrong row width", strings.Replace(valid, "   1 -     10 FJVCP", "  1 -     10 FJVCP", 1), false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-record-sort-lines", Lines: 40, Minimum: 0, Bound: 101}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong record sort line count accepted")
	}
}

func TestRandomIntegerSearchTableOutputContract(t *testing.T) {
	t.Parallel()
	valid := "    1   83\n    2   97\n    3   39\n    4   95\n    5   39\n    6   16\n    7   92\n    8   61\n    9   85\n   10   67\nEntre com o valor de busca (ESC termina) :-1\nNao achei.\n"
	contract := randomOutputExpectation{Kind: "random-int-search-table-lines", Lines: 12, Minimum: 0, Bound: 100}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain and not-found search", valid, true},
		{"value below domain", strings.Replace(valid, "    1   83\n", "    1   -1\n", 1), false},
		{"value at excluded bound", strings.Replace(valid, "    1   83\n", "    1  100\n", 1), false},
		{"leading zero", strings.Replace(valid, "    1   83\n", "    1  083\n", 1), false},
		{"wrong row number", strings.Replace(valid, "    1   83\n", "    2   83\n", 1), false},
		{"wrong row width", strings.Replace(valid, "    1   83\n", "   1   83\n", 1), false},
		{"wrong search input", strings.Replace(valid, "Entre com o valor de busca (ESC termina) :-1", "Entre com o valor de busca (ESC termina) :0", 1), false},
		{"found result", strings.Replace(valid, "Nao achei.", "Achei -1.", 1), false},
		{"missing row", strings.Replace(valid, "   10   67\n", "", 1), false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-int-search-table-lines", Lines: 10, Bound: 100}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong search table line count accepted")
	}
	if err := (randomOutputExpectation{Kind: "random-int-search-table-lines", Lines: 12, Bound: 0}).compare([]byte(valid)); err == nil {
		t.Fatal("zero search table bound accepted")
	}
}

func TestRandomIntegerSearchOutputContract(t *testing.T) {
	t.Parallel()
	valid := "71\n41\n25\n92\n76\n26\n73\n37\n84\n39\n44\n6\n87\n100\n43\n95\n49\n62\n57\n99\nValor para busca (ESC ou menor que 0 termina) : -1\n"
	contract := randomOutputExpectation{Kind: "random-int-search-lines", Lines: 21, Minimum: 0, Bound: 101}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain and exit prompt", valid, true},
		{"value below domain", strings.Replace(valid, "71\n", "-1\n", 1), false},
		{"value above domain", strings.Replace(valid, "71\n", "101\n", 1), false},
		{"leading zero", strings.Replace(valid, "71\n", "071\n", 1), false},
		{"wrong prompt input", strings.Replace(valid, " : -1\n", " : 0\n", 1), false},
		{"unexpected search result", valid + "Nao achei.\n", false},
		{"missing generated line", strings.Replace(valid, "99\nValor para busca", "Valor para busca", 1), false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-int-search-lines", Lines: 20, Bound: 101}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong search line count accepted")
	}
	if err := (randomOutputExpectation{Kind: "random-int-search-lines", Lines: 21, Bound: 0}).compare([]byte(valid)); err == nil {
		t.Fatal("zero search bound accepted")
	}
}

func TestRandomIntegerSearchSortedOutputContract(t *testing.T) {
	t.Parallel()
	valid := "    1    1\n    2   10\n    3   17\n    4   20\n    5   27\n    6   30\n    7   34\n    8   44\n    9   48\n   10   60\n   11   61\n   12   64\n   13   74\n   14   74\n   15   76\n   16   83\n   17   86\n   18   92\n   19   92\n   20   97\nValor para busca (ESC ou menor que 0 termina) : -1\n"
	contract := randomOutputExpectation{Kind: "random-int-search-sorted-lines", Lines: 21, Minimum: 0, Bound: 101}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain, order and framing", valid, true},
		{"value below domain", strings.Replace(valid, "    1    1\n", "    1   -1\n", 1), false},
		{"value at excluded bound", strings.Replace(valid, "   20   97\n", "   20  101\n", 1), false},
		{"leading zero", strings.Replace(valid, "    1    1\n", "    1   01\n", 1), false},
		{"unsorted values", strings.Replace(valid, "    2   10\n", "    2    0\n", 1), false},
		{"wrong row number", strings.Replace(valid, "    2   10\n", "    3   10\n", 1), false},
		{"wrong column spacing", strings.Replace(valid, "    1    1\n", "    1     1\n", 1), false},
		{"wrong prompt input", strings.Replace(valid, " : -1\n", " : 0\n", 1), false},
		{"wrong prompt framing", strings.Replace(valid, " : -1\n", ": -1\n", 1), false},
		{"unexpected search result", valid + "Nao achei.\n", false},
		{"missing generated line", strings.Replace(valid, "   20   97\n", "", 1), false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-int-search-sorted-lines", Lines: 20, Minimum: 0, Bound: 101}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong sorted search line count accepted")
	}
	if err := (randomOutputExpectation{Kind: "random-int-search-sorted-lines", Lines: 21, Minimum: 1, Bound: 101}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong sorted search minimum accepted")
	}
	if err := (randomOutputExpectation{Kind: "random-int-search-sorted-lines", Lines: 21, Minimum: 0, Bound: 100}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong sorted search bound accepted")
	}
}

func TestRandomRecordSearchOutputContract(t *testing.T) {
	t.Parallel()
	valid := "    1   82    597.84\n    2   83    294.13\n    3   61    311.13\n    4   83    455.65\n    5   11    920.21\n    6   19    452.77\n    7   21    502.31\n    8   74    416.78\n    9   16    356.91\n   10   14    262.12\nEntre com o valor de busca (ESC termina) :-1\nNao achei.\n"
	contract := randomOutputExpectation{Kind: "random-record-search-lines", Lines: 12, Minimum: 0, Bound: 100, Decimals: 2}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain, framing and not-found result", valid, true},
		{"code below domain", strings.Replace(valid, "    1   82", "    1   -1", 1), false},
		{"code at excluded bound", strings.Replace(valid, "    1   82", "    1  100", 1), false},
		{"code leading zero", strings.Replace(valid, "    1   82", "    1  082", 1), false},
		{"salary below domain", strings.Replace(valid, "    597.84", "     -1.00", 1), false},
		{"salary at excluded bound", strings.Replace(valid, "    597.84", "   1000.00", 1), false},
		{"salary wrong precision", strings.Replace(valid, "    597.84", "     597.8", 1), false},
		{"salary leading zero", strings.Replace(valid, "    597.84", "   0597.84", 1), false},
		{"wrong row number", strings.Replace(valid, "    1   82", "    2   82", 1), false},
		{"wrong row spacing", strings.Replace(valid, "    1   82", "   1   82", 1), false},
		{"wrong search input", strings.Replace(valid, ":-1\n", ":0\n", 1), false},
		{"found result", strings.Replace(valid, "Nao achei.", "Achei.", 1), false},
		{"missing record", strings.Replace(valid, "   10   14    262.12\n", "", 1), false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-record-search-lines", Lines: 11, Minimum: 0, Bound: 100, Decimals: 2}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong record-search line count accepted")
	}
	if err := (randomOutputExpectation{Kind: "random-record-search-lines", Lines: 12, Minimum: 1, Bound: 100, Decimals: 2}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong record-search minimum accepted")
	}
	if err := (randomOutputExpectation{Kind: "random-record-search-lines", Lines: 12, Minimum: 0, Bound: 100, Decimals: 0}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong record-search precision accepted")
	}
}

func TestRandomIntegerCountSortOutputContract(t *testing.T) {
	t.Parallel()
	valid := "91\n26\n52\n78\n41\n5\n79\n69\n81\n32\n14\n27\n85\n45\n5\n46\n63\n34\n33\n62\n\nCronômetro iniciado.\n\nCronômetro terminado. Tempo decorrido: 0 segundo(s).\nv2[ 1] =  5\nv2[ 2] =  5\nv2[ 3] =  14\nv2[ 4] =  26\nv2[ 5] =  27\nv2[ 6] =  32\nv2[ 7] =  33\nv2[ 8] =  34\nv2[ 9] =  41\nv2[ 10] =  45\nv2[ 11] =  46\nv2[ 12] =  52\nv2[ 13] =  62\nv2[ 14] =  63\nv2[ 15] =  69\nv2[ 16] =  78\nv2[ 17] =  79\nv2[ 18] =  81\nv2[ 19] =  85\nv2[ 20] =  91\n"
	contract := randomOutputExpectation{Kind: "random-int-countsort-lines", Lines: 44, Minimum: 1, Bound: 101}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain, sorting and framing", valid, true},
		{"milliseconds elapsed shape", strings.Replace(valid, "0 segundo(s).", "16 ms.", 1), true},
		{"fractional seconds elapsed shape", strings.Replace(valid, "0 segundo(s).", "2 segundo(s) e 47 ms.", 1), true},
		{"exact seconds elapsed shape", strings.Replace(valid, "0 segundo(s).", "2 segundo(s).", 1), true},
		{"elapsed zero milliseconds rejected", strings.Replace(valid, "0 segundo(s).", "0 ms.", 1), false},
		{"elapsed zero remainder rejected", strings.Replace(valid, "0 segundo(s).", "1 segundo(s) e 0 ms.", 1), false},
		{"elapsed over one second rejected as milliseconds", strings.Replace(valid, "0 segundo(s).", "1000 ms.", 1), false},
		{"input below domain", strings.Replace(valid, "91\n", "0\n", 1), false},
		{"input above domain", strings.Replace(valid, "91\n", "101\n", 1), false},
		{"input leading zero", strings.Replace(valid, "91\n", "091\n", 1), false},
		{"unsorted output", strings.Replace(valid, "v2[ 2] =  5\n", "v2[ 2] =  4\n", 1), false},
		{"wrong output permutation", strings.Replace(valid, "v2[ 2] =  5\n", "v2[ 2] =  6\n", 1), false},
		{"wrong row number", strings.Replace(valid, "v2[ 2] =  5\n", "v2[ 3] =  5\n", 1), false},
		{"wrong row spacing", strings.Replace(valid, "v2[ 1] =  5\n", "v2[1] =  5\n", 1), false},
		{"missing start blank line", strings.Replace(valid, "62\n\nCronômetro", "62\nCronômetro", 1), false},
		{"wrong start message", strings.Replace(valid, "Cronômetro iniciado.", "Cronômetro começou.", 1), false},
		{"wrong stop message", strings.Replace(valid, "Cronômetro terminado.", "Cronômetro concluído.", 1), false},
		{"missing output row", strings.Replace(valid, "v2[ 20] =  91\n", "", 1), false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"missing final newline", strings.TrimSuffix(valid, "\n"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-int-countsort-lines", Lines: 43, Minimum: 1, Bound: 101}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong countsort line count accepted")
	}
	if err := (randomOutputExpectation{Kind: "random-int-countsort-lines", Lines: 44, Minimum: 0, Bound: 101}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong countsort minimum accepted")
	}
	if err := (randomOutputExpectation{Kind: "random-int-countsort-lines", Lines: 44, Minimum: 1, Bound: 100}).compare([]byte(valid)); err == nil {
		t.Fatal("wrong countsort bound accepted")
	}
}

func TestRandomRandiRepeatOutputContract(t *testing.T) {
	t.Parallel()
	valid := " \n ============================================== \nQUANTOS NUMEROS (1-10): 1\nDigite o destaque: 0\n \nA SEQUENCIA É \n 1 \nO NUMERO DE REPETIÇÕES FOI:  0\n \nRESTOU A SEQUENCIA: \n 1"
	contract := randomOutputExpectation{Kind: "random-randi-repeat-lines", Lines: 11, Minimum: 1, Bound: 10}
	for _, tt := range []struct {
		name   string
		output string
		valid  bool
	}{
		{"domain and repeated value", valid, true},
		{"value below domain", strings.Replace(valid, " 1 \n", " 0 \n", 1), false},
		{"value at excluded bound", strings.Replace(valid, " 1 \n", " 10 \n", 1), false},
		{"repeated value mismatch", strings.TrimSuffix(valid, " 1") + " 2", false},
		{"wrong prompt input", strings.Replace(valid, "QUANTOS NUMEROS (1-10): 1", "QUANTOS NUMEROS (1-10): 2", 1), false},
		{"wrong repeat count", strings.Replace(valid, "REPETIÇÕES FOI:  0", "REPETIÇÕES FOI:  1", 1), false},
		{"missing value spacing", strings.Replace(valid, " 1 \n", "1 \n", 1), false},
		{"wrong static line", strings.Replace(valid, "A SEQUENCIA É ", "A SEQUENCIA", 1), false},
		{"wrong line count", strings.Replace(valid, "Digite o destaque: 0\n", "", 1), false},
		{"carriage returns", strings.ReplaceAll(valid, "\n", "\r\n"), false},
		{"unexpected final newline", valid + "\n", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := contract.compare([]byte(tt.output)); (err == nil) != tt.valid {
				t.Fatalf("error=%v, want valid=%t", err, tt.valid)
			}
		})
	}
	if err := (randomOutputExpectation{Kind: "random-randi-repeat-lines", Lines: 11, Minimum: 1, Bound: 1}).compare([]byte(valid)); err == nil {
		t.Fatal("empty repeated-value domain accepted")
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
