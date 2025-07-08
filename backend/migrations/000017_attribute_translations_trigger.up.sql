-- Добавляем триггер для автоматической переиндексации при изменении переводов атрибутов
CREATE OR REPLACE FUNCTION trigger_update_listings_on_attribute_translation_change()
RETURNS TRIGGER AS $$
BEGIN
    -- Отмечаем все объявления с этим атрибутом как требующие переиндексации
    UPDATE marketplace_listings
    SET needs_reindex = true
    WHERE id IN (
        SELECT DISTINCT lav.listing_id 
        FROM listing_attribute_values lav
        JOIN category_attributes ca ON lav.attribute_id = ca.id
        WHERE ca.name = NEW.attribute_name
    );
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер для таблицы attribute_option_translations
DROP TRIGGER IF EXISTS update_listings_on_attribute_translation_change ON attribute_option_translations;
DROP TRIGGER IF EXISTS update_listings_on_attribute_translation_change ON update_listings_on_attribute_translation_change;
CREATE TRIGGER update_listings_on_attribute_translation_change
AFTER INSERT OR UPDATE ON attribute_option_translations
FOR EACH ROW
EXECUTE FUNCTION trigger_update_listings_on_attribute_translation_change();

-- Добавляем колонку для отслеживания необходимости переиндексации
ALTER TABLE marketplace_listings ADD COLUMN IF NOT EXISTS needs_reindex BOOLEAN DEFAULT false;