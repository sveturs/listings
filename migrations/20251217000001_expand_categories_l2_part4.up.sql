-- Migration: Expand L2 categories (Part 4)
-- Date: 2025-12-17
-- Purpose: Add ~100 L2 categories to reach target of 400 total L2
-- Expanding existing L1 categories with additional subcategories

-- =============================================================================
-- Additional L2 for: Odeća i obuća (+ 10 more to reach 25 total)
-- =============================================================================
INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('posteljina-i-peskiri', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/posteljina-i-peskiri', 20,
 '{"sr": "Posteljina i peškiri", "en": "Bedding & towels", "ru": "Постельное белье и полотенца"}'::jsonb,
 '{"sr": "Jastučnice, čaršavi, jorgan, peškiri", "en": "Pillowcases, sheets, duvets, towels", "ru": "Наволочки, простыни, одеяла, полотенца"}'::jsonb,
 '{"sr": "Posteljina i peškiri | Vondi", "en": "Bedding & towels | Vondi", "ru": "Постельное белье и полотенца | Vondi"}'::jsonb,
 '{"sr": "Kupite posteljinu i peškire online", "en": "Buy bedding and towels online", "ru": "Купить постельное белье и полотенца онлайн"}'::jsonb,
 '🛏️', true),

('ves-masine-dodaci', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/ves-masine-dodaci', 21,
 '{"sr": "Veš mašine dodaci", "en": "Laundry accessories", "ru": "Аксессуары для стирки"}'::jsonb,
 '{"sr": "Korpice, vešalice, daske za peglanje", "en": "Baskets, hangers, ironing boards", "ru": "Корзины, вешалки, гладильные доски"}'::jsonb,
 '{"sr": "Veš mašine dodaci | Vondi", "en": "Laundry accessories | Vondi", "ru": "Аксессуары для стирки | Vondi"}'::jsonb,
 '{"sr": "Kupite dodatke za veš online", "en": "Buy laundry accessories online", "ru": "Купить аксессуары для стирки онлайн"}'::jsonb,
 '🧺', true),

('kosulјe-kratkih-rukava', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/kosulje-kratkih-rukava', 22,
 '{"sr": "Košulje kratkih rukava", "en": "Short sleeve shirts", "ru": "Рубашки с короткими рукавами"}'::jsonb,
 '{"sr": "Letnje košulje, polo majice, Hawaiian", "en": "Summer shirts, polo shirts, Hawaiian", "ru": "Летние рубашки, поло, гавайские"}'::jsonb,
 '{"sr": "Košulje kratkih rukava | Vondi", "en": "Short sleeve shirts | Vondi", "ru": "Рубашки с короткими рукавами | Vondi"}'::jsonb,
 '{"sr": "Kupite košulje kratkih rukava online", "en": "Buy short sleeve shirts online", "ru": "Купить рубашки с короткими рукавами онлайн"}'::jsonb,
 '👔', true),

('elegantna-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/elegantna-odeca', 23,
 '{"sr": "Elegantna odeća", "en": "Formal wear", "ru": "Элегантная одежда"}'::jsonb,
 '{"sr": "Večernje haljine, smokingzi, odela", "en": "Evening dresses, tuxedos, suits", "ru": "Вечерние платья, смокинги, костюмы"}'::jsonb,
 '{"sr": "Elegantna odeća | Vondi", "en": "Formal wear | Vondi", "ru": "Элегантная одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite elegantnu odeću online", "en": "Buy formal wear online", "ru": "Купить элегантную одежду онлайн"}'::jsonb,
 '🎩', true),

