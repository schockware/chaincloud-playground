# Build context: src/playlist-engine/
# podman build -f containers/go.standard.Dockerfile -t playlist-engine:standard src/playlist-engine
#
# Standard alternative to cgr.dev/chainguard/go (see IMAGE-DETAILS/go/)
# Final stage uses debian:bookworm-slim to include standard OS packages (OpenSSL, glibc)
# for maximum CVE surface contrast vs the Chainguard static variant.
ARG BASE_IMAGE=debian:bookworm-slim

FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -trimpath -o playlist-engine .

FROM ${BASE_IMAGE}
COPY --from=builder /app/playlist-engine /playlist-engine
EXPOSE 5100
ENTRYPOINT ["/playlist-engine"]
