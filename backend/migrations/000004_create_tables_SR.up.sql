-- Сербские переводы категорий для таблицы translations
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES
-- Основные категории
('category', 1, 'sr', 'name', 'Nekretnine', true, true, NOW(), NOW()),
('category', 2, 'sr', 'name', 'Automobili', true, true, NOW(), NOW()),
('category', 3, 'sr', 'name', 'Elektronika', true, true, NOW(), NOW()),
('category', 4, 'sr', 'name', 'Sve za kuću i stan', true, true, NOW(), NOW()),
('category', 5, 'sr', 'name', 'Sve za baštu', true, true, NOW(), NOW()),
('category', 6, 'sr', 'name', 'Hobi i odmor', true, true, NOW(), NOW()),
('category', 7, 'sr', 'name', 'Životinje', true, true, NOW(), NOW()),
('category', 8, 'sr', 'name', 'Gotov biznis i oprema', true, true, NOW(), NOW()),
('category', 9, 'sr', 'name', 'Ostalo', true, true, NOW(), NOW()),
('category', 10, 'sr', 'name', 'Poslovi', true, true, NOW(), NOW()),
('category', 11, 'sr', 'name', 'Odeća, obuća, aksesori', true, true, NOW(), NOW()),
('category', 12, 'sr', 'name', 'Dečja odeća i obuća, aksesori', true, true, NOW(), NOW()),
-- Недвижимость (подкатегории)
('category', 13, 'sr', 'name', 'Stanovi', true, true, NOW(), NOW()),
('category', 14, 'sr', 'name', 'Kuće', true, true, NOW(), NOW()),
('category', 15, 'sr', 'name', 'Poslovni prostori', true, true, NOW(), NOW()),
('category', 16, 'sr', 'name', 'Zemljište', true, true, NOW(), NOW()),
('category', 17, 'sr', 'name', 'Garaže', true, true, NOW(), NOW()),
('category', 18, 'sr', 'name', 'Sobe', true, true, NOW(), NOW()),
('category', 19, 'sr', 'name', 'Nekretnine u inostranstvu', true, true, NOW(), NOW()),
('category', 20, 'sr', 'name', 'Prodaja poslovnih objekata', true, true, NOW(), NOW()),
('category', 21, 'sr', 'name', 'Iznajmljivanje poslovnih objekata', true, true, NOW(), NOW()),
-- Подкатегории автотранспорта
('category', 22, 'sr', 'name', 'Putnički automobili', true, true, NOW(), NOW()),
('category', 23, 'sr', 'name', 'Teretna vozila', true, true, NOW(), NOW()),
('category', 24, 'sr', 'name', 'Specijalna tehnika', true, true, NOW(), NOW()),
('category', 25, 'sr', 'name', 'Poljoprivredna mehanizacija', true, true, NOW(), NOW()),
('category', 26, 'sr', 'name', 'Iznajmljivanje automobila i specijalne tehnike', true, true, NOW(), NOW()),
('category', 27, 'sr', 'name', 'Motocikli i mototehnika', true, true, NOW(), NOW()),
('category', 28, 'sr', 'name', 'Plovila', true, true, NOW(), NOW()),
('category', 29, 'sr', 'name', 'Rezervni delovi i oprema', true, true, NOW(), NOW()),

-- Грузовые автомобили
('category', 30, 'sr', 'name', 'Kamioni', true, true, NOW(), NOW()),
('category', 31, 'sr', 'name', 'Poluprikolice', true, true, NOW(), NOW()),
('category', 32, 'sr', 'name', 'Laka komercijalna vozila', true, true, NOW(), NOW()),
('category', 33, 'sr', 'name', 'Autobusi', true, true, NOW(), NOW()),

-- Спецтехника
('category', 34, 'sr', 'name', 'Bageri', true, true, NOW(), NOW()),
('category', 35, 'sr', 'name', 'Utovarivači', true, true, NOW(), NOW()),
('category', 36, 'sr', 'name', 'Rovokopači-utovarivači', true, true, NOW(), NOW()),
('category', 37, 'sr', 'name', 'Auto-dizalice', true, true, NOW(), NOW()),
('category', 38, 'sr', 'name', 'Auto-betonske mešalice', true, true, NOW(), NOW()),
('category', 39, 'sr', 'name', 'Drumski valjci', true, true, NOW(), NOW()),
('category', 40, 'sr', 'name', 'Cisterne za pranje ulica', true, true, NOW(), NOW()),
('category', 41, 'sr', 'name', 'Vozila za odvoz smeća', true, true, NOW(), NOW()),
('category', 42, 'sr', 'name', 'Auto-platforme', true, true, NOW(), NOW()),
('category', 43, 'sr', 'name', 'Buldožeri', true, true, NOW(), NOW()),
('category', 44, 'sr', 'name', 'Grejderi', true, true, NOW(), NOW()),
('category', 45, 'sr', 'name', 'Bušilice', true, true, NOW(), NOW()),

