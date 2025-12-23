-- Migration: Add meta_keywords for all Elektronika subcategories (L2/L3)
-- Generated: 2025-12-22
-- Categories: 135 subcategories under elektronika

-- L2 Categories

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', '3d printer, fdm printer, resin printer, 3d printing, filament, pla, abs, petg, creality, prusa, anycubic, ender 3, 3d print',
    'sr', '3d štampač, fdm štampač, resin štampač, 3d štampa, filament, pla, abs, petg, creality, prusa, anycubic, ender 3',
    'ru', '3д принтер, fdm принтер, фотополимерный принтер, 3д печать, филамент, pla, abs, petg, creality, prusa, anycubic'
) WHERE slug = '3d-stampaci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'audio equipment, speakers, amplifiers, headphones, hi-fi, sound system, stereo, home audio, professional audio, studio monitors',
    'sr', 'audio oprema, zvučnici, pojačala, slušalice, hi-fi sistem, stereo, kućni audio, profesionalni audio, studijski monitori',
    'ru', 'аудио оборудование, колонки, усилители, наушники, hi-fi, стерео система, домашний аудио, профессиональный звук'
) WHERE slug = 'audio-oprema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'professional bluetooth headphones, wireless headphones, noise cancelling, studio headphones, sony wh-1000xm5, bose 700, sennheiser momentum',
    'sr', 'profesionalne bluetooth slušalice, bežične slušalice, noise cancelling, studijske slušalice, sony wh-1000xm5, bose 700, sennheiser',
    'ru', 'профессиональные bluetooth наушники, беспроводные наушники, шумоподавление, студийные наушники, sony wh-1000xm5, bose 700'
) WHERE slug = 'bluetooth-slusalice-profesionalne';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bluetooth keyboard, wireless keyboard, mechanical keyboard, ergonomic keyboard, logitech mx keys, keychron, anne pro',
    'sr', 'bluetooth tastatura, bežična tastatura, mehanička tastatura, ergonomska tastatura, logitech mx keys, keychron, anne pro',
    'ru', 'bluetooth клавиатура, беспроводная клавиатура, механическая клавиатура, эргономичная клавиатура, logitech mx keys'
) WHERE slug = 'bluetooth-tastature';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'premium bluetooth speaker, wireless speaker, portable speaker, jbl, bose soundlink, ultimate ears, sony srs, waterproof speaker',
    'sr', 'premium bluetooth zvučnik, bežični zvučnik, prenosivi zvučnik, jbl, bose soundlink, ultimate ears, sony srs, vodootporni zvučnik',
    'ru', 'премиум bluetooth колонка, беспроводная колонка, портативная колонка, jbl, bose soundlink, ultimate ears, sony srs'
) WHERE slug = 'bluetooth-zvucnici-premium';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'e-ink reader, e-book reader, kindle, kobo, pocketbook, onyx boox, e-reader, digital book reader, e-ink display',
    'sr', 'e-ink čitač, elektronska knjiga, kindle, kobo, pocketbook, onyx boox, digitalni čitač, e-ink ekran',
    'ru', 'e-ink ридер, электронная книга, kindle, kobo, pocketbook, onyx boox, читалка, e-ink дисплей'
) WHERE slug = 'citaci-e-ink';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'desktop computer, pc, gaming pc, workstation, all-in-one, desktop tower, custom pc, office computer',
    'sr', 'desktop računar, pc, gaming pc, radna stanica, all-in-one, kućište, custom pc, kancelarijski računar',
    'ru', 'настольный компьютер, пк, игровой компьютер, рабочая станция, моноблок, системный блок, custom pc'
) WHERE slug = 'desktop-racunari';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'electronics accessories, phone accessories, computer accessories, cables, adapters, chargers, cases, screen protectors',
    'sr', 'dodatna oprema za elektroniku, dodaci za telefon, dodaci za računar, kablovi, adapteri, punjači, maske, zaštita',
    'ru', 'аксессуары для электроники, аксессуары для телефона, аксессуары для компьютера, кабели, адаптеры, зарядки'
) WHERE slug = 'dodatna-oprema-elektronika';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'drone, quadcopter, fpv drone, camera drone, dji drone, racing drone, mini drone, professional drone',
    'sr', 'dron, quadcopter, fpv dron, dron sa kamerom, dji dron, trkaći dron, mini dron, profesionalni dron',
    'ru', 'дрон, квадрокоптер, fpv дрон, дрон с камерой, dji дрон, гоночный дрон, мини дрон, профессиональный дрон'
) WHERE slug = 'dronovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'e-reader, ebook reader, kindle, kobo, electronic book, digital reader, e-ink, pocketbook',
    'sr', 'elektronski čitač, kindle, kobo, elektronska knjiga, digitalni čitač, e-ink, pocketbook',
    'ru', 'электронная книга, kindle, kobo, ридер, цифровая книга, e-ink, pocketbook'
) WHERE slug = 'e-citaci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'external ssd, portable ssd, external hard drive, usb ssd, samsung t7, sandisk extreme, nvme ssd external',
    'sr', 'eksterni ssd, prenosivi ssd, eksterni hard disk, usb ssd, samsung t7, sandisk extreme, nvme eksterni',
    'ru', 'внешний ssd, портативный ssd, внешний жесткий диск, usb ssd, samsung t7, sandisk extreme'
) WHERE slug = 'eksterni-hard-diskovi-ssd';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ebook reader, digital book reader, kindle, kobo, e-ink reader, electronic books, reading device',
    'sr', 'čitač elektronskih knjiga, kindle, kobo, e-ink čitač, elektronske knjige, uređaj za čitanje',
    'ru', 'читалка электронных книг, kindle, kobo, e-ink ридер, электронные книги, устройство для чтения'
) WHERE slug = 'elektronske-knjige-citalice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'camera, photo camera, video camera, dslr, mirrorless, camcorder, action camera, digital camera, professional camera',
    'sr', 'foto aparat, video kamera, dslr, mirrorless, kamkorder, action kamera, digitalni foto aparat, profesionalna kamera',
    'ru', 'фотоаппарат, видеокамера, dslr, беззеркальная камера, экшн камера, цифровой фотоаппарат, профессиональная камера'
) WHERE slug = 'foto-i-video-kamere';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gaming equipment, gaming gear, gaming accessories, gaming peripherals, rgb lighting, gaming setup',
    'sr', 'gaming oprema, gejmerska oprema, gaming dodaci, gejmerska periferija, rgb rasveta, gaming setup',
    'ru', 'игровое оборудование, геймерское оборудование, игровые аксессуары, игровая периферия, rgb подсветка'
) WHERE slug = 'gaming-oprema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'graphics tablet, drawing tablet, pen tablet, wacom, huion, xp-pen, digital art tablet, display tablet',
    'sr', 'grafička tabla, tabla za crtanje, pen tabla, wacom, huion, xp-pen, digitalna tabla, display tabla',
    'ru', 'графический планшет, планшет для рисования, wacom, huion, xp-pen, дисплейный планшет, цифровое искусство'
) WHERE slug = 'graficke-tablice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'calculator, scientific calculator, graphing calculator, financial calculator, casio, texas instruments, hp calculator',
    'sr', 'kalkulator, naučni kalkulator, grafički kalkulator, finansijski kalkulator, casio, texas instruments, hp kalkulator',
    'ru', 'калькулятор, научный калькулятор, графический калькулятор, финансовый калькулятор, casio, texas instruments'
) WHERE slug = 'kalkulatori';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gaming console, playstation, xbox, nintendo switch, ps5, xbox series x, gaming, video games',
    'sr', 'gaming konzola, playstation, xbox, nintendo switch, ps5, xbox series x, igranje, video igre',
    'ru', 'игровая консоль, playstation, xbox, nintendo switch, ps5, xbox series x, видеоигры'
) WHERE slug = 'konzole-i-gaming';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'console games, video games, ps5 games, xbox games, nintendo games, playstation games, physical games, game discs',
    'sr', 'konzolne igre, video igre, ps5 igre, xbox igre, nintendo igre, playstation igre, fizičke igre, diskovi igara',
    'ru', 'консольные игры, видеоигры, игры ps5, игры xbox, игры nintendo, игры playstation, физические игры'
) WHERE slug = 'konzolne-igre';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'laptop, notebook, laptop computer, portable computer, ultrabook, gaming laptop, business laptop, chromebook',
    'sr', 'laptop, notebook, prenosivi računar, ultrabook, gaming laptop, poslovni laptop, chromebook',
    'ru', 'ноутбук, лаптоп, портативный компьютер, ультрабук, игровой ноутбук, бизнес ноутбук, chromebook'
) WHERE slug = 'laptop-racunari';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'professional memory card, sd card, microsd card, cf card, sdxc, uhs-ii, v90, high speed memory card',
    'sr', 'profesionalna memorijska kartica, sd kartica, microsd kartica, cf kartica, sdxc, uhs-ii, v90, brza memorijska kartica',
    'ru', 'профессиональная карта памяти, sd карта, microsd карта, cf карта, sdxc, uhs-ii, v90, высокоскоростная карта памяти'
) WHERE slug = 'memorijske-kartice-pro';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'microphone, usb microphone, condenser microphone, dynamic microphone, streaming mic, podcast mic, studio microphone, rode, shure, blue yeti',
    'sr', 'mikrofon, usb mikrofon, kondenzatorski mikrofon, dinamički mikrofon, streaming mikrofon, podcast mikrofon, studijski mikrofon, rode, shure, blue yeti',
    'ru', 'микрофон, usb микрофон, конденсаторный микрофон, динамический микрофон, стриминг микрофон, подкаст микрофон, студийный микрофон'
) WHERE slug = 'mikrofoni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'digital microscope, usb microscope, portable microscope, inspection microscope, electronics microscope, soldering microscope',
    'sr', 'digitalni mikroskop, usb mikroskop, prenosivi mikroskop, inspekcijski mikroskop, mikroskop za elektroniku, mikroskop za lemljenje',
    'ru', 'цифровой микроскоп, usb микроскоп, портативный микроскоп, инспекционный микроскоп, микроскоп для электроники'
) WHERE slug = 'mikroskopi-digitalni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'networking, internet equipment, router, switch, network adapter, ethernet, wifi, network cable, modem',
    'sr', 'mrežna oprema, internet oprema, ruter, svič, mrežni adapter, ethernet, wifi, mrežni kabl, modem',
    'ru', 'сетевое оборудование, интернет оборудование, роутер, свитч, сетевой адаптер, ethernet, wifi, сетевой кабель'
) WHERE slug = 'mreza-i-internet';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nas, network storage, network attached storage, synology, qnap, nas server, home server, data storage, raid storage',
    'sr', 'nas, mrežno skladište, network attached storage, synology, qnap, nas server, kućni server, skladištenje podataka, raid',
    'ru', 'nas, сетевое хранилище, network attached storage, synology, qnap, nas сервер, домашний сервер, хранение данных'
) WHERE slug = 'nas-i-storage';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'smart bulb, smart light, led smart bulb, philips hue, smart lighting, wifi bulb, color changing bulb, app controlled light',
    'sr', 'pametna sijalica, pametno svetlo, led pametna sijalica, philips hue, pametno osvetljenje, wifi sijalica, sijalica koja menja boju',
    'ru', 'умная лампочка, умный свет, led умная лампа, philips hue, умное освещение, wifi лампа, лампа с изменением цвета'
) WHERE slug = 'pametne-sijalice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'smartwatch, smart watch, fitness watch, wearable, apple watch, samsung galaxy watch, garmin, huawei watch, wear os',
    'sr', 'pametni sat, smart watch, fitnes sat, nosiva tehnologija, apple watch, samsung galaxy watch, garmin, huawei watch',
    'ru', 'умные часы, смарт часы, фитнес часы, носимая технология, apple watch, samsung galaxy watch, garmin, huawei watch'
) WHERE slug = 'pametni-satovi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'smartphone, mobile phone, cell phone, android phone, iphone, 5g phone, dual sim, unlocked phone, flagship phone',
    'sr', 'pametni telefon, mobilni telefon, android telefon, iphone, 5g telefon, dual sim, otključan telefon, flagship telefon',
    'ru', 'смартфон, мобильный телефон, андроид телефон, айфон, 5g телефон, двухсимочный, разблокированный телефон'
) WHERE slug = 'pametni-telefoni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'computer peripherals, mouse, keyboard, webcam, speakers, headset, monitor, usb hub, gaming peripherals',
    'sr', 'računarska periferija, miš, tastatura, web kamera, zvučnici, slušalice, monitor, usb hub, gaming periferija',
    'ru', 'компьютерная периферия, мышь, клавиатура, веб камера, колонки, наушники, монитор, usb хаб, игровая периферия'
) WHERE slug = 'periferija';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'powerline adapter, powerline networking, ethernet over power, tp-link powerline, devolo, zyxel, home network extender',
    'sr', 'powerline adapter, powerline mreža, internet preko struje, tp-link powerline, devolo, zyxel, proširivač mreže',
    'ru', 'powerline адаптер, powerline сеть, интернет по электросети, tp-link powerline, devolo, zyxel, расширитель сети'
) WHERE slug = 'powerline-adapteri';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'projector, video projector, home theater projector, portable projector, 4k projector, hd projector, mini projector, epson, benq, optoma',
    'sr', 'projektor, video projektor, kućni bioskop projektor, prenosivi projektor, 4k projektor, hd projektor, mini projektor, epson, benq',
    'ru', 'проектор, видеопроектор, домашний кинотеатр проектор, портативный проектор, 4k проектор, hd проектор, мини проектор'
) WHERE slug = 'projektori';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'wireless charger, qi charger, fast wireless charging, magsafe charger, wireless charging pad, wireless charging stand',
    'sr', 'bežični punjač, qi punjač, brzo bežično punjenje, magsafe punjač, bežični punjač podloga, bežični punjač stalak',
    'ru', 'беспроводное зарядное устройство, qi зарядка, быстрая беспроводная зарядка, magsafe зарядка, беспроводная зарядная панель'
) WHERE slug = 'punjaci-bežicni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'computer components, pc parts, cpu, gpu, motherboard, ram, ssd, power supply, cooling, pc build',
    'sr', 'računarske komponente, pc delovi, procesor, grafička, matična ploča, ram, ssd, napajanje, hlađenje, pc build',
    'ru', 'компьютерные комплектующие, pc компоненты, процессор, видеокарта, материнская плата, оперативная память, ssd, блок питания'
) WHERE slug = 'racunarske-komponente';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'scanner, document scanner, flatbed scanner, photo scanner, portable scanner, duplex scanner, canon scanner, epson scanner',
    'sr', 'skener, skener dokumenata, flatbed skener, foto skener, prenosivi skener, duplex skener, canon skener, epson skener',
    'ru', 'сканер, сканер документов, планшетный сканер, фото сканер, портативный сканер, дуплексный сканер, canon, epson'
) WHERE slug = 'skeneri';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'smart home, home automation, smart devices, iot, alexa, google home, smart hub, home assistant, smart home system',
    'sr', 'pametna kuća, automatizacija kuće, pametni uređaji, iot, alexa, google home, smart hub, home assistant, smart home sistem',
    'ru', 'умный дом, домашняя автоматизация, умные устройства, iot, alexa, google home, smart hub, home assistant'
) WHERE slug = 'smart-home';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'fitness tracker, smart band, activity tracker, fitness band, xiaomi mi band, fitbit, honor band, heart rate monitor',
    'sr', 'fitnes narukvica, pametna narukvica, activity tracker, fitnes traker, xiaomi mi band, fitbit, honor band, monitor otkucaja',
    'ru', 'фитнес браслет, умный браслет, трекер активности, фитнес трекер, xiaomi mi band, fitbit, honor band, пульсометр'
) WHERE slug = 'smart-narukvice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'solar charger, solar power bank, portable solar charger, solar panel charger, camping solar charger, outdoor charger',
    'sr', 'solarni punjač, solarni power bank, prenosivi solarni punjač, punjač sa solarnim panelom, punjač za kampovanje',
    'ru', 'солнечное зарядное устройство, солнечный power bank, портативная солнечная зарядка, солнечная панель зарядка'
) WHERE slug = 'solarni-punjaci';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gimbal stabilizer, camera stabilizer, smartphone gimbal, dji osmo, zhiyun, moza, handheld stabilizer, 3-axis gimbal',
    'sr', 'gimbal stabilizator, stabilizator kamere, stabilizator za telefon, dji osmo, zhiyun, moza, ručni stabilizator, 3-osni gimbal',
    'ru', 'гимбал стабилизатор, стабилизатор камеры, стабилизатор для смартфона, dji osmo, zhiyun, moza, 3-осевой гимбал'
) WHERE slug = 'stabilizatori-gimbal';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'streaming equipment, streaming gear, stream setup, elgato, capture card, streaming camera, studio equipment, twitch setup',
    'sr', 'streaming oprema, oprema za strimovanje, stream setup, elgato, capture card, streaming kamera, studijska oprema, twitch',
    'ru', 'стриминговое оборудование, оборудование для стрима, стрим сетап, elgato, карта захвата, стриминг камера, twitch'
) WHERE slug = 'streaming-oprema';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tablet, ipad, android tablet, samsung tablet, drawing tablet, tablet computer, 2-in-1 tablet, kids tablet',
    'sr', 'tablet, ipad, android tablet, samsung tablet, tabla za crtanje, tablet računar, 2-u-1 tablet, dečiji tablet',
    'ru', 'планшет, ipad, android планшет, samsung планшет, графический планшет, планшетный компьютер, 2-в-1 планшет'
) WHERE slug = 'tableti';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'television, tv, smart tv, 4k tv, oled tv, qled tv, led tv, home theater, uhd tv, android tv',
    'sr', 'televizor, tv, smart tv, 4k tv, oled tv, qled tv, led tv, kućni bioskop, uhd tv, android tv',
    'ru', 'телевизор, тв, smart tv, 4k телевизор, oled телевизор, qled телевизор, led телевизор, домашний кинотеатр'
) WHERE slug = 'tv-i-video';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'powered usb hub, usb hub with power adapter, multi-port usb hub, charging hub, usb 3.0 hub, anker hub, powered hub',
    'sr', 'usb hub sa napajanjem, napajani usb hub, multi-port usb hub, punjač hub, usb 3.0 hub, anker hub, powered hub',
    'ru', 'usb хаб с питанием, питаемый usb хаб, многопортовый usb хаб, зарядный хаб, usb 3.0 хаб, anker hub'
) WHERE slug = 'usb-hub-powered';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'webcam, web camera, hd webcam, 4k webcam, streaming webcam, logitech webcam, usb camera, video conferencing camera',
    'sr', 'web kamera, web cam, hd web kamera, 4k web kamera, streaming kamera, logitech kamera, usb kamera, video konferencijska kamera',
    'ru', 'веб камера, вебкамера, hd веб камера, 4k веб камера, стриминг камера, logitech камера, usb камера, видеоконференция'
) WHERE slug = 'web-kamere';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'mesh wifi router, mesh network, wifi 6 router, mesh system, whole home wifi, tp-link deco, google wifi, eero mesh',
    'sr', 'mesh wifi ruter, mesh mreža, wifi 6 ruter, mesh sistem, wifi za celu kuću, tp-link deco, google wifi, eero mesh',
    'ru', 'mesh wifi роутер, mesh сеть, wifi 6 роутер, mesh система, wifi для всего дома, tp-link deco, google wifi, eero'
) WHERE slug = 'wi-fi-routeri-mesh';

