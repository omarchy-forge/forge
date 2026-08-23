import { createAgentsRoute } from "@vercel/geistdocs/routes/agents";
import type { NextRequest } from "next/server";
import { config } from "@/lib/geistdocs/config";

const route = createAgentsRoute({
  config,
});

export const { generateStaticParams } = route;

export function GET(
  request: NextRequest,
  context: { params: Promise<{ lang: string }> },
) {
  return route.GET(request, context);
}
