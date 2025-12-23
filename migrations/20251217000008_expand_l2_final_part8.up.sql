-- Migration: Expand L2 categories (Part 8 - FINAL)
-- Date: 2025-12-17
-- Purpose: Add final 62 L2 categories - Services, Pets, Books, Misc
-- Expanding: Usluge (23), Kucni ljubimci (12), Knjige (10), Ostalo (17)

-- =============================================================================
-- L2 for: Usluge (23 new categories)
-- =============================================================================

INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('elektricar', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/elektricar', 100,
 '{"sr": "Električar", "en": "Electrician", "ru": "Электрик"}'::jsonb,
 '{"sr": "Usluge profesionalnog električara, instalacije, popravke", "en": "Professional electrician services, installations, repairs", "ru": "Услуги профессионального электрика, установка, ремонт"}'::jsonb,
 '{"sr": "Električar | Vondi", "en": "Electrician | Vondi", "ru": "Электрик | Vondi"}'::jsonb,
 '{"sr": "Pronađite električara - usluge električara online", "en": "Find electrician services online", "ru": "Найти электрика онлайн"}'::jsonb,
 '⚡', true),

('vodoinstalater', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/vodoinstalater', 101,
 '{"sr": "Vodoinstalater", "en": "Plumber", "ru": "Сантехник"}'::jsonb,
 '{"sr": "Usluge vodoinstalatera, popravke, instalacije", "en": "Plumbing services, repairs, installations", "ru": "Услуги сантехника, ремонт, установка"}'::jsonb,
 '{"sr": "Vodoinstalater | Vondi", "en": "Plumber | Vondi", "ru": "Сантехник | Vondi"}'::jsonb,
 '{"sr": "Pronađite vodoinstalatera - brze intervencije", "en": "Find plumber online", "ru": "Найти сантехника онлайн"}'::jsonb,
 '🚰', true),

('moler-farbanje', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/moler-farbanje', 102,
 '{"sr": "Moler i farbanje", "en": "Painter and painting", "ru": "Маляр и покраска"}'::jsonb,
 '{"sr": "Usluge molera, farbanje zidova, fasada, stanova", "en": "Painter services, wall painting, facades, apartments", "ru": "Услуги маляра, покраска стен, фасадов, квартир"}'::jsonb,
 '{"sr": "Moler i farbanje | Vondi", "en": "Painter and painting | Vondi", "ru": "Маляр и покраска | Vondi"}'::jsonb,
 '{"sr": "Pronađite molera - farbanje zidova i stanova", "en": "Find painter online", "ru": "Найти маляра онлайн"}'::jsonb,
 '🎨', true),

('stolar', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/stolar', 103,
 '{"sr": "Stolar", "en": "Carpenter", "ru": "Столяр"}'::jsonb,
 '{"sr": "Usluge stolara, izrada nameštaja, drveni radovi", "en": "Carpenter services, furniture making, woodwork", "ru": "Услуги столяра, изготовление мебели, деревообработка"}'::jsonb,
 '{"sr": "Stolar | Vondi", "en": "Carpenter | Vondi", "ru": "Столяр | Vondi"}'::jsonb,
 '{"sr": "Pronađite stolara - izrada nameštaja po meri", "en": "Find carpenter online", "ru": "Найти столяра онлайн"}'::jsonb,
 '🪚', true),

('bravar', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/bravar', 104,
 '{"sr": "Bravar", "en": "Locksmith", "ru": "Слесарь"}'::jsonb,
 '{"sr": "Usluge bravara, brave, kapije, ograde, metalni radovi", "en": "Locksmith services, locks, gates, fences, metalwork", "ru": "Услуги слесаря, замки, ворота, ограждения, металлоконструкции"}'::jsonb,
 '{"sr": "Bravar | Vondi", "en": "Locksmith | Vondi", "ru": "Слесарь | Vondi"}'::jsonb,
 '{"sr": "Pronađite bravara - metalni radovi i brave", "en": "Find locksmith online", "ru": "Найти слесаря онлайн"}'::jsonb,
 '🔒', true),

('krovopokrivac', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/krovopokrivac', 105,
 '{"sr": "Krovopokrivač", "en": "Roofer", "ru": "Кровельщик"}'::jsonb,
 '{"sr": "Usluge krovopokrivača, popravka i postavljanje krovova", "en": "Roofing services, roof repair and installation", "ru": "Кровельные услуги, ремонт и установка крыш"}'::jsonb,
 '{"sr": "Krovopokrivač | Vondi", "en": "Roofer | Vondi", "ru": "Кровельщик | Vondi"}'::jsonb,
 '{"sr": "Pronađite krovopokrivača - postavljanje krovova", "en": "Find roofer online", "ru": "Найти кровельщика онлайн"}'::jsonb,
 '🏠', true),

('keramicar', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/keramicar', 106,
 '{"sr": "Keramičar", "en": "Tile setter", "ru": "Плиточник"}'::jsonb,
 '{"sr": "Usluge keramičara, postavljanje pločica, kupatila", "en": "Tile setting services, tiling, bathrooms", "ru": "Услуги плиточника, укладка плитки, ванные комнаты"}'::jsonb,
 '{"sr": "Keramičar | Vondi", "en": "Tile setter | Vondi", "ru": "Плиточник | Vondi"}'::jsonb,
 '{"sr": "Pronađite keramičara - postavljanje pločica", "en": "Find tile setter online", "ru": "Найти плиточника онлайн"}'::jsonb,
 '🔲', true),

('parketar', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/parketar', 107,
 '{"sr": "Parketar", "en": "Flooring installer", "ru": "Паркетчик"}'::jsonb,
 '{"sr": "Usluge parketara, postavljanje parketa, laminata", "en": "Flooring installation, parquet, laminate", "ru": "Услуги паркетчика, укладка паркета, ламината"}'::jsonb,
 '{"sr": "Parketar | Vondi", "en": "Flooring installer | Vondi", "ru": "Паркетчик | Vondi"}'::jsonb,
 '{"sr": "Pronađite parketara - postavljanje podova", "en": "Find flooring installer online", "ru": "Найти паркетчика онлайн"}'::jsonb,
 '🪵', true),

