package source

import (
	"bytes"
	"slices"
	"testing"

	"github.com/ncode/portugol-go/internal/token"
)

func TestOriginalPos(t *testing.T) {
	for _, tt := range []struct {
		name, original, text string
		offsets              []int
	}{
		{"empty", "", "", []int{0}},
		{"UTF8", "é€", "é€", []int{0, 1, 2, 3, 4, 5}},
		{"BOM", "\ufeffé€", "é€", []int{3, 4, 5, 6, 7, 8}},
		{"BOM only", "\ufeff", "", []int{3}},
		{"CP1252", "\xe9\x80", "é€", []int{0, 0, 1, 1, 1, 2}},
		{"BOM CP1252", "\ufeff\xe9\x80", "é€", []int{3, 3, 4, 4, 4, 5}},
		{"undefined CP1252", "\x81\x8d", "��", []int{0, 0, 0, 1, 1, 1, 2}},
		{"mixed bytes use CP1252", "\xc3\xa9\x80", "Ã©€", []int{0, 0, 1, 1, 2, 2, 2, 3}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := []byte(tt.original)
			f, err := DecodeFile("positions.alg", input)
			if err != nil {
				t.Fatal(err)
			}
			if f.Text != tt.text || !bytes.Equal(f.Original, input) || f.Positions.Size != len(input) {
				t.Fatalf("decoded text or original bytes changed: %+v", f)
			}
			var positions []int
			for i := 0; i <= len(f.Text); i++ {
				positions = append(positions, int(f.OriginalPos(token.Pos(i))))
			}
			if !slices.Equal(positions, tt.offsets) {
				t.Fatalf("positions = %v, want %v", positions, tt.offsets)
			}
			if f.OriginalPos(-1) != token.Pos(tt.offsets[0]) || f.OriginalPos(token.Pos(len(f.Text)+1)) != token.Pos(len(input)) {
				t.Fatal("out-of-range positions were not clamped")
			}
			clear(input)
			if string(f.Original) != tt.original {
				t.Fatal("decoder retained mutable caller bytes")
			}
		})
	}
}

func TestOriginalLinePositions(t *testing.T) {
	f, err := DecodeFile("positions.alg", []byte("\ufeff\xe9\x80\r\n\xe9\n"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ decoded, original, line, column int }{
		{0, 3, 1, 4}, {5, 5, 1, 6}, {7, 7, 2, 1}, {9, 8, 2, 2}, {10, 9, 3, 1},
	} {
		pos := f.Positions.Position(f.OriginalPos(token.Pos(tt.decoded)))
		if pos.Offset != tt.original || pos.Line != tt.line || pos.Column != tt.column {
			t.Errorf("decoded position %d = %+v; want %d at %d:%d", tt.decoded, pos, tt.original, tt.line, tt.column)
		}
	}
}
