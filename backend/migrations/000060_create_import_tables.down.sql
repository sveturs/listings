-- Drop triggers
DROP TRIGGER IF EXISTS trigger_category_mappings_updated_at ON category_mappings;
DROP TRIGGER IF EXISTS trigger_import_jobs_updated_at ON import_jobs;

-- Drop function
DROP FUNCTION IF EXISTS update_import_tables_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_storefront_products_sku;
DROP INDEX IF EXISTS idx_storefront_products_external_id;
DROP INDEX IF EXISTS idx_category_mappings_import_categories;
DROP INDEX IF EXISTS idx_category_mappings_storefront_id;
DROP INDEX IF EXISTS idx_import_errors_line_number;
DROP INDEX IF EXISTS idx_import_errors_job_id;
DROP INDEX IF EXISTS idx_import_jobs_created_at;
DROP INDEX IF EXISTS idx_import_jobs_status;
DROP INDEX IF EXISTS idx_import_jobs_storefront_id;

-- Remove external_id column if it was added
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'storefront_products' 
        AND column_name = 'external_id'
    ) THEN
        ALTER TABLE storefront_products DROP COLUMN external_id;
    END IF;
END $$;

-- Drop tables
DROP TABLE IF EXISTS category_mappings;
DROP TABLE IF EXISTS import_errors;
DROP TABLE IF EXISTS import_jobs;