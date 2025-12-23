#!/usr/bin/env python3
"""
Генератор L3 категорий для Vondi Marketplace

Использование:
  python3 generate_l3_categories.py --output-dir migrations/
"""

import argparse
import json
from typing import List, Tuple, Dict

# L3 категории для различных L2 родителей
L3_CATEGORIES = {
    # ELEKTRONIKA L2 parents
    'pametni-telefoni': [
        ('samsung-galaxy-s', 'Samsung Galaxy S serija', 'Samsung Galaxy S Series', 'Samsung Galaxy S серия', '📱'),
        ('samsung-galaxy-a', 'Samsung Galaxy A serija', 'Samsung Galaxy A Series', 'Samsung Galaxy A серия', '📱'),
        ('apple-iphone-14-15', 'Apple iPhone 14-15', 'Apple iPhone 14-15', 'Apple iPhone 14-15', '📱'),
        ('apple-iphone-11-12-13', 'Apple iPhone 11-12-13', 'Apple iPhone 11-12-13', 'Apple iPhone 11-12-13', '📱'),
        ('xiaomi-redmi', 'Xiaomi Redmi', 'Xiaomi Redmi', 'Xiaomi Redmi', '📱'),
        ('xiaomi-poco', 'Xiaomi POCO', 'Xiaomi POCO', 'Xiaomi POCO', '📱'),
        ('huawei-p-serija', 'Huawei P serija', 'Huawei P Series', 'Huawei P серия', '📱'),
        ('oneplus-telefoni', 'OnePlus telefoni', 'OnePlus Phones', 'OnePlus телефоны', '📱'),
        ('google-pixel', 'Google Pixel', 'Google Pixel', 'Google Pixel', '📱'),
        ('telefoni-do-150eur', 'Telefoni do 150 EUR', 'Phones under 150 EUR', 'Телефоны до 150 EUR', '💰'),
        ('telefoni-150-300eur', 'Telefoni 150-300 EUR', 'Phones 150-300 EUR', 'Телефоны 150-300 EUR', '💰'),
        ('telefoni-preko-300eur', 'Telefoni preko 300 EUR', 'Phones over 300 EUR', 'Телефоны свыше 300 EUR', '💰'),
        ('maske-samsung', 'Maske za Samsung', 'Samsung Cases', 'Чехлы для Samsung', '📱'),
        ('maske-iphone', 'Maske za iPhone', 'iPhone Cases', 'Чехлы для iPhone', '📱'),
        ('zastitno-staklo', 'Zaštitno staklo', 'Screen Protectors', 'Защитные стекла', '🛡️'),
    ],
    'laptop-racunari': [
        ('laptop-gejming', 'Gejming laptopovi', 'Gaming Laptops', 'Игровые ноутбуки', '🎮'),
        ('laptop-poslovni', 'Poslovni laptopovi', 'Business Laptops', 'Бизнес ноутбуки', '💼'),
        ('ultrabook', 'Ultrabukovi', 'Ultrabooks', 'Ультрабуки', '💻'),
        ('macbook', 'MacBook', 'MacBook', 'MacBook', '🍎'),
        ('laptop-15-inch', 'Laptopovi 15"', 'Laptops 15"', 'Ноутбуки 15"', '💻'),
        ('laptop-17-inch', 'Laptopovi 17"', 'Laptops 17"', 'Ноутбуки 17"', '💻'),
        ('laptop-do-500eur', 'Laptopovi do 500 EUR', 'Laptops under 500 EUR', 'Ноутбуки до 500 EUR', '💰'),
        ('laptop-500-1000eur', 'Laptopovi 500-1000 EUR', 'Laptops 500-1000 EUR', 'Ноутбуки 500-1000 EUR', '💰'),
        ('laptop-preko-1000eur', 'Laptopovi preko 1000 EUR', 'Laptops over 1000 EUR', 'Ноутбуки свыше 1000 EUR', '💰'),
        ('laptop-dodaci', 'Dodaci za laptop', 'Laptop Accessories', 'Аксессуары для ноутбуков', '💼'),
        ('torbe-laptop', 'Torbe za laptop', 'Laptop Bags', 'Сумки для ноутбуков', '👜'),
        ('hladnjaci-laptop', 'Hladnjaci za laptop', 'Laptop Cooling Pads', 'Охлаждение для ноутбуков', '❄️'),
    ],
    'tv-i-video': [
        ('smart-tv-55', 'Smart TV 55"', 'Smart TV 55"', 'Smart TV 55"', '📺'),
        ('smart-tv-65', 'Smart TV 65"', 'Smart TV 65"', 'Smart TV 65"', '📺'),
        ('smart-tv-75', 'Smart TV 75"', 'Smart TV 75"', 'Smart TV 75"', '📺'),
        ('oled-tv', 'OLED TV', 'OLED TV', 'OLED TV', '📺'),
        ('qled-tv', 'QLED TV', 'QLED TV', 'QLED TV', '📺'),
        ('4k-tv', '4K Ultra HD TV', '4K Ultra HD TV', '4K Ultra HD TV', '📺'),
        ('8k-tv', '8K TV', '8K TV', '8K TV', '📺'),
        ('projektori-4k', 'Projektori 4K', '4K Projectors', 'Проекторы 4K', '🎬'),
        ('soundbar', 'Soundbar sistemi', 'Soundbar Systems', 'Саундбары', '🔊'),
        ('home-cinema', 'Home cinema sistemi', 'Home Cinema Systems', 'Домашний кинотеатр', '🎭'),
        ('streaming-box', 'Streaming box uređaji', 'Streaming Boxes', 'Стриминг-приставки', '📦'),
    ],
    'audio-oprema': [
        ('bluetooth-zvucnici', 'Bluetooth zvučnici', 'Bluetooth Speakers', 'Bluetooth колонки', '🔊'),
        ('slusalice-wireless', 'Wireless slušalice', 'Wireless Headphones', 'Беспроводные наушники', '🎧'),
        ('slusalice-anc', 'Slušalice sa ANC', 'ANC Headphones', 'Наушники с ANC', '🎧'),
        ('slusalice-gaming', 'Gaming slušalice', 'Gaming Headsets', 'Игровые наушники', '🎮'),
        ('slusalice-in-ear', 'In-ear slušalice', 'In-ear Headphones', 'Внутриканальные наушники', '🎧'),
        ('zvucnici-polica', 'Zvučnici za policu', 'Bookshelf Speakers', 'Полочные колонки', '🔊'),
        ('karaoke-sistemi', 'Karaoke sistemi', 'Karaoke Systems', 'Караоке системы', '🎤'),
        ('cd-plejeri', 'CD plejeri', 'CD Players', 'CD плееры', '💿'),
        ('vinilploce-plejeri', 'Plejeri za vinilploče', 'Turntables', 'Проигрыватели винила', '🎵'),
    ],
    'foto-i-video-kamere': [
        ('dslr-kamere', 'DSLR kamere', 'DSLR Cameras', 'DSLR камеры', '📷'),
        ('mirrorless-kamere', 'Mirrorless kamere', 'Mirrorless Cameras', 'Беззеркальные камеры', '📷'),
        ('akcione-kamere', 'Akcione kamere', 'Action Cameras', 'Экшн-камеры', '🎬'),
        ('gopro', 'GoPro', 'GoPro', 'GoPro', '📹'),
        ('video-kamere-4k', 'Video kamere 4K', '4K Video Cameras', 'Видеокамеры 4K', '🎥'),
        ('objektivi-canon', 'Objektivi za Canon', 'Canon Lenses', 'Объективы Canon', '🔍'),
        ('objektivi-nikon', 'Objektivi za Nikon', 'Nikon Lenses', 'Объективы Nikon', '🔍'),
        ('objektivi-sony', 'Objektivi za Sony', 'Sony Lenses', 'Объективы Sony', '🔍'),
        ('tronošci', 'Tronošci', 'Tripods', 'Штативы', '📷'),
        ('foto-blicevi', 'Foto blicevi', 'Camera Flashes', 'Вспышки для камеры', '💡'),
    ],
    'pametni-satovi': [
        ('apple-watch', 'Apple Watch', 'Apple Watch', 'Apple Watch', '⌚'),
        ('samsung-galaxy-watch', 'Samsung Galaxy Watch', 'Samsung Galaxy Watch', 'Samsung Galaxy Watch', '⌚'),
        ('xiaomi-mi-watch', 'Xiaomi Mi Watch', 'Xiaomi Mi Watch', 'Xiaomi Mi Watch', '⌚'),
        ('garmin-satovi', 'Garmin satovi', 'Garmin Watches', 'Часы Garmin', '⌚'),
        ('huawei-watch', 'Huawei Watch', 'Huawei Watch', 'Huawei Watch', '⌚'),
    ],
    'gaming-oprema': [
        ('gaming-tastature', 'Gaming tastature', 'Gaming Keyboards', 'Игровые клавиатуры', '⌨️'),
        ('gaming-misevi', 'Gaming miševi', 'Gaming Mice', 'Игровые мыши', '🖱️'),
        ('gaming-slusalice', 'Gaming slušalice', 'Gaming Headsets', 'Игровые наушники', '🎧'),
        ('gaming-stolice', 'Gaming stolice', 'Gaming Chairs', 'Игровые кресла', '🪑'),
        ('gaming-monitori', 'Gaming monitori', 'Gaming Monitors', 'Игровые мониторы', '🖥️'),
        ('racing-volani', 'Racing volani', 'Racing Wheels', 'Рулевые системы', '🎮'),
        ('joystick', 'Joystick-ovi', 'Joysticks', 'Джойстики', '🕹️'),
        ('vr-naocari', 'VR naočari', 'VR Headsets', 'VR гарнитуры', '🥽'),
    ],

    # ODECA-I-OBUCA L2 parents
    'muska-odeca': [
        ('majice-polo', 'Polo majice', 'Polo Shirts', 'Поло футболки', '👕'),
        ('majice-kratkih-rukava', 'Majice kratkih rukava', 'Short Sleeve Shirts', 'Футболки с коротким рукавом', '👕'),
        ('kosulje-dugi-rukav', 'Košulje dugi rukav', 'Long Sleeve Shirts', 'Рубашки с длинным рукавом', '👔'),
        ('kosulje-kratki-rukav', 'Košulje kratki rukav', 'Short Sleeve Shirts', 'Рубашки с коротким рукавом', '👔'),
        ('farmerke-slim', 'Farmerke slim', 'Slim Jeans', 'Джинсы slim', '👖'),
        ('farmerke-regular', 'Farmerke regular', 'Regular Jeans', 'Джинсы regular', '👖'),
        ('pantalone-elegantne', 'Elegantne pantalone', 'Dress Pants', 'Классические брюки', '👔'),
        ('trenerke', 'Trenerke', 'Track Suits', 'Спортивные костюмы', '🏃'),
        ('dzemperi-vuneni', 'Džemperi vuneni', 'Wool Sweaters', 'Шерстяные свитеры', '🧶'),
        ('jakne-kožne', 'Kožne jakne', 'Leather Jackets', 'Кожаные куртки', '🧥'),
        ('perjane-jakne', 'Perjane jakne', 'Down Jackets', 'Пуховики', '🧥'),
        ('odela-muska', 'Muška odela', 'Men''s Suits', 'Мужские костюмы', '🤵'),
        ('smokingzi', 'Smokingzi', 'Tuxedos', 'Смокинги', '🎩'),
        ('sako-muski', 'Muški sako', 'Men''s Blazers', 'Мужские пиджаки', '🧥'),
        ('sorc-muski', 'Muški šortsevi', 'Men''s Shorts', 'Мужские шорты', '🩳'),
    ],
    'zenska-odeca': [
        ('haljine-svecane', 'Svečane haljine', 'Evening Dresses', 'Вечерние платья', '👗'),
        ('haljine-letnje', 'Letnje haljine', 'Summer Dresses', 'Летние платья', '👗'),
        ('haljine-kokteli', 'Kokteil haljine', 'Cocktail Dresses', 'Коктейльные платья', '👗'),
        ('bluze-svilene', 'Svilene bluze', 'Silk Blouses', 'Шелковые блузы', '👚'),
        ('majice-zenske', 'Ženske majice', 'Women''s T-shirts', 'Женские футболки', '👕'),
        ('farmerke-zenske-skinny', 'Ženske farmerke skinny', 'Women''s Skinny Jeans', 'Женские джинсы скинни', '👖'),
        ('farmerke-zenske-mom', 'Ženske farmerke mom', 'Women''s Mom Jeans', 'Женские джинсы mom', '👖'),
        ('suknje-mini', 'Mini suknje', 'Mini Skirts', 'Мини юбки', '👗'),
        ('suknje-midi', 'Midi suknje', 'Midi Skirts', 'Миди юбки', '👗'),
        ('suknje-duge', 'Duge suknje', 'Long Skirts', 'Длинные юбки', '👗'),
        ('pantalone-zenske', 'Ženske pantalone', 'Women''s Pants', 'Женские брюки', '👖'),
        ('zenske-jakne', 'Ženske jakne', 'Women''s Jackets', 'Женские куртки', '🧥'),
        ('kardigani', 'Kardigani', 'Cardigans', 'Кардиганы', '🧶'),
        ('blejzeri-zenske', 'Ženski blejzeri', 'Women''s Blazers', 'Женские пиджаки', '🧥'),
        ('kombinezoni', 'Kombinezoni', 'Jumpsuits', 'Комбинезоны', '👗'),
    ],
    'decija-odeca': [
        ('odeca-bebe-0-2', 'Odeća za bebe 0-2 god', 'Baby Clothing 0-2 years', 'Одежда для малышей 0-2 года', '👶'),
        ('odeca-devojcice-3-6', 'Odeća za devojčice 3-6 god', 'Girls Clothing 3-6 years', 'Одежда для девочек 3-6 лет', '👧'),
        ('odeca-devojcice-7-12', 'Odeća za devojčice 7-12 god', 'Girls Clothing 7-12 years', 'Одежда для девочек 7-12 лет', '👧'),
        ('odeca-decaci-3-6', 'Odeća za dečake 3-6 god', 'Boys Clothing 3-6 years', 'Одежда для мальчиков 3-6 лет', '👦'),
        ('odeca-decaci-7-12', 'Odeća za dečake 7-12 god', 'Boys Clothing 7-12 years', 'Одежда для мальчиков 7-12 лет', '👦'),
        ('skolske-uniforme', 'Školske uniforme', 'School Uniforms', 'Школьная форма', '🎒'),
        ('decije-trenerke', 'Dečije trenerke', 'Kids Track Suits', 'Детские спортивные костюмы', '🏃'),
        ('decije-jakne', 'Dečije jakne', 'Kids Jackets', 'Детские куртки', '🧥'),
        ('decije-kupace-kostimi', 'Dečiji kupaći kostimi', 'Kids Swimwear', 'Детские купальники', '🏊'),
    ],
    'muska-obuca': [
        ('cipele-kožne', 'Kožne cipele', 'Leather Shoes', 'Кожаные туфли', '👞'),
        ('cipele-sportske', 'Sportske cipele', 'Sports Shoes', 'Спортивная обувь', '👟'),
        ('patike-running', 'Patike za trčanje', 'Running Shoes', 'Беговые кроссовки', '👟'),
        ('patike-lifestyle', 'Lifestyle patike', 'Lifestyle Sneakers', 'Lifestyle кроссовки', '👟'),
        ('cizme-muske', 'Muške čizme', 'Men''s Boots', 'Мужские ботинки', '🥾'),
        ('sandale-muske', 'Muške sandale', 'Men''s Sandals', 'Мужские сандалии', '👡'),
        ('papuce-muske', 'Muške papuče', 'Men''s Slippers', 'Мужские тапочки', '🩴'),
    ],
    'zenska-obuca': [
        ('cipele-stikla', 'Cipele na štiklu', 'High Heels', 'Туфли на каблуке', '👠'),
        ('cipele-ravne', 'Ravne cipele', 'Flat Shoes', 'Балетки', '🥿'),
        ('baletanke', 'Baletanke', 'Ballet Flats', 'Балетки', '🥿'),
        ('patike-zenske', 'Ženske patike', 'Women''s Sneakers', 'Женские кроссовки', '👟'),
        ('cizme-zenske', 'Ženske čizme', 'Women''s Boots', 'Женские сапоги', '👢'),
        ('cizme-iznad-kolena', 'Čizme iznad kolena', 'Over-the-Knee Boots', 'Ботфорты', '👢'),
        ('gležnjace', 'Gležnjače', 'Ankle Boots', 'Ботильоны', '👢'),
        ('sandale-zenske', 'Ženske sandale', 'Women''s Sandals', 'Женские сандалии', '👡'),
        ('espadrile', 'Espadrile', 'Espadrilles', 'Эспадрильи', '👟'),
    ],
    'decija-obuca': [
        ('patike-decije-0-2', 'Dečije patike 0-2 god', 'Baby Shoes 0-2 years', 'Детская обувь 0-2 года', '👶'),
        ('patike-decije-3-6', 'Dečije patike 3-6 god', 'Kids Shoes 3-6 years', 'Детская обувь 3-6 лет', '👧'),
        ('patike-decije-7-12', 'Dečije patike 7-12 god', 'Kids Shoes 7-12 years', 'Детская обувь 7-12 лет', '👦'),
        ('sandale-decije', 'Dečije sandale', 'Kids Sandals', 'Детские сандалии', '👡'),
        ('cizme-decije', 'Dečije čizme', 'Kids Boots', 'Детские ботинки', '🥾'),
        ('papuce-decije', 'Dečije papuče', 'Kids Slippers', 'Детские тапочки', '🩴'),
    ],

    # DOM-I-BASTA L2 parents
    'namestaj-dnevna-soba': [
        ('trosed-kauci', 'Trosed kauči', 'Three-Seater Sofas', 'Трехместные диваны', '🛋️'),
        ('dvosed-kauci', 'Dvosed kauči', 'Two-Seater Sofas', 'Двухместные диваны', '🛋️'),
        ('ugaone-garniture', 'Ugaone garniture', 'Corner Sofas', 'Угловые диваны', '🛋️'),
        ('fotelje', 'Fotelje', 'Armchairs', 'Кресла', '🪑'),
        ('klupa-za-hodnik', 'Klupe za hodnik', 'Hallway Benches', 'Скамейки для прихожей', '🪑'),
        ('tv-komode', 'TV komode', 'TV Stands', 'ТВ тумбы', '📺'),
        ('police-za-knjige', 'Police za knjige', 'Bookshelves', 'Книжные полки', '📚'),
        ('vitrini', 'Vitrini', 'Display Cabinets', 'Витрины', '🏛️'),
        ('stolici-za-trpezariju', 'Stolice za trpezariju', 'Dining Chairs', 'Стулья для столовой', '🪑'),
        ('trpezarijski-stolovi', 'Trpezarijski stolovi', 'Dining Tables', 'Обеденные столы', '🍽️'),
    ],
    'namestaj-spavaca-soba': [
        ('bracni-kreveti', 'Bračni kreveti', 'Double Beds', 'Двуспальные кровати', '🛏️'),
        ('boks-kreveti', 'Box kreveti', 'Box Spring Beds', 'Кровати с боксом', '🛏️'),
        ('kreveti-jednostruki', 'Jednostruki kreveti', 'Single Beds', 'Односпальные кровати', '🛏️'),
        ('ormani-garderoberi', 'Ormani garderoberi', 'Wardrobes', 'Шкафы', '🚪'),
        ('nocni-ormancici', 'Noćni ormarići', 'Nightstands', 'Прикроватные тумбочки', '🛋️'),
        ('toaletni-stolovi', 'Toaletni stolovi', 'Vanity Tables', 'Туалетные столики', '💄'),
    ],
    'kupatilo': [
        ('lavaboi', 'Lavaboi', 'Sinks', 'Раковины', '🚰'),
        ('kadе-kupatilo', 'Kade', 'Bathtubs', 'Ванны', '🛁'),
        ('tus-kabine', 'Tuš kabine', 'Shower Cabins', 'Душевые кабины', '🚿'),
        ('toalet-skoljke', 'Toalet školjke', 'Toilets', 'Унитазы', '🚽'),
        ('ogledala-kupatilo', 'Ogledala za kupatilo', 'Bathroom Mirrors', 'Зеркала для ванной', '🪞'),
        ('ormarici-kupatilo', 'Ormarići za kupatilo', 'Bathroom Cabinets', 'Шкафы для ванной', '🚪'),
        ('slavine', 'Slavine', 'Faucets', 'Смесители', '🚰'),
    ],
    'rasveta': [
        ('plafonijere-led', 'LED plafonijere', 'LED Ceiling Lights', 'LED потолочные светильники', '💡'),
        ('lusteri', 'Lusteri', 'Chandeliers', 'Люстры', '✨'),
        ('podne-lampe-moderne', 'Moderne podne lampe', 'Modern Floor Lamps', 'Современные напольные лампы', '🕯️'),
        ('stone-lampe-nocne', 'Noćne stone lampe', 'Night Table Lamps', 'Ночные настольные лампы', '🛋️'),
        ('led-trake', 'LED trake', 'LED Strips', 'LED ленты', '💡'),
        ('spoljna-rasveta-led', 'LED spoljna rasveta', 'LED Outdoor Lighting', 'LED уличное освещение', '💡'),
    ],

    # SPORT-I-TURIZAM L2 parents
    'fitnes-i-teretana': [
        ('tegovi-bučice', 'Tegovi i bučice', 'Weights and Dumbbells', 'Гантели и гири', '🏋️'),
        ('tegovi-tegovi-disk', 'Tegovi disk', 'Weight Plates', 'Диски для штанги', '⚖️'),
        ('bencevi', 'Bencevi za vežbanje', 'Weight Benches', 'Скамьи для жима', '🛋️'),
        ('tred-milin', 'Traka za trčanje', 'Treadmills', 'Беговые дорожки', '🏃'),
        ('bicikl-sobni', 'Sobni bicikl', 'Exercise Bikes', 'Велотренажеры', '🚴'),
        ('veslacki-masine', 'Veslački mašine', 'Rowing Machines', 'Гребные тренажеры', '🚣'),
        ('elipticne-masine', 'Eliptične mašine', 'Elliptical Machines', 'Эллиптические тренажеры', '🏃'),
        ('joga-strunjace', 'Joga strunjače', 'Yoga Mats', 'Йога коврики', '🧘'),
        ('pilates-lopte', 'Pilates lopte', 'Pilates Balls', 'Пилатес мячи', '⚽'),
    ],
    'bicikli-i-trotineti': [
        ('mtb-bicikli', 'MTB bicikli', 'Mountain Bikes', 'Горные велосипеды', '🚵'),
        ('drumski-bicikli', 'Drumski bicikli', 'Road Bikes', 'Шоссейные велосипеды', '🚴'),
        ('gradski-bicikli', 'Gradski bicikli', 'City Bikes', 'Городские велосипеды', '🚲'),
        ('bmx-bicikli', 'BMX bicikli', 'BMX Bikes', 'BMX велосипеды', '🚴'),
        ('deciji-bicikli', 'Dečiji bicikli', 'Kids Bikes', 'Детские велосипеды', '🚲'),
        ('elektricni-bicikli', 'Električni bicikli', 'Electric Bikes', 'Электровелосипеды', '⚡'),
        ('elektricni-trotineti', 'Električni trotineti', 'Electric Scooters', 'Электросамокаты', '🛴'),
        ('trotineti-deca', 'Trotineti za decu', 'Kids Scooters', 'Детские самокаты', '🛴'),
        ('bicikl-delovi', 'Delovi za bicikl', 'Bicycle Parts', 'Запчасти для велосипедов', '🔧'),
    ],
    'kampovanje': [
        ('satori-2-osobe', 'Šatori za 2 osobe', '2-Person Tents', 'Палатки на 2 человека', '⛺'),
        ('satori-4-osobe', 'Šatori za 4 osobe', '4-Person Tents', 'Палатки на 4 человека', '⛺'),
        ('satori-porodicni', 'Porodični šatori', 'Family Tents', 'Семейные палатки', '⛺'),
        ('vrece-spavanja', 'Vreće za spavanje', 'Sleeping Bags', 'Спальные мешки', '😴'),
        ('prostirke-kamp', 'Kamp prostirke', 'Camping Mats', 'Туристические коврики', '🏕️'),
        ('kamp-stolice', 'Kamp stolice', 'Camping Chairs', 'Кемпинговые стулья', '🪑'),
        ('kamp-stolovi', 'Kamp stolovi', 'Camping Tables', 'Кемпинговые столы', '🍽️'),
        ('dzepne-lampe', 'Džepne lampe', 'Flashlights', 'Фонарики', '🔦'),
        ('prenosne-rostilji', 'Prenosni roštilji', 'Portable Grills', 'Портативные грили', '🔥'),
    ],
}


