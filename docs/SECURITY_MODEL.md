# Security model

## Assets

Forge must protect the user's files, source repositories, desktop session,
credentials, and Omarchy installation. Generated projects and reports must not
leak secrets or personal absolute paths.

## Trust boundaries

Omarchy plugin source is untrusted input. Enabled QML executes unsandboxed in
the long-lived `omarchy-shell` process, so a plugin can affect the desktop
session with the user's authority. Static Forge checks must parse files without
executing QML. Runtime testing belongs only in trusted or disposable
environments.

The installed Omarchy implementation is authoritative but owned externally.
Forge may inspect it and invoke documented commands; it must never modify those
files. User configuration such as `~/.config/omarchy/shell.json` also remains
outside Forge's write scope unless a user explicitly requests a documented
operation.

## Required controls

- No telemetry, secret collection, account, or cloud requirement.
- No automatic `sudo`, package installation, or installer hooks.
- No arbitrary remote QML execution in CI or a future website.
- No destructive overwrite of an existing project.
- Validate output paths before writing and prevent directory escape.
- Escape user input for every generated file format.
- Use timeouts and argument-safe process execution for external commands.
- Perform no network access during normal checking unless explicitly requested.
- Report unknown compatibility as unknown, not success.
- Never run untrusted pull-request code on a persistent personal Omarchy host.

## Current milestone

The Milestone 0 CLI reads no project files, runs no external commands, performs
no network access, and writes no user data. Its only commands print static help
or build metadata.
