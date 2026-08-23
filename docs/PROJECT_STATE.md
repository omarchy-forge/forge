# Project state

Last updated: 2026-08-23

This is the durable execution checkpoint for future maintainers and coding
agents. `handoff.md` remains the product and safety source of truth; this file
records what has happened since that original handoff. Architecture rationale
lives in `DECISIONS.md`, and verified Omarchy facts live in
`UPSTREAM_COMPATIBILITY.md`.

## Repository state

- Repository: `omarchy-forge/forge` (public).
- Default branch: `main`.
- Current completed product checkpoint: public Forge `v0.4.0`, including the
  explicit installer/update path, two-project Handoff/Omaudit owner catalog,
  `OF305` NUL-byte rejection, and Qt-directory `qmllint` discovery. Release
  metadata pull request `#48` merged at
  `3584e769e490f9bcf3ed773d9f142a191af046d3`; release workflow run
  `32633386500` passed both build and publish jobs. Pull request `#49` updated
  current-facing installation and Action guidance at
  `6b5c4ab18aabfa7d9840cbc68141c125533ee434`; production deployment
  `dpl_7HppBciD4boi1SUbjfmy23W7WQx8` completed all 42 routes and was verified at
  the canonical domain.
- Milestone 4 implementation commit: `e5f3ab8` (`docs: add documentation
  website`).
- Milestone 4 pull request: `#4` (merged on 2026-08-22 after all required CI
  checks passed).
- Forge's public visibility, documentation deployment, and custom-domain
  connection have been authorized and applied. Handoff was made public on
  2026-08-23 after its release, checks, history scan, metadata, and protections
  were verified. No marketplace submission or public launch announcement has
  been performed.

Always verify Git and GitHub state instead of assuming this snapshot is still
current.

Forge `v0.4.0` is published from exact commit
`3584e769e490f9bcf3ed773d9f142a191af046d3`. Its tag workflow passed, and
independent downloads verified both archive checksums, amd64 and arm64 ELF
architectures, archive contents, and embedded amd64 build metadata. The public
release is at `https://github.com/omarchy-forge/forge/releases/tag/v0.4.0` and
uses the reviewed `docs/RELEASE_NOTES_0.4.0.md` description. The canonical
installer resolved and installed this release successfully in both an isolated
home and the user's normal installation. A generated project then passed
`omaforge doctor`, including exact `qmllint` discovery and official Omarchy
validation.

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

- Separate public repository: `omarchy-forge/handoff`.
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
- Handoff is public, is not installed for ongoing use, and has not been
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

### Milestone 7 — public owner-project catalog

- Handoff is public with protected `main`, secret scanning, public release
  `v0.1.1`, `omaforge-project` metadata, and the official Omarchy installation
  command derived by Forge rather than supplied as arbitrary project text.
- The website Projects page lists Handoff with a checked-in local preview,
  version, compatibility, source and release links, and a copyable install
  command. The homepage, top navigation, sitemap, README, roadmap, and catalog
  documentation link to or describe the feature.
- A strict generator accepts only public, active, non-fork repositories owned
  by `omarchy-forge` with the explicit topic, valid manifest and metadata, and
  a matching stable release. JSON and preview downloads are bounded; no QML is
  installed or executed.
- A daily and manually dispatched workflow regenerates checked-in artifacts,
  opens a protected pull request only when they change, starts CI explicitly,
  and requests auto-merge. Manual run `32630618972` passed as a no-op against
  the initial catalog.
- Live verification returned 200 for `/projects`, the 704×787 Handoff preview,
  and Next's optimized image route; it also confirmed page content, exact
  installation command, sitemap entry, and public Handoff repository access.
- Omaudit was fully reviewed, hardened, tested, released as `v0.1.0`, and moved
  from the owner's personal account to the public `omarchy-forge/omaudit`
  repository. Its protected `main` is at
  `c4dc6ac965cfe23b2e60dfd23473d5cbc33565f6`; release assets and the explicit
  checksum-verifying installer passed independent isolated verification.
- The catalog supports both Omarchy plugins and CLI tools. CLI entries require
  matching PEP 621 package metadata, an identity-bound strict-Bash installer,
  and a matching stable release. Forge derives an immutable exact-version
  command and never executes the installer during synchronization or CI.
