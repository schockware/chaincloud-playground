# Roadmap

---

## Phase 1 — Network / TLS CVE Stress Test (current scope)

Target CVEs that exercise naturally through outbound HTTP calls with no special host
configuration. The application makes real HTTPS calls to OpenWeatherMap and Spotify;
any CVE in the TLS or HTTP stack is exercised automatically under load.

**In scope:**
- openssl CVEs (all three images) — exercised via every HTTPS request
- glibc network-path CVEs (e.g. iconv buffer overflow under high throughput, nscd)
- zlib CVEs — triggered via HTTP response decompression
- CVE-2026-21218 (`dotnet-10`) — runtime-level, exercised under any application load

Experiment manifests live in [`specs/experiments/`](specs/experiments/).

---

## Phase 2 — CVEs Requiring Alternate Host / Container Setup (deferred)

These CVEs require conditions that go beyond a standard Podman-hosted service —
privilege escalation paths, specific kernel configurations, shell access, or
multi-container exploit chains. Deferring until a dedicated host environment can
be stood up for them.

**Deferred CVEs (dotnet-runtime / dotnet-sdk):**

| CVE | Package | Reason deferred |
|-----|---------|----------------|
| CVE-2023-0286 | openssl | X.400/CRL type confusion — requires `X509_V_FLAG_CRL_CHECK` explicitly enabled + mock PKI with X.400 CRL distribution points; off by default in .NET HttpClient and Go net/http. Needs controlled TLS terminator + malicious cert chain. Phase 2 setup: mock CA + custom CRL server. |
| CVE-2023-4911 | glibc | "Looney Tunables" — requires LD_PRELOAD privilege escalation scenario |
| CVE-2023-4039 | libgcc / libstdc++ | Stack protection bypass — requires controlled binary execution environment |
| CVE-2023-6246 | glibc | syslog heap overflow — requires specific process communication setup |
| CVE-2023-6779 / 6780 | glibc | `__vsyslog_internal` — same class as CVE-2023-6246 |

**Deferred CVEs (go image):**

| CVE | Package | Reason deferred |
|-----|---------|----------------|
| CVE-2026-3441 / 3442 | binutils | Under investigation; build-toolchain CVE, no runtime HTTP path |
| CVE-2026-4647 | binutils | Under investigation; same class |
| CVE-2026-6844 / 6846 | binutils | Under investigation; build-toolchain CVE |

**What Phase 2 will need:**
- A dedicated host environment (VM or bare metal) with configurable kernel settings
- Mock CA + CRL server for CVE-2023-0286 (X.400 cert chain + malicious CRL delivery)
- Possibly a second Podman host to simulate container escape / lateral movement
- Scripted exploit harness separate from the k6 stress test

---

## Phase 3 — Future (not yet scoped)

- `GET /playlist/{playlistId}` — playlist persistence and retrieval (cut from v0.1 spec)
- Full Spotify OAuth flow in the UI (cut from v0.1 for simplicity)
- SDK image supply-chain CVE testing (build pipeline attack surface)
- Automated nightly Chainguard image re-scan to catch new advisories
