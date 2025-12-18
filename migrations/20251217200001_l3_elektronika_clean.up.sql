-- Migration: L3 Elektronika Categories (90 new L3)
-- Parent: Elektronika L1 categories
-- Total L3 after: 26 + 90 = 116

-- ============================================================================
-- PAMETNI TELEFONI (parent: pametni-telefoni) - 5 new L3
-- Already exists: apple-iphone, samsung-telefoni, xiaomi-telefoni, huawei-telefoni
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Google Pixel
('google-pixel-telefoni',
 '{"sr": "Google Pixel telefoni", "en": "Google Pixel Phones", "ru": "Телефоны Google Pixel"}',
 (SELECT id FROM categories WHERE slug = 'pametni-telefoni' AND level = 2),
 3, 'elektronika/pametni-telefoni/google-pixel-telefoni', 5, true, NOW(), NOW()),

-- OnePlus
('oneplus-telefoni',
 '{"sr": "OnePlus telefoni", "en": "OnePlus Phones", "ru": "Телефоны OnePlus"}',
 (SELECT id FROM categories WHERE slug = 'pametni-telefoni' AND level = 2),
 3, 'elektronika/pametni-telefoni/oneplus-telefoni', 6, true, NOW(), NOW()),

-- Realme
('realme-telefoni',
 '{"sr": "Realme telefoni", "en": "Realme Phones", "ru": "Телефоны Realme"}',
 (SELECT id FROM categories WHERE slug = 'pametni-telefoni' AND level = 2),
 3, 'elektronika/pametni-telefoni/realme-telefoni', 7, true, NOW(), NOW()),

-- Oppo
('oppo-telefoni',
 '{"sr": "Oppo telefoni", "en": "Oppo Phones", "ru": "Телефоны Oppo"}',
 (SELECT id FROM categories WHERE slug = 'pametni-telefoni' AND level = 2),
 3, 'elektronika/pametni-telefoni/oppo-telefoni', 8, true, NOW(), NOW()),

-- Budget phones
('telefoni-budget-do-200e',
 '{"sr": "Budget telefoni do 200 evra", "en": "Budget Phones Under 200 EUR", "ru": "Бюджетные телефоны до 200 евро"}',
 (SELECT id FROM categories WHERE slug = 'pametni-telefoni' AND level = 2),
 3, 'elektronika/pametni-telefoni/telefoni-budget-do-200e', 9, true, NOW(), NOW());

-- ============================================================================
-- LAPTOP RACUNARI (parent: laptop-racunari) - 7 new L3
-- Already exists: laptop-gaming, laptop-profesionalni, laptop-ultrabook
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- 2-in-1
('laptop-2-u-1',
 '{"sr": "Laptop 2 u 1", "en": "2-in-1 Laptops", "ru": "Ноутбуки 2-в-1"}',
 (SELECT id FROM categories WHERE slug = 'laptop-racunari' AND level = 2),
 3, 'elektronika/laptop-racunari/laptop-2-u-1', 4, true, NOW(), NOW()),

-- Chromebook
('laptop-chromebook',
 '{"sr": "Chromebook laptop", "en": "Chromebook Laptops", "ru": "Chromebook ноутбуки"}',
 (SELECT id FROM categories WHERE slug = 'laptop-racunari' AND level = 2),
 3, 'elektronika/laptop-racunari/laptop-chromebook', 5, true, NOW(), NOW()),

-- MacBook
('laptop-macbook',
 '{"sr": "MacBook laptop", "en": "MacBook Laptops", "ru": "MacBook ноутбуки"}',
 (SELECT id FROM categories WHERE slug = 'laptop-racunari' AND level = 2),
 3, 'elektronika/laptop-racunari/laptop-macbook', 6, true, NOW(), NOW()),

-- Lenovo ThinkPad
('laptop-lenovo-thinkpad',
 '{"sr": "Lenovo ThinkPad", "en": "Lenovo ThinkPad", "ru": "Lenovo ThinkPad"}',
 (SELECT id FROM categories WHERE slug = 'laptop-racunari' AND level = 2),
 3, 'elektronika/laptop-racunari/laptop-lenovo-thinkpad', 7, true, NOW(), NOW()),

