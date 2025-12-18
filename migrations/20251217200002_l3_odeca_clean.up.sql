-- Migration: L3 Odeca i Obuca Categories (90 new L3)
-- Parent: Moda L1 categories
-- Total L3 after: 116 + 90 = 206

-- ============================================================================
-- MUSKA ODECA (parent: muska-odeca) - 12 new L3
-- Already exists: muske-jakne, muske-kosulje, muske-pantalone
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Sweaters & Hoodies
('muski-dzemperi',
 '{"sr": "Muški džemperi", "en": "Men Sweaters", "ru": "Мужские свитеры"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muski-dzemperi', 4, true, NOW(), NOW()),

('muski-duksevi',
 '{"sr": "Muški duksevi", "en": "Men Hoodies", "ru": "Мужские худи"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muski-duksevi', 5, true, NOW(), NOW()),

-- Sportswear
('muska-sportska-odeca',
 '{"sr": "Muška sportska odeća", "en": "Men Sportswear", "ru": "Мужская спортивная одежда"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muska-sportska-odeca', 6, true, NOW(), NOW()),

('muske-trenerke',
 '{"sr": "Muške trenerke", "en": "Men Tracksuits", "ru": "Мужские спортивные костюмы"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muske-trenerke', 7, true, NOW(), NOW()),

-- Shorts
('muski-sorcevi',
 '{"sr": "Muški šorcevi", "en": "Men Shorts", "ru": "Мужские шорты"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muski-sorcevi', 8, true, NOW(), NOW()),

-- Coats
('muski-kaputi',
 '{"sr": "Muški kaputi", "en": "Men Coats", "ru": "Мужские пальто"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muski-kaputi', 9, true, NOW(), NOW()),

-- Vests
('muski-prsluci',
 '{"sr": "Muški prsluci", "en": "Men Vests", "ru": "Мужские жилеты"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muski-prsluci', 10, true, NOW(), NOW()),

-- Formal
('muska-odela',
 '{"sr": "Muška odela", "en": "Men Suits", "ru": "Мужские костюмы"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muska-odela', 11, true, NOW(), NOW()),

('muski-smokingzi',
 '{"sr": "Muški smokingzi", "en": "Men Tuxedos", "ru": "Мужские смокинги"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muski-smokingzi', 12, true, NOW(), NOW()),

-- Accessories
('muske-kravate',
 '{"sr": "Muške kravate", "en": "Men Ties", "ru": "Мужские галстуки"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muske-kravate', 13, true, NOW(), NOW()),

('muski-kaisevi',
 '{"sr": "Muški kaiševi", "en": "Men Belts", "ru": "Мужские ремни"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muski-kaisevi', 14, true, NOW(), NOW()),

('muski-salovi',
 '{"sr": "Muški šalovi", "en": "Men Scarves", "ru": "Мужские шарфы"}',
 (SELECT id FROM categories WHERE slug = 'muska-odeca' AND level = 2),
 3, 'moda/muska-odeca/muski-salovi', 15, true, NOW(), NOW());

-- ============================================================================
-- ZENSKA ODECA (parent: zenska-odeca) - 11 new L3
-- Already exists: zenske-bluze, zenske-haljine, zenske-jakne, zenske-suknje
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Pants
('zenske-pantalone-jeans',
 '{"sr": "Ženske pantalone farmerke", "en": "Women Jeans", "ru": "Женские джинсы"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenske-pantalone-jeans', 5, true, NOW(), NOW()),

('zenske-pantalone-elegantne',
 '{"sr": "Ženske elegantne pantalone", "en": "Women Dress Pants", "ru": "Женские классические брюки"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenske-pantalone-elegantne', 6, true, NOW(), NOW()),

-- Sweaters & Hoodies
('zenski-dzemperi',
 '{"sr": "Ženski džemperi", "en": "Women Sweaters", "ru": "Женские свитеры"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenski-dzemperi', 7, true, NOW(), NOW()),

('zenski-duksevi',
 '{"sr": "Ženski duksevi", "en": "Women Hoodies", "ru": "Женские худи"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenski-duksevi', 8, true, NOW(), NOW()),