('majstor-general', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/majstor-general', 108,
 '{"sr": "Majstor - opšti radovi", "en": "Handyman - general work", "ru": "Мастер на все руки"}'::jsonb,
 '{"sr": "Majstor za sve radove, sitne popravke, kućni poslovi", "en": "Handyman for all work, small repairs, household tasks", "ru": "Мастер для всех работ, мелкий ремонт, домашние задачи"}'::jsonb,
 '{"sr": "Majstor - opšti radovi | Vondi", "en": "Handyman | Vondi", "ru": "Мастер на все руки | Vondi"}'::jsonb,
 '{"sr": "Pronađite majstora - svi kućni poslovi", "en": "Find handyman online", "ru": "Найти мастера онлайн"}'::jsonb,
 '🛠️', true),

('servis-bele-tehnike', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/servis-bele-tehnike', 109,
 '{"sr": "Servis bele tehnike", "en": "Appliance repair", "ru": "Ремонт бытовой техники"}'::jsonb,
 '{"sr": "Popravka mašina za pranje, frižidera, šporeta", "en": "Repair of washing machines, fridges, stoves", "ru": "Ремонт стиральных машин, холодильников, плит"}'::jsonb,
 '{"sr": "Servis bele tehnike | Vondi", "en": "Appliance repair | Vondi", "ru": "Ремонт бытовой техники | Vondi"}'::jsonb,
 '{"sr": "Popravka bele tehnike - brzo i kvalitetno", "en": "Appliance repair online", "ru": "Ремонт бытовой техники онлайн"}'::jsonb,
 '🔧', true),

('servis-laptopa', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/servis-laptopa', 110,
 '{"sr": "Servis laptopa", "en": "Laptop repair", "ru": "Ремонт ноутбуков"}'::jsonb,
 '{"sr": "Popravka laptopa, čišćenje, zamena delova", "en": "Laptop repair, cleaning, parts replacement", "ru": "Ремонт ноутбуков, чистка, замена деталей"}'::jsonb,
 '{"sr": "Servis laptopa | Vondi", "en": "Laptop repair | Vondi", "ru": "Ремонт ноутбуков | Vondi"}'::jsonb,
 '{"sr": "Popravka laptopa - profesionalni servis", "en": "Laptop repair online", "ru": "Ремонт ноутбуков онлайн"}'::jsonb,
 '💻', true),

('servis-telefona', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/servis-telefona', 111,
 '{"sr": "Servis telefona", "en": "Phone repair", "ru": "Ремонт телефонов"}'::jsonb,
 '{"sr": "Popravka mobilnih telefona, zamena ekrana, baterije", "en": "Mobile phone repair, screen replacement, battery", "ru": "Ремонт мобильных телефонов, замена экранов, батарей"}'::jsonb,
 '{"sr": "Servis telefona | Vondi", "en": "Phone repair | Vondi", "ru": "Ремонт телефонов | Vondi"}'::jsonb,
 '{"sr": "Popravka telefona - brza zamena ekrana", "en": "Phone repair online", "ru": "Ремонт телефонов онлайн"}'::jsonb,
 '📱', true),

('servis-tv', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/servis-tv', 112,
 '{"sr": "Servis televizora", "en": "TV repair", "ru": "Ремонт телевизоров"}'::jsonb,
 '{"sr": "Popravka TV uređaja, LED, LCD, OLED", "en": "TV repair, LED, LCD, OLED", "ru": "Ремонт телевизоров, LED, LCD, OLED"}'::jsonb,
 '{"sr": "Servis televizora | Vondi", "en": "TV repair | Vondi", "ru": "Ремонт телевизоров | Vondi"}'::jsonb,
 '{"sr": "Popravka televizora - svi brendovi", "en": "TV repair online", "ru": "Ремонт телевизоров онлайн"}'::jsonb,
 '📺', true),

('fotograf-event', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/fotograf-event', 113,
 '{"sr": "Fotograf za događaje", "en": "Event photographer", "ru": "Фотограф мероприятий"}'::jsonb,
 '{"sr": "Fotografisanje venčanja, krštenja, rođendana, event-a", "en": "Photography for weddings, christenings, birthdays, events", "ru": "Фотография свадеб, крестин, дней рождения, мероприятий"}'::jsonb,
 '{"sr": "Fotograf za događaje | Vondi", "en": "Event photographer | Vondi", "ru": "Фотограф мероприятий | Vondi"}'::jsonb,
 '{"sr": "Profesionalni fotograf za sve događaje", "en": "Professional event photographer", "ru": "Профессиональный фотограф мероприятий"}'::jsonb,
 '📸', true),

('fotograf-portret', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/fotograf-portret', 114,
 '{"sr": "Fotograf - portret", "en": "Portrait photographer", "ru": "Фотограф портретов"}'::jsonb,
 '{"sr": "Portretno fotografisanje, studijske slike, profesionalno", "en": "Portrait photography, studio photos, professional", "ru": "Портретная фотография, студийные снимки, профессионально"}'::jsonb,
 '{"sr": "Fotograf - portret | Vondi", "en": "Portrait photographer | Vondi", "ru": "Фотограф портретов | Vondi"}'::jsonb,
 '{"sr": "Profesionalni portretni fotograf", "en": "Professional portrait photographer", "ru": "Профессиональный портретный фотограф"}'::jsonb,
 '👤', true),

