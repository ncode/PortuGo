package interp

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/parser"
	"github.com/ncode/portugol-go/internal/sema"
	"github.com/ncode/portugol-go/internal/testprocess"
	"github.com/ncode/portugol-go/internal/token"
)

func analyzed(t *testing.T, src string) (*ast.Program, *sema.Info) {
	t.Helper()
	_, toks, ds := lexer.Scan("test.alg", src)
	if diag.HasErrors(ds) {
		t.Fatal(ds)
	}
	p, ds := parser.Parse(toks)
	if diag.HasErrors(ds) {
		t.Fatal(ds)
	}
	info, ds := sema.Analyze(p)
	if diag.HasErrors(ds) {
		t.Fatal(ds)
	}
	return p, info
}

func TestExecutionDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		name, declarations, body, at string
		code                         diag.Code
		options                      Options
	}{
		{"arithmetic", "", "escreva(\"before\")\nescreval(1 / 0)", "/", diag.RArithmetic, Options{}},
		{"storage", "var v: vetor[1..2] de inteiro", "v[3] <- 1", "v[3]", diag.RStorage, Options{}},
		{"input", "var x: inteiro", "leia(x)", "x)", diag.RInput, Options{}},
		{"loop", "var x: inteiro", "para x de 1 ate 2 passo 0 faca\nfimpara", "para", diag.RLoop, Options{}},
		{"builtin", "", "escreval(aleatorio(0))", "aleatorio", diag.RBuiltin, Options{}},
		{"text conversion", "", "escreval(numpcarac(10 ^ 400))", "numpcarac", diag.RBuiltin, Options{}},
		{"output", "", "escreva(1)", "1)", diag.RHost, Options{Output: failingWriter{}}},
		{"input echo", "var x: inteiro", "leia(x)", "x)", diag.RHost, Options{Input: strings.NewReader("7\n"), Output: failingWriter{}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			src := "algoritmo \"diagnostic\"\n" + tt.declarations + "\ninicio\n" + tt.body + "\nfimalgoritmo"
			p, info := analyzed(t, src)
			var out bytes.Buffer
			if tt.options.Output == nil {
				tt.options.Output = &out
			}
			ds := New(tt.options).Run(p, info)
			if len(ds) != 1 || ds[0].Code != tt.code || ds[0].Pos != token.Pos(strings.Index(src, tt.at)) {
				t.Fatalf("got %v, want positioned %s at %d", ds, tt.code, strings.Index(src, tt.at))
			}
			if tt.name == "arithmetic" && out.String() != "before" {
				t.Fatalf("lost partial output: %q", &out)
			}
		})
	}
}

func TestRejectMissingOrStaleInfo(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"x\"\ninicio\nescreval(1)\nfimalgoritmo")
	other, _ := analyzed(t, "algoritmo \"x\"\ninicio\nfimalgoritmo")
	for _, pair := range []struct {
		p    *ast.Program
		info *sema.Info
	}{{nil, nil}, {p, nil}, {other, info}} {
		ds := New(Options{}).Run(pair.p, pair.info)
		if len(ds) != 1 || ds[0].Code != diag.RType {
			t.Fatalf("unchecked run: %v", ds)
		}
	}
	p.Body[0].(*ast.WriteStmt).Args[0].Expr = &ast.LiteralExpr{At: 10, Kind: ast.IntLiteral, Int: 2}
	ds := New(Options{}).Run(p, info)
	if len(ds) != 1 || ds[0].Code != diag.RType {
		t.Fatalf("missing fact accepted: %v", ds)
	}
}

func TestStepBudgetAndReuse(t *testing.T) {
	testprocess.Run(t, func() {
		p, info := analyzed(t, "algoritmo \"x\"\ninicio\nescreva(1)\nfimalgoritmo")
		for _, limit := range []uint64{0, 1, 2} {
			var out bytes.Buffer
			i := New(Options{MaxSteps: limit, Output: &out})
			for run := 0; run < 2; run++ {
				out.Reset()
				ds := i.Run(p, info)
				if limit == 1 {
					if len(ds) != 1 || ds[0].Code != diag.RLoop || out.Len() != 0 {
						t.Fatalf("side effect after exhaustion: %v %q", ds, &out)
					}
				} else if len(ds) != 0 || out.String() != " 1" {
					t.Fatalf("counter not reset: %v %q", ds, &out)
				}
			}
		}
		for _, body := range []string{"enquanto verdadeiro faca\nfimenquanto", "repita\nate falso", "para n de 1 ate 1000000 faca\nfimpara"} {
			p, info := analyzed(t, "algoritmo \"x\"\nvar n: inteiro\ninicio\n"+body+"\nfimalgoritmo")
			ds := New(Options{MaxSteps: 10}).Run(p, info)
			if len(ds) != 1 || ds[0].Code != diag.RLoop {
				t.Fatalf("unbounded empty loop: %v", ds)
			}
		}
	})
}

