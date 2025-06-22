-- Create import_jobs table
CREATE TABLE import_jobs (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(10) NOT NULL CHECK (file_type IN ('xml', 'csv', 'zip')),
    file_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    total_records INTEGER DEFAULT 0,
    processed_records INTEGER DEFAULT 0,
    successful_records INTEGER DEFAULT 0,
    failed_records INTEGER DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create import_errors table
CREATE TABLE import_errors (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES import_jobs(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,
    field_name VARCHAR(100) NOT NULL,
    error_message TEXT NOT NULL,
    raw_data TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create category_mappings table
CREATE TABLE category_mappings (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    import_category1 VARCHAR(255) NOT NULL,
    import_category2 VARCHAR(255),
    import_category3 VARCHAR(255),
    local_category_id INTEGER NOT NULL REFERENCES marketplace_categories(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure unique mapping per storefront
    UNIQUE(storefront_id, import_category1, import_category2, import_category3)
);

-- Add external_id to storefront_products if not exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'storefront_products' 
        AND column_name = 'external_id'
    ) THEN
        ALTER TABLE storefront_products ADD COLUMN external_id VARCHAR(100);
    END IF;
END $$;

-- Create indexes for better performance
CREATE INDEX idx_import_jobs_storefront_id ON import_jobs(storefront_id);
CREATE INDEX idx_import_jobs_status ON import_jobs(status);
CREATE INDEX idx_import_jobs_created_at ON import_jobs(created_at);

CREATE INDEX idx_import_errors_job_id ON import_errors(job_id);
CREATE INDEX idx_import_errors_line_number ON import_errors(line_number);

CREATE INDEX idx_category_mappings_storefront_id ON category_mappings(storefront_id);
CREATE INDEX idx_category_mappings_import_categories ON category_mappings(import_category1, import_category2, import_category3);

CREATE INDEX IF NOT EXISTS idx_storefront_products_external_id ON storefront_products(external_id) WHERE external_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_storefront_products_sku ON storefront_products(sku) WHERE sku IS NOT NULL;

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_import_tables_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_import_jobs_updated_at
    BEFORE UPDATE ON import_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_import_tables_updated_at();

CREATE TRIGGER trigger_category_mappings_updated_at
    BEFORE UPDATE ON category_mappings
    FOR EACH ROW
    EXECUTE FUNCTION update_import_tables_updated_at();