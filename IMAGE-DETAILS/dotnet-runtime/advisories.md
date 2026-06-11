# Chainguard dotnet-runtime Image — Advisory Log

**Image:** `cgr.dev/chainguard/dotnet-runtime`
**Standard alternative:** `mcr.microsoft.com/dotnet/runtime:latest`
**Source:** https://images.chainguard.dev/directory/image/dotnet-runtime/advisories
**Captured:** 2026-06-10

---

> **Scope note:** CVEs requiring privilege escalation, specific kernel configuration, or
> controlled binary execution environments are deferred to Phase 2. See [`ROADMAP.md`](../../ROADMAP.md).
> Phase 1 targets openssl, glibc network-path, zlib, and dotnet-10 CVEs only.

---

## Summary

| Status | CVE count (unique) |
|--------|-------------------|
| Fixed | 84 |
| Not affected | 30 |
| **Total unique CVEs** | **84** |

> Severity ratings are not exposed in the Chainguard advisory UI fragment.
> Cross-reference individual CVE IDs against NVD (https://nvd.nist.gov/vuln/detail/CVE-XXXX-XXXXX) for CVSS scores.

---

## Experiment Relevance

Our .NET service makes outbound HTTPS calls to OpenWeatherMap and orchestrates calls to the Go playlist-engine.
The CVE surface most relevant to this experiment:

| Package | Why it matters |
|---------|---------------|
| **openssl** | TLS stack for all HTTPS calls to OpenWeatherMap — directly in the critical path |
| **glibc** | Core C library underlying .NET's socket and HTTP stack; exercised under every outbound call |
| **dotnet-10** | Runtime-level CVE (CVE-2026-21218) — only CVE specific to the .NET runtime itself |
| **zlib** | HTTP response decompression — used by OpenWeatherMap compressed responses |

**Key finding:** `CVE-2026-21218` is the only .NET runtime-native CVE in this image.
All others are base OS package CVEs shared with the Go image.

---

## Packages Under CVE Advisories

| Package | CVE count | Statuses |
|---------|-----------|---------|
| openssl | 32 CVEs | Fixed / Not affected |
| glibc + related (glibc-locale-posix, ld-linux, libcrypt1) | ~12 CVEs × 4 packages | Fixed / Not affected |
| busybox | 4 CVEs | Fixed |
| libgcc / libstdc++ | 1 CVE (CVE-2023-4039) × 2 packages | Fixed |
| zlib | 2 CVEs | Fixed |
| dotnet-10 | 1 CVE | Fixed |

---

## Notable CVEs

| CVE | Package | Status | Significance |
|-----|---------|--------|-------------|
| **CVE-2026-21218** | dotnet-10 | Fixed | **Only .NET-runtime-native CVE** — directly in our application runtime |
| CVE-2026-2673 | openssl / libssl3 / libcrypto3 | Fixed | OpenSSL TLS — affects all outbound HTTPS |
| CVE-2026-27171 | zlib | Fixed | Decompression — affects HTTP response handling |
| CVE-2023-4911 | glibc | Fixed | "Looney Tunables" — privilege escalation via LD_PRELOAD |
| CVE-2024-2961 | glibc | Fixed | Buffer overflow in iconv — relevant under high-throughput |

---

## Full Advisory Table

### Fixed

| CVE | Package |
|-----|---------|
| CVE-2022-3358 | openssl |
| CVE-2022-3602 | openssl |
| CVE-2022-3786 | openssl |
| CVE-2022-39046 | glibc |
| CVE-2022-3996 | openssl |
| CVE-2022-4203 | openssl |
| CVE-2022-4304 | openssl |
| CVE-2022-4450 | openssl |
| CVE-2023-0215 | openssl |
| CVE-2023-0216 | openssl |
| CVE-2023-0217 | openssl |
| CVE-2023-0286 | openssl |
| CVE-2023-0401 | openssl |
| CVE-2023-0464 | openssl |
| CVE-2023-0465 | openssl |
| CVE-2023-1255 | openssl |
| CVE-2023-25139 | glibc |
| CVE-2023-2650 | openssl |
| CVE-2023-2975 | openssl |
| CVE-2023-3446 | openssl |
| CVE-2023-3817 | openssl |
| CVE-2023-39810 | busybox |
| CVE-2023-4039 | libgcc |
| CVE-2023-4039 | libstdc++ |
| CVE-2023-4527 | glibc |
| CVE-2023-4911 | glibc |
| CVE-2023-5156 | glibc |
| CVE-2023-5363 | openssl |
| CVE-2023-5678 | openssl |
| CVE-2023-6246 | glibc |
| CVE-2023-6779 | glibc |
| CVE-2023-6780 | glibc |
| CVE-2024-0727 | openssl |
| CVE-2024-12797 | openssl |
| CVE-2024-13176 | openssl |
| CVE-2024-2511 | openssl |
| CVE-2024-2961 | glibc |
| CVE-2024-33599 | glibc |
| CVE-2024-33600 | glibc |
| CVE-2024-33601 | glibc |
| CVE-2024-33602 | glibc |
| CVE-2024-4603 | openssl |
| CVE-2024-5535 | openssl |
| CVE-2024-58251 | busybox |
| CVE-2024-6119 | openssl |
| CVE-2025-0395 | glibc |
| CVE-2025-11187 | openssl |
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
| CVE-2025-8058 | glibc |
| CVE-2025-9230 | openssl |
| CVE-2025-9232 | openssl |
| CVE-2026-0861 | glibc |
| CVE-2026-21218 | dotnet-10 |
| CVE-2026-22184 | zlib |
| CVE-2026-22795 | openssl |
| CVE-2026-22796 | openssl |
| CVE-2026-2673 | libcrypto3 |
| CVE-2026-2673 | libssl3 |
| CVE-2026-2673 | openssl |
| CVE-2026-27171 | zlib |
| CVE-2026-4046 | glibc |
| CVE-2026-4046 | glibc-locale-posix |
| CVE-2026-4046 | ld-linux |
| CVE-2026-4046 | libcrypt1 |
| CVE-2026-4437 | glibc |
| CVE-2026-4437 | glibc-locale-posix |
| CVE-2026-4437 | ld-linux |
| CVE-2026-4437 | libcrypt1 |
| CVE-2026-4438 | glibc |
| CVE-2026-4438 | glibc-locale-posix |
| CVE-2026-4438 | ld-linux |
| CVE-2026-4438 | libcrypt1 |
| CVE-2026-5358 | ld-linux |
| CVE-2026-5450 | glibc |
| CVE-2026-5450 | glibc-locale-posix |
| CVE-2026-5450 | ld-linux |
| CVE-2026-5450 | libcrypt1 |
| CVE-2026-5928 | glibc |
| CVE-2026-5928 | glibc-locale-posix |
| CVE-2026-5928 | ld-linux |
| CVE-2026-5928 | libcrypt1 |

### Not Affected

| CVE | Package |
|-----|---------|
| CVE-2010-4756 | glibc |
| CVE-2019-1010022 | glibc |
| CVE-2019-1010023 | glibc |
| CVE-2019-1010024 | glibc |
| CVE-2019-1010025 | glibc |
| CVE-2023-0466 | openssl |
| CVE-2023-0687 | glibc |
| CVE-2023-4807 | openssl |
| CVE-2025-5702 | glibc |
| CVE-2025-5745 | glibc |
| CVE-2026-0915 | glibc |
| CVE-2026-5358 | glibc |
| CVE-2026-5358 | glibc-locale-posix |
| CVE-2026-5358 | libcrypt1 |
| CVE-2026-5435 | glibc |
| CVE-2026-5435 | glibc-locale-posix |
| CVE-2026-5435 | ld-linux |
| CVE-2026-5435 | libcrypt1 |
| CVE-2026-6238 | glibc |
| CVE-2026-6238 | glibc-locale-posix |
| CVE-2026-6238 | ld-linux |
| CVE-2026-6238 | libcrypt1 |
