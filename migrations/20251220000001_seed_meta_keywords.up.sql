-- Migration: Seed meta_keywords for all categories
-- Description: Add SEO-optimized meta keywords in Serbian, English, and Russian for all 172 categories

-- Level 1 Categories

-- Automobilizam
UPDATE categories SET meta_keywords = '{
  "sr": "automobilizam, auto delovi, vozila, auto oprema, auto servis, auto kozmetika, gume, akumulatori, auto elektronika",
  "en": "automotive, auto parts, vehicles, car equipment, car service, auto cosmetics, tires, batteries, car electronics",
  "ru": "автотовары, автозапчасти, транспорт, автооборудование, автосервис, автокосметика, шины, аккумуляторы, автоэлектроника"
}'::jsonb WHERE id = '00a848bb-27c9-4598-b6e0-68c02fa1a797';

-- Događaji i karte
UPDATE categories SET meta_keywords = '{
  "sr": "događaji, karte, koncerti, predstave, festival, bioskop, pozorište, sport događaji, ulaznice",
  "en": "events, tickets, concerts, shows, festivals, cinema, theater, sports events, admission",
  "ru": "события, билеты, концерты, представления, фестивали, кино, театр, спортивные события, входные билеты"
}'::jsonb WHERE id = 'a1b2c3d4-1001-4001-8001-000000001020';

-- Dom i bašta
UPDATE categories SET meta_keywords = '{
  "sr": "dom, bašta, nameštaj, dekoracija, vrt, baštenska oprema, alati, rasveta, tekstil",
  "en": "home, garden, furniture, decoration, yard, garden equipment, tools, lighting, textiles",
  "ru": "дом, сад, мебель, декор, двор, садовое оборудование, инструменты, освещение, текстиль"
}'::jsonb WHERE id = '556bce4b-ec45-4fc3-ba7a-ed9aa5c9475d';

-- Elektronika
UPDATE categories SET meta_keywords = '{
  "sr": "elektronika, tehnika, gadžeti, uređaji, telefoni, računari, tv, audio, tableti, konzole",
  "en": "electronics, tech, gadgets, devices, phones, computers, tv, audio, tablets, consoles",
  "ru": "электроника, техника, гаджеты, устройства, телефоны, компьютеры, тв, аудио, планшеты, консоли"
}'::jsonb WHERE id = '2014ea09-7151-46a0-9f78-d4caa33a3833';

-- Hobiji i zabava
UPDATE categories SET meta_keywords = '{
  "sr": "hobiji, zabava, igre, kolekcionarstvo, knjige, muzika, sport, outdoor, gaming, kreativnost",
  "en": "hobbies, entertainment, games, collecting, books, music, sports, outdoor, gaming, creativity",
  "ru": "хобби, развлечения, игры, коллекционирование, книги, музыка, спорт, активный отдых, геймин г, творчество"
}'::jsonb WHERE id = 'a1b2c3d4-1001-4001-8001-000000001015';

-- Hrana i piće
UPDATE categories SET meta_keywords = '{
  "sr": "hrana, piće, namirnice, delikatesi, kafa, čaj, sokovi, alkohol, organska hrana, supleme nti",
  "en": "food, beverages, groceries, delicacies, coffee, tea, juices, alcohol, organic food, supplements",
  "ru": "продукты, напитки, продовольствие, деликатесы, кофе, чай, соки, алкоголь, органическая еда, добавки"
}'::jsonb WHERE id = '53c99f20-8b91-482a-9445-99fa7c481051';

-- Industrija i alati
UPDATE categories SET meta_keywords = '{
  "sr": "industrija, alati, oprema, mašine, električni alati, ručni alati, powertools, radni alati, građevina",
  "en": "industrial, tools, equipment, machines, power tools, hand tools, powertools, work tools, construction",
  "ru": "промтовары, инструменты, оборудование, машины, электроинструменты, ручные инструменты, строительство"
}'::jsonb WHERE id = '2cc9a6dd-c41a-4e9e-a7cc-8695447f9e36';

-- Kancelarijski materijal
UPDATE categories SET meta_keywords = '{
  "sr": "kancelarija, kancelarijski materijal, papir, olovke, fascikle, organizacija, školski pribor, pisaći pribor",
  "en": "office, office supplies, paper, pens, folders, organization, school supplies, stationery",
  "ru": "канцтовары, офис, бумага, ручки, папки, организация, школьные принадлежности, канцелярия"
}'::jsonb WHERE id = '533f1842-1efd-4870-afbe-47761883f9b3';

-- Knjige i mediji
UPDATE categories SET meta_keywords = '{
  "sr": "knjige, mediji, literatura, udžbenici, stripovi, časopisi, e-knjige, audio knjige, izdavaštvo",
  "en": "books, media, literature, textbooks, comics, magazines, e-books, audiobooks, publishing",
  "ru": "книги, медиа, литература, учебники, комиксы, журналы, электронные книги, аудиокниги, издательство"
}'::jsonb WHERE id = 'f6260278-7bb7-4507-9625-c77e383ff975';

-- Kućni aparati
UPDATE categories SET meta_keywords = '{
  "sr": "kućni aparati, bela tehnika, frižideri, veš mašine, šporet, usisivači, mikr otalasne, klima",
  "en": "appliances, white goods, refrigerators, washing machines, stoves, vacuum cleaners, microwaves, air conditioning",
  "ru": "бытовая техника, холодильники, стиральные машины, плиты, пылесосы, микроволновки, кондиционеры"
}'::jsonb WHERE id = '5bb29612-c67b-4989-8f27-4696f4c04f76';

-- Kućni ljubimci
UPDATE categories SET meta_keywords = '{
  "sr": "kućni ljubimci, psi, mačke, hrana za ljubimce, igračke, aksesoari, akvarijumi, ptice, glodari",
  "en": "pets, dogs, cats, pet food, toys, accessories, aquariums, birds, rodents",
  "ru": "зоотовары, собаки, кошки, корм для животных, игрушки, аксессуары, аквариумы, птицы, грызуны"
}'::jsonb WHERE id = '7b11cb87-8590-48a6-82ce-e14766063fe5';

-- Lepota i zdravlje
UPDATE categories SET meta_keywords = '{
  "sr": "lepota, zdravlje, kozmetika, nega kože, nega kose, parfemi, makeup, higijena, wellness",
  "en": "beauty, health, cosmetics, skincare, hair care, perfumes, makeup, hygiene, wellness",
  "ru": "красота, здоровье, косметика, уход за кожей, уход за волосами, парфюмерия, макияж, гигиена"
}'::jsonb WHERE id = '8dad6570-6e10-499c-a167-2cc07fa62a5b';

-- Muzički instrumenti
UPDATE categories SET meta_keywords = '{
  "sr": "muzički instrumenti, gitara, klavir, bubnjevi, bas, sintisajzer, harmonika, violina, orkestar",
  "en": "musical instruments, guitar, piano, drums, bass, synthesizer, accordion, violin, orchestra",
  "ru": "музыкальные инструменты, гитара, пианино, барабаны, бас, синтезатор, аккордеон, скрипка"
}'::jsonb WHERE id = '6cabb351-fe1a-4168-abd5-b4f2e88c1f78';

-- Nakit i satovi
UPDATE categories SET meta_keywords = '{
  "sr": "nakit, satovi, zlato, srebro, dijamanti, burme, minđuše, ogrlice, prstenje, luksuz",
  "en": "jewelry, watches, gold, silver, diamonds, wedding rings, earrings, necklaces, rings, luxury",
  "ru": "украшения, часы, золото, серебро, бриллианты, обручальные кольца, серьги, ожерелья, роскошь"
}'::jsonb WHERE id = '22e97c40-1400-45c8-8042-e3ed328e2d5b';

-- Nekretnine
UPDATE categories SET meta_keywords = '{
  "sr": "nekretnine, stanovi, kuće, zemljište, izdavanje, prodaja, zakup, poslovni prostor, investment",
  "en": "real estate, apartments, houses, land, rental, sale, lease, commercial space, investment",
  "ru": "недвижимость, квартиры, дома, земля, аренда, продажа, коммерческая недвижимость, инвестиции"
}'::jsonb WHERE id = 'a1b2c3d4-1001-4001-8001-000000001004';

-- Obrazovanje
UPDATE categories SET meta_keywords = '{
  "sr": "obrazovanje, kursevi, obuka, skola, fakultet, online učenje, tutorijali, certifikati, edukacija",
  "en": "education, courses, training, school, university, online learning, tutorials, certificates",
  "ru": "образование, курсы, обучение, школа, университет, онлайн обучение, учебники, сертификаты"
}'::jsonb WHERE id = 'a1b2c3d4-1001-4001-8001-000000001019';

-- Odeća i obuća
UPDATE categories SET meta_keywords = '{
  "sr": "odeća, obuća, moda, garderoba, patike, cipele, jakne, pantalone, haljine, majice",
  "en": "clothing, footwear, fashion, wardrobe, sneakers, shoes, jackets, pants, dresses, shirts",
  "ru": "одежда, обувь, мода, гардероб, кроссовки, туфли, куртки, брюки, платья, футболки"
}'::jsonb WHERE id = '41ddcc2f-10e9-4486-99e3-649f582e8a18';

-- Ostalo
UPDATE categories SET meta_keywords = '{
  "sr": "ostalo, razno, nekategorisano, mešovito, opšte, različito, sve ostalo",
  "en": "other, miscellaneous, uncategorized, mixed, general, various, everything else",
  "ru": "прочее, разное, без категории, смешанное, общее, разнообразное"
}'::jsonb WHERE id = '49f38b62-767b-4882-8b39-b7d1b13297bc';

-- Poljoprivreda
UPDATE categories SET meta_keywords = '{
  "sr": "poljoprivreda, farma, ratarstvo, stočarstvo, voćarstvo, povrtarstvo, traktor, oprema, setva",
  "en": "agriculture, farm, crop farming, livestock, fruit farming, vegetable farming, tractor, equipment, planting",
  "ru": "сельское хозяйство, ферма, растениеводство, животноводство, садоводство, овощеводство, трактор"
}'::jsonb WHERE id = 'a1b2c3d4-1001-4001-8001-000000001006';

-- Poslovi
UPDATE categories SET meta_keywords = '{
  "sr": "poslovi, zapošljavanje, karijera, posao, konkursi, rad, vacancy, job offers, employment",
  "en": "jobs, employment, career, work, vacancies, recruitment, job offers, hiring",
  "ru": "работа, вакансии, карьера, трудоустройство, найм, предложения работы, рекрутинг"
}'::jsonb WHERE id = 'a1b2c3d4-1001-4001-8001-000000001018';

-- Prirodni materijali
UPDATE categories SET meta_keywords = '{
  "sr": "prirodni materijali, drvo, kamen, glina, bambus, ekologija, održivost, organski",
  "en": "natural materials, wood, stone, clay, bamboo, ecology, sustainability, organic",
  "ru": "природные материалы, дерево, камень, глина, бамбук, экология, устойчивость, органика"
}'::jsonb WHERE id = 'a1b2c3d4-1001-4001-8001-000000010207';

