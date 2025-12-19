-- Migration: Seed L1 (top-level) categories
-- Date: 2025-12-16
-- Purpose: Insert 18 main categories with multilingual support (sr/en/ru)
-- Reference: docs/marketplace-categories-implementation/02-CATEGORY-CATALOG.md

-- Insert 18 L1 categories
INSERT INTO categories (slug, level, path, sort_order, name, description, meta_title, meta_description, icon, google_category_id, is_active) VALUES

-- 1. Одежда и обувь / Clothing & Footwear
('odeca-i-obuca', 1, 'odeca-i-obuca', 1,
 '{"sr": "Odeća i obuća", "en": "Clothing & Footwear", "ru": "Одежда и обувь"}'::jsonb,
 '{"sr": "Muška, ženska i dečja odeća, obuća i modni dodaci", "en": "Men''s, women''s and children''s clothing, footwear and fashion accessories", "ru": "Мужская, женская и детская одежда, обувь и модные аксессуары"}'::jsonb,
 '{"sr": "Odeća i obuća | Vondi Marketplace", "en": "Clothing & Footwear | Vondi Marketplace", "ru": "Одежда и обувь | Vondi Marketplace"}'::jsonb,
 '{"sr": "Kupite odeću i obuću online. Širok izbor muške, ženske i dečje odeće.", "en": "Buy clothing and footwear online. Wide selection of men''s, women''s and children''s clothing.", "ru": "Купить одежду и обувь онлайн. Большой выбор мужской, женской и детской одежды."}'::jsonb,
 '👕', 166, true),

-- 2. Электроника / Electronics
('elektronika', 1, 'elektronika', 2,
 '{"sr": "Elektronika", "en": "Electronics", "ru": "Электроника"}'::jsonb,
 '{"sr": "Pametni telefoni, računari, TV, audio oprema i dodaci", "en": "Smartphones, computers, TVs, audio equipment and accessories", "ru": "Смартфоны, компьютеры, ТВ, аудио оборудование и аксессуары"}'::jsonb,
 '{"sr": "Elektronika | Vondi Marketplace", "en": "Electronics | Vondi Marketplace", "ru": "Электроника | Vondi Marketplace"}'::jsonb,
 '{"sr": "Kupite elektroniku online. Telefoni, laptopovi, TV i audio oprema.", "en": "Buy electronics online. Phones, laptops, TVs and audio equipment.", "ru": "Купить электронику онлайн. Телефоны, ноутбуки, ТВ и аудио техника."}'::jsonb,
 '📱', 222, true),

-- 3. Дом и сад / Home & Garden
('dom-i-basta', 1, 'dom-i-basta', 3,
 '{"sr": "Dom i bašta", "en": "Home & Garden", "ru": "Дом и сад"}'::jsonb,
 '{"sr": "Nameštaj, dekoracije, kuhinja, bašta i sve za uređenje doma", "en": "Furniture, decorations, kitchen, garden and everything for home improvement", "ru": "Мебель, декор, кухня, сад и всё для обустройства дома"}'::jsonb,
 '{"sr": "Dom i bašta | Vondi Marketplace", "en": "Home & Garden | Vondi Marketplace", "ru": "Дом и сад | Vondi Marketplace"}'::jsonb,
 '{"sr": "Nameštaj, dekoracije i sve za vaš dom i baštu.", "en": "Furniture, decorations and everything for your home and garden.", "ru": "Мебель, декор и всё для вашего дома и сада."}'::jsonb,
 '🏠', 536, true),

-- 4. Красота и здоровье / Beauty & Health
('lepota-i-zdravlje', 1, 'lepota-i-zdravlje', 4,
 '{"sr": "Lepota i zdravlje", "en": "Beauty & Health", "ru": "Красота и здоровье"}'::jsonb,
 '{"sr": "Kozmetika, nega kože, parfemi, nega kose i zdravlje", "en": "Cosmetics, skincare, perfumes, hair care and health", "ru": "Косметика, уход за кожей, парфюмерия, уход за волосами и здоровье"}'::jsonb,
 '{"sr": "Lepota i zdravlje | Vondi Marketplace", "en": "Beauty & Health | Vondi Marketplace", "ru": "Красота и здоровье | Vondi Marketplace"}'::jsonb,
 '{"sr": "Kozmetika, parfemi, nega kože i kose, zdravlje.", "en": "Cosmetics, perfumes, skin and hair care, health.", "ru": "Косметика, парфюмерия, уход за кожей и волосами, здоровье."}'::jsonb,
 '💄', 469, true),