func TestCallAndValueLimits(t *testing.T) {
	testprocess.Run(t, func() {
		for _, tt := range []struct {
			name, src string
			code      diag.Code
		}{
			{"recursion", "algoritmo \"x\"\nprocedimento p()\ninicio\np()\nfimprocedimento\ninicio\np()\nfimalgoritmo", diag.RCall},
			{"bare function recursion", "algoritmo \"x\"\nfuncao f: inteiro\ninicio\nretorne f\nfimfuncao\ninicio\nescreval(f)\nfimalgoritmo", diag.RCall},
			{"text", "algoritmo \"x\"\nvar s: caractere\ninicio\ns <- \"a\"\nenquanto verdadeiro faca\ns <- s+s\nfimenquanto\nfimalgoritmo", diag.RStorage},
			{"format", "algoritmo \"x\"\ninicio\nescreva(1:9223372036854775807)\nfimalgoritmo", diag.RStorage},
		} {
			t.Run(tt.name, func(t *testing.T) {
				p, info := analyzed(t, tt.src)
				i := New(Options{MaxSteps: 10000})
				ds := i.Run(p, info)
				if len(ds) != 1 || ds[0].Code != tt.code {
					t.Fatalf("missing limit diagnostic: %v", ds)
				}
				p, info = analyzed(t, "algoritmo \"ok\"\ninicio\nescreval(1)\nfimalgoritmo")
				if ds := i.Run(p, info); len(ds) != 0 {
					t.Fatalf("failed run polluted interpreter: %v", ds)
				}
			})
		}
	})
}

type scriptedRandom struct{ bounds []uint64 }

func (*scriptedRandom) Float64() float64          { return .25 }
func (r *scriptedRandom) Uint64N(n uint64) uint64 { r.bounds = append(r.bounds, n); return n - 1 }

func TestOptionsRandomAndDefaults(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"random\"\ninicio\nescreval(aleatorio(), aleatorio(4), aleatorio(3, 5))\nfimalgoritmo")
	r := &scriptedRandom{}
	var out bytes.Buffer
	i := New(Options{Random: r, Output: &out})
	if ds := i.Run(p, info); len(ds) != 0 || out.String() != " 0.25 3 5\n" {
		t.Fatalf("random injection: %v %q", ds, &out)
	}
	if len(r.bounds) != 2 || r.bounds[0] != 4 || r.bounds[1] != 3 {
		t.Fatalf("random bounds: %v", r.bounds)
	}
	if ds := New(Options{}).Run(p, info); len(ds) != 0 {
		t.Fatal(ds)
	}
	h := HeadlessHost{}
	if h.Now().IsZero() {
		t.Fatal("missing clock")
	}
	for _, err := range []error{h.Delay(0 * time.Second), h.Breakpoint(Breakpoint{}), h.ClearScreen(), h.SetDisplay(DisplayState{})} {
		if err != nil {
			t.Fatal(err)
		}
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("writer failed") }

func TestInputBufferOwnership(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"read\"\nvar s: caractere\ninicio\nleia(s)\nescreva(s)\nfimalgoritmo")
	var out bytes.Buffer
	i := New(Options{Input: strings.NewReader("first\nsecond"), Output: &out})
	for run := 0; run < 2; run++ {
		if ds := i.Run(p, info); len(ds) != 0 {
			t.Fatal(ds)
		}
	}
	if out.String() != "first\nfirstsecond\nsecond" {
		t.Fatalf("lost buffered input: %q", &out)
	}
	i = New(Options{Input: io.LimitReader(strings.NewReader(strings.Repeat("x", 16<<20+1)), 16<<20+1)})
	if ds := i.Run(p, info); len(ds) != 1 || ds[0].Code != diag.RStorage {
		t.Fatalf("unbounded input: %v", ds)
	}
}

