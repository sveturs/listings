-- Добавляем поддержку многоязычных адресов через систему переводов
-- Для каждого объявления теперь можно хранить переводы полей location, address_city, address_country

-- Добавляем индекс для быстрого поиска переводов по entity_id и field_name
CREATE INDEX IF NOT EXISTS idx_translations_entity_field 
ON translations(entity_type, entity_id, field_name, language)
WHERE entity_type = 'listing';

-- Добавляем составной индекс для оптимизации выборки всех переводов объявления
CREATE INDEX IF NOT EXISTS idx_translations_listing_all 
ON translations(entity_id, language)
WHERE entity_type = 'listing';

-- Данные будут заполняться через API при создании/обновлении объявлений

-- Комментарий для разработчиков
COMMENT ON INDEX idx_translations_entity_field IS 'Индекс для быстрого поиска переводов адресных полей объявлений';
COMMENT ON INDEX idx_translations_listing_all IS 'Индекс для выборки всех переводов объявления одним запросом';