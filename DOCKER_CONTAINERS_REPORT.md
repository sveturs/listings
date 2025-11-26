# 📊 Полный отчет по Docker контейнерам Listings/Svetu (10 контейнеров)

## 🗂️ Общая структура

Listings микросервис имеет **2 docker-compose файла**:
1. **База данных** (`docker-compose.yml`) - 3 контейнера
2. **Мониторинг** (`deployment/prometheus/docker-compose.yml`) - 7 контейнеров

---

## 1️⃣ БАЗА ДАННЫХ И КЭШИРОВАНИЕ (3 контейнера)

### 1.1 listings_postgres

**Образ:** `postgres:15-alpine`
**Порты:** `35434:5432` (внешний:внутренний)
**Автозапуск:** `restart: unless-stopped`

**Назначение:** Основная база данных микросервиса Orders (корзина, заказы, резервы)

**Переменные окружения:**
```bash
POSTGRES_DB=listings_dev_db
POSTGRES_USER=listings_user
POSTGRES_PASSWORD=listings_secret
```

**Volumes:**
- `postgres_data:/var/lib/postgresql/data` - данные БД
- `./migrations:/docker-entrypoint-initdb.d/migrations:ro` - миграции

**Healthcheck:**
```bash
pg_isready -U listings_user -d listings_db
# interval: 10s, timeout: 5s, retries: 5
```

**Подключение:**
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"
```

---

### 1.2 listings_redis

**Образ:** `redis:7-alpine`
**Порты:** `36380:6379`
**Автозапуск:** `restart: unless-stopped`

**Назначение:** Кэширование для микросервиса Orders

**Команда запуска:**
```bash
redis-server --requirepass redis_password --appendonly yes
```

**Volumes:**
- `redis_data:/data` - персистентные данные

**Healthcheck:**
```bash
redis-cli --raw incr ping
# interval: 10s, timeout: 5s, retries: 5
```

**Подключение:**
```bash
redis-cli -h localhost -p 36380 -a redis_password
docker exec listings_redis redis-cli FLUSHALL  # очистка
```

---

### 1.3 listings_app (НЕ запущен в данный момент)

**Образ:** Собирается из `Dockerfile`
**Container name:** `listings_app`
**Network mode:** `host` (использует host сеть для доступа к localhost PostgreSQL)
**Автозапуск:** `restart: unless-stopped`

**Назначение:** Само приложение микросервиса (gRPC сервер на порту 50053)

**Статус:** `Exited (0) 29 hours ago` - ОСТАНОВЛЕН

**Альтернативный запуск:** Через screen-сессию (см. скрипты управления ниже)

---

## 2️⃣ СИСТЕМА МОНИТОРИНГА (7 контейнеров)

### 2.1 listings-prometheus ⚠️ ПРОБЛЕМА

**Образ:** `prom/prometheus:v2.48.0`
**Порты:** `9090:9090`
**Автозапуск:** `restart: unless-stopped`
**Статус:** `Restarting (2)` - ПОСТОЯННО ПЕРЕЗАПУСКАЕТСЯ!

**Назначение:** Сбор и хранение метрик

**Проблема:**
```
Error loading config: yaml: unmarshal errors:
line 160: field relabel_configs already set in type config.ScrapeConfig
```

**Конфигурация:**
```yaml
command:
  - '--config.file=/etc/prometheus/prometheus.yml'
  - '--storage.tsdb.retention.time=15d'
  - '--storage.tsdb.retention.size=100GB'
  - '--web.enable-lifecycle'  # Позволяет reload без рестарта
  - '--web.enable-admin-api'
