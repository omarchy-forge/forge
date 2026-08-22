# Architecture

## Current implementation

The project is a Go module with one executable entry point in
`cmd/omaforge`. Command parsing and output live in `internal/cli`, allowing
behavior and exit codes to be tested without spawning subprocesses.

The CLI uses only the Go standard library. Release metadata is held
in package variables in `cmd/omaforge` so a future release build can inject a
version, commit, and build date using Go linker flags. Local builds report an
explicit development version and unknown metadata rather than implying a
release provenance.

`internal/scaffold` owns input and target validation, collision planning, and
filesystem writes. It validates the complete request and renders every file
before writing. `templates` embeds the reviewed bar-widget template and renders
validated values with format-specific escaping. Generated paths are fixed by
the template rather than derived from user strings.

`internal/dev` validates the requested project and its regular-file runtime
harness, then delegates to that project-owned harness only after the CLI has
received the literal trust acknowledgement. Isolation, timeout, and cleanup
remain owned by the generated harness so its behavior is reviewable alongside
the plugin version it executes.

Force mode overwrites only colliding generated paths, preserves unrelated
files, rejects symlinks, and prints its complete plan before writing. Dry-run
uses the same validation and rendering path but performs no writes.

`internal/checks` owns the versioned finding model, deterministic static rule
engine, summary calculation, and text, JSON, and SARIF serialization. It never
executes plugin code or accesses the network. Every finding has a stable rule
ID, severity, evidence path when available, remediation, and a source boundary:
`official-parity` or `forge`.

`internal/doctor` composes the same project report with read-only local
environment probes. External commands have short timeouts. Doctor may invoke
the installed official validator, but deterministic `check` never depends on
locally installed tooling.

`action/` is a composite GitHub Action with no vendored package dependencies.
It accepts workflow inputs through environment variables, downloads only an
exact semantic version, verifies the selected archive against the released
checksum file, and then invokes `check`. A small Node script converts the JSON
report to escaped GitHub workflow annotations.

`scripts/build-release.sh` is the single local and CI packaging path. It builds
static Linux `amd64` and `arm64` binaries, injects version provenance, creates
normalized archives, and emits SHA-256 checksums. The tag workflow separates a
read-only build job from the write-authorized release publication job.

`website/` is an undeployed Next.js 16 and MDX application backed by Geistdocs and
Fumadocs. Checked-in content is the source of truth. The website retains static
documentation rendering, navigation, local corpus search, raw Markdown, and
edit links while explicitly disabling AI chat, feedback submission, analytics,
Markdown-request tracking, and Next.js telemetry. It has no database,
authentication, account, or user-data persistence layer.

## Future boundaries

As later milestones require them, focused internal packages may own additional
compatibility capabilities. Core packages do not depend on terminal
interaction, so the same deterministic behavior can serve the CLI and CI.

Official Omarchy validation remains authoritative for installed environments.
Any internal structural rules needed for network-free CI must be versioned,
traceable to upstream evidence, and tested for parity. Static checking must not
execute plugin QML.

Templates and compatibility fixtures will be added only when they contain real,
tested content. A documentation website and reusable Action are deferred until
the corresponding CLI functionality exists.

## Data and network

Normal CLI operation is local and network-independent. The architecture has no
database, account, authentication, telemetry, or background service. Commands
that may later invoke official local tools must expose failures rather than
silently claiming compatibility.
