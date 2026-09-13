//go:build !windows

package main

import (
	"errors"
	"fmt"
	"os"
)

func applyUnavailableFixtureAccess(name string) (func() error, error) {
	info, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	permissions := info.Mode().Perm()
	if err := os.Chmod(name, 0); err != nil {
		return nil, err
	}
	restore := func() error {
		return os.Chmod(name, permissions)
	}
	f, err := os.Open(name)
	if err == nil {
		_ = f.Close()
		_ = restore()
		return nil, fmt.Errorf("fixture access restriction was bypassed")
	}
	if !errors.Is(err, os.ErrPermission) {
		_ = restore()
		return nil, fmt.Errorf("fixture access restriction returned %w", err)
	}
	return restore, nil
}
