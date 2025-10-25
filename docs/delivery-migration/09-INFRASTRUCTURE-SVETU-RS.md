## 🚀 Инфраструктура: Развертывание на svetu.rs

> **Источник данных**: Реальный анализ сервера svetu.rs (2025-10-22)
> **Метод**: SSH анализ через Claude Code с полными правами доступа

### 📊 Текущая архитектура сервера

**Существующие окружения**:
```
/opt/
├── svetu-authpreprod/     # Auth микросервис (Go + gRPC)
├── svetu-dev/             # Dev окружение (монолит)
└── svetu-preprod/         # Preprod окружение (монолит)
```

**Паттерн развертывания**: Docker Compose с изолированными сервисами

### 🔌 Распределение портов

**Занятые порты по окружениям**:

| Окружение | PostgreSQL | Redis | OpenSearch | HTTP | gRPC | Metrics | Health |
|-----------|------------|-------|------------|------|------|---------|--------|
| **svetu-dev** | 5433 | 6380 | 9201 | - | - | - | - |
| **svetu-preprod** | 5489 | 6382 | 9203 | 3012 | - | - | - |
| **svetu-authpreprod** | 25432 | 26379 | - | 28080 | **20051** | 29090 | 28081 |

**Свободные gRPC порты** (диапазон 50050-50060):
- ✅ `50050, 50052, 50053, 54, 55, 56, 57, 58, 59, 60`
- ❌ `50051` (занят auth-service)

**Рекомендуемые порты для delivery-preprod**:

| Сервис | Порт | Назначение |
|--------|------|------------|
| PostgreSQL | `35432` | База данных delivery |
| Redis | `36379` | Кэш и очереди |
| HTTP API | `38080` | REST API (если нужен) |
| **gRPC API** | `30051` | **Основной gRPC сервис** |
| Health Check | `38081` | Healthcheck endpoint |
| Metrics | `39090` | Prometheus metrics |

> **Примечание**: Порты в диапазоне 30000-39999 выбраны для избежания конфликтов

### 📂 Структура директории (по образцу auth-service)

```
/opt/svetu-delivery-preprod/
├── cmd/
│   └── server/              # Точка входа gRPC сервера
│       └── main.go
├── internal/
│   ├── app/                 # Инициализация приложения
│   ├── transport/           # gRPC handlers
│   │   └── grpc/
│   ├── domain/              # Доменные модели
│   ├── repository/          # PostgreSQL repos
│   │   └── postgres/
│   ├── service/             # Бизнес-логика
│   │   ├── delivery.go
│   │   ├── calculator.go
│   │   └── tracking.go
│   ├── gateway/             # Интеграции с внешними API
│   │   └── provider/
│   │       ├── interface.go
│   │       ├── factory.go
│   │       ├── postexpress/
│   │       ├── dex/
│   │       └── mock/
│   └── config/              # Конфигурация
├── pkg/                     # Публичные библиотеки
│   ├── client/              # gRPC клиент для монолита
│   └── service/             # Высокоуровневая обертка
├── deployments/
│   └── docker/
│       └── Dockerfile
├── migrations/              # SQL миграции
├── fixtures/                # Тестовые данные
├── nginx/                   # Nginx конфигурация
├── .env                     # Переменные окружения
├── .env.example             # Шаблон .env
├── docker-compose.yml       # Для локальной разработки
└── docker-compose.preprod.yml  # Для production
```

### 🐳 Docker Compose конфигурация

**Файл**: `/opt/svetu-delivery-preprod/docker-compose.preprod.yml`

