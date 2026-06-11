# Stress Test Specification

Tool: **k6** — https://k6.io

---

## Container Variants

Each experiment runs twice — once against Chainguard images, once against standard
alternatives. The two builds are explicit and separate.

### Dockerfile naming convention

```
containers/
  dotnet-runtime.chainguard.Dockerfile   — cgr.dev/chainguard/dotnet-runtime
  dotnet-runtime.standard.Dockerfile     — mcr.microsoft.com/dotnet/runtime
  go.chainguard.Dockerfile               — cgr.dev/chainguard/go
  go.standard.Dockerfile                 — golang (Docker Hub official)
```

### Build commands

```powershell
# Chainguard variant
podman build -f containers/dotnet-runtime.chainguard.Dockerfile -t weather-api:chainguard src/WeatherPlaylist.Api
podman build -f containers/go.chainguard.Dockerfile             -t playlist-engine:chainguard src/playlist-engine

# Standard variant
podman build -f containers/dotnet-runtime.standard.Dockerfile   -t weather-api:standard src/WeatherPlaylist.Api
podman build -f containers/go.standard.Dockerfile               -t playlist-engine:standard src/playlist-engine
```

### Running a variant

Swap the image tags in the Aspire AppHost configuration or override via environment
before launching. Both variants use identical configuration — only the base image differs.

---

## Results Layout

```
stress/
  EXP-001.js                       — k6 scenario script for EXP-001
  results/
    EXP-001-chainguard.json        — k6 JSON output, Chainguard run
    EXP-001-standard.json          — k6 JSON output, standard run
    EXP-001-comparison.md          — written after both runs; metric diff and findings
```

Results are committed to the repo. `EXP-001-comparison.md` is the human-readable
record of what the experiment found.

---

## Run Procedure

### Pre-flight checks (run before every variant)

1. Both services healthy:
   ```
   curl -s http://localhost:5000/health | ConvertFrom-Json
   curl -s http://localhost:5100/health | ConvertFrom-Json
   ```
2. **Mock mode guard** — assert `X-ARBITRARY-MOCK` is absent from weather responses.
   If present, the OWM TLS path is not being exercised. A run in mock mode is a
   smoke test only and must not be recorded as a valid EXP-001 result.
   ```
   $h = (Invoke-WebRequest http://localhost:5000/weather?lat=51.5&lon=-0.1).Headers
   if ($h['X-ARBITRARY-MOCK']) { Write-Warning "Mock mode active — abort EXP-001" }
   ```
3. Record the OpenSSL version in the running container:
   ```
   podman exec <container-id> openssl version
   ```
   This confirms whether the standard image carries a patched or vulnerable OpenSSL.

### Running k6

```powershell
# Chainguard run
k6 run --out json=stress/results/EXP-001-chainguard.json stress/EXP-001.js

# Standard run (after swapping images and restarting services)
k6 run --out json=stress/results/EXP-001-standard.json stress/EXP-001.js
```

Tag each run with the experiment ID via environment variable — the k6 script reads it:
```powershell
$env:EXPERIMENT_ID = "EXP-001"
$env:IMAGE_VARIANT  = "chainguard"   # or "standard"
```

---

## k6 Scenario Design

Each scenario script follows this structure:

```
1. options block     — VU count, duration, thresholds
2. setup()           — assert mock mode is off; record OpenSSL version header if exposed
3. default function  — the request loop (one iteration = one experiment payload)
4. teardown()        — log summary tags
```

### EXP-001 load profile

| Phase | Duration | VUs | Purpose |
|-------|----------|-----|---------|
| Ramp up | 30s | 0 → 20 | Warm connections, avoid cold-start noise |
| Sustained | 5 min | 20 | Primary measurement window |
| Ramp down | 30s | 20 → 0 | Clean drain |

**Think time:** 0s — maximise TLS handshake rate; that is the CVE surface.

**Why 20 VUs:** Enough to create concurrent TLS handshakes across both services without
overwhelming Spotify's rate limits on the Go leg. Adjust if 429s exceed 5% of Go-direct
calls during the sustained phase.

### Request mix (EXP-001)

| Weight | Endpoint | Service | Rationale |
|--------|----------|---------|-----------|
| 60% | `POST /playlist/generate` | .NET (full end-to-end) | Primary — exercises both images |
| 25% | `GET /weather` | .NET | Isolates .NET → OWM TLS leg |
| 15% | `POST /recipe` | Go (direct) | Baseline — Go compute, no Spotify, no rate-limit noise |

The `POST /playlist/generate` on Go direct is excluded from the mix to avoid hitting
Spotify rate limits. Use it as a separate targeted run if needed.

### Thresholds (pass/fail gates)

```javascript
thresholds: {
  http_req_failed:                ['rate<0.01'],   // <1% errors overall
  'http_req_duration{service:dotnet}': ['p(95)<2000'],  // .NET p95 < 2s
  'http_req_duration{service:go}':     ['p(95)<3000'],  // Go p95 < 3s (Spotify adds latency)
}
```

---

## Metrics to Compare

After both runs, compare these from the two JSON output files:

| Metric | k6 key | What it tells us |
|--------|--------|-----------------|
| HTTP error rate | `http_req_failed` rate | Crashes / connection resets from vulnerable code |
| .NET p95 latency | `http_req_duration` (service:dotnet) | Performance under TLS load |
| Go p95 latency | `http_req_duration` (service:go) | Performance under Spotify TLS load |
| TLS handshake duration | `http_req_tls_handshaking` | Direct measurement of the CVE code path |
| Iteration failures | `iterations` failed count | Complete request failures (likely crashes) |
| `X-Api-Source: GO-API-PASSTHROUGH` frequency | custom counter | How often Go errors propagate to .NET |

The **`http_req_tls_handshaking`** metric is the most direct signal for EXP-001 — it
measures time spent in exactly the code path where CVE-2024-5535, CVE-2022-3602, and
CVE-2022-3786 live.

---

## Interpreting Results

| Observation | Likely meaning |
|-------------|---------------|
| Standard image has higher `http_req_failed` rate | Unpatched OpenSSL crashing under load |
| Standard image has higher `http_req_tls_handshaking` | Slower/unstable TLS path in vulnerable version |
| Both images identical | Standard image carries patched OpenSSL — record version; note finding |
| 429s on Go leg exceed 5% | Spotify rate-limiting before CVE surface is the bottleneck; switch to `POST /recipe` as primary target |
| `X-ARBITRARY-MOCK` present in any response | Abort — OWM mock is active; .NET TLS leg not exercised |

---

## Adding a New Experiment

1. Create `specs/experiments/EXP-NNN.md` (see manifest format in `specs/experiments/README.md`)
2. Create `stress/EXP-NNN.js` following the structure above
3. Run against both image variants; save results to `stress/results/EXP-NNN-chainguard.json`
   and `stress/results/EXP-NNN-standard.json`
4. Write `stress/results/EXP-NNN-comparison.md`
5. Update `specs/experiments/README.md` index
