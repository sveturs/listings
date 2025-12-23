-- Заполнение brand_category_mapping популярными брендами
-- Формат: brand_name, brand_aliases[], category_slug, confidence, is_verified

-- === ELEKTRONIKA ===

-- Smartphones & Tech
INSERT INTO brand_category_mapping (brand_name, brand_aliases, category_slug, confidence, is_verified)
VALUES
    ('Apple', ARRAY['apple', 'iphone', 'ipad', 'macbook', 'mac', 'airpods'], 'elektronika', 0.98, true),
    ('Samsung', ARRAY['samsung', 'galaxy'], 'elektronika', 0.98, true),
    ('Xiaomi', ARRAY['xiaomi', 'redmi', 'poco', 'mi'], 'elektronika', 0.98, true),
    ('Huawei', ARRAY['huawei', 'honor'], 'elektronika', 0.98, true),
    ('OnePlus', ARRAY['oneplus', 'one plus'], 'elektronika', 0.98, true),
    ('Google', ARRAY['google', 'pixel'], 'elektronika', 0.95, true),
    ('Sony', ARRAY['sony', 'xperia', 'playstation', 'ps5', 'ps4'], 'elektronika', 0.95, true),
    ('LG', ARRAY['lg'], 'elektronika', 0.95, true),
    ('Motorola', ARRAY['motorola', 'moto'], 'elektronika', 0.95, true),
    ('Nokia', ARRAY['nokia'], 'elektronika', 0.95, true),
    ('Realme', ARRAY['realme'], 'elektronika', 0.95, true),
    ('Oppo', ARRAY['oppo'], 'elektronika', 0.95, true),
    ('Vivo', ARRAY['vivo'], 'elektronika', 0.95, true),
    ('Lenovo', ARRAY['lenovo', 'thinkpad'], 'elektronika', 0.95, true),
    ('Asus', ARRAY['asus', 'rog'], 'elektronika', 0.95, true),
    ('Acer', ARRAY['acer', 'predator'], 'elektronika', 0.95, true),
    ('HP', ARRAY['hp', 'hewlett packard'], 'elektronika', 0.95, true),
    ('Dell', ARRAY['dell', 'alienware'], 'elektronika', 0.95, true),
    ('Microsoft', ARRAY['microsoft', 'surface', 'xbox'], 'elektronika', 0.95, true),
    ('JBL', ARRAY['jbl'], 'elektronika', 0.90, true),
    ('Bose', ARRAY['bose'], 'elektronika', 0.90, true),
    ('Beats', ARRAY['beats', 'beats by dre'], 'elektronika', 0.90, true);

-- === ODECA I OBUCA ===

-- Fashion & Clothing
INSERT INTO brand_category_mapping (brand_name, brand_aliases, category_slug, confidence, is_verified)
VALUES
    ('Nike', ARRAY['nike', 'air jordan', 'jordan'], 'odeca-i-obuca', 0.98, true),
    ('Adidas', ARRAY['adidas', 'yeezy'], 'odeca-i-obuca', 0.98, true),
    ('Puma', ARRAY['puma'], 'odeca-i-obuca', 0.98, true),
    ('Reebok', ARRAY['reebok'], 'odeca-i-obuca', 0.95, true),
    ('New Balance', ARRAY['new balance', 'nb'], 'odeca-i-obuca', 0.95, true),
    ('Converse', ARRAY['converse', 'chuck taylor'], 'odeca-i-obuca', 0.95, true),
    ('Vans', ARRAY['vans'], 'odeca-i-obuca', 0.95, true),
    ('Zara', ARRAY['zara'], 'odeca-i-obuca', 0.98, true),
    ('H&M', ARRAY['h&m', 'hm', 'h and m'], 'odeca-i-obuca', 0.98, true),
    ('Pull&Bear', ARRAY['pull&bear', 'pull and bear', 'pullbear'], 'odeca-i-obuca', 0.95, true),
    ('Bershka', ARRAY['bershka'], 'odeca-i-obuca', 0.95, true),
    ('Mango', ARRAY['mango'], 'odeca-i-obuca', 0.95, true),
    ('Reserved', ARRAY['reserved'], 'odeca-i-obuca', 0.95, true),
    ('Tommy Hilfiger', ARRAY['tommy hilfiger', 'tommy', 'th'], 'odeca-i-obuca', 0.95, true),
    ('Calvin Klein', ARRAY['calvin klein', 'ck'], 'odeca-i-obuca', 0.95, true),
    ('Levi''s', ARRAY['levis', 'levi''s', 'levi'], 'odeca-i-obuca', 0.95, true),
    ('Guess', ARRAY['guess'], 'odeca-i-obuca', 0.95, true),
    ('Lacoste', ARRAY['lacoste'], 'odeca-i-obuca', 0.95, true);

