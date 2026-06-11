/**
 * EXP-001 — OpenSSL TLS Handshake Bundle
 * CVEs: CVE-2024-5535 (ALPN), CVE-2022-3602, CVE-2022-3786 (X.509 cert verify)
 * See: specs/experiments/EXP-001.md
 *
 * Run:
 *   $env:EXPERIMENT_ID = "EXP-001"; $env:IMAGE_VARIANT = "chainguard"
 *   k6 run --out json=stress/results/EXP-001-chainguard.json stress/EXP-001.js
 *
 * Key comparison metric after both runs: http_req_tls_handshaking (p95/p99)
 * That is the code path where all three CVEs live.
 */

import http    from 'k6/http';
import { check, fail } from 'k6';
import { Counter }     from 'k6/metrics';

// ---------------------------------------------------------------------------
// Configuration — override via env vars before running
// ---------------------------------------------------------------------------

const BASE_DOTNET = __ENV.DOTNET_BASE_URL || 'http://localhost:5000';
const BASE_GO     = __ENV.GO_BASE_URL     || 'http://localhost:5100';
const EXPERIMENT_ID  = __ENV.EXPERIMENT_ID  || 'EXP-001';
const IMAGE_VARIANT  = __ENV.IMAGE_VARIANT  || 'unknown';

// Set USE_SPOTIFY_MOCK=true to add X-ARBITRARY-MOCK: spotify to every request.
// Requires SPOTIFY_MOCK_BASE_URL to be set on the playlist-engine — header alone
// has no effect without the env var (safety gate). See specs/mock-services.md.
const USE_SPOTIFY_MOCK = __ENV.USE_SPOTIFY_MOCK === 'true';

// ---------------------------------------------------------------------------
// Custom metrics
// ---------------------------------------------------------------------------

// Counts responses where .NET is passing through a Go service failure.
// Lets us distinguish which service degrades first under load.
const goPassthroughErrors = new Counter('go_passthrough_errors');

// Hard-fail gate: any mock-mode response invalidates the run.
// OWM TLS leg is not exercised when mock is active.
const mockModeViolations = new Counter('mock_mode_violations');

// ---------------------------------------------------------------------------
// Load profile
// ---------------------------------------------------------------------------

export const options = {
  stages: [
    { duration: '30s', target: 20 }, // ramp up — avoid cold-start noise
    { duration: '5m',  target: 20 }, // sustained — primary measurement window
    { duration: '30s', target: 0  }, // ramp down — clean drain
  ],
  thresholds: {
    // Baseline scenarios must be near-error-free
    'http_req_failed{scenario:weather_only}':    ['rate<0.01'],
    'http_req_failed{scenario:recipe_baseline}': ['rate<0.01'],
    // Full end-to-end allows up to 25% errors — Spotify 429s expected under load
    'http_req_failed{scenario:full_generate}':   ['rate<0.25'],
    // Latency gates
    'http_req_duration{service:dotnet}': ['p(95)<2000'],
    'http_req_duration{service:go}':     ['p(95)<3000'],
    // Hard fail if mock mode activates mid-run (invalidates TLS surface measurement)
    mock_mode_violations: ['count<1'],
    // NOTE: http_req_tls_handshaking has no threshold gate — it is environment-dependent.
    // It is the PRIMARY comparison metric between chainguard and standard runs.
    // Compare p95 and p99 values across the two result JSON files after both runs complete.
  },
};

// ---------------------------------------------------------------------------
// Setup — pre-flight checks (runs once before VUs start)
// ---------------------------------------------------------------------------

export function setup() {
  const res = http.get(`${BASE_DOTNET}/weather?lat=51.5074&lon=-0.1278`, {
    tags: { service: 'dotnet', scenario: 'preflight' },
  });

  if (res.status === 0) {
    fail(`Pre-flight: .NET service unreachable at ${BASE_DOTNET} — is AppHost running?`);
  }

  // Mock mode guard: if X-ARBITRARY-MOCK is present, the OWM TLS path is not
  // exercised. The run must not proceed as a valid EXP-001 measurement.
  if (res.headers['X-ARBITRARY-MOCK']) {
    fail(
      'Pre-flight: X-ARBITRARY-MOCK header present — mock weather active. ' +
      'OWM TLS leg (CVE-2024-5535 / CVE-2022-3602 / CVE-2022-3786) will not fire. ' +
      'Abort. Fix OPENWEATHERMAP_API_KEY and restart services before running EXP-001.'
    );
  }

  if (res.status !== 200) {
    fail(`Pre-flight: GET /weather returned ${res.status} — expected 200`);
  }

  console.log(`[EXP-001] Pre-flight passed. Variant=${IMAGE_VARIANT} Experiment=${EXPERIMENT_ID}`);
  return { experimentId: EXPERIMENT_ID, variant: IMAGE_VARIANT };
}

