# Upstream compatibility evidence

## Inspection snapshot

- Inspection date: 2026-08-22.
- Installed package: `omarchy 4.0.0-1`, reported by `pacman -Q omarchy`.
- Installation path: `/usr/share/omarchy`, resolved from the installed
  `omarchy` executable and environment.
- Git revision: unavailable. The installed package directory is not a Git
  checkout, so no commit claim can be made.
- Quickshell: `/usr/bin/quickshell`.
- QML linter: not found on `PATH` during the initial inspection. A follow-up on
  2026-08-23 verified `qmllint 6.11.1` at
  `/usr/lib/qt6/bin/qmllint`, owned by installed package
  `qt6-declarative 6.11.1-3`.
- Shell health: `omarchy-shell shell ping` returned `ok`.
- `docs/omarchy-shell.md`: not present in this package.
- Upstream test files for plugin validation, registry loading, configuration,
  or hot reload: not present in this installed package snapshot.

## Commands executed

The inspection used these read-only commands or their direct equivalents:

```text
git status --short --branch
git remote -v
omarchy --version
command -v omarchy
command -v omarchy-shell
command -v quickshell
command -v qmllint
printf the OMARCHY_PATH environment value
omarchy plugin list --json
omarchy-shell shell ping
pacman -Q omarchy
git metadata queries against /usr/share/omarchy
file discovery and text inspection under /usr/share/omarchy
read-only inspection of installed third-party manifests
```

`omarchy --version` is not a supported command in this installation; it
returned an unknown-command message. `OMARCHY_PATH` was already set to
`/usr/share/omarchy`.

## Omarchy coding-agent command

Agent-ready scaffolding work on 2026-08-22 verified the installed command
contract without launching an agent. `omarchy agent --help` lists
`omarchy agent prompt [--inline] <prompt...>`, and the installed
`/usr/share/omarchy/bin/omarchy-agent-prompt` requires a nonempty prompt before
delegating to the user's default `omarchy-agent` configuration. Forge's
scaffold-only `--agent-ready` path can document that local handoff command
without selecting a provider, handling credentials, or invoking it
automatically.

A 2026-08-23 follow-up for the guided-agent design verified that
`omarchy-default-agent` prints the configured agent without changing it and
that this system currently reports `claude`. The installed `omarchy-agent`
launcher checks that the configured executable exists, preserves the launch
working directory except for its documented `$HOME` fallback, and passes
provider-specific unattended or approval-bypass flags. `omarchy agent prompt`
forwards its prompt as one argument to that launcher. Guided Forge launches
must therefore disclose the unattended-permission boundary, create a rollback
commit first, and treat generated instructions as non-sandboxed defense in
depth. No Omarchy-owned file or default-agent setting was modified during this
inspection.

## Source files inspected

- `/usr/share/omarchy/shell/README.md`
- `/usr/share/omarchy/shell/plugins/README.md`
- `/usr/share/omarchy/shell/services/PluginRegistry.qml`
- `/usr/share/omarchy/shell/services/BarWidgetRegistry.qml`
- `/usr/share/omarchy/shell/shell.qml`
- `/usr/share/omarchy/bin/omarchy-plugin-validate`
- `/usr/share/omarchy/shell/Commons/Color.qml`
- `/usr/share/omarchy/shell/Commons/Style.qml`
- `/usr/share/omarchy/shell/Commons/Border.qml`
- First-party manifests and QML examples covering full bars, bar widgets,
  menus, overlays, panels, services, and multi-kind plugins.
- Installed third-party manifests as additional local examples.

Paths above identify the package contract without embedding a personal home
directory. No upstream or user-owned file was modified.

## Confirmed manifest contract

`manifest.json` is required at the repository root. The official validator
requires these fields: `schemaVersion`, `id`, `name`, `version`, `kinds`, and
`entryPoints`. Schema version is exactly the JSON number `1`. `kinds` must be a
nonempty array and `entryPoints` must be an object.

The six documented kinds and fixed entry-point mappings are:

| Kind | Required entry-point key |
| --- | --- |
| `bar` | `bar` |
| `bar-widget` | `barWidget` |
| `menu` | `menu` |
| `overlay` | `overlay` |
| `panel` | `panel` |
| `service` | `service` |

A manifest may declare multiple kinds. Every declared known kind needs its
corresponding entry point. Entry-point values must be nonempty relative paths,
must not contain `..` or a newline, and must identify existing regular files.
The official validator rejects any symlink below the plugin folder except
within `.git`.

Third-party IDs must match `[A-Za-z0-9][A-Za-z0-9._-]*`, may not contain `..`,
and may not use the reserved `omarchy.*` namespace. The QML registry performs
some overlapping but less strict runtime checks; the CLI validator deliberately
enforces the safer installation boundary.

Optional `barWidget.defaultSection`, when present, is one of `left`, `center`,
or `right`; omission falls back to `center`. Observed bar-widget metadata also
includes `displayName`, `description`, `category`, `allowMultiple`, `defaults`,
and a `schema` array. Observed schema entries include boolean, string, integer,
and number controls with labels and optional ranges, steps, and defaults.

## Runtime, settings, and theme contracts

Third-party plugins are Git repositories installed beneath the user plugin
directory. They run as unsandboxed code in the single long-lived
`omarchy-shell` Quickshell process. Services and plugins marked
`keepLoaded: true` can remain mounted; other panels, overlays, and menus load on
demand.

