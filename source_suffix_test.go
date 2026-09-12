package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedIgnoredSuffixes(t *testing.T) {
	for _, name := range []string{"unterminated-next-line", "invalid-symbols", "unclosed-block", "notes"} {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", "source-suffix-"+name)
			src, err := source.ReadFile(filepath.Join(dir, "source.alg"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(dir, "stdout.txt"))
			if err != nil {
				t.Fatal(err)
			}
			checkFormattingPreservesExecution(t, src, want)
		})
	}
}
