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

Generated bar-widget projects include an optional local Quickshell harness.
It refuses to run without the literal `--trust-plugin-code` acknowledgement
because loading the entry point also permits that QML to invoke its documented
local commands. The harness copies the reviewed project without `.git` into a
temporary tree, isolates HOME and XDG config/data/cache/runtime paths, uses the
installed read-only Omarchy QML contract, runs it in a dedicated process group,
enforces a timeout, and deletes the tree afterward. It does not install or
enable the plugin, edit `shell.json`, or connect to the live `omarchy-shell`
process.

The harness is excluded from `omaforge check`, generated CI, and Forge's normal
automated tests. Those paths remain static and do not execute plugin QML.

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

## Current implementation

The CLI's check path remains static, deterministic, and network-free. Local
environment probes are read-only. Scaffold writes are collision-protected, and
runtime QML execution is confined to the separately invoked trusted harness
described above.
