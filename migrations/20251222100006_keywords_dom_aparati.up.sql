-- Миграция: Добавление meta_keywords для категорий Dom i bašta и Kućni aparati
-- Создано: 2025-12-22
-- Описание: Заполнение SEO ключевых слов для всех L2/L3 подкатегорий

-- ============================================
-- КАТЕГОРИЯ: Kućni aparati (Бытовая техника)
-- ============================================

-- Frižideri (Холодильники)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'refrigerator, fridge, freezer, side by side, french door, mini fridge, no frost, energy efficient, built-in refrigerator, combined fridge, double door, stainless steel',
    'sr', 'frižider, zamrzivač, kombinovani frižider, no frost, side by side, mini frižider, ugradni, dupla vrata, inox, rashladni uređaj, french door, energetska klasa',
    'ru', 'холодильник, морозильник, двухкамерный, side by side, no frost, встраиваемый, мини холодильник, french door, нержавеющая сталь, энергоэффективный, комбинированный'
) WHERE slug = 'frizideri';

-- Zamrzivači vertikalni (Вертикальные морозильники)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'vertical freezer, upright freezer, chest freezer, no frost freezer, energy efficient freezer, deep freezer, drawer freezer, frost free',
    'sr', 'vertikalni zamrzivač, samostalni zamrzivač, no frost, zamrzivač sa fiokama, duboki zamrzivač, energetska klasa A, frost free, veliki zamrzivač',
    'ru', 'вертикальная морозильная камера, морозильник с ящиками, no frost, энергоэффективный, глубокая заморозка, frost free, большой морозильник'
) WHERE slug = 'zamrzivaci-vertikalni';

-- Vinske vitrine (Винные шкафы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'wine cooler, wine fridge, wine cabinet, wine cellar, dual zone wine cooler, temperature control, wine storage, climate controlled',
    'sr', 'vinska vitrina, frižider za vino, vinski podrum, temperaturna kontrola, hladnjak za vino, dva temperaturna zona, čuvanje vina, klimatizirano',
    'ru', 'винный шкаф, винный холодильник, винный погреб, контроль температуры, двухзонный, хранение вина, климатический шкаф'
) WHERE slug = 'vinske-vitrine';

-- Mašine za pranje (Стиральные машины)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'washing machine, front load washer, top load washer, automatic washing machine, energy efficient, washing machine 7kg, 8kg, 9kg, quick wash, silent wash',
    'sr', 'mašina za pranje veša, automatska mašina, frontalno punjenje, gornje punjenje, energetska klasa, 7kg, 8kg, 9kg, brzo pranje, tiho pranje, A+++',
    'ru', 'стиральная машина, автоматическая стиральная машина, фронтальная загрузка, вертикальная загрузка, энергоэффективная, 7кг, 8кг, 9кг, быстрая стирка'
) WHERE slug = 'masine-za-pranje';

-- Mašine za pranje i sušenje (Стирально-сушильные машины)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'washer dryer combo, washing drying machine, 2-in-1 washer dryer, combined washer dryer, all in one washer dryer, compact washer dryer',
    'sr', 'mašina za pranje i sušenje, kombinovana mašina, 2u1 mašina, kompaktna mašina, automatska sušilica, pranje i sušenje veša',
    'ru', 'стирально-сушильная машина, комбинированная машина, 2в1, компактная машина, автоматическая сушка, стирка и сушка'
) WHERE slug = 'masine-pranje-susenje';

-- Mašine za suđe (Посудомоечные машины)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'dishwasher, built-in dishwasher, freestanding dishwasher, compact dishwasher, 60cm dishwasher, 45cm dishwasher, energy efficient, silent dishwasher',
    'sr', 'mašina za suđe, ugradna mašina za suđe, samostalna, kompaktna, 60cm, 45cm, tiha mašina za suđe, energetska klasa, automatska mašina',
    'ru', 'посудомоечная машина, встраиваемая посудомоечная машина, отдельностоящая, компактная, 60см, 45см, тихая, энергоэффективная'
) WHERE slug = 'masine-sudje';

-- Sudopere i mašine (Мойки)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kitchen sink, stainless steel sink, granite sink, double sink, single sink, undermount sink, top mount sink, sink with drainer',
    'sr', 'sudopera, inox sudopera, granitna sudopera, dvodelna sudopera, ugradna sudopera, sa ocednom površinom, kuhinjska sudopera',
    'ru', 'кухонная мойка, мойка из нержавеющей стали, гранитная мойка, двойная мойка, врезная мойка, накладная мойка, мойка с крылом'
) WHERE slug = 'sudopere-i-masine';

-- Šporet i rerna (Плиты и духовки)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'stove, oven, electric stove, gas stove, combined stove, built-in oven, freestanding cooker, ceramic hob, induction hob, convection oven',
    'sr', 'šporet, rerna, električni šporet, plinski šporet, kombinovani šporet, ugradna rerna, samostalni šporet, staklokeramička ploča, indukciona ploča, konvekciona rerna',
    'ru', 'плита, духовка, электрическая плита, газовая плита, комбинированная плита, встраиваемая духовка, отдельностоящая плита, индукционная плита, конвекционная духовка'
) WHERE slug = 'sporet-i-rerna';

-- Indukcione ploče (Индукционные плиты)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'induction cooktop, induction hob, built-in induction, 4 zone induction, touch control, portable induction, energy efficient, fast heating',
    'sr', 'indukciona ploča, ugradna indukcija, 4 zone, senzorska kontrola, prenosiva indukcija, energetski efikasna, brzo zagrevanje, indukcionо kuvanje',
    'ru', 'индукционная плита, встраиваемая индукционная панель, 4 зоны, сенсорное управление, портативная индукция, энергоэффективная, быстрый нагрев'
) WHERE slug = 'indukcionе-ploce';

