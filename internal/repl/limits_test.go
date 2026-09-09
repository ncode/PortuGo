package repl

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/interp"
	"github.com/ncode/portugol-go/internal/source"
)

func TestSubmissionLimit(t *testing.T) {
	prefix, suffix := "algoritmo \"limit\"\ninicio\n//", "\nfimalgoritmo\n"
	for _, extra := range []int{0, 1} {
		text := prefix + strings.Repeat("x", source.MaxBytes-len(prefix)-len(suffix)+extra) + suffix
		var stderr bytes.Buffer
		ok, err := Run(interp.Options{Input: strings.NewReader(text + "\n:sair\n")}, &stderr)
		if extra == 0 {
			if !ok || err != nil || stderr.Len() != 0 {
				t.Fatalf("boundary rejected: %v, %q", err, &stderr)
			}
		} else {
			if ok || err != nil || strings.Count(stderr.String(), ": E900:") != 1 || !strings.Contains(stderr.String(), "<repl>:") {
				t.Fatalf("unbounded submission: %v, %q", err, &stderr)
			}
		}
	}
}
