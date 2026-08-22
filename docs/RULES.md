# Check rule reference

Forge report schema 1 uses stable rule IDs. `official-parity` rules mirror the
locally inspected Omarchy 4 validator contract. `forge` rules add static quality
and publish-readiness guidance; rules described as heuristics require human
review and do not prove a vulnerability or runtime defect.

## Project boundary and official parity

| Rule | Meaning |
| --- | --- |
| `OF001` | Plugin directory is unavailable. |
| `OF100`–`OF103` | Manifest presence, JSON, schema, and required fields. |
| `OF104` | Plugin ID syntax and reserved namespace. |
| `OF105` | Nonempty string kinds. |
| `OF106` | No symlinks outside `.git`. |
| `OF108`–`OF111` | Entry-point object, kind mapping, safe paths, and files. |
| `OF112` | Valid bar-widget default section. |

## Forge quality rules

| Rule | Meaning |
| --- | --- |
| `OF107` | Kind is known to the pinned compatibility contract. |
| `OF200` | Plugin version uses semantic versioning. |
| `OF201`–`OF205` | README, license, compatibility, required sections, and preview readiness. |
| `OF210`–`OF214` | Heuristic theme, keyboard, loading, empty, and error checks for bar widgets. |
| `OF300`–`OF304` | Heuristic packaged-path, privilege, package installation, command, and color checks. |

Every emitted finding includes its severity, evidence location when available,
specific remediation, and an explanation of the rule boundary.

## Doctor diagnostics

`OD100`–`OD105` cover Omarchy detection and package version, Quickshell,
Omarchy Shell IPC, optional QML lint tooling, and the official validator.
Doctor probes are read-only and time-bounded. The official validator result is
reported separately from Forge's deterministic static checks.