('snimanje-video', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/snimanje-video', 115,
 '{"sr": "Video snimanje", "en": "Video recording", "ru": "Видеосъемка"}'::jsonb,
 '{"sr": "Profesionalno video snimanje događaja, reklama", "en": "Professional video recording of events, commercials", "ru": "Профессиональная видеосъемка мероприятий, рекламы"}'::jsonb,
 '{"sr": "Video snimanje | Vondi", "en": "Video recording | Vondi", "ru": "Видеосъемка | Vondi"}'::jsonb,
 '{"sr": "Profesionalno video snimanje - svi eventi", "en": "Professional video recording", "ru": "Профессиональная видеосъемка"}'::jsonb,
 '🎥', true),

('montaza-video', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/montaza-video', 116,
 '{"sr": "Montaža videa", "en": "Video editing", "ru": "Монтаж видео"}'::jsonb,
 '{"sr": "Profesionalna montaža videa, color grading, efekti", "en": "Professional video editing, color grading, effects", "ru": "Профессиональный монтаж видео, цветокоррекция, эффекты"}'::jsonb,
 '{"sr": "Montaža videa | Vondi", "en": "Video editing | Vondi", "ru": "Монтаж видео | Vondi"}'::jsonb,
 '{"sr": "Profesionalna montaža video zapisa", "en": "Professional video editing", "ru": "Профессиональный монтаж видео"}'::jsonb,
 '🎬', true),

('web-dizajn', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/web-dizajn', 117,
 '{"sr": "Web dizajn", "en": "Web design", "ru": "Веб-дизайн"}'::jsonb,
 '{"sr": "Izrada web sajtova, dizajn, responsivnost", "en": "Website creation, design, responsiveness", "ru": "Создание веб-сайтов, дизайн, адаптивность"}'::jsonb,
 '{"sr": "Web dizajn | Vondi", "en": "Web design | Vondi", "ru": "Веб-дизайн | Vondi"}'::jsonb,
 '{"sr": "Profesionalan web dizajn - moderan sajt", "en": "Professional web design", "ru": "Профессиональный веб-дизайн"}'::jsonb,
 '🌐', true),

('graficki-dizajn', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/graficki-dizajn', 118,
 '{"sr": "Grafički dizajn", "en": "Graphic design", "ru": "Графический дизайн"}'::jsonb,
 '{"sr": "Logo, branding, vizitke, plakati, grafički materijal", "en": "Logo, branding, business cards, posters, graphics", "ru": "Логотип, брендинг, визитки, плакаты, графика"}'::jsonb,
 '{"sr": "Grafički dizajn | Vondi", "en": "Graphic design | Vondi", "ru": "Графический дизайн | Vondi"}'::jsonb,
 '{"sr": "Profesionalan grafički dizajn - svi formati", "en": "Professional graphic design", "ru": "Профессиональный графический дизайн"}'::jsonb,
 '🎨', true),

('prevod-jezik', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/prevod-jezik', 119,
 '{"sr": "Prevod jezika", "en": "Translation services", "ru": "Переводческие услуги"}'::jsonb,
 '{"sr": "Profesionalni prevod dokumenata, overen prevod", "en": "Professional document translation, certified", "ru": "Профессиональный перевод документов, заверенный"}'::jsonb,
 '{"sr": "Prevod jezika | Vondi", "en": "Translation services | Vondi", "ru": "Переводческие услуги | Vondi"}'::jsonb,
 '{"sr": "Profesionalni prevod - svi jezici", "en": "Professional translation", "ru": "Профессиональный перевод"}'::jsonb,
 '🌍', true),

('prepis-teksta', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/prepis-teksta', 120,
 '{"sr": "Prepis teksta", "en": "Text transcription", "ru": "Транскрипция текста"}'::jsonb,
 '{"sr": "Prepis audio, video zapisa, dokumenata", "en": "Transcription of audio, video, documents", "ru": "Транскрипция аудио, видео, документов"}'::jsonb,
 '{"sr": "Prepis teksta | Vondi", "en": "Text transcription | Vondi", "ru": "Транскрипция текста | Vondi"}'::jsonb,
 '{"sr": "Prepis teksta - brzo i kvalitetno", "en": "Text transcription services", "ru": "Услуги транскрипции текста"}'::jsonb,
 '📝', true),

('korepetitor-matematika', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/korepetitor-matematika', 121,
 '{"sr": "Korepetitor matematike", "en": "Math tutor", "ru": "Репетитор математики"}'::jsonb,
 '{"sr": "Časovi matematike, priprema za ispite, osnovci i srednjoškolci", "en": "Math lessons, exam preparation, primary and secondary students", "ru": "Уроки математики, подготовка к экзаменам, школьники"}'::jsonb,
 '{"sr": "Korepetitor matematike | Vondi", "en": "Math tutor | Vondi", "ru": "Репетитор математики | Vondi"}'::jsonb,
 '{"sr": "Časovi matematike - iskusni profesori", "en": "Math tutoring", "ru": "Репетиторство по математике"}'::jsonb,
 '🧮', true),

('korepetitor-engleski', (SELECT id FROM categories WHERE slug = 'usluge' AND level = 1), 2, 'usluge/korepetitor-engleski', 122,
 '{"sr": "Korepetitor engleskog", "en": "English tutor", "ru": "Репетитор английского"}'::jsonb,
 '{"sr": "Časovi engleskog jezika, konverzacija, priprema za ispit", "en": "English lessons, conversation, exam preparation", "ru": "Уроки английского, разговорная практика, подготовка к экзамену"}'::jsonb,
 '{"sr": "Korepetitor engleskog | Vondi", "en": "English tutor | Vondi", "ru": "Репетитор английского | Vondi"}'::jsonb,
 '{"sr": "Časovi engleskog - svi nivoi", "en": "English tutoring", "ru": "Репетиторство по английскому"}'::jsonb,
 '🇬🇧', true);

-- =============================================================================
-- L2 for: Kucni ljubimci (12 new categories)
-- =============================================================================

INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('hrana-za-pse', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/hrana-za-pse', 100,
 '{"sr": "Hrana za pse", "en": "Dog food", "ru": "Корм для собак"}'::jsonb,
 '{"sr": "Suva i vlažna hrana za pse, sve rase i uzrasti", "en": "Dry and wet dog food, all breeds and ages", "ru": "Сухой и влажный корм для собак, все породы и возрасты"}'::jsonb,
 '{"sr": "Hrana za pse | Vondi", "en": "Dog food | Vondi", "ru": "Корм для собак | Vondi"}'::jsonb,
 '{"sr": "Kupite hranu za pse online - premium kvalitet", "en": "Buy dog food online", "ru": "Купить корм для собак онлайн"}'::jsonb,
 '🐕', true),

('hrana-za-macke', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/hrana-za-macke', 101,
 '{"sr": "Hrana za mačke", "en": "Cat food", "ru": "Корм для кошек"}'::jsonb,
 '{"sr": "Suva i vlažna hrana za mačke, različiti ukusi", "en": "Dry and wet cat food, various flavors", "ru": "Сухой и влажный корм для кошек, различные вкусы"}'::jsonb,
 '{"sr": "Hrana za mačke | Vondi", "en": "Cat food | Vondi", "ru": "Корм для кошек | Vondi"}'::jsonb,
 '{"sr": "Kupite hranu za mačke online - sve vrste", "en": "Buy cat food online", "ru": "Купить корм для кошек онлайн"}'::jsonb,
 '🐈', true),

('hrana-za-ptice', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/hrana-za-ptice', 102,
 '{"sr": "Hrana za ptice", "en": "Bird food", "ru": "Корм для птиц"}'::jsonb,
 '{"sr": "Semenke, smeše za papagaje, kanarince, tigrice", "en": "Seeds, mixes for parrots, canaries, budgies", "ru": "Семена, смеси для попугаев, канареек, волнистых попугаев"}'::jsonb,
 '{"sr": "Hrana za ptice | Vondi", "en": "Bird food | Vondi", "ru": "Корм для птиц | Vondi"}'::jsonb,
 '{"sr": "Kupite hranu za ptice online - kvalitetna", "en": "Buy bird food online", "ru": "Купить корм для птиц онлайн"}'::jsonb,
 '🐦', true),

('hrana-za-glodare', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/hrana-za-glodare', 103,
 '{"sr": "Hrana za glodare", "en": "Rodent food", "ru": "Корм для грызунов"}'::jsonb,
 '{"sr": "Hrana za hrčke, zamorce, kunići, činčile", "en": "Food for hamsters, guinea pigs, rabbits, chinchillas", "ru": "Корм для хомяков, морских свинок, кроликов, шиншилл"}'::jsonb,
 '{"sr": "Hrana za glodare | Vondi", "en": "Rodent food | Vondi", "ru": "Корм для грызунов | Vondi"}'::jsonb,
 '{"sr": "Kupite hranu za glodare online", "en": "Buy rodent food online", "ru": "Купить корм для грызунов онлайн"}'::jsonb,
 '🐹', true),

('igracke-za-pse', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/igracke-za-pse', 104,
 '{"sr": "Igračke za pse", "en": "Dog toys", "ru": "Игрушки для собак"}'::jsonb,
 '{"sr": "Loptice, užadi, gumene igračke, interaktivne igračke", "en": "Balls, ropes, rubber toys, interactive toys", "ru": "Мячики, канаты, резиновые игрушки, интерактивные игрушки"}'::jsonb,
 '{"sr": "Igračke za pse | Vondi", "en": "Dog toys | Vondi", "ru": "Игрушки для собак | Vondi"}'::jsonb,
 '{"sr": "Kupite igračke za pse - sve vrste", "en": "Buy dog toys online", "ru": "Купить игрушки для собак онлайн"}'::jsonb,
 '🎾', true),

('igracke-za-macke', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/igracke-za-macke', 105,
 '{"sr": "Igračke za mačke", "en": "Cat toys", "ru": "Игрушки для кошек"}'::jsonb,
 '{"sr": "Miševi, štapići, perušci, laser igračke", "en": "Mice, wands, feathers, laser toys", "ru": "Мышки, палочки, перья, лазерные игрушки"}'::jsonb,
 '{"sr": "Igračke za mačke | Vondi", "en": "Cat toys | Vondi", "ru": "Игрушки для кошек | Vondi"}'::jsonb,
 '{"sr": "Kupite igračke za mačke - zabava", "en": "Buy cat toys online", "ru": "Купить игрушки для кошек онлайн"}'::jsonb,
 '🐁', true),

('oprema-za-setnju', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/oprema-za-setnju', 106,
 '{"sr": "Oprema za šetnju", "en": "Walking gear", "ru": "Снаряжение для прогулок"}'::jsonb,
 '{"sr": "Povodnici, oprtenjači, flexi, torbe za izmet", "en": "Leashes, harnesses, flexi, poop bags", "ru": "Поводки, шлейки, рулетки, пакеты для отходов"}'::jsonb,
 '{"sr": "Oprema za šetnju | Vondi", "en": "Walking gear | Vondi", "ru": "Снаряжение для прогулок | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za šetnju pasa", "en": "Buy walking gear online", "ru": "Купить снаряжение для прогулок онлайн"}'::jsonb,
 '🦮', true),

('kreveti-za-ljubimce', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/kreveti-za-ljubimce', 107,
 '{"sr": "Kreveti za ljubimce", "en": "Pet beds", "ru": "Лежанки для питомцев"}'::jsonb,
 '{"sr": "Kreveti, jastučići, kućice za pse i mačke", "en": "Beds, cushions, houses for dogs and cats", "ru": "Лежанки, подушки, домики для собак и кошек"}'::jsonb,
 '{"sr": "Kreveti za ljubimce | Vondi", "en": "Pet beds | Vondi", "ru": "Лежанки для питомцев | Vondi"}'::jsonb,
 '{"sr": "Kupite krevete za pse i mačke online", "en": "Buy pet beds online", "ru": "Купить лежанки для питомцев онлайн"}'::jsonb,
 '🛏️', true),

