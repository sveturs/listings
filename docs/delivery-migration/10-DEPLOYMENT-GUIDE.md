# Пошаговая инструкция развертывания
        return 204;
    }
}

# HTTP redirect to HTTPS
server {
    listen 80;
    server_name deliverypreprod.svetu.rs;
    return 301 https://$server_name$request_uri;
}
```

### 📝 Пошаговая инструкция развертывания

#### 1. Подготовка сервера

```bash
# SSH на сервер
ssh svetu@svetu.rs

# Создание директории
sudo mkdir -p /opt/svetu-delivery-preprod
sudo chown svetu:svetu /opt/svetu-delivery-preprod
cd /opt/svetu-delivery-preprod

# Клонирование репозитория
git clone git@github.com:sveturs/delivery.git .
git checkout main
```

#### 2. Настройка переменных окружения

```bash
# Копирование шаблона
cp .env.example .env

# Генерация паролей
DB_PASSWORD=$(openssl rand -base64 32)
REDIS_PASSWORD=$(openssl rand -base64 32)

# Обновление .env
sed -i "s/SVETUDELIVERY_DB_PASSWORD=.*/SVETUDELIVERY_DB_PASSWORD=$DB_PASSWORD/" .env
sed -i "s/SVETUDELIVERY_REDIS_PASSWORD=.*/SVETUDELIVERY_REDIS_PASSWORD=$REDIS_PASSWORD/" .env

# Добавление API ключей вручную
nano .env
```

#### 3. Запуск Docker Compose

```bash
# Сборка образа
docker-compose -f docker-compose.preprod.yml build

# Запуск сервисов
docker-compose -f docker-compose.preprod.yml up -d

# Проверка статуса
docker-compose -f docker-compose.preprod.yml ps

# Логи
docker-compose -f docker-compose.preprod.yml logs -f delivery-service
```

#### 4. Применение миграций

```bash
# Подключение к контейнеру
docker exec -it svetudelivery-service sh

# Применение миграций (из контейнера)
/app/migrator up

# Или через docker exec
docker exec svetudelivery-service /app/migrator up
```

#### 5. Настройка Nginx

```bash
# Копирование конфигурации
sudo cp nginx/deliverypreprod.svetu.rs.conf /etc/nginx/sites-available/
sudo ln -s /etc/nginx/sites-available/deliverypreprod.svetu.rs.conf /etc/nginx/sites-enabled/

# Получение SSL сертификата
sudo certbot certonly --nginx -d deliverypreprod.svetu.rs

# Проверка конфигурации
sudo nginx -t

# Перезагрузка Nginx
sudo systemctl reload nginx
```

#### 6. Проверка работоспособности

```bash
# Health check
curl http://localhost:38081/health

# Metrics
curl http://localhost:39090/metrics

# gRPC endpoint (через grpcurl)
grpcurl -plaintext localhost:30051 list
grpcurl -plaintext localhost:30051 delivery.v1.DeliveryService/GetShipment
```

#### 7. Настройка автозапуска

```bash
# Создание systemd service
sudo nano /etc/systemd/system/delivery-preprod.service
```

**Содержимое**:
```ini
[Unit]
Description=Delivery Microservice (Preprod)
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/svetu-delivery-preprod
ExecStart=/usr/bin/docker-compose -f docker-compose.preprod.yml up -d
ExecStop=/usr/bin/docker-compose -f docker-compose.preprod.yml down
