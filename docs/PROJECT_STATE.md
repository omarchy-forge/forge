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
- Latest completed checkpoint: agent-ready hardening pull request `#30`, merged
  as `d42bc38` after both Go and website CI passed. Its documentation source was
  deployed to the production site and verified at the custom domain.
- Milestone 4 implementation commit: `e5f3ab8` (`docs: add documentation
  website`).
- Milestone 4 pull request: `#4` (merged on 2026-08-22 after all required CI
  checks passed).
- Forge's public visibility, documentation deployment, and custom-domain
  connection have been authorized and applied. Handoff remains private; no
  marketplace submission or public launch announcement has been authorized or
  performed.

Always verify Git and GitHub state instead of assuming this snapshot is still
current.

Forge `v0.2.0` is published from exact commit
`f8f5cb9971ac4966afc418229d9e9902dbb5e09d`. Its tag workflow passed, and
independent downloads verified both archive checksums, amd64 and arm64 ELF
architectures, archive contents, and embedded amd64 build metadata. The public
release is at `https://github.com/omarchy-forge/forge/releases/tag/v0.2.0` and
uses the reviewed `docs/RELEASE_NOTES_0.2.0.md` description.

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
- Private maintenance release `v0.1.1` is tagged at
  `6038185070e580696326a6e8591c761d700dd63f` with the Forge `v0.2.0`
  development workflow. Its private CI passed, and an independent exact-tag
  installation exercised ready/empty/error before complete removal, test-data
  cleanup, byte-for-byte `shell.json` restoration, and healthy shell IPC.

### Milestone 6 — local development and visual tooling

- The first increment adds `omaforge dev <directory> --trust-plugin-code` as a
  one-shot entry point for the generated project's isolated Quickshell harness.
- The CLI validates the project, manifest, and regular-file harness before
  executing it with argument-safe process invocation.
- A freshly generated plugin passed static Forge checks and the real isolated
  runtime through the new command on the verified local Wayland session.
- The second increment adds validated `ready`, `empty`, and `error` state
  selection. Generated plugins apply these fictional states in memory inside
  the temporary runtime and must explicitly accept the requested state.
- Handoff exposes the same `setDemoState` vocabulary but its released tree
  predates the generated harness. A disposable exact working-tree copy using
  the current harness exercised all three states without modifying the private
  repository, installing the plugin, or touching live shell configuration.
- Handoff subsequently adopted the released Forge `v0.2.0` harness, explicit
  screenshot target, refresh contract, and public Forge Action on its private
  `main` branch at `84a7a791f23099ae19e877f61ea830fcaca57301`.
  Its private CI and local ready/empty/error runtime plus panel-only screenshot
  verification passed without creating a new Handoff release.
- The third increment adds `omaforge screenshot` using a template-declared Qt
  Quick item. Ready, empty, and error produced plugin-only RGBA PNGs in the real
  isolated runtime; no compositor or desktop capture API is used and existing
  files are not overwritten.
- The fourth increment adds `omaforge dev --watch`: an initial isolated run,
  deterministic content fingerprinting outside `.git`, bounded polling, and
  fresh isolated reruns after changes. Failed runs do not end the session, and
  Ctrl-C cancels it cleanly.
- Milestone 6's planned dev-state, screenshot, and watch slices are now
  implemented. None installs or enables plugins or mutates live shell
  configuration.

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

The three previously reviewed transitive advisories are resolved by overriding
the Geistdocs-pinned Fumadocs family to current compatible releases. This
replaces vulnerable `image-size` and upgrades esbuild. The production audit
reports zero advisories, and `pnpm audit:known` now fails on any finding; no
advisory allowlist remains. See `website/SECURITY.md`.

## Documentation deployment checkpoint

The documentation site is live at `https://www.omarchyforge.com` in the Vercel
project `omarchy-forge-docs`; `https://omarchyforge.com` redirects to that
canonical origin. Production sets `NEXT_PUBLIC_SITE_URL` to the custom origin.
The latest verified production checkpoint is deployment
`dpl_AMiBYWPoMx7oqG56Qut7QkiY3nGM`, built from clean Forge commit `d42bc38` on
2026-08-22. Its Vercel build generated 40 routes and completed successfully.
HTTP verification confirmed that the live template guide contains the
human-controlled `Draft` to `Ready for implementation` agent handoff and that
the apex domain redirects to the canonical `www` origin. The prior broader
homepage, search, sitemap, robots, LLM-readable content, and browser checks
remain recorded evidence for the same site architecture.

## Beginner onboarding checkpoint

The root README and website Quickstart lead with the same numbered first-plugin
journey: enter an Omarchy terminal, install the exact Forge release without
`sudo`, create a plugin interactively, check and officially validate it, preview
fictional states in isolation, enable reviewed code in Omarchy, remove it after
testing, and capture a plugin-only screenshot. The installation block detects
the supported CPU architecture, downloads only the matching archive, and
verifies the published checksum. The exact commands were exercised against the
public `v0.2.0` release. Go and source-build instructions are separated as a
contributor-only path. Pull request `#27` passed both copies of the Go and
website CI suites, merged as `fbaa31c`, and the rendered production Quickstart
was browser-verified with no framework overlay, console errors, or page errors.

## Agent-ready scaffolding checkpoint

The next-release `omaforge init --agent-ready` option conditionally adds exactly
two files to the unchanged bar-widget scaffold: `FORGE_SPEC.md` captures the
product, data, privacy, state, failure, timeout, and acceptance boundary;
`AGENTS.md` constrains a coding agent to the project and prohibits QML execution,
installation, privileged operations, and unapproved capability expansion.
Forge does not select or invoke an agent and adds no provider dependency,
credential, or network path. Focused CLI and scaffold tests pass, and a real
agent-ready output passed both `omaforge check` and the installed official
Omarchy validator. This capability is implemented on `main` for the next
release and is not part of public `v0.2.0`.

An end-to-end static dogfood pass generated a disposable local-only disk-space
widget specification and implemented it while following the generated agent
contract. The project tests, Forge check, and installed official validator all
passed without executing QML, installing the plugin, or changing shell
configuration. Dogfooding found that the original instruction to stop while
`FORGE_SPEC.md` contained the word `TODO` was self-contradictory because the
file's explanatory text retained that word after completion. The template now
uses an explicit human-controlled `Draft` to `Ready for implementation` status
gate, and its allowed static commands include the generated `./tests/run`
entry point. The installed `omarchy agent prompt` command was inspected but no
agent was launched automatically.

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

1. Keep Handoff private. The documentation website and custom domain are live
   at `https://www.omarchyforge.com`.
2. Prepare and review the next Forge release, including release notes and the
   full baseline, before documenting agent-ready scaffolding as available in a
   public binary. Tagging and publication still require explicit authorization.
3. Ongoing local installation, marketplace submission, announcement, and any
   Handoff visibility change remain distinct approval-gated work.
