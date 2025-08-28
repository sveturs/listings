-- Add new automotive categories under 1003 (Automobili)

-- 1. Domestic production (Serbian brands)
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10100, 'Domaća proizvodnja', 1003, 'domaca-proizvodnja', 'factory', 100),
(10101, 'Zastava vozila', 10100, 'zastava-vozila', 'car', 10),
(10102, 'Yugo klasici', 10100, 'yugo-klasici', 'car', 20),
(10103, 'FAP kamioni', 10100, 'fap-kamioni', 'truck', 30),
(10104, 'IMT traktori', 10100, 'imt-traktori', 'tractor', 40);

-- 2. Import vehicles
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10110, 'Uvozna vozila', 1003, 'uvozna-vozila', 'globe', 110),
(10111, 'EU uvoz', 10110, 'eu-uvoz', 'flag', 10),
(10112, 'Švajcarski uvoz', 10110, 'svajcarski-uvoz', 'flag', 20),
(10113, 'Vozila sa stranim tablicama', 10110, 'vozila-sa-stranim-tablicama', 'id-card', 30);

-- 3. Commercial vehicles
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10120, 'Komercijalna vozila', 1003, 'komercijalna-vozila', 'truck', 120),
(10121, 'Kamioni', 10120, 'kamioni', 'truck', 10),
(10122, 'Autobusi', 10120, 'autobusi', 'bus', 20),
(10123, 'Furgoni', 10120, 'furgoni', 'van', 30),
(10124, 'Prikolice', 10120, 'prikolice', 'trailer', 40),
(10125, 'Specijalna vozila', 10120, 'specijalna-vozila', 'cog', 50);

-- 4. Agricultural machinery
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10130, 'Poljoprivredna tehnika', 1003, 'poljoprivredna-tehnika', 'tractor', 130),
(10131, 'Traktori', 10130, 'traktori', 'tractor', 10),
(10132, 'Kombajni', 10130, 'kombajni', 'wheat', 20),
(10133, 'Priključne mašine', 10130, 'prikljucne-masine', 'tools', 30),
(10134, 'Ostala poljoprivredna tehnika', 10130, 'ostala-poljoprivredna-tehnika', 'gear', 40);

-- 5. Water transport
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10140, 'Vodni transport', 1003, 'vodni-transport', 'ship', 140),
(10141, 'Čamci', 10140, 'camci', 'anchor', 10),
(10142, 'Jahte', 10140, 'jahte', 'sailboat', 20),
(10143, 'Jet ski', 10140, 'jet-ski', 'water', 30),
(10144, 'Motori za čamce', 10140, 'motori-za-camce', 'engine', 40),
(10145, 'Prikolice za čamce', 10140, 'prikolice-za-camce', 'trailer', 50);

-- 6. Alternative transport
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10150, 'Alternativni transport', 1003, 'alternativni-transport', 'bolt', 150),
(10151, 'Električni skuteri', 10150, 'elektricni-skuteri', 'scooter', 10),
(10152, 'Električni bicikli', 10150, 'elektricni-bicikli', 'bicycle', 20),
(10153, 'Kvadovi', 10150, 'kvadovi', 'quad-bike', 30),
(10154, 'Motorne sanke', 10150, 'motorne-sanke', 'snowflake', 40),
(10155, 'Golf vozila', 10150, 'golf-vozila', 'golf', 50);

-- 7. Classic vehicles
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10160, 'Klasična vozila', 1003, 'klasicna-vozila', 'star', 160),
(10161, 'Oldtajmeri', 10160, 'oldtajmeri', 'clock', 10),
(10162, 'Youngtajmeri', 10160, 'youngtajmeri', 'calendar', 20),
(10163, 'Kolekcijski automobili', 10160, 'kolekcijski-automobili', 'gem', 30),
(10164, 'Restaurirani automobili', 10160, 'restaurirani-automobili', 'wrench', 40);

-- 8. Sub-categories for personal cars (under 1301)
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10170, 'Električni automobili', 1301, 'elektricni-automobili', 'battery', 10),
(10171, 'Hibridni automobili', 1301, 'hibridni-automobili', 'leaf', 20),
(10172, 'Luksuzni automobili', 1301, 'luksuzni-automobili', 'crown', 30),
(10173, 'Sportski automobili', 1301, 'sportski-automobili', 'racing', 40),
(10174, 'SUV vozila', 1301, 'suv-vozila', 'mountain', 50),
(10175, 'Karavan vozila', 1301, 'karavan-vozila', 'car-side', 60),
(10176, 'Gradski automobili', 1301, 'gradski-automobili', 'city', 70),
(10177, 'Kamp vozila', 1301, 'kamp-vozila', 'caravan', 80);

