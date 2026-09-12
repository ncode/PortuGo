package interp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func TestCorruptedVectorStorage(t *testing.T) {
	for _, body := range []string{"escreval(v[2])", "v[2] <- 5", "leia(v[2])", "p(v[2])"} {
		t.Run(body, func(t *testing.T) {
			src := `algoritmo "storage"
var v: vetor[1..2] de inteiro
procedimento p(var n: inteiro)
inicio
  n <- 5
fimprocedimento
inicio
  v[1] <- 7
  v[2] <- 9
  escreval("BEFORE")
  ` + body + `
  escreval("AFTER")
fimalgoritmo`
			prog, info := analyzed(t, src)
			binding, ok := info.Binding(prog.Globals[0].Names[0])
			if !ok {
				t.Fatal("missing test binding")
			}
			var i *Interpreter
			var out bytes.Buffer
			var original []runtime.Cell
			writer := storageWriter(func(p []byte) (int, error) {
				if original == nil {
					cell, ok := i.global.lookup(binding.ID)
					if !ok {
						t.Fatal("missing test storage")
					}
					original = cell.Value.Vec.Elements
					cell.Value.Vec.Elements = original[:1]
				}
				return out.Write(p)
			})
			i = New(Options{Output: writer, Input: strings.NewReader("13\n")})
			ds := i.Run(prog, info)
			if len(ds) != 1 || ds[0].Code != diag.RStorage || ds[0].Pos != token.Pos(strings.LastIndex(src, "v[2]")) {
				t.Fatalf("diagnostics = %v, want R003 at the failed indexing expression", ds)
			}
			if out.String() != "BEFORE\n" || original[0].Value.Int != 7 || original[1].Value.Int != 9 {
				t.Fatalf("failed access changed output or elements: %q, %v", &out, original)
			}
			if input, err := i.in.Peek(2); err != nil || string(input) != "13" {
				t.Fatalf("failed access consumed input: %q, %v", input, err)
			}
		})
	}
}

type storageWriter func([]byte) (int, error)

func (w storageWriter) Write(p []byte) (int, error) { return w(p) }
