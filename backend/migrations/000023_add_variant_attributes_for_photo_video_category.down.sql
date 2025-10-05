-- Rollback variant attributes for photo-video category

DELETE FROM variant_attribute_mappings
WHERE category_id = 1106
  AND variant_attribute_id IN (144, 146, 112);
