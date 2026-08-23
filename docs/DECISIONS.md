# Decision record

## D-001: Go and the standard library

Status: accepted, 2026-08-22.

The CLI uses Go and initially has no third-party Go dependencies. This supports
a small native binary, deterministic builds, and direct unit testing without a
framework. A dependency may be introduced later only for a concrete correctness
or maintainability benefit.

## D-002: executable name

Status: accepted, 2026-08-22.

The executable is `omaforge`. The shorter name avoids collisions associated
with the generic executable name `forge` while retaining the product identity.

## D-003: honest development metadata

Status: accepted, 2026-08-22.

Local builds report version `dev`, with commit and build date as `unknown` when
not injected. Release automation may set those values with linker flags later.

## D-004: upstream validation boundary

Status: accepted, 2026-08-22.

Forge will complement, not replace, `omarchy plugin validate`. Future local
diagnostics should call the official validator when available. Deterministic
CI rules must identify their upstream capability snapshot and distinguish
structural parity checks from Forge quality rules.

## D-005: CI action references

Status: accepted, 2026-08-22.

Foundation CI uses only GitHub-maintained actions pinned to stable major
versions. Workflow permissions are explicitly read-only. Before a public or
release-oriented security boundary is introduced, action references should be
reconsidered for immutable commit pinning and automated update handling.

## D-006: embedded, deterministic templates

Status: accepted, 2026-08-22.

Plugin template files are embedded in the Go binary and rendered in sorted path
order. User values never determine generated filenames. JSON and QML strings
use JSON encoding, while prose values are escaped for Markdown. A golden tree
records each output path, mode, and content digest.

## D-007: overwrite model

Status: accepted, 2026-08-22.

The default is to refuse any nonempty target. `--force` preserves unrelated
files and overwrites only paths owned by the selected template after printing a
complete plan. Symlinks anywhere in the target are rejected before forced
writes. `--dry-run` performs the same validation and rendering without writing.

## D-008: Milestone 1 upstream pins

Status: accepted, 2026-08-22.

Template design was checked against official Omarchy `quattro` revision
`ed7bae4ac5a570e9df307486e0202fdafcc6ee24` and the official Basecamp plugin
revision `abc1ba72aaf47db530d2a0c6901d99f0f98e6aa7`. The installed Omarchy package
remains authoritative for local validation and runtime contract checks.

## D-009: Versioned findings and exit codes

Status: accepted, 2026-08-22.

Check reports use schema version 1 and stable rule IDs. Rules explicitly name
their source as `official-parity` or `forge`; heuristic Forge rules say so in
their explanation. Findings are sorted independently of filesystem traversal.
Text, JSON, and SARIF serialize the same report.

`omaforge check` returns 0 when there are no error-severity findings, including
when warnings exist; 1 for project errors or report-write failures; and 2 for
invalid usage. This allows publish-readiness warnings to remain visible without
making early CI adoption unexpectedly blocking.

## D-010: Static check and local doctor boundary

Status: accepted, 2026-08-22.

`check` is reproducible, noninteractive, network-free, and never runs QML or
external validators. `doctor` is intentionally environment-aware and may run
read-only installed commands with bounded timeouts, including the authoritative
Omarchy validator. Both share Forge's project rule report.

## D-011: Exact-version Action installation

Status: accepted, 2026-08-22.

The composite Action requires an exact semantic Forge version and rejects
moving values such as `latest`. It downloads the architecture-specific archive
and release checksum file, verifies SHA-256 before extraction, and invokes only
the released `omaforge check` command. Workflow inputs enter shell scripts
through quoted environment variables rather than expression interpolation.

GitHub annotations are derived from the schema-v1 JSON report. Workflow command
and property metacharacters are escaped before output. The Action does not run
plugin QML and has no package-install, telemetry, account, or network behavior
beyond downloading the explicitly selected release assets.

## D-012: Tag-only, least-privilege releases

Status: accepted, 2026-08-22.

One reviewed shell script owns local and CI release packaging. It produces
normalized static Linux `amd64` and `arm64` archives and a checksum manifest,
refusing nonempty output directories. The release workflow runs only for exact
`vX.Y.Z` tags. Its build job has read-only contents access; only the dependent
publication job receives `contents: write`.

Third-party release actions are avoided. Official GitHub Actions are pinned to
immutable commit SHAs, and GitHub CLI creates the release from the already
verified workflow artifact. Tags and releases remain explicit external actions
and are not created by normal CI or local packaging.

## D-013: Static-first Geistdocs website boundary

