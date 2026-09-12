package source

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/ncode/portugol-go/internal/cp1252"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

// MaxBytes is the source-file and REPL-submission byte limit before decoding.
const MaxBytes = 4 << 20

// File retains decoded text and its original source bytes and positions.
type File struct {
	Text      string
	Original  []byte
	Positions *token.File
	prefix    int
	offsets   []int // Decoded offsets at each Windows-1252 byte boundary.
}

// OriginalPos maps a decoded byte offset to its original source byte offset.
// Positions within an expanded character refer to that character's first byte.
func (f *File) OriginalPos(pos token.Pos) token.Pos {
	offset := max(0, min(int(pos), len(f.Text)))
	if len(f.offsets) != 0 {
		offset = sort.Search(len(f.offsets), func(i int) bool { return f.offsets[i] > offset }) - 1
	}
	return token.Pos(f.prefix + offset)
}

// LoadFile reads a bounded source file, retaining its bytes and position mapping.
func LoadFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, readErr := io.ReadAll(io.LimitReader(f, MaxBytes+1))
	if err := errors.Join(readErr, f.Close()); err != nil {
		return nil, err
	}
	decoded, err := DecodeFile(path, b)
	if err != nil {
		prefix := b[:MaxBytes]
		line := bytes.Count(prefix, []byte{'\n'}) + 1
		column := MaxBytes - bytes.LastIndexByte(prefix, '\n')
		return nil, fmt.Errorf("%s:%d:%d: %w", path, line, column, err)
	}
	return decoded, nil
}

// ReadFile reads decoded text. Use LoadFile when original positions are needed.
func ReadFile(path string) (string, error) {
	f, err := LoadFile(path)
	if err != nil {
		return "", err
	}
	return f.Text, nil
}

// Decode converts source bytes to text. Use DecodeFile to retain source positions.
func Decode(b []byte) (string, error) {
	f, err := DecodeFile("", b)
	if err != nil {
		return "", err
	}
	return f.Text, nil
}

// DecodeFile converts UTF-8 or Windows-1252 source and retains original positions.
// The original-byte limit is checked before any transcoding allocation.
// Undefined Windows-1252 bytes become U+FFFD, retaining their original byte.
func DecodeFile(name string, b []byte) (*File, error) {
	if len(b) > MaxBytes {
		return nil, diag.Diagnostic{Code: diag.EResource, Pos: token.Pos(MaxBytes), End: token.Pos(MaxBytes + 1), Message: "source size limit exceeded"}
	}
	f := &File{Original: bytes.Clone(b), Positions: token.NewFile(name, len(b))}
	for i, c := range b {
		if c == '\n' {
			f.Positions.AddLine(i + 1)
		}
	}
	if len(b) >= 3 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		b = b[3:]
		f.prefix = 3
	}
	if utf8.Valid(b) {
		f.Text = string(b)
		return f, nil
	}
	size := 0
	for _, c := range b {
		size += utf8.RuneLen(cp1252.DecodeByte(c))
	}
	var text strings.Builder
	text.Grow(size) // At most three decoded bytes per original byte.
	f.offsets = make([]int, len(b)+1)
	for i, c := range b {
		f.offsets[i] = text.Len()
		text.WriteRune(cp1252.DecodeByte(c))
	}
	f.offsets[len(b)] = text.Len()
	f.Text = text.String()
	return f, nil
}
