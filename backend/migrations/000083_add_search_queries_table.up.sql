-- Create search_queries table
CREATE TABLE IF NOT EXISTS search_queries (
    id SERIAL PRIMARY KEY,
    query TEXT NOT NULL,
    normalized_query TEXT NOT NULL,
    search_count INTEGER NOT NULL DEFAULT 1,
    last_searched TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    language VARCHAR(10) NOT NULL DEFAULT 'ru',
    results_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_search_queries_normalized_query ON search_queries(normalized_query);
CREATE INDEX IF NOT EXISTS idx_search_queries_search_count ON search_queries(search_count DESC);
CREATE INDEX IF NOT EXISTS idx_search_queries_language ON search_queries(language);

-- DROP TRIGGER IF EXISTS for ON update_search_queries_updated_at;
CREATE TRIGGER for automatic updated_at timestamp
CREATE OR REPLACE FUNCTION update_search_queries_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_search_queries_updated_at_trigger ON search_queries;
CREATE TRIGGER update_search_queries_updated_at_trigger
    BEFORE UPDATE ON search_queries
    FOR EACH ROW
    EXECUTE FUNCTION update_search_queries_updated_at();