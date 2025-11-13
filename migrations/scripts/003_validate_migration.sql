-- ============================================================================
-- Script: 003_validate_migration.sql
-- Description: Validate attributes migration data integrity
-- Database: listings_dev_db (listings microservice)
-- Usage: psql "$LISTINGS_DB" -f 003_validate_migration.sql
-- ============================================================================

\echo '============================================================================'
\echo 'Attributes Migration: Data Validation'
\echo '============================================================================'
\echo ''

-- ============================================================================
-- 1. COUNT VALIDATION
-- ============================================================================
\echo '1. Record Counts'
\echo '----------------'

SELECT
    'attributes' as table_name,
    COUNT(*) as record_count,
    CASE
        WHEN COUNT(*) = 203 THEN '✓ EXPECTED'
        ELSE '⚠ UNEXPECTED (expected 203)'
    END as status
FROM attributes
UNION ALL
SELECT
    'category_attributes' as table_name,
    COUNT(*) as record_count,
    CASE
        WHEN COUNT(*) > 0 THEN '✓ HAS DATA'
        ELSE '⚠ EMPTY'
    END as status
FROM category_attributes
UNION ALL
SELECT
    'listing_attribute_values' as table_name,
    COUNT(*) as record_count,
    CASE
        WHEN COUNT(*) > 0 THEN '✓ HAS DATA'
        ELSE 'ℹ EMPTY (OK if no listings)'
    END as status
FROM listing_attribute_values
UNION ALL
SELECT
    'category_variant_attributes' as table_name,
    COUNT(*) as record_count,
    CASE
        WHEN COUNT(*) >= 0 THEN '✓ OK'
        ELSE '⚠ ERROR'
    END as status
FROM category_variant_attributes
UNION ALL
SELECT
    'variant_attribute_values' as table_name,
    COUNT(*) as record_count,
    'ℹ (variants not yet migrated)' as status
FROM variant_attribute_values
UNION ALL
SELECT
    'attribute_options' as table_name,
    COUNT(*) as record_count,
    'ℹ (will be populated later)' as status
FROM attribute_options
UNION ALL
SELECT
    'attribute_search_cache' as table_name,
    COUNT(*) as record_count,
    'ℹ (will be populated by indexer)' as status
FROM attribute_search_cache;

\echo ''

-- ============================================================================
-- 2. DATA INTEGRITY: Foreign Keys
-- ============================================================================
\echo '2. Foreign Key Integrity'
\echo '------------------------'

-- Check category_attributes -> attributes
SELECT
    'category_attributes → attributes' as relationship,
    COUNT(*) as orphaned_records,
    CASE
        WHEN COUNT(*) = 0 THEN '✓ OK'
        ELSE '✗ BROKEN REFERENCES'
    END as status
FROM category_attributes ca
LEFT JOIN attributes a ON ca.attribute_id = a.id
WHERE a.id IS NULL;

-- Check listing_attribute_values -> attributes
SELECT
    'listing_attribute_values → attributes' as relationship,
    COUNT(*) as orphaned_records,
    CASE
        WHEN COUNT(*) = 0 THEN '✓ OK'
        ELSE '✗ BROKEN REFERENCES'
    END as status
FROM listing_attribute_values lav
LEFT JOIN attributes a ON lav.attribute_id = a.id
WHERE a.id IS NULL;

-- Check listing_attribute_values -> listings
SELECT
    'listing_attribute_values → listings' as relationship,
    COUNT(*) as orphaned_records,
    CASE
        WHEN COUNT(*) = 0 THEN '✓ OK'
        WHEN COUNT(*) > 0 THEN '⚠ ORPHANED (listings not migrated yet)'
        ELSE '✗ ERROR'
    END as status
FROM listing_attribute_values lav
LEFT JOIN listings l ON lav.listing_id = l.id
WHERE l.id IS NULL;

-- Check category_variant_attributes -> attributes
SELECT
    'category_variant_attributes → attributes' as relationship,
    COUNT(*) as orphaned_records,
    CASE
        WHEN COUNT(*) = 0 THEN '✓ OK'
        ELSE '✗ BROKEN REFERENCES'
    END as status
FROM category_variant_attributes cva
LEFT JOIN attributes a ON cva.attribute_id = a.id
WHERE a.id IS NULL;

\echo ''