-- Sportswear
('zenska-sportska-odeca',
 '{"sr": "Ženska sportska odeća", "en": "Women Sportswear", "ru": "Женская спортивная одежда"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenska-sportska-odeca', 9, true, NOW(), NOW()),

('zenske-trenerke',
 '{"sr": "Ženske trenerke", "en": "Women Tracksuits", "ru": "Женские спортивные костюмы"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenske-trenerke', 10, true, NOW(), NOW()),

-- Shorts
('zenski-sorcevi',
 '{"sr": "Ženski šorcevi", "en": "Women Shorts", "ru": "Женские шорты"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenski-sorcevi', 11, true, NOW(), NOW()),

-- Coats
('zenski-kaputi',
 '{"sr": "Ženski kaputi", "en": "Women Coats", "ru": "Женские пальто"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenski-kaputi', 12, true, NOW(), NOW()),

('zenske-mantile',
 '{"sr": "Ženske mantile", "en": "Women Trench Coats", "ru": "Женские тренчи"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenske-mantile', 13, true, NOW(), NOW()),

-- Formal
('zenska-vecernja-garderoba',
 '{"sr": "Ženska večernja garderoba", "en": "Women Evening Wear", "ru": "Женская вечерняя одежда"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenska-vecernja-garderoba', 14, true, NOW(), NOW()),

('zenska-poslovna-garderoba',
 '{"sr": "Ženska poslovna garderoba", "en": "Women Business Wear", "ru": "Женская деловая одежда"}',
 (SELECT id FROM categories WHERE slug = 'zenska-odeca' AND level = 2),
 3, 'moda/zenska-odeca/zenska-poslovna-garderoba', 15, true, NOW(), NOW());

-- ============================================================================
-- DECIJA ODECA (parent: decija-odeca) - 15 new L3
-- Already exists: majice-decake, majice-devojcice, pantalone-decake
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Baby clothes by age
('bebe-odeca-0-3-meseca',
 '{"sr": "Bebe odeća 0-3 meseca", "en": "Baby Clothes 0-3 Months", "ru": "Одежда для новорожденных 0-3 месяца"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/bebe-odeca-0-3-meseca', 4, true, NOW(), NOW()),

('bebe-odeca-3-6-meseci',
 '{"sr": "Bebe odeća 3-6 meseci", "en": "Baby Clothes 3-6 Months", "ru": "Одежда для малышей 3-6 месяцев"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/bebe-odeca-3-6-meseci', 5, true, NOW(), NOW()),

('bebe-odeca-6-12-meseci',
 '{"sr": "Bebe odeća 6-12 meseci", "en": "Baby Clothes 6-12 Months", "ru": "Одежда для малышей 6-12 месяцев"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/bebe-odeca-6-12-meseci', 6, true, NOW(), NOW()),

-- Boys by age
('decaci-odeca-1-3-godine',
 '{"sr": "Dečaci odeća 1-3 godine", "en": "Boys Clothes 1-3 Years", "ru": "Одежда для мальчиков 1-3 года"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/decaci-odeca-1-3-godine', 7, true, NOW(), NOW()),

('decaci-odeca-4-7-godina',
 '{"sr": "Dečaci odeća 4-7 godina", "en": "Boys Clothes 4-7 Years", "ru": "Одежда для мальчиков 4-7 лет"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/decaci-odeca-4-7-godina', 8, true, NOW(), NOW()),

('decaci-odeca-8-12-godina',
 '{"sr": "Dečaci odeća 8-12 godina", "en": "Boys Clothes 8-12 Years", "ru": "Одежда для мальчиков 8-12 лет"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/decaci-odeca-8-12-godina', 9, true, NOW(), NOW()),

-- Girls by age
('devojcice-odeca-1-3-godine',
 '{"sr": "Devojčice odeća 1-3 godine", "en": "Girls Clothes 1-3 Years", "ru": "Одежда для девочек 1-3 года"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/devojcice-odeca-1-3-godine', 10, true, NOW(), NOW()),

