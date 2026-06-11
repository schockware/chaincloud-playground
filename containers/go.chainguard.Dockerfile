# Build context: src/playlist-engine/
# podman build -f containers/go.chainguard.Dockerfile -t playlist-engine:chainguard src/playlist-engine
#
# Mirrors src/playlist-engine/Dockerfile — kept here for experiment naming parity.
# Researched image: cgr.dev/chainguard/go (see IMAGE-DETAILS/go/)
ARG BASE_IMAGE=cgr.dev/chainguard/static:latest

FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -trimpath -o playlist-engine .

FROM ${BASE_IMAGE}
COPY --from=builder /app/playlist-engine /playlist-engine
EXPOSE 5100
ENTRYPOINT ["/playlist-engine"]
