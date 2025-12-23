-- Migration: 20251222120001_move_smartwatches
-- Purpose: Move smartwatch categories from Elektronika to Nakit i satovi
-- Fixes Bug #3: Smartwatches should be under Jewelry & Watches, not Electronics

-- Get the target parent category ID
DO $$
DECLARE
    nakit_category_id UUID;
    elektronika_category_id UUID;
    pametni_satovi_id UUID;
BEGIN
    -- Get category IDs
    SELECT id INTO nakit_category_id FROM categories WHERE slug = 'nakit-i-satovi' LIMIT 1;
    SELECT id INTO elektronika_category_id FROM categories WHERE slug = 'elektronika' LIMIT 1;
    SELECT id INTO pametni_satovi_id FROM categories WHERE slug = 'pametni-satovi' AND parent_id = elektronika_category_id LIMIT 1;

    IF nakit_category_id IS NULL THEN
        RAISE EXCEPTION 'Category nakit-i-satovi not found';
    END IF;

    IF pametni_satovi_id IS NULL THEN
        RAISE NOTICE 'Smartwatch category not found under Elektronika, skipping migration';
        RETURN;
    END IF;

    -- Move pametni-satovi (Level 2) from elektronika to nakit-i-satovi
    UPDATE categories
    SET
        parent_id = nakit_category_id,
        path = 'nakit-i-satovi/pametni-satovi',
        updated_at = now()
    WHERE id = pametni_satovi_id;

    RAISE NOTICE 'Moved pametni-satovi to nakit-i-satovi';

    -- Update all Level 3 smartwatch subcategories paths
    UPDATE categories
    SET
        path = REPLACE(path, 'elektronika/pametni-satovi', 'nakit-i-satovi/pametni-satovi'),
        updated_at = now()
    WHERE parent_id = pametni_satovi_id;

    RAISE NOTICE 'Updated paths for smartwatch subcategories';

END $$;

-- Verify the migration
DO $$
DECLARE
    moved_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO moved_count
    FROM categories
    WHERE slug = 'pametni-satovi'
    AND path = 'nakit-i-satovi/pametni-satovi';

    IF moved_count > 0 THEN
        RAISE NOTICE 'Migration successful: pametni-satovi moved to nakit-i-satovi';
    ELSE
        RAISE NOTICE 'Migration skipped or category structure different';
    END IF;
END $$;
