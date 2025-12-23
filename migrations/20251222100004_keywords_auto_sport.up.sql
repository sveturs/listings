-- +migrate Up
-- Добавляем meta_keywords для категорий Automobilizam и Sport i turizam
-- Дата: 2025-12-22

-- ========================================
-- AUTOMOBILIZAM - L2 категории
-- ========================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car batteries, accumulator, 12V battery, start stop battery, AGM, EFB, maintenance free',
    'sr', 'akumulatori za auto, akumulator 12V, start stop, bezodržavajući, punjivi, AGM, EFB',
    'ru', 'автомобильные аккумуляторы, 12V батарея, start stop, необслуживаемый, AGM, EFB'
) WHERE slug = 'akumulatori';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car tools, mechanic tools, torque wrench, car jack, socket set, oil filter wrench',
    'sr', 'alati za auto, mehaničarski alati, dizalica, ključevi, gedore, moment ključ, filter ključ',
    'ru', 'инструменты для авто, механический инструмент, домкрат, ключи, динамометрический ключ'
) WHERE slug = 'alati-za-automobile';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ambient lighting, car interior lights, footwell lights, LED strips, RGB lighting',
    'sr', 'ambijentalno svetlo, unutrašnje osvetljenje, LED trake auto, RGB svetla, atmosfera',
    'ru', 'амбиентная подсветка, внутреннее освещение авто, LED ленты, RGB подсветка'
) WHERE slug = 'ambijentalno-osvetljenje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car audio, GPS navigation, android head unit, car stereo, bluetooth, CarPlay, Android Auto',
    'sr', 'auto radio, navigacija, android multimedija, bluetooth, CarPlay, Android Auto, GPS',
    'ru', 'автомагнитола, навигация, android мультимедиа, bluetooth, CarPlay, Android Auto'
) WHERE slug = 'audio-i-navigacija';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car vacuum cleaner, portable vacuum, wet dry vacuum, handheld vacuum, 12V vacuum',
    'sr', 'auto usisivač, usisivač 12V, prenosivi usisivač, mokro suvo usisavanje',
    'ru', 'автомобильный пылесос, портативный пылесос, 12V пылесос, влажная уборка'
) WHERE slug = 'auto-aspiratori';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car accessories, car gadgets, air freshener, sun shade, steering wheel cover, seat cover',
    'sr', 'auto dodaci, miomiris, senilo, navlaka volana, navlaka sedišta, dodatna oprema',
    'ru', 'автоаксессуары, ароматизатор, солнцезащитная шторка, чехол руля, чехлы сидений'
) WHERE slug = 'auto-dodaci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car cosmetics, car polish, wax, car shampoo, tire shine, glass cleaner, detailing',
    'sr', 'auto kozmetika, polirna pasta, vosak, šampon, sjaj guma, sredstvo staklo, detailing',
    'ru', 'автокосметика, полироль, воск, шампунь, чернитель шин, очиститель стекол'
) WHERE slug = 'auto-kozmetika';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car organizer, trunk organizer, backseat organizer, console organizer, storage bag',
    'sr', 'auto organizator, organizator prtljažnika, organizator sedišta, torba za prtljažnik',
    'ru', 'автомобильный органайзер, органайзер багажника, органайзер сиденья, сумка'
) WHERE slug = 'auto-organizatori';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car cover, waterproof cover, UV protection, scratch protection, outdoor cover, indoor cover',
    'sr', 'pokrivač auta, cerada, zaštita od kiše, UV zaštita, zaštita od ogrebotina',
    'ru', 'чехол для автомобиля, водонепроницаемый, UV защита, защита от царапин'
) WHERE slug = 'auto-pokrivaci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car seat baby, child seat, infant seat, ISOFIX, group 0+, group 1, group 2/3, booster',
    'sr', 'auto sedište beba, dečije sedište, ISOFIX, grupa 0+, grupa 1, grupa 2/3, booster',
    'ru', 'детское автокресло, автолюлька, ISOFIX, группа 0+, группа 1, группа 2/3, бустер'
) WHERE slug = 'auto-sedista-bebe';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'dash cam, dashcam, car DVR, front rear camera, night vision, loop recording, G sensor',
    'sr', 'dash kamera, auto kamera, DVR, prednja zadnja kamera, noćni vid, snimanje, G senzor',
    'ru', 'видеорегистратор, dash камера, авторегистратор, двойная камера, ночная съёмка'
) WHERE slug = 'dash-kamere';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car parts, spare parts, brakes, filters, clutch, suspension, engine parts, transmission',
    'sr', 'auto delovi, rezervni delovi, kočnice, filteri, kvačilo, amortizeri, motor, menjač',
    'ru', 'автозапчасти, тормоза, фильтры, сцепление, амортизаторы, двигатель, коробка передач'
) WHERE slug = 'delovi-za-automobile';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'motorcycle parts, bike parts, chain, sprocket, brake pads, oil filter, spark plug',
    'sr', 'moto delovi, lanac, lančanik, diskovi, pločice, uljni filter, svećica',
    'ru', 'мотозапчасти, цепь, звёздочка, тормозные колодки, масляный фильтр, свеча'
) WHERE slug = 'delovi-za-motocikle';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'car phone holder, phone mount, magnetic holder, dashboard mount, vent mount, wireless charging',
    'sr', 'držač telefona auto, magnetni držač, nosač na tabla, držač ventilacija, bežično punjenje',
    'ru', 'держатель телефона авто, магнитный держатель, крепление на панель, в дефлектор'
) WHERE slug = 'drzaci-telefona-auto';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tuning, car tuning, performance parts, exhaust, air filter, ECU remap, turbo, intercooler',
    'sr', 'tuniranje, tuning delovi, auspuh, sportski filter, čipovanje, turbo, intercooler',
    'ru', 'тюнинг, выхлопная система, спортивный фильтр, чип-тюнинг, турбина, интеркулер'
) WHERE slug = 'tuniranje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'GPS tracker car, vehicle tracker, anti theft, real time tracking, OBD tracker, hardwired',
    'sr', 'GPS treker auto, praćenje vozila, antikrađa, lokator, OBD treker, ugrađeni treker',
    'ru', 'GPS трекер авто, маяк, противоугонная система, мониторинг, OBD трекер'
) WHERE slug = 'gps-trekeri-auto';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tires, car tires, summer tires, winter tires, all season, run flat, alloy wheels, rims',
    'sr', 'gume, auto gume, letnje gume, zimske gume, all season, run flat, felne, alu felge',
    'ru', 'шины, автошины, летние, зимние, всесезонные, run flat, диски, литые диски'
) WHERE slug = 'gume-i-felne';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'portable air conditioner, 12V car cooler, mini AC, car fan, cooling fan, summer',
    'sr', 'prenosiva klima, auto klima 12V, ventilator, rashladni uređaj, mini klima, letnji dodatak',
    'ru', 'портативный кондиционер, авто кондиционер 12V, вентилятор, охлаждение, мини AC'
) WHERE slug = 'klime-uređaji-prenosni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'LED strips car, underglow, interior LED, RGB strip, ambient light, decorative lighting',
    'sr', 'LED trake auto, RGB trake, unutrašnje svetlo, ambijentalno svetlo, dekorativno',
    'ru', 'LED ленты авто, RGB подсветка, внутренняя подсветка, декоративное освещение'
) WHERE slug = 'led-trake-auto';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'motorcycle gear, helmet, jacket, gloves, pants, boots, protectors, riding gear',
    'sr', 'moto oprema, kaciga, jakna, rukavice, pantalone, čizme, štitnici, oprema za vožnju',
    'ru', 'мотоэкипировка, шлем, куртка, перчатки, штаны, ботинки, защита, экипировка'
) WHERE slug = 'moto-oprema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'parking sensors, PDC, rear sensors, front sensors, buzzer, LED display, parking assist',
    'sr', 'parking senzori, senzori za parking, zadnji senzori, prednji senzori, zvučna signalizacija',
    'ru', 'парктроник, датчики парковки, задние датчики, передние датчики, звуковой сигнал'
) WHERE slug = 'parking-senzori';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'EV charger, electric car charger, Type 2, CCS, CHAdeMO, home charger, portable charger',
    'sr', 'punjač električnih vozila, EV punjač, Type 2, CCS, kućni punjač, prenosivi punjač',
    'ru', 'зарядка электромобиля, EV зарядка, Type 2, CCS, домашняя зарядка, портативная'
) WHERE slug = 'punjaci-elektricni-auto';

