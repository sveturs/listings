#!/usr/bin/env python3
"""
Генератор миграций для category tree

Использование:
  python3 generate_categories.py --level 2 --count 81 --output 20251217100002_complete_l2_final.up.sql
  python3 generate_categories.py --level 3 --parent elektronika --count 90 --output 20251217100003_l3_elektronika.up.sql
"""

import argparse
import json
import os
import sys
from typing import List, Dict

# Список СУЩЕСТВУЮЩИХ L2 slugs (301 + 18 = 319)
EXISTING_L2_SLUGS = set("""
akumulatori, akumulatori-profesionalni, akvarijumi, alati-i-popravke, alati-za-automobile, alati-za-basta, anti-aging,
arhiviranje, audio-i-navigacija, audio-knjige, audio-oprema, audio-oprema-za-muziku, auto-dodaci, auto-kozmetika,
bas-gitare, bastenska-garnitura, bastenska-oprema, bastenska-rasveta, bastenske-ukrase, bazeni-i-spa,
bele-table-i-prezentacije, beležnice-blokovi, bez-laktoze, bezglutenska-hrana, bicikli-i-trotineti, biografije,
biser, biserni-nakit, boje-akrilne, boje-u-olovci, boje-uljane, boje-vodene, bojleri, bravar, brusilice-ugaone,
brusilice-vibracione, bubnjevi-setovi, busilice-profesionalne, caj-premium, casopisi, cetkice-za-slikanje, cirkular,
dash-kamere, decija-kozmetika, decija-obuca, decija-obuca-bebe, decija-odeca, decija-odeca-bebe, decije-knjige,
deciji-namestaj, decoupage-setovi, dekoracije, dekorativna-kozmetika, delovi-za-automobile, delovi-za-motocikle,
depilacija, desktop-racunari, dijamanti, dj-oprema, dodaci-i-aksesoari, dodatna-oprema-elektronika, donji-ves,
dronovi, duvacki-instrumenti, dzonovanje, e-citaci, e-knjige, efekti-i-procesori, elegantna-odeca, elektricar,
elektricne-gitare, elektronika-za-decu, esarpe-i-salovi, eseji, eterična-ulja, fantastika, filmovi,
fitnes-i-teretana, foto-i-video-kamere, fotograf-event, fotograf-portret, friteze, frizideri, fudbal,
gaming-oprema, gelovi-za-tusiranje, generatori-benzin, generatori-dizel, gips, gitare-akusticne,
glina-za-modelovanje, graficki-dizajn, grncari ja-i-biljke, gume-i-felne, higijena, horor, hrana-bez-secera,
hrana-za-bebe, hrana-za-glodare, hrana-za-macke, hrana-za-pse, hrana-za-ptice, igracke-za-bebe, igracke-za-decu,
igracke-za-macke, igracke-za-pse, intimna-higijena, istorijske-knjige, kafa-premium, kalkul atori,
kalkulator-i-oprema, kampovanje, kancelarijska-hartija, kancelarijski-nameštaj, kancelarijski-pribor, keramicar,
keto-dijeta, klavijature-i-pianos, knjige-beletristika, knjige-za-roditelje, kompostiranje, kompresori-klipni,
kompresori-ventilatorski, konac-za-vez, konzole-i-gaming, konzolne-igre, korepetitor-engleski, korepetitor-matematika,
korporativni-pokloni, kosarka, kosulјe-kratkih-rukava, kreveti-za-ljubimce, kriminalisticki-romani, krovopokrivac,
kuhinjski-pribor, kupaci-kostimi, kupatilo, laminating, laptop-racunari, lepota-aparati, ljubavni-romani,
lomači-dokumenta, lov, luksuzna-kozmetika, luksuzni-satovi, majstor-general, makeup-cetkice, mali-kucni-aparati,
manikir-i-pedikir, masine-za-kafu, masine-za-pranje, maslinovo-ulje-premium, med-organski, medicinski-proizvodi,
meraci-laserni, merila-digitalna, mikotalasne-rerne, mikrofoni, mikrofoni-za-muziku, mlecni-proizvodi-organik,
moderni-nakit, molberti, moler-farbanje, montaza-video, moto-oprema, mreza-i-internet, muska-nega, muska-obuca,
muska-odeca, muskarci-veliki-brojevi, muski-satovi, muski-stil, muzicki-dodaci, muzika, namestaj-dnevna-soba,
namestaj-kancelarija, namestaj-kuhinja, namestaj-spavaca-soba, namestaj-za-bebe, naocari-i-dodaci, nas-i-storage,
nega-i-higijena-beba, nega-kose, nega-koze, niveliri, nosilje-za-ljubimce, odela-i-smokingzi, odvijaci-akumulatorski,
ogledala, olovke-i-markeri, olovke-za-crtanje, oprema-za-bebe, oprema-za-setnju, oralna-higijena,
organizacija-i-skladistenje, organizacija-stola, organska-kozmetika, organski-sokovi, organsko-povrce, organsko-voce,
orgulje-i-harmonijum, paleo-dijeta, pametni-satovi, pametni-telefoni, parfemi, parketar, parking-senzori, periferija,
perle-za-nakit, planinarenje, platna-za-slikanje, platno-za-vez, plivanje, plus-size-odeca, poezija,
posteljina-i-peskiri, precistaci-vazduha, pregradni-zidovi, prepis-teksta, prevod-jezik, projektna-oprema,
projektori, protein-pločice, proteini-u-prahu, racunarske-komponente, radna-odeca, rasveta, registratori-fascikle,
rende-elektricne, retke-knjige, ribolov, sat i-za-zid, satovski-dodaci, servis-bele-tehnike, servis-laptopa,
servis-telefona, servis-tv, sf-knjige, skener-i-kopir, skeneri, skolski-pribor, smart-home, smart-narukvice,
snimanje-i-produkcija, snimanje-video, spa-i-relaksacija, sporet-i-rerna, sportska-ishrana-dodaci, sportska-odeca,
srebrni-nakit, stampaci-toneri, stege, stolar, stripovi, strucne-knjige, sudopere-i-masine, superfood, tableti,
tekstil-za-dom, tenis, tepihi-i-prostirke, terarijumi, tesle, testenine-integralne, testeri-kruzne, testeri-lancane,
toaleta-za-ljubimce, torbice-i-novcanici, trudnicka-odeca, tuniranje, tv-i-video, udarni-odvijaci, ukulele-i-mandoline,
usisivaci, vaze-i-dekor, veganska-hrana, vegetarijanska-hrana, vencana-odeca, ventilacija-i-klimatizacija,
ventilatori-i-grejalice, verenicko-prstenje, ves-masine-dodaci, violina-i-zicani, vitamini-i-suplementi,
vodoinstalater, vodopasi, web-dizajn, web-kamere, zastita-od-sunca, zavarivaci-inverter, zavarivaci-mig,
zavarivaci-tig, zdrave-grickalice, zene-veliki-brojevi, zenska-obuca, zenska-odeca, zenski-satovi, zimska-garderoba,
zimski-sportovi, zitarice-organik, zlatni-nakit,
bluetooth-zvucnici-premium, bluetooth-slusalice-profesionalne, wi-fi-routeri-mesh, powerline-adapteri,
usb-hub-powered, eksterni-hard-diskovi-ssd, memorijske-kartice-pro, citaci-e-ink, elektronske-knjige-citalice,
bluetooth-tastature, graficke-tablice, 3d-stampaci, streaming-oprema, pametne-sijalice, punjaci-bezicni,
stabilizatori-gimbal, mikroskopi-digitalni, solarni-punjaci
""".replace('\n', ' ').split(', '))

