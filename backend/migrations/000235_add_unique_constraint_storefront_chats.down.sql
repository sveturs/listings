-- Удаляем уникальный constraint для чатов с товарами витрин
ALTER TABLE marketplace_chats
DROP CONSTRAINT IF EXISTS unique_storefront_product_chat;