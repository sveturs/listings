-- Откат миграции: убираем DEFAULT значение
ALTER TABLE storefront_staff 
ALTER COLUMN id DROP DEFAULT;