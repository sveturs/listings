-- Migration: Expand L2 categories (Part 6)
-- Date: 2025-12-17
-- Purpose: Add 80 L2 categories - Empty/minimal L1 categories expansion
-- Expanding: Kancelarija, Muzički instrumenti, Hrana, Igračke, Umetnost

-- =============================================================================
-- Kancelarijski materijal: +15 L2 (beyond existing)
-- =============================================================================

INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('kancelarijska-hartija', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/kancelarijska-hartija', 100,
 '{"sr": "Kancelarijska hartija", "en": "Office paper", "ru": "Канцелярская бумага"}'::jsonb,
 '{"sr": "A4 papir, fotokopir papir, color papir", "en": "A4 paper, photocopy paper, color paper", "ru": "Бумага A4, бумага для копирования, цветная бумага"}'::jsonb,
 '{"sr": "Kancelarijska hartija | Vondi", "en": "Office paper | Vondi", "ru": "Канцелярская бумага | Vondi"}'::jsonb,
 '{"sr": "Kupite kancelarijsku hartiju online", "en": "Buy office paper online", "ru": "Купить канцелярскую бумагу онлайн"}'::jsonb,
 '📄', true),

('stampaci-toneri', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/stampaci-toneri', 101,
 '{"sr": "Štampači i toneri", "en": "Printers & toners", "ru": "Принтеры и тонеры"}'::jsonb,
 '{"sr": "Laserski, inkjet, multifunkcionalni, toner cartridge", "en": "Laser, inkjet, multifunction, toner cartridges", "ru": "Лазерные, струйные, МФУ, тонер-картриджи"}'::jsonb,
 '{"sr": "Štampači i toneri | Vondi", "en": "Printers & toners | Vondi", "ru": "Принтеры и тонеры | Vondi"}'::jsonb,
 '{"sr": "Kupite štampače i tonere online", "en": "Buy printers and toners online", "ru": "Купить принтеры и тонеры онлайн"}'::jsonb,
 '🖨️', true),

('registratori-fascikle', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/registratori-fascikle', 102,
 '{"sr": "Registratori i fascikle", "en": "Binders & folders", "ru": "Папки и скоросшиватели"}'::jsonb,
 '{"sr": "Ring binder, lever arch, plastic folders", "en": "Ring binders, lever arch, plastic folders", "ru": "Кольцевые папки, папки-регистраторы, пластиковые папки"}'::jsonb,
 '{"sr": "Registratori i fascikle | Vondi", "en": "Binders & folders | Vondi", "ru": "Папки и скоросшиватели | Vondi"}'::jsonb,
 '{"sr": "Kupite registratore i fascikle online", "en": "Buy binders and folders online", "ru": "Купить папки и скоросшиватели онлайн"}'::jsonb,
 '📁', true),

('olovke-i-markeri', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/olovke-i-markeri', 103,
 '{"sr": "Olovke i markeri", "en": "Pens & markers", "ru": "Ручки и маркеры"}'::jsonb,
 '{"sr": "Hemijske, roler, gel, permanentni, highlighter", "en": "Ballpoint, rollerball, gel, permanent, highlighter", "ru": "Шариковые, роллерболы, гелевые, перманентные, текстовыделители"}'::jsonb,
 '{"sr": "Olovke i markeri | Vondi", "en": "Pens & markers | Vondi", "ru": "Ручки и маркеры | Vondi"}'::jsonb,
 '{"sr": "Kupite olovke i markere online", "en": "Buy pens and markers online", "ru": "Купить ручки и маркеры онлайн"}'::jsonb,
 '🖊️', true),

('beležnice-blokovi', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/beleznice-blokovi', 104,
 '{"sr": "Beležnice i blokovi", "en": "Notebooks & notepads", "ru": "Блокноты и тетради"}'::jsonb,
 '{"sr": "Spiralne, hard cover, sticky notes, planeri", "en": "Spiral, hardcover, sticky notes, planners", "ru": "Спиральные, с твердой обложкой, стикеры, планировщики"}'::jsonb,
 '{"sr": "Beležnice i blokovi | Vondi", "en": "Notebooks & notepads | Vondi", "ru": "Блокноты и тетради | Vondi"}'::jsonb,
 '{"sr": "Kupite beležnice i blokove online", "en": "Buy notebooks and notepads online", "ru": "Купить блокноты и тетради онлайн"}'::jsonb,
 '📓', true),

