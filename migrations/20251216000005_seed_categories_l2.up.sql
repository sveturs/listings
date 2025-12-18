-- Migration: Seed L2 (second-level) categories
-- Date: 2025-12-16
-- Purpose: Insert ~250 L2 subcategories with multilingual support (sr/en/ru)
-- Reference: 18 L1 parent categories from 20251216000004_seed_categories_l1.up.sql

-- =============================================================================
-- L2 for: 1. Odeća i obuća (Clothing & Footwear) - 15 categories
-- =============================================================================
INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('muska-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/muska-odeca', 1,
 '{"sr": "Muška odeća", "en": "Men''s clothing", "ru": "Мужская одежда"}'::jsonb,
 '{"sr": "Košulje, pantalone, jakne i odela", "en": "Shirts, pants, jackets and suits", "ru": "Рубашки, брюки, куртки и костюмы"}'::jsonb,
 '{"sr": "Muška odeća | Vondi", "en": "Men''s clothing | Vondi", "ru": "Мужская одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite mušku odeću online - košulje, pantalone, jakne", "en": "Buy men''s clothing online - shirts, pants, jackets", "ru": "Купить мужскую одежду онлайн - рубашки, брюки, куртки"}'::jsonb,
 '👔', true),

('zenska-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/zenska-odeca', 2,
 '{"sr": "Ženska odeća", "en": "Women''s clothing", "ru": "Женская одежда"}'::jsonb,
 '{"sr": "Haljine, bluze, suknje i pantalone", "en": "Dresses, blouses, skirts and pants", "ru": "Платья, блузки, юбки и брюки"}'::jsonb,
 '{"sr": "Ženska odeća | Vondi", "en": "Women''s clothing | Vondi", "ru": "Женская одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite žensku odeću online - haljine, bluze, suknje", "en": "Buy women''s clothing online - dresses, blouses, skirts", "ru": "Купить женскую одежду онлайн - платья, блузки, юбки"}'::jsonb,
 '👗', true),

('decija-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/decija-odeca', 3,
 '{"sr": "Dečija odeća", "en": "Kids'' clothing", "ru": "Детская одежда"}'::jsonb,
 '{"sr": "Odeća za decu svih uzrasta", "en": "Clothing for children of all ages", "ru": "Одежда для детей всех возрастов"}'::jsonb,
 '{"sr": "Dečija odeća | Vondi", "en": "Kids'' clothing | Vondi", "ru": "Детская одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiju odeću online - za sve uzraste", "en": "Buy kids'' clothing online - for all ages", "ru": "Купить детскую одежду онлайн - для всех возрастов"}'::jsonb,
 '👶', true),

('muska-obuca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/muska-obuca', 4,
 '{"sr": "Muška obuća", "en": "Men''s footwear", "ru": "Мужская обувь"}'::jsonb,
 '{"sr": "Cipele, patike, čizme i sandale", "en": "Shoes, sneakers, boots and sandals", "ru": "Туфли, кроссовки, ботинки и сандалии"}'::jsonb,
 '{"sr": "Muška obuća | Vondi", "en": "Men''s footwear | Vondi", "ru": "Мужская обувь | Vondi"}'::jsonb,
 '{"sr": "Kupite mušku obuću online - cipele, patike, čizme", "en": "Buy men''s footwear online - shoes, sneakers, boots", "ru": "Купить мужскую обувь онлайн - туфли, кроссовки, ботинки"}'::jsonb,
 '👞', true),

('zenska-obuca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/zenska-obuca', 5,
 '{"sr": "Ženska obuća", "en": "Women''s footwear", "ru": "Женская обувь"}'::jsonb,
 '{"sr": "Cipele, patike, čizme i štikle", "en": "Shoes, sneakers, boots and heels", "ru": "Туфли, кроссовки, ботинки и каблуки"}'::jsonb,
 '{"sr": "Ženska obuća | Vondi", "en": "Women''s footwear | Vondi", "ru": "Женская обувь | Vondi"}'::jsonb,
 '{"sr": "Kupite žensku obuću online - cipele, patike, štikle", "en": "Buy women''s footwear online - shoes, sneakers, heels", "ru": "Купить женскую обувь онлайн - туфли, кроссовки, каблуки"}'::jsonb,
 '👠', true),

