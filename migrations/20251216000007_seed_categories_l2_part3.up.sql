-- Migration: Seed L2 categories (Part 3 of 3 - Final)
-- Date: 2025-12-16
-- Purpose: Complete L2 subcategories for remaining 12 L1 categories
-- Previous: 20251216000006_seed_categories_l2_part2.up.sql

-- =============================================================================
-- L2 for: 7. Automobilizam (Automotive) - 12 categories
-- =============================================================================
INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('delovi-za-automobile', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/delovi-za-automobile', 1,
 '{"sr": "Delovi za automobile", "en": "Auto parts", "ru": "Автозапчасти"}'::jsonb,
 '{"sr": "Motori, kočnice, filteri, svetla", "en": "Engines, brakes, filters, lights", "ru": "Двигатели, тормоза, фильтры, фары"}'::jsonb,
 '{"sr": "Delovi za automobile | Vondi", "en": "Auto parts | Vondi", "ru": "Автозапчасти | Vondi"}'::jsonb,
 '{"sr": "Kupite delove za automobile online", "en": "Buy auto parts online", "ru": "Купить автозапчасти онлайн"}'::jsonb,
 '🔧', true),

('gume-i-felne', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/gume-i-felne', 2,
 '{"sr": "Gume i felne", "en": "Tires & rims", "ru": "Шины и диски"}'::jsonb,
 '{"sr": "Zimske, letnje gume, aluminijumske felne", "en": "Winter, summer tires, alloy rims", "ru": "Зимние, летние шины, литые диски"}'::jsonb,
 '{"sr": "Gume i felne | Vondi", "en": "Tires & rims | Vondi", "ru": "Шины и диски | Vondi"}'::jsonb,
 '{"sr": "Kupite gume i felne online", "en": "Buy tires and rims online", "ru": "Купить шины и диски онлайн"}'::jsonb,
 '🛞', true),

('auto-kozmetika', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/auto-kozmetika', 3,
 '{"sr": "Auto kozmetika", "en": "Car care products", "ru": "Автокосметика"}'::jsonb,
 '{"sr": "Šamponi, voskovi, sredstva za poliranje", "en": "Shampoos, waxes, polishing products", "ru": "Шампуни, воски, средства для полировки"}'::jsonb,
 '{"sr": "Auto kozmetika | Vondi", "en": "Car care products | Vondi", "ru": "Автокосметика | Vondi"}'::jsonb,
 '{"sr": "Kupite auto kozmetiku online", "en": "Buy car care products online", "ru": "Купить автокосметику онлайн"}'::jsonb,
 '🧽', true),

('audio-i-navigacija', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/audio-i-navigacija', 4,
 '{"sr": "Audio i navigacija", "en": "Audio & navigation", "ru": "Аудио и навигация"}'::jsonb,
 '{"sr": "Auto radio, zvučnici, GPS navigacija", "en": "Car radios, speakers, GPS navigation", "ru": "Автомагнитолы, колонки, GPS навигация"}'::jsonb,
 '{"sr": "Audio i navigacija | Vondi", "en": "Audio & navigation | Vondi", "ru": "Аудио и навигация | Vondi"}'::jsonb,
 '{"sr": "Kupite auto audio i GPS online", "en": "Buy car audio and GPS online", "ru": "Купить автомагнитолы и GPS онлайн"}'::jsonb,
 '📻', true),

('auto-dodaci', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/auto-dodaci', 5,
 '{"sr": "Auto dodaci", "en": "Car accessories", "ru": "Автоаксессуары"}'::jsonb,
 '{"sr": "Držači telefona, punjači, osveživači", "en": "Phone holders, chargers, air fresheners", "ru": "Держатели телефонов, зарядки, освежители"}'::jsonb,
 '{"sr": "Auto dodaci | Vondi", "en": "Car accessories | Vondi", "ru": "Автоаксессуары | Vondi"}'::jsonb,
 '{"sr": "Kupite auto dodatke online", "en": "Buy car accessories online", "ru": "Купить автоаксессуары онлайн"}'::jsonb,
 '🚗', true),

('moto-oprema', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/moto-oprema', 6,
 '{"sr": "Moto oprema", "en": "Motorcycle gear", "ru": "Мотоснаряжение"}'::jsonb,
 '{"sr": "Kacige, jakne, rukavice, čizme", "en": "Helmets, jackets, gloves, boots", "ru": "Шлемы, куртки, перчатки, ботинки"}'::jsonb,
 '{"sr": "Moto oprema | Vondi", "en": "Motorcycle gear | Vondi", "ru": "Мотоснаряжение | Vondi"}'::jsonb,
 '{"sr": "Kupite moto opremu online", "en": "Buy motorcycle gear online", "ru": "Купить мотоснаряжение онлайн"}'::jsonb,
 '🏍️', true),

('delovi-za-motocikle', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/delovi-za-motocikle', 7,
 '{"sr": "Delovi za motocikle", "en": "Motorcycle parts", "ru": "Запчасти для мотоциклов"}'::jsonb,
 '{"sr": "Delovi motora, kočnice, izduvni sistem", "en": "Engine parts, brakes, exhaust systems", "ru": "Детали двигателя, тормоза, выхлопные системы"}'::jsonb,
 '{"sr": "Delovi za motocikle | Vondi", "en": "Motorcycle parts | Vondi", "ru": "Запчасти для мотоциклов | Vondi"}'::jsonb,
 '{"sr": "Kupite delove za motocikle online", "en": "Buy motorcycle parts online", "ru": "Купить запчасти для мотоциклов онлайн"}'::jsonb,
 '⚙️', true),

('alati-za-automobile', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/alati-za-automobile', 8,
 '{"sr": "Alati za automobile", "en": "Auto tools", "ru": "Автоинструменты"}'::jsonb,
 '{"sr": "Ključevi, dizalice, kompresori", "en": "Wrenches, jacks, compressors", "ru": "Ключи, домкраты, компрессоры"}'::jsonb,
 '{"sr": "Alati za automobile | Vondi", "en": "Auto tools | Vondi", "ru": "Автоинструменты | Vondi"}'::jsonb,
 '{"sr": "Kupite auto alate online", "en": "Buy auto tools online", "ru": "Купить автоинструменты онлайн"}'::jsonb,
 '🔨', true),

('tuniranje', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/tuniranje', 9,
 '{"sr": "Tuniranje", "en": "Tuning", "ru": "Тюнинг"}'::jsonb,
 '{"sr": "Sportski izduvni sistemi, chip tuning", "en": "Sports exhaust systems, chip tuning", "ru": "Спортивные выхлопные системы, чип-тюнинг"}'::jsonb,
 '{"sr": "Tuniranje | Vondi", "en": "Tuning | Vondi", "ru": "Тюнинг | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za tuniranje online", "en": "Buy tuning equipment online", "ru": "Купить оборудование для тюнинга онлайн"}'::jsonb,
 '🏎️', true),

('dash-kamere', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/dash-kamere', 10,
 '{"sr": "Dash kamere", "en": "Dash cameras", "ru": "Видеорегистраторы"}'::jsonb,
 '{"sr": "Video rekorderi za vozila", "en": "Video recorders for vehicles", "ru": "Видеорегистраторы для автомобилей"}'::jsonb,
 '{"sr": "Dash kamere | Vondi", "en": "Dash cameras | Vondi", "ru": "Видеорегистраторы | Vondi"}'::jsonb,
 '{"sr": "Kupite dash kamere online", "en": "Buy dash cameras online", "ru": "Купить видеорегистраторы онлайн"}'::jsonb,
 '📹', true),

('parking-senzori', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/parking-senzori', 11,
 '{"sr": "Parking senzori", "en": "Parking sensors", "ru": "Парктроники"}'::jsonb,
 '{"sr": "Parking senzori i kamere", "en": "Parking sensors and cameras", "ru": "Парктроники и камеры заднего вида"}'::jsonb,
 '{"sr": "Parking senzori | Vondi", "en": "Parking sensors | Vondi", "ru": "Парктроники | Vondi"}'::jsonb,
 '{"sr": "Kupite parking senzore online", "en": "Buy parking sensors online", "ru": "Купить парктроники онлайн"}'::jsonb,
 '📡', true),

('akumulatori', (SELECT id FROM categories WHERE slug = 'automobilizam'), 2, 'automobilizam/akumulatori', 12,
 '{"sr": "Akumulatori", "en": "Car batteries", "ru": "Аккумуляторы"}'::jsonb,
 '{"sr": "Auto akumulatori i punjači", "en": "Car batteries and chargers", "ru": "Автомобильные аккумуляторы и зарядные устройства"}'::jsonb,
 '{"sr": "Akumulatori | Vondi", "en": "Car batteries | Vondi", "ru": "Аккумуляторы | Vondi"}'::jsonb,
 '{"sr": "Kupite akumulatore online", "en": "Buy car batteries online", "ru": "Купить аккумуляторы онлайн"}'::jsonb,
 '🔋', true),

-- =============================================================================
-- L2 for: 8. Kućni aparati (Appliances) - 12 categories
-- =============================================================================

('frizideri', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/frizideri', 1,
 '{"sr": "Frižideri", "en": "Refrigerators", "ru": "Холодильники"}'::jsonb,
 '{"sr": "Frižideri sa zamrzivačem, side by side", "en": "Refrigerators with freezer, side by side", "ru": "Холодильники с морозильной камерой, side by side"}'::jsonb,
 '{"sr": "Frižideri | Vondi", "en": "Refrigerators | Vondi", "ru": "Холодильники | Vondi"}'::jsonb,
 '{"sr": "Kupite frižidere online", "en": "Buy refrigerators online", "ru": "Купить холодильники онлайн"}'::jsonb,
 '🧊', true),

('masine-za-pranje', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/masine-za-pranje', 2,
 '{"sr": "Mašine za pranje", "en": "Washing machines", "ru": "Стиральные машины"}'::jsonb,
 '{"sr": "Mašine za pranje veša, sušilice", "en": "Washing machines, dryers", "ru": "Стиральные машины, сушилки"}'::jsonb,
 '{"sr": "Mašine za pranje | Vondi", "en": "Washing machines | Vondi", "ru": "Стиральные машины | Vondi"}'::jsonb,
 '{"sr": "Kupite mašine za pranje online", "en": "Buy washing machines online", "ru": "Купить стиральные машины онлайн"}'::jsonb,
 '🧺', true),

