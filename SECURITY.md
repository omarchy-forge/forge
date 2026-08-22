# Security Policy

## Reporting a vulnerability

Do not open a public issue for a suspected vulnerability. Use GitHub's private
security reporting feature for this repository. Include affected versions,
reproduction steps, impact, and any suggested mitigation. Do not include real
credentials or sensitive personal data.

No released version is supported yet. Maintainers will acknowledge a report as
soon as practical, investigate it, and coordinate disclosure after a fix or
mitigation is available.

## Security posture

Forge is local-first, has no telemetry, and requires no account or network for
normal CLI behavior. It does not execute plugin QML during static checks.
Omarchy plugins themselves run unsandboxed inside the desktop shell, so users
must review and trust plugin code before enabling it. See
`docs/SECURITY_MODEL.md` for boundaries and future requirements.

Generated projects may include a separately invoked runtime smoke test. It
requires `--trust-plugin-code`, runs only after explicit local invocation, and
is not executed by Forge checks or generated CI.

The documentation website's currently unpatched transitive advisories and
bounded mitigations are recorded in `website/SECURITY.md`.