```

**Volumes:**
- `./prometheus.yml:/etc/prometheus/prometheus.yml:ro`
- `./alerts.yml:/etc/prometheus/alerts.yml:ro`
- `./recording_rules.yml:/etc/prometheus/recording_rules.yml:ro`
- `prometheus-data:/prometheus`

**Файл конфигурации:** `/p/github.com/sveturs/listings/deployment/prometheus/prometheus.yml:160` - ОШИБКА!

---

### 2.2 listings-grafana ✅

**Образ:** `grafana/grafana:10.2.2`
**Порты:** `3030:3000`
**Автозапуск:** `restart: unless-stopped`
**Статус:** Healthy

**Назначение:** Визуализация метрик, дашборды

**Доступ:** http://localhost:3030
**Логин:** `admin` / `admin123`

**Переменные окружения:**
```bash
GF_SECURITY_ADMIN_USER=admin
GF_SECURITY_ADMIN_PASSWORD=admin123
GF_SERVER_ROOT_URL=http://localhost:3030
GF_INSTALL_PLUGINS=grafana-piechart-panel,grafana-worldmap-panel
```

**Volumes:**
- `grafana-data:/var/lib/grafana` - данные
- `./grafana/provisioning:/etc/grafana/provisioning:ro` - datasources
- `./grafana/dashboards:/var/lib/grafana/dashboards:ro` - дашборды

---

### 2.3 listings-alertmanager ✅

**Образ:** `prom/alertmanager:v0.26.0`
**Порты:** `9093:9093`
**Автозапуск:** `restart: unless-stopped`
**Статус:** Healthy

**Назначение:** Управление алертами (роутинг, группировка, уведомления)

**Доступ:** http://localhost:9093

**Volumes:**
- `./alertmanager.yml:/etc/alertmanager/alertmanager.yml:ro`
- `alertmanager-data:/alertmanager`

---

### 2.4 listings-node-exporter ✅

**Образ:** `prom/node-exporter:v1.7.0`
**Порты:** `9100:9100`
**Автозапуск:** `restart: unless-stopped`
**Статус:** Running

**Назначение:** Экспорт системных метрик (CPU, Memory, Disk, Network)

**Команда:**
```bash
--path.procfs=/host/proc
--path.sysfs=/host/sys
--path.rootfs=/rootfs
--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)
```

**Volumes:**
- `/proc:/host/proc:ro`
- `/sys:/host/sys:ro`
- `/:/rootfs:ro`

**Метрики:** http://localhost:9100/metrics

---

### 2.5 listings-postgres-exporter ✅

**Образ:** `prometheuscommunity/postgres-exporter:v0.15.0`
**Порты:** `9187:9187`
**Автозапуск:** `restart: unless-stopped`
**Статус:** Running

**Назначение:** Экспорт метрик PostgreSQL (connections, queries, locks, transactions)

**Подключение:**
```bash
DATA_SOURCE_NAME=postgresql://postgres:mX3g1XGhMRUZEX3l@host.docker.internal:5432/svetubd?sslmode=disable
```

**Volumes:**
- `./postgres-exporter-queries.yml:/etc/postgres_exporter/queries.yml:ro`

**Метрики:** http://localhost:9187/metrics

---

### 2.6 listings-redis-exporter ✅

**Образ:** `oliver006/redis_exporter:v1.55.0`
**Порты:** `9121:9121`
**Автозапуск:** `restart: unless-stopped`
**Статус:** Running

**Назначение:** Экспорт метрик Redis (memory, commands, keys, clients)

**Подключение:**
```bash
REDIS_ADDR=host.docker.internal:6379
REDIS_PASSWORD=
```

**Метрики:** http://localhost:9121/metrics

---

### 2.7 listings-blackbox-exporter ✅

**Образ:** `prom/blackbox-exporter:v0.24.0`
**Порты:** `9115:9115`
**Автозапуск:** `restart: unless-stopped`
**Статус:** Running

**Назначение:** HTTP probing и мониторинг uptime (проверка доступности endpoints)

**Volumes:**
- `./blackbox.yml:/etc/blackbox_exporter/blackbox.yml:ro`

**Метрики:** http://localhost:9115/metrics

---

## 🛠️ УПРАВЛЕНИЕ КОНТЕЙНЕРАМИ

### Способ 1: Docker Compose (основная БД)

**Директория:** `/p/github.com/sveturs/listings/`
**Файл:** `docker-compose.yml`

```bash
# Запуск БД и Redis
cd /p/github.com/sveturs/listings
docker compose up -d

# Остановка
docker compose down

# Перезапуск
docker compose restart

# Логи
docker compose logs -f postgres
docker compose logs -f redis
```

---

### Способ 2: Docker Compose (мониторинг)

**Директория:** `/p/github.com/sveturs/listings/deployment/prometheus/`
**Файл:** `docker-compose.yml`

```bash
cd /p/github.com/sveturs/listings/deployment/prometheus

# Запуск всего стека мониторинга
make start

# Остановка
make stop

# Перезапуск
make restart

# Проверка статуса
make status

# Логи
make logs-prometheus
make logs-grafana
make logs-alertmanager

# Валидация конфигурации (ОБЯЗАТЕЛЬНО перед запуском!)
make validate

# Проверка health всех сервисов
make status

# Открыть UI в браузере
make open

# Показать URLs
make urls
```

**Другие полезные команды:**
```bash
# Проверить targets
make check-targets

# Проверить правила
make check-rules

# Тест алертов
make test-alerts

# Reload конфига без рестарта
make reload-prometheus
make reload-alertmanager

# Backup метрик
make backup

# Backup дашбордов
make backup-grafana
```

---

### Способ 3: Screen сессии (микросервис приложения)

**Скрипты:** `/home/dim/.local/bin/`

```bash
# Запуск микросервиса
/home/dim/.local/bin/start-listings-microservice.sh

# Остановка микросервиса
/home/dim/.local/bin/stop-listings-microservice.sh

# Убить процесс на порту
/home/dim/.local/bin/kill-port-50053.sh
```

**Screen сессия:**
```bash
# Подключиться к сессии
screen -r listings-microservice-50053

# Отключиться (не останавливая)
# Нажми: Ctrl+A, затем D

# Логи
tail -f /tmp/listings-microservice.log
```

---

### Способ 4: Прямое управление Docker

```bash
# Остановить контейнер
docker stop listings_postgres
docker stop listings_redis
docker stop listings-grafana

# Запустить контейнер
docker start listings_postgres
docker start listings_redis
docker start listings-grafana

# Перезапустить
docker restart listings_postgres

# Логи
docker logs -f listings_postgres
docker logs -f listings-prometheus --tail 50

