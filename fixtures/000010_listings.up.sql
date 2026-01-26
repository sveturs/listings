-- Listings fixture data (C2C and B2C products)

-- Clean up existing fixture data first
DELETE FROM listings WHERE id >= 1001 AND id <= 2010;

-- C2C Listings (user-to-user)
INSERT INTO listings (id, uuid, user_id, storefront_id, title, description, price, currency, status, visibility, quantity, category_id, source_type, stock_status, original_language)
VALUES
    -- Electronics - Laptops
    (1001, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaa001'::uuid, 100, NULL, 'MacBook Pro 14" M3', 'Apple MacBook Pro 14 inch with M3 chip, 16GB RAM, 512GB SSD. Perfect condition, barely used.', 189990.00, 'RSD', 'active', 'public', 1, '11111111-1111-1111-1111-111111111002'::uuid, 'c2c', 'in_stock', 'en'),
    (1002, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaa002'::uuid, 101, NULL, 'Lenovo ThinkPad X1 Carbon', 'Business laptop, Intel i7, 16GB RAM. Great for work.', 129900.00, 'RSD', 'active', 'public', 1, '11111111-1111-1111-1111-111111111002'::uuid, 'c2c', 'in_stock', 'en'),

    -- Electronics - Smartphones
    (1003, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaa003'::uuid, 102, NULL, 'iPhone 15 Pro Max 256GB', 'Brand new iPhone 15 Pro Max, Natural Titanium color. Sealed box.', 179900.00, 'RSD', 'active', 'public', 1, '11111111-1111-1111-1111-111111111003'::uuid, 'c2c', 'in_stock', 'en'),
    (1004, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaa004'::uuid, 103, NULL, 'Samsung Galaxy S24 Ultra', 'Samsung flagship phone, 512GB, Titanium Black. With box and accessories.', 159900.00, 'RSD', 'active', 'public', 1, '11111111-1111-1111-1111-111111111003'::uuid, 'c2c', 'in_stock', 'en'),

    -- Home & Garden - Furniture
    (1005, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaa005'::uuid, 104, NULL, 'IKEA Malm Desk', 'White IKEA Malm desk, 140x65cm. Used for 1 year, excellent condition.', 15000.00, 'RSD', 'active', 'public', 1, '11111111-1111-1111-1111-111111112002'::uuid, 'c2c', 'in_stock', 'en'),
    (1006, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaa006'::uuid, 105, NULL, 'Leather Sofa 3-seater', 'Brown leather sofa, comfortable, minor wear. Must pick up.', 45000.00, 'RSD', 'active', 'public', 1, '11111111-1111-1111-1111-111111112002'::uuid, 'c2c', 'in_stock', 'sr'),

    -- Automotive - Car Parts
    (1007, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaa007'::uuid, 106, NULL, 'BMW E46 Headlights', 'Original BMW E46 angel eyes headlights. Left and right pair.', 25000.00, 'RSD', 'active', 'public', 2, '11111111-1111-1111-1111-111111113002'::uuid, 'c2c', 'in_stock', 'en'),

    -- Clothing
    (1008, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaa008'::uuid, 107, NULL, 'Nike Air Max 90', 'Nike Air Max 90, size 43, white/black. Worn twice.', 8500.00, 'RSD', 'active', 'public', 1, '11111111-1111-1111-1111-111111114002'::uuid, 'c2c', 'in_stock', 'en');

-- B2C Listings (store products)
INSERT INTO listings (id, uuid, user_id, storefront_id, title, description, price, currency, status, visibility, quantity, category_id, source_type, stock_status, sku, original_language)
VALUES
    -- Tech Store products
    (2001, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb001'::uuid, 1, 1001, 'Dell XPS 15 (2024)', 'Dell XPS 15 with Intel Core Ultra 7, 32GB RAM, 1TB SSD, OLED display.', 249900.00, 'RSD', 'active', 'public', 10, '11111111-1111-1111-1111-111111111002'::uuid, 'b2c', 'in_stock', 'TECH-DELL-XPS15-2024', 'en'),
    (2002, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb002'::uuid, 1, 1001, 'Sony WH-1000XM5', 'Sony wireless noise-cancelling headphones. Black color.', 49900.00, 'RSD', 'active', 'public', 25, '11111111-1111-1111-1111-111111111001'::uuid, 'b2c', 'in_stock', 'TECH-SONY-WH1000XM5', 'en'),
    (2003, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb003'::uuid, 1, 1001, 'iPad Pro 12.9" M4', 'Apple iPad Pro with M4 chip, 256GB, Space Black. Latest model.', 169900.00, 'RSD', 'active', 'public', 15, '11111111-1111-1111-1111-111111111004'::uuid, 'b2c', 'in_stock', 'TECH-IPAD-PRO-M4', 'en'),

    -- Home Depot products
    (2004, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb004'::uuid, 2, 1002, 'Office Chair Ergonomic', 'Ergonomic office chair with lumbar support. Mesh back, adjustable height.', 29900.00, 'RSD', 'active', 'public', 50, '11111111-1111-1111-1111-111111112002'::uuid, 'b2c', 'in_stock', 'HOME-CHAIR-ERGO-01', 'en'),
    (2005, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb005'::uuid, 2, 1002, 'Garden Hose 30m', 'Durable garden hose, 30 meters, with spray nozzle included.', 4500.00, 'RSD', 'active', 'public', 100, '11111111-1111-1111-1111-111111112003'::uuid, 'b2c', 'in_stock', 'HOME-HOSE-30M', 'en'),

    -- Auto Parts Plus products
    (2006, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb006'::uuid, 3, 1003, 'Michelin Pilot Sport 4', 'Michelin Pilot Sport 4 tires, 225/45 R17. Price per tire.', 18500.00, 'RSD', 'active', 'public', 40, '11111111-1111-1111-1111-111111113003'::uuid, 'b2c', 'in_stock', 'AUTO-MICH-PS4-22545', 'en'),
    (2007, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb007'::uuid, 3, 1003, 'Bosch Brake Pads Set', 'Bosch brake pads for VW Golf 7. Front axle set.', 7500.00, 'RSD', 'active', 'public', 30, '11111111-1111-1111-1111-111111113002'::uuid, 'b2c', 'in_stock', 'AUTO-BOSCH-BRAKE-VW7', 'en'),

    -- Fashion Hub products
    (2008, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb008'::uuid, 4, 1004, 'Adidas Ultraboost 23', 'Adidas Ultraboost running shoes. Multiple sizes available.', 21900.00, 'RSD', 'active', 'public', 60, '11111111-1111-1111-1111-111111114002'::uuid, 'b2c', 'in_stock', 'FASH-ADIDAS-UB23', 'en'),
    (2009, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb009'::uuid, 4, 1004, 'Zara Summer Dress', 'Elegant summer dress, floral pattern. Sizes S-XL.', 6900.00, 'RSD', 'active', 'public', 45, '11111111-1111-1111-1111-111111114003'::uuid, 'b2c', 'in_stock', 'FASH-ZARA-DRESS-01', 'en'),
    (2010, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb010'::uuid, 4, 1004, 'Levi''s 501 Jeans', 'Classic Levi''s 501 original fit jeans. Dark blue wash.', 12900.00, 'RSD', 'active', 'public', 80, '11111111-1111-1111-1111-111111114002'::uuid, 'b2c', 'in_stock', 'FASH-LEVIS-501-DB', 'en');

-- Update sequences to avoid conflicts
SELECT setval('listings_id_seq', (SELECT MAX(id) FROM listings));
