-- Migration: Add meta_keywords for Lepota i zdravlje, Za bebe i decu, Nakit i satovi categories
-- Date: 2025-12-22
-- Description: SEO keywords in sr/en/ru for L2/L3 beauty, kids, jewelry categories

BEGIN;

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Nega kože
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'anti-aging, anti-wrinkle, rejuvenation, retinol, collagen, youth serum, lifting',
    'sr', 'protiv starenja, protiv bora, pomlađivanje, retinol, kolagen, serum mladosti, lifting',
    'ru', 'антивозрастной, против морщин, омоложение, ретинол, коллаген, сыворотка молодости, лифтинг'
) WHERE slug = 'anti-aging';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'anti-age cream, wrinkle cream, night cream, day cream, firming cream, rejuvenating cream',
    'sr', 'anti-age krema, krema protiv bora, noćna krema, dnevna krema, krema za zatezanje, krema za pomlađivanje',
    'ru', 'антивозрастной крем, крем от морщин, ночной крем, дневной крем, укрепляющий крем, омолаживающий крем'
) WHERE slug = 'anti-age-krema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'facial care, skin care, moisturizer, cleanser, toner, face cream, serum',
    'sr', 'nega lica, nega kože, hidratacija, sredstvo za čišćenje, tonik, krema za lice, serum',
    'ru', 'уход за лицом, уход за кожей, увлажнение, очищающее средство, тоник, крем для лица, сыворотка'
) WHERE slug = 'nega-koze';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'face cream, moisturizing cream, nourishing cream, hydrating cream, facial moisturizer',
    'sr', 'krema za lice, hidratantna krema, hranljiva krema, krema za hidrataciju, negovatelj lica',
    'ru', 'крем для лица, увлажняющий крем, питательный крем, гидратирующий крем, увлажнитель для лица'
) WHERE slug = 'krema-za-lice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'face serum, vitamin C serum, hyaluronic acid, brightening serum, anti-aging serum, concentrated treatment',
    'sr', 'serum za lice, serum sa vitaminom C, hijaluronska kiselina, serum za osvetljavanje, anti-age serum, koncentrisani tretman',
    'ru', 'сыворотка для лица, сыворотка с витамином C, гиалуроновая кислота, осветляющая сыворотка, антивозрастная сыворотка, концентрированное средство'
) WHERE slug = 'serum-za-lice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'face mask, sheet mask, clay mask, peel-off mask, sleeping mask, purifying mask, hydrating mask',
    'sr', 'maska za lice, sheet maska, glinena maska, maska koja se skida, noćna maska, maska za prečišćavanje, hidratantna maska',
    'ru', 'маска для лица, тканевая маска, глиняная маска, маска-пленка, ночная маска, очищающая маска, увлажняющая маска'
) WHERE slug = 'maska-za-lice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'toner, facial toner, refreshing toner, exfoliating toner, balancing toner, pH balancing',
    'sr', 'tonik, tonik za lice, osvežavajući tonik, piling tonik, balansirajući tonik, pH balansiranje',
    'ru', 'тоник, тоник для лица, освежающий тоник, отшелушивающий тоник, балансирующий тоник, pH балансирование'
) WHERE slug = 'tonik-za-lice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'cleansing milk, makeup remover, gentle cleanser, milky cleanser, face milk, cleansing lotion',
    'sr', 'mleko za čišćenje, sredstvo za skidanje šminke, blago sredstvo za čišćenje, mlečno sredstvo, mleko za lice, losion za čišćenje',
    'ru', 'молочко для очищения, средство для снятия макияжа, мягкое очищающее средство, молочко для лица, очищающий лосьон'
) WHERE slug = 'mleko-za-ciscenje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'body cream, body lotion, moisturizer, nourishing cream, firming body cream, body butter',
    'sr', 'krema za telo, losion za telo, hidratacija, hranljiva krema, krema za zatezanje tela, puter za telo',
    'ru', 'крем для тела, лосьон для тела, увлажнитель, питательный крем, укрепляющий крем для тела, масло для тела'
) WHERE slug = 'krema-za-telo';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hand cream, nourishing hand cream, protective hand cream, anti-aging hand cream, cuticle cream',
    'sr', 'krema za ruke, hranljiva krema za ruke, zaštitna krema za ruke, anti-age krema za ruke, krema za zanokitice',
    'ru', 'крем для рук, питательный крем для рук, защитный крем для рук, антивозрастной крем для рук, крем для кутикулы'
) WHERE slug = 'krema-za-ruke';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'sun protection, sunscreen, SPF cream, sunblock, tanning cream, sun lotion, UV protection',
    'sr', 'zaštita od sunca, krema za sunčanje, SPF krema, bloker za sunce, krema za preplanulost, losion za sunce, UV zaštita',
    'ru', 'защита от солнца, солнцезащитный крем, крем SPF, солнцезащитный блокатор, крем для загара, лосьон для загара, UV защита'
) WHERE slug = 'zastita-od-sunca';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'sunscreen, suntan lotion, tanning oil, after sun, beach cream, waterproof sunscreen',
    'sr', 'krema za sunčanje, losion za sunčanje, ulje za preplanulost, krema posle sunca, krema za plažu, vodootporna krema',
    'ru', 'крем для загара, лосьон для загара, масло для загара, средство после загара, пляжный крем, водостойкий крем'
) WHERE slug = 'krema-za-suncanje';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Higijena i tuš
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hygiene, personal care, intimate hygiene, body wash, soap, cleanliness, sanitation',
    'sr', 'higijena, lična nega, intimna higijena, gel za tuširanje, sapun, čistoća, sanitacija',
    'ru', 'гигиена, личная гигиена, интимная гигиена, гель для душа, мыло, чистота, санитария'
) WHERE slug = 'higijena';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'shower gel, body wash, bath gel, liquid soap, moisturizing shower gel, refreshing gel',
    'sr', 'gel za tuširanje, gel za kupanje, tečni sapun, hidratantni gel za tuširanje, osvežavajući gel',
    'ru', 'гель для душа, гель для купания, жидкое мыло, увлажняющий гель для душа, освежающий гель'
) WHERE slug = 'gelovi-za-tusiranje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'shower gel, body wash gel, bathing gel, cleansing gel, aromatic shower gel',
    'sr', 'gel za tuširanje, gel za telo, gel za kupanje, gel za čišćenje, aromatični gel za tuširanje',
    'ru', 'гель для душа, гель для тела, гель для купания, очищающий гель, ароматный гель для душа'
) WHERE slug = 'gel-za-tusiranje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'intimate hygiene, feminine wash, pH balanced, intimate gel, daily care, freshness',
    'sr', 'intimna higijena, gel za intimnu negu, pH balansiran, intimni gel, dnevna nega, svežina',
    'ru', 'интимная гигиена, гель для интимной гигиены, pH сбалансированный, интимный гель, ежедневный уход, свежесть'
) WHERE slug = 'intimna-higijena';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'oral hygiene, toothpaste, mouthwash, dental floss, teeth whitening, fresh breath, dental care',
    'sr', 'oralna higijena, pasta za zube, vodica za ispiranje usta, konac za zube, izbeljivanje zuba, svež dah, dentalna nega',
    'ru', 'оральная гигиена, зубная паста, ополаскиватель, зубная нить, отбеливание зубов, свежее дыхание, уход за зубами'
) WHERE slug = 'oralna-higijena';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Dekorativna kozmetika
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'decorative cosmetics, makeup, beauty products, face makeup, eye makeup, lip makeup, cosmetics',
    'sr', 'dekorativna kozmetika, šminka, beauty proizvodi, šminka za lice, šminka za oči, šminka za usne, kozmetika',
    'ru', 'декоративная косметика, макияж, косметика, макияж для лица, макияж для глаз, макияж для губ, косметические средства'
) WHERE slug = 'dekorativna-kozmetika';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'makeup brushes, cosmetic brushes, foundation brush, eyeshadow brush, blending brush, brush set',
    'sr', 'makeup četkice, kozmetičke četkice, četkica za puder, četkica za senke, četkica za blendovanje, set četkica',
    'ru', 'кисти для макияжа, косметические кисти, кисть для тонального крема, кисть для теней, кисть для растушевки, набор кистей'
) WHERE slug = 'makeup-cetkice';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Nega kose
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hair care, shampoo, conditioner, hair mask, hair oil, styling products, hair treatment',
    'sr', 'nega kose, šampon, regenerator, maska za kosu, ulje za kosu, styling proizvodi, tretman za kosu',
    'ru', 'уход за волосами, шампунь, кондиционер, маска для волос, масло для волос, средства для укладки, уход за волосами'
) WHERE slug = 'nega-kose';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hair dryer, professional hair dryer, blow dryer, ionic dryer, salon dryer, fast drying',
    'sr', 'fen, profesionalni fen, sušilo za kosu, jonski fen, salonski fen, brzo sušenje',
    'ru', 'фен, профессиональный фен, фен для волос, ионный фен, салонный фен, быстрая сушка'
) WHERE slug = 'fene-profesionalni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hair straightener, flat iron, tourmaline iron, ceramic plates, hair styling, straightening iron',
    'sr', 'pegla za kosu, ravnalica, turmalinska pegla, keramičke ploče, styling kose, presa za kosu',
    'ru', 'утюжок для волос, выпрямитель, турмалиновый утюжок, керамические пластины, укладка волос, выпрямитель для волос'
) WHERE slug = 'pegla-kosa-turmalin';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ionic hair brush, heated brush, straightening brush, styling brush, anti-frizz brush',
    'sr', 'jonska četkica za kosu, grejana četkica, četkica za ravnanje, styling četkica, četkica protiv pušenja',
    'ru', 'ионная расческа, нагреваемая расческа, выпрямляющая расческа, стайлинг расческа, расческа против пушистости'
) WHERE slug = 'ionske-cetkice-kosa';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Muška nega
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'men grooming, men care, beard care, shaving, aftershave, men cosmetics, male skincare',
    'sr', 'muška nega, nega za muškarce, nega brade, brijanje, losion posle brijanja, muška kozmetika, nega kože za muškarce',
    'ru', 'мужской уход, уход за мужчинами, уход за бородой, бритье, лосьон после бритья, мужская косметика, уход за кожей для мужчин'
) WHERE slug = 'muska-nega';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'men style, men fashion, men accessories, grooming tools, beard trimmer, electric razor',
    'sr', 'muški stil, muška moda, muški aksesoari, alati za negu, trimer za bradu, električni brijač',
    'ru', 'мужской стиль, мужская мода, мужские аксессуары, инструменты для ухода, триммер для бороды, электробритва'
) WHERE slug = 'muski-stil';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Parfemi
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'perfume, fragrance, eau de parfum, cologne, eau de toilette, scent, luxury perfume, designer fragrance',
    'sr', 'parfem, miris, eau de parfum, kolonjska voda, eau de toilette, miris, luksuzni parfem, dizajnerski miris',
    'ru', 'парфюм, аромат, парфюмерная вода, одеколон, туалетная вода, запах, роскошный парфюм, дизайнерский аромат'
) WHERE slug = 'parfemi';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Luksuzna kozmetika
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'luxury cosmetics, premium beauty, high-end makeup, designer cosmetics, luxury skincare, prestige beauty',
    'sr', 'luksuzna kozmetika, premium beauty, visoko kvalitetna šminka, dizajnerska kozmetika, luksuzna nega kože, prestižna lepota',
    'ru', 'люксовая косметика, премиум косметика, высококачественный макияж, дизайнерская косметика, люксовый уход за кожей, престижная красота'
) WHERE slug = 'luksuzna-kozmetika';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'organic cosmetics, natural beauty, eco cosmetics, vegan makeup, chemical-free, bio products',
    'sr', 'organska kozmetika, prirodna lepota, eko kozmetika, veganska šminka, bez hemikalija, bio proizvodi',
    'ru', 'органическая косметика, натуральная косметика, эко косметика, веганский макияж, без химикатов, био продукты'
) WHERE slug = 'organska-kozmetika';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Beauty aparati
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'beauty devices, skin care devices, facial tools, beauty tech, home salon, professional tools',
    'sr', 'aparati za lepotu, uređaji za negu kože, alati za lice, beauty tehnologija, kućni salon, profesionalni alati',
    'ru', 'устройства для красоты, приборы для ухода за кожей, инструменты для лица, бьюти технологии, домашний салон, профессиональные инструменты'
) WHERE slug = 'lepota-aparati';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'dermaroller, microneedling, skin needling, collagen induction, facial roller, anti-aging tool',
    'sr', 'dermaroller, mikroiglice, tretman iglicama, indukcija kolagena, rola za lice, anti-age alat',
    'ru', 'дермароллер, микроиглы, мезороллер, индукция коллагена, ролик для лица, антивозрастной инструмент'
) WHERE slug = 'dermaroller-mikronedling';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'LED face mask, light therapy, phototherapy, red light therapy, blue light, anti-aging mask',
    'sr', 'LED maska za lice, svetlosna terapija, fototerapija, crvena svetlost, plava svetlost, anti-age maska',
    'ru', 'LED маска для лица, светотерапия, фототерапия, красный свет, синий свет, антивозрастная маска'
) WHERE slug = 'led-maske-lice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ultrasonic cleaner, ultrasonic facial, deep cleansing, pore cleanser, skin scrubber, sonic cleaning',
    'sr', 'ultrazvučni čistač, ultrazvučno čišćenje lica, dubinsko čišćenje, čistač pora, uređaj za čišćenje kože, sonično čišćenje',
    'ru', 'ультразвуковой очиститель, ультразвуковая чистка лица, глубокая очистка, очиститель пор, скраббер для кожи, звуковая очистка'
) WHERE slug = 'ultrazvucni-cistaci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'laser hair removal, IPL, permanent hair removal, home laser, hair reduction, smooth skin',
    'sr', 'lasersko uklanjanje dlaka, IPL, trajno uklanjanje dlaka, kućni laser, redukcija dlaka, glatka koža',
    'ru', 'лазерная эпиляция, IPL, постоянное удаление волос, домашний лазер, уменьшение волос, гладкая кожа'
) WHERE slug = 'laseri-uklanjanje-dlaka';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'light therapy, LED therapy, chromotherapy, color therapy, healing light, wellness device',
    'sr', 'terapija svetlom, LED terapija, hromoterapija, terapija bojama, isceljujuća svetlost, wellness uređaj',
    'ru', 'светотерапия, LED терапия, хромотерапия, цветотерапия, исцеляющий свет, велнес устройство'
) WHERE slug = 'terapije-svetlom';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Depilacija
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hair removal, depilation, waxing, epilator, shaving, smooth skin, hair removal cream',
    'sr', 'uklanjanje dlaka, depilacija, vosak, epilator, brijanje, glatka koža, krema za depilaciju',
    'ru', 'удаление волос, депиляция, воск, эпилятор, бритье, гладкая кожа, крем для депиляции'
) WHERE slug = 'depilacija';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Manikir i pedikir
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'manicure, pedicure, nail care, nail polish, gel nails, nail art, nail tools',
    'sr', 'manikir, pedikir, nega noktiju, lak za nokte, gel nokti, nail art, alati za nokte',
    'ru', 'маникюр, педикюр, уход за ногтями, лак для ногтей, гель лак, нейл арт, инструменты для ногтей'
) WHERE slug = 'manikir-i-pedikir';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Spa i relaksacija
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'spa, relaxation, massage, aromatherapy, wellness, self-care, home spa, relaxing bath',
    'sr', 'spa, relaksacija, masaža, aromaterapija, wellness, briga o sebi, kućni spa, opuštajuće kupanje',
    'ru', 'спа, релаксация, массаж, ароматерапия, велнес, забота о себе, домашний спа, расслабляющая ванна'
) WHERE slug = 'spa-i-relaksacija';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Vitamini i suplementi
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'vitamins, supplements, dietary supplements, multivitamins, minerals, health supplements, immunity',
    'sr', 'vitamini, suplementi, dodaci ishrani, multivitamini, minerali, suplementi za zdravlje, imunitet',
    'ru', 'витамины, добавки, БАДы, мультивитамины, минералы, добавки для здоровья, иммунитет'
) WHERE slug = 'vitamini-i-suplementi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'collagen, collagen powder, hydrolyzed collagen, marine collagen, skin health, joint health',
    'sr', 'kolagen, kolagen u prahu, hidrolizovani kolagen, morski kolagen, zdravlje kože, zdravlje zglobova',
    'ru', 'коллаген, коллаген в порошке, гидролизованный коллаген, морской коллаген, здоровье кожи, здоровье суставов'
) WHERE slug = 'kolagen-preparati';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'biotin, biotin complex, hair vitamins, nail vitamins, vitamin B7, hair growth, beauty supplement',
    'sr', 'biotin, biotin kompleks, vitamini za kosu, vitamini za nokte, vitamin B7, rast kose, suplement za lepotu',
    'ru', 'биотин, биотиновый комплекс, витамины для волос, витамины для ногтей, витамин B7, рост волос, добавка для красоты'
) WHERE slug = 'biotin-kompleksi';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Eterična ulja
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'essential oils, aromatherapy oils, pure essential oils, lavender oil, eucalyptus, tea tree oil',
    'sr', 'eterična ulja, ulja za aromaterapiju, čista eterična ulja, lavanda ulje, eukaliptus, čajevac ulje',
    'ru', 'эфирные масла, масла для ароматерапии, чистые эфирные масла, лавандовое масло, эвкалипт, масло чайного дерева'
) WHERE slug = 'eterična-ulja';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Medicinski proizvodi
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'medical products, health devices, medical equipment, first aid, blood pressure monitor, thermometer',
    'sr', 'medicinski proizvodi, uređaji za zdravlje, medicinska oprema, prva pomoć, merač pritiska, termometar',
    'ru', 'медицинские изделия, устройства для здоровья, медицинское оборудование, первая помощь, тонометр, термометр'
) WHERE slug = 'medicinski-proizvodi';

