-- Migration: L3 Ostalo Categories (64 new L3)
-- Parent: Lepota, Bebe, Auto, Aparati, Knjige, Ljubimci
-- Total L3 after: 286 + 64 = 350 FINAL

-- ============================================================================
-- KOSA I FRIZURA (parent: kosa-i-frizura) - 10 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Hair dryers
('fen-za-kosu',
 '{"sr": "Fen za kosu", "en": "Hair Dryer", "ru": "Фен для волос"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/fen-za-kosu', 1, true, NOW(), NOW()),

-- Straighteners
('peglazakosu',
 '{"sr": "Pegla za kosu", "en": "Hair Straightener", "ru": "Утюжок для волос"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/peglazakosu', 2, true, NOW(), NOW()),

-- Curlers
('figaro',
 '{"sr": "Figaro", "en": "Curling Iron", "ru": "Плойка для завивки"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/figaro', 3, true, NOW(), NOW()),

-- Trimmers
('trimer-za-kosu',
 '{"sr": "Trimer za kosu", "en": "Hair Trimmer", "ru": "Триммер для волос"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/trimer-za-kosu', 4, true, NOW(), NOW()),

('masnica-za-kosu',
 '{"sr": "Mašinica za kosu", "en": "Hair Clipper", "ru": "Машинка для стрижки"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/masnica-za-kosu', 5, true, NOW(), NOW()),

-- Hair care
('sampon',
 '{"sr": "Šampon", "en": "Shampoo", "ru": "Шампунь"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/sampon', 6, true, NOW(), NOW()),

('regenerator-za-kosu',
 '{"sr": "Regenerator za kosu", "en": "Hair Conditioner", "ru": "Кондиционер для волос"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/regenerator-za-kosu', 7, true, NOW(), NOW()),

('maska-za-kosu',
 '{"sr": "Maska za kosu", "en": "Hair Mask", "ru": "Маска для волос"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/maska-za-kosu', 8, true, NOW(), NOW()),

-- Hair coloring
('farba-za-kosu',
 '{"sr": "Farba za kosu", "en": "Hair Dye", "ru": "Краска для волос"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/farba-za-kosu', 9, true, NOW(), NOW()),

-- Styling
('sprej-za-kosu',
 '{"sr": "Sprej za kosu", "en": "Hair Spray", "ru": "Лак для волос"}',
 (SELECT id FROM categories WHERE slug = 'kosa-i-frizura' AND level = 2),
 3, 'lepota/kosa-i-frizura/sprej-za-kosu', 10, true, NOW(), NOW());

-- ============================================================================
-- NEGA KOZE (parent: nega-koze) - 10 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Cleansers
('mleko-za-ciscenje',
 '{"sr": "Mleko za čišćenje", "en": "Cleansing Milk", "ru": "Молочко для очищения"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/mleko-za-ciscenje', 1, true, NOW(), NOW()),

('tonik-za-lice',
 '{"sr": "Tonik za lice", "en": "Face Toner", "ru": "Тоник для лица"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/tonik-za-lice', 2, true, NOW(), NOW()),

-- Moisturizers
('krema-za-lice',
 '{"sr": "Krema za lice", "en": "Face Cream", "ru": "Крем для лица"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/krema-za-lice', 3, true, NOW(), NOW()),

('serum-za-lice',
 '{"sr": "Serum za lice", "en": "Face Serum", "ru": "Сыворотка для лица"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/serum-za-lice', 4, true, NOW(), NOW()),

-- Masks
('maska-za-lice',
 '{"sr": "Maska za lice", "en": "Face Mask", "ru": "Маска для лица"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/maska-za-lice', 5, true, NOW(), NOW()),

-- Sun protection
('krema-za-suncanje',
 '{"sr": "Krema za sunčanje", "en": "Sunscreen", "ru": "Солнцезащитный крем"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/krema-za-suncanje', 6, true, NOW(), NOW()),

-- Anti-aging
('anti-age-krema',
 '{"sr": "Anti-age krema", "en": "Anti-Aging Cream", "ru": "Антивозрастной крем"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/anti-age-krema', 7, true, NOW(), NOW()),

-- Body care
('krema-za-telo',
 '{"sr": "Krema za telo", "en": "Body Lotion", "ru": "Крем для тела"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/krema-za-telo', 8, true, NOW(), NOW()),

