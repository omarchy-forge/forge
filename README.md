# Omarchy Forge

Developer tools for building, testing, and shipping Omarchy plugins.

> Omarchy Forge is an independent community project. It is not affiliated with
> or endorsed by Omarchy, Basecamp, 37signals, or DHH.

Omarchy Forge is in early development. The first goal is a reliable plugin
scaffold and validation workflow that complements the official
`omarchy plugin validate` command.

Milestone 0 provides the repository foundation and a minimal `omaforge` CLI.
At this stage, the CLI only exposes help and development build information; it
does not yet generate or validate plugins.

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