-- ============================================================================
-- LEPOTA I ZDRAVLJE - Dečija kozmetika
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kids cosmetics, children makeup, safe cosmetics, hypoallergenic, gentle products, baby care',
    'sr', 'dečija kozmetika, šminka za decu, bezbedna kozmetika, hipoalergena, nežni proizvodi, nega beba',
    'ru', 'детская косметика, макияж для детей, безопасная косметика, гипоаллергенная, мягкие продукты, уход за младенцами'
) WHERE slug = 'decija-kozmetika';

-- ============================================================================
-- ZA BEBE I DECU - Oprema za bebe
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'baby gear, baby equipment, baby essentials, nursery equipment, infant products, newborn gear',
    'sr', 'oprema za bebe, bebi oprema, osnovna oprema, oprema za dečiju sobu, proizvodi za bebe, oprema za novorođenče',
    'ru', 'товары для младенцев, детское оборудование, необходимые товары, оборудование для детской, товары для новорожденных'
) WHERE slug = 'oprema-za-bebe';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'baby monitor, video monitor, audio monitor, baby camera, night vision, two-way audio, WiFi monitor',
    'sr', 'bebi monitor, video monitor, audio monitor, kamera za bebu, noćni vid, dvosmerni zvuk, WiFi monitor',
    'ru', 'радионяня, видеоняня, аудио монитор, камера для ребенка, ночное видение, двусторонняя связь, WiFi монитор'
) WHERE slug = 'bebi-monitori-video';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'automatic cradle, smart cradle, rocking cradle, electric rocker, soothing motion, baby swing',
    'sr', 'automatska kolevka, pametna kolevka, ljuljajuća kolevka, električna ljuljaška, umirujući pokret, ljuljaška za bebe',
    'ru', 'автоматическая колыбель, умная колыбель, качающаяся колыбель, электрическая качалка, успокаивающее движение, качели для младенцев'
) WHERE slug = 'kolevke-automatske';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bottle warmer, baby bottle heater, milk warmer, portable warmer, quick warming, temperature control',
    'sr', 'grejač flašica, grejač za bebe, grejač mleka, prenosivi grejač, brzo zagrevanje, kontrola temperature',
    'ru', 'подогреватель бутылочек, нагреватель молока, портативный подогреватель, быстрый подогрев, контроль температуры'
) WHERE slug = 'grejaci-flasica';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'sterilizer, baby bottle sterilizer, steam sterilizer, UV sterilizer, electric sterilizer, microwave sterilizer',
    'sr', 'sterilizator, sterilizator za flašice, parni sterilizator, UV sterilizator, električni sterilizator, sterilizator za mikrotalasnu',
    'ru', 'стерилизатор, стерилизатор бутылочек, паровой стерилизатор, UV стерилизатор, электрический стерилизатор, стерилизатор для микроволновки'
) WHERE slug = 'sterilizatori';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'baby jumper, jumperoo, activity center, bouncer, entertainment, jumping toy, developmental toy',
    'sr', 'bebi jumperoo, skakavac, centar aktivnosti, ljuljaška, zabava, skakačka igračka, razvojna igračka',
    'ru', 'детский джамперу, прыгунки, центр активности, качалка, развлечение, прыгающая игрушка, развивающая игрушка'
) WHERE slug = 'bebi-jumperoo';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'baby bath, foldable bath, collapsible tub, travel bath, baby bathing, portable bath, compact bath',
    'sr', 'kada za bebe, sklopiva kada, kada za putovanje, kupanje bebe, prenosiva kada, kompaktna kada',
    'ru', 'детская ванночка, складная ванночка, дорожная ванночка, купание младенца, портативная ванночка, компактная ванночка'
) WHERE slug = 'kada-bebe-sklop ljiva';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'night light, baby nightlight, nursery lamp, bedside lamp, soft light, sleep light, calming light',
    'sr', 'noćnik, noćno svetlo za bebe, lampa za dečiju sobu, lampa za krevet, blago svetlo, svetlo za spavanje, umirujuće svetlo',
    'ru', 'ночник, детский ночник, лампа для детской, прикроватная лампа, мягкий свет, свет для сна, успокаивающий свет'
) WHERE slug = 'nocnici';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'baby carrier, ergonomic carrier, front carrier, back carrier, sling, baby wrap, comfortable carrying',
    'sr', 'nosiljka za bebe, ergonomska nosiljka, prednja nosiljka, zadnja nosiljka, sling, marma za bebe, udobno nošenje',
    'ru', 'переноска для ребенка, эргономичная переноска, передняя переноска, задняя переноска, слинг, перевязь, удобное ношение'
) WHERE slug = 'nosiljke-ergonomske';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'safety gate, baby gate, stair gate, pressure mounted gate, hardware mounted, pet gate, child safety',
    'sr', 'zaštitne ograde, ograde za bebe, ograde za stepenice, pritisna ograda, hardverska montaža, ograda za kućne ljubimce, bezbednost dece',
    'ru', 'защитные ворота, ворота для детей, ворота для лестниц, раздвижные ворота, крепление на стену, ворота для животных, детская безопасность'
) WHERE slug = 'zaštitne-ograde';