-- Dell XPS
('laptop-dell-xps',
 '{"sr": "Dell XPS", "en": "Dell XPS", "ru": "Dell XPS"}',
 (SELECT id FROM categories WHERE slug = 'laptop-racunari' AND level = 2),
 3, 'elektronika/laptop-racunari/laptop-dell-xps', 8, true, NOW(), NOW()),

-- HP Pavilion
('laptop-hp-pavilion',
 '{"sr": "HP Pavilion", "en": "HP Pavilion", "ru": "HP Pavilion"}',
 (SELECT id FROM categories WHERE slug = 'laptop-racunari' AND level = 2),
 3, 'elektronika/laptop-racunari/laptop-hp-pavilion', 9, true, NOW(), NOW()),

-- Asus ROG
('laptop-asus-rog',
 '{"sr": "Asus ROG", "en": "Asus ROG", "ru": "Asus ROG"}',
 (SELECT id FROM categories WHERE slug = 'laptop-racunari' AND level = 2),
 3, 'elektronika/laptop-racunari/laptop-asus-rog', 10, true, NOW(), NOW());

-- ============================================================================
-- TV I VIDEO (parent: tv-i-video) - 15 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- LED TV by size
('led-tv-32-43-inca',
 '{"sr": "LED TV 32-43 inča", "en": "LED TV 32-43 inch", "ru": "LED TV 32-43 дюйма"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/led-tv-32-43-inca', 1, true, NOW(), NOW()),

('led-tv-50-55-inca',
 '{"sr": "LED TV 50-55 inča", "en": "LED TV 50-55 inch", "ru": "LED TV 50-55 дюймов"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/led-tv-50-55-inca', 2, true, NOW(), NOW()),

('led-tv-65-75-inca',
 '{"sr": "LED TV 65-75 inča", "en": "LED TV 65-75 inch", "ru": "LED TV 65-75 дюймов"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/led-tv-65-75-inca', 3, true, NOW(), NOW()),

-- OLED/QLED
('oled-tv',
 '{"sr": "OLED televizori", "en": "OLED TVs", "ru": "OLED телевизоры"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/oled-tv', 4, true, NOW(), NOW()),

('qled-tv',
 '{"sr": "QLED televizori", "en": "QLED TVs", "ru": "QLED телевизоры"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/qled-tv', 5, true, NOW(), NOW()),

-- Smart TV
('smart-tv-android',
 '{"sr": "Smart TV Android", "en": "Android Smart TV", "ru": "Smart TV Android"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/smart-tv-android', 6, true, NOW(), NOW()),

('smart-tv-tizen',
 '{"sr": "Smart TV Tizen", "en": "Tizen Smart TV", "ru": "Smart TV Tizen"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/smart-tv-tizen', 7, true, NOW(), NOW()),

-- Soundbars
('soundbar-2-1',
 '{"sr": "Soundbar 2.1", "en": "Soundbar 2.1", "ru": "Саундбар 2.1"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/soundbar-2-1', 8, true, NOW(), NOW()),

('soundbar-5-1',
 '{"sr": "Soundbar 5.1", "en": "Soundbar 5.1", "ru": "Саундбар 5.1"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/soundbar-5-1', 9, true, NOW(), NOW()),

-- Speakers
('bluetooth-zvucnici',
 '{"sr": "Bluetooth zvučnici", "en": "Bluetooth Speakers", "ru": "Bluetooth колонки"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/bluetooth-zvucnici', 10, true, NOW(), NOW()),

('hi-fi-sistemi',
 '{"sr": "Hi-Fi sistemi", "en": "Hi-Fi Systems", "ru": "Hi-Fi системы"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/hi-fi-sistemi', 11, true, NOW(), NOW()),

-- Projectors
('projektori-full-hd',
 '{"sr": "Projektori Full HD", "en": "Full HD Projectors", "ru": "Проекторы Full HD"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/projektori-full-hd', 12, true, NOW(), NOW()),

('projektori-4k',
 '{"sr": "Projektori 4K", "en": "4K Projectors", "ru": "Проекторы 4K"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/projektori-4k', 13, true, NOW(), NOW()),

