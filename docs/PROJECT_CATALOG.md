# Owner project catalog

The [Projects page](https://www.omarchyforge.com/projects) lists optional public
projects maintained by the `omarchy-forge` GitHub organization. It is an
owner-project directory, not a general marketplace or an endorsement program.

## Add a project

A repository is eligible when it is public, active, not a fork, owned by
`omarchy-forge`, and has the `omaforge-project` topic. Its default branch must
contain a matching stable GitHub release, a preview image under `assets/`, and
`forge-project.json`. Plugin projects also need a valid Forge `manifest.json`:

```json
{
  "schemaVersion": 1,
  "projectType": "plugin",
  "tagline": "A short plain-text description.",
  "previewPath": "assets/preview.png",
  "compatibility": ["omarchy-4"],
  "featured": true,
  "order": 10
}
```

`projectType` may be `plugin` (the default) or `cli`. A CLI project must have a
PEP 621 `pyproject.toml` whose package name matches the repository, and an
explicit strict-Bash `install.sh` tied to that organization repository. Its
stable release tag must match the package version. The catalog then derives an
exact-version installer command from those verified fields; the project cannot
supply arbitrary commands or URLs.

The metadata schema rejects unknown fields and unsafe preview paths. Plugin
installation uses the official Omarchy command derived from repository
identity. CLI installation uses the reviewed release installer at an immutable
tag and passes the same exact version explicitly.

After a qualifying repository is tagged and released, manually run **Sync owner
projects** in GitHub Actions or wait for its daily schedule. The workflow opens
a normal pull request containing only the generated catalog and copied preview
images, explicitly starts the protected CI workflow, and enables auto-merge.
The website changes only after all required checks pass.

The generator downloads bounded metadata, installer text, and image bytes only.
It never executes an installer, installs a project, or loads project QML.