-- Сельхозтехника
('category', 46, 'sr', 'name', 'Traktori', true, true, NOW(), NOW()),
('category', 47, 'sr', 'name', 'Mini-traktori', true, true, NOW(), NOW()),
('category', 48, 'sr', 'name', 'Prese za baliranje', true, true, NOW(), NOW()),
('category', 49, 'sr', 'name', 'Drljače', true, true, NOW(), NOW()),
('category', 50, 'sr', 'name', 'Kosačice', true, true, NOW(), NOW()),

-- Всё для сада (подкатегории)
('category', 51, 'sr', 'name', 'Baštenski nameštaj', true, true, NOW(), NOW()),
('category', 52, 'sr', 'name', 'Baštenski alati', true, true, NOW(), NOW()),
('category', 53, 'sr', 'name', 'Biljke za baštu', true, true, NOW(), NOW()),
('category', 54, 'sr', 'name', 'Semena i sadnice', true, true, NOW(), NOW()),
('category', 55, 'sr', 'name', 'Roštilji i oprema', true, true, NOW(), NOW()),
('category', 56, 'sr', 'name', 'Bazeni i oprema', true, true, NOW(), NOW()),
('category', 57, 'sr', 'name', 'Navodnjavanje', true, true, NOW(), NOW()),
('category', 58, 'sr', 'name', 'Kompostiranje', true, true, NOW(), NOW()),
('category', 59, 'sr', 'name', 'Staklenici i plastenici', true, true, NOW(), NOW()),
('category', 60, 'sr', 'name', 'Gnojiva i zemlja', true, true, NOW(), NOW()),

-- Хобби и отдых (подкатегории)
('category', 101, 'sr', 'name', 'Muzički instrumenti', true, true, NOW(), NOW()),
('category', 102, 'sr', 'name', 'Knjige', true, true, NOW(), NOW()),
('category', 103, 'sr', 'name', 'Sportska oprema', true, true, NOW(), NOW()),
('category', 104, 'sr', 'name', 'Putovanja', true, true, NOW(), NOW()),
('category', 105, 'sr', 'name', 'Kolekcionarstvo', true, true, NOW(), NOW()),
('category', 106, 'sr', 'name', 'Umetnički predmeti', true, true, NOW(), NOW()),
('category', 107, 'sr', 'name', 'Igračke', true, true, NOW(), NOW()),
('category', 108, 'sr', 'name', 'Bicikli', true, true, NOW(), NOW()),
('category', 109, 'sr', 'name', 'Lov i ribolov', true, true, NOW(), NOW()),
('category', 110, 'sr', 'name', 'Kampovanje', true, true, NOW(), NOW()),
('category', 111, 'sr', 'name', 'Antikviteti', true, true, NOW(), NOW()),
('category', 112, 'sr', 'name', 'Rukotvorine', true, true, NOW(), NOW()),
('category', 113, 'sr', 'name', 'Ulaznice za događaje', true, true, NOW(), NOW()),

-- Музыкальные инструменты (подкатегории)
('category', 114, 'sr', 'name', 'Žičani instrumenti', true, true, NOW(), NOW()),
('category', 115, 'sr', 'name', 'Klaviri i klavijature', true, true, NOW(), NOW()),
('category', 116, 'sr', 'name', 'Udaraljke', true, true, NOW(), NOW()),
('category', 117, 'sr', 'name', 'Duvački instrumenti', true, true, NOW(), NOW()),
('category', 118, 'sr', 'name', 'Harmonike i harmonijumi', true, true, NOW(), NOW()),
('category', 119, 'sr', 'name', 'Audio oprema', true, true, NOW(), NOW()),
('category', 120, 'sr', 'name', 'Dodatna oprema za instrumente', true, true, NOW(), NOW()),

-- Животные
('category', 200, 'sr', 'name', 'Psi', true, true, NOW(), NOW()),
('category', 201, 'sr', 'name', 'Mačke', true, true, NOW(), NOW()),
('category', 202, 'sr', 'name', 'Ptice', true, true, NOW(), NOW()),
('category', 203, 'sr', 'name', 'Akvarijum', true, true, NOW(), NOW()),
('category', 204, 'sr', 'name', 'Druge životinje', true, true, NOW(), NOW()),

