-- Migration: Fill meta_keywords for remaining L2/L3 categories (Knjige, Hrana, Ljubimci, Industrija, Kancelarija, Muzika, Umetnost, Usluge, Ostalo)
-- Generated: 2025-12-22

-- ============================================================================
-- KNJIGE I MEDIJI (Books and Media)
-- ============================================================================

-- Knjige - Beletristika
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'fiction, novels, literature, contemporary fiction, bestsellers, prose, modern literature',
    'sr', 'beletristika, romani, književnost, savremena beletristika, bestseleri, proza, moderna književnost',
    'ru', 'беллетристика, романы, литература, современная проза, бестселлеры, художественная литература'
) WHERE slug = 'knjige-beletristika';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'children''s books, picture books, bedtime stories, educational books, kids literature, fairy tales',
    'sr', 'dečije knjige, slikovnice, priče za laku noć, edukativne knjige, dečija književnost, bajke',
    'ru', 'детские книги, книжки с картинками, сказки на ночь, развивающие книги, детская литература'
) WHERE slug = 'decije-knjige';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'biography, life story, memoirs, famous people, historical figures, celebrity biographies',
    'sr', 'biografija, životna priča, memoare, poznate ličnosti, istorijske figure, biografije slavnih',
    'ru', 'биография, история жизни, мемуары, знаменитости, исторические личности, биографии звезд'
) WHERE slug = 'biografije';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'autobiography, personal memoir, self biography, life experience, personal story, confessions',
    'sr', 'autobiografija, lični memoari, sopstvena biografija, životno iskustvo, lična priča, ispovesti',
    'ru', 'автобиография, личные воспоминания, история своей жизни, личный опыт, исповеди'
) WHERE slug = 'autobiografije';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'poetry, poems, verse, contemporary poetry, classic poetry, poets, rhymes',
    'sr', 'poezija, pesme, stihovi, savremena poezija, klasična poezija, pesnici, rime',
    'ru', 'поэзия, стихи, стихотворения, современная поэзия, классическая поэзия, поэты, рифмы'
) WHERE slug = 'poezija';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'essays, articles, non-fiction prose, literary essays, critical essays, personal essays',
    'sr', 'eseji, članci, nefikcionalna proza, književni eseji, kritički eseji, lični eseji',
    'ru', 'эссе, статьи, публицистика, литературные эссе, критические эссе, личные размышления'
) WHERE slug = 'eseji';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'detective, crime fiction, mystery, thriller, whodunit, noir, police procedural',
    'sr', 'detektivski romani, kriminalistički romani, misterija, triler, noir, policijski romani',
    'ru', 'детективы, криминальные романы, мистика, триллеры, детективные истории, нуар'
) WHERE slug = 'kriminalisticki-romani';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'romance novels, love stories, romantic fiction, contemporary romance, historical romance',
    'sr', 'ljubavni romani, ljubavne priče, romantična fikcija, savremeni ljubavni romani, istorijski romani',
    'ru', 'любовные романы, романтика, любовные истории, современные романы о любви, исторические романы'
) WHERE slug = 'ljubavni-romani';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'fantasy, epic fantasy, urban fantasy, magical realism, dragons, wizards, sword and sorcery',
    'sr', 'fantastika, epska fantastika, urbana fantastika, magični realizam, zmajevi, čarobnjaci',
    'ru', 'фэнтези, эпическое фэнтези, городское фэнтези, магический реализм, драконы, волшебники'
) WHERE slug = 'fantastika';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'science fiction, sci-fi, space opera, cyberpunk, time travel, dystopia, aliens, robots',
    'sr', 'naučna fantastika, SF, svemirska opera, cyberpunk, putovanje kroz vreme, distopija, vanzemaljci',
    'ru', 'научная фантастика, фантастика, космическая опера, киберпанк, путешествия во времени, дистопия'
) WHERE slug = 'sf-knjige';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'horror, supernatural horror, psychological horror, ghost stories, vampires, zombies, Stephen King',
    'sr', 'horor, natprirodni horor, psihološki horor, priče o duhovima, vampiri, zombiji, stravične priče',
    'ru', 'ужасы, мистический хоррор, психологический хоррор, истории о призраках, вампиры, зомби'
) WHERE slug = 'horor';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'history books, historical non-fiction, world history, war history, ancient history, medieval',
    'sr', 'istorijske knjige, istorija, svetska istorija, ratna istorija, antička istorija, srednji vek',
    'ru', 'исторические книги, история, мировая история, военная история, древняя история, средневековье'
) WHERE slug = 'istorijske-knjige';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'professional books, textbooks, academic books, study guides, technical books, manuals',
    'sr', 'stručne knjige, udžbenici, akademske knjige, priručnici, tehničke knjige, vodiči za učenje',
    'ru', 'профессиональные книги, учебники, научные книги, руководства, технические книги, пособия'
) WHERE slug = 'strucne-knjige';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'cookbooks, recipes, cooking, baking, cuisine, culinary books, chef recipes, meal planning',
    'sr', 'kuvari, recepti, kuvanje, pečenje, kuhinja, kulinarske knjige, šefovski recepti, planiranje obroka',
    'ru', 'кулинарные книги, рецепты, готовка, выпечка, кухня, книги о еде, рецепты шефов'
) WHERE slug = 'kuvari-recepti';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'parenting books, child development, raising kids, parenting tips, family advice, baby care',
    'sr', 'knjige za roditelje, razvoj dece, vaspitanje dece, saveti za roditelje, porodični saveti, nega beba',
    'ru', 'книги для родителей, развитие детей, воспитание, советы родителям, уход за детьми'
) WHERE slug = 'knjige-za-roditelje';

-- Media
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'audiobooks, audio literature, narrated books, spoken word, book recordings, audible',
    'sr', 'audio knjige, audio književnost, pričane knjige, snimljene knjige, knjige za slušanje',
    'ru', 'аудиокниги, аудио литература, озвученные книги, записи книг, книги для прослушивания'
) WHERE slug = 'audio-knjige';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ebooks, digital books, kindle, epub, pdf books, electronic books, online reading',
    'sr', 'e-knjige, digitalne knjige, kindle, epub, pdf knjige, elektronske knjige, čitanje online',
    'ru', 'электронные книги, цифровые книги, kindle, epub, pdf книги, онлайн чтение'
) WHERE slug = 'e-knjige';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'magazines, periodicals, monthly magazines, weekly magazines, fashion magazines, news magazines',
    'sr', 'časopisi, periodika, mesečni časopisi, nedeljni časopisi, modni časopisi, novinski časopisi',
    'ru', 'журналы, периодика, ежемесячные журналы, еженедельные журналы, модные журналы'
) WHERE slug = 'casopisi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'comics, graphic novels, manga, superheroes, Marvel, DC Comics, comic books, illustrated stories',
    'sr', 'stripovi, grafički romani, manga, superheroji, Marvel, DC Comics, strip albumi, ilustrovane priče',
    'ru', 'комиксы, графические романы, манга, супергерои, Marvel, DC Comics, иллюстрированные истории'
) WHERE slug = 'stripovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'movies, DVD, Blu-ray, cinema, films, classic movies, new releases, movie collection',
    'sr', 'filmovi, DVD, Blu-ray, bioskop, film kolekcija, klasični filmovi, nove izdaje, filmoteka',
    'ru', 'фильмы, DVD, Blu-ray, кино, коллекция фильмов, классические фильмы, новинки кино'
) WHERE slug = 'filmovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'music, CDs, vinyl, albums, music collection, genres, artists, soundtrack, recordings',
    'sr', 'muzika, CD-ovi, vinil, albumi, muzička kolekcija, žanrovi, umetnici, soundtrack, snimci',
    'ru', 'музыка, диски, винил, альбомы, музыкальная коллекция, жанры, артисты, саундтреки'
) WHERE slug = 'muzika';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'rare books, first editions, antique books, collectible books, signed copies, limited editions',
    'sr', 'retke knjige, prva izdanja, antikne knjige, kolekcionarske knjige, potpisane kopije, limitirana izdanja',
    'ru', 'редкие книги, первые издания, антикварные книги, коллекционные книги, подписанные экземпляры'
) WHERE slug = 'retke-knjige';