-- Mikrotalasne rerne (Микроволновые печи)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'microwave oven, microwave, built-in microwave, countertop microwave, grill microwave, convection microwave, solo microwave, inverter microwave',
    'sr', 'mikrotalasna rerna, mikrotalasna, ugradna mikrotalasna, samostalna, sa grilom, konvekciona, inverter mikrotalasna, solo mikrotalasna',
    'ru', 'микроволновая печь, микроволновка, встраиваемая микроволновка, с грилем, конвекционная, инверторная, соло микроволновка'
) WHERE slug = 'mikotalasne-rerne';

-- Parne rerne (Паровые духовки)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'steam oven, combi steam oven, built-in steam oven, convection steam, healthy cooking, steam cooking, combination oven',
    'sr', 'parna rerna, kombinovana parna rerna, ugradna parna rerna, konvekciona para, zdravo kuvanje, kuvanje na pari, kombinovana rerna',
    'ru', 'паровая духовка, комбинированная паровая духовка, встраиваемая паровая духовка, конвекционный пар, здоровое приготовление, готовка на пару'
) WHERE slug = 'parne-rerne';

-- Usisivači (Пылесосы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'vacuum cleaner, bagless vacuum, upright vacuum, canister vacuum, stick vacuum, handheld vacuum, wet dry vacuum, HEPA filter, cordless vacuum',
    'sr', 'usisivač, usisivač bez kese, vertikalni usisivač, ručni usisivač, bežični usisivač, usisivač sa vodom, HEPA filter, mokro suvo usisavanje',
    'ru', 'пылесос, пылесос без мешка, вертикальный пылесос, ручной пылесос, беспроводной пылесос, моющий пылесос, HEPA фильтр, сухая и влажная уборка'
) WHERE slug = 'usisivaci';

-- Aspiratori bez kesice (Безмешковые пылесосы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bagless vacuum, cyclone vacuum, multi-cyclone technology, transparent dustbin, washable filter, energy efficient vacuum, powerful suction',
    'sr', 'usisivač bez kese, ciklonski usisivač, multi-ciklon tehnologija, providna posuda, perivi filter, snažno usisavanje, energetski efikasan',
    'ru', 'безмешковый пылесос, циклонный пылесос, мультициклонная технология, прозрачный контейнер, моющийся фильтр, мощное всасывание'
) WHERE slug = 'aspiratori-bez-kesice';

-- Roboti usisivači (Роботы-пылесосы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'robot vacuum, robotic vacuum cleaner, smart vacuum, automatic vacuum, mapping robot, mopping robot, wifi vacuum, app controlled vacuum',
    'sr', 'robot usisivač, pametni usisivač, automatski usisivač, mapiranje prostora, usisivač sa brisanjem, wifi kontrola, aplikacija za upravljanje',
    'ru', 'робот-пылесос, умный пылесос, автоматический пылесос, картирование помещения, робот с влажной уборкой, wifi управление, управление через приложение'
) WHERE slug = 'roboti-usisivaci';

-- Ventilacija i klimatizacija (Вентиляция и кондиционирование)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'air conditioning, ventilation, AC unit, split system, portable AC, window AC, inverter AC, climate control, cooling heating',
    'sr', 'klimatizacija, ventilacija, klima uređaj, split sistem, pokretna klima, prozorska klima, inverter klima, kontrola klime, hlađenje grejanje',
    'ru', 'кондиционирование, вентиляция, кондиционер, сплит-система, мобильный кондиционер, оконный кондиционер, инверторный кондиционер, климат-контроль'
) WHERE slug = 'ventilacija-i-klimatizacija';

-- Ventilatori i grejalice (Вентиляторы и обогреватели)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'fan, heater, electric heater, tower fan, pedestal fan, ceiling fan, infrared heater, oil heater, convector heater, portable heater',
    'sr', 'ventilator, grejalica, električna grejalica, toranj ventilator, stalak ventilator, plafon ventilator, infracrvena grejalica, uljana grejalica, konvektor',
    'ru', 'вентилятор, обогреватель, электрический обогреватель, башенный вентилятор, напольный вентилятор, потолочный вентилятор, инфракрасный обогреватель, масляный обогреватель, конвектор'
) WHERE slug = 'ventilatori-i-grejalice';

-- Prečistači vazduha (Очистители воздуха)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'air purifier, air cleaner, HEPA air purifier, ionizer, humidifier purifier, smart air purifier, allergen remover, dust remover',
    'sr', 'prečistač vazduha, prečišćivač vazduha, HEPA prečistač, jonizator, ovlaživač sa prečistačem, pametni prečistač, alergeni, prečišćavanje vazduha',
    'ru', 'очиститель воздуха, воздухоочиститель, HEPA очиститель, ионизатор, увлажнитель с очистителем, умный очиститель, удаление аллергенов'
) WHERE slug = 'precistaci-vazduha';

-- Bojleri (Бойлеры)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'water heater, boiler, electric water heater, gas water heater, instant water heater, storage water heater, tankless water heater',
    'sr', 'bojler, grejač vode, električni bojler, gasni bojler, protočni bojler, akumulacioni bojler, trenutni zagrevač vode',
    'ru', 'бойлер, водонагреватель, электрический водонагреватель, газовый водонагреватель, проточный водонагреватель, накопительный водонагреватель'
) WHERE slug = 'bojleri';