-- ============================================================================
-- ZA BEBE I DECU - Nameštaj za bebe
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'baby furniture, nursery furniture, crib, changing table, baby dresser, nursery decor',
    'sr', 'nameštaj za bebe, nameštaj za dečiju sobu, krevetac, sto za presvlačenje, komoda za bebe, dekoracija dečije sobe',
    'ru', 'детская мебель, мебель для детской, кроватка, пеленальный стол, детский комод, декор детской'
) WHERE slug = 'namestaj-za-bebe';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kids furniture, children furniture, kids room, kids bed, kids desk, kids chair, playroom furniture',
    'sr', 'dečiji nameštaj, nameštaj za decu, dečija soba, dečiji krevet, dečiji sto, dečija stolica, nameštaj za igraonicu',
    'ru', 'детская мебель, мебель для детей, детская комната, детская кровать, детский стол, детский стул, мебель для игровой'
) WHERE slug = 'deciji-namestaj';

-- ============================================================================
-- ZA BEBE I DECU - Nega i higijena
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'baby care, baby hygiene, diapers, baby wipes, diaper cream, baby lotion, gentle care',
    'sr', 'nega beba, higijena beba, pelene, vlažne maramice, krema za pelene, losion za bebe, nežna nega',
    'ru', 'уход за младенцем, гигиена младенца, подгузники, влажные салфетки, крем под подгузник, лосьон для младенца, мягкий уход'
) WHERE slug = 'nega-i-higijena-beba';

