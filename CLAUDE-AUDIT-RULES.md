# Audit Rules

**Every turn where Claude modifies, creates, or deletes files must produce an audit entry.**

Audit files live at `CLAUDE-AUDITS/yyyy-MM-dd-audit.md` (e.g. `CLAUDE-AUDITS/2026-06-10-audit.md`).
Append — never overwrite. Create the file if it does not exist.

## Entry Format

```
## HH:MM — <one-line summary of what changed>

**Prompt:** <exact user prompt that triggered this turn>

**Files modified:**
- `path/to/file` — <what changed and why>
- `path/to/other` — <what changed and why>

---
```

## Rules

- Use 24-hour local system time for the timestamp.
- List every file touched: created, edited, or deleted. Deleted files: note as `<path>` — deleted.
- If a turn involved no file changes (read-only, explanation, question), **do not** write an audit entry.
- Truncate the Prompt field at 500 characters if very long (append `…`).
- **Prepend** — insert new entries immediately after the `# Audit Log — yyyy-MM-dd` header line.
  The header is always the unique anchor; no file read is needed to find the insertion point.
  Newest entries appear at the top.
- Write the audit entry as the very last action of the turn, after all file work is done.