('toaleta-za-ljubimce', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/toaleta-za-ljubimce', 108,
 '{"sr": "Toaleta za ljubimce", "en": "Pet toilets", "ru": "Туалеты для питомцев"}'::jsonb,
 '{"sr": "WC posude za mačke, posipke, podloge za pse", "en": "Cat litter boxes, litters, dog pads", "ru": "Лотки для кошек, наполнители, пеленки для собак"}'::jsonb,
 '{"sr": "Toaleta za ljubimce | Vondi", "en": "Pet toilets | Vondi", "ru": "Туалеты для питомцев | Vondi"}'::jsonb,
 '{"sr": "Kupite WC opremu za ljubimce", "en": "Buy pet toilet supplies online", "ru": "Купить туалеты для питомцев онлайн"}'::jsonb,
 '🚽', true),

('akvarijumi', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/akvarijumi', 109,
 '{"sr": "Akvarijumi", "en": "Aquariums", "ru": "Аквариумы"}'::jsonb,
 '{"sr": "Akvarijumi, filteri, pumpe, osvetljenje, dekoracije", "en": "Aquariums, filters, pumps, lighting, decorations", "ru": "Аквариумы, фильтры, насосы, освещение, декорации"}'::jsonb,
 '{"sr": "Akvarijumi | Vondi", "en": "Aquariums | Vondi", "ru": "Аквариумы | Vondi"}'::jsonb,
 '{"sr": "Kupite akvarijume i opremu online", "en": "Buy aquariums online", "ru": "Купить аквариумы онлайн"}'::jsonb,
 '🐠', true),

('terarijumi', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/terarijumi', 110,
 '{"sr": "Terarijumi", "en": "Terrariums", "ru": "Террариумы"}'::jsonb,
 '{"sr": "Terarijumi za gmizavce, osvetljenje, grejači", "en": "Terrariums for reptiles, lighting, heaters", "ru": "Террариумы для рептилий, освещение, обогреватели"}'::jsonb,
 '{"sr": "Terarijumi | Vondi", "en": "Terrariums | Vondi", "ru": "Террариумы | Vondi"}'::jsonb,
 '{"sr": "Kupite terarijume i opremu online", "en": "Buy terrariums online", "ru": "Купить террариумы онлайн"}'::jsonb,
 '🦎', true),

('nosilje-za-ljubimce', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' AND level = 1), 2, 'kucni-ljubimci/nosilje-za-ljubimce', 111,
 '{"sr": "Nosiljke za ljubimce", "en": "Pet carriers", "ru": "Переноски для питомцев"}'::jsonb,
 '{"sr": "Transportne torbe, kavezi, nosiljke za pse i mačke", "en": "Transport bags, cages, carriers for dogs and cats", "ru": "Транспортные сумки, клетки, переноски для собак и кошек"}'::jsonb,
 '{"sr": "Nosiljke za ljubimce | Vondi", "en": "Pet carriers | Vondi", "ru": "Переноски для питомцев | Vondi"}'::jsonb,
 '{"sr": "Kupite nosiljke za transport ljubimaca", "en": "Buy pet carriers online", "ru": "Купить переноски для питомцев онлайн"}'::jsonb,
 '🧳', true);

-- =============================================================================
-- L2 for: Knjige i mediji (10 new categories)
-- =============================================================================

INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('biografije', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/biografije', 100,
 '{"sr": "Biografije", "en": "Biographies", "ru": "Биографии"}'::jsonb,
 '{"sr": "Životne priče poznatih ličnosti, istorijske figure", "en": "Life stories of famous people, historical figures", "ru": "Жизненные истории известных людей, исторические фигуры"}'::jsonb,
 '{"sr": "Biografije | Vondi", "en": "Biographies | Vondi", "ru": "Биографии | Vondi"}'::jsonb,
 '{"sr": "Kupite biografije poznatih ljudi online", "en": "Buy biographies online", "ru": "Купить биографии онлайн"}'::jsonb,
 '📖', true),

('istorijske-knjige', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/istorijske-knjige', 101,
 '{"sr": "Istorijske knjige", "en": "History books", "ru": "Исторические книги"}'::jsonb,
 '{"sr": "Knjige o istoriji, ratovi, civilizacije, događaji", "en": "Books about history, wars, civilizations, events", "ru": "Книги об истории, войнах, цивилизациях, событиях"}'::jsonb,
 '{"sr": "Istorijske knjige | Vondi", "en": "History books | Vondi", "ru": "Исторические книги | Vondi"}'::jsonb,
 '{"sr": "Kupite istorijske knjige online", "en": "Buy history books online", "ru": "Купить исторические книги онлайн"}'::jsonb,
 '🏛️', true),

('fantastika', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/fantastika', 102,
 '{"sr": "Fantastika", "en": "Fantasy", "ru": "Фантастика"}'::jsonb,
 '{"sr": "Fantasy romani, epicka fantastika, vampiri, vilenjaci", "en": "Fantasy novels, epic fantasy, vampires, elves", "ru": "Фэнтезийные романы, эпическое фэнтези, вампиры, эльфы"}'::jsonb,
 '{"sr": "Fantastika | Vondi", "en": "Fantasy | Vondi", "ru": "Фантастика | Vondi"}'::jsonb,
 '{"sr": "Kupite fantasy knjige online", "en": "Buy fantasy books online", "ru": "Купить фэнтези книги онлайн"}'::jsonb,
 '🧙', true),

('sf-knjige', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/sf-knjige', 103,
 '{"sr": "Naučna fantastika", "en": "Science fiction", "ru": "Научная фантастика"}'::jsonb,
 '{"sr": "SF romani, distopija, space opera, cyberpunk", "en": "SF novels, dystopia, space opera, cyberpunk", "ru": "НФ романы, дистопия, космическая опера, киберпанк"}'::jsonb,
 '{"sr": "Naučna fantastika | Vondi", "en": "Science fiction | Vondi", "ru": "Научная фантастика | Vondi"}'::jsonb,
 '{"sr": "Kupite SF knjige online", "en": "Buy science fiction online", "ru": "Купить научную фантастику онлайн"}'::jsonb,
 '🚀', true),