-- ============================================================================
-- HRANA I PICE (Food and Beverages)
-- ============================================================================

-- Organska hrana
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'organic fruits, fresh organic produce, seasonal fruits, eco fruits, bio fruits, natural',
    'sr', 'organsko voće, sveže organsko voće, sezonsko voće, eko voće, bio voće, prirodno',
    'ru', 'органические фрукты, свежие фрукты, сезонные фрукты, био фрукты, натуральные'
) WHERE slug = 'organsko-voce';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'organic vegetables, fresh organic veggies, eco vegetables, bio vegetables, farm vegetables',
    'sr', 'organsko povrće, sveže organsko povrće, eko povrće, bio povrće, povrće sa farme',
    'ru', 'органические овощи, свежие овощи, био овощи, овощи с фермы, натуральные овощи'
) WHERE slug = 'organsko-povrce';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'organic dairy, milk, cheese, yogurt, butter, organic dairy products, lactose-free',
    'sr', 'organski mlečni proizvodi, mleko, sir, jogurt, maslac, bio mlečni proizvodi, bez laktoze',
    'ru', 'органические молочные продукты, молоко, сыр, йогурт, масло, био молочка'
) WHERE slug = 'mlecni-proizvodi-organik';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'organic grains, whole grains, quinoa, oats, rice, buckwheat, organic cereals, wheat',
    'sr', 'organske žitarice, cela zrna, kvinoja, ovas, pirinač, heljda, organske žitarice, pšenica',
    'ru', 'органические злаки, цельное зерно, киноа, овес, рис, гречка, органические крупы'
) WHERE slug = 'zitarice-organik';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'organic honey, raw honey, natural honey, bee products, propolis, bee pollen, manuka honey',
    'sr', 'organski med, sirov med, prirodni med, pčelinji proizvodi, propolis, cvetni prah, manuka med',
    'ru', 'органический мёд, сырой мёд, натуральный мёд, продукты пчеловодства, прополис, манука'
) WHERE slug = 'med-organski';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'organic juices, fresh pressed juice, cold pressed, fruit juice, vegetable juice, detox juice',
    'sr', 'organski sokovi, sveže ceđeni sok, cold pressed, voćni sok, povrtni sok, detox sok',
    'ru', 'органические соки, свежевыжатые соки, холодный отжим, фруктовые соки, детокс соки'
) WHERE slug = 'organski-sokovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'premium olive oil, extra virgin olive oil, cold pressed, Greek olive oil, Italian olive oil',
    'sr', 'premium maslinovo ulje, ekstra devičansko maslinovo ulje, hladno ceđeno, grčko, italijansko',
    'ru', 'премиум оливковое масло, экстра вирджин, холодный отжим, греческое, итальянское'
) WHERE slug = 'maslinovo-ulje-premium';

-- Dijetalna hrana
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'vegan food, plant based, vegan products, dairy free, egg free, animal free, vegan meals',
    'sr', 'veganska hrana, biljni proizvodi, veganski proizvodi, bez mleka, bez jaja, biljni obroci',
    'ru', 'веганская еда, растительные продукты, без молока, без яиц, веганские блюда'
) WHERE slug = 'veganska-hrana';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'vegetarian food, meatless, vegetarian meals, plant protein, veggie products, no meat',
    'sr', 'vegetarijanska hrana, bez mesa, vegetarijanski obroci, biljni protein, proizvodi bez mesa',
    'ru', 'вегетарианская еда, без мяса, вегетарианские блюда, растительный белок'
) WHERE slug = 'vegetarijanska-hrana';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gluten free, celiac safe, wheat free, gluten free products, no gluten, celiac diet',
    'sr', 'bezglutenska hrana, bez glutena, celiac safe, bez pšenice, proizvodi bez glutena',
    'ru', 'безглютеновая еда, без глютена, целиакия, без пшеницы, продукты без глютена'
) WHERE slug = 'bezglutenska-hrana';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'lactose free, dairy free, no lactose, lactose intolerance, milk free, plant milk',
    'sr', 'bez laktoze, bez mlečnih proizvoda, intolerancija na laktozu, biljno mleko',
    'ru', 'без лактозы, без молока, непереносимость лактозы, растительное молоко'
) WHERE slug = 'bez-laktoze';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'sugar free, no sugar, sugar free products, diabetic friendly, stevia, sugar alternatives',
    'sr', 'bez šećera, proizvodi bez šećera, za dijabetičare, stevija, zamena za šećer',
    'ru', 'без сахара, продукты без сахара, для диабетиков, стевия, заменители сахара'
) WHERE slug = 'hrana-bez-secera';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'keto diet, low carb, ketogenic, keto products, high fat low carb, ketosis, keto meals',
    'sr', 'keto dijeta, nisko ugljenih hidrata, ketogena ishrana, keto proizvodi, visoko masti malo UH',
    'ru', 'кето диета, низкоуглеводная, кетогенная диета, кето продукты, кетоз, кето еда'
) WHERE slug = 'keto-dijeta';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'paleo diet, caveman diet, primal diet, paleo products, grain free, paleo meals',
    'sr', 'paleo dijeta, dijeta pećinskog čoveka, paleo proizvodi, bez žitarica, paleo obroci',
    'ru', 'палео диета, диета пещерного человека, палео продукты, без зерновых, палео еда'
) WHERE slug = 'paleo-dijeta';

-- Premium piće
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'premium coffee, specialty coffee, arabica, robusta, single origin, coffee beans, espresso',
    'sr', 'premium kafa, specijalna kafa, arabica, robusta, single origin, kafa u zrnu, espresso',
    'ru', 'премиум кофе, специальный кофе, арабика, робуста, моносорт, кофе в зёрнах, эспрессо'
) WHERE slug = 'kafa-premium';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'premium tea, loose leaf tea, green tea, black tea, herbal tea, oolong, matcha, tea collection',
    'sr', 'premium čaj, čaj u listovima, zeleni čaj, crni čaj, biljni čaj, oolong, matcha, čajna kolekcija',
    'ru', 'премиум чай, листовой чай, зелёный чай, чёрный чай, травяной чай, улун, матча'
) WHERE slug = 'caj-premium';

-- Sportska ishrana
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'protein powder, whey protein, casein, plant protein, protein supplements, isolate, muscle gain',
    'sr', 'proteini u prahu, whey protein, kazein, biljni protein, suplementi proteina, izolat, masa',
    'ru', 'протеин в порошке, сывороточный протеин, казеин, растительный протеин, изолят'
) WHERE slug = 'proteini-u-prahu';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'protein bars, energy bars, protein snacks, meal replacement bars, fitness bars, low sugar',
    'sr', 'protein pločice, energetske pločice, protein grickalice, zamena za obrok, fitness pločice',
    'ru', 'протеиновые батончики, энергетические батончики, фитнес батончики, замена еды'
) WHERE slug = 'protein-pločice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'sports supplements, pre-workout, creatine, BCAA, amino acids, fat burners, mass gainers',
    'sr', 'sportski dodaci ishrani, pre-workout, kreatin, BCAA, aminokiseline, fat burner, gaineri',
    'ru', 'спортивные добавки, предтреники, креатин, BCAA, аминокислоты, жиросжигатели'
) WHERE slug = 'sportska-ishrana-dodaci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'superfood, superfoods, chia seeds, spirulina, chlorella, acai, goji berries, hemp seeds',
    'sr', 'superfood, super hrana, chia seme, spirulina, chlorella, acai, goji bobice, konoplja',
    'ru', 'суперфуды, чиа, спирулина, хлорелла, асаи, годжи, семена конопли'
) WHERE slug = 'superfood';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'healthy snacks, fitness snacks, low calorie snacks, nuts, dried fruits, granola, energy balls',
    'sr', 'zdrave grickalice, fitness grickalice, niskokalorijske grickalice, oraščići, suvo voće, granola',
    'ru', 'здоровые снеки, фитнес снеки, низкокалорийные снеки, орехи, сухофрукты, гранола'
) WHERE slug = 'zdrave-grickalice';

