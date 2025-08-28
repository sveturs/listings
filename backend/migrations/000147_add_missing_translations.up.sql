-- Переводы для новых категорий
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
-- Животные и зоотовары
('category', 1011, 'sr', 'name', 'Kućni ljubimci'),
('category', 1011, 'en', 'name', 'Pets & Pet Supplies'),
('category', 1011, 'ru', 'name', 'Животные и зоотовары'),
-- Книги и канцелярия
('category', 1012, 'sr', 'name', 'Knjige i kancelarija'),
('category', 1012, 'en', 'name', 'Books & Stationery'),
('category', 1012, 'ru', 'name', 'Книги и канцелярия'),
-- Детские товары
('category', 1013, 'sr', 'name', 'Dečiji proizvodi'),
('category', 1013, 'en', 'name', 'Kids & Baby'),
('category', 1013, 'ru', 'name', 'Детские товары'),
-- Здоровье и красота
('category', 1014, 'sr', 'name', 'Zdravlje i lepota'),
('category', 1014, 'en', 'name', 'Health & Beauty'),
('category', 1014, 'ru', 'name', 'Здоровье и красота'),
-- Хобби и развлечения
('category', 1015, 'sr', 'name', 'Hobi i zabava'),
('category', 1015, 'en', 'name', 'Hobbies & Entertainment'),
('category', 1015, 'ru', 'name', 'Хобби и развлечения'),
-- Музыкальные инструменты
('category', 1016, 'sr', 'name', 'Muzički instrumenti'),
('category', 1016, 'en', 'name', 'Musical Instruments'),
('category', 1016, 'ru', 'name', 'Музыкальные инструменты'),
-- Антиквариат и искусство
('category', 1017, 'sr', 'name', 'Antikviteti i umetnost'),
('category', 1017, 'en', 'name', 'Antiques & Art'),
('category', 1017, 'ru', 'name', 'Антиквариат и искусство'),
-- Работа и вакансии
('category', 1018, 'sr', 'name', 'Poslovi'),
('category', 1018, 'en', 'name', 'Jobs'),
('category', 1018, 'ru', 'name', 'Работа и вакансии'),
-- Образование
('category', 1019, 'sr', 'name', 'Obrazovanje'),
('category', 1019, 'en', 'name', 'Education'),
('category', 1019, 'ru', 'name', 'Образование'),
-- События и билеты
('category', 1020, 'sr', 'name', 'Događaji i karte'),
('category', 1020, 'en', 'name', 'Events & Tickets'),
('category', 1020, 'ru', 'name', 'События и билеты');

-- Переводы для подкатегорий электроники
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('category', 1105, 'sr', 'name', 'Konzole za igre'),
('category', 1105, 'en', 'name', 'Gaming Consoles'),
('category', 1105, 'ru', 'name', 'Игровые консоли'),
('category', 1106, 'sr', 'name', 'Foto i video'),
('category', 1106, 'en', 'name', 'Photo & Video'),
('category', 1106, 'ru', 'name', 'Фото и видео'),
('category', 1107, 'sr', 'name', 'Pametna kuća'),
('category', 1107, 'en', 'name', 'Smart Home'),
('category', 1107, 'ru', 'name', 'Умный дом'),
('category', 1108, 'sr', 'name', 'Pribor'),
('category', 1108, 'en', 'name', 'Accessories'),
('category', 1108, 'ru', 'name', 'Аксессуары');

-- Переводы для подкатегорий моды
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('category', 1205, 'sr', 'name', 'Dečija odeća'),
('category', 1205, 'en', 'name', 'Kids Clothing'),
('category', 1205, 'ru', 'name', 'Детская одежда'),
('category', 1206, 'sr', 'name', 'Sportska odeća'),
('category', 1206, 'en', 'name', 'Sports Clothing'),
('category', 1206, 'ru', 'name', 'Спортивная одежда'),
('category', 1207, 'sr', 'name', 'Satovi'),
('category', 1207, 'en', 'name', 'Watches'),
('category', 1207, 'ru', 'name', 'Часы'),
('category', 1208, 'sr', 'name', 'Torbe'),
('category', 1208, 'en', 'name', 'Bags'),
('category', 1208, 'ru', 'name', 'Сумки');

