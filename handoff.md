# Omarchy Forge — Codex CLI Handoff

> Build, test, and ship better Omarchy plugins.

## 1. Handoff purpose

This document is the durable product, architecture, safety, and execution handoff for building **Omarchy Forge** with Codex CLI on an actual Omarchy machine.

It is intentionally self-contained. Treat it as the project source of truth unless the repository later contains a newer explicit decision record. Do not rely on assumptions from model memory about Omarchy, Quickshell, GitHub, or the local machine. Inspect the installed system and official upstream source before implementing behavior that depends on them.

## 2. Current project status

- Product name: **Omarchy Forge**
- Tagline: **Build, test, and ship Omarchy plugins.**
- GitHub organization: `omarchy-forge`
- Main repository: `omarchy-forge/forge`
- Domain: `omarchyforge.com`
- Repository status at handoff: **private and empty**
- Intended license: MIT
- Initial development environment: the owner's Omarchy machine
- Initial audience: developers creating Omarchy Quattro shell plugins
- Relationship to Omarchy: independent community project; not official, affiliated with, sponsored by, or endorsed by Omarchy, Basecamp, 37signals, or DHH

The repository must remain private until it contains a deliberate foundation and one working, tested artifact. Do not make it public, change GitHub visibility, publish packages, register marketplace entries, deploy production services, or announce the project without explicit user approval.

## 3. Product thesis

Omarchy already has an operating system, documentation, community, plugin directory, marketplace, and installation commands. Forge must not duplicate them.

Forge's job is to take a plugin developer from:

```text
empty directory
    -> valid scaffold
    -> productive local development
    -> meaningful diagnostics
    -> repeatable CI
    -> publish-ready plugin
```

The narrow initial promise is:

> **The fastest reliable way to build a polished Omarchy bar widget.**

The broader long-term promise is:

> **The open-source developer toolchain for creating, testing, validating, previewing, and shipping Omarchy plugins.**

The CLI provides the recurring value. The website is the documentation, reference, and visual home for the tools.

## 4. Product boundaries

### Forge is

- A local-first CLI.
- A set of high-quality plugin starter templates.
- An extended diagnostic and release-readiness checker.
- A deterministic CI checker and GitHub Action.
- Version-aware Omarchy plugin documentation.
- Eventually, a local demo, preview, and screenshot workflow.
- An open-source project other plugin authors can inspect and contribute to.

### Forge is not

- Another marketplace or plugin directory.
- Another `awesome-omarchy` list.
- A theme gallery.
- A Discord, forum, social feed, or community replacement.
- A replacement for `omarchy plugin add`, `omarchy plugin validate`, or other official commands.
- A hosted service that executes arbitrary third-party QML.
- A cloud dashboard requiring accounts.
- A billing or subscription product in the initial roadmap.
- A generic Git GUI or project-management application.
- A collection of unrelated Omarchy plugins inside the main repository.

Avoid scope creep toward marketplace, discovery, user profiles, submissions, likes, comments, or hosted execution.

## 5. Critical upstream facts to verify locally

The following reflects the official Omarchy Quattro plugin model at planning time, but the local installation and current upstream repository are authoritative.

Expected current behavior:

- Omarchy runs one long-lived Quickshell process named `omarchy-shell`.
- First-party and third-party components are loaded as plugins.
- Third-party plugins live under `~/.config/omarchy/plugins/<id>/`.
- A third-party plugin is a Git repository with `manifest.json` at its root.
- Plugin manifest `schemaVersion` is currently `1`.
- Supported kinds currently include:
  - `bar-widget`
  - `panel`
  - `overlay`
  - `menu`
  - `service`
  - `bar`