-- Testenine
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'whole grain pasta, wholemeal pasta, brown pasta, fiber pasta, healthy pasta, spelt pasta',
    'sr', 'integralne testenine, testenine od celih zrna, smeđe testenine, vlakna, zdrave testenine',
    'ru', 'цельнозерновые макароны, паста из цельного зерна, коричневая паста, клетчатка'
) WHERE slug = 'testenine-integralne';

-- ============================================================================
-- KUCNI LJUBIMCI (Pets)
-- ============================================================================

-- Hrana za ljubimce
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'dog food, dry dog food, wet dog food, puppy food, senior dog food, grain free, premium brands',
    'sr', 'hrana za pse, suva hrana za pse, vlažna hrana, hrana za štence, za starije pse, bez žitarica',
    'ru', 'корм для собак, сухой корм, влажный корм, для щенков, для пожилых собак, беззерновой'
) WHERE slug = 'hrana-za-pse';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'cat food, dry cat food, wet cat food, kitten food, senior cat food, grain free, premium brands',
    'sr', 'hrana za mačke, suva hrana za mačke, vlažna hrana, hrana za mačiće, za starije mačke',
    'ru', 'корм для кошек, сухой корм, влажный корм, для котят, для пожилых кошек, премиум'
) WHERE slug = 'hrana-za-macke';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bird food, bird seed, parrot food, canary food, finch food, bird treats, millet',
    'sr', 'hrana za ptice, seme za ptice, hrana za papagaje, za kanarince, proso, poslastice za ptice',
    'ru', 'корм для птиц, семена для птиц, корм для попугаев, для канареек, просо, лакомства'
) WHERE slug = 'hrana-za-ptice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'rodent food, hamster food, guinea pig food, rabbit food, chinchilla food, hay, pellets',
    'sr', 'hrana za glodare, hrana za hrčke, zamorče, zečeve, činčile, seno, peleti',
    'ru', 'корм для грызунов, корм для хомяков, морских свинок, кроликов, шиншилл, сено'
) WHERE slug = 'hrana-za-glodare';

-- Oprema za ljubimce
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'dog toys, cat toys, chew toys, interactive toys, squeaky toys, ball toys, plush toys',
    'sr', 'igračke za pse, igračke za žvakanje, interaktivne igračke, lopte, plišane igračke',
    'ru', 'игрушки для собак, жевательные игрушки, интерактивные игрушки, мячики, плюшевые'
) WHERE slug = 'igracke-za-pse';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'cat toys, catnip toys, feather toys, laser pointer, scratching posts, ball toys, interactive',
    'sr', 'igračke za mačke, mačja meta, perje, laser pokazivač, grebalice, lopte, interaktivne',
    'ru', 'игрушки для кошек, с кошачьей мятой, с перьями, лазер, когтеточки, интерактивные'
) WHERE slug = 'igracke-za-macke';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'leash, collar, harness, dog walking equipment, retractable leash, training leash, puppy leash',
    'sr', 'povodac, ogrlica, am, oprema za šetnju, automatski povodac, dresura povodac, za štence',
    'ru', 'поводок, ошейник, шлейка, снаряжение для выгула, рулетка, тренировочный поводок'
) WHERE slug = 'oprema-za-setnju';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'pet beds, dog beds, cat beds, orthopedic beds, heated beds, cushions, sleeping mats',
    'sr', 'kreveti za ljubimce, kreveti za pse, mačke, ortopedski kreveti, grejani kreveti, prostirke',
    'ru', 'лежанки для питомцев, лежанки для собак, кошек, ортопедические, с подогревом, коврики'
) WHERE slug = 'kreveti-za-ljubimce';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'litter box, cat toilet, litter tray, self-cleaning litter box, litter scoop, odor control',
    'sr', 'toaleta za mačke, pesak za mačke, posuda za pesak, samočisteća, lopatica, kontrola mirisa',
    'ru', 'лоток для кошек, туалет для кошек, наполнитель, самоочищающийся, совок, контроль запаха'
) WHERE slug = 'toaleta-za-ljubimce';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'pet carrier, travel carrier, pet backpack, airline approved, soft carrier, hard carrier',
    'sr', 'nosiljke za ljubimce, transportne torbe, rančevi za ljubimce, odobreno od aviokompanija',
    'ru', 'переноски для животных, дорожные переноски, рюкзаки для питомцев, авиапереноски'
) WHERE slug = 'nosilje-za-ljubimce';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'aquarium, fish tank, aquarium setup, saltwater aquarium, freshwater aquarium, nano tank',
    'sr', 'akvarijumi, akvarijum oprema, morski akvarijum, slatkovodni akvarijum, nano akvarijum',
    'ru', 'аквариумы, аквариумное оборудование, морской аквариум, пресноводный, нано аквариум'
) WHERE slug = 'akvarijumi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'terrarium, reptile enclosure, lizard tank, gecko terrarium, snake terrarium, vivarium',
    'sr', 'terarijumi, ograđeni prostor za reptile, za guštera, za geka, za zmiju, vivarijum',
    'ru', 'террариумы, вольер для рептилий, для ящериц, для гекконов, для змей, виварий'
) WHERE slug = 'terarijumi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'water bowls, pet water dispenser, automatic waterer, fountain, food bowls, feeding mat',
    'sr', 'vodopasi, posude za vodu, automatski dozatori, fontana, posude za hranu, podloga za hranjenje',
    'ru', 'поилки, миски для воды, автопоилки, фонтан, миски для еды, коврики для кормления'
) WHERE slug = 'vodopasi';

-- ============================================================================
-- INDUSTRIJA I ALATI (Industry and Tools)
-- ============================================================================

-- Ručni alati električni
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'professional drills, cordless drills, hammer drills, impact drills, drill sets, Bosch, Makita',
    'sr', 'profesionalne bušilice, akumulatorske bušilice, udarne bušilice, Bosch, Makita, DeWalt',
    'ru', 'профессиональные дрели, аккумуляторные дрели, ударные дрели, Bosch, Makita, DeWalt'
) WHERE slug = 'busilice-profesionalne';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'impact drivers, cordless screwdrivers, electric screwdrivers, power screwdriver, drill driver',
    'sr', 'udarni odvijači, akumulatorski odvijači, električni odvijači, udarna bušilica odvijač',
    'ru', 'ударные шуруповерты, аккумуляторные отвертки, электроотвертки, дрель-шуруповерт'
) WHERE slug = 'odvijaci-akumulatorski';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'impact wrenches, pneumatic impact wrench, cordless impact wrench, air wrench, torque wrench',
    'sr', 'udarni odvijači, pneumatski udarni odvijač, akumulatorski, pištolj udarni, okretni moment',
    'ru', 'ударные гайковерты, пневматические гайковерты, аккумуляторные, динамометрический'
) WHERE slug = 'udarni-odvijaci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'angle grinders, disc grinders, grinding, cutting, polishing, 125mm, 230mm, Bosch, Makita',
    'sr', 'ugaone brusilice, flekseri, brusenje, sečenje, poliranje, 125mm, 230mm, Bosch, Makita',
    'ru', 'угловые шлифмашины, болгарки, шлифовка, резка, полировка, 125мм, 230мм, Bosch'
) WHERE slug = 'brusilice-ugaone';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'orbital sanders, vibration sanders, sanding machines, wood sanding, finishing, smooth surface',
    'sr', 'vibracione brusilice, orbitalne brusilice, mašine za brušenje, brušenje drveta, fina obrada',
    'ru', 'виброшлифовальные машины, орбитальные шлифмашины, шлифовка дерева, финишная обработка'
) WHERE slug = 'brusilice-vibracione';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'circular saws, circular saw blades, wood cutting, rip cuts, cross cuts, miter saw, table saw',
    'sr', 'kružne testere, kružne sečice, sečenje drveta, testera listovi, precizno sečenje',
    'ru', 'циркулярные пилы, дисковые пилы, резка дерева, пильные диски, торцовочная пила'
) WHERE slug = 'testeri-kruzne';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'chainsaws, electric chainsaw, gas chainsaw, battery chainsaw, chain, bar, cutting trees, Stihl',
    'sr', 'lančane testere, električna lančanica, benzinska, akumulatorska, lanac, šina, Stihl',
    'ru', 'цепные пилы, электропилы, бензопилы, аккумуляторные, цепь, шина, валка деревьев'
) WHERE slug = 'testeri-lancane';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'electric planers, wood planer, power planer, thickness planer, hand planer, carpentry tools',
    'sr', 'električne rende, rende za drvo, debljačica, strug, alat za stolariju, blanjanje',
    'ru', 'электрорубанки, рубанки по дереву, рейсмусовые станки, строгальные инструменты'
) WHERE slug = 'rende-elektricne';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'circular saw, table saw, miter saw, cutting machine, wood cutting machine, precision cuts',
    'sr', 'cirkular, kružna testera, kombinovana testera, mašina za sečenje drveta, precizno sečenje',
    'ru', 'циркулярка, настольная пила, торцовочная пила, машина для резки дерева, точная резка'
) WHERE slug = 'cirkular';

