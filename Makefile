GO ?= go
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || printf 'dev-%s' "$$(git rev-parse --short HEAD)")
DIST_DIR ?= dist
DIST_TARGETS ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

export GO VERSION DIST_DIR DIST_TARGETS

.PHONY: dist

# Build versioned archives and SHA256SUMS, following ncode/facts.
dist:
	bash scripts/dist.sh