('plus-size-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/plus-size-odeca', 24,
 '{"sr": "Plus size odeća", "en": "Plus size clothing", "ru": "Одежда больших размеров"}'::jsonb,
 '{"sr": "Odeća većih veličina za muškarce i žene", "en": "Larger sizes for men and women", "ru": "Одежда больших размеров для мужчин и женщин"}'::jsonb,
 '{"sr": "Plus size odeća | Vondi", "en": "Plus size clothing | Vondi", "ru": "Одежда больших размеров | Vondi"}'::jsonb,
 '{"sr": "Kupite plus size odeću online", "en": "Buy plus size clothing online", "ru": "Купить одежду больших размеров онлайн"}'::jsonb,
 '👗', true),

('trudnicka-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/trudnicka-odeca', 25,
 '{"sr": "Trudnička odeća", "en": "Maternity wear", "ru": "Одежда для беременных"}'::jsonb,
 '{"sr": "Odeća za trudnice, dojenje, postporođajna", "en": "Maternity clothing, nursing, postpartum", "ru": "Одежда для беременных, кормления, послеродовая"}'::jsonb,
 '{"sr": "Trudnička odeća | Vondi", "en": "Maternity wear | Vondi", "ru": "Одежда для беременных | Vondi"}'::jsonb,
 '{"sr": "Kupite trudničku odeću online", "en": "Buy maternity wear online", "ru": "Купить одежду для беременных онлайн"}'::jsonb,
 '🤰', true),

('naocari-i-dodaci', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/naocari-i-dodaci', 26,
 '{"sr": "Naočari i dodaci", "en": "Glasses & accessories", "ru": "Очки и аксессуары"}'::jsonb,
 '{"sr": "Sunčane naočare, dioptrijske, futrole", "en": "Sunglasses, prescription glasses, cases", "ru": "Солнечные очки, диоптрические, футляры"}'::jsonb,
 '{"sr": "Naočari i dodaci | Vondi", "en": "Glasses & accessories | Vondi", "ru": "Очки и аксессуары | Vondi"}'::jsonb,
 '{"sr": "Kupite naočare i dodatke online", "en": "Buy glasses and accessories online", "ru": "Купить очки и аксессуары онлайн"}'::jsonb,
 '🕶️', true),

('esarpe-i-salovi', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/esarpe-i-salovi', 27,
 '{"sr": "Ešarpe i šalovi", "en": "Scarves & shawls", "ru": "Шарфы и палантины"}'::jsonb,
 '{"sr": "Zimske ešarpe, svileni šalovi, kašmir", "en": "Winter scarves, silk shawls, cashmere", "ru": "Зимние шарфы, шелковые палантины, кашемир"}'::jsonb,
 '{"sr": "Ešarpe i šalovi | Vondi", "en": "Scarves & shawls | Vondi", "ru": "Шарфы и палантины | Vondi"}'::jsonb,
 '{"sr": "Kupite ešarpe i šalove online", "en": "Buy scarves and shawls online", "ru": "Купить шарфы и палантины онлайн"}'::jsonb,
 '🧣', true),

('muskarci-veliki-brojevi', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/muskarci-veliki-brojevi', 28,
 '{"sr": "Muškarci veliki brojevi", "en": "Men big sizes", "ru": "Мужчины большие размеры"}'::jsonb,
 '{"sr": "Muška odeća i obuća velikih brojeva", "en": "Men''s clothing and footwear in large sizes", "ru": "Мужская одежда и обувь больших размеров"}'::jsonb,
 '{"sr": "Muškarci veliki brojevi | Vondi", "en": "Men big sizes | Vondi", "ru": "Мужчины большие размеры | Vondi"}'::jsonb,
 '{"sr": "Kupite odeću za muškarce velikih brojeva online", "en": "Buy big size men''s clothing online", "ru": "Купить одежду для мужчин больших размеров онлайн"}'::jsonb,
 '👔', true),

