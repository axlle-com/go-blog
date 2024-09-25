FROM golang:1.22-alpine AS builder

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
COPY .env .env

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-app -ldflags="-s -w" cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-cli -ldflags="-s -w" cmd/cli/cli.go

# Используем минимальный образ для запуска
FROM alpine:latest AS app

# Устанавливаем сертификаты для Sentry
RUN apk --no-cache add ca-certificates

# Создаем директорию для логов
RUN mkdir -p /var/log/app

# Копируем собранное приложение
COPY --from=builder /go-app /go-app
#COPY --from=builder /go-cli /go-cli

# Устанавливаем рабочую директорию
WORKDIR /app

# Запуск приложения
CMD ["/go-app"]

# Образ для CLI
FROM alpine:latest AS cli

# Устанавливаем сертификаты для Sentry
RUN apk --no-cache add ca-certificates

# Создаем директорию для логов
RUN mkdir -p /var/log/app

# Копируем собранное CLI
COPY --from=builder /go-cli /go-cli

# Устанавливаем рабочую директорию
WORKDIR /app

ENTRYPOINT ["/go-cli", "-command=refill"]