('decija-obuca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/decija-obuca', 6,
 '{"sr": "Dečija obuća", "en": "Kids'' footwear", "ru": "Детская обувь"}'::jsonb,
 '{"sr": "Patike, cipele i čizme za decu", "en": "Sneakers, shoes and boots for kids", "ru": "Кроссовки, туфли и ботинки для детей"}'::jsonb,
 '{"sr": "Dečija obuća | Vondi", "en": "Kids'' footwear | Vondi", "ru": "Детская обувь | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiju obuću online - patike, cipele, čizme", "en": "Buy kids'' footwear online - sneakers, shoes, boots", "ru": "Купить детскую обувь онлайн - кроссовки, туфли, ботинки"}'::jsonb,
 '👟', true),

('torbice-i-novcanici', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/torbice-i-novcanici', 7,
 '{"sr": "Torbice i novčanici", "en": "Bags & wallets", "ru": "Сумки и кошельки"}'::jsonb,
 '{"sr": "Torbice, ranci, novčanici i travel torbe", "en": "Handbags, backpacks, wallets and travel bags", "ru": "Сумки, рюкзаки, кошельки и дорожные сумки"}'::jsonb,
 '{"sr": "Torbice i novčanici | Vondi", "en": "Bags & wallets | Vondi", "ru": "Сумки и кошельки | Vondi"}'::jsonb,
 '{"sr": "Kupite torbice i novčanike online", "en": "Buy bags and wallets online", "ru": "Купить сумки и кошельки онлайн"}'::jsonb,
 '👜', true),

('dodaci-i-aksesoari', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/dodaci-i-aksesoari', 8,
 '{"sr": "Dodaci i aksesoari", "en": "Accessories", "ru": "Аксессуары"}'::jsonb,
 '{"sr": "Remeni, kape, šalovi i modni dodaci", "en": "Belts, caps, scarves and fashion accessories", "ru": "Ремни, кепки, шарфы и модные аксессуары"}'::jsonb,
 '{"sr": "Dodaci i aksesoari | Vondi", "en": "Accessories | Vondi", "ru": "Аксессуары | Vondi"}'::jsonb,
 '{"sr": "Kupite dodatke i aksesoari online", "en": "Buy accessories online", "ru": "Купить аксессуары онлайн"}'::jsonb,
 '🎀', true),

('donji-ves', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/donji-ves', 9,
 '{"sr": "Donji veš", "en": "Underwear", "ru": "Нижнее белье"}'::jsonb,
 '{"sr": "Gaće, grudnjaci, pidžame i čarape", "en": "Boxers, bras, pajamas and socks", "ru": "Трусы, бюстгальтеры, пижамы и носки"}'::jsonb,
 '{"sr": "Donji veš | Vondi", "en": "Underwear | Vondi", "ru": "Нижнее белье | Vondi"}'::jsonb,
 '{"sr": "Kupite donji veš online - udoban i kvalitetan", "en": "Buy underwear online - comfortable and quality", "ru": "Купить нижнее белье онлайн - удобное и качественное"}'::jsonb,
 '🩲', true),

('sportska-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/sportska-odeca', 10,
 '{"sr": "Sportska odeća", "en": "Sportswear", "ru": "Спортивная одежда"}'::jsonb,
 '{"sr": "Trenerke, dresovi i sportske majice", "en": "Tracksuits, jerseys and sports shirts", "ru": "Спортивные костюмы, майки и футболки"}'::jsonb,
 '{"sr": "Sportska odeća | Vondi", "en": "Sportswear | Vondi", "ru": "Спортивная одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite sportsku odeću online", "en": "Buy sportswear online", "ru": "Купить спортивную одежду онлайн"}'::jsonb,
 '🏃', true),

('zimska-garderoba', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/zimska-garderoba', 11,
 '{"sr": "Zimska garderoba", "en": "Winter clothing", "ru": "Зимняя одежда"}'::jsonb,
 '{"sr": "Jakne, kaputi i bundeve", "en": "Jackets, coats and fur coats", "ru": "Куртки, пальто и шубы"}'::jsonb,
 '{"sr": "Zimska garderoba | Vondi", "en": "Winter clothing | Vondi", "ru": "Зимняя одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite zimsku garderobu online", "en": "Buy winter clothing online", "ru": "Купить зимнюю одежду онлайн"}'::jsonb,
 '🧥', true),