```yaml
version: '3.8'

volumes:
  svetudelivery_postgres_data:
    driver: local
  svetudelivery_redis_data:
    driver: local

networks:
  svetudelivery-network:
    driver: bridge

services:
  delivery-postgres:
    image: postgres:15-alpine
    container_name: svetudelivery-postgres
    environment:
      POSTGRES_DB: ${SVETUDELIVERY_DB_NAME:-delivery_db}
      POSTGRES_USER: ${SVETUDELIVERY_DB_USER:-delivery_user}
      POSTGRES_PASSWORD: ${SVETUDELIVERY_DB_PASSWORD}
      POSTGRES_INITDB_ARGS: "--encoding=UTF8 --lc-collate=C --lc-ctype=C"
    volumes:
      - svetudelivery_postgres_data:/var/lib/postgresql/data
    ports:
      - "35432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${SVETUDELIVERY_DB_USER:-delivery_user}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    networks:
      - svetudelivery-network

  delivery-redis:
    image: redis:7-alpine
    container_name: svetudelivery-redis
    command: >
      redis-server
      --requirepass ${SVETUDELIVERY_REDIS_PASSWORD}
      --maxmemory 512mb
      --maxmemory-policy allkeys-lru
      --save 900 1
      --save 300 10
      --save 60 10000
    volumes:
      - svetudelivery_redis_data:/data
    ports:
      - "36379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "--no-auth-warning", "-a", "${SVETUDELIVERY_REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    networks:
      - svetudelivery-network

  delivery-service:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile
      args:
        GO_VERSION: "1.23"
    container_name: svetudelivery-service
    environment:
      # Service
      SVETUDELIVERY_SERVICE_NAME: ${SVETUDELIVERY_SERVICE_NAME:-delivery-service}
      SVETUDELIVERY_SERVICE_ENV: ${SVETUDELIVERY_SERVICE_ENV:-preprod}
      SVETUDELIVERY_SERVICE_LOG_LEVEL: ${SVETUDELIVERY_LOG_LEVEL:-info}

      # Server
      SVETUDELIVERY_SERVER_HTTP_PORT: 8080
      SVETUDELIVERY_SERVER_GRPC_PORT: 50052
      SVETUDELIVERY_SERVER_HEALTH_PORT: 8081
      SVETUDELIVERY_SERVER_METRICS_PORT: 9090

      # Database
      SVETUDELIVERY_DB_HOST: delivery-postgres
      SVETUDELIVERY_DB_PORT: 5432
      SVETUDELIVERY_DB_NAME: ${SVETUDELIVERY_DB_NAME:-delivery_db}
      SVETUDELIVERY_DB_USER: ${SVETUDELIVERY_DB_USER:-delivery_user}
      SVETUDELIVERY_DB_PASSWORD: ${SVETUDELIVERY_DB_PASSWORD}
      SVETUDELIVERY_DB_SSLMODE: disable

      # Redis
      SVETUDELIVERY_REDIS_HOST: delivery-redis
      SVETUDELIVERY_REDIS_PORT: 6379
      SVETUDELIVERY_REDIS_PASSWORD: ${SVETUDELIVERY_REDIS_PASSWORD}
      SVETUDELIVERY_REDIS_DB: 0

      # External APIs
      SVETUDELIVERY_POSTEXPRESS_API_URL: ${SVETUDELIVERY_POSTEXPRESS_API_URL:-https://api.postexpress.rs}
      SVETUDELIVERY_POSTEXPRESS_API_KEY: ${SVETUDELIVERY_POSTEXPRESS_API_KEY}
    ports:
      - "38080:8080"    # HTTP API (опционально)
      - "30051:50052"   # gRPC API (ОСНОВНОЙ!)
      - "38081:8081"    # Health Check
      - "39090:9090"    # Prometheus Metrics
    depends_on:
      delivery-postgres:
        condition: service_healthy
      delivery-redis:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://127.0.0.1:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    restart: unless-stopped
    networks:
      - svetudelivery-network
```

### 🔐 Переменные окружения (.env)

**Файл**: `/opt/svetu-delivery-preprod/.env`

```bash
# Service Configuration
SVETUDELIVERY_SERVICE_NAME=delivery-service
SVETUDELIVERY_SERVICE_ENV=preprod
SVETUDELIVERY_LOG_LEVEL=info

# Database Configuration
SVETUDELIVERY_DB_NAME=delivery_db
SVETUDELIVERY_DB_USER=delivery_user
SVETUDELIVERY_DB_PASSWORD=GENERATE_STRONG_PASSWORD_HERE

# Redis Configuration
SVETUDELIVERY_REDIS_PASSWORD=GENERATE_STRONG_PASSWORD_HERE

# External APIs
SVETUDELIVERY_POSTEXPRESS_API_KEY=YOUR_POST_EXPRESS_API_KEY
SVETUDELIVERY_POSTEXPRESS_API_URL=https://api.postexpress.rs

# Monitoring (optional)
SVETUDELIVERY_PROMETHEUS_ENABLED=true
SVETUDELIVERY_JAEGER_ENABLED=false
```

### 🌐 Nginx конфигурация