- A plugin can declare more than one kind.
- Third-party IDs must be namespaced; `omarchy.*` is reserved.
- Each declared kind requires its corresponding entry point.
- Entry points must be safe relative paths and exist.
- Symlinks are rejected inside plugin folders.
- Bar-widget metadata can include `displayName`, `category`, `defaultSection`, `allowMultiple`, defaults, and settings schema.
- Omarchy already provides `omarchy plugin validate <folder>`.
- `omarchy plugin add <local-path> --enable` is an official local-development path.
- Plugins execute as unsandboxed code inside the user's long-lived shell process.
- The shell exposes theme and structural tokens through QML singletons such as `Color`, `Style`, and `Border`.
- The shell exposes IPC through `omarchy-shell` and plugin-specific targets.
- User customization is persisted in `~/.config/omarchy/shell.json`.
- Saving files beneath the user plugin directory may trigger hot reload, but runtime behavior is young and can change.

Relevant official references:

- `https://github.com/basecamp/omarchy/blob/quattro/manual/32-shell-plugins.md`
- `https://github.com/basecamp/omarchy/blob/quattro/shell/README.md`
- `https://github.com/basecamp/omarchy/blob/quattro/shell/plugins/README.md`
- `https://github.com/basecamp/omarchy/blob/quattro/bin/omarchy-plugin-validate`
- `https://github.com/basecamp/omarchy-basecamp-plugin`

### Mandatory discovery before implementation

On the Omarchy machine, inspect and record at least:

```bash
omarchy --version || true
command -v omarchy
command -v omarchy-shell
command -v quickshell
command -v qmllint
printf '%s\n' "${OMARCHY_PATH:-}"
omarchy plugin list --json
omarchy-shell shell ping
```

Also inspect:

- The actual `$OMARCHY_PATH` and current Git revision/version when available.
- `shell/README.md`.
- `docs/omarchy-shell.md` if present.
- `shell/services/PluginRegistry.qml`.
- `bin/omarchy-plugin-validate`.
- `shell/plugins/README.md`.
- Current first-party manifests and at least one implementation of every relevant plugin kind.
- Current bar-widget implementation and settings contract.
- Current theme tokens and common components.
- Current tests for plugin validation, registry loading, configuration, and hot reload.
- The official Basecamp plugin as a rich third-party example.

Record the findings in `docs/UPSTREAM_COMPATIBILITY.md`. Include the exact Omarchy version, Git commit when available, paths inspected, commands executed, and any differences from this handoff.

Do not copy assumptions into templates when the installed implementation can be inspected.

## 6. Key differentiation: installable versus publish-ready

Omarchy already provides structural validation. Forge must not market a duplicate validator as its primary innovation.

Use this distinction:

```text
Official Omarchy validation:
Can the shell safely recognize and load this plugin structure?

Forge checking:
Is this plugin polished, documented, testable, compatible,
maintainable, and ready for other people to install?
```

Whenever the official validator is available, Forge should call or compare against it rather than silently replacing it. Forge may maintain its own deterministic schema/rules for CI environments, but those rules must be traceable to a known upstream version and tested for parity.

## 7. Naming decision

The product, organization, domain, and repository keep the **Omarchy Forge** name.

Recommended executable and package name:

```text
omaforge
```

Recommended commands:

```bash
omaforge init
omaforge doctor
omaforge check
omaforge version
```

Reason: the generic binary name `forge` collides with many established developer tools, and public Omarchy repositories already use `omarchy-forge` in the context of Laravel Forge. Confirm package-name availability before publishing, but do not let that block local development.

## 8. Repository strategy

Start with one monorepo:

```text
omarchy-forge/forge
```

It initially owns:

- Go CLI.
- Shared rule engine.
- Plugin templates.
- Compatibility fixtures.
- Tests.
- GitHub Action.
- Documentation.
- Early website source.

Split a component into another repository only when it has an independent installation path, release cycle, contributor group, or security boundary.

Examples of future separate repositories:

- `omarchy-forge/handoff` — flagship user-facing plugin.
- Editor extension — only when independently releasable.
- A substantial website or hosted service — only if the monorepo becomes a burden.

## 9. Target repository structure

Use a structure close to this, adjusting only when the chosen Go and website tooling requires it:

