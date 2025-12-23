-- Заполнение meta_keywords для L1 категорий
-- Формат: {"en": "keywords", "sr": "keywords", "ru": "keywords"}

-- 1. Elektronika
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'electronics, phone, smartphone, tablet, laptop, computer, gadget, tech, technology, mobile',
    'sr', 'elektronika, telefon, smartphone, tablet, laptop, računar, računari, gadget, tehnologija, mobilni',
    'ru', 'электроника, телефон, смартфон, планшет, ноутбук, компьютер, гаджет, технология, мобильный'
)
WHERE slug = 'elektronika' AND level = 1;

-- 2. Odeca i obuca
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'clothing, clothes, fashion, shirt, dress, pants, shoes, sneakers, boots, jacket, coat',
    'sr', 'odeća, obuća, moda, košulja, haljina, pantalone, cipele, patike, čizme, jakna, kaput',
    'ru', 'одежда, обувь, мода, рубашка, платье, брюки, туфли, кроссовки, сапоги, куртка, пальто'
)
WHERE slug = 'odeca-i-obuca' AND level = 1;

-- 3. Dom i basta
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'home, garden, furniture, decor, decoration, tool, repair, renovation, household, kitchen',
    'sr', 'dom, bašta, nameštaj, dekor, dekoracija, alat, popravka, renoviranje, kućanstvo, kuhinja',
    'ru', 'дом, сад, мебель, декор, украшение, инструмент, ремонт, реновация, хозяйство, кухня'
)
WHERE slug = 'dom-i-basta' AND level = 1;

-- 4. Lepota i zdravlje
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'beauty, health, cosmetics, makeup, skincare, perfume, care, wellness, hygiene',
    'sr', 'lepota, zdravlje, kozmetika, šminka, nega, parfem, briga, zdravlje, higijena',
    'ru', 'красота, здоровье, косметика, макияж, уход, духи, забота, благополучие, гигиена'
)
WHERE slug = 'lepota-i-zdravlje' AND level = 1;

-- 5. Za bebe i decu
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'baby, kids, children, toy, toys, stroller, crib, diaper, clothing, game',
    'sr', 'bebe, deca, dete, igračka, igračke, kolica, krevetac, pelena, odeća, igra',
    'ru', 'дети, ребёнок, игрушка, игрушки, коляска, кроватка, подгузник, одежда, игра'
)
WHERE slug = 'za-bebe-i-decu' AND level = 1;

-- 6. Sport i turizam
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'sport, fitness, gym, exercise, equipment, outdoor, camping, hiking, bicycle, bike',
    'sr', 'sport, fitnes, teretana, vežbanje, oprema, outdoor, kampovanje, pešačenje, bicikl, bicikla',
    'ru', 'спорт, фитнес, тренажерный зал, упражнение, оборудование, на открытом воздухе, кемпинг, походы, велосипед'
)
WHERE slug = 'sport-i-turizam' AND level = 1;

-- 7. Automobilizam
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'car, vehicle, automobile, auto, parts, accessories, tire, wheel, motor, motorcycle',
    'sr', 'automobil, auto, vozilo, delovi, oprema, guma, točak, motor, motocikl',
    'ru', 'автомобиль, авто, транспорт, запчасти, аксессуары, шина, колесо, мотор, мотоцикл'
)
WHERE slug = 'automobilizam' AND level = 1;

-- 8. Kucni aparati
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'appliances, home appliances, kitchen appliances, fridge, washing machine, oven, dishwasher, vacuum',
    'sr', 'kućni aparati, kućanski aparati, aparati za kuhinju, frižider, veš mašina, rerna, mašina za sudove, usisivač',
    'ru', 'бытовая техника, кухонная техника, холодильник, стиральная машина, духовка, посудомоечная машина, пылесос'
)
WHERE slug = 'kucni-aparati' AND level = 1;

-- 9. Nakit i satovi
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'jewelry, jewellery, watch, watches, ring, necklace, earrings, bracelet, gold, silver',
    'sr', 'nakit, sat, satovi, prsten, ogrlica, minđuše, narukvica, zlato, srebro',
    'ru', 'украшения, часы, кольцо, ожерелье, серьги, браслет, золото, серебро'
)
WHERE slug = 'nakit-i-satovi' AND level = 1;

