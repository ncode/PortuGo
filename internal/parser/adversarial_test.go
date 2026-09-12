package parser

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/testprocess"
)

func TestTruncatedEncodedPrograms(t *testing.T) {
	program := "algoritmo \"ação\"\nvar\nprocedimento P(n: inteiro)\ninicio\nescreval(n)\nfimprocedimento\ninicio\nse verdadeiro entao\nP(1 + // ação\n2)\nfimse\nfimalgoritmo\n\"ignored"
	for _, ending := range []string{"\n", "\r\n", "\r\r\n"} {
		for _, cp1252 := range []bool{false, true} {
			name := strings.ReplaceAll(ending, "\r", "CR")
			name = strings.ReplaceAll(name, "\n", "LF")
			text := strings.ReplaceAll(program, "\n", ending)
			if cp1252 {
				name += "-CP1252"
				text = strings.ReplaceAll(text, "ação", "a\xe7\xe3o")
			}
			t.Run(name, func(t *testing.T) {
				testprocess.Run(t, func() {
					data := []byte("\xef\xbb\xbf" + text)
					for end := 0; end <= len(data); end++ {
						decoded, err := source.DecodeFile("truncated.alg", data[:end])
						if err != nil {
							t.Fatal(err)
						}
						_, tokens, lexDiags := lexer.ScanFile(decoded)
						_, parseDiags := Parse(tokens)
						checkDiagnosticPositions(t, end, lexDiags, parseDiags)
						if end == len(data) && len(lexDiags)+len(parseDiags) != 0 {
							t.Fatal("complete fixture was rejected")
						}
					}
				})
			})
		}
	}
}

func TestEncodedSourceLimitSubprocess(t *testing.T) {
	for _, extra := range []int{0, 1} {
		name := "boundary"
		timeout := 30 * time.Second // Allow race-instrumented decoding of 4 MiB under load.
		if extra != 0 {
			name = "excessive"
			timeout = 5 * time.Second
		}
		t.Run(name, func(t *testing.T) {
			testprocess.RunWithTimeout(t, timeout, func() {
				prefix := []byte("\xef\xbb\xbfalgoritmo \"limit\"\r\r\ninicio\r\r\nfimalgoritmo\r\r\n")
				data := append(prefix, bytes.Repeat([]byte{0x80}, source.MaxBytes+extra-len(prefix))...)
				decoded, err := source.DecodeFile("limit.alg", data)
				if extra != 0 {
					var d diag.Diagnostic
					if decoded != nil || !errors.As(err, &d) || d.Code != diag.EResource || int(d.Pos) != source.MaxBytes || int(d.End) != source.MaxBytes+1 {
						t.Fatal("missing positioned source limit rejection")
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				_, tokens, lexDiags := lexer.ScanFile(decoded)
				program, parseDiags := Parse(tokens)
				checkDiagnosticPositions(t, len(data), lexDiags, parseDiags)
				if program == nil || len(lexDiags)+len(parseDiags) != 0 {
					t.Fatal("encoded source boundary was rejected")
				}
			})
		})
	}
}