-- 5. Детские товары / Baby & Kids
('za-bebe-i-decu', 1, 'za-bebe-i-decu', 5,
 '{"sr": "Za bebe i decu", "en": "Baby & Kids", "ru": "Детские товары"}'::jsonb,
 '{"sr": "Oprema za bebe, igračke, odeća, obuća i sve za decu", "en": "Baby gear, toys, clothing, footwear and everything for kids", "ru": "Детское снаряжение, игрушки, одежда, обувь и всё для детей"}'::jsonb,
 '{"sr": "Za bebe i decu | Vondi Marketplace", "en": "Baby & Kids | Vondi Marketplace", "ru": "Детские товары | Vondi Marketplace"}'::jsonb,
 '{"sr": "Oprema za bebe, igračke, odeća i sve za vašu decu.", "en": "Baby gear, toys, clothing and everything for your kids.", "ru": "Детское снаряжение, игрушки, одежда и всё для ваших детей."}'::jsonb,
 '👶', 561, true),

-- 6. Спорт и туризм / Sports & Outdoors
('sport-i-turizam', 1, 'sport-i-turizam', 6,
 '{"sr": "Sport i turizam", "en": "Sports & Outdoors", "ru": "Спорт и туризм"}'::jsonb,
 '{"sr": "Fitnes, bicikli, kampovanje, sportska oprema i outdoor aktivnosti", "en": "Fitness, bicycles, camping, sports equipment and outdoor activities", "ru": "Фитнес, велосипеды, кемпинг, спортивное снаряжение и outdoor активности"}'::jsonb,
 '{"sr": "Sport i turizam | Vondi Marketplace", "en": "Sports & Outdoors | Vondi Marketplace", "ru": "Спорт и туризм | Vondi Marketplace"}'::jsonb,
 '{"sr": "Sportska oprema, fitnes, bicikli, kampovanje i outdoor.", "en": "Sports equipment, fitness, bicycles, camping and outdoor.", "ru": "Спортивное снаряжение, фитнес, велосипеды, кемпинг и outdoor."}'::jsonb,
 '⚽', 499, true),

-- 7. Автотовары / Automotive
('automobilizam', 1, 'automobilizam', 7,
 '{"sr": "Automobilizam", "en": "Automotive", "ru": "Автотовары"}'::jsonb,
 '{"sr": "Delovi za automobile, moto oprema, gume, dodaci i alati", "en": "Auto parts, motorcycle gear, tires, accessories and tools", "ru": "Автозапчасти, мото снаряжение, шины, аксессуары и инструменты"}'::jsonb,
 '{"sr": "Automobilizam | Vondi Marketplace", "en": "Automotive | Vondi Marketplace", "ru": "Автотовары | Vondi Marketplace"}'::jsonb,
 '{"sr": "Delovi za automobile, gume, dodaci i alati za vozila.", "en": "Auto parts, tires, accessories and tools for vehicles.", "ru": "Автозапчасти, шины, аксессуары и инструменты для автомобилей."}'::jsonb,
 '🚗', 899, true),

-- 8. Бытовая техника / Appliances
('kucni-aparati', 1, 'kucni-aparati', 8,
 '{"sr": "Kućni aparati", "en": "Appliances", "ru": "Бытовая техника"}'::jsonb,
 '{"sr": "Frižideri, mašine za pranje, usisivači i sva bela tehnika", "en": "Refrigerators, washing machines, vacuum cleaners and all white goods", "ru": "Холодильники, стиральные машины, пылесосы и вся белая техника"}'::jsonb,
 '{"sr": "Kućni aparati | Vondi Marketplace", "en": "Appliances | Vondi Marketplace", "ru": "Бытовая техника | Vondi Marketplace"}'::jsonb,
 '{"sr": "Bela tehnika, frižideri, mašine za pranje, usisivači.", "en": "White goods, refrigerators, washing machines, vacuum cleaners.", "ru": "Белая техника, холодильники, стиральные машины, пылесосы."}'::jsonb,
 '🔌', 604, true),

-- 9. Украшения и часы / Jewelry & Watches
('nakit-i-satovi', 1, 'nakit-i-satovi', 9,
 '{"sr": "Nakit i satovi", "en": "Jewelry & Watches", "ru": "Украшения и часы"}'::jsonb,
 '{"sr": "Nakit, satovi, biserno nakit, vereničko prstenje i luksuzni dodaci", "en": "Jewelry, watches, pearl jewelry, engagement rings and luxury accessories", "ru": "Украшения, часы, жемчужные украшения, обручальные кольца и люксовые аксессуары"}'::jsonb,
 '{"sr": "Nakit i satovi | Vondi Marketplace", "en": "Jewelry & Watches | Vondi Marketplace", "ru": "Украшения и часы | Vondi Marketplace"}'::jsonb,
 '{"sr": "Nakit, satovi, vereničko prstenje i luksuzni dodaci.", "en": "Jewelry, watches, engagement rings and luxury accessories.", "ru": "Украшения, часы, обручальные кольца и люксовые аксессуары."}'::jsonb,
 '💍', 188, true),

