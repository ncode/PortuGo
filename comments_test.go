package portugol_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/source"
)

func TestRecordedComments(t *testing.T) {
	for _, id := range []string{
		"slash-comment", "brace-comment", "multiline-brace-comment",
		"c-style-comment", "unterminated-comment", "brace-comment-inline",
		"brace-comment-cross-line", "c-comment-inline", "c-comment-cross-line",
		"c-comment-unclosed", "brace-comment-after-statement", "slash-comment-after-statement",
		"single-slash-comment", "single-star-comment", "closing-brace-comment", "closing-c-comment",
		"brace-inline-trailing-statement", "closing-brace-inline", "indented-slash-line", "indented-star-line",
		"string-single-slash", "string-c-opening", "string-c-closing", "string-brace-opening",
		"string-brace-closing", "string-brace-content",
	} {
		t.Run(id, func(t *testing.T) {
			dir := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id)
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

func TestRecordedCommentInsideString(t *testing.T) {
	for _, id := range []string{"comment-symbols-in-string", "string-double-slash", "string-comment-after-content"} {
		t.Run(id, func(t *testing.T) {
			path := filepath.Join("testdata/conformance/visualg-3.0.7/probes", id, "source.alg")
			src, err := source.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			file, _, ds := lexer.Scan("source.alg", src)
			if len(ds) != 1 || ds[0].Code != diag.ELexer || file.Position(ds[0].Pos).Line != 3 {
				t.Fatalf("diagnostics = %v, want L001 on line 3", ds)
			}
		})
	}
}
