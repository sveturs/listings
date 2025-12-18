-- Migration: Seed L2 categories (Part 2 of 2)
-- Date: 2025-12-16
-- Purpose: Continue inserting L2 subcategories for remaining 15 L1 categories
-- Previous: 20251216000005_seed_categories_l2.up.sql (first 3 L1 categories)

-- =============================================================================
-- L2 for: 4. Lepota i zdravlje (Beauty & Health) - 12 categories
-- =============================================================================
INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('nega-koze', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/nega-koze', 1,
 '{"sr": "Nega kože", "en": "Skincare", "ru": "Уход за кожей"}'::jsonb,
 '{"sr": "Kreme, serumi, maske za lice", "en": "Creams, serums, facial masks", "ru": "Кремы, сыворотки, маски для лица"}'::jsonb,
 '{"sr": "Nega kože | Vondi", "en": "Skincare | Vondi", "ru": "Уход за кожей | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za negu kože online", "en": "Buy skincare products online", "ru": "Купить средства для ухода за кожей онлайн"}'::jsonb,
 '🧴', true),

('nega-kose', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/nega-kose', 2,
 '{"sr": "Nega kose", "en": "Hair care", "ru": "Уход за волосами"}'::jsonb,
 '{"sr": "Šamponi, balzami, maske za kosu", "en": "Shampoos, conditioners, hair masks", "ru": "Шампуни, бальзамы, маски для волос"}'::jsonb,
 '{"sr": "Nega kose | Vondi", "en": "Hair care | Vondi", "ru": "Уход за волосами | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za negu kose online", "en": "Buy hair care products online", "ru": "Купить средства для ухода за волосами онлайн"}'::jsonb,
 '💇', true),

('parfemi', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/parfemi', 3,
 '{"sr": "Parfemi", "en": "Perfumes", "ru": "Парфюмерия"}'::jsonb,
 '{"sr": "Parfemi za žene i muškarce", "en": "Perfumes for women and men", "ru": "Духи для женщин и мужчин"}'::jsonb,
 '{"sr": "Parfemi | Vondi", "en": "Perfumes | Vondi", "ru": "Парфюмерия | Vondi"}'::jsonb,
 '{"sr": "Kupite parfeme online - originalni brendovi", "en": "Buy perfumes online - original brands", "ru": "Купить парфюмерию онлайн - оригинальные бренды"}'::jsonb,
 '🌸', true),

('dekorativna-kozmetika', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/dekorativna-kozmetika', 4,
 '{"sr": "Dekorativna kozmetika", "en": "Makeup", "ru": "Декоративная косметика"}'::jsonb,
 '{"sr": "Šminka, ruž, senke, maskara", "en": "Makeup, lipstick, eyeshadow, mascara", "ru": "Макияж, помада, тени, тушь"}'::jsonb,
 '{"sr": "Dekorativna kozmetika | Vondi", "en": "Makeup | Vondi", "ru": "Декоративная косметика | Vondi"}'::jsonb,
 '{"sr": "Kupite dekorativnu kozmetiku online", "en": "Buy makeup online", "ru": "Купить декоративную косметику онлайн"}'::jsonb,
 '💄', true),

('manikir-i-pedikir', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/manikir-i-pedikir', 5,
 '{"sr": "Manikir i pedikir", "en": "Manicure & pedicure", "ru": "Маникюр и педикюр"}'::jsonb,
 '{"sr": "Lakovi za nokte, gel lakovi, pribor", "en": "Nail polish, gel polish, tools", "ru": "Лак для ногтей, гель-лак, инструменты"}'::jsonb,
 '{"sr": "Manikir i pedikir | Vondi", "en": "Manicure & pedicure | Vondi", "ru": "Маникюр и педикюр | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za manikir online", "en": "Buy manicure products online", "ru": "Купить продукцию для маникюра онлайн"}'::jsonb,
 '💅', true),

('muska-nega', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/muska-nega', 6,
 '{"sr": "Muška nega", "en": "Men''s grooming", "ru": "Мужской уход"}'::jsonb,
 '{"sr": "Aparati za brijanje, pene, losioni", "en": "Shavers, shaving foam, lotions", "ru": "Бритвы, пена для бритья, лосьоны"}'::jsonb,
 '{"sr": "Muška nega | Vondi", "en": "Men''s grooming | Vondi", "ru": "Мужской уход | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za mušku negu online", "en": "Buy men''s grooming products online", "ru": "Купить средства мужского ухода онлайн"}'::jsonb,
 '🧔', true),