('devojcice-odeca-4-7-godina',
 '{"sr": "Devojčice odeća 4-7 godina", "en": "Girls Clothes 4-7 Years", "ru": "Одежда для девочек 4-7 лет"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/devojcice-odeca-4-7-godina', 11, true, NOW(), NOW()),

('devojcice-odeca-8-12-godina',
 '{"sr": "Devojčice odeća 8-12 godina", "en": "Girls Clothes 8-12 Years", "ru": "Одежда для девочек 8-12 лет"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/devojcice-odeca-8-12-godina', 12, true, NOW(), NOW()),

-- General kids categories
('decije-jakne',
 '{"sr": "Dečije jakne", "en": "Kids Jackets", "ru": "Детские куртки"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/decije-jakne', 13, true, NOW(), NOW()),

('deciji-duksevi',
 '{"sr": "Dečiji duksevi", "en": "Kids Hoodies", "ru": "Детские худи"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/deciji-duksevi', 14, true, NOW(), NOW()),

('decije-trenerke',
 '{"sr": "Dečije trenerke", "en": "Kids Tracksuits", "ru": "Детские спортивные костюмы"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/decije-trenerke', 15, true, NOW(), NOW()),

('decije-kaputi',
 '{"sr": "Dečiji kaputi", "en": "Kids Coats", "ru": "Детские пальто"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/decije-kaputi', 16, true, NOW(), NOW()),

-- Special occasion
('skolska-uniforma',
 '{"sr": "Školska uniforma", "en": "School Uniforms", "ru": "Школьная форма"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/skolska-uniforma', 17, true, NOW(), NOW()),

('svecana-decija-odeca',
 '{"sr": "Svečana dečija odeća", "en": "Kids Formal Wear", "ru": "Детская праздничная одежда"}',
 (SELECT id FROM categories WHERE slug = 'decija-odeca' AND level = 2),
 3, 'moda/decija-odeca/svecana-decija-odeca', 18, true, NOW(), NOW());

-- ============================================================================
-- MUSKA OBUCA (parent: muska-obuca) - 15 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Sneakers
('muske-patike-casual',
 '{"sr": "Muške patike casual", "en": "Men Casual Sneakers", "ru": "Мужские кроссовки casual"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-patike-casual', 1, true, NOW(), NOW()),

('muske-patike-running',
 '{"sr": "Muške patike za trčanje", "en": "Men Running Shoes", "ru": "Мужские беговые кроссовки"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-patike-running', 2, true, NOW(), NOW()),

('muske-patike-basketball',
 '{"sr": "Muške košarkaške patike", "en": "Men Basketball Shoes", "ru": "Мужские баскетбольные кроссовки"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-patike-basketball', 3, true, NOW(), NOW()),

('muske-patike-football',
 '{"sr": "Muške fudbalske kopačke", "en": "Men Football Boots", "ru": "Мужские бутсы"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-patike-football', 4, true, NOW(), NOW()),

-- Dress shoes
('muske-cipele-koza',
 '{"sr": "Muške kožne cipele", "en": "Men Leather Shoes", "ru": "Мужские кожаные туфли"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-cipele-koza', 5, true, NOW(), NOW()),

('muske-cipele-oxford',
 '{"sr": "Muške Oxford cipele", "en": "Men Oxford Shoes", "ru": "Мужские оксфорды"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-cipele-oxford', 6, true, NOW(), NOW()),

('muske-cipele-derby',
 '{"sr": "Muške Derby cipele", "en": "Men Derby Shoes", "ru": "Мужские дерби"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-cipele-derby', 7, true, NOW(), NOW()),

-- Boots
('muske-cizme-duboke',
 '{"sr": "Muške duboke čizme", "en": "Men High Boots", "ru": "Мужские высокие ботинки"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-cizme-duboke', 8, true, NOW(), NOW()),

('muske-cizme-chelsea',
 '{"sr": "Muške Chelsea čizme", "en": "Men Chelsea Boots", "ru": "Мужские челси"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-cizme-chelsea', 9, true, NOW(), NOW()),

