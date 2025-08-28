-- Удаляем ограничение проверки
ALTER TABLE marketplace_messages
DROP CONSTRAINT IF EXISTS check_message_target;

-- Удаляем внешний ключ
ALTER TABLE marketplace_messages
DROP CONSTRAINT IF EXISTS fk_marketplace_messages_storefront_product;

-- Удаляем индекс
DROP INDEX IF EXISTS idx_marketplace_messages_storefront_product_id;

-- Удаляем колонку
ALTER TABLE marketplace_messages
DROP COLUMN IF EXISTS storefront_product_id;