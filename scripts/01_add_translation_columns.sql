-- ============================================================================
-- Migration: Add translation columns to c2c_categories
-- Description: Adds title_en, title_ru, title_sr columns for multi-language support
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

-- Add translation columns
ALTER TABLE c2c_categories
  ADD COLUMN IF NOT EXISTS title_en VARCHAR(255),
  ADD COLUMN IF NOT EXISTS title_ru VARCHAR(255),
  ADD COLUMN IF NOT EXISTS title_sr VARCHAR(255);

-- Add indexes for translation columns (for faster searches)
CREATE INDEX IF NOT EXISTS idx_c2c_categories_title_en ON c2c_categories(title_en);
CREATE INDEX IF NOT EXISTS idx_c2c_categories_title_ru ON c2c_categories(title_ru);
CREATE INDEX IF NOT EXISTS idx_c2c_categories_title_sr ON c2c_categories(title_sr);

-- Add comments
COMMENT ON COLUMN c2c_categories.title_en IS 'Category title in English';
COMMENT ON COLUMN c2c_categories.title_ru IS 'Category title in Russian';
COMMENT ON COLUMN c2c_categories.title_sr IS 'Category title in Serbian (Latin)';
