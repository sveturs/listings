-- ============================================================================
-- Migration: Add more demo listings for investor demonstration
-- Date: 2025-01-27
-- Purpose: Create more diverse marketplace listings
-- ============================================================================

-- ============================================================================
-- REAL ESTATE LISTINGS
-- ============================================================================

-- Apartments for sale (using 1401 - Stanovi)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1401, 'Stan 65m2 Liman 3 - novogradnja',
 'Prelep stan u novogradnji na Limanu 3. Dva spavaća sobe, dnevni boravak, kuhinja, terasa 12m2. Parking mesto u ceni. Useljiv odmah.',
 95000.00, 'active', 45.2485, 19.8335, 'Balzakova 45', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '2 days', NOW(), 3456),

(7, 1401, 'Lux penthouse 120m2 Centar',
 'Ekskluzivan penthouse u samom centru grada. Pogled na Dunav, 3 spavaće sobe, 2 kupatila, velika terasa 40m2. Kompletno namešten.',
 220000.00, 'active', 44.8078, 20.4448, 'Kneza Miloša 10', 'Beograd', 'Srbija',
 NOW() - INTERVAL '5 days', NOW(), 2890),

-- Houses for sale (using 1402 - Kuće)
(7, 1402, 'Kuća 200m2 sa bazenom Sremska Kamenica',
 'Moderna kuća sa 4 spavaće sobe, garažom za 2 automobila, bazenom i velikim dvorištem 800m2. Mirna lokacija, blizu škole.',
 280000.00, 'active', 45.2207, 19.8461, 'Fruška Gora 25', 'Sremska Kamenica', 'Srbija',
 NOW() - INTERVAL '8 days', NOW(), 1567);

-- ============================================================================
-- AUTOMOTIVE LISTINGS
-- ============================================================================

-- Cars (using 1301 - Lični automobili)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES
(7, 1301, 'BMW X5 3.0d 2021 - kao nov',
 'BMW X5 xDrive30d, 2021 godište, prešao samo 35000km. Full oprema, M paket, panorama, HEAD-UP display. Servisna knjiga, garancija.',
 75000.00, 'active', 45.2551, 19.8451, 'Bulevar Cara Lazara 100', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '3 days', NOW(), 4567),

(7, 1301, 'Mercedes-Benz E220d 2022',
 'E klasa limuzina, automatik 9G-tronic, AMG line paket. Crna metalik boja, bež enterier. Prvi vlasnik, garažiran.',
 65000.00, 'active', 44.8125, 20.4612, 'Autokomanda', 'Beograd', 'Srbija',
 NOW() - INTERVAL '6 days', NOW(), 3234),

(7, 1301, 'Volkswagen Golf 8 2.0 TDI 2023',
 'Golf VIII generacija, 2.0 TDI 150ks, DSG automatik. Style oprema, virtual cockpit, LED farovi. Fabricka garancija do 2026.',
 32000.00, 'active', 45.2467, 19.8515, 'Temerinska 50', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '4 days', NOW(), 2876),

-- Motorcycles (using 1302 - Motocikli)
(7, 1302, 'Yamaha MT-07 2022 - perfektno stanje',
 'Yamaha MT-07, 689ccm, ABS, quick shifter. Akrapovic izduvni sistem. Prešla samo 8000km. Redovno servisirana.',
 7500.00, 'active', 44.7988, 20.4685, 'Voždovac', 'Beograd', 'Srbija',
 NOW() - INTERVAL '7 days', NOW(), 1234);

-- ============================================================================
-- SERVICES LISTINGS
-- ============================================================================

INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES
-- IT Services (using 1902 - IT usluge)
(7, 1902, 'Izrada web sajtova - WordPress, React',
 'Profesionalna izrada sajtova po meri. WordPress, React, Next.js. SEO optimizacija, responsive dizajn. Portfolio sa 50+ projekata.',
 500.00, 'active', 45.2551, 19.8451, 'Online', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '10 days', NOW(), 890),

-- Construction services (using 1901 - Građevinske usluge)
(7, 1901, 'Renoviranje stanova - kompletan servis',
 'Tim majstora za kompletno renoviranje. Gipsani radovi, keramika, parket, molerski radovi. Garancija na radove 5 godina.',
 50.00, 'active', 44.8078, 20.4448, 'Ceo Beograd', 'Beograd', 'Srbija',
 NOW() - INTERVAL '12 days', NOW(), 567),

-- Beauty services (using 1903 - Lepota i wellness)
(7, 1903, 'Masaža i wellness tretmani',
 'Profesionalne masaže: relaks, medicinska, sportska. Wellness tretmani za lice i telo. Sertifikovani terapeut.',
 40.00, 'active', 45.2485, 19.8335, 'Centar', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '15 days', NOW(), 432);

