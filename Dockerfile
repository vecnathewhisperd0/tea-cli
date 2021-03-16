ARG GOVERSION="1.16.2"

FROM golang:${GOVERSION}-alpine AS buildenv

ARG VERSION="0.7.0"
ENV TEA_VERSION="${VERSION}"

ARG CGO_ENABLED="0"
ARG GOOS="linux"

COPY . $GOPATH/src/
WORKDIR $GOPATH/src

RUN	go get -v . && \
	go build -v -a -ldflags "-X main.Version=${TEA_VERSION}" -o /tea .

FROM scratch
ARG VERSION="0.7.0"
LABEL org.opencontainers.image.title="tea - CLI for Gitea - git with a cup of tea"
LABEL org.opencontainers.image.description="A command line tool to interact with Gitea servers"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.authors="Tamás Gérczei <tamas@gerczei.eu>"
LABEL org.opencontainers.image.vendor="The Gitea Authors"
COPY --from=buildenv /tea /
ENV HOME="/app"
ENTRYPOINT ["/tea"]
