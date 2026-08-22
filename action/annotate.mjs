import { readFileSync } from "node:fs";

const reportPath = process.argv[2];
if (!reportPath) {
  console.error("usage: node annotate.mjs REPORT.json");
  process.exit(2);
}

const report = JSON.parse(readFileSync(reportPath, "utf8"));
const escapeCommand = (value) => String(value).replaceAll("%", "%25").replaceAll("\r", "%0D").replaceAll("\n", "%0A");
const escapeProperty = (value) => escapeCommand(value).replaceAll(":", "%3A").replaceAll(",", "%2C");

for (const finding of report.findings ?? []) {
  const level = finding.severity === "error" ? "error" : finding.severity === "warning" ? "warning" : "notice";
  const properties = [];
  if (finding.path) properties.push(`file=${escapeProperty(finding.path)}`);
  if (finding.line) properties.push(`line=${Number(finding.line)}`);
  properties.push(`title=${escapeProperty(`${finding.ruleId} (${finding.source})`)}`);
  const message = `${finding.message} Fix: ${finding.remediation}`;
  console.log(`::${level} ${properties.join(",")}::${escapeCommand(message)}`);
}

const summary = report.summary ?? {};
console.log(`Forge: ${summary.errors ?? 0} error(s), ${summary.warnings ?? 0} warning(s), ${summary.notes ?? 0} note(s)`);