('usisivaci', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/usisivaci', 3,
 '{"sr": "Usisivači", "en": "Vacuum cleaners", "ru": "Пылесосы"}'::jsonb,
 '{"sr": "Usisivači, robotski usisivači", "en": "Vacuum cleaners, robot vacuums", "ru": "Пылесосы, роботы-пылесосы"}'::jsonb,
 '{"sr": "Usisivači | Vondi", "en": "Vacuum cleaners | Vondi", "ru": "Пылесосы | Vondi"}'::jsonb,
 '{"sr": "Kupite usisivače online", "en": "Buy vacuum cleaners online", "ru": "Купить пылесосы онлайн"}'::jsonb,
 '🧹', true),

('sporet-i-rerna', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/sporet-i-rerna', 4,
 '{"sr": "Šporet i rerna", "en": "Stoves & ovens", "ru": "Плиты и духовки"}'::jsonb,
 '{"sr": "Električni i gasni šporeti, ugradbene rerne", "en": "Electric and gas stoves, built-in ovens", "ru": "Электрические и газовые плиты, встраиваемые духовки"}'::jsonb,
 '{"sr": "Šporet i rerna | Vondi", "en": "Stoves & ovens | Vondi", "ru": "Плиты и духовки | Vondi"}'::jsonb,
 '{"sr": "Kupite šporet i rerne online", "en": "Buy stoves and ovens online", "ru": "Купить плиты и духовки онлайн"}'::jsonb,
 '🍳', true),

('mikotalasne-rerne', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/mikotalasne-rerne', 5,
 '{"sr": "Mikrotalasne rerne", "en": "Microwave ovens", "ru": "Микроволновые печи"}'::jsonb,
 '{"sr": "Mikrotalasne rerne sa grilom i konvekcijom", "en": "Microwave ovens with grill and convection", "ru": "Микроволновые печи с грилем и конвекцией"}'::jsonb,
 '{"sr": "Mikrotalasne rerne | Vondi", "en": "Microwave ovens | Vondi", "ru": "Микроволновые печи | Vondi"}'::jsonb,
 '{"sr": "Kupite mikrotalasne rerne online", "en": "Buy microwave ovens online", "ru": "Купить микроволновые печи онлайн"}'::jsonb,
 '📻', true),

('sudopere-i-masine', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/sudopere-i-masine', 6,
 '{"sr": "Sudopere i mašine", "en": "Dishwashers", "ru": "Посудомоечные машины"}'::jsonb,
 '{"sr": "Mašine za pranje sudova, ugradbene i samostojeće", "en": "Dishwashers, built-in and freestanding", "ru": "Посудомоечные машины, встраиваемые и отдельностоящие"}'::jsonb,
 '{"sr": "Sudopere i mašine | Vondi", "en": "Dishwashers | Vondi", "ru": "Посудомоечные машины | Vondi"}'::jsonb,
 '{"sr": "Kupite mašine za pranje sudova online", "en": "Buy dishwashers online", "ru": "Купить посудомоечные машины онлайн"}'::jsonb,
 '🍽️', true),

('mali-kucni-aparati', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/mali-kucni-aparati', 7,
 '{"sr": "Mali kućni aparati", "en": "Small appliances", "ru": "Малая бытовая техника"}'::jsonb,
 '{"sr": "Blenderi, tosteri, kafe aparati, pegla", "en": "Blenders, toasters, coffee makers, irons", "ru": "Блендеры, тостеры, кофеварки, утюги"}'::jsonb,
 '{"sr": "Mali kućni aparati | Vondi", "en": "Small appliances | Vondi", "ru": "Малая бытовая техника | Vondi"}'::jsonb,
 '{"sr": "Kupite male kućne aparate online", "en": "Buy small appliances online", "ru": "Купить малую бытовую технику онлайн"}'::jsonb,
 '☕', true),

('bojleri', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/bojleri', 8,
 '{"sr": "Bojleri", "en": "Water heaters", "ru": "Водонагреватели"}'::jsonb,
 '{"sr": "Električni i gasni bojleri, protočni", "en": "Electric and gas water heaters, tankless", "ru": "Электрические и газовые водонагреватели, проточные"}'::jsonb,
 '{"sr": "Bojleri | Vondi", "en": "Water heaters | Vondi", "ru": "Водонагреватели | Vondi"}'::jsonb,
 '{"sr": "Kupite bojlere online", "en": "Buy water heaters online", "ru": "Купить водонагреватели онлайн"}'::jsonb,
 '🚿', true),

('ventilatori-i-grejalice', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/ventilatori-i-grejalice', 9,
 '{"sr": "Ventilatori i grejalice", "en": "Fans & heaters", "ru": "Вентиляторы и обогреватели"}'::jsonb,
 '{"sr": "Ventilatori, grejalice, klime", "en": "Fans, heaters, air conditioners", "ru": "Вентиляторы, обогреватели, кондиционеры"}'::jsonb,
 '{"sr": "Ventilatori i grejalice | Vondi", "en": "Fans & heaters | Vondi", "ru": "Вентиляторы и обогреватели | Vondi"}'::jsonb,
 '{"sr": "Kupite ventilatoren i grejalice online", "en": "Buy fans and heaters online", "ru": "Купить вентиляторы и обогреватели онлайн"}'::jsonb,
 '🌀', true),

('precistaci-vazduha', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/precistaci-vazduha', 10,
 '{"sr": "Prečistači vazduha", "en": "Air purifiers", "ru": "Очистители воздуха"}'::jsonb,
 '{"sr": "Prečistači vazduha, ovlaživači, odvlaživači", "en": "Air purifiers, humidifiers, dehumidifiers", "ru": "Очистители воздуха, увлажнители, осушители"}'::jsonb,
 '{"sr": "Prečistači vazduha | Vondi", "en": "Air purifiers | Vondi", "ru": "Очистители воздуха | Vondi"}'::jsonb,
 '{"sr": "Kupite prečistače vazduha online", "en": "Buy air purifiers online", "ru": "Купить очистители воздуха онлайн"}'::jsonb,
 '💨', true),

('friteze', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/friteze', 11,
 '{"sr": "Friteze", "en": "Fryers", "ru": "Фритюрницы"}'::jsonb,
 '{"sr": "Air fryer, električne friteze", "en": "Air fryers, electric fryers", "ru": "Аэрогрили, электрофритюрницы"}'::jsonb,
 '{"sr": "Friteze | Vondi", "en": "Fryers | Vondi", "ru": "Фритюрницы | Vondi"}'::jsonb,
 '{"sr": "Kupite friteze online", "en": "Buy fryers online", "ru": "Купить фритюрницы онлайн"}'::jsonb,
 '🍟', true),

('masine-za-kafu', (SELECT id FROM categories WHERE slug = 'kucni-aparati'), 2, 'kucni-aparati/masine-za-kafu', 12,
 '{"sr": "Mašine za kafu", "en": "Coffee machines", "ru": "Кофемашины"}'::jsonb,
 '{"sr": "Espresso mašine, kafe aparati, kapsulne", "en": "Espresso machines, coffee makers, capsule", "ru": "Эспрессо-машины, кофеварки, капсульные"}'::jsonb,
 '{"sr": "Mašine za kafu | Vondi", "en": "Coffee machines | Vondi", "ru": "Кофемашины | Vondi"}'::jsonb,
 '{"sr": "Kupite mašine za kafu online", "en": "Buy coffee machines online", "ru": "Купить кофемашины онлайн"}'::jsonb,
 '☕', true),

-- =============================================================================
-- L2 for: 9. Nakit i satovi (Jewelry & Watches) - 10 categories
-- =============================================================================

('zlatni-nakit', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/zlatni-nakit', 1,
 '{"sr": "Zlatni nakit", "en": "Gold jewelry", "ru": "Золотые украшения"}'::jsonb,
 '{"sr": "Zlatne ogrlice, narukvice, prstenje", "en": "Gold necklaces, bracelets, rings", "ru": "Золотые ожерелья, браслеты, кольца"}'::jsonb,
 '{"sr": "Zlatni nakit | Vondi", "en": "Gold jewelry | Vondi", "ru": "Золотые украшения | Vondi"}'::jsonb,
 '{"sr": "Kupite zlatni nakit online", "en": "Buy gold jewelry online", "ru": "Купить золотые украшения онлайн"}'::jsonb,
 '💍', true),

('srebrni-nakit', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/srebrni-nakit', 2,
 '{"sr": "Srebrni nakit", "en": "Silver jewelry", "ru": "Серебряные украшения"}'::jsonb,
 '{"sr": "Srebrne ogrlice, narukvice, minđuše", "en": "Silver necklaces, bracelets, earrings", "ru": "Серебряные ожерелья, браслеты, серьги"}'::jsonb,
 '{"sr": "Srebrni nakit | Vondi", "en": "Silver jewelry | Vondi", "ru": "Серебряные украшения | Vondi"}'::jsonb,
 '{"sr": "Kupite srebrni nakit online", "en": "Buy silver jewelry online", "ru": "Купить серебряные украшения онлайн"}'::jsonb,
 '🪙', true),

('muski-satovi', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/muski-satovi', 3,
 '{"sr": "Muški satovi", "en": "Men''s watches", "ru": "Мужские часы"}'::jsonb,
 '{"sr": "Ručni satovi za muškarce, sportski, elegantni", "en": "Wristwatches for men, sports, elegant", "ru": "Наручные часы для мужчин, спортивные, элегантные"}'::jsonb,
 '{"sr": "Muški satovi | Vondi", "en": "Men''s watches | Vondi", "ru": "Мужские часы | Vondi"}'::jsonb,
 '{"sr": "Kupite muške satove online", "en": "Buy men''s watches online", "ru": "Купить мужские часы онлайн"}'::jsonb,
 '⌚', true),

('zenski-satovi', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/zenski-satovi', 4,
 '{"sr": "Ženski satovi", "en": "Women''s watches", "ru": "Женские часы"}'::jsonb,
 '{"sr": "Ručni satovi za žene, elegantni, casual", "en": "Wristwatches for women, elegant, casual", "ru": "Наручные часы для женщин, элегантные, повседневные"}'::jsonb,
 '{"sr": "Ženski satovi | Vondi", "en": "Women''s watches | Vondi", "ru": "Женские часы | Vondi"}'::jsonb,
 '{"sr": "Kupite ženske satove online", "en": "Buy women''s watches online", "ru": "Купить женские часы онлайн"}'::jsonb,
 '⌚', true),

('biserni-nakit', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/biserni-nakit', 5,
 '{"sr": "Biserni nakit", "en": "Pearl jewelry", "ru": "Жемчужные украшения"}'::jsonb,
 '{"sr": "Biserne ogrlice, narukvice, minđuše", "en": "Pearl necklaces, bracelets, earrings", "ru": "Жемчужные ожерелья, браслеты, серьги"}'::jsonb,
 '{"sr": "Biserni nakit | Vondi", "en": "Pearl jewelry | Vondi", "ru": "Жемчужные украшения | Vondi"}'::jsonb,
 '{"sr": "Kupite biserni nakit online", "en": "Buy pearl jewelry online", "ru": "Купить жемчужные украшения онлайн"}'::jsonb,
 '📿', true),

