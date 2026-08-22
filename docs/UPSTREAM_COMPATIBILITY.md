# Upstream compatibility evidence

## Inspection snapshot

- Inspection date: 2026-08-22.
- Installed package: `omarchy 4.0.0-1`, reported by `pacman -Q omarchy`.
- Installation path: `/usr/share/omarchy`, resolved from the installed
  `omarchy` executable and environment.
- Git revision: unavailable. The installed package directory is not a Git
  checkout, so no commit claim can be made.
- Quickshell: `/usr/bin/quickshell`.
- QML linter: not found on `PATH` during inspection.
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
- `qmllint`, the referenced shell document, and upstream registry/validator
  tests were absent locally. They cannot be claimed as verified here.
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
