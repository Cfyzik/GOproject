# Билд стадии
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Билдим приложение
RUN go build -o main .

# Финальная стадия
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник из builder стадии
COPY --from=builder /app/main .

# Копируем статические файлы
COPY static ./static/

# Создаем директорию для БД
RUN mkdir -p /data

# Настройка порта
EXPOSE 8080

# Запуск приложения
CMD ["./main"]