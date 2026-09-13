package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	maxArtifactBytes    = 16 << 20
	maxObservationBytes = 1 << 20
)

var errArtifactTooLarge = errors.New("artifact exceeds size limit")

func safePath(root, name string) (string, error) {
	if name == "." || !fs.ValidPath(name) || strings.ContainsAny(name, "\\:\x00") {
		return "", fmt.Errorf("unsafe path %q", name)
	}
	current := root
	for _, part := range strings.Split(name, "/") {
		if strings.EqualFold(part, ".git") {
			return "", fmt.Errorf("unsafe path %q", name)
		}
		if strings.HasSuffix(part, ".") || strings.HasSuffix(part, " ") {
			return "", fmt.Errorf("unsafe path %q", name)
		}
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
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", name, err)
	}
	defer f.Close()
	b, err := readBounded(f, maxArtifactBytes)
	if errors.Is(err, errArtifactTooLarge) {
		return nil, fmt.Errorf("invalid file size or type: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", name, err)
	}
	return b, nil
}

func readBounded(r io.Reader, limit int) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(r, int64(limit)+1))
	if err != nil {
		return nil, err
	}
	if len(b) > limit {
		return nil, errArtifactTooLarge
	}
	return b, nil
}

func checkAbsent(root, name string) error {
	p, err := safePath(root, name)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(p); !os.IsNotExist(err) {
		return fmt.Errorf("expected absent file: %s", name)
	}
	return nil
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

func readReplayOutputArtifact(root string, a artifact) ([]byte, error) {
	b, err := readArtifact(root, a)
	if err != nil {
		return nil, err
	}
	if len(b) > maxObservationBytes {
		return nil, fmt.Errorf("replay output exceeds size limit: %s", a.Path)
	}
	return b, nil
}

func checkProhibited(root string, r reference) error {
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Clean(filepath.Dir(p)) == filepath.Clean(root) && strings.EqualFold(d.Name(), ".git") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
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