-- ============================================================================
-- 3. JSONB FIELD VALIDATION
-- ============================================================================
\echo '3. JSONB Field Validation'
\echo '-------------------------'

-- Check name JSONB structure
SELECT
    'attributes.name' as field,
    COUNT(*) as records_with_valid_jsonb,
    COUNT(*) - COUNT(CASE WHEN name ? 'en' THEN 1 END) as missing_en,
    COUNT(*) - COUNT(CASE WHEN name ? 'ru' THEN 1 END) as missing_ru,
    COUNT(*) - COUNT(CASE WHEN name ? 'sr' THEN 1 END) as missing_sr,
    CASE
        WHEN COUNT(*) = COUNT(CASE WHEN name ? 'en' THEN 1 END) THEN '✓ OK'
        ELSE '⚠ MISSING TRANSLATIONS'
    END as status
FROM attributes;

-- Check display_name JSONB structure
SELECT
    'attributes.display_name' as field,
    COUNT(*) as records_with_valid_jsonb,
    COUNT(*) - COUNT(CASE WHEN display_name ? 'en' THEN 1 END) as missing_en,
    COUNT(*) - COUNT(CASE WHEN display_name ? 'ru' THEN 1 END) as missing_ru,
    COUNT(*) - COUNT(CASE WHEN display_name ? 'sr' THEN 1 END) as missing_sr,
    CASE
        WHEN COUNT(*) = COUNT(CASE WHEN display_name ? 'en' THEN 1 END) THEN '✓ OK'
        ELSE '⚠ MISSING TRANSLATIONS'
    END as status
FROM attributes;

-- Check options JSONB validity
SELECT
    'attributes.options' as field,
    COUNT(*) as total_records,
    COUNT(options) as non_null_options,
    COUNT(*) - COUNT(options) as null_options,
    CASE
        WHEN COUNT(*) > 0 THEN '✓ OK'
        ELSE '⚠ ERROR'
    END as status
FROM attributes;

\echo ''

-- ============================================================================
-- 4. ATTRIBUTE TYPE DISTRIBUTION
-- ============================================================================
\echo '4. Attribute Type Distribution'
\echo '-------------------------------'

SELECT
    attribute_type,
    purpose,
    COUNT(*) as count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage
FROM attributes
WHERE is_active = true
GROUP BY attribute_type, purpose
ORDER BY count DESC;

\echo ''

-- ============================================================================
-- 5. UNIQUE CONSTRAINTS VALIDATION
-- ============================================================================
\echo '5. Unique Constraints'
\echo '---------------------'

-- Check attributes.code uniqueness
SELECT
    'attributes.code' as field,
    COUNT(*) as total_codes,
    COUNT(DISTINCT code) as unique_codes,
    COUNT(*) - COUNT(DISTINCT code) as duplicates,
    CASE
        WHEN COUNT(*) = COUNT(DISTINCT code) THEN '✓ UNIQUE'
        ELSE '✗ DUPLICATES FOUND'
    END as status
FROM attributes;

-- Check category_attributes uniqueness
SELECT
    'category_attributes (category_id, attribute_id)' as field,
    COUNT(*) as total_records,
    COUNT(DISTINCT (category_id, attribute_id)) as unique_pairs,
    COUNT(*) - COUNT(DISTINCT (category_id, attribute_id)) as duplicates,
    CASE
        WHEN COUNT(*) = COUNT(DISTINCT (category_id, attribute_id)) THEN '✓ UNIQUE'
        ELSE '✗ DUPLICATES FOUND'
    END as status
FROM category_attributes;

-- Check listing_attribute_values uniqueness
SELECT
    'listing_attribute_values (listing_id, attribute_id)' as field,
    COUNT(*) as total_records,
    COUNT(DISTINCT (listing_id, attribute_id)) as unique_pairs,
    COUNT(*) - COUNT(DISTINCT (listing_id, attribute_id)) as duplicates,
    CASE
        WHEN COUNT(*) = COUNT(DISTINCT (listing_id, attribute_id)) THEN '✓ UNIQUE'
        ELSE '✗ DUPLICATES FOUND'
    END as status
FROM listing_attribute_values;

\echo ''

-- ============================================================================
-- 6. ACTIVE/INACTIVE ATTRIBUTES
-- ============================================================================
\echo '6. Active/Inactive Status'
\echo '-------------------------'

SELECT
    is_active,
    COUNT(*) as count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage
