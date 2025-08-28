-- Remove translations added for categories 2006 and 2007
DELETE FROM translations WHERE entity_type = 'category' AND entity_id IN (2006, 2007);

-- Restore original category names
UPDATE marketplace_categories SET name = 'photo' WHERE id = 2006;
UPDATE marketplace_categories SET name = 'wifi-routery' WHERE id = 2007;