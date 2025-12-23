# Database Architecture - Listings Microservice

**Дата:** 2025-11-21
**Версия:** 1.0
**Статус:** ✅ Production Ready

---

## 🎯 Архитектура БД

### ❌ НЕПРАВИЛЬНАЯ конфигурация (до исправления):

```
Listings Microservice → vondi_db (порт 5433) - МОНОЛИТ
                        └─ Пустая таблица listings
                        └─ Таблица c2c_favorites
```

**Проблема:** Микросервис подключался к монолитной БД, где таблица `listings` была пуста после миграции.

### ✅ ПРАВИЛЬНАЯ конфигурация (после исправления):

```
Monolith Backend → vondi_db (порт 5433)
                   └─ Legacy tables: c2c_favorites, c2c_categories, etc.
                   └─ Shared tables: users, balance_transactions, etc.

Listings Microservice → listings_dev_db (порт 35434)
                        ├─ listings (2 записи)
                        ├─ listing_favorites (3 записи)
                        ├─ listing_images
                        ├─ listing_locations
                        ├─ listing_attributes
                        ├─ chats
                        └─ messages
```

---

## 🗄️ База данных микросервиса

### Подключение

```bash
# PostgreSQL в Docker
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"

# Через docker exec
docker exec -it listings_postgres psql -U listings_user -d listings_dev_db
```

### Переменные окружения (.env)

```bash
# ВАЖНО: Микросервис использует ОТДЕЛЬНУЮ БД (НЕ монолит vondi_db!)
VONDILISTINGS_DB_HOST=localhost
VONDILISTINGS_DB_PORT=35434              # НЕ 5433!
VONDILISTINGS_DB_USER=listings_user       # НЕ postgres!
VONDILISTINGS_DB_PASSWORD=listings_secret
VONDILISTINGS_DB_NAME=listings_dev_db     # НЕ vondi_db!
VONDILISTINGS_DB_SSLMODE=disable
```

### Docker Container

```bash
# Проверка статуса
docker ps | grep listings_postgres

# Вывод:
# listings_postgres   postgres:15-alpine   0.0.0.0:35434->5432/tcp
```

---

## 📊 Структура таблиц

### Основные таблицы

1. **listings** - унифицированная таблица объявлений (C2C + B2C)
   - `id` BIGSERIAL PRIMARY KEY
   - `source_type` VARCHAR(10) - "c2c" или "b2c"
   - `user_id`, `title`, `description`, `price`
   - `status`, `created_at`, `updated_at`

2. **listing_favorites** - избранное пользователей
   - `user_id` + `listing_id` - composite PK
   - FK constraint на `listings.id`

3. **listing_images** - изображения объявлений
   - FK constraint на `listings.id`
   - MinIO integration

4. **chats** - чаты по объявлениям
   - FK constraint на `listings.id`

5. **messages** - сообщения в чатах
   - FK constraint на `chats.id`

---

## 🔄 Миграция данных

### История миграции

1. **Phase 11:** C2C/B2C унификация - создана таблица `listings`
2. **Phase 5:** Миграция данных из `c2c_listings` → `listings`
3. **Phase 7:** Удаление legacy таблиц из монолита

### Миграционный скрипт

```bash
cd /p/github.com/sveturs/listings
./migrate_data.sh

# Скрипт мигрирует:
# - c2c_categories
# - c2c_listings → listings
# - c2c_favorites → listing_favorites
# - c2c_images → listing_images
# - c2c_chats → chats
# - c2c_messages → messages
```

---

## 🚨 Частые проблемы

### Проблема: "listing not found"

**Причина:** Микросервис подключен к монолитной БД (vondi_db:5433) вместо своей БД.

**Решение:**
1. Проверить `.env`:
   ```bash
   cat /p/github.com/sveturs/listings/.env | grep DB_PORT
   # Должно быть: VONDILISTINGS_DB_PORT=35434
   ```

2. Перезапустить микросервис:
   ```bash
   /home/dim/.local/bin/stop-listings-microservice.sh
   /home/dim/.local/bin/start-listings-microservice.sh
   ```

### Проблема: "relation listing_favorites does not exist"

**Причина:** Таблица `listing_favorites` не создана или неправильная БД.

**Решение:**
1. Проверить таблицы:
   ```bash
   psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
     -c "\dt listing*"
   ```

2. Применить миграции (если нужно):
   ```bash
   cd /p/github.com/sveturs/listings
   # TODO: добавить migrator для микросервиса
   ```

---

## 📚 Связанные документы

- [Migration Plan](/p/github.com/sveturs/svetu/docs/migration/MIGRATION_PLAN_TO_MICROSERVICE.md)
- [Progress Tracker](/p/github.com/sveturs/svetu/docs/migration/PROGRESS.md)
- [Chat Microservice Design](/p/github.com/sveturs/CHAT_MICROSERVICE_DESIGN.md)

---

## ✅ Чеклист проверки

Перед запуском микросервиса убедись:

- [ ] `.env` указывает на порт 35434 (НЕ 5433)
- [ ] `.env` указывает на БД `listings_dev_db` (НЕ `vondi_db`)
- [ ] Docker container `listings_postgres` запущен
- [ ] Таблица `listings` содержит данные
- [ ] Таблица `listing_favorites` существует
- [ ] Микросервис запущен на порту 50053 (gRPC)

---

**Последнее обновление:** 2025-11-21
**Автор:** Database Architecture Fix (commit b89f785a)