-- ========================================
-- SPORT I TURIZAM - L2 категории
-- ========================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'badminton, racket, shuttlecock, net, badminton shoes, badminton bag, grip tape',
    'sr', 'badminton, reket, loptica, mreža, badminton patike, torba, grip traka',
    'ru', 'бадминтон, ракетка, волан, сетка, обувь для бадминтона, сумка, намотка'
) WHERE slug = 'badminton-oprema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bike, bicycle, scooter, electric scooter, mountain bike, road bike, BMX, kids bike',
    'sr', 'bicikl, trotinet, električni trotinet, brdski bicikl, drumski, BMX, dečiji bicikl',
    'ru', 'велосипед, самокат, электросамокат, горный велосипед, шоссейный, BMX, детский'
) WHERE slug = 'bicikli-i-trotineti';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'martial arts protection, shin guards, headgear, mouth guard, groin protector, chest guard',
    'sr', 'zaštita borilačke veštine, štitnici za noge, kaciga, štitnik za zube, suspenzor, prsluk',
    'ru', 'защита единоборства, щитки, шлем, капа, бандаж, защита корпуса'
) WHERE slug = 'borilacke-vestine-zastita';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'jongleren, juggling balls, devil sticks, poi, clubs, rings, diabolo, juggling equipment',
    'sr', 'džonovanje, loptice, đavolje štapove, poi, keglje, prstenovi, diabolo, oprema',
    'ru', 'жонглирование, мячи, стики, пои, кегли, кольца, диаболо, жонглёрские принадлежности'
) WHERE slug = 'dzonovanje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'fitness, gym equipment, dumbbells, kettlebell, barbell, weight plates, bench, rack',
    'sr', 'fitnes, teretana, bučice, girje, šipka, tegovi, klupa, stalak, sprave',
    'ru', 'фитнес, тренажёрный зал, гантели, гири, штанга, диски, скамья, стойка'
) WHERE slug = 'fitnes-i-teretana';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'football, soccer ball, cleats, shin guards, goalkeeper gloves, football jersey, training',
    'sr', 'fudbal, lopta, kopačke, štitnici, golmanske rukavice, dres, trening',
    'ru', 'футбол, мяч, бутсы, щитки, вратарские перчатки, форма, тренировка'
) WHERE slug = 'fudbal';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'yoga blocks, foam blocks, cork blocks, yoga props, support, stretching, flexibility',
    'sr', 'joga blokovi, blokovi od pene, pluteni blokovi, rekviziti, podrška, istezanje',
    'ru', 'блоки для йоги, пеноблоки, пробковые блоки, пропсы, поддержка, растяжка'
) WHERE slug = 'joga-blokovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kayak, canoe, paddle, kayak accessories, inflatable kayak, sit on top, touring kayak',
    'sr', 'kajak, kanu, veslo, kajak oprema, naduvani kajak, touring, oprema za veslanje',
    'ru', 'каяк, каноэ, весло, надувная байдарка, туристический каяк, аксессуары'
) WHERE slug = 'kajak-kanu-oprema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'camping, tent, sleeping bag, camping chair, camp stove, cooler, camping gear, outdoor',
    'sr', 'kampovanje, šator, spavaća vreća, stolica, roštilj, frižider, kamp oprema, priroda',
    'ru', 'кемпинг, палатка, спальный мешок, кресло, горелка, холодильник, походное снаряжение'
) WHERE slug = 'kampovanje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'basketball, basketball ball, basketball shoes, jersey, shorts, hoop, backboard, net',
    'sr', 'košarka, košarkaška lopta, patike, dres, šorc, koš, tabla, mreža',
    'ru', 'баскетбол, баскетбольный мяч, кроссовки, майка, шорты, кольцо, щит, сетка'
) WHERE slug = 'kosarka';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hunting, rifle, shotgun, ammunition, hunting knife, camouflage, hunting gear, optics',
    'sr', 'lov, puška, sačmara, municija, lovački nož, kamuflaža, oprema, optika',
    'ru', 'охота, ружьё, дробовик, патроны, охотничий нож, камуфляж, снаряжение, оптика'
) WHERE slug = 'lov';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'pilates ball, exercise ball, stability ball, swiss ball, gym ball, pregnancy ball',
    'sr', 'pilates lopta, lopta vežbanje, stabilna lopta, swiss lopta, gimnastička lopta',
    'ru', 'мяч пилатес, фитбол, гимнастический мяч, швейцарский мяч, мяч для беременных'
) WHERE slug = 'pilates-lopte';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hiking, trekking, backpack, hiking boots, trekking poles, outdoor gear, mountain',
    'sr', 'planinarenje, trekking, ranac, planinske cipele, štapovi, oprema, planina',
    'ru', 'походы, треккинг, рюкзак, трекинговые ботинки, палки, снаряжение, горы'
) WHERE slug = 'planinarenje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'swimming, swimsuit, goggles, swim cap, fins, training equipment, pool accessories',
    'sr', 'plivanje, kupaći, naočare, kapa, peraje, trening oprema, bazen',
    'ru', 'плавание, купальник, очки, шапочка, ласты, тренировочное снаряжение, бассейн'
) WHERE slug = 'plivanje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'fishing, rod, reel, tackle, bait, lures, fishing gear, hooks, line, landing net',
    'sr', 'ribolov, štap, mašinica, mamci, oprema, udice, konac, podmetač, ribolovna oprema',
    'ru', 'рыбалка, удочка, катушка, снасти, приманки, крючки, леска, подсачек'
) WHERE slug = 'ribolov';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'diving, scuba diving, wetsuit, dive computer, regulator, BCD, fins, mask, snorkel',
    'sr', 'ronjenje, ronilačka oprema, odelo, kompjuter, regulator, BCD, peraje, maska',
    'ru', 'дайвинг, гидрокостюм, компьютер, регулятор, BCD, ласты, маска, трубка'
) WHERE slug = 'ronilacka-oprema-profesionalna';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'table tennis, ping pong paddle, racket, balls, table tennis table, net, rubber, blade',
    'sr', 'stoni tenis, ping pong reket, loptica, sto, mreža, guma, drvena osnova',
    'ru', 'настольный теннис, ракетка пинг-понг, мячик, стол, сетка, накладка, основание'
) WHERE slug = 'stoni-tenis-reketi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'surfing, surfboard, longboard, shortboard, bodyboard, surf wax, leash, fins, wetsuit',
    'sr', 'surfovanje, daska, longboard, shortboard, bodyboard, vosak, povodac, peraje',
    'ru', 'сёрфинг, доска для серфа, лонгборд, шортборд, бодиборд, воск, лиш, плавники'
) WHERE slug = 'surfovanje-daske';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tennis, tennis racket, tennis balls, tennis shoes, tennis bag, strings, grip, court',
    'sr', 'tenis, teniski reket, loptice, patike, torba, žice, grip, teren',
    'ru', 'теннис, теннисная ракетка, мячи, кроссовки, сумка, струны, намотка, корт'
) WHERE slug = 'tenis';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'professional tennis racket, head heavy, wilson pro staff, RF97, speed, prestige, ultra',
    'sr', 'profesionalni tenis reket, head heavy, wilson pro staff, RF97, speed, prestige',
    'ru', 'профессиональная теннисная ракетка, head heavy, wilson pro staff, RF97'
) WHERE slug = 'tenis-reketi-pro';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'resistance bands, exercise bands, stretch bands, loop bands, therapy bands, latex bands',
    'sr', 'trake istezanje, trake vežbanje, elastične trake, loop trake, rehabilitacione trake',
    'ru', 'эспандеры, ленты для упражнений, резинки для фитнеса, петли, латексные ленты'
) WHERE slug = 'trake-istezanje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'winter sports, skiing, snowboard, ski boots, ski poles, snowboard boots, bindings',
    'sr', 'zimski sportovi, skije, snowboard, ski cipele, štapovi, vezovi, ski oprema',
    'ru', 'зимние виды спорта, лыжи, сноуборд, ботинки, палки, крепления, снаряжение'
) WHERE slug = 'zimski-sportovi';

