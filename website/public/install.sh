#!/usr/bin/env bash
set -euo pipefail

repository="omarchy-forge/forge"
install_dir="${OMAFORGE_INSTALL_DIR:-$HOME/.local/bin}"
requested_version="${OMAFORGE_VERSION:-}"

usage() {
  cat <<'EOF'
Usage: install.sh [--version vX.Y.Z]

Install or update Omarchy Forge for the current user. Without --version, the
script resolves the latest published GitHub release. It never uses sudo or a
system package manager.
EOF
}

while (($#)); do
  case "$1" in
    --version)
      [[ $# -ge 2 ]] || { echo "--version requires a value" >&2; exit 2; }
      requested_version="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

case "$(uname -m)" in
  x86_64) arch="amd64" ;;
  aarch64|arm64) arch="arm64" ;;
  *) echo "unsupported Linux architecture: $(uname -m)" >&2; exit 1 ;;
esac

curl_args=(--fail --location --silent --show-error --retry 3)
if [[ -z "$requested_version" ]]; then
  latest_url="$(curl "${curl_args[@]}" --output /dev/null --write-out '%{url_effective}' \
    "https://github.com/${repository}/releases/latest")"
  requested_version="${latest_url##*/}"
fi

if [[ ! "$requested_version" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z.-]+)?$ ]]; then
  echo "latest release did not resolve to an exact version: $requested_version" >&2
  exit 1
fi

destination="$install_dir/omaforge"
if [[ -L "$destination" ]]; then
  echo "refusing to replace symlinked installation target: $destination" >&2
  exit 1
fi
if [[ -e "$destination" && ! -f "$destination" ]]; then
  echo "installation target is not a regular file: $destination" >&2
  exit 1
fi
if [[ -x "$destination" ]]; then
  installed_version="$("$destination" version 2>/dev/null | awk 'NR == 1 && $1 == "omaforge" { print $2 }' || true)"
  if [[ "$installed_version" == "$requested_version" ]]; then
    echo "Omarchy Forge $requested_version is already installed at $destination"
    exit 0
  fi
fi

release_version="${requested_version#v}"
archive="omaforge_${release_version}_linux_${arch}.tar.gz"
release_root="${OMAFORGE_RELEASE_BASE_URL:-https://github.com/${repository}/releases/download}"
base_url="$release_root/$requested_version"
install_tmp="$(mktemp -d)"
stage=""
cleanup() {
  rm -rf -- "$install_tmp"
  [[ -z "$stage" ]] || rm -f -- "$stage"
}
trap cleanup EXIT

curl "${curl_args[@]}" --output "$install_tmp/$archive" "$base_url/$archive"
curl "${curl_args[@]}" --output "$install_tmp/checksums.txt" "$base_url/checksums.txt"

expected="$(awk -v name="$archive" '$2 == name { print $1 }' "$install_tmp/checksums.txt")"
if [[ ! "$expected" =~ ^[0-9a-f]{64}$ ]]; then
  echo "checksums.txt has no valid checksum for $archive" >&2
  exit 1
fi
actual="$(sha256sum "$install_tmp/$archive" | cut -d' ' -f1)"
if [[ "$actual" != "$expected" ]]; then
  echo "checksum verification failed for $archive" >&2
  exit 1
fi

mkdir -p -- "$install_dir"
tar -xzf "$install_tmp/$archive" -C "$install_tmp" omaforge
[[ -f "$install_tmp/omaforge" && ! -L "$install_tmp/omaforge" ]] || {
  echo "release archive does not contain a regular omaforge binary" >&2
  exit 1
}
stage="$(mktemp "$install_dir/.omaforge.install.XXXXXX")"
install -m 0755 "$install_tmp/omaforge" "$stage"
mv -f -- "$stage" "$destination"
stage=""

"$destination" version
echo "Installed Omarchy Forge at $destination"