('zene-veliki-brojevi', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca'), 2, 'odeca-i-obuca/zene-veliki-brojevi', 29,
 '{"sr": "Žene veliki brojevi", "en": "Women big sizes", "ru": "Женщины большие размеры"}'::jsonb,
 '{"sr": "Ženska odeća i obuća velikih brojeva", "en": "Women''s clothing and footwear in large sizes", "ru": "Женская одежда и обувь больших размеров"}'::jsonb,
 '{"sr": "Žene veliki brojevi | Vondi", "en": "Women big sizes | Vondi", "ru": "Женщины большие размеры | Vondi"}'::jsonb,
 '{"sr": "Kupite odeću za žene velikih brojeva online", "en": "Buy big size women''s clothing online", "ru": "Купить одежду для женщин больших размеров онлайн"}'::jsonb,
 '👗', true),

-- =============================================================================
-- Additional L2 for: Elektronika (+ 10 more to reach 25 total)
-- =============================================================================

('gaming-oprema', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/gaming-oprema', 21,
 '{"sr": "Gaming oprema", "en": "Gaming gear", "ru": "Игровое оборудование"}'::jsonb,
 '{"sr": "Gaming tastature, miševi, slušalice, stolice", "en": "Gaming keyboards, mice, headsets, chairs", "ru": "Игровые клавиатуры, мыши, наушники, кресла"}'::jsonb,
 '{"sr": "Gaming oprema | Vondi", "en": "Gaming gear | Vondi", "ru": "Игровое оборудование | Vondi"}'::jsonb,
 '{"sr": "Kupite gaming opremu online", "en": "Buy gaming gear online", "ru": "Купить игровое оборудование онлайн"}'::jsonb,
 '🎮', true),

('smart-home', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/smart-home', 22,
 '{"sr": "Smart home", "en": "Smart home", "ru": "Умный дом"}'::jsonb,
 '{"sr": "Pametni prekidači, sijalice, kamere, senzori", "en": "Smart switches, bulbs, cameras, sensors", "ru": "Умные выключатели, лампочки, камеры, датчики"}'::jsonb,
 '{"sr": "Smart home | Vondi", "en": "Smart home | Vondi", "ru": "Умный дом | Vondi"}'::jsonb,
 '{"sr": "Kupite smart home uređaje online", "en": "Buy smart home devices online", "ru": "Купить устройства умного дома онлайн"}'::jsonb,
 '🏠', true),

('projektori', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/projektori', 23,
 '{"sr": "Projektori", "en": "Projectors", "ru": "Проекторы"}'::jsonb,
 '{"sr": "Projektori za dom, bioskop, prezentacije", "en": "Projectors for home, cinema, presentations", "ru": "Проекторы для дома, кинотеатра, презентаций"}'::jsonb,
 '{"sr": "Projektori | Vondi", "en": "Projectors | Vondi", "ru": "Проекторы | Vondi"}'::jsonb,
 '{"sr": "Kupite projektore online", "en": "Buy projectors online", "ru": "Купить проекторы онлайн"}'::jsonb,
 '📽️', true),

('web-kamere', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/web-kamere', 24,
 '{"sr": "Web kamere", "en": "Webcams", "ru": "Веб-камеры"}'::jsonb,
 '{"sr": "HD web kamere za online sastanke i streaming", "en": "HD webcams for online meetings and streaming", "ru": "HD веб-камеры для онлайн-встреч и стриминга"}'::jsonb,
 '{"sr": "Web kamere | Vondi", "en": "Webcams | Vondi", "ru": "Веб-камеры | Vondi"}'::jsonb,
 '{"sr": "Kupite web kamere online", "en": "Buy webcams online", "ru": "Купить веб-камеры онлайн"}'::jsonb,
 '📹', true),

('skeneri', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/skeneri', 25,
 '{"sr": "Skeneri", "en": "Scanners", "ru": "Сканеры"}'::jsonb,
 '{"sr": "Dokumentni skeneri, foto skeneri, 3D", "en": "Document scanners, photo scanners, 3D", "ru": "Сканеры документов, фотосканеры, 3D"}'::jsonb,
 '{"sr": "Skeneri | Vondi", "en": "Scanners | Vondi", "ru": "Сканеры | Vondi"}'::jsonb,
 '{"sr": "Kupite skenere online", "en": "Buy scanners online", "ru": "Купить сканеры онлайн"}'::jsonb,
 '🖨️', true),