-- Sport i turizam
UPDATE categories SET meta_keywords = '{
  "sr": "sport, turizam, fitnes, outdoor, planinarenje, kampovanje, skijanje, fudbal, košarka",
  "en": "sports, tourism, fitness, outdoor, hiking, camping, skiing, football, basketball",
  "ru": "спорт, туризм, фитнес, активный отдых, пеший туризм, кемпинг, лыжи, футбол, баскетбол"
}'::jsonb WHERE id = 'f24e428c-e7eb-449b-8738-68eda1a677be';

-- Umetnost i rukotvorine
UPDATE categories SET meta_keywords = '{
  "sr": "umetnost, rukotvorine, handmade, zanati, slike, skulpture, kreativnost, DIY, craft",
  "en": "art, crafts, handmade, handicrafts, paintings, sculptures, creativity, DIY, craft",
  "ru": "искусство, рукоделие, ручная работа, ремесла, картины, скульптуры, творчество"
}'::jsonb WHERE id = 'cdddd6c5-9f12-4fd1-81f4-d7d74eed2aaf';

-- Usluge
UPDATE categories SET meta_keywords = '{
  "sr": "usluge, servisi, popravke, majstori, čišćenje, transport, dostava, montaža, održavanje",
  "en": "services, repairs, handymen, cleaning, transport, delivery, assembly, maintenance",
  "ru": "услуги, ремонт, мастера, уборка, транспорт, доставка, монтаж, обслуживание"
}'::jsonb WHERE id = '641c0b94-b993-4e02-aee7-0cf330a2485b';

-- Za bebe i decu
UPDATE categories SET meta_keywords = '{
  "sr": "bebe, deca, igračke, dečija odeća, krevetići, kolica, hranili ce, dečija oprema",
  "en": "babies, kids, toys, children clothing, cribs, strollers, feeders, baby gear",
  "ru": "младенцы, дети, игрушки, детская одежда, кроватки, коляски, питание, детское оборудование"
}'::jsonb WHERE id = '9848acbd-8e19-4232-8de5-aa0b93146113';

-- Level 2 Categories

-- Alati i popravke
UPDATE categories SET meta_keywords = '{
  "sr": "alati, popravke, majstorski alat, powertools, električni alati, ručni alati, bušilice, brusilice",
  "en": "tools, repairs, hand tools, powertools, electric tools, drills, grinders, workshop",
  "ru": "инструменты, ремонт, ручной инструмент, электроинструменты, дрели, болгарки, мастерская"
}'::jsonb WHERE id = '9f909c82-8582-4dad-9174-3396887da0ba';

-- Alati za baštu
UPDATE categories SET meta_keywords = '{
  "sr": "alati za baštu, lopata, grablje, motika, kosačica, trimer, bašta, vrt, baštovanstvo",
  "en": "garden tools, shovel, rake, hoe, lawnmower, trimmer, garden, yard, gardening",
  "ru": "садовый инвентарь, лопата, грабли, мотыга, газонокосилка, триммер, сад, садоводство"
}'::jsonb WHERE id = '551676da-9f15-49d6-94cd-911f1664695e';

-- Anti-aging
UPDATE categories SET meta_keywords = '{
  "sr": "anti-aging, protiv starenja, kreme, serumi, kolagen, bore, nega lica, mladost, lifting",
  "en": "anti-aging, anti-wrinkle, creams, serums, collagen, wrinkles, facial care, youth, lifting",
  "ru": "антивозрастной, против старения, кремы, сыворотки, коллаген, морщины, уход за лицом"
}'::jsonb WHERE id = '8505b35c-8138-4fda-a902-ae35e399bb3e';

-- Audio oprema
UPDATE categories SET meta_keywords = '{
  "sr": "audio oprema, zvučnici, slušalice, amplifier, receiver, subwoofer, soundbar, hi-fi, stereo",
  "en": "audio equipment, speakers, headphones, amplifier, receiver, subwoofer, soundbar, hi-fi, stereo",
  "ru": "аудио техника, колонки, наушники, усилитель, ресивер, сабвуфер, саундбар, стерео"
}'::jsonb WHERE id = 'c819ede8-a132-4232-add9-f3ba6eeec097';

-- Baštenka garnitura
UPDATE categories SET meta_keywords = '{
  "sr": "baštenka garnitura, spoljni nameštaj, stolice, sto, terasa, vrt, ratan, drvo, aluminijum",
  "en": "outdoor furniture, garden set, chairs, table, terrace, patio, rattan, wood, aluminum",
  "ru": "садовая мебель, комплект, стулья, стол, терраса, патио, ротанг, дерево, алюминий"
}'::jsonb WHERE id = 'c416df52-781f-41ad-8faa-95ecac43a819';

-- Baštenka oprema
UPDATE categories SET meta_keywords = '{
  "sr": "baštenka oprema, creva, prskalice, saksije, đubrivo, navodnjavanje, vrtna oprema",
  "en": "garden equipment, hoses, sprinklers, pots, fertilizer, irrigation, yard equipment",
  "ru": "садовое оборудование, шланги, разбрызгиватели, горшки, удобрение, полив"
}'::jsonb WHERE id = 'd863707b-aefd-4c8a-90bc-73bcc6fbccac';

-- Baštenka rasveta
UPDATE categories SET meta_keywords = '{
  "sr": "baštenka rasveta, solarne lampe, LED, svetla, rasveta dvorišta, outdoor lighting, ambijent",
  "en": "garden lighting, solar lamps, LED, lights, yard lighting, outdoor lighting, ambience",
  "ru": "садовое освещение, солнечные лампы, LED, светильники, дворовое освещение, амбиент"
}'::jsonb WHERE id = '3b021fc8-ade9-4d2d-a082-c338f84e699a';

-- Baštenke ukrase
UPDATE categories SET meta_keywords = '{
  "sr": "baštenke ukrase, vrtne figure, fontane, dekoracija, skulpture, kućice za ptice, ukrasni elementi",
  "en": "garden decorations, garden figures, fountains, decoration, sculptures, birdhouses, ornaments",
  "ru": "садовые украшения, садовые фигуры, фонтаны, декор, скульптуры, скворечники"
}'::jsonb WHERE id = 'acf16ee9-af37-4056-8aeb-0060974476bf';

-- Bazeni i spa
UPDATE categories SET meta_keywords = '{
  "sr": "bazeni, spa, jacuzzi, nadzemni bazen, hidromasaža, bazen oprema, bazenska hemija, filteri",
  "en": "pools, spa, jacuzzi, above ground pool, hydromassage, pool equipment, pool chemicals, filters",
  "ru": "бассейны, спа, джакузи, наземный бассейн, гидромассаж, оборудование для бассейна, химия"
}'::jsonb WHERE id = '3f13c642-2632-44c4-be5f-8518d9ab4513';

-- Bicikli i trotineti
UPDATE categories SET meta_keywords = '{
  "sr": "bicikli, trotineti, brdski bicikl, elektricni bicikl, e-scooter, BMX, oprema za bicikl",
  "en": "bicycles, scooters, mountain bike, electric bike, e-scooter, BMX, bike accessories",
  "ru": "велосипеды, самокаты, горный велосипед, электровелосипед, электросамокат, BMX"
}'::jsonb WHERE id = '181f7df7-3375-42b4-812d-98b52b5e3af7';

-- Dečija kozmetika
UPDATE categories SET meta_keywords = '{
  "sr": "dečija kozmetika, baby kozmetika, šampon za decu, krema za bebe, bebi ulje, nega beba",
  "en": "kids cosmetics, baby cosmetics, kids shampoo, baby cream, baby oil, baby care",
  "ru": "детская косметика, детский шампунь, крем для младенцев, детское масло, уход за младенцами"
}'::jsonb WHERE id = '2696706a-8138-4f4e-acb9-164f4cc46c09';

-- Dečija obuća
UPDATE categories SET meta_keywords = '{
  "sr": "dečija obuća, dečije patike, cipele za decu, cipelice, papuče za decu, sportske patike",
  "en": "kids footwear, kids sneakers, children shoes, booties, kids slippers, sports shoes",
  "ru": "детская обувь, детские кроссовки, детские туфли, пинетки, тапочки для детей"
}'::jsonb WHERE id = '472b8c4f-c811-479a-9709-e1bcea5b8dff';

-- Dečija obuća i bebe
UPDATE categories SET meta_keywords = '{
  "sr": "dečija obuća, bebi obuća, prve cipele, papučice, patike za bebe, obuća za novorođenče",
  "en": "kids footwear, baby footwear, first shoes, baby slippers, baby sneakers, newborn shoes",
  "ru": "детская обувь, обувь для новорожденных, первые туфли, пинетки, кроссовки для младенцев"
}'::jsonb WHERE id = '0ca144eb-a9ca-4590-8167-df1287851ba7';

-- Dečija odeća
UPDATE categories SET meta_keywords = '{
  "sr": "dečija odeća, garderoba za decu, majice za decu, pantalone za decu, jakne, haljine",
  "en": "kids clothing, children wardrobe, kids t-shirts, kids pants, jackets, dresses",
  "ru": "детская одежда, гардероб для детей, футболки для детей, брюки для детей, куртки, платья"
}'::jsonb WHERE id = '29c0a4aa-89cc-4200-a245-21b0f0c470cd';

-- Dečija odeća i bebe
UPDATE categories SET meta_keywords = '{
  "sr": "dečija odeća, bebi odeća, bodici, pidžame, kombinezon, odeća za novorođenče, bebe garderoba",
  "en": "kids clothing, baby clothing, bodysuits, pajamas, rompers, newborn clothing, baby wardrobe",
  "ru": "детская одежда, одежда для новорожденных, боди, пижамы, комбинезоны, гардероб для младенцев"
}'::jsonb WHERE id = '9274260f-ec54-4441-888d-c73f84541f89';

-- Dečiji nameštaj
UPDATE categories SET meta_keywords = '{
  "sr": "dečiji nameštaj, dečiji krevetići, stolice za decu, sto za decu, ormar, dečija soba",
  "en": "kids furniture, children cribs, kids chairs, kids table, wardrobe, children room",
  "ru": "детская мебель, детские кроватки, детские стулья, детский стол, шкаф, детская комната"
}'::jsonb WHERE id = 'b3c71d14-9890-4005-bdb5-dde6249ad586';

-- Dekoracije
UPDATE categories SET meta_keywords = '{
  "sr": "dekoracije, ukrasni predmeti, dekor, slike, ram ovi, ukras, figure, vaze, ukrašavanje",
  "en": "decorations, ornaments, decor, paintings, frames, decoration, figures, vases, embellishment",
  "ru": "декор, украшения, декорации, картины, рамки, фигуры, вазы, украшение интерьера"
}'::jsonb WHERE id = '18e0b517-4e8c-4317-82b3-bda5da7ecca2';

