-- Распределяем тестовые объявления по разным координатам в районе Белграда
-- для лучшей демонстрации работы карты и кластеризации

-- Обновляем координаты для объявлений одежды и аксессуаров (ID 222-228)
-- Распределяем их по кругу в радиусе 2-4 км от центра

-- Tommy Hilfiger Polo majica - 2км на север
UPDATE marketplace_listings
SET latitude = 44.8356, longitude = 20.4649
WHERE id = 222;

-- Calvin Klein džins - 2.5км на северо-восток
UPDATE marketplace_listings
SET latitude = 44.8326, longitude = 20.4899
WHERE id = 223;

-- Zara haljina - 3км на восток
UPDATE marketplace_listings
SET latitude = 44.8176, longitude = 20.4989
WHERE id = 224;

-- Adidas Ultraboost 22 - 2км на юго-восток
UPDATE marketplace_listings
SET latitude = 44.8006, longitude = 20.4839
WHERE id = 225;

-- Nike Air Force 1 - 3км на юг
UPDATE marketplace_listings
SET latitude = 44.7906, longitude = 20.4649
WHERE id = 226;

-- Guess sat za dame - 2.5км на юго-запад
UPDATE marketplace_listings
SET latitude = 44.7976, longitude = 20.4399
WHERE id = 227;

-- Ray-Ban Aviator naočare - 2км на запад
UPDATE marketplace_listings
SET latitude = 44.8176, longitude = 20.4309
WHERE id = 228;

-- Также обновим координаты для некоторых других объявлений, чтобы они были в радиусе поиска

-- BMW X5 - 1.5км на северо-запад
UPDATE marketplace_listings
SET latitude = 44.8296, longitude = 20.4479
WHERE id = 253;

-- Chicco kolica - остается близко к центру
UPDATE marketplace_listings
SET latitude = 44.8196, longitude = 20.4679
WHERE id = 266;

-- Trek električni bicikl - 1км на северо-восток
UPDATE marketplace_listings
SET latitude = 44.8256, longitude = 20.4749
WHERE id = 261;