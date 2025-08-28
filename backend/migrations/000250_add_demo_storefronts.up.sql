-- ============================================================================
-- Migration: Add demo storefronts for investor demonstration
-- Date: 2025-01-27
-- Purpose: Create realistic storefronts and products
-- ============================================================================

-- ============================================================================
-- STOREFRONT 1: TechNova Electronics Store
-- ============================================================================
INSERT INTO storefronts (
  user_id, slug, name, description, 
  phone, email, website, address, city, country,
  latitude, longitude, is_active, is_verified,
  rating, reviews_count, products_count, sales_count, views_count,
  created_at, updated_at
) VALUES (
  7, 'technova-electronics', 'TechNova Electronics', 
  'Vodeći prodavac elektronike i tehnike u Novom Sadu. Nudimo širok asortiman proizvoda od pametnih telefona do kućnih aparata.',
  '+381214567890', 'info@technova.rs', 'https://technova.rs',
  'Bulevar Oslobodjenja 125', 'Novi Sad', 'RS',
  45.2551, 19.8451, true, true,
  4.8, 234, 0, 567, 8934,
  NOW() - INTERVAL '45 days', NOW()
);

-- Get the storefront ID for products
DO $$
DECLARE 
  technova_id INTEGER;
  fashion_house_id INTEGER;
  home_garden_id INTEGER;
  agro_shop_id INTEGER;
BEGIN
  SELECT id INTO technova_id FROM storefronts WHERE slug = 'technova-electronics';
  
  -- Products for TechNova
  INSERT INTO storefront_products (
    storefront_id, name, description, price, currency, category_id, 
    sku, stock_quantity, stock_status, is_active,
    created_at, updated_at
  ) VALUES 
  -- Smartphones
  (technova_id, 'iPhone 15 Pro Max 256GB', 
   'Najnoviji iPhone sa titanijumskim kućištem i A17 Pro čipom. Garancija 2 godine.',
   1499.00, 'EUR', 1101, 'IPH15PM256', 5, 'in_stock', true,
   NOW() - INTERVAL '10 days', NOW()),
  
  (technova_id, 'Samsung Galaxy S24 Ultra', 
   'Flagship Samsung telefon sa S Pen, 200MP kamerom i AI funkcionalnostima.',
   1299.00, 'EUR', 1101, 'SGS24U512', 8, 'in_stock', true,
   NOW() - INTERVAL '9 days', NOW()),
  
  (technova_id, 'Google Pixel 8 Pro', 
   'Google telefon sa najboljom kamerom i čistim Android iskustvom.',
   999.00, 'EUR', 1101, 'GP8P256', 3, 'in_stock', true,
   NOW() - INTERVAL '8 days', NOW()),
  
  -- Laptops (assuming category exists or using general electronics)
  (technova_id, 'MacBook Pro 14" M3', 
   'Apple laptop sa M3 čipom, 16GB RAM, 512GB SSD. Idealan za profesionalce.',
   2199.00, 'EUR', 1001, 'MBP14M3', 4, 'in_stock', true,
   NOW() - INTERVAL '7 days', NOW()),
  
  (technova_id, 'Dell XPS 15', 
   'Premium Windows laptop sa 4K OLED ekranom i Intel Core i9 procesorom.',
   1899.00, 'EUR', 1001, 'DXPS15I9', 2, 'in_stock', true,
   NOW() - INTERVAL '6 days', NOW()),
  
  -- TVs
  (technova_id, 'LG OLED 65" C3', 
   'OLED TV sa 4K rezolucijom, 120Hz, HDMI 2.1 za gaming. Smart TV sa webOS.',
   1799.00, 'EUR', 1103, 'LGC365', 6, 'in_stock', true,
   NOW() - INTERVAL '5 days', NOW()),
  
  (technova_id, 'Sony Bravia XR 55"', 
   'Sony TV sa Cognitive Processor XR i Google TV. Dolby Vision i Atmos.',
   1299.00, 'EUR', 1103, 'SNYXR55', 4, 'in_stock', true,
   NOW() - INTERVAL '4 days', NOW());

  -- Update products count
  UPDATE storefronts SET products_count = 7 WHERE id = technova_id;