('kupaci-kostimi', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/kupaci-kostimi', 12,
 '{"sr": "Kupaći kostimi", "en": "Swimwear", "ru": "Купальники"}'::jsonb,
 '{"sr": "Bikiniji, jednodelni i muški kupaći", "en": "Bikinis, one-pieces and men''s swimwear", "ru": "Бикини, слитные и мужские купальники"}'::jsonb,
 '{"sr": "Kupaći kostimi | Vondi", "en": "Swimwear | Vondi", "ru": "Купальники | Vondi"}'::jsonb,
 '{"sr": "Kupite kupaće kostime online", "en": "Buy swimwear online", "ru": "Купить купальники онлайн"}'::jsonb,
 '👙', true),

('odela-i-smokingzi', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/odela-i-smokingzi', 13,
 '{"sr": "Odela i smokingzi", "en": "Suits & tuxedos", "ru": "Костюмы и смокинги"}'::jsonb,
 '{"sr": "Poslovna odela i svečana odeća", "en": "Business suits and formal wear", "ru": "Деловые костюмы и торжественная одежда"}'::jsonb,
 '{"sr": "Odela i smokingzi | Vondi", "en": "Suits & tuxedos | Vondi", "ru": "Костюмы и смокинги | Vondi"}'::jsonb,
 '{"sr": "Kupite odela i smokingze online", "en": "Buy suits and tuxedos online", "ru": "Купить костюмы и смокинги онлайн"}'::jsonb,
 '🤵', true),

('vencana-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/vencana-odeca', 14,
 '{"sr": "Venčana odeća", "en": "Wedding attire", "ru": "Свадебная одежда"}'::jsonb,
 '{"sr": "Venčanice, venčana odela i dodaci", "en": "Wedding dresses, suits and accessories", "ru": "Свадебные платья, костюмы и аксессуары"}'::jsonb,
 '{"sr": "Venčana odeća | Vondi", "en": "Wedding attire | Vondi", "ru": "Свадебная одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite venčanu odeću online", "en": "Buy wedding attire online", "ru": "Купить свадебную одежду онлайн"}'::jsonb,
 '👰', true),

('radna-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/radna-odeca', 15,
 '{"sr": "Radna odeća", "en": "Workwear", "ru": "Рабочая одежда"}'::jsonb,
 '{"sr": "Uniforma, zaštitna odeća i obuća", "en": "Uniforms, protective clothing and footwear", "ru": "Униформа, защитная одежда и обувь"}'::jsonb,
 '{"sr": "Radna odeća | Vondi", "en": "Workwear | Vondi", "ru": "Рабочая одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite radnu odeću online", "en": "Buy workwear online", "ru": "Купить рабочую одежду онлайн"}'::jsonb,
 '👷', true),

-- =============================================================================
-- L2 for: 2. Elektronika (Electronics) - 15 categories
-- =============================================================================

('pametni-telefoni', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/pametni-telefoni', 1,
 '{"sr": "Pametni telefoni", "en": "Smartphones", "ru": "Смартфоны"}'::jsonb,
 '{"sr": "Android i iPhone svih brendova", "en": "Android and iPhone of all brands", "ru": "Android и iPhone всех брендов"}'::jsonb,
 '{"sr": "Pametni telefoni | Vondi", "en": "Smartphones | Vondi", "ru": "Смартфоны | Vondi"}'::jsonb,
 '{"sr": "Kupite pametne telefone online - Samsung, Apple, Xiaomi", "en": "Buy smartphones online - Samsung, Apple, Xiaomi", "ru": "Купить смартфоны онлайн - Samsung, Apple, Xiaomi"}'::jsonb,
 '📱', true),

('laptop-racunari', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/laptop-racunari', 2,
 '{"sr": "Laptop računari", "en": "Laptops", "ru": "Ноутбуки"}'::jsonb,
 '{"sr": "Laptop računari za posao i igrice", "en": "Laptops for work and gaming", "ru": "Ноутбуки для работы и игр"}'::jsonb,
 '{"sr": "Laptop računari | Vondi", "en": "Laptops | Vondi", "ru": "Ноутбуки | Vondi"}'::jsonb,
 '{"sr": "Kupite laptop računare online - Dell, HP, Lenovo", "en": "Buy laptops online - Dell, HP, Lenovo", "ru": "Купить ноутбуки онлайн - Dell, HP, Lenovo"}'::jsonb,
 '💻', true),