-- === AUTOMOBILIZAM ===

-- Cars & Automotive
INSERT INTO brand_category_mapping (brand_name, brand_aliases, category_slug, confidence, is_verified)
VALUES
    ('BMW', ARRAY['bmw', 'bayerische motoren werke'], 'automobilizam', 0.98, true),
    ('Mercedes', ARRAY['mercedes', 'mercedes-benz', 'mercedes benz', 'benz'], 'automobilizam', 0.98, true),
    ('Audi', ARRAY['audi'], 'automobilizam', 0.98, true),
    ('Volkswagen', ARRAY['volkswagen', 'vw'], 'automobilizam', 0.98, true),
    ('Toyota', ARRAY['toyota'], 'automobilizam', 0.98, true),
    ('Honda', ARRAY['honda'], 'automobilizam', 0.98, true),
    ('Ford', ARRAY['ford'], 'automobilizam', 0.98, true),
    ('Opel', ARRAY['opel'], 'automobilizam', 0.98, true),
    ('Renault', ARRAY['renault'], 'automobilizam', 0.98, true),
    ('Peugeot', ARRAY['peugeot'], 'automobilizam', 0.98, true),
    ('Citroen', ARRAY['citroen', 'citroën'], 'automobilizam', 0.98, true),
    ('Fiat', ARRAY['fiat'], 'automobilizam', 0.98, true),
    ('Skoda', ARRAY['skoda', 'škoda'], 'automobilizam', 0.98, true),
    ('Seat', ARRAY['seat'], 'automobilizam', 0.95, true),
    ('Mazda', ARRAY['mazda'], 'automobilizam', 0.95, true),
    ('Nissan', ARRAY['nissan'], 'automobilizam', 0.95, true),
    ('Hyundai', ARRAY['hyundai'], 'automobilizam', 0.95, true),
    ('Kia', ARRAY['kia'], 'automobilizam', 0.95, true),
    ('Suzuki', ARRAY['suzuki'], 'automobilizam', 0.95, true),
    ('Bosch', ARRAY['bosch'], 'automobilizam', 0.90, true),
    ('Michelin', ARRAY['michelin'], 'automobilizam', 0.90, true),
    ('Bridgestone', ARRAY['bridgestone'], 'automobilizam', 0.90, true),
    ('Continental', ARRAY['continental'], 'automobilizam', 0.90, true);

-- === KUCNI APARATI ===

-- Home Appliances
INSERT INTO brand_category_mapping (brand_name, brand_aliases, category_slug, confidence, is_verified)
VALUES
    ('Bosch', ARRAY['bosch'], 'kucni-aparati', 0.95, true),
    ('Siemens', ARRAY['siemens'], 'kucni-aparati', 0.95, true),
    ('Whirlpool', ARRAY['whirlpool'], 'kucni-aparati', 0.95, true),
    ('Beko', ARRAY['beko'], 'kucni-aparati', 0.95, true),
    ('Gorenje', ARRAY['gorenje'], 'kucni-aparati', 0.95, true),
    ('Samsung', ARRAY['samsung'], 'kucni-aparati', 0.95, true),
    ('LG', ARRAY['lg'], 'kucni-aparati', 0.95, true),
    ('Miele', ARRAY['miele'], 'kucni-aparati', 0.95, true),
    ('Electrolux', ARRAY['electrolux'], 'kucni-aparati', 0.95, true),
    ('Philips', ARRAY['philips'], 'kucni-aparati', 0.95, true),
    ('Tefal', ARRAY['tefal'], 'kucni-aparati', 0.90, true),
    ('Braun', ARRAY['braun'], 'kucni-aparati', 0.90, true),
    ('Dyson', ARRAY['dyson'], 'kucni-aparati', 0.90, true);

-- === SPORT I TURIZAM ===