END $$;

-- ============================================================================
-- STOREFRONT 2: Fashion House Belgrade
-- ============================================================================
INSERT INTO storefronts (
  user_id, slug, name, description,
  phone, email, address, city, country,
  latitude, longitude, is_active, is_verified,
  rating, reviews_count, products_count, sales_count, views_count,
  created_at, updated_at
) VALUES (
  7, 'fashion-house-belgrade', 'Fashion House Belgrade',
  'Ekskluzivna modna kuća sa brendiranom odećom i aksesoarema. Zastupamo svetske brendove.',
  '+381113456789', 'contact@fashionhouse.rs',
  'Knez Mihailova 28', 'Beograd', 'RS',
  44.8125, 20.4612, true, true,
  4.7, 456, 0, 892, 15234,
  NOW() - INTERVAL '60 days', NOW()
);

DO $$
DECLARE
  fashion_id INTEGER;
BEGIN
  SELECT id INTO fashion_id FROM storefronts WHERE slug = 'fashion-house-belgrade';
  
  INSERT INTO storefront_products (
    storefront_id, name, description, price, currency, category_id,
    sku, stock_quantity, stock_status, is_active,
    created_at, updated_at
  ) VALUES
  -- Men's clothing
  (fashion_id, 'Tommy Hilfiger Polo majica',
   'Originalna Tommy Hilfiger polo majica, 100% pamuk. Dostupne sve veličine.',
   89.00, 'EUR', 1201, 'THPOLO01', 15, 'in_stock', true,
   NOW() - INTERVAL '12 days', NOW()),
  
  (fashion_id, 'Calvin Klein džins',
   'CK Jeans slim fit model. Premium denim, udoban i moderan kroj.',
   129.00, 'EUR', 1201, 'CKJEANS01', 20, 'in_stock', true,
   NOW() - INTERVAL '11 days', NOW()),
  
  -- Women's clothing (assuming category exists)
  (fashion_id, 'Zara haljina',
   'Elegantna večernja haljina, crna boja. Idealna za svečane prilike.',
   149.00, 'EUR', 1202, 'ZRDRESS01', 8, 'in_stock', true,
   NOW() - INTERVAL '10 days', NOW()),
  
  -- Shoes
  (fashion_id, 'Adidas Ultraboost 22',
   'Trčanje patike sa Boost đonom. Maksimalna udobnost i performanse.',
   179.00, 'EUR', 1203, 'ADUB22', 12, 'in_stock', true,
   NOW() - INTERVAL '9 days', NOW()),
  
  (fashion_id, 'Nike Air Force 1',
   'Klasične Nike patike, bele boje. Ikonski model koji ne izlazi iz mode.',
   139.00, 'EUR', 1203, 'NAF1W', 25, 'in_stock', true,
   NOW() - INTERVAL '8 days', NOW()),
  
  -- Accessories
  (fashion_id, 'Guess sat za dame',
   'Elegantan ženski sat sa kristalima. Rose gold boja, vodootporan.',
   289.00, 'EUR', 1204, 'GUESSW01', 5, 'in_stock', true,
   NOW() - INTERVAL '7 days', NOW()),
  
  (fashion_id, 'Ray-Ban Aviator naočare',
   'Originalne Ray-Ban sunčane naočare. Klasični aviator model.',
   189.00, 'EUR', 1204, 'RBAV01', 10, 'in_stock', true,
   NOW() - INTERVAL '6 days', NOW());
  
  UPDATE storefronts SET products_count = 7 WHERE id = fashion_id;
END $$;

