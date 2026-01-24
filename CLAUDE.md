# CLAUDE.md - Listings Microservice

## 🚀 Быстрый старт

### Запуск сервиса (рекомендуется)

```bash
make status       # Проверить статус сервиса (что запущено, на каких портах)
make deps-up      # Запустить PostgreSQL и Redis
make migrate-up   # Применить миграции
make start        # Запустить сервис в фоне (автоматически остановит предыдущий)
make stop         # Остановить сервис
```

Сервис запускается в фоне, логи пишутся в `logs/listings-service.log`.

### Управление Docker зависимостями (PostgreSQL + Redis)

```bash
make deps-up      # Запустить PostgreSQL и Redis
make deps-down    # Остановить (данные сохраняются)
make deps-clean   # Полностью удалить контейнеры и volumes (все данные!)
make deps-reset   # Чистый перезапуск: удалить всё + запустить заново
```

### Чистый перезапуск всего окружения

Если нужно начать с чистого листа (удалить все данные БД):

```bash
make reset-all    # Одна команда: stop + deps-reset + migrate-up + start
```

Или по шагам:
```bash
make stop         # Остановить сервис
make deps-reset   # Удалить и пересоздать БД
make migrate-up   # Применить миграции
make start        # Запустить сервис
```

### Порты

| Сервис | Порт | Протокол |
|--------|------|----------|
| Listings Service | 50053 | gRPC |
| Listings Service | 48086 | HTTP |
| PostgreSQL | 35434 | TCP |
| Redis | 36380 | TCP |

### Проверка кода

```bash
make format       # Форматирование
make lint         # Линтер
make test         # Тесты
make ci           # Полная проверка (deps + lint + test + build)
```

---

## 🔴 ПРАВИЛО №16: ZERO HALLUCINATION POLICY

**АБСОЛЮТНО ЗАПРЕЩЕНО использовать ЛЮБЫЕ названия без явной проверки их существования!**

### 🔍 ОБЯЗАТЕЛЬНЫЙ WORKFLOW перед написанием кода:

**1️⃣ IDENTIFY** → что нужно использовать
**2️⃣ SEARCH** → найти в коде (Grep/Read/Bash)
**3️⃣ VERIFY** → подтвердить ТОЧНОЕ название
**4️⃣ USE** → использовать ТОЛЬКО проверенное

**ЕСЛИ НЕ НАШЁЛ → НЕ ПРИДУМЫВАЙ → СПРОСИ!**

### ⚠️ ЗАПРЕЩЕНО:
- ❌ Придумывать названия функций, методов, структур
- ❌ Придумывать названия таблиц, колонок БД
- ❌ Придумывать endpoints, компоненты, переменные
- ❌ Использовать названия из других проектов
- ❌ "Улучшать" существующие названия

### 📋 CHECKLIST перед использованием:
- [ ] Я искал это в коде? (Grep/Read)
- [ ] Я прочитал файл где это используется?
- [ ] Я подтвердил ТОЧНОЕ написание?
- [ ] Это реально существует или я придумал?

**Если хоть один ответ "НЕТ" - СТОП! СНАЧАЛА НАЙДИ!**

### ❌ Примеры ГАЛЛЮЦИНАЦИЙ:

```go
// ❌ НЕПРАВИЛЬНО (придумано)
user := authSvc.GetUserByID(userID)         // Может не существовать!
SELECT user_email FROM users                 // Колонка может называться просто email!
resp, err := client.CreateNewListing(...)    // Может быть просто CreateListing!

// ✅ ПРАВИЛЬНО (проверено в коде)
user := authSvc.FetchUser(ctx, userID)       // Найдено в коде
SELECT email FROM users                      // Проверено в миграциях
resp, err := client.CreateListing(...)       // Найдено в proto файлах
```

**Цель:** Полное устранение галлюцинаций через обязательную проверку существования.

---


## 🎯 О микросервисе

**Vondi Listings Service** - микросервис для управления объявлениями, заказами, корзиной, избранным и чатами.

- **Порт gRPC:** 50053
- **Порт HTTP:** 8086
- **База данных:** `listings_dev_db` (PostgreSQL порт 35434)
- **Redis:** порт 36380
- **Директория:** `/p/github.com/sveturs/listings`

---

## 🔴 КРИТИЧЕСКИ ВАЖНО: База данных

**Микросервис использует ОТДЕЛЬНУЮ БД, а НЕ монолитную!**

### ✅ ПРАВИЛЬНАЯ конфигурация (.env):