-- ============================================================================
-- ZA BEBE I DECU - Hrana za bebe
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'baby food, infant nutrition, baby formula, purees, organic baby food, baby snacks, feeding',
    'sr', 'hrana za bebe, ishrana za bebe, formula za bebe, pirei, organska hrana za bebe, užina za bebe, hranjenje',
    'ru', 'детское питание, питание для младенцев, детская смесь, пюре, органическое детское питание, закуски для младенцев, кормление'
) WHERE slug = 'hrana-za-bebe';

-- ============================================================================
-- ZA BEBE I DECU - Igračke za bebe
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'baby toys, infant toys, teething toys, rattles, soft toys, sensory toys, developmental toys',
    'sr', 'igračke za bebe, igračke za novorođenče, igračke za zubiće, zvečke, meke igračke, senzorne igračke, razvojne igračke',
    'ru', 'игрушки для младенцев, игрушки для новорожденных, прорезыватели, погремушки, мягкие игрушки, сенсорные игрушки, развивающие игрушки'
) WHERE slug = 'igracke-za-bebe';

-- ============================================================================
-- ZA BEBE I DECU - Igračke za decu
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kids toys, children toys, educational toys, board games, puzzles, dolls, action figures, outdoor toys',
    'sr', 'dečije igračke, igračke za decu, edukativne igračke, društvene igre, slagalice, lutke, akcione figure, igračke za napolju',
    'ru', 'детские игрушки, игрушки для детей, образовательные игрушки, настольные игры, пазлы, куклы, фигурки, уличные игрушки'
) WHERE slug = 'igracke-za-decu';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'play mat, educational mat, activity mat, learning mat, kids carpet, interactive mat',
    'sr', 'edukativni tepih, tepih za igru, aktivnostni tepih, tepih za učenje, dečiji tepih, interaktivni tepih',
    'ru', 'развивающий коврик, игровой коврик, активный коврик, обучающий коврик, детский ковер, интерактивный коврик'
) WHERE slug = 'edukativni-tepisi';