-- ============================================================================
-- STOREFRONT 3: Home & Garden Center
-- ============================================================================
INSERT INTO storefronts (
  user_id, slug, name, description,
  phone, email, address, city, country,
  latitude, longitude, is_active, is_verified,
  rating, reviews_count, products_count, sales_count, views_count,
  created_at, updated_at
) VALUES (
  7, 'home-garden-center', 'Home & Garden Center',
  'Sve za vaš dom i baštu. Od nameštaja do baštenskog alata.',
  '+381214445556', 'info@homegarden.rs',
  'Rumenački put 88', 'Novi Sad', 'RS',
  45.2671, 19.8335, true, true,
  4.6, 234, 0, 456, 7890,
  NOW() - INTERVAL '90 days', NOW()
);

DO $$
DECLARE
  home_id INTEGER;
BEGIN
  SELECT id INTO home_id FROM storefronts WHERE slug = 'home-garden-center';
  
  INSERT INTO storefront_products (
    storefront_id, name, description, price, currency, category_id,
    sku, stock_quantity, stock_status, is_active,
    created_at, updated_at
  ) VALUES
  -- Garden tools
  (home_id, 'Bosch električna kosačica',
   'Bosch ARM 34 kosačica, 1300W. Širina košenja 34cm, zapremina korpe 40L.',
   189.00, 'EUR', 1502, 'BARM34', 8, 'in_stock', true,
   NOW() - INTERVAL '14 days', NOW()),
  
  (home_id, 'Gardena set za navodnjavanje',
   'Komplet sistem za automatsko zalivanje bašte. Tajmer i prskalice.',
   149.00, 'EUR', 1502, 'GRDIRR01', 12, 'in_stock', true,
   NOW() - INTERVAL '13 days', NOW()),
  
  -- Furniture (assuming category exists)
  (home_id, 'IKEA trosed KIVIK',
   'Udoban trosed sa mogućnošću skidanja navlake. Siva boja.',
   599.00, 'EUR', 1501, 'IKKIVIK', 3, 'in_stock', true,
   NOW() - INTERVAL '12 days', NOW()),
  
  (home_id, 'Simpo radni sto',
   'Radni sto 140x60cm sa fiokama. Hrast dekor, metal noge.',
   249.00, 'EUR', 1501, 'SMPDESK01', 6, 'in_stock', true,
   NOW() - INTERVAL '11 days', NOW()),
  
  -- Decoration
  (home_id, 'LED dekorativna rasveta',
   'Smart LED traka 5m sa WiFi kontrolom. RGB boje, aplikacija za telefon.',
   79.00, 'EUR', 1503, 'LEDRGB5M', 20, 'in_stock', true,
   NOW() - INTERVAL '10 days', NOW()),
  
  (home_id, 'Zidni sat vintage',
   'Retro dizajn zidni sat, prečnik 60cm. Rimski brojevi, tihi mehanizam.',
   89.00, 'EUR', 1503, 'CLKVINT60', 15, 'in_stock', true,
   NOW() - INTERVAL '9 days', NOW());
  
  UPDATE storefronts SET products_count = 6 WHERE id = home_id;
END $$;

-- ============================================================================
-- STOREFRONT 4: AgroShop - Agricultural Supplies
-- ============================================================================
INSERT INTO storefronts (
  user_id, slug, name, description,
  phone, email, website, address, city, country,
  latitude, longitude, is_active, is_verified,
  rating, reviews_count, products_count, sales_count, views_count,
  created_at, updated_at
) VALUES (
  7, 'agroshop-supplies', 'AgroShop',
  'Poljoprivredna oprema, semena, đubriva i zaštitna sredstva. Partner farmera već 20 godina.',
  '+381215556667', 'prodaja@agroshop.rs', 'https://agroshop.rs',
  'Industrijska zona bb', 'Novi Sad', 'RS',
  45.2485, 19.8335, true, true,
  4.9, 567, 0, 1234, 23456,
  NOW() - INTERVAL '180 days', NOW()
);

DO $$
DECLARE
  agro_id INTEGER;
