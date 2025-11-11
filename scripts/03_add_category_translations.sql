-- ============================================================================
-- Translation Update: Add translations for all c2c_categories
-- Description: Adds English, Russian, and Serbian translations
-- Author: System
-- Date: 2025-11-10
-- Total categories: 77
-- ============================================================================

-- Begin transaction
BEGIN;

-- ============================================================================
-- ROOT CATEGORIES (Level 0)
-- ============================================================================

-- ID: 1001 - Electronics
UPDATE c2c_categories SET
  title_en = 'Electronics',
  title_ru = 'Электроника',
  title_sr = 'Elektronika'
WHERE slug = 'electronics';

-- ID: 1002 - Fashion
UPDATE c2c_categories SET
  title_en = 'Fashion',
  title_ru = 'Мода',
  title_sr = 'Moda'
WHERE slug = 'fashion';

-- ID: 1003 - Automotive
UPDATE c2c_categories SET
  title_en = 'Automotive',
  title_ru = 'Автомобили',
  title_sr = 'Automobili'
WHERE slug = 'automotive';

-- ID: 1004 - Real Estate
UPDATE c2c_categories SET
  title_en = 'Real Estate',
  title_ru = 'Недвижимость',
  title_sr = 'Nekretnine'
WHERE slug = 'real-estate';

-- ID: 1005 - Home & Garden
UPDATE c2c_categories SET
  title_en = 'Home & Garden',
  title_ru = 'Дом и сад',
  title_sr = 'Dom i bašta'
WHERE slug = 'home-garden';

-- ID: 1006 - Agriculture
UPDATE c2c_categories SET
  title_en = 'Agriculture',
  title_ru = 'Сельское хозяйство',
  title_sr = 'Poljoprivreda'
WHERE slug = 'agriculture';

-- ID: 1007 - Industrial
UPDATE c2c_categories SET
  title_en = 'Industrial',
  title_ru = 'Промышленность',
  title_sr = 'Industrija'
WHERE slug = 'industrial';

-- ID: 1008 - Food & Beverages
UPDATE c2c_categories SET
  title_en = 'Food & Beverages',
  title_ru = 'Еда и напитки',
  title_sr = 'Hrana i piće'
WHERE slug = 'food-beverages';

-- ID: 1009 - Services
UPDATE c2c_categories SET
  title_en = 'Services',
  title_ru = 'Услуги',
  title_sr = 'Usluge'
WHERE slug = 'services';

-- ID: 1010 - Sports & Recreation
UPDATE c2c_categories SET
  title_en = 'Sports & Recreation',
  title_ru = 'Спорт и отдых',
  title_sr = 'Sport i rekreacija'
WHERE slug = 'sports-recreation';

-- ID: 1011 - Pets
UPDATE c2c_categories SET
  title_en = 'Pets',
  title_ru = 'Животные',
  title_sr = 'Kućni ljubimci'
WHERE slug = 'pets';

-- ID: 1012 - Books & Stationery
UPDATE c2c_categories SET
  title_en = 'Books & Stationery',
  title_ru = 'Книги и канцтовары',
  title_sr = 'Knjige i kancelarija'
WHERE slug = 'books-stationery';

-- ID: 1013 - Kids & Baby
UPDATE c2c_categories SET
  title_en = 'Kids & Baby',
  title_ru = 'Детские товары',
  title_sr = 'Deca i bebe'
WHERE slug = 'kids-baby';

-- ID: 1014 - Health & Beauty
UPDATE c2c_categories SET
  title_en = 'Health & Beauty',
  title_ru = 'Здоровье и красота',
  title_sr = 'Zdravlje i lepota'
WHERE slug = 'health-beauty';

-- ID: 1015 - Hobbies & Entertainment
UPDATE c2c_categories SET
  title_en = 'Hobbies & Entertainment',
  title_ru = 'Хобби и развлечения',
  title_sr = 'Hobiji i zabava'
WHERE slug = 'hobbies-entertainment';

-- ID: 1016 - Musical Instruments
UPDATE c2c_categories SET
  title_en = 'Musical Instruments',
  title_ru = 'Музыкальные инструменты',
  title_sr = 'Muzički instrumenti'
WHERE slug = 'musical-instruments';

