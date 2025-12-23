-- Migration: Complete L2 Categories (Part 1/2) - добавить 50 новых L2 категорий
-- Date: 2025-12-17
-- Target: Увеличить L2 с 301 до 351 (+50)

INSERT INTO categories (slug, name, description, meta_title, meta_description, parent_id, level, path, sort_order, icon, is_active)
VALUES
  -- ===== ELEKTRONIKA (18 новых L2) =====
  (
    'bluetooth-zvucnici-premium',
    '{"sr": "Bluetooth zvučnici premium", "en": "Premium Bluetooth Speakers", "ru": "Премиум Bluetooth колонки"}',
    '{"sr": "Visokokvalitetni bluetooth zvučnici sa naprednim funkcijama", "en": "High-quality Bluetooth speakers with advanced features", "ru": "Высококачественные Bluetooth колонки с расширенными функциями"}',
    '{"sr": "Premium bluetooth zvučnici | Vondi", "en": "Premium Bluetooth Speakers | Vondi", "ru": "Премиум Bluetooth колонки | Vondi"}',
    '{"sr": "Najbolji premium bluetooth zvučnici sa odličnim kvalitetom zvuka", "en": "Best premium Bluetooth speakers with excellent sound quality", "ru": "Лучшие премиальные Bluetooth колонки с отличным качеством звука"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/bluetooth-zvucnici-premium',
    201,
    '🔊',
    true
  ),
  (
    'bluetooth-slusalice-profesionalne',
    '{"sr": "Bluetooth slušalice profesionalne", "en": "Professional Bluetooth Headphones", "ru": "Профессиональные Bluetooth наушники"}',
    '{"sr": "Profesionalne bluetooth slušalice za studio i DJ", "en": "Professional Bluetooth headphones for studio and DJ", "ru": "Профессиональные Bluetooth наушники для студии и DJ"}',
    '{"sr": "Profesionalne bluetooth slušalice | Vondi", "en": "Professional Bluetooth Headphones | Vondi", "ru": "Профессиональные Bluetooth наушники | Vondi"}',
    '{"sr": "Najbolje profesionalne bluetooth slušalice za studijske uslove", "en": "Best professional Bluetooth headphones for studio conditions", "ru": "Лучшие профессиональные Bluetooth наушники для студийных условий"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/bluetooth-slusalice-profesionalne',
    202,
    '🎧',
    true
  ),
  (
    'wi-fi-routeri-mesh',
    '{"sr": "Wi-Fi routeri mesh", "en": "Mesh Wi-Fi Routers", "ru": "Mesh Wi-Fi роутеры"}',
    '{"sr": "Mesh Wi-Fi sistemi za celokupno pokrivanje doma", "en": "Mesh Wi-Fi systems for whole-home coverage", "ru": "Mesh Wi-Fi системы для покрытия всего дома"}',
    '{"sr": "Mesh Wi-Fi routeri | Vondi", "en": "Mesh Wi-Fi Routers | Vondi", "ru": "Mesh Wi-Fi роутеры | Vondi"}',
    '{"sr": "Najbolji mesh Wi-Fi sistemi za snažan signal u celom domu", "en": "Best mesh Wi-Fi systems for strong signal throughout the home", "ru": "Лучшие mesh Wi-Fi системы для сильного сигнала во всем доме"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/wi-fi-routeri-mesh',
    203,
    '📡',
    true
  ),
  (
    'powerline-adapteri',
    '{"sr": "Powerline adapteri", "en": "Powerline Adapters", "ru": "Powerline адаптеры"}',
    '{"sr": "Internet preko struje - powerline adapteri", "en": "Internet over power lines - powerline adapters", "ru": "Интернет по электропроводке - powerline адаптеры"}',
    '{"sr": "Powerline adapteri | Vondi", "en": "Powerline Adapters | Vondi", "ru": "Powerline адаптеры | Vondi"}',
    '{"sr": "Powerline adapteri za internet preko strujne mreže", "en": "Powerline adapters for internet over electrical network", "ru": "Powerline адаптеры для интернета по электросети"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/powerline-adapteri',
    204,
    '⚡',
    true
  ),
  (
    'usb-hub-powered',
    '{"sr": "USB hub sa napajanjem", "en": "Powered USB Hub", "ru": "USB хаб с питанием"}',
    '{"sr": "USB hub-ovi sa eksternim napajanjem", "en": "USB hubs with external power supply", "ru": "USB хабы с внешним питанием"}',
    '{"sr": "USB hub sa napajanjem | Vondi", "en": "Powered USB Hub | Vondi", "ru": "USB хаб с питанием | Vondi"}',
    '{"sr": "Powered USB hubovi za priključivanje više uređaja", "en": "Powered USB hubs for connecting multiple devices", "ru": "Питаемые USB хабы для подключения нескольких устройств"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/usb-hub-powered',
    205,
    '🔌',
    true
  ),
  (
    'eksterni-hard-diskovi-ssd',
    '{"sr": "Eksterni SSD diskovi", "en": "External SSD Drives", "ru": "Внешние SSD диски"}',
    '{"sr": "Brzi eksterni SSD diskovi za prenos podataka", "en": "Fast external SSD drives for data transfer", "ru": "Быстрые внешние SSD диски для передачи данных"}',
    '{"sr": "Eksterni SSD diskovi | Vondi", "en": "External SSD Drives | Vondi", "ru": "Внешние SSD диски | Vondi"}',
    '{"sr": "Najbolji eksterni SSD diskovi sa velikom brzinom prenosa", "en": "Best external SSD drives with high transfer speed", "ru": "Лучшие внешние SSD диски с высокой скоростью передачи"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/eksterni-hard-diskovi-ssd',
    206,
    '💾',
    true
  ),
  (
    'memorijske-kartice-pro',
    '{"sr": "Memorijske kartice profesionalne", "en": "Professional Memory Cards", "ru": "Профессиональные карты памяти"}',
    '{"sr": "Profesionalne SD i microSD kartice", "en": "Professional SD and microSD cards", "ru": "Профессиональные SD и microSD карты"}',
    '{"sr": "Profesionalne memorijske kartice | Vondi", "en": "Professional Memory Cards | Vondi", "ru": "Профессиональные карты памяти | Vondi"}',
    '{"sr": "Pro memorijske kartice za kamere i dronove", "en": "Pro memory cards for cameras and drones", "ru": "Профессиональные карты памяти для камер и дронов"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/memorijske-kartice-pro',
    207,
    '💳',
    true
  ),
  (
    'citaci-e-ink',
    '{"sr": "E-ink čitači", "en": "E-ink Readers", "ru": "E-ink ридеры"}',
    '{"sr": "E-ink čitači elektronskih knjiga", "en": "E-ink electronic book readers", "ru": "E-ink ридеры электронных книг"}',
    '{"sr": "E-ink čitači knjiga | Vondi", "en": "E-ink Book Readers | Vondi", "ru": "E-ink ридеры книг | Vondi"}',
    '{"sr": "Najbolji e-ink čitači za čitanje elektronskih knjiga", "en": "Best e-ink readers for reading electronic books", "ru": "Лучшие e-ink ридеры для чтения электронных книг"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/citaci-e-ink',
    208,
    '📖',
    true
  ),
  (
    'elektronske-knjige-citalice',
    '{"sr": "Elektronske knjige čitalice", "en": "Electronic Book Readers", "ru": "Электронные книги"}',
    '{"sr": "Čitalice elektronskih knjiga", "en": "E-book readers and devices", "ru": "Читалки электронных книг"}',
    '{"sr": "Elektronske knjige čitalice | Vondi", "en": "E-book Readers | Vondi", "ru": "Электронные книги | Vondi"}',
    '{"sr": "Čitalice za elektronske knjige sa e-ink ekranom", "en": "E-book readers with e-ink screen", "ru": "Читалки электронных книг с e-ink экраном"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/elektronske-knjige-citalice',
    209,
    '📚',
    true
  ),
  (
    'bluetooth-tastature',
    '{"sr": "Bluetooth tastature", "en": "Bluetooth Keyboards", "ru": "Bluetooth клавиатуры"}',
    '{"sr": "Bežične bluetooth tastature", "en": "Wireless Bluetooth keyboards", "ru": "Беспроводные Bluetooth клавиатуры"}',
    '{"sr": "Bluetooth tastature | Vondi", "en": "Bluetooth Keyboards | Vondi", "ru": "Bluetooth клавиатуры | Vondi"}',
    '{"sr": "Najbolje bluetooth tastature za kompjuter i tablet", "en": "Best Bluetooth keyboards for computer and tablet", "ru": "Лучшие Bluetooth клавиатуры для компьютера и планшета"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/bluetooth-tastature',
    210,
    '⌨️',
    true
  ),
  (
    'graficke-tablice',
    '{"sr": "Grafičke table", "en": "Graphics Tablets", "ru": "Графические планшеты"}',
    '{"sr": "Profesionalne grafičke table za crtanje", "en": "Professional graphics tablets for drawing", "ru": "Профессиональные графические планшеты для рисования"}',
    '{"sr": "Grafičke table | Vondi", "en": "Graphics Tablets | Vondi", "ru": "Графические планшеты | Vondi"}',
    '{"sr": "Grafičke table za digitalno crtanje i dizajn", "en": "Graphics tablets for digital drawing and design", "ru": "Графические планшеты для цифрового рисования и дизайна"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/graficke-tablice',
    211,
    '🖊️',
    true
  ),
  (
    '3d-stampaci',
    '{"sr": "3D štampači", "en": "3D Printers", "ru": "3D принтеры"}',
    '{"sr": "3D štampači za kućnu upotrebu", "en": "3D printers for home use", "ru": "3D принтеры для домашнего использования"}',
    '{"sr": "3D štampači | Vondi", "en": "3D Printers | Vondi", "ru": "3D принтеры | Vondi"}',
    '{"sr": "Najbolji 3D štampači za hobiste i profesionalce", "en": "Best 3D printers for hobbyists and professionals", "ru": "Лучшие 3D принтеры для любителей и профессионалов"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/3d-stampaci',
    212,
    '🖨️',
    true
  ),
  (
    'streaming-oprema',
    '{"sr": "Streaming oprema", "en": "Streaming Equipment", "ru": "Стриминг оборудование"}',
    '{"sr": "Oprema za strimovanje i live broadcasting", "en": "Equipment for streaming and live broadcasting", "ru": "Оборудование для стриминга и прямых трансляций"}',
    '{"sr": "Streaming oprema | Vondi", "en": "Streaming Equipment | Vondi", "ru": "Стриминг оборудование | Vondi"}',
    '{"sr": "Profesionalna oprema za streaming na Twitch i YouTube", "en": "Professional equipment for streaming on Twitch and YouTube", "ru": "Профессиональное оборудование для стриминга на Twitch и YouTube"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/streaming-oprema',
    213,
    '🎥',
    true
  ),
  (
    'pametne-sijalice',
    '{"sr": "Pametne sijalice", "en": "Smart Light Bulbs", "ru": "Умные лампочки"}',
    '{"sr": "LED pametne sijalice sa WiFi kontrolom", "en": "LED smart bulbs with WiFi control", "ru": "LED умные лампочки с WiFi управлением"}',
    '{"sr": "Pametne LED sijalice | Vondi", "en": "Smart LED Bulbs | Vondi", "ru": "Умные LED лампочки | Vondi"}',
    '{"sr": "Pametne sijalice koje se kontrolišu preko telefona", "en": "Smart bulbs controlled via phone", "ru": "Умные лампочки управляемые через телефон"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/pametne-sijalice',
    214,
    '💡',
    true
  ),
  (
    'punjaci-bežicni',
    '{"sr": "Bežični punjači", "en": "Wireless Chargers", "ru": "Беспроводные зарядки"}',
    '{"sr": "Wireless punjači za telefone i satove", "en": "Wireless chargers for phones and watches", "ru": "Беспроводные зарядки для телефонов и часов"}',
    '{"sr": "Bežični punjači | Vondi", "en": "Wireless Chargers | Vondi", "ru": "Беспроводные зарядки | Vondi"}',
    '{"sr": "Najbolji wireless punjači sa brzim punjenjem", "en": "Best wireless chargers with fast charging", "ru": "Лучшие беспроводные зарядки с быстрой зарядкой"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/punjaci-bezicni',
    215,
    '🔋',
    true
  ),
  (
    'stabilizatori-gimbal',
    '{"sr": "Stabilizatori gimbal", "en": "Gimbal Stabilizers", "ru": "Стабилизаторы gimbal"}',
    '{"sr": "Gimbal stabilizatori za telefone i kamere", "en": "Gimbal stabilizers for phones and cameras", "ru": "Gimbal стабилизаторы для телефонов и камер"}',
    '{"sr": "Gimbal stabilizatori | Vondi", "en": "Gimbal Stabilizers | Vondi", "ru": "Gimbal стабилизаторы | Vondi"}',
    '{"sr": "Gimbal stabilizatori za glatko snimanje videa", "en": "Gimbal stabilizers for smooth video recording", "ru": "Gimbal стабилизаторы для плавной видеосъемки"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/stabilizatori-gimbal',
    216,
    '📹',
    true
  ),
  (
    'mikroskopi-digitalni',
    '{"sr": "Digitalni mikroskopi", "en": "Digital Microscopes", "ru": "Цифровые микроскопы"}',
    '{"sr": "USB digitalni mikroskopi", "en": "USB digital microscopes", "ru": "USB цифровые микроскопы"}',
    '{"sr": "Digitalni mikroskopi | Vondi", "en": "Digital Microscopes | Vondi", "ru": "Цифровые микроскопы | Vondi"}',
    '{"sr": "Digitalni mikroskopi za edukaciju i hobi", "en": "Digital microscopes for education and hobby", "ru": "Цифровые микроскопы для образования и хобби"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/mikroskopi-digitalni',
    217,
    '🔬',
    true
  ),
  (
    'solarni-punjaci',
    '{"sr": "Solarni punjači", "en": "Solar Chargers", "ru": "Солнечные зарядки"}',
    '{"sr": "Solarni punjači za telefone i uređaje", "en": "Solar chargers for phones and devices", "ru": "Солнечные зарядки для телефонов и устройств"}',
    '{"sr": "Solarni punjači | Vondi", "en": "Solar Chargers | Vondi", "ru": "Солнечные зарядки | Vondi"}',
    '{"sr": "Solarni punjači za kampovanje i putovanja", "en": "Solar chargers for camping and travel", "ru": "Солнечные зарядки для кемпинга и путешествий"}',
    (SELECT id FROM categories WHERE slug = 'elektronika' AND level = 1),
    2,
    'elektronika/solarni-punjaci',
    218,
    '☀️',
    true
  );

-- ===== ПРОДОЛЖЕНИЕ СЛЕДУЕТ В PART 2 =====

-- Verification: проверка на дубликаты
DO $$
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
  RAISE NOTICE 'Successfully added 18 new L2 categories';
END $$;