('verenicko-prstenje', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/verenicko-prstenje', 6,
 '{"sr": "Vereničko prstenje", "en": "Engagement rings", "ru": "Обручальные кольца"}'::jsonb,
 '{"sr": "Vereničko i burme sa dijamantima", "en": "Engagement and wedding rings with diamonds", "ru": "Помолвочные и обручальные кольца с бриллиантами"}'::jsonb,
 '{"sr": "Vereničko prstenje | Vondi", "en": "Engagement rings | Vondi", "ru": "Обручальные кольца | Vondi"}'::jsonb,
 '{"sr": "Kupite vereničko prstenje online", "en": "Buy engagement rings online", "ru": "Купить обручальные кольца онлайн"}'::jsonb,
 '💎', true),

('luksuzni-satovi', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/luksuzni-satovi', 7,
 '{"sr": "Luksuzni satovi", "en": "Luxury watches", "ru": "Люксовые часы"}'::jsonb,
 '{"sr": "Rolex, Omega, TAG Heuer premium satovi", "en": "Rolex, Omega, TAG Heuer premium watches", "ru": "Rolex, Omega, TAG Heuer премиум часы"}'::jsonb,
 '{"sr": "Luksuzni satovi | Vondi", "en": "Luxury watches | Vondi", "ru": "Люксовые часы | Vondi"}'::jsonb,
 '{"sr": "Kupite luksuzne satove online", "en": "Buy luxury watches online", "ru": "Купить люксовые часы онлайн"}'::jsonb,
 '👑', true),

('dijamanti', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/dijamanti', 8,
 '{"sr": "Dijamanti", "en": "Diamonds", "ru": "Бриллианты"}'::jsonb,
 '{"sr": "Dijamantski nakit, sertifikovani dijamanti", "en": "Diamond jewelry, certified diamonds", "ru": "Бриллиантовые украшения, сертифицированные бриллианты"}'::jsonb,
 '{"sr": "Dijamanti | Vondi", "en": "Diamonds | Vondi", "ru": "Бриллианты | Vondi"}'::jsonb,
 '{"sr": "Kupite dijamante i dijamantski nakit online", "en": "Buy diamonds and diamond jewelry online", "ru": "Купить бриллианты и бриллиантовые украшения онлайн"}'::jsonb,
 '💎', true),

('moderni-nakit', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/moderni-nakit', 9,
 '{"sr": "Moderni nakit", "en": "Fashion jewelry", "ru": "Модные украшения"}'::jsonb,
 '{"sr": "Modni nakit, biserne imitacije, avantura", "en": "Fashion jewelry, pearl imitations, costume jewelry", "ru": "Модные украшения, имитации жемчуга, бижутерия"}'::jsonb,
 '{"sr": "Moderni nakit | Vondi", "en": "Fashion jewelry | Vondi", "ru": "Модные украшения | Vondi"}'::jsonb,
 '{"sr": "Kupite moderni nakit online", "en": "Buy fashion jewelry online", "ru": "Купить модные украшения онлайн"}'::jsonb,
 '✨', true),

('satovski-dodaci', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi'), 2, 'nakit-i-satovi/satovski-dodaci', 10,
 '{"sr": "Satovski dodaci", "en": "Watch accessories", "ru": "Аксессуары для часов"}'::jsonb,
 '{"sr": "Narukvice za satove, kutije, winder-i", "en": "Watch straps, boxes, winders", "ru": "Ремешки для часов, коробки, виндеры"}'::jsonb,
 '{"sr": "Satovski dodaci | Vondi", "en": "Watch accessories | Vondi", "ru": "Аксессуары для часов | Vondi"}'::jsonb,
 '{"sr": "Kupite satovske dodatke online", "en": "Buy watch accessories online", "ru": "Купить аксессуары для часов онлайн"}'::jsonb,
 '⏱️', true),

-- =============================================================================
-- L2 for: 10. Knjige i mediji (Books & Media) - 10 categories
-- =============================================================================

('knjige-beletristika', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/knjige-beletristika', 1,
 '{"sr": "Knjige beletristika", "en": "Fiction books", "ru": "Художественная литература"}'::jsonb,
 '{"sr": "Romani, novele, poezija", "en": "Novels, short stories, poetry", "ru": "Романы, рассказы, поэзия"}'::jsonb,
 '{"sr": "Knjige beletristika | Vondi", "en": "Fiction books | Vondi", "ru": "Художественная литература | Vondi"}'::jsonb,
 '{"sr": "Kupite knjige beletristiku online", "en": "Buy fiction books online", "ru": "Купить художественную литературу онлайн"}'::jsonb,
 '📖', true),

('strucne-knjige', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/strucne-knjige', 2,
 '{"sr": "Stručne knjige", "en": "Non-fiction books", "ru": "Научная литература"}'::jsonb,
 '{"sr": "Priručnici, udžbenici, biografije", "en": "Handbooks, textbooks, biographies", "ru": "Справочники, учебники, биографии"}'::jsonb,
 '{"sr": "Stručne knjige | Vondi", "en": "Non-fiction books | Vondi", "ru": "Научная литература | Vondi"}'::jsonb,
 '{"sr": "Kupite stručne knjige online", "en": "Buy non-fiction books online", "ru": "Купить научную литературу онлайн"}'::jsonb,
 '📚', true),

('decije-knjige', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/decije-knjige', 3,
 '{"sr": "Dečije knjige", "en": "Children''s books", "ru": "Детские книги"}'::jsonb,
 '{"sr": "Bajke, slikovnice, edukativne knjige", "en": "Fairy tales, picture books, educational books", "ru": "Сказки, книжки с картинками, развивающие книги"}'::jsonb,
 '{"sr": "Dečije knjige | Vondi", "en": "Children''s books | Vondi", "ru": "Детские книги | Vondi"}'::jsonb,
 '{"sr": "Kupite dečije knjige online", "en": "Buy children''s books online", "ru": "Купить детские книги онлайн"}'::jsonb,
 '📕', true),

('casopisi', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/casopisi', 4,
 '{"sr": "Časopisi", "en": "Magazines", "ru": "Журналы"}'::jsonb,
 '{"sr": "Modni, sportski, naučni časopisi", "en": "Fashion, sports, science magazines", "ru": "Модные, спортивные, научные журналы"}'::jsonb,
 '{"sr": "Časopisi | Vondi", "en": "Magazines | Vondi", "ru": "Журналы | Vondi"}'::jsonb,
 '{"sr": "Kupite časopise online", "en": "Buy magazines online", "ru": "Купить журналы онлайн"}'::jsonb,
 '📰', true),

('stripovi', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/stripovi', 5,
 '{"sr": "Stripovi", "en": "Comics", "ru": "Комиксы"}'::jsonb,
 '{"sr": "Manga, Marvel, DC stripovi", "en": "Manga, Marvel, DC comics", "ru": "Манга, Marvel, DC комиксы"}'::jsonb,
 '{"sr": "Stripovi | Vondi", "en": "Comics | Vondi", "ru": "Комиксы | Vondi"}'::jsonb,
 '{"sr": "Kupite stripove online", "en": "Buy comics online", "ru": "Купить комиксы онлайн"}'::jsonb,
 '📜', true),

('filmovi', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/filmovi', 6,
 '{"sr": "Filmovi", "en": "Movies", "ru": "Фильмы"}'::jsonb,
 '{"sr": "DVD, Blu-ray, digitalni filmovi", "en": "DVDs, Blu-rays, digital movies", "ru": "DVD, Blu-ray, цифровые фильмы"}'::jsonb,
 '{"sr": "Filmovi | Vondi", "en": "Movies | Vondi", "ru": "Фильмы | Vondi"}'::jsonb,
 '{"sr": "Kupite filmove online", "en": "Buy movies online", "ru": "Купить фильмы онлайн"}'::jsonb,
 '🎬', true),

('muzika', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/muzika', 7,
 '{"sr": "Muzika", "en": "Music", "ru": "Музыка"}'::jsonb,
 '{"sr": "CD, vinili, digitalna muzika", "en": "CDs, vinyl, digital music", "ru": "CD, винил, цифровая музыка"}'::jsonb,
 '{"sr": "Muzika | Vondi", "en": "Music | Vondi", "ru": "Музыка | Vondi"}'::jsonb,
 '{"sr": "Kupite muziku online", "en": "Buy music online", "ru": "Купить музыку онлайн"}'::jsonb,
 '🎵', true),

('audio-knjige', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/audio-knjige', 8,
 '{"sr": "Audio knjige", "en": "Audiobooks", "ru": "Аудиокниги"}'::jsonb,
 '{"sr": "Audio knjige, podkasti na CD", "en": "Audiobooks, podcasts on CD", "ru": "Аудиокниги, подкасты на CD"}'::jsonb,
 '{"sr": "Audio knjige | Vondi", "en": "Audiobooks | Vondi", "ru": "Аудиокниги | Vondi"}'::jsonb,
 '{"sr": "Kupite audio knjige online", "en": "Buy audiobooks online", "ru": "Купить аудиокниги онлайн"}'::jsonb,
 '🎧', true),

('e-knjige', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/e-knjige', 9,
 '{"sr": "E-knjige", "en": "E-books", "ru": "Электронные книги"}'::jsonb,
 '{"sr": "Elektronske knjige, PDF, ePub", "en": "Electronic books, PDF, ePub", "ru": "Электронные книги, PDF, ePub"}'::jsonb,
 '{"sr": "E-knjige | Vondi", "en": "E-books | Vondi", "ru": "Электронные книги | Vondi"}'::jsonb,
 '{"sr": "Kupite e-knjige online", "en": "Buy e-books online", "ru": "Купить электронные книги онлайн"}'::jsonb,
 '📱', true),

('retke-knjige', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji'), 2, 'knjige-i-mediji/retke-knjige', 10,
 '{"sr": "Retke knjige", "en": "Rare books", "ru": "Редкие книги"}'::jsonb,
 '{"sr": "Antikvarne, kolekcionar ske, prvo izdanje", "en": "Antique, collectible, first editions", "ru": "Антикварные, коллекционные, первые издания"}'::jsonb,
 '{"sr": "Retke knjige | Vondi", "en": "Rare books | Vondi", "ru": "Редкие книги | Vondi"}'::jsonb,
 '{"sr": "Kupite retke knjige online", "en": "Buy rare books online", "ru": "Купить редкие книги онлайн"}'::jsonb,
 '📜', true);

