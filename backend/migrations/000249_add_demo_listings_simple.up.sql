-- ============================================================================
-- Migration: Add demo listings for investor demonstration (simplified version)
-- Date: 2025-01-27
-- Purpose: Create realistic marketplace listings with existing categories
-- ============================================================================

-- First ensure we have a test user (ID 7)
INSERT INTO users (id, name, email, phone, created_at, updated_at)
VALUES (7, 'Demo User', 'demo@svetu.rs', '+381641234567', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  email = EXCLUDED.email,
  phone = EXCLUDED.phone;

-- ============================================================================
-- ELECTRONICS LISTINGS (category 1001 - Elektronika)
-- ============================================================================

-- iPhone 14 Pro Max (smartphones 1101)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1101, 'iPhone 14 Pro Max 256GB Deep Purple', 
 'Novi iPhone 14 Pro Max u fabrickoj foliji. Deep Purple boja, 256GB memorije. Kompletna garancija 2 godine. Moguca rate.',
 1200.00, 'active', 45.2551, 19.8451, 'Dunavska 15', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '3 days', NOW(), 2341);

-- Samsung Galaxy S23 Ultra 
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1101, 'Samsung Galaxy S23 Ultra 512GB', 
 'Samsung Galaxy S23 Ultra sa S Pen olovkom. 512GB memorije, 12GB RAM. Fantastican telefon za fotografiju. Kao nov.',
 950.00, 'active', 44.8125, 20.4612, 'Knez Mihailova 30', 'Beograd', 'Srbija',
 NOW() - INTERVAL '5 days', NOW(), 1567);

-- TV Samsung QLED (TV i audio 1103)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1103, 'Samsung QLED TV 65" 4K Smart TV', 
 'Samsung QLED televizor 65 inca sa 4K rezolucijom. Quantum Dot tehnologija, 120Hz, HDR10+ podrska. Smart TV sa Tizen OS.',
 1100.00, 'active', 45.2467, 19.8515, 'Bulevar Oslobodjenja 88', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '7 days', NOW() - INTERVAL '2 days', 876);

-- ============================================================================
-- FASHION LISTINGS (category 1002 - not shown, using subcategories)
-- ============================================================================

-- Muska odeca (1201)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1201, 'Hugo Boss odelo - tamno plavo', 
 'Elegantno Hugo Boss odelo, tamno plava boja. Velicina 50, slim fit kroj. Noseno nekoliko puta na svecanostima.',
 450.00, 'active', 44.7988, 20.4685, 'Bulevar Kralja Aleksandra 45', 'Beograd', 'Srbija',
 NOW() - INTERVAL '10 days', NOW() - INTERVAL '3 days', 543);

-- Cipele (1203)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1203, 'Nike Air Max 90 - original patike', 
 'Original Nike Air Max 90 patike, velicina 43. Kupljene u Sport Vision, sa racunom. Nosene svega par puta.',
 120.00, 'active', 45.2485, 19.8335, 'Futoski put 45', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '4 days', NOW(), 789);

-- Modni dodaci - torba (1204) 
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1204, 'Michael Kors torba original', 
 'Originalna Michael Kors torba, model Jet Set. Prirodna koza, crna boja. Sa dokumentima i dust bag. Odlicno stanje.',
 280.00, 'active', 44.8078, 20.4448, 'Knez Mihailova 15', 'Beograd', 'Srbija',
 NOW() - INTERVAL '2 days', NOW(), 890);

-- ============================================================================
-- HOME & GARDEN (category 1005)
-- ============================================================================

-- Alat za bastu (1502)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1502, 'Bosch kosilica za travu + trimer', 
 'Bosch elektricna kosilica ARM 37 i trimer ART 26. Odlicno stanje, redovno servisirana. Idealno za srednje baste.',
 250.00, 'active', 45.2671, 19.8335, 'Veternik', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '8 days', NOW() - INTERVAL '4 days', 234);

