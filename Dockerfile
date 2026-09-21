FROM golang:1.27-alpine AS builder

WORKDIR /build

ENV CGO_ENABLED=0

COPY . .

RUN go mod download
RUN go build .



FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=builder /build/miniflux-sieve /usr/local/bin/miniflux-sieve

ENTRYPOINT ["/usr/local/bin/miniflux-sieve"]
