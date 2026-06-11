# Chaincloud Playground

Experimental workspace for evaluating Chainguard hardened container images. Goal: stress-test
Chainguard images vs standard alternatives, focusing on network/TLS CVE surface exercised by
real outbound HTTP calls (OpenWeatherMap → .NET service, Spotify → Go service).

See [ROADMAP.md](ROADMAP.md) for phase scope and [SESSION.md](SESSION.md) for current state.

## Research References

| Resource | URL pattern | Use |
|----------|-------------|-----|
| NVD CVE detail | `https://nvd.nist.gov/vuln/detail/{CVE-ID}` | Severity (CVSS), description, affected versions |
| Chainguard advisory | `https://images.chainguard.dev/security/{CGA-ID}` | Fix status and patch notes |
| Chainguard image advisories | `https://images.chainguard.dev/directory/image/{image}/advisories` | Full advisory list per image |
| Image catalog | https://images.chainguard.dev/directory | Browse available images |

## Logging

Every file-modifying turn → append to `CLAUDE-AUDITS/yyyy-MM-dd-audit.md`. See [CLAUDE-AUDIT-RULES.md](CLAUDE-AUDIT-RULES.md).

Every accepted recommendation → append to `CLAUDE-CONSULTATIONS/yyyy-MM-dd-consultations.md`. See [CLAUDE-CONSULTATION-RULES.md](CLAUDE-CONSULTATION-RULES.md).