('gel-za-tusiranje',
 '{"sr": "Gel za tuširanje", "en": "Shower Gel", "ru": "Гель для душа"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/gel-za-tusiranje', 9, true, NOW(), NOW()),

-- Hand care
('krema-za-ruke',
 '{"sr": "Krema za ruke", "en": "Hand Cream", "ru": "Крем для рук"}',
 (SELECT id FROM categories WHERE slug = 'nega-koze' AND level = 2),
 3, 'lepota/nega-koze/krema-za-ruke', 10, true, NOW(), NOW());

-- ============================================================================
-- BEBE OPREMA (parent: bebe-oprema) - 12 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Strollers
('kolica-za-bebe',
 '{"sr": "Kolica za bebe", "en": "Baby Stroller", "ru": "Детская коляска"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/kolica-za-bebe', 1, true, NOW(), NOW()),

('nosiljka-za-bebe',
 '{"sr": "Nosiljka za bebe", "en": "Baby Carrier", "ru": "Слинг"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/nosiljka-za-bebe', 2, true, NOW(), NOW()),

-- Car seats
('auto-sediste-za-bebe',
 '{"sr": "Auto sedište za bebe", "en": "Baby Car Seat", "ru": "Автокресло для младенцев"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/auto-sediste-za-bebe', 3, true, NOW(), NOW()),

-- Cribs
('krevetac-za-bebe',
 '{"sr": "Krevetac za bebe", "en": "Baby Crib", "ru": "Детская кроватка"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/krevetac-za-bebe', 4, true, NOW(), NOW()),

-- High chairs
('hranilica',
 '{"sr": "Hranilica", "en": "High Chair", "ru": "Стульчик для кормления"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/hranilica', 5, true, NOW(), NOW()),

-- Baby monitors
('baby-alarm',
 '{"sr": "Baby alarm", "en": "Baby Monitor", "ru": "Радионяня"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/baby-alarm', 6, true, NOW(), NOW()),

-- Bottles
('bocice-za-bebe',
 '{"sr": "Bočice za bebe", "en": "Baby Bottles", "ru": "Детские бутылочки"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/bocice-za-bebe', 7, true, NOW(), NOW()),

-- Pumps
('pumpa-za-mleko',
 '{"sr": "Pumpa za mleko", "en": "Breast Pump", "ru": "Молокоотсос"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/pumpa-za-mleko', 8, true, NOW(), NOW()),

-- Diapers
('pelene',
 '{"sr": "Pelene", "en": "Diapers", "ru": "Подгузники"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/pelene', 9, true, NOW(), NOW()),

-- Bath
('kadica-za-kupanje',
 '{"sr": "Kadica za kupanje", "en": "Baby Bath Tub", "ru": "Ванночка для купания"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/kadica-za-kupanje', 10, true, NOW(), NOW()),

-- Toys
('bebe-igracke',
 '{"sr": "Bebe igračke", "en": "Baby Toys", "ru": "Игрушки для младенцев"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/bebe-igracke', 11, true, NOW(), NOW()),

-- Playpen
('ogradica-za-bebe',
 '{"sr": "Ogradica za bebe", "en": "Baby Playpen", "ru": "Манеж"}',
 (SELECT id FROM categories WHERE slug = 'bebe-oprema' AND level = 2),
 3, 'bebe/bebe-oprema/ogradica-za-bebe', 12, true, NOW(), NOW());

-- ============================================================================
-- AUTO DELOVI (parent: auto-delovi) - 10 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Tires
('gume-ljetne',
 '{"sr": "Gume letnje", "en": "Summer Tires", "ru": "Летние шины"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/gume-ljetne', 1, true, NOW(), NOW()),

('gume-zimske',
 '{"sr": "Gume zimske", "en": "Winter Tires", "ru": "Зимние шины"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/gume-zimske', 2, true, NOW(), NOW()),

-- Wheels
('aluminijske-felne',
 '{"sr": "Aluminijske felne", "en": "Alloy Wheels", "ru": "Литые диски"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/aluminijske-felne', 3, true, NOW(), NOW()),

-- Batteries
('akumulator',
 '{"sr": "Akumulator", "en": "Car Battery", "ru": "Автомобильный аккумулятор"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/akumulator', 4, true, NOW(), NOW()),

