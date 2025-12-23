-- Migration: Expand L2 categories (Part 5)
-- Date: 2025-12-17
-- Purpose: Add 80 L2 categories - Popular categories expansion
-- Expanding: Odeća, Elektronika, Dom i bašta, Lepota, Bebe, Sport, Auto, Aparati

-- =============================================================================
-- Additional L2 for: Odeća i obuća (+ 10 more)
-- =============================================================================

INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('spec-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/spec-odeca', 100,
 '{"sr": "Specijalna odeća", "en": "Specialized clothing", "ru": "Специальная одежда"}'::jsonb,
 '{"sr": "Radna odeća, uniforma, zaštitna oprema", "en": "Work clothing, uniforms, protective gear", "ru": "Рабочая одежда, униформа, защитное снаряжение"}'::jsonb,
 '{"sr": "Specijalna odeća | Vondi", "en": "Specialized clothing | Vondi", "ru": "Специальная одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite specijalnu odeću online - radna, uniforma, zaštitna", "en": "Buy specialized clothing online", "ru": "Купить специальную одежду онлайн"}'::jsonb,
 '🦺', true),

('uniforma', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/uniforma', 101,
 '{"sr": "Uniforme", "en": "Uniforms", "ru": "Униформа"}'::jsonb,
 '{"sr": "Medicinska, kuhinjska, hotelska uniforma", "en": "Medical, kitchen, hotel uniforms", "ru": "Медицинская, кухонная, гостиничная униформа"}'::jsonb,
 '{"sr": "Uniforme | Vondi", "en": "Uniforms | Vondi", "ru": "Униформа | Vondi"}'::jsonb,
 '{"sr": "Kupite uniforme za sve profesije online", "en": "Buy professional uniforms online", "ru": "Купить профессиональную униформу онлайн"}'::jsonb,
 '👔', true),

('tradicionalna-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/tradicionalna-odeca', 102,
 '{"sr": "Tradicionalna odeća", "en": "Traditional clothing", "ru": "Традиционная одежда"}'::jsonb,
 '{"sr": "Narodna nošnja, folklorni kostimi", "en": "National costumes, folk outfits", "ru": "Национальные костюмы, фольклорная одежда"}'::jsonb,
 '{"sr": "Tradicionalna odeća | Vondi", "en": "Traditional clothing | Vondi", "ru": "Традиционная одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite tradicionalnu odeću i narodnu nošnju", "en": "Buy traditional and folk clothing online", "ru": "Купить традиционную и народную одежду онлайн"}'::jsonb,
 '👘', true),

('premium-aksesoari', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/premium-aksesoari', 103,
 '{"sr": "Premium aksesoari", "en": "Premium accessories", "ru": "Премиум аксессуары"}'::jsonb,
 '{"sr": "Dizajnerski dodaci, luksuzne torbe, premium aksesoari", "en": "Designer accessories, luxury bags, premium items", "ru": "Дизайнерские аксессуары, люксовые сумки, премиум изделия"}'::jsonb,
 '{"sr": "Premium aksesoari | Vondi", "en": "Premium accessories | Vondi", "ru": "Премиум аксессуары | Vondi"}'::jsonb,
 '{"sr": "Kupite premium i dizajnerske aksesoare online", "en": "Buy premium and designer accessories online", "ru": "Купить премиум и дизайнерские аксессуары онлайн"}'::jsonb,
 '👜', true),

('vintage-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/vintage-odeca', 104,
 '{"sr": "Vintage odeća", "en": "Vintage clothing", "ru": "Винтажная одежда"}'::jsonb,
 '{"sr": "Retro komadi, second-hand vintage, vintage butici", "en": "Retro pieces, second-hand vintage, vintage boutiques", "ru": "Ретро вещи, винтажный секонд-хенд, винтажные бутики"}'::jsonb,
 '{"sr": "Vintage odeća | Vondi", "en": "Vintage clothing | Vondi", "ru": "Винтажная одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite vintage odeću i retro komade online", "en": "Buy vintage clothing and retro pieces online", "ru": "Купить винтажную одежду и ретро вещи онлайн"}'::jsonb,
 '🕰️', true),

('sportske-uniforme', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/sportske-uniforme', 105,
 '{"sr": "Sportske uniforme", "en": "Sports uniforms", "ru": "Спортивная форма"}'::jsonb,
 '{"sr": "Dresovi, timska oprema, sportske majice sa brojem", "en": "Team kits, sports jerseys with numbers", "ru": "Командная форма, спортивные майки с номерами"}'::jsonb,
 '{"sr": "Sportske uniforme | Vondi", "en": "Sports uniforms | Vondi", "ru": "Спортивная форма | Vondi"}'::jsonb,
 '{"sr": "Kupite sportske uniforme za timove online", "en": "Buy sports team uniforms online", "ru": "Купить командную спортивную форму онлайн"}'::jsonb,
 '⚽', true),

('kucni-mantili', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/kucni-mantili', 106,
 '{"sr": "Kućni mantili", "en": "Home robes", "ru": "Домашние халаты"}'::jsonb,
 '{"sr": "Bade mantili, pidžame, kućna garderoba", "en": "Bath robes, pajamas, home wear", "ru": "Банные халаты, пижамы, домашняя одежда"}'::jsonb,
 '{"sr": "Kućni mantili | Vondi", "en": "Home robes | Vondi", "ru": "Домашние халаты | Vondi"}'::jsonb,
 '{"sr": "Kupite kućne mantile i pidžame online", "en": "Buy home robes and pajamas online", "ru": "Купить домашние халаты и пижамы онлайн"}'::jsonb,
 '🛀', true),

('pletena-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/pletena-odeca', 107,
 '{"sr": "Pletena odeća", "en": "Knitwear", "ru": "Трикотажная одежда"}'::jsonb,
 '{"sr": "Ručno pleteni džemperi, šalovi, pletene haljine", "en": "Hand-knit sweaters, scarves, knit dresses", "ru": "Ручная вязка, шарфы, вязаные платья"}'::jsonb,
 '{"sr": "Pletena odeća | Vondi", "en": "Knitwear | Vondi", "ru": "Трикотажная одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite pletenu odeću i ručne radove online", "en": "Buy knitwear and handmade items online", "ru": "Купить вязаную одежду и ручную работу онлайн"}'::jsonb,
 '🧶', true),

('funkcionalna-odeca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/funkcionalna-odeca', 108,
 '{"sr": "Funkcionalna odeća", "en": "Functional clothing", "ru": "Функциональная одежда"}'::jsonb,
 '{"sr": "Outdoor, hiking, planinarska odeća", "en": "Outdoor, hiking, mountain clothing", "ru": "Одежда для активного отдыха, походов, гор"}'::jsonb,
 '{"sr": "Funkcionalna odeća | Vondi", "en": "Functional clothing | Vondi", "ru": "Функциональная одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite funkcionalnu odeću za outdoor online", "en": "Buy functional outdoor clothing online", "ru": "Купить функциональную одежду для активного отдыха онлайн"}'::jsonb,
 '🏔️', true),

('modna-obuca', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' AND level = 1), 2, 'odeca-i-obuca/modna-obuca', 109,
 '{"sr": "Modna obuća", "en": "Fashion footwear", "ru": "Модная обувь"}'::jsonb,
 '{"sr": "Dizajnerska obuća, trendy patike, modne cipele", "en": "Designer footwear, trendy sneakers, fashion shoes", "ru": "Дизайнерская обувь, трендовые кроссовки, модная обувь"}'::jsonb,
 '{"sr": "Modna obuća | Vondi", "en": "Fashion footwear | Vondi", "ru": "Модная обувь | Vondi"}'::jsonb,
 '{"sr": "Kupite modnu i dizajnersku obuću online", "en": "Buy fashion and designer footwear online", "ru": "Купить модную и дизайнерскую обувь онлайн"}'::jsonb,
 '👠', true),