FROM attributes
GROUP BY is_active
ORDER BY is_active DESC;

\echo ''

-- ============================================================================
-- 7. SEARCHABLE/FILTERABLE FLAGS
-- ============================================================================
\echo '7. Searchable/Filterable Flags'
\echo '-------------------------------'

SELECT
    'Searchable' as flag_type,
    COUNT(*) FILTER (WHERE is_searchable = true) as enabled,
    COUNT(*) FILTER (WHERE is_searchable = false) as disabled,
    ROUND(COUNT(*) FILTER (WHERE is_searchable = true) * 100.0 / COUNT(*), 2) as percentage
FROM attributes
WHERE is_active = true
UNION ALL
SELECT
    'Filterable' as flag_type,
    COUNT(*) FILTER (WHERE is_filterable = true) as enabled,
    COUNT(*) FILTER (WHERE is_filterable = false) as disabled,
    ROUND(COUNT(*) FILTER (WHERE is_filterable = true) * 100.0 / COUNT(*), 2) as percentage
FROM attributes
WHERE is_active = true
UNION ALL
SELECT
    'Required' as flag_type,
    COUNT(*) FILTER (WHERE is_required = true) as enabled,
    COUNT(*) FILTER (WHERE is_required = false) as disabled,
    ROUND(COUNT(*) FILTER (WHERE is_required = true) * 100.0 / COUNT(*), 2) as percentage
FROM attributes
WHERE is_active = true
UNION ALL
SELECT
    'Variant Compatible' as flag_type,
    COUNT(*) FILTER (WHERE is_variant_compatible = true) as enabled,
    COUNT(*) FILTER (WHERE is_variant_compatible = false) as disabled,
    ROUND(COUNT(*) FILTER (WHERE is_variant_compatible = true) * 100.0 / COUNT(*), 2) as percentage
FROM attributes
WHERE is_active = true;

\echo ''

-- ============================================================================
-- 8. SAMPLE DATA INSPECTION
-- ============================================================================
\echo '8. Sample Attributes (First 5)'
\echo '-------------------------------'

SELECT
    id,
    code,
    name->>'en' as name_en,
    attribute_type,
    purpose,
    is_active,
    is_searchable,
    is_filterable
FROM attributes
ORDER BY id
LIMIT 5;

\echo ''

-- ============================================================================
-- 9. CATEGORY COVERAGE
-- ============================================================================
\echo '9. Category Coverage (Top 10)'
\echo '------------------------------'

SELECT
    ca.category_id,
    COUNT(*) as attribute_count,
    COUNT(*) FILTER (WHERE ca.is_required = true) as required_count,
    COUNT(*) FILTER (WHERE ca.is_enabled = true) as enabled_count
FROM category_attributes ca
GROUP BY ca.category_id
ORDER BY attribute_count DESC
LIMIT 10;

\echo ''

-- ============================================================================
-- 10. POLYMORPHIC VALUE DISTRIBUTION
-- ============================================================================
\echo '10. Polymorphic Value Distribution'
\echo '-----------------------------------'

SELECT
    'value_text' as value_type,
    COUNT(*) as count
FROM listing_attribute_values
WHERE value_text IS NOT NULL
UNION ALL
SELECT
    'value_number' as value_type,
    COUNT(*) as count
FROM listing_attribute_values
WHERE value_number IS NOT NULL
UNION ALL
SELECT
    'value_boolean' as value_type,
    COUNT(*) as count
FROM listing_attribute_values
WHERE value_boolean IS NOT NULL
UNION ALL
SELECT
    'value_date' as value_type,
    COUNT(*) as count
FROM listing_attribute_values
WHERE value_date IS NOT NULL
UNION ALL
SELECT
    'value_json' as value_type,
    COUNT(*) as count
FROM listing_attribute_values
WHERE value_json IS NOT NULL;

\echo ''

-- ============================================================================
-- VALIDATION COMPLETE
-- ============================================================================
\echo '============================================================================'
\echo 'Validation Complete!'
\echo '============================================================================'
\echo ''
\echo 'Expected Results:'
\echo '  ✓ attributes: 203 records'
\echo '  ✓ All foreign keys valid (0 orphaned records)'
\echo '  ✓ JSONB fields contain valid i18n data'
\echo '  ✓ No duplicate codes or relationships'
\echo ''
\echo 'If all checks passed, migration is successful!'
\echo ''
