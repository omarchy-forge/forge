#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
test_root="$(mktemp -d /tmp/omaforge-installer-test.XXXXXX)"
trap 'rm -rf -- "$test_root"' EXIT

build_release() {
  local version="$1"
  local output="$test_root/releases/$version"
  mkdir -p "$output"
  "$repo_root/scripts/build-release.sh" "$version" "commit-${version#v}" "2026-08-23T12:00:00Z" "$output"
}

build_release v0.1.0
build_release v0.2.0
build_release v0.3.0

install_dir="$test_root/bin"
run_installer() {
  OMAFORGE_VERSION="$1" \
    OMAFORGE_RELEASE_BASE_URL="file://$test_root/releases" \
    OMAFORGE_INSTALL_DIR="$install_dir" \
    bash "$repo_root/website/public/install.sh"
}

first_output="$(run_installer v0.1.0)"
grep -Fq "omaforge v0.1.0" <<<"$first_output"
grep -Fq "Installed Omarchy Forge" <<<"$first_output"
test -x "$install_dir/omaforge"

current_output="$(run_installer v0.1.0)"
grep -Fq "already installed" <<<"$current_output"

update_output="$(run_installer v0.2.0)"
grep -Fq "omaforge v0.2.0" <<<"$update_output"
grep -Fq "Installed Omarchy Forge" <<<"$update_output"
test "$("$install_dir/omaforge" version | sed -n '1s/^omaforge //p')" = "v0.2.0"

tampered="$test_root/tampered/v0.3.0"
mkdir -p "$tampered"
cp "$test_root/releases/v0.3.0"/* "$tampered/"
printf 'tampered\n' >> "$tampered/omaforge_0.3.0_linux_amd64.tar.gz"
if OMAFORGE_VERSION=v0.3.0 OMAFORGE_RELEASE_BASE_URL="file://$test_root/tampered" \
  OMAFORGE_INSTALL_DIR="$install_dir" \
  bash "$repo_root/website/public/install.sh" >/dev/null 2>&1; then
  echo "installer accepted a tampered archive" >&2
  exit 1
fi
test "$("$install_dir/omaforge" version | sed -n '1s/^omaforge //p')" = "v0.2.0"

if OMAFORGE_VERSION=latest OMAFORGE_INSTALL_DIR="$test_root/invalid-bin" \
  bash "$repo_root/website/public/install.sh" >/dev/null 2>&1; then
  echo "installer accepted a non-version release" >&2
  exit 1
fi

test ! -e "$test_root/invalid-bin/omaforge"

symlink_dir="$test_root/symlink-bin"
mkdir -p "$symlink_dir"
ln -s "$install_dir/omaforge" "$symlink_dir/omaforge"
if OMAFORGE_VERSION=v0.3.0 OMAFORGE_RELEASE_BASE_URL="file://$test_root/releases" \
  OMAFORGE_INSTALL_DIR="$symlink_dir" \
  bash "$repo_root/website/public/install.sh" >/dev/null 2>&1; then
  echo "installer replaced a symlinked destination" >&2
  exit 1
fi
test -L "$symlink_dir/omaforge"
echo "install script checks passed"
