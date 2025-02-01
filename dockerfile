# Geliştirme aşamasında kullanacağın imaj
FROM golang:1.23-alpine as dev

# Çalışma dizini
WORKDIR /app

# Mod dosyalarını ve bağımlılıkları yükle
COPY go.mod go.sum ./
RUN go mod download

# Uygulama dosyalarını kopyala
COPY . .

# Uygulamayı derle
RUN go build -o app .

# Container içinde çalıştırma
CMD ["./app"]

