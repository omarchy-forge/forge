# Omarchy Forge documentation website

This directory contains the private, undeployed Omarchy Forge documentation
website. It uses Next.js, TypeScript, MDX, and Geistdocs/Fumadocs.

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

Set `NEXT_PUBLIC_SITE_URL` to the canonical origin when a deployment is
explicitly authorized. No Vercel project or domain is configured by this tree.
