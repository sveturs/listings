-- Добавление переводов для объявлений без переводов
-- Миграция добавляет переводы на английский и русский языки для сербских объявлений

-- Функция для безопасного добавления переводов (избегаем дубликатов)
CREATE OR REPLACE FUNCTION add_translation_if_not_exists(
    p_entity_type VARCHAR,
    p_entity_id INTEGER,
    p_language VARCHAR,
    p_field_name VARCHAR,
    p_translated_text TEXT
) RETURNS VOID AS $$
BEGIN
    INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
    VALUES (p_entity_type, p_entity_id, p_language, p_field_name, p_translated_text, true, false)
    ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;
END;
$$ LANGUAGE plpgsql;

-- Объявление 183: Домашний мёд
SELECT add_translation_if_not_exists('marketplace_listing', 183, 'en', 'title', 'Domestic Acacia Honey 1kg');
SELECT add_translation_if_not_exists('marketplace_listing', 183, 'en', 'description', 'Pure acacia honey from our own apiary. 100% natural, no additives. Crystal clear, light color.');
SELECT add_translation_if_not_exists('marketplace_listing', 183, 'ru', 'title', 'Домашний мёд акация 1кг');
SELECT add_translation_if_not_exists('marketplace_listing', 183, 'ru', 'description', 'Чистый акациевый мёд с собственной пасеки. 100% натуральный, без добавок. Кристально чистый, светлого цвета.');

-- Объявление 250: Квартира Лиман 3
SELECT add_translation_if_not_exists('marketplace_listing', 250, 'en', 'title', 'Apartment 65m2 Liman 3 - New Building');
SELECT add_translation_if_not_exists('marketplace_listing', 250, 'en', 'description', 'Beautiful apartment in a new building on Liman 3. Fully furnished, ready to move in. Terrace, parking space.');
SELECT add_translation_if_not_exists('marketplace_listing', 250, 'ru', 'title', 'Квартира 65м2 Лиман 3 - новостройка');
SELECT add_translation_if_not_exists('marketplace_listing', 250, 'ru', 'description', 'Прекрасная квартира в новостройке на Лимане 3. Полностью меблирована, готова к заселению. Терраса, парковочное место.');

-- Объявление 251: Люкс пентхаус
SELECT add_translation_if_not_exists('marketplace_listing', 251, 'en', 'title', 'Luxury Penthouse 120m2 Center');
SELECT add_translation_if_not_exists('marketplace_listing', 251, 'en', 'description', 'Luxury penthouse in the city center. Danube view, 2 terraces, jacuzzi. Fully equipped.');
SELECT add_translation_if_not_exists('marketplace_listing', 251, 'ru', 'title', 'Люкс пентхаус 120м2 Центр');
SELECT add_translation_if_not_exists('marketplace_listing', 251, 'ru', 'description', 'Роскошный пентхаус в центре города. Вид на Дунай, 2 террасы, джакузи. Полностью оборудован.');

-- Объявление 252: Дом с бассейном
SELECT add_translation_if_not_exists('marketplace_listing', 252, 'en', 'title', 'House 200m2 with Pool Sremska Kamenica');
SELECT add_translation_if_not_exists('marketplace_listing', 252, 'en', 'description', 'Modern house with pool. 3 bedrooms, large living room, garage for 2 cars.');
SELECT add_translation_if_not_exists('marketplace_listing', 252, 'ru', 'title', 'Дом 200м2 с бассейном Сремска Каменица');
SELECT add_translation_if_not_exists('marketplace_listing', 252, 'ru', 'description', 'Современный дом с бассейном. 3 спальни, большая гостиная, гараж на 2 автомобиля.');

-- Объявление 253: BMW X5
SELECT add_translation_if_not_exists('marketplace_listing', 253, 'en', 'title', 'BMW X5 3.0d 2021 - Like New');
SELECT add_translation_if_not_exists('marketplace_listing', 253, 'en', 'description', 'BMW X5 xDrive30d, M package, full equipment. First owner, service book, warranty until 2025.');
SELECT add_translation_if_not_exists('marketplace_listing', 253, 'ru', 'title', 'BMW X5 3.0d 2021 - как новый');
SELECT add_translation_if_not_exists('marketplace_listing', 253, 'ru', 'description', 'BMW X5 xDrive30d, M пакет, полная комплектация. Первый владелец, сервисная книжка, гарантия до 2025.');