-- Continue in next section due to length constraints
-- Progress: 54 L2 categories added (Automobilizam: 12, Kućni aparati: 12, Nakit: 10, Knjige: 10)

DO $$
DECLARE
    l2_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO l2_count FROM categories WHERE level = 2;
    RAISE NOTICE 'Part 3 section 1: % total L2 categories created', l2_count;
END $$;

-- =============================================================================
-- L2 for: 11. Kućni ljubimci (Pet Supplies) - 10 categories
-- =============================================================================

('hrana-za-pse', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/hrana-za-pse', 1,
 '{"sr": "Hrana za pse", "en": "Dog food", "ru": "Корм для собак"}'::jsonb,
 '{"sr": "Suva i konzervisana hrana za pse", "en": "Dry and canned dog food", "ru": "Сухой и консервированный корм для собак"}'::jsonb,
 '{"sr": "Hrana za pse | Vondi", "en": "Dog food | Vondi", "ru": "Корм для собак | Vondi"}'::jsonb,
 '{"sr": "Kupite hranu za pse online", "en": "Buy dog food online", "ru": "Купить корм для собак онлайн"}'::jsonb,
 '🐕', true),

('hrana-za-macke', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/hrana-za-macke', 2,
 '{"sr": "Hrana za mačke", "en": "Cat food", "ru": "Корм для кошек"}'::jsonb,
 '{"sr": "Suva i konzervisana hrana za mačke", "en": "Dry and canned cat food", "ru": "Сухой и консервированный корм для кошек"}'::jsonb,
 '{"sr": "Hrana za mačke | Vondi", "en": "Cat food | Vondi", "ru": "Корм для кошек | Vondi"}'::jsonb,
 '{"sr": "Kupite hranu za mačke online", "en": "Buy cat food online", "ru": "Купить корм для кошек онлайн"}'::jsonb,
 '🐈', true),

('igracke-za-ljubimce', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/igracke-za-ljubimce', 3,
 '{"sr": "Igračke za ljubimce", "en": "Pet toys", "ru": "Игрушки для питомцев"}'::jsonb,
 '{"sr": "Igračke za pse, mačke i druge ljubimce", "en": "Toys for dogs, cats and other pets", "ru": "Игрушки для собак, кошек и других питомцев"}'::jsonb,
 '{"sr": "Igračke za ljubimce | Vondi", "en": "Pet toys | Vondi", "ru": "Игрушки для питомцев | Vondi"}'::jsonb,
 '{"sr": "Kupite igračke za ljubimce online", "en": "Buy pet toys online", "ru": "Купить игрушки для питомцев онлайн"}'::jsonb,
 '🎾', true),

('oprema-za-pse', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/oprema-za-pse', 4,
 '{"sr": "Oprema za pse", "en": "Dog supplies", "ru": "Товары для собак"}'::jsonb,
 '{"sr": "Ogrlice, povodci, koševi, odeća", "en": "Collars, leashes, beds, clothing", "ru": "Ошейники, поводки, лежанки, одежда"}'::jsonb,
 '{"sr": "Oprema za pse | Vondi", "en": "Dog supplies | Vondi", "ru": "Товары для собак | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za pse online", "en": "Buy dog supplies online", "ru": "Купить товары для собак онлайн"}'::jsonb,
 '🦴', true),

('oprema-za-macke', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/oprema-za-macke', 5,
 '{"sr": "Oprema za mačke", "en": "Cat supplies", "ru": "Товары для кошек"}'::jsonb,
 '{"sr": "Kućice, grebalice, posude, pesak", "en": "Houses, scratching posts, bowls, litter", "ru": "Домики, когтеточки, миски, наполнитель"}'::jsonb,
 '{"sr": "Oprema za mačke | Vondi", "en": "Cat supplies | Vondi", "ru": "Товары для кошек | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za mačke online", "en": "Buy cat supplies online", "ru": "Купить товары для кошек онлайн"}'::jsonb,
 '🐱', true),

('akvarijumi', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/akvarijumi', 6,
 '{"sr": "Akvarijumi", "en": "Aquariums", "ru": "Аквариумы"}'::jsonb,
 '{"sr": "Akvarijumi, ribice, oprema, hrana", "en": "Aquariums, fish, equipment, food", "ru": "Аквариумы, рыбки, оборудование, корм"}'::jsonb,
 '{"sr": "Akvarijumi | Vondi", "en": "Aquariums | Vondi", "ru": "Аквариумы | Vondi"}'::jsonb,
 '{"sr": "Kupite akvarijume i opremu online", "en": "Buy aquariums and equipment online", "ru": "Купить аквариумы и оборудование онлайн"}'::jsonb,
 '🐠', true),

('ptice', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/ptice', 7,
 '{"sr": "Ptice", "en": "Birds", "ru": "Птицы"}'::jsonb,
 '{"sr": "Kavezi, hrana, igračke za ptice", "en": "Cages, food, toys for birds", "ru": "Клетки, корм, игрушки для птиц"}'::jsonb,
 '{"sr": "Ptice | Vondi", "en": "Birds | Vondi", "ru": "Птицы | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za ptice online", "en": "Buy bird supplies online", "ru": "Купить товары для птиц онлайн"}'::jsonb,
 '🦜', true),

('glodari', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/glodari', 8,
 '{"sr": "Glodari", "en": "Rodents", "ru": "Грызуны"}'::jsonb,
 '{"sr": "Kavezi, hrana za hrčke, zamorce", "en": "Cages, food for hamsters, guinea pigs", "ru": "Клетки, корм для хомяков, морских свинок"}'::jsonb,
 '{"sr": "Glodari | Vondi", "en": "Rodents | Vondi", "ru": "Грызуны | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za glodareon line", "en": "Buy rodent supplies online", "ru": "Купить товары для грызунов онлайн"}'::jsonb,
 '🐹', true),

('nega-ljubimaca', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/nega-ljubimaca', 9,
 '{"sr": "Nega ljubimaca", "en": "Pet grooming", "ru": "Уход за питомцами"}'::jsonb,
 '{"sr": "Šamponi, četke, makaze, veterinarski proizvodi", "en": "Shampoos, brushes, scissors, veterinary products", "ru": "Шампуни, щетки, ножницы, ветеринарные товары"}'::jsonb,
 '{"sr": "Nega ljubimaca | Vondi", "en": "Pet grooming | Vondi", "ru": "Уход за питомцами | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za negu ljubimaca online", "en": "Buy pet grooming products online", "ru": "Купить товары для ухода за питомцами онлайн"}'::jsonb,
 '✂️', true),

('terarijumi', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci'), 2, 'kucni-ljubimci/terarijumi', 10,
 '{"sr": "Terarijumi", "en": "Terrariums", "ru": "Террариумы"}'::jsonb,
 '{"sr": "Terarijumi, gmazovi, oprema", "en": "Terrariums, reptiles, equipment", "ru": "Террариумы, рептилии, оборудование"}'::jsonb,
 '{"sr": "Terarijumi | Vondi", "en": "Terrariums | Vondi", "ru": "Террариумы | Vondi"}'::jsonb,
 '{"sr": "Kupite terarijume i opremu online", "en": "Buy terrariums and equipment online", "ru": "Купить террариумы и оборудование онлайн"}'::jsonb,
 '🦎', true),

-- =============================================================================
-- L2 for: 12. Kancelarijski materijal (Office Supplies) - 8 categories
-- =============================================================================

('sveske-i-papir', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal'), 2, 'kancelarijski-materijal/sveske-i-papir', 1,
 '{"sr": "Sveske i papir", "en": "Notebooks & paper", "ru": "Тетради и бумага"}'::jsonb,
 '{"sr": "Sveske, blokovi, hartija za štampač", "en": "Notebooks, pads, printer paper", "ru": "Тетради, блокноты, бумага для принтера"}'::jsonb,
 '{"sr": "Sveske i papir | Vondi", "en": "Notebooks & paper | Vondi", "ru": "Тетради и бумага | Vondi"}'::jsonb,
 '{"sr": "Kupite sveske i papir online", "en": "Buy notebooks and paper online", "ru": "Купить тетради и бумагу онлайн"}'::jsonb,
 '📓', true),

('olovke-i-hemijske', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal'), 2, 'kancelarijski-materijal/olovke-i-hemijske', 2,
 '{"sr": "Olovke i hemijske", "en": "Pens & pencils", "ru": "Ручки и карандаши"}'::jsonb,
 '{"sr": "Hemijske olovke, grafitne, flomasters", "en": "Ballpoint pens, pencils, markers", "ru": "Шариковые ручки, карандаши, маркеры"}'::jsonb,
 '{"sr": "Olovke i hemijske | Vondi", "en": "Pens & pencils | Vondi", "ru": "Ручки и карандаши | Vondi"}'::jsonb,
 '{"sr": "Kupite olovke i hemijske online", "en": "Buy pens and pencils online", "ru": "Купить ручки и карандаши онлайн"}'::jsonb,
 '✏️', true),

('fascikle-i-registratori', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal'), 2, 'kancelarijski-materijal/fascikle-i-registratori', 3,
 '{"sr": "Fascikle i registratori", "en": "Folders & binders", "ru": "Папки и регистраторы"}'::jsonb,
 '{"sr": "Fascikle, registratori, fascikle sa klipom", "en": "Folders, binders, clipboards", "ru": "Папки, регистраторы, папки с зажимом"}'::jsonb,
 '{"sr": "Fascikle i registratori | Vondi", "en": "Folders & binders | Vondi", "ru": "Папки и регистраторы | Vondi"}'::jsonb,
 '{"sr": "Kupite fascikle i registratoreon line", "en": "Buy folders and binders online", "ru": "Купить папки и регистраторы онлайн"}'::jsonb,
 '📁', true),

('kancelarijski-pribor', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal'), 2, 'kancelarijski-materijal/kancelarijski-pribor', 4,
 '{"sr": "Kancelarijski pribor", "en": "Office supplies", "ru": "Канцелярские принадлежности"}'::jsonb,
 '{"sr": "Lepak, makaze, spajalice, klameri", "en": "Glue, scissors, staplers, clips", "ru": "Клей, ножницы, степлеры, скрепки"}'::jsonb,
 '{"sr": "Kancelarijski pribor | Vondi", "en": "Office supplies | Vondi", "ru": "Канцелярские принадлежности | Vondi"}'::jsonb,
 '{"sr": "Kupite kancelarijski pribor online", "en": "Buy office supplies online", "ru": "Купить канцелярские принадлежности онлайн"}'::jsonb,
 '📎', true),

