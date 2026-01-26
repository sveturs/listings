-- Storefronts fixture data (for B2C listings)

-- Clean up existing fixture data first
DELETE FROM storefronts WHERE id IN (1001, 1002, 1003, 1004);

INSERT INTO storefronts (id, user_id, slug, name, description, phone, email, city, country, is_active, is_verified, legal_entity_type, business_category)
VALUES
    (1001, 1, 'tech-store', 'Tech Store', 'Your favorite electronics store', '+381111234567', 'info@techstore.rs', 'Belgrade', 'RS', true, true, 'doo', 'retail'),
    (1002, 2, 'home-depot', 'Home Depot Serbia', 'Everything for your home', '+381112345678', 'info@homedepot.rs', 'Novi Sad', 'RS', true, true, 'doo', 'retail'),
    (1003, 3, 'auto-parts-plus', 'Auto Parts Plus', 'Quality car parts', '+381113456789', 'sales@autoparts.rs', 'Nis', 'RS', true, false, 'preduzetnik', 'retail'),
    (1004, 4, 'fashion-hub', 'Fashion Hub', 'Latest fashion trends', '+381114567890', 'hello@fashionhub.rs', 'Belgrade', 'RS', true, true, 'doo', 'retail');

-- Update sequence to avoid conflicts
SELECT setval('storefronts_id_seq', (SELECT MAX(id) FROM storefronts));
