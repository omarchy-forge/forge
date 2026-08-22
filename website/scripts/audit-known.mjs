import { spawnSync } from "node:child_process";

const allowedAdvisories = new Set([
  "GHSA-5p2g-fcmc-qvqq",
  "GHSA-g7r4-m6w7-qqqr",
  "GHSA-w3rx-r6r6-pgpr",
]);

const audit = spawnSync("pnpm", ["audit", "--prod", "--json"], {
  encoding: "utf8",
});

if (audit.error) {
  console.error(`Could not run pnpm audit: ${audit.error.message}`);
  process.exit(1);
}

let report;
try {
  report = JSON.parse(audit.stdout);
} catch {
  console.error("pnpm audit did not return valid JSON.");
  console.error(audit.stderr.trim());
  process.exit(1);
}

const advisories = Object.values(report.advisories ?? {});
const unexpected = advisories.filter(
  ({ github_advisory_id: id }) => !allowedAdvisories.has(id),
);

if (unexpected.length > 0) {
  console.error("Unexpected production dependency advisories:");
  for (const advisory of unexpected) {
    console.error(
      `- ${advisory.github_advisory_id ?? advisory.id}: ${advisory.module_name} (${advisory.severity})`,
    );
  }
  process.exit(1);
}

if (advisories.length === 0) {
  console.log("No production dependency advisories found.");
} else {
  console.warn("Only reviewed advisories listed in SECURITY.md were found:");
  for (const advisory of advisories) {
    console.warn(
      `- ${advisory.github_advisory_id}: ${advisory.module_name} (${advisory.severity})`,
    );
  }
}