('kancelarijski-pribor', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/kancelarijski-pribor', 105,
 '{"sr": "Kancelarijski pribor", "en": "Office supplies", "ru": "Канцелярские принадлежности"}'::jsonb,
 '{"sr": "Makaze, lepkovi, spajalice, čiode, klameri", "en": "Scissors, glue, staples, pins, clips", "ru": "Ножницы, клей, скрепки, кнопки, зажимы"}'::jsonb,
 '{"sr": "Kancelarijski pribor | Vondi", "en": "Office supplies | Vondi", "ru": "Канцелярские принадлежности | Vondi"}'::jsonb,
 '{"sr": "Kupite kancelarijski pribor online", "en": "Buy office supplies online", "ru": "Купить канцелярские принадлежности онлайн"}'::jsonb,
 '✂️', true),

('organizacija-stola', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/organizacija-stola', 106,
 '{"sr": "Organizacija stola", "en": "Desk organization", "ru": "Организация стола"}'::jsonb,
 '{"sr": "Posude za olovke, stalci, podmetači, kalendari", "en": "Pen holders, stands, desk mats, calendars", "ru": "Подставки для ручек, органайзеры, коврики, календари"}'::jsonb,
 '{"sr": "Organizacija stola | Vondi", "en": "Desk organization | Vondi", "ru": "Организация стола | Vondi"}'::jsonb,
 '{"sr": "Kupite organizaciju za sto online", "en": "Buy desk organization online", "ru": "Купить организацию для стола онлайн"}'::jsonb,
 '🗂️', true),

('kalkulator-i-oprema', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/kalkulator-i-oprema', 107,
 '{"sr": "Kalkulatori i oprema", "en": "Calculators & equipment", "ru": "Калькуляторы и оборудование"}'::jsonb,
 '{"sr": "Naučni kalkulatori, finansijski, džepni", "en": "Scientific calculators, financial, pocket", "ru": "Научные калькуляторы, финансовые, карманные"}'::jsonb,
 '{"sr": "Kalkulatori i oprema | Vondi", "en": "Calculators & equipment | Vondi", "ru": "Калькуляторы и оборудование | Vondi"}'::jsonb,
 '{"sr": "Kupite kalkulatore online", "en": "Buy calculators online", "ru": "Купить калькуляторы онлайн"}'::jsonb,
 '🧮', true),

('bele-table-i-prezentacije', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/bele-table-i-prezentacije', 108,
 '{"sr": "Bele table i prezentacije", "en": "Whiteboards & presentations", "ru": "Доски и презентации"}'::jsonb,
 '{"sr": "Magnetne table, flipchart, projektori, pointeri", "en": "Magnetic boards, flipcharts, projectors, pointers", "ru": "Магнитные доски, флипчарты, проекторы, указки"}'::jsonb,
 '{"sr": "Bele table i prezentacije | Vondi", "en": "Whiteboards & presentations | Vondi", "ru": "Доски и презентации | Vondi"}'::jsonb,
 '{"sr": "Kupite bele table online", "en": "Buy whiteboards online", "ru": "Купить доски онлайн"}'::jsonb,
 '📊', true),

('lomači-dokumenta', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/lomaci-dokumenta', 109,
 '{"sr": "Lomači dokumenta", "en": "Shredders", "ru": "Уничтожители документов"}'::jsonb,
 '{"sr": "Cross-cut shredders, micro-cut, heavy duty", "en": "Cross-cut shredders, micro-cut, heavy duty", "ru": "Шредеры крест-накрест, микро-резка, профессиональные"}'::jsonb,
 '{"sr": "Lomači dokumenta | Vondi", "en": "Shredders | Vondi", "ru": "Уничтожители документов | Vondi"}'::jsonb,
 '{"sr": "Kupite lomače dokumenta online", "en": "Buy shredders online", "ru": "Купить уничтожители документов онлайн"}'::jsonb,
 '🗑️', true),

