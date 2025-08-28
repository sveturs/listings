-- Add missing standard marketplace categories that are commonly found on marketplaces

-- Add Business & Industrial parent category with subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, level, sort_order, is_active, description) VALUES
-- Parent category: Business & Industrial
(2100, 'Poslovanje i industrija', 'business-industrial', NULL, 0, 21, true, 'Oprema i usluge za poslovanje i industriju'),

-- Office Supplies & Equipment
(2101, 'Kancelarijska oprema', 'office-supplies', 2100, 1, 1, true, 'Oprema za kancelariju i ured'),
(2102, 'Štamparija i grafika', 'printing-graphics', 2100, 1, 2, true, 'Usluge štampe i grafičkog dizajna'),
(2103, 'Bezbednost i zaštita', 'security-safety', 2100, 1, 3, true, 'Sistemi bezbednosti i zaštitna oprema'),

-- Collectibles & Hobby
(2200, 'Kolekcionarstvo i hobi', 'collectibles-hobby', NULL, 0, 22, true, 'Kolekcionarski predmeti i hobi artikli'),
(2201, 'Kovanice i novčanice', 'coins-banknotes', 2200, 1, 1, true, 'Numizmatika i kolekcioniranje novca'),
(2202, 'Poštanske marke', 'stamps', 2200, 1, 2, true, 'Filatelija i kolekcioniranje marki'),
(2203, 'Modeli i makete', 'models-miniatures', 2200, 1, 3, true, 'Modeli vozila, brodova, aviona'),

-- Travel & Tourism
(2300, 'Putovanja i turizam', 'travel-tourism', NULL, 0, 23, true, 'Turističke usluge i putovanja'),
(2301, 'Hoteli i smeštaj', 'hotels-accommodation', 2300, 1, 1, true, 'Rezervacija hotela i apartmana'),
(2302, 'Transport i vožnja', 'transport-rides', 2300, 1, 2, true, 'Prevoz i taksi usluge'),
(2303, 'Turistički vodiči', 'tour-guides', 2300, 1, 3, true, 'Vodiči i turističke ture');

-- Add Electronics subcategories that might be missing
INSERT INTO marketplace_categories (id, name, slug, parent_id, level, sort_order, is_active, description) VALUES
-- Electronics subcategories
(2400, 'Dronovi i RC modeli', 'drones-rc', 1001, 1, 10, true, 'Dronovi i radio kontrolisani modeli'),
(2401, 'VR i AR oprema', 'vr-ar-equipment', 1001, 1, 11, true, 'Virtual i Augmented Reality oprema'),
(2402, 'Pametni satovi', 'smartwatches', 1001, 1, 12, true, 'Pametni satovi i fitness trackeri');

-- Add Fashion subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, level, sort_order, is_active, description) VALUES
-- Fashion subcategories
(2500, 'Vintage odeća', 'vintage-clothing', 1002, 1, 10, true, 'Vintage i retro odeća'),
(2501, 'Uniforme', 'uniforms', 1002, 1, 11, true, 'Radne i školske uniforme'),
(2502, 'Kostimi', 'costumes', 1002, 1, 12, true, 'Kostimi za maškare i pozorište');

-- Update sequence for marketplace_categories to handle new IDs
SELECT setval('marketplace_categories_id_seq', (SELECT MAX(id) FROM marketplace_categories) + 1, false);