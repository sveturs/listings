-- Migration: L3 Dom i Sport Categories (80 new L3)
-- Parent: Dom, Sport categories
-- Total L3 after: 206 + 80 = 286

-- ============================================================================
-- NAMESTAJ (parent: namestaj) - 15 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Living Room
('trosed-dvosed',
 '{"sr": "Trosed i dvosed", "en": "3-Seater and 2-Seater Sofas", "ru": "Трехместные и двухместные диваны"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/trosed-dvosed', 1, true, NOW(), NOW()),

('ugaona-garnitura',
 '{"sr": "Ugaona garnitura", "en": "Corner Sofas", "ru": "Угловые диваны"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/ugaona-garnitura', 2, true, NOW(), NOW()),

('fotelja',
 '{"sr": "Fotelja", "en": "Armchairs", "ru": "Кресла"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/fotelja', 3, true, NOW(), NOW()),

('tv-komoda',
 '{"sr": "TV komoda", "en": "TV Stands", "ru": "Тумбы под TV"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/tv-komoda', 4, true, NOW(), NOW()),

('sto-za-dnevnu-sobu',
 '{"sr": "Sto za dnevnu sobu", "en": "Coffee Tables", "ru": "Журнальные столы"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/sto-za-dnevnu-sobu', 5, true, NOW(), NOW()),

-- Bedroom
('bracni-krevet',
 '{"sr": "Bračni krevet", "en": "Double Bed", "ru": "Двуспальная кровать"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/bracni-krevet', 6, true, NOW(), NOW()),

('samacki-krevet',
 '{"sr": "Samački krevet", "en": "Single Bed", "ru": "Односпальная кровать"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/samacki-krevet', 7, true, NOW(), NOW()),

('sprat-krevet',
 '{"sr": "Sprat krevet", "en": "Bunk Bed", "ru": "Двухъярусная кровать"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/sprat-krevet', 8, true, NOW(), NOW()),

('orman',
 '{"sr": "Orman", "en": "Wardrobe", "ru": "Шкаф"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/orman', 9, true, NOW(), NOW()),

('komode',
 '{"sr": "Komode", "en": "Chest of Drawers", "ru": "Комод"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/komode', 10, true, NOW(), NOW()),

('nocni-stocic',
 '{"sr": "Noćni stočić", "en": "Nightstand", "ru": "Прикроватная тумбочка"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/nocni-stocic', 11, true, NOW(), NOW()),

-- Dining Room
('trpezarijski-sto',
 '{"sr": "Trpezarijski sto", "en": "Dining Table", "ru": "Обеденный стол"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/trpezarijski-sto', 12, true, NOW(), NOW()),

('trpezarijske-stolice',
 '{"sr": "Trpezarijske stolice", "en": "Dining Chairs", "ru": "Обеденные стулья"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/trpezarijske-stolice', 13, true, NOW(), NOW()),

-- Office
('kancelarijski-sto',
 '{"sr": "Kancelarijski sto", "en": "Office Desk", "ru": "Письменный стол"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/kancelarijski-sto', 14, true, NOW(), NOW()),

('kancelarijska-stolica',
 '{"sr": "Kancelarijska stolica", "en": "Office Chair", "ru": "Офисное кресло"}',
 (SELECT id FROM categories WHERE slug = 'namestaj' AND level = 2),
 3, 'dom/namestaj/kancelarijska-stolica', 15, true, NOW(), NOW());

-- ============================================================================
-- KUPATILO (parent: kupatilo) - 12 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Sinks
('lavabo',
 '{"sr": "Lavabo", "en": "Bathroom Sink", "ru": "Раковина"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/lavabo', 1, true, NOW(), NOW()),

('ugradni-lavabo',
 '{"sr": "Ugradni lavabo", "en": "Built-in Sink", "ru": "Встраиваемая раковина"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/ugradni-lavabo', 2, true, NOW(), NOW()),

-- Bathtubs
('kada',
 '{"sr": "Kada", "en": "Bathtub", "ru": "Ванна"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/kada', 3, true, NOW(), NOW()),

('hidromasazna-kada',
 '{"sr": "Hidromasažna kada", "en": "Whirlpool Bathtub", "ru": "Гидромассажная ванна"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/hidromasazna-kada', 4, true, NOW(), NOW()),

-- Showers
('tus-kabina',
 '{"sr": "Tuš kabina", "en": "Shower Cabin", "ru": "Душевая кабина"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/tus-kabina', 5, true, NOW(), NOW()),

('tus-set',
 '{"sr": "Tuš set", "en": "Shower Set", "ru": "Душевой набор"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/tus-set', 6, true, NOW(), NOW()),

