import { ExternalLink, PackageCheck } from "lucide-react";
import type { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";

import { CopyInstallCommand } from "@/components/copy-install-command";
import catalog from "@/data/projects.json";

export const metadata: Metadata = {
  title: "Projects built by Omarchy Forge",
  description: "Optional Omarchy projects maintained by the Omarchy Forge organization.",
};

const ProjectsPage = () => (
  <main className="forge-projects">
    <header className="forge-projects-header">
      <div className="forge-kicker"><span /> BUILT BY OMARCHY FORGE</div>
      <h1>Projects from the forge</h1>
      <p>Optional tools and plugins maintained by our organization. Review each project before installing it.</p>
    </header>
    <section className="forge-project-list" aria-label="Owner-built projects">
      {catalog.projects.map((project) => (
        <article className="forge-project-card" key={project.slug}>
          <div className="forge-project-preview">
            <Image alt={`${project.name} preview`} fill sizes="(max-width: 760px) 100vw, 38vw" src={project.preview} />
          </div>
          <div className="forge-project-copy">
            <div className="forge-project-meta"><PackageCheck aria-hidden="true" /> {project.version} · {project.compatibility.join(", ").replaceAll("-", " ")}</div>
            <h2>{project.name}</h2>
            <p>{project.tagline}</p>
            <CopyInstallCommand command={project.installCommand} />
            <nav aria-label={`${project.name} links`}>
              <Link href={project.repositoryUrl}>Source <ExternalLink aria-hidden="true" /></Link>
              <Link href={project.releaseUrl}>Release notes <ExternalLink aria-hidden="true" /></Link>
            </nav>
          </div>
        </article>
      ))}
    </section>
    <p className="forge-projects-note">This is an owner-maintained project list, not a marketplace or an endorsement of community software.</p>
  </main>
);

export default ProjectsPage;
