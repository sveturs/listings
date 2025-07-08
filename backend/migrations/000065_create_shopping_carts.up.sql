-- Создание таблицы корзин
CREATE TABLE shopping_carts (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    session_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Ограничения: корзина должна принадлежать либо пользователю, либо сессии
    CONSTRAINT check_cart_owner CHECK (
        (user_id IS NOT NULL AND session_id IS NULL) OR
        (user_id IS NULL AND session_id IS NOT NULL)
    ),
    
    -- Уникальность: один пользователь может иметь только одну корзину на витрину
    CONSTRAINT unique_user_storefront_cart UNIQUE (user_id, storefront_id),
    
    -- Уникальность: одна сессия может иметь только одну корзину на витрину  
    CONSTRAINT unique_session_storefront_cart UNIQUE (session_id, storefront_id)
);

-- Создание таблицы позиций корзины
CREATE TABLE shopping_cart_items (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES shopping_carts(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES storefront_product_variants(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price_per_unit DECIMAL(10,2) NOT NULL CHECK (price_per_unit >= 0),
    total_price DECIMAL(10,2) NOT NULL CHECK (total_price >= 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Уникальность: один товар (с вариантом) может быть только один раз в корзине
    CONSTRAINT unique_cart_product_variant UNIQUE (cart_id, product_id, variant_id)
);

-- Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_shopping_carts_user_id ON shopping_carts(user_id);
CREATE INDEX IF NOT EXISTS idx_shopping_carts_session_id ON shopping_carts(session_id);
CREATE INDEX IF NOT EXISTS idx_shopping_carts_storefront_id ON shopping_carts(storefront_id);
CREATE INDEX IF NOT EXISTS idx_shopping_cart_items_cart_id ON shopping_cart_items(cart_id);
CREATE INDEX IF NOT EXISTS idx_shopping_cart_items_product_id ON shopping_cart_items(product_id);

-- Функция обновления updated_at
CREATE OR REPLACE FUNCTION update_shopping_cart_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггеры для автоматического обновления updated_at
CREATE TRIGGER trigger_shopping_carts_updated_at
    BEFORE UPDATE ON shopping_carts
    FOR EACH ROW
    EXECUTE FUNCTION update_shopping_cart_updated_at();

CREATE TRIGGER trigger_shopping_cart_items_updated_at
    BEFORE UPDATE ON shopping_cart_items
    FOR EACH ROW
    EXECUTE FUNCTION update_shopping_cart_updated_at();