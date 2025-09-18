-- Создаём функцию для обновления счётчика товаров в витрине
CREATE OR REPLACE FUNCTION update_storefront_products_count()
RETURNS TRIGGER AS $$
BEGIN
    -- При вставке или удалении товара обновляем счётчик
    IF TG_OP = 'INSERT' THEN
        UPDATE storefronts
        SET products_count = (
            SELECT COUNT(*)
            FROM storefront_products
            WHERE storefront_id = NEW.storefront_id
            AND is_active = true
        )
        WHERE id = NEW.storefront_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE storefronts
        SET products_count = (
            SELECT COUNT(*)
            FROM storefront_products
            WHERE storefront_id = OLD.storefront_id
            AND is_active = true
        )
        WHERE id = OLD.storefront_id;
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        -- При обновлении проверяем изменение is_active или storefront_id
        IF OLD.is_active != NEW.is_active OR OLD.storefront_id != NEW.storefront_id THEN
            -- Обновляем счётчик для старой витрины
            IF OLD.storefront_id != NEW.storefront_id THEN
                UPDATE storefronts
                SET products_count = (
                    SELECT COUNT(*)
                    FROM storefront_products
                    WHERE storefront_id = OLD.storefront_id
                    AND is_active = true
                )
                WHERE id = OLD.storefront_id;
            END IF;

            -- Обновляем счётчик для новой витрины
            UPDATE storefronts
            SET products_count = (
                SELECT COUNT(*)
                FROM storefront_products
                WHERE storefront_id = NEW.storefront_id
                AND is_active = true
            )
            WHERE id = NEW.storefront_id;
        END IF;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Создаём триггер для автообновления счётчика товаров
DROP TRIGGER IF EXISTS trigger_update_storefront_products_count ON storefront_products;
CREATE TRIGGER trigger_update_storefront_products_count
AFTER INSERT OR UPDATE OR DELETE ON storefront_products
FOR EACH ROW EXECUTE FUNCTION update_storefront_products_count();

-- Обновляем текущие счётчики для всех витрин
UPDATE storefronts s
SET products_count = (
    SELECT COUNT(*)
    FROM storefront_products sp
    WHERE sp.storefront_id = s.id
    AND sp.is_active = true
);

-- Добавляем комментарий к колонке для документации
COMMENT ON COLUMN storefronts.products_count IS 'Автоматически обновляемый счётчик активных товаров в витрине';