('kancelarijski-nameštaj', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/kancelarijski-namestaj', 110,
 '{"sr": "Kancelarijski nameštaj", "en": "Office furniture", "ru": "Офисная мебель"}'::jsonb,
 '{"sr": "Radni stolovi, ergonomske stolice, ormari, police", "en": "Desks, ergonomic chairs, cabinets, shelves", "ru": "Рабочие столы, эргономичные стулья, шкафы, полки"}'::jsonb,
 '{"sr": "Kancelarijski nameštaj | Vondi", "en": "Office furniture | Vondi", "ru": "Офисная мебель | Vondi"}'::jsonb,
 '{"sr": "Kupite kancelarijski nameštaj online", "en": "Buy office furniture online", "ru": "Купить офисную мебель онлайн"}'::jsonb,
 '🪑', true),

('arhiviranje', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/arhiviranje', 111,
 '{"sr": "Arhiviranje", "en": "Archiving", "ru": "Архивирование"}'::jsonb,
 '{"sr": "Arhivske kutije, etikete, numeratori, pečati", "en": "Archive boxes, labels, numberers, stamps", "ru": "Архивные коробки, этикетки, нумераторы, печати"}'::jsonb,
 '{"sr": "Arhiviranje | Vondi", "en": "Archiving | Vondi", "ru": "Архивирование | Vondi"}'::jsonb,
 '{"sr": "Kupite arhivske materijale online", "en": "Buy archiving materials online", "ru": "Купить архивные материалы онлайн"}'::jsonb,
 '📦', true),

('laminating', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/laminating', 112,
 '{"sr": "Laminiranje", "en": "Laminating", "ru": "Ламинирование"}'::jsonb,
 '{"sr": "Laminator mašine, folije za laminiranje", "en": "Laminator machines, laminating pouches", "ru": "Ламинаторы, пленки для ламинирования"}'::jsonb,
 '{"sr": "Laminiranje | Vondi", "en": "Laminating | Vondi", "ru": "Ламинирование | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za laminiranje online", "en": "Buy laminating equipment online", "ru": "Купить оборудование для ламинирования онлайн"}'::jsonb,
 '📋', true),

('korporativni-pokloni', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/korporativni-pokloni', 113,
 '{"sr": "Korporativni pokloni", "en": "Corporate gifts", "ru": "Корпоративные подарки"}'::jsonb,
 '{"sr": "Branded olovke, USB, beleške, calendar, poklon setovi", "en": "Branded pens, USB, notebooks, calendars, gift sets", "ru": "Брендированные ручки, USB, блокноты, календари, подарочные наборы"}'::jsonb,
 '{"sr": "Korporativni pokloni | Vondi", "en": "Corporate gifts | Vondi", "ru": "Корпоративные подарки | Vondi"}'::jsonb,
 '{"sr": "Kupite korporativne poklone online", "en": "Buy corporate gifts online", "ru": "Купить корпоративные подарки онлайн"}'::jsonb,
 '🎁', true),

('skener-i-kopir', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/skener-i-kopir', 114,
 '{"sr": "Skeneri i kopir mašine", "en": "Scanners & copiers", "ru": "Сканеры и копиры"}'::jsonb,
 '{"sr": "Flatbed skeneri, portable, multifunkcijski kopiri", "en": "Flatbed scanners, portable, multifunction copiers", "ru": "Планшетные сканеры, портативные, МФУ копиры"}'::jsonb,
 '{"sr": "Skeneri i kopir mašine | Vondi", "en": "Scanners & copiers | Vondi", "ru": "Сканеры и копиры | Vondi"}'::jsonb,
 '{"sr": "Kupite skenere i kopir mašine online", "en": "Buy scanners and copiers online", "ru": "Купить сканеры и копиры онлайн"}'::jsonb,
 '🖨️', true),

('projektna-oprema', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' AND level = 1), 2, 'kancelarijski-materijal/projektna-oprema', 115,
 '{"sr": "Projektna oprema", "en": "Project equipment", "ru": "Проектное оборудование"}'::jsonb,
 '{"sr": "Gantt charts, project boards, planer setovi", "en": "Gantt charts, project boards, planner sets", "ru": "Диаграммы Ганта, проектные доски, наборы планировщиков"}'::jsonb,
 '{"sr": "Projektna oprema | Vondi", "en": "Project equipment | Vondi", "ru": "Проектное оборудование | Vondi"}'::jsonb,
 '{"sr": "Kupite projektnu opremu online", "en": "Buy project equipment online", "ru": "Купить проектное оборудование онлайн"}'::jsonb,
 '📈', true);

