-- Добавление категории "Принтеры и сканеры" в раздел "Электроника"

-- Добавляем категорию "Office Equipment" как подкатегорию Electronics
INSERT INTO marketplace_categories (id, parent_id, name, slug, description, sort_order, created_at)
VALUES
    (1109, 1001, 'Office Equipment', 'office-equipment', 'Printers, scanners, and office devices', 9, NOW());

-- Добавляем подкатегории
INSERT INTO marketplace_categories (id, parent_id, name, slug, description, sort_order, created_at)
VALUES
    (2007, 1109, 'Printers & Scanners', 'printers-scanners', 'Printers, scanners, and multifunction devices', 1, NOW()),
    (2008, 1109, 'Office Supplies', 'office-supplies', 'Paper, ink, toner, and other supplies', 2, NOW());

-- Добавляем переводы для новых категорий
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, created_at, updated_at)
VALUES
    -- Office Equipment
    ('category', 1109, 'name', 'ru', 'Офисная техника', NOW(), NOW()),
    ('category', 1109, 'name', 'en', 'Office Equipment', NOW(), NOW()),
    ('category', 1109, 'name', 'sr', 'Канцеларијска опрема', NOW(), NOW()),
    ('category', 1109, 'description', 'ru', 'Принтеры, сканеры и офисные устройства', NOW(), NOW()),
    ('category', 1109, 'description', 'en', 'Printers, scanners, and office devices', NOW(), NOW()),
    ('category', 1109, 'description', 'sr', 'Штампачи, скенери и канцеларијски уређаји', NOW(), NOW()),

    -- Printers & Scanners
    ('category', 2007, 'name', 'ru', 'Принтеры и сканеры', NOW(), NOW()),
    ('category', 2007, 'name', 'en', 'Printers & Scanners', NOW(), NOW()),
    ('category', 2007, 'name', 'sr', 'Штампачи и скенери', NOW(), NOW()),
    ('category', 2007, 'description', 'ru', 'Принтеры, сканеры и МФУ', NOW(), NOW()),
    ('category', 2007, 'description', 'en', 'Printers, scanners, and multifunction devices', NOW(), NOW()),
    ('category', 2007, 'description', 'sr', 'Штампачи, скенери и мултифункционални уређаји', NOW(), NOW()),

    -- Office Supplies
    ('category', 2008, 'name', 'ru', 'Расходные материалы', NOW(), NOW()),
    ('category', 2008, 'name', 'en', 'Office Supplies', NOW(), NOW()),
    ('category', 2008, 'name', 'sr', 'Канцеларијски материјал', NOW(), NOW()),
    ('category', 2008, 'description', 'ru', 'Бумага, чернила, тонер и другие расходники', NOW(), NOW()),
    ('category', 2008, 'description', 'en', 'Paper, ink, toner, and other supplies', NOW(), NOW()),
    ('category', 2008, 'description', 'sr', 'Папир, мастило, тонер и други потрошни материјал', NOW(), NOW());

-- Добавляем ключевые слова для категории принтеров
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source, created_at, updated_at)
VALUES
    -- Английские ключевые слова
    (2007, 'printer', 'en', 100.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'scanner', 'en', 95.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'mfu', 'en', 90.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'multifunction', 'en', 85.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'inkjet', 'en', 80.0, 'secondary', 'system', NOW(), NOW()),
    (2007, 'laser', 'en', 80.0, 'secondary', 'system', NOW(), NOW()),
    (2007, 'copier', 'en', 75.0, 'secondary', 'system', NOW(), NOW()),
    (2007, 'fax', 'en', 70.0, 'secondary', 'system', NOW(), NOW()),
    (2007, 'print', 'en', 85.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'scan', 'en', 85.0, 'primary', 'system', NOW(), NOW()),

    -- Русские ключевые слова
    (2007, 'принтер', 'ru', 100.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'сканер', 'ru', 95.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'мфу', 'ru', 90.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'струйный', 'ru', 80.0, 'secondary', 'system', NOW(), NOW()),
    (2007, 'лазерный', 'ru', 80.0, 'secondary', 'system', NOW(), NOW()),
    (2007, 'копир', 'ru', 75.0, 'secondary', 'system', NOW(), NOW()),
    (2007, 'факс', 'ru', 70.0, 'secondary', 'system', NOW(), NOW()),
    (2007, 'печать', 'ru', 85.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'сканирование', 'ru', 85.0, 'primary', 'system', NOW(), NOW()),

    -- Сербские ключевые слова
    (2007, 'štampač', 'sr', 100.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'skener', 'sr', 95.0, 'primary', 'system', NOW(), NOW()),
    (2007, 'multifunkcijski', 'sr', 90.0, 'primary', 'system', NOW(), NOW()),

    -- Бренды принтеров с уменьшенным весом для Canon (чтобы не путать с камерами)
    (2007, 'canon', 'en', 50.0, 'brand', 'system', NOW(), NOW()),
    (2007, 'hp', 'en', 60.0, 'brand', 'system', NOW(), NOW()),
    (2007, 'epson', 'en', 60.0, 'brand', 'system', NOW(), NOW()),
    (2007, 'brother', 'en', 60.0, 'brand', 'system', NOW(), NOW()),
    (2007, 'xerox', 'en', 60.0, 'brand', 'system', NOW(), NOW()),
    (2007, 'lexmark', 'en', 60.0, 'brand', 'system', NOW(), NOW()),
    (2007, 'samsung', 'en', 50.0, 'brand', 'system', NOW(), NOW());

-- Обновляем последний товар с правильной категорией
UPDATE marketplace_listings
SET category_id = 2007
WHERE id = 323;