-- Dekorativna kozmetika
UPDATE categories SET meta_keywords = '{
  "sr": "dekorativna kozmetika, makeup, šminka, ruž, maskara, puder, rumenilo, olovke za oči",
  "en": "makeup, cosmetics, lipstick, mascara, powder, blush, eyeliner, foundation",
  "ru": "декоративная косметика, макияж, помада, тушь, пудра, румяна, подводка для глаз"
}'::jsonb WHERE id = '610e6d3f-17ab-4187-904a-21e959be4b4b';

-- Depilacija
UPDATE categories SET meta_keywords = '{
  "sr": "depilacija, epilacija, vosak, kreme za depilaciju, brijači, laser, uklanjanje dlaka",
  "en": "hair removal, depilation, wax, hair removal cream, razors, laser, epilation",
  "ru": "депиляция, эпиляция, воск, крем для депиляции, бритвы, лазер, удаление волос"
}'::jsonb WHERE id = 'fa710775-c416-4c24-878f-210c00cd98d5';

-- Desktop računari
UPDATE categories SET meta_keywords = '{
  "sr": "desktop računari, PC, desktop, kućni računar, gaming PC, radne stanice, konfiguracija",
  "en": "desktop computers, PC, desktop, home computer, gaming PC, workstations, configuration",
  "ru": "настольные компьютеры, ПК, домашний компьютер, игровой ПК, рабочие станции"
}'::jsonb WHERE id = 'a0cc8b8a-4018-4c80-8af3-b4e6da709010';

-- Dodaci i aksesoari
UPDATE categories SET meta_keywords = '{
  "sr": "dodaci, aksesoari, modni dodaci, kaiši, šeširi, rukavice, kravate, šalovi",
  "en": "accessories, fashion accessories, belts, hats, gloves, ties, scarves",
  "ru": "аксессуары, модные аксессуары, ремни, шляпы, перчатки, галстуки, шарфы"
}'::jsonb WHERE id = 'cc5e27e5-07fa-412a-a236-41ee92aec4f7';

-- Dodatna oprema (Elektronika)
UPDATE categories SET meta_keywords = '{
  "sr": "dodatna oprema, aksesoari za elektroniku, kablovi, punjači, futrole, nastavci, adapteri",
  "en": "accessories, electronics accessories, cables, chargers, cases, extensions, adapters",
  "ru": "аксессуары, аксессуары для электроники, кабели, зарядные устройства, чехлы, адаптеры"
}'::jsonb WHERE id = 'b72eaab6-7826-4c89-b482-1113a92477e7';

-- Donji veš
UPDATE categories SET meta_keywords = '{
  "sr": "donji veš, gaćice, grudnjaci, bokserice, kupaći veš, intimo, bralette, haljinice",
  "en": "underwear, panties, bras, boxers, lingerie, intimates, bralette, slips",
  "ru": "нижнее белье, трусики, бюстгальтеры, боксеры, белье, интимное белье"
}'::jsonb WHERE id = 'a94be020-9b14-46bb-9f3c-d21bef5771c6';

-- Dronovi
UPDATE categories SET meta_keywords = '{
  "sr": "dronovi, quadcopter, bespilotne letelice, dron kamere, DJI, FPV, RC dronovi, zračna fotografija",
  "en": "drones, quadcopter, UAV, drone cameras, DJI, FPV, RC drones, aerial photography",
  "ru": "дроны, квадрокоптеры, беспилотники, камеры на дронах, DJI, FPV, аэрофотосъемка"
}'::jsonb WHERE id = '48646e1c-bdda-4f31-b8d9-65322d3ee01c';

-- Džonovanje
UPDATE categories SET meta_keywords = '{
  "sr": "džonovanje, trčanje, jogging, trke, patike za trčanje, sprint, maraton, running",
  "en": "jogging, running, races, running shoes, sprint, marathon, running gear",
  "ru": "бег, джоггинг, забеги, кроссовки для бега, спринт, марафон"
}'::jsonb WHERE id = '7ea0ae1b-69cd-4790-ab33-37c527e0c322';

-- E-čitači
UPDATE categories SET meta_keywords = '{
  "sr": "e-čitači, e-reader, kindle, elektronske knjige, e-ink, čitači knjiga, digitalne knjige",
  "en": "e-readers, e-reader, kindle, electronic books, e-ink, book readers, digital books",
  "ru": "электронные книги, ридеры, kindle, e-ink, читалки, цифровые книги"
}'::jsonb WHERE id = '4cd308ee-a94a-4343-8256-476e3eb68439';

-- Elegantna odeća
UPDATE categories SET meta_keywords = '{
  "sr": "elegantna odeća, formalna odeća, odela, smokingzi, koktel haljine, svečana odeća, poslovna",
  "en": "formal wear, formal clothing, suits, tuxedos, cocktail dresses, evening wear, business",
  "ru": "элегантная одежда, формальная одежда, костюмы, смокинги, коктейльные платья, вечерняя одежда"
}'::jsonb WHERE id = 'feff7afc-16bd-4ca1-9af7-ed5bf282dd54';

-- Elektronika za decu
UPDATE categories SET meta_keywords = '{
  "sr": "elektronika za decu, dečije igračke, edukativni tablet, elektronske igračke, robotika, STEM",
  "en": "kids electronics, children toys, educational tablet, electronic toys, robotics, STEM",
  "ru": "электроника для детей, детские игрушки, образовательный планшет, электронные игрушки, STEM"
}'::jsonb WHERE id = 'cc9d2678-9444-48c5-af59-feaafa7fa974';

-- Ešarpe i šalovi
UPDATE categories SET meta_keywords = '{
  "sr": "ešarpe, šalovi, marame, kašmir, vuneni šalovi, zimski dodaci, halstuk",
  "en": "scarves, shawls, bandanas, cashmere, wool scarves, winter accessories",
  "ru": "шарфы, палантины, платки, кашемир, шерстяные шарфы, зимние аксессуары"
}'::jsonb WHERE id = '588940f8-20e9-4922-b9bd-625c32fe7a48';

-- Eterična ulja
UPDATE categories SET meta_keywords = '{
  "sr": "eterična ulja, aromaterapija, esencijalna ulja, difuzori, lavanda, kadulja, mirisna ulja",
  "en": "essential oils, aromatherapy, essential oils, diffusers, lavender, sage, fragrance oils",
  "ru": "эфирные масла, ароматерапия, диффузоры, лаванда, шалфей, ароматические масла"
}'::jsonb WHERE id = '47a847a8-8c15-47e6-a2f1-993e1694fcba';

-- Fitnes i teretana
UPDATE categories SET meta_keywords = '{
  "sr": "fitnes, teretana, gym, tegovi, bučice, traka za trčanje, oprema za vežbanje, cardio",
  "en": "fitness, gym, weights, dumbbells, treadmill, exercise equipment, cardio, workout",
  "ru": "фитнес, тренажерный зал, веса, гантели, беговая дорожка, оборудование для тренировок"
}'::jsonb WHERE id = 'dbd60b13-2f3b-4464-b5c7-c764c0c1b86e';

-- Foto i video kamere
UPDATE categories SET meta_keywords = '{
  "sr": "foto kamere, video kamere, DSLR, mirrorless, GoPro, objektivi, Canon, Nikon, Sony",
  "en": "cameras, camcorders, DSLR, mirrorless, GoPro, lenses, Canon, Nikon, Sony",
  "ru": "фотокамеры, видеокамеры, зеркальные, беззеркальные, GoPro, объективы, Canon, Nikon"
}'::jsonb WHERE id = '00d2e0ce-f505-451d-bac3-cce231edcaf9';

-- Fudbal
UPDATE categories SET meta_keywords = '{
  "sr": "fudbal, football, kopačke, lopta, dresovi, fudbalska oprema, golovi, treninzi",
  "en": "football, soccer, cleats, ball, jerseys, football equipment, goals, training",
  "ru": "футбол, бутсы, мяч, форма, футбольное оборудование, ворота, тренировки"
}'::jsonb WHERE id = '8ffc3b82-24ba-4506-8d72-0d81e1cf7ddd';

-- Gaming oprema
UPDATE categories SET meta_keywords = '{
  "sr": "gaming oprema, gaming miš, tastatura, slušalice, gaming stolica, RGB, gejming, stream",
  "en": "gaming gear, gaming mouse, keyboard, headset, gaming chair, RGB, gaming, stream",
  "ru": "игровое оборудование, игровая мышь, клавиатура, наушники, игровое кресло, RGB"
}'::jsonb WHERE id = 'd99e4a84-cbdb-4735-ae58-ac35e5633bd7';

-- Gelovi za tuširanje
UPDATE categories SET meta_keywords = '{
  "sr": "gelovi za tuširanje, gel, sapun, kupka, tuš, higijena, shower gel, prirodni gelovi",
  "en": "shower gels, gel, soap, bath, shower, hygiene, shower gel, natural gels",
  "ru": "гели для душа, гель, мыло, ванна, душ, гигиена, натуральные гели"
}'::jsonb WHERE id = '90876996-8963-47d1-84ca-a0221a24bbed';

-- Grnčarija i biljke
UPDATE categories SET meta_keywords = '{
  "sr": "grnčarija, biljke, saksije, sobne biljke, kaktusi, sukulenti, zelenilo, baštenska keramika",
  "en": "pots, plants, planters, houseplants, cacti, succulents, greenery, garden ceramics",
  "ru": "горшки, растения, комнатные растения, кактусы, суккуленты, зелень, садовая керамика"
}'::jsonb WHERE id = '1d91aacd-b69f-446c-a586-5443d9b9e4c5';

-- Higijena
UPDATE categories SET meta_keywords = '{
  "sr": "higijena, sapun, šampon, dezodorans, pasta za zube, vlažne maramice, četka, lična higijena",
  "en": "hygiene, soap, shampoo, deodorant, toothpaste, wet wipes, toothbrush, personal hygiene",
  "ru": "гигиена, мыло, шампунь, дезодорант, зубная паста, влажные салфетки, личная гигиена"
}'::jsonb WHERE id = '4e6d66f1-c468-405d-91c1-da69a8e3d6b5';

-- Hrana za bebe
UPDATE categories SET meta_keywords = '{
  "sr": "hrana za bebe, bebi hrana, kašice, dohranu, mleko za bebe, formula, organska hrana",
  "en": "baby food, baby meals, porridge, complementary food, baby milk, formula, organic food",
  "ru": "детское питание, каши, прикорм, детское молоко, смеси, органическое питание"
}'::jsonb WHERE id = '43f0588a-8560-4066-a963-900d00bb23a3';

-- Igračke za bebe
UPDATE categories SET meta_keywords = '{
  "sr": "igračke za bebe, bebi igračke, zvečke, plišane igračke, edukativne igračke, krevetske igračke",
  "en": "baby toys, rattles, plush toys, educational toys, crib toys, infant toys",
  "ru": "игрушки для новорожденных, погремушки, плюшевые игрушки, развивающие игрушки"
}'::jsonb WHERE id = '165337d4-dce5-4d3b-966a-04648b84aefd';

