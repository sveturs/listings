-- Удаление триггеров
DROP TRIGGER IF EXISTS trigger_shopping_cart_items_updated_at ON shopping_cart_items;
DROP TRIGGER IF EXISTS trigger_shopping_carts_updated_at ON shopping_carts;

-- Удаление функции
DROP FUNCTION IF EXISTS update_shopping_cart_updated_at();

-- Удаление индексов
DROP INDEX IF EXISTS idx_shopping_cart_items_product_id;
DROP INDEX IF EXISTS idx_shopping_cart_items_cart_id;
DROP INDEX IF EXISTS idx_shopping_carts_storefront_id;
DROP INDEX IF EXISTS idx_shopping_carts_session_id;
DROP INDEX IF EXISTS idx_shopping_carts_user_id;

-- Удаление таблиц (в обратном порядке из-за зависимостей)
DROP TABLE IF EXISTS shopping_cart_items;
DROP TABLE IF EXISTS shopping_carts;