('nas-i-storage', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/nas-i-storage', 26,
 '{"sr": "NAS i storage", "en": "NAS & storage", "ru": "NAS и хранение данных"}'::jsonb,
 '{"sr": "Mrežni diskovi, eksterni HDD, SSD", "en": "Network drives, external HDDs, SSDs", "ru": "Сетевые диски, внешние HDD, SSD"}'::jsonb,
 '{"sr": "NAS i storage | Vondi", "en": "NAS & storage | Vondi", "ru": "NAS и хранение данных | Vondi"}'::jsonb,
 '{"sr": "Kupite NAS i storage uređaje online", "en": "Buy NAS and storage devices online", "ru": "Купить NAS и устройства хранения онлайн"}'::jsonb,
 '💾', true),

('kalkul atori', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/kalkulatori', 27,
 '{"sr": "Kalkulatori", "en": "Calculators", "ru": "Калькуляторы"}'::jsonb,
 '{"sr": "Naučni, grafički, finansijski kalkulatori", "en": "Scientific, graphing, financial calculators", "ru": "Научные, графические, финансовые калькуляторы"}'::jsonb,
 '{"sr": "Kalkulatori | Vondi", "en": "Calculators | Vondi", "ru": "Калькуляторы | Vondi"}'::jsonb,
 '{"sr": "Kupite kalkulatore online", "en": "Buy calculators online", "ru": "Купить калькуляторы онлайн"}'::jsonb,
 '🔢', true),

('mikrofoni', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/mikrofoni', 28,
 '{"sr": "Mikrofoni", "en": "Microphones", "ru": "Микрофоны"}'::jsonb,
 '{"sr": "USB mikrofoni, kondenzatorski, bezžični", "en": "USB microphones, condenser, wireless", "ru": "USB микрофоны, конденсаторные, беспроводные"}'::jsonb,
 '{"sr": "Mikrofoni | Vondi", "en": "Microphones | Vondi", "ru": "Микрофоны | Vondi"}'::jsonb,
 '{"sr": "Kupite mikrofone online", "en": "Buy microphones online", "ru": "Купить микрофоны онлайн"}'::jsonb,
 '🎙️', true),

('smart-narukvice', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/smart-narukvice', 29,
 '{"sr": "Smart narukvice", "en": "Smart bands", "ru": "Умные браслеты"}'::jsonb,
 '{"sr": "Fitness narukvice, trackers aktivnosti", "en": "Fitness bands, activity trackers", "ru": "Фитнес-браслеты, трекеры активности"}'::jsonb,
 '{"sr": "Smart narukvice | Vondi", "en": "Smart bands | Vondi", "ru": "Умные браслеты | Vondi"}'::jsonb,
 '{"sr": "Kupite smart narukvice online", "en": "Buy smart bands online", "ru": "Купить умные браслеты онлайн"}'::jsonb,
 '⌚', true),

('konzolne-igre', (SELECT id FROM categories WHERE slug = 'elektronika'), 2, 'elektronika/konzolne-igre', 30,
 '{"sr": "Konzolne igre", "en": "Console games", "ru": "Игры для консолей"}'::jsonb,
 '{"sr": "PS5, Xbox, Nintendo Switch igre", "en": "PS5, Xbox, Nintendo Switch games", "ru": "Игры для PS5, Xbox, Nintendo Switch"}'::jsonb,
 '{"sr": "Konzolne igre | Vondi", "en": "Console games | Vondi", "ru": "Игры для консолей | Vondi"}'::jsonb,
 '{"sr": "Kupite konzolne igre online", "en": "Buy console games online", "ru": "Купить игры для консолей онлайн"}'::jsonb,
 '🎮', true),