-- =============================================================================
-- Additional L2 for: Elektronika (+ 10 more)
-- =============================================================================

('dronovi', (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1), 2, 'elektronika/dronovi', 100,
 '{"sr": "Dronovi", "en": "Drones", "ru": "Дроны"}'::jsonb,
 '{"sr": "DJI dronovi, kamere za dronove, FPV dronovi", "en": "DJI drones, drone cameras, FPV drones", "ru": "Дроны DJI, камеры для дронов, FPV дроны"}'::jsonb,
 '{"sr": "Dronovi | Vondi", "en": "Drones | Vondi", "ru": "Дроны | Vondi"}'::jsonb,
 '{"sr": "Kupite dronove sa kamerom online", "en": "Buy drones with cameras online", "ru": "Купить дроны с камерой онлайн"}'::jsonb,
 '🚁', true),

('vr-ar-oprema', (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1), 2, 'elektronika/vr-ar-oprema', 101,
 '{"sr": "VR i AR oprema", "en": "VR & AR equipment", "ru": "VR и AR оборудование"}'::jsonb,
 '{"sr": "Virtuelna realnost, proširena realnost, VR naočare", "en": "Virtual reality, augmented reality, VR headsets", "ru": "Виртуальная реальность, дополненная реальность, VR шлемы"}'::jsonb,
 '{"sr": "VR i AR oprema | Vondi", "en": "VR & AR equipment | Vondi", "ru": "VR и AR оборудование | Vondi"}'::jsonb,
 '{"sr": "Kupite VR i AR opremu online", "en": "Buy VR and AR equipment online", "ru": "Купить VR и AR оборудование онлайн"}'::jsonb,
 '🥽', true),

('roboti', (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1), 2, 'elektronika/roboti', 102,
 '{"sr": "Roboti", "en": "Robots", "ru": "Роботы"}'::jsonb,
 '{"sr": "Roboti usisivači, edukativni roboti, robot igrački", "en": "Robot vacuums, educational robots, toy robots", "ru": "Роботы-пылесосы, образовательные роботы, роботы-игрушки"}'::jsonb,
 '{"sr": "Roboti | Vondi", "en": "Robots | Vondi", "ru": "Роботы | Vondi"}'::jsonb,
 '{"sr": "Kupite robote i robot opremu online", "en": "Buy robots and robot equipment online", "ru": "Купить роботов и роботизированное оборудование онлайн"}'::jsonb,
 '🤖', true),

-- REMOVED DUPLICATE: smart-home (already exists in part4)

('3d-printeri', (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1), 2, 'elektronika/3d-printeri', 104,
 '{"sr": "3D printeri", "en": "3D printers", "ru": "3D принтеры"}'::jsonb,
 '{"sr": "3D štampači, filamenti, resin printeri", "en": "3D printers, filaments, resin printers", "ru": "3D принтеры, филаменты, смоляные принтеры"}'::jsonb,
 '{"sr": "3D printeri | Vondi", "en": "3D printers | Vondi", "ru": "3D принтеры | Vondi"}'::jsonb,
 '{"sr": "Kupite 3D printere i materijale online", "en": "Buy 3D printers and materials online", "ru": "Купить 3D принтеры и материалы онлайн"}'::jsonb,
 '🖨️', true),

('gaming-stolice', (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1), 2, 'elektronika/gaming-stolice', 105,
 '{"sr": "Gaming stolice", "en": "Gaming chairs", "ru": "Игровые кресла"}'::jsonb,
 '{"sr": "Gamer fotelje, racing stolice, ergonomske stolice", "en": "Gamer chairs, racing chairs, ergonomic chairs", "ru": "Геймерские кресла, гоночные кресла, эргономичные кресла"}'::jsonb,
 '{"sr": "Gaming stolice | Vondi", "en": "Gaming chairs | Vondi", "ru": "Игровые кресла | Vondi"}'::jsonb,
 '{"sr": "Kupite gaming stolice online", "en": "Buy gaming chairs online", "ru": "Купить игровые кресла онлайн"}'::jsonb,
 '🪑', true),

('streaming-oprema', (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1), 2, 'elektronika/streaming-oprema', 106,
 '{"sr": "Streaming oprema", "en": "Streaming equipment", "ru": "Оборудование для стриминга"}'::jsonb,
 '{"sr": "Mikrofoni, kamere, capture kartice, osvetljenje", "en": "Microphones, cameras, capture cards, lighting", "ru": "Микрофоны, камеры, карты захвата, освещение"}'::jsonb,
 '{"sr": "Streaming oprema | Vondi", "en": "Streaming equipment | Vondi", "ru": "Оборудование для стриминга | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za streaming online", "en": "Buy streaming equipment online", "ru": "Купить оборудование для стриминга онлайн"}'::jsonb,
 '🎙️', true),

('power-stanice', (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1), 2, 'elektronika/power-stanice', 107,
 '{"sr": "Power stanice", "en": "Power stations", "ru": "Портативные электростанции"}'::jsonb,
 '{"sr": "Prenosive baterije, generatori, solarni paneli", "en": "Portable batteries, generators, solar panels", "ru": "Портативные батареи, генераторы, солнечные панели"}'::jsonb,
 '{"sr": "Power stanice | Vondi", "en": "Power stations | Vondi", "ru": "Портативные электростанции | Vondi"}'::jsonb,
 '{"sr": "Kupite power stanice i solarne panele online", "en": "Buy power stations and solar panels online", "ru": "Купить портативные электростанции и солнечные панели онлайн"}'::jsonb,
 '🔋', true),

('projektori-prenosni', (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1), 2, 'elektronika/projektori-prenosni', 108,
 '{"sr": "Prenosni projektori", "en": "Portable projectors", "ru": "Портативные проекторы"}'::jsonb,
 '{"sr": "Mini projektori, pocket projektori, LED projektori", "en": "Mini projectors, pocket projectors, LED projectors", "ru": "Мини проекторы, карманные проекторы, LED проекторы"}'::jsonb,
 '{"sr": "Prenosni projektori | Vondi", "en": "Portable projectors | Vondi", "ru": "Портативные проекторы | Vondi"}'::jsonb,
 '{"sr": "Kupite prenosne projektore online", "en": "Buy portable projectors online", "ru": "Купить портативные проекторы онлайн"}'::jsonb,
 '📽️', true),

('elektro-transport', (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1), 2, 'elektronika/elektro-transport', 109,
 '{"sr": "Elektro transport", "en": "Electric transport", "ru": "Электротранспорт"}'::jsonb,
 '{"sr": "Električni trotineti, segway, hoverboard, elektro bicikli", "en": "Electric scooters, segway, hoverboard, e-bikes", "ru": "Электросамокаты, сегвеи, гироскутеры, электровелосипеды"}'::jsonb,
 '{"sr": "Elektro transport | Vondi", "en": "Electric transport | Vondi", "ru": "Электротранспорт | Vondi"}'::jsonb,
 '{"sr": "Kupite elektro transport online", "en": "Buy electric transport online", "ru": "Купить электротранспорт онлайн"}'::jsonb,
 '🛴', true),

-- =============================================================================
-- Additional L2 for: Dom i bašta (+ 10 more)
-- =============================================================================

('premium-tekstil', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/premium-tekstil', 100,
 '{"sr": "Premium tekstil", "en": "Premium textiles", "ru": "Премиум текстиль"}'::jsonb,
 '{"sr": "Svilena posteljina, pamučni peškiri, luksuzne tkanine", "en": "Silk bedding, cotton towels, luxury fabrics", "ru": "Шелковое белье, хлопковые полотенца, люксовые ткани"}'::jsonb,
 '{"sr": "Premium tekstil | Vondi", "en": "Premium textiles | Vondi", "ru": "Премиум текстиль | Vondi"}'::jsonb,
 '{"sr": "Kupite premium tekstil za dom online", "en": "Buy premium home textiles online", "ru": "Купить премиум текстиль для дома онлайн"}'::jsonb,
 '🧵', true),

