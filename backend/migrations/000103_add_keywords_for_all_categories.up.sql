-- Добавление ключевых слов для всех категорий на всех языках (en, ru, sr)

-- Очищаем существующие ключевые слова для обновления
DELETE FROM category_keywords WHERE source = 'manual';

-- ==========================================
-- ОСНОВНЫЕ КАТЕГОРИИ
-- ==========================================

-- Elektronika (1001)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1001, 'electronics', 'en', 10.0, 'main', 'manual'),
(1001, 'gadgets', 'en', 8.0, 'synonym', 'manual'),
(1001, 'devices', 'en', 7.0, 'synonym', 'manual'),
(1001, 'tech', 'en', 7.0, 'synonym', 'manual'),
(1001, 'technology', 'en', 6.0, 'context', 'manual'),
-- Russian
(1001, 'электроника', 'ru', 10.0, 'main', 'manual'),
(1001, 'гаджеты', 'ru', 8.0, 'synonym', 'manual'),
(1001, 'устройства', 'ru', 7.0, 'synonym', 'manual'),
(1001, 'техника', 'ru', 7.0, 'synonym', 'manual'),
(1001, 'технологии', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1001, 'elektronika', 'sr', 10.0, 'main', 'manual'),
(1001, 'gedžeti', 'sr', 8.0, 'synonym', 'manual'),
(1001, 'uređaji', 'sr', 7.0, 'synonym', 'manual'),
(1001, 'tehnika', 'sr', 7.0, 'synonym', 'manual'),
(1001, 'tehnologija', 'sr', 6.0, 'context', 'manual');

-- Moda (1002)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1002, 'fashion', 'en', 10.0, 'main', 'manual'),
(1002, 'clothing', 'en', 9.0, 'synonym', 'manual'),
(1002, 'apparel', 'en', 8.0, 'synonym', 'manual'),
(1002, 'style', 'en', 7.0, 'context', 'manual'),
(1002, 'wear', 'en', 6.0, 'context', 'manual'),
-- Russian
(1002, 'мода', 'ru', 10.0, 'main', 'manual'),
(1002, 'одежда', 'ru', 9.0, 'synonym', 'manual'),
(1002, 'стиль', 'ru', 7.0, 'context', 'manual'),
(1002, 'наряд', 'ru', 6.0, 'synonym', 'manual'),
(1002, 'гардероб', 'ru', 5.0, 'context', 'manual'),
-- Serbian
(1002, 'moda', 'sr', 10.0, 'main', 'manual'),
(1002, 'odeća', 'sr', 9.0, 'synonym', 'manual'),
(1002, 'stil', 'sr', 7.0, 'context', 'manual'),
(1002, 'garderoba', 'sr', 6.0, 'context', 'manual'),
(1002, 'odevanje', 'sr', 5.0, 'context', 'manual');

-- Automobili (1003)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1003, 'automotive', 'en', 10.0, 'main', 'manual'),
(1003, 'cars', 'en', 9.0, 'synonym', 'manual'),
(1003, 'vehicles', 'en', 8.0, 'synonym', 'manual'),
(1003, 'auto', 'en', 7.0, 'synonym', 'manual'),
(1003, 'motors', 'en', 6.0, 'context', 'manual'),
-- Russian
(1003, 'автомобили', 'ru', 10.0, 'main', 'manual'),
(1003, 'машины', 'ru', 9.0, 'synonym', 'manual'),
(1003, 'авто', 'ru', 8.0, 'synonym', 'manual'),
(1003, 'транспорт', 'ru', 7.0, 'context', 'manual'),
(1003, 'автотранспорт', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1003, 'automobili', 'sr', 10.0, 'main', 'manual'),
(1003, 'kola', 'sr', 9.0, 'synonym', 'manual'),
(1003, 'vozila', 'sr', 8.0, 'synonym', 'manual'),
(1003, 'auto', 'sr', 7.0, 'synonym', 'manual'),
(1003, 'motori', 'sr', 6.0, 'context', 'manual');

-- Nekretnine (1004)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1004, 'real estate', 'en', 10.0, 'main', 'manual'),
(1004, 'property', 'en', 9.0, 'synonym', 'manual'),
(1004, 'realty', 'en', 8.0, 'synonym', 'manual'),
(1004, 'housing', 'en', 7.0, 'context', 'manual'),
(1004, 'properties', 'en', 6.0, 'synonym', 'manual'),
-- Russian
(1004, 'недвижимость', 'ru', 10.0, 'main', 'manual'),
(1004, 'имущество', 'ru', 8.0, 'synonym', 'manual'),
(1004, 'жилье', 'ru', 7.0, 'context', 'manual'),
(1004, 'помещения', 'ru', 6.0, 'context', 'manual'),
(1004, 'объекты', 'ru', 5.0, 'context', 'manual'),
-- Serbian
(1004, 'nekretnine', 'sr', 10.0, 'main', 'manual'),
(1004, 'imovina', 'sr', 8.0, 'synonym', 'manual'),
(1004, 'nepokretnosti', 'sr', 7.0, 'synonym', 'manual'),
(1004, 'stanovanje', 'sr', 6.0, 'context', 'manual'),
(1004, 'objekti', 'sr', 5.0, 'context', 'manual');

