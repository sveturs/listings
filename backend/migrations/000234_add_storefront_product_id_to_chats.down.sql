-- Удаляем внешний ключ
ALTER TABLE marketplace_chats
DROP CONSTRAINT IF EXISTS fk_marketplace_chats_storefront_product;

-- Удаляем constraint
ALTER TABLE marketplace_chats
DROP CONSTRAINT IF EXISTS check_chat_target;

-- Удаляем индекс
DROP INDEX IF EXISTS idx_marketplace_chats_storefront_product_id;

-- Удаляем столбец
ALTER TABLE marketplace_chats
DROP COLUMN IF EXISTS storefront_product_id;