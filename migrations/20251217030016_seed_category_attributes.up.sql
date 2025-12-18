-- Migration: 20251217030016_seed_category_attributes
-- Description: Seed category-attribute relationships (Phase 2, Task BE-2.7)
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17
-- Note: ~200 relationships between categories and attributes

-- ============================================================================
-- PART 1: GLOBAL ATTRIBUTES FOR ALL L1 CATEGORIES
-- ============================================================================

INSERT INTO category_attributes (
    category_id,
    attribute_id,
    is_enabled,
    is_required,
    is_searchable,
    is_filterable,
    sort_order
)
SELECT
    c.id AS category_id,
    a.id AS attribute_id,
    true AS is_enabled,
    false AS is_required,  -- Global attributes are not required by default
    a.is_searchable,
    a.is_filterable,
    a.sort_order
FROM categories c
CROSS JOIN attributes a
WHERE c.level = 1  -- All L1 categories
  AND a.code IN (
      'brand',
      'condition',
      'country_of_origin',
      'material',
      'weight',
      'dimensions',
      'warranty_months',
      'model_number',
      'year_of_manufacture',
      'energy_class'
  )
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ============================================================================
-- PART 2: CLOTHING ATTRIBUTES FOR "Odeća i obuća" (L1)
-- ============================================================================

INSERT INTO category_attributes (
    category_id,
    attribute_id,
    is_enabled,
    is_required,
    is_searchable,
    is_filterable,
    sort_order
)
SELECT
    c.id AS category_id,
    a.id AS attribute_id,
    true AS is_enabled,
    CASE
        WHEN a.code = 'clothing_size' THEN true  -- Size is required
        WHEN a.code = 'color' THEN true          -- Color is required
        ELSE false
    END AS is_required,
    a.is_searchable,
    a.is_filterable,
    CASE a.code
        WHEN 'clothing_size' THEN 1
        WHEN 'color' THEN 2
        WHEN 'gender' THEN 3
        WHEN 'fit' THEN 4
        WHEN 'style' THEN 5
        WHEN 'season' THEN 6
        WHEN 'neckline' THEN 7
        WHEN 'sleeve_length' THEN 8
        WHEN 'pattern' THEN 9
        WHEN 'closure_type' THEN 10
        ELSE 99
    END AS sort_order
FROM categories c
CROSS JOIN attributes a
WHERE c.slug = 'odeca-i-obuca'  -- L1 category for clothing
  AND a.code IN (
      'clothing_size',
      'color',
      'gender',
      'fit',
      'style',
      'season',
      'neckline',
      'sleeve_length',
      'pattern',
      'closure_type'
  )
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ============================================================================
-- PART 3: CLOTHING ATTRIBUTES FOR L2 CHILDREN OF "Odeća i obuća"
-- ============================================================================

INSERT INTO category_attributes (
    category_id,
    attribute_id,
    is_enabled,
    is_required,
    is_searchable,
    is_filterable,
    sort_order
)
SELECT DISTINCT
    c.id AS category_id,
    a.id AS attribute_id,
    true AS is_enabled,
    false AS is_required,  -- L2 can inherit requirements from L1
    a.is_searchable,
    a.is_filterable,
    0 AS sort_order
FROM categories c
CROSS JOIN attributes a
WHERE c.parent_id = (SELECT id FROM categories WHERE slug = 'odeca-i-obuca')
  AND c.level = 2
  AND a.code IN (
      'clothing_size',
      'color',
      'gender',
      'fit',
      'style',
      'season',
      'neckline',
      'sleeve_length'
  )
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ============================================================================
-- PART 4: ELECTRONICS ATTRIBUTES FOR "Elektronika" (L1)
-- ============================================================================

