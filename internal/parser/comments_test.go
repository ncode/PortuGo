package parser

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/lexer"
	"github.com/ncode/portugol-go/internal/source"
	"github.com/ncode/portugol-go/internal/token"
)

func TestCommentAnchorsRoundTrip(t *testing.T) {
	for _, name := range []string{"empty_blocks", "declarations", "subprograms", "comment_forms", "multiline_comments", "default_semicolon_comments"} {
		wantBytes, err := os.ReadFile(filepath.Join("..", "..", "testdata", "format", name+".alg"))
		if err != nil {
			t.Fatal(err)
		}
		want := string(wantBytes)
		for _, ending := range []string{"\n", "\r\n"} {
			input := strings.ReplaceAll(want, "\n", ending)
			var previous *ast.Program
			for pass := range 2 {
				_, tokens, ds := lexer.Scan("comments.alg", input)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				program, ds := Parse(tokens)
				if len(ds) != 0 {
					t.Fatal(ds)
				}
				var output bytes.Buffer
				if err := ast.Fprint(&output, program); err != nil {
					t.Fatal(err)
				}
				if output.String() != want {
					t.Fatalf("pass %d:\ngot:\n%s\nwant:\n%s", pass, &output, want)
				}
				input = output.String()
				program.Fragments = nil // Formatting trivia is checked by the byte golden.
				program.Suffix.Text = strings.ReplaceAll(program.Suffix.Text, "\r\n", "\n")
				clearSyntaxPositions(reflect.ValueOf(program))
				if previous != nil && !reflect.DeepEqual(previous, program) {
					t.Fatalf("pass %d changed syntax or comment contents/order/anchors", pass)
				}
				previous = program
			}
		}
	}
}

func clearSyntaxPositions(v reflect.Value) {
	if !v.IsValid() {
		return
	}
	if v.Type() == reflect.TypeFor[token.Pos]() && v.CanSet() {
		v.SetInt(0)
		return
	}
	if v.Type() == reflect.TypeFor[token.Token]() && v.CanSet() {
		v.FieldByName("Raw").SetString("")
	}
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !v.IsNil() {
			clearSyntaxPositions(v.Elem())
		}
	case reflect.Struct:
		for i := range v.NumField() {
			clearSyntaxPositions(v.Field(i))
		}
	case reflect.Slice:
		for i := range v.Len() {
			clearSyntaxPositions(v.Index(i))
		}
	}
}

func TestCommentSourceSpans(t *testing.T) {
	src := "// leading\r\nalgoritmo \"comments\" // header\r\ninicio\r\n  { body\r\nfimalgoritmo // suffix\r\n"
	_, tokens, ds := lexer.Scan("comments.alg", src)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	program, ds := Parse(tokens)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	want := []struct {
		text   string
		inline bool
	}{
		{"// leading", false}, {"// header", true}, {"{ body", false},
	}
	if len(program.Comments) != len(want) {
		t.Fatalf("comments = %+v", program.Comments)
	}
	for i, comment := range program.Comments {
		if comment.Text != want[i].text || comment.Inline != want[i].inline || comment.Pos != token.Pos(strings.Index(src, want[i].text)) || src[comment.Pos:comment.End] != comment.Text {
			t.Errorf("comment %d lost its span or anchor: %+v", i, comment)
		}
	}
	if program.Suffix.Text != " // suffix\r\n" {
		t.Fatalf("suffix = %q", program.Suffix.Text)
	}
}

func TestCP1252CommentsRoundTrip(t *testing.T) {
	original := []byte("algoritmo \"comments\"\ninicio\n// a\xe7\xe3o\nfimalgoritmo\n")
	input, err := source.DecodeFile("comments.alg", original)
	if err != nil {
		t.Fatal(err)
	}
	_, tokens, ds := lexer.ScanFile(input)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	program, ds := Parse(tokens)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	if len(program.Comments) != 1 || program.Comments[0].Pos != token.Pos(bytes.Index(original, []byte("//"))) || program.Comments[0].End-program.Comments[0].Pos != token.Pos(len("// a\xe7\xe3o")) {
		t.Fatalf("comment lost its original-byte span: %+v", program.Comments)
	}
	var output bytes.Buffer
	if err := ast.Fprint(&output, program); err != nil {
		t.Fatal(err)
	}
	if want := "algoritmo \"comments\"\ninicio\n  // ação\nfimalgoritmo\n"; output.String() != want {
		t.Fatalf("got %q, want %q", output.String(), want)
	}
}
