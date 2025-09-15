-- Откат миграции - возвращаем оригинальные координаты

-- Возвращаем одинаковые координаты для объявлений одежды (ID 222-228)
UPDATE marketplace_listings
SET latitude = 44.8125, longitude = 20.4612
WHERE id IN (222, 223, 224, 225, 226, 227, 228);

-- Возвращаем оригинальные координаты для других объявлений
UPDATE marketplace_listings
SET latitude = 44.8125, longitude = 20.4612
WHERE id = 253;

UPDATE marketplace_listings
SET latitude = 44.8176, longitude = 20.4633
WHERE id = 266;

UPDATE marketplace_listings
SET latitude = 44.8176, longitude = 20.4633
WHERE id = 261;