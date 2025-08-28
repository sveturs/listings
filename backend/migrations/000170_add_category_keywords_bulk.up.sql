-- Добавление ключевых слов для 20 самых популярных категорий

-- AUTOMOTIVE
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1003, 'автомобиль', 'ru', 'main', 10.0, 'manual'),
(1003, 'car', 'en', 'main', 10.0, 'manual'),
(1003, 'automobil', 'sr', 'main', 10.0, 'manual'),
(1003, 'машина', 'ru', 'main', 9.0, 'manual'),
(1003, 'vehicle', 'en', 'main', 9.0, 'manual'),
(1003, 'vozilo', 'sr', 'main', 9.0, 'manual'),
(1003, 'авто', 'ru', 'synonym', 8.0, 'manual'),
(1003, 'automobile', 'en', 'synonym', 8.0, 'manual'),
(1003, 'kola', 'sr', 'synonym', 8.0, 'manual'),
(1003, 'toyota', '', 'brand', 5.0, 'manual'),
(1003, 'bmw', '', 'brand', 5.0, 'manual'),
(1003, 'mercedes', '', 'brand', 5.0, 'manual'),
(1003, 'audi', '', 'brand', 5.0, 'manual'),
(1003, 'volkswagen', '', 'brand', 5.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- ELECTRONICS
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1001, 'электроника', 'ru', 'main', 10.0, 'manual'),
(1001, 'electronics', 'en', 'main', 10.0, 'manual'),
(1001, 'elektronika', 'sr', 'main', 10.0, 'manual'),
(1001, 'техника', 'ru', 'main', 9.0, 'manual'),
(1001, 'technology', 'en', 'main', 9.0, 'manual'),
(1001, 'гаджеты', 'ru', 'synonym', 8.0, 'manual'),
(1001, 'gadgets', 'en', 'synonym', 8.0, 'manual'),
(1001, 'samsung', '', 'brand', 5.0, 'manual'),
(1001, 'apple', '', 'brand', 5.0, 'manual'),
(1001, 'sony', '', 'brand', 5.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- SMARTPHONES
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1101, 'смартфон', 'ru', 'main', 10.0, 'manual'),
(1101, 'smartphone', 'en', 'main', 10.0, 'manual'),
(1101, 'pametni telefon', 'sr', 'main', 10.0, 'manual'),
(1101, 'телефон', 'ru', 'main', 9.0, 'manual'),
(1101, 'phone', 'en', 'main', 9.0, 'manual'),
(1101, 'мобильный', 'ru', 'synonym', 8.0, 'manual'),
(1101, 'mobile', 'en', 'synonym', 8.0, 'manual'),
(1101, 'iphone', '', 'brand', 6.0, 'manual'),
(1101, 'samsung', '', 'brand', 6.0, 'manual'),
(1101, 'xiaomi', '', 'brand', 6.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- COMPUTERS
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1102, 'компьютер', 'ru', 'main', 10.0, 'manual'),
(1102, 'computer', 'en', 'main', 10.0, 'manual'),
(1102, 'računar', 'sr', 'main', 10.0, 'manual'),
(1102, 'ноутбук', 'ru', 'main', 9.0, 'manual'),
(1102, 'laptop', 'en', 'main', 9.0, 'manual'),
(1102, 'пк', 'ru', 'synonym', 8.0, 'manual'),
(1102, 'pc', 'en', 'synonym', 8.0, 'manual'),
(1102, 'dell', '', 'brand', 5.0, 'manual'),
(1102, 'hp', '', 'brand', 5.0, 'manual'),
(1102, 'lenovo', '', 'brand', 5.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- TIRES-AND-WHEELS
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1304, 'шины', 'ru', 'main', 10.0, 'manual'),
(1304, 'tires', 'en', 'main', 10.0, 'manual'),
(1304, 'gume', 'sr', 'main', 10.0, 'manual'),
(1304, 'колеса', 'ru', 'main', 9.0, 'manual'),
(1304, 'wheels', 'en', 'main', 9.0, 'manual'),
(1304, 'покрышки', 'ru', 'synonym', 8.0, 'manual'),
(1304, 'michelin', '', 'brand', 5.0, 'manual'),
(1304, 'bridgestone', '', 'brand', 5.0, 'manual'),
(1304, 'continental', '', 'brand', 5.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- MENS-CLOTHING
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1201, 'мужская одежда', 'ru', 'main', 10.0, 'manual'),
(1201, 'mens clothing', 'en', 'main', 10.0, 'manual'),
(1201, 'muška odeća', 'sr', 'main', 10.0, 'manual'),
(1201, 'одежда для мужчин', 'ru', 'main', 9.0, 'manual'),
(1201, 'menswear', 'en', 'synonym', 8.0, 'manual'),
(1201, 'костюм', 'ru', 'synonym', 7.0, 'manual'),
(1201, 'suit', 'en', 'synonym', 7.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- INDUSTRIAL
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1007, 'промышленность', 'ru', 'main', 10.0, 'manual'),
(1007, 'industrial', 'en', 'main', 10.0, 'manual'),
(1007, 'industrija', 'sr', 'main', 10.0, 'manual'),
(1007, 'оборудование', 'ru', 'main', 9.0, 'manual'),
(1007, 'equipment', 'en', 'main', 9.0, 'manual'),
(1007, 'производство', 'ru', 'synonym', 8.0, 'manual'),
(1007, 'manufacturing', 'en', 'synonym', 8.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- FARM-MACHINERY
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1601, 'сельхозтехника', 'ru', 'main', 10.0, 'manual'),
(1601, 'farm machinery', 'en', 'main', 10.0, 'manual'),
(1601, 'poljoprivredne mašine', 'sr', 'main', 10.0, 'manual'),
(1601, 'трактор', 'ru', 'main', 9.0, 'manual'),
(1601, 'tractor', 'en', 'main', 9.0, 'manual'),
(1601, 'комбайн', 'ru', 'synonym', 8.0, 'manual'),
(1601, 'john deere', '', 'brand', 5.0, 'manual'),
(1601, 'case', '', 'brand', 5.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- FOOD-BEVERAGES
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1008, 'еда', 'ru', 'main', 10.0, 'manual'),
(1008, 'food', 'en', 'main', 10.0, 'manual'),
(1008, 'hrana', 'sr', 'main', 10.0, 'manual'),
(1008, 'напитки', 'ru', 'main', 9.0, 'manual'),
(1008, 'beverages', 'en', 'main', 9.0, 'manual'),
(1008, 'продукты', 'ru', 'synonym', 8.0, 'manual'),
(1008, 'groceries', 'en', 'synonym', 8.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- ORGANIC-FOOD
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1801, 'органическая еда', 'ru', 'main', 10.0, 'manual'),
(1801, 'organic food', 'en', 'main', 10.0, 'manual'),
(1801, 'organska hrana', 'sr', 'main', 10.0, 'manual'),
(1801, 'эко продукты', 'ru', 'main', 9.0, 'manual'),
(1801, 'биопродукты', 'ru', 'synonym', 8.0, 'manual'),
(1801, 'натуральная еда', 'ru', 'synonym', 7.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- ELECTRONICS-ACCESSORIES
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1108, 'аксессуары для электроники', 'ru', 'main', 10.0, 'manual'),
(1108, 'electronics accessories', 'en', 'main', 10.0, 'manual'),
(1108, 'зарядное устройство', 'ru', 'main', 9.0, 'manual'),
(1108, 'charger', 'en', 'main', 9.0, 'manual'),
(1108, 'кабель', 'ru', 'synonym', 8.0, 'manual'),
(1108, 'cable', 'en', 'synonym', 8.0, 'manual'),
(1108, 'чехол', 'ru', 'synonym', 7.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- PLUMBING
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1508, 'сантехника', 'ru', 'main', 10.0, 'manual'),
(1508, 'plumbing', 'en', 'main', 10.0, 'manual'),
(1508, 'водопровод', 'ru', 'main', 9.0, 'manual'),
(1508, 'смеситель', 'ru', 'synonym', 8.0, 'manual'),
(1508, 'faucet', 'en', 'synonym', 8.0, 'manual'),
(1508, 'труба', 'ru', 'synonym', 7.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- KIDS-CLOTHING  
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1205, 'детская одежда', 'ru', 'main', 10.0, 'manual'),
(1205, 'kids clothing', 'en', 'main', 10.0, 'manual'),
(1205, 'dečja odeća', 'sr', 'main', 10.0, 'manual'),
(1205, 'одежда для детей', 'ru', 'main', 9.0, 'manual'),
(1205, 'детская мода', 'ru', 'synonym', 8.0, 'manual'),
(1205, 'школьная форма', 'ru', 'synonym', 7.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- EVENTS-TICKETS
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES  
(1020, 'билеты', 'ru', 'main', 10.0, 'manual'),
(1020, 'tickets', 'en', 'main', 10.0, 'manual'),
(1020, 'karte', 'sr', 'main', 10.0, 'manual'),
(1020, 'мероприятия', 'ru', 'main', 9.0, 'manual'),
(1020, 'events', 'en', 'main', 9.0, 'manual'),
(1020, 'концерт', 'ru', 'synonym', 8.0, 'manual'),
(1020, 'concert', 'en', 'synonym', 8.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- BAGS
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1208, 'сумки', 'ru', 'main', 10.0, 'manual'),
(1208, 'bags', 'en', 'main', 10.0, 'manual'),
(1208, 'torbe', 'sr', 'main', 10.0, 'manual'),
(1208, 'рюкзак', 'ru', 'main', 9.0, 'manual'),
(1208, 'backpack', 'en', 'main', 9.0, 'manual'),
(1208, 'портфель', 'ru', 'synonym', 8.0, 'manual'),
(1208, 'кошелек', 'ru', 'synonym', 7.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- WINTER-TIRES
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(1315, 'зимние шины', 'ru', 'main', 10.0, 'manual'),
(1315, 'winter tires', 'en', 'main', 10.0, 'manual'),
(1315, 'zimske gume', 'sr', 'main', 10.0, 'manual'),
(1315, 'зимняя резина', 'ru', 'main', 9.0, 'manual'),
(1315, 'снежные шины', 'ru', 'synonym', 8.0, 'manual'),
(1315, 'шипованные', 'ru', 'synonym', 7.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;

-- HUNTING-FISHING
INSERT INTO category_keywords (category_id, keyword, language, keyword_type, weight, source) VALUES
(2013, 'охота', 'ru', 'main', 10.0, 'manual'),
(2013, 'hunting', 'en', 'main', 10.0, 'manual'),
(2013, 'lov', 'sr', 'main', 10.0, 'manual'),
(2013, 'рыбалка', 'ru', 'main', 10.0, 'manual'),
(2013, 'fishing', 'en', 'main', 10.0, 'manual'),
(2013, 'удочка', 'ru', 'synonym', 8.0, 'manual'),
(2013, 'ружье', 'ru', 'synonym', 7.0, 'manual')
ON CONFLICT (category_id, keyword, language) DO NOTHING;