-- Generatori i kompresori
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gasoline generators, portable generator, backup power, emergency power, 2kW, 5kW, 10kW',
    'sr', 'benzinski generatori, prenosni generator, rezervna struja, agregat, 2kW, 5kW, 10kW',
    'ru', 'бензиновые генераторы, портативные генераторы, резервное питание, 2кВт, 5кВт, 10кВт'
) WHERE slug = 'generatori-benzin';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'diesel generators, industrial generator, heavy duty, 10kW, 20kW, 50kW, three phase generator',
    'sr', 'dizel generatori, industrijski agregat, profesionalni, 10kW, 20kW, 50kW, trofazni',
    'ru', 'дизельные генераторы, промышленные генераторы, 10кВт, 20кВт, 50кВт, трёхфазные'
) WHERE slug = 'generatori-dizel';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'piston compressors, oil compressor, belt driven, 50L, 100L, 200L, air compressor, workshop',
    'sr', 'klipni kompresori, uljni kompresor, remenski, 50L, 100L, 200L, vazdušni kompresor, radionica',
    'ru', 'поршневые компрессоры, масляные компрессоры, ременные, 50л, 100л, 200л, воздушные'
) WHERE slug = 'kompresori-klipni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'fan compressors, oil-free compressor, portable compressor, quiet compressor, mobile workshop',
    'sr', 'ventilatorski kompresori, bez ulja, prenosni kompresor, tih kompresor, mobilna radionica',
    'ru', 'вентиляторные компрессоры, безмасляные, портативные, тихие компрессоры, мобильные'
) WHERE slug = 'kompresori-ventilatorski';

-- Zavarivači
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'inverter welders, MMA welding, ARC welder, portable welder, 160A, 200A, electrode welding',
    'sr', 'inverter zavarivači, MMA zavarivanje, ručno elektrolučno, prenosni, 160A, 200A, elektroda',
    'ru', 'инверторные сварочные аппараты, ММА сварка, ручная дуговая, 160А, 200А, электродная'
) WHERE slug = 'zavarivaci-inverter';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'MIG/MAG welders, semi-automatic welder, wire welding, gas welding, CO2 welding, flux core',
    'sr', 'MIG/MAG zavarivači, poluautomatski zavarivač, žica, zavarivanje u gasu, CO2, flux core',
    'ru', 'MIG/MAG сварка, полуавтоматы, проволочная сварка, газовая сварка, CO2, флюс'
) WHERE slug = 'zavarivaci-mig';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'TIG welders, argon welding, tungsten electrode, AC/DC welder, aluminum welding, stainless steel',
    'sr', 'TIG zavarivači, argon zavarivanje, volfram elektroda, AC/DC, aluminijum, nerđajući čelik',
    'ru', 'TIG сварка, аргоновая сварка, вольфрамовый электрод, AC/DC, алюминий, нержавейка'
) WHERE slug = 'zavarivaci-tig';

-- Merila i nivelliri
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'laser level, cross line laser, rotary laser, self-leveling laser, construction laser, green beam',
    'sr', 'niveliri, laserski niveler, krstasti laser, rotacioni laser, gradjevinski laser, zeleni zrak',
    'ru', 'нивелиры, лазерные уровни, кросслайнер, ротационные, строительные лазеры, зелёный луч'
) WHERE slug = 'niveliri';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'laser measure, laser distance meter, rangefinder, 50m, 100m, area calculation, Bosch, Leica',
    'sr', 'laserski merači, laserski daljinomer, merač razdaljine, 50m, 100m, površina, Bosch, Leica',
    'ru', 'лазерные рулетки, дальномеры, измеритель расстояний, 50м, 100м, площадь, Bosch, Leica'
) WHERE slug = 'meraci-laserni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'digital calipers, vernier calipers, micrometer, measuring tools, precision measurement, 0.01mm',
    'sr', 'digitalna merila, pomoćno merilo, mikrometar, alati za merenje, precizno merenje, 0.01mm',
    'ru', 'цифровые штангенциркули, микрометры, измерительные инструменты, точность 0.01мм'
) WHERE slug = 'merila-digitalna';

-- Stezanje
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'clamps, quick clamps, bar clamps, F clamps, corner clamps, woodworking clamps, pipe clamps',
    'sr', 'stege, brze stege, šipkaste stege, F stege, ugaone stege, stolarijske stege, cevne stege',
    'ru', 'струбцины, быстрозажимные, стержневые, F-образные, угловые, столярные струбцины'
) WHERE slug = 'stege';

-- Baterije profesionalne
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'professional batteries, 18V battery, 20V battery, power tool batteries, lithium ion, Ah capacity',
    'sr', 'profesionalni akumulatori, 18V baterija, 20V, baterije za alate, litijum jonski, Ah kapacitet',
    'ru', 'профессиональные аккумуляторы, 18В батарея, 20В, литий-ионные, батареи для инструмента'
) WHERE slug = 'akumulatori-profesionalni';

