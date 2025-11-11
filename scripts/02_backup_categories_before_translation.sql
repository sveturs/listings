-- ============================================================================
-- Backup: Create backup of c2c_categories before adding translations
-- Description: Creates a backup table with current state
-- Author: System
-- Date: 2025-11-10
-- Usage: Run before applying translations
-- ============================================================================

-- Create backup table
DROP TABLE IF EXISTS c2c_categories_backup_20251110;

CREATE TABLE c2c_categories_backup_20251110 AS
SELECT * FROM c2c_categories;

-- Verify backup
SELECT
  'Original table count' as info,
  COUNT(*) as count
FROM c2c_categories
UNION ALL
SELECT
  'Backup table count' as info,
  COUNT(*) as count
FROM c2c_categories_backup_20251110;

-- Show categories without translations
SELECT
  COUNT(*) as categories_without_translations
FROM c2c_categories
WHERE title_en IS NULL OR title_ru IS NULL OR title_sr IS NULL;
