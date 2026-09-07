package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const maxArtifactBytes = 16 << 20

func safePath(root, name string) (string, error) {
	if name == "." || !fs.ValidPath(name) || strings.ContainsAny(name, "\\:\x00") {
		return "", fmt.Errorf("unsafe path %q", name)
	}
	current := root
	for _, part := range strings.Split(name, "/") {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("inspect %s: %w", name, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("symlink path %s", name)
		}
	}
	return current, nil
}

func readFile(root, name string) ([]byte, error) {
	p, err := safePath(root, name)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(p)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", name, err)
	}
	if !info.Mode().IsRegular() || info.Size() > maxArtifactBytes {
		return nil, fmt.Errorf("invalid file size or type: %s", name)
	}
	return os.ReadFile(p)
}

func validHash(s string) bool {
	b, err := hex.DecodeString(s)
	return err == nil && len(b) == sha256.Size && s == strings.ToLower(s)
}

func hashBytes(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }

func readArtifact(root string, a artifact) ([]byte, error) {
	if !validHash(a.SHA256) {
		return nil, fmt.Errorf("invalid artifact hash: %s", a.Path)
	}
	b, err := readFile(root, a.Path)
	if err != nil {
		return nil, err
	}
	if hashBytes(b) != a.SHA256 {
		return nil, fmt.Errorf("hash mismatch: %s", a.Path)
	}
	return b, nil
}

func checkProhibited(root string, r reference) error {
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		name, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		name = filepath.ToSlash(name)
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink path %s", name)
		}
		switch strings.ToLower(path.Ext(name)) {
		case ".exe", ".dll", ".msi", ".rar", ".zip", ".7z", ".cab", ".iso", ".com":
			return fmt.Errorf("prohibited artifact: %s", name)
		}
		b, err := readFile(root, name)
		if err != nil {
			return err
		}
		hash := hashBytes(b)
		if hash == r.ArchiveSHA256 || hash == r.ExecutableSHA256 {
			return fmt.Errorf("prohibited artifact: %s", name)
		}
		return nil
	})
}
