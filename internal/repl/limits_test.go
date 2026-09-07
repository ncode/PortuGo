package repl

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
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
			var d diag.Diagnostic
			if ok || !errors.As(err, &d) || d.Code != diag.EResource || !strings.Contains(err.Error(), "<repl>:") {
				t.Fatalf("unbounded submission: %v", err)
			}
		}
	}
}