-- 10. Книги и медиа / Books & Media
('knjige-i-mediji', 1, 'knjige-i-mediji', 10,
 '{"sr": "Knjige i mediji", "en": "Books & Media", "ru": "Книги и медиа"}'::jsonb,
 '{"sr": "Knjige, časopisi, filmovi, muzika i audio knjige", "en": "Books, magazines, movies, music and audiobooks", "ru": "Книги, журналы, фильмы, музыка и аудиокниги"}'::jsonb,
 '{"sr": "Knjige i mediji | Vondi Marketplace", "en": "Books & Media | Vondi Marketplace", "ru": "Книги и медиа | Vondi Marketplace"}'::jsonb,
 '{"sr": "Knjige, časopisi, filmovi, muzika i audio knjige.", "en": "Books, magazines, movies, music and audiobooks.", "ru": "Книги, журналы, фильмы, музыка и аудиокниги."}'::jsonb,
 '📚', 783, true),

-- 11. Зоотовары / Pet Supplies
('kucni-ljubimci', 1, 'kucni-ljubimci', 11,
 '{"sr": "Kućni ljubimci", "en": "Pet Supplies", "ru": "Зоотовары"}'::jsonb,
 '{"sr": "Hrana za ljubimce, igračke, oprema za pse, mačke i druge ljubimce", "en": "Pet food, toys, supplies for dogs, cats and other pets", "ru": "Корм для питомцев, игрушки, товары для собак, кошек и других питомцев"}'::jsonb,
 '{"sr": "Kućni ljubimci | Vondi Marketplace", "en": "Pet Supplies | Vondi Marketplace", "ru": "Зоотовары | Vondi Marketplace"}'::jsonb,
 '{"sr": "Hrana, igračke i oprema za vaše ljubimce.", "en": "Food, toys and supplies for your pets.", "ru": "Корм, игрушки и товары для ваших питомцев."}'::jsonb,
 '🐕', 6, true),

-- 12. Канцтовары / Office Supplies
('kancelarijski-materijal', 1, 'kancelarijski-materijal', 12,
 '{"sr": "Kancelarijski materijal", "en": "Office Supplies", "ru": "Канцтовары"}'::jsonb,
 '{"sr": "Kancelarijski materijal, papir, olovke, sveske i oprema za kancelariju", "en": "Office supplies, paper, pens, notebooks and office equipment", "ru": "Канцелярские товары, бумага, ручки, тетради и офисное оборудование"}'::jsonb,
 '{"sr": "Kancelarijski materijal | Vondi Marketplace", "en": "Office Supplies | Vondi Marketplace", "ru": "Канцтовары | Vondi Marketplace"}'::jsonb,
 '{"sr": "Kancelarijski materijal, papir, sveske i oprema.", "en": "Office supplies, paper, notebooks and equipment.", "ru": "Канцелярские товары, бумага, тетради и оборудование."}'::jsonb,
 '📎', 922, true),

-- 13. Продукты питания / Food & Beverages
('hrana-i-pice', 1, 'hrana-i-pice', 13,
 '{"sr": "Hrana i piće", "en": "Food & Beverages", "ru": "Продукты питания"}'::jsonb,
 '{"sr": "Organska hrana, poslastice, piće, kafa, čaj i delikatesi", "en": "Organic food, sweets, beverages, coffee, tea and delicacies", "ru": "Органическая еда, сладости, напитки, кофе, чай и деликатесы"}'::jsonb,
 '{"sr": "Hrana i piće | Vondi Marketplace", "en": "Food & Beverages | Vondi Marketplace", "ru": "Продукты питания | Vondi Marketplace"}'::jsonb,
 '{"sr": "Organska hrana, poslastice, kafa, čaj i delikatesi.", "en": "Organic food, sweets, coffee, tea and delicacies.", "ru": "Органическая еда, сладости, кофе, чай и деликатесы."}'::jsonb,
 '🍕', 427, true),