-- Аренда авто и спецтехники
('category', 205, 'sr', 'name', 'Automobili', true, true, NOW(), NOW()),
('category', 206, 'sr', 'name', 'Podizna oprema', true, true, NOW(), NOW()),
('category', 207, 'sr', 'name', 'Zemljokopačka oprema', true, true, NOW(), NOW()),
('category', 208, 'sr', 'name', 'Komunalna oprema', true, true, NOW(), NOW()),
('category', 209, 'sr', 'name', 'Oprema za putogradnju', true, true, NOW(), NOW()),
('category', 210, 'sr', 'name', 'Teretna vozila', true, true, NOW(), NOW()),
('category', 211, 'sr', 'name', 'Utovarivačka oprema', true, true, NOW(), NOW()),
('category', 212, 'sr', 'name', 'Dodatna oprema', true, true, NOW(), NOW()),
('category', 213, 'sr', 'name', 'Prikolice', true, true, NOW(), NOW()),
('category', 214, 'sr', 'name', 'Poljoprivredna mehanizacija', true, true, NOW(), NOW()),
('category', 215, 'sr', 'name', 'Kamperi', true, true, NOW(), NOW()),

-- Мотоциклы и мототехника
('category', 216, 'sr', 'name', 'Terenska vozila', true, true, NOW(), NOW()),
('category', 217, 'sr', 'name', 'Karting', true, true, NOW(), NOW()),
('category', 218, 'sr', 'name', 'Kvadovi i bagiji', true, true, NOW(), NOW()),
('category', 219, 'sr', 'name', 'Mopedi', true, true, NOW(), NOW()),
('category', 220, 'sr', 'name', 'Skuteri', true, true, NOW(), NOW()),
('category', 221, 'sr', 'name', 'Motocikli', true, true, NOW(), NOW()),
('category', 222, 'sr', 'name', 'Motorne sanke', true, true, NOW(), NOW()),

-- Водный транспорт
('category', 223, 'sr', 'name', 'Čamci na vesla', true, true, NOW(), NOW()),
('category', 224, 'sr', 'name', 'Kajaci', true, true, NOW(), NOW()),
('category', 225, 'sr', 'name', 'Skuteri na vodi', true, true, NOW(), NOW()),
('category', 226, 'sr', 'name', 'Jahte i brodovi', true, true, NOW(), NOW()),
('category', 227, 'sr', 'name', 'Čamci sa motorom i motori', true, true, NOW(), NOW()),

-- Запчасти и аксессуары
('category', 228, 'sr', 'name', 'Rezervni delovi', true, true, NOW(), NOW()),
('category', 229, 'sr', 'name', 'Gume, felne i točkovi', true, true, NOW(), NOW()),
('category', 230, 'sr', 'name', 'Audio i video oprema', true, true, NOW(), NOW()),
('category', 231, 'sr', 'name', 'Dodatna oprema', true, true, NOW(), NOW()),
('category', 232, 'sr', 'name', 'Ulja i auto-hemija', true, true, NOW(), NOW()),
('category', 233, 'sr', 'name', 'Alati', true, true, NOW(), NOW()),
('category', 234, 'sr', 'name', 'Krovni nosači i kuke', true, true, NOW(), NOW()),
('category', 235, 'sr', 'name', 'Prikolice', true, true, NOW(), NOW()),
('category', 236, 'sr', 'name', 'Oprema', true, true, NOW(), NOW()),
('category', 237, 'sr', 'name', 'Protivprovalni uređaji', true, true, NOW(), NOW()),
('category', 238, 'sr', 'name', 'GPS-navigatori', true, true, NOW(), NOW()),

-- Электроника
('category', 239, 'sr', 'name', 'Telefoni', true, true, NOW(), NOW()),
('category', 240, 'sr', 'name', 'Audio i video', true, true, NOW(), NOW()),
('category', 241, 'sr', 'name', 'Računarska oprema', true, true, NOW(), NOW()),
('category', 242, 'sr', 'name', 'Igre, konzole i programi', true, true, NOW(), NOW()),
('category', 243, 'sr', 'name', 'Laptopi', true, true, NOW(), NOW()),
('category', 244, 'sr', 'name', 'Foto-tehnika', true, true, NOW(), NOW()),
('category', 245, 'sr', 'name', 'Tableti i e-knjige', true, true, NOW(), NOW()),
('category', 246, 'sr', 'name', 'Kancelarijska oprema i potrošni materijal', true, true, NOW(), NOW()),
('category', 247, 'sr', 'name', 'Kućni aparati', true, true, NOW(), NOW()),

-- Подкатегории Телефонов
('category', 248, 'sr', 'name', 'Mobilni telefoni', true, true, NOW(), NOW()),
('category', 249, 'sr', 'name', 'Dodatna oprema', true, true, NOW(), NOW()),
('category', 250, 'sr', 'name', 'Radio stanice', true, true, NOW(), NOW()),
('category', 251, 'sr', 'name', 'Fiksni telefoni', true, true, NOW(), NOW()),

