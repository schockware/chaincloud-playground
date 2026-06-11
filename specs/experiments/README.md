# Experiment Manifests

Each file in this directory defines one CVE stress test experiment.
Filename = experiment ID. IDs are referenced in the `experiment_id` field of every
API request and the `X-Experiment-Id` header on every response.

---

## Index

| ID | CVE(s) targeted | Image variant | Status |
|----|----------------|---------------|--------|
| _(pending — add as CVEs are selected)_ | | | |

---

## Manifest Format

```markdown
# {EXP-ID} — {Short title}

## CVE Targets

| CVE | Package | NVD link | Why targeted |
|-----|---------|----------|-------------|
| CVE-XXXX-XXXXX | package | https://nvd.nist.gov/vuln/detail/CVE-XXXX-XXXXX | reason |

## Image Variants Under Test

| Role | Chainguard image | Standard image |
|------|-----------------|----------------|
| .NET runtime | cgr.dev/chainguard/dotnet-runtime:latest | mcr.microsoft.com/dotnet/runtime:latest |
| Go | cgr.dev/chainguard/go:latest | golang:latest |

## Stress Test Scenario

**Endpoint:** `POST /playlist/generate` (or specify)
**Tool:** k6
**Script:** `stress/{EXP-ID}.js`

| Parameter | Value |
|-----------|-------|
| Virtual users (VUs) | TBD |
| Duration | TBD |
| Ramp-up | TBD |
| Think time | TBD |

## CVE Surface Exercised

Explain which packages in the call path are affected by the targeted CVEs
and why the chosen endpoint exercises them.

## Success Criteria

- [ ] Chainguard image shows lower error rate than standard image under identical load
- [ ] No crashes or unexpected 5xx on Chainguard image
- [ ] Define specific thresholds (p95 latency, error rate %, etc.)

## Notes

Any additional context, gotchas, or links to raw advisory data.
```

---

## Naming Convention

`EXP-{zero-padded 3-digit number}.md` — e.g. `EXP-001.md`, `EXP-002.md`.

Experiments are immutable once a stress test run has been recorded against them.
To vary parameters, create a new experiment ID.
