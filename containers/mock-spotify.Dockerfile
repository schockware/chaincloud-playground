ARG BASE_IMAGE=cgr.dev/chainguard/static:latest

FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /src
COPY . .
RUN go build -o /mock-spotify .

FROM $BASE_IMAGE
COPY --from=builder /mock-spotify /mock-spotify
EXPOSE 5200
ENTRYPOINT ["/mock-spotify"]