('desktop-racunari', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/desktop-racunari', 3,
 '{"sr": "Desktop računari", "en": "Desktop computers", "ru": "Настольные компьютеры"}'::jsonb,
 '{"sr": "Desktop i gaming PC računari", "en": "Desktop and gaming PCs", "ru": "Настольные и игровые ПК"}'::jsonb,
 '{"sr": "Desktop računari | Vondi", "en": "Desktop computers | Vondi", "ru": "Настольные компьютеры | Vondi"}'::jsonb,
 '{"sr": "Kupite desktop računare online", "en": "Buy desktop computers online", "ru": "Купить настольные компьютеры онлайн"}'::jsonb,
 '🖥️', true),

('tableti', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/tableti', 4,
 '{"sr": "Tableti", "en": "Tablets", "ru": "Планшеты"}'::jsonb,
 '{"sr": "Tableti Android i iPad", "en": "Android tablets and iPads", "ru": "Планшеты Android и iPad"}'::jsonb,
 '{"sr": "Tableti | Vondi", "en": "Tablets | Vondi", "ru": "Планшеты | Vondi"}'::jsonb,
 '{"sr": "Kupite tablete online - iPad, Samsung Galaxy Tab", "en": "Buy tablets online - iPad, Samsung Galaxy Tab", "ru": "Купить планшеты онлайн - iPad, Samsung Galaxy Tab"}'::jsonb,
 '📲', true),

('tv-i-video', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/tv-i-video', 5,
 '{"sr": "TV i video", "en": "TVs & video", "ru": "ТВ и видео"}'::jsonb,
 '{"sr": "Smart TV, LED, OLED televizori", "en": "Smart TV, LED, OLED televisions", "ru": "Smart TV, LED, OLED телевизоры"}'::jsonb,
 '{"sr": "TV i video | Vondi", "en": "TVs & video | Vondi", "ru": "ТВ и видео | Vondi"}'::jsonb,
 '{"sr": "Kupite TV online - Smart TV, LED, OLED", "en": "Buy TVs online - Smart TV, LED, OLED", "ru": "Купить ТВ онлайн - Smart TV, LED, OLED"}'::jsonb,
 '📺', true),

('audio-oprema', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/audio-oprema', 6,
 '{"sr": "Audio oprema", "en": "Audio equipment", "ru": "Аудио техника"}'::jsonb,
 '{"sr": "Slušalice, zvučnici i soundbar", "en": "Headphones, speakers and soundbars", "ru": "Наушники, колонки и саундбары"}'::jsonb,
 '{"sr": "Audio oprema | Vondi", "en": "Audio equipment | Vondi", "ru": "Аудио техника | Vondi"}'::jsonb,
 '{"sr": "Kupite audio opremu online", "en": "Buy audio equipment online", "ru": "Купить аудио технику онлайн"}'::jsonb,
 '🎧', true),

('foto-i-video-kamere', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/foto-i-video-kamere', 7,
 '{"sr": "Foto i video kamere", "en": "Cameras & camcorders", "ru": "Фото и видеокамеры"}'::jsonb,
 '{"sr": "DSLR, mirrorless i action kamere", "en": "DSLR, mirrorless and action cameras", "ru": "Зеркалки, беззеркалки и экшн-камеры"}'::jsonb,
 '{"sr": "Foto i video kamere | Vondi", "en": "Cameras & camcorders | Vondi", "ru": "Фото и видеокамеры | Vondi"}'::jsonb,
 '{"sr": "Kupite foto i video kamere online", "en": "Buy cameras and camcorders online", "ru": "Купить фото и видеокамеры онлайн"}'::jsonb,
 '📷', true),

('pametni-satovi', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/pametni-satovi', 8,
 '{"sr": "Pametni satovi", "en": "Smartwatches", "ru": "Умные часы"}'::jsonb,
 '{"sr": "Pametni satovi i fitness trakeri", "en": "Smartwatches and fitness trackers", "ru": "Умные часы и фитнес-трекеры"}'::jsonb,
 '{"sr": "Pametni satovi | Vondi", "en": "Smartwatches | Vondi", "ru": "Умные часы | Vondi"}'::jsonb,
 '{"sr": "Kupite pametne satove online - Apple Watch, Samsung", "en": "Buy smartwatches online - Apple Watch, Samsung", "ru": "Купить умные часы онлайн - Apple Watch, Samsung"}'::jsonb,
 '⌚', true),