```text
forge/
├── README.md
├── LICENSE
├── CHANGELOG.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── SECURITY.md
├── AGENTS.md
├── .editorconfig
├── .gitignore
├── go.mod
│
├── cmd/
│   └── omaforge/
│       └── main.go
│
├── internal/
│   ├── cli/
│   ├── scaffold/
│   ├── check/
│   ├── doctor/
│   ├── manifest/
│   ├── upstream/
│   └── report/
│
├── templates/
│   └── bar-widget/
│       ├── template.yaml
│       └── files/
│
├── schemas/
│   └── manifest/
│       └── v1.schema.json
│
├── compat/
│   └── omarchy-4/
│       ├── capabilities.json
│       └── fixtures/
│
├── action/
│   ├── action.yml
│   └── README.md
│
├── apps/
│   └── web/
│
├── examples/
│   └── README.md
│
├── docs/
│   ├── PRODUCT.md
│   ├── ROADMAP.md
│   ├── ARCHITECTURE.md
│   ├── DECISIONS.md
│   ├── UPSTREAM_COMPATIBILITY.md
│   ├── SECURITY_MODEL.md
│   └── RELEASES.md
│
├── test/
│   ├── fixtures/
│   └── golden/
│
└── .github/
    ├── ISSUE_TEMPLATE/
    ├── PULL_REQUEST_TEMPLATE.md
    └── workflows/
        ├── ci.yml
        └── release.yml
```

Do not create empty placeholder directories just to mimic this diagram. Add a directory when it receives real content.

## 10. Technology decisions

### CLI: Go

Use Go unless local discovery reveals a compelling technical blocker.

Reasons:

- Produces one native Linux binary.
- No Node, Python, or Rust toolchain required for end users.
- Straightforward cross-compilation for `linux/amd64` and `linux/arm64`.
- Good standard-library support for files, JSON, templates, processes, and testing.
- Fast startup and predictable distribution.

Prefer the Go standard library initially. Avoid a large CLI framework until command complexity proves it necessary. A small dependency is acceptable when it materially improves correctness, but every dependency must be justified.

### Website: Next.js + TypeScript + MDX

The website should be documentation-first and mostly static. Use Next.js and MDX because they align with the owner's experience and can live cleanly in the monorepo.

Do not begin the website until the CLI produces a real plugin scaffold.

Do not introduce a database, Supabase, authentication, or server-side user data for the MVP.

### GitHub Action

Keep the Action in `action/` initially. It should download or invoke a pinned Forge release and run deterministic checks. It should not execute arbitrary plugin QML.

## 11. CLI contract

### `omaforge version`

Output the Forge version, commit, build date when embedded, supported manifest schemas, and optionally detected Omarchy version.

It must be deterministic and script-friendly.

### `omaforge init`

Purpose: generate a clean, valid, publishable starting plugin.

Interactive example:

```text
Plugin name: Project Pulse
Plugin ID: eddie.project-pulse
Template: Bar widget with popout
Description: Your active project at a glance
Author: Eddie Ortega
License: MIT
Default section: right
Initialize Git repository: yes
Include Forge CI: yes
```

Noninteractive example:

```bash
omaforge init project-pulse \
  --id eddie.project-pulse \
  --template bar-widget \
  --section right \
  --author "Eddie Ortega" \
  --license MIT
```

Initial requirements:

- Validate names, IDs, output paths, and collisions before writing.
- Never overwrite existing files without explicit `--force` plus a clear preview.
- Refuse dangerous targets such as `/`, the user's home directory, `$OMARCHY_PATH`, or an existing nonempty directory unless the intended semantics are unambiguous.
- Support `--dry-run`.
- Support fully noninteractive use.
- Produce deterministic output for the same inputs.
- Escape user-provided values safely in JSON, QML, Markdown, YAML, and shell contexts.
- Generate a plugin that passes the official validator.
- Generate honest documentation; do not claim features that are not in the template.
- Optionally initialize Git, but never create a remote, push, publish, or make a repository public.

### `omaforge doctor`

Purpose: explain local environment and project problems in human-readable form.

