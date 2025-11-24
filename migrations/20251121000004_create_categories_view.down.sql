-- Rollback: Drop categories view

BEGIN;

DROP VIEW IF EXISTS categories CASCADE;

COMMIT;