('organizacija-stola', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal'), 2, 'kancelarijski-materijal/organizacija-stola', 5,
 '{"sr": "Organizacija stola", "en": "Desk organization", "ru": "Организация рабочего стола"}'::jsonb,
 '{"sr": "Držači, organajzeri, podmetači", "en": "Holders, organizers, desk pads", "ru": "Держатели, органайзеры, подставки"}'::jsonb,
 '{"sr": "Organizacija stola | Vondi", "en": "Desk organization | Vondi", "ru": "Организация рабочего стола | Vondi"}'::jsonb,
 '{"sr": "Kupite organizatore za sto online", "en": "Buy desk organizers online", "ru": "Купить органайзеры для стола онлайн"}'::jsonb,
 '🗂️', true),

('stampaci-i-toneri', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal'), 2, 'kancelarijski-materijal/stampaci-i-toneri', 6,
 '{"sr": "Štampači i toneri", "en": "Printers & toners", "ru": "Принтеры и тонеры"}'::jsonb,
 '{"sr": "Štampači, toneri, kertridži", "en": "Printers, toners, cartridges", "ru": "Принтеры, тонеры, картриджи"}'::jsonb,
 '{"sr": "Štampači i toneri | Vondi", "en": "Printers & toners | Vondi", "ru": "Принтеры и тонеры | Vondi"}'::jsonb,
 '{"sr": "Kupite štampače i tonere online", "en": "Buy printers and toners online", "ru": "Купить принтеры и тонеры онлайн"}'::jsonb,
 '🖨️', true),

('kalendari-i-planeri', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal'), 2, 'kancelarijski-materijal/kalendari-i-planeri', 7,
 '{"sr": "Kalendari i planeri", "en": "Calendars & planners", "ru": "Календари и планировщики"}'::jsonb,
 '{"sr": "Zidni kalendari, planeri, dnevnici", "en": "Wall calendars, planners, diaries", "ru": "Настенные календари, планировщики, ежедневники"}'::jsonb,
 '{"sr": "Kalendari i planeri | Vondi", "en": "Calendars & planners | Vondi", "ru": "Календари и планировщики | Vondi"}'::jsonb,
 '{"sr": "Kupite kalendare i planere online", "en": "Buy calendars and planners online", "ru": "Купить календари и планировщики онлайн"}'::jsonb,
 '📅', true),

('table-i-stikeri', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal'), 2, 'kancelarijski-materijal/table-i-stikeri', 8,
 '{"sr": "Table i stikeri", "en": "Boards & stickers", "ru": "Доски и стикеры"}'::jsonb,
 '{"sr": "Bele table, magnetne table, post-it", "en": "Whiteboards, magnetic boards, post-its", "ru": "Белые доски, магнитные доски, стикеры"}'::jsonb,
 '{"sr": "Table i stikeri | Vondi", "en": "Boards & stickers | Vondi", "ru": "Доски и стикеры | Vondi"}'::jsonb,
 '{"sr": "Kupite table i stikere online", "en": "Buy boards and stickers online", "ru": "Купить доски и стикеры онлайн"}'::jsonb,
 '📋', true),

-- =============================================================================
-- L2 for: 13. Muzički instrumenti (Musical Instruments) - 8 categories
-- =============================================================================

('gitare', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti'), 2, 'muzicki-instrumenti/gitare', 1,
 '{"sr": "Gitare", "en": "Guitars", "ru": "Гитары"}'::jsonb,
 '{"sr": "Akustične, električne, bas gitare", "en": "Acoustic, electric, bass guitars", "ru": "Акустические, электрические, бас-гитары"}'::jsonb,
 '{"sr": "Gitare | Vondi", "en": "Guitars | Vondi", "ru": "Гитары | Vondi"}'::jsonb,
 '{"sr": "Kupite gitare online", "en": "Buy guitars online", "ru": "Купить гитары онлайн"}'::jsonb,
 '🎸', true),

('klavijature', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti'), 2, 'muzicki-instrumenti/klavijature', 2,
 '{"sr": "Klavijature", "en": "Keyboards", "ru": "Клавишные"}'::jsonb,
 '{"sr": "Klavijature, sintisajzeri, pianina", "en": "Keyboards, synthesizers, pianos", "ru": "Клавиатуры, синтезаторы, пианино"}'::jsonb,
 '{"sr": "Klavijature | Vondi", "en": "Keyboards | Vondi", "ru": "Клавишные | Vondi"}'::jsonb,
 '{"sr": "Kupite klavijature online", "en": "Buy keyboards online", "ru": "Купить клавишные онлайн"}'::jsonb,
 '🎹', true),

('bubnjevi', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti'), 2, 'muzicki-instrumenti/bubnjevi', 3,
 '{"sr": "Bubnjevi", "en": "Drums", "ru": "Барабаны"}'::jsonb,
 '{"sr": "Akustični i elektronski bubnjevi, činele", "en": "Acoustic and electronic drums, cymbals", "ru": "Акустические и электронные барабаны, тарелки"}'::jsonb,
 '{"sr": "Bubnjevi | Vondi", "en": "Drums | Vondi", "ru": "Барабаны | Vondi"}'::jsonb,
 '{"sr": "Kupite bubnjeve online", "en": "Buy drums online", "ru": "Купить барабаны онлайн"}'::jsonb,
 '🥁', true),

('duvacki-instrumenti', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti'), 2, 'muzicki-instrumenti/duvacki-instrumenti', 4,
 '{"sr": "Duvački instrumenti", "en": "Wind instruments", "ru": "Духовые инструменты"}'::jsonb,
 '{"sr": "Saksofoni, flaute, trube, klarineti", "en": "Saxophones, flutes, trumpets, clarinets", "ru": "Саксофоны, флейты, трубы, кларнеты"}'::jsonb,
 '{"sr": "Duvački instrumenti | Vondi", "en": "Wind instruments | Vondi", "ru": "Духовые инструменты | Vondi"}'::jsonb,
 '{"sr": "Kupite duvačke instrumente online", "en": "Buy wind instruments online", "ru": "Купить духовые инструменты онлайн"}'::jsonb,
 '🎺', true),

('violina-i-gudacki', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti'), 2, 'muzicki-instrumenti/violina-i-gudacki', 5,
 '{"sr": "Violina i gudački", "en": "Violin & strings", "ru": "Скрипка и струнные"}'::jsonb,
 '{"sr": "Violina, viola, violončelo, kontrabas", "en": "Violin, viola, cello, double bass", "ru": "Скрипка, альт, виолончель, контрабас"}'::jsonb,
 '{"sr": "Violina i gudački | Vondi", "en": "Violin & strings | Vondi", "ru": "Скрипка и струнные | Vondi"}'::jsonb,
 '{"sr": "Kupite violinu i gudačke instrumente online", "en": "Buy violin and string instruments online", "ru": "Купить скрипку и струнные инструменты онлайн"}'::jsonb,
 '🎻', true),

('muzicka-oprema', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti'), 2, 'muzicki-instrumenti/muzicka-oprema', 6,
 '{"sr": "Muzička oprema", "en": "Music equipment", "ru": "Музыкальное оборудование"}'::jsonb,
 '{"sr": "Pojačala, efekti, mikrofoni, kablovi", "en": "Amplifiers, effects, microphones, cables", "ru": "Усилители, эффекты, микрофоны, кабели"}'::jsonb,
 '{"sr": "Muzička oprema | Vondi", "en": "Music equipment | Vondi", "ru": "Музыкальное оборудование | Vondi"}'::jsonb,
 '{"sr": "Kupite muzičku opremu online", "en": "Buy music equipment online", "ru": "Купить музыкальное оборудование онлайн"}'::jsonb,
 '🎤', true),

('dj-oprema', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti'), 2, 'muzicki-instrumenti/dj-oprema', 7,
 '{"sr": "DJ oprema", "en": "DJ equipment", "ru": "DJ оборудование"}'::jsonb,
 '{"sr": "Gramofoni, mikšete, kontroleri", "en": "Turntables, mixers, controllers", "ru": "Вертушки, микшеры, контроллеры"}'::jsonb,
 '{"sr": "DJ oprema | Vondi", "en": "DJ equipment | Vondi", "ru": "DJ оборудование | Vondi"}'::jsonb,
 '{"sr": "Kupite DJ opremu online", "en": "Buy DJ equipment online", "ru": "Купить DJ оборудование онлайн"}'::jsonb,
 '🎧', true),

('note-i-priručnici', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti'), 2, 'muzicki-instrumenti/note-i-prirucnici', 8,
 '{"sr": "Note i priručnici", "en": "Sheet music & guides", "ru": "Ноты и учебники"}'::jsonb,
 '{"sr": "Note, udžbenici, priručnici za muziku", "en": "Sheet music, textbooks, music guides", "ru": "Ноты, учебники, музыкальные пособия"}'::jsonb,
 '{"sr": "Note i priručnici | Vondi", "en": "Sheet music & guides | Vondi", "ru": "Ноты и учебники | Vondi"}'::jsonb,
 '{"sr": "Kupite note i priručnike online", "en": "Buy sheet music and guides online", "ru": "Купить ноты и учебники онлайн"}'::jsonb,
 '🎼', true);

-- Continue with remaining 5 L1 categories...
-- Progress: 90 L2 categories (Ljubimci: 10, Kancelarija: 8, Muzika: 8)

DO $$
DECLARE
    l2_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO l2_count FROM categories WHERE level = 2;
    RAISE NOTICE 'Part 3 section 2: % total L2 categories created', l2_count;
END $$;

-- =============================================================================
-- L2 for: 14. Hrana i piće (Food & Beverages) - 10 categories
-- =============================================================================

('organska-hrana', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/organska-hrana', 1,
 '{"sr": "Organska hrana", "en": "Organic food", "ru": "Органическая еда"}'::jsonb,
 '{"sr": "Organski proizvodi, zdrava hrana", "en": "Organic products, healthy food", "ru": "Органические продукты, здоровая еда"}'::jsonb,
 '{"sr": "Organska hrana | Vondi", "en": "Organic food | Vondi", "ru": "Органическая еда | Vondi"}'::jsonb,
 '{"sr": "Kupite organsku hranu online", "en": "Buy organic food online", "ru": "Купить органическую еду онлайн"}'::jsonb,
 '🥬', true),

('kafa-i-caj', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/kafa-i-caj', 2,
 '{"sr": "Kafa i čaj", "en": "Coffee & tea", "ru": "Кофе и чай"}'::jsonb,
 '{"sr": "Kafa u zrnu, mleta, čajevi", "en": "Coffee beans, ground coffee, teas", "ru": "Кофе в зернах, молотый, чаи"}'::jsonb,
 '{"sr": "Kafa i čaj | Vondi", "en": "Coffee & tea | Vondi", "ru": "Кофе и чай | Vondi"}'::jsonb,
 '{"sr": "Kupite kafu i čaj online", "en": "Buy coffee and tea online", "ru": "Купить кофе и чай онлайн"}'::jsonb,
 '☕', true),