-- ============================================================================
-- ZA BEBE I DECU - Elektronika za decu
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kids electronics, educational electronics, kids tablet, learning device, electronic toys, smart toys',
    'sr', 'elektronika za decu, edukativna elektronika, dečiji tablet, uređaj za učenje, elektronske igračke, pametne igračke',
    'ru', 'детская электроника, образовательная электроника, детский планшет, обучающее устройство, электронные игрушки, умные игрушки'
) WHERE slug = 'elektronika-za-decu';

-- ============================================================================
-- ZA BEBE I DECU - Školski pribor
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'school supplies, stationery, backpack, pencil case, notebooks, pens, crayons, art supplies',
    'sr', 'školski pribor, kancelarija, ranac, pernica, sveske, olovke, bojice, pribor za crtanje',
    'ru', 'школьные принадлежности, канцелярия, рюкзак, пенал, тетради, ручки, карандаши, принадлежности для рисования'
) WHERE slug = 'skolski-pribor';

-- ============================================================================
-- ZA BEBE I DECU - Dečija odeća
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kids clothing, children clothes, baby clothes, toddler clothing, kids fashion, infant wear',
    'sr', 'dečija odeća, odeća za decu, odeća za bebe, odeća za malu decu, dečija moda, odeća za novorođenče',
    'ru', 'детская одежда, одежда для детей, одежда для младенцев, одежда для малышей, детская мода, одежда для новорожденных'
) WHERE slug = 'decija-odeca-bebe';

