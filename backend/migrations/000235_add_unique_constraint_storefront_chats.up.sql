-- Добавляем уникальный constraint для чатов с товарами витрин
-- Это позволит использовать ON CONFLICT при создании чатов
ALTER TABLE marketplace_chats
ADD CONSTRAINT unique_storefront_product_chat
UNIQUE (storefront_product_id, buyer_id, seller_id);

-- Комментарий: этот constraint гарантирует, что для каждого товара витрины
-- может быть только один чат между конкретным покупателем и продавцом