-- ========================================
-- KAMPOVANJE - L3 категории (дочерние для kampovanje)
-- ========================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'camping stove, portable gas stove, butane stove, propane stove, camp cooker, burner',
    'sr', 'kamp rešo, prenosivi rešo, gasni rešo, butan rešo, roštilj, gorionik',
    'ru', 'кемпинговая горелка, газовая плита, туристическая горелка, портативная плита'
) WHERE slug = 'gas-reaud';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'camping lamp, lantern, LED lantern, rechargeable, battery powered, solar lamp, headlamp',
    'sr', 'kamp lampa, fenjer, LED lampa, punjiva, baterijska, solarna lampa, čeona lampa',
    'ru', 'кемпинговая лампа, фонарь, светодиодный фонарь, аккумуляторный, солнечная лампа'
) WHERE slug = 'kamp-lampa';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'camping backpack, hiking backpack, 40L, 60L, 80L, trekking backpack, outdoor backpack',
    'sr', 'kamp ranac, planinski ranac, 40L, 60L, 80L, trekking ranac, outdoor ranac',
    'ru', 'туристический рюкзак, походный рюкзак, 40L, 60L, 80L, треккинговый рюкзак'
) WHERE slug = 'kamp-ranac';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'camping grill, portable BBQ, charcoal grill, gas grill, folding grill, outdoor BBQ',
    'sr', 'kamp roštilj, prenosivi roštilj, roštilj na ugalj, gasni, sklopivi, roštilj priroda',
    'ru', 'кемпинговый гриль, портативный мангал, угольный гриль, газовый, складной'
) WHERE slug = 'kamp-roštilj';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'camping chair, folding chair, lightweight chair, reclining chair, camp stool, outdoor chair',
    'sr', 'kamp stolica, sklopiva stolica, lagana stolica, stolica sa naslonom, outdoor stolica',
    'ru', 'кемпинговое кресло, складной стул, раскладное кресло, туристический стул'
) WHERE slug = 'kamp-stolica';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'portable fridge, camping cooler, 12V fridge, electric cooler, car fridge, ice box',
    'sr', 'prenosivi frižider, kamp frižider, 12V frižider, električni, auto frižider, kutija',
    'ru', 'портативный холодильник, автохолодильник, 12V холодильник, сумка холодильник'
) WHERE slug = 'prenosivi-frižider';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'self inflating mat, sleeping pad, camping mat, air mat, foam mat, insulated mat',
    'sr', 'samoduvajući dušek, podloga spavanje, kamp dušek, vazdušni dušek, izolaciona podloga',
    'ru', 'самонадувающийся коврик, спальный коврик, туристический коврик, каремат'
) WHERE slug = 'samoduvajuci-dušek';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tent 2-3 person, small tent, camping tent, lightweight tent, waterproof, dome tent',
    'sr', 'šator 2-3 osobe, mali šator, kamp šator, lagan šator, vodootporan, kupola šator',
    'ru', 'палатка 2-3 человека, маленькая палатка, лёгкая палатка, водонепроницаемая'
) WHERE slug = 'sator-2-3-osobe';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tent 4-6 person, family tent, large tent, camping tent, waterproof, tunnel tent, cabin tent',
    'sr', 'šator 4-6 osoba, porodični šator, veliki šator, kamp šator, vodootporan, tunel šator',
    'ru', 'палатка 4-6 человек, семейная палатка, большая палатка, водонепроницаемая'
) WHERE slug = 'sator-4-6-osoba';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'sleeping bag, mummy sleeping bag, rectangular bag, winter bag, summer bag, down bag',
    'sr', 'spavaća vreća, zimska vreća, letnja vreća, mumija vreća, pravougaona, puh vreća',
    'ru', 'спальный мешок, зимний мешок, летний мешок, кокон, пуховый мешок'
) WHERE slug = 'spavaca-vreća';