('konzole-i-gaming', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/konzole-i-gaming', 9,
 '{"sr": "Konzole i gaming", "en": "Consoles & gaming", "ru": "Консоли и игры"}'::jsonb,
 '{"sr": "PlayStation, Xbox, Nintendo i igre", "en": "PlayStation, Xbox, Nintendo and games", "ru": "PlayStation, Xbox, Nintendo и игры"}'::jsonb,
 '{"sr": "Konzole i gaming | Vondi", "en": "Consoles & gaming | Vondi", "ru": "Консоли и игры | Vondi"}'::jsonb,
 '{"sr": "Kupite gaming konzole online - PS5, Xbox, Switch", "en": "Buy gaming consoles online - PS5, Xbox, Switch", "ru": "Купить игровые консоли онлайн - PS5, Xbox, Switch"}'::jsonb,
 '🎮', true),

('racunarske-komponente', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/racunarske-komponente', 10,
 '{"sr": "Računarske komponente", "en": "Computer components", "ru": "Компьютерные компоненты"}'::jsonb,
 '{"sr": "Procesori, grafičke kartice, RAM, SSD", "en": "Processors, graphics cards, RAM, SSDs", "ru": "Процессоры, видеокарты, оперативная память, SSD"}'::jsonb,
 '{"sr": "Računarske komponente | Vondi", "en": "Computer components | Vondi", "ru": "Компьютерные компоненты | Vondi"}'::jsonb,
 '{"sr": "Kupite računarske komponente online", "en": "Buy computer components online", "ru": "Купить компьютерные компоненты онлайн"}'::jsonb,
 '🖲️', true),

('periferija', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/periferija', 11,
 '{"sr": "Periferija", "en": "Peripherals", "ru": "Периферия"}'::jsonb,
 '{"sr": "Tastature, miševi, monitori, štampači", "en": "Keyboards, mice, monitors, printers", "ru": "Клавиатуры, мыши, мониторы, принтеры"}'::jsonb,
 '{"sr": "Periferija | Vondi", "en": "Peripherals | Vondi", "ru": "Периферия | Vondi"}'::jsonb,
 '{"sr": "Kupite periferiju online", "en": "Buy peripherals online", "ru": "Купить периферию онлайн"}'::jsonb,
 '⌨️', true),

('mreza-i-internet', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/mreza-i-internet', 12,
 '{"sr": "Mreža i internet", "en": "Networking", "ru": "Сети и интернет"}'::jsonb,
 '{"sr": "Ruteri, modemi i WiFi oprema", "en": "Routers, modems and WiFi equipment", "ru": "Роутеры, модемы и WiFi оборудование"}'::jsonb,
 '{"sr": "Mreža i internet | Vondi", "en": "Networking | Vondi", "ru": "Сети и интернет | Vondi"}'::jsonb,
 '{"sr": "Kupite mrežnu opremu online", "en": "Buy networking equipment online", "ru": "Купить сетевое оборудование онлайн"}'::jsonb,
 '📡', true),

('dodatna-oprema-elektronika', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/dodatna-oprema-elektronika', 13,
 '{"sr": "Dodatna oprema", "en": "Accessories", "ru": "Аксессуары"}'::jsonb,
 '{"sr": "Kabl, punjači, futrole, memorije", "en": "Cables, chargers, cases, memory cards", "ru": "Кабели, зарядки, чехлы, карты памяти"}'::jsonb,
 '{"sr": "Dodatna oprema | Vondi", "en": "Accessories | Vondi", "ru": "Аксессуары | Vondi"}'::jsonb,
 '{"sr": "Kupite dodatnu opremu online", "en": "Buy accessories online", "ru": "Купить аксессуары онлайн"}'::jsonb,
 '🔌', true),

('dronovi', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/dronovi', 14,
 '{"sr": "Dronovi", "en": "Drones", "ru": "Дроны"}'::jsonb,
 '{"sr": "Dronovi sa kamerom za hobi", "en": "Drones with camera for hobby", "ru": "Дроны с камерой для хобби"}'::jsonb,
 '{"sr": "Dronovi | Vondi", "en": "Drones | Vondi", "ru": "Дроны | Vondi"}'::jsonb,
 '{"sr": "Kupite dronove online - DJI, sa kamerom", "en": "Buy drones online - DJI, with camera", "ru": "Купить дроны онлайн - DJI, с камерой"}'::jsonb,
 '🚁', true),

