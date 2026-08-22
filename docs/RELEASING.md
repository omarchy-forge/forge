# Release process

Forge releases are created only from exact semantic-version tags such as
`v1.0.0`. Pushing a matching tag triggers the release workflow; ordinary
branches and pull requests never publish releases.

## Before tagging

1. Confirm `main` is clean, reviewed, and green in CI.
2. Update the changelog and release documentation.
3. Run the complete local baseline, including `tests/release-action.sh`.
4. Review the tag target and obtain explicit authorization to create a release.

For `v0.1.0`, also verify `docs/RELEASE_NOTES_0.1.0.md` and confirm that the
Action guide pins both the Action reference and downloaded version to
`v0.1.0`.

## Artifacts

The workflow cross-compiles static Linux binaries for `amd64` and `arm64` and
packages each with the license and README:

```text
omaforge_<version>_linux_amd64.tar.gz
omaforge_<version>_linux_arm64.tar.gz
checksums.txt
```

Build metadata is injected into `omaforge version`. Archives normalize file
order, ownership, timestamps, and gzip metadata. `checksums.txt` is generated
after both archives and verified before workflow upload.

The build job has read-only repository permissions. A separate publish job
receives `contents: write`, downloads the immutable workflow artifact, and uses
GitHub CLI to create the release for the existing tag. The workflow uses
official GitHub-maintained Actions pinned to immutable commit SHAs.

## Local packaging check

```bash
scripts/build-release.sh v0.1.0 "$(git rev-parse HEAD)" "$(git show -s --format=%cI HEAD)" ./dist
cd dist
sha256sum --check checksums.txt
```

The output directory must be absent or empty; the script refuses to overwrite
an existing release bundle.

Creating tags, pushing tags, and creating GitHub Releases are external
publishing actions and always require explicit approval.
