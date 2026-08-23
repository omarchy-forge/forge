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
`omaforge dev <directory> --trust-plugin-code` is an explicit convenience entry
point for that same project-owned harness; it does not weaken or infer the trust
acknowledgement.
The optional `--state` value is restricted to `ready`, `empty`, or `error` and
is passed as a fixed argument into the harness. Generated plugins apply it only
to their fictional in-memory demo contract inside the temporary runtime.

`omaforge screenshot` uses Qt Quick item-level rendering against a
template-declared `forgeScreenshotTarget`. It never invokes monitor or window
screen-copy tools. Output must be a new `.png` path, preventing accidental
replacement of an existing asset.
Watch mode performs local reads only, excludes `.git`, and executes only the
already acknowledged project harness. Each change receives a fresh isolated
runtime rather than keeping plugin QML resident in the live shell.

The installed Omarchy implementation is authoritative but owned externally.
Forge may inspect it and invoke documented commands; it must never modify those
files. User configuration such as `~/.config/omarchy/shell.json` also remains
outside Forge's write scope unless a user explicitly requests a documented
operation.

`omaforge init --agent-ready` generates instructions but never invokes an AI
agent. Its project-scoped `AGENTS.md` prohibits QML execution, plugin
installation, live-shell changes, privileged commands, package installation,
and unapproved expansion into networking, persistence, authentication, or
secret handling. `FORGE_SPEC.md` tells users not to record secrets and requires
a human-reviewed acceptance boundary. These are defense-in-depth instructions,
not a sandbox: users must review agent changes and retain a clean Git rollback
point before running an unattended Omarchy agent.

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

The public `install.sh` is an explicit networked installation boundary, not
part of normal local-first checking. It resolves a published release only when
the user invokes it, accepts only an exact semantic-version tag, downloads the
single matching architecture archive, verifies its SHA-256 value against the
release checksum manifest, and replaces only `~/.local/bin/omaforge` by default.
It uses no `sudo`, package manager, shell startup mutation, telemetry, scheduled
task, or background update check. The same audited path handles a later
user-requested update and exits early when the installed version is current.

## Current implementation

The CLI's check path remains static, deterministic, and network-free. Local
environment probes are read-only. Scaffold writes are collision-protected, and
runtime QML execution is confined to the separately invoked trusted harness
described above.