('kriminalisticki-romani', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/kriminalisticki-romani', 104,
 '{"sr": "Kriminalistički romani", "en": "Crime novels", "ru": "Детективные романы"}'::jsonb,
 '{"sr": "Detektivske priče, triler, misteriozne istrage", "en": "Detective stories, thrillers, mystery investigations", "ru": "Детективные истории, триллеры, загадочные расследования"}'::jsonb,
 '{"sr": "Kriminalistički romani | Vondi", "en": "Crime novels | Vondi", "ru": "Детективные романы | Vondi"}'::jsonb,
 '{"sr": "Kupite kriminalističke romane online", "en": "Buy crime novels online", "ru": "Купить детективные романы онлайн"}'::jsonb,
 '🕵️', true),

('ljubavni-romani', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/ljubavni-romani', 105,
 '{"sr": "Ljubavni romani", "en": "Romance novels", "ru": "Любовные романы"}'::jsonb,
 '{"sr": "Romantika, ljubavne priče, emotivni romani", "en": "Romance, love stories, emotional novels", "ru": "Романтика, любовные истории, эмоциональные романы"}'::jsonb,
 '{"sr": "Ljubavni romani | Vondi", "en": "Romance novels | Vondi", "ru": "Любовные романы | Vondi"}'::jsonb,
 '{"sr": "Kupite ljubavne romane online", "en": "Buy romance novels online", "ru": "Купить любовные романы онлайн"}'::jsonb,
 '💕', true),

('horor', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/horor', 106,
 '{"sr": "Horor", "en": "Horror", "ru": "Ужасы"}'::jsonb,
 '{"sr": "Horor romani, strašne priče, Stephen King", "en": "Horror novels, scary stories, Stephen King", "ru": "Ужастики, страшные истории, Стивен Кинг"}'::jsonb,
 '{"sr": "Horor | Vondi", "en": "Horror | Vondi", "ru": "Ужасы | Vondi"}'::jsonb,
 '{"sr": "Kupite horor knjige online", "en": "Buy horror books online", "ru": "Купить ужастики онлайн"}'::jsonb,
 '👻', true),

('poezija', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/poezija', 107,
 '{"sr": "Poezija", "en": "Poetry", "ru": "Поэзия"}'::jsonb,
 '{"sr": "Zbirke pesama, srpski i svetski pesnici", "en": "Poetry collections, Serbian and world poets", "ru": "Сборники стихов, сербские и мировые поэты"}'::jsonb,
 '{"sr": "Poezija | Vondi", "en": "Poetry | Vondi", "ru": "Поэзия | Vondi"}'::jsonb,
 '{"sr": "Kupite poeziju online - klasici i moderni", "en": "Buy poetry online", "ru": "Купить поэзию онлайн"}'::jsonb,
 '📜', true),

('eseji', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/eseji', 108,
 '{"sr": "Eseji", "en": "Essays", "ru": "Эссе"}'::jsonb,
 '{"sr": "Filozofski eseji, društvena tematika, eseistika", "en": "Philosophical essays, social topics, essay writing", "ru": "Философские эссе, социальная тематика, эссеистика"}'::jsonb,
 '{"sr": "Eseji | Vondi", "en": "Essays | Vondi", "ru": "Эссе | Vondi"}'::jsonb,
 '{"sr": "Kupite eseje online - filosofija i društvo", "en": "Buy essays online", "ru": "Купить эссе онлайн"}'::jsonb,
 '📝', true),

('knjige-za-roditelje', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' AND level = 1), 2, 'knjige-i-mediji/knjige-za-roditelje', 109,
 '{"sr": "Knjige za roditelje", "en": "Parenting books", "ru": "Книги для родителей"}'::jsonb,
 '{"sr": "Vodič za roditelje, vaspitanje dece, saveti", "en": "Parenting guide, child rearing, advice", "ru": "Руководство для родителей, воспитание детей, советы"}'::jsonb,
 '{"sr": "Knjige za roditelje | Vondi", "en": "Parenting books | Vondi", "ru": "Книги для родителей | Vondi"}'::jsonb,
 '{"sr": "Kupite knjige o roditeljstvu online", "en": "Buy parenting books online", "ru": "Купить книги о родительстве онлайн"}'::jsonb,
 '👶', true);

-- =============================================================================
-- L2 for: Ostalo (17 new categories)
-- =============================================================================

INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('antikviteti-namestaj', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/antikviteti-namestaj', 100,
 '{"sr": "Antikvitetni nameštaj", "en": "Antique furniture", "ru": "Антикварная мебель"}'::jsonb,
 '{"sr": "Stari nameštaj, barokni, klasični, retro komadi", "en": "Old furniture, baroque, classic, retro pieces", "ru": "Старинная мебель, барокко, классика, ретро предметы"}'::jsonb,
 '{"sr": "Antikvitetni nameštaj | Vondi", "en": "Antique furniture | Vondi", "ru": "Антикварная мебель | Vondi"}'::jsonb,
 '{"sr": "Kupite antikvitetni nameštaj online", "en": "Buy antique furniture online", "ru": "Купить антикварную мебель онлайн"}'::jsonb,
 '🪑', true),

('antikviteti-satovi', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/antikviteti-satovi', 101,
 '{"sr": "Antikvitetni satovi", "en": "Antique clocks", "ru": "Антикварные часы"}'::jsonb,
 '{"sr": "Stari zidni satovi, džepni satovi, kolekcionarski", "en": "Old wall clocks, pocket watches, collectibles", "ru": "Старинные настенные часы, карманные часы, коллекционные"}'::jsonb,
 '{"sr": "Antikvitetni satovi | Vondi", "en": "Antique clocks | Vondi", "ru": "Антикварные часы | Vondi"}'::jsonb,
 '{"sr": "Kupite antikvitetne satove online", "en": "Buy antique clocks online", "ru": "Купить антикварные часы онлайн"}'::jsonb,
 '🕰️', true),