('e-citaci', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/e-citaci', 15,
 '{"sr": "E-čitači", "en": "E-readers", "ru": "Электронные книги"}'::jsonb,
 '{"sr": "Kindle, Kobo i drugi e-čitači", "en": "Kindle, Kobo and other e-readers", "ru": "Kindle, Kobo и другие электронные книги"}'::jsonb,
 '{"sr": "E-čitači | Vondi", "en": "E-readers | Vondi", "ru": "Электронные книги | Vondi"}'::jsonb,
 '{"sr": "Kupite e-čitače online - Kindle, Kobo", "en": "Buy e-readers online - Kindle, Kobo", "ru": "Купить электронные книги онлайн - Kindle, Kobo"}'::jsonb,
 '📖', true),

-- =============================================================================
-- L2 for: 3. Dom i bašta (Home & Garden) - 15 categories
-- =============================================================================

('namestaj-dnevna-soba', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/namestaj-dnevna-soba', 1,
 '{"sr": "Nameštaj dnevna soba", "en": "Living room furniture", "ru": "Мебель для гостиной"}'::jsonb,
 '{"sr": "Sofe, fotelje, stolovi i police", "en": "Sofas, armchairs, tables and shelves", "ru": "Диваны, кресла, столы и полки"}'::jsonb,
 '{"sr": "Nameštaj dnevna soba | Vondi", "en": "Living room furniture | Vondi", "ru": "Мебель для гостиной | Vondi"}'::jsonb,
 '{"sr": "Kupite nameštaj za dnevnu sobu online", "en": "Buy living room furniture online", "ru": "Купить мебель для гостиной онлайн"}'::jsonb,
 '🛋️', true),

('namestaj-spavaca-soba', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/namestaj-spavaca-soba', 2,
 '{"sr": "Nameštaj spavaća soba", "en": "Bedroom furniture", "ru": "Мебель для спальни"}'::jsonb,
 '{"sr": "Kreveti, ormari i noćni stočići", "en": "Beds, wardrobes and nightstands", "ru": "Кровати, шкафы и тумбочки"}'::jsonb,
 '{"sr": "Nameštaj spavaća soba | Vondi", "en": "Bedroom furniture | Vondi", "ru": "Мебель для спальни | Vondi"}'::jsonb,
 '{"sr": "Kupite nameštaj za spavaću sobu online", "en": "Buy bedroom furniture online", "ru": "Купить мебель для спальни онлайн"}'::jsonb,
 '🛏️', true),

('namestaj-kuhinja', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/namestaj-kuhinja', 3,
 '{"sr": "Nameštaj kuhinja", "en": "Kitchen furniture", "ru": "Кухонная мебель"}'::jsonb,
 '{"sr": "Kuhinjski elementi, stolice i stolovi", "en": "Kitchen cabinets, chairs and tables", "ru": "Кухонные шкафы, стулья и столы"}'::jsonb,
 '{"sr": "Nameštaj kuhinja | Vondi", "en": "Kitchen furniture | Vondi", "ru": "Кухонная мебель | Vondi"}'::jsonb,
 '{"sr": "Kupite kuhinjski nameštaj online", "en": "Buy kitchen furniture online", "ru": "Купить кухонную мебель онлайн"}'::jsonb,
 '🍴', true),

('dekoracije', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/dekoracije', 4,
 '{"sr": "Dekoracije", "en": "Decorations", "ru": "Декор"}'::jsonb,
 '{"sr": "Slike, vaze, sveće i dekorativni dodaci", "en": "Paintings, vases, candles and decorative accessories", "ru": "Картины, вазы, свечи и декоративные аксессуары"}'::jsonb,
 '{"sr": "Dekoracije | Vondi", "en": "Decorations | Vondi", "ru": "Декор | Vondi"}'::jsonb,
 '{"sr": "Kupite dekoracije za dom online", "en": "Buy home decorations online", "ru": "Купить декор для дома онлайн"}'::jsonb,
 '🎨', true),