-- Объявление 254: Mercedes-Benz E220d
SELECT add_translation_if_not_exists('marketplace_listing', 254, 'en', 'title', 'Mercedes-Benz E220d 2022');
SELECT add_translation_if_not_exists('marketplace_listing', 254, 'en', 'description', 'Mercedes E class, AMG line, automatic. Navigation, leather seats, panoramic roof.');
SELECT add_translation_if_not_exists('marketplace_listing', 254, 'ru', 'title', 'Mercedes-Benz E220d 2022');
SELECT add_translation_if_not_exists('marketplace_listing', 254, 'ru', 'description', 'Mercedes E класс, AMG линия, автомат. Навигация, кожаные сиденья, панорамная крыша.');

-- Объявление 255: Volkswagen Golf 8
SELECT add_translation_if_not_exists('marketplace_listing', 255, 'en', 'title', 'Volkswagen Golf 8 2.0 TDI 2023');
SELECT add_translation_if_not_exists('marketplace_listing', 255, 'en', 'description', 'Golf VIII generation, 2.0 TDI 150hp, DSG automatic. Style equipment, virtual cockpit, LED headlights. Factory warranty until 2026.');
SELECT add_translation_if_not_exists('marketplace_listing', 255, 'ru', 'title', 'Volkswagen Golf 8 2.0 TDI 2023');
SELECT add_translation_if_not_exists('marketplace_listing', 255, 'ru', 'description', 'Golf VIII поколение, 2.0 TDI 150л.с., DSG автомат. Комплектация Style, виртуальная панель, LED фары. Заводская гарантия до 2026.');

-- Объявление 256: Yamaha MT-07
SELECT add_translation_if_not_exists('marketplace_listing', 256, 'en', 'title', 'Yamaha MT-07 2022 - Perfect Condition');
SELECT add_translation_if_not_exists('marketplace_listing', 256, 'en', 'description', 'Yamaha MT-07, ABS, new tires, serviced. No damage, garaged. 8500km.');
SELECT add_translation_if_not_exists('marketplace_listing', 256, 'ru', 'title', 'Yamaha MT-07 2022 - идеальное состояние');
SELECT add_translation_if_not_exists('marketplace_listing', 256, 'ru', 'description', 'Yamaha MT-07, ABS, новая резина, обслужен. Без повреждений, гаражное хранение. 8500км.');

-- Объявление 257: Создание веб-сайтов
SELECT add_translation_if_not_exists('marketplace_listing', 257, 'en', 'title', 'Website Development - WordPress, React');
SELECT add_translation_if_not_exists('marketplace_listing', 257, 'en', 'description', 'Professional website development. WordPress, React, Node.js. SEO optimization included.');
SELECT add_translation_if_not_exists('marketplace_listing', 257, 'ru', 'title', 'Создание веб-сайтов - WordPress, React');
SELECT add_translation_if_not_exists('marketplace_listing', 257, 'ru', 'description', 'Профессиональная разработка веб-сайтов. WordPress, React, Node.js. SEO оптимизация включена.');

-- Объявление 258: Ремонт квартир
SELECT add_translation_if_not_exists('marketplace_listing', 258, 'en', 'title', 'Apartment Renovation - Complete Service');
SELECT add_translation_if_not_exists('marketplace_listing', 258, 'en', 'description', 'Complete apartment renovation. Tiles, parquet, plastering, painting. Work guarantee.');
SELECT add_translation_if_not_exists('marketplace_listing', 258, 'ru', 'title', 'Ремонт квартир - полный сервис');
SELECT add_translation_if_not_exists('marketplace_listing', 258, 'ru', 'description', 'Комплексный ремонт квартир. Плитка, паркет, штукатурка, покраска. Гарантия на работы.');