-- 9. Sub-categories for motorcycles (under 1302)
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10180, 'Sportski motocikli', 1302, 'sportski-motocikli', 'motorcycle', 10),
(10181, 'Touring motocikli', 1302, 'touring-motocikli', 'road', 20),
(10182, 'Cruiser motocikli', 1302, 'cruiser-motocikli', 'route', 30),
(10183, 'Enduro/Cross motocikli', 1302, 'enduro-cross-motocikli', 'mountain', 40),
(10184, 'Skuteri', 1302, 'skuteri', 'scooter', 50),
(10185, 'Mopedi', 1302, 'mopedi', 'bicycle', 60),
(10186, 'Tricikli', 1302, 'tricikli', 'triangle', 70),
(10187, 'Električni motocikli', 1302, 'elektricni-motocikli', 'battery', 80);

-- 10. Additional auto parts categories (under 1303)
INSERT INTO marketplace_categories (id, name, parent_id, slug, icon, sort_order) VALUES
(10190, 'Akumulatori i punjači', 1303, 'akumulatori-i-punjaci', 'battery', 140),
(10191, 'Audio i video oprema', 1303, 'audio-i-video-oprema', 'music', 150),
(10192, 'GPS i navigacija', 1303, 'gps-i-navigacija', 'map', 160),
(10193, 'Alarmni sistemi', 1303, 'alarmni-sistemi', 'shield', 170),
(10194, 'Tuning delovi', 1303, 'tuning-delovi', 'speed', 180),
(10195, 'Delovi za oldtajmere', 1303, 'delovi-za-oldtajmere', 'vintage', 190);

-- Add translations for new categories
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text) VALUES
-- Domestic production
('category', '10100', 'name', 'sr', 'Domaća proizvodnja'),
('category', '10100', 'name', 'en', 'Domestic Production'),
('category', '10100', 'name', 'ru', 'Отечественное производство'),

('category', '10101', 'name', 'sr', 'Zastava vozila'),
('category', '10101', 'name', 'en', 'Zastava Vehicles'),
('category', '10101', 'name', 'ru', 'Автомобили Застава'),

('category', '10102', 'name', 'sr', 'Yugo klasici'),
('category', '10102', 'name', 'en', 'Yugo Classics'),
('category', '10102', 'name', 'ru', 'Классические Юго'),

('category', '10103', 'name', 'sr', 'FAP kamioni'),
('category', '10103', 'name', 'en', 'FAP Trucks'),
('category', '10103', 'name', 'ru', 'Грузовики ФАП'),

('category', '10104', 'name', 'sr', 'IMT traktori'),
('category', '10104', 'name', 'en', 'IMT Tractors'),
('category', '10104', 'name', 'ru', 'Тракторы ИМТ'),

-- Import vehicles
('category', '10110', 'name', 'sr', 'Uvozna vozila'),
('category', '10110', 'name', 'en', 'Import Vehicles'),
('category', '10110', 'name', 'ru', 'Импортные автомобили'),

('category', '10111', 'name', 'sr', 'EU uvoz'),
('category', '10111', 'name', 'en', 'EU Import'),
('category', '10111', 'name', 'ru', 'Импорт из ЕС'),

('category', '10112', 'name', 'sr', 'Švajcarski uvoz'),
('category', '10112', 'name', 'en', 'Swiss Import'),
('category', '10112', 'name', 'ru', 'Швейцарский импорт'),

('category', '10113', 'name', 'sr', 'Vozila sa stranim tablicama'),
('category', '10113', 'name', 'en', 'Foreign Plates Vehicles'),
('category', '10113', 'name', 'ru', 'Авто на иностранных номерах'),

-- Commercial vehicles
('category', '10120', 'name', 'sr', 'Komercijalna vozila'),
('category', '10120', 'name', 'en', 'Commercial Vehicles'),
('category', '10120', 'name', 'ru', 'Коммерческий транспорт'),

('category', '10121', 'name', 'sr', 'Kamioni'),
('category', '10121', 'name', 'en', 'Trucks'),
('category', '10121', 'name', 'ru', 'Грузовики'),

('category', '10122', 'name', 'sr', 'Autobusi'),
('category', '10122', 'name', 'en', 'Buses'),
('category', '10122', 'name', 'ru', 'Автобусы'),

('category', '10123', 'name', 'sr', 'Furgoni'),
('category', '10123', 'name', 'en', 'Vans'),
('category', '10123', 'name', 'ru', 'Фургоны'),

('category', '10124', 'name', 'sr', 'Prikolice'),
('category', '10124', 'name', 'en', 'Trailers'),
('category', '10124', 'name', 'ru', 'Прицепы'),

('category', '10125', 'name', 'sr', 'Specijalna vozila'),
('category', '10125', 'name', 'en', 'Special Vehicles'),
('category', '10125', 'name', 'ru', 'Спецтехника'),

