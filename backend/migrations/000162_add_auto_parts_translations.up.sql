-- Add translations for auto parts categories

-- First, fix incorrect translations for category 1304
DELETE FROM translations WHERE entity_type = 'category' AND entity_id = 1304;

-- Add translations for auto parts main category and subcategories
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text) VALUES
-- Auto delovi (1303) - already exists, but let's update English and Russian
('category', 1303, 'name', 'en', 'Auto Parts'),
('category', 1303, 'name', 'ru', 'Автозапчасти'),
('category', 1303, 'seo_title', 'en', 'Auto Parts and Accessories'),
('category', 1303, 'seo_title', 'ru', 'Автозапчасти и аксессуары'),
('category', 1303, 'seo_description', 'en', 'Spare parts and automotive accessories'),
('category', 1303, 'seo_description', 'ru', 'Запасные части и автомобильные аксессуары'),

-- Gume i točkovi (1304)
('category', 1304, 'name', 'en', 'Tires and Wheels'),
('category', 1304, 'name', 'ru', 'Шины и колеса'),
('category', 1304, 'name', 'sr', 'Гуме и точкови'),
('category', 1304, 'seo_title', 'en', 'Tires, Wheels and Rims'),
('category', 1304, 'seo_title', 'ru', 'Шины, колеса и диски'),
('category', 1304, 'seo_title', 'sr', 'Гуме, точкови и фелне'),
('category', 1304, 'seo_description', 'en', 'Summer, winter and all-season tires, wheels and rims'),
('category', 1304, 'seo_description', 'ru', 'Летние, зимние и всесезонные шины, колеса и диски'),
('category', 1304, 'seo_description', 'sr', 'Летње, зимске и целогодишње гуме, точкови и фелне'),

-- Motor i delovi motora (1305)
('category', 1305, 'name', 'en', 'Engine and Engine Parts'),
('category', 1305, 'name', 'ru', 'Двигатель и детали двигателя'),
('category', 1305, 'name', 'sr', 'Мотор и делови мотора'),
('category', 1305, 'seo_title', 'en', 'Engine Parts and Components'),
('category', 1305, 'seo_title', 'ru', 'Детали и компоненты двигателя'),
('category', 1305, 'seo_title', 'sr', 'Делови и компоненте мотора'),

-- Karoserija i delovi (1306)
('category', 1306, 'name', 'en', 'Body Parts'),
('category', 1306, 'name', 'ru', 'Кузовные детали'),
('category', 1306, 'name', 'sr', 'Каросерија и делови'),
('category', 1306, 'seo_title', 'en', 'Car Body Parts and Accessories'),
('category', 1306, 'seo_title', 'ru', 'Детали кузова и аксессуары'),
('category', 1306, 'seo_title', 'sr', 'Делови каросерије и прибор'),

-- Električni i elektronski delovi (1307)
('category', 1307, 'name', 'en', 'Electrical and Electronic Parts'),
('category', 1307, 'name', 'ru', 'Электрические и электронные детали'),
('category', 1307, 'name', 'sr', 'Електрични и електронски делови'),

-- Sistem za kočenje (1308)
('category', 1308, 'name', 'en', 'Brake System'),
('category', 1308, 'name', 'ru', 'Тормозная система'),
('category', 1308, 'name', 'sr', 'Систем за кочење'),

-- Sistem vešanja (1309)
('category', 1309, 'name', 'en', 'Suspension System'),
('category', 1309, 'name', 'ru', 'Подвеска'),
('category', 1309, 'name', 'sr', 'Систем вешања'),

-- Sistem hlađenja (1310)
('category', 1310, 'name', 'en', 'Cooling System'),
('category', 1310, 'name', 'ru', 'Система охлаждения'),
('category', 1310, 'name', 'sr', 'Систем хлађења'),