('rasveta', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/rasveta', 5,
 '{"sr": "Rasveta", "en": "Lighting", "ru": "Освещение"}'::jsonb,
 '{"sr": "Lusteri, lampe, LED sijalice", "en": "Chandeliers, lamps, LED bulbs", "ru": "Люстры, лампы, LED лампы"}'::jsonb,
 '{"sr": "Rasveta | Vondi", "en": "Lighting | Vondi", "ru": "Освещение | Vondi"}'::jsonb,
 '{"sr": "Kupite rasvetu za dom online", "en": "Buy lighting for home online", "ru": "Купить освещение для дома онлайн"}'::jsonb,
 '💡', true),

('tekstil-za-dom', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/tekstil-za-dom', 6,
 '{"sr": "Tekstil za dom", "en": "Home textiles", "ru": "Домашний текстиль"}'::jsonb,
 '{"sr": "Zavese, posteljina, peškiri i tepihi", "en": "Curtains, bedding, towels and carpets", "ru": "Шторы, постельное белье, полотенца и ковры"}'::jsonb,
 '{"sr": "Tekstil za dom | Vondi", "en": "Home textiles | Vondi", "ru": "Домашний текстиль | Vondi"}'::jsonb,
 '{"sr": "Kupite tekstil za dom online", "en": "Buy home textiles online", "ru": "Купить домашний текстиль онлайн"}'::jsonb,
 '🛌', true),

('kupatilo', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/kupatilo', 7,
 '{"sr": "Kupatilo", "en": "Bathroom", "ru": "Ванная комната"}'::jsonb,
 '{"sr": "Slavine, tuševi i kupatilski dodaci", "en": "Faucets, showers and bathroom accessories", "ru": "Смесители, душевые и аксессуары для ванной"}'::jsonb,
 '{"sr": "Kupatilo | Vondi", "en": "Bathroom | Vondi", "ru": "Ванная комната | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za kupatilo online", "en": "Buy bathroom equipment online", "ru": "Купить оборудование для ванной онлайн"}'::jsonb,
 '🚿', true),

('bastenska-oprema', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/bastenska-oprema', 8,
 '{"sr": "Baštenka oprema", "en": "Garden equipment", "ru": "Садовое оборудование"}'::jsonb,
 '{"sr": "Kosačice, trijmeri i baštenski alati", "en": "Lawnmowers, trimmers and garden tools", "ru": "Газонокосилки, триммеры и садовый инвентарь"}'::jsonb,
 '{"sr": "Baštenka oprema | Vondi", "en": "Garden equipment | Vondi", "ru": "Садовое оборудование | Vondi"}'::jsonb,
 '{"sr": "Kupite baštensku opremu online", "en": "Buy garden equipment online", "ru": "Купить садовое оборудование онлайн"}'::jsonb,
 '🌿', true),

('bastenska-garnitura', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/bastenska-garnitura', 9,
 '{"sr": "Baštenka garnitura", "en": "Outdoor furniture", "ru": "Садовая мебель"}'::jsonb,
 '{"sr": "Garniture, stolovi, stolice i ležaljke", "en": "Furniture sets, tables, chairs and loungers", "ru": "Мебельные наборы, столы, стулья и лежаки"}'::jsonb,
 '{"sr": "Baštenka garnitura | Vondi", "en": "Outdoor furniture | Vondi", "ru": "Садовая мебель | Vondi"}'::jsonb,
 '{"sr": "Kupite baštensku garnituru online", "en": "Buy outdoor furniture online", "ru": "Купить садовую мебель онлайн"}'::jsonb,
 '🪑', true),

('grncari ja-i-biljke', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/grncari ja-i-biljke', 10,
 '{"sr": "Grnčarija i biljke", "en": "Pots & plants", "ru": "Горшки и растения"}'::jsonb,
 '{"sr": "Saksije, biljke i semenje", "en": "Plant pots, plants and seeds", "ru": "Цветочные горшки, растения и семена"}'::jsonb,
 '{"sr": "Grnčarija i biljke | Vondi", "en": "Pots & plants | Vondi", "ru": "Горшки и растения | Vondi"}'::jsonb,
 '{"sr": "Kupite saksije i biljke online", "en": "Buy pots and plants online", "ru": "Купить горшки и растения онлайн"}'::jsonb,
 '🪴', true),