func TestFormattedItemBoundary(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"width\"\ninicio\nescreva(1:1048576)\nfimalgoritmo")
	var out bytes.Buffer
	ds := New(Options{Output: &out}).Run(p, info)
	if len(ds) != 0 || out.Len() != 1<<20 || out.Bytes()[out.Len()-1] != '1' {
		t.Fatalf("valid width corrupted: diagnostics %v, bytes %d", ds, out.Len())
	}
}

func TestBuiltinTextLimits(t *testing.T) {
	testprocess.Run(t, func() {
		src := "algoritmo \"case\"\nvar s: caractere\nn: inteiro\ninicio\ns <- \"\u023f\"\npara n de 1 ate 23 faca\ns <- s+s\nfimpara\nescreva(maiusc(s))\nfimalgoritmo"
		p, info := analyzed(t, src)
		ds := New(Options{MaxSteps: 10000}).Run(p, info)
		if len(ds) != 1 || ds[0].Code != diag.RStorage || ds[0].Pos != token.Pos(strings.Index(src, "maiusc")) {
			t.Fatalf("case expansion bypassed value limit: %v", ds)
		}
		p, info = analyzed(t, "algoritmo \"copy\"\ninicio\nescreva(copia(\"abc\", 2, 9223372036854775807))\nfimalgoritmo")
		var out bytes.Buffer
		ds = New(Options{Output: &out}).Run(p, info)
		if len(ds) != 0 || out.String() != "bc" {
			t.Fatalf("overflowing copy length: %v %q", ds, &out)
		}
	})
}

func TestStateObservationIsIndependent(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"state\"\nvar v: vetor[1..2] de inteiro\ninicio\nv[2] <- 7\nfimalgoritmo")
	i := New(Options{})
	if ds := i.Run(p, info); len(ds) != 0 {
		t.Fatal(ds)
	}
	state := i.State()
	v := state["v"]
	v.Vec.Elements[1].Value.Int = 9
	v.Vec.Type.Ranges[0].High = 999
	if got := i.State()["v"]; got.Vec.Elements[1].Value.Int != 7 || got.Vec.Type.Ranges[0].High != 2 {
		t.Fatal("observation mutated live state")
	}
}

func TestEachExpressionChargedOnce(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"budget\"\nvar v: vetor[1..1] de inteiro\ninicio\nescreva(v[1])\nfimalgoritmo")
	var out bytes.Buffer
	ds := New(Options{MaxSteps: 4, Output: &out}).Run(p, info)
	if len(ds) != 0 || out.String() != " 0" {
		t.Fatalf("index charged twice: %v %q", ds, &out)
	}
}

func TestBuiltinCallBindingWithLocalName(t *testing.T) {
	p, info := analyzed(t, "algoritmo \"binding\"\nprocedimento p(abs: inteiro)\ninicio\nescreva(abs(-1), abs)\nfimprocedimento\ninicio\np(7)\nfimalgoritmo")
	var out bytes.Buffer
	ds := New(Options{Output: &out}).Run(p, info)
	if len(ds) != 0 || out.String() != " 1 7" {
		t.Fatalf("builtin binding changed: %v %q", ds, &out)
	}
}

func TestDefensiveEvaluationDepth(t *testing.T) {
	testprocess.Run(t, func() {
		p, info := analyzed(t, "algoritmo \"cycle\"\ninicio\nescreva(1+1)\nfimalgoritmo")
		expr := p.Body[0].(*ast.WriteStmt).Args[0].Expr.(*ast.BinaryExpr)
		expr.Left = expr // Corrupted syntax must not bypass the runtime guard.
		i := New(Options{MaxSteps: 10000})
		if ds := i.Run(p, info); len(ds) != 1 || ds[0].Code != diag.RStorage {
			t.Fatalf("missing evaluation guard: %v", ds)
		}
		p, info = analyzed(t, "algoritmo \"ok\"\ninicio\nescreva(1)\nfimalgoritmo")
		if ds := i.Run(p, info); len(ds) != 0 {
			t.Fatalf("failed traversal polluted next run: %v", ds)
		}
	})
}
