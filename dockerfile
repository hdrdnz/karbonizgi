# Geliştirme aşamasında kullanacağın imaj
FROM golang:1.23-alpine AS builder

# Çalışma dizini
WORKDIR /app

# Mod dosyalarını ve bağımlılıkları yükle
COPY go.mod go.sum ./
RUN go mod download

# Uygulama dosyalarını kopyala
COPY . .
COPY  data /app/data
COPY config /app/config
COPY secret.key /app/secret.key

# Uygulamayı derle
RUN go build -ldflags="-s -w" -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/config /app/config
COPY --from=builder /app/data /app/data
CMD ["./app"]
