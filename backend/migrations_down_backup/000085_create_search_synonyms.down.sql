-- Drop analytics table
DROP TABLE IF EXISTS search_analytics;

-- Drop synonym expansion function
DROP FUNCTION IF EXISTS expand_search_query(TEXT, VARCHAR);

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_search_synonyms_updated_at ON search_synonyms;
DROP FUNCTION IF EXISTS update_search_synonyms_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_search_synonyms_unique;
DROP INDEX IF EXISTS idx_search_synonyms_active;
DROP INDEX IF EXISTS idx_search_synonyms_language;
DROP INDEX IF EXISTS idx_search_synonyms_synonym;
DROP INDEX IF EXISTS idx_search_synonyms_term;

-- Drop synonyms table
DROP TABLE IF EXISTS search_synonyms;