-- Igračke za decu
UPDATE categories SET meta_keywords = '{
  "sr": "igračke za decu, dečije igračke, lutke, autići, LEGO, edukativne igračke, društvene igre",
  "en": "kids toys, children toys, dolls, cars, LEGO, educational toys, board games",
  "ru": "детские игрушки, куклы, машинки, LEGO, развивающие игрушки, настольные игры"
}'::jsonb WHERE id = '993187b7-2b0a-4789-8de1-246eafb59abb';

-- Intimna higijena
UPDATE categories SET meta_keywords = '{
  "sr": "intimna higijena, intimni gelovi, ulosci, tamponi, menstrualni proizvodi, higijena žena",
  "en": "intimate hygiene, intimate gels, pads, tampons, menstrual products, women hygiene",
  "ru": "интимная гигиена, интимные гели, прокладки, тампоны, менструальные продукты"
}'::jsonb WHERE id = '9e51d900-069f-4934-aa18-757f2660c36c';

-- Kalkulatori
UPDATE categories SET meta_keywords = '{
  "sr": "kalkulatori, digitalni kalkulator, naučni kalkulator, grafički kalkulator, računaljka",
  "en": "calculators, digital calculator, scientific calculator, graphing calculator",
  "ru": "калькуляторы, цифровой калькулятор, научный калькулятор, графический калькулятор"
}'::jsonb WHERE id = 'c2af7023-43b7-44fd-a2dd-ae3c7cc43488';

-- Kampovanje
UPDATE categories SET meta_keywords = '{
  "sr": "kampovanje, šatori, vreće za spavanje, oprema za kampovanje, outdoor, kamp, camping",
  "en": "camping, tents, sleeping bags, camping gear, outdoor, camp, backpacking",
  "ru": "кемпинг, палатки, спальные мешки, снаряжение для кемпинга, outdoor"
}'::jsonb WHERE id = 'e02aa6f1-68ec-4680-b75c-0bf677970468';

-- Kompostiranje
UPDATE categories SET meta_keywords = '{
  "sr": "kompostiranje, kompost, organski otpad, ekologija, đubrivo, biorazgradivo, reciklaža",
  "en": "composting, compost, organic waste, ecology, fertilizer, biodegradable, recycling",
  "ru": "компостирование, компост, органические отходы, экология, удобрение, биоразлагаемое"
}'::jsonb WHERE id = '13e68c8b-344d-4cdc-abf1-68254d34c176';

-- Konzole i gaming
UPDATE categories SET meta_keywords = '{
  "sr": "konzole, gaming, PlayStation, Xbox, Nintendo Switch, PS5, igre, gejmeri",
  "en": "consoles, gaming, PlayStation, Xbox, Nintendo Switch, PS5, games, gamers",
  "ru": "консоли, гейминг, PlayStation, Xbox, Nintendo Switch, PS5, игры, геймеры"
}'::jsonb WHERE id = '88de0de9-36e0-43a9-977e-c49d4937b139';

-- Konzolne igre
UPDATE categories SET meta_keywords = '{
  "sr": "konzolne igre, video igre, PS5 igre, Xbox igre, Nintendo igre, gaming, gejmeri",
  "en": "console games, video games, PS5 games, Xbox games, Nintendo games, gaming, gamers",
  "ru": "игры для консолей, видеоигры, игры для PS5, игры для Xbox, Nintendo игры"
}'::jsonb WHERE id = '34721b20-c683-478b-a468-182aaaddd68f';

-- Košarka
UPDATE categories SET meta_keywords = '{
  "sr": "košarka, basketball, patike za košarku, lopta, košarkaška oprema, NBA, dresovi",
  "en": "basketball, basketball shoes, ball, basketball gear, NBA, jerseys, court",
  "ru": "баскетбол, баскетбольные кроссовки, мяч, баскетбольное оборудование, NBA"
}'::jsonb WHERE id = '93e4487d-1405-419a-aee1-0458d54e0f07';

-- Košulje kratkih rukava
UPDATE categories SET meta_keywords = '{
  "sr": "košulje kratkih rukava, kratke košulje, letnje košulje, casual shirts, muške košulje, polo",
  "en": "short sleeve shirts, short shirts, summer shirts, casual shirts, men shirts, polo",
  "ru": "рубашки с короткими рукавами, короткие рубашки, летние рубашки, мужские рубашки"
}'::jsonb WHERE id = '57dce48b-f47a-4c83-ad39-a6124f8d9289';

-- Kuhinjski pribor
UPDATE categories SET meta_keywords = '{
  "sr": "kuhinjski pribor, noževi, tanjiri, šolje, kašike, viljuške, posuđe, kuhinja",
  "en": "kitchenware, knives, plates, cups, spoons, forks, dishes, kitchen utensils",
  "ru": "кухонная утварь, ножи, тарелки, чашки, ложки, вилки, посуда, кухня"
}'::jsonb WHERE id = 'c1b7054d-01e3-42f2-a0e2-3450b49a940a';

-- Kupaći kostimi
UPDATE categories SET meta_keywords = '{
  "sr": "kupaći kostimi, kupaće gaće, bikini, jednodelni kupaći, swim wear, plaža, bazen",
  "en": "swimwear, swim trunks, bikini, one-piece swimsuit, swim wear, beach, pool",
  "ru": "купальники, плавки, бикини, цельный купальник, пляж, бассейн"
}'::jsonb WHERE id = 'c153baa9-8276-4257-8966-84ad0a177bf3';

-- Kupatilo
UPDATE categories SET meta_keywords = '{
  "sr": "kupatilo, tuš kabina, umivaonik, toaletna školjka, kupka, sanitarije, kupatilski nameštaj",
  "en": "bathroom, shower cabin, sink, toilet bowl, bathtub, sanitary ware, bathroom furniture",
  "ru": "ванная комната, душевая кабина, раковина, унитаз, ванна, сантехника, мебель для ванной"
}'::jsonb WHERE id = '4e73c1f5-d393-450b-8448-647542059cf6';

-- Laptop računari
UPDATE categories SET meta_keywords = '{
  "sr": "laptop računari, laptopi, notebook, gaming laptop, poslovno laptop, ultrabook, HP, Dell, Lenovo",
  "en": "laptops, notebook, gaming laptop, business laptop, ultrabook, HP, Dell, Lenovo",
  "ru": "ноутбуки, лаптопы, игровой ноутбук, бизнес ноутбук, ультрабук, HP, Dell, Lenovo"
}'::jsonb WHERE id = '227d356f-6128-4aaf-9bd9-96e2382bdb59';

-- Lepota aparati
UPDATE categories SET meta_keywords = '{
  "sr": "lepota aparati, fen za kosu, pegla za kosu, epilator, aparati za negu, beauty tech",
  "en": "beauty devices, hair dryer, hair straightener, epilator, care devices, beauty tech",
  "ru": "приборы для красоты, фен, утюжок для волос, эпилятор, техника для ухода"
}'::jsonb WHERE id = 'ad3067a9-fb4b-4c48-8c27-b35d8f07d260';

-- Lov
UPDATE categories SET meta_keywords = '{
  "sr": "lov, lovačka oprema, puške, municija, lovački psi, lovačka odeća, outdoor, hunter",
  "en": "hunting, hunting gear, rifles, ammunition, hunting dogs, hunting clothing, outdoor, hunter",
  "ru": "охота, охотничье снаряжение, ружья, боеприпасы, охотничьи собаки, охотничья одежда"
}'::jsonb WHERE id = '24969378-e5c7-4b00-aa1c-0cea90780462';

-- Luksuzna kozmetika
UPDATE categories SET meta_keywords = '{
  "sr": "luksuzna kozmetika, premium kozmetika, Dior, Chanel, luxury beauty, designer kozmetika",
  "en": "luxury cosmetics, premium cosmetics, Dior, Chanel, luxury beauty, designer cosmetics",
  "ru": "люксовая косметика, премиум косметика, Dior, Chanel, дизайнерская косметика"
}'::jsonb WHERE id = '54bf1249-7311-4d80-8501-5a9bfa569c7f';

-- Makeup četkice
UPDATE categories SET meta_keywords = '{
  "sr": "makeup četkice, kistovi za šminku, četke za makeup, beauty brushes, profesionalne četke",
  "en": "makeup brushes, cosmetic brushes, beauty brushes, professional brushes, brush set",
  "ru": "кисти для макияжа, кисточки для косметики, профессиональные кисти, набор кистей"
}'::jsonb WHERE id = '48c096d5-3be0-49ee-9ae4-13a60e0a365e';

-- Manikir i pedikir
UPDATE categories SET meta_keywords = '{
  "sr": "manikir, pedikir, lak za nokte, gel lak, nail art, pila za nokte, nail care",
  "en": "manicure, pedicure, nail polish, gel polish, nail art, nail file, nail care",
  "ru": "маникюр, педикюр, лак для ногтей, гель-лак, нейл-арт, пилка для ногтей"
}'::jsonb WHERE id = 'f1ffb531-6ccd-4459-addd-cf3a14292bd6';

-- Medicinski proizvodi
UPDATE categories SET meta_keywords = '{
  "sr": "medicinski proizvodi, aparat za pritisak, termometri, ortopedski proizvodi, medical supplies, prva pomoć",
  "en": "medical products, blood pressure monitor, thermometers, orthopedic products, medical supplies, first aid",
  "ru": "медицинские товары, тонометр, термометры, ортопедические товары, медтехника, первая помощь"
}'::jsonb WHERE id = 'c14d43c8-0300-4594-bdcb-5773545570fb';

-- Mikrofoni
UPDATE categories SET meta_keywords = '{
  "sr": "mikrofoni, bežični mikrofon, studio mikrofon, USB mikrofon, podcast, streaming, recording",
  "en": "microphones, wireless microphone, studio microphone, USB microphone, podcast, streaming, recording",
  "ru": "микрофоны, беспроводной микрофон, студийный микрофон, USB микрофон, подкаст, стриминг"
}'::jsonb WHERE id = '72c8ab1c-5d03-45c0-bda1-cb040fbc26c9';

-- Mreža i internet
UPDATE categories SET meta_keywords = '{
  "sr": "mreža, internet, ruteri, switch, WiFi, ethernet kablovi, networking, modem, router",
  "en": "networking, internet, routers, switch, WiFi, ethernet cables, networking, modem, router",
  "ru": "сети, интернет, роутеры, свитчи, WiFi, ethernet кабели, сетевое оборудование, модем"
}'::jsonb WHERE id = '7e76e831-18b7-4ee4-939b-eaf85d356666';