- The generated catalog now contains Handoff and Omaudit, including a checked-in
  Omaudit preview and a visible project-type label on the Projects page.

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
The agent-ready hardening production checkpoint is deployment
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

Post-release installation dogfooding found that the copy-paste block left the
interactive terminal inside its `mktemp` download directory. The root README
and website Quickstart now isolate installation in a subshell, remove that
directory through an exit trap on both success and failure, and finish in
`$HOME`. The exact public `v0.3.0` block was verified with an isolated home;
the installed binary reported the expected release metadata, temporary files
were removed, and both successful and forced-failure paths returned home.
Pull request `#35` merged this correction as `eb97ab0` after all Go and website
CI checks passed. Production deployment `dpl_FULYL3kk7p5hz776VP5EPp95joDR`
completed a 40-route build, is aliased to `https://www.omarchyforge.com`, and
had no error-level production logs in the post-deploy scan.

## Agent-ready scaffolding checkpoint

The released `omaforge init --agent-ready` option conditionally adds exactly
two files to the unchanged bar-widget scaffold: `FORGE_SPEC.md` captures the
product, data, privacy, state, failure, timeout, and acceptance boundary;
`AGENTS.md` constrains a coding agent to the project and prohibits QML execution,
installation, privileged operations, and unapproved capability expansion.
Forge does not select or invoke an agent and adds no provider dependency,
credential, or network path. Focused CLI and scaffold tests pass, and a real
agent-ready output passed both `omaforge check` and the installed official
Omarchy validator. This capability is published in Forge `v0.3.0`.

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

The public `v0.3.0` amd64 binary independently generated an agent-ready sample
from the downloaded release archive. Both generated contract files contained
the expected readiness gate, and the sample passed `omaforge check` plus the
installed official validator without executing QML.

## Guided-agent development checkpoint (unreleased)

Project Clock dogfooding exposed two workflow gaps: users could not move from
an idea to the existing specification without manual editing, and generated
documentation listed live `demo/run` commands before explaining that they need
an installed plugin. A local, unmerged implementation adds interactive
`omaforge init --agent` while preserving `--agent-ready` unchanged.

The guided mode reports the configured Omarchy agent, gathers product and
access boundaries, shows an editable summary and unattended-permission warning,
and writes nothing on cancellation. Confirmation generates a completed
`FORGE_SPEC.md`, project `AGENTS.md`, and `AGENT_PROMPT.md`; initializes Git;
creates a hook- and signing-disabled baseline commit; and invokes the configured
agent from the project. Generated completion instructions require a local
implementation commit and order human review, isolated ready/empty/error tests,
deliberate local installation, installed demos, removal, and optional pushing.
Forge still never executes QML, installs the plugin, configures a remote, or
pushes automatically.

Focused tests simulate agent discovery, commit, launch, cancellation, missing
configuration, and unsafe option combinations without launching an agent.
The guided flow also accepts up to ten optional local text, Markdown, PNG,
JPEG, WebP, or non-rendered XML-validated SVG references. It validates bounded regular files, captures their
bytes from a generated `references/` drop directory, records SHA-256
digests, and warns that the configured agent/provider may receive the contents.
Dogfooding the guided flow with a local control-center project showed that
answers such as `icon` and `dashboard` can erase the functionality stated in a
detailed reference. The flow now offers `references` or `questionnaire`.
Reference mode makes confirmed text and image files the primary product brief,
requires a complete functional and visual inventory, and asks only for local
command, network, persistence, and failure boundaries. Questionnaire mode keeps
the detailed bar, dashboard, action, data-source, and command prompts. Generated
prompts forbid treating reference controls as decoration or silently reducing
them to placeholder UI.

Reference-mode dogfooding also exposed an initialization-order regression:
Forge attempted full specification validation before collecting the author,
causing an empty-input retry loop. Two test runs remained alive, exhausted the
`/tmp` tmpfs quota, and approached system OOM. The prompt now collects common
identity fields before reference staging and exits safely when no reference is
added. A regression test covers that cancellation path. The exact runaway test
processes and their three stale Go build directories were removed; the machine
returned to roughly 12 GiB available memory and `/tmp` returned to 2% usage.

