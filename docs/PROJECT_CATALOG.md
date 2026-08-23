# Owner project catalog

The [Projects page](https://www.omarchyforge.com/projects) lists optional public
projects maintained by the `omarchy-forge` GitHub organization. It is an
owner-project directory, not a general marketplace or an endorsement program.

## Add a project

A repository is eligible when it is public, active, not a fork, owned by
`omarchy-forge`, and has the `omaforge-project` topic. Its default branch must
contain a valid Forge `manifest.json`, a matching stable GitHub release, a
preview image under `assets/`, and `forge-project.json`:

```json
{
  "schemaVersion": 1,
  "tagline": "A short plain-text description.",
  "previewPath": "assets/preview.png",
  "compatibility": ["omarchy-4"],
  "featured": true,
  "order": 10
}
```

The metadata schema rejects unknown fields and unsafe preview paths. The sync
job constructs installation and repository links from the verified repository
identity; projects cannot supply arbitrary commands or URLs.

After a qualifying repository is tagged and released, manually run **Sync owner
projects** in GitHub Actions or wait for its daily schedule. The workflow opens
a normal pull request containing only the generated catalog and copied preview
images, explicitly starts the protected CI workflow, and enables auto-merge.
The website changes only after all required checks pass.

The generator downloads JSON metadata and image bytes only. It never installs a
project or executes its QML.