# Cleanup: remove empty strings
EXISTING_L2_SLUGS = {s.strip() for s in EXISTING_L2_SLUGS if s.strip()}

# L1 категории и их slugs
L1_CATEGORIES = {
    'elektronika': 'Elektronika',
    'odeca-i-obuca': 'Odeća i obuća',
    'dom-i-basta': 'Dom i bašta',
    'lepota-i-zdravlje': 'Lepota i zdravlje',
    'za-bebe-i-decu': 'Za bebe i decu',
    'sport-i-turizam': 'Sport i turizam',
    'automobilizam': 'Automobilizam',
    'kucni-aparati': 'Kućni aparati',
    'nakit-i-satovi': 'Nakit i satovi',
    'knjige-i-mediji': 'Knjige i mediji',
    'kucni-ljubimci': 'Kućni ljubimci',
    'kancelarijski-materijal': 'Kancelarijski materijal',
    'hrana-i-pice': 'Hrana i piće',
    'umetnost-i-rukotvorine': 'Umetnost i rukotvorine',
    'muzicki-instrumenti': 'Muzički instrumenti',
    'industrija-i-alati': 'Industrija i alati',
    'usluge': 'Usluge',
    'ostalo': 'Ostalo'
}


# Новые L2 категории для разных L1
NEW_L2_CATEGORIES = {
    'odeca-i-obuca': [
        ('kosmeticke-uniforme', 'Kozmetičke uniforme', 'Cosmetic Uniforms', 'Косметическая униформа', '👗'),
        ('pidzhame-svila', 'Pidžame od svile', 'Silk Pajamas', 'Шелковые пижамы', '💤'),
        ('kucne-halјine', 'Kućne haljine', 'House Dresses', 'Домашние платья', '🏡'),
        ('bade-mantili', 'Bade mantili', 'Bathrobes', 'Халаты', '🛁'),
        ('negl ize-i-spavacice', 'Negliže i spavaćice', 'Negligees and Nightgowns', 'Неглиже и ночные рубашки', '😴'),
        ('odeca-trudnice-svecana', 'Svečana odeća za trudnice', 'Formal Maternity Wear', 'Вечерняя одежда для беременных', '🤰'),
        ('muske-velike-velicine', 'Muške velike veličine', 'Men''s Plus Size', 'Мужские большие размеры', '👔'),
        ('zenske-velike-velicine', 'Ženske velike veličine', 'Women''s Plus Size', 'Женские большие размеры', '👗'),
        ('termo-ves', 'Termo veš', 'Thermal Underwear', 'Термобелье', '🥶'),
        ('odeca-ronjenje', 'Odeća za ronjenje', 'Diving Wear', 'Одежда для дайвинга', '🤿'),
    ],
    'dom-i-basta': [
        ('kuhinjski-set-nozeva', 'Kuhinjski set noževa', 'Kitchen Knife Sets', 'Кухонные наборы ножей', '🔪'),
        ('drveno-posude-kuhinja', 'Drveno posuđe za kuhinju', 'Wooden Kitchenware', 'Деревянная посуда', '🥄'),
        ('keramika-posude-rucno', 'Keramičko posuđe ručno rađeno', 'Handmade Ceramic Dishes', 'Керамическая посуда ручной работы', '🏺'),
        ('staklo-kristal-posude', 'Staklo i kristal posuđe', 'Glass and Crystal Dishes', 'Стеклянная и хрустальная посуда', '🍷'),
        ('tacne-posluzavnici', 'Tacne i poslužavnici', 'Trays and Serving Platters', 'Подносы и сервировочные блюда', '🍽️'),
        ('kozarci-case-luksuz', 'Kozarci i čaše luksuzni', 'Luxury Glasses', 'Роскошные бокалы', '🥂'),
        ('kafe-servisi', 'Servisi za kafu', 'Coffee Sets', 'Кофейные сервизы', '☕'),
        ('caj-servisi', 'Servisi za čaj', 'Tea Sets', 'Чайные сервизы', '🍵'),
        ('zamrzivaci-vertikalni', 'Zamrzivači vertikalni', 'Upright Freezers', 'Вертикальные морозильники', '🧊'),
        ('vinske-vitrine', 'Vinske vitrine', 'Wine Coolers', 'Винные холодильники', '🍷'),
        ('bar-oprema-dom', 'Bar oprema za dom', 'Home Bar Equipment', 'Барное оборудование для дома', '🍸'),
    ],
    'sport-i-turizam': [
        ('joga-blokovi', 'Joga blokovi', 'Yoga Blocks', 'Блоки для йоги', '🧘'),
        ('pilates-lopte', 'Pilates lopte', 'Pilates Balls', 'Мячи для пилатеса', '⚽'),
        ('trake-istezanje', 'Trake za istezanje', 'Resistance Bands', 'Эластичные ленты', '🏋️'),
        ('tenis-reketi-pro', 'Tenis reketi profesionalni', 'Professional Tennis Rackets', 'Профессиональные теннисные ракетки', '🎾'),
        ('badminton-oprema', 'Badminton oprema', 'Badminton Equipment', 'Оборудование для бадминтона', '🏸'),
        ('stoni-tenis-reketi', 'Stoni tenis reketi', 'Table Tennis Rackets', 'Ракетки для настольного тенниса', '🏓'),
        ('borilacke-vestine-zastita', 'Borilačke veštine zaštita', 'Martial Arts Protective Gear', 'Защита для боевых искусств', '🥋'),
        ('ronilacka-oprema-profesionalna', 'Ronilačka oprema profesionalna', 'Professional Diving Equipment', 'Профессиональное оборудование для дайвинга', '🤿'),
        ('surfovanje-daske', 'Surfovanje daske', 'Surfboards', 'Доски для серфинга', '🏄'),
        ('kajak-kanu-oprema', 'Kajak i kanu oprema', 'Kayak and Canoe Equipment', 'Оборудование для каяка и каноэ', '🛶'),
    ],
    'automobilizam': [
        ('auto-sedista-bebe', 'Auto sedišta za bebe', 'Baby Car Seats', 'Детские автокресла', '👶'),
        ('auto-organizatori', 'Auto organizatori', 'Car Organizers', 'Автомобильные органайзеры', '📦'),
        ('ambijentalno-osvetljenje', 'Ambijentalno osvetljenje', 'Ambient Lighting', 'Ambient освещение', '💡'),
        ('drzaci-telefona-auto', 'Držači telefona za auto', 'Car Phone Holders', 'Держатели телефонов для авто', '📱'),
        ('auto-pokrivaci', 'Auto pokrivači', 'Car Covers', 'Автомобильные чехлы', '🚗'),
        ('klime-uređaji-prenosni', 'Klima uređaji prenosni', 'Portable Air Conditioners', 'Портативные кондиционеры', '❄️'),
        ('gps-trekeri-auto', 'GPS trekeri za auto', 'GPS Trackers for Cars', 'GPS трекеры для авто', '📍'),
        ('punjaci-elektricni-auto', 'Punjači za električne auto', 'EV Chargers', 'Зарядки для электромобилей', '🔌'),
        ('auto-aspiratori', 'Auto aspiratori', 'Car Vacuum Cleaners', 'Автомобильные пылесосы', '💨'),
        ('led-trake-auto', 'LED trake za auto', 'LED Strips for Cars', 'LED ленты для авто', '💡'),
    ],
    'kucni-aparati': [
        ('indukcionе-ploce', 'Indukcione ploče', 'Induction Cooktops', 'Индукционные плиты', '🔥'),
        ('parne-rerne', 'Parne rerne', 'Steam Ovens', 'Паровые духовки', '♨️'),
        ('aspiratori-bez-kesice', 'Aspiratori bez kesice', 'Bagless Vacuum Cleaners', 'Пылесосы без мешков', '🌪️'),
        ('roboti-usisivaci', 'Roboti usisivači', 'Robot Vacuum Cleaners', 'Роботы-пылесосы', '🤖'),
        ('masine-pranje-susenje', 'Mašine za pranje i sušenje', 'Washer-Dryer Combos', 'Стирально-сушильные машины', '🧺'),
        ('masine-sudjе', 'Mašine za suđe', 'Dishwashers', 'Посудомоечные машины', '🍽️'),
        ('friteze-sa-vrucim-vazduhom', 'Friteze sa vrućim vazduhom', 'Air Fryers', 'Аэрофритюрницы', '🍟'),
        ('instant-lonci', 'Instant lonci', 'Instant Pots', 'Мультиварки Instant Pot', '🍲'),
        ('multi-kukeri', 'Multi kukeri', 'Multi-Cookers', 'Мультиварки', '👨\u200d🍳'),
        ('blenderi-profesionalni', 'Blenderi profesionalni', 'Professional Blenders', 'Профессиональные блендеры', '🥤'),
    ],
    'lepota-i-zdravlje': [
        ('led-maske-lice', 'LED maske za lice', 'LED Face Masks', 'LED маски для лица', '💡'),
        ('dermaroller-mikronedling', 'Dermaroller mikroneedling', 'Dermaroller Microneedling', 'Дермароллер микронидлинг', '💉'),
        ('ultrazvucni-cistaci', 'Ultrazvučni čistači', 'Ultrasonic Cleansers', 'Ультразвуковые очистители', '🔊'),
        ('ionske-cetkice-kosa', 'Jonske četkice za kosu', 'Ionic Hair Brushes', 'Ионные щетки для волос', '💇'),
        ('laseri-uklanjanje-dlaka', 'Laseri za uklanjanje dlaka', 'Laser Hair Removal', 'Лазеры для удаления волос', '✨'),
        ('fene-profesionalni', 'Fēnovi profesionalni', 'Professional Hair Dryers', 'Профессиональные фены', '💨'),
        ('pegla-kosa-turmalin', 'Pegle za kosu turmalin', 'Tourmaline Hair Straighteners', 'Турмалиновые утюжки для волос', '🦱'),
        ('biotin-kompleksi', 'Biotin kompleksi', 'Biotin Complexes', 'Биотиновые комплексы', '💊'),
        ('kolagen-preparati', 'Kolagen preparati', 'Collagen Supplements', 'Коллагеновые добавки', '✨'),
        ('terapije-svetlom', 'Terapije svetlom', 'Light Therapy', 'Светотерапия', '💡'),
    ],
    'za-bebe-i-decu': [
        ('bebi-monitori-video', 'Bebi monitori video', 'Video Baby Monitors', 'Видеоняни', '📹'),
        ('nosiljke-ergonomske', 'Nosiljke ergonomske', 'Ergonomic Baby Carriers', 'Эргономичные переноски', '👶'),
        ('kolevke-automatske', 'Kolevke automatske', 'Automatic Cradles', 'Автоматические колыбели', '😴'),
        ('sterilizatori', 'Sterilizatori', 'Sterilizers', 'Стерилизаторы', '🧼'),
        ('grejaci-flasica', 'Grejači flašica', 'Bottle Warmers', 'Подогреватели бутылочек', '🍼'),
        ('kada-bebe-sklop ljiva', 'Kade za bebe sklopljive', 'Foldable Baby Bathtubs', 'Складные детские ванночки', '🛁'),
        ('nocnici', 'Noćnici', 'Potty Chairs', 'Детские горшки', '🚽'),
        ('zaštitne-ograde', 'Zaštitne ograde', 'Safety Gates', 'Защитные ограждения', '🚪'),
        ('bebi-jumperoo', 'Bebi jumperoo', 'Baby Jumpers', 'Детские прыгунки', '🦘'),
        ('edukativni-tepisi', 'Edukativni tepisi', 'Educational Play Mats', 'Развивающие коврики', '🧸'),
    ],
    'nakit-i-satovi': [
        ('prstenje-veridba-dijamant', 'Prstenje za veridbu dijamant', 'Diamond Engagement Rings', 'Помолвочные кольца с бриллиантами', '💍'),
        ('luksuzni-satovi-swiss', 'Luksuzni satovi švajcarski', 'Swiss Luxury Watches', 'Швейцарские люкс часы', '⌚'),
        ('pametni-satovi-fitness', 'Pametni satovi fitness', 'Fitness Smartwatches', 'Фитнес-часы', '⌚'),
        ('perlice-narukvice', 'Perlice narukvice', 'Pearl Bracelets', 'Жемчужные браслеты', '📿'),
        ('zlatne-ogrlice', 'Zlatne ogrlice', 'Gold Necklaces', 'Золотые ожерелья', '✨'),
        ('srebrne-minđuše', 'Srebrne minđuše', 'Silver Earrings', 'Серебряные серьги', '💎'),
        ('broše-vintage', 'Broševi vintage', 'Vintage Brooches', 'Винтажные броши', '🌹'),
        ('privesci-lanci', 'Privesci i lanci', 'Pendants and Chains', 'Кулоны и цепочки', '⛓️'),
    ],
    'knjige-i-mediji': [
        ('autobiografije', 'Autobiografije', 'Autobiographies', 'Автобиографии', '📚'),
        ('kuvari-recepti', 'Kuvari i recepti', 'Cookbooks and Recipes', 'Кулинарные книги', '👨\u200d🍳'),
        ('putopisi', 'Putopisi', 'Travel Books', 'Путеводители', '✈️'),
        ('enciklopedije', 'Enciklopedije', 'Encyclopedias', 'Энциклопедии', '📕'),
        ('knjige-za-decu-slikovnice', 'Dečije slikovnice', 'Children''s Picture Books', 'Детские книжки с картинками', '📖'),
        ('grafički-romani', 'Grafički romani', 'Graphic Novels', 'Графические романы', '📚'),
        ('audio-knjige-cd', 'Audio knjige CD', 'Audio Books CD', 'Аудиокниги на CD', '💿'),
        ('blu-ray-diskovi', 'Blu-ray diskovi', 'Blu-ray Discs', 'Blu-ray диски', '💿'),
        ('vinil-ploce', 'Vinil ploče', 'Vinyl Records', 'Виниловые пластинки', '🎵'),
    ],
    'kucni-ljubimci': [
        ('automatski-hranjivaci', 'Automatski hranjivači', 'Automatic Feeders', 'Автоматические кормушки', '🍽️'),
        ('fontane-voda-ljubimci', 'Fontane za vodu', 'Pet Water Fountains', 'Фонтанчики для воды', '💧'),
        ('ogrlice-gps-trekeri', 'Ogrlice GPS trekeri', 'GPS Tracker Collars', 'GPS ошейники', '📍'),
        ('kucice-ljubimce', 'Kućice za ljubimce', 'Pet Houses', 'Домики для питомцев', '🏠'),
        ('torbe-nosiljke-ljubimci', 'Torbe i nosiljke', 'Pet Carriers and Bags', 'Переноски для питомцев', '👜'),
        ('oprema-dresura', 'Oprema za dresuru', 'Training Equipment', 'Оборудование для дрессировки', '🎓'),
        ('kostimi-ljubimci', 'Kostimi za ljubimce', 'Pet Costumes', 'Костюмы для питомцев', '🎭'),
    ],
    'hrana-i-pice': [
        ('domace-slatko', 'Domaće slatko', 'Homemade Jam', 'Домашнее варенье', '🍓'),
        ('domaci-med', 'Domaći med', 'Local Honey', 'Местный мёд', '🍯'),
        ('domace-rakije', 'Domaće rakije', 'Homemade Spirits', 'Домашние настойки', '🍷'),
        ('sirevi-lokalni', 'Sirevi lokalni', 'Local Cheeses', 'Местные сыры', '🧀'),
        ('organske-zitarice', 'Organske žitarice', 'Organic Grains', 'Органические злаки', '🌾'),
        ('glutenfree-testenine', 'Gluten-free testenine', 'Gluten-Free Pasta', 'Безглютеновые макароны', '🍝'),
        ('cajevi-biljni', 'Čajevi biljni', 'Herbal Teas', 'Травяные чаи', '🍵'),
        ('protein-barovi', 'Protein barovi', 'Protein Bars', 'Протеиновые батончики', '🍫'),
        ('ovsene-pahuljice', 'Ovsene pahuljice', 'Oat Flakes', 'Овсяные хлопья', '🥣'),
    ],
}