```bash
# Database - Отдельная БД микросервиса (НЕ монолит!)
VONDILISTINGS_DB_HOST=localhost
VONDILISTINGS_DB_PORT=35434              # НЕ 5433!
VONDILISTINGS_DB_USER=listings_user      # НЕ postgres!
VONDILISTINGS_DB_PASSWORD=listings_secret
VONDILISTINGS_DB_NAME=listings_dev_db    # НЕ vondi_db!
VONDILISTINGS_DB_SSLMODE=disable
```

### ❌ НЕПРАВИЛЬНАЯ конфигурация:

```bash
# НЕ ДЕЛАЙ ТАК - это монолитная БД!
VONDILISTINGS_DB_PORT=5433     # ❌ Это монолит!
VONDILISTINGS_DB_NAME=vondi_db  # ❌ Это монолит!
VONDILISTINGS_DB_USER=postgres # ❌ Это монолит!
```

### 🐳 Docker контейнер БД:

```bash
# Проверить контейнер
docker ps | grep listings_postgres
# Вывод: listings_postgres   postgres:15-alpine   0.0.0.0:35434->5432/tcp

# Подключиться к БД
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"

# Проверить таблицы
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -c "\dt"
```

---

## 🚀 Локальный запуск с нуля

**Полный чеклист для запуска listings с чистой базы:**

```bash
# 1. Запустить инфраструктуру (PostgreSQL, Redis, OpenSearch)
docker-compose up -d postgres redis

# 2. Применить миграции
make migrate-up

# 3. Наполнить данными (категории и товары)
# Вариант А: Вставить тестовые данные напрямую в PostgreSQL
# Вариант Б: Синхронизировать из монолита (если монолит имеет данные)
# python3 scripts/sync_listings_data.py

# 4. ⚠️ ВАЖНО: Создать OpenSearch индекс и проиндексировать товары
pip3 install psycopg2-binary requests rich  # зависимости
python3 scripts/create_opensearch_index.py
python3 scripts/reindex_listings.py --target-port 35434 --target-password listings_secret --target-db listings_dev_db

# 5. Запустить микросервис
make run

# 6. Проверить
curl http://localhost:48086/health
curl http://localhost:9200/listings_microservice/_count  # должен показать количество товаров
```

### ⚠️ Важно: OpenSearch индексация

**Frontend ищет товары через OpenSearch, а не напрямую из PostgreSQL!**

Если товары есть в БД, но не отображаются на frontend:
1. Проверь индекс: `curl http://localhost:9200/listings_microservice/_count`
2. Если индекс пустой или не существует - запусти индексацию:
   ```bash
   python3 scripts/create_opensearch_index.py
   python3 scripts/reindex_listings.py --target-port 35434 --target-password listings_secret --target-db listings_dev_db
   ```

### Конфигурация монолита для OpenSearch

В `vondi/backend/.env` (монолит):
```bash
OPENSEARCH_MARKETPLACE_INDEX=listings_microservice  # индекс listings
```

---

## 🚀 Запуск и остановка

### Быстрый запуск:

```bash
# Запустить микросервис
/home/dim/.local/bin/start-listings-microservice.sh

# Остановить микросервис
/home/dim/.local/bin/stop-listings-microservice.sh

# Проверить статус
netstat -tlnp | grep ":50053"
tail -f /tmp/listings-microservice.log
```

### Ручной запуск:

```bash
# 1. Остановить старые процессы
/home/dim/.local/bin/kill-port-50053.sh

# 2. Закрыть screen сессии
screen -ls | grep listings-microservice | awk '{print $1}' | xargs -I {} screen -S {} -X quit
screen -wipe

# 3. Запустить
cd /p/github.com/sveturs/listings
screen -dmS listings-microservice-50053 bash -c 'go run ./cmd/server/main.go 2>&1 | tee /tmp/listings-microservice.log'

# 4. Проверить
netstat -tlnp | grep ":50053"
```

---

## 📋 Структура БД

### Основные таблицы:

- **listings** - унифицированная таблица объявлений (C2C + B2C)
- **listing_favorites** - избранное пользователей
- **listing_images** - изображения объявлений
- **listing_locations** - геолокация объявлений
- **listing_attributes** - атрибуты объявлений
- **categories** - категории объявлений (синхронизируется из монолита)
- **attributes** - атрибуты категорий (синхронизируется из монолита, JSONB для многоязычности)
- **category_attributes** - связи категорий и атрибутов (синхронизируется из монолита)
- **chats** - чаты по объявлениям
- **messages** - сообщения в чатах
- **chat_attachments** - вложения в сообщениях
- **orders** - заказы
- **cart_items** - корзина покупок
- **storefronts** - витрины магазинов

### 🔄 Синхронизация справочных данных

**Справочники (categories, attributes, category_attributes) синхронизируются из монолита!**

#### Быстрая синхронизация:

```bash
python3 /p/github.com/sveturs/listings/scripts/sync_listings_data.py
```