-- Mali kućni aparati (Малая бытовая техника)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'small appliances, kitchen appliances, toaster, kettle, blender, mixer, coffee maker, food processor, hand mixer, immersion blender',
    'sr', 'mali kućni aparati, kuhinjski aparati, toster, kuvalo za vodu, blender, mikser, aparat za kafu, kuhinjski robot, ručni mikser, štapni mikser',
    'ru', 'малая бытовая техника, кухонная техника, тостер, чайник, блендер, миксер, кофеварка, кухонный комбайн, ручной миксер, погружной блендер'
) WHERE slug = 'mali-kucni-aparati';

-- Mašine za kafu (Кофемашины)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'coffee machine, espresso machine, coffee maker, automatic coffee machine, capsule coffee machine, bean to cup, drip coffee maker, turkish coffee machine',
    'sr', 'aparat za kafu, espresso aparat, automatski aparat za kafu, kafe aparat na kapsule, aparat za tursku kafu, filter kafa, aparat sa mlinom',
    'ru', 'кофемашина, эспрессо машина, кофеварка, автоматическая кофемашина, капсульная кофемашина, зерновая кофемашина, капельная кофеварка, турка электрическая'
) WHERE slug = 'masine-za-kafu';

-- Blenderi profesionalni (Профессиональные блендеры)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'professional blender, high speed blender, commercial blender, smoothie blender, powerful blender, vacuum blender, glass jar blender',
    'sr', 'profesionalni blender, snažni blender, komercijalni blender, blender za smuti, vakuum blender, blender sa staklenom posudom, high speed blender',
    'ru', 'профессиональный блендер, мощный блендер, коммерческий блендер, блендер для смузи, вакуумный блендер, блендер со стеклянной чашей'
) WHERE slug = 'blenderi-profesionalni';

-- Friteze (Фритюрницы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'deep fryer, electric fryer, oil fryer, commercial fryer, home fryer, stainless steel fryer, temperature control',
    'sr', 'friteza, električna friteza, friteza sa uljem, domaća friteza, inox friteza, kontrola temperature, duboko prženje',
    'ru', 'фритюрница, электрическая фритюрница, фритюрница с маслом, домашняя фритюрница, из нержавеющей стали, контроль температуры'
) WHERE slug = 'friteze';

-- Friteze sa vrućim vazduhom (Аэрофритюрницы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'air fryer, hot air fryer, oil free fryer, healthy fryer, digital air fryer, multi-function air fryer, air fryer oven',
    'sr', 'air fryer, friteza na vrući vazduh, friteza bez ulja, zdrava friteza, digitalna friteza, multi-funkciona friteza, air fryer rerna',
    'ru', 'аэрофритюрница, фритюрница с горячим воздухом, фритюрница без масла, здоровая фритюрница, цифровая фритюрница, мультифункциональная'
) WHERE slug = 'friteze-sa-vrucim-vazduhom';

-- Multi kukeri (Мультиварки)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'multi cooker, pressure cooker, slow cooker, rice cooker, instant pot, programmable cooker, electric pressure cooker',
    'sr', 'multi kuker, multifunkcionalni lonac, kuvalo pod pritiskom, sporогрејач, lonac za pirinač, programabilni kuker, električni lonac',
    'ru', 'мультиварка, скороварка, медленноварка, рисоварка, программируемая мультиварка, электрическая скороварка, instant pot'
) WHERE slug = 'multi-kukeri';

-- Instant lonci (Скороварки)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'instant pot, electric pressure cooker, programmable pressure cooker, multi-use cooker, fast cooking, one pot cooking',
    'sr', 'instant lonac, električni lonac pod pritiskom, programabilni lonac, multifunkcionalni lonac, brzo kuvanje, jedan lonac',
    'ru', 'скороварка электрическая, instant pot, программируемая скороварка, мультифункциональная скороварка, быстрое приготовление'
) WHERE slug = 'instant-lonci';

-- Kuhinjski pribor (Кухонные принадлежности)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kitchen utensils, cooking tools, spatula, whisk, ladle, tongs, kitchen gadgets, cooking accessories, utensil set',
    'sr', 'kuhinjski pribor, kuvarski alati, lopatica, mešač, kutlača, hvataljka, kuhinjski alati, pribor za kuvanje, set alata',
    'ru', 'кухонные принадлежности, кухонные инструменты, лопатка, венчик, половник, щипцы, кухонные гаджеты, набор инструментов'
) WHERE slug = 'kuhinjski-pribor';

-- Kuhinjski set noževa (Наборы ножей)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'knife set, kitchen knives, chef knife, stainless steel knives, ceramic knives, knife block, professional knives, sharp knives',
    'sr', 'set noževa, kuhinjski noževi, šef nož, inox noževi, keramički noževi, blok za noževe, profesionalni noževi, oštri noževi',
    'ru', 'набор ножей, кухонные ножи, шеф-нож, ножи из нержавеющей стали, керамические ножи, подставка для ножей, профессиональные ножи'
) WHERE slug = 'kuhinjski-set-nozeva';

-- Drveno posuđe za kuhinju (Деревянная посуда)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'wooden kitchenware, wooden utensils, wooden spoons, wooden cutting board, bamboo utensils, wooden bowls, eco friendly',
    'sr', 'drveno posuđe, drveni kuhinjski pribor, drvene kašike, daska za sečenje, bambus pribor, drvene činije, ekološki',
    'ru', 'деревянная посуда, деревянные кухонные принадлежности, деревянные ложки, разделочная доска, бамбуковые принадлежности, деревянные миски'
) WHERE slug = 'drveno-posude-kuhinja';

