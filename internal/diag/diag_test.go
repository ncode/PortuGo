package diag

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/token"
)

func TestDiagnosticContract(t *testing.T) {
	cause := errors.New("underlying failure")
	d := Diagnostic{Code: RHost, Pos: 4, End: 7, Message: "write failed", Cause: cause}
	if !errors.Is(d, cause) || !HasErrors([]Diagnostic{d}) {
		t.Fatal("error severity or wrapped cause was lost")
	}
	if HasErrors([]Diagnostic{{Severity: Warning}, {Severity: Note}}) {
		t.Fatal("non-errors prevent execution")
	}
	diags := []Diagnostic{d, {Code: RArithmetic, Pos: 1, Message: "first"}, {Code: RInput, Pos: 1, Message: "second"}}
	var out bytes.Buffer
	Render(&out, token.NewFile("sample.alg", 10), diags)
	want := "sample.alg:1:2: R002: first\nsample.alg:1:2: R004: second\nsample.alg:1:5: R008: write failed\n"
	if out.String() != want || diags[0].Pos != 4 || strings.Contains(out.String(), cause.Error()) {
		t.Fatalf("unstable or leaking diagnostic rendering: %q", &out)
	}
	codes := []Code{RType, RArithmetic, RStorage, RInput, RCall, RLoop, RBuiltin, RHost}
	for n, code := range codes {
		if code != Code("R00"+string(rune('1'+n))) {
			t.Fatalf("runtime code %d = %q", n, code)
		}
	}
}
