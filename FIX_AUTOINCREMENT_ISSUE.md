# Исправление проблемы с автоинкрементом в базе данных

## Дата: 22 июля 2025

## Проблема
В базе данных множество таблиц потеряли DEFAULT значения для автоинкремента полей ID. Это привело к ошибкам при создании новых записей:
- `ERROR: null value in column "id" of relation "XXX" violates not-null constraint`

## Затронутые функции
1. **Google OAuth авторизация** - не создавались новые пользователи
2. **Создание витрин (storefronts)** - не создавались витрины и связанные записи
3. **Уведомления** - не создавались уведомления

## Исправленные таблицы

### Миграция 000129
- `users` - исправлен автоинкремент для создания новых пользователей

### Миграция 000130
- `storefronts` - исправлен автоинкремент для создания витрин

### Миграция 000131
- `storefront_staff` - исправлен автоинкремент для добавления персонала

### Миграция 000132 (комплексная)
- `storefront_delivery_options`
- `storefront_hours`
- `storefront_inventory_movements`
- `storefront_orders`
- `storefront_payment_methods`
- `storefront_product_attributes`
- `storefront_product_variants`

### Миграция 000133
- `notifications` - исправлен автоинкремент для уведомлений

## Таблицы, требующие исправления
При проверке обнаружены другие таблицы без автоинкремента:
- `address_change_log`
- `admin_users`
- `attribute_group_items`
- `attribute_groups`
- `attribute_option_translations`
- `balance_transactions`
- `category_attribute_groups`
- `category_attributes`
- `chat_attachments`
- `component_templates`
- `custom_ui_component_usage`
- `custom_ui_components`
- `custom_ui_templates`
- `escrow_payments`
- `geocoding_cache`
- `gis_filter_analytics`
- `gis_isochrone_cache`
- `gis_poi_cache`
- `import_history`
- `import_sources`
- И другие...

## Команда для проверки таблиц без автоинкремента
```sql
SELECT 
    c.table_name,
    c.column_name,
    CASE 
        WHEN c.column_default IS NULL THEN 'MISSING'
        ELSE 'OK'
    END as status
FROM information_schema.columns c
JOIN pg_class pc ON pc.relname = c.table_name
JOIN pg_namespace pn ON pn.oid = pc.relnamespace AND pn.nspname = c.table_schema
WHERE c.table_schema = 'public' 
AND c.column_name = 'id'
AND c.data_type IN ('integer', 'bigint')
AND c.column_default IS NULL
AND pc.relkind = 'r'
ORDER BY c.table_name;
```

## Шаблон миграции для исправления
```sql
-- up.sql
ALTER TABLE table_name 
ALTER COLUMN id SET DEFAULT nextval('table_name_id_seq');
ALTER SEQUENCE table_name_id_seq OWNED BY table_name.id;
SELECT setval('table_name_id_seq', COALESCE((SELECT MAX(id) FROM table_name), 1), true);

-- down.sql
ALTER TABLE table_name ALTER COLUMN id DROP DEFAULT;
```

## Текущий статус
- ✅ Google OAuth работает
- ✅ Создание витрин работает
- ✅ Уведомления должны работать
- ⚠️ Другие таблицы могут требовать исправления при возникновении ошибок

## Рекомендации
1. При возникновении похожих ошибок проверить наличие DEFAULT для id в соответствующей таблице
2. Создать миграцию по шаблону выше
3. Рассмотреть создание общей миграции для исправления всех оставшихся таблиц