// ---------------------------------------------------------------------------
// Default function — request loop (one iteration per VU per tick)
// ---------------------------------------------------------------------------

export default function (data) {
  const roll = Math.random();

  if (roll < 0.60) {
    // 60% — full end-to-end: .NET → OWM (TLS) → Go → Spotify (TLS)
    // Exercises CVE surface on both images in one request.
    fullPlaylistGenerate(data);
  } else if (roll < 0.85) {
    // 25% — isolates the .NET → OWM TLS leg only
    weatherOnly(data);
  } else {
    // 15% — Go compute baseline: no Spotify calls, no rate-limit noise
    // Isolates Go CVE surface (glibc) from Spotify TLS variability.
    goRecipeBaseline(data);
  }
}

// ---------------------------------------------------------------------------
// Scenario implementations
// ---------------------------------------------------------------------------

function fullPlaylistGenerate(data) {
  const res = http.post(
    `${BASE_DOTNET}/playlist/generate`,
    JSON.stringify({
      lat: 51.5074,
      lon: -0.1278,
      location_label: 'London',
      experiment_id: data.experimentId,
    }),
    {
      headers: requestHeaders(data),
      tags: { service: 'dotnet', scenario: 'full_generate' },
    }
  );

  checkMockViolation(res);

  // Count Go-origin failures propagated through .NET
  if (res.headers['X-Api-Source'] === 'GO-API-PASSTHROUGH') {
    goPassthroughErrors.add(1);
  }

  check(res, {
    'full_generate: 2xx or 429': (r) =>
      (r.status >= 200 && r.status < 300) || r.status === 429,
    'full_generate: has X-Correlation-Id': (r) =>
      r.headers['X-Correlation-Id'] !== undefined,
    'full_generate: has X-Experiment-Id': (r) =>
      r.headers['X-Experiment-Id'] !== undefined,
  });
}

function weatherOnly(data) {
  const res = http.get(
    `${BASE_DOTNET}/weather?lat=51.5074&lon=-0.1278`,
    {
      headers: requestHeaders(data),
      tags: { service: 'dotnet', scenario: 'weather_only' },
    }
  );

  checkMockViolation(res);

  check(res, {
    'weather_only: status 200': (r) => r.status === 200,
    'weather_only: no mock header': (r) => !r.headers['X-ARBITRARY-MOCK'],
    'weather_only: has condition': (r) => {
      if (r.status !== 200) return true; // skip body check on non-200
      try {
        return JSON.parse(r.body).condition !== undefined;
      } catch {
        return false;
      }
    },
  });
}

function goRecipeBaseline(data) {
  // Direct Go call — compute only, no outbound HTTP, no Spotify.
  // Establishes the glibc CVE baseline and separates it from TLS noise.
  const res = http.post(
    `${BASE_GO}/recipe`,
    JSON.stringify({
      weather: {
        condition: 'Clear',
        temperature_c: 18.5,
        humidity_pct: 55,
        time_of_day: 'afternoon',
        location_label: 'London',
      },
      location_label: 'London',
      experiment_id: data.experimentId,
    }),
    {
      headers: requestHeaders(data),
      tags: { service: 'go', scenario: 'recipe_baseline' },
    }
  );

  check(res, {
    'recipe_baseline: status 200': (r) => r.status === 200,
    'recipe_baseline: has X-Correlation-Id': (r) =>
      r.headers['X-Correlation-Id'] !== undefined,
  });
}

// ---------------------------------------------------------------------------
// Teardown — runs once after all VUs finish
// ---------------------------------------------------------------------------

export function teardown(data) {
  console.log(
    `[EXP-001] Run complete. Variant=${data.variant} Experiment=${data.experimentId}`
  );
  console.log(
    '[EXP-001] Key metric to compare: http_req_tls_handshaking (p95, p99) ' +
    'against the opposite variant result file.'
  );
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function requestHeaders(data) {
  const headers = {
    'Content-Type': 'application/json',
    'X-Experiment-Id': data.experimentId,
    'X-Correlation-Id': uuidv4(),
  };
  if (USE_SPOTIFY_MOCK) {
    headers['X-ARBITRARY-MOCK'] = 'spotify';
  }
  return headers;
}

function checkMockViolation(res) {
  if (res.headers['X-ARBITRARY-MOCK']) {
    mockModeViolations.add(1);
  }
}

// UUID v4 — no external dependency needed; used for X-Correlation-Id per request
function uuidv4() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16);
  });
}
