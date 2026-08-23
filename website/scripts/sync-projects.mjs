import { mkdir, readdir, readFile, unlink, writeFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import {
  buildCatalogEntry,
  hasSupportedImageSignature,
  imageExtension,
  sortCatalog,
  validateProjectMetadata,
} from "./project-catalog.mjs";

const scriptDirectory = dirname(fileURLToPath(import.meta.url));
const websiteRoot = resolve(scriptDirectory, "..");
const outputPath = resolve(websiteRoot, "data/projects.json");
const imageDirectory = resolve(websiteRoot, "public/project-images");
const token = process.env.GITHUB_TOKEN ?? "";
const headers = {
  Accept: "application/vnd.github+json",
  "User-Agent": "omarchy-forge-project-catalog",
  "X-GitHub-Api-Version": "2022-11-28",
  ...(token ? { Authorization: `Bearer ${token}` } : {}),
};

const github = async (path, options = {}) => {
  const response = await fetch(`https://api.github.com${path}`, {
    ...options,
    headers: { ...headers, ...options.headers },
  });
  if (!response.ok) throw new Error(`GitHub ${path} returned ${response.status}`);
  return response;
};

const json = async (path) => github(path).then((response) => response.json());
const raw = async (repository, path, ref) => github(
  `/repos/${repository}/contents/${encodeURIComponent(path).replaceAll("%2F", "/")}?ref=${encodeURIComponent(ref)}`,
  { headers: { Accept: "application/vnd.github.raw+json" } },
).then((response) => response.arrayBuffer()).then((buffer) => Buffer.from(buffer));

const parseJSON = (bytes, label) => {
  if (bytes.length > 100_000) throw new Error(`${label} exceeds 100 KB`);
  try {
    return JSON.parse(bytes.toString("utf8"));
  } catch (error) {
    throw new Error(`${label} is not valid JSON: ${error.message}`);
  }
};

const repositories = await json("/orgs/omarchy-forge/repos?type=public&per_page=100");
const eligible = repositories.filter((repository) =>
  !repository.fork && !repository.archived && !repository.disabled && !repository.private &&
  repository.topics?.includes("omaforge-project"));

const projects = [];
const images = new Map();
for (const repository of eligible) {
  const fullName = repository.full_name;
  const commit = await json(`/repos/${fullName}/commits/${encodeURIComponent(repository.default_branch)}`);
  const ref = commit.sha;
  const [release, manifestBytes, metadataBytes] = await Promise.all([
    json(`/repos/${fullName}/releases/latest`),
    raw(fullName, "manifest.json", ref),
    raw(fullName, "forge-project.json", ref),
  ]);
  const manifest = parseJSON(manifestBytes, `${fullName}/manifest.json`);
  const metadata = parseJSON(metadataBytes, `${fullName}/forge-project.json`);
  const checkedMetadata = validateProjectMetadata(metadata);
  const entry = buildCatalogEntry({ repository, commit, release, manifest, metadata });
  const preview = await raw(fullName, checkedMetadata.previewPath, ref);
  if (preview.length > 2_000_000) throw new Error(`${fullName} preview exceeds 2 MB`);
  const extension = imageExtension(checkedMetadata.previewPath);
  if (!hasSupportedImageSignature(preview, extension)) throw new Error(`${fullName} preview signature is invalid`);
  images.set(`${entry.slug}.${extension}`, preview);
  projects.push(entry);
}

if (projects.length === 0) throw new Error("no eligible owner projects were found");

await mkdir(dirname(outputPath), { recursive: true });
await mkdir(imageDirectory, { recursive: true });
await writeFile(outputPath, `${JSON.stringify({ schemaVersion: 1, projects: sortCatalog(projects) }, null, 2)}\n`);

const expectedImages = new Set(images.keys());
for (const existing of await readdir(imageDirectory)) {
  if (!expectedImages.has(existing)) await unlink(resolve(imageDirectory, existing));
}
for (const [name, bytes] of images) await writeFile(resolve(imageDirectory, name), bytes);

const catalog = JSON.parse(await readFile(outputPath, "utf8"));
console.log(`synchronized ${catalog.projects.length} owner project(s)`);