It should eventually cover:

- Omarchy detection and version.
- Quickshell detection.
- Shell health/IPC availability.
- Official validator availability and result.
- QML tooling availability.
- Manifest and entry-point checks.
- Theme integration.
- Keyboard and focus behavior heuristics.
- Loading, empty, and error states.
- External command timeouts.
- Unsafe command construction.
- Writes into Omarchy-owned paths.
- Use of `sudo`, installer hooks, or unexpected package installation.
- Missing README, license, preview, requirements, privacy, update, and removal sections.
- Compatibility declaration.

Every rule needs:

- Stable rule ID.
- Severity.
- Clear message.
- Evidence location when available.
- Suggested remediation.
- Documentation link or explanation.

### `omaforge check`

Purpose: deterministic noninteractive CI checking.

Expected modes:

```bash
omaforge check .
omaforge check . --format text
omaforge check . --format json
omaforge check . --format sarif
```

Requirements:

- Stable exit codes.
- No prompts.
- No network by default.
- No QML execution.
- Reproducible results.
- Machine-readable output with schema/version.
- Ability to pin a target Omarchy compatibility version.

### Deferred commands

Do not implement these in the first milestone:

- `omaforge dev`
- `omaforge preview`
- `omaforge screenshot`
- `omaforge publish`
- `omaforge gallery`
- `omaforge login`

They remain roadmap ideas, not promises.

## 12. First template: bar widget with popout

Build one excellent template before adding more kinds.

The generated plugin should include approximately:

```text
my-plugin/
├── manifest.json
├── README.md
├── LICENSE
├── CHANGELOG.md
├── Widget.qml
├── components/
│   ├── DetailsPopout.qml
│   ├── EmptyState.qml
│   ├── ErrorState.qml
│   └── LoadingState.qml
├── services/
│   └── DataService.qml
├── assets/
│   └── icon.svg
├── demo/
│   └── run
├── tests/
│   └── run
└── .github/
    └── workflows/
        └── forge.yml
```

Adjust exact QML filenames and component boundaries to match the verified official runtime.

The template must demonstrate:

- Official manifest schema and current entry-point contract.
- Theme-aware colors, typography, spacing, surfaces, and borders.
- Top, bottom, left, and right bar orientation where the runtime supports it.
- Mouse and keyboard operation.
- Visible keyboard focus.
- Loading, success, empty, and error states.
- Async local command handling without blocking the UI.
- Timeouts and explicit error handling.
- User-configurable refresh interval or comparable setting.
- Safe quoting/argument passing.
- Clean reload and shutdown behavior.
- Local fictional demo data.
- Tests that do not require credentials or network access.
- Privacy/security documentation.
- Requirements, install, update, configuration, development, and removal instructions.

Use official components and contracts when available instead of creating parallel abstractions.

## 13. Template expansion order

After the first template is stable:

1. `service + bar-widget`
2. `panel`
3. `overlay`
4. `menu`
5. Full replacement `bar`

Do not create weak templates merely to advertise broader coverage.

## 14. Validation and quality rule categories

Forge checking should eventually include:

### Structural

- Valid JSON.
- Supported schema version.
- Required fields and types.
- Namespaced nonreserved ID.
- Known kinds only.
- Entry point for each kind.
- Safe relative paths.
- Existing files.
- No symlinks.
- Semver-compatible plugin version.

### QML/static

- QML parser/linter result when available.
- Imports resolvable for the target Omarchy environment.
- Expected injected properties and lifecycle functions.
- No obvious synchronous/blocking process patterns.
- No unsafe string interpolation into shell commands.
- Timer/polling bounds.
- Cleanup of processes, timers, temporary files, and connections.

### Omarchy integration

- Use of current theme tokens instead of hard-coded shell appearance.
- Correct orientation behavior.
- Correct settings schema and defaults.
- Correct plugin enable/load semantics.
- No direct editing of Omarchy-owned source or default configuration.
- Compatibility range recorded.

### UX/accessibility

