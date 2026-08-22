# Initial Codex CLI Prompt — Omarchy Forge

Copy and paste everything below into Codex CLI from the root of the local clone of the private `omarchy-forge/forge` repository on the Omarchy machine.

---

You are beginning implementation of **Omarchy Forge**, an independent open-source developer toolchain for building, testing, and shipping Omarchy plugins.

Read `handoff.md` completely before taking any project action. Treat it as the authoritative product, architecture, safety, scope, and roadmap handoff. If it conflicts with the actual installed Omarchy implementation, do not guess: inspect the local installation and official upstream source, document the difference, and follow the real current contract.

## Current situation

- GitHub organization: `omarchy-forge`
- Repository: `omarchy-forge/forge`
- The remote repository is private and currently empty.
- Product name: Omarchy Forge
- Recommended executable name: `omaforge`
- Domain reserved for later: `omarchyforge.com`
- The repository must remain private during this task.
- This task is **Milestone 0 only: private repository foundation**.

## Product positioning

Omarchy Forge is a CLI-first toolchain that will eventually provide:

- Plugin scaffolding.
- High-quality starter templates.
- Extended diagnostics and release-readiness checks.
- Deterministic CI checks and a GitHub Action.
- Documentation and compatibility references.
- Later, trusted local preview/demo/screenshot tooling.

It is not another marketplace, plugin directory, theme gallery, community site, plugin manager, cloud dashboard, or hosted QML execution service.

Omarchy already provides `omarchy plugin validate`. Forge must complement official tooling rather than duplicate or replace it.

## Task authorization

You are authorized to inspect the local repository and installed Omarchy environment, create and edit files inside this repository, run safe local formatting/build/test/read-only diagnostic commands, and initialize the local Git history if the repository is truly unborn.

You are **not** authorized to:

- Push commits or branches.
- Change repository visibility or GitHub settings.
- Open or merge pull requests.
- Publish a release or package.
- Deploy the website or change DNS/domain settings.
- Submit anything to a marketplace.
- Modify files owned by the installed Omarchy distribution.
- Modify `~/.config/omarchy/shell.json`.
- Install packages or use `sudo` without stopping and asking first.
- Add telemetry, analytics, accounts, authentication, billing, Supabase, or a database.
- Execute untrusted third-party QML.
- Begin later milestones such as the full scaffolder, doctor, GitHub Action, website, or Handoff plugin.

Do not make a commit unless I explicitly approve the final diff. Prepare the repository and stop for review.

## Required workflow

### 1. Inspect before editing

Inspect:

- `git status`, configured remotes, current/default branch state, and whether the repository is unborn.
- Any existing files, including hidden repository instructions.
- The installed Omarchy version and paths.
- The actual local plugin contract and examples.

Run safe discovery commands such as:

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

Resolve `$OMARCHY_PATH` from the actual system if the environment variable is not populated.

Inspect the relevant local/upstream files when present:

- `shell/README.md`
- `docs/omarchy-shell.md`
- `shell/services/PluginRegistry.qml`
- `bin/omarchy-plugin-validate`
- `shell/plugins/README.md`
- Current plugin manifests.
- Current bar-widget and service examples.
- Plugin validation and registry tests.
- Theme tokens and shared QML components.

Do not modify those upstream files.

### 2. Reconcile the handoff with reality

Create `docs/UPSTREAM_COMPATIBILITY.md` with:

- Inspection date.
- Installed Omarchy version.
- Git commit/revision when available.
- Paths and source files inspected.
- Commands executed.
- Confirmed manifest schema and plugin kinds.
- Required entry-point mapping.
- Plugin ID rules.
- Theme/settings/runtime contracts relevant to future templates.
- Official validator behavior.
- Local-development and hot-reload behavior.
- Differences from `handoff.md` or uncertainties that need later validation.

Use exact evidence from the machine. Do not invent version claims.

### 3. Create the repository foundation

Create only useful files with real content:

```text
README.md
LICENSE
CHANGELOG.md
CONTRIBUTING.md
CODE_OF_CONDUCT.md
SECURITY.md
AGENTS.md
.editorconfig
.gitignore
docs/PRODUCT.md
docs/ROADMAP.md
docs/ARCHITECTURE.md
docs/DECISIONS.md
docs/UPSTREAM_COMPATIBILITY.md
docs/SECURITY_MODEL.md
go.mod
cmd/omaforge/main.go
internal/cli/...
.github/workflows/ci.yml
.github/PULL_REQUEST_TEMPLATE.md
.github/ISSUE_TEMPLATE/...
```