-- Muška nega
UPDATE categories SET meta_keywords = '{
  "sr": "muška nega, muška kozmetika, brijaći aparati, kreme za lice, parfemi za muškarce, grooming",
  "en": "men grooming, men cosmetics, razors, face creams, men perfumes, grooming, shaving",
  "ru": "мужской уход, мужская косметика, бритвы, кремы для лица, мужские парфюмы"
}'::jsonb WHERE id = '0fe74127-7bbf-455a-bb2f-9da84437db27';

-- Muška obuća
UPDATE categories SET meta_keywords = '{
  "sr": "muška obuća, muške patike, cipele, sportske patike, elegantne cipele, Nike, Adidas",
  "en": "men footwear, men sneakers, shoes, sports shoes, formal shoes, Nike, Adidas",
  "ru": "мужская обувь, мужские кроссовки, туфли, спортивная обувь, деловая обувь, Nike, Adidas"
}'::jsonb WHERE id = 'c66fefa1-f4d0-4772-a14b-61f288dff767';

-- Muška odeća
UPDATE categories SET meta_keywords = '{
  "sr": "muška odeća, garderoba za muškarce, majice, jakne, pantalone, košulje, moda, stilova",
  "en": "men clothing, men wardrobe, t-shirts, jackets, pants, shirts, fashion, styles",
  "ru": "мужская одежда, мужской гардероб, футболки, куртки, брюки, рубашки, мода, стили"
}'::jsonb WHERE id = '03f55418-c46f-42d1-a91d-dab950352040';

-- Muškarci veliki brojevi
UPDATE categories SET meta_keywords = '{
  "sr": "muškarci veliki brojevi, XXL, XXXL, plus size, veće veličine, odeća za krupnije",
  "en": "men big sizes, XXL, XXXL, plus size, larger sizes, clothing for bigger men",
  "ru": "мужчины большие размеры, XXL, XXXL, большие размеры, одежда для крупных мужчин"
}'::jsonb WHERE id = '13bbc6e0-14c9-4c60-9c64-f228325d3e93';

-- Muški stil
UPDATE categories SET meta_keywords = '{
  "sr": "muški stil, men fashion, trendi, outfit, muška moda, stil, muški dodaci, kombinacija",
  "en": "men style, men fashion, trends, outfit, men fashion, style, men accessories, combination",
  "ru": "мужской стиль, мужская мода, тренды, наряд, стиль, мужские аксессуары"
}'::jsonb WHERE id = '8dda66e1-dadd-426b-88a6-577d5829b579';

-- Nameštaj dnevna soba
UPDATE categories SET meta_keywords = '{
  "sr": "nameštaj dnevna soba, garniture, sofa, fotelje, trpezarijski sto, komoda, TV komoda",
  "en": "living room furniture, sofa sets, sofa, armchairs, dining table, dresser, TV stand",
  "ru": "мебель для гостиной, диваны, кресла, обеденный стол, комод, тумба под ТВ"
}'::jsonb WHERE id = '135ec0ea-132d-4822-bd05-0cc36ea51797';

-- Nameštaj kancelarija
UPDATE categories SET meta_keywords = '{
  "sr": "nameštaj kancelarija, kancelarijski sto, kancelarijska stolica, polica, orman, radni sto",
  "en": "office furniture, office desk, office chair, shelf, cabinet, work desk",
  "ru": "офисная мебель, офисный стол, офисное кресло, полка, шкаф, рабочий стол"
}'::jsonb WHERE id = 'dd1c3d9f-1ecc-4754-9a3d-95434f181f0e';

-- Nameštaj kuhinja
UPDATE categories SET meta_keywords = '{
  "sr": "nameštaj kuhinja, kuhinjski elementi, kuhinjski sto, stolice, kuhinjske police, plakar",
  "en": "kitchen furniture, kitchen cabinets, kitchen table, chairs, kitchen shelves, cupboard",
  "ru": "кухонная мебель, кухонные шкафы, кухонный стол, стулья, кухонные полки, буфет"
}'::jsonb WHERE id = 'c775d481-6b21-4192-8dfe-4f44454f3887';

-- Nameštaj spavaća soba
UPDATE categories SET meta_keywords = '{
  "sr": "nameštaj spavaća soba, kreveti, garderoberi, noćni stolići, komoda, dvospalni krevet",
  "en": "bedroom furniture, beds, wardrobes, nightstands, dresser, double bed",
  "ru": "мебель для спальни, кровати, шкафы для одежды, тумбочки, комоды, двуспальная кровать"
}'::jsonb WHERE id = 'ebe25e54-d91b-4586-8697-4b31ae0c6d34';

-- Nameštaj za bebe
UPDATE categories SET meta_keywords = '{
  "sr": "nameštaj za bebe, krevetići, komode za bebe, stolice za hranjenje, sto za presvlačenje",
  "en": "baby furniture, cribs, baby dressers, highchairs, changing table, nursery",
  "ru": "мебель для новорожденных, кроватки, комоды для младенцев, стульчики для кормления"
}'::jsonb WHERE id = '167e3f8d-eb36-454e-95d7-48b428ee8756';

-- Naočari i dodaci
UPDATE categories SET meta_keywords = '{
  "sr": "naočari, sunčane naočare, dioptrijske naočare, futrole, okviri, Ray-Ban, Oakley",
  "en": "glasses, sunglasses, prescription glasses, cases, frames, Ray-Ban, Oakley",
  "ru": "очки, солнцезащитные очки, очки с диоптриями, футляры, оправы, Ray-Ban, Oakley"
}'::jsonb WHERE id = '7e9cf4e2-19d0-454c-8f76-510058e6d432';

-- NAS i storage
UPDATE categories SET meta_keywords = '{
  "sr": "NAS, storage, network storage, eksterni hard disk, SSD, cloud storage, backup, data",
  "en": "NAS, storage, network storage, external hard drive, SSD, cloud storage, backup, data",
  "ru": "NAS, хранилище, сетевое хранилище, внешний жесткий диск, SSD, облачное хранилище"
}'::jsonb WHERE id = '4808c8c2-2c1c-4a0d-9382-2a42447e7cd6';

-- Nega i higijena beba
UPDATE categories SET meta_keywords = '{
  "sr": "nega i higijena beba, pelene, vlažne maramice, bebi puder, kupanje beba, bebi ulje",
  "en": "baby care, baby hygiene, diapers, wet wipes, baby powder, baby bathing, baby oil",
  "ru": "уход за младенцами, гигиена младенцев, подгузники, влажные салфетки, детская присыпка"
}'::jsonb WHERE id = 'f406e2a7-f3d5-46c9-a9d7-206e521ca7c3';

-- Nega kose
UPDATE categories SET meta_keywords = '{
  "sr": "nega kose, šampon, regenerator, maska za kosu, ulje za kosu, tretmani, hair care",
  "en": "hair care, shampoo, conditioner, hair mask, hair oil, treatments, hair products",
  "ru": "уход за волосами, шампунь, кондиционер, маска для волос, масло для волос"
}'::jsonb WHERE id = '8beab479-446b-4079-8639-4e5e271e44c6';

-- Nega kože
UPDATE categories SET meta_keywords = '{
  "sr": "nega kože, kreme za lice, serumi, hidratacija, anti-age, cleanser, tonik, skincare",
  "en": "skincare, face creams, serums, hydration, anti-age, cleanser, toner, skin care",
  "ru": "уход за кожей, кремы для лица, сыворотки, увлажнение, анти-эйдж, очищение"
}'::jsonb WHERE id = '7c1b5733-29cd-4f35-9f4b-e2163cf069fa';

-- Odela i smokingzi
UPDATE categories SET meta_keywords = '{
  "sr": "odela, smokingzi, frakovi, poslovna odela, venčana odela, crno odelo, suit, tuxedo",
  "en": "suits, tuxedos, tailcoats, business suits, wedding suits, black suit, suit, tuxedo",
  "ru": "костюмы, смокинги, фраки, деловые костюмы, свадебные костюмы, черный костюм"
}'::jsonb WHERE id = 'f392e7ed-a5b4-43b5-aad4-802aa88c5108';

-- Ogledala
UPDATE categories SET meta_keywords = '{
  "sr": "ogledala, ogledalo za kupatilo, ogledalo za hodnik, ukrasna ogledala, zidna ogledala",
  "en": "mirrors, bathroom mirror, hallway mirror, decorative mirrors, wall mirrors",
  "ru": "зеркала, зеркало для ванной, зеркало для прихожей, декоративные зеркала, настенные зеркала"
}'::jsonb WHERE id = 'b3876bd1-07b3-4722-a08e-3c0b27a896f0';

-- Oprema za bebe
UPDATE categories SET meta_keywords = '{
  "sr": "oprema za bebe, kolica, autosedišta, nosiljke, ljuljaške, hranilice, bebi oprema",
  "en": "baby gear, strollers, car seats, carriers, swings, feeding accessories, baby equipment",
  "ru": "оборудование для новорожденных, коляски, автокресла, переноски, качели, кормление"
}'::jsonb WHERE id = '05abfc8b-9963-436a-94ac-d5f6f2af318e';

-- Oralna higijena
UPDATE categories SET meta_keywords = '{
  "sr": "oralna higijena, četkica za zube, pasta za zube, konac za zube, vodica za ispiranje, dental",
  "en": "oral hygiene, toothbrush, toothpaste, dental floss, mouthwash, dental care",
  "ru": "гигиена полости рта, зубная щетка, зубная паста, зубная нить, ополаскиватель"
}'::jsonb WHERE id = '50428ca2-d361-4d14-9232-fb8fde879a33';

-- Organizacija i skladištenje
UPDATE categories SET meta_keywords = '{
  "sr": "organizacija, skladištenje, kutije, police, vešalice, korpe, organaizeri, storage",
  "en": "organization, storage, boxes, shelves, hangers, baskets, organizers, storage solutions",
  "ru": "организация, хранение, коробки, полки, вешалки, корзины, органайзеры"
}'::jsonb WHERE id = '60174245-1b9c-4685-aba3-7e27e1fa6b1d';

-- Organska kozmetika
UPDATE categories SET meta_keywords = '{
  "sr": "organska kozmetika, prirodna kozmetika, bio kozmetika, vegan kozmetika, eko, prirodno",
  "en": "organic cosmetics, natural cosmetics, bio cosmetics, vegan cosmetics, eco, natural",
  "ru": "органическая косметика, натуральная косметика, био косметика, веганская косметика"
}'::jsonb WHERE id = '40644f9e-9b53-4e46-bf1f-94f41cc8967a';

-- Pametni satovi
UPDATE categories SET meta_keywords = '{
  "sr": "pametni satovi, smartwatch, Apple Watch, Samsung Galaxy Watch, fitness tracker, wearable",
  "en": "smartwatches, smartwatch, Apple Watch, Samsung Galaxy Watch, fitness tracker, wearable",
  "ru": "умные часы, смарт-часы, Apple Watch, Samsung Galaxy Watch, фитнес-трекер"
}'::jsonb WHERE id = 'bbb97e8d-1ea5-48c7-91f0-36476ea38c31';

