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

In development on `feat/github-action-releases`:

- Checksum-verifying reusable GitHub Action with PR annotations.
- Tag-driven Linux `amd64` and `arm64` archives with checksums.
- Least-privilege release workflow with immutable Action references.
- Release and Action documentation plus end-to-end packaging tests.

## Later milestones

- Documentation-first website, deployed only with explicit approval.
- A separately released Handoff reference plugin.
- Additional templates after the first template is proven.
- Trusted local preview, demo, and screenshot tooling.

Marketplace features, hosted plugin execution, accounts, and cloud dashboards
are not planned product scope.