-- Аксессуары для телефонов
('category', 252, 'sr', 'name', 'Baterije', true, true, NOW(), NOW()),
('category', 253, 'sr', 'name', 'Slušalice', true, true, NOW(), NOW()),
('category', 254, 'sr', 'name', 'Punjači', true, true, NOW(), NOW()),
('category', 255, 'sr', 'name', 'Kablovi i adapteri', true, true, NOW(), NOW()),
('category', 256, 'sr', 'name', 'Modemi i ruteri', true, true, NOW(), NOW()),
('category', 257, 'sr', 'name', 'Maske i folije', true, true, NOW(), NOW()),
('category', 258, 'sr', 'name', 'Rezervni delovi', true, true, NOW(), NOW()),

-- Аудио и видео
('category', 259, 'sr', 'name', 'Televizori i projektori', true, true, NOW(), NOW()),
('category', 260, 'sr', 'name', 'Slušalice', true, true, NOW(), NOW()),
('category', 261, 'sr', 'name', 'Zvučnici, koloone, subvuferi', true, true, NOW(), NOW()),
('category', 262, 'sr', 'name', 'Dodatna oprema', true, true, NOW(), NOW()),
('category', 263, 'sr', 'name', 'Muzički centri, radio aparati', true, true, NOW(), NOW()),
('category', 264, 'sr', 'name', 'Pojačala i risiveri', true, true, NOW(), NOW()),
('category', 265, 'sr', 'name', 'Video kamere', true, true, NOW(), NOW()),
('category', 266, 'sr', 'name', 'Video, DVD i Blu-ray plejeri', true, true, NOW(), NOW()),
('category', 267, 'sr', 'name', 'Kablovi i adapteri', true, true, NOW(), NOW()),
('category', 268, 'sr', 'name', 'Muzika i filmovi', true, true, NOW(), NOW()),
('category', 269, 'sr', 'name', 'Mikrofoni', true, true, NOW(), NOW()),
('category', 270, 'sr', 'name', 'MP3-plejeri', true, true, NOW(), NOW()),

-- Товары для компьютера
('category', 271, 'sr', 'name', 'Kućišta računara', true, true, NOW(), NOW()),
('category', 272, 'sr', 'name', 'All-in-One računari', true, true, NOW(), NOW()),
('category', 273, 'sr', 'name', 'Komponente', true, true, NOW(), NOW()),
('category', 274, 'sr', 'name', 'Monitori i delovi', true, true, NOW(), NOW()),
('category', 275, 'sr', 'name', 'Mrežna oprema', true, true, NOW(), NOW()),
('category', 276, 'sr', 'name', 'Tastature i miševi', true, true, NOW(), NOW()),
('category', 277, 'sr', 'name', 'Džojstici i volani', true, true, NOW(), NOW()),
('category', 278, 'sr', 'name', 'USB i memorijske kartice', true, true, NOW(), NOW()),
('category', 279, 'sr', 'name', 'Web-kamere', true, true, NOW(), NOW()),
('category', 280, 'sr', 'name', 'TV-tjuneri', true, true, NOW(), NOW()),

-- Комплектующие для компьютера
('category', 281, 'sr', 'name', 'CD, DVD i Blu-ray uređaji', true, true, NOW(), NOW()),
('category', 282, 'sr', 'name', 'Napajanja', true, true, NOW(), NOW()),
('category', 283, 'sr', 'name', 'Grafičke karte', true, true, NOW(), NOW()),
('category', 284, 'sr', 'name', 'Hard diskovi', true, true, NOW(), NOW()),
('category', 285, 'sr', 'name', 'Zvučne karte', true, true, NOW(), NOW()),
('category', 286, 'sr', 'name', 'Kontroleri', true, true, NOW(), NOW()),
('category', 287, 'sr', 'name', 'Kućišta', true, true, NOW(), NOW()),
('category', 288, 'sr', 'name', 'Matične ploče', true, true, NOW(), NOW()),
('category', 289, 'sr', 'name', 'RAM memorija', true, true, NOW(), NOW()),
('category', 290, 'sr', 'name', 'Procesori', true, true, NOW(), NOW()),
('category', 291, 'sr', 'name', 'Sistemi hlađenja', true, true, NOW(), NOW()),

-- Игры, приставки и программы
('category', 292, 'sr', 'name', 'Konzole za igru', true, true, NOW(), NOW()),
('category', 293, 'sr', 'name', 'Igre za konzole', true, true, NOW(), NOW()),
('category', 294, 'sr', 'name', 'Računarske igre', true, true, NOW(), NOW()),
('category', 295, 'sr', 'name', 'Programi', true, true, NOW(), NOW()),