-- ============================================================================
-- ZA BEBE I DECU - Dečija obuća
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'kids shoes, children footwear, baby shoes, toddler shoes, sneakers, sandals, boots',
    'sr', 'dečija obuća, obuća za decu, cipele za bebe, cipele za malu decu, patike, sandale, čizme',
    'ru', 'детская обувь, обувь для детей, обувь для младенцев, обувь для малышей, кроссовки, сандалии, ботинки'
) WHERE slug = 'decija-obuca-bebe';

-- ============================================================================
-- NAKIT I SATOVI - Zlatni nakit
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gold jewelry, gold necklace, gold bracelet, gold ring, gold earrings, 14k gold, 18k gold, yellow gold',
    'sr', 'zlatni nakit, zlatna ogrlica, zlatna narukvica, zlatni prsten, zlatne minđuše, zlato 14k, zlato 18k, žuto zlato',
    'ru', 'золотые украшения, золотое ожерелье, золотой браслет, золотое кольцо, золотые серьги, золото 14k, золото 18k, желтое золото'
) WHERE slug = 'zlatni-nakit';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gold necklace, gold chain, pendant necklace, statement necklace, gold choker, luxury necklace',
    'sr', 'zlatna ogrlica, zlatni lanac, ogrlica sa priveskom, statement ogrlica, zlatna čoker ogrlica, luksuzna ogrlica',
    'ru', 'золотое ожерелье, золотая цепочка, ожерелье с подвеской, массивное ожерелье, золотой чокер, роскошное ожерелье'
) WHERE slug = 'zlatne-ogrlice';

