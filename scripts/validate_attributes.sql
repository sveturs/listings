-- ============================================================================
-- Validation Script for Attributes Migration
-- ============================================================================
-- Run on target database: listings_dev_db (localhost:35434)
-- ============================================================================

\echo '=== Attributes Migration Validation ==='
\echo ''

-- 1. Count comparison
\echo '1. Record count check:'
SELECT
    'Target DB (microservice)' as database,
    COUNT(*) as total_attributes
FROM attributes;

-- Expected: 203 (or more if some already existed)
\echo ''
\echo 'Expected: 203 records from source DB'
\echo ''

-- 2. JSONB structure validation
\echo '2. JSONB structure validation (first 10 records):'
SELECT
    id,
    code,
    name->>'en' as name_en,
    name->>'ru' as name_ru,
    name->>'sr' as name_sr,
    display_name->>'en' as display_name_en
FROM attributes
ORDER BY id
LIMIT 10;

\echo ''

-- 3. Check for NULL values in critical fields
\echo '3. NULL value check:'
SELECT
    'name has nulls' as check_name,
    COUNT(*) as count
FROM attributes
WHERE name IS NULL
   OR name->>'en' IS NULL
   OR name->>'ru' IS NULL
   OR name->>'sr' IS NULL
UNION ALL
SELECT
    'display_name has nulls',
    COUNT(*)
FROM attributes
WHERE display_name IS NULL
   OR display_name->>'en' IS NULL
   OR display_name->>'ru' IS NULL
   OR display_name->>'sr' IS NULL;

\echo 'Expected: 0 records with NULL values'
\echo ''

-- 4. Data type validation
\echo '4. Attribute type distribution:'
SELECT
    attribute_type,
    COUNT(*) as count
FROM attributes
GROUP BY attribute_type
ORDER BY count DESC;

\echo ''

-- 5. Purpose distribution
\echo '5. Purpose distribution:'
SELECT
    purpose,
    COUNT(*) as count
FROM attributes
GROUP BY purpose
ORDER BY count DESC;

\echo ''

-- 6. Boolean flags summary
\echo '6. Boolean flags summary:'
SELECT
    'is_searchable' as flag,
    COUNT(CASE WHEN is_searchable THEN 1 END) as true_count,
    COUNT(CASE WHEN NOT is_searchable THEN 1 END) as false_count
FROM attributes
UNION ALL
SELECT 'is_filterable',
    COUNT(CASE WHEN is_filterable THEN 1 END),
    COUNT(CASE WHEN NOT is_filterable THEN 1 END)
FROM attributes
UNION ALL
SELECT 'is_required',
    COUNT(CASE WHEN is_required THEN 1 END),
    COUNT(CASE WHEN NOT is_required THEN 1 END)
FROM attributes
UNION ALL
SELECT 'is_active',
    COUNT(CASE WHEN is_active THEN 1 END),
    COUNT(CASE WHEN NOT is_active THEN 1 END)
FROM attributes;

\echo ''

-- 7. Legacy ID mappings
\echo '7. Legacy ID mappings:'
SELECT
    COUNT(CASE WHEN legacy_category_attribute_id IS NOT NULL THEN 1 END) as has_legacy_category_id,
    COUNT(CASE WHEN legacy_product_variant_attribute_id IS NOT NULL THEN 1 END) as has_legacy_variant_id
FROM attributes;

\echo ''

-- 8. Check for duplicate codes (should be 0)
\echo '8. Duplicate code check:'
SELECT
    code,
    COUNT(*) as occurrences
FROM attributes
GROUP BY code
HAVING COUNT(*) > 1;

\echo 'Expected: 0 duplicate codes'
\echo ''

-- 9. ID sequence validation
\echo '9. Sequence validation:'
SELECT
    'attributes_id_seq' as sequence_name,
    last_value as current_value,
    (SELECT MAX(id) FROM attributes) as max_id_in_table
FROM attributes_id_seq;

\echo 'Current value should be >= max_id_in_table'
\echo ''

-- 10. Search vector check
\echo '10. Search vector check:'
SELECT
    COUNT(CASE WHEN search_vector IS NOT NULL THEN 1 END) as has_search_vector,
    COUNT(CASE WHEN search_vector IS NULL THEN 1 END) as missing_search_vector
FROM attributes;

\echo ''
\echo '=== Validation Complete ==='