-- Pametni telefoni
UPDATE categories SET meta_keywords = '{
  "sr": "pametni telefoni, smartphone, iPhone, Samsung, Xiaomi, Huawei, Android, mobilni telefoni",
  "en": "smartphones, smartphone, iPhone, Samsung, Xiaomi, Huawei, Android, mobile phones",
  "ru": "смартфоны, телефоны, iPhone, Samsung, Xiaomi, Huawei, Android, мобильные телефоны"
}'::jsonb WHERE id = 'ba829402-0e4e-467b-b117-e1da1e2b51ce';

-- Parfemi
UPDATE categories SET meta_keywords = '{
  "sr": "parfemi, mirisi, eau de parfum, toaletna voda, muški parfemi, ženski parfemi, Chanel, Dior",
  "en": "perfumes, fragrances, eau de parfum, eau de toilette, men perfumes, women perfumes, Chanel, Dior",
  "ru": "парфюмерия, ароматы, парфюмерная вода, туалетная вода, мужские духи, женские духи"
}'::jsonb WHERE id = '13653900-58b7-4bb3-bcb2-6177e2f50d61';

-- Periferija
UPDATE categories SET meta_keywords = '{
  "sr": "periferija, miš, tastatura, monitori, printere, skeneri, USB, računarska oprema",
  "en": "peripherals, mouse, keyboard, monitors, printers, scanners, USB, computer accessories",
  "ru": "периферия, мышь, клавиатура, мониторы, принтеры, сканеры, USB, компьютерные аксессуары"
}'::jsonb WHERE id = '3718ce2c-0364-4941-8874-5fd8562be1a2';

-- Planinarenje
UPDATE categories SET meta_keywords = '{
  "sr": "planinarenje, hiking, trekking, planinske cipele, ruksaci, outdoor, planinska oprema",
  "en": "hiking, trekking, hiking boots, backpacks, outdoor, mountain gear, trails",
  "ru": "пеший туризм, хайкинг, трекинг, треккинговая обувь, рюкзаки, горное снаряжение"
}'::jsonb WHERE id = '314c4525-dd31-4db0-95fc-1cd51b89b3f4';

-- Plivanje
UPDATE categories SET meta_keywords = '{
  "sr": "plivanje, kupaći kostimi, naočare za plivanje, kapa za plivanje, swim, bazen, plivačka oprema",
  "en": "swimming, swimsuits, swimming goggles, swim cap, swim, pool, swimming gear",
  "ru": "плавание, купальники, очки для плавания, шапочка для плавания, бассейн"
}'::jsonb WHERE id = 'db3af33d-5651-408c-a0f4-8e01184ee180';

-- Plus size odeća
UPDATE categories SET meta_keywords = '{
  "sr": "plus size odeća, velike veličine, curvy, odeća za krupnije, XXL, XXXL, body positive",
  "en": "plus size clothing, large sizes, curvy, clothing for bigger, XXL, XXXL, body positive",
  "ru": "одежда больших размеров, большие размеры, curvy, одежда для полных, XXL, XXXL"
}'::jsonb WHERE id = '19ad330f-1195-466c-b7d5-d70d9f758c54';

-- Posteljina i peškiri
UPDATE categories SET meta_keywords = '{
  "sr": "posteljina, peškiri, čaršafi, jorgani, jastuci, pamučna posteljina, kvalitetno rublje",
  "en": "bedding, towels, sheets, comforters, pillows, cotton bedding, quality linen",
  "ru": "постельное белье, полотенца, простыни, одеяла, подушки, хлопковое белье"
}'::jsonb WHERE id = '5dff2d36-c1d8-4c09-b3b6-6deeb0417abc';

-- Pregradni zidovi
UPDATE categories SET meta_keywords = '{
  "sr": "pregradni zidovi, paravani, sobni pregradni zidovi, separatori, privacy, dekorativni zidovi",
  "en": "room dividers, screens, room partitions, separators, privacy, decorative walls",
  "ru": "перегородки, ширмы, комнатные перегородки, сепараторы, приватность, декоративные стены"
}'::jsonb WHERE id = 'f072b85c-a210-4a61-bfc2-2e64923bc38b';

-- Projektori
UPDATE categories SET meta_keywords = '{
  "sr": "projektori, video projektori, prezentacioni projektori, home theater, 4K, HDMI, ekrani",
  "en": "projectors, video projectors, presentation projectors, home theater, 4K, HDMI, screens",
  "ru": "проекторы, видеопроекторы, презентационные проекторы, домашний кинотеатр, 4K"
}'::jsonb WHERE id = '877a5d01-c9a8-49a4-958a-2f1f6bedea58';

-- Računarske komponente
UPDATE categories SET meta_keywords = '{
  "sr": "računarske komponente, CPU, GPU, RAM, matična ploča, napajanje, SSD, HDD, cooling",
  "en": "computer components, CPU, GPU, RAM, motherboard, power supply, SSD, HDD, cooling",
  "ru": "компьютерные компоненты, процессор, видеокарта, оперативная память, материнская плата"
}'::jsonb WHERE id = '619909cf-8e95-46a1-b664-8af503ac1e2f';

-- Radna odeća
UPDATE categories SET meta_keywords = '{
  "sr": "radna odeća, radne pantalone, radna jakna, uniforma, zaštitna odeća, workwear, radna obuća",
  "en": "workwear, work pants, work jacket, uniform, protective clothing, workwear, work shoes",
  "ru": "рабочая одежда, рабочие брюки, рабочая куртка, униформа, защитная одежда"
}'::jsonb WHERE id = '49d6d6e9-3222-4780-94f7-98a7a604be4b';

-- Rasveta
UPDATE categories SET meta_keywords = '{
  "sr": "rasveta, lampe, lusteri, LED sijalice, podna svetla, zidne lampe, svetiljke, osvetljenje",
  "en": "lighting, lamps, chandeliers, LED bulbs, floor lamps, wall lamps, lights, illumination",
  "ru": "освещение, лампы, люстры, LED лампочки, торшеры, настенные светильники, свет"
}'::jsonb WHERE id = 'f493ba94-806c-4983-adc7-6b2e7a58b1a8';

-- Ribolov
UPDATE categories SET meta_keywords = '{
  "sr": "ribolov, pecanje, štapovi, mamci, roleri, ribolovačka oprema, fishing, рибарење",
  "en": "fishing, fishing rods, lures, reels, fishing gear, angling, tackle",
  "ru": "рыбалка, удочки, приманки, катушки, рыболовное снаряжение, снасти"
}'::jsonb WHERE id = 'f97c75e5-31df-4bbd-b937-75a31fbd1f59';

-- Satovi za zid
UPDATE categories SET meta_keywords = '{
  "sr": "satovi za zid, zidni satovi, kućni satovi, dekorativni satovi, quartz, digitalni satovi",
  "en": "wall clocks, wall watches, home clocks, decorative clocks, quartz, digital clocks",
  "ru": "настенные часы, домашние часы, декоративные часы, кварцевые, цифровые часы"
}'::jsonb WHERE id = 'c1bb224c-06af-4161-9072-5e0a87780951';

-- Skeneri
UPDATE categories SET meta_keywords = '{
  "sr": "skeneri, scanneri, dokumenti skener, flatbed skener, wireless skener, OCR, digitalizacija",
  "en": "scanners, document scanner, flatbed scanner, wireless scanner, OCR, digitization",
  "ru": "сканеры, сканер документов, планшетный сканер, беспроводной сканер, OCR, оцифровка"
}'::jsonb WHERE id = '79749ca4-d5b3-463e-9455-d6657db3623c';

-- Školski pribor
UPDATE categories SET meta_keywords = '{
  "sr": "školski pribor, sveske, olovke, ranac, pribor za školu, školske torbe, školski set",
  "en": "school supplies, notebooks, pens, backpack, school accessories, school bags, school set",
  "ru": "школьные принадлежности, тетради, ручки, рюкзак, школьные сумки, школьный набор"
}'::jsonb WHERE id = 'ce68a877-1881-4d25-9562-50254820656a';

-- Smart home
UPDATE categories SET meta_keywords = '{
  "sr": "smart home, pametna kuća, automatizacija, IoT, pametne sijalice, smart rasveta, Alexa, Google Home",
  "en": "smart home, smart house, automation, IoT, smart bulbs, smart lighting, Alexa, Google Home",
  "ru": "умный дом, домашняя автоматизация, IoT, умные лампочки, умное освещение, Alexa"
}'::jsonb WHERE id = 'da3aafa2-86a4-48ed-92a6-034fd6200a59';

-- Smart narukvice
UPDATE categories SET meta_keywords = '{
  "sr": "smart narukvice, fitness narukvica, activity tracker, Mi Band, Xiaomi, zdravlje, koraci",
  "en": "smart bands, fitness band, activity tracker, Mi Band, Xiaomi, health, steps",
  "ru": "умные браслеты, фитнес-браслет, трекер активности, Mi Band, Xiaomi, здоровье"
}'::jsonb WHERE id = '33c94cc8-527e-4b3e-8787-16de281c3490';

-- Spa i relaksacija
UPDATE categories SET meta_keywords = '{
  "sr": "spa, relaksacija, masaža, aromaterapija, wellness, sauna, relax, spa tretmani",
  "en": "spa, relaxation, massage, aromatherapy, wellness, sauna, relax, spa treatments",
  "ru": "спа, релаксация, массаж, ароматерапия, велнес, сауна, расслабление, спа-процедуры"
}'::jsonb WHERE id = 'b52997f7-4800-43f7-89bb-bc40e5a621cf';

-- Sportska odeća
UPDATE categories SET meta_keywords = '{
  "sr": "sportska odeća, aktivna odeća, dres, trenerka, helanke, sportski dukser, Nike, Adidas",
  "en": "sportswear, activewear, jersey, tracksuit, leggings, sport hoodie, Nike, Adidas",
  "ru": "спортивная одежда, активная одежда, спортивный костюм, леггинсы, спортивная толстовка"
}'::jsonb WHERE id = '19a40b88-3dd5-4b6e-9937-09462e07c3c4';

-- Tableti
UPDATE categories SET meta_keywords = '{
  "sr": "tableti, tablet, iPad, Samsung tablet, Android tablet, tablet računari, touch screen",
  "en": "tablets, tablet, iPad, Samsung tablet, Android tablet, tablet computers, touch screen",
  "ru": "планшеты, планшет, iPad, Samsung планшет, Android планшет, планшетные компьютеры"
}'::jsonb WHERE id = '14773b71-2c74-4f1c-b709-cbf008cb038a';

-- Tekstil za dom
UPDATE categories SET meta_keywords = '{
  "sr": "tekstil za dom, zavese, stolnjaci, jastučnice, ukrasni jastuk, kućni tekstil, dekorativni tekstil",
  "en": "home textiles, curtains, tablecloths, pillowcases, decorative pillow, house textiles",
  "ru": "домашний текстиль, шторы, скатерти, наволочки, декоративные подушки, текстиль для дома"
}'::jsonb WHERE id = 'efb32041-a42f-43ae-a0d4-443431d179ef';

