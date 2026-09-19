FROM gcr.io/distroless/static-debian12:nonroot

COPY miniflux-sieve /usr/local/bin/miniflux-sieve

ENTRYPOINT ["/usr/local/bin/miniflux-sieve"]