-- Progress notification
DO $$
BEGIN
  RAISE NOTICE 'Migration 20251217000003 (Part 1/4): Added 15 L2 categories (Kancelarija)';
END $$;

-- =============================================================================
-- Muzički instrumenti: +15 L2
-- =============================================================================

INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, icon, is_active) VALUES

('gitare-akusticne', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/gitare-akusticne', 100,
 '{"sr": "Akustične gitare", "en": "Acoustic guitars", "ru": "Акустические гитары"}'::jsonb,
 '{"sr": "Klasične, steel string, folk gitare", "en": "Classical, steel string, folk guitars", "ru": "Классические, стальные струны, фолк гитары"}'::jsonb,
 '{"sr": "Akustične gitare | Vondi", "en": "Acoustic guitars | Vondi", "ru": "Акустические гитары | Vondi"}'::jsonb,
 '{"sr": "Kupite akustične gitare online", "en": "Buy acoustic guitars online", "ru": "Купить акустические гитары онлайн"}'::jsonb,
 '🎸', true),

('elektricne-gitare', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/elektricne-gitare', 101,
 '{"sr": "Električne gitare", "en": "Electric guitars", "ru": "Электрогитары"}'::jsonb,
 '{"sr": "Stratocaster, Les Paul, semi-hollow, baritone", "en": "Stratocaster, Les Paul, semi-hollow, baritone", "ru": "Страт, Лес Пол, полуакустические, баритоны"}'::jsonb,
 '{"sr": "Električne gitare | Vondi", "en": "Electric guitars | Vondi", "ru": "Электрогитары | Vondi"}'::jsonb,
 '{"sr": "Kupite električne gitare online", "en": "Buy electric guitars online", "ru": "Купить электрогитары онлайн"}'::jsonb,
 '🎸', true),

('bas-gitare', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/bas-gitare', 102,
 '{"sr": "Bas gitare", "en": "Bass guitars", "ru": "Бас-гитары"}'::jsonb,
 '{"sr": "4-string, 5-string, fretless, upright bass", "en": "4-string, 5-string, fretless, upright bass", "ru": "4-струнные, 5-струнные, безладовые, контрабасы"}'::jsonb,
 '{"sr": "Bas gitare | Vondi", "en": "Bass guitars | Vondi", "ru": "Бас-гитары | Vondi"}'::jsonb,
 '{"sr": "Kupite bas gitare online", "en": "Buy bass guitars online", "ru": "Купить бас-гитары онлайн"}'::jsonb,
 '🎸', true),

('bubnjevi-setovi', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/bubnjevi-setovi', 103,
 '{"sr": "Bubnjevi setovi", "en": "Drum kits", "ru": "Барабанные установки"}'::jsonb,
 '{"sr": "Akustični, elektronski, jazz, rock setovi", "en": "Acoustic, electronic, jazz, rock sets", "ru": "Акустические, электронные, джаз, рок установки"}'::jsonb,
 '{"sr": "Bubnjevi setovi | Vondi", "en": "Drum kits | Vondi", "ru": "Барабанные установки | Vondi"}'::jsonb,
 '{"sr": "Kupite bubnjeve setove online", "en": "Buy drum kits online", "ru": "Купить барабанные установки онлайн"}'::jsonb,
 '🥁', true),

('klavijature-i-pianos', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/klavijature-i-pianos', 104,
 '{"sr": "Klavijature i pianoi", "en": "Keyboards & pianos", "ru": "Клавиатуры и пианино"}'::jsonb,
 '{"sr": "Digitalni piano, synthesizeri, MIDI klavijature", "en": "Digital pianos, synthesizers, MIDI keyboards", "ru": "Цифровые пианино, синтезаторы, MIDI клавиатуры"}'::jsonb,
 '{"sr": "Klavijature i pianoi | Vondi", "en": "Keyboards & pianos | Vondi", "ru": "Клавиатуры и пианино | Vondi"}'::jsonb,
 '{"sr": "Kupite klavijature i pianoe online", "en": "Buy keyboards and pianos online", "ru": "Купить клавиатуры и пианино онлайн"}'::jsonb,
 '🎹', true),