-- Agricultural machinery
('category', '10130', 'name', 'sr', 'Poljoprivredna tehnika'),
('category', '10130', 'name', 'en', 'Agricultural Machinery'),
('category', '10130', 'name', 'ru', 'Сельхозтехника'),

('category', '10131', 'name', 'sr', 'Traktori'),
('category', '10131', 'name', 'en', 'Tractors'),
('category', '10131', 'name', 'ru', 'Тракторы'),

('category', '10132', 'name', 'sr', 'Kombajni'),
('category', '10132', 'name', 'en', 'Harvesters'),
('category', '10132', 'name', 'ru', 'Комбайны'),

('category', '10133', 'name', 'sr', 'Priključne mašine'),
('category', '10133', 'name', 'en', 'Attachments'),
('category', '10133', 'name', 'ru', 'Навесное оборудование'),

('category', '10134', 'name', 'sr', 'Ostala poljoprivredna tehnika'),
('category', '10134', 'name', 'en', 'Other Agricultural Equipment'),
('category', '10134', 'name', 'ru', 'Прочая сельхозтехника'),

-- Water transport
('category', '10140', 'name', 'sr', 'Vodni transport'),
('category', '10140', 'name', 'en', 'Water Transport'),
('category', '10140', 'name', 'ru', 'Водный транспорт'),

('category', '10141', 'name', 'sr', 'Čamci'),
('category', '10141', 'name', 'en', 'Boats'),
('category', '10141', 'name', 'ru', 'Лодки'),

('category', '10142', 'name', 'sr', 'Jahte'),
('category', '10142', 'name', 'en', 'Yachts'),
('category', '10142', 'name', 'ru', 'Яхты'),

('category', '10143', 'name', 'sr', 'Jet ski'),
('category', '10143', 'name', 'en', 'Jet Skis'),
('category', '10143', 'name', 'ru', 'Гидроциклы'),

('category', '10144', 'name', 'sr', 'Motori za čamce'),
('category', '10144', 'name', 'en', 'Boat Motors'),
('category', '10144', 'name', 'ru', 'Лодочные моторы'),

('category', '10145', 'name', 'sr', 'Prikolice za čamce'),
('category', '10145', 'name', 'en', 'Boat Trailers'),
('category', '10145', 'name', 'ru', 'Прицепы для лодок'),

-- Alternative transport
('category', '10150', 'name', 'sr', 'Alternativni transport'),
('category', '10150', 'name', 'en', 'Alternative Transport'),
('category', '10150', 'name', 'ru', 'Альтернативный транспорт'),

('category', '10151', 'name', 'sr', 'Električni skuteri'),
('category', '10151', 'name', 'en', 'Electric Scooters'),
('category', '10151', 'name', 'ru', 'Электросамокаты'),

('category', '10152', 'name', 'sr', 'Električni bicikli'),
('category', '10152', 'name', 'en', 'Electric Bicycles'),
('category', '10152', 'name', 'ru', 'Электровелосипеды'),

('category', '10153', 'name', 'sr', 'Kvadovi'),
('category', '10153', 'name', 'en', 'ATVs'),
('category', '10153', 'name', 'ru', 'Квадроциклы'),

('category', '10154', 'name', 'sr', 'Motorne sanke'),
('category', '10154', 'name', 'en', 'Snowmobiles'),
('category', '10154', 'name', 'ru', 'Снегоходы'),

('category', '10155', 'name', 'sr', 'Golf vozila'),
('category', '10155', 'name', 'en', 'Golf Carts'),
('category', '10155', 'name', 'ru', 'Гольф-кары'),

-- Classic vehicles
('category', '10160', 'name', 'sr', 'Klasična vozila'),
('category', '10160', 'name', 'en', 'Classic Vehicles'),
('category', '10160', 'name', 'ru', 'Классические автомобили'),

('category', '10161', 'name', 'sr', 'Oldtajmeri'),
('category', '10161', 'name', 'en', 'Oldtimers'),
('category', '10161', 'name', 'ru', 'Олдтаймеры'),

('category', '10162', 'name', 'sr', 'Youngtajmeri'),
('category', '10162', 'name', 'en', 'Youngtimers'),
('category', '10162', 'name', 'ru', 'Янгтаймеры'),

('category', '10163', 'name', 'sr', 'Kolekcijski automobili'),
('category', '10163', 'name', 'en', 'Collectible Cars'),
('category', '10163', 'name', 'ru', 'Коллекционные автомобили'),

('category', '10164', 'name', 'sr', 'Restaurirani automobili'),
('category', '10164', 'name', 'en', 'Restored Cars'),
('category', '10164', 'name', 'ru', 'Реставрированные автомобили'),

