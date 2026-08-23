# Website dependency security notes

The 2026-08-22 refresh resolves the three previously reviewed production
advisories by overriding the Geistdocs-pinned Fumadocs dependency family to
current compatible releases:

- `fumadocs-core@16.14.5` replaces vulnerable `image-size` with
  `@fumari/image-size`, resolving `GHSA-w3rx-r6r6-pgpr` and
  `GHSA-5p2g-fcmc-qvqq`.
- `fumadocs-mdx@15.3.0` uses patched `esbuild@0.28.2`, resolving
  `GHSA-g7r4-m6w7-qqqr`.
- A local `fumadocs-ui@16.2.2` package patch rewrites its destructured sidebar
  exports as explicit bindings. This preserves Geistdocs' pinned UI API while
  allowing Next.js production bundlers to resolve those exports statically.

The overrides are required because Geistdocs still pins the older Fumadocs
versions. Typecheck and the full production build verify compatibility. CI runs
`pnpm audit:known` on every change and now
fails on any production advisory at any severity; no advisory allowlist
remains.