('muske-cizme-radne',
 '{"sr": "Muške radne čizme", "en": "Men Work Boots", "ru": "Мужские рабочие ботинки"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-cizme-radne', 10, true, NOW(), NOW()),

-- Sandals
('muske-sandale',
 '{"sr": "Muške sandale", "en": "Men Sandals", "ru": "Мужские сандалии"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-sandale', 11, true, NOW(), NOW()),

('muske-papuce',
 '{"sr": "Muške papuče", "en": "Men Slippers", "ru": "Мужские тапочки"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-papuce', 12, true, NOW(), NOW()),

-- Loafers
('muske-mokasine',
 '{"sr": "Muške mokasine", "en": "Men Loafers", "ru": "Мужские мокасины"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-mokasine', 13, true, NOW(), NOW()),

-- Espadrilles
('muske-espadrile',
 '{"sr": "Muške espadrile", "en": "Men Espadrilles", "ru": "Мужские эспадрильи"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-espadrile', 14, true, NOW(), NOW()),

-- Boat shoes
('muske-cipele-brodske',
 '{"sr": "Muške brodske cipele", "en": "Men Boat Shoes", "ru": "Мужские топсайдеры"}',
 (SELECT id FROM categories WHERE slug = 'muska-obuca' AND level = 2),
 3, 'moda/muska-obuca/muske-cipele-brodske', 15, true, NOW(), NOW());

-- ============================================================================
-- ZENSKA OBUCA (parent: zenska-obuca) - 15 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Sneakers
('zenske-patike-casual',
 '{"sr": "Ženske patike casual", "en": "Women Casual Sneakers", "ru": "Женские кроссовки casual"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-patike-casual', 1, true, NOW(), NOW()),

('zenske-patike-running',
 '{"sr": "Ženske patike za trčanje", "en": "Women Running Shoes", "ru": "Женские беговые кроссовки"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-patike-running', 2, true, NOW(), NOW()),

('zenske-patike-fitness',
 '{"sr": "Ženske fitnes patike", "en": "Women Fitness Shoes", "ru": "Женские фитнес кроссовки"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-patike-fitness', 3, true, NOW(), NOW()),

-- Heels
('zenske-stikle',
 '{"sr": "Ženske štikle", "en": "Women High Heels", "ru": "Женские туфли на шпильке"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-stikle', 4, true, NOW(), NOW()),

('zenske-cipele-potpetica',
 '{"sr": "Ženske cipele sa potpeticom", "en": "Women Heeled Shoes", "ru": "Женские туфли на каблуке"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-cipele-potpetica', 5, true, NOW(), NOW()),

-- Flats
('zenske-balerinke',
 '{"sr": "Ženske balerinke", "en": "Women Ballet Flats", "ru": "Женские балетки"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-balerinke', 6, true, NOW(), NOW()),

('zenske-mokasine',
 '{"sr": "Ženske mokasine", "en": "Women Loafers", "ru": "Женские мокасины"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-mokasine', 7, true, NOW(), NOW()),

-- Boots
('zenske-cizme-duboke',
 '{"sr": "Ženske duboke čizme", "en": "Women High Boots", "ru": "Женские высокие сапоги"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-cizme-duboke', 8, true, NOW(), NOW()),

('zenske-cizme-gležnjače',
 '{"sr": "Ženske gležnjače", "en": "Women Ankle Boots", "ru": "Женские ботильоны"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-cizme-gleznjace', 9, true, NOW(), NOW()),

('zenske-cizme-preko-kolena',
 '{"sr": "Ženske čizme preko kolena", "en": "Women Over Knee Boots", "ru": "Женские сапоги выше колена"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-cizme-preko-kolena', 10, true, NOW(), NOW()),

-- Sandals
('zenske-sandale',
 '{"sr": "Ženske sandale", "en": "Women Sandals", "ru": "Женские сандалии"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-sandale', 11, true, NOW(), NOW()),

('zenske-papuce',
 '{"sr": "Ženske papuče", "en": "Women Slippers", "ru": "Женские тапочки"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-papuce', 12, true, NOW(), NOW()),

