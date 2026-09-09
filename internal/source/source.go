package source

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/ncode/portugol-go/internal/cp1252"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

// MaxBytes is the source-file and REPL-submission byte limit before decoding.
const MaxBytes = 4 << 20

// ReadFile reads and decodes a Portugol source file.
func ReadFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	b, readErr := io.ReadAll(io.LimitReader(f, MaxBytes+1))
	if err := errors.Join(readErr, f.Close()); err != nil {
		return "", err
	}
	text, err := Decode(b)
	if err != nil {
		prefix := b[:MaxBytes]
		line := bytes.Count(prefix, []byte{'\n'}) + 1
		column := MaxBytes - bytes.LastIndexByte(prefix, '\n')
		return "", fmt.Errorf("%s:%d:%d: %w", path, line, column, err)
	}
	return text, nil
}

// Decode converts UTF-8 or Windows-1252 bytes into a Go string.
// The original-byte limit is checked before any transcoding allocation.
func Decode(b []byte) (string, error) {
	if len(b) > MaxBytes {
		return "", diag.Diagnostic{Code: diag.EResource, Pos: token.Pos(MaxBytes), End: token.Pos(MaxBytes + 1), Message: "source size limit exceeded"}
	}
	if len(b) >= 3 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		b = b[3:]
	}
	if utf8.Valid(b) {
		return string(b), nil
	}
	size := 0
	for _, c := range b {
		size += utf8.RuneLen(cp1252.DecodeByte(c))
	}
	var text strings.Builder
	text.Grow(size) // At most three decoded bytes per original byte.
	for _, c := range b {
		text.WriteRune(cp1252.DecodeByte(c))
	}
	return text.String(), nil
}