-- AV Receivers
('av-resiveri',
 '{"sr": "AV resiveri", "en": "AV Receivers", "ru": "AV ресиверы"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/av-resiveri', 14, true, NOW(), NOW()),

-- Home Theater
('kucni-bioskop',
 '{"sr": "Kućni bioskop", "en": "Home Theater", "ru": "Домашний кинотеатр"}',
 (SELECT id FROM categories WHERE slug = 'tv-i-video' AND level = 2),
 3, 'elektronika/tv-i-video/kucni-bioskop', 15, true, NOW(), NOW());

-- ============================================================================
-- RACUNARSKE KOMPONENTE (parent: racunarske-komponente) - 20 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Graphics Cards NVIDIA
('graficke-rtx-4000',
 '{"sr": "Grafičke RTX 4000 serija", "en": "RTX 4000 Series Graphics", "ru": "Видеокарты RTX 4000"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/graficke-rtx-4000', 1, true, NOW(), NOW()),

('graficke-rtx-3000',
 '{"sr": "Grafičke RTX 3000 serija", "en": "RTX 3000 Series Graphics", "ru": "Видеокарты RTX 3000"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/graficke-rtx-3000', 2, true, NOW(), NOW()),

('graficke-gtx-1000',
 '{"sr": "Grafičke GTX 1000 serija", "en": "GTX 1000 Series Graphics", "ru": "Видеокарты GTX 1000"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/graficke-gtx-1000', 3, true, NOW(), NOW()),

-- Graphics Cards AMD
('graficke-amd-rx-7000',
 '{"sr": "Grafičke AMD RX 7000 serija", "en": "AMD RX 7000 Series Graphics", "ru": "Видеокарты AMD RX 7000"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/graficke-amd-rx-7000', 4, true, NOW(), NOW()),

('graficke-amd-rx-6000',
 '{"sr": "Grafičke AMD RX 6000 serija", "en": "AMD RX 6000 Series Graphics", "ru": "Видеокарты AMD RX 6000"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/graficke-amd-rx-6000', 5, true, NOW(), NOW()),

-- Processors Intel
('procesor-intel-i9',
 '{"sr": "Procesor Intel Core i9", "en": "Intel Core i9 Processor", "ru": "Процессор Intel Core i9"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/procesor-intel-i9', 6, true, NOW(), NOW()),

('procesor-intel-i7',
 '{"sr": "Procesor Intel Core i7", "en": "Intel Core i7 Processor", "ru": "Процессор Intel Core i7"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/procesor-intel-i7', 7, true, NOW(), NOW()),

('procesor-intel-i5',
 '{"sr": "Procesor Intel Core i5", "en": "Intel Core i5 Processor", "ru": "Процессор Intel Core i5"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/procesor-intel-i5', 8, true, NOW(), NOW()),

-- Processors AMD
('procesor-amd-ryzen-9',
 '{"sr": "Procesor AMD Ryzen 9", "en": "AMD Ryzen 9 Processor", "ru": "Процессор AMD Ryzen 9"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/procesor-amd-ryzen-9', 9, true, NOW(), NOW()),

('procesor-amd-ryzen-7',
 '{"sr": "Procesor AMD Ryzen 7", "en": "AMD Ryzen 7 Processor", "ru": "Процессор AMD Ryzen 7"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/procesor-amd-ryzen-7', 10, true, NOW(), NOW()),

('procesor-amd-ryzen-5',
 '{"sr": "Procesor AMD Ryzen 5", "en": "AMD Ryzen 5 Processor", "ru": "Процессор AMD Ryzen 5"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/procesor-amd-ryzen-5', 11, true, NOW(), NOW()),

-- RAM DDR5
('ram-ddr5-32gb',
 '{"sr": "RAM DDR5 32GB", "en": "RAM DDR5 32GB", "ru": "Оперативная память DDR5 32ГБ"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/ram-ddr5-32gb', 12, true, NOW(), NOW()),

('ram-ddr5-16gb',
 '{"sr": "RAM DDR5 16GB", "en": "RAM DDR5 16GB", "ru": "Оперативная память DDR5 16ГБ"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/ram-ddr5-16gb', 13, true, NOW(), NOW()),

