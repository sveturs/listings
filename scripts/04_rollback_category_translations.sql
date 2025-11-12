-- ============================================================================
-- Rollback: Restore categories from backup or clear translations
-- Description: Provides options to rollback translation changes
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

-- Option 1: Restore from backup table (if backup exists)
-- ============================================================================
-- Uncomment to restore from backup:
/*
BEGIN;

-- Restore translation columns from backup
UPDATE c2c_categories c
SET
  title_en = b.title_en,
  title_ru = b.title_ru,
  title_sr = b.title_sr
FROM c2c_categories_backup_20251110 b
WHERE c.id = b.id;

-- Verify restoration
SELECT
  COUNT(*) as total_categories,
  COUNT(title_en) as with_english,
  COUNT(title_ru) as with_russian,
  COUNT(title_sr) as with_serbian
FROM c2c_categories;

COMMIT;
*/

-- Option 2: Clear all translations (set to NULL)
-- ============================================================================
-- Uncomment to clear all translations:
/*
BEGIN;

UPDATE c2c_categories
SET
  title_en = NULL,
  title_ru = NULL,
  title_sr = NULL;

-- Verify clearing
SELECT
  COUNT(*) as total_categories,
  COUNT(title_en) as with_english,
  COUNT(title_ru) as with_russian,
  COUNT(title_sr) as with_serbian
FROM c2c_categories;

COMMIT;
*/

-- Option 3: Drop translation columns completely
-- ============================================================================
-- WARNING: This will permanently remove translation columns!
-- Uncomment to drop columns:
/*
BEGIN;

ALTER TABLE c2c_categories
  DROP COLUMN IF EXISTS title_en,
  DROP COLUMN IF EXISTS title_ru,
  DROP COLUMN IF EXISTS title_sr;

-- Drop indexes
DROP INDEX IF EXISTS idx_c2c_categories_title_en;
DROP INDEX IF EXISTS idx_c2c_categories_title_ru;
DROP INDEX IF EXISTS idx_c2c_categories_title_sr;

COMMIT;
*/

-- ============================================================================
-- INSTRUCTIONS
-- ============================================================================
\echo '⚠️  Rollback script loaded but not executed.'
\echo ''
\echo 'To rollback translations, choose ONE option:'
\echo ''
\echo 'Option 1: Restore from backup (requires backup table)'
\echo '  - Uncomment the "Option 1" section'
\echo '  - Run this script again'
\echo ''
\echo 'Option 2: Clear all translations'
\echo '  - Uncomment the "Option 2" section'
\echo '  - Run this script again'
\echo ''
\echo 'Option 3: Drop translation columns'
\echo '  - Uncomment the "Option 3" section'
\echo '  - Run this script again'
\echo '  - WARNING: This is irreversible!'