-- Объявление 259: Массаж и wellness
SELECT add_translation_if_not_exists('marketplace_listing', 259, 'en', 'title', 'Massage and Wellness Treatments');
SELECT add_translation_if_not_exists('marketplace_listing', 259, 'en', 'description', 'Professional massages: relaxation, medical, sports. Wellness center in the city center.');
SELECT add_translation_if_not_exists('marketplace_listing', 259, 'ru', 'title', 'Массаж и wellness процедуры');
SELECT add_translation_if_not_exists('marketplace_listing', 259, 'ru', 'description', 'Профессиональный массаж: расслабляющий, медицинский, спортивный. Wellness центр в центре города.');

-- Объявление 260: Оборудование для спортзала
SELECT add_translation_if_not_exists('marketplace_listing', 260, 'en', 'title', 'Weights and Gym Equipment - Complete Set');
SELECT add_translation_if_not_exists('marketplace_listing', 260, 'en', 'description', 'Complete home gym equipment set. 200kg weights, bench, rack, bars. Everything like new.');
SELECT add_translation_if_not_exists('marketplace_listing', 260, 'ru', 'title', 'Гантели и оборудование для спортзала - комплект');
SELECT add_translation_if_not_exists('marketplace_listing', 260, 'ru', 'description', 'Полный комплект оборудования для домашнего спортзала. Веса 200кг, скамья, стойка, штанги. Всё как новое.');

-- Объявление 261: Электрический велосипед Trek
SELECT add_translation_if_not_exists('marketplace_listing', 261, 'en', 'title', 'Trek Electric Bicycle 2023');
SELECT add_translation_if_not_exists('marketplace_listing', 261, 'en', 'description', 'Trek e-bike, Bosch motor, 100km range. Hydraulic brakes, 10 speeds. Warranty.');
SELECT add_translation_if_not_exists('marketplace_listing', 261, 'ru', 'title', 'Электрический велосипед Trek 2023');
SELECT add_translation_if_not_exists('marketplace_listing', 261, 'ru', 'description', 'Trek e-bike, мотор Bosch, запас хода 100км. Гидравлические тормоза, 10 скоростей. Гарантия.');

-- Объявление 262: Щенки Golden Retriever
SELECT add_translation_if_not_exists('marketplace_listing', 262, 'en', 'title', 'Golden Retriever Puppies with Papers');
SELECT add_translation_if_not_exists('marketplace_listing', 262, 'en', 'description', 'Purebred Golden Retriever puppies. Pedigree, vaccinated, chipped. Excellent with children.');
SELECT add_translation_if_not_exists('marketplace_listing', 262, 'ru', 'title', 'Щенки Golden Retriever с документами');
SELECT add_translation_if_not_exists('marketplace_listing', 262, 'ru', 'description', 'Чистокровные щенки Golden Retriever. Родословная, привиты, чипированы. Отлично ладят с детьми.');

-- Объявление 263: Британские котята
SELECT add_translation_if_not_exists('marketplace_listing', 263, 'en', 'title', 'British Shorthair Kittens');
SELECT add_translation_if_not_exists('marketplace_listing', 263, 'en', 'description', 'British Shorthair kittens, blue color. Litter trained, vaccinated. 3 months old.');
SELECT add_translation_if_not_exists('marketplace_listing', 263, 'ru', 'title', 'Британские короткошёрстные котята');
SELECT add_translation_if_not_exists('marketplace_listing', 263, 'ru', 'description', 'Британские короткошёрстные котята, голубой окрас. Приучены к лотку, привиты. Возраст 3 месяца.');

-- Объявление 264: Медицинские книги
SELECT add_translation_if_not_exists('marketplace_listing', 264, 'en', 'title', 'Medical Books Set - 50 Titles');
SELECT add_translation_if_not_exists('marketplace_listing', 264, 'en', 'description', 'Medical literature, 50 books. Anatomy, physiology, internal medicine. Excellent condition.');
SELECT add_translation_if_not_exists('marketplace_listing', 264, 'ru', 'title', 'Комплект медицинских книг - 50 наименований');
SELECT add_translation_if_not_exists('marketplace_listing', 264, 'ru', 'description', 'Медицинская литература, 50 книг. Анатомия, физиология, внутренние болезни. Отличное состояние.');