-- Dom i bašta (1005)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1005, 'home', 'en', 10.0, 'main', 'manual'),
(1005, 'garden', 'en', 9.0, 'main', 'manual'),
(1005, 'household', 'en', 8.0, 'synonym', 'manual'),
(1005, 'furniture', 'en', 7.0, 'context', 'manual'),
(1005, 'decor', 'en', 6.0, 'context', 'manual'),
-- Russian
(1005, 'дом', 'ru', 10.0, 'main', 'manual'),
(1005, 'сад', 'ru', 9.0, 'main', 'manual'),
(1005, 'быт', 'ru', 8.0, 'synonym', 'manual'),
(1005, 'мебель', 'ru', 7.0, 'context', 'manual'),
(1005, 'интерьер', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1005, 'dom', 'sr', 10.0, 'main', 'manual'),
(1005, 'bašta', 'sr', 9.0, 'main', 'manual'),
(1005, 'kuća', 'sr', 8.0, 'synonym', 'manual'),
(1005, 'nameštaj', 'sr', 7.0, 'context', 'manual'),
(1005, 'dvorište', 'sr', 6.0, 'context', 'manual');

-- Poljoprivreda (1006)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1006, 'agriculture', 'en', 10.0, 'main', 'manual'),
(1006, 'farming', 'en', 9.0, 'synonym', 'manual'),
(1006, 'farm', 'en', 8.0, 'synonym', 'manual'),
(1006, 'crops', 'en', 7.0, 'context', 'manual'),
(1006, 'livestock', 'en', 6.0, 'context', 'manual'),
-- Russian
(1006, 'сельское хозяйство', 'ru', 10.0, 'main', 'manual'),
(1006, 'фермерство', 'ru', 9.0, 'synonym', 'manual'),
(1006, 'агро', 'ru', 8.0, 'synonym', 'manual'),
(1006, 'скот', 'ru', 7.0, 'context', 'manual'),
(1006, 'урожай', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1006, 'poljoprivreda', 'sr', 10.0, 'main', 'manual'),
(1006, 'zemljoradnja', 'sr', 9.0, 'synonym', 'manual'),
(1006, 'farma', 'sr', 8.0, 'synonym', 'manual'),
(1006, 'stoka', 'sr', 7.0, 'context', 'manual'),
(1006, 'usevi', 'sr', 6.0, 'context', 'manual');

-- Industrija (1007)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1007, 'industrial', 'en', 10.0, 'main', 'manual'),
(1007, 'industry', 'en', 9.0, 'synonym', 'manual'),
(1007, 'manufacturing', 'en', 8.0, 'synonym', 'manual'),
(1007, 'machinery', 'en', 7.0, 'context', 'manual'),
(1007, 'equipment', 'en', 6.0, 'context', 'manual'),
-- Russian
(1007, 'промышленность', 'ru', 10.0, 'main', 'manual'),
(1007, 'индустрия', 'ru', 9.0, 'synonym', 'manual'),
(1007, 'производство', 'ru', 8.0, 'synonym', 'manual'),
(1007, 'оборудование', 'ru', 7.0, 'context', 'manual'),
(1007, 'станки', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1007, 'industrija', 'sr', 10.0, 'main', 'manual'),
(1007, 'proizvodnja', 'sr', 8.0, 'synonym', 'manual'),
(1007, 'mašine', 'sr', 7.0, 'context', 'manual'),
(1007, 'oprema', 'sr', 6.0, 'context', 'manual'),
(1007, 'fabrika', 'sr', 5.0, 'context', 'manual');

-- Hrana i piće (1008)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1008, 'food', 'en', 10.0, 'main', 'manual'),
(1008, 'beverages', 'en', 9.0, 'main', 'manual'),
(1008, 'drinks', 'en', 8.0, 'synonym', 'manual'),
(1008, 'cuisine', 'en', 7.0, 'context', 'manual'),
(1008, 'groceries', 'en', 6.0, 'context', 'manual'),
-- Russian
(1008, 'еда', 'ru', 10.0, 'main', 'manual'),
(1008, 'напитки', 'ru', 9.0, 'main', 'manual'),
(1008, 'продукты', 'ru', 8.0, 'synonym', 'manual'),
(1008, 'питание', 'ru', 7.0, 'context', 'manual'),
(1008, 'кухня', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1008, 'hrana', 'sr', 10.0, 'main', 'manual'),
(1008, 'piće', 'sr', 9.0, 'main', 'manual'),
(1008, 'namirnice', 'sr', 8.0, 'synonym', 'manual'),
(1008, 'ishrana', 'sr', 7.0, 'context', 'manual'),
(1008, 'kuhinja', 'sr', 6.0, 'context', 'manual');

