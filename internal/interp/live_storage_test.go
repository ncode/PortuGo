package interp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/PortuGo/internal/diag"
	"github.com/ncode/PortuGo/internal/runtime"
	"github.com/ncode/PortuGo/internal/token"
)

func TestLiveStorageLimit(t *testing.T) {
	src := "algoritmo \"storage\"\nvar a, b: vetor[1..1048576] de inteiro\ninicio\nfimalgoritmo\n"
	p, info := analyzed(t, src)
	i := New(Options{MaxSteps: 1})
	ds := i.Run(p, info)
	if len(ds) != 1 || ds[0].Code != diag.RStorage || ds[0].Pos != token.Pos(strings.Index(src, "b:")) {
		t.Fatalf("diagnostics = %v, want R003 before allocating b", ds)
	}
}

func TestLiveStorageAccounting(t *testing.T) {
	for _, tt := range []struct {
		name, src, at string
		limit, used   int
	}{
		{"cumulative", "var a, b: vetor[1..3] de inteiro\ninicio", "b:", 9, 5},
		{"nested", "var v: vetor[1..2] de vetor[1..2] de inteiro\ninicio", "v:", 9, 0},
		{"constant bound", "const n = 3\nvar v: vetor[1..n] de inteiro\ninicio", "v:", 5, 1},
		{"exact boundary", "var v: vetor[1..3] de inteiro\ninicio", "", 5, 5},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := "algoritmo \"storage\"\n" + tt.src + "\nfimalgoritmo\n"
			p, info := analyzed(t, src)
			i := New(Options{})
			i.storageLimit = tt.limit
			ds := i.Run(p, info)
			if tt.at == "" {
				if len(ds) != 0 {
					t.Fatal(ds)
				}
			} else if len(ds) != 1 || ds[0].Code != diag.RStorage || ds[0].Pos != token.Pos(strings.Index(src, tt.at)) {
				t.Fatalf("diagnostics = %v, want R003 at %q", ds, tt.at)
			}
			if i.storageUsed != tt.used {
				t.Fatalf("used = %d, want %d", i.storageUsed, tt.used)
			}
		})
	}
}

func TestLiveStorageCallLifetimes(t *testing.T) {
	for _, tt := range []struct {
		name, body string
		limit      int
		code       diag.Code
	}{
		{"sequential", "p(0)\np(0)\np(0)", 5, ""},
		{"recursive", "p(1)", 9, diag.RStorage},
		{"failure", "p(2)", 20, diag.RArithmetic},
		{"halt", "p(3)", 20, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := `algoritmo "frames"
procedimento stop(x: inteiro)
inicio
  escreval("missing:", x)
fimprocedimento
procedimento p(n: inteiro)
var v: vetor[1..2] de inteiro
inicio
  se n = 1 entao
    p(0)
  fimse
  se n = 2 entao
    escreval(1 / 0)
  fimse
  se n = 3 entao
    stop()
  fimse
fimprocedimento
inicio
` + tt.body + "\nfimalgoritmo\n"
			p, info := analyzed(t, src)
			i := New(Options{MaxSteps: 1000})
			i.storageLimit = tt.limit
			ds := i.Run(p, info)
			if tt.code == "" {
				if len(ds) != 0 {
					t.Fatal(ds)
				}
			} else if len(ds) != 1 || ds[0].Code != tt.code {
				t.Fatalf("diagnostics = %v, want %s", ds, tt.code)
			}
			if i.storageUsed != 0 || len(i.retainedStorage) != 0 {
				t.Fatalf("frame storage leaked: used=%d retained=%d", i.storageUsed, len(i.retainedStorage))
			}
			fresh, freshInfo := analyzed(t, "algoritmo \"fresh\"\nvar v: vetor[1..3] de inteiro\ninicio\nfimalgoritmo\n")
			i.storageLimit = 5
			if ds := i.Run(fresh, freshInfo); len(ds) != 0 || i.storageUsed != 5 {
				t.Fatalf("new run retained old charges: %v used=%d", ds, i.storageUsed)
			}
		})
	}
}

func TestLiveStorageRetainedCalleeRecords(t *testing.T) {
	src := `algoritmo "retained"
funcao f(): inteiro
tipo t = registro
  x: inteiro
fimregistro
var local: t
inicio
  escreva(local = 1)
  retorne 1
fimfuncao
inicio
  escreval(exp(f(), f()))
fimalgoritmo
`
	p, info := analyzed(t, src)
	for _, limit := range []int{5, 6} {
		var out bytes.Buffer
		i := New(Options{Output: &out})
		i.storageLimit = limit
		ds := i.Run(p, info)
		if limit == 5 {
			if len(ds) != 1 || ds[0].Code != diag.RStorage || out.String() != " 1\n" {
				t.Fatalf("retained record lost its charge: %v output=%q", ds, &out)
			}
		} else if len(ds) != 0 {
			t.Fatal(ds)
		}
		if i.storageUsed != 0 || len(i.retainedStorage) != 0 {
			t.Fatalf("retained storage leaked: used=%d retained=%d", i.storageUsed, len(i.retainedStorage))
		}
	}
}