-- =============================================================================
-- Additional L2 for: Dom i bašta (+ 10 more to reach 25 total)
-- =============================================================================

('namestaj-kancelarija', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/namestaj-kancelarija', 16,
 '{"sr": "Nameštaj kancelarija", "en": "Office furniture", "ru": "Офисная мебель"}'::jsonb,
 '{"sr": "Radni stolovi, kancelarijske stolice, police", "en": "Desks, office chairs, shelves", "ru": "Письменные столы, офисные кресла, полки"}'::jsonb,
 '{"sr": "Nameštaj kancelarija | Vondi", "en": "Office furniture | Vondi", "ru": "Офисная мебель | Vondi"}'::jsonb,
 '{"sr": "Kupite nameštaj za kancelariju online", "en": "Buy office furniture online", "ru": "Купить офисную мебель онлайн"}'::jsonb,
 '🪑', true),

('kuhinjski-pribor', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/kuhinjski-pribor', 17,
 '{"sr": "Kuhinjski pribor", "en": "Kitchenware", "ru": "Кухонная утварь"}'::jsonb,
 '{"sr": "Noževi, šerpe, tiganj, posuđe", "en": "Knives, pots, pans, dishes", "ru": "Ножи, кастрюли, сковороды, посуда"}'::jsonb,
 '{"sr": "Kuhinjski pribor | Vondi", "en": "Kitchenware | Vondi", "ru": "Кухонная утварь | Vondi"}'::jsonb,
 '{"sr": "Kupite kuhinjski pribor online", "en": "Buy kitchenware online", "ru": "Купить кухонную утварь онлайн"}'::jsonb,
 '🔪', true),

('tepihi-i-prostirke', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/tepihi-i-prostirke', 18,
 '{"sr": "Tepisi i prostirke", "en": "Carpets & rugs", "ru": "Ковры и коврики"}'::jsonb,
 '{"sr": "Tepisi, prostirke, staze, protivklizni", "en": "Carpets, rugs, runners, anti-slip", "ru": "Ковры, коврики, дорожки, противоскользящие"}'::jsonb,
 '{"sr": "Tepisi i prostirke | Vondi", "en": "Carpets & rugs | Vondi", "ru": "Ковры и коврики | Vondi"}'::jsonb,
 '{"sr": "Kupite tepise i prostirke online", "en": "Buy carpets and rugs online", "ru": "Купить ковры и коврики онлайн"}'::jsonb,
 '🪆', true),

('ogledala', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/ogledala', 19,
 '{"sr": "Ogledala", "en": "Mirrors", "ru": "Зеркала"}'::jsonb,
 '{"sr": "Zidna ogledala, stojeća, sa osvetljenjem", "en": "Wall mirrors, standing, with lighting", "ru": "Настенные зеркала, напольные, с подсветкой"}'::jsonb,
 '{"sr": "Ogledala | Vondi", "en": "Mirrors | Vondi", "ru": "Зеркала | Vondi"}'::jsonb,
 '{"sr": "Kupite ogledala online", "en": "Buy mirrors online", "ru": "Купить зеркала онлайн"}'::jsonb,
 '🪞', true),

('sat i-za-zid', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/satovi-za-zid', 20,
 '{"sr": "Satovi za zid", "en": "Wall clocks", "ru": "Настенные часы"}'::jsonb,
 '{"sr": "Zidni satovi, alarmni, klatna", "en": "Wall clocks, alarm clocks, pendulum", "ru": "Настенные часы, будильники, маятниковые"}'::jsonb,
 '{"sr": "Satovi za zid | Vondi", "en": "Wall clocks | Vondi", "ru": "Настенные часы | Vondi"}'::jsonb,
 '{"sr": "Kupite satove za zid online", "en": "Buy wall clocks online", "ru": "Купить настенные часы онлайн"}'::jsonb,
 '🕰️', true),

