-- Fix real estate category attributes
-- Update room attribute to have proper options
UPDATE category_attributes 
SET options = '{"values": ["garsonjera", "1", "1.5", "2", "2.5", "3", "3.5", "4", "4.5", "5+"]}'::jsonb
WHERE id = 2302;

-- Remove irrelevant attributes from main real estate category
DELETE FROM category_attribute_mapping 
WHERE category_id = 1004 
AND attribute_id IN (2002, 2003, 2004); -- condition, brand, color

-- Add proper real estate attributes to main category
INSERT INTO category_attribute_mapping (category_id, attribute_id) 
VALUES 
    (1004, 2301), -- area
    (1004, 2302), -- rooms
    (1004, 2303), -- floor
    (1004, 2304), -- furnished
    (1004, 2305), -- parking
    (1004, 2306)  -- balcony
ON CONFLICT DO NOTHING;