def escape_sql(text: str) -> str:
    """Escape single quotes for SQL"""
    return text.replace("'", "''")


def generate_l2_migration(count: int, output_file: str) -> None:
    """Generate L2 migration file"""

    sql_content = []
    sql_content.append(f"-- Migration: Complete L2 Categories - добавить {count} новых L2 категорий")
    sql_content.append(f"-- Date: 2025-12-17")
    sql_content.append(f"-- Target: Увеличить L2 до 400")
    sql_content.append("")
    sql_content.append("INSERT INTO categories (slug, name, description, meta_title, meta_description, parent_id, level, path, sort_order, icon, is_active)")
    sql_content.append("VALUES")

    values = []
    sort_order = 300

    for parent_slug, categories in NEW_L2_CATEGORIES.items():
        for cat in categories:
            if len(values) >= count:
                break

            slug, name_sr, name_en, name_ru, icon = cat

            # Skip if already exists
            if slug in EXISTING_L2_SLUGS:
                print(f"⚠️ Skipping duplicate: {slug}")
                continue

            sort_order += 1

            name_sr_escaped = escape_sql(name_sr)
            name_en_escaped = escape_sql(name_en)
            name_ru_escaped = escape_sql(name_ru)

            value = f"""  (
    '{slug}',
    '{{"sr": "{name_sr_escaped}", "en": "{name_en_escaped}", "ru": "{name_ru_escaped}"}}',
    '{{"sr": "{name_sr_escaped} - širok asortiman", "en": "{name_en_escaped} - wide range", "ru": "{name_ru_escaped} - широкий ассортимент"}}',
    '{{"sr": "{name_sr_escaped} | Vondi", "en": "{name_en_escaped} | Vondi", "ru": "{name_ru_escaped} | Vondi"}}',
    '{{"sr": "Pronađite najbolje {name_sr_escaped} na Vondi marketplace", "en": "Find the best {name_en_escaped} on Vondi marketplace", "ru": "Найдите лучшие {name_ru_escaped} на Vondi marketplace"}}',
    (SELECT id FROM categories WHERE slug = '{parent_slug}' AND level = 1),
    2,
    '{parent_slug}/{slug}',
    {sort_order},
    '{icon}',
    true
  )"""
            values.append(value)

        if len(values) >= count:
            break

    sql_content.append(",\n".join(values))
    sql_content.append(";")
    sql_content.append("")
    sql_content.append("-- Verification: проверка на дубликаты")
    sql_content.append("""DO $$
DECLARE
  duplicate_count INT;
BEGIN
  SELECT COUNT(*) INTO duplicate_count FROM (
    SELECT slug, COUNT(*) FROM categories GROUP BY slug HAVING COUNT(*) > 1
  ) dup;
  IF duplicate_count > 0 THEN
    RAISE EXCEPTION 'Found % duplicate slugs!', duplicate_count;
  END IF;
  RAISE NOTICE 'No duplicates found ✅';
  RAISE NOTICE 'Successfully added new L2 categories';
END $$;""")

    with open(output_file, 'w', encoding='utf-8') as f:
        f.write('\n'.join(sql_content))

    print(f"✅ Generated {output_file} with {len(values)} categories")


def generate_down_migration(slugs: List[str], output_file: str) -> None:
    """Generate DOWN migration"""
    down_content = []
    down_content.append("-- Rollback: Delete newly added L2 categories")
    down_content.append("")
    down_content.append("DELETE FROM categories WHERE slug IN (")
    down_content.append(",\n".join([f"  '{slug}'" for slug in slugs]))
    down_content.append(") AND level = 2;")

    down_file = output_file.replace('.up.sql', '.down.sql')
    with open(down_file, 'w', encoding='utf-8') as f:
        f.write('\n'.join(down_content))

    print(f"✅ Generated {down_file}")


def main():
    parser = argparse.ArgumentParser(description='Generate category migrations')
    parser.add_argument('--level', type=int, choices=[2, 3], required=True, help='Category level')
    parser.add_argument('--count', type=int, required=True, help='Number of categories to generate')
    parser.add_argument('--output', type=str, required=True, help='Output SQL file name')
    parser.add_argument('--parent', type=str, help='Parent category slug (for L3)')

    args = parser.parse_args()

    if args.level == 2:
        generate_l2_migration(args.count, args.output)
    else:
        print("L3 generation not implemented yet")
        sys.exit(1)


if __name__ == '__main__':
    main()