-- Keramičko posuđe ručno rađeno (Керамическая посуда ручной работы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'handmade ceramic, ceramic dishes, pottery, artisan ceramics, handcrafted bowls, ceramic plates, unique ceramics',
    'sr', 'ručno rađena keramika, keramičko posuđe, glina, zanatska keramika, ručno rađene činije, keramički tanjiri, jedinstvena keramika',
    'ru', 'керамика ручной работы, керамическая посуда, гончарные изделия, авторская керамика, керамические миски, керамические тарелки'
) WHERE slug = 'keramika-posude-rucno';

-- Staklo i kristal posuđe (Стекло и хрусталь)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'glassware, crystal glass, glass dishes, crystal dishes, wine glasses, champagne glasses, glass bowls, lead crystal',
    'sr', 'stakleno posuđe, kristal, stakleni tanjiri, kristalno posuđe, čaše za vino, čaše za šampanjac, staklene činije, olovni kristal',
    'ru', 'стеклянная посуда, хрусталь, стеклянные тарелки, хрустальная посуда, бокалы для вина, бокалы для шампанского, хрустальные миски'
) WHERE slug = 'staklo-kristal-posude';

-- Kozarci i čaše luksuzni (Бокалы премиум)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'luxury glasses, premium glassware, crystal glasses, wine glasses, whiskey glasses, champagne flutes, handmade glasses',
    'sr', 'luksuzni kozarci, premium čaše, kristalni kozarci, čaše za vino, čaše za viski, čaše za šampanjac, ručno rađeni kozarci',
    'ru', 'роскошные бокалы, премиум стекло, хрустальные бокалы, бокалы для вина, бокалы для виски, бокалы для шампанского, ручной работы'
) WHERE slug = 'kozarci-case-luksuz';

-- Servisi za kafu (Кофейные сервизы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'coffee set, espresso cups, coffee service, porcelain coffee set, turkish coffee set, coffee pot, coffee cups saucers',
    'sr', 'servis za kafu, espresso šoljice, set za kafu, porcelanski servis, servis za tursku kafu, džezva, šoljice sa tacnama',
    'ru', 'кофейный сервиз, чашки для эспрессо, кофейный набор, фарфоровый сервиз, сервиз для турецкого кофе, кофейник, чашки с блюдцами'
) WHERE slug = 'kafe-servisi';

-- Servisi za čaj (Чайные сервизы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tea set, tea service, porcelain tea set, teapot, tea cups, english tea set, chinese tea set, ceramic tea set',
    'sr', 'servis za čaj, set za čaj, porcelanski servis, čajnik, šoljice za čaj, engleski servis, kineski servis, keramički servis',
    'ru', 'чайный сервиз, чайный набор, фарфоровый сервиз, заварочный чайник, чашки для чая, английский сервиз, китайский сервиз'
) WHERE slug = 'caj-servisi';

-- Tacne i poslužavnici (Подносы и сервировочные блюда)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'serving tray, serving platter, wooden tray, metal tray, decorative tray, serving dishes, food platter, cake stand',
    'sr', 'tacna, poslužavnik, drvena tacna, metalna tacna, dekorativna tacna, posuđe za serviranje, tacna za hranu, stalak za torte',
    'ru', 'поднос, сервировочное блюдо, деревянный поднос, металлический поднос, декоративный поднос, блюда для подачи, подставка для торта'
) WHERE slug = 'tacne-posluzavnici';

-- Bar oprema za dom (Барная техника)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'home bar equipment, cocktail shaker, wine opener, bar tools, cocktail set, bartender kit, ice bucket, wine aerator',
    'sr', 'kućna bar oprema, šejker, otvarač za vino, bar alati, koktel set, bartender set, kofa za led, aerator za vino',
    'ru', 'домашнее барное оборудование, шейкер, штопор, барные инструменты, коктейльный набор, бармен набор, ведерко для льда, аэратор для вина'
) WHERE slug = 'bar-oprema-dom';

-- ============================================
-- КАТЕГОРИЯ: Dom i bašta (Дом и сад)
-- ============================================

-- Nameštaj dnevna soba (Мебель для гостиной)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'living room furniture, sofa, couch, sectional, coffee table, tv stand, armchair, bookshelf, entertainment center, corner sofa',
    'sr', 'nameštaj za dnevnu sobu, sofa, kauč, ugaona garnitura, klub sto, tv komoda, fotelja, polica za knjige, dnevna soba garnitura',
    'ru', 'мебель для гостиной, диван, угловой диван, журнальный столик, тумба под телевизор, кресло, книжный шкаф, стенка'
) WHERE slug = 'namestaj-dnevna-soba';

-- Nameštaj spavaća soba (Мебель для спальни)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bedroom furniture, bed, mattress, wardrobe, dresser, nightstand, bedside table, bed frame, king size bed, queen bed',
    'sr', 'nameštaj za spavaću sobu, krevet, dušek, ormar, komoda, noćni sto, okvir kreveta, bračni krevet, dvokrevetna garnitura',
    'ru', 'мебель для спальни, кровать, матрас, шкаф, комод, прикроватная тумбочка, каркас кровати, двуспальная кровать'
) WHERE slug = 'namestaj-spavaca-soba';

-- Nameštaj kuhinja (Кухонная мебель)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kitchen furniture, kitchen cabinets, dining table, kitchen island, bar stool, pantry, kitchen storage, modular kitchen',
    'sr', 'kuhinjski nameštaj, kuhinjski elementi, trpezarijski sto, kuhinjski otok, bar stolica, ostava, kuhinjsko skladište, modularna kuhinja',
    'ru', 'кухонная мебель, кухонные шкафы, обеденный стол, кухонный остров, барный стул, кладовая, кухонное хранение, модульная кухня'
) WHERE slug = 'namestaj-kuhinja';

