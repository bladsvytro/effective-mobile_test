# Этап сборки
FROM golang:1.26.2-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Генерируем Swagger документацию
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go -o docs

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Финальный этап
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/.env .

# Копируем документацию Swagger
COPY --from=builder /app/docs ./docs

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]