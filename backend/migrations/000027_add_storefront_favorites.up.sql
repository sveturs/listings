-- Создание таблицы для избранных товаров витрин
CREATE TABLE IF NOT EXISTS storefront_favorites (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, product_id)
);

-- Создание индексов для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_storefront_favorites_user_id ON storefront_favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_storefront_favorites_product_id ON storefront_favorites(product_id);
CREATE INDEX IF NOT EXISTS idx_storefront_favorites_created_at ON storefront_favorites(created_at DESC);

-- Добавляем комментарии
COMMENT ON TABLE storefront_favorites IS 'Избранные товары витрин пользователей';
COMMENT ON COLUMN storefront_favorites.user_id IS 'ID пользователя';
COMMENT ON COLUMN storefront_favorites.product_id IS 'ID товара витрины';
COMMENT ON COLUMN storefront_favorites.created_at IS 'Дата добавления в избранное';