-- Transmisija i delovi (1311)
('category', 1311, 'name', 'en', 'Transmission and Parts'),
('category', 1311, 'name', 'ru', 'Трансмиссия и детали'),
('category', 1311, 'name', 'sr', 'Трансмисија и делови'),

-- Unutrašnjost (1312)
('category', 1312, 'name', 'en', 'Interior Parts'),
('category', 1312, 'name', 'ru', 'Интерьер'),
('category', 1312, 'name', 'sr', 'Унутрашњост'),

-- Dodatna oprema (1313)
('category', 1313, 'name', 'en', 'Auto Accessories'),
('category', 1313, 'name', 'ru', 'Автомобильные аксессуары'),
('category', 1313, 'name', 'sr', 'Додатна опрема'),

-- Tire subcategories
-- Letnje gume (1314)
('category', 1314, 'name', 'en', 'Summer Tires'),
('category', 1314, 'name', 'ru', 'Летние шины'),
('category', 1314, 'name', 'sr', 'Летње гуме'),

-- Zimske gume (1315)
('category', 1315, 'name', 'en', 'Winter Tires'),
('category', 1315, 'name', 'ru', 'Зимние шины'),
('category', 1315, 'name', 'sr', 'Зимске гуме'),

-- Celogodišnje gume (1316)
('category', 1316, 'name', 'en', 'All-Season Tires'),
('category', 1316, 'name', 'ru', 'Всесезонные шины'),
('category', 1316, 'name', 'sr', 'Целогодишње гуме'),

-- Felne (1317)
('category', 1317, 'name', 'en', 'Rims'),
('category', 1317, 'name', 'ru', 'Диски'),
('category', 1317, 'name', 'sr', 'Фелне'),

-- Kompletni točkovi (1318)
('category', 1318, 'name', 'en', 'Complete Wheels'),
('category', 1318, 'name', 'ru', 'Колеса в сборе'),
('category', 1318, 'name', 'sr', 'Комплетни точкови'),

-- Ratkapne (1319)
('category', 1319, 'name', 'en', 'Wheel Covers'),
('category', 1319, 'name', 'ru', 'Колпаки'),
('category', 1319, 'name', 'sr', 'Раткапне'),

-- Vijci za točkove (1320)
('category', 1320, 'name', 'en', 'Wheel Bolts'),
('category', 1320, 'name', 'ru', 'Болты для колес'),
('category', 1320, 'name', 'sr', 'Вијци за точкове'),

-- Third level tire categories
-- Putničke letnje gume (1321)
('category', 1321, 'name', 'en', 'Passenger Summer Tires'),
('category', 1321, 'name', 'ru', 'Летние шины для легковых авто'),
('category', 1321, 'name', 'sr', 'Путничке летње гуме'),

-- SUV letnje gume (1322)
('category', 1322, 'name', 'en', 'SUV Summer Tires'),
('category', 1322, 'name', 'ru', 'Летние шины для внедорожников'),
('category', 1322, 'name', 'sr', 'SUV летње гуме'),

-- Kamionske letnje gume (1323)
('category', 1323, 'name', 'en', 'Truck Summer Tires'),
('category', 1323, 'name', 'ru', 'Летние грузовые шины'),
('category', 1323, 'name', 'sr', 'Камионске летње гуме'),

-- Putničke zimske gume (1324)
('category', 1324, 'name', 'en', 'Passenger Winter Tires'),
('category', 1324, 'name', 'ru', 'Зимние шины для легковых авто'),
('category', 1324, 'name', 'sr', 'Путничке зимске гуме'),

-- SUV zimske gume (1325)
('category', 1325, 'name', 'en', 'SUV Winter Tires'),
('category', 1325, 'name', 'ru', 'Зимние шины для внедорожников'),
('category', 1325, 'name', 'sr', 'SUV зимске гуме'),

-- Kamionske zimske gume (1326)
('category', 1326, 'name', 'en', 'Truck Winter Tires'),
('category', 1326, 'name', 'ru', 'Зимние грузовые шины'),
('category', 1326, 'name', 'sr', 'Камионске зимске гуме'),

