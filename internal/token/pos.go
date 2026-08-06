package token

import "sort"

// Pos is a byte offset in a decoded source file.
type Pos int

// NoPos is the zero value for an unknown source position.
const NoPos Pos = 0

// Position is a rendered file position.
type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

// File tracks line offsets for one source file.
type File struct {
	Name  string
	Size  int
	lines []int
}

// NewFile creates a source position table for a file.
func NewFile(name string, size int) *File {
	return &File{Name: name, Size: size, lines: []int{0}}
}

// AddLine records the offset of a new line.
func (f *File) AddLine(offset int) {
	if offset <= 0 || offset > f.Size {
		return
	}
	if last := f.lines[len(f.lines)-1]; offset <= last {
		return
	}
	f.lines = append(f.lines, offset)
}

// Position returns the line and column for pos.
func (f *File) Position(pos Pos) Position {
	offset := int(pos)
	if offset < 0 {
		offset = 0
	}
	if offset > f.Size {
		offset = f.Size
	}
	line := sort.Search(len(f.lines), func(i int) bool { return f.lines[i] > offset })
	if line == 0 {
		line = 1
	}
	lineStart := f.lines[line-1]
	return Position{
		Filename: f.Name,
		Offset:   offset,
		Line:     line,
		Column:   offset - lineStart + 1,
	}
}
