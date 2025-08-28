-- Rollback real estate category attributes changes
-- Restore empty options for rooms attribute
UPDATE category_attributes 
SET options = '{}'::jsonb
WHERE id = 2302;

-- Remove real estate attributes from main category
DELETE FROM category_attribute_mapping 
WHERE category_id = 1004 
AND attribute_id IN (2301, 2302, 2303, 2304, 2305, 2306);

-- Restore original attributes
INSERT INTO category_attribute_mapping (category_id, attribute_id) 
VALUES 
    (1004, 2002), -- condition
    (1004, 2003), -- brand
    (1004, 2004)  -- color
ON CONFLICT DO NOTHING;