-- Фототехника
('category', 296, 'sr', 'name', 'Oprema i dodaci', true, true, NOW(), NOW()),
('category', 297, 'sr', 'name', 'Objektivi', true, true, NOW(), NOW()),
('category', 298, 'sr', 'name', 'Kompaktni fotoaparati', true, true, NOW(), NOW()),
('category', 299, 'sr', 'name', 'Filmski fotoaparati', true, true, NOW(), NOW()),
('category', 300, 'sr', 'name', 'DSLR fotoaparati', true, true, NOW(), NOW()),

-- Планшеты и электронные книги
('category', 301, 'sr', 'name', 'Tableti', true, true, NOW(), NOW()),
('category', 302, 'sr', 'name', 'E-knjige', true, true, NOW(), NOW()),
('category', 303, 'sr', 'name', 'Dodatna oprema', true, true, NOW(), NOW()),

-- Оргтехника и расходники
('category', 304, 'sr', 'name', 'MFU, kopir aparati i skeneri', true, true, NOW(), NOW()),
('category', 305, 'sr', 'name', 'Štampači', true, true, NOW(), NOW()),
('category', 306, 'sr', 'name', 'Kancelarijski materijal', true, true, NOW(), NOW()),
('category', 307, 'sr', 'name', 'UPS, mrežni filteri', true, true, NOW(), NOW()),
('category', 308, 'sr', 'name', 'Telefonija', true, true, NOW(), NOW()),
('category', 309, 'sr', 'name', 'Uništavači papira', true, true, NOW(), NOW()),
('category', 310, 'sr', 'name', 'Potrošni materijal', true, true, NOW(), NOW()),

-- Бытовая техника
('category', 311, 'sr', 'name', 'Za kuhinju', true, true, NOW(), NOW()),
('category', 312, 'sr', 'name', 'Za domaćinstvo', true, true, NOW(), NOW()),
('category', 313, 'sr', 'name', 'Klimatizaciona oprema', true, true, NOW(), NOW()),
('category', 314, 'sr', 'name', 'Za ličnu negu', true, true, NOW(), NOW()),

-- Техника для кухни
('category', 315, 'sr', 'name', 'Aspiratori', true, true, NOW(), NOW()),
('category', 316, 'sr', 'name', 'Mali kuhinjski aparati', true, true, NOW(), NOW()),
('category', 317, 'sr', 'name', 'Mikrotalasne pećnice', true, true, NOW(), NOW()),
('category', 318, 'sr', 'name', 'Šporeti i rerne', true, true, NOW(), NOW()),
('category', 319, 'sr', 'name', 'Mašine za pranje sudova', true, true, NOW(), NOW()),
('category', 320, 'sr', 'name', 'Frižideri i zamrzivači', true, true, NOW(), NOW()),

-- Техника для дома
('category', 321, 'sr', 'name', 'Usisivači i delovi', true, true, NOW(), NOW()),
('category', 322, 'sr', 'name', 'Mašine za pranje i sušenje veša', true, true, NOW(), NOW()),
('category', 323, 'sr', 'name', 'Pegle', true, true, NOW(), NOW()),
('category', 324, 'sr', 'name', 'Oprema za šivenje', true, true, NOW(), NOW()),

-- Климатическое оборудование
('category', 325, 'sr', 'name', 'Ventilatori', true, true, NOW(), NOW()),
('category', 326, 'sr', 'name', 'Klima uređaji i delovi', true, true, NOW(), NOW()),
('category', 327, 'sr', 'name', 'Grejalice', true, true, NOW(), NOW()),
('category', 328, 'sr', 'name', 'Prečišćivači vazduha', true, true, NOW(), NOW()),
('category', 329, 'sr', 'name', 'Termometri i meteo stanice', true, true, NOW(), NOW()),

-- Техника для индивидуального ухода
('category', 330, 'sr', 'name', 'Brijači i trimeri', true, true, NOW(), NOW()),
('category', 331, 'sr', 'name', 'Mašinice za šišanje', true, true, NOW(), NOW()),
('category', 332, 'sr', 'name', 'Fenovi i aparati za stilizovanje', true, true, NOW(), NOW()),
('category', 333, 'sr', 'name', 'Epilatori', true, true, NOW(), NOW()),

-- Для дома и квартиры
('category', 334, 'sr', 'name', 'Remont i izgradnja', true, true, NOW(), NOW()),
('category', 335, 'sr', 'name', 'Nameštaj i enterijer', true, true, NOW(), NOW()),
('category', 336, 'sr', 'name', 'Prehrambeni proizvodi', true, true, NOW(), NOW()),
('category', 337, 'sr', 'name', 'Posuđe i kuhinjski pribor', true, true, NOW(), NOW()),