('premium-posudje', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/premium-posudje', 101,
 '{"sr": "Premium posuđe", "en": "Premium tableware", "ru": "Премиум посуда"}'::jsonb,
 '{"sr": "Porcelan, kristal, dizajnersko posuđe", "en": "Porcelain, crystal, designer tableware", "ru": "Фарфор, хрусталь, дизайнерская посуда"}'::jsonb,
 '{"sr": "Premium posuđe | Vondi", "en": "Premium tableware | Vondi", "ru": "Премиум посуда | Vondi"}'::jsonb,
 '{"sr": "Kupite premium posuđe online", "en": "Buy premium tableware online", "ru": "Купить премиум посуду онлайн"}'::jsonb,
 '🍽️', true),

('dekor-zidova', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/dekor-zidova', 102,
 '{"sr": "Dekor zidova", "en": "Wall decor", "ru": "Декор стен"}'::jsonb,
 '{"sr": "Slike, tapete, zidne nalepnice, ogledala", "en": "Paintings, wallpapers, wall stickers, mirrors", "ru": "Картины, обои, настенные наклейки, зеркала"}'::jsonb,
 '{"sr": "Dekor zidova | Vondi", "en": "Wall decor | Vondi", "ru": "Декор стен | Vondi"}'::jsonb,
 '{"sr": "Kupite dekor za zidove online", "en": "Buy wall decor online", "ru": "Купить декор для стен онлайн"}'::jsonb,
 '🖼️', true),

('tepisi-i-circi', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/tepisi-i-circi', 103,
 '{"sr": "Tepisi i ćilimi", "en": "Carpets & rugs", "ru": "Ковры и коврики"}'::jsonb,
 '{"sr": "Persijski tepisi, moderna, ručno tkani ćilimi", "en": "Persian carpets, modern rugs, hand-woven kilims", "ru": "Персидские ковры, современные, ручные килимы"}'::jsonb,
 '{"sr": "Tepisi i ćilimi | Vondi", "en": "Carpets & rugs | Vondi", "ru": "Ковры и коврики | Vondi"}'::jsonb,
 '{"sr": "Kupite tepise i ćilime online", "en": "Buy carpets and rugs online", "ru": "Купить ковры и коврики онлайн"}'::jsonb,
 '🧶', true),

('zavese-i-zastori', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/zavese-i-zastori', 104,
 '{"sr": "Zavese i zastori", "en": "Curtains & blinds", "ru": "Шторы и жалюзи"}'::jsonb,
 '{"sr": "Zavese, rolletne, rimske zavese, tamne zavese", "en": "Curtains, roller blinds, roman shades, blackout curtains", "ru": "Шторы, рулонные, римские, блэкаут"}'::jsonb,
 '{"sr": "Zavese i zastori | Vondi", "en": "Curtains & blinds | Vondi", "ru": "Шторы и жалюзи | Vondi"}'::jsonb,
 '{"sr": "Kupite zavese i zastorne online", "en": "Buy curtains and blinds online", "ru": "Купить шторы и жалюзи онлайн"}'::jsonb,
 '🪟', true),

('pametna-rasveta', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/pametna-rasveta', 105,
 '{"sr": "Pametna rasveta", "en": "Smart lighting", "ru": "Умное освещение"}'::jsonb,
 '{"sr": "LED pametne sijalice, RGB osvetljenje, smart lampe", "en": "LED smart bulbs, RGB lighting, smart lamps", "ru": "Умные LED лампы, RGB освещение, умные светильники"}'::jsonb,
 '{"sr": "Pametna rasveta | Vondi", "en": "Smart lighting | Vondi", "ru": "Умное освещение | Vondi"}'::jsonb,
 '{"sr": "Kupite pametnu rasvetu online", "en": "Buy smart lighting online", "ru": "Купить умное освещение онлайн"}'::jsonb,
 '💡', true),

('baštenska-premium-oprema', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/bastenska-premium-oprema', 106,
 '{"sr": "Baštenska premium oprema", "en": "Garden premium equipment", "ru": "Премиум садовая техника"}'::jsonb,
 '{"sr": "Luksuzni ležaljke, terasna premium oprema", "en": "Luxury loungers, premium terrace equipment", "ru": "Люксовые шезлонги, премиум терасное оборудование"}'::jsonb,
 '{"sr": "Baštenska premium oprema | Vondi", "en": "Garden premium equipment | Vondi", "ru": "Премиум садовая техника | Vondi"}'::jsonb,
 '{"sr": "Kupite premium opremu za baštu online", "en": "Buy premium garden equipment online", "ru": "Купить премиум садовую технику онлайн"}'::jsonb,
 '🌿', true),

('gril-i-bbq', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/gril-i-bbq', 107,
 '{"sr": "Gril i BBQ", "en": "Grill & BBQ", "ru": "Гриль и барбекю"}'::jsonb,
 '{"sr": "Roštilj, kamado gril, gas gril, pribor za roštilj", "en": "Barbecue, kamado grill, gas grill, BBQ accessories", "ru": "Барбекю, камадо гриль, газовый гриль, аксессуары для барбекю"}'::jsonb,
 '{"sr": "Gril i BBQ | Vondi", "en": "Grill & BBQ | Vondi", "ru": "Гриль и барбекю | Vondi"}'::jsonb,
 '{"sr": "Kupite gril i BBQ opremu online", "en": "Buy grill and BBQ equipment online", "ru": "Купить гриль и барбекю оборудование онлайн"}'::jsonb,
 '🔥', true),

('bazeni-i-dzakuzi', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/bazeni-i-dzakuzi', 108,
 '{"sr": "Bazeni i džakuzi", "en": "Pools & jacuzzis", "ru": "Бассейны и джакузи"}'::jsonb,
 '{"sr": "Naduvni bazeni,Frame pool, Hot tub, oprema za bazene", "en": "Inflatable pools, Frame pools, Hot tubs, pool equipment", "ru": "Надувные бассейны, каркасные бассейны, джакузи, оборудование для бассейнов"}'::jsonb,
 '{"sr": "Bazeni i džakuzi | Vondi", "en": "Pools & jacuzzis | Vondi", "ru": "Бассейны и джакузи | Vondi"}'::jsonb,
 '{"sr": "Kupite bazene i džakuzi online", "en": "Buy pools and jacuzzis online", "ru": "Купить бассейны и джакузи онлайн"}'::jsonb,
 '🏊', true),

('sigurnost-za-dom', (SELECT id FROM categories WHERE slug = 'dom-i-basta' AND level = 1), 2, 'dom-i-basta/sigurnost-za-dom', 109,
 '{"sr": "Sigurnost za dom", "en": "Home security", "ru": "Безопасность для дома"}'::jsonb,
 '{"sr": "Kamere, senzori, alarmi, video interfoni, brave", "en": "Cameras, sensors, alarms, video intercoms, locks", "ru": "Камеры, датчики, сигнализации, видеодомофоны, замки"}'::jsonb,
 '{"sr": "Sigurnost za dom | Vondi", "en": "Home security | Vondi", "ru": "Безопасность для дома | Vondi"}'::jsonb,
 '{"sr": "Kupite sisteme sigurnosti za dom online", "en": "Buy home security systems online", "ru": "Купить системы безопасности для дома онлайн"}'::jsonb,
 '🔐', true),

-- =============================================================================
-- Additional L2 for: Lepota i zdravlje (+ 10 more)
-- =============================================================================