-- ============================================================================
-- SPORTS & RECREATION
-- ============================================================================

INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES
-- Fitness equipment
(7, 1010, 'Tegovi i oprema za teretanu - komplet',
 'Profesionalna oprema: bench klupa, stalak za čučanj, olimpijska šipka, 200kg tegova, bučice set. Odlično stanje.',
 1200.00, 'active', 45.2671, 19.8335, 'Veternik', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '9 days', NOW(), 678),

-- Bicycles
(7, 1010, 'Trek električni bicikl 2023',
 'Trek Powerfly 5 e-bike, Bosch motor 500W, baterija 625Wh. Domet do 120km. Shimano Deore oprema. Kao nov.',
 2500.00, 'active', 44.8125, 20.4612, 'Vračar', 'Beograd', 'Srbija',
 NOW() - INTERVAL '11 days', NOW(), 456);

-- ============================================================================
-- PETS & ANIMALS
-- ============================================================================

INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES
-- Dogs
(7, 1011, 'Golden Retriever štenci sa papirima',
 'Čistokrvni golden retriver štenci, 2 meseca stari. Čipovani, vakcinisani, sa pedigreom. Roditelji šampioni.',
 600.00, 'active', 45.2551, 19.8451, 'Rumenka', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '4 days', NOW(), 2345),

-- Cats
(7, 1011, 'Britanska kratkodlaka mačića',
 'Predivna britanska kratkodlaka mačića, plava boja. 3 meseca, naučeni na pesak. Očišćeni od parazita.',
 400.00, 'active', 44.7866, 20.4489, 'Zvezdara', 'Beograd', 'Srbija',
 NOW() - INTERVAL '6 days', NOW(), 1234);

-- ============================================================================
-- BOOKS & EDUCATION
-- ============================================================================

INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES
-- Books (assuming category)
(7, 1012, 'Komplet knjiga za medicinu - 50 naslova',
 'Medicinska literatura, anatomija, fiziologija, farmakologija. Sve knjige u odličnom stanju. Idealno za studente.',
 200.00, 'active', 45.2467, 19.8515, 'Medicinski fakultet okolina', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '13 days', NOW(), 345),

-- Musical instruments
(7, 1016, 'Yamaha gitara C40 sa futrolom',
 'Klasična gitara Yamaha C40, odličan izbor za početnike. Sa kvalitetnom futrolom i štimerom. Kao nova.',
 150.00, 'active', 44.8078, 20.4448, 'Stari grad', 'Beograd', 'Srbija',
 NOW() - INTERVAL '14 days', NOW(), 456);

-- ============================================================================
-- BABY & KIDS
-- ============================================================================

INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES
-- Baby equipment
(7, 1013, 'Chicco kolica 3u1 - kao nova',
 'Chicco Trio sistem: kolica, nosiljka, auto sedište. Korišćeno 6 meseci, u perfektnom stanju. Sa svim dodacima.',
 350.00, 'active', 45.2485, 19.8335, 'Novo naselje', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '7 days', NOW(), 567),

-- Toys
(7, 1013, 'LEGO kolekcija - 15 kompleta',
 'Velika kolekcija LEGO kockica: City, Technic, Star Wars. Svi kompleti kompletni sa uputstvima. Ukupno preko 10000 delova.',
 500.00, 'active', 44.8125, 20.4612, 'Banovo brdo', 'Beograd', 'Srbija',
 NOW() - INTERVAL '9 days', NOW(), 890);

-- Add images for all new listings (using same placeholder approach)
INSERT INTO marketplace_images (
  listing_id, file_name, file_path, file_size, content_type, is_main, created_at
)
SELECT
  id,
  'listing_' || id || '_main.jpg',
  'listings/' || id || '/main.jpg',
  1024000,
  'image/jpeg',
  true,
  NOW()
FROM marketplace_listings
WHERE created_at > NOW() - INTERVAL '16 days'
AND id NOT IN (SELECT DISTINCT listing_id FROM marketplace_images);

-- Add secondary images for some listings
INSERT INTO marketplace_images (
  listing_id, file_name, file_path, file_size, content_type, is_main, created_at
)
SELECT
  id,
  'listing_' || id || '_2.jpg',
  'listings/' || id || '/image2.jpg',
  1024000,
  'image/jpeg',
  false,
  NOW()
FROM marketplace_listings
WHERE created_at > NOW() - INTERVAL '10 days'
AND id NOT IN (SELECT DISTINCT listing_id FROM marketplace_images WHERE is_main = false)
LIMIT 10;

-- Output summary
SELECT
  'Added ' || COUNT(*) || ' more marketplace listings' as summary
FROM marketplace_listings
WHERE created_at > NOW() - INTERVAL '16 days';