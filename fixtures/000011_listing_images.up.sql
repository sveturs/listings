-- Listing images fixture data

-- Clean up existing fixture images
DELETE FROM listing_images WHERE listing_id >= 1001 AND listing_id <= 2010;

-- C2C Listings images
INSERT INTO listing_images (listing_id, url, thumbnail_url, display_order, is_primary, mime_type)
VALUES
    -- MacBook Pro 14" M3 (1001)
    (1001, 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),
    (1001, 'https://images.unsplash.com/photo-1611186871348-b1ce696e52c9?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1611186871348-b1ce696e52c9?w=400&h=300&fit=crop', 1, false, 'image/jpeg'),

    -- Lenovo ThinkPad X1 Carbon (1002)
    (1002, 'https://images.unsplash.com/photo-1588872657578-7efd1f1555ed?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1588872657578-7efd1f1555ed?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- iPhone 15 Pro Max (1003)
    (1003, 'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),
    (1003, 'https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5?w=400&h=300&fit=crop', 1, false, 'image/jpeg'),

    -- Samsung Galaxy S24 Ultra (1004)
    (1004, 'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- IKEA Malm Desk (1005)
    (1005, 'https://images.unsplash.com/photo-1518455027359-f3f8164ba6bd?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1518455027359-f3f8164ba6bd?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- Leather Sofa 3-seater (1006)
    (1006, 'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),
    (1006, 'https://images.unsplash.com/photo-1493663284031-b7e3aefcae8e?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1493663284031-b7e3aefcae8e?w=400&h=300&fit=crop', 1, false, 'image/jpeg'),

    -- BMW E46 Headlights (1007)
    (1007, 'https://images.unsplash.com/photo-1489824904134-891ab64532f1?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1489824904134-891ab64532f1?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- Nike Air Max 90 (1008)
    (1008, 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),
    (1008, 'https://images.unsplash.com/photo-1600185365926-3a2ce3cdb9eb?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1600185365926-3a2ce3cdb9eb?w=400&h=300&fit=crop', 1, false, 'image/jpeg');

-- B2C Listings images
INSERT INTO listing_images (listing_id, url, thumbnail_url, display_order, is_primary, mime_type)
VALUES
    -- Dell XPS 15 (2001)
    (2001, 'https://images.unsplash.com/photo-1593642632559-0c6d3fc62b89?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1593642632559-0c6d3fc62b89?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),
    (2001, 'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=400&h=300&fit=crop', 1, false, 'image/jpeg'),

    -- Sony WH-1000XM5 (2002)
    (2002, 'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- iPad Pro 12.9" M4 (2003)
    (2003, 'https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- Office Chair Ergonomic (2004)
    (2004, 'https://images.unsplash.com/photo-1580480055273-228ff5388ef8?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1580480055273-228ff5388ef8?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- Garden Hose 30m (2005)
    (2005, 'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- Michelin Pilot Sport 4 (2006)
    (2006, 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- Bosch Brake Pads Set (2007)
    (2007, 'https://images.unsplash.com/photo-1486262715619-67b85e0b08d3?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1486262715619-67b85e0b08d3?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- Adidas Ultraboost 23 (2008)
    (2008, 'https://images.unsplash.com/photo-1608231387042-66d1773070a5?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1608231387042-66d1773070a5?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),
    (2008, 'https://images.unsplash.com/photo-1606107557195-0e29a4b5b4aa?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1606107557195-0e29a4b5b4aa?w=400&h=300&fit=crop', 1, false, 'image/jpeg'),

    -- Zara Summer Dress (2009)
    (2009, 'https://images.unsplash.com/photo-1595777457583-95e059d581b8?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1595777457583-95e059d581b8?w=400&h=300&fit=crop', 0, true, 'image/jpeg'),

    -- Levi's 501 Jeans (2010)
    (2010, 'https://images.unsplash.com/photo-1542272604-787c3835535d?w=800&h=600&fit=crop', 'https://images.unsplash.com/photo-1542272604-787c3835535d?w=400&h=300&fit=crop', 0, true, 'image/jpeg');