Status: accepted, 2026-08-22.

The documentation website uses Next.js 16, TypeScript, MDX, and the
Geistdocs/Fumadocs rendering foundation in `website/`. The package-backed
structure provides navigation, source-backed search, raw Markdown routes, and
edit links without creating a parallel documentation engine.

The stock initializer's AI chat, feedback submission, analytics, and Markdown
request tracking are explicitly removed or disabled. Next.js telemetry is also
disabled in every package script. The site has no account, authentication,
database, telemetry, or user-data persistence. Only `esbuild` is explicitly
allowed to run an install lifecycle script in pnpm's workspace policy.

The planned `omarchyforge.com` origin is not hard-coded as a deployed fact. A
canonical URL is supplied through `NEXT_PUBLIC_SITE_URL` when deployment is
separately authorized; local builds fall back to localhost. Building the site
does not create a Vercel project, deployment, domain, or environment variable.

## D-014: Separate, local-first Handoff repository

Status: accepted, 2026-08-22.

The flagship Handoff plugin began in the separate private repository
`omarchy-forge/handoff`. It has an independent installation and eventual
release boundary and does not belong inside the Forge monorepo.

The MVP is one Omarchy `bar-widget` with a native popout. It invokes Git and
the terminal through array-form process arguments, never interpolated shell
commands. It stores canonical project paths, one note, and Git metadata in an
atomically written versioned JSON file under the user's XDG data directory.
It does not execute pinned-project code, hooks, build scripts, or QML and has
no server, account, API key, telemetry, analytics, or required network access.

Repository implementation and private CI do not imply release. Tags, GitHub
releases, marketplace submission, installation for ongoing use, visibility
changes, and announcements remain separate explicit actions.

## D-015: Git-only scaffold targets

Status: accepted, 2026-08-22.

`omaforge init` treats an existing directory containing exactly one real
`.git` directory as an empty scaffold target. This supports the normal workflow
of creating or cloning an empty repository before generating its plugin files.

A `.git` file, symlink, or any additional directory entry still makes the
target nonempty and requires `--force`. Normal generation does not write inside
existing Git metadata; `--git` remains the explicit option for running Git
initialization when requested.

## D-016: Explicitly trusted isolated runtime harness

Status: accepted, 2026-08-22.

Generated bar-widget projects include `tests/runtime`, an optional Omarchy-host
entry-point smoke test. It requires the literal `--trust-plugin-code` argument
because QML is executable code. It runs a copy of the plugin in temporary HOME
and XDG directories, links only to the installed read-only Omarchy QML contract,
checks lifecycle methods and finite geometry, and cleans up a dedicated process
group after a bounded run.

The harness is never called by static checks, generated CI, or the normal test
script. This preserves the rule that untrusted pull-request QML is not executed
automatically while giving a developer an explicit local runtime check after
reviewing their own code.

## D-017: Defer the service-plus-widget template

Status: accepted, 2026-08-22.

Forge will not add a service-plus-widget template from the single installed
first-party example. The service lifecycle is distinct and valuable for shared
multi-monitor state, but the widget-to-service access path is not documented as
a stable third-party contract and no second independent plugin proves a common
shape.

The evaluation in `SERVICE_WIDGET_EVALUATION.md` records the evidence and
revisit criteria. Until one criterion is met, the supported bar-widget template
remains intentionally singular.

## D-018: Public Forge, initially private Handoff

Status: accepted, 2026-08-22.

The `omarchy-forge/forge` repository and its `v0.1.0` release are public. The
separate `omarchy-forge/handoff` repository and release remain private.

Forge visibility does not authorize a website deployment, connection of
`omarchyforge.com`, marketplace submission, Handoff visibility change, or a
public launch announcement. Those remain separate approval-gated actions.

Handoff's visibility boundary was later changed with explicit approval on
2026-08-23; D-026 records its public catalog boundary.

## D-019: Project-owned trusted development harness

Status: accepted, 2026-08-22.

The first `omaforge dev` increment delegates to the generated project's
versioned `tests/runtime` harness. Forge requires the literal
`--trust-plugin-code` acknowledgement, verifies that the project, manifest, and
harness exist, rejects a symlinked harness, and invokes it without a shell-built
command string.

The harness retains responsibility for temporary HOME/XDG isolation, read-only
Omarchy module links, process-group cleanup, and its timeout. This one-shot
runtime check does not install or enable the plugin, mutate shell configuration,
watch files, capture screenshots, or imply that arbitrary plugin QML is safe.

