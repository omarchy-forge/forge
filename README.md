# Omarchy Forge

Developer tools for building, testing, and shipping Omarchy plugins.

> Omarchy Forge is an independent community project. It is not affiliated with
> or endorsed by Omarchy, Basecamp, 37signals, or DHH.

Omarchy Forge is in early development. Its first working artifact is a safe,
deterministic scaffold for an Omarchy bar widget with a popout. Generated
projects complement the official `omarchy plugin validate` command rather than
replacing it.

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
and [contribution guide](CONTRIBUTING.md) for the current scope.

## License

Omarchy Forge is licensed under the [MIT License](LICENSE).
