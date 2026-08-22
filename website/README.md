# Omarchy Forge documentation website

This directory contains the Omarchy Forge documentation website, deployed at
[www.omarchyforge.com](https://www.omarchyforge.com). It uses
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

Production sets `NEXT_PUBLIC_SITE_URL=https://www.omarchyforge.com` so metadata,
the sitemap, and robots file use the canonical custom origin. The apex domain
redirects to `www`.
