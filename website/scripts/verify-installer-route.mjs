import { readFile, stat } from "node:fs/promises";

const installer = new URL("../public/install.sh", import.meta.url);
const proxy = new URL("../proxy.ts", import.meta.url);
const quickstart = new URL("../content/docs/quickstart.mdx", import.meta.url);

const [installerStat, installerSource, proxySource, quickstartSource] =
  await Promise.all([
    stat(installer),
    readFile(installer, "utf8"),
    readFile(proxy, "utf8"),
    readFile(quickstart, "utf8"),
  ]);

if ((installerStat.mode & 0o111) === 0) {
  throw new Error("public/install.sh must remain executable");
}
if (!installerSource.startsWith("#!/usr/bin/env bash\nset -euo pipefail\n")) {
  throw new Error("public/install.sh is missing its strict Bash entry point");
}
if (!proxySource.includes("|install.sh)")) {
  throw new Error("the locale proxy must bypass the canonical /install.sh path");
}
if (!quickstartSource.includes("https://www.omarchyforge.com/install.sh")) {
  throw new Error("the Quickstart must document the canonical installer URL");
}

console.log("installer route boundary checks passed");
