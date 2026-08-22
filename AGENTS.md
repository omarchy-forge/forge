# Agent guidance

- Read `handoff.md` and all relevant repository documents before acting.
- Inspect the installed Omarchy implementation; do not assume its APIs.
- Never edit Omarchy-owned files, including files under `$OMARCHY_PATH`.
- Keep the CLI local-first, deterministic, and network-independent by default.
- Do not add telemetry, accounts, authentication, databases, `sudo`, automatic
  package installation, or installer hooks.
- Do not modify `~/.config/omarchy/shell.json` without explicit authorization.
- Add tests for implemented behavior and run the applicable baseline checks.
- Preserve user changes and unrelated files; avoid destructive Git commands.
- Do not describe planned behavior as shipped behavior.
- Never push, publish, deploy, change repository visibility, create releases,
  or make external announcements without explicit approval.
- Record verified upstream facts in `docs/UPSTREAM_COMPATIBILITY.md` and
  meaningful architectural choices in `docs/DECISIONS.md`.
- Treat plugins and their QML as untrusted. Do not execute third-party QML in
  Forge checks or on persistent self-hosted CI runners.
