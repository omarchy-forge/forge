# Decision record

## D-001: Go and the standard library

Status: accepted, 2026-08-22.

The CLI uses Go and initially has no third-party Go dependencies. This supports
a small native binary, deterministic builds, and direct unit testing without a
framework. A dependency may be introduced later only for a concrete correctness
or maintainability benefit.

## D-002: executable name

Status: accepted, 2026-08-22.

The executable is `omaforge`. The shorter name avoids collisions associated
with the generic executable name `forge` while retaining the product identity.

## D-003: honest development metadata

Status: accepted, 2026-08-22.

Local builds report version `dev`, with commit and build date as `unknown` when
not injected. Release automation may set those values with linker flags later.

## D-004: upstream validation boundary

Status: accepted, 2026-08-22.

Forge will complement, not replace, `omarchy plugin validate`. Future local
diagnostics should call the official validator when available. Deterministic
CI rules must identify their upstream capability snapshot and distinguish
structural parity checks from Forge quality rules.

## D-005: CI action references

Status: accepted, 2026-08-22.

Foundation CI uses only GitHub-maintained actions pinned to stable major
versions. Workflow permissions are explicitly read-only. Before a public or
release-oriented security boundary is introduced, action references should be
reconsidered for immutable commit pinning and automated update handling.

## D-006: embedded, deterministic templates

Status: accepted, 2026-08-22.

Plugin template files are embedded in the Go binary and rendered in sorted path
order. User values never determine generated filenames. JSON and QML strings
use JSON encoding, while prose values are escaped for Markdown. A golden tree
records each output path, mode, and content digest.

## D-007: overwrite model

Status: accepted, 2026-08-22.

The default is to refuse any nonempty target. `--force` preserves unrelated
files and overwrites only paths owned by the selected template after printing a
complete plan. Symlinks anywhere in the target are rejected before forced
writes. `--dry-run` performs the same validation and rendering without writing.

## D-008: Milestone 1 upstream pins

Status: accepted, 2026-08-22.

Template design was checked against official Omarchy `quattro` revision
`ed7bae4ac5a570e9df307486e0202fdafcc6ee24` and the official Basecamp plugin
revision `abc1ba72aaf47db530d2a0c6901d99f0f98e6aa7`. The installed Omarchy package
remains authoritative for local validation and runtime contract checks.

## D-009: Versioned findings and exit codes

Status: accepted, 2026-08-22.

Check reports use schema version 1 and stable rule IDs. Rules explicitly name
their source as `official-parity` or `forge`; heuristic Forge rules say so in
their explanation. Findings are sorted independently of filesystem traversal.
Text, JSON, and SARIF serialize the same report.

`omaforge check` returns 0 when there are no error-severity findings, including
when warnings exist; 1 for project errors or report-write failures; and 2 for
invalid usage. This allows publish-readiness warnings to remain visible without
making early CI adoption unexpectedly blocking.

## D-010: Static check and local doctor boundary

Status: accepted, 2026-08-22.

`check` is reproducible, noninteractive, network-free, and never runs QML or
external validators. `doctor` is intentionally environment-aware and may run
read-only installed commands with bounded timeouts, including the authoritative
Omarchy validator. Both share Forge's project rule report.