def escape_sql(text: str) -> str:
    """Escape single quotes for SQL"""
    return text.replace("'", "''")


def generate_l3_migration_file(parent_l2_slugs: List[str], output_file: str, start_sort_order: int = 1) -> int:
    """Generate L3 migration for specified L2 parents"""

    sql_content = []
    sql_content.append(f"-- Migration: L3 Categories")
    sql_content.append(f"-- Date: 2025-12-17")
    sql_content.append("")
    sql_content.append("INSERT INTO categories (slug, name, description, meta_title, meta_description, parent_id, level, path, sort_order, icon, is_active)")
    sql_content.append("VALUES")

    values = []
    sort_order = start_sort_order
    added_slugs = []

    for parent_slug in parent_l2_slugs:
        if parent_slug not in L3_CATEGORIES:
            print(f"⚠️ No L3 data for parent: {parent_slug}")
            continue

        for cat in L3_CATEGORIES[parent_slug]:
            slug, name_sr, name_en, name_ru, icon = cat

            sort_order += 1

            name_sr_escaped = escape_sql(name_sr)
            name_en_escaped = escape_sql(name_en)
            name_ru_escaped = escape_sql(name_ru)

            value = f"""  (
    '{slug}',
    '{{"sr": "{name_sr_escaped}", "en": "{name_en_escaped}", "ru": "{name_ru_escaped}"}}',
    '{{"sr": "{name_sr_escaped} - najbolja ponuda", "en": "{name_en_escaped} - best selection", "ru": "{name_ru_escaped} - лучший выбор"}}',
    '{{"sr": "{name_sr_escaped} | Vondi", "en": "{name_en_escaped} | Vondi", "ru": "{name_ru_escaped} | Vondi"}}',
    '{{"sr": "Kupite {name_sr_escaped} na Vondi", "en": "Buy {name_en_escaped} on Vondi", "ru": "Купите {name_ru_escaped} на Vondi"}}',
    (SELECT id FROM categories WHERE slug = '{parent_slug}' AND level = 2),
    3,
    (SELECT path FROM categories WHERE slug = '{parent_slug}' AND level = 2) || '/{slug}',
    {sort_order},
    '{icon}',
    true
  )"""
            values.append(value)
            added_slugs.append(slug)

    sql_content.append(",\n".join(values))
    sql_content.append(";")
    sql_content.append("")
    sql_content.append("-- Verification")
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
  RAISE NOTICE 'Successfully added L3 categories';