('slatkisi', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/slatkisi', 3,
 '{"sr": "Slatkiši", "en": "Sweets", "ru": "Сладости"}'::jsonb,
 '{"sr": "Čokolada, bomboni, keks i, torte", "en": "Chocolate, candies, cookies, cakes", "ru": "Шоколад, конфеты, печенье, торты"}'::jsonb,
 '{"sr": "Slatkiši | Vondi", "en": "Sweets | Vondi", "ru": "Сладости | Vondi"}'::jsonb,
 '{"sr": "Kupite slatkiše online", "en": "Buy sweets online", "ru": "Купить сладости онлайн"}'::jsonb,
 '🍫', true),

('sokovi-i-napici', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/sokovi-i-napici', 4,
 '{"sr": "Sokovi i napici", "en": "Juices & drinks", "ru": "Соки и напитки"}'::jsonb,
 '{"sr": "Prirodni sokovi, gazirani napici", "en": "Natural juices, carbonated drinks", "ru": "Натуральные соки, газированные напитки"}'::jsonb,
 '{"sr": "Sokovi i napici | Vondi", "en": "Juices & drinks | Vondi", "ru": "Соки и напитки | Vondi"}'::jsonb,
 '{"sr": "Kupite sokove i napitke online", "en": "Buy juices and drinks online", "ru": "Купить соки и напитки онлайн"}'::jsonb,
 '🧃', true),

('zacini-i-dodaci', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/zacini-i-dodaci', 5,
 '{"sr": "Začini i dodaci", "en": "Spices & condiments", "ru": "Специи и приправы"}'::jsonb,
 '{"sr": "Začini, ulja, sosevi, sirće", "en": "Spices, oils, sauces, vinegar", "ru": "Специи, масла, соусы, уксус"}'::jsonb,
 '{"sr": "Začini i dodaci | Vondi", "en": "Spices & condiments | Vondi", "ru": "Специи и приправы | Vondi"}'::jsonb,
 '{"sr": "Kupite začine i dodatke online", "en": "Buy spices and condiments online", "ru": "Купить специи и приправы онлайн"}'::jsonb,
 '🧂', true),

('tjestenina-i-zitarice', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/tjestenina-i-zitarice', 6,
 '{"sr": "Tjestenina i žitarice", "en": "Pasta & cereals", "ru": "Макароны и крупы"}'::jsonb,
 '{"sr": "Tjestenina, pirinač, kaše", "en": "Pasta, rice, porridge", "ru": "Макароны, рис, каши"}'::jsonb,
 '{"sr": "Tjestenina i žitarice | Vondi", "en": "Pasta & cereals | Vondi", "ru": "Макароны и крупы | Vondi"}'::jsonb,
 '{"sr": "Kupite tjesteninu i žitarice online", "en": "Buy pasta and cereals online", "ru": "Купить макароны и крупы онлайн"}'::jsonb,
 '🍝', true),

('konzerve', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/konzerve', 7,
 '{"sr": "Konzerve", "en": "Canned food", "ru": "Консервы"}'::jsonb,
 '{"sr": "Konzervisana riba, povrće, voće", "en": "Canned fish, vegetables, fruits", "ru": "Консервированная рыба, овощи, фрукты"}'::jsonb,
 '{"sr": "Konzerve | Vondi", "en": "Canned food | Vondi", "ru": "Консервы | Vondi"}'::jsonb,
 '{"sr": "Kupite konzerve online", "en": "Buy canned food online", "ru": "Купить консервы онлайн"}'::jsonb,
 '🥫', true),

('mlecni-proizvodi', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/mlecni-proizvodi', 8,
 '{"sr": "Mlečni proizvodi", "en": "Dairy products", "ru": "Молочные продукты"}'::jsonb,
 '{"sr": "Mleko, sir, jogurt, kajmak", "en": "Milk, cheese, yogurt, cream", "ru": "Молоко, сыр, йогурт, сливки"}'::jsonb,
 '{"sr": "Mlečni proizvodi | Vondi", "en": "Dairy products | Vondi", "ru": "Молочные продукты | Vondi"}'::jsonb,
 '{"sr": "Kupite mlečne proizvode online", "en": "Buy dairy products online", "ru": "Купить молочные продукты онлайн"}'::jsonb,
 '🥛', true),

('peciva', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/peciva', 9,
 '{"sr": "Peciva", "en": "Bakery", "ru": "Выпечка"}'::jsonb,
 '{"sr": "Hleb, kifle, pekarski proizvodi", "en": "Bread, rolls, bakery products", "ru": "Хлеб, булочки, хлебобулочные изделия"}'::jsonb,
 '{"sr": "Peciva | Vondi", "en": "Bakery | Vondi", "ru": "Выпечка | Vondi"}'::jsonb,
 '{"sr": "Kupite peciva online", "en": "Buy bakery products online", "ru": "Купить выпечку онлайн"}'::jsonb,
 '🥐', true),

('delikatesi', (SELECT id FROM categories WHERE slug = 'hrana-i-pice'), 2, 'hrana-i-pice/delikatesi', 10,
 '{"sr": "Delikatesi", "en": "Delicacies", "ru": "Деликатесы"}'::jsonb,
 '{"sr": "Pršuta, sir, maslinovo ulje", "en": "Prosciutto, cheese, olive oil", "ru": "Прошутто, сыр, оливковое масло"}'::jsonb,
 '{"sr": "Delikatesi | Vondi", "en": "Delicacies | Vondi", "ru": "Деликатесы | Vondi"}'::jsonb,
 '{"sr": "Kupite delikatese online", "en": "Buy delicacies online", "ru": "Купить деликатесы онлайн"}'::jsonb,
 '🧀', true),

-- =============================================================================
-- L2 for: 15. Umetnost i rukotvorine (Art & Crafts) - 8 categories
-- =============================================================================

('materijali-za-slikanje', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine'), 2, 'umetnost-i-rukotvorine/materijali-za-slikanje', 1,
 '{"sr": "Materijali za slikanje", "en": "Painting supplies", "ru": "Материалы для рисования"}'::jsonb,
 '{"sr": "Boje, četkice, platna, moleri", "en": "Paints, brushes, canvases, easels", "ru": "Краски, кисти, холсты, мольберты"}'::jsonb,
 '{"sr": "Materijali za slikanje | Vondi", "en": "Painting supplies | Vondi", "ru": "Материалы для рисования | Vondi"}'::jsonb,
 '{"sr": "Kupite materijale za slikanje online", "en": "Buy painting supplies online", "ru": "Купить материалы для рисования онлайн"}'::jsonb,
 '🎨', true),

('rucni-rad', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine'), 2, 'umetnost-i-rukotvorine/rucni-rad', 2,
 '{"sr": "Ručni rad", "en": "Handmade crafts", "ru": "Ручная работа"}'::jsonb,
 '{"sr": "Vez, heklanje, pletenje, DIY", "en": "Embroidery, crochet, knitting, DIY", "ru": "Вышивка, вязание крючком, вязание, DIY"}'::jsonb,
 '{"sr": "Ručni rad | Vondi", "en": "Handmade crafts | Vondi", "ru": "Ручная работа | Vondi"}'::jsonb,
 '{"sr": "Kupite materijale za ručni rad online", "en": "Buy handmade craft supplies online", "ru": "Купить материалы для ручной работы онлайн"}'::jsonb,
 '🧶', true),

('skulptura', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine'), 2, 'umetnost-i-rukotvorine/skulptura', 3,
 '{"sr": "Skulptura", "en": "Sculpture", "ru": "Скульптура"}'::jsonb,
 '{"sr": "Glina, alati za modelovanje, gips", "en": "Clay, modeling tools, plaster", "ru": "Глина, инструменты для лепки, гипс"}'::jsonb,
 '{"sr": "Skulptura | Vondi", "en": "Sculpture | Vondi", "ru": "Скульптура | Vondi"}'::jsonb,
 '{"sr": "Kupite materijale za skulpturu online", "en": "Buy sculpture supplies online", "ru": "Купить материалы для скульптуры онлайн"}'::jsonb,
 '🗿', true),

('umetnicke-slike', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine'), 2, 'umetnost-i-rukotvorine/umetnicke-slike', 4,
 '{"sr": "Umetničke slike", "en": "Artwork", "ru": "Художественные картины"}'::jsonb,
 '{"sr": "Slike na platnu, plakati, grafike", "en": "Canvas paintings, posters, prints", "ru": "Картины на холсте, постеры, принты"}'::jsonb,
 '{"sr": "Umetničke slike | Vondi", "en": "Artwork | Vondi", "ru": "Художественные картины | Vondi"}'::jsonb,
 '{"sr": "Kupite umetničke slike online", "en": "Buy artwork online", "ru": "Купить художественные картины онлайн"}'::jsonb,
 '🖼️', true),

('papir-i-karton', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine'), 2, 'umetnost-i-rukotvorine/papir-i-karton', 5,
 '{"sr": "Papir i karton", "en": "Paper & cardboard", "ru": "Бумага и картон"}'::jsonb,
 '{"sr": "Papir za crtanje, karton, origami", "en": "Drawing paper, cardboard, origami", "ru": "Бумага для рисования, картон, оригами"}'::jsonb,
 '{"sr": "Papir i karton | Vondi", "en": "Paper & cardboard | Vondi", "ru": "Бумага и картон | Vondi"}'::jsonb,
 '{"sr": "Kupite papir i karton online", "en": "Buy paper and cardboard online", "ru": "Купить бумагу и картон онлайн"}'::jsonb,
 '📄', true),

('kreativni-alati', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine'), 2, 'umetnost-i-rukotvorine/kreativni-alati', 6,
 '{"sr": "Kreativni alati", "en": "Creative tools", "ru": "Творческие инструменты"}'::jsonb,
 '{"sr": "Makaze, lepak, sekateri, noževi", "en": "Scissors, glue, punches, knives", "ru": "Ножницы, клей, дыроколы, ножи"}'::jsonb,
 '{"sr": "Kreativni alati | Vondi", "en": "Creative tools | Vondi", "ru": "Творческие инструменты | Vondi"}'::jsonb,
 '{"sr": "Kupite kreativne alate online", "en": "Buy creative tools online", "ru": "Купить творческие инструменты онлайн"}'::jsonb,
 '✂️', true),

