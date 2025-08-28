-- Откат миграции для многоязычных адресов

-- Удаляем индексы
DROP INDEX IF EXISTS idx_translations_entity_field;
DROP INDEX IF EXISTS idx_translations_listing_all;

-- Переводы удаляются каскадно при удалении объявлений