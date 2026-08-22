# Website guidance

- The root `AGENTS.md` applies to this directory.
- Keep the website documentation-first, mostly static, and honest about Forge's
  current release and deployment state.
- Do not add AI chat, analytics, tracking, feedback submission, authentication,
  databases, telemetry, or server-side user data without explicit approval.
- All MDX pages require `title` and `description` frontmatter.
- Preserve the independent-project disclaimer in the shared layout.
- Use the custom Forge identity; never use the official Omarchy logo as Forge
  branding.
- Run `pnpm typecheck` and `pnpm build` after website changes.
- Do not deploy, connect a domain, or create a Vercel project without explicit
  approval.

<!-- BEGIN:nextjs-agent-rules -->

# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` (resolved from this file's directory; in monorepos the `next` package may not be visible from the repo root) before writing any code. Heed deprecation notices.

This block is written and re-added by `next dev` — verify at `node_modules/next/dist/server/lib/generate-agent-files.js`. Removing it from a diff only re-creates the uncommitted change; committing it with your work keeps the tree clean.

<!-- END:nextjs-agent-rules -->