-- Usluge (1009)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1009, 'services', 'en', 10.0, 'main', 'manual'),
(1009, 'service', 'en', 9.0, 'synonym', 'manual'),
(1009, 'professional', 'en', 7.0, 'context', 'manual'),
(1009, 'assistance', 'en', 6.0, 'context', 'manual'),
(1009, 'help', 'en', 5.0, 'context', 'manual'),
-- Russian
(1009, 'услуги', 'ru', 10.0, 'main', 'manual'),
(1009, 'сервис', 'ru', 9.0, 'synonym', 'manual'),
(1009, 'обслуживание', 'ru', 8.0, 'synonym', 'manual'),
(1009, 'помощь', 'ru', 6.0, 'context', 'manual'),
(1009, 'работы', 'ru', 5.0, 'context', 'manual'),
-- Serbian
(1009, 'usluge', 'sr', 10.0, 'main', 'manual'),
(1009, 'servis', 'sr', 9.0, 'synonym', 'manual'),
(1009, 'pomoć', 'sr', 6.0, 'context', 'manual'),
(1009, 'podrška', 'sr', 5.0, 'context', 'manual'),
(1009, 'rad', 'sr', 4.0, 'context', 'manual');

-- Sport i rekreacija (1010)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1010, 'sports', 'en', 10.0, 'main', 'manual'),
(1010, 'recreation', 'en', 9.0, 'main', 'manual'),
(1010, 'fitness', 'en', 8.0, 'synonym', 'manual'),
(1010, 'exercise', 'en', 7.0, 'context', 'manual'),
(1010, 'gym', 'en', 6.0, 'context', 'manual'),
-- Russian
(1010, 'спорт', 'ru', 10.0, 'main', 'manual'),
(1010, 'отдых', 'ru', 9.0, 'main', 'manual'),
(1010, 'фитнес', 'ru', 8.0, 'synonym', 'manual'),
(1010, 'тренировки', 'ru', 7.0, 'context', 'manual'),
(1010, 'активность', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1010, 'sport', 'sr', 10.0, 'main', 'manual'),
(1010, 'rekreacija', 'sr', 9.0, 'main', 'manual'),
(1010, 'fitnes', 'sr', 8.0, 'synonym', 'manual'),
(1010, 'vežbanje', 'sr', 7.0, 'context', 'manual'),
(1010, 'aktivnost', 'sr', 6.0, 'context', 'manual');

-- Pets (1011)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1011, 'pets', 'en', 10.0, 'main', 'manual'),
(1011, 'animals', 'en', 9.0, 'synonym', 'manual'),
(1011, 'pet supplies', 'en', 8.0, 'context', 'manual'),
(1011, 'dogs', 'en', 7.0, 'context', 'manual'),
(1011, 'cats', 'en', 6.0, 'context', 'manual'),
-- Russian
(1011, 'питомцы', 'ru', 10.0, 'main', 'manual'),
(1011, 'животные', 'ru', 9.0, 'synonym', 'manual'),
(1011, 'домашние животные', 'ru', 8.0, 'synonym', 'manual'),
(1011, 'собаки', 'ru', 7.0, 'context', 'manual'),
(1011, 'кошки', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1011, 'kućni ljubimci', 'sr', 10.0, 'main', 'manual'),
(1011, 'ljubimci', 'sr', 9.0, 'synonym', 'manual'),
(1011, 'životinje', 'sr', 8.0, 'synonym', 'manual'),
(1011, 'psi', 'sr', 7.0, 'context', 'manual'),
(1011, 'mačke', 'sr', 6.0, 'context', 'manual');

-- Books & Stationery (1012)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1012, 'books', 'en', 10.0, 'main', 'manual'),
(1012, 'stationery', 'en', 9.0, 'main', 'manual'),
(1012, 'literature', 'en', 8.0, 'synonym', 'manual'),
(1012, 'office supplies', 'en', 7.0, 'context', 'manual'),
(1012, 'reading', 'en', 6.0, 'context', 'manual'),
-- Russian
(1012, 'книги', 'ru', 10.0, 'main', 'manual'),
(1012, 'канцтовары', 'ru', 9.0, 'main', 'manual'),
(1012, 'литература', 'ru', 8.0, 'synonym', 'manual'),
(1012, 'канцелярия', 'ru', 7.0, 'synonym', 'manual'),
(1012, 'чтение', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1012, 'knjige', 'sr', 10.0, 'main', 'manual'),
(1012, 'kancelarijski materijal', 'sr', 9.0, 'main', 'manual'),
(1012, 'literatura', 'sr', 8.0, 'synonym', 'manual'),
(1012, 'školski pribor', 'sr', 7.0, 'context', 'manual'),
(1012, 'čitanje', 'sr', 6.0, 'context', 'manual');