**Файл**: `/etc/nginx/sites-available/deliverypreprod.svetu.rs`

```nginx
# Upstream для delivery gRPC service
upstream delivery_grpc_backend {
    server 127.0.0.1:30051;
    keepalive 32;
}

# HTTP/2 для gRPC (требуется SSL)
server {
    listen 443 ssl http2;
    server_name deliverypreprod.svetu.rs;

    # SSL сертификаты (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/deliverypreprod.svetu.rs/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/deliverypreprod.svetu.rs/privkey.pem;

    # SSL оптимизация
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # gRPC специфичные настройки
    grpc_read_timeout 300s;
    grpc_send_timeout 300s;
    client_body_timeout 300s;

    # Логирование
    access_log /var/log/nginx/deliverypreprod.access.log;
    error_log /var/log/nginx/deliverypreprod.error.log;

    # gRPC endpoint
    location / {
        grpc_pass grpc://delivery_grpc_backend;

        # Headers
        grpc_set_header Host $host;
        grpc_set_header X-Real-IP $remote_addr;
        grpc_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        grpc_set_header X-Forwarded-Proto $scheme;

        # Error handling
        error_page 502 = /error502grpc;
        error_page 503 = /error503grpc;
        error_page 504 = /error504grpc;
    }

    # Health check (HTTP, не gRPC)
    location /health {
        proxy_pass http://127.0.0.1:38081/health;
        access_log off;
    }

    # Metrics (HTTP, не gRPC) - для внутреннего использования
    location /metrics {
        proxy_pass http://127.0.0.1:39090/metrics;
        allow 127.0.0.1;
        deny all;
    }

    # gRPC error responses
    location = /error502grpc {
        internal;
        default_type application/grpc;
        add_header grpc-status 14;  # UNAVAILABLE
        add_header grpc-message "Bad Gateway";
        return 204;
    }

    location = /error503grpc {
        internal;
        default_type application/grpc;
        add_header grpc-status 14;  # UNAVAILABLE
        add_header grpc-message "Service Temporarily Unavailable";
        return 204;
    }

    location = /error504grpc {
        internal;
        default_type application/grpc;
        add_header grpc-status 4;   # DEADLINE_EXCEEDED
        add_header grpc-message "Gateway Timeout";
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
TimeoutStartSec=300

[Install]
WantedBy=multi-user.target
```

```bash
# Активация
sudo systemctl daemon-reload
sudo systemctl enable delivery-preprod.service
sudo systemctl start delivery-preprod.service
```

### 🔍 Мониторинг и отладка

#### Логи

```bash
# Все сервисы
docker-compose -f docker-compose.preprod.yml logs -f

# Только delivery-service
docker-compose -f docker-compose.preprod.yml logs -f delivery-service

# PostgreSQL
docker-compose -f docker-compose.preprod.yml logs -f delivery-postgres

# Redis
docker-compose -f docker-compose.preprod.yml logs -f delivery-redis
```

#### Проверка портов

```bash
# Занятые порты
sudo netstat -tlnp | grep -E "30051|35432|36379|38080|38081|39090"

# Процессы Docker
docker ps | grep svetudelivery
```

#### Подключение к базе данных

```bash
# Из хоста
psql "postgres://delivery_user:PASSWORD@localhost:35432/delivery_db"

# Или через docker exec
docker exec -it svetudelivery-postgres psql -U delivery_user -d delivery_db
```

#### Проверка Redis

```bash
# Ping
docker exec svetudelivery-redis redis-cli -a PASSWORD ping

# Мониторинг команд
docker exec svetudelivery-redis redis-cli -a PASSWORD monitor
```

### 🚨 Troubleshooting

#### Проблема: Порт 30051 занят

```bash
# Найти процесс
sudo lsof -i :30051

# Остановить конфликтующий сервис
docker-compose -f /opt/OTHER_SERVICE/docker-compose.yml stop
```

#### Проблема: БД не поднимается

```bash
# Проверка логов
docker logs svetudelivery-postgres

# Проверка прав доступа
docker exec svetudelivery-postgres ls -la /var/lib/postgresql/data

# Пересоздание volume
docker-compose -f docker-compose.preprod.yml down -v
docker-compose -f docker-compose.preprod.yml up -d
```

#### Проблема: gRPC недоступен

```bash
# Проверка Nginx конфигурации
sudo nginx -t

# Проверка SSL сертификата
sudo certbot certificates

# Проверка firewall
sudo ufw status
```