('masazeri', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' AND level = 1), 2, 'lepota-i-zdravlje/masazeri', 100,
 '{"sr": "Masažeri", "en": "Massagers", "ru": "Массажеры"}'::jsonb,
 '{"sr": "Električni masažeri, masažne stolice, masažni jastuk", "en": "Electric massagers, massage chairs, massage pillows", "ru": "Электрические массажеры, массажные кресла, массажные подушки"}'::jsonb,
 '{"sr": "Masažeri | Vondi", "en": "Massagers | Vondi", "ru": "Массажеры | Vondi"}'::jsonb,
 '{"sr": "Kupite masažere online", "en": "Buy massagers online", "ru": "Купить массажеры онлайн"}'::jsonb,
 '💆', true),

('medicinski-aparati', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' AND level = 1), 2, 'lepota-i-zdravlje/medicinski-aparati', 101,
 '{"sr": "Medicinski aparati", "en": "Medical devices", "ru": "Медицинские приборы"}'::jsonb,
 '{"sr": "Aparati za pritisak, glukometri, oksimetri, termometri", "en": "Blood pressure monitors, glucometers, oximeters, thermometers", "ru": "Тонометры, глюкометры, оксиметры, термометры"}'::jsonb,
 '{"sr": "Medicinski aparati | Vondi", "en": "Medical devices | Vondi", "ru": "Медицинские приборы | Vondi"}'::jsonb,
 '{"sr": "Kupite medicinske aparate online", "en": "Buy medical devices online", "ru": "Купить медицинские приборы онлайн"}'::jsonb,
 '🩺', true),

('pametne-vage', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' AND level = 1), 2, 'lepota-i-zdravlje/pametne-vage', 102,
 '{"sr": "Pametne vage", "en": "Smart scales", "ru": "Умные весы"}'::jsonb,
 '{"sr": "Body composition scale, Bluetooth vage, analiza tela", "en": "Body composition scales, Bluetooth scales, body analysis", "ru": "Весы с анализом состава тела, Bluetooth весы, анализ тела"}'::jsonb,
 '{"sr": "Pametne vage | Vondi", "en": "Smart scales | Vondi", "ru": "Умные весы | Vondi"}'::jsonb,
 '{"sr": "Kupite pametne vage online", "en": "Buy smart scales online", "ru": "Купить умные весы онлайн"}'::jsonb,
 '⚖️', true),

('fitnes-trakeri-zdravlje', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' AND level = 1), 2, 'lepota-i-zdravlje/fitnes-trakeri-zdravlje', 103,
 '{"sr": "Fitnes trakeri za zdravlje", "en": "Health fitness trackers", "ru": "Фитнес трекеры для здоровья"}'::jsonb,
 '{"sr": "Pametni satovi za praćenje zdravlja, puls, san, kalorije", "en": "Smart watches for health tracking, pulse, sleep, calories", "ru": "Умные часы для здоровья, пульс, сон, калории"}'::jsonb,
 '{"sr": "Fitnes trakeri za zdravlje | Vondi", "en": "Health fitness trackers | Vondi", "ru": "Фитнес трекеры для здоровья | Vondi"}'::jsonb,
 '{"sr": "Kupite fitnes trakere za zdravlje online", "en": "Buy health fitness trackers online", "ru": "Купить фитнес трекеры для здоровья онлайн"}'::jsonb,
 '⌚', true),

('vitamini-premium', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' AND level = 1), 2, 'lepota-i-zdravlje/vitamini-premium', 104,
 '{"sr": "Vitamini premium", "en": "Premium vitamins", "ru": "Премиум витамины"}'::jsonb,
 '{"sr": "Dodaci ishrani, suplementi, proteini, omega-3", "en": "Dietary supplements, protein, omega-3", "ru": "Пищевые добавки, протеины, омега-3"}'::jsonb,
 '{"sr": "Vitamini premium | Vondi", "en": "Premium vitamins | Vondi", "ru": "Премиум витамины | Vondi"}'::jsonb,
 '{"sr": "Kupite premium vitamine i suplemente online", "en": "Buy premium vitamins and supplements online", "ru": "Купить премиум витамины и добавки онлайн"}'::jsonb,
 '💊', true),

('ajurveda', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' AND level = 1), 2, 'lepota-i-zdravlje/ajurveda', 105,
 '{"sr": "Ajurveda", "en": "Ayurveda", "ru": "Аюрведа"}'::jsonb,
 '{"sr": "Ajurvedski čajevi, ulja, kozmetika, prirodni lekovi", "en": "Ayurvedic teas, oils, cosmetics, natural remedies", "ru": "Аюрведические чаи, масла, косметика, натуральные средства"}'::jsonb,
 '{"sr": "Ajurveda | Vondi", "en": "Ayurveda | Vondi", "ru": "Аюрведа | Vondi"}'::jsonb,
 '{"sr": "Kupite ajurvedske proizvode online", "en": "Buy ayurvedic products online", "ru": "Купить аюрведические продукты онлайн"}'::jsonb,
 '🌿', true),

-- REMOVED DUPLICATE: organska-kozmetika (already exists in part4)

('anti-age', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' AND level = 1), 2, 'lepota-i-zdravlje/anti-age', 107,
 '{"sr": "Anti-age", "en": "Anti-aging", "ru": "Антивозрастной уход"}'::jsonb,
 '{"sr": "Kreme protiv starenja, serumi, tretmani protiv bora", "en": "Anti-aging creams, serums, anti-wrinkle treatments", "ru": "Антивозрастные кремы, сыворотки, средства против морщин"}'::jsonb,
 '{"sr": "Anti-age | Vondi", "en": "Anti-aging | Vondi", "ru": "Антивозрастной уход | Vondi"}'::jsonb,
 '{"sr": "Kupite anti-age kozmetiku online", "en": "Buy anti-aging cosmetics online", "ru": "Купить антивозрастную косметику онлайн"}'::jsonb,
 '✨', true),

('muska-nega', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' AND level = 1), 2, 'lepota-i-zdravlje/muska-nega', 108,
 '{"sr": "Muška nega", "en": "Men's care", "ru": "Мужской уход"}'::jsonb,
 '{"sr": "Brijanje, parfemi, nega kože za muškarce", "en": "Shaving, perfumes, skincare for men", "ru": "Бритье, парфюмерия, уход за кожей для мужчин"}'::jsonb,
 '{"sr": "Muška nega | Vondi", "en": "Men's care | Vondi", "ru": "Мужской уход | Vondi"}'::jsonb,
 '{"sr": "Kupite kozmetiku za muškarce online", "en": "Buy men's cosmetics online", "ru": "Купить мужскую косметику онлайн"}'::jsonb,
 '🧔', true),

('salonska-oprema', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' AND level = 1), 2, 'lepota-i-zdravlje/salonska-oprema', 109,
 '{"sr": "Salonska oprema", "en": "Salon equipment", "ru": "Салонное оборудование"}'::jsonb,
 '{"sr": "Frizerska oprema, kozmetički aparati, profesionalna oprema", "en": "Hairdressing equipment, cosmetic devices, professional equipment", "ru": "Парикмахерское оборудование, косметические приборы, профессиональное оборудование"}'::jsonb,
 '{"sr": "Salonska oprema | Vondi", "en": "Salon equipment | Vondi", "ru": "Салонное оборудование | Vondi"}'::jsonb,
 '{"sr": "Kupite salonsku opremu online", "en": "Buy salon equipment online", "ru": "Купить салонное оборудование онлайн"}'::jsonb,
 '💇', true);

-- Progress notification
DO $$
BEGIN
  RAISE NOTICE 'Migration 20251217000002: Added 40 L2 categories (Odeća, Elektronika, Dom i bašta, Lepota)';
END $$;
-- Migration: Expand L2 categories (Part 5 - Second Half)
-- Date: 2025-12-17
-- Purpose: Add remaining 40 L2 categories (Bebe, Sport, Auto, Aparati)

