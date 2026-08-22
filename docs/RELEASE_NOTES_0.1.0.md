# Omarchy Forge 0.1.0 release notes

Omarchy Forge `v0.1.0` is the first private source release of the local-first
developer toolchain for building, checking, and shipping Omarchy plugins.

## Included

- `omaforge init` with one polished bar-widget-with-popout template,
  deterministic output, dry runs, collision protection, `.git`-only target
  support, optional CI, and optional local Git initialization.
- `omaforge check` with deterministic text, JSON, and SARIF reports, stable
  findings, official-parity structural rules, and Forge quality rules.
- `omaforge doctor` with read-only local Omarchy, Quickshell, shell IPC, QML
  tooling, and official-validator diagnostics.
- A checksum-verifying composite GitHub Action that downloads an exact release
  and never executes plugin QML.
- Reproducible static Linux archives for `amd64` and `arm64` plus SHA-256
  checksums.
- Generated fictional demo states and an explicitly trusted, isolated local
  Quickshell entry-point smoke test.
- The undeployed documentation website source and complete initial product,
  compatibility, security, and contributor documentation.

## Security and privacy

Forge has no telemetry, account, authentication, database, hosted execution,
or required network access for normal CLI use. Static checks do not execute
plugin QML. The optional generated runtime test requires the literal
`--trust-plugin-code` acknowledgement and isolates HOME, XDG state, logs, and
the process group.

The release does not deploy the website, connect `omarchyforge.com`, change
repository visibility, submit anything to a marketplace, or announce Forge
publicly.

## Compatibility

This release targets the inspected Omarchy 4 manifest schema and was verified
against installed package `omarchy 4.0.0-1`. Compatibility evidence and known
uncertainties are recorded in `docs/UPSTREAM_COMPATIBILITY.md`.

## Private release use

While the repository remains private, downloading release assets and using the
composite Action requires GitHub credentials with access to
`omarchy-forge/forge`. Verify `checksums.txt` before using a downloaded archive.