INSERT INTO category_attributes (
    category_id,
    attribute_id,
    is_enabled,
    is_required,
    is_searchable,
    is_filterable,
    sort_order
)
SELECT
    c.id AS category_id,
    a.id AS attribute_id,
    true AS is_enabled,
    false AS is_required,
    a.is_searchable,
    a.is_filterable,
    CASE a.code
        WHEN 'screen_size' THEN 1
        WHEN 'processor' THEN 2
        WHEN 'ram' THEN 3
        WHEN 'storage_capacity' THEN 4
        WHEN 'operating_system' THEN 5
        WHEN 'connectivity' THEN 6
        WHEN 'battery_capacity' THEN 7
        WHEN 'camera_resolution' THEN 8
        WHEN 'refresh_rate' THEN 9
        WHEN 'resolution' THEN 10
        ELSE 99
    END AS sort_order
FROM categories c
CROSS JOIN attributes a
WHERE c.slug = 'elektronika'  -- L1 category for electronics
  AND a.code IN (
      'screen_size',
      'processor',
      'ram',
      'storage_capacity',
      'operating_system',
      'connectivity',
      'battery_capacity',
      'camera_resolution',
      'refresh_rate',
      'resolution'
  )
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ============================================================================
-- PART 5: ELECTRONICS ATTRIBUTES FOR L2 CHILDREN OF "Elektronika"
-- ============================================================================

INSERT INTO category_attributes (
    category_id,
    attribute_id,
    is_enabled,
    is_required,
    is_searchable,
    is_filterable,
    sort_order
)
SELECT DISTINCT
    c.id AS category_id,
    a.id AS attribute_id,
    true AS is_enabled,
    false AS is_required,
    a.is_searchable,
    a.is_filterable,
    0 AS sort_order
FROM categories c
CROSS JOIN attributes a
WHERE c.parent_id = (SELECT id FROM categories WHERE slug = 'elektronika')
  AND c.level = 2
  AND a.code IN (
      'screen_size',
      'processor',
      'ram',
      'storage_capacity',
      'operating_system',
      'connectivity',
      'battery_capacity',
      'camera_resolution'
  )
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- ============================================================================
-- PROGRESS NOTIFICATION
-- ============================================================================

DO $$
DECLARE
    total_links INT;
    l1_links INT;
    clothing_links INT;
    electronics_links INT;
    l1_count INT;
BEGIN
    -- Count all links
    SELECT COUNT(*) INTO total_links FROM category_attributes;

    -- Count L1 global attribute links
    SELECT COUNT(*) INTO l1_links
    FROM category_attributes ca
    INNER JOIN categories c ON ca.category_id = c.id
    WHERE c.level = 1;

    -- Count L1 categories
    SELECT COUNT(*) INTO l1_count FROM categories WHERE level = 1;

    -- Count clothing links
    SELECT COUNT(*) INTO clothing_links
    FROM category_attributes ca
    INNER JOIN categories c ON ca.category_id = c.id
    INNER JOIN attributes a ON ca.attribute_id = a.id
    WHERE c.slug = 'odeca-i-obuca' OR c.parent_id = (SELECT id FROM categories WHERE slug = 'odeca-i-obuca');

    -- Count electronics links
    SELECT COUNT(*) INTO electronics_links
    FROM category_attributes ca
    INNER JOIN categories c ON ca.category_id = c.id
    INNER JOIN attributes a ON ca.attribute_id = a.id
    WHERE c.slug = 'elektronika' OR c.parent_id = (SELECT id FROM categories WHERE slug = 'elektronika');

    RAISE NOTICE '';
    RAISE NOTICE '✅ Category-Attribute links seed complete!';
    RAISE NOTICE '   Total links: %', total_links;
    RAISE NOTICE '   L1 categories with global attributes: % (% categories x 10 attributes = %)',
        l1_count, l1_count, l1_count * 10;
    RAISE NOTICE '   Clothing category links: %', clothing_links;
    RAISE NOTICE '   Electronics category links: %', electronics_links;
    RAISE NOTICE '';

    IF total_links < 100 THEN
        RAISE WARNING 'Expected at least 100 links, but found %. Some links may be missing.', total_links;
    END IF;
END $$;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
