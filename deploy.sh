#!/bin/bash
echo "🚀 Deploy başlatılıyor..."

# Docker Hub’dan en güncel image’ı çek
docker pull elifgider/karbonizgi:latest

# Çalışan eski container’ı durdur ve kaldır
docker stop karbonizgi-container
docker rm karbonizgi-container

# Yeni container’ı başlat
docker run -d --name karbonizgi-container -p 8080:8080 elifgider/karbonizgi:latest

echo "✅ Deploy tamamlandı!"