-- RAM DDR4
('ram-ddr4-16gb',
 '{"sr": "RAM DDR4 16GB", "en": "RAM DDR4 16GB", "ru": "Оперативная память DDR4 16ГБ"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/ram-ddr4-16gb', 14, true, NOW(), NOW()),

('ram-ddr4-8gb',
 '{"sr": "RAM DDR4 8GB", "en": "RAM DDR4 8GB", "ru": "Оперативная память DDR4 8ГБ"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/ram-ddr4-8gb', 15, true, NOW(), NOW()),

-- SSD NVMe
('ssd-nvme-1tb',
 '{"sr": "SSD NVMe 1TB", "en": "SSD NVMe 1TB", "ru": "SSD NVMe 1ТБ"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/ssd-nvme-1tb', 16, true, NOW(), NOW()),

('ssd-nvme-500gb',
 '{"sr": "SSD NVMe 500GB", "en": "SSD NVMe 500GB", "ru": "SSD NVMe 500ГБ"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/ssd-nvme-500gb', 17, true, NOW(), NOW()),

-- SSD SATA
('ssd-sata-1tb',
 '{"sr": "SSD SATA 1TB", "en": "SSD SATA 1TB", "ru": "SSD SATA 1ТБ"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/ssd-sata-1tb', 18, true, NOW(), NOW()),

-- Motherboards
('maticna-ploca-intel-z790',
 '{"sr": "Matična ploča Intel Z790", "en": "Intel Z790 Motherboard", "ru": "Материнская плата Intel Z790"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/maticna-ploca-intel-z790', 19, true, NOW(), NOW()),

('maticna-ploca-amd-x670',
 '{"sr": "Matična ploča AMD X670", "en": "AMD X670 Motherboard", "ru": "Материнская плата AMD X670"}',
 (SELECT id FROM categories WHERE slug = 'racunarske-komponente' AND level = 2),
 3, 'elektronika/racunarske-komponente/maticna-ploca-amd-x670', 20, true, NOW(), NOW());

-- ============================================================================
-- GAMING OPREMA (parent: gaming-oprema) - 15 new L3
-- Already exists: gaming-tastature, gaming-misevi
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Consoles
('playstation-5',
 '{"sr": "PlayStation 5 konzole", "en": "PlayStation 5 Consoles", "ru": "Консоли PlayStation 5"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/playstation-5', 3, true, NOW(), NOW()),

('playstation-5-igre',
 '{"sr": "PlayStation 5 igre", "en": "PlayStation 5 Games", "ru": "Игры для PlayStation 5"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/playstation-5-igre', 4, true, NOW(), NOW()),

('xbox-series-x',
 '{"sr": "Xbox Series X konzole", "en": "Xbox Series X Consoles", "ru": "Консоли Xbox Series X"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/xbox-series-x', 5, true, NOW(), NOW()),

('xbox-series-s',
 '{"sr": "Xbox Series S konzole", "en": "Xbox Series S Consoles", "ru": "Консоли Xbox Series S"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/xbox-series-s', 6, true, NOW(), NOW()),

('xbox-igre',
 '{"sr": "Xbox igre", "en": "Xbox Games", "ru": "Игры для Xbox"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/xbox-igre', 7, true, NOW(), NOW()),

('nintendo-switch',
 '{"sr": "Nintendo Switch konzole", "en": "Nintendo Switch Consoles", "ru": "Консоли Nintendo Switch"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/nintendo-switch', 8, true, NOW(), NOW()),

('nintendo-igre',
 '{"sr": "Nintendo Switch igre", "en": "Nintendo Switch Games", "ru": "Игры для Nintendo Switch"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/nintendo-igre', 9, true, NOW(), NOW()),

-- Controllers
('gejmpad-ps5',
 '{"sr": "Gejmpad za PS5", "en": "PS5 Controllers", "ru": "Геймпады для PS5"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/gejmpad-ps5', 10, true, NOW(), NOW()),

('gejmpad-xbox',
 '{"sr": "Gejmpad za Xbox", "en": "Xbox Controllers", "ru": "Геймпады для Xbox"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/gejmpad-xbox', 11, true, NOW(), NOW()),

-- Gaming Peripherals
('gaming-slusalice',
 '{"sr": "Gaming slušalice", "en": "Gaming Headsets", "ru": "Игровые наушники"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/gaming-slusalice', 12, true, NOW(), NOW()),

