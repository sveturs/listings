# Задача: Паспортизация системы Sve Tu Platform

## Дата начала: 29.06.2025
## Статус: В процессе (34% выполнено)

## Цель
Создать полную систему паспортов для всех компонентов проекта, чтобы для создания любых новых подсистем достаточно было смотреть в паспорт, а не листать весь код.

## Выполненная работа

### ✅ Завершено:
1. **Анализ существующей системы паспортов**
   - Изучена структура в `/memory-bank/system-passports/`
   - Выявлено покрытие: только 2 модели из десятков компонентов

2. **Создан план паспортизации**
   - Файл: `/memory-bank/system-passports/PASSPORT_PLAN.md`
   - Определены приоритеты и структура паспортов

3. **Создано 13 паспортов таблиц БД (34%)**:
   
   **Пользователи:**
   - `/database/tables/users.md` - основная таблица пользователей
   
   **Маркетплейс:**
   - `/database/tables/marketplace_categories.md` - категории товаров
   - `/database/tables/marketplace_listings.md` - объявления
   - `/database/tables/marketplace_images.md` - изображения товаров
   - `/database/tables/marketplace_favorites.md` - избранные объявления
   
   **Коммуникации:**
   - `/database/tables/marketplace_chats.md` - чаты между пользователями
   - `/database/tables/marketplace_messages.md` - сообщения в чатах
   - `/database/tables/notifications.md` - системные уведомления
   - `/database/tables/notification_settings.md` - настройки уведомлений
   
   **Финансы:**
   - `/database/tables/user_balances.md` - балансы пользователей
   - `/database/tables/balance_transactions.md` - транзакции баланса
   - `/database/tables/payment_methods.md` - способы оплаты
   - `/database/tables/payment_transactions.md` - платежные транзакции

4. **Обновлены индексы:**
   - `/database/index.md` - главный индекс таблиц БД
   - `/index.md` - основной индекс системы паспортов

## Структура паспорта таблицы

Каждый паспорт содержит:
1. Назначение
2. Полную структуру CREATE TABLE
3. Описание всех полей
4. Индексы и их назначение
5. Триггеры
6. Связи с другими таблицами (прямые и обратные)
7. Бизнес-правила и ограничения
8. Примеры SQL запросов
9. Структуру JSONB полей
10. API интеграцию
11. Известные особенности
12. Информацию о миграциях

## Что осталось сделать

### Таблицы БД (25 таблиц):
**Высокий приоритет:**
- reviews, review_responses, review_votes
- user_contacts, user_privacy_settings
- listing_views, price_history

**Средний приоритет:**
- category_attributes, listing_attribute_values
- payment_gateways, escrow_payments
- storefronts, import_sources

**Низкий приоритет:**
- translations, admin_users
- schema_migrations

### Другие компоненты:
1. **API Endpoints** (12 модулей)
2. **Backend Handlers** (12 модулей)
3. **Frontend компоненты**
4. **OpenSearch индексы**
5. **MinIO структура**
6. **Бизнес-процессы**

## Текущий TODO список

```
1. ✅ Изучить существующую структуру паспортов
2. ⏳ Создать паспорта для всех таблиц БД (34% done)
   2.1 ✅ Таблицы коммуникаций
   2.2 ✅ Финансовые таблицы
   2.3 ⏳ Таблицы атрибутов
   2.4 ⏳ Таблицы отзывов
3. ⏳ Backend handlers
4. ⏳ Frontend компоненты
5. ⏳ MinIO структура
6. ⏳ OpenSearch индексы
7. ⏳ API endpoints
8. ⏳ Обновить навигацию
```

## Следующие шаги

1. **Продолжить с таблицами отзывов:**
   - reviews
   - review_responses
   - review_votes

2. **Затем перейти к API endpoints** - это критично для разработки

## Команда для продолжения в новой сессии

```
Продолжи паспортизацию системы Sve Tu Platform. 
Текущий прогресс: 34% (13 из 38 таблиц БД документировано).
Следующая задача: создать паспорта для таблиц отзывов (reviews, review_responses, review_votes).
Вся информация о задаче в файле: /memory-bank/tasks/001-system-passports.md
```