('pregradni-zidovi', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/pregradni-zidovi', 21,
 '{"sr": "Pregradni zidovi", "en": "Room dividers", "ru": "Перегородки"}'::jsonb,
 '{"sr": "Paravani, police, klizni paneli", "en": "Screens, shelves, sliding panels", "ru": "Ширмы, полки, раздвижные панели"}'::jsonb,
 '{"sr": "Pregradni zidovi | Vondi", "en": "Room dividers | Vondi", "ru": "Перегородки | Vondi"}'::jsonb,
 '{"sr": "Kupite pregradne zidove online", "en": "Buy room dividers online", "ru": "Купить перегородки онлайн"}'::jsonb,
 '🚪', true),

('vaze-i-dekor', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/vaze-i-dekor', 22,
 '{"sr": "Vaze i dekor", "en": "Vases & decor", "ru": "Вазы и декор"}'::jsonb,
 '{"sr": "Staklene vaze, keramičke, sveće, ukrasne figurice", "en": "Glass vases, ceramic, candles, decorative figures", "ru": "Стеклянные вазы, керамические, свечи, декоративные фигурки"}'::jsonb,
 '{"sr": "Vaze i dekor | Vondi", "en": "Vases & decor | Vondi", "ru": "Вазы и декор | Vondi"}'::jsonb,
 '{"sr": "Kupite vaze i dekoracije online", "en": "Buy vases and decorations online", "ru": "Купить вазы и украшения онлайн"}'::jsonb,
 '🏺', true),

('bastenska-rasveta', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/bastenska-rasveta', 23,
 '{"sr": "Baštenka rasveta", "en": "Garden lighting", "ru": "Садовое освещение"}'::jsonb,
 '{"sr": "Solarne lampe, LED rasveta, reflektori", "en": "Solar lamps, LED lighting, floodlights", "ru": "Солнечные лампы, LED освещение, прожекторы"}'::jsonb,
 '{"sr": "Baštenka rasveta | Vondi", "en": "Garden lighting | Vondi", "ru": "Садовое освещение | Vondi"}'::jsonb,
 '{"sr": "Kupite baštensku rasvetu online", "en": "Buy garden lighting online", "ru": "Купить садовое освещение онлайн"}'::jsonb,
 '💡', true),

('bastenske-ukrase', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/bastenske-ukrase', 24,
 '{"sr": "Baštenke ukrase", "en": "Garden decorations", "ru": "Садовые украшения"}'::jsonb,
 '{"sr": "Patuljci, figure životinja, fontane, vetruške", "en": "Gnomes, animal figures, fountains, wind spinners", "ru": "Гномы, фигурки животных, фонтаны, ветряные спиннеры"}'::jsonb,
 '{"sr": "Baštenke ukrase | Vondi", "en": "Garden decorations | Vondi", "ru": "Садовые украшения | Vondi"}'::jsonb,
 '{"sr": "Kupite baštenke ukrase online", "en": "Buy garden decorations online", "ru": "Купить садовые украшения онлайн"}'::jsonb,
 '🌻', true),

('kompostiranje', (SELECT id FROM categories WHERE slug = 'dom-i-basta'), 2, 'dom-i-basta/kompostiranje', 25,
 '{"sr": "Kompostiranje", "en": "Composting", "ru": "Компостирование"}'::jsonb,
 '{"sr": "Komposteri, biorazgradive kese, alati", "en": "Composters, biodegradable bags, tools", "ru": "Компостеры, биоразлагаемые пакеты, инструменты"}'::jsonb,
 '{"sr": "Kompostiranje | Vondi", "en": "Composting | Vondi", "ru": "Компостирование | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za kompostiranje online", "en": "Buy composting equipment online", "ru": "Купить оборудование для компостирования онлайн"}'::jsonb,
 '♻️', true),

-- =============================================================================
-- Additional L2 for: Lepota i zdravlje (+ 10 more to reach 22 total)
-- =============================================================================

