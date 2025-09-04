-- Миграция 000040: Завершение недостающих переводов для категорий
-- Дата: 03.09.2025
-- Цель: Добавить переводы для оставшихся 36 категорий (в основном автозапчасти)

-- Переводы для категорий автозапчастей (уровень 2)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
-- Основные категории автозапчастей
('category', '1313', 'ru', 'name', 'Дополнительное оборудование'),
('category', '1313', 'en', 'name', 'Additional Equipment'),
('category', '1307', 'ru', 'name', 'Электрические и электронные детали'),
('category', '1307', 'en', 'name', 'Electrical and Electronic Parts'),
('category', '1304', 'ru', 'name', 'Шины и колёса'),
('category', '1304', 'en', 'name', 'Tires and Wheels'),
('category', '1306', 'ru', 'name', 'Кузов и детали'),
('category', '1306', 'en', 'name', 'Body and Parts'),
('category', '1305', 'ru', 'name', 'Двигатель и детали двигателя'),
('category', '1305', 'en', 'name', 'Engine and Engine Parts'),
('category', '1310', 'ru', 'name', 'Система охлаждения'),
('category', '1310', 'en', 'name', 'Cooling System'),
('category', '1309', 'ru', 'name', 'Подвеска'),
('category', '1309', 'en', 'name', 'Suspension System'),
('category', '1308', 'ru', 'name', 'Тормозная система'),
('category', '1308', 'en', 'name', 'Braking System'),
('category', '1311', 'ru', 'name', 'Трансмиссия и детали'),
('category', '1311', 'en', 'name', 'Transmission and Parts'),
('category', '1312', 'ru', 'name', 'Интерьер'),
('category', '1312', 'en', 'name', 'Interior'),

-- Технические категории
('category', '2006', 'ru', 'name', 'Фотоаппараты'),
('category', '2006', 'en', 'name', 'Cameras'),
('category', '2006', 'sr', 'name', 'Foto aparati'),
('category', '2007', 'ru', 'name', 'Wi-Fi роутеры'),
('category', '2007', 'en', 'name', 'Wi-Fi Routers'),
('category', '2007', 'sr', 'name', 'Wi-Fi ruteri'),

-- Подкатегории автозапчастей (уровень 3)
('category', '1338', 'ru', 'name', 'Крылья'),
('category', '1338', 'en', 'name', 'Fenders'),
('category', '1316', 'ru', 'name', 'Всесезонные шины'),
('category', '1316', 'en', 'name', 'All-Season Tires'),
('category', '1317', 'ru', 'name', 'Диски'),
('category', '1317', 'en', 'name', 'Rims'),
('category', '1330', 'ru', 'name', 'Фильтры'),
('category', '1330', 'en', 'name', 'Filters'),
('category', '1337', 'ru', 'name', 'Капоты'),
('category', '1337', 'en', 'name', 'Hoods'),
('category', '1334', 'ru', 'name', 'Выхлопная система'),
('category', '1334', 'en', 'name', 'Exhaust System'),
('category', '1318', 'ru', 'name', 'Колёса в сборе'),
('category', '1318', 'en', 'name', 'Complete Wheels'),
('category', '1319', 'ru', 'name', 'Колпаки'),
('category', '1319', 'en', 'name', 'Hubcaps'),

-- Более глубокие подкатегории (уровень 3-4)
('category', '1320', 'ru', 'name', 'Зимние шины'),
('category', '1320', 'en', 'name', 'Winter Tires'),
('category', '1333', 'ru', 'name', 'Свечи зажигания'),
('category', '1333', 'en', 'name', 'Spark Plugs'),
('category', '1342', 'ru', 'name', 'Стёкла'),
('category', '1342', 'en', 'name', 'Glass'),
('category', '1339', 'ru', 'name', 'Двери'),
('category', '1339', 'en', 'name', 'Doors'),
('category', '1314', 'ru', 'name', 'Летние шины'),
('category', '1314', 'en', 'name', 'Summer Tires'),
('category', '1344', 'ru', 'name', 'Задние стёкла'),
('category', '1344', 'en', 'name', 'Rear Windows'),
('category', '1340', 'ru', 'name', 'Багажники'),
('category', '1340', 'en', 'name', 'Trunks'),
('category', '1343', 'ru', 'name', 'Ветровые стёкла'),
('category', '1343', 'en', 'name', 'Windshields'),
('category', '1341', 'ru', 'name', 'Бамперы'),
('category', '1341', 'en', 'name', 'Bumpers'),
('category', '1323', 'ru', 'name', 'Тормозные диски'),
('category', '1323', 'en', 'name', 'Brake Discs'),
('category', '1324', 'ru', 'name', 'Тормозные колодки'),
('category', '1324', 'en', 'name', 'Brake Pads'),
('category', '1327', 'ru', 'name', 'Амортизаторы'),
('category', '1327', 'en', 'name', 'Shock Absorbers'),
('category', '1321', 'ru', 'name', 'Аккумуляторы'),
('category', '1321', 'en', 'name', 'Batteries'),
('category', '1322', 'ru', 'name', 'Фары'),
('category', '1322', 'en', 'name', 'Headlights'),
('category', '1346', 'ru', 'name', 'Кондиционеры'),
('category', '1346', 'en', 'name', 'Air Conditioners'),
('category', '1348', 'ru', 'name', 'Радиаторы'),
('category', '1348', 'en', 'name', 'Radiators'),
('category', '1347', 'ru', 'name', 'Обогреватели'),
('category', '1347', 'en', 'name', 'Heaters')
ON CONFLICT (entity_type, entity_id, language, field_name) DO UPDATE 
SET translated_text = EXCLUDED.translated_text,
    updated_at = CURRENT_TIMESTAMP;

-- Обновляем счетчик переводов
DO $$
DECLARE
    total_categories INTEGER;
    translated_categories INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_categories FROM marketplace_categories;
    SELECT COUNT(DISTINCT entity_id::integer) INTO translated_categories 
    FROM translations 
    WHERE entity_type = 'category' 
      AND field_name = 'name' 
      AND language = 'ru';
    
    RAISE NOTICE 'Переводы добавлены. Всего категорий: %, Переведено на русский: % (%)', 
                 total_categories, translated_categories, 
                 ROUND((translated_categories::numeric / total_categories) * 100, 1);
END $$;