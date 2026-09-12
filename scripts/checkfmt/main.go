// Command checkfmt verifies gofmt without rewriting source files.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func check(root string) error {
	var unformatted []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && (strings.HasPrefix(entry.Name(), ".") || entry.Name() == "vendor") {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		formatted, err := format.Source(src)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if !bytes.Equal(src, formatted) {
			unformatted = append(unformatted, filepath.ToSlash(path))
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(unformatted) > 0 {
		return fmt.Errorf("gofmt required:\n%s", strings.Join(unformatted, "\n"))
	}
	return nil
}

func main() {
	if err := check("."); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