('makeup-cetkice', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/makeup-cetkice', 13,
 '{"sr": "Makeup četkice", "en": "Makeup brushes", "ru": "Кисти для макияжа"}'::jsonb,
 '{"sr": "Četkice za puder, senke, korektor, kompleti", "en": "Brushes for powder, eyeshadow, concealer, sets", "ru": "Кисти для пудры, теней, корректора, наборы"}'::jsonb,
 '{"sr": "Makeup četkice | Vondi", "en": "Makeup brushes | Vondi", "ru": "Кисти для макияжа | Vondi"}'::jsonb,
 '{"sr": "Kupite makeup četkice online", "en": "Buy makeup brushes online", "ru": "Купить кисти для макияжа онлайн"}'::jsonb,
 '🖌️', true),

('lepota-aparati', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/lepota-aparati', 14,
 '{"sr": "Lepota aparati", "en": "Beauty devices", "ru": "Приборы для красоты"}'::jsonb,
 '{"sr": "Fen za kosu, presa, epilatori, manikir set", "en": "Hair dryers, straighteners, epilators, manicure sets", "ru": "Фены, выпрямители, эпиляторы, маникюрные наборы"}'::jsonb,
 '{"sr": "Lepota aparati | Vondi", "en": "Beauty devices | Vondi", "ru": "Приборы для красоты | Vondi"}'::jsonb,
 '{"sr": "Kupite aparate za lepotu online", "en": "Buy beauty devices online", "ru": "Купить приборы для красоты онлайн"}'::jsonb,
 '💆', true),

('muski-stil', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/muski-stil', 15,
 '{"sr": "Muški stil", "en": "Men''s style", "ru": "Мужской стиль"}'::jsonb,
 '{"sr": "Aparati za brijanje, trim eri, balzami nakon brijanja", "en": "Shavers, trimmers, aftershave balms", "ru": "Бритвы, триммеры, бальзамы после бритья"}'::jsonb,
 '{"sr": "Muški stil | Vondi", "en": "Men''s style | Vondi", "ru": "Мужской стиль | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za muški stil online", "en": "Buy men''s style products online", "ru": "Купить продукты для мужского стиля онлайн"}'::jsonb,
 '🧔', true),

('anti-aging', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/anti-aging', 16,
 '{"sr": "Anti-aging", "en": "Anti-aging", "ru": "Антивозрастной уход"}'::jsonb,
 '{"sr": "Kreme protiv bora, serumi, tretmani", "en": "Anti-wrinkle creams, serums, treatments", "ru": "Кремы против морщин, сыворотки, процедуры"}'::jsonb,
 '{"sr": "Anti-aging | Vondi", "en": "Anti-aging | Vondi", "ru": "Антивозрастной уход | Vondi"}'::jsonb,
 '{"sr": "Kupite anti-aging proizvode online", "en": "Buy anti-aging products online", "ru": "Купить антивозрастные продукты онлайн"}'::jsonb,
 '✨', true),

('organska-kozmetika', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/organska-kozmetika', 17,
 '{"sr": "Organska kozmetika", "en": "Organic cosmetics", "ru": "Органическая косметика"}'::jsonb,
 '{"sr": "Prirodna kozmetika, veganska, bez parabena", "en": "Natural cosmetics, vegan, paraben-free", "ru": "Натуральная косметика, веганская, без парабенов"}'::jsonb,
 '{"sr": "Organska kozmetika | Vondi", "en": "Organic cosmetics | Vondi", "ru": "Органическая косметика | Vondi"}'::jsonb,
 '{"sr": "Kupite organsku kozmetiku online", "en": "Buy organic cosmetics online", "ru": "Купить органическую косметику онлайн"}'::jsonb,
 '🌿', true),