('nakit-rucni-rad', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine'), 2, 'umetnost-i-rukotvorine/nakit-rucni-rad', 7,
 '{"sr": "Nakit ručni rad", "en": "Handmade jewelry", "ru": "Украшения ручной работы"}'::jsonb,
 '{"sr": "Biseri, sagovi, materijali za nakit", "en": "Beads, wires, jewelry-making supplies", "ru": "Бусины, проволока, материалы для украшений"}'::jsonb,
 '{"sr": "Nakit ručni rad | Vondi", "en": "Handmade jewelry | Vondi", "ru": "Украшения ручной работы | Vondi"}'::jsonb,
 '{"sr": "Kupite materijale za izradu nakita online", "en": "Buy jewelry-making supplies online", "ru": "Купить материалы для создания украшений онлайн"}'::jsonb,
 '💍', true),

('dekorativne-tehnike', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine'), 2, 'umetnost-i-rukotvorine/dekorativne-tehnike', 8,
 '{"sr": "Dekorativne tehnike", "en": "Decorative techniques", "ru": "Декоративные техники"}'::jsonb,
 '{"sr": "Decoupage, scrapbooking, pirogravura", "en": "Decoupage, scrapbooking, pyrography", "ru": "Декупаж, скрапбукинг, пирография"}'::jsonb,
 '{"sr": "Dekorativne tehnike | Vondi", "en": "Decorative techniques | Vondi", "ru": "Декоративные техники | Vondi"}'::jsonb,
 '{"sr": "Kupite materijale za dekorativne tehnike online", "en": "Buy decorative technique supplies online", "ru": "Купить материалы для декоративных техник онлайн"}'::jsonb,
 '🎭', true),

-- =============================================================================
-- L2 for: 16. Industrija i alati (Industrial & Tools) - 10 categories
-- =============================================================================

('rucni-alati', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/rucni-alati', 1,
 '{"sr": "Ručni alati", "en": "Hand tools", "ru": "Ручные инструменты"}'::jsonb,
 '{"sr": "Čekići, klješta, odvijači, pile", "en": "Hammers, pliers, screwdrivers, saws", "ru": "Молотки, плоскогубцы, отвертки, пилы"}'::jsonb,
 '{"sr": "Ručni alati | Vondi", "en": "Hand tools | Vondi", "ru": "Ручные инструменты | Vondi"}'::jsonb,
 '{"sr": "Kupite ručne alate online", "en": "Buy hand tools online", "ru": "Купить ручные инструменты онлайн"}'::jsonb,
 '🔨', true),

('elektricni-alati', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/elektricni-alati', 2,
 '{"sr": "Električni alati", "en": "Power tools", "ru": "Электроинструменты"}'::jsonb,
 '{"sr": "Bušilice, brusilice, testere", "en": "Drills, grinders, saws", "ru": "Дрели, болгарки, пилы"}'::jsonb,
 '{"sr": "Električni alati | Vondi", "en": "Power tools | Vondi", "ru": "Электроинструменты | Vondi"}'::jsonb,
 '{"sr": "Kupite električne alate online", "en": "Buy power tools online", "ru": "Купить электроинструменты онлайн"}'::jsonb,
 '⚡', true),

('gradjevinski-materijali', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/gradjevinski-materijali', 3,
 '{"sr": "Građevinski materijali", "en": "Building materials", "ru": "Строительные материалы"}'::jsonb,
 '{"sr": "Cement, gips, opeka, pločice", "en": "Cement, plaster, bricks, tiles", "ru": "Цемент, гипс, кирпичи, плитка"}'::jsonb,
 '{"sr": "Građevinski materijali | Vondi", "en": "Building materials | Vondi", "ru": "Строительные материалы | Vondi"}'::jsonb,
 '{"sr": "Kupite građevinske materijale online", "en": "Buy building materials online", "ru": "Купить строительные материалы онлайн"}'::jsonb,
 '🏗️', true),

('boje-i-lakovi', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/boje-i-lakovi', 4,
 '{"sr": "Boje i lakovi", "en": "Paints & varnishes", "ru": "Краски и лаки"}'::jsonb,
 '{"sr": "Zidne boje, lakovi, farbe", "en": "Wall paints, varnishes, coatings", "ru": "Краски для стен, лаки, покрытия"}'::jsonb,
 '{"sr": "Boje i lakovi | Vondi", "en": "Paints & varnishes | Vondi", "ru": "Краски и лаки | Vondi"}'::jsonb,
 '{"sr": "Kupite boje i lakove online", "en": "Buy paints and varnishes online", "ru": "Купить краски и лаки онлайн"}'::jsonb,
 '🎨', true),

('meraci-i-instrumenti', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/meraci-i-instrumenti', 5,
 '{"sr": "Merači i instrumenti", "en": "Measuring tools", "ru": "Измерительные инструменты"}'::jsonb,
 '{"sr": "Metar, libela, laserski merači", "en": "Tape measure, level, laser measurers", "ru": "Рулетка, уровень, лазерные измерители"}'::jsonb,
 '{"sr": "Merači i instrumenti | Vondi", "en": "Measuring tools | Vondi", "ru": "Измерительные инструменты | Vondi"}'::jsonb,
 '{"sr": "Kupite merače i instrumente online", "en": "Buy measuring tools online", "ru": "Купить измерительные инструменты онлайн"}'::jsonb,
 '📐', true),

('zastita-na-radu', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/zastita-na-radu', 6,
 '{"sr": "Zaštita na radu", "en": "Safety equipment", "ru": "Защита на работе"}'::jsonb,
 '{"sr": "Kacige, rukavice, naočare, zaštitna odeća", "en": "Helmets, gloves, glasses, protective clothing", "ru": "Каски, перчатки, очки, защитная одежда"}'::jsonb,
 '{"sr": "Zaštita na radu | Vondi", "en": "Safety equipment | Vondi", "ru": "Защита на работе | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za zaštitu na radu online", "en": "Buy safety equipment online", "ru": "Купить средства защиты на работе онлайн"}'::jsonb,
 '🦺', true),

('hidraulika-i-sanitarije', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/hidraulika-i-sanitarije', 7,
 '{"sr": "Hidraulika i sanitarije", "en": "Plumbing & sanitary", "ru": "Гидравлика и сантехника"}'::jsonb,
 '{"sr": "Cevi, slavine, ventili, pumpe", "en": "Pipes, faucets, valves, pumps", "ru": "Трубы, смесители, клапаны, насосы"}'::jsonb,
 '{"sr": "Hidraulika i sanitarije | Vondi", "en": "Plumbing & sanitary | Vondi", "ru": "Гидравлика и сантехника | Vondi"}'::jsonb,
 '{"sr": "Kupite hidrauliku i sanitarije online", "en": "Buy plumbing and sanitary online", "ru": "Купить гидравлику и сантехнику онлайн"}'::jsonb,
 '🚰', true),

('elektromaterijal', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/elektromaterijal', 8,
 '{"sr": "Elektromaterijal", "en": "Electrical materials", "ru": "Электроматериалы"}'::jsonb,
 '{"sr": "Kablovi, utičnice, prekidači, osigurači", "en": "Cables, sockets, switches, fuses", "ru": "Кабели, розетки, выключатели, предохранители"}'::jsonb,
 '{"sr": "Elektromaterijal | Vondi", "en": "Electrical materials | Vondi", "ru": "Электроматериалы | Vondi"}'::jsonb,
 '{"sr": "Kupite elektromaterijal online", "en": "Buy electrical materials online", "ru": "Купить электроматериалы онлайн"}'::jsonb,
 '🔌', true),

('radne-masine', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/radne-masine', 9,
 '{"sr": "Radne mašine", "en": "Industrial machinery", "ru": "Рабочие машины"}'::jsonb,
 '{"sr": "Generatori, kompresori, pumpe", "en": "Generators, compressors, pumps", "ru": "Генераторы, компрессоры, насосы"}'::jsonb,
 '{"sr": "Radne mašine | Vondi", "en": "Industrial machinery | Vondi", "ru": "Рабочие машины | Vondi"}'::jsonb,
 '{"sr": "Kupite radne mašine online", "en": "Buy industrial machinery online", "ru": "Купить рабочие машины онлайн"}'::jsonb,
 '⚙️', true),

('pricvrscivaњi-elementi', (SELECT id FROM categories WHERE slug = 'industrija-i-alati'), 2, 'industrija-i-alati/pricvrscivaњi-elementi', 10,
 '{"sr": "Pričvrsćivači elementi", "en": "Fasteners", "ru": "Крепежные элементы"}'::jsonb,
 '{"sr": "Šrafovi, ekseri, viljuške, podloške", "en": "Screws, nails, bolts, washers", "ru": "Винты, гвозди, болты, шайбы"}'::jsonb,
 '{"sr": "Pričvrsćivači elementi | Vondi", "en": "Fasteners | Vondi", "ru": "Крепежные элементы | Vondi"}'::jsonb,
 '{"sr": "Kupite pričvrsćivače elemente online", "en": "Buy fasteners online", "ru": "Купить крепежные элементы онлайн"}'::jsonb,
 '🔩', true),

-- =============================================================================
-- L2 for: 17. Usluge (Services) - 10 categories
-- =============================================================================

('popravke', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/popravke', 1,
 '{"sr": "Popravke", "en": "Repairs", "ru": "Ремонт"}'::jsonb,
 '{"sr": "Popravke aparata, odeće, obuće", "en": "Appliance, clothing, shoe repairs", "ru": "Ремонт техники, одежды, обуви"}'::jsonb,
 '{"sr": "Popravke | Vondi", "en": "Repairs | Vondi", "ru": "Ремонт | Vondi"}'::jsonb,
 '{"sr": "Usluge popravki online", "en": "Repair services online", "ru": "Услуги ремонта онлайн"}'::jsonb,
 '🔧', true),

('ciscenje', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/ciscenje', 2,
 '{"sr": "Čišćenje", "en": "Cleaning", "ru": "Уборка"}'::jsonb,
 '{"sr": "Čišćenje stanova, kancelarija, pranje prozora", "en": "Apartment, office cleaning, window washing", "ru": "Уборка квартир, офисов, мытье окон"}'::jsonb,
 '{"sr": "Čišćenje | Vondi", "en": "Cleaning | Vondi", "ru": "Уборка | Vondi"}'::jsonb,
 '{"sr": "Usluge čišćenja online", "en": "Cleaning services online", "ru": "Услуги уборки онлайн"}'::jsonb,
 '🧹', true),