-- Переводы для подкатегорий дома и сада
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('category', 1505, 'sr', 'name', 'Kuhinjski pribor'),
('category', 1505, 'en', 'name', 'Kitchenware'),
('category', 1505, 'ru', 'name', 'Посуда'),
('category', 1506, 'sr', 'name', 'Tekstil'),
('category', 1506, 'en', 'name', 'Textiles'),
('category', 1506, 'ru', 'name', 'Текстиль'),
('category', 1507, 'sr', 'name', 'Rasveta'),
('category', 1507, 'en', 'name', 'Lighting'),
('category', 1507, 'ru', 'name', 'Освещение'),
('category', 1508, 'sr', 'name', 'Vodovod i sanitarije'),
('category', 1508, 'en', 'name', 'Plumbing'),
('category', 1508, 'ru', 'name', 'Сантехника');

-- Переводы для подкатегорий спорта
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('category', 2011, 'sr', 'name', 'Bicikli'),
('category', 2011, 'en', 'name', 'Bicycles'),
('category', 2011, 'ru', 'name', 'Велосипеды'),
('category', 2012, 'sr', 'name', 'Vodeni sportovi'),
('category', 2012, 'en', 'name', 'Water Sports'),
('category', 2012, 'ru', 'name', 'Водный спорт'),
('category', 2013, 'sr', 'name', 'Lov i ribolov'),
('category', 2013, 'en', 'name', 'Hunting & Fishing'),
('category', 2013, 'ru', 'name', 'Охота и рыбалка'),
('category', 2014, 'sr', 'name', 'Kampovanje'),
('category', 2014, 'en', 'name', 'Camping & Hiking'),
('category', 2014, 'ru', 'name', 'Кемпинг и туризм');

-- Переводы для атрибутов
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
-- Базовые атрибуты
('attribute', 2001, 'sr', 'name', 'Cena'),
('attribute', 2001, 'en', 'name', 'Price'),
('attribute', 2001, 'ru', 'name', 'Цена'),
('attribute', 2002, 'sr', 'name', 'Stanje'),
('attribute', 2002, 'en', 'name', 'Condition'),
('attribute', 2002, 'ru', 'name', 'Состояние'),
('attribute', 2003, 'sr', 'name', 'Brend'),
('attribute', 2003, 'en', 'name', 'Brand'),
('attribute', 2003, 'ru', 'name', 'Бренд'),
('attribute', 2004, 'sr', 'name', 'Boja'),
('attribute', 2004, 'en', 'name', 'Color'),
('attribute', 2004, 'ru', 'name', 'Цвет');

-- Атрибуты электроники
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 2101, 'sr', 'name', 'Memorija'),
('attribute', 2101, 'en', 'name', 'Storage'),
('attribute', 2101, 'ru', 'name', 'Память'),
('attribute', 2102, 'sr', 'name', 'Operativni sistem'),
('attribute', 2102, 'en', 'name', 'Operating System'),
('attribute', 2102, 'ru', 'name', 'Операционная система'),
('attribute', 2103, 'sr', 'name', 'Procesor'),
('attribute', 2103, 'en', 'name', 'Processor'),
('attribute', 2103, 'ru', 'name', 'Процессор'),
('attribute', 2104, 'sr', 'name', 'RAM memorija'),
('attribute', 2104, 'en', 'name', 'RAM'),
('attribute', 2104, 'ru', 'name', 'Оперативная память'),
('attribute', 2105, 'sr', 'name', 'Tip skladišta'),
('attribute', 2105, 'en', 'name', 'Storage Type'),
('attribute', 2105, 'ru', 'name', 'Тип накопителя');

-- Атрибуты автомобилей
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 2201, 'sr', 'name', 'Model'),
('attribute', 2201, 'en', 'name', 'Model'),
('attribute', 2201, 'ru', 'name', 'Модель'),
('attribute', 2202, 'sr', 'name', 'Godište'),
('attribute', 2202, 'en', 'name', 'Year'),
('attribute', 2202, 'ru', 'name', 'Год выпуска'),
('attribute', 2203, 'sr', 'name', 'Kilometraža'),
('attribute', 2203, 'en', 'name', 'Mileage'),
('attribute', 2203, 'ru', 'name', 'Пробег'),
('attribute', 2204, 'sr', 'name', 'Gorivo'),
('attribute', 2204, 'en', 'name', 'Fuel Type'),
('attribute', 2204, 'ru', 'name', 'Тип топлива'),
('attribute', 2205, 'sr', 'name', 'Menjač'),
('attribute', 2205, 'en', 'name', 'Transmission'),
('attribute', 2205, 'ru', 'name', 'Коробка передач'),
('attribute', 2206, 'sr', 'name', 'Tip karoserije'),
('attribute', 2206, 'en', 'name', 'Body Type'),
('attribute', 2206, 'ru', 'name', 'Тип кузова');

