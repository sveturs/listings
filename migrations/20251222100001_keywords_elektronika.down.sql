-- Rollback: Remove meta_keywords for all Elektronika subcategories (L2/L3)
-- Generated: 2025-12-22

-- Clear meta_keywords for all Elektronika subcategories
WITH RECURSIVE elektronika_tree AS (
    SELECT id FROM categories WHERE slug = 'elektronika'
    UNION ALL
    SELECT c.id FROM categories c
    JOIN elektronika_tree et ON c.parent_id = et.id
)
UPDATE categories
SET meta_keywords = '{}'::jsonb
WHERE id IN (SELECT id FROM elektronika_tree)
  AND slug != 'elektronika';

-- End of rollback migration