('gaming-monitor-144hz',
 '{"sr": "Gaming monitor 144Hz", "en": "144Hz Gaming Monitors", "ru": "Игровые мониторы 144Гц"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/gaming-monitor-144hz', 13, true, NOW(), NOW()),

('gaming-monitor-240hz',
 '{"sr": "Gaming monitor 240Hz", "en": "240Hz Gaming Monitors", "ru": "Игровые мониторы 240Гц"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/gaming-monitor-240hz', 14, true, NOW(), NOW()),

('gaming-stolice',
 '{"sr": "Gaming stolice", "en": "Gaming Chairs", "ru": "Игровые кресла"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/gaming-stolice', 15, true, NOW(), NOW()),

('rgb-rasveta',
 '{"sr": "RGB rasveta", "en": "RGB Lighting", "ru": "RGB подсветка"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/rgb-rasveta', 16, true, NOW(), NOW()),

('streaming-mikrofon',
 '{"sr": "Streaming mikrofoni", "en": "Streaming Microphones", "ru": "Микрофоны для стриминга"}',
 (SELECT id FROM categories WHERE slug = 'gaming-oprema' AND level = 2),
 3, 'elektronika/gaming-oprema/streaming-mikrofon', 17, true, NOW(), NOW());

-- ============================================================================
-- FOTO I VIDEO KAMERE (parent: foto-i-video-kamere) - 12 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- DSLR Cameras
('dslr-canon',
 '{"sr": "DSLR Canon fotoaparati", "en": "Canon DSLR Cameras", "ru": "Зеркальные камеры Canon"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/dslr-canon', 1, true, NOW(), NOW()),

('dslr-nikon',
 '{"sr": "DSLR Nikon fotoaparati", "en": "Nikon DSLR Cameras", "ru": "Зеркальные камеры Nikon"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/dslr-nikon', 2, true, NOW(), NOW()),

-- Mirrorless Cameras
('mirrorless-sony',
 '{"sr": "Mirrorless Sony fotoaparati", "en": "Sony Mirrorless Cameras", "ru": "Беззеркальные камеры Sony"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/mirrorless-sony', 3, true, NOW(), NOW()),

('mirrorless-fujifilm',
 '{"sr": "Mirrorless Fujifilm fotoaparati", "en": "Fujifilm Mirrorless Cameras", "ru": "Беззеркальные камеры Fujifilm"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/mirrorless-fujifilm', 4, true, NOW(), NOW()),

-- Action Cameras
('action-kamera-gopro',
 '{"sr": "GoPro action kamere", "en": "GoPro Action Cameras", "ru": "Экшн-камеры GoPro"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/action-kamera-gopro', 5, true, NOW(), NOW()),

('action-kamera-dji',
 '{"sr": "DJI action kamere", "en": "DJI Action Cameras", "ru": "Экшн-камеры DJI"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/action-kamera-dji', 6, true, NOW(), NOW()),

-- Drones
('dron-sa-kamerom',
 '{"sr": "Dronovi sa kamerom", "en": "Camera Drones", "ru": "Дроны с камерой"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/dron-sa-kamerom', 7, true, NOW(), NOW()),

-- Video Cameras
('4k-video-kamera',
 '{"sr": "4K video kamere", "en": "4K Video Cameras", "ru": "4K видеокамеры"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/4k-video-kamera', 8, true, NOW(), NOW()),

-- Lenses
('objektiv-za-canon',
 '{"sr": "Objektivi za Canon", "en": "Canon Lenses", "ru": "Объективы для Canon"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/objektiv-za-canon', 9, true, NOW(), NOW()),

('objektiv-za-nikon',
 '{"sr": "Objektivi za Nikon", "en": "Nikon Lenses", "ru": "Объективы для Nikon"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/objektiv-za-nikon', 10, true, NOW(), NOW()),

-- Tripods
('stativ-foto-video',
 '{"sr": "Stativi za foto i video", "en": "Photo Video Tripods", "ru": "Штативы для фото и видео"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/stativ-foto-video', 11, true, NOW(), NOW()),

