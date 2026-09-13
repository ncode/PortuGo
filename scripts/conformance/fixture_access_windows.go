//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

func applyUnavailableFixtureAccess(name string) (func() error, error) {
	path, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	h, err := syscall.CreateFile(path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return nil, err
	}
	restore := func() error {
		if h == syscall.InvalidHandle {
			return nil
		}
		err := syscall.CloseHandle(h)
		h = syscall.InvalidHandle
		return err
	}
	f, err := os.Open(name)
	if err == nil {
		_ = f.Close()
		_ = restore()
		return nil, fmt.Errorf("fixture access restriction was bypassed")
	}
	if !errors.Is(err, syscall.Errno(32)) && !errors.Is(err, os.ErrPermission) {
		_ = restore()
		return nil, fmt.Errorf("fixture access restriction returned %w", err)
	}
	return restore, nil
}
