-- Откат миграции: убираем DEFAULT значение
ALTER TABLE users 
ALTER COLUMN id DROP DEFAULT;