('vitamini-i-suplementi', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/vitamini-i-suplementi', 7,
 '{"sr": "Vitamini i suplementi", "en": "Vitamins & supplements", "ru": "Витамины и добавки"}'::jsonb,
 '{"sr": "Vitamini, minerali, proteini", "en": "Vitamins, minerals, proteins", "ru": "Витамины, минералы, протеины"}'::jsonb,
 '{"sr": "Vitamini i suplementi | Vondi", "en": "Vitamins & supplements | Vondi", "ru": "Витамины и добавки | Vondi"}'::jsonb,
 '{"sr": "Kupite vitamine i suplemente online", "en": "Buy vitamins and supplements online", "ru": "Купить витамины и добавки онлайн"}'::jsonb,
 '💊', true),

('medicinski-proizvodi', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/medicinski-proizvodi', 8,
 '{"sr": "Medicinski proizvodi", "en": "Medical products", "ru": "Медицинские товары"}'::jsonb,
 '{"sr": "Termometri, tonometri, prve pomoći", "en": "Thermometers, tonometers, first aid", "ru": "Термометры, тонометры, первая помощь"}'::jsonb,
 '{"sr": "Medicinski proizvodi | Vondi", "en": "Medical products | Vondi", "ru": "Медицинские товары | Vondi"}'::jsonb,
 '{"sr": "Kupite medicinske proizvode online", "en": "Buy medical products online", "ru": "Купить медицинские товары онлайн"}'::jsonb,
 '🩺', true),

('eterična-ulja', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/eterična-ulja', 9,
 '{"sr": "Eterična ulja", "en": "Essential oils", "ru": "Эфирные масла"}'::jsonb,
 '{"sr": "Prirodna eterična ulja i aromater apija", "en": "Natural essential oils and aromatherapy", "ru": "Натуральные эфирные масла и ароматерапия"}'::jsonb,
 '{"sr": "Eterična ulja | Vondi", "en": "Essential oils | Vondi", "ru": "Эфирные масла | Vondi"}'::jsonb,
 '{"sr": "Kupite eterična ulja online", "en": "Buy essential oils online", "ru": "Купить эфирные масла онлайн"}'::jsonb,
 '🌿', true),

('spa-i-relaksacija', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/spa-i-relaksacija', 10,
 '{"sr": "Spa i relaksacija", "en": "Spa & relaxation", "ru": "Спа и релаксация"}'::jsonb,
 '{"sr": "Masažeri, difuzeri, spa proizvodi", "en": "Massagers, diffusers, spa products", "ru": "Массажеры, диффузоры, спа продукты"}'::jsonb,
 '{"sr": "Spa i relaksacija | Vondi", "en": "Spa & relaxation | Vondi", "ru": "Спа и релаксация | Vondi"}'::jsonb,
 '{"sr": "Kupite spa proizvode online", "en": "Buy spa products online", "ru": "Купить спа продукты онлайн"}'::jsonb,
 '🧖', true),

('higijena', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/higijena', 11,
 '{"sr": "Higijena", "en": "Hygiene", "ru": "Гигиена"}'::jsonb,
 '{"sr": "Sapuni, gelovi za tuširanje, dezodoransi", "en": "Soaps, shower gels, deodorants", "ru": "Мыло, гели для душа, дезодоранты"}'::jsonb,
 '{"sr": "Higijena | Vondi", "en": "Hygiene | Vondi", "ru": "Гигиена | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za higijenu online", "en": "Buy hygiene products online", "ru": "Купить средства гигиены онлайн"}'::jsonb,
 '🧼', true),

('oralna-higijena', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje'), 2, 'lepota-i-zdravlje/oralna-higijena', 12,
 '{"sr": "Oralna higijena", "en": "Oral hygiene", "ru": "Гигиена полости рта"}'::jsonb,
 '{"sr": "Paste za zube, četkice, vodice", "en": "Toothpaste, toothbrushes, mouthwash", "ru": "Зубные пасты, щетки, ополаскиватели"}'::jsonb,
 '{"sr": "Oralna higijena | Vondi", "en": "Oral hygiene | Vondi", "ru": "Гигиена полости рта | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za oralnu higijenu online", "en": "Buy oral hygiene products online", "ru": "Купить средства для гигиены полости рта онлайн"}'::jsonb,
 '🦷', true),

-- =============================================================================
-- L2 for: 5. Za bebe i decu (Baby & Kids) - 12 categories
-- =============================================================================