-- Personal car sub-categories
('category', '10170', 'name', 'sr', 'Električni automobili'),
('category', '10170', 'name', 'en', 'Electric Cars'),
('category', '10170', 'name', 'ru', 'Электромобили'),

('category', '10171', 'name', 'sr', 'Hibridni automobili'),
('category', '10171', 'name', 'en', 'Hybrid Cars'),
('category', '10171', 'name', 'ru', 'Гибридные автомобили'),

('category', '10172', 'name', 'sr', 'Luksuzni automobili'),
('category', '10172', 'name', 'en', 'Luxury Cars'),
('category', '10172', 'name', 'ru', 'Люксовые автомобили'),

('category', '10173', 'name', 'sr', 'Sportski automobili'),
('category', '10173', 'name', 'en', 'Sports Cars'),
('category', '10173', 'name', 'ru', 'Спортивные автомобили'),

('category', '10174', 'name', 'sr', 'SUV vozila'),
('category', '10174', 'name', 'en', 'SUV Vehicles'),
('category', '10174', 'name', 'ru', 'Внедорожники'),

('category', '10175', 'name', 'sr', 'Karavan vozila'),
('category', '10175', 'name', 'en', 'Station Wagons'),
('category', '10175', 'name', 'ru', 'Универсалы'),

('category', '10176', 'name', 'sr', 'Gradski automobili'),
('category', '10176', 'name', 'en', 'City Cars'),
('category', '10176', 'name', 'ru', 'Городские автомобили'),

('category', '10177', 'name', 'sr', 'Kamp vozila'),
('category', '10177', 'name', 'en', 'Camper Vans'),
('category', '10177', 'name', 'ru', 'Кемперы'),

-- Motorcycle sub-categories
('category', '10180', 'name', 'sr', 'Sportski motocikli'),
('category', '10180', 'name', 'en', 'Sport Motorcycles'),
('category', '10180', 'name', 'ru', 'Спортбайки'),

('category', '10181', 'name', 'sr', 'Touring motocikli'),
('category', '10181', 'name', 'en', 'Touring Motorcycles'),
('category', '10181', 'name', 'ru', 'Туристические мотоциклы'),

('category', '10182', 'name', 'sr', 'Cruiser motocikli'),
('category', '10182', 'name', 'en', 'Cruiser Motorcycles'),
('category', '10182', 'name', 'ru', 'Круизеры'),

('category', '10183', 'name', 'sr', 'Enduro/Cross motocikli'),
('category', '10183', 'name', 'en', 'Enduro/Cross Motorcycles'),
('category', '10183', 'name', 'ru', 'Эндуро/Кросс'),

('category', '10184', 'name', 'sr', 'Skuteri'),
('category', '10184', 'name', 'en', 'Scooters'),
('category', '10184', 'name', 'ru', 'Скутеры'),

('category', '10185', 'name', 'sr', 'Mopedi'),
('category', '10185', 'name', 'en', 'Mopeds'),
('category', '10185', 'name', 'ru', 'Мопеды'),

('category', '10186', 'name', 'sr', 'Tricikli'),
('category', '10186', 'name', 'en', 'Trikes'),
('category', '10186', 'name', 'ru', 'Трициклы'),

('category', '10187', 'name', 'sr', 'Električni motocikli'),
('category', '10187', 'name', 'en', 'Electric Motorcycles'),
('category', '10187', 'name', 'ru', 'Электромотоциклы'),

-- Auto parts sub-categories
('category', '10190', 'name', 'sr', 'Akumulatori i punjači'),
('category', '10190', 'name', 'en', 'Batteries and Chargers'),
('category', '10190', 'name', 'ru', 'Аккумуляторы и зарядные'),

('category', '10191', 'name', 'sr', 'Audio i video oprema'),
('category', '10191', 'name', 'en', 'Audio and Video Equipment'),
('category', '10191', 'name', 'ru', 'Аудио и видео оборудование'),

('category', '10192', 'name', 'sr', 'GPS i navigacija'),
('category', '10192', 'name', 'en', 'GPS and Navigation'),
('category', '10192', 'name', 'ru', 'GPS и навигация'),

('category', '10193', 'name', 'sr', 'Alarmni sistemi'),
('category', '10193', 'name', 'en', 'Alarm Systems'),
('category', '10193', 'name', 'ru', 'Сигнализации'),

('category', '10194', 'name', 'sr', 'Tuning delovi'),
('category', '10194', 'name', 'en', 'Tuning Parts'),
('category', '10194', 'name', 'ru', 'Тюнинг'),

('category', '10195', 'name', 'sr', 'Delovi za oldtajmere'),
('category', '10195', 'name', 'en', 'Oldtimer Parts'),
('category', '10195', 'name', 'ru', 'Запчасти для олдтаймеров');