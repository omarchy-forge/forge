# Omarchy Forge documentation website

This directory contains the Omarchy Forge documentation website, deployed at
[omarchy-forge-docs.vercel.app](https://omarchy-forge-docs.vercel.app). It uses
Next.js, TypeScript, MDX, and Geistdocs/Fumadocs.

```bash
corepack enable
pnpm install --frozen-lockfile
pnpm typecheck
pnpm dev
```

The site intentionally disables AI chat, feedback submission, analytics,
tracking, and Next.js telemetry. It has no database, authentication, account,
or server-side user-data feature. Search and Markdown/LLM-readable routes operate
only on the checked-in documentation corpus.

The production Vercel origin is the canonical origin. `NEXT_PUBLIC_SITE_URL`
can override it for a separately authorized custom domain; no custom domain is
currently connected.
