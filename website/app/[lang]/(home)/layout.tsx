import { GeistdocsHomeLayout } from "@vercel/geistdocs/home-layout";
import { config } from "@/lib/geistdocs/config";
import { getRootLang } from "@/lib/geistdocs/root-params";
import { source } from "@/lib/geistdocs/source";

const Layout = async ({ children }: LayoutProps<"/[lang]">) => {
  const lang = await getRootLang();

  return (
    <GeistdocsHomeLayout config={config} tree={source.pageTree[lang]}>
      <div className="pt-0 pb-32">{children}</div>
    </GeistdocsHomeLayout>
  );
};

export default Layout;
