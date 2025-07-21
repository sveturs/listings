-- Script to fix autoincrement for all tables with id columns
DO $$
DECLARE
    tbl RECORD;
    seq_name TEXT;
    max_id INTEGER;
BEGIN
    -- Get all tables with id column without default
    FOR tbl IN 
        SELECT 
            t.table_name,
            t.column_name
        FROM information_schema.columns t
        WHERE t.table_schema = 'public' 
        AND t.column_name = 'id' 
        AND t.data_type IN ('integer', 'bigint')
        AND t.column_default IS NULL
        AND EXISTS (
            -- Check if sequence exists
            SELECT 1 
            FROM pg_sequences s 
            WHERE s.schemaname = 'public' 
            AND s.sequencename = t.table_name || '_id_seq'
        )
        ORDER BY t.table_name
    LOOP
        seq_name := tbl.table_name || '_id_seq';
        
        -- Set default for id column
        EXECUTE format('ALTER TABLE %I ALTER COLUMN id SET DEFAULT nextval(%L::regclass)', 
            tbl.table_name, seq_name);
        
        -- Set sequence owner
        EXECUTE format('ALTER SEQUENCE %I OWNED BY %I.id', 
            seq_name, tbl.table_name);
        
        -- Get max id from table
        EXECUTE format('SELECT COALESCE(MAX(id), 0) FROM %I', tbl.table_name) INTO max_id;
        
        -- Reset sequence value
        EXECUTE format('SELECT setval(%L, %s + 1, false)', seq_name, max_id);
        
        RAISE NOTICE 'Fixed autoincrement for table: %', tbl.table_name;
    END LOOP;
END $$;