-- Tenis
UPDATE categories SET meta_keywords = '{
  "sr": "tenis, teniski reket, teniski lopti ce, teniski patike, teniski dresovi, tennis, tereni",
  "en": "tennis, tennis racket, tennis balls, tennis shoes, tennis apparel, tennis, courts",
  "ru": "теннис, теннисная ракетка, теннисные мячи, теннисная обувь, теннисная одежда"
}'::jsonb WHERE id = 'cb9d4c17-150c-4e82-8d6e-8e35f0304a3b';

-- Tepisi i prostirke
UPDATE categories SET meta_keywords = '{
  "sr": "tepisi, prostirke, ćilimi, otirači, tepih za dnevnu sobu, hodničke staze, dekori",
  "en": "carpets, rugs, kilims, doormats, living room carpet, hallway runners, decor",
  "ru": "ковры, коврики, килимы, придверные коврики, ковер для гостиной, дорожки"
}'::jsonb WHERE id = '5829311e-48fb-49ed-87fc-7c6786350d26';

-- Torbice i novčanici
UPDATE categories SET meta_keywords = '{
  "sr": "torbice, novčanici, torbe, kožne torbice, muški novčanici, ženski novčanici, luksuzne torbe",
  "en": "bags, wallets, handbags, leather bags, men wallets, women wallets, luxury bags",
  "ru": "сумки, кошельки, сумочки, кожаные сумки, мужские кошельки, женские кошельки"
}'::jsonb WHERE id = 'eb4d2951-2bc1-43c5-93ed-51d5bebf3a55';

-- Trudnička odeća
UPDATE categories SET meta_keywords = '{
  "sr": "trudnička odeća, odeća za trudnice, majice za trudnice, pantalone za trudnice, haljine za trudnice",
  "en": "maternity wear, maternity clothing, maternity tops, maternity pants, maternity dresses",
  "ru": "одежда для беременных, майки для беременных, брюки для беременных, платья для беременных"
}'::jsonb WHERE id = 'f90c1634-2e4e-4b2c-9513-70b5c0016b63';

-- TV i video
UPDATE categories SET meta_keywords = '{
  "sr": "TV, televizori, Smart TV, 4K TV, OLED, QLED, video, Samsung, LG, Sony",
  "en": "TV, televisions, Smart TV, 4K TV, OLED, QLED, video, Samsung, LG, Sony",
  "ru": "ТВ, телевизоры, Smart TV, 4K телевизор, OLED, QLED, видео, Samsung, LG, Sony"
}'::jsonb WHERE id = '81d5c8da-0dbe-47bb-a2c5-38dde686ecb9';

-- Vaze i dekor
UPDATE categories SET meta_keywords = '{
  "sr": "vaze, dekoracija, ukrasni predmeti, vaze za cveće, dekorativne vaze, dekor za dom",
  "en": "vases, decoration, ornaments, flower vases, decorative vases, home decor",
  "ru": "вазы, декор, украшения, вазы для цветов, декоративные вазы, декор для дома"
}'::jsonb WHERE id = 'b380ef49-e9ee-41d9-8324-b8c7337ec60f';

-- Venčana odeća
UPDATE categories SET meta_keywords = '{
  "sr": "venčana odeća, venčanice, venčano odelo, svadba, venčani dodaci, mladoženja, mladina",
  "en": "wedding attire, wedding dress, wedding suit, wedding, wedding accessories, groom, bride",
  "ru": "свадебная одежда, свадебное платье, свадебный костюм, свадьба, невеста, жених"
}'::jsonb WHERE id = 'fe2e875f-df0c-4ed6-b6c0-a26c3bf47097';

-- Ventilacija i klimatizacija
UPDATE categories SET meta_keywords = '{
  "sr": "ventilacija, klimatizacija, klima uređaj, ventilatori, klima sistemi, hlađenje, grejanje",
  "en": "ventilation, air conditioning, air conditioner, fans, AC systems, cooling, heating",
  "ru": "вентиляция, кондиционирование, кондиционер, вентиляторы, системы кондиционирования"
}'::jsonb WHERE id = 'ef792439-12b4-4047-a7d4-4ffe40447d96';

-- Veš mašine dodaci
UPDATE categories SET meta_keywords = '{
  "sr": "veš mašine dodaci, deterdžent, omekšivač, antikalc, filter za veš mašinu, pranje veša",
  "en": "laundry accessories, detergent, fabric softener, descaler, washing machine filter, laundry",
  "ru": "аксессуары для стирки, моющее средство, кондиционер для белья, фильтр для стиральной машины"
}'::jsonb WHERE id = 'f309d99d-d2d5-43ed-b004-df4ee345ad07';

-- Vitamini i suplementi
UPDATE categories SET meta_keywords = '{
  "sr": "vitamini, suplementi, dodaci ishrani, vitamin C, vitamin D, omega 3, proteini, wellness",
  "en": "vitamins, supplements, dietary supplements, vitamin C, vitamin D, omega 3, proteins, wellness",
  "ru": "витамины, добавки, пищевые добавки, витамин C, витамин D, омега 3, протеины"
}'::jsonb WHERE id = '655906b4-000b-4b55-9806-a5c6efbc51e0';

-- Web kamere
UPDATE categories SET meta_keywords = '{
  "sr": "web kamere, webcam, kamera za računar, streaming kamera, video poziv, Logitech, 1080p",
  "en": "webcams, webcam, computer camera, streaming camera, video call, Logitech, 1080p",
  "ru": "веб-камеры, вебкамера, камера для компьютера, стрим-камера, видеозвонок, Logitech"
}'::jsonb WHERE id = 'f7e7319c-96e8-4864-a5e0-cfe421b04ccc';

-- Zaštita od sunca
UPDATE categories SET meta_keywords = '{
  "sr": "zaštita od sunca, kreme za sunčanje, SPF, losion za sunčanje, after sun, UV zaštita",
  "en": "sun protection, sunscreen, SPF, suntan lotion, after sun, UV protection, sun care",
  "ru": "защита от солнца, солнцезащитный крем, SPF, лосьон для загара, после солнца"
}'::jsonb WHERE id = 'c5ce23d3-3e70-4b87-81d6-9954462f4519';

-- Žene veliki brojevi
UPDATE categories SET meta_keywords = '{
  "sr": "žene veliki brojevi, plus size, XXL, XXXL, veće veličine, odeća za krupnije žene",
  "en": "women big sizes, plus size, XXL, XXXL, larger sizes, clothing for bigger women",
  "ru": "женщины большие размеры, плюс сайз, XXL, XXXL, большие размеры, одежда для полных женщин"
}'::jsonb WHERE id = '48d7c16e-2bff-4506-8b86-8b1c4baea643';

-- Ženska obuća
UPDATE categories SET meta_keywords = '{
  "sr": "ženska obuća, ženske patike, cipele, sandale, čizme, štikle, baletanke, Nike, Adidas",
  "en": "women footwear, women sneakers, shoes, sandals, boots, heels, flats, Nike, Adidas",
  "ru": "женская обувь, женские кроссовки, туфли, сандалии, сапоги, каблуки, балетки"
}'::jsonb WHERE id = '38756e71-1b6a-44e8-9702-fad06761f3b0';

-- Ženska odeća
UPDATE categories SET meta_keywords = '{
  "sr": "ženska odeća, garderoba za žene, haljine, majice, jakne, pantalone, bluze, moda",
  "en": "women clothing, women wardrobe, dresses, shirts, jackets, pants, blouses, fashion",
  "ru": "женская одежда, женский гардероб, платья, футболки, куртки, брюки, блузки, мода"
}'::jsonb WHERE id = 'e4cfe7a3-ef8f-4ff7-a8cf-1716f01b3c2f';

-- Zimska garderoba
UPDATE categories SET meta_keywords = '{
  "sr": "zimska garderoba, zimske jakne, puf jakne, kape, rukavice, šalovi, zimska odeća",
  "en": "winter clothing, winter jackets, puffer jackets, hats, gloves, scarves, winter wear",
  "ru": "зимняя одежда, зимние куртки, пуховики, шапки, перчатки, шарфы, зимний гардероб"
}'::jsonb WHERE id = '374851cb-b7a3-499f-bcfe-6196a69abfaa';

-- Zimski sportovi
UPDATE categories SET meta_keywords = '{
  "sr": "zimski sportovi, skijanje, snowboarding, skije, snowboard, zimska oprema, ski pass",
  "en": "winter sports, skiing, snowboarding, skis, snowboard, winter gear, ski pass",
  "ru": "зимние виды спорта, лыжи, сноуборд, лыжное снаряжение, ски-пасс"
}'::jsonb WHERE id = 'f866494f-e9d0-4585-b1cf-f4f38e29efcb';

-- Level 3 Categories

-- Apple iPhone
UPDATE categories SET meta_keywords = '{
  "sr": "apple iphone, iPhone, iPhone 15, iPhone Pro, iOS, Apple, pametni telefon, premium telefon",
  "en": "apple iphone, iPhone, iPhone 15, iPhone Pro, iOS, Apple, smartphone, premium phone",
  "ru": "apple iphone, айфон, iPhone 15, iPhone Pro, iOS, Apple, смартфон, премиум телефон"
}'::jsonb WHERE id = 'dbd48c81-ab21-42f0-b74d-45e2a408dc5e';

-- Bežične slušalice
UPDATE categories SET meta_keywords = '{
  "sr": "bežične slušalice, AirPods, bluetooth slušalice, wireless headphones, earbuds, TWS",
  "en": "wireless headphones, AirPods, bluetooth headphones, wireless earphones, earbuds, TWS",
  "ru": "беспроводные наушники, AirPods, bluetooth наушники, беспроводные наушники, TWS"
}'::jsonb WHERE id = '95816182-7ba7-450d-ba70-114f281e70fa';

-- Držači za telefon
UPDATE categories SET meta_keywords = '{
  "sr": "držači za telefon, auto držač, držač za bicikl, magnetni držač, phone holder, stalak",
  "en": "phone holders, car holder, bike holder, magnetic holder, phone mount, stand",
  "ru": "держатели для телефона, автомобильный держатель, велосипедный держатель, подставка"
}'::jsonb WHERE id = '4326ea65-579e-42f7-9f13-91d8e02344a4';

-- Maske za telefone
UPDATE categories SET meta_keywords = '{
  "sr": "maske za telefone, futrole, zaštitne maske, silikon maske, kožne maske, phone cases",
  "en": "phone cases, phone covers, protective cases, silicone cases, leather cases, cases",
  "ru": "чехлы для телефонов, защитные чехлы, силиконовые чехлы, кожаные чехлы"
}'::jsonb WHERE id = 'd1744e94-f222-4066-be99-861f5370534f';

