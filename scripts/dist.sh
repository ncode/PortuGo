#!/usr/bin/env bash
set -euo pipefail

# Values are exported by Makefile, including command-line overrides.
if [[ ! "$VERSION" =~ ^[A-Za-z0-9][A-Za-z0-9._+-]*$ ]]; then
  printf 'invalid version: %s\n' "$VERSION" >&2
  exit 1
fi
read -r -a targets <<< "$DIST_TARGETS"
if [[ ${#targets[@]} -eq 0 ]]; then
  printf 'no distribution targets\n' >&2
  exit 1
fi
for target in "${targets[@]}"; do
  case "$target" in
    linux/amd64|linux/arm64|darwin/amd64|darwin/arm64|windows/amd64|windows/arm64) ;;
    *) printf 'invalid target: %s\n' "$target" >&2; exit 1 ;;
  esac
done

dist_work=$(mktemp -d)
trap 'rm -rf "$dist_work"' EXIT
archives=()
for target in "${targets[@]}"; do
  goos=${target%/*}
  goarch=${target#*/}
  name="portugo-$VERSION-$goos-$goarch"
  binary=portugo
  if [[ "$goos" == windows ]]; then binary=portugo.exe; fi
  mkdir -p "$dist_work/$name"
  printf 'building %s\n' "$name"
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" "$GO" build -trimpath \
    -ldflags "-X main.version=$VERSION" -o "$dist_work/$name/$binary" ./cmd/portugo
  if [[ "$goos" == windows ]]; then
    archive="$name.zip"
    (cd "$dist_work" && zip -q -r "$archive" "$name")
  else
    archive="$name.tar.gz"
    COPYFILE_DISABLE=1 tar -czf "$dist_work/$archive" -C "$dist_work" "$name"
  fi
  archives+=("$archive")
done

if command -v sha256sum >/dev/null 2>&1; then
  checksum=(sha256sum)
else
  checksum=(shasum -a 256)
fi
(cd "$dist_work" && "${checksum[@]}" "${archives[@]}" > SHA256SUMS)
mkdir -p "$DIST_DIR"
for archive in "${archives[@]}"; do
  mv -f "$dist_work/$archive" "$DIST_DIR/$archive"
done
mv -f "$dist_work/SHA256SUMS" "$DIST_DIR/SHA256SUMS"
cat "$DIST_DIR/SHA256SUMS"