-- Nameštaj kancelarija (Офисная мебель)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'office furniture, desk, office chair, ergonomic chair, filing cabinet, bookcase, standing desk, computer desk, executive desk',
    'sr', 'kancelarijski nameštaj, radni sto, kancelarijska stolica, ergonomska stolica, ormar za dokumenta, polica, sto za rad stojećи, kompjuterski sto',
    'ru', 'офисная мебель, письменный стол, офисное кресло, эргономичное кресло, шкаф для документов, книжная полка, стол для работы стоя'
) WHERE slug = 'namestaj-kancelarija';

-- Organizacija i skladištenje (Организация и хранение)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'storage solutions, organizer, closet organizer, storage bins, shelving units, pantry organizer, drawer organizer, storage boxes',
    'sr', 'organizacija prostora, organizer, ormar organizer, kutije za skladištenje, police, organizer za ostavu, organizer za fioke, skladišne kutije',
    'ru', 'системы хранения, органайзер, органайзер для шкафа, контейнеры для хранения, стеллажи, органайзер для кладовой, органайзер для ящиков'
) WHERE slug = 'organizacija-i-skladistenje';

-- Tekstil za dom (Домашний текстиль)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'home textiles, bedding, curtains, pillows, blankets, towels, bed linen, table cloth, cushion covers, duvet cover',
    'sr', 'tekstil za dom, posteljina, zavese, jastuci, ćebad, peškiri, posteljno rublje, stolnjak, jastučnice, navlaka za jorgan',
    'ru', 'домашний текстиль, постельное белье, шторы, подушки, одеяла, полотенца, столовое белье, скатерть, наволочки, пододеяльник'
) WHERE slug = 'tekstil-za-dom';

-- Tepisi i prostirke (Ковры и коврики)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'rugs, carpets, area rugs, runner rugs, Persian rugs, modern rugs, outdoor rugs, bathroom mat, doormat',
    'sr', 'tepisi, prostirke, tepih za dnevnu sobu, hodnik tepih, persijski tepisi, moderni tepisi, spoljni tepisi, kupatilska prostirka, otirač',
    'ru', 'ковры, ковровые дорожки, ковры для гостиной, персидские ковры, современные ковры, уличные коврики, коврик для ванной, придверный коврик'
) WHERE slug = 'tepihi-i-prostirke';

-- Rasveta (Освещение)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'lighting, lamps, ceiling lights, wall lights, floor lamps, table lamps, LED lights, chandelier, pendant lights, smart lighting',
    'sr', 'rasveta, lampe, plafon rasveta, zidna svetla, podne lampe, stona lampa, LED rasveta, luster, viseća svetla, pametna rasveta',
    'ru', 'освещение, лампы, потолочные светильники, настенные светильники, торшеры, настольные лампы, LED освещение, люстра, подвесные светильники'
) WHERE slug = 'rasveta';

-- Luster (Люстры)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'chandelier, crystal chandelier, modern chandelier, dining room chandelier, large chandelier, luxury chandelier, LED chandelier',
    'sr', 'luster, kristalni luster, moderni luster, luster za trpezariju, veliki luster, luksuzni luster, LED luster',
    'ru', 'люстра, хрустальная люстра, современная люстра, люстра для столовой, большая люстра, роскошная люстра, LED люстра'
) WHERE slug = 'luster';

-- Plafonjera (Потолочные светильники)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ceiling light, flush mount light, close to ceiling light, LED ceiling light, modern ceiling light, kitchen ceiling light',
    'sr', 'plafonjera, plafon svetlo, ugradna plafonjera, LED plafonjera, moderna plafonjera, plafonjera za kuhinju',
    'ru', 'потолочный светильник, накладной светильник, LED потолочный светильник, современный потолочный светильник, светильник для кухни'
) WHERE slug = 'plafonjera';

-- Zidna lampa (Настенные светильники)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'wall lamp, wall sconce, wall light, bedside wall lamp, decorative wall light, LED wall lamp, modern wall sconce',
    'sr', 'zidna lampa, zidno svetlo, zidna lampa za spavaću sobu, dekorativna zidna lampa, LED zidna lampa, moderna zidna lampa',
    'ru', 'настенный светильник, бра, прикроватный светильник, декоративный настенный светильник, LED настенный светильник'
) WHERE slug = 'zidna-lampa';

-- Stojeća lampa (Торшеры)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'floor lamp, standing lamp, tall lamp, reading lamp, arc floor lamp, tripod floor lamp, LED floor lamp, modern floor lamp',
    'sr', 'stojeća lampa, podna lampa, visoka lampa, lampa za čitanje, lučna lampa, tripod lampa, LED stojeća lampa, moderna stojeća lampa',
    'ru', 'торшер, напольная лампа, высокая лампа, лампа для чтения, арочный торшер, торшер на треноге, LED торшер'
) WHERE slug = 'stojeća-lampa';

-- Stona lampa (Настольные лампы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'table lamp, desk lamp, bedside lamp, reading lamp, LED table lamp, modern table lamp, decorative lamp, touch lamp',
    'sr', 'stona lampa, lampa za radni sto, noćna lampa, lampa za čitanje, LED stona lampa, moderna stona lampa, dekorativna lampa',
    'ru', 'настольная лампа, лампа для письменного стола, прикроватная лампа, лампа для чтения, LED настольная лампа, современная настольная лампа'
) WHERE slug = 'stonalampa';