('duvacki-instrumenti', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/duvacki-instrumenti', 105,
 '{"sr": "Duvački instrumenti", "en": "Wind instruments", "ru": "Духовые инструменты"}'::jsonb,
 '{"sr": "Saksofoni, trube, klarineti, flaute, harmonika", "en": "Saxophones, trumpets, clarinets, flutes, harmonicas", "ru": "Саксофоны, трубы, кларнеты, флейты, гармоники"}'::jsonb,
 '{"sr": "Duvački instrumenti | Vondi", "en": "Wind instruments | Vondi", "ru": "Духовые инструменты | Vondi"}'::jsonb,
 '{"sr": "Kupite duvačke instrumente online", "en": "Buy wind instruments online", "ru": "Купить духовые инструменты онлайн"}'::jsonb,
 '🎺', true),

('violina-i-zicani', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/violina-i-zicani', 106,
 '{"sr": "Violina i žičani", "en": "Violin & strings", "ru": "Скрипка и струнные"}'::jsonb,
 '{"sr": "Violina, viola, čelo, gudala, koferi", "en": "Violin, viola, cello, bows, cases", "ru": "Скрипка, альт, виолончель, смычки, кофры"}'::jsonb,
 '{"sr": "Violina i žičani | Vondi", "en": "Violin & strings | Vondi", "ru": "Скрипка и струнные | Vondi"}'::jsonb,
 '{"sr": "Kupite violine i žičane instrumente online", "en": "Buy violins and string instruments online", "ru": "Купить скрипки и струнные инструменты онлайн"}'::jsonb,
 '🎻', true),

('audio-oprema-za-muziku', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/audio-oprema-za-muziku', 107,
 '{"sr": "Audio oprema za muziku", "en": "Music audio equipment", "ru": "Аудиооборудование для музыки"}'::jsonb,
 '{"sr": "Pojačala, zvučnici, mikser pultovi, efekti", "en": "Amplifiers, speakers, mixer consoles, effects", "ru": "Усилители, колонки, микшерные пульты, эффекты"}'::jsonb,
 '{"sr": "Audio oprema za muziku | Vondi", "en": "Music audio equipment | Vondi", "ru": "Аудиооборудование для музыки | Vondi"}'::jsonb,
 '{"sr": "Kupite audio opremu za muziku online", "en": "Buy music audio equipment online", "ru": "Купить аудиооборудование для музыки онлайн"}'::jsonb,
 '🔊', true),

('mikrofoni-za-muziku', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/mikrofoni-za-muziku', 108,
 '{"sr": "Mikrofoni za muziku", "en": "Music microphones", "ru": "Микрофоны для музыки"}'::jsonb,
 '{"sr": "Dinamički, kondenzatorski, USB, bežični", "en": "Dynamic, condenser, USB, wireless", "ru": "Динамические, конденсаторные, USB, беспроводные"}'::jsonb,
 '{"sr": "Mikrofoni za muziku | Vondi", "en": "Music microphones | Vondi", "ru": "Микрофоны для музыки | Vondi"}'::jsonb,
 '{"sr": "Kupite mikrofone za muziku online", "en": "Buy music microphones online", "ru": "Купить микрофоны для музыки онлайн"}'::jsonb,
 '🎤', true),

('efekti-i-procesori', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/efekti-i-procesori', 109,
 '{"sr": "Efekti i procesori", "en": "Effects & processors", "ru": "Эффекты и процессоры"}'::jsonb,
 '{"sr": "Gitarske pedale, multi-efekti, reverb, delay", "en": "Guitar pedals, multi-effects, reverb, delay", "ru": "Гитарные педали, мультиэффекты, реверберация, задержка"}'::jsonb,
 '{"sr": "Efekti i procesori | Vondi", "en": "Effects & processors | Vondi", "ru": "Эффекты и процессоры | Vondi"}'::jsonb,
 '{"sr": "Kupite efekte i procesore online", "en": "Buy effects and processors online", "ru": "Купить эффекты и процессоры онлайн"}'::jsonb,
 '🎛️', true),