-- Kids & Baby (1013)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1013, 'kids', 'en', 10.0, 'main', 'manual'),
(1013, 'baby', 'en', 9.0, 'main', 'manual'),
(1013, 'children', 'en', 8.0, 'synonym', 'manual'),
(1013, 'toys', 'en', 7.0, 'context', 'manual'),
(1013, 'infant', 'en', 6.0, 'synonym', 'manual'),
-- Russian
(1013, 'дети', 'ru', 10.0, 'main', 'manual'),
(1013, 'малыши', 'ru', 9.0, 'synonym', 'manual'),
(1013, 'детские товары', 'ru', 8.0, 'context', 'manual'),
(1013, 'игрушки', 'ru', 7.0, 'context', 'manual'),
(1013, 'младенцы', 'ru', 6.0, 'synonym', 'manual'),
-- Serbian
(1013, 'deca', 'sr', 10.0, 'main', 'manual'),
(1013, 'bebe', 'sr', 9.0, 'main', 'manual'),
(1013, 'dečiji', 'sr', 8.0, 'synonym', 'manual'),
(1013, 'igračke', 'sr', 7.0, 'context', 'manual'),
(1013, 'mališani', 'sr', 6.0, 'synonym', 'manual');

-- Health & Beauty (1014)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1014, 'health', 'en', 10.0, 'main', 'manual'),
(1014, 'beauty', 'en', 9.0, 'main', 'manual'),
(1014, 'cosmetics', 'en', 8.0, 'synonym', 'manual'),
(1014, 'wellness', 'en', 7.0, 'context', 'manual'),
(1014, 'care', 'en', 6.0, 'context', 'manual'),
-- Russian
(1014, 'здоровье', 'ru', 10.0, 'main', 'manual'),
(1014, 'красота', 'ru', 9.0, 'main', 'manual'),
(1014, 'косметика', 'ru', 8.0, 'synonym', 'manual'),
(1014, 'уход', 'ru', 7.0, 'context', 'manual'),
(1014, 'велнес', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1014, 'zdravlje', 'sr', 10.0, 'main', 'manual'),
(1014, 'lepota', 'sr', 9.0, 'main', 'manual'),
(1014, 'kozmetika', 'sr', 8.0, 'synonym', 'manual'),
(1014, 'nega', 'sr', 7.0, 'context', 'manual'),
(1014, 'wellness', 'sr', 6.0, 'context', 'manual');

-- Hobbies & Entertainment (1015)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1015, 'hobbies', 'en', 10.0, 'main', 'manual'),
(1015, 'entertainment', 'en', 9.0, 'main', 'manual'),
(1015, 'leisure', 'en', 8.0, 'synonym', 'manual'),
(1015, 'fun', 'en', 7.0, 'context', 'manual'),
(1015, 'games', 'en', 6.0, 'context', 'manual'),
-- Russian
(1015, 'хобби', 'ru', 10.0, 'main', 'manual'),
(1015, 'развлечения', 'ru', 9.0, 'main', 'manual'),
(1015, 'досуг', 'ru', 8.0, 'synonym', 'manual'),
(1015, 'увлечения', 'ru', 7.0, 'synonym', 'manual'),
(1015, 'игры', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1015, 'hobiji', 'sr', 10.0, 'main', 'manual'),
(1015, 'zabava', 'sr', 9.0, 'main', 'manual'),
(1015, 'razonoda', 'sr', 8.0, 'synonym', 'manual'),
(1015, 'slobodno vreme', 'sr', 7.0, 'context', 'manual'),
(1015, 'igre', 'sr', 6.0, 'context', 'manual');

-- Musical Instruments (1016)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1016, 'musical instruments', 'en', 10.0, 'main', 'manual'),
(1016, 'instruments', 'en', 9.0, 'synonym', 'manual'),
(1016, 'music', 'en', 8.0, 'context', 'manual'),
(1016, 'guitar', 'en', 7.0, 'context', 'manual'),
(1016, 'piano', 'en', 6.0, 'context', 'manual'),
-- Russian
(1016, 'музыкальные инструменты', 'ru', 10.0, 'main', 'manual'),
(1016, 'инструменты', 'ru', 9.0, 'synonym', 'manual'),
(1016, 'музыка', 'ru', 8.0, 'context', 'manual'),
(1016, 'гитара', 'ru', 7.0, 'context', 'manual'),
(1016, 'пианино', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1016, 'muzički instrumenti', 'sr', 10.0, 'main', 'manual'),
(1016, 'instrumenti', 'sr', 9.0, 'synonym', 'manual'),
(1016, 'muzika', 'sr', 8.0, 'context', 'manual'),
(1016, 'gitara', 'sr', 7.0, 'context', 'manual'),
(1016, 'klavir', 'sr', 6.0, 'context', 'manual');

