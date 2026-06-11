# Accelerants: Spec-First Implementation

Observed during Phase 1 Aspire scaffolding (2026-06-11).
Context: full implementation built in one session from `specs/public-api.yaml` v0.2 and
`specs/playlist-engine.yaml` v0.2.

---

## 1. Hard edge cases were pre-decided, not discovered mid-code

The passthrough error behavior — `X-Api-Source: GO-API-PASSTHROUGH`, forwarding `Retry-After`
on 429s, preserving the Go service's status code and body verbatim — would have required a
design decision during implementation without the spec. Instead it was read and built exactly
once. Same for the correlation ID rules: generate at .NET if absent, propagate to Go,
echo on all responses including errors.

**Pattern:** Behavioral rules that touch multiple layers (header names, status codes, error
bodies) are the highest-cost things to discover mid-implementation. Capturing them in the
contract up front eliminates the "stop and figure out what this should do" interrupts.

---

## 2. Topology decisions unlocked the entire wiring

The fan-out direction (.NET calls Go, Go calls Spotify, never reversed) was a single
decision recorded in the spec. It determined:
- Which project gets the `IHttpClientFactory` for Spotify (Go only)
- How the Aspire AppHost wires `WithReference` (Api → playlist-engine endpoint)
- Where `X-Correlation-Id` originates (.NET generates it, Go echoes it)
- The response assembly shape (PlaslistResponse = GoPlaylistResult + WeatherSnapshot)

Without that decision locked in, the implementation would have required either guessing
or stopping to design. One decision in the spec prevented four implementation forks.

---

## 3. Endpoints that exist for non-obvious reasons were preserved

`POST /recipe` on the Go service — derive playlist recipe without calling Spotify — is
easy to omit. Its purpose is CVE isolation: a k6 target that exercises Go + glibc without
openssl, so Spotify API latency and rate limits don't contaminate the CVE signal. That
context lives in the spec description. Without reading it, the endpoint looks redundant and
gets cut. With it, the handler and routing are built correctly on the first pass.

**Pattern:** Endpoints that exist for operational or testing reasons (not user-facing
features) are the ones most at risk of being dropped during implementation. The spec
description field is the right place to record the reason they must exist.

---

## 4. Models were derivable, not designable

Every type in both services — `WeatherCondition`, `PlaylistRecipe`, `TempoRange`,
`GeneratePlaylistRequest`, `GoPlaylistResult`, `ProblemDetails` — fell directly out of the
OpenAPI schemas. No type design decisions were made during implementation. Field names,
required vs optional, enum values, nested object shapes: all specified. This is roughly
30% of the file count done before a line of code is written.

**Pattern:** The cost of model design (naming, nullability, nesting) is invisible when it
happens in the spec but very visible when it happens during implementation as a series of
small interrupts. Moving it earlier compounds the saving across every file that uses the type.

---

## 5. Where the spec did not help: unspecified behavior

The weather-to-recipe mapping table (condition + time_of_day → mood, genres, tempo, energy,
valence, track_count) had no specification — only the enum values for inputs and outputs.
The actual mapping required judgment during implementation. This was the one place where
forward progress stopped to make design decisions.

**Pattern:** A spec that defines inputs and outputs but not the transformation is a partial
spec. For lookup tables and business rules, a decision table or mapping matrix in the spec
(or a linked ADR) moves that judgment to design time and avoids the implementation interrupt.

---

## Summary

| Accelerant | Mechanism | Approx. cost without spec |
|------------|-----------|--------------------------|
| Hard edge cases pre-decided | No mid-implementation design stops | 2–4 interrupts per edge case |
| Topology locked | Wiring, client placement, header origin all derivable | Would require rework if guessed wrong |
| Non-obvious endpoints preserved | Spec description explains the why | Likely dropped and added later |
| Models derivable from schemas | No type design during implementation | ~30% of file count requires design |
| Unspecified mapping table | N/A — still required judgment | Unavoidable without a decision table |
