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

ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED}

ARG GOOS=linux
ENV GOOS=${GOOS}

ARG GOARCH=amd64
ENV GOARCH=${GOARCH}

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
