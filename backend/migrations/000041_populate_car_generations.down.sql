-- Откат миграции: удаляем добавленные поколения
TRUNCATE TABLE car_generations RESTART IDENTITY CASCADE;