# Exec команды
docker exec listings_postgres psql -U listings_user -d listings_dev_db -c "SELECT version();"
docker exec listings_redis redis-cli -a redis_password PING

# Инспекция
docker inspect listings_postgres
docker stats listings_redis
```

---

## 🔧 НАСТРОЙКА И КОНФИГУРАЦИЯ

### Переменные окружения

**Файл:** `/p/github.com/sveturs/listings/.env`

```bash
# База данных
VONDILISTINGS_DB_NAME=listings_dev_db
VONDILISTINGS_DB_USER=listings_user
VONDILISTINGS_DB_PASSWORD=listings_secret

# Redis
VONDILISTINGS_REDIS_PASSWORD=redis_password
VONDILISTINGS_REDIS_PORT=36380

# Микросервис
GRPC_PORT=50053
```

---

### Важные файлы конфигурации

**Мониторинг:**
```
/p/github.com/sveturs/listings/deployment/prometheus/
├── prometheus.yml           # Scraping config ⚠️ ОШИБКА на строке 160!
├── alerts.yml              # Alert rules
├── recording_rules.yml     # Recording rules
├── alertmanager.yml        # Alerting config
├── blackbox.yml            # HTTP probing config
└── postgres-exporter-queries.yml
```

**База данных:**
```
/p/github.com/sveturs/listings/
├── docker-compose.yml      # БД и Redis
├── .env                    # Переменные окружения
└── migrations/             # SQL миграции
```

---

## ⚠️ КРИТИЧНЫЕ ПРОБЛЕМЫ

### 1. Prometheus постоянно падает

**Причина:** Ошибка в конфигурации `/p/github.com/sveturs/listings/deployment/prometheus/prometheus.yml:160`

```
field relabel_configs already set in type config.ScrapeConfig
```

**Решение:**
```bash
cd /p/github.com/sveturs/listings/deployment/prometheus

# 1. Валидировать конфигурацию
make validate

# 2. Исправить дубликат relabel_configs на строке 160
# Открой prometheus.yml и найди дублированное поле

# 3. После исправления - перезапустить
docker restart listings-prometheus

# Или reload без рестарта (если Prometheus запустится)
make reload-prometheus
```

---

## 📊 АРХИТЕКТУРА МОНИТОРИНГА

```
┌─────────────────────────────────────────────────────────────┐
│                    МОНИТОРИНГ LISTINGS                      │
└─────────────────────────────────────────────────────────────┘

┌─────────────────┐      ┌──────────────────┐      ┌─────────┐
│  Listings App   │─────→│   Prometheus     │─────→│ Grafana │
│  :50053         │      │   :9090          │      │  :3030  │
└─────────────────┘      └──────────────────┘      └─────────┘
                                ↑                        ↑
                                │                        │
        ┌───────────────────────┼────────────────────────┤
        │                       │                        │
┌───────▼──────┐   ┌───────────▼────┐   ┌───────────────▼──────┐
│ Node Exporter│   │ Redis Exporter │   │ Postgres Exporter    │
│   :9100      │   │    :9121       │   │      :9187           │
└──────────────┘   └────────────────┘   └──────────────────────┘
                                │
                    ┌───────────▼───────────┐
                    │  Blackbox Exporter    │
                    │      :9115            │
                    └───────────────────────┘
                                │
                    ┌───────────▼───────────┐
                    │   Alertmanager        │
                    │      :9093            │
                    └───────────────────────┘
```

---

## 📝 ПОЛЕЗНЫЕ ССЫЛКИ И ЭНДПОИНТЫ

### Web UI:
- **Grafana:** http://localhost:3030 (`admin` / `admin123`)
- **Prometheus:** http://localhost:9090 (НЕ РАБОТАЕТ - падает!)
- **Alertmanager:** http://localhost:9093

### Метрики:
- **Node Exporter:** http://localhost:9100/metrics
- **Redis Exporter:** http://localhost:9121/metrics
- **Postgres Exporter:** http://localhost:9187/metrics
- **Blackbox Exporter:** http://localhost:9115/metrics

### Health checks:
- **Prometheus:** http://localhost:9090/-/healthy
- **Grafana:** http://localhost:3030/api/health
- **Alertmanager:** http://localhost:9093/-/healthy

---

## 🎯 РЕЗЮМЕ

### ✅ Работающие контейнеры (9 из 10):
1. listings_postgres
2. listings_redis
3. listings-grafana
4. listings-alertmanager
5. listings-node-exporter
6. listings-postgres-exporter
7. listings-redis-exporter
8. listings-blackbox-exporter
9. (listings_app - остановлен, но не удалён)

### ❌ Проблемные контейнеры:
10. **listings-prometheus** - постоянно падает из-за ошибки в `prometheus.yml:160`

### 🛠️ Необходимые действия:
1. **Исправить prometheus.yml** - убрать дублирующийся `relabel_configs`
2. **Перезапустить Prometheus** после исправления
3. Решить нужен ли `listings_app` контейнер или достаточно screen-сессии

---

**Дата создания:** 2025-11-18
**Автор:** Claude Code Analysis
