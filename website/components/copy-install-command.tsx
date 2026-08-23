"use client";

import { Check, Copy } from "lucide-react";
import { useState } from "react";

export const CopyInstallCommand = ({ command }: { command: string }) => {
  const [status, setStatus] = useState<"idle" | "copied" | "failed">("idle");

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(command);
      setStatus("copied");
      window.setTimeout(() => setStatus("idle"), 2000);
    } catch {
      setStatus("failed");
    }
  };

  return (
    <div className="forge-project-command">
      <code>{command}</code>
      <button aria-label="Copy install command" onClick={copy} type="button">
        {status === "copied" ? <Check aria-hidden="true" /> : <Copy aria-hidden="true" />}
      </button>
      <span aria-live="polite" className="sr-only">
        {status === "copied" ? "Install command copied." : status === "failed" ? "Could not copy the install command." : ""}
      </span>
    </div>
  );
};