-- Объявление 265: Гитара Yamaha
SELECT add_translation_if_not_exists('marketplace_listing', 265, 'en', 'title', 'Yamaha Guitar C40 with Case');
SELECT add_translation_if_not_exists('marketplace_listing', 265, 'en', 'description', 'Classical guitar Yamaha C40. Case, stand, capo. Ideal for beginners. Like new.');
SELECT add_translation_if_not_exists('marketplace_listing', 265, 'ru', 'title', 'Гитара Yamaha C40 с чехлом');
SELECT add_translation_if_not_exists('marketplace_listing', 265, 'ru', 'description', 'Классическая гитара Yamaha C40. Чехол, стойка, каподастр. Идеально для начинающих. Как новая.');

-- Объявление 266: Детская коляска Chicco
SELECT add_translation_if_not_exists('marketplace_listing', 266, 'en', 'title', 'Chicco Stroller 3in1 - Like New');
SELECT add_translation_if_not_exists('marketplace_listing', 266, 'en', 'description', 'Chicco 3in1 system: stroller, carrier, car seat. Used for 6 months. All accessories.');
SELECT add_translation_if_not_exists('marketplace_listing', 266, 'ru', 'title', 'Коляска Chicco 3в1 - как новая');
SELECT add_translation_if_not_exists('marketplace_listing', 266, 'ru', 'description', 'Система Chicco 3в1: коляска, переноска, автокресло. Использовалась 6 месяцев. Все аксессуары.');

-- Объявление 267: Коллекция LEGO
SELECT add_translation_if_not_exists('marketplace_listing', 267, 'en', 'title', 'LEGO Collection - 15 Sets');
SELECT add_translation_if_not_exists('marketplace_listing', 267, 'en', 'description', 'Large LEGO collection. Star Wars, Technic, City. Complete sets with instructions.');
SELECT add_translation_if_not_exists('marketplace_listing', 267, 'ru', 'title', 'Коллекция LEGO - 15 наборов');
SELECT add_translation_if_not_exists('marketplace_listing', 267, 'ru', 'description', 'Большая коллекция LEGO. Star Wars, Technic, City. Полные наборы с инструкциями.');

-- Объявление 268: Роутер Huawei (уже на русском, добавим английский и сербский)
SELECT add_translation_if_not_exists('marketplace_listing', 268, 'en', 'title', 'Optical Router Huawei HG8546M • 2.4/5GHz • Like New');
SELECT add_translation_if_not_exists('marketplace_listing', 268, 'en', 'description', '🌐 POWERFUL OPTICAL ROUTER FOR YOUR HOME!

✨ ADVANTAGES:
- GPON technology support for optical internet
- Dual-band WiFi (2.4GHz + 5GHz)
- 4 Gigabit Ethernet ports
- High speed up to 300 Mbps

📱 SPECIFICATIONS:
- 2 external antennas for stable signal
- IPTV support
- Easy setup via web interface
- Parental control

🛡️ CONDITION:
- Fully working
- All ports in perfect condition
- Factory firmware

📦 INCLUDED:
- Huawei HG8546M router
- Power adapter
- Network cable

🔥 Get a powerful router for high-speed internet! Call now! 📞');

SELECT add_translation_if_not_exists('marketplace_listing', 268, 'sr', 'title', 'Optički ruter Huawei HG8546M • 2.4/5GHz • Kao nov');
SELECT add_translation_if_not_exists('marketplace_listing', 268, 'sr', 'description', '🌐 MOĆAN OPTIČKI RUTER ZA VAŠ DOM!

✨ PREDNOSTI:
- Podrška GPON tehnologije za optički internet
- Dual-band WiFi (2.4GHz + 5GHz)
- 4 Gigabit Ethernet porta
- Velika brzina do 300 Mbps

📱 KARAKTERISTIKE:
- 2 spoljne antene za stabilan signal
- Podrška za IPTV
- Jednostavno podešavanje preko web interfejsa
- Roditeljska kontrola

🛡️ STANJE:
- Potpuno ispravan
- Svi portovi u savršenom stanju
- Fabrički firmware

📦 U KOMPLETU:
- Ruter Huawei HG8546M
- Adapter za napajanje
- Mrežni kabl

🔥 Nabavite moćan ruter za brzi internet! Pozovite odmah! 📞');

-- Удаляем временную функцию
DROP FUNCTION IF EXISTS add_translation_if_not_exists;