- Keyboard reachability.
- Visible focus state.
- Useful tooltip or accessible description where applicable.
- Loading, empty, and error states.
- Long-label/overflow behavior.
- No essential meaning encoded only by color.

### Security/privacy

- No embedded secrets, tokens, keys, or credentials.
- No unexpected telemetry.
- No network access unless documented and essential.
- No `sudo` or package installation performed by plugin code.
- No installer hooks.
- No writes to Omarchy-owned directories.
- External commands and files documented.
- Credential use delegated to the owning CLI when possible.

### Publish readiness

- README.
- License.
- Requirements.
- Install/update/remove instructions.
- Configuration.
- Development/testing instructions.
- Privacy/security explanation.
- Preview image or explicit warning when absent.
- Compatibility declaration.
- Changelog/version consistency.

Heuristic checks must be labeled as heuristics. Do not present an inconclusive pattern match as proof of vulnerability.

## 15. Compatibility strategy

Omarchy's plugin system is new and moving quickly. Compatibility must be a first-class capability rather than a hard-coded afterthought.

Maintain:

- A machine-readable capability snapshot for each supported Omarchy release line.
- Fixtures copied or derived only when licensing permits, with source and commit recorded.
- Tests showing parity with the official validator for supported cases.
- A document explaining which checks come from upstream and which are Forge quality rules.

Prefer capability detection over vague version comparisons when possible.

When Forge cannot confidently determine compatibility, report `unknown`; do not silently claim support.

## 16. Security model

Omarchy plugins run as unsandboxed code in a privileged position relative to the user's desktop session. Treat safety as a core product feature.

Non-negotiable rules:

- No Forge telemetry by default.
- No user account or cloud requirement.
- No secret collection.
- No arbitrary remote QML execution on the website or CI.
- No automatic `sudo`.
- No automatic package installation.
- No installer hooks in generated plugins.
- No editing `$OMARCHY_PATH` or other Omarchy-owned files.
- No silent changes to `~/.config/omarchy/shell.json`.
- No destructive overwrite of an existing project.
- No publishing, GitHub visibility changes, marketplace submission, domain changes, or external announcement without explicit approval.
- Never run untrusted public pull-request code on the owner's persistent self-hosted Omarchy machine.

If real shell-runtime testing is later automated, use trusted branches or disposable/ephemeral environments.

Generated projects must ignore at least:

```gitignore
.env
.env.*
!.env.example
*.pem
*.key
```

Never commit credentials, tokens, SSH keys, registrar details, local personal paths, or screenshots containing private information.

## 17. Testing strategy

### CLI unit tests

- Input validation.
- Plugin ID validation.
- Path safety.
- Template rendering.
- Escaping across output formats.
- Exit codes.
- Report serialization.

### Golden tests

For known inputs, generate the full plugin tree and compare it with reviewed golden files. Provide a deliberate update mechanism for golden fixtures.

### Integration tests

- Generate into a temporary directory.
- Run internal manifest checks.
- Run official `omarchy plugin validate` when available.
- Run QML linting when available.
- Verify no file escapes the target directory.
- Verify an existing nonempty target is not overwritten.
- Verify `--dry-run` writes nothing.
- Verify deterministic output.

### Runtime smoke tests

On the actual Omarchy machine, manually or through a trusted local script:

- Install the generated plugin from a local checkout.
- Enable it.
- Confirm the shell stays healthy.
- Confirm the widget renders.
- Test keyboard and pointer interactions.
- Test horizontal and vertical bar positions where applicable.
- Test multiple themes and text/spacing scaling.
- Test loading, empty, and error demo states.
- Edit a file and observe reload behavior.
- Disable and remove the plugin.
- Confirm user configuration and shell state are preserved.

Document verified scenarios and failures. Do not claim runtime compatibility from static checks alone.

## 18. Git and GitHub workflow

The repository is empty, so the first bootstrap is an exception to the later PR-only workflow.

Recommended sequence:

