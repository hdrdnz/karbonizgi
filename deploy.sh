#!/bin/bash
set -e

#Konfigürasyon dosyasından değişkenleri oku
REMOTE_USER=$(grep '"remoteUser"' ./config/config.json | awk -F'"' '{print $4}')
REMOTE_HOST=$(grep '"remoteHost"' ./config/config.json | awk -F'"' '{print $4}')
APP_NAME=$(grep '"containerName"' ./config/config.json | awk -F'"' '{print $4}')
PORT_MAPPING=$(grep '"port"' ./config/config.json | awk -F'"' '{print $4}')

echo "APP_NAME: $APP_NAME"

# Go projesini derle
# echo "Go projesi derleniyor..."
GOOS=linux GOARCH=amd64 go build -o $APP_NAME .

# Kopyalanacak dosyalar
COPY_FILES="$APP_NAME data config upload Dockerfile docs .env"

#  Uzak sunucuda uygulama dizini oluştur
echo "Uzak sunucuda klasör oluşturuluyor..."
ssh $REMOTE_USER@$REMOTE_HOST "mkdir -p /home/$APP_NAME"

# Dosyaları rsync ile kopyala
echo "Uygulama dosyaları rsync ile gönderiliyor..."
rsync -avz --delete $COPY_FILES $REMOTE_USER@$REMOTE_HOST:/home/$APP_NAME/

# Eski konteyner durduruluyor ve siliniyor
echo "Eski konteyner temizleniyor..."
ssh $REMOTE_USER@$REMOTE_HOST << EOF
if docker ps -a --format '{{.Names}}' | grep -Eq "^$APP_NAME\$"; then
    docker stop $APP_NAME || true
    docker rm $APP_NAME || true
else
    echo "Önceki konteyner bulunamadı, atlanıyor."
fi
EOF

#Docker image oluşturulup çalıştırılıyor
echo "Docker imajı oluşturuluyor ve konteyner başlatılıyor..."
ssh $REMOTE_USER@$REMOTE_HOST << EOF
echo "APP_NAME2: $APP_NAME"
cd /home/$APP_NAME
docker build -t $APP_NAME . || {
  echo "Build başarısız!"
  exit 1
}
docker run -dti --restart=on-failure -p $PORT_MAPPING --name $APP_NAME $APP_NAME || {
    echo "Docker konteyner başlatılamadı, image yok veya sorun var."
    exit 1
}
CONTAINER_ID=\$(docker ps -q --filter "name=$APP_NAME")
if [ ! -z "\$CONTAINER_ID" ]; then
    docker update --restart unless-stopped \$CONTAINER_ID
fi
EOF

echo "✅ Deploy işlemi başarıyla tamamlandı!"