-- Основные категории и другие важные для демонстрации
('category', 338, 'sr', 'name', 'Vrata', true, true, NOW(), NOW()),
('category', 339, 'sr', 'name', 'Alati', true, true, NOW(), NOW()),
('category', 340, 'sr', 'name', 'Kamini i grejalice', true, true, NOW(), NOW()),
('category', 341, 'sr', 'name', 'Prozori i balkoni', true, true, NOW(), NOW()),
('category', 342, 'sr', 'name', 'Plafoni', true, true, NOW(), NOW()),
('category', 343, 'sr', 'name', 'Za baštu i vikendicu', true, true, NOW(), NOW()),
('category', 344, 'sr', 'name', 'Sanitarije, vodovod i sauna', true, true, NOW(), NOW()),
('category', 345, 'sr', 'name', 'Gotove građevine i brvnare', true, true, NOW(), NOW()),
('category', 346, 'sr', 'name', 'Kapije, ograde i barijere', true, true, NOW(), NOW()),
('category', 347, 'sr', 'name', 'Zaštita i alarmi', true, true, NOW(), NOW()),
('category', 348, 'sr', 'name', 'Kreveti, sofe i fotelje', true, true, NOW(), NOW()),
('category', 349, 'sr', 'name', 'Tekstil i tepisi', true, true, NOW(), NOW()),
('category', 350, 'sr', 'name', 'Osvetljenje', true, true, NOW(), NOW()),
('category', 351, 'sr', 'name', 'Kompjuterski stolovi i stolice', true, true, NOW(), NOW()),
('category', 352, 'sr', 'name', 'Ormari, komode i police', true, true, NOW(), NOW()),
('category', 353, 'sr', 'name', 'Kuhinjski elementi', true, true, NOW(), NOW()),
('category', 354, 'sr', 'name', 'Stolovi i stolice', true, true, NOW(), NOW()),
('category', 355, 'sr', 'name', 'Sobne biljke', true, true, NOW(), NOW()),
('category', 356, 'sr', 'name', 'Ukrasno lišće', true, true, NOW(), NOW()),
('category', 357, 'sr', 'name', 'Cvetne biljke', true, true, NOW(), NOW()),
('category', 358, 'sr', 'name', 'Palme i fikusi', true, true, NOW(), NOW()),
('category', 359, 'sr', 'name', 'Kaktusi i sukulenti', true, true, NOW(), NOW()),
('category', 360, 'sr', 'name', 'Čaj, kafa, kakao', true, true, NOW(), NOW()),
('category', 361, 'sr', 'name', 'Napici', true, true, NOW(), NOW()),
('category', 362, 'sr', 'name', 'Riba, morski plodovi, ikra', true, true, NOW(), NOW()),
('category', 363, 'sr', 'name', 'Meso, živina, iznutrice', true, true, NOW(), NOW()),
('category', 364, 'sr', 'name', 'Konditorski proizvodi', true, true, NOW(), NOW()),
('category', 366, 'sr', 'name', 'Posuđe', true, true, NOW(), NOW()),
('category', 367, 'sr', 'name', 'Kuhinjski pribor', true, true, NOW(), NOW()),
('category', 368, 'sr', 'name', 'Postavljanje stola', true, true, NOW(), NOW()),
('category', 369, 'sr', 'name', 'Priprema hrane', true, true, NOW(), NOW()),
('category', 370, 'sr', 'name', 'Čuvanje namirnica', true, true, NOW(), NOW()),
('category', 371, 'sr', 'name', 'Priprema napitaka', true, true, NOW(), NOW()),
('category', 372, 'sr', 'name', 'Proizvodi za domaćinstvo', true, true, NOW(), NOW()),

-- Все для сада (родитель id=5)
('category', 373, 'sr', 'name', 'Baštenski nameštaj', true, true, NOW(), NOW()),
('category', 374, 'sr', 'name', 'Osvetljenje', true, true, NOW(), NOW()),
('category', 375, 'sr', 'name', 'Uređenje enterijera', true, true, NOW(), NOW()),
('category', 376, 'sr', 'name', 'Biljke i semena', true, true, NOW(), NOW()),
('category', 377, 'sr', 'name', 'Bašta i vrt', true, true, NOW(), NOW()),
('category', 378, 'sr', 'name', 'Baštensko bilje', true, true, NOW(), NOW()),
('category', 379, 'sr', 'name', 'Semena, lukovice, krtole', true, true, NOW(), NOW()),
('category', 380, 'sr', 'name', 'Proizvodi za negu biljaka', true, true, NOW(), NOW()),
('category', 381, 'sr', 'name', 'Panoi i veštačke biljke', true, true, NOW(), NOW()),

