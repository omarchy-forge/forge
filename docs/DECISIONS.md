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

## D-011: Exact-version Action installation

Status: accepted, 2026-08-22.

The composite Action requires an exact semantic Forge version and rejects
moving values such as `latest`. It downloads the architecture-specific archive
and release checksum file, verifies SHA-256 before extraction, and invokes only
the released `omaforge check` command. Workflow inputs enter shell scripts
through quoted environment variables rather than expression interpolation.

GitHub annotations are derived from the schema-v1 JSON report. Workflow command
and property metacharacters are escaped before output. The Action does not run
plugin QML and has no package-install, telemetry, account, or network behavior
beyond downloading the explicitly selected release assets.

## D-012: Tag-only, least-privilege releases

Status: accepted, 2026-08-22.

One reviewed shell script owns local and CI release packaging. It produces
normalized static Linux `amd64` and `arm64` archives and a checksum manifest,
refusing nonempty output directories. The release workflow runs only for exact
`vX.Y.Z` tags. Its build job has read-only contents access; only the dependent
publication job receives `contents: write`.

Third-party release actions are avoided. Official GitHub Actions are pinned to
immutable commit SHAs, and GitHub CLI creates the release from the already
verified workflow artifact. Tags and releases remain explicit external actions
and are not created by normal CI or local packaging.

## D-013: Static-first Geistdocs website boundary

Status: accepted, 2026-08-22.

The documentation website uses Next.js 16, TypeScript, MDX, and the
Geistdocs/Fumadocs rendering foundation in `website/`. The package-backed
structure provides navigation, source-backed search, raw Markdown routes, and
edit links without creating a parallel documentation engine.

The stock initializer's AI chat, feedback submission, analytics, and Markdown
request tracking are explicitly removed or disabled. Next.js telemetry is also
disabled in every package script. The site has no account, authentication,
database, telemetry, or user-data persistence. Only `esbuild` is explicitly
allowed to run an install lifecycle script in pnpm's workspace policy.

The planned `omarchyforge.com` origin is not hard-coded as a deployed fact. A
canonical URL is supplied through `NEXT_PUBLIC_SITE_URL` when deployment is
separately authorized; local builds fall back to localhost. Building the site
does not create a Vercel project, deployment, domain, or environment variable.

## D-014: Separate, local-first Handoff repository

Status: accepted, 2026-08-22.

The flagship Handoff plugin lives in the separate private repository
`omarchy-forge/handoff`. It has an independent installation and eventual
release boundary and does not belong inside the Forge monorepo.

The MVP is one Omarchy `bar-widget` with a native popout. It invokes Git and
the terminal through array-form process arguments, never interpolated shell
commands. It stores canonical project paths, one note, and Git metadata in an
atomically written versioned JSON file under the user's XDG data directory.
It does not execute pinned-project code, hooks, build scripts, or QML and has
no server, account, API key, telemetry, analytics, or required network access.

Repository implementation and private CI do not imply release. Tags, GitHub
releases, marketplace submission, installation for ongoing use, visibility
changes, and announcements remain separate explicit actions.

## D-015: Git-only scaffold targets

Status: accepted, 2026-08-22.

`omaforge init` treats an existing directory containing exactly one real
`.git` directory as an empty scaffold target. This supports the normal workflow
of creating or cloning an empty repository before generating its plugin files.

A `.git` file, symlink, or any additional directory entry still makes the
target nonempty and requires `--force`. Normal generation does not write inside
existing Git metadata; `--git` remains the explicit option for running Git
initialization when requested.

## D-016: Explicitly trusted isolated runtime harness

Status: accepted, 2026-08-22.

Generated bar-widget projects include `tests/runtime`, an optional Omarchy-host
entry-point smoke test. It requires the literal `--trust-plugin-code` argument
because QML is executable code. It runs a copy of the plugin in temporary HOME
and XDG directories, links only to the installed read-only Omarchy QML contract,
checks lifecycle methods and finite geometry, and cleans up a dedicated process
group after a bounded run.

The harness is never called by static checks, generated CI, or the normal test
script. This preserves the rule that untrusted pull-request QML is not executed
automatically while giving a developer an explicit local runtime check after
reviewing their own code.

## D-017: Defer the service-plus-widget template

Status: accepted, 2026-08-22.

Forge will not add a service-plus-widget template from the single installed
first-party example. The service lifecycle is distinct and valuable for shared
multi-monitor state, but the widget-to-service access path is not documented as
a stable third-party contract and no second independent plugin proves a common
shape.

The evaluation in `SERVICE_WIDGET_EVALUATION.md` records the evidence and
revisit criteria. Until one criterion is met, the supported bar-widget template
remains intentionally singular.