END $$;""")

    with open(output_file, 'w', encoding='utf-8') as f:
        f.write('\n'.join(sql_content))

    # Generate down migration
    down_content = """-- Rollback: Delete L3 categories

DELETE FROM categories WHERE slug IN (
"""
    down_content += ',\n'.join([f"  '{slug}'" for slug in added_slugs])
    down_content += "\n) AND level = 3;\n"

    down_file = output_file.replace('.up.sql', '.down.sql')
    with open(down_file, 'w', encoding='utf-8') as f:
        f.write(down_content)

    print(f"✅ Generated {output_file} with {len(values)} L3 categories")
    print(f"✅ Generated {down_file}")

    return len(values)


def main():
    parser = argparse.ArgumentParser(description='Generate L3 category migrations')
    parser.add_argument('--output-dir', type=str, default='migrations/', help='Output directory')
    args = parser.parse_args()

    # Part 1: Elektronika (90 L3)
    elektronika_parents = [
        'pametni-telefoni', 'laptop-racunari', 'tv-i-video', 'audio-oprema',
        'foto-i-video-kamere', 'pametni-satovi', 'gaming-oprema'
    ]
    count1 = generate_l3_migration_file(
        elektronika_parents,
        f"{args.output_dir}20251217100003_l3_elektronika.up.sql",
        1
    )

    # Part 2: Odeca i Obuca (90 L3)
    odeca_parents = [
        'muska-odeca', 'zenska-odeca', 'decija-odeca',
        'muska-obuca', 'zenska-obuca', 'decija-obuca'
    ]
    count2 = generate_l3_migration_file(
        odeca_parents,
        f"{args.output_dir}20251217100004_l3_odeca.up.sql",
        1000
    )

    # Part 3: Dom i Sport (80 L3)
    dom_sport_parents = [
        'namestaj-dnevna-soba', 'namestaj-spavaca-soba', 'kupatilo', 'rasveta',
        'fitnes-i-teretana', 'bicikli-i-trotineti', 'kampovanje'
    ]
    count3 = generate_l3_migration_file(
        dom_sport_parents,
        f"{args.output_dir}20251217100005_l3_dom_sport.up.sql",
        2000
    )

    print(f"\n📊 Summary:")
    print(f"  Elektronika: {count1} L3")
    print(f"  Odeca: {count2} L3")
    print(f"  Dom & Sport: {count3} L3")
    print(f"  TOTAL: {count1 + count2 + count3} L3")


if __name__ == '__main__':
    main()