-- =============================================================================
-- Additional L2 for: Za bebe i decu (+ 10 more)
-- =============================================================================

INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('decija-oprema-nameštaj', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/decija-oprema-namestaj', 100,
 '{"sr": "Dečija oprema i nameštaj", "en": "Children's furniture", "ru": "Детская мебель"}'::jsonb,
 '{"sr": "Krevetci, ormari, stolice za hranjenje, komode", "en": "Cribs, wardrobes, high chairs, chests", "ru": "Кроватки, шкафы, стульчики для кормления, комоды"}'::jsonb,
 '{"sr": "Dečija oprema i nameštaj | Vondi", "en": "Children's furniture | Vondi", "ru": "Детская мебель | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiji nameštaj online", "en": "Buy children's furniture online", "ru": "Купить детскую мебель онлайн"}'::jsonb,
 '🛏️', true),

('školski-ranaci', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/skolski-ranaci', 101,
 '{"sr": "Školski ranci i torbe", "en": "School backpacks & bags", "ru": "Школьные рюкзаки и сумки"}'::jsonb,
 '{"sr": "Školski rančevi, sportske torbe, pernice", "en": "School backpacks, sports bags, pencil cases", "ru": "Школьные рюкзаки, спортивные сумки, пеналы"}'::jsonb,
 '{"sr": "Školski ranci i torbe | Vondi", "en": "School backpacks & bags | Vondi", "ru": "Школьные рюкзаки и сумки | Vondi"}'::jsonb,
 '{"sr": "Kupite školske rančeve online", "en": "Buy school backpacks online", "ru": "Купить школьные рюкзаки онлайн"}'::jsonb,
 '🎒', true),

('muzičke-igračke', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/muzicke-igracke', 102,
 '{"sr": "Muzičke igračke", "en": "Musical toys", "ru": "Музыкальные игрушки"}'::jsonb,
 '{"sr": "Klavijature, bubnjevi, gitare za decu", "en": "Keyboards, drums, guitars for children", "ru": "Клавиатуры, барабаны, гитары для детей"}'::jsonb,
 '{"sr": "Muzičke igračke | Vondi", "en": "Musical toys | Vondi", "ru": "Музыкальные игрушки | Vondi"}'::jsonb,
 '{"sr": "Kupite muzičke igračke online", "en": "Buy musical toys online", "ru": "Купить музыкальные игрушки онлайн"}'::jsonb,
 '🎹', true),

('obrazovne-igre', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/obrazovne-igre', 103,
 '{"sr": "Obrazovne igre", "en": "Educational games", "ru": "Образовательные игры"}'::jsonb,
 '{"sr": "STEM igračke, edukativne slagalice, matematičke igre", "en": "STEM toys, educational puzzles, math games", "ru": "STEM игрушки, образовательные головоломки, математические игры"}'::jsonb,
 '{"sr": "Obrazovne igre | Vondi", "en": "Educational games | Vondi", "ru": "Образовательные игры | Vondi"}'::jsonb,
 '{"sr": "Kupite obrazovne igre za decu online", "en": "Buy educational games for children online", "ru": "Купить образовательные игры для детей онлайн"}'::jsonb,
 '📚', true),

('decija-elektronika', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/decija-elektronika', 104,
 '{"sr": "Dečija elektronika", "en": "Children's electronics", "ru": "Детская электроника"}'::jsonb,
 '{"sr": "Tableti za decu, smart satovi, kamere za decu", "en": "Tablets for kids, smart watches, cameras for children", "ru": "Планшеты для детей, умные часы, камеры для детей"}'::jsonb,
 '{"sr": "Dečija elektronika | Vondi", "en": "Children's electronics | Vondi", "ru": "Детская электроника | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiju elektroniku online", "en": "Buy children's electronics online", "ru": "Купить детскую электронику онлайн"}'::jsonb,
 '📱', true),

('decija-kozmetika', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/decija-kozmetika', 105,
 '{"sr": "Dečija kozmetika", "en": "Children's cosmetics", "ru": "Детская косметика"}'::jsonb,
 '{"sr": "Dečiji šamponi, sapuni, kreme, losioni", "en": "Children's shampoos, soaps, creams, lotions", "ru": "Детские шампуни, мыло, кремы, лосьоны"}'::jsonb,
 '{"sr": "Dečija kozmetika | Vondi", "en": "Children's cosmetics | Vondi", "ru": "Детская косметика | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiju kozmetiku online", "en": "Buy children's cosmetics online", "ru": "Купить детскую косметику онлайн"}'::jsonb,
 '🧴', true),

('artikli-za-novorodjence', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/artikli-za-novorodjence', 106,
 '{"sr": "Artikli za novorođenčad", "en": "Newborn items", "ru": "Товары для новорожденных"}'::jsonb,
 '{"sr": "Pelene, flašice, dude, bebi monitore", "en": "Diapers, bottles, pacifiers, baby monitors", "ru": "Подгузники, бутылочки, пустышки, радионяни"}'::jsonb,
 '{"sr": "Artikli za novorođenčad | Vondi", "en": "Newborn items | Vondi", "ru": "Товары для новорожденных | Vondi"}'::jsonb,
 '{"sr": "Kupite artikle za novorođenčad online", "en": "Buy newborn items online", "ru": "Купить товары для новорожденных онлайн"}'::jsonb,
 '👶', true),

('kupanje-bebe', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/kupanje-bebe', 107,
 '{"sr": "Kupanje bebe", "en": "Baby bathing", "ru": "Купание малыша"}'::jsonb,
 '{"sr": "Kadice, termometri za vodu, peškirići, plaštiči", "en": "Bathtubs, water thermometers, towels, raincoats", "ru": "Ванночки, термометры для воды, полотенца, плащи"}'::jsonb,
 '{"sr": "Kupanje bebe | Vondi", "en": "Baby bathing | Vondi", "ru": "Купание малыша | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za kupanje bebe online", "en": "Buy baby bathing equipment online", "ru": "Купить оборудование для купания малыша онлайн"}'::jsonb,
 '🛁', true),

('deciji-tekstil', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/deciji-tekstil', 108,
 '{"sr": "Dečiji tekstil", "en": "Children's textiles", "ru": "Детский текстиль"}'::jsonb,
 '{"sr": "Posteljina, ćebad, jastučići, peškiri za decu", "en": "Bedding, blankets, pillows, towels for children", "ru": "Постельное белье, одеяла, подушки, полотенца для детей"}'::jsonb,
 '{"sr": "Dečiji tekstil | Vondi", "en": "Children's textiles | Vondi", "ru": "Детский текстиль | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiji tekstil online", "en": "Buy children's textiles online", "ru": "Купить детский текстиль онлайн"}'::jsonb,
 '🛌', true),

('dekoracija-dečije-sobe', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' AND level = 1), 2, 'za-bebe-i-decu/dekoracija-decije-sobe', 109,
 '{"sr": "Dekoracija dečije sobe", "en": "Children's room decor", "ru": "Декор детской комнаты"}'::jsonb,
 '{"sr": "Poster, nalepnice, lampice, tepihi za decu", "en": "Posters, stickers, night lights, children's rugs", "ru": "Постеры, наклейки, ночники, детские ковры"}'::jsonb,
 '{"sr": "Dekoracija dečije sobe | Vondi", "en": "Children's room decor | Vondi", "ru": "Декор детской комнаты | Vondi"}'::jsonb,
 '{"sr": "Kupite dekoracije za dečiju sobu online", "en": "Buy children's room decorations online", "ru": "Купить декор для детской комнаты онлайн"}'::jsonb,
 '🎨', true),

-- =============================================================================
-- Additional L2 for: Sport i turizam (+ 10 more)
-- =============================================================================

