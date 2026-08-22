# Website dependency security notes

The 2026-08-22 production dependency audit reports:

- [`GHSA-w3rx-r6r6-pgpr`](https://github.com/advisories/GHSA-w3rx-r6r6-pgpr)
  and
  [`GHSA-5p2g-fcmc-qvqq`](https://github.com/advisories/GHSA-5p2g-fcmc-qvqq)
  in transitive
  `image-size@2.0.2`. GitHub's reviewed advisories list no patched version, and
  npm currently publishes no newer release.
- [`GHSA-g7r4-m6w7-qqqr`](https://github.com/advisories/GHSA-g7r4-m6w7-qqqr)
  in transitive `esbuild@0.27.7`. The issue affects the
  esbuild development server on Windows. Forge uses Next.js Turbopack and does
  not invoke esbuild's development server.

The `image-size` issue is an infinite-loop denial of service from specially
crafted ICNS, JXL, or HEIF input. This website has no image upload, remote image
host, CMS, feedback ingestion, or untrusted content pipeline. Builds process
only reviewed files committed to this repository. Do not add an upload or
external-content path until these advisories are patched and the threat model
is revisited.

CI runs `pnpm audit:known` on every change. It permits only these reviewed
advisory IDs and fails on any new finding at any severity. Remove each exception
as soon as a patched compatible release is available.
