FROM golang:1.23-alpine

WORKDIR /app

# Встановлення Air
RUN go install github.com/air-verse/air@latest

# Копіюємо лише файли залежностей
COPY go.mod go.sum* ./
RUN go mod download

# Не копіюємо весь код - він буде монтуватися через volumes

EXPOSE 8080

# Запуск через Air
CMD ["air", "-c", ".air.toml"]