-- Espadrilles
('zenske-espadrile',
 '{"sr": "Ženske espadrile", "en": "Women Espadrilles", "ru": "Женские эспадрильи"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-espadrile', 13, true, NOW(), NOW()),

-- Wedges
('zenske-cipele-platforme',
 '{"sr": "Ženske cipele platforme", "en": "Women Platform Shoes", "ru": "Женские туфли на платформе"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-cipele-platforme', 14, true, NOW(), NOW()),

-- Mules
('zenske-natikace',
 '{"sr": "Ženske natikače", "en": "Women Mules", "ru": "Женские мюли"}',
 (SELECT id FROM categories WHERE slug = 'zenska-obuca' AND level = 2),
 3, 'moda/zenska-obuca/zenske-natikace', 15, true, NOW(), NOW());

-- ============================================================================
-- DECIJA OBUCA (parent: decija-obuca) - 15 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Baby shoes by age
('bebe-obuca-0-6-meseci',
 '{"sr": "Bebe obuća 0-6 meseci", "en": "Baby Shoes 0-6 Months", "ru": "Обувь для новорожденных 0-6 месяцев"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/bebe-obuca-0-6-meseci', 1, true, NOW(), NOW()),

('bebe-obuca-6-12-meseci',
 '{"sr": "Bebe obuća 6-12 meseci", "en": "Baby Shoes 6-12 Months", "ru": "Обувь для малышей 6-12 месяцев"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/bebe-obuca-6-12-meseci', 2, true, NOW(), NOW()),

('decija-obuca-prve-korake',
 '{"sr": "Dečija obuća za prve korake", "en": "First Steps Shoes", "ru": "Обувь для первых шагов"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decija-obuca-prve-korake', 3, true, NOW(), NOW()),

-- Kids sneakers by age
('decije-patike-1-3-godine',
 '{"sr": "Dečije patike 1-3 godine", "en": "Kids Sneakers 1-3 Years", "ru": "Детские кроссовки 1-3 года"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-patike-1-3-godine', 4, true, NOW(), NOW()),

('decije-patike-4-7-godina',
 '{"sr": "Dečije patike 4-7 godina", "en": "Kids Sneakers 4-7 Years", "ru": "Детские кроссовки 4-7 лет"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-patike-4-7-godina', 5, true, NOW(), NOW()),

('decije-patike-8-12-godina',
 '{"sr": "Dečije patike 8-12 godina", "en": "Kids Sneakers 8-12 Years", "ru": "Детские кроссовки 8-12 лет"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-patike-8-12-godina', 6, true, NOW(), NOW()),

-- Sports shoes
('decije-sportske-patike',
 '{"sr": "Dečije sportske patike", "en": "Kids Sports Shoes", "ru": "Детские спортивные кроссовки"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-sportske-patike', 7, true, NOW(), NOW()),

('decije-fudbalske-kopacke',
 '{"sr": "Dečije fudbalske kopačke", "en": "Kids Football Boots", "ru": "Детские бутсы"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-fudbalske-kopacke', 8, true, NOW(), NOW()),

-- Boots
('decije-cizme-gumene',
 '{"sr": "Dečije gumene čizme", "en": "Kids Rain Boots", "ru": "Детские резиновые сапоги"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-cizme-gumene', 9, true, NOW(), NOW()),

('decije-cizme-zimske',
 '{"sr": "Dečije zimske čizme", "en": "Kids Winter Boots", "ru": "Детские зимние сапоги"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-cizme-zimske', 10, true, NOW(), NOW()),

('decije-cizme-duboke',
 '{"sr": "Dečije duboke čizme", "en": "Kids High Boots", "ru": "Детские высокие ботинки"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-cizme-duboke', 11, true, NOW(), NOW()),

-- Sandals
('decije-sandale',
 '{"sr": "Dečije sandale", "en": "Kids Sandals", "ru": "Детские сандалии"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-sandale', 12, true, NOW(), NOW()),

('decije-papuce',
 '{"sr": "Dečije papuče", "en": "Kids Slippers", "ru": "Детские тапочки"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-papuce', 13, true, NOW(), NOW()),