-- Атрибуты недвижимости
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 2301, 'sr', 'name', 'Površina (m²)'),
('attribute', 2301, 'en', 'name', 'Area (m²)'),
('attribute', 2301, 'ru', 'name', 'Площадь (м²)'),
('attribute', 2302, 'sr', 'name', 'Broj soba'),
('attribute', 2302, 'en', 'name', 'Number of Rooms'),
('attribute', 2302, 'ru', 'name', 'Количество комнат'),
('attribute', 2303, 'sr', 'name', 'Sprat'),
('attribute', 2303, 'en', 'name', 'Floor'),
('attribute', 2303, 'ru', 'name', 'Этаж'),
('attribute', 2304, 'sr', 'name', 'Namešteno'),
('attribute', 2304, 'en', 'name', 'Furnished'),
('attribute', 2304, 'ru', 'name', 'Меблировано'),
('attribute', 2305, 'sr', 'name', 'Parking'),
('attribute', 2305, 'en', 'name', 'Parking'),
('attribute', 2305, 'ru', 'name', 'Парковка'),
('attribute', 2306, 'sr', 'name', 'Balkon/terasa'),
('attribute', 2306, 'en', 'name', 'Balcony/Terrace'),
('attribute', 2306, 'ru', 'name', 'Балкон/терраса');

-- Переводы для новых общих атрибутов
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 2601, 'sr', 'name', 'Lokacija'),
('attribute', 2601, 'en', 'name', 'Location'),
('attribute', 2601, 'ru', 'name', 'Местоположение'),
('attribute', 2602, 'sr', 'name', 'Dostava'),
('attribute', 2602, 'en', 'name', 'Delivery Available'),
('attribute', 2602, 'ru', 'name', 'Доставка'),
('attribute', 2603, 'sr', 'name', 'Dogovor'),
('attribute', 2603, 'en', 'name', 'Negotiable'),
('attribute', 2603, 'ru', 'name', 'Торг'),
('attribute', 2604, 'sr', 'name', 'Garancija'),
('attribute', 2604, 'en', 'name', 'Warranty'),
('attribute', 2604, 'ru', 'name', 'Гарантия'),
('attribute', 2605, 'sr', 'name', 'Povrat'),
('attribute', 2605, 'en', 'name', 'Return Policy'),
('attribute', 2605, 'ru', 'name', 'Возврат');

-- Переводы для новых атрибутов электроники
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 2701, 'sr', 'name', 'Veličina ekrana'),
('attribute', 2701, 'en', 'name', 'Screen Size'),
('attribute', 2701, 'ru', 'name', 'Размер экрана'),
('attribute', 2702, 'sr', 'name', 'Trajanje baterije'),
('attribute', 2702, 'en', 'name', 'Battery Life'),
('attribute', 2702, 'ru', 'name', 'Время работы батареи'),
('attribute', 2703, 'sr', 'name', 'Povezivanje'),
('attribute', 2703, 'en', 'name', 'Connectivity'),
('attribute', 2703, 'ru', 'name', 'Подключение'),
('attribute', 2704, 'sr', 'name', 'Rezolucija'),
('attribute', 2704, 'en', 'name', 'Resolution'),
('attribute', 2704, 'ru', 'name', 'Разрешение');

-- Переводы для новых атрибутов моды
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 2801, 'sr', 'name', 'Veličina'),
('attribute', 2801, 'en', 'name', 'Size'),
('attribute', 2801, 'ru', 'name', 'Размер'),
('attribute', 2802, 'sr', 'name', 'Materijal'),
('attribute', 2802, 'en', 'name', 'Material'),
('attribute', 2802, 'ru', 'name', 'Материал'),
('attribute', 2803, 'sr', 'name', 'Pol'),
('attribute', 2803, 'en', 'name', 'Gender'),
('attribute', 2803, 'ru', 'name', 'Пол'),
('attribute', 2804, 'sr', 'name', 'Sezona'),
('attribute', 2804, 'en', 'name', 'Season'),
('attribute', 2804, 'ru', 'name', 'Сезон');