-- Čelične felne (1327)
('category', 1327, 'name', 'en', 'Steel Rims'),
('category', 1327, 'name', 'ru', 'Стальные диски'),
('category', 1327, 'name', 'sr', 'Челичне фелне'),

-- Aluminijumske felne (1328)
('category', 1328, 'name', 'en', 'Aluminum Rims'),
('category', 1328, 'name', 'ru', 'Алюминиевые диски'),
('category', 1328, 'name', 'sr', 'Алуминијумске фелне'),

-- Sportske felne (1329)
('category', 1329, 'name', 'en', 'Sport Rims'),
('category', 1329, 'name', 'ru', 'Спортивные диски'),
('category', 1329, 'name', 'sr', 'Спортске фелне'),

-- Engine subcategories
-- Filtri (1330)
('category', 1330, 'name', 'en', 'Filters'),
('category', 1330, 'name', 'ru', 'Фильтры'),
('category', 1330, 'name', 'sr', 'Филтри'),

-- Remeni i lančanici (1331)
('category', 1331, 'name', 'en', 'Belts and Chains'),
('category', 1331, 'name', 'ru', 'Ремни и цепи'),
('category', 1331, 'name', 'sr', 'Ремени и ланчаници'),

-- Ulje i tečnosti (1332)
('category', 1332, 'name', 'en', 'Oils and Fluids'),
('category', 1332, 'name', 'ru', 'Масла и жидкости'),
('category', 1332, 'name', 'sr', 'Уље и течности'),

-- Svećice (1333)
('category', 1333, 'name', 'en', 'Spark Plugs'),
('category', 1333, 'name', 'ru', 'Свечи зажигания'),
('category', 1333, 'name', 'sr', 'Свећице'),

-- Izduvni sistem (1334)
('category', 1334, 'name', 'en', 'Exhaust System'),
('category', 1334, 'name', 'ru', 'Выхлопная система'),
('category', 1334, 'name', 'sr', 'Издувни систем'),

-- Body parts subcategories
-- Branici (1335)
('category', 1335, 'name', 'en', 'Bumpers'),
('category', 1335, 'name', 'ru', 'Бамперы'),
('category', 1335, 'name', 'sr', 'Браници'),

-- Vrata (1336)
('category', 1336, 'name', 'en', 'Doors'),
('category', 1336, 'name', 'ru', 'Двери'),
('category', 1336, 'name', 'sr', 'Врата'),

-- Haube (1337)
('category', 1337, 'name', 'en', 'Hoods'),
('category', 1337, 'name', 'ru', 'Капоты'),
('category', 1337, 'name', 'sr', 'Хаубе'),

-- Blatobrani (1338)
('category', 1338, 'name', 'en', 'Fenders'),
('category', 1338, 'name', 'ru', 'Крылья'),
('category', 1338, 'name', 'sr', 'Блатобрани'),

-- Retrovizori (1339)
('category', 1339, 'name', 'en', 'Mirrors'),
('category', 1339, 'name', 'ru', 'Зеркала'),
('category', 1339, 'name', 'sr', 'Ретровизори'),

-- Stakla (1340)
('category', 1340, 'name', 'en', 'Windows'),
('category', 1340, 'name', 'ru', 'Стекла'),
('category', 1340, 'name', 'sr', 'Стакла')
ON CONFLICT (entity_type, entity_id, field_name, language) 
DO UPDATE SET translated_text = EXCLUDED.translated_text;

-- Update existing translations for category 1303 where needed
UPDATE translations 
SET translated_text = 'Auto Parts' 
WHERE entity_type = 'category' AND entity_id = 1303 AND field_name = 'name' AND language = 'en';

UPDATE translations 
SET translated_text = 'Автозапчасти' 
WHERE entity_type = 'category' AND entity_id = 1303 AND field_name = 'name' AND language = 'ru';