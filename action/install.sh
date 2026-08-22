#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
repository="${2:-omarchy-forge/forge}"
token="${OMAFORGE_GITHUB_TOKEN:-}"

if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([+-][0-9A-Za-z.-]+)?$ ]]; then
  echo "Forge Action requires an exact version such as v1.2.3" >&2
  exit 2
fi
case "$(uname -m)" in
  x86_64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo "unsupported Linux architecture: $(uname -m)" >&2; exit 1 ;;
esac

release_version="${version#v}"
archive="omaforge_${release_version}_linux_${arch}.tar.gz"
base_url="${OMAFORGE_RELEASE_BASE_URL:-https://github.com/${repository}/releases/download}/${version}"
install_root="${RUNNER_TEMP:?RUNNER_TEMP is required}/omaforge-action/${version}"
download_root="$install_root/download"
bin_root="$install_root/bin"
mkdir -p "$download_root" "$bin_root"

curl_args=(--fail --location --silent --show-error --retry 3)
if [[ -n "$token" ]]; then curl_args+=(--header "Authorization: Bearer $token"); fi
curl "${curl_args[@]}" --output "$download_root/$archive" "$base_url/$archive"
curl "${curl_args[@]}" --output "$download_root/checksums.txt" "$base_url/checksums.txt"

expected="$(awk -v name="$archive" '$2 == name { print $1 }' "$download_root/checksums.txt")"
if [[ ! "$expected" =~ ^[0-9a-f]{64}$ ]]; then
  echo "checksums.txt has no valid checksum for $archive" >&2
  exit 1
fi
actual="$(sha256sum "$download_root/$archive" | cut -d' ' -f1)"
if [[ "$actual" != "$expected" ]]; then
  echo "checksum verification failed for $archive" >&2
  exit 1
fi

tar -xzf "$download_root/$archive" -C "$bin_root" omaforge
chmod 0755 "$bin_root/omaforge"
printf '%s\n' "$bin_root" >> "${GITHUB_PATH:?GITHUB_PATH is required}"
printf 'binary=%s/omaforge\n' "$bin_root" >> "${GITHUB_OUTPUT:?GITHUB_OUTPUT is required}"