-- Toilets
('wc-solja',
 '{"sr": "WC šolja", "en": "Toilet Bowl", "ru": "Унитаз"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/wc-solja', 7, true, NOW(), NOW()),

('bidе',
 '{"sr": "Bide", "en": "Bidet", "ru": "Биде"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/bide', 8, true, NOW(), NOW()),

-- Faucets
('slavina-za-lavabo',
 '{"sr": "Slavina za lavabo", "en": "Sink Faucet", "ru": "Смеситель для раковины"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/slavina-za-lavabo', 9, true, NOW(), NOW()),

('slavina-za-kadu',
 '{"sr": "Slavina za kadu", "en": "Bath Faucet", "ru": "Смеситель для ванны"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/slavina-za-kadu', 10, true, NOW(), NOW()),

-- Cabinets
('ogledalo-sa-ormarićem',
 '{"sr": "Ogledalo sa ormarićem", "en": "Mirror Cabinet", "ru": "Зеркальный шкаф"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/ogledalo-sa-ormarićem', 11, true, NOW(), NOW()),

('kupatilski-namestaj',
 '{"sr": "Kupatilski nameštaj", "en": "Bathroom Furniture", "ru": "Мебель для ванной"}',
 (SELECT id FROM categories WHERE slug = 'kupatilo' AND level = 2),
 3, 'dom/kupatilo/kupatilski-namestaj', 12, true, NOW(), NOW());

-- ============================================================================
-- KUHINJA (parent: kuhinja) - 13 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Kitchen units
('kuhinjski-elementi',
 '{"sr": "Kuhinjski elementi", "en": "Kitchen Units", "ru": "Кухонные шкафы"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/kuhinjski-elementi', 1, true, NOW(), NOW()),

('kuhinjska-radna-ploca',
 '{"sr": "Kuhinjska radna ploča", "en": "Kitchen Countertop", "ru": "Кухонная столешница"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/kuhinjska-radna-ploca', 2, true, NOW(), NOW()),

-- Sinks
('kuhinjska-sudopera',
 '{"sr": "Kuhinjska sudopera", "en": "Kitchen Sink", "ru": "Кухонная мойка"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/kuhinjska-sudopera', 3, true, NOW(), NOW()),

('kuhinjska-slavina',
 '{"sr": "Kuhinjska slavina", "en": "Kitchen Faucet", "ru": "Кухонный смеситель"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/kuhinjska-slavina', 4, true, NOW(), NOW()),

-- Appliances
('aspirator',
 '{"sr": "Aspirator", "en": "Range Hood", "ru": "Вытяжка"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/aspirator', 5, true, NOW(), NOW()),

('stednjak',
 '{"sr": "Štednjak", "en": "Stove", "ru": "Плита"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/stednjak', 6, true, NOW(), NOW()),

('ugradna-rerna',
 '{"sr": "Ugradna rerna", "en": "Built-in Oven", "ru": "Встраиваемая духовка"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/ugradna-rerna', 7, true, NOW(), NOW()),

('ploca-za-kuvanje',
 '{"sr": "Ploča za kuvanje", "en": "Cooktop", "ru": "Варочная панель"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/ploca-za-kuvanje', 8, true, NOW(), NOW()),

('masina-za-pranje-sudova',
 '{"sr": "Mašina za pranje sudova", "en": "Dishwasher", "ru": "Посудомоечная машина"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/masina-za-pranje-sudova', 9, true, NOW(), NOW()),

('frizider',
 '{"sr": "Frižider", "en": "Refrigerator", "ru": "Холодильник"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/frizider', 10, true, NOW(), NOW()),

('zamrzivac',
 '{"sr": "Zamrzivač", "en": "Freezer", "ru": "Морозильник"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/zamrzivac', 11, true, NOW(), NOW()),

('mikrotalasna-rerna',
 '{"sr": "Mikrotalasna rerna", "en": "Microwave Oven", "ru": "Микроволновая печь"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/mikrotalasna-rerna', 12, true, NOW(), NOW()),

('mini-bojler',
 '{"sr": "Mini bojler", "en": "Water Heater", "ru": "Водонагреватель"}',
 (SELECT id FROM categories WHERE slug = 'kuhinja' AND level = 2),
 3, 'dom/kuhinja/mini-bojler', 13, true, NOW(), NOW());

-- ============================================================================
-- RASVETA (parent: rasveta) - 10 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Ceiling lights
('luster',
 '{"sr": "Luster", "en": "Chandelier", "ru": "Люстра"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/luster', 1, true, NOW(), NOW()),

('plafonjera',
 '{"sr": "Plafonjera", "en": "Ceiling Lamp", "ru": "Потолочный светильник"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/plafonjera', 2, true, NOW(), NOW()),

