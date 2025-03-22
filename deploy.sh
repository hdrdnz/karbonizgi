#!/bin/bash

# Hata oluşursa betiği durdur
set -e
REMOTE_USER=$(grep '"remoteUser"' ./config/config.json | awk -F'"' '{print $4}')
REMOTE_HOST=$(grep '"remoteHost"' ./config/config.json | awk -F'"' '{print $4}')
DOCKER_PASSWORD=$(grep '"password"' ./config/config.json | awk -F'"' '{print $4}')
DOCKER_USERNAME=$(grep '"userName"' ./config/config.json | awk -F'"' '{print $4}')
DOCKER_REPO=$(grep '"repo"' ./config/config.json | awk -F'"' '{print $4}')
CONTAINER_NAME=$(grep '"containerName"' ./config/config.json | awk -F'"' '{print $4}')
PORT_MAPPING=$(grep '"port"' ./config/config.json | awk -F'"' '{print $4}')

# Kullanıcıdan Docker Hub bilgileri al
echo "🔐 Docker Hub'a giriş yapılıyor..."
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

echo "🚀 Docker imajı oluşturuluyor..."
docker build -t $DOCKER_REPO .

echo "📤 Docker Hub'a imaj yükleniyor..."
 docker tag  $DOCKER_REPO  $DOCKER_REPO
docker push $DOCKER_REPO

# Sunucuya bağlan ve yeni imajı çekip çalıştır
ssh $REMOTE_USER@$REMOTE_HOST << EOF
set -e
echo "🚀 Sunucuya bağlanıldı."

echo "📥 Docker Hub'dan yeni imaj çekiliyor..."
docker pull $DOCKER_REPO

echo "🛑 Eski konteyner durduruluyor ve kaldırılıyor..."
docker stop $CONTAINER_NAME || true
docker rm $CONTAINER_NAME || true

echo "🚀 Yeni konteyner başlatılıyor..."
docker run -d --name $CONTAINER_NAME -p $PORT_MAPPING $DOCKER_REPO

echo "✅ Deploy işlemi başarıyla tamamlandı!"
EOF
