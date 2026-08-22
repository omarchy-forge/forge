# Product

## Purpose

Omarchy Forge is a local-first developer toolchain for building, testing, and
shipping Omarchy plugins. Its narrow initial promise is to provide the fastest
reliable path from an empty directory to a polished Omarchy bar widget.

Forge complements official Omarchy tooling. Official validation answers
whether the shell can safely recognize and load a plugin structure. Future
Forge checks will address broader quality: documentation, maintainability,
compatibility, testability, security, accessibility, and release readiness.

## Principles

- Local-first and useful without a network connection.
- Deterministic, scriptable behavior suitable for CI.
- Evidence-based compatibility with identified Omarchy versions.
- Safe handling of paths, generated content, and untrusted plugin source.
- One high-quality template before broad template coverage.
- Honest reporting when compatibility is unknown.

## Current state

Milestone 0 contains only the repository foundation and the `omaforge --help`
and `omaforge version` commands. Plugin generation and checking are planned but
not implemented.

## Non-goals

Forge is not a marketplace, plugin directory, theme gallery, community site,
plugin manager, cloud dashboard, account system, or hosted QML execution
service. It does not replace `omarchy plugin add` or
`omarchy plugin validate`. Telemetry, billing, authentication, and a database
are outside the initial roadmap.

## Independence

Omarchy Forge is an independent community project. It is not affiliated with
or endorsed by Omarchy, Basecamp, 37signals, or DHH.