Что синхронизируется:
- ✅ **Categories:** c2c_categories → categories (прямое копирование)
- ✅ **Attributes:** unified_attributes → attributes (с трансформацией VARCHAR → JSONB)
- ✅ **Category Attributes:** unified_category_attributes → category_attributes (прямое копирование)

**Важно:** Скрипт выполняет трансформацию VARCHAR полей `name` и `display_name` в JSONB для поддержки многоязычности:
```
"Brand" → {"en": "Brand", "sr": "Brand", "ru": "Brand"}
```

📚 **Подробная документация:** [scripts/README_SYNC.md](scripts/README_SYNC.md)

#### Когда запускать синхронизацию:
- После добавления новых категорий в монолите
- После добавления новых атрибутов в монолите
- После изменения связей категория-атрибут
- При необходимости полной пересинхронизации

### Проверка данных:

```bash
# Количество объявлений
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "SELECT COUNT(*) FROM listings;"

# Количество избранных
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "SELECT COUNT(*) FROM listing_favorites;"

# Список таблиц
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "\dt"
```

---

## 🔧 Feature Flags (в монолите)

Включение микросервиса в монолите через переменные окружения:

```bash
# В /p/github.com/sveturs/svetu/backend/.env
USE_LISTINGS_MICROSERVICE=true
USE_ORDERS_MICROSERVICE=true
USE_SEARCH_MICROSERVICE=true
USE_ANALYTICS_MICROSERVICE=true
USE_CHAT_MICROSERVICE=true

LISTINGS_GRPC_URL=localhost:50053
LISTINGS_GRPC_TIMEOUT=10s
```

---

## 🧪 Тестирование

### Проверка доступности:

```bash
# Health check
curl http://localhost:8086/health

# Метрики
curl http://localhost:8086/metrics
```

### Проверка избранного:

```bash
# Получить токен
TOKEN=$(cat /tmp/token)

# Получить список избранного
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/favorites | jq '.'

# Добавить в избранное
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/favorites/11 | jq '.'

# Удалить из избранного
curl -s -X DELETE -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/marketplace/favorites/11 | jq '.'
```

---

## 🐛 Troubleshooting

### Проблема: "listing not found"

**Причина:** Микросервис подключен к монолитной БД вместо своей.

**Решение:**
1. Проверь `.env`:
   ```bash
   cat .env | grep DB_PORT
   # Должно быть: VONDILISTINGS_DB_PORT=35434
   ```

2. Исправь конфигурацию (см. раздел "База данных" выше)

3. Перезапусти микросервис:
   ```bash
   /home/dim/.local/bin/stop-listings-microservice.sh
   /home/dim/.local/bin/start-listings-microservice.sh
   ```

### Проблема: "relation listing_favorites does not exist"

**Причина:** Неправильная БД или миграции не применены.

**Решение:**
1. Проверь подключение к правильной БД (порт 35434)
2. Проверь наличие таблиц:
   ```bash
   psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
     -c "\dt listing*"
   ```

---

## 📚 Документация

- **Architecture:** [docs/DATABASE_ARCHITECTURE.md](docs/DATABASE_ARCHITECTURE.md)
- **Migration Plan:** [/p/github.com/sveturs/svetu/docs/migration/MIGRATION_PLAN_TO_MICROSERVICE.md](/p/github.com/sveturs/svetu/docs/migration/MIGRATION_PLAN_TO_MICROSERVICE.md)
- **Chat Design:** [/p/github.com/sveturs/CHAT_MICROSERVICE_DESIGN.md](/p/github.com/sveturs/CHAT_MICROSERVICE_DESIGN.md)

---

## ✅ Чеклист перед запуском

### Конфигурация
- [ ] `.env` указывает на порт 35434 (НЕ 5433)
- [ ] `.env` указывает на БД `listings_dev_db` (НЕ `vondi_db`)

### Инфраструктура
- [ ] Docker container `listings_postgres` запущен
- [ ] Redis доступен на порту 36380
- [ ] OpenSearch доступен на порту 9200
- [ ] MinIO доступен на `s3.vondi.rs`

### Данные
- [ ] Миграции применены (`make migrate-up`)
- [ ] Таблица `categories` содержит категории
- [ ] Таблица `listings` содержит товары
- [ ] ⚠️ **OpenSearch индекс создан и заполнен** (см. раздел "Локальный запуск с нуля")

### Проверка
```bash
# PostgreSQL - должны быть товары
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -c "SELECT COUNT(*) FROM listings;"

# OpenSearch - должен быть индекс с товарами
curl -s http://localhost:9200/listings_microservice/_count | jq .count

# Health check
curl http://localhost:48086/health
```

---

**Последнее обновление:** 2026-01-09