('ugradna-led-rasveta',
 '{"sr": "Ugradna LED rasveta", "en": "Built-in LED Lights", "ru": "Встраиваемые LED светильники"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/ugradna-led-rasveta', 3, true, NOW(), NOW()),

-- Wall lights
('zidna-lampa',
 '{"sr": "Zidna lampa", "en": "Wall Lamp", "ru": "Настенный светильник"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/zidna-lampa', 4, true, NOW(), NOW()),

-- Floor lamps
('stojeća-lampa',
 '{"sr": "Stojeća lampa", "en": "Floor Lamp", "ru": "Торшер"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/stojeca-lampa', 5, true, NOW(), NOW()),

-- Table lamps
('stonalampa',
 '{"sr": "Stona lampa", "en": "Table Lamp", "ru": "Настольная лампа"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/stonalampa', 6, true, NOW(), NOW()),

-- Outdoor lights
('spoljna-rasveta',
 '{"sr": "Spoljna rasveta", "en": "Outdoor Lighting", "ru": "Уличное освещение"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/spoljna-rasveta', 7, true, NOW(), NOW()),

-- Bulbs
('led-sijalice',
 '{"sr": "LED sijalice", "en": "LED Bulbs", "ru": "LED лампочки"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/led-sijalice', 8, true, NOW(), NOW()),

-- Smart lights
('pametna-rasveta',
 '{"sr": "Pametna rasveta", "en": "Smart Lighting", "ru": "Умное освещение"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/pametna-rasveta', 9, true, NOW(), NOW()),

-- Decorative lights
('dekorativna-rasveta',
 '{"sr": "Dekorativna rasveta", "en": "Decorative Lighting", "ru": "Декоративное освещение"}',
 (SELECT id FROM categories WHERE slug = 'rasveta' AND level = 2),
 3, 'dom/rasveta/dekorativna-rasveta', 10, true, NOW(), NOW());

-- ============================================================================
-- FITNES OPREMA (parent: fitnes-oprema) - 10 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Cardio
('traka-za-trcanje',
 '{"sr": "Traka za trčanje", "en": "Treadmill", "ru": "Беговая дорожка"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/traka-za-trcanje', 1, true, NOW(), NOW()),

('bicikl-sobni',
 '{"sr": "Bicikl sobni", "en": "Exercise Bike", "ru": "Велотренажер"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/bicikl-sobni', 2, true, NOW(), NOW()),

('elipticni-trenažer',
 '{"sr": "Eliptični trenažer", "en": "Elliptical Trainer", "ru": "Эллиптический тренажер"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/elipticni-trenazer', 3, true, NOW(), NOW()),

('veslo-masina',
 '{"sr": "Veslo mašina", "en": "Rowing Machine", "ru": "Гребной тренажер"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/veslo-masina', 4, true, NOW(), NOW()),

-- Weights
('bučice',
 '{"sr": "Bučice", "en": "Dumbbells", "ru": "Гантели"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/bucice', 5, true, NOW(), NOW()),

('sipka-i-tegovi',
 '{"sr": "Šipka i tegovi", "en": "Barbell and Weights", "ru": "Штанга и диски"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/sipka-i-tegovi', 6, true, NOW(), NOW()),

-- Benches
('klupa-za-vezbanje',
 '{"sr": "Klupa za vežbanje", "en": "Exercise Bench", "ru": "Скамья для тренировок"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/klupa-za-vezbanje', 7, true, NOW(), NOW()),

-- Accessories
('podloga-za-jogu',
 '{"sr": "Podloga za jogu", "en": "Yoga Mat", "ru": "Коврик для йоги"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/podloga-za-jogu', 8, true, NOW(), NOW()),

('lopta-za-pilates',
 '{"sr": "Lopta za pilates", "en": "Pilates Ball", "ru": "Мяч для пилатеса"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/lopta-za-pilates', 9, true, NOW(), NOW()),

('trx-trake',
 '{"sr": "TRX trake", "en": "TRX Straps", "ru": "TRX петли"}',
 (SELECT id FROM categories WHERE slug = 'fitnes-oprema' AND level = 2),
 3, 'sport/fitnes-oprema/trx-trake', 10, true, NOW(), NOW());

-- ============================================================================
-- BICIKLI (parent: bicikli) - 10 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Mountain bikes
('brdski-bicikl',
 '{"sr": "Brdski bicikl", "en": "Mountain Bike", "ru": "Горный велосипед"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/brdski-bicikl', 1, true, NOW(), NOW()),

-- Road bikes
('drumski-bicikl',
 '{"sr": "Drumski bicikl", "en": "Road Bike", "ru": "Шоссейный велосипед"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/drumski-bicikl', 2, true, NOW(), NOW()),

