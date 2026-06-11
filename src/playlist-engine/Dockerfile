# Build arg lets you swap the final stage image for CVE comparison experiments:
#   podman build -t playlist-engine:chainguard .
#   podman build --build-arg BASE_IMAGE=gcr.io/distroless/static-debian12 -t playlist-engine:standard .
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