-- ============================================================================
-- KANCELARIJSKI MATERIJAL (Office Supplies)
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'office supplies, stationery, office accessories, desk accessories, office essentials, business supplies',
    'sr', 'kancelarijski pribor, kancelarija, pribor za pisanje, oprema za stoni, poslovno',
    'ru', 'канцелярия, офисные принадлежности, настольные аксессуары, офисные товары'
) WHERE slug = 'kancelarijski-pribor';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'pens, ballpoint pen, gel pens, markers, highlighters, permanent markers, fine tip, writing tools',
    'sr', 'olovke, hemijske olovke, gel olovke, markeri, highlighteri, trajni markeri, tanki vrh',
    'ru', 'ручки, шариковые ручки, гелевые ручки, маркеры, хайлайтеры, перманентные маркеры'
) WHERE slug = 'olovke-i-markeri';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'notebooks, notepads, spiral notebooks, lined notebooks, dot grid, bullet journal, A4, A5',
    'sr', 'beležnice, blokovi, spiralni blokovi, linirani blokovi, tačkasti, bullet journal, A4, A5',
    'ru', 'блокноты, записные книжки, спиральные блокноты, линованные, точечная сетка, А4, А5'
) WHERE slug = 'beležnice-blokovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'office paper, copy paper, printer paper, A4 paper, A3 paper, white paper, colored paper, 80gsm',
    'sr', 'kancelarijska hartija, papir za kopiranje, papir za štampač, A4 papir, A3, beli, 80gsm',
    'ru', 'офисная бумага, бумага для копирования, для принтера, А4, А3, белая, цветная, 80г'
) WHERE slug = 'kancelarijska-hartija';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'binders, file folders, ring binders, lever arch files, folders, organizing documents, A4 files',
    'sr', 'registratori, fascikle, registratori sa obručima, fascikle sa mehanizmom, organizacija dokumenata',
    'ru', 'папки-регистраторы, скоросшиватели, папки с кольцами, папки с арочным механизмом, А4'
) WHERE slug = 'registratori-fascikle';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'desk organizer, pen holder, file tray, document tray, desk accessories, office organization',
    'sr', 'organizacija stola, postolje za olovke, fioke za dokumenta, držači, kancelarijska organizacija',
    'ru', 'органайзеры для стола, подставки для ручек, лотки для документов, офисная организация'
) WHERE slug = 'organizacija-stola';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'calculators, scientific calculator, financial calculator, graphing calculator, office calculator',
    'sr', 'kalkulatori, naučni kalkulator, finansijski kalkulator, grafički kalkulator, kancelarijski',
    'ru', 'калькуляторы, научный калькулятор, финансовый, графический, офисный калькулятор'
) WHERE slug = 'kalkulator-i-oprema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'whiteboard, dry erase board, magnetic whiteboard, glass board, presentation board, markers included',
    'sr', 'bele table, tabla za prezentacije, magnetne table, staklena tabla, tabla za markere',
    'ru', 'маркерные доски, доски для презентаций, магнитные доски, стеклянные доски, white board'
) WHERE slug = 'bele-table-i-prezentacije';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'projectors, presentation equipment, laser pointer, presentation remote, screen projector, HDMI',
    'sr', 'projektna oprema, projektori, laser pokazivač, daljinski za prezentacije, ekran projektor',
    'ru', 'проекционное оборудование, проекторы, лазерная указка, пульт для презентаций, HDMI'
) WHERE slug = 'projektna-oprema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'laminating machine, thermal laminator, laminating pouches, cold laminator, A4 laminator',
    'sr', 'laminatori, termički laminator, folije za laminiranje, hladni laminator, A4 laminator',
    'ru', 'ламинаторы, термические ламинаторы, пленка для ламинирования, холодные ламинаторы'
) WHERE slug = 'laminating';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'paper shredder, document shredder, cross-cut shredder, micro-cut, office shredder, secure',
    'sr', 'lomači dokumenta, uništavač papira, cross-cut, mikro sečenje, kancelarijski lomač, siguran',
    'ru', 'шредеры, уничтожители документов, поперечная резка, микрорезка, офисный шредер'
) WHERE slug = 'lomaci-dokumenta';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'scanners, document scanner, flatbed scanner, photo scanner, portable scanner, duplex scanning',
    'sr', 'skeneri, skener za dokumente, skener ravnog kreveta, foto skener, prenosni, dvostrani',
    'ru', 'сканеры, сканеры документов, планшетные, фотосканеры, портативные, двусторонние'
) WHERE slug = 'skener-i-kopir';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'printers, laser printers, inkjet printers, toner cartridges, ink cartridges, wireless printing',
    'sr', 'štampači, laserski štampači, inkjet, toneri, kertridži, bežično štampanje, all-in-one',
    'ru', 'принтеры, лазерные принтеры, струйные, тонеры, картриджи, беспроводная печать'
) WHERE slug = 'stampaci-toneri';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'archive boxes, storage boxes, document storage, cardboard boxes, file boxes, archiving solutions',
    'sr', 'arhiviranje, kutije za arhiviranje, skladištenje dokumenata, kartonske kutije, arhivske kutije',
    'ru', 'архивные коробки, коробки для хранения, хранение документов, картонные коробки'
) WHERE slug = 'arhiviranje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'office furniture, office desk, office chair, filing cabinet, bookshelf, ergonomic furniture',
    'sr', 'kancelarijski nameštaj, radni sto, kancelarijska stolica, orman, polica, ergonomski',
    'ru', 'офисная мебель, офисный стол, офисное кресло, картотечный шкаф, полки, эргономичная'
) WHERE slug = 'kancelarijski-nameštaj';

-- ============================================================================
-- MUZICKI INSTRUMENTI (Musical Instruments)
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'acoustic guitars, classical guitar, steel string guitar, 6 string, 12 string, dreadnought, Yamaha',
    'sr', 'akustične gitare, klasična gitara, čelične žice, 6 žica, 12 žica, dreadnought, Yamaha',
    'ru', 'акустические гитары, классическая гитара, стальные струны, 6 струн, 12 струн, дредноут'
) WHERE slug = 'gitare-akusticne';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'electric guitars, solid body, semi-hollow, stratocaster, telecaster, les paul, Fender, Gibson',
    'sr', 'električne gitare, solid body, polu-šuplje, stratocaster, telecaster, les paul, Fender, Gibson',
    'ru', 'электрогитары, solid body, полуакустические, стратокастер, телекастер, лес пол, Fender'
) WHERE slug = 'elektricne-gitare';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bass guitars, 4 string bass, 5 string bass, electric bass, fretless bass, precision bass, jazz bass',
    'sr', 'bas gitare, 4 žice, 5 žica, električna bas, fretless bas, precision bass, jazz bass',
    'ru', 'бас-гитары, 4 струны, 5 струн, электробас, безладовый бас, precision, jazz bass'
) WHERE slug = 'bas-gitare';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ukulele, soprano ukulele, concert ukulele, tenor ukulele, mandolin, folk instruments',
    'sr', 'ukulele, sopran ukulele, koncertno ukulele, tenor, mandolina, narodni instrumenti',
    'ru', 'укулеле, сопрано укулеле, концертное укулеле, тенор, мандолина, народные инструменты'
) WHERE slug = 'ukulele-i-mandoline';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'keyboards, digital pianos, synthesizers, MIDI keyboards, 61 keys, 88 keys, weighted keys, Yamaha',
    'sr', 'klavijature, digitalni pianos, sintisajzeri, MIDI, 61 dirki, 88 dirki, otežane dirke, Yamaha',
    'ru', 'клавишные, цифровые пианино, синтезаторы, MIDI клавиатуры, 61 клавиша, 88 клавиш, Yamaha'
) WHERE slug = 'klavijature-i-pianos';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'organs, church organ, harmonium, accordion, pipe organ, electric organ, liturgical instruments',
    'sr', 'orgulje, crkvene orgulje, harmonijum, harmonika, cevne orgulje, električne orgulje',
    'ru', 'органы, церковный орган, гармониум, аккордеон, трубный орган, электрический орган'
) WHERE slug = 'orgulje-i-harmonijum';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'drum sets, acoustic drums, 5 piece drums, cymbals, snare drum, bass drum, drum hardware, Pearl, Tama',
    'sr', 'bubnjevi setovi, akustični bubnjevi, 5 delova, činele, snare, bas bubanj, Pearl, Tama',
    'ru', 'барабанные установки, акустические барабаны, 5 частей, тарелки, малый барабан, Pearl'
) WHERE slug = 'bubnjevi-setovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'wind instruments, trumpet, saxophone, clarinet, flute, trombone, French horn, brass, woodwinds',
    'sr', 'duvački instrumenti, truba, saksofon, klarinet, flauta, trombon, rog, limeni, drveni',
    'ru', 'духовые инструменты, труба, саксофон, кларнет, флейта, тромбон, валторна, медные, деревянные'
) WHERE slug = 'duvacki-instrumenti';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'violin, viola, cello, double bass, string instruments, bow, rosin, strings, classical strings',
    'sr', 'violina, viola, violončelo, kontrabas, žičani instrumenti, gudalo, kalafon, žice, klasični',
    'ru', 'скрипка, альт, виолончель, контрабас, струнные инструменты, смычок, канифоль, струны'
) WHERE slug = 'violina-i-zicani';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'microphones, studio microphone, condenser mic, dynamic mic, USB microphone, XLR, pop filter, Shure',
    'sr', 'mikrofoni za muziku, studijski mikrofon, kondenzatorski, dinamički, USB, XLR, pop filter, Shure',
    'ru', 'микрофоны, студийные микрофоны, конденсаторные, динамические, USB, XLR, поп-фильтр'
) WHERE slug = 'mikrofoni-za-muziku';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'audio equipment, amplifiers, audio interfaces, mixing console, monitors, preamps, studio gear',
    'sr', 'audio oprema za muziku, pojačala, audio interfejsi, mikseta, monitori, predpojačala, studijska',
    'ru', 'аудио оборудование, усилители, аудиоинтерфейсы, микшеры, мониторы, предусилители'
) WHERE slug = 'audio-oprema-za-muziku';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'effects pedals, guitar pedals, distortion, overdrive, delay, reverb, chorus, wah, multi-effects',
    'sr', 'efekti i procesori, gitarske pedale, distorzija, overdrive, delay, reverb, chorus, wah',
    'ru', 'гитарные процессоры, педали эффектов, дисторшн, овердрайв, делей, реверб, хорус'
) WHERE slug = 'efekti-i-procesori';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'DJ equipment, DJ controller, turntables, mixer, CDJ, DJ software, Pioneer, Serato, Traktor',
    'sr', 'DJ oprema, DJ kontroler, gramofoni, mikseta, CDJ, DJ softver, Pioneer, Serato, Traktor',
    'ru', 'DJ оборудование, DJ контроллеры, вертушки, микшеры, CDJ, DJ софт, Pioneer, Serato'
) WHERE slug = 'dj-oprema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'music accessories, guitar strings, drum sticks, reeds, picks, cables, tuners, capos, straps',
    'sr', 'muzički dodaci, žice za gitaru, palice za bubnjeve, trstike, trzalice, kablovi, tuneri, kapo',
    'ru', 'музыкальные аксессуары, струны для гитары, барабанные палочки, медиаторы, кабели, тюнеры'
) WHERE slug = 'muzicki-dodaci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'recording, music production, home studio, audio recording, mixing, mastering, DAW, Pro Tools, Logic',
    'sr', 'snimanje i produkcija, muzička produkcija, kućni studio, audio snimanje, miksovanje, mastering',
    'ru', 'запись, музыкальная продукция, домашняя студия, аудиозапись, сведение, мастеринг, DAW'
) WHERE slug = 'snimanje-i-produkcija';

