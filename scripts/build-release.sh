#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
commit="${2:-}"
build_date="${3:-}"
dist="${4:-dist}"

if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([+-][0-9A-Za-z.-]+)?$ ]]; then
  echo "usage: scripts/build-release.sh vX.Y.Z COMMIT RFC3339_DATE [DIST]" >&2
  exit 2
fi
if [[ -z "$commit" || -z "$build_date" ]]; then
  echo "commit and build date are required" >&2
  exit 2
fi
if [[ -e "$dist" ]] && [[ -n "$(find "$dist" -mindepth 1 -print -quit 2>/dev/null)" ]]; then
  echo "release output directory must be empty: $dist" >&2
  exit 1
fi

mkdir -p "$dist"
dist="$(cd "$dist" && pwd)"
release_version="${version#v}"
ldflags="-s -w -X main.version=$version -X main.commit=$commit -X main.buildDate=$build_date"

for arch in amd64 arm64; do
  work="$dist/.work-$arch"
  mkdir -p "$work"
  CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -trimpath -buildvcs=false -ldflags "$ldflags" -o "$work/omaforge" ./cmd/omaforge
  cp LICENSE README.md "$work/"
  archive="$dist/omaforge_${release_version}_linux_${arch}.tar.gz"
  tar --sort=name --owner=0 --group=0 --numeric-owner --mtime='UTC 1970-01-01' -C "$work" -cf - LICENSE README.md omaforge | gzip -n > "$archive"
  rm -f "$work/LICENSE" "$work/README.md" "$work/omaforge"
  rmdir "$work"
done

(
  cd "$dist"
  sha256sum omaforge_*.tar.gz > checksums.txt
)
