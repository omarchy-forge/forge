# Changelog

## Unreleased

- Add `omaforge dev <directory> --trust-plugin-code` as a one-shot entry point
  for a reviewed generated plugin's isolated Quickshell runtime harness.

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
