# Mock Services Specification

Two dedicated mock servers — `mock-spotify` and `mock-owm` — that mimic the external
APIs used by the real services. Both accept genuine HTTPS connections, so the TLS
handshake code path (and therefore the CVE surface) is fully exercised without
depending on external API availability or rate limits.

---

## Why dedicated external mocks

The in-process `MockWeatherService` (activated by absent `OPENWEATHERMAP_API_KEY`)
already handles OWM fallback, but it does so without making any HTTP calls. That is
useful for local dev but inappropriate for CVE stress testing — the openssl TLS
handshake never fires.

Dedicated mock servers solve this: the real services make genuine HTTPS connections
to a local process that happens to return canned data. From the openssl perspective,
there is no difference between connecting to `api.spotify.com` and connecting to
`https://localhost:5200` — the same TLS handshake code path fires on every connection.

---

## Header convention

### On requests (caller → .NET → Go)

`X-ARBITRARY-MOCK: spotify`

- Set by k6 or any caller.
- The .NET service passes it through to the Go service unchanged.
- The Go service reads it and routes Spotify API calls for **this request only** to
  `SPOTIFY_MOCK_BASE_URL` instead of `api.spotify.com`.
- Has no effect if `SPOTIFY_MOCK_BASE_URL` is not set in the environment (safety gate
  — you cannot accidentally enable mock routing without deploying the mock server).

### On responses (service → caller)

`X-ARBITRARY-MOCK` reflects what was actually mocked for this response:

| Value | Meaning |
|-------|---------|
| `weather` | OWM mock active at startup (existing behaviour) |
| `spotify` | Spotify mock used for this request |
| `weather,spotify` | Both were mocked |

The response header aggregates both signals so callers and k6 can inspect a single
header to know the full mock state.

---

## mock-spotify

**Location:** `src/mock-spotify/`
**Language:** Go (stdlib only — consistent with playlist-engine)
**Port:** 5200 (HTTPS)
**Dockerfile:** `containers/mock-spotify.Dockerfile`

### Endpoints (Spotify API surface we use)

All responses are deterministic — same input always produces the same canned output.

| Method | Path | Canned response |
|--------|------|----------------|
| `POST` | `/api/token` | `{"access_token":"mock-token","token_type":"Bearer","expires_in":3600}` |
| `GET`  | `/v1/me` | `{"id":"mock-user","display_name":"Mock User"}` |
| `GET`  | `/v1/search` | 5 canned tracks (see below) |
| `POST` | `/v1/users/{user_id}/playlists` | `{"id":"mock-playlist-id","external_urls":{"spotify":"https://open.spotify.com/playlist/mock"}}` |
| `POST` | `/v1/playlists/{playlist_id}/tracks` | `{"snapshot_id":"mock-snapshot"}` |

**Canned track set** (returned by `/v1/search` regardless of query):
```json
{
  "tracks": {
    "items": [
      {"uri": "spotify:track:mock001", "name": "Mock Track 1", "artists": [{"name": "Mock Artist"}]},
      {"uri": "spotify:track:mock002", "name": "Mock Track 2", "artists": [{"name": "Mock Artist"}]},
      {"uri": "spotify:track:mock003", "name": "Mock Track 3", "artists": [{"name": "Mock Artist"}]},
      {"uri": "spotify:track:mock004", "name": "Mock Track 4", "artists": [{"name": "Mock Artist"}]},
      {"uri": "spotify:track:mock005", "name": "Mock Track 5", "artists": [{"name": "Mock Artist"}]}
    ]
  }
}
```

### Environment variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `MOCK_TLS_CERT_FILE` | `certs/mock-spotify.crt` | Server certificate path |
| `MOCK_TLS_KEY_FILE`  | `certs/mock-spotify.key` | Server private key path |
| `MOCK_PORT`          | `5200` | Listening port |

---

## mock-owm

**Location:** `src/mock-owm/`
**Language:** Go (stdlib only)
**Port:** 5300 (HTTPS)
**Dockerfile:** `containers/mock-owm.Dockerfile`

### Endpoints (OWM API surface we use)

| Method | Path | Behaviour |
|--------|------|-----------|
| `GET`  | `/data/2.5/weather` | Returns one of 7 canned conditions, selected deterministically by `lat` + `lon` + UTC hour. Same 7-condition set as the in-process `MockWeatherService`. |

This determinism means k6 runs on the same coordinates always produce the same
weather condition, making result comparison reproducible.

### Environment variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `MOCK_TLS_CERT_FILE` | `certs/mock-owm.crt` | Server certificate path |
| `MOCK_TLS_KEY_FILE`  | `certs/mock-owm.key` | Server private key path |
| `MOCK_PORT`          | `5300` | Listening port |

---

## TLS setup

Both mock servers use HTTPS with certificates signed by a local CA. The CA is
generated once and checked in to `containers/certs/`. Client services load the
CA certificate to trust the mock servers without disabling cert verification.