### 📊 Ресурсы сервера

**Текущее состояние** (2025-10-22):
- **Диск**: 22GB свободно из 193GB (90% использовано)
- **Docker**: версия 27.5.1
- **Go**: версия 1.25.0

**Рекомендации**:
1. ⚠️ Мониторить место на диске (осталось мало!)
2. Настроить ротацию логов Docker
3. Очистить старые образы: `docker system prune -a`

### 🔄 Интеграция с монолитом

После развертывания микросервиса, монолит будет обращаться к нему через:

**gRPC адрес (внутренний)**: `localhost:30051`
**gRPC адрес (внешний)**: `deliverypreprod.svetu.rs:443`

**Конфигурация в монолите** (`backend/internal/config/config.go`):
```go
type DeliveryConfig struct {
    GRPCAddress string `env:"DELIVERY_GRPC_ADDRESS" envDefault:"localhost:30051"`
    UseTLS      bool   `env:"DELIVERY_USE_TLS" envDefault:"false"`
}
```

**Для preprod окружения**:
```bash
# В .env монолита
DELIVERY_GRPC_ADDRESS=localhost:30051
DELIVERY_USE_TLS=false
```

**Для production**:
```bash
DELIVERY_GRPC_ADDRESS=deliverypreprod.svetu.rs:443
DELIVERY_USE_TLS=true
```

---

## 📋 Обновленный чеклист с учетом инфраструктуры

### Фаза 0: Подготовка инфраструктуры (Week 0)
- [ ] Создать директорию `/opt/svetu-delivery-preprod`
- [ ] Получить SSL сертификат для `deliverypreprod.svetu.rs`
- [ ] Настроить Nginx конфигурацию
- [ ] Проверить свободные порты (30051, 35432, 36379, 38080-81, 39090)
- [ ] Сгенерировать пароли для БД и Redis
- [ ] Создать `.env` файл с конфигурацией
- [ ] Настроить systemd service для автозапуска

### Фаза 1: Разработка (Week 1-2)
- [ ] Proto код сгенерирован
- [ ] Domain models созданы
- [ ] Repository реализован
- [ ] Provider factory создан
- [ ] Post Express интеграция перенесена
- [ ] Service layer реализован
- [ ] gRPC handlers реализованы
- [ ] pkg/client библиотека готова
- [ ] Dockerfile создан
- [ ] docker-compose.preprod.yml настроен
- [ ] Микросервис запускается локально

### Фаза 2: Тестирование (Week 3)
- [ ] Unit tests написаны (coverage > 80%)
- [ ] Integration tests написаны
- [ ] gRPC client test работает
- [ ] Локальное тестирование пройдено
- [ ] Health checks работают
- [ ] Metrics endpoint функционирует
- [ ] Docker образ собирается успешно

### Фаза 3: Развертывание (Week 4)
- [ ] Код выгружен на сервер `/opt/svetu-delivery-preprod`
- [ ] Docker Compose запущен
- [ ] Миграции применены
- [ ] Nginx перезагружен
- [ ] Health check доступен: `curl http://localhost:38081/health`
- [ ] gRPC доступен: `grpcurl localhost:30051 list`
- [ ] Metrics доступны: `curl http://localhost:39090/metrics`
- [ ] SSL работает: `curl https://deliverypreprod.svetu.rs/health`
- [ ] Systemd service активирован
- [ ] Мониторинг логов настроен

### Фаза 4: Миграция монолита (Week 4-5)
- [ ] Старый код удален из монолита
- [ ] gRPC клиент интегрирован
- [ ] Handlers обновлены (proxy в микросервис)
- [ ] Routes обновлены
- [ ] Config обновлен (DELIVERY_GRPC_ADDRESS)
- [ ] Монолит перезапущен
- [ ] Интеграционное тестирование пройдено
- [ ] Frontend работает без изменений
- [ ] Старые таблицы удалены из БД монолита

### Фаза 5: Финализация (Week 5)
- [ ] Документация обновлена
- [ ] Runbook создан
- [ ] Smoke tests пройдены
- [ ] Метрики в Prometheus настроены
- [ ] Grafana dashboard создан
- [ ] Алерты настроены
- [ ] Резервное копирование БД настроено

---

**Обновлено**: 2025-10-22 (добавлена реальная инфраструктура svetu.rs)
