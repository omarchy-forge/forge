import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";
import {
  buildCatalogEntry,
  hasSupportedImageSignature,
  imageExtension,
  sortCatalog,
  validateCliInstaller,
  validateProjectMetadata,
  validatePythonProject,
} from "./project-catalog.mjs";

const metadata = {
  schemaVersion: 1,
  tagline: "Pick up work where you left off.",
  previewPath: "assets/preview.png",
  compatibility: ["omarchy-4"],
  featured: true,
  order: 10,
};

test("builds a fixed install command from validated organization metadata", () => {
  const project = buildCatalogEntry({
    repository: {
      archived: false,
      disabled: false,
      html_url: "https://github.com/omarchy-forge/handoff",
      name: "handoff",
      owner: { login: "omarchy-forge" },
      private: false,
      topics: ["omaforge-project"],
    },
    commit: { sha: "a".repeat(40) },
    release: {
      draft: false,
      prerelease: false,
      html_url: "https://github.com/omarchy-forge/handoff/releases/tag/v0.1.1",
      tag_name: "v0.1.1",
    },
    manifest: { schemaVersion: 1, id: "org.omarchyforge.handoff", name: "Handoff", version: "0.1.1" },
    metadata,
  });
  assert.equal(project.installCommand, "omarchy plugin add https://github.com/omarchy-forge/handoff.git --enable");
  assert.equal(project.projectType, "plugin");
  assert.equal(project.preview, "/project-images/handoff.png");
});

test("builds an immutable CLI installer command from validated project metadata", () => {
  const project = buildCatalogEntry({
    repository: {
      archived: false,
      disabled: false,
      html_url: "https://github.com/omarchy-forge/omaudit",
      name: "omaudit",
      owner: { login: "omarchy-forge" },
      private: false,
      topics: ["omaforge-project"],
    },
    commit: { sha: "b".repeat(40) },
    release: {
      draft: false,
      prerelease: false,
      html_url: "https://github.com/omarchy-forge/omaudit/releases/tag/v0.1.0",
      tag_name: "v0.1.0",
    },
    pythonProject: '[project]\nname = "omaudit"\nversion = "0.1.0"\n\n[build-system]\nrequires = []\n',
    metadata: { ...metadata, projectType: "cli", order: 20 },
  });
  assert.equal(project.installCommand, "curl -fsSL https://raw.githubusercontent.com/omarchy-forge/omaudit/v0.1.0/install.sh | bash -s -- --version v0.1.0");
  assert.equal(project.name, "Omaudit");
  assert.equal(project.pluginId, null);
  assert.equal(project.projectType, "cli");
});

test("validates CLI package metadata and installer identity", () => {
  assert.deepEqual(validatePythonProject('[project]\nname = "omaudit"\nversion = "0.1.0"\n'), { name: "omaudit", version: "0.1.0" });
  assert.doesNotThrow(() => validateCliInstaller(
    Buffer.from('#!/usr/bin/env bash\nset -euo pipefail\nrepository="omarchy-forge/omaudit"\n'),
    "omarchy-forge/omaudit",
  ));
  assert.throws(() => validateCliInstaller(
    Buffer.from('#!/usr/bin/env bash\nset -euo pipefail\nrepository="someone/else"\n'),
    "omarchy-forge/omaudit",
  ), /identity/);
  assert.throws(() => validateCliInstaller(
    Buffer.from('#!/usr/bin/env bash\nset -euo pipefail\n# repository="omarchy-forge/omaudit"\nrepository="someone/else"\n'),
    "omarchy-forge/omaudit",
  ), /identity/);
  assert.throws(() => validateProjectMetadata({ ...metadata, projectType: "service" }), /plugin or cli/);
});

test("rejects unsafe preview paths and mismatched releases", () => {
  assert.throws(() => validateProjectMetadata({ ...metadata, previewPath: "../secret.png" }), /safe image path/);
  assert.throws(() => buildCatalogEntry({
    repository: { archived: false, disabled: false, html_url: "x", name: "handoff", owner: { login: "omarchy-forge" }, private: false, topics: ["omaforge-project"] },
    commit: { sha: "a".repeat(40) },
    release: { draft: false, prerelease: false, html_url: "x", tag_name: "v9.9.9" },
    manifest: { schemaVersion: 1, id: "org.omarchyforge.handoff", name: "Handoff", version: "0.1.1" },
    metadata,
  }), /must match/);
});

test("recognizes supported image signatures and deterministic ordering", () => {
  assert.equal(hasSupportedImageSignature(Buffer.from([137, 80, 78, 71, 13, 10, 26, 10]), "png"), true);
  assert.equal(hasSupportedImageSignature(Buffer.from("not an image"), "png"), false);
  assert.deepEqual(sortCatalog([
    { featured: false, order: 1, name: "Later" },
    { featured: true, order: 10, name: "Featured" },
    { featured: false, order: 0, name: "Earlier" },
  ]).map((item) => item.name), ["Featured", "Earlier", "Later"]);
});

test("ships a valid checked-in catalog and preview", async () => {
  const catalog = JSON.parse(await readFile(new URL("../data/projects.json", import.meta.url), "utf8"));
  assert.equal(catalog.schemaVersion, 1);
  assert.ok(catalog.projects.length > 0);
  for (const project of catalog.projects) {
    if ((project.projectType ?? "plugin") === "plugin") {
      assert.equal(project.installCommand, `omarchy plugin add https://github.com/${project.repository}.git --enable`);
    } else {
      assert.equal(project.installCommand, `curl -fsSL https://raw.githubusercontent.com/${project.repository}/v${project.version}/install.sh | bash -s -- --version v${project.version}`);
    }
    const preview = await readFile(new URL(`../public${project.preview}`, import.meta.url));
    assert.equal(hasSupportedImageSignature(preview, imageExtension(project.preview)), true);
  }
});

test("keeps generated previews outside locale proxy handling", async () => {
  const proxy = await readFile(new URL("../proxy.ts", import.meta.url), "utf8");
  assert.match(proxy, /project-images\(\?:\/\|\$\)/);
});
