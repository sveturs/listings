-- Categories fixture data
-- Root categories (level 1) and subcategories (level 2, 3)

-- Clean up existing fixture data first
DELETE FROM categories WHERE id IN (
    '11111111-1111-1111-1111-111111111001'::uuid,
    '11111111-1111-1111-1111-111111111002'::uuid,
    '11111111-1111-1111-1111-111111111003'::uuid,
    '11111111-1111-1111-1111-111111111004'::uuid,
    '11111111-1111-1111-1111-111111112001'::uuid,
    '11111111-1111-1111-1111-111111112002'::uuid,
    '11111111-1111-1111-1111-111111112003'::uuid,
    '11111111-1111-1111-1111-111111113001'::uuid,
    '11111111-1111-1111-1111-111111113002'::uuid,
    '11111111-1111-1111-1111-111111113003'::uuid,
    '11111111-1111-1111-1111-111111114001'::uuid,
    '11111111-1111-1111-1111-111111114002'::uuid,
    '11111111-1111-1111-1111-111111114003'::uuid
);

-- Electronics (root)
INSERT INTO categories (id, slug, parent_id, level, path, sort_order, name, description, is_active)
VALUES
    ('11111111-1111-1111-1111-111111111001'::uuid, 'electronics', NULL, 1, '/electronics', 1,
     '{"en": "Electronics", "sr": "Elektronika", "ru": "Электроника"}',
     '{"en": "Electronic devices and gadgets", "sr": "Elektronski uređaji i gadžeti", "ru": "Электронные устройства и гаджеты"}',
     true);

-- Electronics subcategories (level 2)
INSERT INTO categories (id, slug, parent_id, level, path, sort_order, name, description, is_active)
VALUES
    ('11111111-1111-1111-1111-111111111002'::uuid, 'laptops', '11111111-1111-1111-1111-111111111001'::uuid, 2, '/electronics/laptops', 1,
     '{"en": "Laptops", "sr": "Laptopovi", "ru": "Ноутбуки"}',
     '{"en": "Portable computers", "sr": "Prenosivi računari", "ru": "Портативные компьютеры"}',
     true),
    ('11111111-1111-1111-1111-111111111003'::uuid, 'smartphones', '11111111-1111-1111-1111-111111111001'::uuid, 2, '/electronics/smartphones', 2,
     '{"en": "Smartphones", "sr": "Pametni telefoni", "ru": "Смартфоны"}',
     '{"en": "Mobile phones", "sr": "Mobilni telefoni", "ru": "Мобильные телефоны"}',
     true),
    ('11111111-1111-1111-1111-111111111004'::uuid, 'tablets', '11111111-1111-1111-1111-111111111001'::uuid, 2, '/electronics/tablets', 3,
     '{"en": "Tablets", "sr": "Tableti", "ru": "Планшеты"}',
     '{"en": "Tablet devices", "sr": "Tablet uređaji", "ru": "Планшетные устройства"}',
     true);

-- Home & Garden (root)
INSERT INTO categories (id, slug, parent_id, level, path, sort_order, name, description, is_active)
VALUES
    ('11111111-1111-1111-1111-111111112001'::uuid, 'home-garden', NULL, 1, '/home-garden', 2,
     '{"en": "Home & Garden", "sr": "Dom i bašta", "ru": "Дом и сад"}',
     '{"en": "Home and garden products", "sr": "Proizvodi za dom i baštu", "ru": "Товары для дома и сада"}',
     true);

-- Home subcategories (level 2)
INSERT INTO categories (id, slug, parent_id, level, path, sort_order, name, description, is_active)
VALUES
    ('11111111-1111-1111-1111-111111112002'::uuid, 'furniture', '11111111-1111-1111-1111-111111112001'::uuid, 2, '/home-garden/furniture', 1,
     '{"en": "Furniture", "sr": "Nameštaj", "ru": "Мебель"}',
     '{"en": "Home furniture", "sr": "Kućni nameštaj", "ru": "Домашняя мебель"}',
     true),
    ('11111111-1111-1111-1111-111111112003'::uuid, 'garden-tools', '11111111-1111-1111-1111-111111112001'::uuid, 2, '/home-garden/garden-tools', 2,
     '{"en": "Garden Tools", "sr": "Baštenske alatke", "ru": "Садовые инструменты"}',
     '{"en": "Tools for gardening", "sr": "Alati za baštovanstvo", "ru": "Инструменты для садоводства"}',
     true);

-- Automotive (root)
INSERT INTO categories (id, slug, parent_id, level, path, sort_order, name, description, is_active)
VALUES
    ('11111111-1111-1111-1111-111111113001'::uuid, 'automotive', NULL, 1, '/automotive', 3,
     '{"en": "Automotive", "sr": "Automobili", "ru": "Автомобили"}',
     '{"en": "Cars and automotive parts", "sr": "Automobili i auto delovi", "ru": "Автомобили и запчасти"}',
     true);

-- Automotive subcategories (level 2)
INSERT INTO categories (id, slug, parent_id, level, path, sort_order, name, description, is_active)
VALUES
    ('11111111-1111-1111-1111-111111113002'::uuid, 'car-parts', '11111111-1111-1111-1111-111111113001'::uuid, 2, '/automotive/car-parts', 1,
     '{"en": "Car Parts", "sr": "Auto delovi", "ru": "Автозапчасти"}',
     '{"en": "Spare parts for cars", "sr": "Rezervni delovi za automobile", "ru": "Запасные части для автомобилей"}',
     true),
    ('11111111-1111-1111-1111-111111113003'::uuid, 'tires', '11111111-1111-1111-1111-111111113001'::uuid, 2, '/automotive/tires', 2,
     '{"en": "Tires", "sr": "Gume", "ru": "Шины"}',
     '{"en": "Car tires", "sr": "Automobilske gume", "ru": "Автомобильные шины"}',
     true);

-- Clothing (root)
INSERT INTO categories (id, slug, parent_id, level, path, sort_order, name, description, is_active)
VALUES
    ('11111111-1111-1111-1111-111111114001'::uuid, 'clothing', NULL, 1, '/clothing', 4,
     '{"en": "Clothing", "sr": "Odeća", "ru": "Одежда"}',
     '{"en": "Clothes and fashion", "sr": "Odeća i moda", "ru": "Одежда и мода"}',
     true);

-- Clothing subcategories (level 2)
INSERT INTO categories (id, slug, parent_id, level, path, sort_order, name, description, is_active)
VALUES
    ('11111111-1111-1111-1111-111111114002'::uuid, 'mens-clothing', '11111111-1111-1111-1111-111111114001'::uuid, 2, '/clothing/mens-clothing', 1,
     '{"en": "Men''s Clothing", "sr": "Muška odeća", "ru": "Мужская одежда"}',
     '{"en": "Clothing for men", "sr": "Odeća za muškarce", "ru": "Одежда для мужчин"}',
     true),
    ('11111111-1111-1111-1111-111111114003'::uuid, 'womens-clothing', '11111111-1111-1111-1111-111111114001'::uuid, 2, '/clothing/womens-clothing', 2,
     '{"en": "Women''s Clothing", "sr": "Ženska odeća", "ru": "Женская одежда"}',
     '{"en": "Clothing for women", "sr": "Odeća za žene", "ru": "Одежда для женщин"}',
     true);
