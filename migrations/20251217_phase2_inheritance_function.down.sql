-- Migration rollback: 20251217_phase2_inheritance_function
-- Description: Drop inheritance function
-- Date: 2025-12-17

DROP FUNCTION IF EXISTS get_category_attributes_with_inheritance(UUID);