func TestLiveStorageRecordCopy(t *testing.T) {
	src := `algoritmo "copy"
tipo t = registro
  x: inteiro
fimregistro
var a, b: t
inicio
  a.x <- 7
  b.x <- 9
  b <- a
fimalgoritmo
`
	p, info := analyzed(t, src)
	for _, limit := range []int{7, 8} {
		i := New(Options{})
		i.storageLimit = limit
		ds := i.Run(p, info)
		want := int64(7)
		if limit == 7 {
			want = 9
			if len(ds) != 1 || ds[0].Code != diag.RStorage || ds[0].Pos != token.Pos(strings.Index(src, "b <-")) {
				t.Fatalf("copy not guarded: %v", ds)
			}
		} else if len(ds) != 0 {
			t.Fatal(ds)
		}
		if i.State()["b"].Rec.Fields[0].Value.Int != want || i.storageUsed != 6 {
			t.Fatalf("copy changed destination or leaked a temporary: used=%d state=%v", i.storageUsed, i.State())
		}
	}
}

func TestLiveStorageCopyBackRetainsPendingValue(t *testing.T) {
	typ := runtime.Type{Kind: runtime.RecordType}
	source := &runtime.Cell{Type: typ, Value: runtime.Zero(typ)}
	destination := &runtime.Cell{Type: typ, Value: runtime.Zero(typ)}
	i := New(Options{})
	i.storageUsed, i.storageLimit = 4, 5
	i.operands = append(i.operands, destination.Value)
	if err := i.copyCell(destination, source); err != nil || i.storageUsed != 5 {
		t.Fatalf("pending value lost its charge: %v used=%d", err, i.storageUsed)
	}
	i.operands = nil
	i.releaseRetainedStorage()
	if i.storageUsed != 4 || len(i.retainedStorage) != 0 {
		t.Fatalf("pending value charge leaked: used=%d", i.storageUsed)
	}
}

func TestLiveStorageRecordFieldRetention(t *testing.T) {
	typ := runtime.Type{Kind: runtime.RecordType, RecordID: 1, Fields: []runtime.Field{{Name: "x", Type: runtime.Type{Kind: runtime.BoolType}}}}
	source := &runtime.Cell{Type: typ, Value: runtime.Zero(typ)}
	destination := &runtime.Cell{Type: typ, Value: runtime.Zero(typ)}
	field := &destination.Value.Rec.Fields[0]
	field.Type = runtime.Type{Kind: runtime.RecordType}
	field.Value = runtime.Zero(field.Type)
	i := New(Options{})
	i.storageUsed, i.storageLimit = 7, 9
	i.operands = append(i.operands, field.Value)
	if err := i.assign(destination, source.Value); err != nil || i.storageUsed != 7 || field != &destination.Value.Rec.Fields[0] {
		t.Fatalf("field reference or pending charge lost: %v used=%d", err, i.storageUsed)
	}
	i.operands = nil
	i.releaseRetainedStorage()
	if i.storageUsed != 6 || len(i.retainedStorage) != 0 {
		t.Fatalf("field's pending charge leaked: used=%d", i.storageUsed)
	}
}

func TestLiveStorageBuiltinArguments(t *testing.T) {
	src := `algoritmo "pending arguments"
tipo t = registro
  b: logico
fimregistro
var r, s: t
procedimento g()
var v: vetor[1..6] de inteiro
inicio
fimprocedimento
funcao f(): inteiro
inicio
  r <- s
  g()
  retorne 1
fimfuncao
inicio
  r.b <- 1 = r
  escreval(exp(r.b, f()))
fimalgoritmo
`
	p, info := analyzed(t, src)
	for _, limit := range []int{15, 16} {
		i := New(Options{})
		i.storageLimit = limit
		ds := i.Run(p, info)
		if limit == 15 {
			if len(ds) != 1 || ds[0].Code != diag.RStorage || ds[0].Pos != token.Pos(strings.Index(src, "v:")) {
				t.Fatalf("builtin argument lost its charge: %v", ds)
			}
		} else if len(ds) != 1 || ds[0].Code != diag.RBuiltin {
			t.Fatalf("expected numeric argument diagnostic after successful calls: %v", ds)
		}
		if i.storageUsed != 6 || len(i.retainedStorage) != 0 || len(i.argumentRoots) != 0 {
			t.Fatalf("builtin argument charge leaked: used=%d retained=%d", i.storageUsed, len(i.retainedStorage))
		}
	}
}