-- Lights
('auto-farovi',
 '{"sr": "Auto farovi", "en": "Headlights", "ru": "Автомобильные фары"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/auto-farovi', 5, true, NOW(), NOW()),

('led-auto-sijalice',
 '{"sr": "LED auto sijalice", "en": "LED Car Bulbs", "ru": "LED автолампы"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/led-auto-sijalice', 6, true, NOW(), NOW()),

-- Brakes
('kocione-plocice',
 '{"sr": "Kočione pločice", "en": "Brake Pads", "ru": "Тормозные колодки"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/kocione-plocice', 7, true, NOW(), NOW()),

-- Filters
('uljni-filter',
 '{"sr": "Uljni filter", "en": "Oil Filter", "ru": "Масляный фильтр"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/uljni-filter', 8, true, NOW(), NOW()),

('vazdusni-filter',
 '{"sr": "Vazdušni filter", "en": "Air Filter", "ru": "Воздушный фильтр"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/vazdusni-filter', 9, true, NOW(), NOW()),

-- Wipers
('brisaci',
 '{"sr": "Brisači", "en": "Windshield Wipers", "ru": "Стеклоочистители"}',
 (SELECT id FROM categories WHERE slug = 'auto-delovi' AND level = 2),
 3, 'auto/auto-delovi/brisaci', 10, true, NOW(), NOW());

-- ============================================================================
-- APARATI ZA KUCU (parent: aparati-za-kucu) - 10 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Washing machines
('masina-za-pranje-vesa',
 '{"sr": "Mašina za pranje veša", "en": "Washing Machine", "ru": "Стиральная машина"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/masina-za-pranje-vesa', 1, true, NOW(), NOW()),

('masina-za-susenje',
 '{"sr": "Mašina za sušenje", "en": "Dryer", "ru": "Сушильная машина"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/masina-za-susenje', 2, true, NOW(), NOW()),

-- Vacuum cleaners
('usisivac',
 '{"sr": "Usisivač", "en": "Vacuum Cleaner", "ru": "Пылесос"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/usisivac', 3, true, NOW(), NOW()),

('robotski-usisivac',
 '{"sr": "Robotski usisivač", "en": "Robot Vacuum", "ru": "Робот-пылесос"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/robotski-usisivac', 4, true, NOW(), NOW()),

-- Irons
('pegla',
 '{"sr": "Pegla", "en": "Iron", "ru": "Утюг"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/pegla', 5, true, NOW(), NOW()),

('pegla-sa-parom',
 '{"sr": "Pegla sa parom", "en": "Steam Iron", "ru": "Утюг с паром"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/pegla-sa-parom', 6, true, NOW(), NOW()),

-- Air conditioners
('klima-uredjaj',
 '{"sr": "Klima uređaj", "en": "Air Conditioner", "ru": "Кондиционер"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/klima-uredjaj', 7, true, NOW(), NOW()),

-- Heaters
('grejalica',
 '{"sr": "Grejalica", "en": "Heater", "ru": "Обогреватель"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/grejalica', 8, true, NOW(), NOW()),

-- Air purifiers
('preciscivac-vazduha',
 '{"sr": "Prečišćivač vazduha", "en": "Air Purifier", "ru": "Очиститель воздуха"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/preciscivac-vazduha', 9, true, NOW(), NOW()),

-- Dehumidifiers
('razvlaživač',
 '{"sr": "Razvlaživač", "en": "Dehumidifier", "ru": "Осушитель воздуха"}',
 (SELECT id FROM categories WHERE slug = 'aparati-za-kucu' AND level = 2),
 3, 'aparati/aparati-za-kucu/razvlazivac', 10, true, NOW(), NOW());

-- ============================================================================
-- KNJIGE (parent: knjige) - 6 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Fiction
('romani',
 '{"sr": "Romani", "en": "Novels", "ru": "Романы"}',
 (SELECT id FROM categories WHERE slug = 'knjige' AND level = 2),
 3, 'knjige/knjige/romani', 1, true, NOW(), NOW()),

('price',
 '{"sr": "Priče", "en": "Short Stories", "ru": "Рассказы"}',
 (SELECT id FROM categories WHERE slug = 'knjige' AND level = 2),
 3, 'knjige/knjige/price', 2, true, NOW(), NOW()),

