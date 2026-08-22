# Service-plus-widget template evaluation

Evaluation date: 2026-08-22.

## Outcome

Do not add a service-plus-widget template yet. The shape solves a real problem,
but the installed Omarchy contract does not yet provide enough documented,
third-party evidence for Forge to present one implementation as a stable
starter.

Revisit this decision when either:

- Omarchy documents a supported third-party bar-widget-to-service access or
  injection contract; or
- a second independently developed plugin proves a common structure that Forge
  can extract without relying on first-party internals.

## Evidence inspected

The installed `omarchy 4.0.0-1` package has one manifest combining `service`
and `bar-widget`: `omarchy.media`. Its manifest declares `keepLoaded: true`,
`Service.qml`, and `BarWidget.qml`.

The shell creates one enabled service object per plugin ID beneath its hidden
service host. It injects optional shell, manifest, registry, and Omarchy-path
properties. Services are destroyed when their plugin is disabled, removed, or
reloaded. This is materially different from putting a data object inside a bar
widget, because bar-widget instances may exist once per monitor.

The media widget obtains its service through
`bar.shell.firstPartyServiceFor("omarchy.media")`. The shell also implements a
generic `serviceFor(pluginId)` function, but the installed documentation does
not establish widget access through that function as a supported third-party
contract. No installed third-party service-plus-widget example was available.

Handoff demonstrates the motivation but not a reusable service contract. Its
internal `DataService` owns persistence and Git state inside the bar-widget
entry point. That design passed its first private release, and no concrete
multi-monitor defect has been recorded. Migrating it solely to justify a Forge
template would manufacture evidence rather than learn from it.

## Why deferral is safer

A template would need to choose how the widget finds its service, how missing
or asynchronously loaded services appear in the UI, whether `keepLoaded` is
required, and how settings are divided between the service and widget. Those
choices affect lifecycle, reload, multi-monitor behavior, and compatibility.

Forge should not turn the current first-party lookup pattern into a community
contract prematurely. The existing bar-widget template remains the supported
starting point; authors with a verified need can add a service deliberately
against the Omarchy version they target.
