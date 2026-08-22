# Contributing

Omarchy Forge is in early development. Before proposing a change, read
`handoff.md`, `AGENTS.md`, and the documents in `docs/`.

Keep changes focused on the active milestone. Behavior that depends on Omarchy
must be checked against the installed implementation or an identified official
upstream revision rather than assumed. Never edit files owned by an installed
Omarchy distribution as part of Forge development.

For code changes:

1. Add or update tests for behavior.
2. Run `gofmt` on Go files.
3. Run `go test ./...`, `go vet ./...`, and `go build ./cmd/omaforge`.
4. Update compatibility evidence or decision records when applicable.
5. Keep documentation honest about what is implemented.

By participating, you agree to follow the [Code of Conduct](CODE_OF_CONDUCT.md).
Security issues should be reported using [SECURITY.md](SECURITY.md), not a
public issue.
