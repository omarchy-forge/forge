import { createDocsMarkdownRoute } from "@vercel/geistdocs/routes/llms";
import { geistdocsSource } from "@/lib/geistdocs/source";

const route = createDocsMarkdownRoute({
  notFound: {},
  sources: [geistdocsSource],
});

export function generateStaticParams({ params }: { params: { lang: string } }) {
  return route.generateStaticParams({ params: Promise.resolve(params) });
}

export function GET(
  request: Request,
  context: { params: Promise<{ lang: string; slug?: string[] }> },
) {
  return route.GET(request, context);
}