-- ============================================================================
-- UMETNOST I RUKOTVORINE (Art and Crafts)
-- ============================================================================

-- Slikanje
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'oil paints, oil painting, artist oil colors, fine art paints, oil paint sets, titanium white',
    'sr', 'uljane boje, slikanje uljanim bojama, umetničke uljane boje, setovi uljanih boja',
    'ru', 'масляные краски, масляная живопись, художественные масляные краски, наборы масла'
) WHERE slug = 'boje-uljane';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'acrylic paints, acrylic painting, fast drying paints, acrylic paint sets, artist acrylics',
    'sr', 'akrilne boje, slikanje akrilom, brzo sušeće boje, setovi akrilnih boja, umetničke akrilne',
    'ru', 'акриловые краски, акриловая живопись, быстросохнущие краски, наборы акрила'
) WHERE slug = 'boje-akrilne';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'watercolors, watercolor paints, aquarelle, watercolor sets, pan watercolors, tube watercolors',
    'sr', 'vodene boje, akvarel, akvarelne boje, setovi vodenih boja, vodene u tubama, u kockicama',
    'ru', 'акварельные краски, акварель, наборы акварели, акварель в кюветах, в тюбиках'
) WHERE slug = 'boje-vodene';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'colored pencils, artist pencils, drawing pencils, watercolor pencils, professional pencils, sets',
    'sr', 'olovke u boji, umetničke olovke, olovke za crtanje, akvarelne olovke, profesionalne, setovi',
    'ru', 'цветные карандаши, художественные карандаши, для рисования, акварельные, профессиональные'
) WHERE slug = 'boje-u-olovci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'drawing pencils, graphite pencils, sketching pencils, charcoal pencils, 2B, 4B, 6B, 8B, HB, 2H',
    'sr', 'olovke za crtanje, grafitne olovke, olovke za skiciranje, ugalj, 2B, 4B, 6B, 8B, HB, 2H',
    'ru', 'карандаши для рисования, графитовые, для эскизов, угольные, 2B, 4B, 6B, 8B, HB, 2H'
) WHERE slug = 'olovke-za-crtanje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'paint brushes, artist brushes, oil brushes, watercolor brushes, acrylic brushes, synthetic, natural',
    'sr', 'četkice za slikanje, umetničke četkice, za ulje, za akvarel, za akril, sintetičke, prirodne',
    'ru', 'кисти для живописи, художественные кисти, для масла, для акварели, для акрила, синтетика'
) WHERE slug = 'cetkice-za-slikanje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'canvases, stretched canvas, canvas boards, primed canvas, linen canvas, cotton canvas, art canvases',
    'sr', 'platna za slikanje, razvučeno platno, platno na dasci, grundirano platno, lan, pamuk',
    'ru', 'холсты для живописи, натянутые холсты, холсты на картоне, грунтованные, лён, хлопок'
) WHERE slug = 'platna-za-slikanje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'easels, table easel, studio easel, field easel, H-frame easel, portable easel, adjustable',
    'sr', 'molberti, stoni molbert, studio molbert, terenski molbert, H-ram, prenosni molbert, podesivi',
    'ru', 'мольберты, настольный мольберт, студийный, пленэрный, переносной, регулируемый'
) WHERE slug = 'molberti';

-- Rukotvorine
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'polymer clay, modeling clay, air dry clay, oven bake clay, sculpting clay, FIMO, Sculpey',
    'sr', 'glina za modelovanje, polimerna glina, glina na vazduhu, glina za pečenje, grnčarska glina',
    'ru', 'полимерная глина, моделирующая глина, самозатвердевающая, запекаемая, FIMO, Sculpey'
) WHERE slug = 'glina-za-modelovanje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'plaster, casting plaster, plaster of Paris, sculpting plaster, molds, plaster crafts, gypsum',
    'sr', 'gips, odlivni gips, pariski gips, gips za modelovanje, kalupi, gipsana rukotvorina',
    'ru', 'гипс, формовочный гипс, парижский гипс, скульптурный гипс, формы, гипсовое литье'
) WHERE slug = 'gips';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'embroidery thread, cotton thread, silk thread, DMC, cross stitch, embroidery floss, sewing thread',
    'sr', 'konac za vez, pamučni konac, svileni konac, DMC, krst vez, konac za vezenje, konac za šivenje',
    'ru', 'нитки для вышивания, хлопковые нитки, шёлковые, DMC, мулине, нитки для вышивки крестом'
) WHERE slug = 'konac-za-vez';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'embroidery fabric, aida cloth, linen fabric, cross stitch fabric, even-weave, canvas for embroidery',
    'sr', 'platno za vez, aida platno, laneno platno, platno za krst vez, kanvas za vez, etan vez',
    'ru', 'ткань для вышивания, канва aida, льняная ткань, для вышивки крестом, равномерка'
) WHERE slug = 'platno-za-vez';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'beads, seed beads, glass beads, crystal beads, wooden beads, ceramic beads, bead stringing',
    'sr', 'perle za nakit, seme perle, staklene perle, kristalne perle, drvene perle, keramičke perle',
    'ru', 'бусины для украшений, бисер, стеклянные бусины, хрустальные, деревянные, керамические'
) WHERE slug = 'perle-za-nakit';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'pearls, freshwater pearls, cultured pearls, pearl beads, pearl jewelry, baroque pearls',
    'sr', 'biser, slatkovodni biser, gajeni biser, biserne perle, biserni nakit, barokni biser',
    'ru', 'жемчуг, пресноводный жемчуг, культивированный жемчуг, жемчужные бусины, барочный жемчуг'
) WHERE slug = 'biser';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'decoupage, napkin decoupage, decoupage glue, decoupage paper, mod podge, craft projects',
    'sr', 'decoupage setovi, salveta decoupage, decoupage lepak, decoupage papir, mod podge',
    'ru', 'декупаж, салфетки для декупажа, клей для декупажа, бумага для декупажа, mod podge'
) WHERE slug = 'decoupage-setovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'handmade, handcrafted, DIY crafts, craft projects, handmade gifts, artisan crafts, custom made',
    'sr', 'ručni rad, ručna izrada, DIY projekti, craft projekti, ručno pravljeni pokloni, ručno napravljeno',
    'ru', 'ручная работа, ручное изготовление, DIY проекты, крафт проекты, подарки ручной работы'
) WHERE slug = 'rucni-rad';

