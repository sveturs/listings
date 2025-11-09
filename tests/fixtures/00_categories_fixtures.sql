-- Categories fixture for integration tests
-- This file provides minimal category data required by test listings
--
-- IMPORTANT: This file should be loaded FIRST (prefix 00_) before other fixtures
-- because listings table has FK constraint to c2c_categories

BEGIN;

-- Root category: Electronics
INSERT INTO c2c_categories (id, name, slug, parent_id, level, is_active, sort_order)
VALUES (1300, 'Electronics', 'electronics', NULL, 0, TRUE, 1);

-- Sub-category: Laptops (most commonly used in tests)
INSERT INTO c2c_categories (id, name, slug, parent_id, level, is_active, sort_order)
VALUES (1301, 'Laptops', 'laptops', 1300, 1, TRUE, 1);

-- Additional sub-categories for comprehensive testing
INSERT INTO c2c_categories (id, name, slug, parent_id, level, is_active, sort_order)
VALUES
    (1302, 'Smartphones', 'smartphones', 1300, 1, TRUE, 2),
    (1303, 'Tablets', 'tablets', 1300, 1, TRUE, 3),
    (1304, 'Desktop Computers', 'desktop-computers', 1300, 1, TRUE, 4),
    (1305, 'Computer Accessories', 'computer-accessories', 1300, 1, TRUE, 5);

-- Root category: Home & Garden
INSERT INTO c2c_categories (id, name, slug, parent_id, level, is_active, sort_order)
VALUES (1400, 'Home & Garden', 'home-garden', NULL, 0, TRUE, 2);

-- Sub-category: Furniture
INSERT INTO c2c_categories (id, name, slug, parent_id, level, is_active, sort_order)
VALUES (1401, 'Furniture', 'furniture', 1400, 1, TRUE, 1);

-- Root category: Automotive
INSERT INTO c2c_categories (id, name, slug, parent_id, level, is_active, sort_order)
VALUES (1500, 'Automotive', 'automotive', NULL, 0, TRUE, 3);

-- Sub-category: Car Parts
INSERT INTO c2c_categories (id, name, slug, parent_id, level, is_active, sort_order)
VALUES (1501, 'Car Parts', 'car-parts', 1500, 1, TRUE, 1);

COMMIT;

-- Verification query (for manual testing):
-- SELECT id, name, parent_id, level FROM c2c_categories ORDER BY id;
--
-- Expected result:
-- | id   | name                  | parent_id | level |
-- |------|-----------------------|-----------|-------|
-- | 1300 | Electronics           | NULL      | 0     |
-- | 1301 | Laptops               | 1300      | 1     |
-- | 1302 | Smartphones           | 1300      | 1     |
-- | 1303 | Tablets               | 1300      | 1     |
-- | 1304 | Desktop Computers     | 1300      | 1     |
-- | 1305 | Computer Accessories  | 1300      | 1     |
-- | 1400 | Home & Garden         | NULL      | 0     |
-- | 1401 | Furniture             | 1400      | 1     |
-- | 1500 | Automotive            | NULL      | 0     |
-- | 1501 | Car Parts             | 1500      | 1     |
