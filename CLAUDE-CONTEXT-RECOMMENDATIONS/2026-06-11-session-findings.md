# Context Management Findings — Session 2026-06-10/11

Recorded at end of first full working session (research + spec buildout).
Use these to configure future sessions before they start.

---

## What Burned Context Unnecessarily

### 1. HTML advisory extraction output
Three 1MB HTML files were parsed via PowerShell. The regex output (hundreds of
CVE+package+status rows) sat in context for the rest of the session even after
the data was distilled into `IMAGE-DETAILS/*/advisories.md`. Next time: do HTML
parsing in a dedicated sub-agent or a separate throwaway chat. The distilled
markdown is all that survives into the working context.

### 2. NVD fetch responses
Ten parallel NVD page fetches, each returning a full structured summary. Valuable
at research time; dead weight during spec writing. Same fix: research in its own
chat, distilled findings carried forward in the experiment manifest.

### 3. CLAUDE.md audit/consultation rules in every turn
The full format templates and rule lists loaded into context on every turn even
when no logging was happening. Extracted to separate rule files — CLAUDE.md now
references them by name only.

### 4. Log file re-reads to find append point
The audit and consultation logs were re-read multiple times during the session to
find the tail before appending. Fix: **prepend instead of append**. The file header
(`# Audit Log — yyyy-MM-dd`) is always on line 1 and is always unique — the Edit
tool can use it as the insertion anchor with zero file reads. Newest entries at top
is also the more useful reading order.

---

## Structural Changes Made This Session

| Change | File(s) | Why |
|--------|---------|-----|
| Extracted audit rules | `CLAUDE-AUDIT-RULES.md` | Removed verbose format block from CLAUDE.md |
| Extracted consultation rules | `CLAUDE-CONSULTATION-RULES.md` | Same |
| Added SESSION.md (git-ignored) | `SESSION.md` | Lightweight "where we are" primer for session starts |

---

## Chat Splitting Recommendation

Split by phase. Each phase has a clean deliverable and a natural context reset point.

| Phase | Chat focus | Starts with | Ignores |
|-------|-----------|-------------|---------|
| Research | CVE selection, NVD analysis, advisory parsing, experiment manifests | CLAUDE.md, IMAGE-DETAILS/, specs/experiments/ | src/, stress/, containers/ |
| Implementation | Aspire host, service stubs, Dockerfiles, Podman config | CLAUDE.md, specs/ (YAML + README) | IMAGE-DETAILS/, CLAUDE-AUDITS/ |
| Stress test | k6 scripting, result analysis, comparison writeup | CLAUDE.md, specs/experiments/, stress/ | src/ internals |

**Handoff protocol:** Before ending a research chat, verify SESSION.md is current.
The implementation chat reads SESSION.md + specs/ and needs nothing else to start.

---

## MCP Server — Deferred

Not warranted yet. Read + Grep reach any file in one hop and the file structure
is clean. Revisit if:
- Phase 1 produces 5+ experiments with cross-referencing needs, or
- Nightly re-scan automation (Phase 3) generates advisory diffs that need querying

---

## CLAUDE.md Size Budget

Target: under 40 lines. Anything that is a rule or format template belongs in a
referenced file, not inline. The only things that belong in CLAUDE.md directly:
- One-paragraph project description
- Research reference URL table
- Pointers to rule files (audit, consultation)
- Pointers to SESSION.md and ROADMAP.md
