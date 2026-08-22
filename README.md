# Omarchy Forge

Developer tools for building, testing, and shipping Omarchy plugins.

> Omarchy Forge is an independent community project. It is not affiliated with
> or endorsed by Omarchy, Basecamp, 37signals, or DHH.

Omarchy Forge is in early development. Its first working artifact is a safe,
deterministic scaffold for an Omarchy bar widget with a popout. Generated
projects complement the official `omarchy plugin validate` command rather than
replacing it.

## Check a plugin

Forge checks are deterministic, local, and do not execute plugin QML or use the
network:

```bash
omaforge check .
omaforge check . --format json
omaforge check . --format sarif
```

Findings identify whether a rule mirrors the inspected official validator or
is a Forge quality rule. Exit code 0 means no error-severity findings, 1 means
the project has errors, and 2 means command usage is invalid. Warnings are
reported without failing the check.

For human-readable local environment diagnostics, including Omarchy,
Quickshell, shell IPC, QML tooling, and the official validator, run:

```bash
omaforge doctor .
```

## GitHub Action and releases

The repository contains a composite check Action and a tag-driven Linux release
pipeline. The Action downloads an exact release version, verifies its checksum,
and emits pull-request annotations. No public Forge release or stable Action tag
exists yet, so these files are release infrastructure rather than an available
installation promise. See the [Action guide](action/README.md) and
[release process](docs/RELEASING.md).

## Documentation website

The undeployed website lives in `website/` and covers the homepage, quickstart,
templates, commands, plugin anatomy, manifest, compatibility, roadmap, and
contributing guides. It uses a static-first Next.js/MDX architecture without AI
chat, analytics, tracking, accounts, authentication, or a database.

```bash
cd website
corepack enable
pnpm install --frozen-lockfile
pnpm typecheck
pnpm build
```

Deployment and domain connection require separate explicit approval.

## Create a plugin

```bash
omaforge init project-pulse \
  --id dev.example.project-pulse \
  --author "Example Developer"
```

Run `omaforge init --help` for interactive and noninteractive options,
including `--dry-run`, section selection, optional CI, and local Git
initialization. Forge refuses dangerous targets and nonempty directories unless
the generated-file collisions are explicitly previewed with `--force`.

## Build and run

Go 1.23 or newer is required for development.

```bash
go build -o ./tmp/omaforge ./cmd/omaforge
./tmp/omaforge --help
./tmp/omaforge version
./tmp/omaforge check .
```

## Development

Run the same baseline checks used by CI:

```bash
test -z "$(gofmt -l .)"
go test ./...
go vet ./...
go build ./cmd/omaforge
```

See [the product definition](docs/PRODUCT.md), [roadmap](docs/ROADMAP.md),
[rule reference](docs/RULES.md), and [contribution guide](CONTRIBUTING.md) for
the current scope.

## License

Omarchy Forge is licensed under the [MIT License](LICENSE).
