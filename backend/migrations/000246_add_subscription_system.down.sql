-- Откат миграции системы подписок

-- Удаление триггера
DROP TRIGGER IF EXISTS update_storefront_usage ON storefronts;

-- Удаление функций
DROP FUNCTION IF EXISTS update_subscription_usage();
DROP FUNCTION IF EXISTS get_user_subscription(INTEGER);
DROP FUNCTION IF EXISTS check_subscription_limits(INTEGER, VARCHAR, INTEGER);

-- Удаление полей из таблицы storefronts
ALTER TABLE storefronts
DROP COLUMN IF EXISTS subscription_id,
DROP COLUMN IF EXISTS is_subscription_active;

-- Удаление таблиц в правильном порядке (из-за foreign keys)
DROP TABLE IF EXISTS subscription_usage;
DROP TABLE IF EXISTS subscription_payments;
DROP TABLE IF EXISTS subscription_history;
DROP TABLE IF EXISTS user_subscriptions;
DROP TABLE IF EXISTS subscription_plans;