-- Non-fiction
('biografije',
 '{"sr": "Biografije", "en": "Biographies", "ru": "Биографии"}',
 (SELECT id FROM categories WHERE slug = 'knjige' AND level = 2),
 3, 'knjige/knjige/biografije', 3, true, NOW(), NOW()),

('popularna-nauka',
 '{"sr": "Popularna nauka", "en": "Popular Science", "ru": "Научно-популярная литература"}',
 (SELECT id FROM categories WHERE slug = 'knjige' AND level = 2),
 3, 'knjige/knjige/popularna-nauka', 4, true, NOW(), NOW()),

-- Kids
('decije-knjige',
 '{"sr": "Dečije knjige", "en": "Children Books", "ru": "Детские книги"}',
 (SELECT id FROM categories WHERE slug = 'knjige' AND level = 2),
 3, 'knjige/knjige/decije-knjige', 5, true, NOW(), NOW()),

-- Comics
('stripovi',
 '{"sr": "Stripovi", "en": "Comics", "ru": "Комиксы"}',
 (SELECT id FROM categories WHERE slug = 'knjige' AND level = 2),
 3, 'knjige/knjige/stripovi', 6, true, NOW(), NOW());

-- ============================================================================
-- HRANA ZA LJUBIMCE (parent: hrana-za-ljubimce) - 6 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Dog food
('hrana-za-pse-suva',
 '{"sr": "Hrana za pse suva", "en": "Dry Dog Food", "ru": "Сухой корм для собак"}',
 (SELECT id FROM categories WHERE slug = 'hrana-za-ljubimce' AND level = 2),
 3, 'ljubimci/hrana-za-ljubimce/hrana-za-pse-suva', 1, true, NOW(), NOW()),

('hrana-za-pse-konzerve',
 '{"sr": "Hrana za pse konzerve", "en": "Wet Dog Food", "ru": "Влажный корм для собак"}',
 (SELECT id FROM categories WHERE slug = 'hrana-za-ljubimce' AND level = 2),
 3, 'ljubimci/hrana-za-ljubimce/hrana-za-pse-konzerve', 2, true, NOW(), NOW()),

-- Cat food
('hrana-za-macke-suva',
 '{"sr": "Hrana za mačke suva", "en": "Dry Cat Food", "ru": "Сухой корм для кошек"}',
 (SELECT id FROM categories WHERE slug = 'hrana-za-ljubimce' AND level = 2),
 3, 'ljubimci/hrana-za-ljubimce/hrana-za-macke-suva', 3, true, NOW(), NOW()),

('hrana-za-macke-konzerve',
 '{"sr": "Hrana za mačke konzerve", "en": "Wet Cat Food", "ru": "Влажный корм для кошек"}',
 (SELECT id FROM categories WHERE slug = 'hrana-za-ljubimce' AND level = 2),
 3, 'ljubimci/hrana-za-ljubimce/hrana-za-macke-konzerve', 4, true, NOW(), NOW()),

-- Birds
('hrana-za-ptice',
 '{"sr": "Hrana za ptice", "en": "Bird Food", "ru": "Корм для птиц"}',
 (SELECT id FROM categories WHERE slug = 'hrana-za-ljubimce' AND level = 2),
 3, 'ljubimci/hrana-za-ljubimce/hrana-za-ptice', 5, true, NOW(), NOW()),

-- Aquarium
('hrana-za-ribe',
 '{"sr": "Hrana za ribe", "en": "Fish Food", "ru": "Корм для рыб"}',
 (SELECT id FROM categories WHERE slug = 'hrana-za-ljubimce' AND level = 2),
 3, 'ljubimci/hrana-za-ljubimce/hrana-za-ribe', 6, true, NOW(), NOW());

-- ============================================================================
-- Migration Summary
-- ============================================================================
-- Total new L3 categories: 64
-- Categories breakdown:
--   - Kosa i frizura: 10
--   - Nega koze: 10
--   - Bebe oprema: 12
--   - Auto delovi: 10
--   - Aparati za kucu: 10
--   - Knjige: 6
--   - Hrana za ljubimce: 6
-- Total L3 after this migration: 350 FINAL
-- ============================================================================
