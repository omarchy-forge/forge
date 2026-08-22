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

Force mode overwrites only colliding generated paths, preserves unrelated
files, rejects symlinks, and prints its complete plan before writing. Dry-run
uses the same validation and rendering path but performs no writes.

## Future boundaries

As milestones require them, focused internal packages may own manifest parsing,
compatibility capabilities, checks, diagnostics, and report serialization.
Those packages should not depend on terminal interaction so the
same deterministic behavior can serve the CLI and CI.

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
