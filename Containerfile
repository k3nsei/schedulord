# ==================================================
# Builder
# ==================================================
FROM golang:1.26.4-alpine AS builder

RUN apk add --no-cache \
      build-base \
      ca-certificates \
      git

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o bin/schedulord ./cmd/schedulord

# ==================================================
# Runtime
# ==================================================
FROM alpine:3.24.1

ARG VERSION=v0.0.0
ENV VERSION=${VERSION}

ARG USER=schedulord
ENV USER=${USER}

ARG GROUP=schedulord
ENV GROUP=${GROUP}

ARG UID=1000
ENV UID=${UID}

ARG GID=1000
ENV GID=${GID}

LABEL maintainer="Piotr Stępniewski"
LABEL description="Schedulord is a simple, configurable job scheduler written in Go, designed to run as a containerized service."
LABEL version="${VERSION}"

LABEL org.opencontainers.image.title="Schedulord"
LABEL org.opencontainers.image.description="Schedulord is a simple, configurable job scheduler written in Go, designed to run as a containerized service."
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.url="https://github.com/k3nsei/schedulord"
LABEL org.opencontainers.image.documentation="https://github.com/k3nsei/schedulord"
LABEL org.opencontainers.image.vendor="Piotr Stępniewski"

RUN set -eux && \
    addgroup -g ${GID} ${GROUP} && \
    adduser -D -u ${UID} -G ${GROUP} ${USER} && \
    mkdir /app && \
    chown -R ${USER}:${GROUP} /app && \
    true

RUN set -eux && \
    apk add --no-cache \
      ca-certificates \
      tzdata && \
    true

USER ${USER}:${GROUP}

WORKDIR /app

COPY --from=builder --chown=${USER}:${GROUP} --chmod=0755 /workspace/bin/schedulord /app/schedulord

ENTRYPOINT ["/app/schedulord"]
