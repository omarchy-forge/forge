# Project state

Last updated: 2026-08-22

This is the durable execution checkpoint for future maintainers and coding
agents. `handoff.md` remains the product and safety source of truth; this file
records what has happened since that original handoff. Architecture rationale
lives in `DECISIONS.md`, and verified Omarchy facts live in
`UPSTREAM_COMPATIBILITY.md`.

## Repository state

- Repository: `omarchy-forge/forge` (public).
- Default branch: `main`.
- Latest completed checkpoint before this update: public-visibility documentation
  pull request `#15`, merged as `dec9469`.
- Milestone 4 implementation commit: `e5f3ab8` (`docs: add documentation
  website`).
- Milestone 4 pull request: `#4` (merged on 2026-08-22 after all required CI
  checks passed).
- Forge's public visibility has been authorized and applied. No website
  deployment, domain connection, Handoff visibility change, marketplace
  submission, or public launch announcement has been authorized or performed.

Always verify Git and GitHub state instead of assuming this snapshot is still
current.

## Delivered milestones

### Milestone 0 — repository foundation

- Product, architecture, security, contribution, compatibility, and roadmap
  documentation.
- Minimal Go CLI with honest build metadata, help, and version behavior.
- Baseline least-privilege CI and repository governance.

### Milestone 1 — bar-widget scaffolder

- `omaforge init` generates a theme-aware Omarchy bar-widget plugin.
- Deterministic templates, dry-run behavior, collision protection, tests, and
  optional Git initialization.
- Runtime-safe demo states and documentation.

### Milestone 2 — checks and doctor

- `omaforge check` provides deterministic structural, security, UX, and
  publish-readiness diagnostics without executing plugin QML.
- `omaforge doctor` reports local development capabilities and official
  validator availability.
- Human, JSON, and GitHub-oriented output with tested exit semantics.

### Milestone 3 — GitHub Action and release pipeline

- Checksum-verifying reusable Action with pull-request annotations.
- Tag-driven Linux `amd64` and `arm64` release archives and checksums.
- Least-privilege workflows with immutable Action references.
- Release/Action integration tests and documentation.
- Public release `v0.1.0` is tagged at
  `4c9dc57f753b1ba4d2167e3ce7022de1ae106f50` and published at
  `https://github.com/omarchy-forge/forge/releases/tag/v0.1.0`.
- The tag workflow built and published static Linux `amd64` and `arm64`
  archives plus `checksums.txt`. Independent downloads passed checksum,
  archive-content, architecture, and embedded build-metadata verification.
- The release is neither draft nor prerelease and uses the reviewed release
  notes in `docs/RELEASE_NOTES_0.1.0.md`.

### Milestone 4 — documentation website

- Static-first Next.js 16, React 19, Geistdocs, Fumadocs, and MDX website under
  `website/`.
- Distinct responsive homepage plus quickstart, templates, commands, plugin
  anatomy, manifest, compatibility, roadmap, and contributing pages.
- Search and generated machine-readable documentation routes.
- AI chat, analytics, feedback, tracking, accounts, authentication, database,
  uploads, CMS ingestion, and remote image hosts are disabled or absent.
- CI installs from the lockfile, audits dependencies, typechecks, and produces
  the production build.
- Desktop and mobile Chromium renders, route content, dev-server output, and the
  production build were verified locally.

### Milestone 5 — Handoff flagship plugin

- Private separate repository: `omarchy-forge/handoff`.
- Private release: `v0.1.0`, tagged at
  `afa10733292e9d8f2acb9f197eb5802b1eb42422` and published at
  `https://github.com/omarchy-forge/handoff/releases/tag/v0.1.0`.
- Local-first Omarchy bar widget for pinning Git projects, saving one next-step
  note, showing branch/dirty/latest-commit context, and opening a terminal.
- Atomic state under the user's XDG data directory, with no server, account,
  API key, telemetry, analytics, or required network access.
- Official validation, Forge checking, private CI, isolated Quickshell tests,
  and a controlled Omarchy 4 live session passed.