1. Confirm the remote points to the private `omarchy-forge/forge` repository.
2. Create the initial repository foundation locally.
3. Review for secrets, personal paths, generated junk, and misleading claims.
4. Create the first commit on the default branch:

```text
chore: bootstrap Omarchy Forge
```

5. Push only after the user approves the diff and commit.
6. Use feature branches and pull requests for meaningful work after the default branch exists.

Recommended next branches:

```text
feat/bar-widget-scaffolder
test/generated-plugin
feat/plugin-doctor
feat/github-action
feat/website
```

Do not push, open a pull request, merge, create releases, or change repository settings unless the user explicitly asks.

## 19. Public-release gate

Keep the repository private until all of the following are true:

- README clearly explains the real current product.
- MIT license is present.
- Independent-project disclaimer is visible.
- `AGENTS.md`, architecture, roadmap, security, and contribution guidance exist.
- `omaforge --help` works.
- `omaforge version` works.
- `omaforge init demo-plugin` generates a real plugin.
- The generated plugin passes the official validator on the Omarchy machine.
- Tests and CI are green.
- No secrets or personal paths are committed.
- The README includes an honest terminal recording, screenshot, or exact sample output.
- Installation instructions have been tested.

The repository does not need `doctor`, the GitHub Action, Handoff, or the full website before becoming public.

## 20. Milestone roadmap

### Milestone 0 — Private repository foundation

Deliver:

- Core project documents.
- MIT license.
- Security and contribution policies.
- `AGENTS.md` for future agents.
- Go module and minimal CLI.
- `omaforge --help` and `omaforge version`.
- Basic unit test and CI workflow.
- Upstream discovery notes.

### Milestone 1 — Bar-widget scaffolder

Deliver:

- Interactive and noninteractive `omaforge init`.
- `--dry-run`.
- Safe path/collision behavior.
- Excellent bar-widget-with-popout template.
- Deterministic golden tests.
- Official validation integration test.
- Generated README, license, demo, tests, and CI.

### Milestone 2 — Check and Doctor

Deliver:

- Shared rule engine.
- Human-readable `doctor`.
- Deterministic `check`.
- Text and JSON reporting.
- SARIF when the rule model is stable.
- Clear distinction between official and Forge rules.

### Milestone 3 — GitHub Action and releases

Deliver:

- Reusable Action in `action/`.
- PR annotations.
- Versioned Linux release binaries.
- SHA256 checksums.
- Release documentation.

### Milestone 4 — Documentation website

Deliver:

- Homepage.
- Quickstart.
- Template documentation.
- Command reference.
- Plugin anatomy.
- Manifest reference.
- Compatibility page.
- Roadmap and contributing pages.
- Deployment to `omarchyforge.com` only with explicit approval.

### Milestone 5 — Handoff flagship plugin

Create a separate repository only after the scaffolder is useful.

Handoff MVP:

- Select or pin a Git project.
- Save one next-step note per project.
- Record branch, dirty state, last commit, and timestamp.
- Show the note when returning to the project.
- Open the project in a terminal.
- Store data locally, for example under `~/.local/share/omarchy-handoff/`.
- No server, API key, account, telemetry, or required network access.

Use Handoff to expose weaknesses in Forge templates and checks. Later expand it toward Project Pulse, dev-service health, and coding-agent awareness.

### Milestone 6 — Local development and visual tooling

Possible later features:

- `omaforge dev`.
- Trusted local demo harness.
- Local state simulation.
- Screenshot capture.
- API/component explorer.
- Pattern gallery.
- Compatibility-change monitoring.

Do not promise these in the initial README as shipped features.

## 21. Website plan

The initial website should contain:

```text
Home
Quickstart
Templates
Commands
Plugin anatomy
Manifest reference
Compatibility
Roadmap
Contributing
```

Recommended homepage copy:

```text
OMARCHY FORGE

From an empty directory to a publish-ready Omarchy plugin.

Scaffold polished plugins, catch problems early,
and validate every pull request.
```

Required footer disclaimer:

