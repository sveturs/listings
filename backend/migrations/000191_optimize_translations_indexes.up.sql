-- Migration: Optimize translations table indexes
-- Date: 2025-10-13
-- Phase 2, Task 2.13
-- Goal: Remove unused indexes, keep only frequently used ones

-- Drop metadata GIN index (never used, 104 KB)
DROP INDEX IF EXISTS idx_translations_metadata;

-- Drop entity_field partial index (never used, 48 KB, only for listings)
DROP INDEX IF EXISTS idx_translations_entity_field;

-- Keep frequently used indexes:
-- - translations_entity_type_entity_id_language_field_name_key (UNIQUE) - 40,139 scans
-- - idx_translations_lookup - 4,196 scans
-- - idx_translations_listing_all - 3,126 scans
-- - idx_translations_entity_lang_field - 2,130 scans
-- - idx_translations_type_lang - 1,292 scans
-- - idx_translations_composite - 333 scans
