FROM gcr.io/distroless/static-debian13:nonroot

COPY miniflux-sieve /usr/local/bin/miniflux-sieve

ENTRYPOINT ["/usr/local/bin/miniflux-sieve"]
