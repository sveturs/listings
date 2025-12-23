-- Migration: 20251222120001_move_smartwatches (rollback)
-- Purpose: Move smartwatch categories back from Nakit i satovi to Elektronika

-- Rollback: Move pametni-satovi back to elektronika
DO $$
DECLARE
    nakit_category_id UUID;
    elektronika_category_id UUID;
    pametni_satovi_id UUID;
BEGIN
    -- Get category IDs
    SELECT id INTO nakit_category_id FROM categories WHERE slug = 'nakit-i-satovi' LIMIT 1;
    SELECT id INTO elektronika_category_id FROM categories WHERE slug = 'elektronika' LIMIT 1;
    SELECT id INTO pametni_satovi_id FROM categories WHERE slug = 'pametni-satovi' AND parent_id = nakit_category_id LIMIT 1;

    IF elektronika_category_id IS NULL THEN
        RAISE EXCEPTION 'Category elektronika not found';
    END IF;

    IF pametni_satovi_id IS NULL THEN
        RAISE NOTICE 'Smartwatch category not found under Nakit i satovi, skipping rollback';
        RETURN;
    END IF;

    -- Move pametni-satovi back to elektronika
    UPDATE categories
    SET
        parent_id = elektronika_category_id,
        path = 'elektronika/pametni-satovi',
        updated_at = now()
    WHERE id = pametni_satovi_id;

    RAISE NOTICE 'Moved pametni-satovi back to elektronika';

    -- Update all Level 3 smartwatch subcategories paths
    UPDATE categories
    SET
        path = REPLACE(path, 'nakit-i-satovi/pametni-satovi', 'elektronika/pametni-satovi'),
        updated_at = now()
    WHERE parent_id = pametni_satovi_id;

    RAISE NOTICE 'Updated paths for smartwatch subcategories';

END $$;