## D-020: Fixed isolated demo-state vocabulary

Status: accepted, 2026-08-22.

`omaforge dev` accepts only `ready`, `empty`, or `error` through `--state`.
Forge passes the validated value to the project harness as fixed arguments; the
harness exposes it to its temporary QML host and requires the plugin's
`setDemoState` method to accept it before reporting success.

State-driven runs require only the common panel lifecycle and `setDemoState`
contract; they do not impose the starter template's public `refresh` method on
independently implemented plugins. The no-state smoke run retains the refresh
check.

States are fictional and in-memory. They do not read fixture paths supplied by
the caller, write plugin data, call live-shell IPC, or broaden the explicit QML
trust boundary. The short settling interval validates each state contract but
is not yet an interactive watcher or screenshot session.

## D-021: Template-declared plugin-only screenshots

Status: accepted, 2026-08-22.

Screenshot capture uses Qt Quick `grabToImage` on a plugin-declared
`forgeScreenshotTarget` after applying a fictional state in the isolated
runtime. Forge does not invoke compositor, monitor, window, or region capture,
so unrelated desktop pixels cannot enter the result.

The command requires explicit QML trust, one of the fixed demo states, and a
new `.png` output path. Existing files are never overwritten. Plugins without
an explicit capture target fail closed; Forge does not guess at a screen region.

## D-022: Polling watch sessions over fresh isolated runs

Status: accepted, 2026-08-22.

`omaforge dev --watch` computes a deterministic fingerprint of sorted regular
project files while excluding `.git`. A bounded local poll detects changes and
reruns the project-owned harness in a fresh temporary runtime. No filesystem
watch dependency or persistent live-shell installation is introduced.

One explicit trust acknowledgement covers the session. Failed QML runs are
reported without terminating the watcher so the developer can repair source;
Ctrl-C or termination ends the session cleanly. The watcher performs no network
access and does not write inside the project.

## D-023: Vercel-hosted documentation checkpoint

Status: accepted, 2026-08-22.

The documentation website is deployed from the verified Forge `main` branch to
the dedicated Vercel project `omarchy-forge-docs`. Its canonical production
origin is `https://www.omarchyforge.com`, with `https://omarchyforge.com`
redirecting to it. The production `NEXT_PUBLIC_SITE_URL` setting ensures
generated metadata and machine-readable discovery files use that origin.
The deployment retains the static-first, no-account, no-database, no-analytics,
and no-AI-chat boundaries recorded in D-013.

At this checkpoint, the site deployment did not announce the site, submit
anything to a marketplace, or change Handoff's then-private visibility. Those
were separate approval-gated actions.

## D-024: Provider-neutral agent-ready scaffolding

Status: accepted, 2026-08-22.

Forge integrates with Omarchy's first-class coding-agent workflow by generating
a provider-neutral specification and project-scoped agent instructions. The
opt-in `omaforge init --agent-ready` path adds `FORGE_SPEC.md` and `AGENTS.md`
through the existing deterministic template renderer. It does not embed a model
client, choose a provider, collect credentials, make a network request, or
automatically launch an agent.

The generated contract permits agents to edit product-specific code and run
static Forge and official validation. It prohibits QML execution, plugin
installation, live-shell mutation, privileged operations, and unapproved
expansion into networking, persistence, authentication, or secret handling.
The specification starts with an explicit `Draft` status; the generated agent
must not implement it until a human changes that status to
`Ready for implementation` after resolving and reviewing every placeholder.
Because instructions are not a sandbox and Omarchy prompt tasks may run
unattended, a human must complete the specification, commit a rollback point,
review the resulting diff, and explicitly decide whether to trust and execute
the plugin.

## D-025: One explicit checksum-verifying install and update path

Status: accepted, 2026-08-23.

Forge publishes one readable Bash script at the canonical documentation origin
for both first installation and later user-requested updates. With no version
argument it follows GitHub's latest-release redirect, validates the resulting
exact semantic-version tag, and compares it with the installed binary before
downloading. `--version` retains immutable release pinning for reproducible
installation.

The script selects only a supported Linux architecture, verifies the selected
archive against the published SHA-256 manifest, stages the binary in the user
installation directory, and replaces only `~/.local/bin/omaforge` by default.
It never uses privilege escalation, a package manager, shell configuration,
telemetry, scheduled execution, or background network checks. Documentation
keeps an inspect-first path and the full manual exact-version block alongside
the short convenience command because piping a remote script requires explicit
user trust.

## D-026: Checked-in owner-project catalog with protected automated updates