-- ============================================================================
-- NAKIT I SATOVI - Srebrni nakit
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'silver jewelry, sterling silver, silver necklace, silver bracelet, silver ring, silver earrings, 925 silver',
    'sr', 'srebrni nakit, sterling srebro, srebrna ogrlica, srebrna narukvica, srebrni prsten, srebrne minđuše, srebro 925',
    'ru', 'серебряные украшения, стерлинговое серебро, серебряное ожерелье, серебряный браслет, серебряное кольцо, серебряные серьги, серебро 925'
) WHERE slug = 'srebrni-nakit';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'silver earrings, sterling silver earrings, drop earrings, hoop earrings, stud earrings, dangle earrings',
    'sr', 'srebrne minđuše, sterling srebrne minđuše, viseće minđuše, ringla minđuše, naušnice, minđuše sa viskom',
    'ru', 'серебряные серьги, серьги из стерлингового серебра, висячие серьги, серьги кольца, серьги гвоздики, серьги с подвесками'
) WHERE slug = 'srebrne-minđuše';

-- ============================================================================
-- NAKIT I SATOVI - Dijamanti
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'diamonds, diamond jewelry, diamond ring, diamond necklace, diamond earrings, certified diamonds, brilliant cut',
    'sr', 'dijamanti, dijamantski nakit, dijamantski prsten, dijamantska ogrlica, dijamantske minđuše, sertifikovani dijamanti, brilijantni rez',
    'ru', 'бриллианты, бриллиантовые украшения, бриллиантовое кольцо, бриллиантовое ожерелье, бриллиантовые серьги, сертифицированные бриллианты, бриллиантовая огранка'
) WHERE slug = 'dijamanti';

-- ============================================================================
-- NAKIT I SATOVI - Vereničko prstenje
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'engagement ring, proposal ring, diamond engagement ring, solitaire ring, wedding ring, bridal jewelry',
    'sr', 'vereničko prstenje, prsten za prosidbu, dijamantski vereničko prstenje, soliter prsten, venčani prsten, nakit za venčanje',
    'ru', 'помолвочное кольцо, кольцо для предложения, бриллиантовое помолвочное кольцо, солитер, обручальное кольцо, свадебные украшения'
) WHERE slug = 'verenicko-prstenje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'diamond engagement ring, proposal ring with diamond, bridal ring, certified diamond ring, luxury engagement',
    'sr', 'prsten za veridbu sa dijamantom, vereničko prstenje sa dijamantom, prsten za venčanje, sertifikovani dijamant, luksuzna veridba',
    'ru', 'помолвочное кольцо с бриллиантом, кольцо для предложения с бриллиантом, свадебное кольцо, сертифицированное кольцо с бриллиантом, роскошная помолвка'
) WHERE slug = 'prstenje-veridba-dijamant';

-- ============================================================================
-- NAKIT I SATOVI - Biserni nakit
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'pearl jewelry, pearl necklace, pearl earrings, pearl bracelet, freshwater pearls, cultured pearls, Akoya pearls',
    'sr', 'biserni nakit, biserna ogrlica, biserne minđuše, biserna narukvica, slatkovodni biseri, gajeni biseri, Akoya biseri',
    'ru', 'жемчужные украшения, жемчужное ожерелье, жемчужные серьги, жемчужный браслет, пресноводный жемчуг, культивированный жемчуг, жемчуг Akoya'
) WHERE slug = 'biserni-nakit';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'pearl beads, pearl bracelet, beaded bracelet, pearl bead necklace, gemstone beads',
    'sr', 'biserne perlice, biserna narukvica, narukvica od perlica, ogrlica od bisera, perlice od dragog kamenja',
    'ru', 'жемчужные бусины, жемчужный браслет, браслет из бусин, ожерелье из жемчужных бусин, бусины из драгоценных камней'
) WHERE slug = 'perlice-narukvice';

