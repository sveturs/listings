-- Откат миграции: убираем DEFAULT значение
ALTER TABLE storefronts 
ALTER COLUMN id DROP DEFAULT;