-- Rollback: Remove L1 categories
-- Date: 2025-12-16

-- Delete all L1 categories (CASCADE will handle children if any exist)
DELETE FROM categories WHERE level = 1;

-- Verify deletion
DO $$
DECLARE
    category_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO category_count FROM categories WHERE level = 1;

    IF category_count != 0 THEN
        RAISE EXCEPTION 'Failed to delete L1 categories, found % remaining', category_count;
    END IF;

    RAISE NOTICE 'Successfully deleted all L1 categories';
END $$;