-- Antiques & Art (1017)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1017, 'antiques', 'en', 10.0, 'main', 'manual'),
(1017, 'art', 'en', 9.0, 'main', 'manual'),
(1017, 'collectibles', 'en', 8.0, 'synonym', 'manual'),
(1017, 'vintage', 'en', 7.0, 'context', 'manual'),
(1017, 'paintings', 'en', 6.0, 'context', 'manual'),
-- Russian
(1017, 'антиквариат', 'ru', 10.0, 'main', 'manual'),
(1017, 'искусство', 'ru', 9.0, 'main', 'manual'),
(1017, 'коллекционирование', 'ru', 8.0, 'synonym', 'manual'),
(1017, 'винтаж', 'ru', 7.0, 'context', 'manual'),
(1017, 'картины', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1017, 'antikviteti', 'sr', 10.0, 'main', 'manual'),
(1017, 'umetnost', 'sr', 9.0, 'main', 'manual'),
(1017, 'kolekcionarstvo', 'sr', 8.0, 'synonym', 'manual'),
(1017, 'starinarnice', 'sr', 7.0, 'context', 'manual'),
(1017, 'slike', 'sr', 6.0, 'context', 'manual');

-- Jobs (1018)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1018, 'jobs', 'en', 10.0, 'main', 'manual'),
(1018, 'employment', 'en', 9.0, 'synonym', 'manual'),
(1018, 'work', 'en', 8.0, 'synonym', 'manual'),
(1018, 'career', 'en', 7.0, 'context', 'manual'),
(1018, 'vacancy', 'en', 6.0, 'context', 'manual'),
-- Russian
(1018, 'работа', 'ru', 10.0, 'main', 'manual'),
(1018, 'вакансии', 'ru', 9.0, 'synonym', 'manual'),
(1018, 'трудоустройство', 'ru', 8.0, 'synonym', 'manual'),
(1018, 'карьера', 'ru', 7.0, 'context', 'manual'),
(1018, 'должность', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1018, 'poslovi', 'sr', 10.0, 'main', 'manual'),
(1018, 'zaposlenje', 'sr', 9.0, 'synonym', 'manual'),
(1018, 'rad', 'sr', 8.0, 'synonym', 'manual'),
(1018, 'karijera', 'sr', 7.0, 'context', 'manual'),
(1018, 'posao', 'sr', 6.0, 'synonym', 'manual');

-- Education (1019)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1019, 'education', 'en', 10.0, 'main', 'manual'),
(1019, 'training', 'en', 9.0, 'synonym', 'manual'),
(1019, 'courses', 'en', 8.0, 'synonym', 'manual'),
(1019, 'learning', 'en', 7.0, 'context', 'manual'),
(1019, 'school', 'en', 6.0, 'context', 'manual'),
-- Russian
(1019, 'образование', 'ru', 10.0, 'main', 'manual'),
(1019, 'обучение', 'ru', 9.0, 'synonym', 'manual'),
(1019, 'курсы', 'ru', 8.0, 'synonym', 'manual'),
(1019, 'учеба', 'ru', 7.0, 'context', 'manual'),
(1019, 'школа', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1019, 'obrazovanje', 'sr', 10.0, 'main', 'manual'),
(1019, 'obuka', 'sr', 9.0, 'synonym', 'manual'),
(1019, 'kursevi', 'sr', 8.0, 'synonym', 'manual'),
(1019, 'učenje', 'sr', 7.0, 'context', 'manual'),
(1019, 'škola', 'sr', 6.0, 'context', 'manual');

-- Events & Tickets (1020)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1020, 'events', 'en', 10.0, 'main', 'manual'),
(1020, 'tickets', 'en', 9.0, 'main', 'manual'),
(1020, 'concerts', 'en', 8.0, 'context', 'manual'),
(1020, 'shows', 'en', 7.0, 'context', 'manual'),
(1020, 'festival', 'en', 6.0, 'context', 'manual'),
-- Russian
(1020, 'мероприятия', 'ru', 10.0, 'main', 'manual'),
(1020, 'билеты', 'ru', 9.0, 'main', 'manual'),
(1020, 'концерты', 'ru', 8.0, 'context', 'manual'),
(1020, 'события', 'ru', 7.0, 'synonym', 'manual'),
(1020, 'фестиваль', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1020, 'događaji', 'sr', 10.0, 'main', 'manual'),
(1020, 'karte', 'sr', 9.0, 'main', 'manual'),
(1020, 'koncerti', 'sr', 8.0, 'context', 'manual'),
(1020, 'manifestacije', 'sr', 7.0, 'synonym', 'manual'),
(1020, 'festival', 'sr', 6.0, 'context', 'manual');

-- ==========================================
-- ПОДКАТЕГОРИИ
-- ==========================================

