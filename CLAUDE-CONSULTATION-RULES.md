# Consultation Rules

**Any turn where Claude presents options or recommendations and the user accepts one
must produce a consultation entry.**

Consultation files live at `CLAUDE-CONSULTATIONS/yyyy-MM-dd-consultations.md`.
Append — never overwrite. Create the file if it does not exist.

## Entry Format

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

## Rules

- Log the entry at the end of the turn, after any file changes.
- If Claude made only one suggestion and the user agreed, still log it — note as a
  single recommendation accepted rather than a multi-option choice.
- If the user rejects Claude's recommendation or provides their own direction,
  do **not** log a consultation entry.
- If a consultation also results in file changes, write both entries (consultation
  log and audit log) for that turn.
- **Prepend** — insert new entries immediately after the `# Consultation Log — yyyy-MM-dd`
  header line. The header is always the unique anchor; no file read is needed.
  Newest entries appear at the top.
