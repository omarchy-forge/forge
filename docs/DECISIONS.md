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
