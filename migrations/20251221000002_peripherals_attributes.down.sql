-- ============================================================================
-- ROLLBACK PERIPHERALS ATTRIBUTES
-- ============================================================================

-- Step 1: Remove category_attributes links for new peripheral attributes
DELETE FROM category_attributes
WHERE category_id = '2c1f8391-95e3-4a37-baec-f923b4c9e5a1'
  AND attribute_id IN (101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112);

-- Step 2: Remove the new peripheral attributes
DELETE FROM attributes WHERE id IN (101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112);

-- Step 3: Restore original (incorrect) smartphone attributes to Peripherals
-- (This restores the original state before this migration)
INSERT INTO category_attributes (category_id, attribute_id, is_enabled, is_required, is_filterable, sort_order)
VALUES
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 21, true, false, true, 0),  -- screen_size
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 22, true, false, false, 0), -- processor
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 23, true, false, true, 0),  -- ram
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 24, true, false, true, 0),  -- storage_capacity
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 25, true, false, true, 0),  -- operating_system
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 27, true, false, true, 0),  -- battery_capacity
    ('2c1f8391-95e3-4a37-baec-f923b4c9e5a1', 28, true, false, true, 0)   -- camera_resolution
ON CONFLICT (category_id, attribute_id) DO NOTHING;