> **Why not `InsecureSkipVerify`?** Skipping verification bypasses the X.509 cert
> verification code path — exactly where CVE-2022-3602 and CVE-2022-3786 live.
> A trusted local CA ensures the full verification code path fires on every
> connection, preserving the CVE surface being measured.

### Cert files

```
containers/certs/
  ca.crt           — local CA certificate (checked in — public)
  ca.key           — local CA private key  (gitignored)
  mock-spotify.crt — server cert for mock-spotify, signed by ca (checked in)
  mock-spotify.key — server private key for mock-spotify (gitignored)
  mock-owm.crt     — server cert for mock-owm, signed by ca (checked in)
  mock-owm.key     — server private key for mock-owm (gitignored)
  generate.ps1     — one-time cert generation script (see below)
```

### Generating certs (one time)

```powershell
# Run once from the repo root. Requires openssl in PATH.
./containers/certs/generate.ps1
```

The script generates a 2048-bit RSA CA and two leaf certs (one per mock service),
all with 10-year validity (this is a local test CA, not a production CA).
Private keys are never committed — `generate.ps1` is idempotent (skips files that
already exist).

### Client configuration

Services that connect to the mock servers need to trust the local CA:

| Service | Variable | Purpose |
|---------|----------|---------|
| `playlist-engine` (Go) | `SPOTIFY_TLS_CA_FILE` | Path to `ca.crt`; added to Go's root cert pool |
| `WeatherPlaylist.Api` (.NET) | `OWM_TLS_CA_FILE` | Path to `ca.crt`; added to HttpClientHandler trusted certs |

---

## Env vars on client services

These variables control where real services send API calls. Setting them to mock
server URLs enables mock mode at the URL level.

| Service | Variable | Default | Mock value |
|---------|----------|---------|-----------|
| `playlist-engine` | `SPOTIFY_BASE_URL` | `https://api.spotify.com` | (not overridden — header routing handles Spotify) |
| `playlist-engine` | `SPOTIFY_MOCK_BASE_URL` | *(unset)* | `https://localhost:5200` |
| `WeatherPlaylist.Api` | `OWM_BASE_URL` | `https://api.openweathermap.org` | `https://localhost:5300` |
| `WeatherPlaylist.Api` | `OWM_TLS_CA_FILE` | *(unset)* | `containers/certs/ca.crt` |
| `playlist-engine` | `SPOTIFY_TLS_CA_FILE` | *(unset)* | `containers/certs/ca.crt` |

**Spotify routing logic in playlist-engine:**
1. If request has `X-ARBITRARY-MOCK: spotify` **AND** `SPOTIFY_MOCK_BASE_URL` is set → use `SPOTIFY_MOCK_BASE_URL`
2. Otherwise → use `SPOTIFY_BASE_URL`

The env var is the safety gate. You cannot enable per-request mock routing without
explicitly deploying and configuring the mock server.

---

## In-service fallback (no dedicated mock server)

If `SPOTIFY_MOCK_BASE_URL` is **not** set but the request carries `X-ARBITRARY-MOCK: spotify`,
the Go service falls back to an in-process canned response (no HTTP call at all).
The response header will include `X-ARBITRARY-MOCK: spotify` to signal the fallback.

This mode is useful for local dev and schema validation. It does **not** exercise
the TLS CVE surface. Never use it as a valid EXP-001 result.

---

## Aspire integration

Both mock servers are optional Aspire resources. They are not started by default —
only when running a mock-mode stress test.

The AppHost should be extended (in the implementation phase) to accept a
`USE_MOCK_SERVERS=true` environment flag that adds these containers and wires the
env vars into the real services automatically.

---

## k6 usage

### CVE stress run with mock Spotify (no rate limits)

```powershell
# Start services with mock Spotify configured
$env:SPOTIFY_MOCK_BASE_URL = "https://localhost:5200"
$env:SPOTIFY_TLS_CA_FILE   = "containers/certs/ca.crt"
# (restart playlist-engine for env vars to take effect)

# Run k6 — header activates mock routing per request
$env:USE_SPOTIFY_MOCK = "true"   # read by EXP-001.js; adds X-ARBITRARY-MOCK: spotify to requests
k6 run --out json=stress/results/EXP-001-chainguard-mock.json stress/EXP-001.js
```

### Fully mocked run (no external credentials required)

```powershell
$env:OWM_BASE_URL          = "https://localhost:5300"
$env:OWM_TLS_CA_FILE       = "containers/certs/ca.crt"
$env:SPOTIFY_MOCK_BASE_URL = "https://localhost:5200"
$env:SPOTIFY_TLS_CA_FILE   = "containers/certs/ca.crt"
# (restart both services)

$env:USE_SPOTIFY_MOCK = "true"
k6 run --out json=stress/results/EXP-001-chainguard-fullmock.json stress/EXP-001.js
```

In this mode both TLS chains are live (real HTTPS to local mock servers). No external
credentials are required. This is the recommended setup for high-VU stress runs.
