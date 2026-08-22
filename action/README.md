# Omarchy Forge Check Action

This composite Action downloads an exact Forge release, verifies its SHA-256
checksum, runs deterministic static checks, and emits GitHub error, warning, or
notice annotations. It never executes plugin QML.

> No stable Forge release or Action tag exists yet. The version below documents
> the intended contract for the first authorized release.

```yaml
permissions:
  contents: read

steps:
  - uses: actions/checkout@v6
  - uses: omarchy-forge/forge/action@v1
    with:
      version: v1.0.0
      path: .
```

Pin both the Action reference and `version` to reviewed releases. Moving tags
such as `latest` are rejected. While the Forge repository is private, pass a
token that can read its release assets:

```yaml
      github-token: ${{ secrets.GITHUB_TOKEN }}
```

The `report` output contains the path to the schema-v1 JSON report. Error
findings fail the step, warnings annotate the pull request without failing it,
and invalid Action or check usage fails the step.

The Action supports GitHub-hosted Linux `x64` and `arm64` runners. It requires
`bash`, `curl`, `tar`, `sha256sum`, and Node.js, all present on supported
GitHub-hosted Ubuntu runners.