-- ============================================================================
-- USLUGE (Services)
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'handyman, general repairs, home repairs, minor repairs, maintenance, jack of all trades, odd jobs',
    'sr', 'majstor, opšti radovi, popravke u kući, sitni poslovi, održavanje, razni popravci',
    'ru', 'мастер на все руки, общие работы, мелкий ремонт, обслуживание, разнорабочий'
) WHERE slug = 'majstor-general';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'electrician, electrical services, wiring, electrical repairs, circuit breaker, outlet installation',
    'sr', 'električar, električne usluge, ožičenje, električne popravke, osigurači, utičnice, instalacije',
    'ru', 'электрик, электромонтажные работы, проводка, электроремонт, автоматы, розетки'
) WHERE slug = 'elektricar';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'plumber, plumbing services, pipe repair, drain cleaning, leak repair, water heater, faucet installation',
    'sr', 'vodoinstalater, vodoinstalaterske usluge, popravka cevi, otpušavanje, popravka curenja',
    'ru', 'сантехник, сантехнические услуги, ремонт труб, прочистка, устранение протечек'
) WHERE slug = 'vodoinstalater';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'painter, painting services, interior painting, exterior painting, wall painting, house painting',
    'sr', 'moler, farbanje, farbanje zidova, unutrašnje farbanje, spoljašnje farbanje, krečenje kuće',
    'ru', 'маляр, покрасочные работы, покраска стен, внутренняя покраска, наружная покраска'
) WHERE slug = 'moler-farbanje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'locksmith, lock installation, key duplication, lock repair, emergency locksmith, safe opening',
    'sr', 'bravar, bravar usluge, ugradnja brave, kopiranje ključeva, hitna intervencija, otključavanje',
    'ru', 'слесарь, установка замков, изготовление ключей, ремонт замков, вскрытие замков'
) WHERE slug = 'bravar';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'carpenter, carpentry services, custom furniture, woodworking, wood repairs, cabinet installation',
    'sr', 'stolar, stolarija, izrada nameštaja po meri, obrada drveta, popravke, ugradnja elemenata',
    'ru', 'столяр, столярные услуги, мебель на заказ, работа с деревом, ремонт, установка шкафов'
) WHERE slug = 'stolar';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tiling, tile installation, bathroom tiling, floor tiling, wall tiling, ceramic tiles, porcelain',
    'sr', 'keramičar, postavljanje pločica, pločice kupatilo, podne pločice, zidne pločice, keramika',
    'ru', 'плиточник, укладка плитки, плитка в ванной, напольная плитка, настенная, керамика'
) WHERE slug = 'keramicar';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'flooring, parquet installation, wood flooring, laminate flooring, floor restoration, sanding',
    'sr', 'parketar, postavljanje parketa, drveni podovi, laminat, renoviranje podova, brušenje',
    'ru', 'паркетчик, укладка паркета, деревянные полы, ламинат, реставрация полов, шлифовка'
) WHERE slug = 'parketar';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'roofing, roof repair, roof installation, tile roofing, metal roofing, leak repair, gutter installation',
    'sr', 'krovopokrivač, popravka krova, montaža krova, criep, lim, popravka curenja, oluk',
    'ru', 'кровельщик, ремонт крыши, монтаж кровли, черепица, металлочерепица, ремонт протечек'
) WHERE slug = 'krovopokrivac';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'appliance repair, washing machine repair, refrigerator repair, dishwasher repair, oven repair',
    'sr', 'servis bele tehnike, popravka veš mašine, frižider servis, sudopera, rerna, popravke',
    'ru', 'ремонт бытовой техники, ремонт стиральных машин, холодильников, посудомоек, духовок'
) WHERE slug = 'servis-bele-tehnike';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'TV repair, television repair, LCD repair, LED TV repair, plasma TV, screen repair, TV service',
    'sr', 'servis televizora, TV popravke, LCD servis, LED TV, plazma, popravka ekrana, TV servis',
    'ru', 'ремонт телевизоров, ремонт LCD, LED TV, плазменных телевизоров, ремонт экранов'
) WHERE slug = 'servis-tv';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'laptop repair, computer repair, MacBook repair, screen replacement, keyboard repair, virus removal',
    'sr', 'servis laptopa, popravka računara, MacBook servis, zamena ekrana, tastatura, virusi',
    'ru', 'ремонт ноутбуков, ремонт компьютеров, MacBook ремонт, замена экрана, клавиатуры'
) WHERE slug = 'servis-laptopa';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'phone repair, smartphone repair, iPhone repair, screen replacement, battery replacement, water damage',
    'sr', 'servis telefona, popravka smartfona, iPhone servis, zamena ekrana, baterija, voda oštećenje',
    'ru', 'ремонт телефонов, ремонт смартфонов, iPhone ремонт, замена экрана, батареи, после воды'
) WHERE slug = 'servis-telefona';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'photographer, event photography, wedding photographer, party photographer, birthday photography',
    'sr', 'fotograf za događaje, fotografija, venčanja, rođendani, proslave, event foto, profesionalni',
    'ru', 'фотограф на мероприятия, свадебный фотограф, дни рождения, праздники, event фото'
) WHERE slug = 'fotograf-event';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'portrait photographer, family photos, professional portraits, studio photography, headshots',
    'sr', 'fotograf portret, porodične fotografije, profesionalni portreti, studijska fotografija',
    'ru', 'портретный фотограф, семейные фото, профессиональные портреты, студийная съёмка'
) WHERE slug = 'fotograf-portret';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'videographer, event videography, wedding videography, promotional videos, drone videography',
    'sr', 'snimanje video, videograf, venčanja video, promocioni video, dronovi, korporativni video',
    'ru', 'видеосъёмка, видеограф, свадебная видеосъёмка, промо видео, съёмка с дрона'
) WHERE slug = 'snimanje-video';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'video editing, post production, video montage, color grading, motion graphics, YouTube editing',
    'sr', 'montaža videa, postprodukcija, video montaža, korekcija boja, motion graphics, YouTube',
    'ru', 'монтаж видео, постпродакшн, видеомонтаж, цветокоррекция, моушн графика, YouTube'
) WHERE slug = 'montaza-video';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'graphic design, logo design, branding, flyer design, poster design, social media graphics, Adobe',
    'sr', 'grafički dizajn, dizajn logoa, brending, flajeri, posteri, grafika za društvene mreže, Adobe',
    'ru', 'графический дизайн, дизайн логотипа, брендинг, флаеры, постеры, графика для соцсетей'
) WHERE slug = 'graficki-dizajn';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'web design, website design, UX/UI design, landing page, responsive design, Figma, Webflow',
    'sr', 'web dizajn, dizajn veb sajta, UX/UI dizajn, landing page, responsive dizajn, Figma, Webflow',
    'ru', 'веб-дизайн, дизайн сайтов, UX/UI дизайн, лендинг, адаптивный дизайн, Figma, Webflow'
) WHERE slug = 'web-dizajn';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'translation services, language translation, certified translation, document translation, interpreter',
    'sr', 'prevod jezika, usluge prevođenja, overen prevod, prevod dokumenata, tumač, sudski tumač',
    'ru', 'перевод, языковой перевод, заверенный перевод, перевод документов, переводчик, нотариальный'
) WHERE slug = 'prevod-jezik';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'transcription, audio transcription, video transcription, typing services, subtitle creation',
    'sr', 'prepis teksta, transkripcija, audio transkripcija, video transkripcija, kucanje, titlovi',
    'ru', 'транскрибация, расшифровка аудио, расшифровка видео, набор текста, создание субтитров'
) WHERE slug = 'prepis-teksta';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'math tutor, mathematics, algebra, geometry, calculus, trigonometry, private lessons, exam prep',
    'sr', 'korepetitor matematike, instrukcije matematike, algebra, geometrija, kalkulus, priprema ispita',
    'ru', 'репетитор математики, алгебра, геометрия, математический анализ, подготовка к экзаменам'
) WHERE slug = 'korepetitor-matematika';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'English tutor, English lessons, IELTS, TOEFL, grammar, conversation practice, business English',
    'sr', 'korepetitor engleskog, časovi engleskog, IELTS, TOEFL, gramatika, konverzacija, poslovni engleski',
    'ru', 'репетитор английского, уроки английского, IELTS, TOEFL, грамматика, разговорный, деловой'
) WHERE slug = 'korepetitor-engleski';

