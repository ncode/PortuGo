package portugol_test

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func runDist(t *testing.T, dir, version, targets string) ([]byte, error) {
	t.Helper()
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("distribution packaging runs on Linux and macOS")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	t.Cleanup(cancel)
	cmd := exec.CommandContext(ctx, "make", "dist", "DIST_DIR="+dir, "VERSION="+version, "DIST_TARGETS="+targets)
	return cmd.CombinedOutput()
}

func TestReleaseArchives(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "release files")
	version := "v9.8.7-rc.1"
	native := runtime.GOOS + "/" + runtime.GOARCH
	if output, err := runDist(t, dir, version, native+" windows/amd64"); err != nil {
		t.Fatalf("build distribution: %v\n%s", err, output)
	}
	sums, err := os.ReadFile(filepath.Join(dir, "SHA256SUMS"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(sums)), "\n")
	if len(lines) != 2 {
		t.Fatalf("checksum entries = %d, want 2: %s", len(lines), sums)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 3 {
		t.Fatalf("distribution inventory: %v (%v)", entries, err)
	}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 || filepath.Base(fields[1]) != fields[1] {
			t.Fatalf("invalid checksum line: %q", line)
		}
		data, err := os.ReadFile(filepath.Join(dir, fields[1]))
		if err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != fields[0] {
			t.Fatalf("checksum for %s = %s, want %s", fields[1], got, fields[0])
		}
	}
	name := "portugo-" + version + "-" + runtime.GOOS + "-" + runtime.GOARCH
	file, err := os.Open(filepath.Join(dir, name+".tar.gz"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := file.Close(); err != nil {
			t.Error(err)
		}
	})
	gz, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := gz.Close(); err != nil {
			t.Error(err)
		}
	})
	archive := tar.NewReader(gz)
	var binary []byte
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if header.Typeflag == tar.TypeDir && header.Name == name+"/" {
			continue
		}
		if header.Name != name+"/portugo" || header.Typeflag != tar.TypeReg || header.Mode&0111 == 0 || binary != nil {
			t.Fatalf("unexpected archive entry: %+v", header)
		}
		binary, err = io.ReadAll(archive)
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(binary) == 0 {
		t.Fatal("native archive has no executable")
	}
	executable := filepath.Join(t.TempDir(), "portugo")
	if err := os.WriteFile(executable, binary, 0700); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		args []string
		want string
	}{
		{[]string{"--version"}, version + "\n"},
		{[]string{"run", "examples/hello.alg"}, "Ola, Portugol!\n"},
	} {
		cmd := exec.CommandContext(t.Context(), executable, tt.args...)
		output, err := cmd.CombinedOutput()
		if err != nil || string(output) != tt.want {
			t.Fatalf("%v: error=%v output=%q, want %q", tt.args, err, output, tt.want)
		}
	}
	winName := "portugo-" + version + "-windows-amd64"
	win, err := zip.OpenReader(filepath.Join(dir, winName+".zip"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := win.Close(); err != nil {
			t.Error(err)
		}
	})
	var windowsBinary bool
	for _, entry := range win.File {
		if entry.FileInfo().IsDir() && entry.Name == winName+"/" {
			continue
		}
		if entry.Name != winName+"/portugo.exe" || windowsBinary {
			t.Fatalf("unexpected Windows archive entry: %s", entry.Name)
		}
		reader, err := entry.Open()
		if err != nil {
			t.Fatal(err)
		}
		magic := make([]byte, 2)
		_, err = io.ReadFull(reader, magic)
		closeErr := reader.Close()
		if err != nil || closeErr != nil || !bytes.Equal(magic, []byte("MZ")) {
			t.Fatalf("invalid Windows executable: magic=%q, errors=%v/%v", magic, err, closeErr)
		}
		windowsBinary = true
	}
	if !windowsBinary {
		t.Fatal("Windows archive has no executable")
	}
}

func TestReleaseRejectsInvalidInputs(t *testing.T) {
	for _, tt := range []struct{ name, version, targets, diagnostic string }{
		{"version path", "../escape", "linux/amd64", "invalid version"},
		{"version whitespace", "v1.2.3 extra", "linux/amd64", "invalid version"},
		{"target path", "v1.2.3", "../../escape", "invalid target"},
		{"empty targets", "v1.2.3", "", "no distribution targets"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "dist")
			output, err := runDist(t, dir, tt.version, tt.targets)
			if err == nil || !strings.Contains(string(output), tt.diagnostic) {
				t.Fatalf("error=%v output=%s; want %s", err, output, tt.diagnostic)
			}
			if _, err := os.Stat(dir); !os.IsNotExist(err) {
				t.Fatalf("invalid input created output directory: %v", err)
			}
		})
	}
}

func TestReleaseFailedBuildPreservesDistribution(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "SHA256SUMS")
	want := []byte("previous distribution\n")
	if err := os.WriteFile(path, want, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GO", "false")
	if output, err := runDist(t, dir, "v1.2.3", "linux/amd64"); err == nil {
		t.Fatalf("build unexpectedly succeeded: %s", output)
	}
	got, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("previous checksums changed: %q (%v)", got, err)
	}
}
