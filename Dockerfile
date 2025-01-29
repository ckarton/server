# Используем официальный образ Go
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем приложение
RUN go build -o myapp .

# Используем минимальный образ для запуска
FROM alpine:latest

# Устанавливаем необходимые зависимости
RUN apk --no-cache add ca-certificates

# Копируем скомпилированное приложение из образа builder
COPY --from=builder /app/myapp .

# Указываем порт, который будет использоваться приложением
EXPOSE 8080

# Запускаем приложение
CMD ["./myapp"]