-- Sports Equipment
INSERT INTO brand_category_mapping (brand_name, brand_aliases, category_slug, confidence, is_verified)
VALUES
    ('Nike', ARRAY['nike'], 'sport-i-turizam', 0.95, true),
    ('Adidas', ARRAY['adidas'], 'sport-i-turizam', 0.95, true),
    ('Under Armour', ARRAY['under armour', 'underarmour', 'ua'], 'sport-i-turizam', 0.95, true),
    ('Decathlon', ARRAY['decathlon', 'quechua', 'btwin'], 'sport-i-turizam', 0.95, true),
    ('The North Face', ARRAY['the north face', 'north face', 'tnf'], 'sport-i-turizam', 0.95, true),
    ('Columbia', ARRAY['columbia'], 'sport-i-turizam', 0.95, true),
    ('Salomon', ARRAY['salomon'], 'sport-i-turizam', 0.90, true),
    ('Mammut', ARRAY['mammut'], 'sport-i-turizam', 0.90, true);

-- === LEPOTA I ZDRAVLJE ===

-- Beauty & Health
INSERT INTO brand_category_mapping (brand_name, brand_aliases, category_slug, confidence, is_verified)
VALUES
    ('L''Oréal', ARRAY['loreal', 'l''oreal', 'l oreal'], 'lepota-i-zdravlje', 0.95, true),
    ('Nivea', ARRAY['nivea'], 'lepota-i-zdravlje', 0.95, true),
    ('Garnier', ARRAY['garnier'], 'lepota-i-zdravlje', 0.95, true),
    ('Maybelline', ARRAY['maybelline'], 'lepota-i-zdravlje', 0.95, true),
    ('MAC', ARRAY['mac', 'mac cosmetics'], 'lepota-i-zdravlje', 0.95, true),
    ('Estée Lauder', ARRAY['estee lauder', 'estée lauder'], 'lepota-i-zdravlje', 0.95, true),
    ('Clinique', ARRAY['clinique'], 'lepota-i-zdravlje', 0.90, true),
    ('Dior', ARRAY['dior', 'christian dior'], 'lepota-i-zdravlje', 0.90, true),
    ('Chanel', ARRAY['chanel'], 'lepota-i-zdravlje', 0.90, true);

-- === NAKIT I SATOVI ===

-- Watches & Jewelry
INSERT INTO brand_category_mapping (brand_name, brand_aliases, category_slug, confidence, is_verified)
VALUES
    ('Rolex', ARRAY['rolex'], 'nakit-i-satovi', 0.98, true),
    ('Omega', ARRAY['omega'], 'nakit-i-satovi', 0.98, true),
    ('Casio', ARRAY['casio', 'g-shock', 'gshock'], 'nakit-i-satovi', 0.98, true),
    ('Seiko', ARRAY['seiko'], 'nakit-i-satovi', 0.95, true),
    ('Citizen', ARRAY['citizen'], 'nakit-i-satovi', 0.95, true),
    ('Swarovski', ARRAY['swarovski'], 'nakit-i-satovi', 0.95, true),
    ('Pandora', ARRAY['pandora'], 'nakit-i-satovi', 0.95, true),
    ('Fossil', ARRAY['fossil'], 'nakit-i-satovi', 0.90, true),
    ('Michael Kors', ARRAY['michael kors', 'mk'], 'nakit-i-satovi', 0.90, true);

-- === ZA BEBE I DECU ===

-- Baby & Kids
INSERT INTO brand_category_mapping (brand_name, brand_aliases, category_slug, confidence, is_verified)
VALUES
    ('Pampers', ARRAY['pampers'], 'za-bebe-i-decu', 0.95, true),
    ('Huggies', ARRAY['huggies'], 'za-bebe-i-decu', 0.95, true),
    ('Chicco', ARRAY['chicco'], 'za-bebe-i-decu', 0.95, true),
    ('Fisher-Price', ARRAY['fisher-price', 'fisher price'], 'za-bebe-i-decu', 0.95, true),
    ('Lego', ARRAY['lego'], 'za-bebe-i-decu', 0.95, true),
    ('Mattel', ARRAY['mattel', 'barbie', 'hot wheels'], 'za-bebe-i-decu', 0.95, true),
    ('Hasbro', ARRAY['hasbro'], 'za-bebe-i-decu', 0.90, true),
    ('Graco', ARRAY['graco'], 'za-bebe-i-decu', 0.90, true);

-- Добавить timestamp для created_at
UPDATE brand_category_mapping SET created_at = NOW(), updated_at = NOW();
