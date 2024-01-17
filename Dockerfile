FROM cgr.dev/chainguard/go:latest AS build
COPY . /build/
WORKDIR /build
RUN	make build && mkdir -p /app/.config/tea

FROM cgr.dev/chainguard/busybox:latest-glibc
COPY --from=build /build/tea /bin/tea
COPY --from=build --chown=65532:65532 /app /app
VOLUME [ "/app" ]
ENV HOME="/app"
ENTRYPOINT ["/bin/sh", "-c"]
CMD [ "tea" ]