BEGIN
  SELECT id INTO agro_id FROM storefronts WHERE slug = 'agroshop-supplies';
  
  INSERT INTO storefront_products (
    storefront_id, name, description, price, currency, category_id,
    sku, stock_quantity, stock_status, is_active,
    created_at, updated_at
  ) VALUES
  -- Seeds
  (agro_id, 'Seme kukuruza NS hibrid',
   'Visokoprinosni hibrid kukuruza FAO 400. Otporan na sušu. Pakovanje 25kg.',
   145.00, 'EUR', 1602, 'NSCRN400', 100, 'in_stock', true,
   NOW() - INTERVAL '20 days', NOW()),
  
  (agro_id, 'Seme pšenice NS 40S',
   'Seme ozime pšenice, visok sadržaj proteina. Sertifikovano. Vreća 50kg.',
   85.00, 'EUR', 1602, 'NSWHT40S', 150, 'in_stock', true,
   NOW() - INTERVAL '19 days', NOW()),
  
  -- Fertilizers
  (agro_id, 'NPK đubrivo 15-15-15',
   'Mineralno đubrivo sa balansiranim sadržajem azota, fosfora i kalijuma. 25kg.',
   55.00, 'EUR', 1602, 'NPK151515', 200, 'in_stock', true,
   NOW() - INTERVAL '18 days', NOW()),
  
  (agro_id, 'Urea 46% N',
   'Azotno đubrivo sa visokim sadržajem azota. Za prihranu useva. Vreća 50kg.',
   75.00, 'EUR', 1602, 'UREA46', 180, 'in_stock', true,
   NOW() - INTERVAL '17 days', NOW()),
  
  -- Agricultural machinery parts
  (agro_id, 'Plug za traktor IMT',
   'Dvobrazdni obrtni plug za IMT traktore. Kvalitetna domaća izrada.',
   890.00, 'EUR', 1601, 'PLUGIMT2', 5, 'in_stock', true,
   NOW() - INTERVAL '16 days', NOW()),
  
  (agro_id, 'Prskalica 400L',
   'Traktorska prskalica 400 litara sa krilima 12m. Nova, garancija.',
   1450.00, 'EUR', 1601, 'SPRAY400', 3, 'in_stock', true,
   NOW() - INTERVAL '15 days', NOW()),
  
  -- Agricultural products
  (agro_id, 'Bagremov med 3kg',
   'Prirodni bagremov med sa našeg pčelinjaka. Svetle boje, blag ukus.',
   35.00, 'EUR', 1604, 'HONEY3KG', 50, 'in_stock', true,
   NOW() - INTERVAL '14 days', NOW()),
  
  (agro_id, 'Domaća rakija šljiva 1L',
   'Tradicionalna šljivovica, duplo pečena. Jačina 45%. Sa sertifikatom.',
   25.00, 'EUR', 1802, 'RAKSLJV1L', 80, 'in_stock', true,
   NOW() - INTERVAL '13 days', NOW());
  
  UPDATE storefronts SET products_count = 8 WHERE id = agro_id;
END $$;

-- Add images to storefront products (simplified - using placeholders)
INSERT INTO storefront_product_images (
  storefront_product_id, image_url, thumbnail_url, display_order, is_default, created_at
)
SELECT 
  sp.id,
  'https://picsum.photos/800/600?random=' || sp.id,
  'https://picsum.photos/200/150?random=' || sp.id,
  1,
  true,
  NOW()
FROM storefront_products sp
WHERE sp.created_at > NOW() - INTERVAL '30 days';

-- Output summary
SELECT 
  'Created ' || COUNT(DISTINCT s.id) || ' storefronts with ' || 
  COUNT(DISTINCT sp.id) || ' products' as summary
FROM storefronts s
LEFT JOIN storefront_products sp ON s.id = sp.storefront_id
WHERE s.created_at > NOW() - INTERVAL '200 days';