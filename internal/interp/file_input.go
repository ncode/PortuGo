package interp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/cp1252"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/token"
)

type fileInputState struct {
	file   *os.File
	reader *bufio.Reader // A nil reader records console input into a new file.
	writer *bufio.Writer
	at     token.Pos
}

func (i *Interpreter) configureFile(s *ast.FileInputStmt) error {
	if err := i.closeInputFile(false); err != nil {
		return err
	}
	name := filepath.FromSlash(strings.ReplaceAll(s.Path, "\\", "/"))
	if name == "" || strings.ContainsRune(name, 0) || filepath.Separator != '\\' && strings.ContainsRune(name, ':') {
		return diag.Diagnostic{Code: diag.RHost, Pos: s.At, Message: "invalid input filename"}
	}
	if !filepath.IsAbs(name) {
		name = filepath.Join(i.options.WorkingDir, name)
	}
	info, statErr := os.Stat(name)
	if statErr == nil && !info.Mode().IsRegular() {
		return diag.Diagnostic{Code: diag.RHost, Pos: s.At, Message: "input file is not a regular file"}
	}
	f, err := os.Open(name)
	if statErr == nil && unavailableFileRead(err) {
		return nil // Existing unreadable files leave console or random input active.
	}
	recording := errors.Is(err, os.ErrNotExist)
	if recording {
		// Exclusive creation cannot overwrite a file appearing after the read attempt.
		f, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if errors.Is(err, os.ErrNotExist) {
			return nil // A missing parent leaves console or random input active.
		}
	}
	if err != nil {
		return diag.Diagnostic{Code: diag.RHost, Pos: s.At, Message: "cannot open input file", Cause: err}
	}
	i.fileInput = fileInputState{file: f, at: s.At}
	info, err = f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		closeErr := i.closeInputFile(false)
		return diag.Diagnostic{Code: diag.RHost, Pos: s.At, Message: "input file is not a regular file", Cause: errors.Join(err, closeErr)}
	}
	if recording {
		i.fileInput.writer = bufio.NewWriterSize(f, 128)
	} else {
		i.fileInput.reader = bufio.NewReader(f)
	}
	return nil
}

func unavailableFileRead(err error) bool {
	// Windows sharing violations use error 32; that number has other meanings elsewhere.
	return errors.Is(err, os.ErrPermission) || runtime.GOOS == "windows" && errors.Is(err, syscall.Errno(32))
}

func (i *Interpreter) closeInputFile(commit bool) error {
	state := i.fileInput
	i.fileInput = fileInputState{}
	if state.file != nil {
		var err error
		if commit && state.writer != nil {
			err = state.writer.Flush()
		}
		err = errors.Join(err, state.file.Close())
		if err != nil {
			return diag.Diagnostic{Code: diag.RHost, Pos: state.at, Message: "cannot finalize input file", Cause: err}
		}
	}
	return nil
}

func (i *Interpreter) readFileLine() (string, error) {
	state := &i.fileInput
	var text strings.Builder
	for count := 0; ; count++ {
		b, err := state.reader.ReadByte()
		if err == io.EOF || err == nil && b == '\n' {
			break
		}
		if err != nil {
			return "", diag.Diagnostic{Code: diag.RHost, Pos: state.at, Message: "cannot read input file", Cause: err}
		}
		r := cp1252.DecodeByte(b)
		if r == utf8.RuneError {
			return "", diag.Diagnostic{Code: diag.RHost, Pos: state.at, Message: "input file contains an undefined Windows-1252 byte"}
		}
		if count >= maxTextBytes || text.Len() > maxTextBytes-utf8.RuneLen(r) {
			return "", failure(state.at, diag.RStorage, fmt.Errorf("file input line size limit exceeded"))
		}
		if b != '\r' {
			text.WriteRune(r)
		}
	}
	// Even an initially empty file supplies one empty line before fallback.
	if _, err := state.reader.Peek(1); err == io.EOF {
		if err := i.closeInputFile(false); err != nil {
			return "", err
		}
	}
	return text.String(), nil
}

func (i *Interpreter) recordInput(text string) error {
	state := &i.fileInput
	if state.writer == nil {
		return nil
	}
	data := make([]byte, 0, len(text)+2)
	for _, r := range text {
		b, ok := cp1252.EncodeRune(r)
		if !ok {
			return diag.Diagnostic{Code: diag.RHost, Pos: state.at, Message: "input cannot be recorded as Windows-1252"}
		}
		data = append(data, b)
	}
	data = append(data, '\r', '\n')
	// Byte writes preserve the recorded 128-byte flush boundaries. Write(data)
	// could bypass the buffer for a long line and expose an unrecorded suffix.
	for _, b := range data {
		if err := state.writer.WriteByte(b); err != nil {
			return diag.Diagnostic{Code: diag.RHost, Pos: state.at, Message: "cannot record input file", Cause: err}
		}
		if state.writer.Available() == 0 {
			if err := state.writer.Flush(); err != nil {
				return diag.Diagnostic{Code: diag.RHost, Pos: state.at, Message: "cannot record input file", Cause: err}
			}
		}
	}
	return nil
}