-- Pametni telefoni (1101)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1101, 'smartphones', 'en', 10.0, 'main', 'manual'),
(1101, 'phones', 'en', 9.0, 'synonym', 'manual'),
(1101, 'mobile', 'en', 8.0, 'synonym', 'manual'),
(1101, 'iphone', 'en', 7.0, 'brand', 'manual'),
(1101, 'android', 'en', 6.0, 'context', 'manual'),
-- Russian
(1101, 'смартфоны', 'ru', 10.0, 'main', 'manual'),
(1101, 'телефоны', 'ru', 9.0, 'synonym', 'manual'),
(1101, 'мобильные', 'ru', 8.0, 'synonym', 'manual'),
(1101, 'айфон', 'ru', 7.0, 'brand', 'manual'),
(1101, 'андроид', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1101, 'pametni telefoni', 'sr', 10.0, 'main', 'manual'),
(1101, 'telefoni', 'sr', 9.0, 'synonym', 'manual'),
(1101, 'mobilni', 'sr', 8.0, 'synonym', 'manual'),
(1101, 'ajfon', 'sr', 7.0, 'brand', 'manual'),
(1101, 'android', 'sr', 6.0, 'context', 'manual');

-- Računari (1102)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1102, 'computers', 'en', 10.0, 'main', 'manual'),
(1102, 'laptop', 'en', 9.0, 'synonym', 'manual'),
(1102, 'desktop', 'en', 8.0, 'synonym', 'manual'),
(1102, 'pc', 'en', 7.0, 'synonym', 'manual'),
(1102, 'notebook', 'en', 6.0, 'synonym', 'manual'),
-- Russian
(1102, 'компьютеры', 'ru', 10.0, 'main', 'manual'),
(1102, 'ноутбук', 'ru', 9.0, 'synonym', 'manual'),
(1102, 'пк', 'ru', 8.0, 'synonym', 'manual'),
(1102, 'лэптоп', 'ru', 7.0, 'synonym', 'manual'),
(1102, 'десктоп', 'ru', 6.0, 'synonym', 'manual'),
-- Serbian
(1102, 'računari', 'sr', 10.0, 'main', 'manual'),
(1102, 'laptop', 'sr', 9.0, 'synonym', 'manual'),
(1102, 'desktop', 'sr', 8.0, 'synonym', 'manual'),
(1102, 'kompjuter', 'sr', 7.0, 'synonym', 'manual'),
(1102, 'notebook', 'sr', 6.0, 'synonym', 'manual');

-- TV i audio (1103)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1103, 'tv', 'en', 10.0, 'main', 'manual'),
(1103, 'television', 'en', 9.0, 'synonym', 'manual'),
(1103, 'audio', 'en', 8.0, 'main', 'manual'),
(1103, 'speakers', 'en', 7.0, 'context', 'manual'),
(1103, 'sound', 'en', 6.0, 'context', 'manual'),
-- Russian
(1103, 'телевизор', 'ru', 10.0, 'main', 'manual'),
(1103, 'тв', 'ru', 9.0, 'synonym', 'manual'),
(1103, 'аудио', 'ru', 8.0, 'main', 'manual'),
(1103, 'колонки', 'ru', 7.0, 'context', 'manual'),
(1103, 'звук', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1103, 'tv', 'sr', 10.0, 'main', 'manual'),
(1103, 'televizor', 'sr', 9.0, 'synonym', 'manual'),
(1103, 'audio', 'sr', 8.0, 'main', 'manual'),
(1103, 'zvučnici', 'sr', 7.0, 'context', 'manual'),
(1103, 'zvuk', 'sr', 6.0, 'context', 'manual');

-- Kućni aparati (1104)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1104, 'home appliances', 'en', 10.0, 'main', 'manual'),
(1104, 'appliances', 'en', 9.0, 'synonym', 'manual'),
(1104, 'kitchen', 'en', 8.0, 'context', 'manual'),
(1104, 'refrigerator', 'en', 7.0, 'context', 'manual'),
(1104, 'washing machine', 'en', 6.0, 'context', 'manual'),
-- Russian
(1104, 'бытовая техника', 'ru', 10.0, 'main', 'manual'),
(1104, 'техника', 'ru', 9.0, 'synonym', 'manual'),
(1104, 'кухня', 'ru', 8.0, 'context', 'manual'),
(1104, 'холодильник', 'ru', 7.0, 'context', 'manual'),
(1104, 'стиральная машина', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1104, 'kućni aparati', 'sr', 10.0, 'main', 'manual'),
(1104, 'aparati', 'sr', 9.0, 'synonym', 'manual'),
(1104, 'kuhinja', 'sr', 8.0, 'context', 'manual'),
(1104, 'frižider', 'sr', 7.0, 'context', 'manual'),
(1104, 'veš mašina', 'sr', 6.0, 'context', 'manual');