-- LED sijalice (LED лампочки)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'LED bulb, LED light bulb, energy saving bulb, smart LED bulb, dimmable LED, E27 bulb, E14 bulb, warm white, cool white',
    'sr', 'LED sijalica, LED lampa, energetski efikasna sijalica, pametna LED sijalica, dimerizovana LED, E27 sijalica, E14, topla bela, hladna bela',
    'ru', 'LED лампочка, светодиодная лампа, энергосберегающая лампа, умная LED лампа, диммируемая LED, E27 лампа, E14, теплый белый, холодный белый'
) WHERE slug = 'led-sijalice';

-- Ugradna LED rasveta (Встраиваемое LED освещение)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'recessed LED lighting, downlight, ceiling recessed light, spotlights, under cabinet lighting, LED strip, dimmable downlight',
    'sr', 'ugradna LED rasveta, ugradni reflektor, plafon ugradna svetla, spot svetla, rasveta ispod elemenata, LED traka, dimerizovana ugradna',
    'ru', 'встраиваемое LED освещение, точечные светильники, потолочные встраиваемые светильники, споты, подсветка под шкафы, LED лента'
) WHERE slug = 'ugradna-led-rasveta';

-- Spoljna rasveta (Уличное освещение)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'outdoor lighting, garden lights, porch light, pathway lights, solar lights, motion sensor light, waterproof lights, LED outdoor',
    'sr', 'spoljna rasveta, baštenska svetla, svetlo za trem, staza rasveta, solarna svetla, svetlo sa senzorom, vodootporna svetla, LED spoljna',
    'ru', 'уличное освещение, садовые светильники, свет для крыльца, освещение дорожек, солнечные светильники, свет с датчиком движения, водонепроницаемые'
) WHERE slug = 'spoljna-rasveta';

-- Dekorativna rasveta (Декоративное освещение)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'decorative lighting, string lights, fairy lights, neon lights, decorative lamps, accent lighting, mood lighting, LED strips',
    'sr', 'dekorativna rasveta, svetleće girland, vilinske lampice, neon svetla, dekorativne lampe, ambient rasveta, LED trake za dekoraciju',
    'ru', 'декоративное освещение, гирлянды, светодиодные огни, неоновые светильники, декоративные лампы, акцентное освещение, LED ленты'
) WHERE slug = 'dekorativna-rasveta';

-- Pametna rasveta (Умное освещение)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'smart lighting, wifi bulbs, voice control lights, app controlled lighting, color changing lights, smart LED, alexa lights, google home',
    'sr', 'pametna rasveta, wifi sijalice, glasovna kontrola svetla, aplikacija za svetlo, promena boje svetla, pametna LED, alexa, google home',
    'ru', 'умное освещение, wifi лампочки, голосовое управление светом, управление через приложение, изменение цвета, умные LED, alexa, google home'
) WHERE slug = 'pametna-rasveta';

-- Dekoracije (Декор)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'home decor, wall decor, decorative items, home accessories, interior decoration, modern decor, vintage decor, handmade decor',
    'sr', 'dekoracije za dom, zidna dekoracija, dekorativni predmeti, kućni aksesoari, unutrašnja dekoracija, moderna dekoracija, vintage, ručno rađeno',
    'ru', 'декор для дома, настенный декор, декоративные предметы, домашние аксессуары, интерьерное украшение, современный декор, винтаж, ручная работа'
) WHERE slug = 'dekoracije';

-- Vaze i dekor (Вазы и декор)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'vases, decorative vases, flower vases, ceramic vases, glass vases, modern vases, large vases, floor vases, table vases',
    'sr', 'vaze, dekorativne vaze, vaze za cveće, keramičke vaze, staklene vaze, moderne vaze, velike vaze, podne vaze, stona vaza',
    'ru', 'вазы, декоративные вазы, вазы для цветов, керамические вазы, стеклянные вазы, современные вазы, большие вазы, напольные вазы'
) WHERE slug = 'vaze-i-dekor';

-- Ogledala (Зеркала)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'mirrors, wall mirror, floor mirror, decorative mirror, bathroom mirror, round mirror, framed mirror, LED mirror, full length mirror',
    'sr', 'ogledala, zidno ogledalo, podno ogledalo, dekorativno ogledalo, kupatilsko ogledalo, okruglo ogledalo, uramljeno ogledalo, LED ogledalo',
    'ru', 'зеркала, настенное зеркало, напольное зеркало, декоративное зеркало, зеркало для ванной, круглое зеркало, зеркало в раме, LED зеркало'
) WHERE slug = 'ogledala';

-- Satovi za zid (Настенные часы)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'wall clock, decorative clock, large wall clock, modern clock, vintage clock, silent clock, pendulum clock, wooden clock',
    'sr', 'zidni sat, dekorativni sat, veliki zidni sat, moderni sat, vintage sat, tihi sat, sat sa klatnom, drveni sat',
    'ru', 'настенные часы, декоративные часы, большие настенные часы, современные часы, винтажные часы, бесшумные часы, часы с маятником'
) WHERE slug = 'sat i-za-zid';

-- Pregradni zidovi (Перегородки)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'room divider, partition wall, folding screen, privacy screen, decorative partition, sliding partition, office partition',
    'sr', 'pregradni zid, separator prostora, paravent, pregradna vrata, dekorativna pregrada, klizna pregrada, kancelarijska pregrada',
    'ru', 'комнатная перегородка, ширма, раскладная ширма, декоративная перегородка, раздвижная перегородка, офисная перегородка'
) WHERE slug = 'pregradni-zidovi';

-- Kupatilo (Ванная комната)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bathroom, bathroom furniture, sink, toilet, shower, bathtub, bathroom accessories, bathroom fixtures, modern bathroom',
    'sr', 'kupatilo, kupatilski nameštaj, lavabo, wc šolja, tuš, kada, kupatilski aksesoari, sanitarije, moderno kupatilo',
    'ru', 'ванная комната, мебель для ванной, раковина, унитаз, душ, ванна, аксессуары для ванной, сантехника, современная ванная'
) WHERE slug = 'kupatilo';