-- ID: 1017 - Antiques & Art
UPDATE c2c_categories SET
  title_en = 'Antiques & Art',
  title_ru = 'Антиквариат и искусство',
  title_sr = 'Antikviteti i umetnost'
WHERE slug = 'antiques-art';

-- ID: 1018 - Jobs
UPDATE c2c_categories SET
  title_en = 'Jobs',
  title_ru = 'Работа',
  title_sr = 'Poslovi'
WHERE slug = 'jobs';

-- ID: 1019 - Education
UPDATE c2c_categories SET
  title_en = 'Education',
  title_ru = 'Образование',
  title_sr = 'Obrazovanje'
WHERE slug = 'education';

-- ID: 1020 - Events & Tickets
UPDATE c2c_categories SET
  title_en = 'Events & Tickets',
  title_ru = 'События и билеты',
  title_sr = 'Dogadjaji i karte'
WHERE slug = 'events-tickets';

-- ID: 10207 - Natural Materials
UPDATE c2c_categories SET
  title_en = 'Natural Materials',
  title_ru = 'Природные материалы',
  title_sr = 'Prirodni materijali'
WHERE slug = 'natural-materials';

-- ID: 10233 - Test Category
UPDATE c2c_categories SET
  title_en = 'Test Category',
  title_ru = 'Тестовая категория',
  title_sr = 'Test Kategorija'
WHERE slug = 'test-category';

-- ID: 10234 - Test Category Manual
UPDATE c2c_categories SET
  title_en = 'Test Category Manual',
  title_ru = 'Тестовая категория (ручная)',
  title_sr = 'Test Kategorija (ručna)'
WHERE slug = 'test-category-manual';

-- ============================================================================
-- ELECTRONICS SUBCATEGORIES (Parent: 1001)
-- ============================================================================

-- ID: 1101 - Smartphones
UPDATE c2c_categories SET
  title_en = 'Smartphones',
  title_ru = 'Смартфоны',
  title_sr = 'Pametni telefoni'
WHERE slug = 'smartphones';

-- ID: 1102 - Computers
UPDATE c2c_categories SET
  title_en = 'Computers',
  title_ru = 'Компьютеры',
  title_sr = 'Računari'
WHERE slug = 'computers';

-- ID: 1103 - TV & Audio
UPDATE c2c_categories SET
  title_en = 'TV & Audio',
  title_ru = 'ТВ и аудио',
  title_sr = 'TV i audio'
WHERE slug = 'tv-audio';

-- ID: 1104 - Home Appliances
UPDATE c2c_categories SET
  title_en = 'Home Appliances',
  title_ru = 'Бытовая техника',
  title_sr = 'Kućni aparati'
WHERE slug = 'home-appliances';

-- ID: 1105 - Gaming Consoles
UPDATE c2c_categories SET
  title_en = 'Gaming Consoles',
  title_ru = 'Игровые консоли',
  title_sr = 'Gaming konzole'
WHERE slug = 'gaming-consoles';

-- ID: 1106 - Photo & Video
UPDATE c2c_categories SET
  title_en = 'Photo & Video',
  title_ru = 'Фото и видео',
  title_sr = 'Foto i video'
WHERE slug = 'photo-video';

-- ID: 1107 - Smart Home
UPDATE c2c_categories SET
  title_en = 'Smart Home',
  title_ru = 'Умный дом',
  title_sr = 'Pametna kuća'
WHERE slug = 'smart-home';

-- ID: 1108 - Electronics Accessories
UPDATE c2c_categories SET
  title_en = 'Electronics Accessories',
  title_ru = 'Аксессуары для электроники',
  title_sr = 'Elektronski dodaci'
WHERE slug = 'electronics-accessories';

-- ID: 2006 - Photo (subcategory of Photo & Video)
UPDATE c2c_categories SET
  title_en = 'Photo',
  title_ru = 'Фото',
  title_sr = 'Foto'
WHERE slug = 'photo';

-- ID: 2007 - WiFi Routers (subcategory of Electronics Accessories)
UPDATE c2c_categories SET
  title_en = 'WiFi Routers',
  title_ru = 'WiFi роутеры',
  title_sr = 'WiFi ruteri'
WHERE slug = 'wifi-routery';

-- ============================================================================
-- FASHION SUBCATEGORIES (Parent: 1002)
-- ============================================================================

