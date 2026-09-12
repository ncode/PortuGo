package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func writeFormatted(path string, data []byte) error {
	path, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s: cannot format a non-regular file", path)
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".portugol-fmt-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(file.Name()) }()
	if err := file.Chmod(info.Mode().Perm()); err != nil {
		return errors.Join(err, file.Close())
	}
	_, err = file.Write(data)
	if err := errors.Join(err, file.Close()); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}