('dj-oprema', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/dj-oprema', 110,
 '{"sr": "DJ oprema", "en": "DJ equipment", "ru": "DJ оборудование"}'::jsonb,
 '{"sr": "DJ kontroleri, gramofoni, mikser, slušalice DJ", "en": "DJ controllers, turntables, mixers, DJ headphones", "ru": "DJ контроллеры, проигрыватели, микшеры, DJ наушники"}'::jsonb,
 '{"sr": "DJ oprema | Vondi", "en": "DJ equipment | Vondi", "ru": "DJ оборудование | Vondi"}'::jsonb,
 '{"sr": "Kupite DJ opremu online", "en": "Buy DJ equipment online", "ru": "Купить DJ оборудование онлайн"}'::jsonb,
 '🎧', true),

('snimanje-i-produkcija', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/snimanje-i-produkcija', 111,
 '{"sr": "Snimanje i produkcija", "en": "Recording & production", "ru": "Запись и продакшн"}'::jsonb,
 '{"sr": "Audio interfejsi, DAW software, monitore, akustika", "en": "Audio interfaces, DAW software, monitors, acoustics", "ru": "Аудиоинтерфейсы, DAW софт, мониторы, акустика"}'::jsonb,
 '{"sr": "Snimanje i produkcija | Vondi", "en": "Recording & production | Vondi", "ru": "Запись и продакшн | Vondi"}'::jsonb,
 '{"sr": "Kupite opremu za snimanje online", "en": "Buy recording equipment online", "ru": "Купить оборудование для записи онлайн"}'::jsonb,
 '🎚️', true),

('muzicki-dodaci', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/muzicki-dodaci', 112,
 '{"sr": "Muzički dodaci", "en": "Music accessories", "ru": "Музыкальные аксессуары"}'::jsonb,
 '{"sr": "Žice, kablovi, stalci, koferi, torbice", "en": "Strings, cables, stands, cases, bags", "ru": "Струны, кабели, стойки, кофры, сумки"}'::jsonb,
 '{"sr": "Muzički dodaci | Vondi", "en": "Music accessories | Vondi", "ru": "Музыкальные аксессуары | Vondi"}'::jsonb,
 '{"sr": "Kupite muzičke dodatke online", "en": "Buy music accessories online", "ru": "Купить музыкальные аксессуары онлайн"}'::jsonb,
 '🎼', true),

('ukulele-i-mandoline', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/ukulele-i-mandoline', 113,
 '{"sr": "Ukulele i mandoline", "en": "Ukulele & mandolins", "ru": "Укулеле и мандолины"}'::jsonb,
 '{"sr": "Sopran, koncert, tenor ukulele, mandoline", "en": "Soprano, concert, tenor ukulele, mandolins", "ru": "Сопрано, концерт, тенор укулеле, мандолины"}'::jsonb,
 '{"sr": "Ukulele i mandoline | Vondi", "en": "Ukulele & mandolins | Vondi", "ru": "Укулеле и мандолины | Vondi"}'::jsonb,
 '{"sr": "Kupite ukulele i mandoline online", "en": "Buy ukulele and mandolins online", "ru": "Купить укулеле и мандолины онлайн"}'::jsonb,
 '🎸', true),

('orgulje-i-harmonijum', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' AND level = 1), 2, 'muzicki-instrumenti/orgulje-i-harmonijum', 114,
 '{"sr": "Orgulje i harmonijum", "en": "Organs & harmoniums", "ru": "Органы и гармониумы"}'::jsonb,
 '{"sr": "Crkvene orgulje, Hammond, akordeon, harmonijum", "en": "Church organs, Hammond, accordion, harmonium", "ru": "Церковные органы, Хаммонд, аккордеон, гармониум"}'::jsonb,
 '{"sr": "Orgulje i harmonijum | Vondi", "en": "Organs & harmoniums | Vondi", "ru": "Органы и гармониумы | Vondi"}'::jsonb,
 '{"sr": "Kupite orgulje i harmonijum online", "en": "Buy organs and harmoniums online", "ru": "Купить органы и гармониумы онлайн"}'::jsonb,
 '🎹', true);

-- Progress notification
DO $$
BEGIN
  RAISE NOTICE 'Migration 20251217000003 (Part 2/4): Added 15 L2 categories (Muzički instrumenti)';
END $$;