Adjust the exact minimal layout when justified. Do not create empty placeholder directories or claim unimplemented features.

The README must clearly state:

```text
Omarchy Forge
Developer tools for building, testing, and shipping Omarchy plugins.
```

Include this visible disclaimer:

```text
Omarchy Forge is an independent community project. It is not affiliated with
or endorsed by Omarchy, Basecamp, 37signals, or DHH.
```

The README must label the project as early development and state that the first goal is a reliable plugin scaffold and validation workflow. Do not advertise `init`, `doctor`, the Action, website, or preview as already shipped.

Use the MIT license.

### 4. Create the minimal Go CLI

Implement only:

```bash
omaforge --help
omaforge version
```

Requirements:

- Use Go.
- Prefer the standard library.
- Keep command parsing small and testable.
- Development version output must be honest.
- Structure version/build metadata so releases can inject version, commit, and build date later.
- Unknown commands and invalid arguments must produce useful messages and nonzero exit codes.
- Add unit tests for the behavior implemented.

Do not implement `init`, `doctor`, `check`, preview, publishing, networking, or plugin generation during this task.

### 5. Add baseline CI

Add a simple GitHub Actions workflow suitable for the private repository that performs the checks supported by the current foundation, such as:

```text
go test ./...
go vet ./...
format verification
go build ./cmd/omaforge
```

Pin permissions to the minimum required. Do not add deployment, release, package publication, or third-party actions that are unnecessary. Pin third-party Actions to stable major versions or immutable commits according to the repository's documented security decision.

### 6. Add agent and contributor guidance

`AGENTS.md` must tell future coding agents to:

- Read `handoff.md` and repository docs.
- Inspect installed Omarchy rather than assume APIs.
- Never edit Omarchy-owned files.
- Keep the CLI local-first and network-independent.
- Avoid telemetry, auth, databases, sudo, package installation, and installer hooks.
- Add tests for behavior.
- Preserve user changes.
- Avoid misleading documentation.
- Never push, publish, deploy, change visibility, or make external announcements without explicit approval.
- Record compatibility evidence and architectural decisions.

### 7. Verify everything

Run all safe applicable checks, including:

```bash
gofmt
go test ./...
go vet ./...
go build ./cmd/omaforge
./path/to/built/omaforge --help
./path/to/built/omaforge version
git diff --check
```

Use a temporary build-output location or an ignored local path; do not leave compiled binaries in the repository unless deliberately part of the project.

Inspect the final repository for:

- Secrets or credentials.
- Personal absolute paths.
- Misleading shipped-feature claims.
- Unnecessary generated files.
- Accidental edits outside the repository.

## Engineering constraints

- Make the smallest coherent implementation that satisfies Milestone 0.
- Do not build multiple milestones at once.
- Prefer deterministic, testable behavior.
- Preserve the user's existing files and Git configuration.
- Do not use destructive Git commands.
- Do not add dependencies without a concrete need and explanation.
- Do not require network access for normal CLI behavior.
- Do not silently recover from incompatibility by claiming success.
- If a command or API is unavailable, document it and degrade gracefully.
- If a meaningful decision is uncertain, record options in `docs/DECISIONS.md` and choose only when necessary for Milestone 0.

## Stop conditions

Stop and ask before:

- Installing software or using `sudo`.
- Changing system or Omarchy configuration.
- Pushing or committing.
- Changing GitHub repository settings or visibility.
- Publishing or deploying anything.
- Needing credentials, tokens, or other secrets.
- Expanding beyond Milestone 0.

## Completion criteria

Milestone 0 is complete when:

- The private repository contains a coherent foundation.
- The product definition and non-goals are explicit.
- The independent-project disclaimer is present.
- Upstream compatibility inspection is documented with exact evidence.
- The Go CLI builds.
- `omaforge --help` succeeds.
- `omaforge version` succeeds.
- Tests pass.
- Formatting and vet checks pass.
- CI configuration matches the checks run locally.
- No secrets, personal paths, or misleading claims are present.
- No Omarchy-owned or external files were changed.
- Nothing was pushed, published, deployed, or made public.

## Final response format

When finished, do not continue to the scaffolder. Stop and report:

```text
Outcome
Files changed
Omarchy contract findings
Architecture and decisions made
Commands/tests run
Verification results
Known limitations or uncertainties
Recommended next milestone
Git status
```

Include a concise proposed commit message, but do not commit until I approve.

---

