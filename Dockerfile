# syntax=docker/dockerfile:1.7

# ---- Build stage ----
ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS builder

# Базовые пакеты
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache git ca-certificates

WORKDIR /src

# Меньше параллелизма и без CGO, чтобы не упираться в память
ENV CGO_ENABLED=0 \
    GOMAXPROCS=1

# Сначала модули — лучше кэшируется
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Затем исходники
COPY . .
# .env в бинарник не нужен — читай его в рантайме из /app/.env
# COPY .env .env  # <-- обычно НЕ надо в билд-стадию

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -p 1 -trimpath -ldflags="-s -w" -o /go-app ./cmd/main.go

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -p 1 -trimpath -ldflags="-s -w" -o /go-cli ./cmd/cli/cli.go

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
