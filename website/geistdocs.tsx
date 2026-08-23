import { Anvil } from "lucide-react";

export const Logo = () => (
  <span className="flex items-center gap-2 font-semibold text-gray-1000 text-lg leading-none tracking-[-3%]">
    <Anvil aria-hidden="true" className="size-5 text-[var(--forge-accent)]" />
    Omarchy Forge
  </span>
);

export const github = {
  branch: "main",
  editPath: "website/content/docs/{path}",
  owner: "omarchy-forge",
  repo: "forge",
};

export const nav = [
  {
    label: "Projects",
    href: "/projects",
  },
  {
    label: "Quickstart",
    href: "/docs/quickstart",
  },
  {
    label: "Source",
    href: `https://github.com/${github.owner}/${github.repo}/`,
  },
];

export const title = "Omarchy Forge Documentation";

export const translations = {
  en: {
    displayName: "English",
  },
};

export const basePath: string | undefined = undefined;
