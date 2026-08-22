import "../global.css";
import { Footer } from "@vercel/geistdocs/footer";
import { Navbar } from "@vercel/geistdocs/navbar";
import { GeistdocsProvider } from "@/components/geistdocs/provider";
import { config } from "@/lib/geistdocs/config";
import { mono, sans } from "@/lib/geistdocs/fonts";
import { i18n } from "@/lib/geistdocs/i18n";
import { getRootLang } from "@/lib/geistdocs/root-params";
import { cn } from "@/lib/utils";
import type { Metadata, Viewport } from "next";

const siteUrl = process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000";

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title: {
    default: "Omarchy Forge",
    template: "%s | Omarchy Forge",
  },
  description:
    "Local-first developer tools for building, checking, and shipping Omarchy plugins.",
  openGraph: {
    title: "Omarchy Forge",
    description: "From an empty directory to a publish-ready Omarchy plugin.",
    type: "website",
  },
};

export const viewport: Viewport = {
  colorScheme: "dark light",
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "#faf8f5" },
    { media: "(prefers-color-scheme: dark)", color: "#11100f" },
  ],
};

export const generateStaticParams = () =>
  i18n.languages.map((lang) => ({ lang }));

const Layout = async ({ children }: LayoutProps<"/[lang]">) => {
  const lang = await getRootLang();

  return (
    <html
      className={cn(sans.variable, mono.variable, "antialiased")}
      lang={lang}
      suppressHydrationWarning
    >
      <body>
        <GeistdocsProvider basePath={config.basePath} lang={lang}>
          <Navbar config={config} />
          {children}
          <aside className="forge-disclaimer">
            Omarchy Forge is an independent community project. It is not
            affiliated with or endorsed by Omarchy, Basecamp, 37signals, or DHH.
          </aside>
          <Footer />
        </GeistdocsProvider>
      </body>
    </html>
  );
};

export default Layout;