-- ============================================================================
-- NAKIT I SATOVI - Moderni nakit
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'modern jewelry, contemporary jewelry, fashion jewelry, trendy jewelry, designer jewelry, statement pieces',
    'sr', 'moderni nakit, savremeni nakit, modni nakit, trendovski nakit, dizajnerski nakit, statement komadi',
    'ru', 'современные украшения, модные украшения, бижутерия, трендовые украшения, дизайнерские украшения, массивные украшения'
) WHERE slug = 'moderni-nakit';

-- ============================================================================
-- NAKIT I SATOVI - Vintage broševi
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'vintage brooch, antique brooch, retro brooch, collectible brooch, estate jewelry, classic brooch',
    'sr', 'vintage broš, antikni broš, retro broš, kolekcionarski broš, starinski nakit, klasični broš',
    'ru', 'винтажная брошь, антикварная брошь, ретро брошь, коллекционная брошь, старинные украшения, классическая брошь'
) WHERE slug = 'broše-vintage';

-- ============================================================================
-- NAKIT I SATOVI - Privesci i lanci
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'pendants, chains, pendant necklace, charm pendant, gold chain, silver chain, necklace chain',
    'sr', 'privesci, lanci, ogrlica sa priveskom, privezak, zlatni lanac, srebrni lanac, lanac za ogrlicu',
    'ru', 'подвески, цепочки, ожерелье с подвеской, шарм подвеска, золотая цепочка, серебряная цепочка, цепочка для ожерелья'
) WHERE slug = 'privesci-lanci';

-- ============================================================================
-- NAKIT I SATOVI - Muški satovi
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'men watches, mens wristwatch, luxury watch, sports watch, automatic watch, chronograph, dive watch',
    'sr', 'muški satovi, ručni sat za muškarce, luksuzni sat, sportski sat, automatski sat, hronograf, ronilački sat',
    'ru', 'мужские часы, наручные часы для мужчин, роскошные часы, спортивные часы, автоматические часы, хронограф, дайверские часы'
) WHERE slug = 'muski-satovi';

-- ============================================================================
-- NAKIT I SATOVI - Ženski satovi
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'women watches, ladies watch, luxury ladies watch, fashion watch, elegant watch, bracelet watch, dress watch',
    'sr', 'ženski satovi, ručni sat za žene, luksuzni ženski sat, modni sat, elegantan sat, sat narukvica, svečani sat',
    'ru', 'женские часы, наручные часы для женщин, роскошные женские часы, модные часы, элегантные часы, часы-браслет, вечерние часы'
) WHERE slug = 'zenski-satovi';

-- ============================================================================
-- NAKIT I SATOVI - Pametni satovi fitness
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'smartwatch, fitness tracker, smart watch, activity tracker, heart rate monitor, GPS watch, sports watch',
    'sr', 'pametni sat, fitness tracker, sat sa meračem aktivnosti, praćenje aktivnosti, merač pulsa, GPS sat, sportski sat',
    'ru', 'умные часы, фитнес трекер, смарт часы, трекер активности, пульсометр, GPS часы, спортивные часы'
) WHERE slug = 'pametni-satovi-fitness';

-- ============================================================================
-- NAKIT I SATOVI - Luksuzni satovi
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'luxury watches, prestige watches, high-end watches, premium watches, Swiss watches, exclusive timepieces',
    'sr', 'luksuzni satovi, prestižni satovi, vrhunski satovi, premium satovi, švajcarski satovi, ekskluzivni satovi',
    'ru', 'роскошные часы, престижные часы, высококлассные часы, премиум часы, швейцарские часы, эксклюзивные часы'
) WHERE slug = 'luksuzni-satovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'Swiss watches, Swiss luxury watches, Swiss made, Rolex, Omega, Tag Heuer, Breitling, precision timepieces',
    'sr', 'švajcarski satovi, švajcarski luksuzni satovi, napravljeno u Švajcarskoj, Rolex, Omega, Tag Heuer, Breitling, precizni satovi',
    'ru', 'швейцарские часы, швейцарские роскошные часы, сделано в Швейцарии, Rolex, Omega, Tag Heuer, Breitling, точные часы'
) WHERE slug = 'luksuzni-satovi-swiss';

-- ============================================================================
-- NAKIT I SATOVI - Satovski dodaci
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'watch accessories, watch strap, watch band, watch winder, watch box, watch tools, replacement band',
    'sr', 'satovski dodaci, kaiš za sat, narukvica za sat, navijač satova, kutija za satove, alati za satove, rezervni kaiš',
    'ru', 'аксессуары для часов, ремешок для часов, браслет для часов, шкатулка для часов, коробка для часов, инструменты для часов, сменный ремешок'
) WHERE slug = 'satovski-dodaci';

COMMIT;
