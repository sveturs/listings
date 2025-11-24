-- ============================================================================
-- Rollback Script for Attributes Migration
-- ============================================================================
-- Run on target database: listings_dev_db (localhost:35434)
-- DANGER: This will DELETE all migrated attributes!
-- ============================================================================

-- Safety check: Show what will be deleted
\echo '=== ROLLBACK PREVIEW ==='
\echo ''
\echo 'The following attributes will be DELETED:'
\echo ''

SELECT
    COUNT(*) as total_to_delete,
    MIN(id) as min_id,
    MAX(id) as max_id
FROM attributes;

\echo ''
\echo 'Sample records that will be deleted:'
SELECT id, code, name->>'en' as name, display_name->>'en' as display_name
FROM attributes
ORDER BY id
LIMIT 10;

\echo ''
\echo '=== WARNING ==='
\echo 'This will DELETE all attributes from the table!'
\echo 'Press Ctrl+C to cancel, or press Enter to continue...'
\prompt 'Type YES to confirm deletion: ' confirm

-- Conditional deletion (requires \if support in psql 12+)
-- For older psql, comment out \if/\endif and run manually

\if :{?confirm}
    \if :confirm = 'YES'
        BEGIN;

        \echo 'Deleting attributes...'
        DELETE FROM attributes;

        \echo 'Resetting sequence...'
        SELECT setval('attributes_id_seq', 1, false);

        \echo 'Verifying deletion...'
        SELECT COUNT(*) as remaining_records FROM attributes;

        \echo ''
        \echo 'Rollback transaction? (will UNDO the deletion)'
        \prompt 'Type COMMIT to confirm, or ROLLBACK to undo: ' action

        \if :action = 'COMMIT'
            COMMIT;
            \echo '✓ Rollback committed. All attributes deleted.'
        \else
            ROLLBACK;
            \echo '✓ Rollback cancelled. No changes made.'
        \endif
    \else
        \echo 'Confirmation failed. No changes made.'
    \endif
\else
    \echo 'No confirmation provided. No changes made.'
\endif

\echo ''
\echo '=== Rollback Script Complete ==='