('oprema-za-bebe', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/oprema-za-bebe', 1,
 '{"sr": "Oprema za bebe", "en": "Baby gear", "ru": "Оборудование для новорожденных"}'::jsonb,
 '{"sr": "Kolica, autosedišta, nosiljke", "en": "Strollers, car seats, carriers", "ru": "Коляски, автокресла, переноски"}'::jsonb,
 '{"sr": "Oprema za bebe | Vondi", "en": "Baby gear | Vondi", "ru": "Оборудование для новорожденных | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za bebe online", "en": "Buy baby gear online", "ru": "Купить оборудование для новорожденных онлайн"}'::jsonb,
 '🍼', true),

('namestaj-za-bebe', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/namestaj-za-bebe', 2,
 '{"sr": "Nameštaj za bebe", "en": "Baby furniture", "ru": "Мебель для новорожденных"}'::jsonb,
 '{"sr": "Krevetci, komoda, stolice za hranjenje", "en": "Cribs, dressers, high chairs", "ru": "Кроватки, комоды, стульчики для кормления"}'::jsonb,
 '{"sr": "Nameštaj za bebe | Vondi", "en": "Baby furniture | Vondi", "ru": "Мебель для новорожденных | Vondi"}'::jsonb,
 '{"sr": "Kupite nameštaj za bebe online", "en": "Buy baby furniture online", "ru": "Купить мебель для новорожденных онлайн"}'::jsonb,
 '🛏️', true),

('nega-i-higijena-beba', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/nega-i-higijena-beba', 3,
 '{"sr": "Nega i higijena beba", "en": "Baby care & hygiene", "ru": "Уход и гигиена младенцев"}'::jsonb,
 '{"sr": "Pelene, vlažne maramice, kreme", "en": "Diapers, wet wipes, creams", "ru": "Подгузники, влажные салфетки, кремы"}'::jsonb,
 '{"sr": "Nega i higijena beba | Vondi", "en": "Baby care & hygiene | Vondi", "ru": "Уход и гигиена младенцев | Vondi"}'::jsonb,
 '{"sr": "Kupite proizvode za negu beba online", "en": "Buy baby care products online", "ru": "Купить средства ухода за младенцами онлайн"}'::jsonb,
 '👶', true),

('hrana-za-bebe', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/hrana-za-bebe', 4,
 '{"sr": "Hrana za bebe", "en": "Baby food", "ru": "Детское питание"}'::jsonb,
 '{"sr": "Mleko, kaše, sokovi i kašice", "en": "Formula, cereals, juices and purees", "ru": "Смеси, каши, соки и пюре"}'::jsonb,
 '{"sr": "Hrana za bebe | Vondi", "en": "Baby food | Vondi", "ru": "Детское питание | Vondi"}'::jsonb,
 '{"sr": "Kupite hranu za bebe online", "en": "Buy baby food online", "ru": "Купить детское питание онлайн"}'::jsonb,
 '🍼', true),

('igracke-za-bebe', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/igracke-za-bebe', 5,
 '{"sr": "Igračke za bebe", "en": "Baby toys", "ru": "Игрушки для новорожденных"}'::jsonb,
 '{"sr": "Zvečke, plišane igračke, razvojne igre", "en": "Rattles, plush toys, developmental games", "ru": "Погремушки, мягкие игрушки, развивающие игры"}'::jsonb,
 '{"sr": "Igračke za bebe | Vondi", "en": "Baby toys | Vondi", "ru": "Игрушки для новорожденных | Vondi"}'::jsonb,
 '{"sr": "Kupite igračke za bebe online", "en": "Buy baby toys online", "ru": "Купить игрушки для новорожденных онлайн"}'::jsonb,
 '🧸', true),

('igracke-za-decu', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/igracke-za-decu', 6,
 '{"sr": "Igračke za decu", "en": "Kids'' toys", "ru": "Детские игрушки"}'::jsonb,
 '{"sr": "Lutke, autići, konstruktori, edukativne igre", "en": "Dolls, cars, constructors, educational games", "ru": "Куклы, машинки, конструкторы, обучающие игры"}'::jsonb,
 '{"sr": "Igračke za decu | Vondi", "en": "Kids'' toys | Vondi", "ru": "Детские игрушки | Vondi"}'::jsonb,
 '{"sr": "Kupite igračke za decu online", "en": "Buy kids'' toys online", "ru": "Купить детские игрушки онлайн"}'::jsonb,
 '🎲', true),