-- City bikes
('gradski-bicikl',
 '{"sr": "Gradski bicikl", "en": "City Bike", "ru": "Городской велосипед"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/gradski-bicikl', 3, true, NOW(), NOW()),

-- Electric bikes
('električni-bicikl',
 '{"sr": "Električni bicikl", "en": "E-Bike", "ru": "Электровелосипед"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/elektricni-bicikl', 4, true, NOW(), NOW()),

-- Folding bikes
('sklopivi-bicikl',
 '{"sr": "Sklopivi bicikl", "en": "Folding Bike", "ru": "Складной велосипед"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/sklopivi-bicikl', 5, true, NOW(), NOW()),

-- Kids bikes
('deciji-bicikl',
 '{"sr": "Dečiji bicikl", "en": "Kids Bike", "ru": "Детский велосипед"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/deciji-bicikl', 6, true, NOW(), NOW()),

-- BMX
('bmx-bicikl',
 '{"sr": "BMX bicikl", "en": "BMX Bike", "ru": "BMX велосипед"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/bmx-bicikl', 7, true, NOW(), NOW()),

-- Accessories
('kaciga-za-bicikl',
 '{"sr": "Kaciga za bicikl", "en": "Bike Helmet", "ru": "Велосипедный шлем"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/kaciga-za-bicikl', 8, true, NOW(), NOW()),

('biciklisticke-rukavice',
 '{"sr": "Biciklističke rukavice", "en": "Cycling Gloves", "ru": "Велосипедные перчатки"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/biciklisticke-rukavice', 9, true, NOW(), NOW()),

('biciklistička-torba',
 '{"sr": "Biciklistička torba", "en": "Bike Bag", "ru": "Велосипедная сумка"}',
 (SELECT id FROM categories WHERE slug = 'bicikli' AND level = 2),
 3, 'sport/bicikli/biciklisticka-torba', 10, true, NOW(), NOW());

-- ============================================================================
-- KAMPOVANJE (parent: kampovanje) - 10 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Tents
('sator-2-3-osobe',
 '{"sr": "Šator 2-3 osobe", "en": "Tent 2-3 Person", "ru": "Палатка 2-3 человека"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/sator-2-3-osobe', 1, true, NOW(), NOW()),

('sator-4-6-osoba',
 '{"sr": "Šator 4-6 osoba", "en": "Tent 4-6 Person", "ru": "Палатка 4-6 человек"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/sator-4-6-osoba', 2, true, NOW(), NOW()),

-- Sleeping bags
('spavaca-vreća',
 '{"sr": "Spavaća vreća", "en": "Sleeping Bag", "ru": "Спальный мешок"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/spavaca-vreca', 3, true, NOW(), NOW()),

-- Mattresses
('samoduvajuci-dušek',
 '{"sr": "Samoduvajući dušek", "en": "Self-Inflating Mattress", "ru": "Самонадувающийся коврик"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/samoduvajuci-dusek', 4, true, NOW(), NOW()),

-- Backpacks
('kamp-ranac',
 '{"sr": "Kamp ranac", "en": "Camping Backpack", "ru": "Туристический рюкзак"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/kamp-ranac', 5, true, NOW(), NOW()),

-- Stoves
('kamp-roštilj',
 '{"sr": "Kamp roštilj", "en": "Camping Grill", "ru": "Туристический гриль"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/kamp-rostilj', 6, true, NOW(), NOW()),

('gas-reaud',
 '{"sr": "Gas rešo", "en": "Gas Stove", "ru": "Газовая плитка"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/gas-reaud', 7, true, NOW(), NOW()),

-- Lights
('kamp-lampa',
 '{"sr": "Kamp lampa", "en": "Camping Lantern", "ru": "Кемпинговый фонарь"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/kamp-lampa', 8, true, NOW(), NOW()),

-- Coolers
('prenosivi-frižider',
 '{"sr": "Prenosivi frižider", "en": "Portable Cooler", "ru": "Переносной холодильник"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/prenosivi-frizider', 9, true, NOW(), NOW()),

-- Furniture
('kamp-stolica',
 '{"sr": "Kamp stolica", "en": "Camping Chair", "ru": "Складной стул"}',
 (SELECT id FROM categories WHERE slug = 'kampovanje' AND level = 2),
 3, 'sport/kampovanje/kamp-stolica', 10, true, NOW(), NOW());

-- ============================================================================
-- Migration Summary
-- ============================================================================
-- Total new L3 categories: 80
-- Categories breakdown:
--   - Namestaj: 15
--   - Kupatilo: 12
--   - Kuhinja: 13
--   - Rasveta: 10
--   - Fitnes oprema: 10
--   - Bicikli: 10
--   - Kampovanje: 10
-- Total L3 after this migration: 286
-- ============================================================================