-- L3 Categories (detailed subcategories)

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', '4k video camera, 4k camcorder, ultra hd camera, professional video camera, 4k recording, cinema camera, 60fps 4k',
    'sr', '4k video kamera, 4k kamkorder, ultra hd kamera, profesionalna video kamera, 4k snimanje, cinema kamera, 60fps 4k',
    'ru', '4k видеокамера, 4k камкордер, ultra hd камера, профессиональная видеокамера, 4k запись, cinema камера, 60fps 4k'
) WHERE slug = '4k-video-kamera';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'dji action camera, dji osmo action, action cam, 4k action camera, waterproof action camera, dji action 3, dji action 4',
    'sr', 'dji action kamera, dji osmo action, action cam, 4k action kamera, vodootporna action kamera, dji action 3, dji action 4',
    'ru', 'dji экшн камера, dji osmo action, экшн кам, 4k экшн камера, водонепроницаемая экшн камера, dji action 3, dji action 4'
) WHERE slug = 'action-kamera-dji';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gopro, gopro hero, action camera, gopro hero 11, gopro hero 12, waterproof camera, sports camera, 4k gopro, gopro accessories',
    'sr', 'gopro, gopro hero, action kamera, gopro hero 11, gopro hero 12, vodootporna kamera, sportska kamera, 4k gopro, gopro dodaci',
    'ru', 'gopro, gopro hero, экшн камера, gopro hero 11, gopro hero 12, водонепроницаемая камера, спортивная камера, 4k gopro'
) WHERE slug = 'action-kamera-gopro';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'apple iphone, iphone 15, iphone 14, iphone 13, iphone pro, iphone pro max, ios phone, unlocked iphone, new iphone',
    'sr', 'apple iphone, iphone 15, iphone 14, iphone 13, iphone pro, iphone pro max, ios telefon, otključan iphone, novi iphone',
    'ru', 'apple iphone, айфон 15, айфон 14, айфон 13, iphone pro, iphone pro max, ios телефон, разблокированный iphone'
) WHERE slug = 'apple-iphone';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'apple watch series, apple watch 9, apple watch 8, apple watch se, watchos, apple smartwatch, fitness watch, apple watch bands',
    'sr', 'apple watch series, apple watch 9, apple watch 8, apple watch se, watchos, apple pametni sat, fitnes sat, apple watch narukvice',
    'ru', 'apple watch series, apple watch 9, apple watch 8, apple watch se, watchos, apple умные часы, фитнес часы'
) WHERE slug = 'apple-watch-series';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'apple watch ultra, apple watch ultra 2, outdoor smartwatch, titanium watch, diving watch, adventure watch, rugged smartwatch',
    'sr', 'apple watch ultra, apple watch ultra 2, outdoor pametni sat, titanijumski sat, ronilački sat, avanturistički sat, izdržljiv sat',
    'ru', 'apple watch ultra, apple watch ultra 2, outdoor умные часы, титановые часы, часы для дайвинга, adventure часы'
) WHERE slug = 'apple-watch-ultra';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'av receiver, home theater receiver, surround sound receiver, denon, yamaha, onkyo, marantz, 5.1 receiver, 7.2 receiver, dolby atmos',
    'sr', 'av risiver, kućni bioskop risiver, surround sound risiver, denon, yamaha, onkyo, marantz, 5.1 risiver, 7.2 risiver, dolby atmos',
    'ru', 'av ресивер, домашний кинотеатр ресивер, surround sound ресивер, denon, yamaha, onkyo, marantz, 5.1 ресивер, 7.2 ресивер'
) WHERE slug = 'av-resiveri';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'wireless headphones, bluetooth headphones, over-ear wireless, noise cancelling wireless, sony wh, bose quietcomfort, sennheiser wireless',
    'sr', 'bežične slušalice, bluetooth slušalice, preko-ušne bežične, noise cancelling bežične, sony wh, bose quietcomfort, sennheiser bežične',
    'ru', 'беспроводные наушники, bluetooth наушники, накладные беспроводные, шумоподавление беспроводные, sony wh, bose quietcomfort'
) WHERE slug = 'bežične-slušalice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'bluetooth speaker, portable bluetooth speaker, wireless speaker, jbl flip, jbl charge, bose soundlink, sony srs, waterproof speaker',
    'sr', 'bluetooth zvučnik, prenosivi bluetooth zvučnik, bežični zvučnik, jbl flip, jbl charge, bose soundlink, sony srs, vodootporni zvučnik',
    'ru', 'bluetooth колонка, портативная bluetooth колонка, беспроводная колонка, jbl flip, jbl charge, bose soundlink, sony srs'
) WHERE slug = 'bluetooth-zvucnici';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'drone with camera, camera drone, fpv drone, dji drone, aerial photography drone, quadcopter camera, 4k drone, gimbal drone',
    'sr', 'dron sa kamerom, kamera dron, fpv dron, dji dron, dron za aerofotografiju, quadcopter kamera, 4k dron, gimbal dron',
    'ru', 'дрон с камерой, камера дрон, fpv дрон, dji дрон, дрон для аэросъемки, квадрокоптер с камерой, 4k дрон, гимбал дрон'
) WHERE slug = 'dron-sa-kamerom';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'phone holder, car phone mount, phone stand, magnetic phone holder, dashboard mount, windshield mount, vent mount, phone grip',
    'sr', 'držač za telefon, auto držač za telefon, stalak za telefon, magnetni držač, nosač za tablu, nosač za vetrobran, ventilacioni držač',
    'ru', 'держатель для телефона, автомобильный держатель, подставка для телефона, магнитный держатель, крепление на панель'
) WHERE slug = 'držači-za-telefon';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'canon dslr, canon eos, canon 90d, canon 850d, canon rebel, canon camera, dslr camera canon, canon photography, canon aps-c',
    'sr', 'canon dslr, canon eos, canon 90d, canon 850d, canon rebel, canon foto aparat, dslr kamera canon, canon fotografija, canon aps-c',
    'ru', 'canon dslr, canon eos, canon 90d, canon 850d, canon rebel, canon фотоаппарат, dslr камера canon, canon фотография'
) WHERE slug = 'dslr-canon';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nikon dslr, nikon d, nikon d850, nikon d780, nikon d7500, nikon camera, dslr camera nikon, nikon photography, nikon fx',
    'sr', 'nikon dslr, nikon d, nikon d850, nikon d780, nikon d7500, nikon foto aparat, dslr kamera nikon, nikon fotografija, nikon fx',
    'ru', 'nikon dslr, nikon d, nikon d850, nikon d780, nikon d7500, nikon фотоаппарат, dslr камера nikon, nikon фотография'
) WHERE slug = 'dslr-nikon';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'photo lighting, studio lighting, softbox, ring light, led panel, continuous lighting, photography light, video light, godox lighting',
    'sr', 'foto rasveta, studijska rasveta, softbox, ring light, led panel, kontinuirana rasveta, svetlo za fotografiju, svetlo za video, godox',
    'ru', 'фото освещение, студийный свет, софтбокс, кольцевой свет, led панель, постоянный свет, освещение для фотографии'
) WHERE slug = 'foto-rasveta';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gaming monitor 144hz, 144hz monitor, high refresh rate monitor, gaming display, 1080p 144hz, 1440p 144hz, fast monitor, esports monitor',
    'sr', 'gaming monitor 144hz, 144hz monitor, monitor visoke osvežavajuće stope, gaming ekran, 1080p 144hz, 1440p 144hz, brzi monitor',
    'ru', 'игровой монитор 144hz, 144hz монитор, монитор с высокой частотой, игровой дисплей, 1080p 144hz, 1440p 144hz, быстрый монитор'
) WHERE slug = 'gaming-monitor-144hz';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gaming monitor 240hz, 240hz monitor, ultra high refresh rate, competitive gaming monitor, 1080p 240hz, 1440p 240hz, pro gaming monitor',
    'sr', 'gaming monitor 240hz, 240hz monitor, ultra visoka stopa osvežavanja, kompetitivni gaming monitor, 1080p 240hz, 1440p 240hz',
    'ru', 'игровой монитор 240hz, 240hz монитор, ультра высокая частота, конкурентный игровой монитор, 1080p 240hz, 1440p 240hz'
) WHERE slug = 'gaming-monitor-240hz';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gaming headset, gaming headphones, pc gaming headset, ps5 headset, xbox headset, surround sound headset, steelseries, razer, hyperx',
    'sr', 'gaming slušalice, gejmerske slušalice, pc gaming slušalice, ps5 slušalice, xbox slušalice, surround sound slušalice, steelseries, razer',
    'ru', 'игровые наушники, геймерские наушники, pc gaming наушники, ps5 наушники, xbox наушники, surround sound наушники, steelseries'
) WHERE slug = 'gaming-slusalice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'gaming chair, gamer chair, ergonomic gaming chair, racing chair, secretlab, dxracer, noblechairs, office gaming chair',
    'sr', 'gaming stolica, gejmerska stolica, ergonomska gaming stolica, trkaća stolica, secretlab, dxracer, noblechairs, kancelarijska gaming',
    'ru', 'игровое кресло, геймерское кресло, эргономичное игровое кресло, гоночное кресло, secretlab, dxracer, noblechairs'
) WHERE slug = 'gaming-stolice';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'garmin fitness watch, garmin forerunner, garmin venu, garmin vivoactive, gps watch, running watch, triathlon watch, garmin sports watch',
    'sr', 'garmin fitnes sat, garmin forerunner, garmin venu, garmin vivoactive, gps sat, sat za trčanje, triatlon sat, garmin sportski sat',
    'ru', 'garmin фитнес часы, garmin forerunner, garmin venu, garmin vivoactive, gps часы, часы для бега, триатлон часы'
) WHERE slug = 'garmin-fitness-sat';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ps5 controller, dualsense, playstation 5 controller, ps5 gamepad, wireless ps5 controller, dualsense edge, ps5 charging station',
    'sr', 'ps5 kontroler, dualsense, playstation 5 kontroler, ps5 gejmpad, bežični ps5 kontroler, dualsense edge, ps5 punjač stanica',
    'ru', 'ps5 контроллер, dualsense, playstation 5 контроллер, ps5 геймпад, беспроводной ps5 контроллер, dualsense edge'
) WHERE slug = 'gejmpad-ps5';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'xbox controller, xbox wireless controller, xbox gamepad, xbox series controller, xbox elite controller, xbox pc controller',
    'sr', 'xbox kontroler, xbox bežični kontroler, xbox gejmpad, xbox series kontroler, xbox elite kontroler, xbox pc kontroler',
    'ru', 'xbox контроллер, xbox беспроводной контроллер, xbox геймпад, xbox series контроллер, xbox elite контроллер'
) WHERE slug = 'gejmpad-xbox';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'google pixel, pixel 8, pixel 7, pixel 8 pro, pixel 7a, android phone google, pure android, google camera, unlocked pixel',
    'sr', 'google pixel, pixel 8, pixel 7, pixel 8 pro, pixel 7a, android telefon google, čist android, google kamera, otključan pixel',
    'ru', 'google pixel, pixel 8, pixel 7, pixel 8 pro, pixel 7a, android телефон google, чистый android, google камера'
) WHERE slug = 'google-pixel-telefoni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'amd rx 6000 graphics card, amd radeon rx 6000, rx 6600, rx 6700 xt, rx 6800 xt, rx 6900 xt, rdna 2, amd gpu',
    'sr', 'amd rx 6000 grafička kartica, amd radeon rx 6000, rx 6600, rx 6700 xt, rx 6800 xt, rx 6900 xt, rdna 2, amd gpu',
    'ru', 'amd rx 6000 видеокарта, amd radeon rx 6000, rx 6600, rx 6700 xt, rx 6800 xt, rx 6900 xt, rdna 2, amd gpu'
) WHERE slug = 'graficke-amd-rx-6000';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'amd rx 7000 graphics card, amd radeon rx 7000, rx 7600, rx 7700 xt, rx 7800 xt, rx 7900 xt, rx 7900 xtx, rdna 3',
    'sr', 'amd rx 7000 grafička kartica, amd radeon rx 7000, rx 7600, rx 7700 xt, rx 7800 xt, rx 7900 xt, rx 7900 xtx, rdna 3',
    'ru', 'amd rx 7000 видеокарта, amd radeon rx 7000, rx 7600, rx 7700 xt, rx 7800 xt, rx 7900 xt, rx 7900 xtx, rdna 3'
) WHERE slug = 'graficke-amd-rx-7000';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nvidia gtx 1000 series, gtx 1650, gtx 1660, gtx 1660 super, gtx 1660 ti, nvidia graphics card, budget gpu, pascal gpu',
    'sr', 'nvidia gtx 1000 serija, gtx 1650, gtx 1660, gtx 1660 super, gtx 1660 ti, nvidia grafička kartica, budžet gpu, pascal gpu',
    'ru', 'nvidia gtx 1000 серия, gtx 1650, gtx 1660, gtx 1660 super, gtx 1660 ti, nvidia видеокарта, бюджетная gpu, pascal gpu'
) WHERE slug = 'graficke-gtx-1000';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nvidia rtx 3000 series, rtx 3060, rtx 3060 ti, rtx 3070, rtx 3070 ti, rtx 3080, rtx 3090, ampere gpu, ray tracing, dlss',
    'sr', 'nvidia rtx 3000 serija, rtx 3060, rtx 3060 ti, rtx 3070, rtx 3070 ti, rtx 3080, rtx 3090, ampere gpu, ray tracing, dlss',
    'ru', 'nvidia rtx 3000 серия, rtx 3060, rtx 3060 ti, rtx 3070, rtx 3070 ti, rtx 3080, rtx 3090, ampere gpu, трассировка лучей'
) WHERE slug = 'graficke-rtx-3000';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nvidia rtx 4000 series, rtx 4060, rtx 4060 ti, rtx 4070, rtx 4070 ti, rtx 4080, rtx 4090, ada lovelace, dlss 3, ray tracing',
    'sr', 'nvidia rtx 4000 serija, rtx 4060, rtx 4060 ti, rtx 4070, rtx 4070 ti, rtx 4080, rtx 4090, ada lovelace, dlss 3, ray tracing',
    'ru', 'nvidia rtx 4000 серия, rtx 4060, rtx 4060 ti, rtx 4070, rtx 4070 ti, rtx 4080, rtx 4090, ada lovelace, dlss 3'
) WHERE slug = 'graficke-rtx-4000';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hi-fi system, high fidelity audio, stereo system, audiophile system, tube amplifier, turntable, hi-fi speakers, vinyl setup',
    'sr', 'hi-fi sistem, visoka vernost zvuka, stereo sistem, audiofilski sistem, lampsko pojačalo, gramofon, hi-fi zvučnici, vinil setup',
    'ru', 'hi-fi система, высококачественный звук, стерео система, аудиофильская система, ламповый усилитель, проигрыватель виниловых пластинок'
) WHERE slug = 'hi-fi-sistemi';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'huawei watch gt, huawei watch gt 3, huawei watch gt 4, huawei smartwatch, fitness smartwatch, long battery watch, harmonyos watch',
    'sr', 'huawei watch gt, huawei watch gt 3, huawei watch gt 4, huawei pametni sat, fitnes pametni sat, sat sa dugom baterijom, harmonyos',
    'ru', 'huawei watch gt, huawei watch gt 3, huawei watch gt 4, huawei умные часы, фитнес смарт часы, часы с долгой батареей'
) WHERE slug = 'huawei-watch-gt';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'home theater, home cinema, surround sound system, 5.1 system, 7.1 system, dolby atmos, theater speakers, subwoofer system',
    'sr', 'kućni bioskop, home cinema, surround sound sistem, 5.1 sistem, 7.1 sistem, dolby atmos, bioskopski zvučnici, subwoofer sistem',
    'ru', 'домашний кинотеатр, домашний cinema, система объемного звука, 5.1 система, 7.1 система, dolby atmos, кинотеатральные колонки'
) WHERE slug = 'kucni-bioskop';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', '2-in-1 laptop, convertible laptop, laptop tablet, touchscreen laptop, 360 laptop, detachable laptop, hybrid laptop, tablet pc',
    'sr', '2-u-1 laptop, konvertibilni laptop, laptop tablet, laptop sa ekranom na dodir, 360 laptop, odvojivi laptop, hibridni laptop',
    'ru', '2-в-1 ноутбук, трансформер ноутбук, ноутбук планшет, ноутбук с сенсорным экраном, 360 ноутбук, отсоединяемый ноутбук'
) WHERE slug = 'laptop-2-u-1';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'asus rog laptop, rog gaming laptop, republic of gamers, asus rog strix, rog zephyrus, gaming notebook asus, high performance laptop',
    'sr', 'asus rog laptop, rog gaming laptop, republic of gamers, asus rog strix, rog zephyrus, gaming notebook asus, visokih performansi',
    'ru', 'asus rog ноутбук, rog игровой ноутбук, republic of gamers, asus rog strix, rog zephyrus, высокопроизводительный ноутбук'
) WHERE slug = 'laptop-asus-rog';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'chromebook, chrome os laptop, google chromebook, cloud laptop, student laptop, web laptop, acer chromebook, asus chromebook',
    'sr', 'chromebook, chrome os laptop, google chromebook, cloud laptop, laptop za studente, web laptop, acer chromebook, asus chromebook',
    'ru', 'chromebook, chrome os ноутбук, google chromebook, облачный ноутбук, ноутбук для студентов, веб ноутбук'
) WHERE slug = 'laptop-chromebook';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'dell xps, dell xps 13, dell xps 15, dell xps 17, premium laptop, ultrabook dell, business laptop dell, infinity edge display',
    'sr', 'dell xps, dell xps 13, dell xps 15, dell xps 17, premium laptop, ultrabook dell, poslovni laptop dell, infinity edge ekran',
    'ru', 'dell xps, dell xps 13, dell xps 15, dell xps 17, премиум ноутбук, ультрабук dell, бизнес ноутбук dell'
) WHERE slug = 'laptop-dell-xps';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'hp pavilion, hp pavilion laptop, hp notebook, everyday laptop, multimedia laptop, student laptop hp, home laptop hp',
    'sr', 'hp pavilion, hp pavilion laptop, hp notebook, svakodnevni laptop, multimedijalni laptop, laptop za studente hp, kućni laptop hp',
    'ru', 'hp pavilion, hp pavilion ноутбук, hp notebook, повседневный ноутбук, мультимедийный ноутбук, студенческий ноутбук hp'
) WHERE slug = 'laptop-hp-pavilion';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'lenovo thinkpad, thinkpad laptop, business laptop, professional laptop, thinkpad x1, thinkpad t series, corporate laptop',
    'sr', 'lenovo thinkpad, thinkpad laptop, poslovni laptop, profesionalni laptop, thinkpad x1, thinkpad t serija, korporativni laptop',
    'ru', 'lenovo thinkpad, thinkpad ноутбук, бизнес ноутбук, профессиональный ноутбук, thinkpad x1, thinkpad t серия'
) WHERE slug = 'laptop-lenovo-thinkpad';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'macbook, macbook pro, macbook air, apple laptop, m1 macbook, m2 macbook, m3 macbook, macos laptop, retina display',
    'sr', 'macbook, macbook pro, macbook air, apple laptop, m1 macbook, m2 macbook, m3 macbook, macos laptop, retina ekran',
    'ru', 'macbook, macbook pro, macbook air, apple ноутбук, m1 macbook, m2 macbook, m3 macbook, macos ноутбук, retina дисплей'
) WHERE slug = 'laptop-macbook';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'led tv 32 inch, 32 inch tv, led tv 43 inch, small tv, bedroom tv, kitchen tv, compact tv, hd tv, full hd tv',
    'sr', 'led tv 32 inča, tv 32 inča, led tv 43 inča, mali tv, tv za spavaću sobu, tv za kuhinju, kompaktni tv, hd tv, full hd tv',
    'ru', 'led телевизор 32 дюйма, телевизор 32 дюйма, led телевизор 43 дюйма, маленький тв, тв для спальни, компактный тв'
) WHERE slug = 'led-tv-32-43-inca';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'led tv 50 inch, 50 inch tv, led tv 55 inch, 55 inch tv, living room tv, 4k tv, uhd tv, smart tv 55, mid size tv',
    'sr', 'led tv 50 inča, tv 50 inča, led tv 55 inča, tv 55 inča, tv za dnevnu sobu, 4k tv, uhd tv, smart tv 55, srednje veliki tv',
    'ru', 'led телевизор 50 дюймов, телевизор 50 дюймов, led телевизор 55 дюймов, телевизор для гостиной, 4k тв, uhd тв'
) WHERE slug = 'led-tv-50-55-inca';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'led tv 65 inch, 65 inch tv, led tv 75 inch, 75 inch tv, large tv, home theater tv, 4k large tv, 8k tv, cinema size tv',
    'sr', 'led tv 65 inča, tv 65 inča, led tv 75 inča, tv 75 inča, veliki tv, kućni bioskop tv, veliki 4k tv, 8k tv, bioskopska veličina',
    'ru', 'led телевизор 65 дюймов, телевизор 65 дюймов, led телевизор 75 дюймов, большой тв, домашний кинотеатр тв, 4k большой'
) WHERE slug = 'led-tv-65-75-inca';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'phone case, phone cover, protective case, smartphone case, silicone case, hard case, clear case, bumper case, designer case',
    'sr', 'maska za telefon, futrola za telefon, zaštitna maska, maska za pametni telefon, silikonska maska, tvrda maska, providna maska',
    'ru', 'чехол для телефона, защитный чехол, силиконовый чехол, жесткий чехол, прозрачный чехол, бампер чехол, дизайнерский чехол'
) WHERE slug = 'maske-za-telefone';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'amd x670 motherboard, x670e motherboard, am5 motherboard, ryzen 7000 motherboard, ddr5 motherboard, pcie 5.0 motherboard, asus x670',
    'sr', 'amd x670 matična ploča, x670e matična, am5 matična ploča, ryzen 7000 matična, ddr5 matična ploča, pcie 5.0 matična, asus x670',
    'ru', 'amd x670 материнская плата, x670e материнская, am5 материнская плата, ryzen 7000 материнская, ddr5 материнская плата'
) WHERE slug = 'maticna-ploca-amd-x670';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'intel z790 motherboard, lga 1700 motherboard, intel 13th gen motherboard, intel 14th gen, ddr5 motherboard, gaming motherboard, asus z790',
    'sr', 'intel z790 matična ploča, lga 1700 matična, intel 13. gen matična, intel 14. gen, ddr5 matična, gaming matična ploča, asus z790',
    'ru', 'intel z790 материнская плата, lga 1700 материнская, intel 13 поколения материнская, intel 14 поколения, ddr5 материнская'
) WHERE slug = 'maticna-ploca-intel-z790';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'fujifilm mirrorless, fujifilm x, fuji xt5, fuji xs20, fujifilm camera, aps-c mirrorless, retro camera, film simulation',
    'sr', 'fujifilm mirrorless, fujifilm x, fuji xt5, fuji xs20, fujifilm foto aparat, aps-c mirrorless, retro kamera, film simulacija',
    'ru', 'fujifilm беззеркальная, fujifilm x, fuji xt5, fuji xs20, fujifilm камера, aps-c беззеркальная, ретро камера'
) WHERE slug = 'mirrorless-fujifilm';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'sony mirrorless, sony alpha, sony a7, sony a7r, sony a7s, sony a9, full frame mirrorless, sony camera, e-mount camera',
    'sr', 'sony mirrorless, sony alpha, sony a7, sony a7r, sony a7s, sony a9, full frame mirrorless, sony kamera, e-mount kamera',
    'ru', 'sony беззеркальная, sony alpha, sony a7, sony a7r, sony a7s, sony a9, полнокадровая беззеркальная, sony камера'
) WHERE slug = 'mirrorless-sony';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nintendo switch games, switch games, mario, zelda, pokemon, animal crossing, splatoon, physical switch games, nintendo titles',
    'sr', 'nintendo switch igre, switch igre, mario, zelda, pokemon, animal crossing, splatoon, fizičke switch igre, nintendo naslovi',
    'ru', 'nintendo switch игры, switch игры, mario, zelda, pokemon, animal crossing, splatoon, физические switch игры'
) WHERE slug = 'nintendo-igre';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nintendo switch, switch oled, switch lite, nintendo console, portable console, handheld console, nintendo gaming',
    'sr', 'nintendo switch, switch oled, switch lite, nintendo konzola, prenosiva konzola, ručna konzola, nintendo igranje',
    'ru', 'nintendo switch, switch oled, switch lite, nintendo консоль, портативная консоль, карманная консоль'
) WHERE slug = 'nintendo-switch';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'canon lens, canon ef lens, canon rf lens, canon telephoto, canon wide angle, canon prime lens, canon zoom lens, l series lens',
    'sr', 'canon objektiv, canon ef objektiv, canon rf objektiv, canon telefoto, canon široki ugao, canon fiksni objektiv, canon zum, l serija',
    'ru', 'canon объектив, canon ef объектив, canon rf объектив, canon телефото, canon широкоугольный, canon фикс объектив'
) WHERE slug = 'objektiv-za-canon';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nikon lens, nikon f mount, nikon z mount, nikon telephoto, nikon wide angle, nikon prime lens, nikon zoom lens, nikkor lens',
    'sr', 'nikon objektiv, nikon f bajonet, nikon z bajonet, nikon telefoto, nikon široki ugao, nikon fiksni objektiv, nikon zum, nikkor objektiv',
    'ru', 'nikon объектив, nikon f байонет, nikon z байонет, nikon телефото, nikon широкоугольный, nikon фикс объектив'
) WHERE slug = 'objektiv-za-nikon';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'oled tv, oled television, lg oled, sony oled, self-lit oled, 4k oled, oled display, perfect black, infinite contrast',
    'sr', 'oled tv, oled televizor, lg oled, sony oled, samoosvetljavajući oled, 4k oled, oled ekran, savršena crna, beskonačan kontrast',
    'ru', 'oled телевизор, oled тв, lg oled, sony oled, самосветящийся oled, 4k oled, oled дисплей, идеальный черный'
) WHERE slug = 'oled-tv';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'oneplus phone, oneplus 11, oneplus 12, oneplus nord, oxygen os, flagship killer, fast charging phone, oneplus smartphone',
    'sr', 'oneplus telefon, oneplus 11, oneplus 12, oneplus nord, oxygen os, flagship ubica, brzo punjenje telefon, oneplus pametni telefon',
    'ru', 'oneplus телефон, oneplus 11, oneplus 12, oneplus nord, oxygen os, флагманский убийца, быстрая зарядка телефон'
) WHERE slug = 'oneplus-telefoni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'oppo phone, oppo find, oppo reno, coloros, fast charging oppo, oppo smartphone, camera phone oppo, 5g oppo phone',
    'sr', 'oppo telefon, oppo find, oppo reno, coloros, brzo punjenje oppo, oppo pametni telefon, kamera telefon oppo, 5g oppo telefon',
    'ru', 'oppo телефон, oppo find, oppo reno, coloros, быстрая зарядка oppo, oppo смартфон, камера телефон oppo, 5g oppo'
) WHERE slug = 'oppo-telefoni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'budget smartwatch, affordable smartwatch, cheap smartwatch, entry level smartwatch, fitness tracker watch, basic smartwatch',
    'sr', 'budžet pametni sat, pristupačan pametni sat, jeftin pametni sat, početni pametni sat, fitnes traker sat, osnovni pametni sat',
    'ru', 'бюджетные умные часы, доступные смарт часы, дешевые умные часы, начальные смарт часы, фитнес трекер часы'
) WHERE slug = 'pametni-satovi-budget';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'playstation 5, ps5 console, ps5 disc, ps5 digital, sony playstation 5, ps5 gaming, next gen console, dualsense controller',
    'sr', 'playstation 5, ps5 konzola, ps5 disk, ps5 digitalna, sony playstation 5, ps5 igranje, sledeća generacija konzola, dualsense kontroler',
    'ru', 'playstation 5, ps5 консоль, ps5 диск, ps5 цифровая, sony playstation 5, ps5 игры, консоль нового поколения'
) WHERE slug = 'playstation-5';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ps5 games, playstation 5 games, ps5 exclusive, god of war, spider-man, horizon, last of us, physical ps5 games, ps5 titles',
    'sr', 'ps5 igre, playstation 5 igre, ps5 ekskluzive, god of war, spider-man, horizon, last of us, fizičke ps5 igre, ps5 naslovi',
    'ru', 'ps5 игры, playstation 5 игры, ps5 эксклюзивы, god of war, spider-man, horizon, last of us, физические ps5 игры'
) WHERE slug = 'playstation-5-igre';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'used phone, refurbished phone, second hand phone, pre-owned phone, budget phone, cheap smartphone, renewed phone, unlocked used phone',
    'sr', 'polovan telefon, obnovljen telefon, telefon iz druge ruke, korišćen telefon, budžet telefon, jeftin pametni telefon, obnovljeni telefon',
    'ru', 'б/у телефон, восстановленный телефон, телефон с рук, подержанный телефон, бюджетный телефон, дешевый смартфон'
) WHERE slug = 'polovni-telefoni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'power bank, portable charger, battery pack, external battery, fast charging power bank, 20000mah power bank, usb-c power bank',
    'sr', 'power bank, prenosivi punjač, baterija paket, eksterna baterija, brzo punjenje power bank, 20000mah power bank, usb-c power bank',
    'ru', 'power bank, портативное зарядное устройство, батарейный пакет, внешняя батарея, быстрая зарядка power bank, 20000mah'
) WHERE slug = 'power-bank';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'amd ryzen 5, ryzen 5 processor, ryzen 5 5600, ryzen 5 7600, am4 ryzen 5, am5 ryzen 5, mid range cpu, gaming processor',
    'sr', 'amd ryzen 5, ryzen 5 procesor, ryzen 5 5600, ryzen 5 7600, am4 ryzen 5, am5 ryzen 5, srednji rang procesor, gaming procesor',
    'ru', 'amd ryzen 5, процессор ryzen 5, ryzen 5 5600, ryzen 5 7600, am4 ryzen 5, am5 ryzen 5, среднеуровневый процессор'
) WHERE slug = 'procesor-amd-ryzen-5';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'amd ryzen 7, ryzen 7 processor, ryzen 7 5800x, ryzen 7 7700x, ryzen 7 7800x3d, am4 ryzen 7, am5 ryzen 7, high performance cpu',
    'sr', 'amd ryzen 7, ryzen 7 procesor, ryzen 7 5800x, ryzen 7 7700x, ryzen 7 7800x3d, am4 ryzen 7, am5 ryzen 7, visokih performansi',
    'ru', 'amd ryzen 7, процессор ryzen 7, ryzen 7 5800x, ryzen 7 7700x, ryzen 7 7800x3d, am4 ryzen 7, am5 ryzen 7'
) WHERE slug = 'procesor-amd-ryzen-7';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'amd ryzen 9, ryzen 9 processor, ryzen 9 5900x, ryzen 9 7900x, ryzen 9 7950x, am4 ryzen 9, am5 ryzen 9, enthusiast cpu, flagship cpu',
    'sr', 'amd ryzen 9, ryzen 9 procesor, ryzen 9 5900x, ryzen 9 7900x, ryzen 9 7950x, am4 ryzen 9, am5 ryzen 9, entuzijast procesor',
    'ru', 'amd ryzen 9, процессор ryzen 9, ryzen 9 5900x, ryzen 9 7900x, ryzen 9 7950x, am4 ryzen 9, am5 ryzen 9, топовый процессор'
) WHERE slug = 'procesor-amd-ryzen-9';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'intel core i5, core i5 processor, i5 12400f, i5 13400f, i5 14400f, lga 1700, intel processor, mid range processor, gaming i5',
    'sr', 'intel core i5, core i5 procesor, i5 12400f, i5 13400f, i5 14400f, lga 1700, intel procesor, srednji rang procesor, gaming i5',
    'ru', 'intel core i5, процессор core i5, i5 12400f, i5 13400f, i5 14400f, lga 1700, intel процессор, среднеуровневый процессор'
) WHERE slug = 'procesor-intel-i5';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'intel core i7, core i7 processor, i7 12700k, i7 13700k, i7 14700k, lga 1700, high performance intel, gaming i7, overclocking i7',
    'sr', 'intel core i7, core i7 procesor, i7 12700k, i7 13700k, i7 14700k, lga 1700, visokih performansi intel, gaming i7, overclocking i7',
    'ru', 'intel core i7, процессор core i7, i7 12700k, i7 13700k, i7 14700k, lga 1700, высокопроизводительный intel, игровой i7'
) WHERE slug = 'procesor-intel-i7';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'intel core i9, core i9 processor, i9 12900k, i9 13900k, i9 14900k, lga 1700, flagship intel, enthusiast processor, extreme performance',
    'sr', 'intel core i9, core i9 procesor, i9 12900k, i9 13900k, i9 14900k, lga 1700, flagship intel, entuzijast procesor, ekstremne performanse',
    'ru', 'intel core i9, процессор core i9, i9 12900k, i9 13900k, i9 14900k, lga 1700, флагманский intel, топовый процессор'
) WHERE slug = 'procesor-intel-i9';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', '4k projector, ultra hd projector, 3840x2160 projector, home theater 4k, gaming projector 4k, laser projector, hdr projector',
    'sr', '4k projektor, ultra hd projektor, 3840x2160 projektor, kućni bioskop 4k, gaming projektor 4k, laserski projektor, hdr projektor',
    'ru', '4k проектор, ultra hd проектор, 3840x2160 проектор, домашний кинотеатр 4k, игровой проектор 4k, лазерный проектор'
) WHERE slug = 'projektori-4k';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'full hd projector, 1080p projector, 1920x1080 projector, budget projector, home projector, portable full hd, affordable projector',
    'sr', 'full hd projektor, 1080p projektor, 1920x1080 projektor, budžet projektor, kućni projektor, prenosivi full hd, pristupačan projektor',
    'ru', 'full hd проектор, 1080p проектор, 1920x1080 проектор, бюджетный проектор, домашний проектор, портативный full hd'
) WHERE slug = 'projektori-full-hd';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'phone charger, charging cable, usb-c cable, lightning cable, fast charger, wall charger, car charger, wireless charger cable',
    'sr', 'punjač za telefon, kabl za punjenje, usb-c kabl, lightning kabl, brzi punjač, zidni punjač, auto punjač, bežični punjač kabl',
    'ru', 'зарядное устройство для телефона, кабель для зарядки, usb-c кабель, lightning кабель, быстрое зарядное, настенное зарядное'
) WHERE slug = 'punjaci-i-kablovi-telefon';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'qled tv, quantum dot tv, samsung qled, 4k qled, qled television, hdr qled, bright tv, quantum processor',
    'sr', 'qled tv, quantum dot tv, samsung qled, 4k qled, qled televizor, hdr qled, svetli tv, quantum procesor',
    'ru', 'qled телевизор, quantum dot телевизор, samsung qled, 4k qled, qled тв, hdr qled, яркий телевизор, quantum процессор'
) WHERE slug = 'qled-tv';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ddr4 16gb ram, 16gb memory, ddr4 memory, desktop ram 16gb, laptop ram 16gb, gaming ram, dual channel 16gb, high speed ram',
    'sr', 'ddr4 16gb ram, 16gb memorija, ddr4 memorija, desktop ram 16gb, laptop ram 16gb, gaming ram, dual channel 16gb, brza ram',
    'ru', 'ddr4 16gb оперативная память, 16gb память, ddr4 память, desktop ram 16gb, ноутбук ram 16gb, игровая память'
) WHERE slug = 'ram-ddr4-16gb';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ddr4 8gb ram, 8gb memory, ddr4 memory, budget ram, entry level ram, desktop ram 8gb, laptop ram 8gb, affordable memory',
    'sr', 'ddr4 8gb ram, 8gb memorija, ddr4 memorija, budžet ram, početna ram, desktop ram 8gb, laptop ram 8gb, pristupačna memorija',
    'ru', 'ddr4 8gb оперативная память, 8gb память, ddr4 память, бюджетная ram, начальная память, desktop ram 8gb'
) WHERE slug = 'ram-ddr4-8gb';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ddr5 16gb ram, 16gb ddr5, next gen ram, ddr5 memory, high speed ddr5, gaming ddr5, desktop ddr5 16gb, laptop ddr5 16gb',
    'sr', 'ddr5 16gb ram, 16gb ddr5, sledeća generacija ram, ddr5 memorija, brza ddr5, gaming ddr5, desktop ddr5 16gb, laptop ddr5 16gb',
    'ru', 'ddr5 16gb оперативная память, 16gb ddr5, память нового поколения, ddr5 память, высокоскоростная ddr5, игровая ddr5'
) WHERE slug = 'ram-ddr5-16gb';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'ddr5 32gb ram, 32gb ddr5, high capacity ram, professional ram, content creation ram, workstation ram, dual channel 32gb',
    'sr', 'ddr5 32gb ram, 32gb ddr5, veliki kapacitet ram, profesionalna ram, kreatorska ram, radna stanica ram, dual channel 32gb',
    'ru', 'ddr5 32gb оперативная память, 32gb ddr5, высокая емкость ram, профессиональная память, память для рабочей станции'
) WHERE slug = 'ram-ddr5-32gb';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'realme phone, realme smartphone, realme gt, realme 11, realme c series, budget realme, fast charging realme, 5g realme',
    'sr', 'realme telefon, realme pametni telefon, realme gt, realme 11, realme c serija, budžet realme, brzo punjenje realme, 5g realme',
    'ru', 'realme телефон, realme смартфон, realme gt, realme 11, realme c серия, бюджетный realme, быстрая зарядка realme'
) WHERE slug = 'realme-telefoni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'rgb lighting, rgb led strip, addressable rgb, gaming rgb, pc rgb lighting, corsair rgb, razer chroma, rgb fans, rgb controller',
    'sr', 'rgb rasveta, rgb led traka, adresabilna rgb, gaming rgb, pc rgb rasveta, corsair rgb, razer chroma, rgb ventilatori, rgb kontroler',
    'ru', 'rgb подсветка, rgb led лента, адресуемая rgb, игровая rgb, pc rgb освещение, corsair rgb, razer chroma, rgb вентиляторы'
) WHERE slug = 'rgb-rasveta';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'samsung galaxy watch, galaxy watch 6, galaxy watch 5, wear os watch, samsung smartwatch, tizen watch, fitness smartwatch samsung',
    'sr', 'samsung galaxy watch, galaxy watch 6, galaxy watch 5, wear os sat, samsung pametni sat, tizen sat, fitnes pametni sat samsung',
    'ru', 'samsung galaxy watch, galaxy watch 6, galaxy watch 5, wear os часы, samsung умные часы, tizen часы, фитнес смарт часы'
) WHERE slug = 'samsung-galaxy-watch';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'samsung phone, samsung galaxy, galaxy s24, galaxy s23, galaxy a series, galaxy z fold, galaxy z flip, android samsung, 5g samsung',
    'sr', 'samsung telefon, samsung galaxy, galaxy s24, galaxy s23, galaxy a serija, galaxy z fold, galaxy z flip, android samsung, 5g samsung',
    'ru', 'samsung телефон, samsung galaxy, galaxy s24, galaxy s23, galaxy a серия, galaxy z fold, galaxy z flip, android samsung'
) WHERE slug = 'samsung-telefoni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'android tv, smart tv android, google tv, android television, chromecast built-in, google assistant tv, streaming tv android',
    'sr', 'android tv, smart tv android, google tv, android televizor, chromecast ugrađen, google assistant tv, streaming tv android',
    'ru', 'android tv, smart телевизор android, google tv, android телевизор, chromecast встроенный, google assistant тв'
) WHERE slug = 'smart-tv-android';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'tizen tv, samsung smart tv, tizen os, smart hub, samsung television, tizen operating system, samsung tv platform',
    'sr', 'tizen tv, samsung smart tv, tizen os, smart hub, samsung televizor, tizen operativni sistem, samsung tv platforma',
    'ru', 'tizen tv, samsung smart телевизор, tizen os, smart hub, samsung телевизор, tizen операционная система'
) WHERE slug = 'smart-tv-tizen';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'soundbar 2.1, 2.1 channel soundbar, soundbar with subwoofer, tv soundbar, home audio soundbar, wireless subwoofer soundbar',
    'sr', 'soundbar 2.1, 2.1 kanalni soundbar, soundbar sa subwooferom, tv soundbar, kućni audio soundbar, bežični subwoofer soundbar',
    'ru', 'саундбар 2.1, 2.1 канальный саундбар, саундбар с сабвуфером, тв саундбар, домашний аудио саундбар'
) WHERE slug = 'soundbar-2-1';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'soundbar 5.1, 5.1 channel soundbar, surround sound soundbar, home theater soundbar, wireless surround, rear speakers soundbar',
    'sr', 'soundbar 5.1, 5.1 kanalni soundbar, surround sound soundbar, kućni bioskop soundbar, bežični surround, zadnji zvučnici soundbar',
    'ru', 'саундбар 5.1, 5.1 канальный саундбар, объемный звук саундбар, домашний кинотеатр саундбар, беспроводной объемный'
) WHERE slug = 'soundbar-5-1';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nvme ssd 1tb, 1tb nvme, m.2 ssd 1tb, pcie ssd 1tb, fast ssd, samsung 980 pro, wd black sn850, gen4 ssd, storage 1tb',
    'sr', 'nvme ssd 1tb, 1tb nvme, m.2 ssd 1tb, pcie ssd 1tb, brzi ssd, samsung 980 pro, wd black sn850, gen4 ssd, skladište 1tb',
    'ru', 'nvme ssd 1tb, 1тб nvme, m.2 ssd 1тб, pcie ssd 1тб, быстрый ssd, samsung 980 pro, wd black sn850, gen4 ssd'
) WHERE slug = 'ssd-nvme-1tb';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'nvme ssd 500gb, 500gb nvme, m.2 ssd 500gb, pcie ssd 500gb, fast ssd, budget nvme, gaming ssd 500gb, gen3 ssd, storage 500gb',
    'sr', 'nvme ssd 500gb, 500gb nvme, m.2 ssd 500gb, pcie ssd 500gb, brzi ssd, budžet nvme, gaming ssd 500gb, gen3 ssd, skladište 500gb',
    'ru', 'nvme ssd 500gb, 500гб nvme, m.2 ssd 500гб, pcie ssd 500гб, быстрый ssd, бюджетный nvme, игровой ssd 500гб'
) WHERE slug = 'ssd-nvme-500gb';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'sata ssd 1tb, 1tb sata, 2.5 inch ssd, sata 3 ssd, budget ssd, laptop ssd 1tb, desktop ssd 1tb, reliable storage',
    'sr', 'sata ssd 1tb, 1tb sata, 2.5 inča ssd, sata 3 ssd, budžet ssd, laptop ssd 1tb, desktop ssd 1tb, pouzdano skladište',
    'ru', 'sata ssd 1tb, 1тб sata, 2.5 дюйма ssd, sata 3 ssd, бюджетный ssd, ноутбук ssd 1тб, desktop ssd 1тб'
) WHERE slug = 'ssd-sata-1tb';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'camera tripod, photo tripod, video tripod, travel tripod, monopod, ball head tripod, carbon fiber tripod, manfrotto, gitzo',
    'sr', 'stativ za kameru, foto stativ, video stativ, putni stativ, monopod, kuglasta glava stativ, karbon stativ, manfrotto, gitzo',
    'ru', 'штатив для камеры, фото штатив, видео штатив, путешествие штатив, монопод, шаровая головка штатив, карбоновый штатив'
) WHERE slug = 'stativ-foto-video';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'streaming microphone, usb streaming mic, condenser mic streaming, shure sm7b, elgato wave, rode podcaster, broadcast microphone',
    'sr', 'streaming mikrofon, usb streaming mikrofon, kondenzatorski mikrofon streaming, shure sm7b, elgato wave, rode podcaster, broadcast mikrofon',
    'ru', 'стриминг микрофон, usb стриминг микрофон, конденсаторный микрофон streaming, shure sm7b, elgato wave, rode podcaster'
) WHERE slug = 'streaming-mikrofon';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'budget phone, cheap smartphone, phone under 200 euro, affordable phone, entry level phone, starter phone, budget android, value phone',
    'sr', 'budžet telefon, jeftin pametni telefon, telefon do 200 evra, pristupačan telefon, početni telefon, starter telefon, budžet android',
    'ru', 'бюджетный телефон, дешевый смартфон, телефон до 200 евро, доступный телефон, начальный телефон, starter телефон'
) WHERE slug = 'telefoni-budget-do-200e';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'xbox games, xbox series x games, xbox series s games, halo, forza, gears of war, gamepass games, physical xbox games',
    'sr', 'xbox igre, xbox series x igre, xbox series s igre, halo, forza, gears of war, gamepass igre, fizičke xbox igre',
    'ru', 'xbox игры, xbox series x игры, xbox series s игры, halo, forza, gears of war, gamepass игры, физические xbox игры'
) WHERE slug = 'xbox-igre';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'xbox series s, xbox s console, digital xbox, budget xbox, gamepass console, 1440p gaming, microsoft xbox s, white xbox',
    'sr', 'xbox series s, xbox s konzola, digitalna xbox, budžet xbox, gamepass konzola, 1440p igranje, microsoft xbox s, bela xbox',
    'ru', 'xbox series s, xbox s консоль, цифровая xbox, бюджетная xbox, gamepass консоль, 1440p игры, microsoft xbox s'
) WHERE slug = 'xbox-series-s';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'xbox series x, xbox x console, next gen xbox, 4k gaming xbox, disc xbox, powerful console, flagship xbox, microsoft xbox x',
    'sr', 'xbox series x, xbox x konzola, sledeća generacija xbox, 4k igranje xbox, disk xbox, moćna konzola, flagship xbox, microsoft xbox x',
    'ru', 'xbox series x, xbox x консоль, следующее поколение xbox, 4k игры xbox, дисковая xbox, мощная консоль, флагманская xbox'
) WHERE slug = 'xbox-series-x';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'xiaomi mi band, mi band 8, mi band 7, xiaomi fitness tracker, budget fitness band, affordable tracker, heart rate band',
    'sr', 'xiaomi mi band, mi band 8, mi band 7, xiaomi fitnes traker, budžet fitnes narukvica, pristupačan traker, narukvica za otkucaje',
    'ru', 'xiaomi mi band, mi band 8, mi band 7, xiaomi фитнес трекер, бюджетный фитнес браслет, доступный трекер, пульсометр браслет'
) WHERE slug = 'xiaomi-mi-band';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'xiaomi phone, xiaomi smartphone, redmi, poco, mi phone, xiaomi 13, xiaomi 14, miui, budget xiaomi, 5g xiaomi phone',
    'sr', 'xiaomi telefon, xiaomi pametni telefon, redmi, poco, mi telefon, xiaomi 13, xiaomi 14, miui, budžet xiaomi, 5g xiaomi telefon',
    'ru', 'xiaomi телефон, xiaomi смартфон, redmi, poco, mi телефон, xiaomi 13, xiaomi 14, miui, бюджетный xiaomi, 5g xiaomi'
) WHERE slug = 'xiaomi-telefoni';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'xiaomi watch, xiaomi smartwatch, redmi watch, mi watch, wear os xiaomi, fitness watch xiaomi, budget smartwatch xiaomi',
    'sr', 'xiaomi sat, xiaomi pametni sat, redmi sat, mi sat, wear os xiaomi, fitnes sat xiaomi, budžet pametni sat xiaomi',
    'ru', 'xiaomi часы, xiaomi умные часы, redmi часы, mi часы, wear os xiaomi, фитнес часы xiaomi, бюджетные смарт часы xiaomi'
) WHERE slug = 'xiaomi-watch';

UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'screen protector, tempered glass, phone screen protector, 9h glass, protective glass, scratch resistant, oleophobic coating',
    'sr', 'zaštitno staklo, kaljeno staklo, zaštita ekrana telefona, 9h staklo, zaštitno staklo, otporno na ogrebotine, olofobni premaz',
    'ru', 'защитное стекло, закаленное стекло, защита экрана телефона, 9h стекло, защитное стекло, устойчивое к царапинам'
) WHERE slug = 'zastitno-staklo';

-- End of migration
