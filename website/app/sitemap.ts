import type { MetadataRoute } from "next";
import { cacheLife } from "next/cache";

import { source } from "@/lib/geistdocs/source";

const protocol = process.env.NODE_ENV === "production" ? "https" : "http";
const baseUrl = `${protocol}://${process.env.NEXT_PUBLIC_VERCEL_PROJECT_PRODUCTION_URL}`;

// biome-ignore lint/suspicious/useAwait: Next.js requires cached functions to be async.
export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  "use cache";
  cacheLife("max");

  const url = (path: string): string => new URL(path, baseUrl).toString();

  const pages: MetadataRoute.Sitemap = [];

  for (const page of source.getPages() as Array<{
    data: { lastModified?: Date };
    url: string;
  }>) {
    const data = page.data as {
      lastModified?: Date;
    };

    pages.push({
      changeFrequency: "weekly" as const,
      lastModified: data.lastModified ? new Date(data.lastModified) : undefined,
      priority: 0.5,
      url: url(page.url),
    });
  }

  return [
    {
      changeFrequency: "monthly",
      priority: 1,
      url: url("/"),
    },
    {
      changeFrequency: "weekly",
      priority: 0.8,
      url: url("/projects"),
    },
    ...pages,
  ];
}