The shell conditionally injects supported properties into plugin instances,
including `omarchyPath`, `shell`, `manifest`, `barWidgetRegistry`, and
`pluginRegistry`. Bar widgets use the shared `BarWidget` base and receive their
bar context and inline settings through the bar registry. Settings persist on
the matching entry in `shell.json`, not in a nested `config` object. The active
bar ID, three layout sections, enabled non-bar plugins, and disabled first-party
plugins are also represented there.

Shared QML imports expose theme and layout primitives through `qs.Commons` and
`qs.Ui`. Verified tokens include the `Color`, `Style`, and `Border` singletons;
widgets also obtain bar foreground, font, orientation, dimensions, tooltip,
and command helpers from their bar context. Future templates must follow these
installed components instead of hard-coding a parallel visual system.

## Official validation and local development

`omarchy plugin validate <folder>` performs JSON/schema, required-field, ID,
kinds, entry-point, default-section, file-existence, and symlink checks. It does
not establish polish, accessibility, runtime correctness, or publish readiness.

The installed documentation describes `omarchy plugin add <local-path>
--enable` as a development path. Manual installation followed by
`omarchy-shell shell rescanPlugins` is also documented. Saving beneath the user
plugin directory automatically reloads plugin code in this version; the rescan
IPC command remains an explicit fallback. The shell persists customization in
`~/.config/omarchy/shell.json`; Forge must not silently modify it.

## Differences and uncertainties

- The handoff expected an Omarchy version command, but this package does not
  implement `omarchy --version`; package-manager metadata supplied the version.
- No installed Git revision was available, so compatibility is tied only to
  package `4.0.0-1` and the inspected files, not an upstream commit.
- `qmllint` was absent from `PATH` but was later verified in Qt's standard tool
  directory. The referenced shell document and upstream registry/validator
  tests were absent locally and cannot be claimed as verified here.
- The richer official Basecamp third-party plugin was not present locally and
  was not needed for the Milestone 0 executable. It should be inspected at a
  pinned official revision before designing the first template.
- Milestone 0 intentionally contained no template and executed no QML. The
  Milestone 1 checks below supersede that earlier runtime limitation.

## Milestone 1 follow-up inspection

On 2026-08-22, template work also inspected these pinned official revisions:

- Omarchy `quattro`: `ed7bae4ac5a570e9df307486e0202fdafcc6ee24`.
- Basecamp plugin: `abc1ba72aaf47db530d2a0c6901d99f0f98e6aa7`.

The upstream checkout supplied the validator, registry, manifest entry-point,
and bar-widget contract tests missing from the installed package. The Basecamp
plugin confirmed current practices for `Panel`, `BarIconButton`,
`KeyboardPanel`, `PanelKeyCatcher`, array-form `Process.command` values, inline
settings, fictional demo data, local-path installation, privacy documentation,
and QML tests.

A Forge-generated sample passed the installed `omarchy plugin validate`
command. Its actual entry point was then instantiated using the pinned
Omarchy bar-widget contract harness with an isolated temporary home and config.
The harness confirmed delayed bar/settings injection, finite dimensions, safe
refresh/close calls, and horizontal-to-vertical geometry changes.

A guarded manual session then installed a generated Git-backed sample through
`omarchy plugin add`, exercised top, bottom, left, and right bar positions,
keyboard refresh and close behavior, a light theme, disable/re-enable, and
removal. The generated `demo/run` script was also verified visually in its
ready, empty, and error states. Its demo state is applied to the live widget
through the widget's `IpcHandler`; it does not rely on transient environment
variables, restart the shell, or persist demo configuration.

The manual session backed up `shell.json` before each mutation and restored it
byte-for-byte afterward. The pre- and post-test SHA-256 values were identical,
the plugin was absent after cleanup, the original Aether theme and workspace
were restored, and `omarchy-shell shell ping` returned `ok`.

## Milestone 2 diagnostic verification

The deterministic Forge rule engine was checked against the same manifest
schema and kind-to-entry-point mapping implemented by the installed official
validator. Forge labels these structural findings `official-parity`; stricter
quality, security, UX, and publish-readiness findings are labeled `forge` and
heuristic matches explicitly require human review.

On 2026-08-22, the read-only doctor probes detected `omarchy 4.0.0-1`, found
Quickshell, received `ok` from shell IPC, and passed a generated sample through
`omarchy plugin validate`. `qmllint` remained unavailable on `PATH` and was
reported as an optional-tool warning at that checkpoint. A 2026-08-23 follow-up
verified it in `/usr/lib/qt6/bin`, motivating a doctor fallback for Qt's tool
directory. Doctor did not modify user or packaged
configuration, install software, execute plugin QML, or use the network.

## Isolated generated-project runtime harness

On 2026-08-22, a newly generated bar widget passed the opt-in runtime harness
against the installed Omarchy 4 `Commons` and `Ui` modules and Quickshell
`0.3.0` revision `28771c7c74b42e20afca0b1b63980cb46515537c`. The harness
loaded the actual generated `Panel.qml`, confirmed its refresh/close lifecycle
and finite implicit geometry, and rejected deliberately invalid QML.

The run used temporary HOME and XDG config/data/cache/runtime directories, left
no harness process or temporary directory behind, did not install or enable the
plugin, did not edit shell configuration, and did not call the live shell IPC.
This is trusted local runtime evidence, not part of static checking or CI.
