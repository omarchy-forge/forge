# Changelog

## Unreleased

- Add opt-in `omaforge init --agent-ready` scaffolding with a structured
  `FORGE_SPEC.md`, acceptance criteria, and project-scoped `AGENTS.md` safety
  contract. Forge does not contact or automatically launch an AI agent.

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
