# Етап збірки
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Встановлення golangci-lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.0

# Копіювання модулів Go
COPY go.mod go.sum ./
RUN go mod download

# Копіювання коду
COPY . .

# Збірка додатку
RUN go build -o api ./cmd/api

# Етап фінальний
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/api ./

EXPOSE 8080

CMD ["./api"]