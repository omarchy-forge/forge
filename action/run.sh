#!/usr/bin/env bash
set -euo pipefail

binary="${1:?binary path is required}"
target="${2:-.}"
report="${3:?report path is required}"
annotator="${4:?annotator path is required}"

set +e
"$binary" check "$target" --format json > "$report"
status=$?
set -e

node "$annotator" "$report"
exit "$status"