-- 10. Knjige i mediji
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'book, books, media, music, movie, cd, dvd, vinyl, magazine, newspaper',
    'sr', 'knjiga, knjige, mediji, muzika, film, cd, dvd, vinil, časopis, novine',
    'ru', 'книга, книги, медиа, музыка, фильм, диск, журнал, газета'
)
WHERE slug = 'knjige-i-mediji' AND level = 1;

-- 11. Hrana i pice
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'food, drink, beverage, grocery, snacks, coffee, tea, wine, beer, organic',
    'sr', 'hrana, piće, namirnice, grickalice, kafa, čaj, vino, pivo, organski',
    'ru', 'еда, напиток, продукты, закуски, кофе, чай, вино, пиво, органический'
)
WHERE slug = 'hrana-i-pice' AND level = 1;

-- 12. Industrija i alati
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'industry, industrial, tools, equipment, machinery, construction, professional, workshop',
    'sr', 'industrija, industrijski, alati, oprema, mašine, građevina, profesionalni, radionica',
    'ru', 'промышленность, промышленный, инструменты, оборудование, машины, строительство, профессиональный, мастерская'
)
WHERE slug = 'industrija-i-alati' AND level = 1;

-- 13. Kancelarijski materijal
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'office, supplies, stationery, paper, pen, printer, desk, chair, organization',
    'sr', 'kancelarija, materijal, papir, olovka, štampač, sto, stolica, organizacija',
    'ru', 'офис, канцелярия, бумага, ручка, принтер, стол, стул, организация'
)
WHERE slug = 'kancelarijski-materijal' AND level = 1;

-- 14. Kucni ljubimci
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'pet, pets, dog, cat, animal, food, toys, accessories, care, aquarium',
    'sr', 'kućni ljubimac, ljubimci, pas, mačka, životinja, hrana, igračke, oprema, nega, akvarijum',
    'ru', 'домашнее животное, питомец, собака, кошка, корм, игрушки, аксессуары, уход, аквариум'
)
WHERE slug = 'kucni-ljubimci' AND level = 1;

-- 15. Muzicki instrumenti
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'music, musical instruments, instrument, guitar, piano, drums, violin, synthesizer, microphone',
    'sr', 'muzika, muzički instrumenti, instrument, gitara, klavir, bubnjevi, violina, sintisajzer, mikrofon',
    'ru', 'музыка, музыкальные инструменты, инструмент, гитара, пианино, барабаны, скрипка, синтезатор, микрофон'
)
WHERE slug = 'muzicki-instrumenti' AND level = 1;

-- 16. Ostalo
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'other, miscellaneous, various, mixed, general, different, assorted',
    'sr', 'ostalo, razno, različito, mešovito, opšte, raznovrsno',
    'ru', 'другое, разное, различное, смешанное, общее, разнообразное'
)
WHERE slug = 'ostalo' AND level = 1;

-- 17. Umetnost i rukotvorine
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'art, craft, handmade, painting, sculpture, artwork, creative, artist, handcraft',
    'sr', 'umetnost, rukotvorine, ručni rad, slika, skulptura, umetnička dela, kreativno, umetnik',
    'ru', 'искусство, ремесло, ручная работа, картина, скульптура, произведение искусства, творчество, художник'
)
WHERE slug = 'umetnost-i-rukotvorine' AND level = 1;

-- 18. Usluge
UPDATE categories
SET meta_keywords = jsonb_build_object(
    'en', 'services, service, professional, help, assistance, repair, maintenance, consulting',
    'sr', 'usluge, usluga, profesionalne, pomoć, asistencija, popravka, održavanje, konsalting',
    'ru', 'услуги, сервис, профессиональный, помощь, ремонт, обслуживание, консалтинг'
)
WHERE slug = 'usluge' AND level = 1;

-- Обновить updated_at для всех изменённых записей
UPDATE categories
SET updated_at = NOW()
WHERE level = 1 AND meta_keywords IS NOT NULL AND meta_keywords != '{}'::jsonb;