- The live session found and fixed missing-state handling in `1a71b44`, tested
  ready/empty/error states and keyboard behavior, and produced a fictional
  panel-only preview.
- Cleanup restored `shell.json` byte-for-byte, removed the plugin and test data,
  and left `omarchy-shell shell ping` healthy.
- The private `v0.1.0` source release was installed from an independent exact-tag
  clone with the official Omarchy plugin command. The installed checkout matched
  the released commit, the ready/empty/error demos passed, and shell IPC stayed
  healthy.
- Post-release cleanup removed the plugin and its test data, restored
  `shell.json` byte-for-byte (SHA-256
  `802fa2600cac1cd2971c48769661432a8f30eb5beb2eadb63f0356a913172f9f`),
  and left `omarchy-shell shell ping` returning `ok`.
- Handoff remains private, is not installed for ongoing use, and has not been
  submitted to a marketplace or announced publicly.

### Milestone 6 — local development and visual tooling

- The first increment adds `omaforge dev <directory> --trust-plugin-code` as a
  one-shot entry point for the generated project's isolated Quickshell harness.
- The CLI validates the project, manifest, and regular-file harness before
  executing it with argument-safe process invocation.
- A freshly generated plugin passed static Forge checks and the real isolated
  runtime through the new command on the verified local Wayland session.
- File watching, richer state simulation, and screenshot capture remain later
  increments; the command does not install or enable plugins or mutate live
  shell configuration.

## Forge feedback from Handoff

Dogfooding recorded three follow-up opportunities. The first and third are
merged:

1. `omaforge init` treats a freshly cloned, `.git`-only directory as an empty
   target without requiring `--force`, while preserving the Git metadata.
2. Consider a service-plus-widget template for stateful plugins.
3. Generated projects provide an opt-in isolated Quickshell entry-point harness
   with an explicit trusted-code gate and no live-shell installation.

The service-plus-widget opportunity was evaluated after the harness merge and
deferred. The installed package has only one combined example, it is
first-party, and its widget-to-service lookup is not documented as a stable
third-party contract. See `SERVICE_WIDGET_EVALUATION.md` for the evidence and
revisit criteria.

## Current verification baseline

The Milestone 4 branch passed:

- `go test -race ./...`
- `go vet ./...`
- `go build ./cmd/omaforge`
- `tests/release-action.sh`
- Go formatting verification
- workflow/action YAML parsing
- `git diff --check`
- secret and personal-path scans
- `pnpm install --frozen-lockfile`
- `pnpm audit:known`
- `pnpm typecheck`
- `pnpm build` (40 generated routes at the recorded checkpoint)

Re-run applicable checks after any change; this list is evidence, not a waiver.

## Dependency-security checkpoint

`website/SECURITY.md` documents three reviewed transitive advisories:

- `GHSA-w3rx-r6r6-pgpr`
- `GHSA-5p2g-fcmc-qvqq`
- `GHSA-g7r4-m6w7-qqqr`

The first two affect `image-size@2.0.2`, for which no compatible patched release
was available at the checkpoint. The third affects esbuild's Windows
development server, which this project does not invoke. `pnpm audit:known`
allows only these exact advisory IDs and fails for any additional finding.
Remove allowlist entries as patched compatible dependencies become available.

## Standing boundaries

- Keep the CLI local-first, deterministic, and network-independent by default.
- Never execute third-party QML during checks.
- Do not edit Omarchy-owned files or user shell configuration.
- Do not add telemetry, analytics, accounts, authentication, databases, or
  hosted execution without an explicit architecture decision and authorization.
- Do not push, publish, deploy, connect the domain, create a release, change
  visibility, or make announcements without explicit user approval.
- Describe planned and branch-only behavior accurately; do not call it shipped
  until it is merged or released as appropriate.

## Next checkpoint

1. Keep the website undeployed and Handoff private until those boundaries are
   approved separately.
2. Choose the next product milestone from real usage; do not expand the
   template set without the evidence recorded in the roadmap.
3. Ongoing local installation, marketplace submission, announcement, website
   deployment, and connection of `omarchyforge.com` remain distinct
   approval-gated work.