-- Gaming Consoles (1105)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1105, 'gaming consoles', 'en', 10.0, 'main', 'manual'),
(1105, 'playstation', 'en', 9.0, 'brand', 'manual'),
(1105, 'xbox', 'en', 8.0, 'brand', 'manual'),
(1105, 'nintendo', 'en', 7.0, 'brand', 'manual'),
(1105, 'games', 'en', 6.0, 'context', 'manual'),
-- Russian
(1105, 'игровые консоли', 'ru', 10.0, 'main', 'manual'),
(1105, 'приставки', 'ru', 9.0, 'synonym', 'manual'),
(1105, 'плейстейшн', 'ru', 8.0, 'brand', 'manual'),
(1105, 'иксбокс', 'ru', 7.0, 'brand', 'manual'),
(1105, 'игры', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1105, 'gejming konzole', 'sr', 10.0, 'main', 'manual'),
(1105, 'konzole', 'sr', 9.0, 'synonym', 'manual'),
(1105, 'plejstejšn', 'sr', 8.0, 'brand', 'manual'),
(1105, 'iksboks', 'sr', 7.0, 'brand', 'manual'),
(1105, 'igre', 'sr', 6.0, 'context', 'manual');

-- Photo & Video (1106)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1106, 'photo', 'en', 10.0, 'main', 'manual'),
(1106, 'video', 'en', 9.0, 'main', 'manual'),
(1106, 'camera', 'en', 8.0, 'synonym', 'manual'),
(1106, 'photography', 'en', 7.0, 'context', 'manual'),
(1106, 'lens', 'en', 6.0, 'context', 'manual'),
-- Russian
(1106, 'фото', 'ru', 10.0, 'main', 'manual'),
(1106, 'видео', 'ru', 9.0, 'main', 'manual'),
(1106, 'камера', 'ru', 8.0, 'synonym', 'manual'),
(1106, 'фотография', 'ru', 7.0, 'context', 'manual'),
(1106, 'объектив', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1106, 'foto', 'sr', 10.0, 'main', 'manual'),
(1106, 'video', 'sr', 9.0, 'main', 'manual'),
(1106, 'kamera', 'sr', 8.0, 'synonym', 'manual'),
(1106, 'fotografija', 'sr', 7.0, 'context', 'manual'),
(1106, 'objektiv', 'sr', 6.0, 'context', 'manual');

-- Smart Home (1107)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1107, 'smart home', 'en', 10.0, 'main', 'manual'),
(1107, 'automation', 'en', 9.0, 'synonym', 'manual'),
(1107, 'iot', 'en', 8.0, 'context', 'manual'),
(1107, 'alexa', 'en', 7.0, 'brand', 'manual'),
(1107, 'google home', 'en', 6.0, 'brand', 'manual'),
-- Russian
(1107, 'умный дом', 'ru', 10.0, 'main', 'manual'),
(1107, 'автоматизация', 'ru', 9.0, 'synonym', 'manual'),
(1107, 'смарт', 'ru', 8.0, 'context', 'manual'),
(1107, 'алекса', 'ru', 7.0, 'brand', 'manual'),
(1107, 'гугл хоум', 'ru', 6.0, 'brand', 'manual'),
-- Serbian
(1107, 'pametna kuća', 'sr', 10.0, 'main', 'manual'),
(1107, 'automatizacija', 'sr', 9.0, 'synonym', 'manual'),
(1107, 'smart', 'sr', 8.0, 'context', 'manual'),
(1107, 'aleksa', 'sr', 7.0, 'brand', 'manual'),
(1107, 'gugl houm', 'sr', 6.0, 'brand', 'manual');

-- Electronics Accessories (1108)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1108, 'accessories', 'en', 10.0, 'main', 'manual'),
(1108, 'cables', 'en', 9.0, 'context', 'manual'),
(1108, 'chargers', 'en', 8.0, 'context', 'manual'),
(1108, 'cases', 'en', 7.0, 'context', 'manual'),
(1108, 'adapters', 'en', 6.0, 'context', 'manual'),
-- Russian
(1108, 'аксессуары', 'ru', 10.0, 'main', 'manual'),
(1108, 'кабели', 'ru', 9.0, 'context', 'manual'),
(1108, 'зарядки', 'ru', 8.0, 'context', 'manual'),
(1108, 'чехлы', 'ru', 7.0, 'context', 'manual'),
(1108, 'адаптеры', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1108, 'aksesoari', 'sr', 10.0, 'main', 'manual'),
(1108, 'kablovi', 'sr', 9.0, 'context', 'manual'),
(1108, 'punjači', 'sr', 8.0, 'context', 'manual'),
(1108, 'futrole', 'sr', 7.0, 'context', 'manual'),
(1108, 'adapteri', 'sr', 6.0, 'context', 'manual');