('luksuzna-kozmetika', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/luksuzna-kozmetika', 18,
 '{"sr": "Luksuzna kozmetika", "en": "Luxury cosmetics", "ru": "Люксовая косметика"}'::jsonb,
 '{"sr": "Premium brendovi, luksuzna šminka, nege", "en": "Premium brands, luxury makeup, care", "ru": "Премиум бренды, роскошный макияж, уход"}'::jsonb,
 '{"sr": "Luksuzna kozmetika | Vondi", "en": "Luxury cosmetics | Vondi", "ru": "Люксовая косметика | Vondi"}'::jsonb,
 '{"sr": "Kupite luksuznu kozmetiku online", "en": "Buy luxury cosmetics online", "ru": "Купить люксовую косметику онлайн"}'::jsonb,
 '👑', true),

('zastita-od-sunca', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/zastita-od-sunca', 19,
 '{"sr": "Zaštita od sunca", "en": "Sun protection", "ru": "Защита от солнца"}'::jsonb,
 '{"sr": "Kreme za sunčanje, after sun, SPF", "en": "Sunscreen, after sun, SPF products", "ru": "Солнцезащитные кремы, после загара, SPF"}'::jsonb,
 '{"sr": "Zaštita od sunca | Vondi", "en": "Sun protection | Vondi", "ru": "Защита от солнца | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za zaštitu od sunca online", "en": "Buy sun protection products online", "ru": "Купить средства защиты от солнца онлайн"}'::jsonb,
 '☀️', true),

('depilacija', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/depilacija', 20,
 '{"sr": "Depilacija", "en": "Hair removal", "ru": "Депиляция"}'::jsonb,
 '{"sr": "Epilatori, vosak, kreme za depilaciju, laseri", "en": "Epilators, wax, depilatory creams, lasers", "ru": "Эпиляторы, воск, кремы для депиляции, лазеры"}'::jsonb,
 '{"sr": "Depilacija | Vondi", "en": "Hair removal | Vondi", "ru": "Депиляция | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za depilaciju online", "en": "Buy hair removal products online", "ru": "Купить средства для депиляции онлайн"}'::jsonb,
 '🪒', true),

('intimna-higijena', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/intimna-higijena', 21,
 '{"sr": "Intimna higijena", "en": "Intimate hygiene", "ru": "Интимная гигиена"}'::jsonb,
 '{"sr": "Intimni gelovi, vlažne maramice, brendovi", "en": "Intimate gels, wet wipes, brands", "ru": "Интимные гели, влажные салфетки, бренды"}'::jsonb,
 '{"sr": "Intimna higijena | Vondi", "en": "Intimate hygiene | Vondi", "ru": "Интимная гигиена | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za intimnu higijenu online", "en": "Buy intimate hygiene products online", "ru": "Купить средства интимной гигиены онлайн"}'::jsonb,
 '🧴', true),

('gelovi-za-tusiranje', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/gelovi-za-tusiranje', 22,
 '{"sr": "Gelovi za tuširanje", "en": "Shower gels", "ru": "Гели для душа"}'::jsonb,
 '{"sr": "Gelovi, pene, mirisni, hidratantni", "en": "Gels, foams, fragrant, moisturizing", "ru": "Гели, пены, ароматные, увлажняющие"}'::jsonb,
 '{"sr": "Gelovi za tuširanje | Vondi", "en": "Shower gels | Vondi", "ru": "Гели для душа | Vondi"}'::jsonb,
 '{"sr": "Kupite gelove za tuširanje online", "en": "Buy shower gels online", "ru": "Купить гели для душа онлайн"}'::jsonb,
 '🧼', true);

-- Progress: Part 4 adds 50 new L2 categories
-- Total L2 so far: 194 (previous) + 50 (this part) = 244

DO $$
DECLARE
    l2_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO l2_count FROM categories WHERE level = 2;
    RAISE NOTICE 'Part 4 complete: % total L2 categories', l2_count;

    IF l2_count < 244 THEN
        RAISE WARNING 'Expected at least 244 L2 categories, found %', l2_count;
    END IF;
END $$;