Status: accepted, 2026-08-23.

The Forge website exposes a curated directory limited to public, active,
non-fork repositories owned by `omarchy-forge` and explicitly marked with the
`omaforge-project` topic. Each repository supplies a strict, versioned metadata
document, but Forge derives repository, release, preview, and installation
targets from verified GitHub identity and version data. A matching stable
release is required. Plugins use validated Forge manifests and the official
Omarchy install command. CLI tools use validated PEP 621 metadata and a bounded,
strict-Bash installer whose repository identity must match; Forge derives an
immutable release-tag command and never executes the script. The initial
eligible plugin is Handoff, followed by the Omaudit CLI security tool.

The catalog and preview images are generated artifacts committed to Forge, so
ordinary website builds remain deterministic and network-independent. An
explicit daily or manually dispatched GitHub workflow reads only bounded
metadata, installer text, and image data, opens a pull request when artifacts change, invokes protected
CI, and requests auto-merge. Neither the generator nor CI installs projects or
executes project QML. This is an owner-maintained project directory, not a
third-party marketplace or submission system.

## D-027: Confirmed guided handoff to the configured Omarchy agent

Status: accepted, 2026-08-23.

Forge keeps `--agent-ready` as deterministic scaffold-only generation and adds
a separate interactive `omaforge init --agent` path for users who explicitly
want orchestration. The guided path reads but never changes the current Omarchy
agent, gathers product and access boundaries, renders an editable summary, and
writes nothing unless the user confirms that summary and launch.

Confirmation creates a complete ready specification, project safety contract,
and initial implementation prompt, initializes a local Git repository, and
commits the generated baseline before invoking `omarchy agent prompt` from the
project directory. Forge does not embed a provider, handle credentials, push a
remote, or install or execute the plugin. Noninteractive, force, and dry-run
combinations are rejected to keep the confirmation meaningful.

The baseline commit disables repository hooks and commit signing for that one
generated commit so an interactive signing prompt or inherited hook cannot
silently extend the guided workflow. It still uses the user's configured Git
identity and fails with recovery guidance if that identity is unavailable.

The installed Omarchy launcher runs configured agents with unattended
permissions. Forge states that fact before confirmation and treats generated
instructions as defense in depth, not a sandbox. The implementation agent may
edit, run static checks, and commit locally. Its required final handoff orders
human diff review, isolated runtime states, deliberate local installation,
installed-plugin demos, cleanup, and optional remote publication without
claiming that static validation proves runtime safety.

Guided projects may include up to ten explicitly selected local references.
Forge accepts bounded UTF-8 `.txt`/`.md` files, signature-checked PNG, JPEG,
or WebP files, and XML-validated SVG files without rendering them; it creates a project `references/` drop directory, rejects symlinks
and duplicate paths, captures the selected bytes, and records SHA-256 digests.
Reference content is supporting, untrusted product input and cannot override
the spec or agent rules. Forge performs no upload or OCR, but confirmation
warns that the configured external agent/provider may receive the copied
content.

The guided model offers two explicit product-input paths. Reference mode makes
confirmed text and image files the primary product brief, requires the agent to
inventory every function and presentation requirement, and asks the user only
for access boundaries that untrusted files cannot authorize. Questionnaire mode
retains explicit bar, popout, action, data-source, and command questions. Both
paths prohibit substituting generic in-memory placeholder content for requested
behavior.

Reference-driven presentation requirements are fidelity targets rather than
optional inspiration. UI mockups require a complete visual inventory and close
composition match within supported Omarchy APIs; supplied logos and icons must
be used in their intended locations unless the user explicitly marks them as
inspiration-only. Raster branding does not require vector recreation. Material
deviations require user approval, and a successful static/runtime check does not
replace human reference-versus-screenshot review.

## D-028: Patch the Geistdocs-pinned Fumadocs UI export shape

Status: accepted, 2026-08-23.

The website retains Geistdocs' exact `fumadocs-ui@16.2.2` API instead of
overriding it to a newer UI release. Newer UI releases change provider and
internal import contracts that Geistdocs `1.22.0` still depends on. The pinned
UI release, however, publishes several sidebar bindings through a destructured
export that Next.js production bundlers cannot resolve statically.

Forge carries a narrow pnpm package patch that rewrites only those bindings as
explicit exports. Route adapters also expose Next.js 16-compatible handler and
static-parameter signatures. This keeps the compatibility exception reviewed,
checked in, reproducible from the lockfile, and removable when Geistdocs ships
a compatible dependency update.
