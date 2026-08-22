import { ArrowRight, CheckCircle2, FileCode2, ShieldCheck, Terminal } from "lucide-react";
import type { Metadata } from "next";
import Link from "next/link";

const title = "Omarchy Forge";
const description =
  "From an empty directory to a publish-ready Omarchy plugin.";

export const metadata: Metadata = {
  title,
  description,
};

const HomePage = () => (
  <main className="forge-home">
    <section className="forge-hero">
      <div className="forge-kicker"><span /> LOCAL-FIRST PLUGIN TOOLING</div>
      <h1>Build Omarchy plugins<br /><em>with proof.</em></h1>
      <p>{description} Scaffold polished plugins, catch problems early, and validate every pull request.</p>
      <div className="forge-actions">
        <Link className="forge-primary" href="/docs/quickstart">Start building <ArrowRight aria-hidden="true" /></Link>
        <Link className="forge-secondary" href="/docs/commands">Explore commands</Link>
      </div>
      <div className="forge-terminal" aria-label="Example Forge commands">
        <div><i /><i /><i /><span>omarchy-forge / plugin</span></div>
        <code><b>$</b> omaforge init project-pulse --id dev.example.project-pulse</code>
        <code><b>$</b> omaforge check project-pulse</code>
        <code className="forge-pass"><CheckCircle2 aria-hidden="true" /> PASS <span>No error findings.</span></code>
      </div>
    </section>
    <section className="forge-grid" aria-label="Forge capabilities">
      <article><Terminal aria-hidden="true" /><small>01 / SCAFFOLD</small><h2>Start from a tested shape</h2><p>Generate a theme-aware bar widget with safe defaults, demo states, tests, and CI.</p></article>
      <article><ShieldCheck aria-hidden="true" /><small>02 / CHECK</small><h2>Make quality visible</h2><p>Run deterministic structural, security, UX, and publish-readiness rules without executing QML.</p></article>
      <article><FileCode2 aria-hidden="true" /><small>03 / SHIP</small><h2>Annotate every change</h2><p>Use schema-versioned reports and a checksum-verifying Action in pull requests.</p></article>
    </section>
  </main>
);

export default HomePage;