-- 14. Хендмейд / Art & Crafts
('umetnost-i-rukotvorine', 1, 'umetnost-i-rukotvorine', 14,
 '{"sr": "Umetnost i rukotvorine", "en": "Art & Crafts", "ru": "Искусство и рукоделие"}'::jsonb,
 '{"sr": "Ručni rad, materijali za izradu, umetničke slike i kreativne potrepštine", "en": "Handmade items, craft materials, artwork and creative supplies", "ru": "Ручная работа, материалы для творчества, художественные картины и творческие принадлежности"}'::jsonb,
 '{"sr": "Umetnost i rukotvorine | Vondi Marketplace", "en": "Art & Crafts | Vondi Marketplace", "ru": "Искусство и рукоделие | Vondi Marketplace"}'::jsonb,
 '{"sr": "Ručni rad, materijali i umetničke slike.", "en": "Handmade items, materials and artwork.", "ru": "Ручная работа, материалы и художественные картины."}'::jsonb,
 '🎨', 505, true),

-- 15. Музыкальные инструменты / Musical Instruments
('muzicki-instrumenti', 1, 'muzicki-instrumenti', 15,
 '{"sr": "Muzički instrumenti", "en": "Musical Instruments", "ru": "Музыкальные инструменты"}'::jsonb,
 '{"sr": "Gitare, klavijature, bubnjevi, duvački instrumenti i oprema", "en": "Guitars, keyboards, drums, wind instruments and equipment", "ru": "Гитары, клавишные, барабаны, духовые инструменты и оборудование"}'::jsonb,
 '{"sr": "Muzički instrumenti | Vondi Marketplace", "en": "Musical Instruments | Vondi Marketplace", "ru": "Музыкальные инструменты | Vondi Marketplace"}'::jsonb,
 '{"sr": "Gitare, klavijature, bubnjevi i muzička oprema.", "en": "Guitars, keyboards, drums and music equipment.", "ru": "Гитары, клавишные, барабаны и музыкальное оборудование."}'::jsonb,
 '🎸', 55, true),

-- 16. Промтовары и инструменты / Industrial & Tools
('industrija-i-alati', 1, 'industrija-i-alati', 16,
 '{"sr": "Industrija i alati", "en": "Industrial & Tools", "ru": "Промтовары и инструменты"}'::jsonb,
 '{"sr": "Ručni alati, električni alati, industrijska oprema i građevinski materijali", "en": "Hand tools, power tools, industrial equipment and building materials", "ru": "Ручные инструменты, электроинструменты, промышленное оборудование и строительные материалы"}'::jsonb,
 '{"sr": "Industrija i alati | Vondi Marketplace", "en": "Industrial & Tools | Vondi Marketplace", "ru": "Промтовары и инструменты | Vondi Marketplace"}'::jsonb,
 '{"sr": "Alati, industrijska oprema i građevinski materijali.", "en": "Tools, industrial equipment and building materials.", "ru": "Инструменты, промышленное оборудование и строительные материалы."}'::jsonb,
 '🔧', 1167, true),

-- 17. Услуги / Services
('usluge', 1, 'usluge', 17,
 '{"sr": "Usluge", "en": "Services", "ru": "Услуги"}'::jsonb,
 '{"sr": "Profesionalne usluge, popravke, čišćenje, prevoz i ostale usluge", "en": "Professional services, repairs, cleaning, transportation and other services", "ru": "Профессиональные услуги, ремонт, уборка, транспорт и прочие услуги"}'::jsonb,
 '{"sr": "Usluge | Vondi Marketplace", "en": "Services | Vondi Marketplace", "ru": "Услуги | Vondi Marketplace"}'::jsonb,
 '{"sr": "Profesionalne usluge, popravke, čišćenje i prevoz.", "en": "Professional services, repairs, cleaning and transportation.", "ru": "Профессиональные услуги, ремонт, уборка и транспорт."}'::jsonb,
 '⚙️', 2764, true),

-- 18. Прочее / Other
('ostalo', 1, 'ostalo', 18,
 '{"sr": "Ostalo", "en": "Other", "ru": "Прочее"}'::jsonb,
 '{"sr": "Sve ostale kategorije proizvoda koje ne spadaju u gore navedene", "en": "All other product categories not listed above", "ru": "Все прочие категории товаров, не перечисленные выше"}'::jsonb,
 '{"sr": "Ostalo | Vondi Marketplace", "en": "Other | Vondi Marketplace", "ru": "Прочее | Vondi Marketplace"}'::jsonb,
 '{"sr": "Sve ostale kategorije proizvoda.", "en": "All other product categories.", "ru": "Все прочие категории товаров."}'::jsonb,
 '📦', 1, true);

-- Verify insertion
DO $$
DECLARE
    category_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO category_count FROM categories WHERE level = 1;

    IF category_count != 18 THEN
        RAISE EXCEPTION 'Expected 18 L1 categories, but found %', category_count;
    END IF;

    RAISE NOTICE 'Successfully inserted 18 L1 categories';
END $$;
