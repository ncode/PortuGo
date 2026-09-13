package main

import (
	"fmt"
	"os"
)

func prepareFixtureAccess(root string, access *fixtureAccess) (func() error, error) {
	if access == nil {
		return func() error { return nil }, nil
	}
	if access.Mode != "unavailable" {
		return nil, fmt.Errorf("unsupported fixture access mode %q", access.Mode)
	}
	name, err := safePath(root, access.Path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(name)
	if err != nil {
		return nil, fmt.Errorf("fixture access path %s: %w", access.Path, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("fixture access path %s is not a regular file", access.Path)
	}
	return applyUnavailableFixtureAccess(name)
}
