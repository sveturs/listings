-- Добавляем поле storefront_product_id для поддержки чатов по товарам витрин
ALTER TABLE marketplace_messages
ADD COLUMN storefront_product_id INTEGER DEFAULT NULL;

-- Добавляем индекс для быстрого поиска сообщений по товару витрины
CREATE INDEX idx_marketplace_messages_storefront_product_id 
ON marketplace_messages(storefront_product_id) 
WHERE storefront_product_id IS NOT NULL;

-- Добавляем внешний ключ на таблицу storefront_products
ALTER TABLE marketplace_messages
ADD CONSTRAINT fk_marketplace_messages_storefront_product
FOREIGN KEY (storefront_product_id) 
REFERENCES storefront_products(id) 
ON DELETE CASCADE;

-- Добавляем комментарий к полю
COMMENT ON COLUMN marketplace_messages.storefront_product_id IS 'ID товара витрины, если сообщение относится к товару витрины, а не к обычному объявлению';

-- Добавляем проверку: сообщение должно относиться либо к listing_id, либо к storefront_product_id, либо быть прямым сообщением
ALTER TABLE marketplace_messages
ADD CONSTRAINT check_message_target
CHECK (
    (listing_id IS NOT NULL AND storefront_product_id IS NULL) OR
    (listing_id IS NULL AND storefront_product_id IS NOT NULL) OR
    (listing_id IS NULL AND storefront_product_id IS NULL) -- прямое сообщение между контактами
);