-- ID: 1202 - Women's Clothing
UPDATE c2c_categories SET
  title_en = 'Women''s Clothing',
  title_ru = 'Женская одежда',
  title_sr = 'Ženska odeća'
WHERE slug = 'womens-clothing';

-- ID: 1207 - Watches
UPDATE c2c_categories SET
  title_en = 'Watches',
  title_ru = 'Часы',
  title_sr = 'Satovi'
WHERE slug = 'watches';

-- ============================================================================
-- AUTOMOTIVE SUBCATEGORIES (Parent: 1003)
-- ============================================================================

-- ID: 1301 - Cars
UPDATE c2c_categories SET
  title_en = 'Cars',
  title_ru = 'Легковые автомобили',
  title_sr = 'Lični automobili'
WHERE slug = 'cars';

-- ID: 1302 - Motorcycles
UPDATE c2c_categories SET
  title_en = 'Motorcycles',
  title_ru = 'Мотоциклы',
  title_sr = 'Motocikli'
WHERE slug = 'motorcycles';

-- ID: 1303 - Auto Parts
UPDATE c2c_categories SET
  title_en = 'Auto Parts',
  title_ru = 'Автозапчасти',
  title_sr = 'Auto delovi'
WHERE slug = 'auto-parts';

-- ID: 10100 - Domestic Production
UPDATE c2c_categories SET
  title_en = 'Domestic Production',
  title_ru = 'Отечественное производство',
  title_sr = 'Domaća proizvodnja'
WHERE slug = 'domaca-proizvodnja';

-- ID: 10110 - Imported Vehicles
UPDATE c2c_categories SET
  title_en = 'Imported Vehicles',
  title_ru = 'Импортные автомобили',
  title_sr = 'Uvozna vozila'
WHERE slug = 'uvozna-vozila';

-- Cars subcategories
-- ID: 10170 - Electric Cars
UPDATE c2c_categories SET
  title_en = 'Electric Cars',
  title_ru = 'Электромобили',
  title_sr = 'Električni automobili'
WHERE slug = 'elektricni-automobili';

-- ID: 10171 - Hybrid Cars
UPDATE c2c_categories SET
  title_en = 'Hybrid Cars',
  title_ru = 'Гибридные автомобили',
  title_sr = 'Hibridni automobili'
WHERE slug = 'hibridni-automobili';

-- ID: 10172 - Luxury Cars
UPDATE c2c_categories SET
  title_en = 'Luxury Cars',
  title_ru = 'Роскошные автомобили',
  title_sr = 'Luksuzni automobili'
WHERE slug = 'luksuzni-automobili';

-- ID: 10173 - Sports Cars
UPDATE c2c_categories SET
  title_en = 'Sports Cars',
  title_ru = 'Спортивные автомобили',
  title_sr = 'Sportski automobili'
WHERE slug = 'sportski-automobili';

-- ID: 10174 - SUV Vehicles
UPDATE c2c_categories SET
  title_en = 'SUV Vehicles',
  title_ru = 'Внедорожники',
  title_sr = 'SUV vozila'
WHERE slug = 'suv-vozila';

-- ID: 10175 - Station Wagons
UPDATE c2c_categories SET
  title_en = 'Station Wagons',
  title_ru = 'Универсалы',
  title_sr = 'Karavan vozila'
WHERE slug = 'karavan-vozila';

-- ID: 10176 - City Cars
UPDATE c2c_categories SET
  title_en = 'City Cars',
  title_ru = 'Городские автомобили',
  title_sr = 'Gradski automobili'
WHERE slug = 'gradski-automobili';

-- ID: 10177 - Camper Vehicles
UPDATE c2c_categories SET
  title_en = 'Camper Vehicles',
  title_ru = 'Дома на колесах',
  title_sr = 'Kamp vozila'
WHERE slug = 'kamp-vozila';

-- Domestic Production subcategories
-- ID: 10102 - Yugo Classics
UPDATE c2c_categories SET
  title_en = 'Yugo Classics',
  title_ru = 'Классические Yugo',
  title_sr = 'Yugo klasici'
WHERE slug = 'yugo-klasici';

-- ID: 10103 - FAP Trucks
UPDATE c2c_categories SET
  title_en = 'FAP Trucks',
  title_ru = 'Грузовики FAP',
  title_sr = 'FAP kamioni'
