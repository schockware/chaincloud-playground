# Chainguard Go Image — Advisory Log

**Image:** `cgr.dev/chainguard/go`
**Standard alternative:** `golang:latest` (Docker Hub official)
**Source:** https://images.chainguard.dev/directory/image/go/advisories
**Captured:** 2026-06-10

---

> **Scope note:** binutils CVEs (currently "Under investigation") and any CVEs requiring
> build-toolchain or shell-access conditions are deferred to Phase 2. See [`ROADMAP.md`](../../ROADMAP.md).
> Phase 1 targets openssl, glibc, and zlib CVEs exercised via HTTPS calls to Spotify.

---

## Summary

| Status | CVE count (unique) |
|--------|-------------------|
| Fixed | 39 |
| Not affected | 14 |
| Under investigation | 4 |
| **Total unique CVEs** | **46** |

> Severity ratings are not exposed in the Chainguard advisory UI fragment.
> Cross-reference individual CVE IDs against NVD (https://nvd.nist.gov/vuln/detail/CVE-XXXX-XXXXX) for CVSS scores.

---

## Experiment Relevance

Our Go microservice (playlist-engine) makes outbound HTTPS calls to the Spotify API.
The CVE surface most relevant to this experiment:

| Package | Why it matters |
|---------|---------------|
| **openssl** | TLS stack for all HTTPS calls — directly exercised by Spotify API requests |
| **glibc** | Core C library underlying Go's `net/http` and TLS; exercised under every HTTP call |
| **zlib** | HTTP response decompression (gzip) — used by Spotify API responses |
| **binutils** (under investigation) | Build-time toolchain; lower runtime risk but worth monitoring |

The Go image has **no Go stdlib/runtime-level CVEs** in this fragment. All CVEs are base OS packages — confirming the hypothesis that the network-layer packages (openssl, glibc) are the primary risk surface.

---

## Packages Under CVE Advisories

| Package | CVE count | Statuses |
|---------|-----------|---------|
| glibc + related (glibc-dev, glibc-locale-posix, ld-linux, libcrypt1, nss-db, nss-hesiod) | ~18 CVEs × 7 packages | Fixed / Not affected |
| openssl / libssl3 / libcrypto3 | 12 CVEs × 3 packages | Fixed |
| binutils | 10 CVEs | Fixed / Under investigation |
| busybox | 4 CVEs | Fixed |
| libstdc++ / libstdc++-dev / libgcc / libgomp / libquadmath / libatomic | 1 CVE (CVE-2023-4039) × 6 packages | Fixed |
| zlib | 2 CVEs | Fixed |
| git | 1 CVE | Fixed |

---

## Under Investigation (watch list)

| CVE | Package | Notes |
|-----|---------|-------|
| CVE-2026-3441 | binutils | Under investigation |
| CVE-2026-3442 | binutils | Under investigation |
| CVE-2026-4647 | binutils | Under investigation |
| CVE-2026-6844 | binutils | Under investigation |

These are build-toolchain CVEs. No runtime HTTP/TLS risk, but track for patch.

---

## Full Advisory Table

### Fixed

| CVE | Package |
|-----|---------|
| CVE-2023-39810 | busybox |
| CVE-2023-4039 | libatomic |
| CVE-2023-4039 | libgcc |
| CVE-2023-4039 | libgomp |
| CVE-2023-4039 | libquadmath |
| CVE-2023-4039 | libstdc++ |
| CVE-2023-4039 | libstdc++-dev |
| CVE-2024-58251 | busybox |
| CVE-2025-11187 | openssl |
| CVE-2025-11839 | binutils |
| CVE-2025-11840 | binutils |
| CVE-2025-15281 | glibc |
| CVE-2025-15467 | openssl |
| CVE-2025-15468 | openssl |
| CVE-2025-15469 | openssl |
| CVE-2025-46394 | busybox |
| CVE-2025-60876 | busybox |
| CVE-2025-66199 | openssl |
| CVE-2025-68160 | openssl |
| CVE-2025-69418 | openssl |
| CVE-2025-69419 | openssl |
| CVE-2025-69420 | openssl |
| CVE-2025-69421 | openssl |
| CVE-2025-69649 | binutils |
| CVE-2025-69650 | binutils |
| CVE-2025-69651 | binutils |
| CVE-2025-69652 | binutils |
| CVE-2025-7545 | binutils |
| CVE-2025-7546 | binutils |
| CVE-2026-0861 | glibc |
| CVE-2026-0915 | glibc |
| CVE-2026-22184 | zlib |
| CVE-2026-22795 | openssl |
| CVE-2026-22796 | openssl |
| CVE-2026-2673 | libcrypto3 |
| CVE-2026-2673 | libssl3 |
| CVE-2026-2673 | openssl |
| CVE-2026-27171 | zlib |
| CVE-2026-32631 | git |
| CVE-2026-4046 | glibc |
| CVE-2026-4046 | glibc-dev |
| CVE-2026-4046 | glibc-locale-posix |
| CVE-2026-4046 | ld-linux |
| CVE-2026-4046 | libcrypt1 |
| CVE-2026-4046 | nss-db |
| CVE-2026-4046 | nss-hesiod |
| CVE-2026-4437 | glibc |
| CVE-2026-4437 | glibc-dev |
| CVE-2026-4437 | glibc-locale-posix |
| CVE-2026-4437 | ld-linux |
| CVE-2026-4437 | libcrypt1 |
| CVE-2026-4437 | nss-db |
| CVE-2026-4437 | nss-hesiod |
| CVE-2026-4438 | glibc |
| CVE-2026-4438 | glibc-dev |
| CVE-2026-4438 | glibc-locale-posix |
| CVE-2026-4438 | ld-linux |
| CVE-2026-4438 | libcrypt1 |
| CVE-2026-4438 | nss-db |
| CVE-2026-4438 | nss-hesiod |
| CVE-2026-5450 | glibc |
| CVE-2026-5450 | glibc-dev |
| CVE-2026-5450 | glibc-locale-posix |
| CVE-2026-5450 | ld-linux |
| CVE-2026-5450 | libcrypt1 |
| CVE-2026-5450 | nss-db |
| CVE-2026-5450 | nss-hesiod |
| CVE-2026-5928 | glibc |
| CVE-2026-5928 | glibc-dev |
| CVE-2026-5928 | glibc-locale-posix |
| CVE-2026-5928 | ld-linux |
| CVE-2026-5928 | libcrypt1 |
| CVE-2026-5928 | nss-db |
| CVE-2026-5928 | nss-hesiod |

### Not Affected

| CVE | Package |
|-----|---------|
| CVE-2026-5358 | glibc |
| CVE-2026-5358 | glibc-dev |
| CVE-2026-5358 | glibc-locale-posix |
| CVE-2026-5358 | ld-linux |
| CVE-2026-5358 | libcrypt1 |
| CVE-2026-5358 | nss-db |
| CVE-2026-5358 | nss-hesiod |
| CVE-2026-5435 | glibc |
| CVE-2026-5435 | glibc-dev |
| CVE-2026-5435 | glibc-locale-posix |
| CVE-2026-5435 | ld-linux |
| CVE-2026-5435 | libcrypt1 |
| CVE-2026-5435 | nss-db |
| CVE-2026-5435 | nss-hesiod |
| CVE-2026-6238 | glibc |
| CVE-2026-6238 | glibc-dev |
| CVE-2026-6238 | glibc-locale-posix |
| CVE-2026-6238 | ld-linux |
| CVE-2026-6238 | libcrypt1 |
| CVE-2026-6238 | nss-db |
| CVE-2026-6238 | nss-hesiod |

### Under Investigation

| CVE | Package |
|-----|---------|
| CVE-2026-3441 | binutils |
| CVE-2026-3442 | binutils |
| CVE-2026-4647 | binutils |
| CVE-2026-6844 | binutils |
| CVE-2026-6846 | binutils |
