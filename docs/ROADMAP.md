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
- Public Forge `v0.1.0` release with verified Linux `amd64` and `arm64`
  archives and checksums.
- Public Forge `v0.2.0` release with the completed Milestone 6 trusted local
  development workflow and independently verified Linux archives.
- Public Forge `v0.3.0` release with provider-neutral agent-ready scaffolding,
  a human-controlled specification gate, and independently verified Linux
  archives.

## Milestone 4: documentation website

Implemented and merged:

- Distinct homepage and complete initial documentation navigation.
- Quickstart, templates, commands, anatomy, manifest, and compatibility guides.
- Roadmap and contributing pages grounded in current shipped behavior.
- Beginner-first installation and first-plugin flow for users new to Linux.
- Static-first architecture without accounts, tracking, or AI chat.
- Production deployment at `www.omarchyforge.com`, with the apex domain
  redirecting to the canonical `www` origin.

## Milestone 5: Handoff flagship plugin

Implemented in the public, separate `omarchy-forge/handoff` repository:

- Pin a local Git project and save one next-step note.
- Record branch, clean/dirty state, latest commit, and timestamps.
- Open the selected project in the configured terminal.
- Store state atomically under the user's XDG data directory.
- Run without a server, account, telemetry, API key, or required network.
- Pass official validation, Forge checks, private CI, isolated runtime tests,
  and a cleanup-verified Omarchy 4 live session.

Handoff `v0.1.0` was released while the repository was private and verified by
an exact-tag Omarchy installation followed by complete cleanup. The repository
is now public and appears in Forge's owner-project catalog; it remains absent
from marketplaces and is not installed for ongoing use.

Handoff `v0.1.1` was also released initially as a development-tooling
maintenance update adopting Forge `v0.2.0`; it passed exact-tag installation
and cleanup verification without changing those distribution boundaries.

## Forge follow-ups from Handoff

- Accept `.git`-only targets without requiring `--force` (implemented after
  `v0.1.0` dogfooding).
- Evaluate a service-plus-widget template (completed and deferred pending a
  documented third-party service-access contract or a second proven example).
- Provide a reusable isolated Quickshell entry-point harness (implemented after
  `v0.1.0` dogfooding with an explicit trust gate).
- Adopt the released Forge `v0.2.0` development harness and Action in private
  Handoff (implemented; no Handoff visibility or release change).

## Completed later milestones

- Milestone 7 agent-ready scaffolding: opt-in structured specifications,
  acceptance criteria, and project-scoped agent safety guidance are released
  in Forge `v0.3.0` without adding an AI provider or automatic agent launch.
- Milestone 6 local development tooling: the first one-shot trusted
  `omaforge dev` runtime check and fictional ready/empty/error state simulation
  are implemented, along with deterministic plugin-only PNG capture and
  isolated file-watching sessions.
- Production documentation deployment and custom-domain connection are
  complete.

## Still later

- Release the guided `omaforge init --agent` workflow after protected review
  and a supervised end-to-end trial. The implementation is currently
  unreleased and must not be described as part of `v0.4.0`.

- Any Handoff marketplace submission, public announcement, visibility change,
  or ongoing local installation, each only with separate explicit approval.
- Additional templates after the first template is proven.

Marketplace features, hosted plugin execution, accounts, and cloud dashboards
are not planned product scope.