('decija-odeca-bebe', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/decija-odeca-bebe', 7,
 '{"sr": "Dečija odeća i bebe", "en": "Kids'' & baby clothing", "ru": "Детская одежда и для новорожденных"}'::jsonb,
 '{"sr": "Bodići, pidžamice, haljine, pantalone", "en": "Bodysuits, pajamas, dresses, pants", "ru": "Боди, пижамы, платья, брюки"}'::jsonb,
 '{"sr": "Dečija odeća i bebe | Vondi", "en": "Kids'' & baby clothing | Vondi", "ru": "Детская одежда и для новорожденных | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiju odeću i za bebe online", "en": "Buy kids'' & baby clothing online", "ru": "Купить детскую одежду и для новорожденных онлайн"}'::jsonb,
 '👕', true),

('decija-obuca-bebe', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/decija-obuca-bebe', 8,
 '{"sr": "Dečija obuća i bebe", "en": "Kids'' & baby footwear", "ru": "Детская обувь и для новорожденных"}'::jsonb,
 '{"sr": "Patike, sandale, čižmice za bebe", "en": "Sneakers, sandals, baby boots", "ru": "Кроссовки, сандалии, ботиночки для младенцев"}'::jsonb,
 '{"sr": "Dečija obuća i bebe | Vondi", "en": "Kids'' & baby footwear | Vondi", "ru": "Детская обувь и для новорожденных | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiju obuću i za bebe online", "en": "Buy kids'' & baby footwear online", "ru": "Купить детскую обувь и для новорожденных онлайн"}'::jsonb,
 '👟', true),

('skolski-pribor', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/skolski-pribor', 9,
 '{"sr": "Školski pribor", "en": "School supplies", "ru": "Школьные принадлежности"}'::jsonb,
 '{"sr": "Ranci, torbe, sveske, olovke", "en": "Backpacks, bags, notebooks, pens", "ru": "Рюкзаки, сумки, тетради, ручки"}'::jsonb,
 '{"sr": "Školski pribor | Vondi", "en": "School supplies | Vondi", "ru": "Школьные принадлежности | Vondi"}'::jsonb,
 '{"sr": "Kupite školski pribor online", "en": "Buy school supplies online", "ru": "Купить школьные принадлежности онлайн"}'::jsonb,
 '🎒', true),

('deciji-namestaj', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/deciji-namestaj', 10,
 '{"sr": "Dečiji nameštaj", "en": "Kids'' furniture", "ru": "Детская мебель"}'::jsonb,
 '{"sr": "Kreveti, stolovi, stolice, police", "en": "Beds, tables, chairs, shelves", "ru": "Кровати, столы, стулья, полки"}'::jsonb,
 '{"sr": "Dečiji nameštaj | Vondi", "en": "Kids'' furniture | Vondi", "ru": "Детская мебель | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiji nameštaj online", "en": "Buy kids'' furniture online", "ru": "Купить детскую мебель онлайн"}'::jsonb,
 '🪑', true),

('decija-kozmetika', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/decija-kozmetika', 11,
 '{"sr": "Dečija kozmetika", "en": "Kids'' cosmetics", "ru": "Детская косметика"}'::jsonb,
 '{"sr": "Šamponi, paste, kreme za decu", "en": "Shampoos, toothpaste, creams for kids", "ru": "Шампуни, пасты, кремы для детей"}'::jsonb,
 '{"sr": "Dečija kozmetika | Vondi", "en": "Kids'' cosmetics | Vondi", "ru": "Детская косметика | Vondi"}'::jsonb,
 '{"sr": "Kupite dečiju kozmetiku online", "en": "Buy kids'' cosmetics online", "ru": "Купить детскую косметику онлайн"}'::jsonb,
 '🧴', true),

('elektronika-za-decu', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu'), 2, 'za-bebe-i-decu/elektronika-za-decu', 12,
 '{"sr": "Elektronika za decu", "en": "Kids'' electronics", "ru": "Электроника для детей"}'::jsonb,
 '{"sr": "Tableti za decu, igračke, satovi", "en": "Kids'' tablets, toys, watches", "ru": "Детские планшеты, игрушки, часы"}'::jsonb,
 '{"sr": "Elektronika za decu | Vondi", "en": "Kids'' electronics | Vondi", "ru": "Электроника для детей | Vondi"}'::jsonb,
 '{"sr": "Kupite elektroniku za decu online", "en": "Buy kids'' electronics online", "ru": "Купить электронику для детей онлайн"}'::jsonb,
 '📱', true),

