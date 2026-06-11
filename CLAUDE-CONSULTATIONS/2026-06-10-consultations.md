# Consultation Log — 2026-06-10

## 17:30 — Consultation tracking format and location accepted

**Context:** User wanted a trail of turns where they accepted a Claude recommendation, mirroring the file-modification audit log.

**Options presented:**
1. Same format/location pattern as CLAUDE-AUDITS, capturing context, options offered, the chosen option, and accepted rationale — stored in `CLAUDE-CONSULTATIONS/yyyy-MM-dd-consultations.md`

**Chosen:** Option 1 — mirror the audit log pattern for consultations

**Rationale accepted:** Consistent structure across both logs makes the experiment easier to review in git history; separating consultations from file-change audits keeps each log focused on one concern.

---

## 17:35 — Architecture decisions for Chainguard experiment accepted

**Context:** Starting a spec-first, contract-first experiment using Chainguard dotnet and Go images with Aspire in Podman. User needed to lock in contract format, service topology, and local runtime approach before any files were created.

**Options presented:**
1. Contract: OpenAPI 3.1 (REST) — portable, git-diffable, k6-compatible
2. Contract: gRPC (proto) — more performant but adds tooling complexity
3. Topology: Fan-out (.NET calls Go) — realistic multi-image end-to-end stress path
4. Topology: Independent services — clean A/B comparison
5. Topology: Sidecar — Go handles a cross-cutting concern
6. Runtime: Aspire dev-mode + Podman socket via DOCKER_HOST — least friction
7. Runtime: Podman Compose only — simpler but loses Aspire dashboard
8. Runtime: Aspire with Docker fallback — reduces local setup friction

**Chosen:**
- Contract: OpenAPI 3.1 (Option 1)
- Topology: Fan-out — .NET Aspire host calls Go microservice (Option 3)
- Runtime: Aspire dev-mode + Podman socket (Option 6)
- Domain: Weather API + Spotify API → dynamic playlists with crude UI (user-specified)

**Rationale accepted:** OpenAPI keeps the public repo approachable; fan-out gives a realistic multi-image stress scenario; Podman socket integration preserves Aspire's dashboard and service-discovery while honoring the Chainguard/Podman toolchain preference.

---

## 17:42 — Spotify auth approach and weather provider accepted

**Context:** Domain locked as Weather→Spotify playlists. Two decisions needed before writing contracts: how to handle Spotify OAuth, and which weather provider to use.

**Options presented:**
1. Spotify auth: Pre-auth refresh token in config — clean for stress testing, no OAuth redirect flow to build
2. Spotify auth: Full OAuth flow in UI — more realistic but adds auth surface irrelevant to CVE experiment
3. Spotify auth: Spotify mock/stub — removes real HTTP surface, defeats the CVE hypothesis
4. Weather provider: OpenWeatherMap — free tier, API key auth, realistic HTTP surface
5. Weather provider: Open-Meteo — no API key, simpler but no auth headers in HTTP calls

**Chosen:**
- Spotify auth: Pre-auth refresh token in config (Option 1)
- Weather provider: OpenWeatherMap (Option 4)

**Rationale accepted:** Pre-auth token keeps the stress test clean and the contracts simple; OpenWeatherMap's API key auth adds a realistic authentication header to the outbound HTTP calls, which is relevant to the CVE surface being tested.

---

## 18:05 — Spec workshop Q1–Q4 and experiment tracking approach accepted

**Context:** Workshop session to resolve four open spec questions and design the experiment tracking model before Aspire scaffolding begins.

**Options presented:**
1. Q1 Persistence: cut GET /playlist/{id} from v0.1 — add state management later
2. Q1 Persistence: keep endpoint with in-memory store
3. Q2 Error passthrough: status code pass-through only
4. Q2 Error passthrough: pass-through + X-Api-Source header + X-Correlation-Id spanning both services
5. Q3 Idempotency: always create new playlist
6. Q3 Idempotency: cache by weather hash with TTL
7. Q4 UI endpoint: omit from spec entirely
8. Q4 UI endpoint: add GET / returning text/html
9. Experiment tracking: experiment_id in request body + X-Experiment-Id header + specs/experiments/ manifests linking to CVEs

**Chosen:**
- Q1: Option 1 — cut GET /playlist/{id} for now
- Q2: Option 4 — pass-through with X-Api-Source + X-Correlation-Id (user extended with correlation ID requirement)
- Q3: Option 5 — always create new
- Q4: Option 7 — omit from spec
- Experiment tracking: Option 9 — experiment_id in payload, header echo, manifest directory

**Rationale accepted:** Tighter spec = cleaner implementation. Correlation ID spanning both services enables cross-service log tracing without external tooling. Experiment manifests create a durable, repo-auditable link between stress test runs and the specific CVEs being tested.

---

## 18:47 — CVE-2023-0286 Phase 2 placement and EXP-001 CVE bundle accepted

**Context:** After deep-diving CVE-2023-0286 (X.400/CRL type confusion), user confirmed it looks like a cert DDoS vector. Needed a decision on Phase 2 vs Phase 1 placement, and which CVEs to bundle into EXP-001.

**Options presented:**
1. CVE-2023-0286: Phase 1 — configure CRL checking explicitly and build mock PKI to exercise it
2. CVE-2023-0286: Phase 2 — defer; requires mock CA + CRL server, controlled PKI environment
3. EXP-001: CVE-2024-5535 alone (CRITICAL, highest severity first)
4. EXP-001: CVE-2024-5535 + CVE-2022-3602 + CVE-2022-3786 bundled (all fire on TLS handshake, same k6 scenario)
5. EXP-001: CVE-2026-21218 (.NET spoofing — different impact class, separate experiment)

**Chosen:**
- CVE-2023-0286: Option 2 — Phase 2 with mock CA + CRL server infrastructure noted
- EXP-001: Option 4 — TLS handshake bundle (CVE-2024-5535 + CVE-2022-3602 + CVE-2022-3786)
- CVE-2026-21218 implicitly deferred to EXP-002 (different impact class: I:H vs A:H)

**Rationale accepted:** CVE-2023-0286 cannot fire without explicit CRL flag and controlled PKI — wrong fit for Phase 1's "fires naturally over HTTP" criterion. Bundling the three TLS handshake CVEs into EXP-001 maximises coverage per k6 run since they share the same trigger.

---
