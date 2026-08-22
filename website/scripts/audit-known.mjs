import { spawnSync } from "node:child_process";

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
if (advisories.length > 0) {
  console.error("Production dependency advisories found:");
  for (const advisory of advisories) {
    console.error(
      `- ${advisory.github_advisory_id ?? advisory.id}: ${advisory.module_name} (${advisory.severity})`,
    );
  }
  process.exit(1);
}

console.log("No production dependency advisories found.");
