Hi Steven,

Thanks for applying — we're excited you're interested in helping secure the software supply chain with us! Our team is reviewing your background and will be in touch if there's a strong match.

In the meantime, check out Chainguard Catalog Starter https://www.chainguard.dev/unchained/introducing-chainguard-catalog-starter — a free way to explore five trusted container images from our catalog and get hands-on with what we're building.

We'll be in touch soon!

The Chainguard Talent Team

This email is from an unmonitored inbox — please don't reply here

# Experiment: Weather → Spotify Playlist (Chainguard Image Stress Test)

Specification and contract-first implementation. Goal: stress-test Chainguard hardened
container images vs. their standard alternatives, with a focus on the network/HTTP-client
layer where CVEs tend to cluster.

## Hypothesis

Most CVEs in .NET and Go container images live at the network or HTTP client call level.
This experiment exercises those paths under real load by making genuine outbound HTTP calls
to OpenWeatherMap and the Spotify API.

## Architecture

```
Browser UI
    │
    ▼
.NET Aspire Service  ──── OpenWeatherMap API (outbound HTTP, .NET image CVE surface)
    │
    │  POST /playlist/generate  {WeatherCondition}
    ▼
Go Playlist Engine  ──── Spotify API (outbound HTTP, Go image CVE surface)
    │
    ▼  PlaylistResult  {embed_url, recipe, tracks_added}
```

- **.NET service** (`src/WeatherPlaylist.Api`) — weather ingestion, orchestration, UI host
- **Go service** (`src/playlist-engine`) — weather→mood mapping, Spotify track search + playlist creation
- **.NET AppHost** (`src/WeatherPlaylist.AppHost`) — Aspire orchestration via Podman socket

## Contract-First Specs

| File | Serves |
|------|--------|
| [`specs/public-api.yaml`](specs/public-api.yaml) | .NET service — what the UI calls |
| [`specs/playlist-engine.yaml`](specs/playlist-engine.yaml) | Go service — what .NET calls |

Both specs are OpenAPI 3.1.0. **Implementation must conform to the spec, not the other way around.**

## Images Under Test

| Role | Chainguard image | Standard alternative |
|------|-----------------|----------------------|
| .NET runtime | `cgr.dev/chainguard/dotnet-runtime` | `mcr.microsoft.com/dotnet/aspnet` |
| .NET SDK (build) | `cgr.dev/chainguard/dotnet-sdk` | `mcr.microsoft.com/dotnet/sdk` |
| Go | `cgr.dev/chainguard/go` | `golang` (Docker Hub official) |

## Local Setup

Requires: Podman, .NET 9 SDK, Go 1.22+, k6

```
# Point Aspire at Podman's Docker-compatible socket
$env:DOCKER_HOST = "npipe:////./pipe/docker_engine"   # Windows Podman Desktop default

# Copy and fill in credentials
cp .env.example .env

dotnet run --project src/WeatherPlaylist.AppHost
```

## Stress Testing

Load scenarios live in [`stress/`](stress/). Run with:

```
k6 run stress/scenario.js
```

Two container variants are tested: Chainguard images and standard images. Results are
compared for throughput, error rate, and latency under identical load.

## Credentials

Secrets are **never committed**. Use `.env` (git-ignored). See `.env.example` for
required variables:
- `OPENWEATHERMAP_API_KEY`
- `SPOTIFY_CLIENT_ID`
- `SPOTIFY_CLIENT_SECRET`
- `SPOTIFY_REFRESH_TOKEN`

## Audit Trail

All Claude-assisted changes: [`CLAUDE-AUDITS/`](CLAUDE-AUDITS/)
All Claude recommendations accepted: [`CLAUDE-CONSULTATIONS/`](CLAUDE-CONSULTATIONS/)