-- =============================================================================
-- L2 for: 6. Sport i turizam (Sports & Outdoors) - 12 categories
-- =============================================================================

('fitnes-i-teretana', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/fitnes-i-teretana', 1,
 '{"sr": "Fitnes i teretana", "en": "Fitness & gym", "ru": "Фитнес и тренажерный зал"}'::jsonb,
 '{"sr": "Tegovi, sprave, podloge za vežbanje", "en": "Dumbbells, equipment, exercise mats", "ru": "Гантели, оборудование, коврики для упражнений"}'::jsonb,
 '{"sr": "Fitnes i teretana | Vondi", "en": "Fitness & gym | Vondi", "ru": "Фитнес и тренажерный зал | Vondi"}'::jsonb,
 '{"sr": "Kupite fitnes opremu online", "en": "Buy fitness equipment online", "ru": "Купить фитнес оборудование онлайн"}'::jsonb,
 '🏋️', true),

('bicikli-i-trotineti', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/bicikli-i-trotineti', 2,
 '{"sr": "Bicikli i trotineti", "en": "Bicycles & scooters", "ru": "Велосипеды и самокаты"}'::jsonb,
 '{"sr": "Bicikli, električni bicikli, trotineti", "en": "Bicycles, e-bikes, scooters", "ru": "Велосипеды, электровелосипеды, самокаты"}'::jsonb,
 '{"sr": "Bicikli i trotineti | Vondi", "en": "Bicycles & scooters | Vondi", "ru": "Велосипеды и самокаты | Vondi"}'::jsonb,
 '{"sr": "Kupite bicikle i trotinete online", "en": "Buy bicycles and scooters online", "ru": "Купить велосипеды и самокаты онлайн"}'::jsonb,
 '🚴', true),

('kampovanje', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/kampovanje', 3,
 '{"sr": "Kampovanje", "en": "Camping", "ru": "Кемпинг"}'::jsonb,
 '{"sr": "Šatori, vreće za spavanje, oprema", "en": "Tents, sleeping bags, equipment", "ru": "Палатки, спальные мешки, снаряжение"}'::jsonb,
 '{"sr": "Kampovanje | Vondi", "en": "Camping | Vondi", "ru": "Кемпинг | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za kampovanje online", "en": "Buy camping equipment online", "ru": "Купить снаряжение для кемпинга онлайн"}'::jsonb,
 '⛺', true),

('fudbal', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/fudbal', 4,
 '{"sr": "Fudbal", "en": "Football", "ru": "Футбол"}'::jsonb,
 '{"sr": "Lopte, kopačke, dresovi, golovi", "en": "Balls, cleats, jerseys, goals", "ru": "Мячи, бутсы, майки, ворота"}'::jsonb,
 '{"sr": "Fudbal | Vondi", "en": "Football | Vondi", "ru": "Футбол | Vondi"}'::jsonb,
 '{"sr": "Kupite fudbalsku opremu online", "en": "Buy football equipment online", "ru": "Купить футбольное снаряжение онлайн"}'::jsonb,
 '⚽', true),

('kosarka', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/kosarka', 5,
 '{"sr": "Košarka", "en": "Basketball", "ru": "Баскетбол"}'::jsonb,
 '{"sr": "Lopte, koševi, patike, dresovi", "en": "Balls, hoops, sneakers, jerseys", "ru": "Мячи, кольца, кроссовки, майки"}'::jsonb,
 '{"sr": "Košarka | Vondi", "en": "Basketball | Vondi", "ru": "Баскетбол | Vondi"}'::jsonb,
 '{"sr": "Kupite košarkašku opremu online", "en": "Buy basketball equipment online", "ru": "Купить баскетбольное снаряжение онлайн"}'::jsonb,
 '🏀', true),

('tenis', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/tenis', 6,
 '{"sr": "Tenis", "en": "Tennis", "ru": "Теннис"}'::jsonb,
 '{"sr": "Reketi, loptice, patike, mreže", "en": "Rackets, balls, shoes, nets", "ru": "Ракетки, мячи, обувь, сетки"}'::jsonb,
 '{"sr": "Tenis | Vondi", "en": "Tennis | Vondi", "ru": "Теннис | Vondi"}'::jsonb,
 '{"sr": "Kupite tenisku opremu online", "en": "Buy tennis equipment online", "ru": "Купить теннисное снаряжение онлайн"}'::jsonb,
 '🎾', true),

