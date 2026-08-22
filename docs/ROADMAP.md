# Roadmap

Roadmap items describe intent, not shipped functionality.

## Milestone 0: private repository foundation

- Core project, contribution, and security documentation.
- Verified notes about the installed Omarchy plugin contract.
- Minimal Go CLI with help and version output.
- Unit tests and baseline CI.

## Milestone 1: bar-widget scaffolder

Implemented and merged:

- Interactive and fully noninteractive `omaforge init`.
- Safe collision and output-path behavior, plus `--dry-run`.
- One polished bar-widget-with-popout template.
- Deterministic golden tests and official-validator integration.

## Milestone 2: check and doctor

Implemented and merged:

- Shared versioned rule model with stable IDs and severities.
- Human-readable, read-only local diagnostics.
- Deterministic CI checks with text, JSON, and SARIF output.
- A clear boundary between official-parity structural validation and Forge rules.

## Milestone 3: action and releases

Implemented and merged:

- Checksum-verifying reusable GitHub Action with PR annotations.
- Tag-driven Linux `amd64` and `arm64` archives with checksums.
- Least-privilege release workflow with immutable Action references.
- Release and Action documentation plus end-to-end packaging tests.

## Milestone 4: documentation website

Implemented and merged:

- Distinct homepage and complete initial documentation navigation.
- Quickstart, templates, commands, anatomy, manifest, and compatibility guides.
- Roadmap and contributing pages grounded in current shipped behavior.
- Static-first architecture without accounts, tracking, or AI chat.
- No deployment or domain connection without explicit approval.

## Milestone 5: Handoff flagship plugin

Implemented in the private, separate `omarchy-forge/handoff` repository:

- Pin a local Git project and save one next-step note.
- Record branch, clean/dirty state, latest commit, and timestamps.
- Open the selected project in the configured terminal.
- Store state atomically under the user's XDG data directory.
- Run without a server, account, telemetry, API key, or required network.
- Pass official validation, Forge checks, private CI, isolated runtime tests,
  and a cleanup-verified Omarchy 4 live session.

Handoff `v0.1.0` has been released from the private repository and verified by
an exact-tag Omarchy installation followed by complete cleanup. It remains
private, unannounced, absent from marketplaces, and not installed for ongoing
use.

## Forge follow-ups from Handoff

- Accept `.git`-only targets without requiring `--force` (implemented after
  `v0.1.0` dogfooding).
- Evaluate a service-plus-widget template.
- Provide a reusable isolated Quickshell entry-point harness (implemented on
  the current development branch with an explicit trust gate).

## Later milestones

- Documentation website deployment and domain connection, only with explicit
  approval.
- Any Handoff marketplace submission, public announcement, visibility change,
  or ongoing local installation, each only with separate explicit approval.
- Additional templates after the first template is proven.
- Trusted local preview, demo, and screenshot tooling.

Marketplace features, hosted plugin execution, accounts, and cloud dashboards
are not planned product scope.