-- Подкатегории Садовые растения (родитель id=378)
('category', 382, 'sr', 'name', 'Ukrasno žbunje i drveće', true, true, NOW(), NOW()),
('category', 383, 'sr', 'name', 'Četinari', true, true, NOW(), NOW()),
('category', 384, 'sr', 'name', 'Višegodišnje biljke', true, true, NOW(), NOW()),
('category', 385, 'sr', 'name', 'Voćno bilje', true, true, NOW(), NOW()),
('category', 386, 'sr', 'name', 'Travnjak', true, true, NOW(), NOW()),
('category', 387, 'sr', 'name', 'Zeleni začini i začinsko bilje', true, true, NOW(), NOW()),

-- Подкатегории Товары для ухода за растениями (родитель id=380)
('category', 388, 'sr', 'name', 'Zemlja i supstrati', true, true, NOW(), NOW()),
('category', 389, 'sr', 'name', 'Đubriva', true, true, NOW(), NOW()),
('category', 390, 'sr', 'name', 'Sredstva protiv štetočina i korova', true, true, NOW(), NOW()),
('category', 391, 'sr', 'name', 'Saksije i žardinijere', true, true, NOW(), NOW()),
('category', 392, 'sr', 'name', 'Fitolampe', true, true, NOW(), NOW()),
('category', 393, 'sr', 'name', 'Merači vlažnosti', true, true, NOW(), NOW()),
('category', 394, 'sr', 'name', 'Staklenici, gredice, leje', true, true, NOW(), NOW()),

-- Аквариум (подкатегории)
('category', 401, 'sr', 'name', 'Akvarijumi i oprema', true, true, NOW(), NOW()),
('category', 402, 'sr', 'name', 'Ribe', true, true, NOW(), NOW()),
('category', 403, 'sr', 'name', 'Biljke za akvarijum', true, true, NOW(), NOW()),
('category', 404, 'sr', 'name', 'Hrana za ribe', true, true, NOW(), NOW()),
('category', 405, 'sr', 'name', 'Filtri i dodatna oprema', true, true, NOW(), NOW()),
-- Другие животные (подкатегории)
('category', 406, 'sr', 'name', 'Mali glodari', true, true, NOW(), NOW()),
('category', 407, 'sr', 'name', 'Reptili', true, true, NOW(), NOW()),
('category', 408, 'sr', 'name', 'Egzotične životinje', true, true, NOW(), NOW()),
('category', 409, 'sr', 'name', 'Pčele i košnice', true, true, NOW(), NOW()),
('category', 410, 'sr', 'name', 'Domaće životinje', true, true, NOW(), NOW()),
('category', 411, 'sr', 'name', 'Konji', true, true, NOW(), NOW()),
('category', 412, 'sr', 'name', 'Hrana za životinje', true, true, NOW(), NOW()),
('category', 413, 'sr', 'name', 'Oprema za životinje', true, true, NOW(), NOW()),

-- Прочие категории (добавьте остальные категории при необходимости)
('category', 473, 'sr', 'name', 'Rakija i vino', true, true, NOW(), NOW()),
('category', 474, 'sr', 'name', 'Domaći sirevi', true, true, NOW(), NOW()),
('category', 475, 'sr', 'name', 'Kajmak', true, true, NOW(), NOW()),
('category', 476, 'sr', 'name', 'Ajvar', true, true, NOW(), NOW()),
('category', 477, 'sr', 'name', 'Šljivovica', true, true, NOW(), NOW()),
('category', 478, 'sr', 'name', 'Lozovača', true, true, NOW(), NOW()),
('category', 479, 'sr', 'name', 'Voćna rakija', true, true, NOW(), NOW()),
('category', 480, 'sr', 'name', 'Domaće vino', true, true, NOW(), NOW()),
('category', 481, 'sr', 'name', 'Narodni zanati', true, true, NOW(), NOW()),
('category', 482, 'sr', 'name', 'Opanci', true, true, NOW(), NOW()),
('category', 483, 'sr', 'name', 'Keramika', true, true, NOW(), NOW()),
('category', 484, 'sr', 'name', 'Vez', true, true, NOW(), NOW()),
('category', 485, 'sr', 'name', 'Tkanje', true, true, NOW(), NOW()),
('category', 486, 'sr', 'name', 'Narodni instrumenti', true, true, NOW(), NOW()),

-- Сельскохозяйственные категории
('category', 487, 'sr', 'name', 'Pčelarstvo', true, true, NOW(), NOW()),

-- Подкатегории пчеловодства (родитель id=487)
('category', 488, 'sr', 'name', 'Med', true, true, NOW(), NOW()),
('category', 489, 'sr', 'name', 'Pčelinji vosak', true, true, NOW(), NOW()),
('category', 490, 'sr', 'name', 'Propolis', true, true, NOW(), NOW()),
('category', 491, 'sr', 'name', 'Pčelarski pribor', true, true, NOW(), NOW()),

