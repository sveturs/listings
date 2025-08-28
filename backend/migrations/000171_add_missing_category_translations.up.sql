-- Add missing translations for categories without proper translations
-- Categories that need translations: 2006 (photo), 2007 (wifi-routery)

-- Photo category translations (ID: 2006)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
VALUES 
    ('category', 2006, 'en', 'name', 'Photo Equipment', false, true),
    ('category', 2006, 'ru', 'name', 'Фототехника', false, true),
    ('category', 2006, 'sr', 'name', 'Foto oprema', false, true),
    ('category', 2006, 'en', 'seo_title', 'Photo Equipment - Cameras, Lenses, Accessories', false, true),
    ('category', 2006, 'ru', 'seo_title', 'Фототехника - Камеры, объективы, аксессуары', false, true),
    ('category', 2006, 'sr', 'seo_title', 'Foto oprema - Kamere, objektivi, pribor', false, true),
    ('category', 2006, 'en', 'seo_description', 'Digital cameras, lenses, tripods and photo accessories', false, true),
    ('category', 2006, 'ru', 'seo_description', 'Цифровые камеры, объективы, штативы и фотоаксессуары', false, true),
    ('category', 2006, 'sr', 'seo_description', 'Digitalne kamere, objektivi, stativovi i foto pribor', false, true);

-- Wi-Fi routers category translations (ID: 2007)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
VALUES 
    ('category', 2007, 'en', 'name', 'Wi-Fi Routers', false, true),
    ('category', 2007, 'ru', 'name', 'Wi-Fi роутеры', false, true),
    ('category', 2007, 'sr', 'name', 'WiFi ruteri', false, true),
    ('category', 2007, 'en', 'seo_title', 'Wi-Fi Routers - Wireless Network Equipment', false, true),
    ('category', 2007, 'ru', 'seo_title', 'Wi-Fi роутеры - Беспроводное сетевое оборудование', false, true),
    ('category', 2007, 'sr', 'seo_title', 'WiFi ruteri - Bežična mrežna oprema', false, true),
    ('category', 2007, 'en', 'seo_description', 'Wireless routers, modems and network equipment for home and office', false, true),
    ('category', 2007, 'ru', 'seo_description', 'Беспроводные роутеры, модемы и сетевое оборудование для дома и офиса', false, true),
    ('category', 2007, 'sr', 'seo_description', 'Bežični ruteri, modemi i mrežna oprema za dom i kancelariju', false, true);

-- Fix category names to use proper translations instead of technical slugs
UPDATE marketplace_categories SET name = 'Photo Equipment' WHERE id = 2006;
UPDATE marketplace_categories SET name = 'WiFi ruteri' WHERE id = 2007;