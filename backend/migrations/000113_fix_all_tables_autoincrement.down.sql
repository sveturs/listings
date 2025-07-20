-- Remove autoincrement from all tables (revert to previous state)
-- Note: This is not recommended as it will break insert operations

DO $$
DECLARE
    tbl RECORD;
BEGIN
    -- Get all tables with id column that have default nextval
    FOR tbl IN 
        SELECT 
            t.table_name
        FROM information_schema.columns t
        WHERE t.table_schema = 'public' 
        AND t.column_name = 'id' 
        AND t.data_type IN ('integer', 'bigint')
        AND t.column_default LIKE 'nextval%'
        ORDER BY t.table_name
    LOOP
        -- Remove default
        EXECUTE format('ALTER TABLE %I ALTER COLUMN id DROP DEFAULT', tbl.table_name);
        
        RAISE NOTICE 'Removed autoincrement from table: %', tbl.table_name;
    END LOOP;
END $$;