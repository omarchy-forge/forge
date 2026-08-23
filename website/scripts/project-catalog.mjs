const projectKeys = new Set([
  "schemaVersion",
  "projectType",
  "tagline",
  "previewPath",
  "compatibility",
  "featured",
  "order",
]);

const text = (value, label, maximum = 200) => {
  if (typeof value !== "string" || value.trim() === "" || value.length > maximum) {
    throw new Error(`${label} must be a nonempty string of at most ${maximum} characters`);
  }
  return value.trim();
};

export const validateProjectMetadata = (value) => {
  if (!value || typeof value !== "object" || Array.isArray(value)) {
    throw new Error("forge-project.json must be an object");
  }
  for (const key of Object.keys(value)) {
    if (!projectKeys.has(key)) throw new Error(`unknown forge-project.json field: ${key}`);
  }
  if (value.schemaVersion !== 1) throw new Error("forge-project.json schemaVersion must be 1");
  const projectType = value.projectType ?? "plugin";
  if (!new Set(["plugin", "cli"]).has(projectType)) {
    throw new Error("projectType must be plugin or cli");
  }
  const previewPath = text(value.previewPath, "previewPath");
  if (!/^assets\/[A-Za-z0-9._/-]+\.(png|jpe?g|webp)$/.test(previewPath) || previewPath.includes("..")) {
    throw new Error("previewPath must be a safe image path below assets/");
  }
  if (!Array.isArray(value.compatibility) || value.compatibility.length === 0 ||
      value.compatibility.some((item) => !/^omarchy-[0-9]+$/.test(item))) {
    throw new Error("compatibility must contain version keys such as omarchy-4");
  }
  if (typeof value.featured !== "boolean") throw new Error("featured must be boolean");
  if (!Number.isSafeInteger(value.order) || value.order < 0 || value.order > 10000) {
    throw new Error("order must be an integer from 0 through 10000");
  }
  return {
    compatibility: [...new Set(value.compatibility)].sort(),
    featured: value.featured,
    order: value.order,
    projectType,
    previewPath,
    tagline: text(value.tagline, "tagline"),
  };
};

export const validatePythonProject = (source) => {
  if (typeof source !== "string" || source.length > 100_000) {
    throw new Error("pyproject.toml must be text of at most 100 KB");
  }
  const header = source.match(/^\[project\]\s*$/m);
  const section = header
    ? source.slice(header.index + header[0].length).split(/^\[/m, 1)[0]
    : "";
  const name = section.match(/^name\s*=\s*["']([^"']+)["']\s*$/m)?.[1];
  const version = section.match(/^version\s*=\s*["']([^"']+)["']\s*$/m)?.[1];
  if (!name || !/^[a-z0-9][a-z0-9-]*$/.test(name)) throw new Error("Python project name is invalid");
  if (!version || !/^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$/.test(version)) {
    throw new Error("Python project version must be stable semantic versioning");
  }
  return { name, version };
};

export const validateCliInstaller = (bytes, repository) => {
  if (!Buffer.isBuffer(bytes) || bytes.length > 100_000) throw new Error("CLI installer exceeds 100 KB");
  const source = bytes.toString("utf8");
  if (!source.startsWith("#!/usr/bin/env bash\nset -euo pipefail\n")) {
    throw new Error("CLI installer is missing its strict Bash entry point");
  }
  if (!source.split("\n").includes(`repository="${repository}"`)) {
    throw new Error("CLI installer repository identity does not match");
  }
};

export const validateManifest = (value) => {
  if (!value || typeof value !== "object" || Array.isArray(value)) {
    throw new Error("manifest.json must be an object");
  }
  if (value.schemaVersion !== 1) throw new Error("manifest schemaVersion must be 1");
  const id = text(value.id, "manifest id", 100);
  if (!/^[A-Za-z0-9][A-Za-z0-9._-]*$/.test(id) || id.includes("..") || id.startsWith("omarchy.")) {
    throw new Error("manifest id is invalid or reserved");
  }
  const version = text(value.version, "manifest version", 50);
  if (!/^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$/.test(version)) {
    throw new Error("manifest version must be stable semantic versioning");
  }
  return {
    id,
    name: text(value.name, "manifest name", 100),
    version,
  };
};

export const imageExtension = (path) => {
  const extension = path.toLowerCase().split(".").pop();
  return extension === "jpeg" ? "jpg" : extension;
};

export const hasSupportedImageSignature = (bytes, extension) => {
  if (extension === "png") {
    return bytes.length >= 8 && Buffer.from(bytes.subarray(0, 8)).equals(Buffer.from([137, 80, 78, 71, 13, 10, 26, 10]));
  }
  if (extension === "jpg") return bytes.length >= 3 && bytes[0] === 0xff && bytes[1] === 0xd8 && bytes[2] === 0xff;
  if (extension === "webp") {
    return bytes.length >= 12 && Buffer.from(bytes.subarray(0, 4)).toString("ascii") === "RIFF" &&
      Buffer.from(bytes.subarray(8, 12)).toString("ascii") === "WEBP";
  }
  return false;
};

export const buildCatalogEntry = ({ repository, commit, release, manifest, pythonProject, metadata }) => {
  if (repository.owner?.login !== "omarchy-forge" || repository.private || repository.archived || repository.disabled) {
    throw new Error("catalog repository must be an active public omarchy-forge project");
  }
  if (!Array.isArray(repository.topics) || !repository.topics.includes("omaforge-project")) {
    throw new Error("catalog repository is missing the omaforge-project topic");
  }
  const checkedMetadata = validateProjectMetadata(metadata);
  const slug = text(repository.name, "repository name", 100);
  if (!/^[a-z0-9][a-z0-9-]*$/.test(slug)) throw new Error("repository name is not a safe catalog slug");
  if (!/^[0-9a-f]{40}$/.test(commit.sha)) throw new Error("default-branch commit must be a full SHA");
  const software = checkedMetadata.projectType === "plugin"
    ? validateManifest(manifest)
    : validatePythonProject(pythonProject);
  if (checkedMetadata.projectType === "cli" && software.name !== slug) {
    throw new Error("Python project name must match its repository");
  }
  if (release.tag_name !== `v${software.version}` || release.draft || release.prerelease) {
    throw new Error("latest stable release must match the project version");
  }
  const extension = imageExtension(checkedMetadata.previewPath);
  const installCommand = checkedMetadata.projectType === "plugin"
    ? `omarchy plugin add https://github.com/omarchy-forge/${slug}.git --enable`
    : `curl -fsSL https://raw.githubusercontent.com/omarchy-forge/${slug}/v${software.version}/install.sh | bash -s -- --version v${software.version}`;
  return {
    compatibility: checkedMetadata.compatibility,
    featured: checkedMetadata.featured,
    installCommand,
    name: checkedMetadata.projectType === "plugin"
      ? software.name
      : `${slug[0].toUpperCase()}${slug.slice(1)}`,
    order: checkedMetadata.order,
    pluginId: software.id ?? null,
    preview: `/project-images/${slug}.${extension}`,
    releaseUrl: release.html_url,
    repository: `omarchy-forge/${slug}`,
    repositoryUrl: repository.html_url,
    slug,
    sourceCommit: commit.sha,
    tagline: checkedMetadata.tagline,
    projectType: checkedMetadata.projectType,
    version: software.version,
  };
};

export const sortCatalog = (projects) => [...projects].sort((left, right) =>
  Number(right.featured) - Number(left.featured) || left.order - right.order || left.name.localeCompare(right.name));