('antikviteti-slike', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/antikviteti-slike', 102,
 '{"sr": "Antikvitetne slike", "en": "Antique paintings", "ru": "Антикварные картины"}'::jsonb,
 '{"sr": "Stare slike, umetničke slike, ulje na platnu", "en": "Old paintings, art paintings, oil on canvas", "ru": "Старинные картины, произведения искусства, масло на холсте"}'::jsonb,
 '{"sr": "Antikvitetne slike | Vondi", "en": "Antique paintings | Vondi", "ru": "Антикварные картины | Vondi"}'::jsonb,
 '{"sr": "Kupite antikvitetne slike online", "en": "Buy antique paintings online", "ru": "Купить антикварные картины онлайн"}'::jsonb,
 '🖼️', true),

('vintage-odeca-muska', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/vintage-odeca-muska', 103,
 '{"sr": "Vintage muška odeća", "en": "Vintage men''s clothing", "ru": "Винтажная мужская одежда"}'::jsonb,
 '{"sr": "Retro muška odeća, vintage brendovi, kolekcionarske", "en": "Retro men''s clothing, vintage brands, collectibles", "ru": "Ретро мужская одежда, винтажные бренды, коллекционная"}'::jsonb,
 '{"sr": "Vintage muška odeća | Vondi", "en": "Vintage men''s clothing | Vondi", "ru": "Винтажная мужская одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite vintage mušku odeću online", "en": "Buy vintage men''s clothing online", "ru": "Купить винтажную мужскую одежду онлайн"}'::jsonb,
 '👔', true),

('vintage-odeca-zenska', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/vintage-odeca-zenska', 104,
 '{"sr": "Vintage ženska odeća", "en": "Vintage women''s clothing", "ru": "Винтажная женская одежда"}'::jsonb,
 '{"sr": "Retro ženska odeća, vintage haljine, second-hand", "en": "Retro women''s clothing, vintage dresses, second-hand", "ru": "Ретро женская одежда, винтажные платья, секонд-хенд"}'::jsonb,
 '{"sr": "Vintage ženska odeća | Vondi", "en": "Vintage women''s clothing | Vondi", "ru": "Винтажная женская одежда | Vondi"}'::jsonb,
 '{"sr": "Kupite vintage žensku odeću online", "en": "Buy vintage women''s clothing online", "ru": "Купить винтажную женскую одежду онлайн"}'::jsonb,
 '👗', true),

('retro-elektronika', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/retro-elektronika', 105,
 '{"sr": "Retro elektronika", "en": "Retro electronics", "ru": "Ретро электроника"}'::jsonb,
 '{"sr": "Stari radio-aparati, gramofoni, kasete, retro konzole", "en": "Old radios, gramophones, cassettes, retro consoles", "ru": "Старые радиоприемники, граммофоны, кассеты, ретро консоли"}'::jsonb,
 '{"sr": "Retro elektronika | Vondi", "en": "Retro electronics | Vondi", "ru": "Ретро электроника | Vondi"}'::jsonb,
 '{"sr": "Kupite retro elektroniku online", "en": "Buy retro electronics online", "ru": "Купить ретро электронику онлайн"}'::jsonb,
 '📻', true),

('kolekcionarske-karte', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/kolekcionarske-karte', 106,
 '{"sr": "Kolekcionarske karte", "en": "Trading cards", "ru": "Коллекционные карточки"}'::jsonb,
 '{"sr": "Pokemon, Magic, Yu-Gi-Oh, sportske karte, retke", "en": "Pokemon, Magic, Yu-Gi-Oh, sports cards, rare", "ru": "Покемон, Magic, Yu-Gi-Oh, спортивные карточки, редкие"}'::jsonb,
 '{"sr": "Kolekcionarske karte | Vondi", "en": "Trading cards | Vondi", "ru": "Коллекционные карточки | Vondi"}'::jsonb,
 '{"sr": "Kupite kolekcionarske karte online", "en": "Buy trading cards online", "ru": "Купить коллекционные карточки онлайн"}'::jsonb,
 '🃏', true),

('kolekcionarski-novcici', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/kolekcionarski-novcici', 107,
 '{"sr": "Kolekcionarski novčići", "en": "Collectible coins", "ru": "Коллекционные монеты"}'::jsonb,
 '{"sr": "Stari novčići, retki, jugoslovenski, strani", "en": "Old coins, rare, Yugoslav, foreign", "ru": "Старые монеты, редкие, югославские, иностранные"}'::jsonb,
 '{"sr": "Kolekcionarski novčići | Vondi", "en": "Collectible coins | Vondi", "ru": "Коллекционные монеты | Vondi"}'::jsonb,
 '{"sr": "Kupite kolekcionarske novčiće online", "en": "Buy collectible coins online", "ru": "Купить коллекционные монеты онлайн"}'::jsonb,
 '🪙', true),

('markice-filatelija', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/markice-filatelija', 108,
 '{"sr": "Markice - filatelija", "en": "Stamps - philately", "ru": "Марки - филателия"}'::jsonb,
 '{"sr": "Poštanske marke, stare, retke, jugoslovenske", "en": "Postage stamps, old, rare, Yugoslav", "ru": "Почтовые марки, старые, редкие, югославские"}'::jsonb,
 '{"sr": "Markice - filatelija | Vondi", "en": "Stamps - philately | Vondi", "ru": "Марки - филателия | Vondi"}'::jsonb,
 '{"sr": "Kupite poštanske marke online - filatelija", "en": "Buy stamps online", "ru": "Купить почтовые марки онлайн"}'::jsonb,
 '💌', true),

