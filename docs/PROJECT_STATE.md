# Project state

Last updated: 2026-08-22

This is the durable execution checkpoint for future maintainers and coding
agents. `handoff.md` remains the product and safety source of truth; this file
records what has happened since that original handoff. Architecture rationale
lives in `DECISIONS.md`, and verified Omarchy facts live in
`UPSTREAM_COMPATIBILITY.md`.

## Repository state

- Repository: `omarchy-forge/forge` (private at the latest verified checkpoint).
- Default branch: `main`.
- Latest merged commit: `548260d`, merging Milestone 3.
- Active branch: `feat/documentation-website`.
- Active branch commit: `e5f3ab8` (`docs: add documentation website`).
- The active branch is pushed to `origin`; no Milestone 4 pull request has been
  opened or merged yet.
- No website deployment, domain connection, release, package publication, or
  repository-visibility change has been authorized or performed.

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
- Infrastructure exists, but no tag or release has been published.

### Milestone 4 — documentation website (branch complete)

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

1. Review the pushed Milestone 4 branch and open its pull request only with user
   approval.
2. Let required GitHub checks complete; investigate failures before merge.
3. Merge only with explicit user approval.
4. After merge, update this file and `ROADMAP.md` to record the merged commit.
5. Plan the next milestone separately. Website deployment and connection of
   `omarchyforge.com` remain distinct approval-gated work.