WHERE slug = 'fap-kamioni';

-- ID: 10104 - IMT Tractors
UPDATE c2c_categories SET
  title_en = 'IMT Tractors',
  title_ru = 'Тракторы IMT',
  title_sr = 'IMT traktori'
WHERE slug = 'imt-traktori';

-- Imported Vehicles subcategories
-- ID: 10111 - EU Import
UPDATE c2c_categories SET
  title_en = 'EU Import',
  title_ru = 'Импорт из ЕС',
  title_sr = 'EU uvoz'
WHERE slug = 'eu-uvoz';

-- ID: 10112 - Swiss Import
UPDATE c2c_categories SET
  title_en = 'Swiss Import',
  title_ru = 'Швейцарский импорт',
  title_sr = 'Švajcarski uvoz'
WHERE slug = 'svajcarski-uvoz';

-- ID: 10113 - Vehicles with Foreign Plates
UPDATE c2c_categories SET
  title_en = 'Vehicles with Foreign Plates',
  title_ru = 'Автомобили с иностранными номерами',
  title_sr = 'Vozila sa stranim tablicama'
WHERE slug = 'vozila-sa-stranim-tablicama';

-- Motorcycles subcategories
-- ID: 10180 - Sport Bikes
UPDATE c2c_categories SET
  title_en = 'Sport Bikes',
  title_ru = 'Спортбайки',
  title_sr = 'Sportski motocikli'
WHERE slug = 'sportski-motocikli';

-- Auto Parts subcategories
-- ID: 1311 - Transmission & Parts
UPDATE c2c_categories SET
  title_en = 'Transmission & Parts',
  title_ru = 'Трансмиссия и запчасти',
  title_sr = 'Transmisija i delovi'
WHERE slug = 'transmission-parts';

-- ID: 10190 - Batteries & Chargers
UPDATE c2c_categories SET
  title_en = 'Batteries & Chargers',
  title_ru = 'Аккумуляторы и зарядные устройства',
  title_sr = 'Akumulatori i punjači'
WHERE slug = 'akumulatori-i-punjaci';

-- ID: 10191 - Audio & Video Equipment
UPDATE c2c_categories SET
  title_en = 'Audio & Video Equipment',
  title_ru = 'Аудио и видео оборудование',
  title_sr = 'Audio i video oprema'
WHERE slug = 'audio-i-video-oprema';

-- ID: 10192 - GPS & Navigation
UPDATE c2c_categories SET
  title_en = 'GPS & Navigation',
  title_ru = 'GPS и навигация',
  title_sr = 'GPS i navigacija'
WHERE slug = 'gps-i-navigacija';

-- ID: 10193 - Alarm Systems
UPDATE c2c_categories SET
  title_en = 'Alarm Systems',
  title_ru = 'Сигнализационные системы',
  title_sr = 'Alarmni sistemi'
WHERE slug = 'alarmni-sistemi';

-- ID: 10194 - Tuning Parts
UPDATE c2c_categories SET
  title_en = 'Tuning Parts',
  title_ru = 'Тюнинг запчасти',
  title_sr = 'Tuning delovi'
WHERE slug = 'tuning-delovi';

-- ID: 10195 - Parts for Oldtimers
UPDATE c2c_categories SET
  title_en = 'Parts for Oldtimers',
  title_ru = 'Запчасти для раритетов',
  title_sr = 'Delovi za oldtajmere'
WHERE slug = 'delovi-za-oldtajmere';

-- ============================================================================
-- REAL ESTATE SUBCATEGORIES (Parent: 1004)
-- ============================================================================

-- ID: 1401 - Apartments
UPDATE c2c_categories SET
  title_en = 'Apartments',
  title_ru = 'Квартиры',
  title_sr = 'Stanovi'
WHERE slug = 'apartments';

-- ID: 1402 - Houses
UPDATE c2c_categories SET
  title_en = 'Houses',
  title_ru = 'Дома',
  title_sr = 'Kuće'
WHERE slug = 'houses';

-- ID: 1403 - Land
UPDATE c2c_categories SET
  title_en = 'Land',
  title_ru = 'Земельные участки',
  title_sr = 'Zemljište'
WHERE slug = 'land';

