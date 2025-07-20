-- Serbia Marketplace Categories with full hierarchy and SEO

-- Clear existing data
DELETE FROM marketplace_categories WHERE id >= 1000; -- Use high IDs to avoid conflicts

-- Main Categories (Level 1)
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1001, 'Elektronika', 'electronics', NULL, 1, true, 'Elektronika - Svi elektronski uređaji', 'Kupujte i prodajte elektronske uređaje u Srbiji', 'elektronika, telefoni, kompjuteri, tv'),
(1002, 'Moda', 'fashion', NULL, 2, true, 'Moda i odeća - Garderoba za sve', 'Odeća, obuća i modni dodaci za sve uzraste', 'odeća, obuća, moda, garderoba'),
(1003, 'Automobili', 'automotive', NULL, 3, true, 'Automobili i motori - Vozila u Srbiji', 'Kupovina i prodaja vozila, delova i opreme', 'automobili, motori, delovi, auto oprema'),
(1004, 'Nekretnine', 'real-estate', NULL, 4, true, 'Nekretnine - Stanovi, kuće i zemljišta', 'Prodaja i iznajmljivanje nekretnina u Srbiji', 'nekretnine, stanovi, kuće, zemljište'),
(1005, 'Dom i bašta', 'home-garden', NULL, 5, true, 'Dom i bašta - Sve za vaš dom', 'Namještaj, bašta, alati i kućni dodaci', 'dom, bašta, namještaj, alati'),
(1006, 'Poljoprivreda', 'agriculture', NULL, 6, true, 'Poljoprivreda - Proizvodi i oprema za poljoprivrednike', 'Poljoprivredni proizvodi, mašine i oprema', 'poljoprivreda, traktori, semena, đubriva'),
(1007, 'Industrija', 'industrial', NULL, 7, true, 'Industrija - Oprema i mašine za industriju', 'Industrijska oprema, mašine i alati', 'industrija, mašine, oprema, alati'),
(1008, 'Hrana i piće', 'food-beverages', NULL, 8, true, 'Hrana i piće - Lokalni proizvodi', 'Domaći proizvodi, hrana i pića od lokalnih proizvođača', 'hrana, piće, domaći proizvodi, organsko'),
(1009, 'Usluge', 'services', NULL, 9, true, 'Usluge - Sve vrste usluga u Srbiji', 'Profesionalne usluge za privatne i poslovne klijente', 'usluge, servisi, majstori, konsultanti'),
(1010, 'Sport i rekreacija', 'sports-recreation', NULL, 10, true, 'Sport i rekreacija - Sportska oprema', 'Sportska oprema i proizvodi za rekreaciju', 'sport, rekreacija, fitnes, oprema');

-- Electronics Subcategories (Level 2)
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1101, 'Pametni telefoni', 'smartphones', 1001, 1, true, 'Pametni telefoni - iPhone, Samsung, Huawei', 'Novi i polovan pametni telefoni svih brendova', 'telefoni, iphone, samsung, android'),
(1102, 'Računari', 'computers', 1001, 2, true, 'Računari - Desktop, laptop, komponente', 'Računari, laptopi i komponente za IT', 'računari, laptop, desktop, komponente'),
(1103, 'TV i audio', 'tv-audio', 1001, 3, true, 'TV i audio oprema', 'Televizori, audio sistemi i multimedija', 'tv, audio, televizor, muzika'),
(1104, 'Kućni aparati', 'home-appliances', 1001, 4, true, 'Kućni aparati - Frižideri, mašine, kuvanje', 'Velika i mala kućna tehnika', 'frižider, mašina, sporeti, kućni');

-- Fashion Subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1201, 'Muška odeća', 'mens-clothing', 1002, 1, true, 'Muška odeća - Odela, košulje, pantalone', 'Muška garderoba za sve prilike', 'muška odeća, odela, košulje, pantalone'),
(1202, 'Ženska odeća', 'womens-clothing', 1002, 2, true, 'Ženska odeća - Haljine, majice, suknje', 'Ženska garderoba i moda', 'ženska odeća, haljine, majice, suknje'),
(1203, 'Obuća', 'shoes', 1002, 3, true, 'Obuća - Cipele, patike, čizme', 'Obuća za sve uzraste i prilike', 'cipele, patike, čizme, obuća'),
(1204, 'Modni dodaci', 'accessories', 1002, 4, true, 'Modni dodaci - Torbe, satovi, nakit', 'Modni dodaci i ukrasi', 'torbe, satovi, nakit, modni dodaci');

-- Automotive Subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1301, 'Lični automobili', 'cars', 1003, 1, true, 'Lični automobili - Novi i polovni', 'Kupovina i prodaja ličnih vozila', 'automobili, kola, vozila, polovni'),
(1302, 'Motocikli', 'motorcycles', 1003, 2, true, 'Motocikli i skuteri', 'Dvokolesna vozila i oprema', 'motori, motocikli, skuteri, dvokolesna'),
(1303, 'Auto delovi', 'auto-parts', 1003, 3, true, 'Auto delovi i oprema', 'Rezervni delovi i auto oprema', 'delovi, rezervni, auto oprema, gume'),
(1304, 'Komercijalna vozila', 'trucks-commercial', 1003, 4, true, 'Kamioni i komercijalna vozila', 'Teška vozila za transport i rad', 'kamioni, komercijalna, teška vozila');