-- Категории для сезонных работ
('category', 492, 'sr', 'name', 'Sezonski poslovi', true, true, NOW(), NOW()),
-- Подкатегории сезонных работ (родитель id=492)
('category', 493, 'sr', 'name', 'Berba', true, true, NOW(), NOW()),
('category', 494, 'sr', 'name', 'Rad u vinogradu', true, true, NOW(), NOW()),
('category', 495, 'sr', 'name', 'Sezonski građevinski radovi', true, true, NOW(), NOW()),

-- Туристические услуги
('category', 496, 'sr', 'name', 'Seoski turizam', true, true, NOW(), NOW()),

-- Подкатегории сельского туризма (родитель id=496)
('category', 497, 'sr', 'name', 'Etno-sela', true, true, NOW(), NOW()),
('category', 498, 'sr', 'name', 'Vinske ture', true, true, NOW(), NOW()),
('category', 499, 'sr', 'name', 'Agroturizam', true, true, NOW(), NOW()),
('category', 500, 'sr', 'name', 'Planinski turizam', true, true, NOW(), NOW());


-- Обновляем sequence для translations
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);


-- Insert marketplace listings
INSERT INTO marketplace_listings (id, user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, views_count, created_at, updated_at, show_on_map, original_language) VALUES
(8, 2, 22, 'Toyota Corolla 2018', 'Продајем Toyota Corolla 2018 годиште, 80.000 км, одлично стање. Први власник, редовно одржавање, сва документација доступна.', 1150000.00, 'used', 'active', 'Нови Сад, Србија', 45.26710000, 19.83350000, 'Нови Сад', 'Србија', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'sr'),
(9, 3, 248, 'mobile Samsung Galaxy S21', 'Selling Samsung Galaxy S21, 256GB, Deep Purple. Perfect condition, complete set with original box and accessories. AppleCare+ until 2024.', 120000.00, 'used', 'active', 'Novi Sad, Serbia', 45.25510000, 19.84520000, 'Novi Sad', 'Serbia', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'en'),
(10, 4, 272, 'Игровой компьютер RTX 4080', 'Продаю мощный игровой ПК: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеален для любых игр и тяжелых задач.', 350000.00, 'used', 'active', 'Нови-Сад, Сербия', 45.25410000, 19.84010000, 'Нови-Сад', 'Сербия', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'ru'),
(12, 2, 22, 'автомобиль Toyota Corolla 2018', 'Продаю Toyota Corolla 2018 года, 80.000 км, отличное состояние. Первый владелец, регулярное обслуживание, вся документация в наличии.', 1475000.00, 'used', 'active', 'Косте Мајинског 4, Ветерник, Сербия', 45.24755670, 19.76878366, 'Ветерник', 'Сербия', 0, '2025-02-07 17:33:27.680035', '2025-02-07 17:40:23.957971', true, 'ru');

SELECT setval('marketplace_listings_id_seq', 12, true);

-- Insert marketplace images
INSERT INTO marketplace_images (id, listing_id, file_path, file_name, file_size, content_type, is_main, created_at) VALUES
(15, 8, 'toyota_1.jpg', 'toyota_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(16, 8, 'toyota_2.jpg', 'toyota_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(17, 9, 'galaxy_s21_1.jpg', 'galaxy_s21_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(18, 9, 'galaxy_s21_2.jpg', 'galaxy_s21_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(19, 10, 'gaming_pc_1.jpg', 'gaming_pc_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(20, 10, 'gaming_pc_2.jpg', 'gaming_pc_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(21, 12, 'toyota_1.jpg', 'toyota_1.jpg', 454842, 'image/jpeg', true, '2025-02-07 17:35:09.579393'),
(22, 12, 'toyota_2.jpg', 'toyota_2.jpg', 398035, 'image/jpeg', true, '2025-02-07 17:40:24.397595');

SELECT setval('marketplace_images_id_seq', 22, true);

-- Insert reviews and related data
INSERT INTO reviews (id, user_id, entity_type, entity_id, rating, comment, pros, cons, photos, likes_count, is_verified_purchase, status, created_at, updated_at, helpful_votes, not_helpful_votes, original_language) VALUES
(1, 2, 'listing', 8, 5, 'норм', NULL, NULL, NULL, 0, true, 'published', '2025-02-07 07:47:17.001726', '2025-02-07 14:25:23.586871', 0, 1, 'ru');

SELECT setval('reviews_id_seq', 1, true);

INSERT INTO review_responses (id, review_id, user_id, response, created_at, updated_at) VALUES
(1, 1, 3, 'ok', '2025-02-07 07:49:14.935308', '2025-02-07 07:49:14.935308');

SELECT setval('review_responses_id_seq', 1, true);

INSERT INTO review_votes (review_id, user_id, vote_type, created_at) VALUES
(1, 3, 'not_helpful', '2025-02-07 07:48:11.709016');