-- ID: 1404 - Commercial Real Estate
UPDATE c2c_categories SET
  title_en = 'Commercial Real Estate',
  title_ru = 'Коммерческая недвижимость',
  title_sr = 'Poslovne nekretnine'
WHERE slug = 'commercial-real-estate';

-- ============================================================================
-- HOME & GARDEN SUBCATEGORIES (Parent: 1005)
-- ============================================================================

-- ID: 1501 - Furniture
UPDATE c2c_categories SET
  title_en = 'Furniture',
  title_ru = 'Мебель',
  title_sr = 'Nameštaj'
WHERE slug = 'furniture';

-- ID: 1504 - Building Materials
UPDATE c2c_categories SET
  title_en = 'Building Materials',
  title_ru = 'Строительные материалы',
  title_sr = 'Građevinski materijal'
WHERE slug = 'building-materials';

-- ============================================================================
-- AGRICULTURE SUBCATEGORIES (Parent: 1006)
-- ============================================================================

-- ID: 1601 - Farm Machinery
UPDATE c2c_categories SET
  title_en = 'Farm Machinery',
  title_ru = 'Сельскохозяйственная техника',
  title_sr = 'Poljoprivredne mašine'
WHERE slug = 'farm-machinery';

-- ID: 1602 - Seeds & Fertilizers
UPDATE c2c_categories SET
  title_en = 'Seeds & Fertilizers',
  title_ru = 'Семена и удобрения',
  title_sr = 'Seme i đubriva'
WHERE slug = 'seeds-fertilizers';

-- ID: 1603 - Livestock
UPDATE c2c_categories SET
  title_en = 'Livestock',
  title_ru = 'Скот',
  title_sr = 'Stoka'
WHERE slug = 'livestock';

-- ID: 1604 - Farm Products
UPDATE c2c_categories SET
  title_en = 'Farm Products',
  title_ru = 'Сельхозпродукция',
  title_sr = 'Poljoprivredni proizvodi'
WHERE slug = 'farm-products';

-- ============================================================================
-- INDUSTRIAL SUBCATEGORIES (Parent: 1007)
-- ============================================================================

-- ID: 10206 - Construction Materials
UPDATE c2c_categories SET
  title_en = 'Construction Materials',
  title_ru = 'Строительные материалы',
  title_sr = 'Građevinski materijali'
WHERE slug = 'construction-materials';

-- ID: 10232 - Construction Tools
UPDATE c2c_categories SET
  title_en = 'Construction Tools',
  title_ru = 'Строительные инструменты',
  title_sr = 'Građevinski alati'
WHERE slug = 'construction-tools';

-- ============================================================================
-- HOBBIES & ENTERTAINMENT SUBCATEGORIES (Parent: 1015)
-- ============================================================================

-- ID: 10202 - Toys
UPDATE c2c_categories SET
  title_en = 'Toys',
  title_ru = 'Игрушки',
  title_sr = 'Igračke'
WHERE slug = 'toys';

-- ID: 10203 - Puzzles
UPDATE c2c_categories SET
  title_en = 'Puzzles',
  title_ru = 'Пазлы',
  title_sr = 'Slagalice'
WHERE slug = 'puzzles';

-- ID: 10205 - Collectibles
UPDATE c2c_categories SET
  title_en = 'Collectibles',
  title_ru = 'Коллекционирование',
  title_sr = 'Kolekcionarstvo'
WHERE slug = 'collectibles';

-- ============================================================================
-- VERIFICATION QUERIES
-- ============================================================================

-- Show translation statistics
SELECT
  COUNT(*) as total_categories,
  COUNT(title_en) as with_english,
  COUNT(title_ru) as with_russian,
  COUNT(title_sr) as with_serbian,
  COUNT(CASE WHEN title_en IS NULL OR title_ru IS NULL OR title_sr IS NULL THEN 1 END) as missing_translations
FROM c2c_categories;

-- Show any categories still missing translations
SELECT
  id,
  name,
  slug,
  title_en,
  title_ru,
  title_sr
FROM c2c_categories
WHERE title_en IS NULL OR title_ru IS NULL OR title_sr IS NULL
ORDER BY id;

-- Commit transaction
COMMIT;

-- ============================================================================
-- SUCCESS MESSAGE
-- ============================================================================
\echo '✅ Translation update completed successfully!'
\echo '📊 Check the output above to verify all categories have translations.'
