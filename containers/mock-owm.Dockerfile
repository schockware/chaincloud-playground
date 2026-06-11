ARG BASE_IMAGE=cgr.dev/chainguard/static:latest

FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /src
COPY . .
RUN go build -o /mock-owm .

FROM $BASE_IMAGE
COPY --from=builder /mock-owm /mock-owm
EXPOSE 5300
ENTRYPOINT ["/mock-owm"]
