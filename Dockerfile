# syntax=docker/dockerfile:1.7

# ---- Build stage ----
ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS builder

# Опционально: включить BuildKit-кэш для модулей и сборки
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache git ca-certificates

WORKDIR /src

# Сначала модули — лучше кэшируется
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Затем исходники
COPY . .
# .env в бинарник не нужен — читай его в рантайме из /app/.env
# COPY .env .env  # <-- обычно НЕ надо в билд-стадию

# Сборка (по умолчанию на musl, без CGO)
ENV CGO_ENABLED=0
# если используешь buildx и multi-arch:
# ARG TARGETOS TARGETARCH
# ENV GOOS=$TARGETOS GOARCH=$TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /go-app ./cmd/main.go

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /go-cli ./cmd/cli/cli.go

# ---- Runtime (app) ----
FROM alpine:latest AS app
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
# если нужен .env — монтируй его как файл/секрет в рантайме или COPY здесь
# COPY .env .env

COPY --from=builder /go-app /go-app
EXPOSE 3000
ENTRYPOINT ["/go-app"]

# ---- Runtime (cli) ----
FROM alpine:latest AS cli
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /go-cli /go-cli
ENTRYPOINT ["/go-cli","-command=refill"]
