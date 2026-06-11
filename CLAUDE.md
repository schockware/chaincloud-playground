# Chaincloud Playground

Experimental workspace for evaluating Chainguard hardened container images. The goal is to explore free-tier Chainguard images (dotnet-runtime, dotnet-sdk, Go, etc.) as drop-in replacements for standard base images, and document findings.

Image catalog reference: https://images.chainguard.dev/directory

## Research References

| Resource | URL pattern | Use |
|----------|-------------|-----|
| NVD CVE detail | `https://nvd.nist.gov/vuln/detail/{CVE-ID}` | Severity (CVSS), description, affected versions, references |
| Chainguard advisory | `https://images.chainguard.dev/security/{CGA-ID}` | Chainguard-specific fix status and patch notes |
| Chainguard image advisories | `https://images.chainguard.dev/directory/image/{image}/advisories` | Full advisory list per image |

## Consultation Log

**Any turn where Claude presents options or recommendations and the user accepts one must produce a consultation entry.**

Consultation files live at `CLAUDE-CONSULTATIONS/yyyy-MM-dd-consultations.md`. Append — never overwrite. Create the file if it does not exist.

### Consultation Entry Format

```
## HH:MM — <one-line description of the decision made>

**Context:** <what the user was trying to accomplish>

**Options presented:**
1. <option A> — <brief rationale>
2. <option B> — <brief rationale>

**Chosen:** Option <N> — <restate chosen option>

**Rationale accepted:** <why this option was recommended / why it was the right call>

---
```

### Consultation Rules

- Log the entry at the end of the turn, after any file changes.
- If Claude made only one suggestion and the user agreed (no explicit options list), still log it — note it as a single recommendation accepted rather than a multi-option choice.
- If the user rejects Claude's recommendation or provides their own direction, do **not** log a consultation entry.
- If a consultation also results in file changes, write both entries (consultation log and audit log) for that turn.

---

## Auditing Requirement

**Every turn where Claude modifies, creates, or deletes files must produce an audit entry.**

Audit files live at `CLAUDE-AUDITS/yyyy-MM-dd-audit.md` (e.g. `CLAUDE-AUDITS/2026-06-10-audit.md`). Append — never overwrite — using the format below. Create the file if it does not exist.

### Audit Entry Format

```
## HH:MM — <one-line summary of what changed>

**Prompt:** <exact user prompt that triggered this turn>

**Files modified:**
- `path/to/file` — <what changed and why>
- `path/to/other` — <what changed and why>

---
```

### Rules

- Use 24-hour time (local system time) for the timestamp.
- List every file touched: created, edited, or deleted. If a file was deleted, note it as `<path>` — deleted.
- If a turn involved no file changes (read-only, explanation, question), **do not** write an audit entry.
- The "Prompt" field should be the user's message verbatim, truncated to 500 characters if very long (append `…` if truncated).
- Append the audit entry as the very last action of the turn, after all file work is done.