('joga-i-pilates', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/joga-i-pilates', 100,
 '{"sr": "Joga i pilates", "en": "Yoga & pilates", "ru": "Йога и пилатес"}'::jsonb,
 '{"sr": "Podloge za jogu, blokovi, trake, pilates rekviziti", "en": "Yoga mats, blocks, straps, pilates props", "ru": "Коврики для йоги, блоки, ремни, реквизит для пилатеса"}'::jsonb,
 '{"sr": "Joga i pilates | Vondi", "en": "Yoga & pilates | Vondi", "ru": "Йога и пилатес | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za jogu i pilates online", "en": "Buy yoga and pilates equipment online", "ru": "Купить оборудование для йоги и пилатеса онлайн"}'::jsonb,
 '🧘', true),

('boks-i-borilacke-vestine', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/boks-i-borilacke-vestine', 101,
 '{"sr": "Boks i borilačke veštine", "en": "Boxing & martial arts", "ru": "Бокс и единоборства"}'::jsonb,
 '{"sr": "Rukavice za boks, vreća, kimona, štitnici", "en": "Boxing gloves, punching bags, kimonos, guards", "ru": "Боксерские перчатки, мешки, кимоно, щитки"}'::jsonb,
 '{"sr": "Boks i borilačke veštine | Vondi", "en": "Boxing & martial arts | Vondi", "ru": "Бокс и единоборства | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za boks i borilačke veštine online", "en": "Buy boxing and martial arts equipment online", "ru": "Купить оборудование для бокса и единоборств онлайн"}'::jsonb,
 '🥊', true),

('plivanje', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/plivanje', 102,
 '{"sr": "Plivanje", "en": "Swimming", "ru": "Плавание"}'::jsonb,
 '{"sr": "Kupaći kostimi, naočare za plivanje, kape, peraje", "en": "Swimsuits, goggles, caps, fins", "ru": "Купальники, очки для плавания, шапочки, ласты"}'::jsonb,
 '{"sr": "Plivanje | Vondi", "en": "Swimming | Vondi", "ru": "Плавание | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za plivanje online", "en": "Buy swimming equipment online", "ru": "Купить оборудование для плавания онлайн"}'::jsonb,
 '🏊', true),

('tenis', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/tenis', 103,
 '{"sr": "Tenis", "en": "Tennis", "ru": "Теннис"}'::jsonb,
 '{"sr": "Teniski reketi, loptice, torbe, oprema za tenis", "en": "Tennis rackets, balls, bags, tennis equipment", "ru": "Теннисные ракетки, мячи, сумки, оборудование для тенниса"}'::jsonb,
 '{"sr": "Tenis | Vondi", "en": "Tennis | Vondi", "ru": "Теннис | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za tenis online", "en": "Buy tennis equipment online", "ru": "Купить оборудование для тенниса онлайн"}'::jsonb,
 '🎾', true),

('odbojka', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/odbojka', 104,
 '{"sr": "Odbojka", "en": "Volleyball", "ru": "Волейбол"}'::jsonb,
 '{"sr": "Odbojkaške lopte, mreže, oprema za odbojku", "en": "Volleyballs, nets, volleyball equipment", "ru": "Волейбольные мячи, сетки, оборудование для волейбола"}'::jsonb,
 '{"sr": "Odbojka | Vondi", "en": "Volleyball | Vondi", "ru": "Волейбол | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za odbojku online", "en": "Buy volleyball equipment online", "ru": "Купить оборудование для волейбола онлайн"}'::jsonb,
 '🏐', true),

('kosarka', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/kosarka', 105,
 '{"sr": "Košarka", "en": "Basketball", "ru": "Баскетбол"}'::jsonb,
 '{"sr": "Košarkaške lopte, koševi, patike za košarku", "en": "Basketballs, hoops, basketball shoes", "ru": "Баскетбольные мячи, корзины, баскетбольные кроссовки"}'::jsonb,
 '{"sr": "Košarka | Vondi", "en": "Basketball | Vondi", "ru": "Баскетбол | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za košarku online", "en": "Buy basketball equipment online", "ru": "Купить оборудование для баскетбола онлайн"}'::jsonb,
 '🏀', true),

('fudbal-oprema', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/fudbal-oprema', 106,
 '{"sr": "Fudbal oprema", "en": "Football equipment", "ru": "Футбольная экипировка"}'::jsonb,
 '{"sr": "Fudbalske lopte, golovi, štitnici, dresovi", "en": "Footballs, goals, shin guards, jerseys", "ru": "Футбольные мячи, ворота, щитки, футболки"}'::jsonb,
 '{"sr": "Fudbal oprema | Vondi", "en": "Football equipment | Vondi", "ru": "Футбольная экипировка | Vondi"}'::jsonb,
 '{"sr": "Kupite fudbalsku opremu online", "en": "Buy football equipment online", "ru": "Купить футбольную экипировку онлайн"}'::jsonb,
 '⚽', true),

('trcanje-i-atletika', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/trcanje-i-atletika', 107,
 '{"sr": "Trčanje i atletika", "en": "Running & athletics", "ru": "Бег и легкая атлетика"}'::jsonb,
 '{"sr": "Patike za trčanje, GPS satovi, trakice, oprema", "en": "Running shoes, GPS watches, bands, equipment", "ru": "Кроссовки для бега, GPS часы, повязки, оборудование"}'::jsonb,
 '{"sr": "Trčanje i atletika | Vondi", "en": "Running & athletics | Vondi", "ru": "Бег и легкая атлетика | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za trčanje online", "en": "Buy running equipment online", "ru": "Купить оборудование для бега онлайн"}'::jsonb,
 '🏃', true),

('ekstremni-sportovi', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/ekstremni-sportovi', 108,
 '{"sr": "Ekstremni sportovi", "en": "Extreme sports", "ru": "Экстремальные виды спорта"}'::jsonb,
 '{"sr": "Penjanje, skok elastikom, base jumping, paraglajding", "en": "Climbing, bungee jumping, base jumping, paragliding", "ru": "Скалолазание, банджи-джампинг, бейсджампинг, параглайдинг"}'::jsonb,
 '{"sr": "Ekstremni sportovi | Vondi", "en": "Extreme sports | Vondi", "ru": "Экстремальные виды спорта | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za ekstremne sportove online", "en": "Buy extreme sports equipment online", "ru": "Купить оборудование для экстремальных видов спорта онлайн"}'::jsonb,
 '🪂', true),

('lov-i-ribolov', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' AND level = 1), 2, 'sport-i-turizam/lov-i-ribolov', 109,
 '{"sr": "Lov i ribolov", "en": "Hunting & fishing", "ru": "Охота и рыбалка"}'::jsonb,
 '{"sr": "Pecanje štapovi, mamci, lovačka oprema, noževi", "en": "Fishing rods, baits, hunting gear, knives", "ru": "Удочки, приманки, охотничье снаряжение, ножи"}'::jsonb,
 '{"sr": "Lov i ribolov | Vondi", "en": "Hunting & fishing | Vondi", "ru": "Охота и рыбалка | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za lov i ribolov online", "en": "Buy hunting and fishing equipment online", "ru": "Купить оборудование для охоты и рыбалки онлайн"}'::jsonb,
 '🎣', true);

-- Progress notification
DO $$
BEGIN
  RAISE NOTICE 'Migration 20251217000002 (Second Half): Added 20 L2 categories (Bebe, Sport)';
END $$;

-- =============================================================================
-- Additional L2 for: Automobilizam (+ 10 more)
-- =============================================================================

INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('elektromobil-aksesoari', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/elektromobil-aksesoari', 100,
 '{"sr": "Elektromobil aksesoari", "en": "Electric vehicle accessories", "ru": "Аксессуары для электромобилей"}'::jsonb,
 '{"sr": "Punjači za EV, kablov i, adaptori, držači", "en": "EV chargers, cables, adapters, holders", "ru": "Зарядники для EV, кабели, адаптеры, держатели"}'::jsonb,
 '{"sr": "Elektromobil aksesoari | Vondi", "en": "Electric vehicle accessories | Vondi", "ru": "Аксессуары для электромобилей | Vondi"}'::jsonb,
 '{"sr": "Kupite akses oare za elektromobile online", "en": "Buy electric vehicle accessories online", "ru": "Купить аксессуары для электромобилей онлайн"}'::jsonb,
 '🔌', true),

('tjuning', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/tjuning', 101,
 '{"sr": "Tjuning", "en": "Tuning", "ru": "Тюнинг"}'::jsonb,
 '{"sr": "Spojleri, body kits, naplatke, performanse delovi", "en": "Spoilers, body kits, rims, performance parts", "ru": "Спойлеры, обвесы, диски, тюнинг детали"}'::jsonb,
 '{"sr": "Tjuning | Vondi", "en": "Tuning | Vondi", "ru": "Тюнинг | Vondi"}'::jsonb,
 '{"sr": "Kupite tjuning delove online", "en": "Buy tuning parts online", "ru": "Купить тюнинг детали онлайн"}'::jsonb,
 '🏎️', true),

('autozvuk-premium', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/autozvuk-premium', 102,
 '{"sr": "Autozvuk premium", "en": "Premium car audio", "ru": "Премиум автозвук"}'::jsonb,
 '{"sr": "Hi-Fi zvučnici, pojačala, subwoofer, DSP procesori", "en": "Hi-Fi speakers, amplifiers, subwoofers, DSP processors", "ru": "Hi-Fi динамики, усилители, сабвуферы, DSP процессоры"}'::jsonb,
 '{"sr": "Autozvuk premium | Vondi", "en": "Premium car audio | Vondi", "ru": "Премиум автозвук | Vondi"}'::jsonb,
 '{"sr": "Kupite premium autozvuk online", "en": "Buy premium car audio online", "ru": "Купить премиум автозвук онлайн"}'::jsonb,
 '🔊', true),

('autokozmetika-premium', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/autokozmetika-premium', 103,
 '{"sr": "Autokozmetika premium", "en": "Premium car care", "ru": "Премиум автокосметика"}'::jsonb,
 '{"sr": "Detailing, keramičke zaštite, premium voskovi", "en": "Detailing, ceramic coatings, premium waxes", "ru": "Детейлинг, керамические покрытия, премиум воски"}'::jsonb,
 '{"sr": "Autokozmetika premium | Vondi", "en": "Premium car care | Vondi", "ru": "Премиум автокосметика | Vondi"}'::jsonb,
 '{"sr": "Kupite premium autokozmetiku online", "en": "Buy premium car care online", "ru": "Купить премиум автокосметику онлайн"}'::jsonb,
 '✨', true),

('krovni-nosaci', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/krovni-nosaci', 104,
 '{"sr": "Krovni nosači", "en": "Roof racks", "ru": "Багажники"}'::jsonb,
 '{"sr": "Bagažnici, box za krov, nosači bicikla, ski nosači", "en": "Roof racks, roof boxes, bike carriers, ski racks", "ru": "Багажники на крышу, боксы, велокрепления, лыжные крепления"}'::jsonb,
 '{"sr": "Krovni nosači | Vondi", "en": "Roof racks | Vondi", "ru": "Багажники | Vondi"}'::jsonb,
 '{"sr": "Kupite krovne nosače online", "en": "Buy roof racks online", "ru": "Купить багажники онлайн"}'::jsonb,
 '🚗', true),

('auto-elektronika-premium', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/auto-elektronika-premium', 105,
 '{"sr": "Auto elektronika premium", "en": "Premium car electronics", "ru": "Премиум автоэлектроника"}'::jsonb,
 '{"sr": "Android head units, CarPlay, Android Auto sistemi", "en": "Android head units, CarPlay, Android Auto systems", "ru": "Android магнитолы, CarPlay, Android Auto системы"}'::jsonb,
 '{"sr": "Auto elektronika premium | Vondi", "en": "Premium car electronics | Vondi", "ru": "Премиум автоэлектроника | Vondi"}'::jsonb,
 '{"sr": "Kupite premium auto elektroniku online", "en": "Buy premium car electronics online", "ru": "Купить премиум автоэлектронику онлайн"}'::jsonb,
 '📱', true),

('gps-navigacija', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/gps-navigacija', 106,
 '{"sr": "GPS navigacija", "en": "GPS navigation", "ru": "GPS навигация"}'::jsonb,
 '{"sr": "Garmin, TomTom, auto GPS uređaji, mape", "en": "Garmin, TomTom, car GPS devices, maps", "ru": "Garmin, TomTom, автомобильные GPS, карты"}'::jsonb,
 '{"sr": "GPS navigacija | Vondi", "en": "GPS navigation | Vondi", "ru": "GPS навигация | Vondi"}'::jsonb,
 '{"sr": "Kupite GPS navigaciju online", "en": "Buy GPS navigation online", "ru": "Купить GPS навигацию онлайн"}'::jsonb,
 '🗺️', true),

('video-registratori', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/video-registratori', 107,
 '{"sr": "Video registratori", "en": "Dash cams", "ru": "Видеорегистраторы"}'::jsonb,
 '{"sr": "DVR kamere, 4K dash cam, dual kamere, parking mode", "en": "DVR cameras, 4K dash cams, dual cameras, parking mode", "ru": "DVR камеры, 4K видеорегистраторы, двойные камеры, режим парковки"}'::jsonb,
 '{"sr": "Video registratori | Vondi", "en": "Dash cams | Vondi", "ru": "Видеорегистраторы | Vondi"}'::jsonb,
 '{"sr": "Kupite video registratore online", "en": "Buy dash cams online", "ru": "Купить видеорегистраторы онлайн"}'::jsonb,
 '📹', true),

('parking-senzori', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/parking-senzori', 108,
 '{"sr": "Parking senzori", "en": "Parking sensors", "ru": "Парковочные датчики"}'::jsonb,
 '{"sr": "Senzori za parkiranje, kamere za vožnju unazad, 360 kamere", "en": "Parking sensors, rear cameras, 360 cameras", "ru": "Парковочные датчики, камеры заднего вида, 360 камеры"}'::jsonb,
 '{"sr": "Parking senzori | Vondi", "en": "Parking sensors | Vondi", "ru": "Парковочные датчики | Vondi"}'::jsonb,
 '{"sr": "Kupite parking senzore online", "en": "Buy parking sensors online", "ru": "Купить парковочные датчики онлайн"}'::jsonb,
 '📡', true),

('auto-cehli-premium', (SELECT id FROM categories WHERE slug = 'automobilizam' AND level = 1), 2, 'automobilizam/auto-cehli-premium', 109,
 '{"sr": "Auto čehli premium", "en": "Premium car covers", "ru": "Премиум авточехлы"}'::jsonb,
 '{"sr": "Kožne navlake, eko-koža, custom fit, grejane navlake", "en": "Leather covers, eco-leather, custom fit, heated covers", "ru": "Кожаные чехлы, эко-кожа, индивидуальные, с подогревом"}'::jsonb,
 '{"sr": "Auto čehli premium | Vondi", "en": "Premium car covers | Vondi", "ru": "Премиум авточехлы | Vondi"}'::jsonb,
 '{"sr": "Kupite premium auto čehle online", "en": "Buy premium car covers online", "ru": "Купить премиум авточехлы онлайн"}'::jsonb,
 '🪑', true),

-- =============================================================================
-- Additional L2 for: Kućni aparati (+ 10 more)
-- =============================================================================

