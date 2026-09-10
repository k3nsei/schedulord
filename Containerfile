# ==================================================
# Builder
# ==================================================
FROM golang:1.27-alpine AS builder

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

RUN go build -o ./bin/schedulord ./cmd/schedulord

# ==================================================
# Runtime
# ==================================================
FROM gcr.io/distroless/static-debian13:nonroot

ARG VERSION=v0.0.0
ENV VERSION=${VERSION}

ARG UID=65532
ENV UID=${UID}

ARG GID=65532
ENV GID=${GID}

USER ${UID}:${GID}

WORKDIR /app

COPY --from=builder --chown=${UID}:${GID} --chmod=0755 /workspace/bin/schedulord /app/schedulord

ENTRYPOINT ["/app/schedulord"]
