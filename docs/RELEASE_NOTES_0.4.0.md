# Omarchy Forge 0.4.0 release notes

Omarchy Forge `v0.4.0` adds a verified user-controlled install and update path,
an owner-project catalog for optional Forge tools, and two diagnostic safety
improvements. The CLI remains local-first and network-independent by default.

## Install and update

- The canonical website installer resolves the latest stable GitHub release or
  accepts an exact `--version` tag.
- Release archives are selected for supported Linux architectures and verified
  against the published SHA-256 manifest before installation.
- An already-current installation is left unchanged. Updates occur only when
  the user invokes the installer; there are no background checks, package
  manager changes, privilege escalation, or shell-configuration edits.
- The full manual checksum-verification path remains documented for users who
  do not want to pipe a reviewed remote script into Bash.

## Owner projects

- The website Projects page lists public tools maintained by the
  `omarchy-forge` organization, beginning with the Handoff plugin and Omaudit
  capability scanner.
- Catalog data and preview images are checked into Forge so website builds stay
  deterministic and network-independent.
- Forge derives installation commands from validated repository, release, and
  package or manifest data. Projects cannot inject arbitrary commands or URLs.
- The protected synchronizer downloads bounded metadata and images only. It
  never installs projects, invokes installers, or executes plugin QML.

## Diagnostics and safety

- Static checking rejects NUL bytes in plugin source with error rule `OF305`,
  closing a parser ambiguity before downstream tooling receives the file.
- `omaforge doctor` now finds `qmllint` in Qt 6's standard tool directory when
  it is installed but absent from the user's `PATH`, and reports its exact
  location. The linter remains optional and is not executed by Forge checks.

## Compatibility and boundaries

This release continues to target the inspected Omarchy 4 manifest and plugin
contracts. Forge does not install or enable plugins, execute untrusted QML in
static checks or CI, modify Omarchy-owned files, edit `shell.json`, add
telemetry, require an account, or perform automatic package installation.

Handoff and Omaudit are optional separate projects. Their inclusion in the
owner-maintained directory is not a third-party marketplace or an endorsement
program.
