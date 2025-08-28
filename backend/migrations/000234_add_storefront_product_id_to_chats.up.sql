-- Добавляем поле storefront_product_id в таблицу marketplace_chats
ALTER TABLE marketplace_chats
ADD COLUMN storefront_product_id INTEGER DEFAULT NULL;

-- Создаем индекс для быстрого поиска
CREATE INDEX idx_marketplace_chats_storefront_product_id ON marketplace_chats(storefront_product_id);

-- Добавляем constraint - чат должен быть связан либо с listing_id, либо с storefront_product_id
-- Но разрешаем оба поля быть NULL для обратной совместимости (прямые сообщения пользователей)
ALTER TABLE marketplace_chats
ADD CONSTRAINT check_chat_target
CHECK (
    NOT (listing_id IS NOT NULL AND storefront_product_id IS NOT NULL)
);

-- Добавляем внешний ключ на storefront_products
ALTER TABLE marketplace_chats
ADD CONSTRAINT fk_marketplace_chats_storefront_product
FOREIGN KEY (storefront_product_id) 
REFERENCES storefront_products(id) 
ON DELETE CASCADE;