-- Real Estate Subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1401, 'Stanovi', 'apartments', 1004, 1, true, 'Stanovi - Prodaja i izdavanje', 'Stanovi u Beogradu, Novom Sadu i širom Srbije', 'stanovi, prodaja, izdavanje, iznajmljivanje'),
(1402, 'Kuće', 'houses', 1004, 2, true, 'Kuće - Porodične kuće i vikendice', 'Kuće za život i odmor', 'kuće, porodične, vikendice, imanja'),
(1403, 'Zemljište', 'land', 1004, 3, true, 'Zemljište - Građevinsko i poljoprivredno', 'Parcele za izgradnju i poljoprivredu', 'zemljište, parcele, građevinsko, poljoprivredno'),
(1404, 'Poslovne nekretnine', 'commercial-real-estate', 1004, 4, true, 'Poslovne nekretnine - Kancelarije, lokali', 'Poslovni prostor za trgovinu i usluge', 'poslovne nekretnine, kancelarije, lokali');

-- Home & Garden Subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1501, 'Nameštaj', 'furniture', 1005, 1, true, 'Nameštaj - Sofe, stolovi, kreveti', 'Nameštaj za kuću, kancelariju i baštu', 'nameštaj, sofe, stolovi, kreveti'),
(1502, 'Alat za baštu', 'garden-tools', 1005, 2, true, 'Alat za baštu - Kosilice, trimeri', 'Alat i oprema za održavanje bašte', 'bašta, alat, kosilice, trimeri'),
(1503, 'Dekoracija', 'home-decor', 1005, 3, true, 'Dekoracija doma - Slike, vaze, tepisoni', 'Ukrasni predmeti za dom', 'dekoracija, slike, vaze, tepisoni'),
(1504, 'Građevinski materijal', 'building-materials', 1005, 4, true, 'Građevinski materijal - Pločice, farbe', 'Materijali za gradnju i renoviranje', 'građevinski materijal, pločice, farbe');

-- Agriculture Subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1601, 'Poljoprivredne mašine', 'farm-machinery', 1006, 1, true, 'Poljoprivredne mašine - Traktori, kombaini', 'Mašine i oprema za poljoprivredu', 'traktori, kombaini, mašine, poljoprivreda'),
(1602, 'Seme i đubriva', 'seeds-fertilizers', 1006, 2, true, 'Seme i đubriva - Kvalitetno seme i hranjiva', 'Semena, đubriva i preparati za poljoprivredu', 'seme, đubriva, preparati, hranjiva'),
(1603, 'Stoka', 'livestock', 1006, 3, true, 'Stoka - Goveda, svinje, živo perje', 'Prodaja domaćih životinja', 'stoka, goveda, svinje, kokoške'),
(1604, 'Poljoprivredni proizvodi', 'farm-products', 1006, 4, true, 'Poljoprivredni proizvodi - Žitarice, voće', 'Direktna prodaja od proizvođača', 'žitarice, voće, povrće, poljoprivreda');

-- Industrial Subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1701, 'Industrijske mašine', 'industrial-machinery', 1007, 1, true, 'Industrijske mašine - CNC, prese, alati', 'Oprema za industrijsku proizvodnju', 'cnc, prese, alati, industrijske mašine'),
(1702, 'Hemijski proizvodi', 'chemical-products', 1007, 2, true, 'Hemijski proizvodi - Sirovine, reagensi', 'Hemikalije za industrijsku upotrebu', 'hemikalije, reagensi, sirovine, hemija'),
(1703, 'Sirovine', 'raw-materials', 1007, 3, true, 'Sirovine - Metal, plastika, tekstil', 'Sirovine za industrijsku proizvodnju', 'metal, plastika, tekstil, sirovine'),
(1704, 'Bezbednosna oprema', 'safety-equipment', 1007, 4, true, 'Bezbednosna oprema - Zaštitna odela, kacige', 'Oprema za bezbednost na radu', 'zaštitna oprema, kacige, rukavice');

-- Food & Beverages Subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1801, 'Organska hrana', 'organic-food', 1008, 1, true, 'Organska hrana - Prirodni proizvodi', 'Zdrava hrana bez hemikalija', 'organska hrana, prirodno, zdravo, eko'),
(1802, 'Pića', 'beverages', 1008, 2, true, 'Pića - Sokovi, vino, rakija', 'Domaća pića i napici', 'pića, sokovi, vino, rakija'),
(1803, 'Mlečni proizvodi', 'dairy-products', 1008, 3, true, 'Mlečni proizvodi - Sir, kajmak, jogurt', 'Sveži mlečni proizvodi', 'mleko, sir, kajmak, jogurt'),
(1804, 'Mesni proizvodi', 'meat-products', 1008, 4, true, 'Mesni proizvodi - Sveže meso, kobasice', 'Kvalitetni mesni proizvodi', 'meso, kobasice, slanina, domaće');