('prevoz', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/prevoz', 3,
 '{"sr": "Prevoz", "en": "Transportation", "ru": "Транспорт"}'::jsonb,
 '{"sr": "Selidbe, prevoz robe, rent-a-car", "en": "Moving, goods transportation, car rental", "ru": "Переезды, перевозка грузов, аренда авто"}'::jsonb,
 '{"sr": "Prevoz | Vondi", "en": "Transportation | Vondi", "ru": "Транспорт | Vondi"}'::jsonb,
 '{"sr": "Usluge prevoza online", "en": "Transportation services online", "ru": "Транспортные услуги онлайн"}'::jsonb,
 '🚚', true),

('edukacija', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/edukacija', 4,
 '{"sr": "Edukacija", "en": "Education", "ru": "Образование"}'::jsonb,
 '{"sr": "Kursevi, casovi, obuke, online učenje", "en": "Courses, lessons, training, online learning", "ru": "Курсы, уроки, тренинги, онлайн-обучение"}'::jsonb,
 '{"sr": "Edukacija | Vondi", "en": "Education | Vondi", "ru": "Образование | Vondi"}'::jsonb,
 '{"sr": "Edukativne usluge online", "en": "Educational services online", "ru": "Образовательные услуги онлайн"}'::jsonb,
 '📚', true),

('zdravstvene-usluge', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/zdravstvene-usluge', 5,
 '{"sr": "Zdravstvene usluge", "en": "Healthcare services", "ru": "Медицинские услуги"}'::jsonb,
 '{"sr": "Pregledi, terapije, nega, masaže", "en": "Examinations, therapies, care, massages", "ru": "Осмотры, терапии, уход, массажи"}'::jsonb,
 '{"sr": "Zdravstvene usluge | Vondi", "en": "Healthcare services | Vondi", "ru": "Медицинские услуги | Vondi"}'::jsonb,
 '{"sr": "Zdravstvene usluge online", "en": "Healthcare services online", "ru": "Медицинские услуги онлайн"}'::jsonb,
 '🩺', true),

('lepota-usluge', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/lepota-usluge', 6,
 '{"sr": "Lepota usluge", "en": "Beauty services", "ru": "Услуги красоты"}'::jsonb,
 '{"sr": "Frizure, manikir, kozmetika, spa", "en": "Hairstyles, manicures, cosmetics, spa", "ru": "Прически, маникюр, косметика, спа"}'::jsonb,
 '{"sr": "Lepota usluge | Vondi", "en": "Beauty services | Vondi", "ru": "Услуги красоты | Vondi"}'::jsonb,
 '{"sr": "Usluge lepote online", "en": "Beauty services online", "ru": "Услуги красоты онлайн"}'::jsonb,
 '💇', true),

('event-organizacija', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/event-organizacija', 7,
 '{"sr": "Event organizacija", "en": "Event planning", "ru": "Организация мероприятий"}'::jsonb,
 '{"sr": "Organizacija venčanja, rođendana, konferencija", "en": "Wedding, birthday, conference planning", "ru": "Организация свадеб, дней рождения, конференций"}'::jsonb,
 '{"sr": "Event organizacija | Vondi", "en": "Event planning | Vondi", "ru": "Организация мероприятий | Vondi"}'::jsonb,
 '{"sr": "Usluge event organizacije online", "en": "Event planning services online", "ru": "Услуги организации мероприятий онлайн"}'::jsonb,
 '🎉', true),

('foto-i-video', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/foto-i-video', 8,
 '{"sr": "Foto i video", "en": "Photo & video", "ru": "Фото и видео"}'::jsonb,
 '{"sr": "Fotografisanje, snimanje, montaža", "en": "Photography, filming, editing", "ru": "Фотосъемка, видеосъемка, монтаж"}'::jsonb,
 '{"sr": "Foto i video | Vondi", "en": "Photo & video | Vondi", "ru": "Фото и видео | Vondi"}'::jsonb,
 '{"sr": "Usluge foto i video online", "en": "Photo and video services online", "ru": "Услуги фото и видео онлайн"}'::jsonb,
 '📸', true),

('pravne-usluge', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/pravne-usluge', 9,
 '{"sr": "Pravne usluge", "en": "Legal services", "ru": "Юридические услуги"}'::jsonb,
 '{"sr": "Advokati, notari, pravni saveti", "en": "Lawyers, notaries, legal advice", "ru": "Адвокаты, нотариусы, юридические консультации"}'::jsonb,
 '{"sr": "Pravne usluge | Vondi", "en": "Legal services | Vondi", "ru": "Юридические услуги | Vondi"}'::jsonb,
 '{"sr": "Pravne usluge online", "en": "Legal services online", "ru": "Юридические услуги онлайн"}'::jsonb,
 '⚖️', true),

('it-usluge', (SELECT id FROM categories WHERE slug = 'usluge'), 2, 'usluge/it-usluge', 10,
 '{"sr": "IT usluge", "en": "IT services", "ru": "IT услуги"}'::jsonb,
 '{"sr": "Web dizajn, programiranje, hosting, IT podrška", "en": "Web design, programming, hosting, IT support", "ru": "Веб-дизайн, программирование, хостинг, IT поддержка"}'::jsonb,
 '{"sr": "IT usluge | Vondi", "en": "IT services | Vondi", "ru": "IT услуги | Vondi"}'::jsonb,
 '{"sr": "IT usluge online", "en": "IT services online", "ru": "IT услуги онлайн"}'::jsonb,
 '💻', true),

-- =============================================================================
-- L2 for: 18. Ostalo (Other) - 5 categories
-- =============================================================================

('kolekcionarstvo', (SELECT id FROM categories WHERE slug = 'ostalo'), 2, 'ostalo/kolekcionarstvo', 1,
 '{"sr": "Kolekcionarstvo", "en": "Collectibles", "ru": "Коллекционирование"}'::jsonb,
 '{"sr": "Antikvarije, novčići, marke, memorabilije", "en": "Antiques, coins, stamps, memorabilia", "ru": "Антиквариат, монеты, марки, памятные вещи"}'::jsonb,
 '{"sr": "Kolekcionarstvo | Vondi", "en": "Collectibles | Vondi", "ru": "Коллекционирование | Vondi"}'::jsonb,
 '{"sr": "Kupite kolekcionarske predmete online", "en": "Buy collectibles online", "ru": "Купить коллекционные предметы онлайн"}'::jsonb,
 '🏺', true),

('vintage', (SELECT id FROM categories WHERE slug = 'ostalo'), 2, 'ostalo/vintage', 2,
 '{"sr": "Vintage", "en": "Vintage", "ru": "Винтаж"}'::jsonb,
 '{"sr": "Retro odeća, stari predmeti, vintage dodaci", "en": "Retro clothing, old items, vintage accessories", "ru": "Ретро одежда, старые предметы, винтажные аксессуары"}'::jsonb,
 '{"sr": "Vintage | Vondi", "en": "Vintage | Vondi", "ru": "Винтаж | Vondi"}'::jsonb,
 '{"sr": "Kupite vintage predmete online", "en": "Buy vintage items online", "ru": "Купить винтажные предметы онлайн"}'::jsonb,
 '📻', true),

('pokloni-i-suveniri', (SELECT id FROM categories WHERE slug = 'ostalo'), 2, 'ostalo/pokloni-i-suveniri', 3,
 '{"sr": "Pokloni i suveniri", "en": "Gifts & souvenirs", "ru": "Подарки и сувениры"}'::jsonb,
 '{"sr": "Poklon paketi, suveniri, personalizovani pokloni", "en": "Gift sets, souvenirs, personalized gifts", "ru": "Подарочные наборы, сувениры, персонализированные подарки"}'::jsonb,
 '{"sr": "Pokloni i suveniri | Vondi", "en": "Gifts & souvenirs | Vondi", "ru": "Подарки и сувениры | Vondi"}'::jsonb,
 '{"sr": "Kupite poklone i suveniron line", "en": "Buy gifts and souvenirs online", "ru": "Купить подарки и сувениры онлайн"}'::jsonb,
 '🎁', true),

('erotika', (SELECT id FROM categories WHERE slug = 'ostalo'), 2, 'ostalo/erotika', 4,
 '{"sr": "Erotika", "en": "Adult", "ru": "Эротика"}'::jsonb,
 '{"sr": "Erotski proizvodi za odrasle", "en": "Adult products", "ru": "Эротические товары для взрослых"}'::jsonb,
 '{"sr": "Erotika | Vondi", "en": "Adult | Vondi", "ru": "Эротика | Vondi"}'::jsonb,
 '{"sr": "Kupite erotske proizvode online", "en": "Buy adult products online", "ru": "Купить эротические товары онлайн"}'::jsonb,
 '🔞', true),

('razno', (SELECT id FROM categories WHERE slug = 'ostalo'), 2, 'ostalo/razno', 5,
 '{"sr": "Razno", "en": "Miscellaneous", "ru": "Разное"}'::jsonb,
 '{"sr": "Ostali proizvodi i usluge", "en": "Other products and services", "ru": "Прочие товары и услуги"}'::jsonb,
 '{"sr": "Razno | Vondi", "en": "Miscellaneous | Vondi", "ru": "Разное | Vondi"}'::jsonb,
 '{"sr": "Kupite razne proizvode online", "en": "Buy miscellaneous products online", "ru": "Купить разные товары онлайн"}'::jsonb,
 '📦', true);

-- =============================================================================
-- Final verification and summary
-- =============================================================================
DO $$
DECLARE
    l2_count INTEGER;
    l2_by_parent RECORD;
BEGIN
    SELECT COUNT(*) INTO l2_count FROM categories WHERE level = 2;

    RAISE NOTICE '=== L2 Categories Migration Complete ===';
    RAISE NOTICE 'Total L2 categories inserted: %', l2_count;
    RAISE NOTICE '';
    RAISE NOTICE 'L2 categories by parent (L1):';
    
    FOR l2_by_parent IN 
        SELECT 
            p.slug as parent_slug,
            p.name->>'sr' as parent_name,
            COUNT(c.id) as l2_count
        FROM categories p
        LEFT JOIN categories c ON c.parent_id = p.id AND c.level = 2
        WHERE p.level = 1
        GROUP BY p.id, p.slug, p.name
        ORDER BY p.sort_order
    LOOP
        RAISE NOTICE '  % (%) = % L2 categories', 
            l2_by_parent.parent_slug, 
            l2_by_parent.parent_name, 
            l2_by_parent.l2_count;
    END LOOP;

    IF l2_count < 190 THEN
        RAISE WARNING 'Expected at least 190 L2 categories, but found %. Some categories may be missing.', l2_count;
    ELSE
        RAISE NOTICE '';
        RAISE NOTICE 'SUCCESS: L2 categories seed data complete!';
    END IF;
END $$;
