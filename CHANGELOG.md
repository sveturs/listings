# Changelog - Listings Microservice

Все значимые изменения в этом проекте документируются в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/).

## [Unreleased]

### Fixed - 2025-12-29 (c3ade748a)

**Исправлена обработка ошибок удаления из корзины + скрипты очистки + CI lint fix**

#### CI Fix (c3ade748a)
- Удалена unused функция `enrichCategoryFromDB` из `internal/service/category_detection_service.go`
- Исправлен lint warning: "func (*CategoryDetectionService).enrichCategoryFromDB is unused"

#### Original Fix (b8cc01329)

#### Изменения в коде

**Error handling:**
- `internal/service/errors.go`:
  - Добавлен `ErrCartItemNotFound` в функцию `IsNotFoundError()`
  - Теперь при попытке удалить несуществующий cart_item возвращается **404 Not Found** вместо **500 Internal Server Error**
  - Правильная обработка через `mapServiceErrorToGRPC()`

**Новые скрипты:**

1. `scripts/cleanup_anonymous_carts.sql`:
   - SQL скрипт для удаления анонимных корзин (user_id IS NULL) старше 7 дней
   - Предотвращает накопление orphan cart_items в БД
   - Использование: `psql "postgres://..." -f cleanup_anonymous_carts.sql`

2. `scripts/cleanup_anonymous_carts.sh`:
   - Bash wrapper для SQL скрипта
   - Автоматическая загрузка credentials из ENV
   - Может быть добавлен в cron для автоматической очистки:
     ```bash
     # Очистка каждый день в 3:00
     0 3 * * * /path/to/cleanup_anonymous_carts.sh
     ```

#### Изменения в БД

**База данных:** `listings_dev_db` (PostgreSQL порт 35434)

**Выполнено вручную (data cleanup):**
```sql
-- Удалены анонимные корзины и их содержимое
DELETE FROM cart_items WHERE cart_id IN (
  SELECT id FROM shopping_carts WHERE user_id IS NULL
);
DELETE FROM shopping_carts WHERE user_id IS NULL;

-- Результат: удалено 2 корзины, 1 cart_item
```

**Обоснование:** Анонимные корзины больше не создаются на frontend (требуется авторизация). Существующие анонимные корзины - orphan данные из предыдущей версии.

#### Проблемы решены

- ❌ 500 Internal Server Error при попытке удалить несуществующий cart_item
- ❌ Некорректный gRPC error code (Internal вместо NotFound)
- ❌ Накопление orphan анонимных корзин в БД

#### Интеграция с монолитом

**Монолит (vondi) изменения:**
- Frontend теперь требует авторизацию перед добавлением в корзину
- Реализован механизм pending actions для сохранения намерений
- Очистка Redux корзины при логине

**Синхронизация:**
- После логина через OAuth: очистка старых данных корзины из Redux
- Pending action автоматически добавляет товар в новую авторизованную корзину
- Нет конфликта между анонимными и авторизованными корзинами

---