-- WC šolja (Унитаз)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'toilet, WC, wall hung toilet, floor standing toilet, bidet toilet, smart toilet, dual flush toilet, rimless toilet',
    'sr', 'wc šolja, toalet, viseća wc šolja, podna wc, bidet wc, pametna wc, duplo ispiranje, bez obruča',
    'ru', 'унитаз, подвесной унитаз, напольный унитаз, унитаз с биде, умный унитаз, двойной слив, безободковый унитаз'
) WHERE slug = 'wc-solja';

-- Bide (Биде)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bidet, wall hung bidet, floor standing bidet, bidet attachment, electronic bidet, luxury bidet',
    'sr', 'bide, viseći bide, podni bide, dodatak za bide, elektronski bide, luksuzni bide',
    'ru', 'биде, подвесное биде, напольное биде, насадка биде, электронное биде, роскошное биде'
) WHERE slug = 'bidе';

-- Lavabo (Раковина)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bathroom sink, washbasin, lavabo, vessel sink, pedestal sink, wall mounted sink, countertop sink, ceramic sink',
    'sr', 'lavabo, umivaonik, sudopera za kupatilo, nadgradni lavabo, postolj lavabo, zidni lavabo, pult lavabo, keramički lavabo',
    'ru', 'раковина, умывальник, настольная раковина, раковина с пьедесталом, подвесная раковина, накладная раковина, керамическая раковина'
) WHERE slug = 'lavabo';

-- Ugradni lavabo (Встраиваемая раковина)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'built-in sink, undermount sink, drop-in sink, recessed sink, countertop integrated sink, flush mount sink',
    'sr', 'ugradni lavabo, podgradni lavabo, umet lavabo, integrisani lavabo, ravno ugradni lavabo',
    'ru', 'встраиваемая раковина, врезная раковина, подстольная раковина, интегрированная раковина, накладная раковина'
) WHERE slug = 'ugradni-lavabo';

-- Slavina za lavabo (Смеситель для раковины)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bathroom faucet, basin mixer, lavabo tap, single handle faucet, waterfall faucet, modern faucet, chrome faucet',
    'sr', 'slavina za lavabo, mešač za umivaonik, baterija za lavabo, jednoručna slavina, vodopad slavina, moderna slavina, hrom slavina',
    'ru', 'смеситель для раковины, кран для умывальника, однорычажный смеситель, водопад смеситель, современный смеситель, хромированный'
) WHERE slug = 'slavina-za-lavabo';

-- Kada (Ванна)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bathtub, bath, freestanding bathtub, built-in bathtub, acrylic bathtub, whirlpool bath, soaking tub, corner bathtub',
    'sr', 'kada, samostojeća kada, ugradna kada, akrilna kada, hidromasažna kada, kada za kupanje, ugaona kada',
    'ru', 'ванна, отдельностоящая ванна, встраиваемая ванна, акриловая ванна, гидромассажная ванна, угловая ванна'
) WHERE slug = 'kada';

-- Hidromasažna kada (Гидромассажная ванна)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'whirlpool bath, jacuzzi, hydromassage bathtub, spa bath, jet tub, massage tub, luxury bath',
    'sr', 'hidromasažna kada, džakuzi, masažna kada, spa kada, mlazna kada, luksuzna kada',
    'ru', 'гидромассажная ванна, джакузи, спа ванна, массажная ванна, ванна с форсунками, роскошная ванна'
) WHERE slug = 'hidromasazna-kada';

-- Slavina za kadu (Смеситель для ванны)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bath faucet, bathtub mixer, tub filler, bath tap, freestanding bath faucet, wall mounted bath faucet, waterfall bath faucet',
    'sr', 'slavina za kadu, mešač za kadu, baterija za kadu, samostojeća slavina, zidna slavina za kadu, vodopad slavina',
    'ru', 'смеситель для ванны, кран для ванны, напольный смеситель, настенный смеситель для ванны, водопад смеситель'
) WHERE slug = 'slavina-za-kadu';

-- Tuš kabina (Душевая кабина)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'shower cabin, shower enclosure, walk-in shower, glass shower, steam shower, corner shower, quadrant shower',
    'sr', 'tuš kabina, tuš paravan, walk-in tuš, staklena kabina, parna kabina, ugaona tuš kabina, četvrtasti tuš',
    'ru', 'душевая кабина, душевое ограждение, walk-in душ, стеклянная кабина, паровая кабина, угловая душевая кабина'
) WHERE slug = 'tus-kabina';

-- Tuš set (Душевой гарнитур)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'shower set, shower system, rain shower, handheld shower, thermostatic shower, shower mixer, modern shower set',
    'sr', 'tuš set, tuš sistem, kiša tuš, ručni tuš, termostatski tuš, tuš mešač, moderni tuš set',
    'ru', 'душевой гарнитур, душевая система, тропический душ, ручной душ, термостатический душ, смеситель для душа'
) WHERE slug = 'tus-set';

-- Kupatilski nameštaj (Мебель для ванной)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bathroom furniture, vanity unit, bathroom cabinet, storage cabinet, wall hung vanity, floor standing vanity, modern bathroom furniture',
    'sr', 'kupatilski nameštaj, ormar sa umivaonikom, kupatilski ormar, ormar za skladištenje, viseći ormar, podni ormar, moderni kupatilski nameštaj',
    'ru', 'мебель для ванной, тумба с раковиной, шкаф для ванной, шкаф для хранения, подвесная тумба, напольная тумба, современная мебель'
) WHERE slug = 'kupatilski-namestaj';