-- School shoes
('decije-cipele-skolske',
 '{"sr": "Dečije školske cipele", "en": "Kids School Shoes", "ru": "Детские школьные туфли"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-cipele-skolske', 14, true, NOW(), NOW()),

-- Ballet
('decije-baletanke',
 '{"sr": "Dečije baletanke", "en": "Kids Ballet Shoes", "ru": "Детские балетки"}',
 (SELECT id FROM categories WHERE slug = 'decija-obuca' AND level = 2),
 3, 'moda/decija-obuca/decije-baletanke', 15, true, NOW(), NOW());

-- ============================================================================
-- DODATNA KATEGORIJA - KUPACI KOSTIMI (parent: kupaci-kostimi) - 7 new L3
-- Already exists: kupaci-kostimi-bikini, kupaci-kostimi-celi
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Men swimwear
('muski-kupaci-sorcevi',
 '{"sr": "Muški kupaći šorcevi", "en": "Men Swim Shorts", "ru": "Мужские плавки-шорты"}',
 (SELECT id FROM categories WHERE slug = 'kupaci-kostimi' AND level = 2),
 3, 'moda/kupaci-kostimi/muski-kupaci-sorcevi', 3, true, NOW(), NOW()),

('muski-kupaci-slip',
 '{"sr": "Muški kupaći slip", "en": "Men Swim Briefs", "ru": "Мужские плавки"}',
 (SELECT id FROM categories WHERE slug = 'kupaci-kostimi' AND level = 2),
 3, 'moda/kupaci-kostimi/muski-kupaci-slip', 4, true, NOW(), NOW()),

-- Women tankini
('zenski-tankini',
 '{"sr": "Ženski tankini", "en": "Women Tankini", "ru": "Женские танкини"}',
 (SELECT id FROM categories WHERE slug = 'kupaci-kostimi' AND level = 2),
 3, 'moda/kupaci-kostimi/zenski-tankini', 5, true, NOW(), NOW()),

('zenski-monokini',
 '{"sr": "Ženski monokini", "en": "Women Monokini", "ru": "Женские монокини"}',
 (SELECT id FROM categories WHERE slug = 'kupaci-kostimi' AND level = 2),
 3, 'moda/kupaci-kostimi/zenski-monokini', 6, true, NOW(), NOW()),

-- Kids swimwear
('deciji-kupaci-kostimi',
 '{"sr": "Dečiji kupaći kostimi", "en": "Kids Swimwear", "ru": "Детские купальники"}',
 (SELECT id FROM categories WHERE slug = 'kupaci-kostimi' AND level = 2),
 3, 'moda/kupaci-kostimi/deciji-kupaci-kostimi', 7, true, NOW(), NOW()),

-- Rash guards
('kupaći-majice-uv-zaštita',
 '{"sr": "Kupaće majice UV zaštita", "en": "UV Protection Swim Shirts", "ru": "Футболки для плавания UV"}',
 (SELECT id FROM categories WHERE slug = 'kupaci-kostimi' AND level = 2),
 3, 'moda/kupaci-kostimi/kupaci-majice-uv-zastita', 8, true, NOW(), NOW()),

-- Cover-ups
('pareo-tunike',
 '{"sr": "Pareo i tunike", "en": "Pareos and Cover-ups", "ru": "Парео и туники"}',
 (SELECT id FROM categories WHERE slug = 'kupaci-kostimi' AND level = 2),
 3, 'moda/kupaci-kostimi/pareo-tunike', 9, true, NOW(), NOW());

-- ============================================================================
-- Migration Summary
-- ============================================================================
-- Total new L3 categories: 90
-- Categories breakdown:
--   - Muska odeca: 12 (total 15 L3)
--   - Zenska odeca: 11 (total 15 L3)
--   - Decija odeca: 15 (total 18 L3)
--   - Muska obuca: 15
--   - Zenska obuca: 15
--   - Decija obuca: 15
--   - Kupaci kostimi: 7 (total 9 L3)
-- Total L3 after this migration: 206
-- ============================================================================