```text
Omarchy Forge is an independent community project.
It is not affiliated with or endorsed by Omarchy, Basecamp, 37signals, or DHH.
```

Use a distinct visual identity. Do not use the official Omarchy logo as the Forge logo.

## 22. First-build backlog

Prioritize in this order:

1. Private repository foundation.
2. Local upstream contract inspection.
3. Minimal Go CLI.
4. `omaforge init`.
5. Bar-widget-with-popout template.
6. Generated README/demo/tests.
7. Golden and official validation tests.
8. `omaforge doctor`.
9. `omaforge check`.
10. GitHub Action.
11. Release binaries.
12. Documentation website.
13. Handoff flagship plugin.
14. Service + bar-widget template.
15. Panel template.
16. Overlay template.
17. Local preview/demo tooling.
18. Component/API explorer.
19. Compatibility matrix and badges.
20. Editor integration only after real demand.

## 23. Future plugin ideas — not current Forge scope

The research also identified user-facing plugins that could later be built with Forge:

1. Handoff.
2. Project Pulse.
3. Dev Environment Doctor.
4. Agent Control Center.
5. Secret Sentinel.
6. Remote Pairing.
7. PR Landing Zone.
8. Release Checklist.
9. Media Maker Mode.
10. Laptop Health.

Only Handoff belongs in the near-term Forge roadmap. Do not build these inside the main Forge repository.

## 24. Codex working rules

- Begin every milestone with inspection, not implementation guesses.
- Read repository instructions, especially `AGENTS.md`, before editing.
- Preserve existing user changes and unrelated files.
- Use small, reviewable commits.
- Do not build several milestones simultaneously.
- Keep source files focused and testable.
- Add tests with behavior, not after the fact.
- Prefer deterministic local behavior.
- Keep the CLI functional without network access.
- Treat the installed Omarchy source and official upstream code as authoritative.
- Record important decisions in `docs/DECISIONS.md`.
- Record upstream compatibility facts in `docs/UPSTREAM_COMPATIBILITY.md`.
- Update `CHANGELOG.md` for user-visible behavior once releases begin.
- Do not claim a feature is available until it is implemented and verified.
- Do not modify Omarchy-owned files.
- Do not change GitHub visibility or repository settings.
- Do not publish, deploy, push, merge, or announce without explicit permission.
- Stop when a task requires secrets, external account changes, production deployment, destructive system modification, or meaningful scope expansion.

## 25. First Codex session objective

The first Codex CLI session should execute **Milestone 0 only**.

Expected work:

1. Inspect the empty/private repository and Git state.
2. Inspect the installed Omarchy environment and authoritative plugin contract.
3. Create the repository foundation files.
4. Write honest product, architecture, roadmap, security, upstream, and decision documents.
5. Create a minimal Go CLI with `--help` and `version`.
6. Add tests and a basic CI workflow.
7. Run formatting, tests, and any safe validation available.
8. Summarize the exact diff and verification.
9. Stop for user review.

Do not implement the complete scaffolder, doctor, Action, website, or Handoff in the first session.

Do not push or make the repository public.

## 26. Milestone 0 acceptance criteria

- Repository has a coherent, minimal foundation rather than empty promises.
- Product definition and non-goals are explicit.
- Independent-project disclaimer is present.
- Go project builds on the Omarchy machine.
- `omaforge --help` exits successfully.
- `omaforge version` exits successfully and reports a development version.
- Unit tests pass.
- Go formatting and vet/static checks used by the project pass.
- CI configuration is syntactically sound.
- Upstream inspection is documented with exact local evidence.
- No secret, credential, private screenshot, or personal absolute path is committed.
- No Omarchy-owned files are modified.
- No external action is taken without approval.

## 27. Completion reporting format

At the end of each Codex session, report:

```text
Outcome
Files changed
Architecture/decisions made
Commands and tests run
Verification results
Known limitations
Risks or upstream uncertainties
Recommended next step
Git status
```

If anything could not be verified, say so explicitly. Never present planned work as completed work.