-- Переводы для всех остальных атрибутов
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
-- Атрибуты недвижимости (продолжение)
('attribute', 2307, 'sr', 'name', 'Površina kuće (m²)'),
('attribute', 2307, 'en', 'name', 'House Area (m²)'),
('attribute', 2307, 'ru', 'name', 'Площадь дома (м²)'),
('attribute', 2308, 'sr', 'name', 'Površina placa (m²)'),
('attribute', 2308, 'en', 'name', 'Land Area (m²)'),
('attribute', 2308, 'ru', 'name', 'Площадь участка (м²)'),
('attribute', 2309, 'sr', 'name', 'Broj kupatila'),
('attribute', 2309, 'en', 'name', 'Number of Bathrooms'),
('attribute', 2309, 'ru', 'name', 'Количество ванных'),
('attribute', 2310, 'sr', 'name', 'Bašta'),
('attribute', 2310, 'en', 'name', 'Garden'),
('attribute', 2310, 'ru', 'name', 'Сад'),
('attribute', 2311, 'sr', 'name', 'Garaža'),
('attribute', 2311, 'en', 'name', 'Garage'),
('attribute', 2311, 'ru', 'name', 'Гараж'),
-- Новые атрибуты недвижимости
('attribute', 2901, 'sr', 'name', 'Tip grejanja'),
('attribute', 2901, 'en', 'name', 'Heating Type'),
('attribute', 2901, 'ru', 'name', 'Тип отопления'),
('attribute', 2902, 'sr', 'name', 'Godina izgradnje'),
('attribute', 2902, 'en', 'name', 'Construction Year'),
('attribute', 2902, 'ru', 'name', 'Год постройки'),
('attribute', 2903, 'sr', 'name', 'Lift'),
('attribute', 2903, 'en', 'name', 'Elevator'),
('attribute', 2903, 'ru', 'name', 'Лифт'),
('attribute', 2904, 'sr', 'name', 'Obezbeđenje'),
('attribute', 2904, 'en', 'name', 'Security'),
('attribute', 2904, 'ru', 'name', 'Охрана');

-- Атрибуты промышленности
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 2401, 'sr', 'name', 'Radnih sati'),
('attribute', 2401, 'en', 'name', 'Working Hours'),
('attribute', 2401, 'ru', 'name', 'Моточасы'),
('attribute', 2402, 'sr', 'name', 'Snaga (KS)'),
('attribute', 2402, 'en', 'name', 'Power (HP)'),
('attribute', 2402, 'ru', 'name', 'Мощность (л.с.)');

-- Атрибуты услуг
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 2501, 'sr', 'name', 'Tip usluge'),
('attribute', 2501, 'en', 'name', 'Service Type'),
('attribute', 2501, 'ru', 'name', 'Тип услуги'),
('attribute', 2502, 'sr', 'name', 'Godine iskustva'),
('attribute', 2502, 'en', 'name', 'Years of Experience'),
('attribute', 2502, 'ru', 'name', 'Лет опыта'),
('attribute', 2503, 'sr', 'name', 'Dostupnost'),
('attribute', 2503, 'en', 'name', 'Availability'),
('attribute', 2503, 'ru', 'name', 'Доступность'),
('attribute', 2504, 'sr', 'name', 'Oblast rada'),
('attribute', 2504, 'en', 'name', 'Service Area'),
('attribute', 2504, 'ru', 'name', 'Область работы');

-- Новые атрибуты автомобилей
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 3001, 'sr', 'name', 'Zapremina motora'),
('attribute', 3001, 'en', 'name', 'Engine Size'),
('attribute', 3001, 'ru', 'name', 'Объем двигателя'),
('attribute', 3002, 'sr', 'name', 'Broj vrata'),
('attribute', 3002, 'en', 'name', 'Number of Doors'),
('attribute', 3002, 'ru', 'name', 'Количество дверей'),
('attribute', 3003, 'sr', 'name', 'Broj sedišta'),
('attribute', 3003, 'en', 'name', 'Number of Seats'),
('attribute', 3003, 'ru', 'name', 'Количество мест'),
('attribute', 3004, 'sr', 'name', 'Pogon'),
('attribute', 3004, 'en', 'name', 'Drive Type'),
('attribute', 3004, 'ru', 'name', 'Привод');

-- Новые атрибуты услуг
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text) VALUES
('attribute', 3101, 'sr', 'name', 'Tip cene'),
('attribute', 3101, 'en', 'name', 'Price Type'),
('attribute', 3101, 'ru', 'name', 'Тип цены'),
('attribute', 3102, 'sr', 'name', 'Sertifikati'),
('attribute', 3102, 'en', 'name', 'Certifications'),
('attribute', 3102, 'ru', 'name', 'Сертификаты'),
('attribute', 3103, 'sr', 'name', 'Jezici'),
('attribute', 3103, 'en', 'name', 'Languages'),
('attribute', 3103, 'ru', 'name', 'Языки'),
('attribute', 3104, 'sr', 'name', 'Portfolio'),
('attribute', 3104, 'en', 'name', 'Portfolio'),
('attribute', 3104, 'ru', 'name', 'Портфолио');