('minerali-i-kamenje', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/minerali-i-kamenje', 109,
 '{"sr": "Minerali i kamenje", "en": "Minerals and stones", "ru": "Минералы и камни"}'::jsonb,
 '{"sr": "Poludrago kamenje, kristali, ametist, kvarc", "en": "Semi-precious stones, crystals, amethyst, quartz", "ru": "Полудрагоценные камни, кристаллы, аметист, кварц"}'::jsonb,
 '{"sr": "Minerali i kamenje | Vondi", "en": "Minerals and stones | Vondi", "ru": "Минералы и камни | Vondi"}'::jsonb,
 '{"sr": "Kupite minerale i kamenje online", "en": "Buy minerals and stones online", "ru": "Купить минералы и камни онлайн"}'::jsonb,
 '💎', true),

('fosili', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/fosili', 110,
 '{"sr": "Fosili", "en": "Fossils", "ru": "Окаменелости"}'::jsonb,
 '{"sr": "Fosilizovani ostaci, zubi ajkule, amoniti", "en": "Fossilized remains, shark teeth, ammonites", "ru": "Окаменелости, зубы акулы, аммониты"}'::jsonb,
 '{"sr": "Fosili | Vondi", "en": "Fossils | Vondi", "ru": "Окаменелости | Vondi"}'::jsonb,
 '{"sr": "Kupite fosile online - kolekcionarski", "en": "Buy fossils online", "ru": "Купить окаменелости онлайн"}'::jsonb,
 '🦴', true),

('rucni-rad', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/rucni-rad', 111,
 '{"sr": "Ručni rad", "en": "Handmade", "ru": "Ручная работа"}'::jsonb,
 '{"sr": "Ručno rađeni proizvodi, unikatni pokloni", "en": "Handmade products, unique gifts", "ru": "Ручные изделия, уникальные подарки"}'::jsonb,
 '{"sr": "Ručni rad | Vondi", "en": "Handmade | Vondi", "ru": "Ручная работа | Vondi"}'::jsonb,
 '{"sr": "Kupite ručno rađene proizvode online", "en": "Buy handmade products online", "ru": "Купить изделия ручной работы онлайн"}'::jsonb,
 '🤲', true),

('custom-proizvodi', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/custom-proizvodi', 112,
 '{"sr": "Custom proizvodi", "en": "Custom products", "ru": "Кастомные изделия"}'::jsonb,
 '{"sr": "Proizvodi po meri, personalizovani pokloni", "en": "Made-to-order products, personalized gifts", "ru": "Изделия на заказ, персонализированные подарки"}'::jsonb,
 '{"sr": "Custom proizvodi | Vondi", "en": "Custom products | Vondi", "ru": "Кастомные изделия | Vondi"}'::jsonb,
 '{"sr": "Naručite custom proizvode po meri", "en": "Order custom products online", "ru": "Заказать кастомные изделия онлайн"}'::jsonb,
 '🎁', true),

('personalizovani-pokloni', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/personalizovani-pokloni', 113,
 '{"sr": "Personalizovani pokloni", "en": "Personalized gifts", "ru": "Персонализированные подарки"}'::jsonb,
 '{"sr": "Pokloni sa gravurom, fotografijom, imenom", "en": "Gifts with engraving, photo, name", "ru": "Подарки с гравировкой, фотографией, именем"}'::jsonb,
 '{"sr": "Personalizovani pokloni | Vondi", "en": "Personalized gifts | Vondi", "ru": "Персонализированные подарки | Vondi"}'::jsonb,
 '{"sr": "Kupite personalizovane poklone online", "en": "Buy personalized gifts online", "ru": "Купить персонализированные подарки онлайн"}'::jsonb,
 '💝', true),

('korporativni-pokloni-misc', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/korporativni-pokloni-misc', 114,
 '{"sr": "Korporativni pokloni", "en": "Corporate gifts", "ru": "Корпоративные подарки"}'::jsonb,
 '{"sr": "Poslovno poklanjanje, brendirani proizvodi", "en": "Business gifting, branded products", "ru": "Деловые подарки, брендированные товары"}'::jsonb,
 '{"sr": "Korporativni pokloni | Vondi", "en": "Corporate gifts | Vondi", "ru": "Корпоративные подарки | Vondi"}'::jsonb,
 '{"sr": "Kupite korporativne poklone online", "en": "Buy corporate gifts online", "ru": "Купить корпоративные подарки онлайн"}'::jsonb,
 '🎀', true),

('razno', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/razno', 115,
 '{"sr": "Razno", "en": "Miscellaneous", "ru": "Разное"}'::jsonb,
 '{"sr": "Ostali proizvodi koji ne spadaju ni u jednu kategoriju", "en": "Other products that don''t fit any category", "ru": "Прочие товары, не относящиеся ни к одной категории"}'::jsonb,
 '{"sr": "Razno | Vondi", "en": "Miscellaneous | Vondi", "ru": "Разное | Vondi"}'::jsonb,
 '{"sr": "Razni proizvodi - sve ostalo", "en": "Miscellaneous products", "ru": "Разные товары"}'::jsonb,
 '📦', true),

('nerazvrstavano', (SELECT id FROM categories WHERE slug = 'ostalo' AND level = 1), 2, 'ostalo/nerazvrstavano', 116,
 '{"sr": "Nerazvrstavano", "en": "Uncategorized", "ru": "Без категории"}'::jsonb,
 '{"sr": "Proizvodi koji čekaju kategorizaciju", "en": "Products awaiting categorization", "ru": "Товары, ожидающие категоризации"}'::jsonb,
 '{"sr": "Nerazvrstavano | Vondi", "en": "Uncategorized | Vondi", "ru": "Без категории | Vondi"}'::jsonb,
 '{"sr": "Nerazvrstavani proizvodi", "en": "Uncategorized products", "ru": "Товары без категории"}'::jsonb,
 '❓', true);