-- Services Subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(1901, 'Građevinske usluge', 'construction-services', 1009, 1, true, 'Građevinske usluge - Majstori, renoviranje', 'Građevinski radovi i renoviranje', 'majstori, renoviranje, gradnja, radovi'),
(1902, 'IT usluge', 'it-services', 1009, 2, true, 'IT usluge - Web, software, dizajn', 'Informacione tehnologije i digitalne usluge', 'web, software, dizajn, programiranje'),
(1903, 'Lepota i wellness', 'beauty-wellness', 1009, 3, true, 'Lepota i wellness - Frizeri, masaža', 'Usluge lepote i brige o telu', 'frizerski, masaža, kozmetika, wellness'),
(1904, 'Poslovne usluge', 'business-services', 1009, 4, true, 'Poslovne usluge - Računovodstvo, marketing', 'Profesionalne usluge za biznise', 'računovodstvo, marketing, konsalting');

-- Sports & Recreation Subcategories
INSERT INTO marketplace_categories (id, name, slug, parent_id, sort_order, is_active, seo_title, seo_description, seo_keywords) VALUES
(2001, 'Fitnes oprema', 'fitness-equipment', 1010, 1, true, 'Fitnes oprema - Sprave, tegovi', 'Oprema za vežbanje i fitnes', 'fitnes, sprave, tegovi, vežbanje'),
(2002, 'Sportovi na otvorenom', 'outdoor-sports', 1010, 2, true, 'Sporturi na otvorenom - Bicikli, planinarenje', 'Oprema za aktivnosti napolju', 'bicikli, planinarenje, kampovanje'),
(2003, 'Timski sportovi', 'team-sports', 1010, 3, true, 'Timski sportovi - Fudbal, košarka', 'Oprema za kolektivne sportove', 'fudbal, košarka, odbojka, sport'),
(2004, 'Zimski sportovi', 'winter-sports', 1010, 4, true, 'Zimski sportovi - Skije, skejtovanje', 'Oprema za zimske aktivnosti', 'skije, snowboard, klizanje, zima');

-- Translations for main categories
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text, is_machine_translated, is_verified) VALUES
-- Electronics
('category', 1001, 'name', 'sr', 'Електроника', false, true),
('category', 1001, 'name', 'ru', 'Электроника', true, false),
('category', 1001, 'name', 'en', 'Electronics', true, false),
-- Fashion  
('category', 1002, 'name', 'sr', 'Мода', false, true),
('category', 1002, 'name', 'ru', 'Мода', true, false),
('category', 1002, 'name', 'en', 'Fashion', true, false),
-- Automotive
('category', 1003, 'name', 'sr', 'Аутомобили', false, true),
('category', 1003, 'name', 'ru', 'Автомобили', true, false),
('category', 1003, 'name', 'en', 'Automotive', true, false),
-- Real Estate
('category', 1004, 'name', 'sr', 'Некретнине', false, true),
('category', 1004, 'name', 'ru', 'Недвижимость', true, false),
('category', 1004, 'name', 'en', 'Real Estate', true, false),
-- Home & Garden
('category', 1005, 'name', 'sr', 'Дом и башта', false, true),
('category', 1005, 'name', 'ru', 'Дом и сад', true, false),
('category', 1005, 'name', 'en', 'Home & Garden', true, false),
-- Agriculture
('category', 1006, 'name', 'sr', 'Пољопривреда', false, true),
('category', 1006, 'name', 'ru', 'Сельское хозяйство', true, false),
('category', 1006, 'name', 'en', 'Agriculture', true, false),
-- Industrial
('category', 1007, 'name', 'sr', 'Индустрија', false, true),
('category', 1007, 'name', 'ru', 'Промышленность', true, false),
('category', 1007, 'name', 'en', 'Industrial', true, false),
-- Food & Beverages
('category', 1008, 'name', 'sr', 'Храна и пиће', false, true),
('category', 1008, 'name', 'ru', 'Еда и напитки', true, false),
('category', 1008, 'name', 'en', 'Food & Beverages', true, false),
-- Services
('category', 1009, 'name', 'sr', 'Услуге', false, true),
('category', 1009, 'name', 'ru', 'Услуги', true, false),
('category', 1009, 'name', 'en', 'Services', true, false),
-- Sports & Recreation
('category', 1010, 'name', 'sr', 'Спорт и рекреација', false, true),
('category', 1010, 'name', 'ru', 'Спорт и отдых', true, false),
('category', 1010, 'name', 'en', 'Sports & Recreation', true, false);

-- Reset sequence to avoid conflicts
SELECT setval('marketplace_categories_id_seq', 2100, true);