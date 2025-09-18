-- Создаем функцию для обновления агрегированных просмотров витрины
CREATE OR REPLACE FUNCTION update_storefront_views_count()
RETURNS TRIGGER AS $$
BEGIN
    -- Обновляем views_count для витрины при изменении view_count товара
    IF TG_OP = 'UPDATE' THEN
        -- Обновляем для старой витрины (если товар переместили)
        IF OLD.storefront_id IS DISTINCT FROM NEW.storefront_id AND OLD.storefront_id IS NOT NULL THEN
            UPDATE storefronts
            SET views_count = COALESCE((
                SELECT SUM(view_count)
                FROM storefront_products
                WHERE storefront_id = OLD.storefront_id
            ), 0)
            WHERE id = OLD.storefront_id;
        END IF;

        -- Обновляем для новой/текущей витрины
        IF NEW.storefront_id IS NOT NULL THEN
            UPDATE storefronts
            SET views_count = COALESCE((
                SELECT SUM(view_count)
                FROM storefront_products
                WHERE storefront_id = NEW.storefront_id
            ), 0)
            WHERE id = NEW.storefront_id;
        END IF;
    ELSIF TG_OP = 'INSERT' THEN
        -- При добавлении нового товара
        IF NEW.storefront_id IS NOT NULL THEN
            UPDATE storefronts
            SET views_count = COALESCE((
                SELECT SUM(view_count)
                FROM storefront_products
                WHERE storefront_id = NEW.storefront_id
            ), 0)
            WHERE id = NEW.storefront_id;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        -- При удалении товара
        IF OLD.storefront_id IS NOT NULL THEN
            UPDATE storefronts
            SET views_count = COALESCE((
                SELECT SUM(view_count)
                FROM storefront_products
                WHERE storefront_id = OLD.storefront_id
            ), 0)
            WHERE id = OLD.storefront_id;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер для автоматического обновления views_count
DROP TRIGGER IF EXISTS update_storefront_views_count_trigger ON storefront_products;
CREATE TRIGGER update_storefront_views_count_trigger
AFTER INSERT OR UPDATE OF view_count, storefront_id OR DELETE ON storefront_products
FOR EACH ROW
EXECUTE FUNCTION update_storefront_views_count();

-- Пересчитываем текущие значения
UPDATE storefronts
SET views_count = COALESCE((
    SELECT SUM(view_count)
    FROM storefront_products
    WHERE storefront_id = storefronts.id
), 0);

-- Создаем индекс для оптимизации
CREATE INDEX IF NOT EXISTS idx_storefront_products_storefront_id_view_count
ON storefront_products(storefront_id, view_count);