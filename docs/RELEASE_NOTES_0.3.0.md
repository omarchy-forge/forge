# Omarchy Forge 0.3.0 release notes

Omarchy Forge `v0.3.0` adds a provider-neutral workflow for handing a generated
plugin project to the coding agent configured by the user through Omarchy.
Forge still does not select a model, collect credentials, contact an AI
service, or launch an agent automatically.

## Included

- `omaforge init --agent-ready` adds exactly two files to the existing reviewed
  bar-widget scaffold without changing the normal template:
  - `FORGE_SPEC.md` captures product behavior, data and command boundaries,
    required UI states, privacy, failures, timeouts, and acceptance criteria.
  - `AGENTS.md` scopes a coding agent to the plugin repository and preserves
    Forge and Omarchy safety contracts.
- A human-controlled specification gate starts in `Draft` and must be changed
  to `Ready for implementation` only after every placeholder is resolved and
  reviewed.
- The generated agent contract permits deterministic project tests, Forge
  checking, and official Omarchy validation while reserving QML execution,
  screenshot capture, plugin installation, and live-shell changes for an
  explicit human decision.
- CLI completion guidance explains the specification and readiness steps before
  directing users to the official validator.

## Security and privacy

Agent-ready generation is deterministic and local. It adds no model SDK,
provider dependency, prompt collection, credential path, telemetry, account,
database, or network request.

Generated instructions prohibit privileged commands, package installation,
plugin installation or publication, live Omarchy Shell mutation, and unapproved
expansion into networking, persistence, authentication, or secret handling.
Instructions are defense in depth rather than a sandbox: users must commit a
rollback point, review the agent's diff, and explicitly decide whether to trust
and execute plugin QML.

## Dogfooding and compatibility

A disposable local-only disk-space widget exercised the completed specification
and agent contract. Its generated tests, `omaforge check`, and the installed
`omarchy plugin validate` command passed without executing QML, installing the
plugin, or modifying shell configuration. That pass found and corrected an
ambiguous placeholder gate before release.

This release continues to target the inspected Omarchy 4 manifest schema. The
installed `omarchy agent prompt [--inline] <prompt...>` contract was verified
read-only; Forge documents that user-invoked handoff but never invokes it.

## Existing projects

Existing projects and normal `omaforge init` output are unchanged. Agent-ready
files are opt-in and are generated only when `--agent-ready` is supplied. Forge
does not retrofit, overwrite, or migrate existing plugin repositories.

This release does not install or publish plugins, submit anything to a
marketplace, change the private Handoff repository's visibility, or constitute
a public launch announcement.