After the fix, bounded focused guided-agent tests and the full `go test ./...`,
`go vet ./...`, `go build ./cmd/omaforge`, and `git diff --check` pass. A real
CLI walkthrough used a disposable Markdown reference, reached the confirmed
reference-driven review with the configured agent, then cancelled before any
project generation or agent launch; its temporary files were removed. Earlier
checks in this development checkpoint also covered the release action,
installer, website dependency, catalog, installer route, typecheck, official
validator, and actionable live-demo failure. This work remains local,
uncommitted, and unreleased pending review.

A subsequent supervised agent trial implemented a reviewed local Git-status
specification in a disposable generated project. Static review caught a
literal NUL byte that made one QML file appear binary even though the existing
project, Forge, and official checks passed. Forge now reports error rule
`OF305`, including the source path and line, whenever a checked QML, JavaScript,
or shell source file contains a NUL byte. The regression test uses inert file
bytes and does not execute plugin QML.

Control-center dogfooding then exposed a visual-fidelity gap: an implementation
could inventory reference behavior yet dismiss a supplied raster icon as
branding-only and approximate a UI mockup with generic generated composition.
The unreleased guided contract now makes design images fidelity targets,
requires supplied logos/icons in their intended locations, prohibits inventing
a vector-conversion prerequisite, and requires human reference-versus-render
review for visual acceptance.

The website dependency baseline was rebuilt cleanly on 2026-08-23. A narrow
checked-in pnpm patch makes Geistdocs' pinned `fumadocs-ui@16.2.2` sidebar
exports statically visible to Next.js, and local route adapters satisfy the
current Next.js 16 route contract. Type checking and a complete 42-route
production build pass with Webpack. Turbopack cannot be verified in the current
restricted execution environment because its MDX loader attempts to bind a
local process port; the failure is environmental rather than a missing content
file.

## Installer and repository protection checkpoint

The canonical website now owns one explicit `install.sh` path for first install
and later user-requested updates. Its isolated test builds local release
fixtures and verifies fresh installation, current-version no-op behavior,
upgrade replacement, invalid-version rejection, checksum failure without
clobbering the working binary, and symlink-target refusal. The public docs keep
the full exact-version block and add short, inspect-first, and pinned-version
options. The script does not run automatically or weaken Forge's normal
network-independent behavior.

Production verification initially found that the Geistdocs locale proxy
intercepted the otherwise deployed `/install.sh` static asset and returned 404.
The proxy now explicitly bypasses that canonical public path, and website CI
checks the executable asset, strict Bash header, proxy exclusion, and Quickstart
URL together before building.

GitHub `main` protection now requires pull requests, an up-to-date branch, the
`test` and `website` status checks, and resolved review conversations. The rules
apply to administrators and prohibit force pushes and branch deletion. Required
approvals remain zero while the repository has one maintainer, so the owner is
not locked out by an impossible self-approval requirement.

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

1. No owner-catalog implementation work remains. Keep both public releases and
   repository protections healthy. The documentation website and custom domain
   are live at `https://www.omarchyforge.com`; `/projects`, the exact Omaudit
   install command, its 753×544 local preview and optimized image route, and the
   sitemap entry were verified after the production deployment.
2. The hands-on `v0.3.0` onboarding and supervised agent-ready implementation
   trial are complete. Keep future agent-produced plugin execution subject to
   explicit human review and the existing trust acknowledgement. Forge does
   not automatically install the separate Handoff plugin.
3. The user's ongoing local installation is now the verified `v0.4.0` release.
   Marketplace submission and announcement remain distinct approval-gated work.
4. Review the complete guided-agent working-tree diff, rerun the release-sized
   baseline if accepted, and commit it locally before opening any pull request
   or beginning release work.
5. Control Center is public with the `omaforge-project` topic and stable
   `v0.1.0` release at
   `cbf45874c2a344c7fba78837eb19c91343dbaf51`. Forge's generated catalog now
   includes its checked-in 749×1102 preview and derived official plugin install
   command. This catalog update remains local until its Forge commit is pushed.