-- Dekoracija (1503)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1503, 'Persijski tepih 3x4m', 
 'Prekrasan rucno tkan persijski tepih. Dimenzije 3x4 metra, prirodne boje. Izuzetno ocuvan, iz nekadusackog doma.',
 800.00, 'active', 44.7866, 20.4489, 'Terazije 25', 'Beograd', 'Srbija',
 NOW() - INTERVAL '15 days', NOW() - INTERVAL '7 days', 345);

-- ============================================================================
-- AGRICULTURE (category 1006)
-- ============================================================================

-- Poljoprivredne masine (1601)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1601, 'Traktor IMT 539 sa prikljuccima', 
 'IMT 539 traktor u odlicnom stanju. Redovno odrzavan, nova guma. Uz traktor ide plug, drljaca i prskalica.',
 8500.00, 'active', 45.2551, 19.8451, 'Futoski put 120', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '12 days', NOW() - INTERVAL '5 days', 567);

-- Seme i djubriva (1602)  
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1602, 'Seme kukuruza - hibrid NS', 
 'Kvalitetno seme kukuruza, hibrid Novosadskog instituta. Visok prinos, otporan na susu. Pakovanje 25kg.',
 120.00, 'active', 45.2485, 19.8335, 'Rumenacka 15', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '6 days', NOW() - INTERVAL '1 day', 432);

-- Stoka (1603)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1603, 'Mlada junica simentalske rase', 
 'Prodajem mladu junicu simentalske rase, stara 18 meseci. Zdrava, redovno vakcinisana. Odlicna za dalji uzgoj.',
 1200.00, 'active', 45.2671, 19.8335, 'Futog', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '9 days', NOW() - INTERVAL '2 days', 234);

-- Poljoprivredni proizvodi (1604)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1604, 'Domaci med bagrem 1kg', 
 'Cist bagremov med sa sopstvene pceline. 100% prirodno, bez dodataka. Kristalno cist, svetle boje.',
 12.00, 'active', 45.2467, 19.8515, 'Sremski Karlovci', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '3 days', NOW(), 890);

-- ============================================================================
-- INDUSTRY (category 1007)
-- ============================================================================

-- Industrijske masine (1701)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1701, 'CNC glodalica 3-osna', 
 'Profesionalna CNC glodalica, radni prostor 800x600x200mm. Vreteno 2.2kW, vodeno hladjenje. Softer Mach3.',
 4500.00, 'active', 44.8125, 20.4612, 'Rakovica', 'Beograd', 'Srbija',
 NOW() - INTERVAL '14 days', NOW() - INTERVAL '6 days', 123);

-- Bezbednosna oprema (1704)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1704, 'Zastitna oprema komplet', 
 'Komplet zastitne opreme: kaciga, naocale, rukavice, cipele S3, prsluk. Sve novo, sa CE sertifikatima.',
 150.00, 'active', 45.2551, 19.8451, 'Industrijska zona', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '5 days', NOW() - INTERVAL '1 day', 456);

-- ============================================================================
-- FOOD & BEVERAGES (category 1008 not shown, using subcategory 1802)
-- ============================================================================

-- Pica (1802)
INSERT INTO marketplace_listings (
  user_id, category_id, title, description, price, status,
  latitude, longitude, location, address_city, address_country,
  created_at, updated_at, views_count
) VALUES 
(7, 1802, 'Domaca rakija sljiva 10L', 
 'Domaca sljivovica, duplo pecena. Prirodno voce, bez dodataka. Jacina 45%. Moguca degustacija.',
 35.00, 'active', 45.2485, 19.8335, 'Sremska Kamenica', 'Novi Sad', 'Srbija',
 NOW() - INTERVAL '7 days', NOW() - INTERVAL '3 days', 678);

-- Images and translations are skipped for simplicity
-- They would need separate handling based on actual DB structure

-- Output summary
SELECT 
  'Created ' || COUNT(*) || ' demo listings' as summary
FROM marketplace_listings 
WHERE user_id = 7;