-- Lighting
('foto-rasveta',
 '{"sr": "Foto rasveta i studijski kompleti", "en": "Photo Lighting and Studio Kits", "ru": "Фото освещение и студийные комплекты"}',
 (SELECT id FROM categories WHERE slug = 'foto-i-video-kamere' AND level = 2),
 3, 'elektronika/foto-i-video-kamere/foto-rasveta', 12, true, NOW(), NOW());

-- ============================================================================
-- PAMETNI SATOVI (parent: pametni-satovi) - 8 new L3
-- ============================================================================

INSERT INTO categories (slug, name, parent_id, level, path, sort_order, is_active, created_at, updated_at)
VALUES
-- Apple Watch
('apple-watch-series',
 '{"sr": "Apple Watch Series", "en": "Apple Watch Series", "ru": "Apple Watch Series"}',
 (SELECT id FROM categories WHERE slug = 'pametni-satovi' AND level = 2),
 3, 'elektronika/pametni-satovi/apple-watch-series', 1, true, NOW(), NOW()),

('apple-watch-ultra',
 '{"sr": "Apple Watch Ultra", "en": "Apple Watch Ultra", "ru": "Apple Watch Ultra"}',
 (SELECT id FROM categories WHERE slug = 'pametni-satovi' AND level = 2),
 3, 'elektronika/pametni-satovi/apple-watch-ultra', 2, true, NOW(), NOW()),

-- Samsung Galaxy Watch
('samsung-galaxy-watch',
 '{"sr": "Samsung Galaxy Watch", "en": "Samsung Galaxy Watch", "ru": "Samsung Galaxy Watch"}',
 (SELECT id FROM categories WHERE slug = 'pametni-satovi' AND level = 2),
 3, 'elektronika/pametni-satovi/samsung-galaxy-watch', 3, true, NOW(), NOW()),

-- Garmin
('garmin-fitness-sat',
 '{"sr": "Garmin fitnes satovi", "en": "Garmin Fitness Watches", "ru": "Фитнес часы Garmin"}',
 (SELECT id FROM categories WHERE slug = 'pametni-satovi' AND level = 2),
 3, 'elektronika/pametni-satovi/garmin-fitness-sat', 4, true, NOW(), NOW()),

-- Xiaomi
('xiaomi-mi-band',
 '{"sr": "Xiaomi Mi Band narukvice", "en": "Xiaomi Mi Band", "ru": "Xiaomi Mi Band"}',
 (SELECT id FROM categories WHERE slug = 'pametni-satovi' AND level = 2),
 3, 'elektronika/pametni-satovi/xiaomi-mi-band', 5, true, NOW(), NOW()),

('xiaomi-watch',
 '{"sr": "Xiaomi pametni satovi", "en": "Xiaomi Smart Watches", "ru": "Умные часы Xiaomi"}',
 (SELECT id FROM categories WHERE slug = 'pametni-satovi' AND level = 2),
 3, 'elektronika/pametni-satovi/xiaomi-watch', 6, true, NOW(), NOW()),

-- Huawei
('huawei-watch-gt',
 '{"sr": "Huawei Watch GT", "en": "Huawei Watch GT", "ru": "Huawei Watch GT"}',
 (SELECT id FROM categories WHERE slug = 'pametni-satovi' AND level = 2),
 3, 'elektronika/pametni-satovi/huawei-watch-gt', 7, true, NOW(), NOW()),

-- Budget
('pametni-satovi-budget',
 '{"sr": "Budget pametni satovi", "en": "Budget Smart Watches", "ru": "Бюджетные умные часы"}',
 (SELECT id FROM categories WHERE slug = 'pametni-satovi' AND level = 2),
 3, 'elektronika/pametni-satovi/pametni-satovi-budget', 8, true, NOW(), NOW());

-- ============================================================================
-- Migration Summary
-- ============================================================================
-- Total new L3 categories: 90
-- Categories breakdown:
--   - Pametni telefoni: 5 (total 9 L3)
--   - Laptop racunari: 7 (total 10 L3)
--   - TV i video: 15
--   - Racunarske komponente: 20
--   - Gaming oprema: 15 (total 17 L3)
--   - Foto i video kamere: 12
--   - Pametni satovi: 8
-- ============================================================================
