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
    WHERE (
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

-- Create function to get category path with names (simplified version without path column)
CREATE OR REPLACE FUNCTION get_category_path_names(
    category_id INTEGER,
    lang_code VARCHAR(10) DEFAULT 'ru'
) RETURNS TEXT[] AS $$
DECLARE
    path_names TEXT[] := '{}';
    current_id INTEGER := category_id;
    cat_name TEXT;
BEGIN
    -- Build path by traversing parent_id
    WHILE current_id IS NOT NULL LOOP
        SELECT t.translated_text, c.parent_id INTO cat_name, current_id
        FROM marketplace_categories c
        LEFT JOIN translations t ON t.entity_id = c.id
            AND t.entity_type = 'category'
            AND t.field_name = 'name'
            AND t.language = lang_code
        WHERE c.id = current_id;
        
        IF cat_name IS NOT NULL THEN
            path_names := array_prepend(cat_name, path_names);
        END IF;
    END LOOP;
    
    RETURN path_names;
END;
$$ LANGUAGE plpgsql;

-- Add indexes for attribute values in listing_attribute_values (corrected table name)
CREATE INDEX IF NOT EXISTS idx_listing_attribute_values_text_value_trgm 
    ON listing_attribute_values USING gin (text_value gin_trgm_ops)
    WHERE text_value IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_listing_attribute_values_text_value_unaccent_trgm 
    ON listing_attribute_values USING gin (f_unaccent(text_value) gin_trgm_ops)
    WHERE text_value IS NOT NULL;

-- Create materialized view for category statistics (for search relevance)
CREATE MATERIALIZED VIEW IF NOT EXISTS category_listing_stats AS
SELECT 
    c.id as category_id,
    c.slug as category_slug,
    c.parent_id as parent_id,
    COUNT(DISTINCT l.id) as listing_count,
    COUNT(DISTINCT l.user_id) as seller_count,
    AVG(l.price) as avg_price,
    MIN(l.price) as min_price,
    MAX(l.price) as max_price
FROM marketplace_categories c
LEFT JOIN marketplace_listings l ON l.category_id = c.id AND l.status = 'active'
GROUP BY c.id, c.slug, c.parent_id;

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