-- Ogledalo sa ormarićem (Зеркало с шкафчиком)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'mirror cabinet, bathroom mirror with storage, medicine cabinet, LED mirror cabinet, wall mounted mirror cabinet',
    'sr', 'ogledalo sa ormarićem, kupatilsko ogledalo sa skladištem, ormarić za lekove, LED ogledalo sa ormarićem, zidni ormarić',
    'ru', 'зеркальный шкаф, зеркало для ванной с полками, аптечка, LED зеркальный шкаф, настенный зеркальный шкаф'
) WHERE slug = 'ogledalo-sa-ormarićem';

-- Baštenka garnitura (Садовая мебель)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'garden furniture, outdoor furniture, patio furniture, garden set, rattan furniture, wooden garden furniture, sun loungers',
    'sr', 'baštenka garnitura, spoljašnji nameštaj, terasa nameštaj, baštenski set, ratan nameštaj, drveni baštanski nameštaj, ležaljke',
    'ru', 'садовая мебель, уличная мебель, мебель для патио, садовый набор, ротанговая мебель, деревянная садовая мебель, шезлонги'
) WHERE slug = 'bastenska-garnitura';

-- Baštenka oprema (Садовое оборудование)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'garden equipment, lawn mower, trimmer, hedge trimmer, leaf blower, chainsaw, pressure washer, garden tools',
    'sr', 'baštenka oprema, kosilica za travu, trimer, makaze za živicu, puhač lišća, motorna testera, perač pod pritiskom, baštanski alati',
    'ru', 'садовое оборудование, газонокосилка, триммер, кусторез, воздуходувка, бензопила, мойка высокого давления, садовые инструменты'
) WHERE slug = 'bastenska-oprema';

-- Alati za bašta (Садовые инструменты)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'garden tools, shovel, rake, hoe, pruning shears, garden fork, watering can, wheelbarrow, hand tools',
    'sr', 'alati za baštu, ašov, grablje, motika, makaze za orezivanje, viljuška, kanta za zalivanje, kolica, ručni alati',
    'ru', 'садовые инструменты, лопата, грабли, мотыга, секатор, вилы, лейка, тачка, ручные инструменты'
) WHERE slug = 'alati-za-basta';

-- Alati i popravke (Инструменты и ремонт)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tools, power tools, hand tools, drill, saw, hammer, screwdriver, tool set, repair tools, DIY tools',
    'sr', 'alati, električni alati, ručni alati, bušilica, testera, čekić, šrafciger, set alata, alati za popravke, DIY alati',
    'ru', 'инструменты, электроинструменты, ручные инструменты, дрель, пила, молоток, отвертка, набор инструментов, ремонтные инструменты'
) WHERE slug = 'alati-i-popravke';

-- Baštenka rasveta (Садовое освещение)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'garden lighting, outdoor lights, solar garden lights, pathway lights, landscape lighting, garden spotlights, LED outdoor lights',
    'sr', 'baštenka rasveta, spoljna svetla, solarna baštenka svetla, staza rasveta, pejzažna rasveta, baštenski reflektori, LED spoljna svetla',
    'ru', 'садовое освещение, уличное освещение, солнечные садовые светильники, освещение дорожек, ландшафтное освещение, садовые прожекторы'
) WHERE slug = 'bastenska-rasveta';

-- Baštenke ukrase (Садовые украшения)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'garden decorations, garden ornaments, garden statues, wind chimes, garden figures, decorative stones, bird bath, sundial',
    'sr', 'baštenke ukrase, baštenski ukrasi, baštenke statue, zvona vetra, baštenke figure, dekorativno kamenje, kupatilo za ptice, sunčani sat',
    'ru', 'садовые украшения, садовые статуи, музыка ветра, садовые фигуры, декоративные камни, купальня для птиц, солнечные часы'
) WHERE slug = 'bastenske-ukrase';

-- Grnčarija i biljke (Горшки и растения)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'pots and plants, flowerpots, planters, ceramic pots, terracotta pots, indoor plants, outdoor plants, decorative pots',
    'sr', 'grnčarija i biljke, saksije za cveće, posude za sadnju, keramičke saksije, terakota saksije, sobne biljke, spoljašnje biljke, dekorativne saksije',
    'ru', 'горшки и растения, цветочные горшки, кашпо, керамические горшки, терракотовые горшки, комнатные растения, уличные растения'
) WHERE slug = 'grncari ja-i-biljke';

-- Kompostiranje (Компостирование)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'composting, compost bin, composter, organic waste, garden compost, compost tumbler, worm composting, kitchen compost',
    'sr', 'kompostiranje, kontejner za kompost, komposter, organski otpad, baštenski kompost, rotacioni komposter, komposter sa crvima, kuhinjski kompost',
    'ru', 'компостирование, компостный ящик, компостер, органические отходы, садовый компост, барабанный компостер, вермикомпостирование'
) WHERE slug = 'kompostiranje';

-- Bazeni i spa (Бассейны и спа)
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'swimming pool, pool, spa, hot tub, jacuzzi, inflatable pool, above ground pool, pool accessories, pool cleaning',
    'sr', 'bazen, plivački bazen, spa, hidromasažna kada, džakuzi, naduvavajući bazen, nadzemni bazen, pribor za bazen, čišćenje bazena',
    'ru', 'бассейн, плавательный бассейн, спа, гидромассажная ванна, джакузи, надувной бассейн, наземный бассейн, аксессуары для бассейна'
) WHERE slug = 'bazeni-i-spa';

-- Завершение миграции