-- Muška poslovna odeća
UPDATE categories SET meta_keywords = '{
  "sr": "muška poslovna odeća, kos tumi, kravate, košulje, elegantna odeća, formalna odeća, poslovni stil",
  "en": "men business clothing, suits, ties, shirts, formal wear, business attire, business style",
  "ru": "мужская деловая одежда, костюмы, галстуки, рубашки, формальная одежда, деловой стиль"
}'::jsonb WHERE id = '71d73af3-d488-498b-a957-e70fb1187fe5';

-- Muska sportska odeća
UPDATE categories SET meta_keywords = '{
  "sr": "muska sportska odeća, dresovi, trenere, sport majice, muška activewear, Nike, Adidas",
  "en": "men sportswear, tracksuits, sweatpants, sport shirts, men activewear, Nike, Adidas",
  "ru": "мужская спортивная одежда, спортивные костюмы, тренировочные брюки, Nike, Adidas"
}'::jsonb WHERE id = '0a406d4c-8375-4ab7-85d2-9bd14690c106';

-- Muške jakne
UPDATE categories SET meta_keywords = '{
  "sr": "muške jakne, muška jakna, zimske jakne, kožne jakne, puf jakne, vetrovke, bomber jakne",
  "en": "men jackets, men jacket, winter jackets, leather jackets, puffer jackets, windbreakers",
  "ru": "мужские куртки, зимние куртки, кожаные куртки, пуховики, ветровки, бомберы"
}'::jsonb WHERE id = 'f986ac4e-1db1-4df5-9fa5-8ea6f3937369';

-- Muške košulje
UPDATE categories SET meta_keywords = '{
  "sr": "muške košulje, košulja, elegantne košulje, casual košulje, letnje košulje, muške košulje kratkih rukava",
  "en": "men shirts, shirt, formal shirts, casual shirts, summer shirts, men short sleeve shirts",
  "ru": "мужские рубашки, рубашка, деловые рубашки, повседневные рубашки, летние рубашки"
}'::jsonb WHERE id = '7db79909-8be8-48e3-b9a8-6f6176c25de6';

-- Muške majice
UPDATE categories SET meta_keywords = '{
  "sr": "muške majice, muška majica, T-shirt, polo majice, majice kratkih rukava, casual majice",
  "en": "men t-shirts, men t-shirt, T-shirt, polo shirts, short sleeve shirts, casual shirts",
  "ru": "мужские футболки, футболка, поло, футболки с короткими рукавами, casual футболки"
}'::jsonb WHERE id = '5e29727c-6300-4ded-81e8-993bfe1556a6';

-- Muške pantalone
UPDATE categories SET meta_keywords = '{
  "sr": "muške pantalone, pantalone, farmerke, traperice, chino pantalone, elegantne pantalone",
  "en": "men pants, pants, jeans, trousers, chino pants, formal pants, men trousers",
  "ru": "мужские брюки, брюки, джинсы, чинос, деловые брюки, классические брюки"
}'::jsonb WHERE id = 'ef6a98f4-3d97-46c9-ad59-cb5478ed4ca6';

-- Muški džemperi
UPDATE categories SET meta_keywords = '{
  "sr": "muški džemperi, džemper, puloveri, vuneni džemperi, muške puloveri, kašmir",
  "en": "men sweaters, sweater, pullovers, wool sweaters, men pullovers, cashmere",
  "ru": "мужские свитеры, свитер, пуловеры, шерстяные свитеры, кашемир"
}'::jsonb WHERE id = 'b97afd3c-e2e3-4316-b4a0-1829481b32d2';

-- Muški šorcevi
UPDATE categories SET meta_keywords = '{
  "sr": "muški šorcevi, šorcevi, kratke pantalone, letnji šorcevi, sportski šorcevi, plažni šorcevi",
  "en": "men shorts, shorts, short pants, summer shorts, sport shorts, beach shorts",
  "ru": "мужские шорты, шорты, короткие брюки, летние шорты, спортивные шорты, пляжные шорты"
}'::jsonb WHERE id = '6598279e-a369-49ab-98d5-e6e00f652498';

-- Polovni telefoni
UPDATE categories SET meta_keywords = '{
  "sr": "polovni telefoni, second hand telefoni, refurbished, rabljeni telefoni, korišćeni telefoni, polovno",
  "en": "used phones, second hand phones, refurbished, pre-owned phones, used smartphones",
  "ru": "бывшие в употреблении телефоны, б/у телефоны, восстановленные, подержанные телефоны"
}'::jsonb WHERE id = '5fe5ac3a-7d51-4319-b20c-ee8da4ee698c';

-- Power bank
UPDATE categories SET meta_keywords = '{
  "sr": "power bank, eksterni punjač, mobilna baterija, prenosivi punjač, USB power bank",
  "en": "power bank, external battery, portable charger, mobile battery, USB power bank",
  "ru": "повербанк, внешний аккумулятор, портативное зарядное устройство, мобильная батарея"
}'::jsonb WHERE id = 'f053d318-a2bd-49e2-921d-62d72a5f268b';

-- Punjači i kablovi (telefon)
UPDATE categories SET meta_keywords = '{
  "sr": "punjači, kablovi, USB kablovi, brzi punjači, wireless punjač, kabl za iPhone, Type-C",
  "en": "chargers, cables, USB cables, fast chargers, wireless charger, iPhone cable, Type-C",
  "ru": "зарядные устройства, кабели, USB кабели, быстрые зарядки, беспроводная зарядка"
}'::jsonb WHERE id = 'eebb7b72-daa1-4f33-b055-505abc31c818';

-- Samsung telefoni
UPDATE categories SET meta_keywords = '{
  "sr": "samsung telefoni, Samsung Galaxy, Samsung, Android, Galaxy S, Galaxy A, smartphone Samsung",
  "en": "samsung phones, Samsung Galaxy, Samsung, Android, Galaxy S, Galaxy A, Samsung smartphone",
  "ru": "телефоны Samsung, Samsung Galaxy, Самсунг, Android, Galaxy S, Galaxy A"
}'::jsonb WHERE id = '8b75d427-1bf9-4b7d-9446-94b40f085546';

-- Xiaomi telefoni
UPDATE categories SET meta_keywords = '{
  "sr": "xiaomi telefoni, Xiaomi, Redmi, Poco, Mi, Android, xiaomi smartphone, Xiaomi telefon",
  "en": "xiaomi phones, Xiaomi, Redmi, Poco, Mi, Android, xiaomi smartphone, Xiaomi phone",
  "ru": "телефоны Xiaomi, Сяоми, Redmi, Poco, Mi, Android, смартфон Xiaomi"
}'::jsonb WHERE id = 'b587109e-0668-49be-8088-49047f6201bb';

-- Zaštitno staklo
UPDATE categories SET meta_keywords = '{
  "sr": "zaštitno staklo, staklo za ekran, screen protector, tempered glass, zaštita ekrana",
  "en": "screen protectors, tempered glass, screen protection, glass screen protector, display protection",
  "ru": "защитное стекло, стекло для экрана, защита экрана, каленое стекло"
}'::jsonb WHERE id = '20239311-d7ec-46f7-9d4c-e52509994f0b';

-- Ženska poslovna odeća
UPDATE categories SET meta_keywords = '{
  "sr": "ženska poslovna odeća, ženska kostimi, poslovna odeća, formalna odeća, košulje, poslovni stil",
  "en": "women business clothing, women suits, business attire, formal wear, shirts, business style",
  "ru": "женская деловая одежда, женские костюмы, формальная одежда, рубашки, деловой стиль"
}'::jsonb WHERE id = '84adecbb-186a-4e63-aeae-5a9efe892c51';

-- Ženske bluze
UPDATE categories SET meta_keywords = '{
  "sr": "ženske bluze, bluza, elegantne bluze, letnje bluze, košulje, ženske košulje",
  "en": "women blouses, blouse, formal blouses, summer blouses, shirts, women shirts",
  "ru": "женские блузки, блузка, деловые блузки, летние блузки, рубашки, женские рубашки"
}'::jsonb WHERE id = 'bb12c180-b24d-4e70-93a2-f770b4facfad';

-- Ženske haljine
UPDATE categories SET meta_keywords = '{
  "sr": "ženske haljine, haljina, letnje haljine, elegantne haljine, koktel haljine, maxi haljine",
  "en": "women dresses, dress, summer dresses, formal dresses, cocktail dresses, maxi dresses",
  "ru": "женские платья, платье, летние платья, вечерние платья, коктейльные платья, макси платья"
}'::jsonb WHERE id = 'ffbdeb0e-3287-4b89-a278-45191a58a840';

-- Ženske jakne
UPDATE categories SET meta_keywords = '{
  "sr": "ženske jakne, jakna, zimske jakne, kožne jakne, puf jakne, bomber jakne, vetrovke",
  "en": "women jackets, jacket, winter jackets, leather jackets, puffer jackets, bomber jackets",
  "ru": "женские куртки, куртка, зимние куртки, кожаные куртки, пуховики, бомберы, ветровки"
}'::jsonb WHERE id = '2b721e1f-dcd4-4a8d-ae44-c69fc19c1556';

-- Ženske majice
UPDATE categories SET meta_keywords = '{
  "sr": "ženske majice, majica, T-shirt, casual majice, majice kratkih rukava, top",
  "en": "women t-shirts, t-shirt, T-shirt, casual shirts, short sleeve shirts, top",
  "ru": "женские футболки, футболка, casual футболки, футболки с короткими рукавами, топ"
}'::jsonb WHERE id = '4965cc81-6d4e-4cc3-b9b1-1811ad1f3f1b';

-- Ženske pantalone
UPDATE categories SET meta_keywords = '{
  "sr": "ženske pantalone, pantalone, farmerke, traperice, elegantne pantalone, leggings",
  "en": "women pants, pants, jeans, trousers, formal pants, leggings, women trousers",
  "ru": "женские брюки, брюки, джинсы, деловые брюки, леггинсы, классические брюки"
}'::jsonb WHERE id = 'f5282524-1e6d-45c5-806e-3e99b1135349';

-- Ženske suknje
UPDATE categories SET meta_keywords = '{
  "sr": "ženske suknje, suknja, mini suknje, midi suknje, maxi suknje, elegantne suknje",
  "en": "women skirts, skirt, mini skirts, midi skirts, maxi skirts, formal skirts",
  "ru": "женские юбки, юбка, мини юбки, миди юбки, макси юбки, деловые юбки"
}'::jsonb WHERE id = '79129e9f-6343-449a-9a01-e291d2e6a847';

-- Ženski džemperi
UPDATE categories SET meta_keywords = '{
  "sr": "ženski džemperi, džemper, puloveri, vuneni džemperi, kašmir, ženski puloveri",
  "en": "women sweaters, sweater, pullovers, wool sweaters, cashmere, women pullovers",
  "ru": "женские свитеры, свитер, пуловеры, шерстяные свитеры, кашемир, женские пуловеры"
}'::jsonb WHERE id = 'dc26f29f-cfbc-4f8b-936b-e6f6d0397569';