-- Ženska odeća (1202)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1202, 'womens clothing', 'en', 10.0, 'main', 'manual'),
(1202, 'dress', 'en', 9.0, 'context', 'manual'),
(1202, 'skirt', 'en', 8.0, 'context', 'manual'),
(1202, 'blouse', 'en', 7.0, 'context', 'manual'),
(1202, 'fashion', 'en', 6.0, 'context', 'manual'),
-- Russian
(1202, 'женская одежда', 'ru', 10.0, 'main', 'manual'),
(1202, 'платье', 'ru', 9.0, 'context', 'manual'),
(1202, 'юбка', 'ru', 8.0, 'context', 'manual'),
(1202, 'блузка', 'ru', 7.0, 'context', 'manual'),
(1202, 'мода', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1202, 'ženska odeća', 'sr', 10.0, 'main', 'manual'),
(1202, 'haljina', 'sr', 9.0, 'context', 'manual'),
(1202, 'suknja', 'sr', 8.0, 'context', 'manual'),
(1202, 'bluza', 'sr', 7.0, 'context', 'manual'),
(1202, 'moda', 'sr', 6.0, 'context', 'manual');

-- Watches (1207)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1207, 'watches', 'en', 10.0, 'main', 'manual'),
(1207, 'watch', 'en', 9.0, 'synonym', 'manual'),
(1207, 'smartwatch', 'en', 8.0, 'context', 'manual'),
(1207, 'rolex', 'en', 7.0, 'brand', 'manual'),
(1207, 'timepiece', 'en', 6.0, 'synonym', 'manual'),
-- Russian
(1207, 'часы', 'ru', 10.0, 'main', 'manual'),
(1207, 'наручные часы', 'ru', 9.0, 'synonym', 'manual'),
(1207, 'смарт часы', 'ru', 8.0, 'context', 'manual'),
(1207, 'ролекс', 'ru', 7.0, 'brand', 'manual'),
(1207, 'хронометр', 'ru', 6.0, 'synonym', 'manual'),
-- Serbian
(1207, 'satovi', 'sr', 10.0, 'main', 'manual'),
(1207, 'sat', 'sr', 9.0, 'synonym', 'manual'),
(1207, 'pametni sat', 'sr', 8.0, 'context', 'manual'),
(1207, 'roleks', 'sr', 7.0, 'brand', 'manual'),
(1207, 'časovnik', 'sr', 6.0, 'synonym', 'manual');

-- Добавляем ключевые слова для остальных подкатегорий автомобилей

-- Motocikli (1302)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1302, 'motorcycles', 'en', 10.0, 'main', 'manual'),
(1302, 'bikes', 'en', 9.0, 'synonym', 'manual'),
(1302, 'motorbike', 'en', 8.0, 'synonym', 'manual'),
(1302, 'scooter', 'en', 7.0, 'context', 'manual'),
(1302, 'harley', 'en', 6.0, 'brand', 'manual'),
-- Russian
(1302, 'мотоциклы', 'ru', 10.0, 'main', 'manual'),
(1302, 'мотоцикл', 'ru', 9.0, 'synonym', 'manual'),
(1302, 'байк', 'ru', 8.0, 'synonym', 'manual'),
(1302, 'скутер', 'ru', 7.0, 'context', 'manual'),
(1302, 'харлей', 'ru', 6.0, 'brand', 'manual'),
-- Serbian
(1302, 'motocikli', 'sr', 10.0, 'main', 'manual'),
(1302, 'motor', 'sr', 9.0, 'synonym', 'manual'),
(1302, 'bajk', 'sr', 8.0, 'synonym', 'manual'),
(1302, 'skuter', 'sr', 7.0, 'context', 'manual'),
(1302, 'harlej', 'sr', 6.0, 'brand', 'manual');

-- Auto delovi (1303)
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- English
(1303, 'auto parts', 'en', 10.0, 'main', 'manual'),
(1303, 'car parts', 'en', 9.0, 'synonym', 'manual'),
(1303, 'spare parts', 'en', 8.0, 'synonym', 'manual'),
(1303, 'engine', 'en', 7.0, 'context', 'manual'),
(1303, 'tires', 'en', 6.0, 'context', 'manual'),
-- Russian
(1303, 'автозапчасти', 'ru', 10.0, 'main', 'manual'),
(1303, 'запчасти', 'ru', 9.0, 'synonym', 'manual'),
(1303, 'детали', 'ru', 8.0, 'synonym', 'manual'),
(1303, 'двигатель', 'ru', 7.0, 'context', 'manual'),
(1303, 'шины', 'ru', 6.0, 'context', 'manual'),
-- Serbian
(1303, 'auto delovi', 'sr', 10.0, 'main', 'manual'),
(1303, 'delovi', 'sr', 9.0, 'synonym', 'manual'),
(1303, 'rezervni delovi', 'sr', 8.0, 'synonym', 'manual'),
(1303, 'motor', 'sr', 7.0, 'context', 'manual'),
(1303, 'gume', 'sr', 6.0, 'context', 'manual');

-- Update timestamp
UPDATE category_keywords SET updated_at = CURRENT_TIMESTAMP WHERE source = 'manual';

-- Добавляем статистику успешности для основных категорий
UPDATE category_keywords 
SET usage_count = 10, success_rate = 0.9
WHERE category_id IN (1001, 1002, 1003, 1004, 1005) 
  AND keyword_type = 'main';