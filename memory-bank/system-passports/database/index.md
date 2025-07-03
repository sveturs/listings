# Паспорта базы данных PostgreSQL

## Обзор
База данных Sve Tu Platform использует PostgreSQL и содержит 38 таблиц, организованных по функциональным группам.

## Структура паспортов таблиц

Каждый паспорт таблицы содержит:
- Назначение и описание
- Полную структуру CREATE TABLE
- Описание всех полей
- Индексы и их назначение
- Связи с другими таблицами
- Бизнес-правила и ограничения
- Примеры использования
- Известные особенности

## Таблицы по группам

### 👤 Пользователи и аутентификация
- [users](./tables/users.md) - основная таблица пользователей ✅
- [user_contacts](./tables/user_contacts.md) - контактная информация ✅
- [user_privacy_settings](./tables/user_privacy_settings.md) - настройки приватности ✅
- [user_telegram_connections](./tables/user_telegram_connections.md) - Telegram интеграция ✅
- [refresh_tokens](./tables/refresh_tokens.md) - JWT refresh токены ✅

### 🛍️ Маркетплейс
- [marketplace_categories](./tables/marketplace_categories.md) - категории товаров ✅
- [marketplace_listings](./tables/marketplace_listings.md) - объявления ✅
- [marketplace_images](./tables/marketplace_images.md) - изображения товаров ✅
- [marketplace_favorites](./tables/marketplace_favorites.md) - избранные объявления ✅
- [listing_views](./tables/listing_views.md) - просмотры объявлений ✅
- [price_history](./tables/price_history.md) - история изменения цен ✅

### 💬 Коммуникация
- [marketplace_chats](./tables/marketplace_chats.md) - чаты между пользователями ✅
- [marketplace_messages](./tables/marketplace_messages.md) - сообщения в чатах ✅
- [notifications](./tables/notifications.md) - системные уведомления ✅
- [notification_settings](./tables/notification_settings.md) - настройки уведомлений ✅

### 🏷️ Атрибуты и характеристики
- [category_attributes](./category_attributes.md) - определения атрибутов ✅
- [category_attribute_mapping](./category_attribute_mapping.md) - связь атрибутов с категориями ✅
- [listing_attribute_values](./listing_attribute_values.md) - значения атрибутов товаров ✅
- [attribute_groups](./attribute_groups.md) - группировка атрибутов ✅
- [custom_ui_components](./custom_ui_components.md) - кастомные UI компоненты ✅

### 💰 Финансы и платежи
- [user_balances](./tables/user_balances.md) - балансы пользователей ✅
- [balance_transactions](./tables/balance_transactions.md) - транзакции баланса ✅
- [payment_methods](./tables/payment_methods.md) - способы оплаты ✅
- [payment_gateways](./tables/payment_gateways.md) - платежные шлюзы ✅
- [payment_transactions](./tables/payment_transactions.md) - платежные транзакции ✅
- [escrow_payments](./tables/escrow_payments.md) - эскроу платежи ✅
- [merchant_payouts](./tables/merchant_payouts.md) - выплаты продавцам ✅

### ⭐ Отзывы и рейтинги
- [reviews](./reviews.md) - отзывы о сделках ✅
- [review_responses](./review_responses.md) - ответы на отзывы ✅
- [review_votes](./review_votes.md) - голосование за отзывы ✅

### 🏪 Витрины и импорт
- [user_storefronts](./tables/user_storefronts.md) - витрины пользователей (deprecated) ✅
- [storefronts](./tables/storefronts.md) - витрины магазинов ✅
- [storefront_products](./tables/storefront_products.md) - товары витрин ✅
- [storefront_analytics](./tables/storefront_analytics.md) - аналитика витрин ✅
- [import_sources](./tables/import_sources.md) - источники импорта ✅
- [import_history](./tables/import_history.md) - история импорта ✅

### 🌐 Системные таблицы
- [translations](./tables/translations.md) - мультиязычные переводы ✅
- [admin_users](./tables/admin_users.md) - административные пользователи ✅
- [schema_migrations](./tables/schema_migrations.md) - версии миграций БД

## Соглашения о наименованиях

### Таблицы
- Множественное число: `users`, `listings`
- Префиксы по модулям: `marketplace_`, `payment_`
- Snake_case: `user_contacts`, `price_history`

### Поля
- Primary key: всегда `id`
- Foreign keys: `{table}_id` (например, `user_id`, `listing_id`)
- Timestamps: `created_at`, `updated_at`
- Boolean: префикс `is_` или `has_` (например, `is_active`, `has_attachments`)

### Индексы
- Паттерн: `idx_{table}_{field(s)}`
- Составные: `idx_{table}_{field1}_{field2}`
- Частичные: добавляется условие в конце

## Типы данных

### Идентификаторы
- `SERIAL` (INT) - для большинства таблиц
- `BIGSERIAL` (BIGINT) - для таблиц с большим объемом данных

### Строки
- `VARCHAR(n)` - для ограниченных строк
- `TEXT` - для больших текстов без ограничения

### Числа
- `DECIMAL(12,2)` - для денежных сумм
- `INT` - для счетчиков и количества
- `DECIMAL(10,8)` / `DECIMAL(11,8)` - для координат

### Даты
- `TIMESTAMP` - старые таблицы
- `TIMESTAMP WITH TIME ZONE` - новые таблицы (рекомендуется)

### JSON
- `JSONB` - для структурированных данных (settings, metadata)

## Статус документирования

✅ Завершено: 38 таблиц
⏳ В процессе: 0 таблиц  
❌ Не начато: 0 таблиц

**Прогресс: 100%**