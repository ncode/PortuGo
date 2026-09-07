#!/usr/bin/env bash
set -euo pipefail

# Group 2 introduces both paths together. A partial installation fails closed.
manifest=testdata/conformance/visualg-3.0.7/manifest.json
if [[ -e "$manifest" || -d scripts/conformance ]]; then
  test -f "$manifest"
  test -d scripts/conformance
  go run ./scripts/conformance validate --mode evidence
else
  printf '%s\n' 'Oracle gate pending group 2: the complete corpus and conformance claim remain unfinished.'
fi
