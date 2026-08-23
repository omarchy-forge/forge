# Changelog

## Unreleased

- Add the released Control Center process-monitoring plugin to the generated
  owner-project catalog with its checked-in preview and derived install command.
- Require reference-driven projects to treat UI mockups as visual-fidelity
  targets and supplied logos/icons as mandatory assets, with human screenshot
  comparison and explicit approval for material deviations.
- Keep guided reference selection open when a directory scan finds no usable
  files, report the exact directory being scanned, and allow an explicit retry
  after file copies finish without implying that the staging directory was
  never created.
- Accept bounded, XML-validated SVG references without rendering or executing
  them, and make each retry prompt trigger one scan instead of requiring two
  consecutive Returns.
- Add an interactive `omaforge init --agent` workflow that reports the current
  Omarchy agent, gathers and confirms a complete specification, generates an
  initial prompt, creates a local baseline commit, and launches the configured
  agent only after explicit confirmation.
- Create a guided `references/` directory for up to ten confirmed local text,
  Markdown, PNG, JPEG, or WebP files with bounded validation and SHA-256
  provenance.
- Guide users to keep the bar entry minimal while separately specifying exact
  dashboard cards, user actions, data sources, and local commands. Generated
  prompts require functional wiring instead of generic placeholder dashboards.
- Require implementation agents to compare attached reference requirements with
  the confirmed spec and stop on missing or conflicting capabilities.
- Let guided users choose reference-driven creation or the detailed
  questionnaire. Reference mode treats confirmed text and images as the primary
  product brief and asks only for explicit access and failure boundaries.
- Give generated agents an ordered completion handoff covering static checks,
  isolated runtime states, deliberate local installation, installed demos,
  removal, and optional Git publishing.
- Clarify generated pre-install state testing and replace bare live-demo
  `Target not found` failures with actionable isolated-runtime guidance.

## [0.4.0] - 2026-08-23

- Detect `qmllint` in Qt 6's standard tool directory when it is installed but
  not exposed through the user's `PATH`.
- Add a generated owner-project catalog and website Projects page, beginning
  with the public Handoff plugin and Omaudit CLI security tool, plus a daily
  protected pull-request sync. CLI entries use validated package metadata and
  immutable, derived installer commands.
- Add an explicitly invoked, checksum-verifying user installer that resolves
  the latest release, skips an already-current installation, and safely updates
  `~/.local/bin/omaforge` without `sudo` or background checks.
- Reject NUL bytes in plugin source with error rule `OF305`.

## [0.3.0] - 2026-08-22

- Add opt-in `omaforge init --agent-ready` scaffolding with a structured
  `FORGE_SPEC.md`, acceptance criteria, and project-scoped `AGENTS.md` safety
  contract. Forge does not contact or automatically launch an AI agent.
- Gate agent implementation on an explicit human-reviewed specification status
  and permit the generated static test entry point in the agent safety contract.

## [0.2.0] - 2026-08-22

- Add `omaforge dev <directory> --trust-plugin-code` as a one-shot entry point
  for a reviewed generated plugin's isolated Quickshell runtime harness.
- Add isolated fictional `ready`, `empty`, and `error` state simulation to
  `omaforge dev` without installing the plugin or persisting state.
- Add deterministic plugin-only PNG capture through a template-declared Qt
  Quick item, with no desktop capture and no overwrite behavior.
- Add debounced `omaforge dev --watch` sessions that rerun fresh isolated
  harnesses after local project changes and stop cleanly on Ctrl-C.

## [0.1.0] - 2026-08-22

- Add an undeployed documentation-first Next.js/MDX website with the complete
  initial Forge documentation set and a distinct visual identity.
- Add a checksum-verifying composite GitHub Action with native annotations.
- Add reproducible Linux release archives, SHA-256 checksums, and tag-driven
  GitHub Release automation.
- Add deterministic `omaforge check` reporting in text, JSON, and SARIF.
- Add read-only `omaforge doctor` diagnostics for Omarchy plugin development.
- Add a shared, versioned finding model that distinguishes official-parity and
  Forge quality rules.
- Add interactive and noninteractive `omaforge init` scaffolding.
- Add a theme-aware bar-widget-with-popout template with safe local command
  execution, timeouts, configurable refresh, keyboard support, UI states,
  fictional demo data, project tests, and optional CI.
- Add dry-run, collision, dangerous-path, symlink, and input safety controls.
- Allow scaffolding into a repository containing only a real `.git` directory
  without requiring `--force`.
- Add deterministic golden tests and official Omarchy validator integration.
- Add an explicitly trusted, isolated Quickshell entry-point smoke test to
  generated projects without executing QML in static checks or CI.
- Record the Handoff flagship plugin dogfooding evidence and the deliberate
  deferral criteria for a service-plus-widget template.

[0.1.0]: https://github.com/omarchy-forge/forge/releases/tag/v0.1.0
[0.2.0]: https://github.com/omarchy-forge/forge/releases/tag/v0.2.0
[0.3.0]: https://github.com/omarchy-forge/forge/releases/tag/v0.3.0
[0.4.0]: https://github.com/omarchy-forge/forge/releases/tag/v0.4.0
