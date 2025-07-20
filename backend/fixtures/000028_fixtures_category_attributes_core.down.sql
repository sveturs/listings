-- Remove core category attributes fixtures

-- Remove attribute translations
DELETE FROM translations WHERE entity_type = 'attribute' AND entity_id >= 2001 AND entity_id <= 3000;

-- Remove attribute option translations
DELETE FROM attribute_option_translations WHERE attribute_id >= 2001 AND attribute_id <= 3000;

-- Remove category attributes
DELETE FROM category_attributes WHERE id >= 2001 AND id <= 3000;

-- Reset sequences
SELECT setval('category_attributes_id_seq', (SELECT COALESCE(MAX(id), 1) FROM category_attributes), true);
SELECT setval('attribute_option_translations_id_seq', (SELECT COALESCE(MAX(id), 1) FROM attribute_option_translations), true);