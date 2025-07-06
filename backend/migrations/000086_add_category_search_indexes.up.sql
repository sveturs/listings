-- Add text search capabilities to categories and attributes

-- Add trigram indexes for category names in translations table
CREATE INDEX IF NOT EXISTS idx_translations_translated_text_trgm 
    ON translations USING gin (translated_text gin_trgm_ops)
    WHERE entity_type IN ('category', 'attribute', 'attribute_option');

-- Add unaccent trigram indexes for translations
CREATE INDEX IF NOT EXISTS idx_translations_translated_text_unaccent_trgm 
    ON translations USING gin (f_unaccent(translated_text) gin_trgm_ops)
    WHERE entity_type IN ('category', 'attribute', 'attribute_option');

-- Add full-text search indexes for translations
CREATE INDEX IF NOT EXISTS idx_translations_translated_text_fts_ru 
    ON translations USING gin (to_tsvector('russian_unaccent', translated_text))
    WHERE entity_type IN ('category', 'attribute', 'attribute_option') AND language = 'ru';

CREATE INDEX IF NOT EXISTS idx_translations_translated_text_fts_en 
    ON translations USING gin (to_tsvector('english_unaccent', translated_text))
    WHERE entity_type IN ('category', 'attribute', 'attribute_option') AND language = 'en';

-- Create function to search categories by name with fuzzy matching
CREATE OR REPLACE FUNCTION search_categories_fuzzy(
    search_term TEXT,
    lang_code VARCHAR(10) DEFAULT 'ru',
    similarity_threshold FLOAT DEFAULT 0.3
) RETURNS TABLE (
    category_id INTEGER,
    category_slug VARCHAR(255),
    category_name TEXT,
    similarity_score FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT
        c.id as category_id,
        c.slug as category_slug,
        t.translated_text as category_name,
        similarity(f_unaccent(t.translated_text), f_unaccent(search_term)) as similarity_score
    FROM marketplace_categories c
    JOIN translations t ON t.entity_id = c.id 
        AND t.entity_type = 'category' 
        AND t.field_name = 'name'
        AND t.language = lang_code
    WHERE c.is_active = true
        AND (
            -- Exact match (case-insensitive)
            lower(f_unaccent(t.translated_text)) = lower(f_unaccent(search_term))
            -- Fuzzy match using trigrams
            OR similarity(f_unaccent(t.translated_text), f_unaccent(search_term)) >= similarity_threshold
            -- Partial match
            OR f_unaccent(t.translated_text) ILIKE '%' || f_unaccent(search_term) || '%'
        )
    ORDER BY 
        -- Exact matches first
        CASE WHEN lower(f_unaccent(t.translated_text)) = lower(f_unaccent(search_term)) THEN 0 ELSE 1 END,
        -- Then by similarity score
        similarity_score DESC,
        -- Then alphabetically
        category_name;
END;
$$ LANGUAGE plpgsql;

-- Create function to get category path with names
CREATE OR REPLACE FUNCTION get_category_path_names(
    category_id INTEGER,
    lang_code VARCHAR(10) DEFAULT 'ru'
) RETURNS TEXT[] AS $$
DECLARE
    path_ids INTEGER[];
    path_names TEXT[] := '{}';
    cat_id INTEGER;
    cat_name TEXT;
BEGIN
    -- Get category path
    SELECT string_to_array(path, '.')::INTEGER[] INTO path_ids
    FROM marketplace_categories
    WHERE id = category_id;
    
    -- Get names for each category in path
    FOREACH cat_id IN ARRAY path_ids
    LOOP
        SELECT t.translated_text INTO cat_name
        FROM translations t
        WHERE t.entity_id = cat_id
            AND t.entity_type = 'category'
            AND t.field_name = 'name'
            AND t.language = lang_code;
        
        path_names := array_append(path_names, COALESCE(cat_name, ''));
    END LOOP;
    
    RETURN path_names;
END;
$$ LANGUAGE plpgsql;

-- Add indexes for attribute values in marketplace_listing_attributes
CREATE INDEX IF NOT EXISTS idx_listing_attributes_value_text_trgm 
    ON marketplace_listing_attributes USING gin (value_text gin_trgm_ops)
    WHERE value_text IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_listing_attributes_value_text_unaccent_trgm 
    ON marketplace_listing_attributes USING gin (f_unaccent(value_text) gin_trgm_ops)
    WHERE value_text IS NOT NULL;

-- Create materialized view for category statistics (for search relevance)
CREATE MATERIALIZED VIEW IF NOT EXISTS category_listing_stats AS
SELECT 
    c.id as category_id,
    c.slug as category_slug,
    c.path as category_path,
    COUNT(DISTINCT l.id) as listing_count,
    COUNT(DISTINCT l.user_id) as seller_count,
    AVG(l.price) as avg_price,
    MIN(l.price) as min_price,
    MAX(l.price) as max_price
FROM marketplace_categories c
LEFT JOIN marketplace_listings l ON l.category_id = c.id AND l.status = 'active'
WHERE c.is_active = true
GROUP BY c.id, c.slug, c.path;

CREATE UNIQUE INDEX IF NOT EXISTS idx_category_listing_stats_id 
    ON category_listing_stats(category_id);

CREATE INDEX IF NOT EXISTS idx_category_listing_stats_listing_count 
    ON category_listing_stats(listing_count DESC);

-- Function to refresh category statistics
CREATE OR REPLACE FUNCTION refresh_category_stats() 
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY category_listing_stats;
END;
$$ LANGUAGE plpgsql;