-- ============================================================================
-- OSTALO (Other / Miscellaneous)
-- ============================================================================

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'collectibles, collectible cards, trading cards, sports cards, Pokemon, Magic The Gathering, Yu-Gi-Oh',
    'sr', 'kolekcionarske karte, trading cards, sportske karte, Pokemon, Magic The Gathering, Yu-Gi-Oh',
    'ru', 'коллекционные карточки, трейдинговые карты, спортивные карты, Pokemon, Magic, Yu-Gi-Oh'
) WHERE slug = 'kolekcionarske-karte';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'collectible coins, numismatics, rare coins, gold coins, silver coins, commemorative coins',
    'sr', 'kolekcionarski novčići, numizmatika, retki novčići, zlatnici, srebrnjaci, komemorativni',
    'ru', 'коллекционные монеты, нумизматика, редкие монеты, золотые, серебряные, памятные монеты'
) WHERE slug = 'kolekcionarski-novcici';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'stamps, philately, rare stamps, stamp collection, postage stamps, commemorative stamps',
    'sr', 'markice, filatelija, retke markice, kolekcija markica, poštanske markice, komemorativne',
    'ru', 'марки, филателия, редкие марки, коллекция марок, почтовые марки, памятные марки'
) WHERE slug = 'markice-filatelija';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'fossils, dinosaur fossils, ammonites, petrified wood, trilobites, paleontology, prehistoric',
    'sr', 'fosili, dinosaurusi fosili, amoniti, okamenjeno drvo, trilobiti, paleontologija, praistorija',
    'ru', 'окаменелости, ископаемые, динозавры, аммониты, окаменелое дерево, трилобиты, палеонтология'
) WHERE slug = 'fosili';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'minerals, crystals, gemstones, quartz, amethyst, geodes, healing stones, rock collection',
    'sr', 'minerali i kamenje, kristali, drago kamenje, kvarc, ametist, geode, lečenje kamenjem',
    'ru', 'минералы и камни, кристаллы, драгоценные камни, кварц, аметист, жеоды, целебные камни'
) WHERE slug = 'minerali-i-kamenje';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'antique furniture, vintage furniture, antique chairs, antique tables, retro furniture, collectible',
    'sr', 'antikvitetni nameštaj, vintage nameštaj, antikne stolice, antikni stolovi, retro, kolekcionarski',
    'ru', 'антикварная мебель, винтажная мебель, антикварные стулья, столы, ретро, коллекционная'
) WHERE slug = 'antikviteti-namestaj';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'antique clocks, vintage clocks, pocket watch, wall clocks, grandfather clocks, mechanical clocks',
    'sr', 'antikvitetni satovi, vintage satovi, džepni sat, zidni satovi, stari satovi, mehanički satovi',
    'ru', 'антикварные часы, винтажные часы, карманные часы, настенные, напольные, механические'
) WHERE slug = 'antikviteti-satovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'antique paintings, vintage art, oil paintings, collectible art, fine art, classic paintings',
    'sr', 'antikvitetne slike, vintage umetnost, uljane slike, kolekcionarske slike, fina umetnost',
    'ru', 'антикварные картины, винтажное искусство, масляная живопись, коллекционные картины'
) WHERE slug = 'antikviteti-slike';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'vintage clothing men, retro fashion men, vintage jackets, vintage shirts, second hand, thrift',
    'sr', 'vintage muška odeća, retro moda muškarci, vintage jakne, vintage košulje, second hand',
    'ru', 'винтажная мужская одежда, ретро мода мужчины, винтажные куртки, рубашки, секонд хенд'
) WHERE slug = 'vintage-odeca-muska';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'vintage clothing women, retro fashion women, vintage dresses, vintage blouses, second hand',
    'sr', 'vintage ženska odeća, retro moda žene, vintage haljine, vintage bluze, second hand, thrift',
    'ru', 'винтажная женская одежда, ретро мода женщины, винтажные платья, блузки, секонд хенд'
) WHERE slug = 'vintage-odeca-zenska';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'retro electronics, vintage electronics, old radios, vinyl players, cassette players, 80s 90s tech',
    'sr', 'retro elektronika, vintage elektronika, stari radio, gramofoni, kaseta plejeri, 80s 90s',
    'ru', 'ретро электроника, винтажная электроника, старые радио, виниловые проигрыватели, 80х 90х'
) WHERE slug = 'retro-elektronika';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'corporate gifts, business gifts, branded gifts, promotional products, gift sets, executive gifts',
    'sr', 'korporativni pokloni, poslovni pokloni, brendirani pokloni, promotivni proizvodi, seti poklona',
    'ru', 'корпоративные подарки, бизнес-подарки, брендированные подарки, промо-продукция, наборы'
) WHERE slug = 'korporativni-pokloni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'corporate gifts, business gifts, branded gifts, promotional products, gift sets, executive gifts',
    'sr', 'korporativni pokloni, poslovni pokloni, brendirani pokloni, promotivni proizvodi, seti poklona',
    'ru', 'корпоративные подарки, бизнес-подарки, брендированные подарки, промо-продукция, наборы'
) WHERE slug = 'korporativni-pokloni-misc';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'personalized gifts, custom gifts, engraved gifts, photo gifts, monogrammed, name gifts, unique gifts',
    'sr', 'personalizovani pokloni, prilagođeni pokloni, gravirani pokloni, foto pokloni, monogrami, sa imenom',
    'ru', 'персонализированные подарки, индивидуальные, гравировка, фото-подарки, с монограммой'
) WHERE slug = 'personalizovani-pokloni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'custom products, made to order, bespoke, personalized, custom made, tailored, handmade custom',
    'sr', 'custom proizvodi, po meri, bespoke, personalizovano, ručno izrađeno, prilagođeno',
    'ru', 'индивидуальные продукты, на заказ, кастомные, персонализированные, ручная работа'
) WHERE slug = 'custom-proizvodi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'miscellaneous, various items, random, mixed, diverse products, other items, uncategorized',
    'sr', 'razno, razni proizvodi, random, mešovito, raznovrsno, ostali proizvodi, nekategorisano',
    'ru', 'разное, различные товары, смешанное, разнообразные продукты, прочее, без категории'
) WHERE slug = 'razno';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'uncategorized, not categorized, miscellaneous items, other, various, unsorted, uncategorized products',
    'sr', 'nerazvrstavano, nekategorisano, razni proizvodi, ostalo, nesortiran, proizvodi bez kategorije',
    'ru', 'без категории, несортированное, разные товары, прочее, различное, неклассифицированное'
) WHERE slug = 'nerazvrstavano';

-- Migration complete
