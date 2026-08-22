#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
test_root="$(mktemp -d /tmp/omaforge-release-test.XXXXXX)"
trap 'rm -f "$test_root/github-path" "$test_root/github-output"; rm -rf "$test_root"' EXIT

version=v0.1.0
commit=0123456789abcdef
build_date=2026-08-22T12:00:00Z
dist="$test_root/releases/$version"
second_dist="$test_root/second/$version"
mkdir -p "$dist"
mkdir -p "$second_dist"

cd "$repo_root"
scripts/build-release.sh "$version" "$commit" "$build_date" "$dist"
scripts/build-release.sh "$version" "$commit" "$build_date" "$second_dist"
(
  cd "$dist"
  sha256sum --check checksums.txt
)
test "$(find "$dist" -maxdepth 1 -name 'omaforge_*.tar.gz' | wc -l)" -eq 2
cmp "$dist/checksums.txt" "$second_dist/checksums.txt"

export RUNNER_TEMP="$test_root/runner"
export GITHUB_PATH="$test_root/github-path"
export GITHUB_OUTPUT="$test_root/github-output"
export OMAFORGE_RELEASE_BASE_URL="file://$test_root/releases"
mkdir -p "$RUNNER_TEMP"
action/install.sh "$version" example/forge
binary="$(sed -n 's/^binary=//p' "$GITHUB_OUTPUT")"
test -x "$binary"
version_output="$("$binary" version)"
grep -Fx "omaforge $version" <<<"$version_output"
grep -Fx "commit: $commit" <<<"$version_output"
grep -Fx "built: $build_date" <<<"$version_output"

sample="$test_root/sample"
"$binary" init "$sample" --non-interactive --id dev.example.release-test --author Example >/dev/null
report="$test_root/report.json"
set +e
output="$(action/run.sh "$binary" "$sample" "$report" "$repo_root/action/annotate.mjs" 2>&1)"
status=$?
set -e
test "$status" -eq 0
grep -F '"schemaVersion": "1"' "$report"
grep -F '"errors": 0' "$report"
grep -F '"warnings": 0' "$report"

escape_report="$test_root/escape-report.json"
printf '%s\n' '{"findings":[{"ruleId":"OF999","source":"forge","severity":"error","message":"bad\n%message","path":"bad,file:qml","line":7,"remediation":"review"}],"summary":{"errors":1}}' > "$escape_report"
escape_output="$(node "$repo_root/action/annotate.mjs" "$escape_report")"
grep -F 'file=bad%2Cfile%3Aqml,line=7,title=OF999 (forge)' <<<"$escape_output"
grep -F 'bad%0A%25message' <<<"$escape_output"

tampered="$test_root/tampered/v0.1.0"
mkdir -p "$tampered"
cp "$dist"/* "$tampered/"
printf 'tampered\n' >> "$tampered/omaforge_0.1.0_linux_amd64.tar.gz"
export OMAFORGE_RELEASE_BASE_URL="file://$test_root/tampered"
if action/install.sh "$version" example/forge >/dev/null 2>&1; then
  echo "installer accepted a tampered archive" >&2
  exit 1
fi