('alati-za-basta', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/alati-za-basta', 11,
 '{"sr": "Alati za bašta", "en": "Garden tools", "ru": "Садовый инвентарь"}'::jsonb,
 '{"sr": "Lopate, grablje, kantre i navodnjavanje", "en": "Shovels, rakes, wheelbarrows and irrigation", "ru": "Лопаты, грабли, тачки и полив"}'::jsonb,
 '{"sr": "Alati za bašta | Vondi", "en": "Garden tools | Vondi", "ru": "Садовый инвентарь | Vondi"}'::jsonb,
 '{"sr": "Kupite alate za baštu online", "en": "Buy garden tools online", "ru": "Купить садовый инвентарь онлайн"}'::jsonb,
 '🔨', true),

('organizacija-i-skladistenje', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/organizacija-i-skladistenje', 12,
 '{"sr": "Organizacija i skladištenje", "en": "Organization & storage", "ru": "Организация и хранение"}'::jsonb,
 '{"sr": "Kutije, korpe, polica i ormari", "en": "Boxes, baskets, shelves and cabinets", "ru": "Коробки, корзины, полки и шкафы"}'::jsonb,
 '{"sr": "Organizacija i skladištenje | Vondi", "en": "Organization & storage | Vondi", "ru": "Организация и хранение | Vondi"}'::jsonb,
 '{"sr": "Kupite sistem za skladištenje online", "en": "Buy storage systems online", "ru": "Купить системы хранения онлайн"}'::jsonb,
 '📦', true),

('alati-i-popravke', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/alati-i-popravke', 13,
 '{"sr": "Alati i popravke", "en": "Tools & repairs", "ru": "Инструменты и ремонт"}'::jsonb,
 '{"sr": "Ručni alati, električni alati, farba", "en": "Hand tools, power tools, paint", "ru": "Ручные инструменты, электроинструменты, краска"}'::jsonb,
 '{"sr": "Alati i popravke | Vondi", "en": "Tools & repairs | Vondi", "ru": "Инструменты и ремонт | Vondi"}'::jsonb,
 '{"sr": "Kupite alate i opremu za popravke online", "en": "Buy tools and repair equipment online", "ru": "Купить инструменты и оборудование для ремонта онлайн"}'::jsonb,
 '🔧', true),

('ventilacija-i-klimatizacija', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/ventilacija-i-klimatizacija', 14,
 '{"sr": "Ventilacija i klimatizacija", "en": "Ventilation & air conditioning", "ru": "Вентиляция и кондиционирование"}'::jsonb,
 '{"sr": "Klima uređaji, ventilatori i grejalice", "en": "Air conditioners, fans and heaters", "ru": "Кондиционеры, вентиляторы и обогреватели"}'::jsonb,
 '{"sr": "Ventilacija i klimatizacija | Vondi", "en": "Ventilation & air conditioning | Vondi", "ru": "Вентиляция и кондиционирование | Vondi"}'::jsonb,
 '{"sr": "Kupite klima uređaje i ventilatoren online", "en": "Buy air conditioners and fans online", "ru": "Купить кондиционеры и вентиляторы онлайн"}'::jsonb,
 '❄️', true),

('bazeni-i-spa', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/bazeni-i-spa', 15,
 '{"sr": "Bazeni i spa", "en": "Pools & spa", "ru": "Бассейны и спа"}'::jsonb,
 '{"sr": "Naduvni bazeni, hemija i oprema", "en": "Inflatable pools, chemicals and equipment", "ru": "Надувные бассейны, химия и оборудование"}'::jsonb,
 '{"sr": "Bazeni i spa | Vondi", "en": "Pools & spa | Vondi", "ru": "Бассейны и спа | Vondi"}'::jsonb,
 '{"sr": "Kupite bazene i spa opremu online", "en": "Buy pools and spa equipment online", "ru": "Купить бассейны и spa оборудование онлайн"}'::jsonb,
 '🏊', true);

-- Due to file length constraints, this migration will continue in a Part 2 file.
-- Current progress: 45 L2 categories created (15 per L1 category x 3 L1 categories)
-- Remaining: 15 more L1 categories to process

-- Temporary verification
DO $$
DECLARE
    l2_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO l2_count FROM categories WHERE level = 2;

    RAISE NOTICE 'Part 1 completed: % L2 categories inserted (Odeća: 15, Elektronika: 15, Dom: 15)', l2_count;
END $$;