('roboti-usisivaci', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/roboti-usisivaci', 100,
 '{"sr": "Roboti usisivači", "en": "Robot vacuums", "ru": "Роботы-пылесосы"}'::jsonb,
 '{"sr": "Xiaomi, Roborock, iRobot, smart usisivači", "en": "Xiaomi, Roborock, iRobot, smart vacuums", "ru": "Xiaomi, Roborock, iRobot, умные пылесосы"}'::jsonb,
 '{"sr": "Roboti usisivači | Vondi", "en": "Robot vacuums | Vondi", "ru": "Роботы-пылесосы | Vondi"}'::jsonb,
 '{"sr": "Kupite robote usisivače online", "en": "Buy robot vacuums online", "ru": "Купить роботы-пылесосы онлайн"}'::jsonb,
 '🤖', true),

('pametni-frizidera', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/pametni-frizidera', 101,
 '{"sr": "Pametni frižideri", "en": "Smart refrigerators", "ru": "Умные холодильники"}'::jsonb,
 '{"sr": "Smart screen frižideri, WiFi, app kontrola", "en": "Smart screen refrigerators, WiFi, app control", "ru": "Умные холодильники с экраном, WiFi, управление приложением"}'::jsonb,
 '{"sr": "Pametni frižideri | Vondi", "en": "Smart refrigerators | Vondi", "ru": "Умные холодильники | Vondi"}'::jsonb,
 '{"sr": "Kupite pametne frižidere online", "en": "Buy smart refrigerators online", "ru": "Купить умные холодильники онлайн"}'::jsonb,
 '🧊', true),

('vinski-frizideri', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/vinski-frizideri', 102,
 '{"sr": "Vinski frižideri", "en": "Wine coolers", "ru": "Винные шкафы"}'::jsonb,
 '{"sr": "Hladnjaci za vino, wine cellars, dual zone", "en": "Wine coolers, wine cellars, dual zone", "ru": "Винные холодильники, винные погреба, двухзонные"}'::jsonb,
 '{"sr": "Vinski frižideri | Vondi", "en": "Wine coolers | Vondi", "ru": "Винные шкафы | Vondi"}'::jsonb,
 '{"sr": "Kupite vinske frižidere online", "en": "Buy wine coolers online", "ru": "Купить винные шкафы онлайн"}'::jsonb,
 '🍷', true),

('sudopere-premium', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/sudopere-premium', 103,
 '{"sr": "Sudomašine premium", "en": "Premium dishwashers", "ru": "Премиум посудомойки"}'::jsonb,
 '{"sr": "Bosch, Miele, tihe sudomašine, A+++ klasa", "en": "Bosch, Miele, quiet dishwashers, A+++ class", "ru": "Bosch, Miele, тихие посудомойки, класс A+++"}'::jsonb,
 '{"sr": "Sudomašine premium | Vondi", "en": "Premium dishwashers | Vondi", "ru": "Премиум посудомойки | Vondi"}'::jsonb,
 '{"sr": "Kupite premium sudomašine online", "en": "Buy premium dishwashers online", "ru": "Купить премиум посудомойки онлайн"}'::jsonb,
 '🍽️', true),

('ves-masine-premium', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/ves-masine-premium', 104,
 '{"sr": "Veš mašine premium", "en": "Premium washing machines", "ru": "Премиум стиральные машины"}'::jsonb,
 '{"sr": "Front load, inverter motor, A+++ klasa, steam wash", "en": "Front load, inverter motor, A+++ class, steam wash", "ru": "Фронтальная загрузка, инверторный мотор, класс A+++, паровая стирка"}'::jsonb,
 '{"sr": "Veš mašine premium | Vondi", "en": "Premium washing machines | Vondi", "ru": "Премиум стиральные машины | Vondi"}'::jsonb,
 '{"sr": "Kupite premium veš mašine online", "en": "Buy premium washing machines online", "ru": "Купить премиум стиральные машины онлайн"}'::jsonb,
 '🧺', true),

('susare-masine', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/susare-masine', 105,
 '{"sr": "Sušare mašine", "en": "Dryers", "ru": "Сушильные машины"}'::jsonb,
 '{"sr": "Heat pump sušare, kondenzacione, ventilacione", "en": "Heat pump dryers, condensation, ventilation", "ru": "Тепловой насос сушилки, конденсационные, вентиляционные"}'::jsonb,
 '{"sr": "Sušare mašine | Vondi", "en": "Dryers | Vondi", "ru": "Сушильные машины | Vondi"}'::jsonb,
 '{"sr": "Kupite sušare mašine online", "en": "Buy dryers online", "ru": "Купить сушильные машины онлайн"}'::jsonb,
 '🌬️', true),

('parne-stanice', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/parne-stanice', 106,
 '{"sr": "Parne stanice", "en": "Steam stations", "ru": "Паровые станции"}'::jsonb,
 '{"sr": "Pegla sa rezervoarom, profesionalno peglanje", "en": "Iron with tank, professional ironing", "ru": "Утюг с резервуаром, профессиональная глажка"}'::jsonb,
 '{"sr": "Parne stanice | Vondi", "en": "Steam stations | Vondi", "ru": "Паровые станции | Vondi"}'::jsonb,
 '{"sr": "Kupite parne stanice online", "en": "Buy steam stations online", "ru": "Купить паровые станции онлайн"}'::jsonb,
 '♨️', true),

('multivarka', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/multivarka', 107,
 '{"sr": "Multivarka", "en": "Multicookers", "ru": "Мультиварки"}'::jsonb,
 '{"sr": "Instant Pot, pritisak lonci, slow cooker", "en": "Instant Pot, pressure cookers, slow cooker", "ru": "Instant Pot, скороварки, медленноварки"}'::jsonb,
 '{"sr": "Multivarka | Vondi", "en": "Multicookers | Vondi", "ru": "Мультиварки | Vondi"}'::jsonb,
 '{"sr": "Kupite multivarke online", "en": "Buy multicookers online", "ru": "Купить мультиварки онлайн"}'::jsonb,
 '🍲', true),

('aerofriteze', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/aerofriteze', 108,
 '{"sr": "Aerofriteze", "en": "Air fryers", "ru": "Аэрогрили"}'::jsonb,
 '{"sr": "Ninja, Philips, friteza bez ulja, hot air", "en": "Ninja, Philips, oil-free fryer, hot air", "ru": "Ninja, Philips, фритюрница без масла, горячий воздух"}'::jsonb,
 '{"sr": "Aerofriteze | Vondi", "en": "Air fryers | Vondi", "ru": "Аэрогрили | Vondi"}'::jsonb,
 '{"sr": "Kupite aerofriteze online", "en": "Buy air fryers online", "ru": "Купить аэрогрили онлайн"}'::jsonb,
 '🍟', true),

('hlebopekare', (SELECT id FROM categories WHERE slug = 'kucni-aparati' AND level = 1), 2, 'kucni-aparati/hlebopekare', 109,
 '{"sr": "Hlebopekare", "en": "Bread makers", "ru": "Хлебопечки"}'::jsonb,
 '{"sr": "Automatske hlebopekare, gluten-free programi", "en": "Automatic bread makers, gluten-free programs", "ru": "Автоматические хлебопечки, безглютеновые программы"}'::jsonb,
 '{"sr": "Hlebopekare | Vondi", "en": "Bread makers | Vondi", "ru": "Хлебопечки | Vondi"}'::jsonb,
 '{"sr": "Kupite hlebopekare online", "en": "Buy bread makers online", "ru": "Купить хлебопечки онлайн"}'::jsonb,
 '🍞', true);

-- Final progress notification
DO $$
BEGIN
  RAISE NOTICE 'Migration 20251217000002 COMPLETED: Total added 80 L2 categories (Odeća, Elektronika, Dom, Lepota, Bebe, Sport, Auto, Aparati)';
END $$;
