import { defineConfig } from "@vercel/geistdocs/config";
import {
  basePath,
  github,
  Logo,
  nav,
  title,
  translations,
} from "@/geistdocs";

export const config = defineConfig({
  title,
  defaultLanguage: "en",
  logo: <Logo />,
  github,
  nav,
  basePath,
  translations,
  content: [{ id: "docs", label: "Docs", dir: "content/docs", route: "/docs" }],
  ai: { enabled: false },
  feedback: { enabled: false },
  pageActions: {
    askAI: false,
    copyPage: true,
    editSource: true,
    openInChat: false,
    scrollTop: true,
  },
  search: { enabled: true },
});