('plivanje', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/plivanje', 7,
 '{"sr": "Plivanje", "en": "Swimming", "ru": "Плавание"}'::jsonb,
 '{"sr": "Kupaći, naočare, kape, daske", "en": "Swimwear, goggles, caps, boards", "ru": "Купальники, очки, шапочки, доски"}'::jsonb,
 '{"sr": "Plivanje | Vondi", "en": "Swimming | Vondi", "ru": "Плавание | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za plivanje online", "en": "Buy swimming equipment online", "ru": "Купить снаряжение для плавания онлайн"}'::jsonb,
 '🏊', true),

('planinarenje', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/planinarenje', 8,
 '{"sr": "Planinarenje", "en": "Hiking", "ru": "Пеший туризм"}'::jsonb,
 '{"sr": "Ranci, cipele, štapovi, oprema", "en": "Backpacks, boots, poles, equipment", "ru": "Рюкзаки, ботинки, палки, снаряжение"}'::jsonb,
 '{"sr": "Planinarenje | Vondi", "en": "Hiking | Vondi", "ru": "Пеший туризм | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za planinarenje online", "en": "Buy hiking equipment online", "ru": "Купить снаряжение для пешего туризма онлайн"}'::jsonb,
 '🥾', true),

('zimski-sportovi', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/zimski-sportovi', 9,
 '{"sr": "Zimski sportovi", "en": "Winter sports", "ru": "Зимние виды спорта"}'::jsonb,
 '{"sr": "Skije, snowboard, klizaljke", "en": "Skis, snowboards, ice skates", "ru": "Лыжи, сноуборды, коньки"}'::jsonb,
 '{"sr": "Zimski sportovi | Vondi", "en": "Winter sports | Vondi", "ru": "Зимние виды спорта | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za zimske sportove online", "en": "Buy winter sports equipment online", "ru": "Купить снаряжение для зимних видов спорта онлайн"}'::jsonb,
 '⛷️', true),

('ribolov', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/ribolov', 10,
 '{"sr": "Ribolov", "en": "Fishing", "ru": "Рыбалка"}'::jsonb,
 '{"sr": "Štapovi, mašine, mamci, oprema", "en": "Rods, reels, lures, equipment", "ru": "Удочки, катушки, приманки, снаряжение"}'::jsonb,
 '{"sr": "Ribolov | Vondi", "en": "Fishing | Vondi", "ru": "Рыбалка | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za ribolov online", "en": "Buy fishing equipment online", "ru": "Купить снаряжение для рыбалки онлайн"}'::jsonb,
 '🎣', true),

('lov', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/lov', 11,
 '{"sr": "Lov", "en": "Hunting", "ru": "Охота"}'::jsonb,
 '{"sr": "Vazdušni pištolji, oprema, odeća", "en": "Air guns, equipment, clothing", "ru": "Пневматическое оружие, снаряжение, одежда"}'::jsonb,
 '{"sr": "Lov | Vondi", "en": "Hunting | Vondi", "ru": "Охота | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za lov online", "en": "Buy hunting equipment online", "ru": "Купить снаряжение для охоты онлайн"}'::jsonb,
 '🏹', true),

('dzonovanje', (SELECT id FROM categories WHERE slug = 'sport-i-turizam'), 2, 'sport-i-turizam/dzonovanje', 12,
 '{"sr": "Džonovanje", "en": "Jogging & running", "ru": "Бег и джоггинг"}'::jsonb,
 '{"sr": "Patike za trčanje, trenerke, dodaci", "en": "Running shoes, tracksuits, accessories", "ru": "Беговая обувь, костюмы, аксессуары"}'::jsonb,
 '{"sr": "Džonovanje | Vondi", "en": "Jogging & running | Vondi", "ru": "Бег и джоггинг | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za džonovanje online", "en": "Buy jogging equipment online", "ru": "Купить снаряжение для бега онлайн"}'::jsonb,
 '🏃', true);

-- Continue with remaining categories (Automobilizam, Kućni aparati, etc.) in separate file due to length
-- Current progress: 81 L2 categories (15 + 15 + 15 + 12 + 12 + 12)

-- Temporary verification
DO $$
DECLARE
    l2_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO l2_count FROM categories WHERE level = 2;

    RAISE NOTICE 'Part 2 progress: % total L2 categories (Lepota: 12, Bebe: 12, Sport: 12)', l2_count;
END $$;
