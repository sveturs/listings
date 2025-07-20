-- Remove Serbia marketplace categories fixtures

-- Remove translations
DELETE FROM translations WHERE entity_type = 'category' AND entity_id >= 1001 AND entity_id <= 2004;

-- Remove categories (will cascade to children due to foreign key constraints)
DELETE FROM marketplace_categories WHERE id >= 1001;

-- Reset sequence
SELECT setval('marketplace_categories_id_seq', (SELECT COALESCE(MAX(id), 1) FROM marketplace_categories), true);