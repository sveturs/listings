-- Create table for search synonyms
-- This table will store word synonyms to improve search results
CREATE TABLE IF NOT EXISTS search_synonyms (
    id SERIAL PRIMARY KEY,
    term VARCHAR(255) NOT NULL,
    synonym VARCHAR(255) NOT NULL,
    language VARCHAR(10) NOT NULL DEFAULT 'ru',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_search_synonyms_term ON search_synonyms(term);
CREATE INDEX IF NOT EXISTS idx_search_synonyms_synonym ON search_synonyms(synonym);
CREATE INDEX IF NOT EXISTS idx_search_synonyms_language ON search_synonyms(language);
CREATE INDEX IF NOT EXISTS idx_search_synonyms_active ON search_synonyms(is_active) WHERE is_active = true;

-- Create unique constraint to prevent duplicate synonym mappings
CREATE UNIQUE INDEX IF NOT EXISTS idx_search_synonyms_unique 
    ON search_synonyms(term, synonym, language) WHERE is_active = true;

-- DROP TRIGGER IF EXISTS to ON update_search_synonyms_updated_at;
CREATE TRIGGER to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_search_synonyms_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_search_synonyms_updated_at ON search_synonyms;
CREATE TRIGGER trigger_update_search_synonyms_updated_at
    BEFORE UPDATE ON search_synonyms
    FOR EACH ROW
    EXECUTE FUNCTION update_search_synonyms_updated_at();

-- Insert some common synonyms for Serbian/Russian marketplace
INSERT INTO search_synonyms (term, synonym, language) VALUES
-- Russian synonyms
('квартира', 'апартаменты', 'ru'),
('квартира', 'жилье', 'ru'),
('квартира', 'недвижимость', 'ru'),
('апартаменты', 'квартира', 'ru'),
('апартаменты', 'жилье', 'ru'),
('комната', 'спальня', 'ru'),
('дом', 'коттедж', 'ru'),
('дом', 'вилла', 'ru'),
('машина', 'автомобиль', 'ru'),
('машина', 'авто', 'ru'),
('автомобиль', 'машина', 'ru'),
('автомобиль', 'авто', 'ru'),
('телефон', 'смартфон', 'ru'),
('телефон', 'мобильный', 'ru'),
('ноутбук', 'лэптоп', 'ru'),
('ноутбук', 'компьютер', 'ru'),
('работа', 'вакансия', 'ru'),
('работа', 'труд', 'ru'),
('услуги', 'сервис', 'ru'),

-- Serbian synonyms (Latin)
('stan', 'apartman', 'sr'),
('stan', 'smeštaj', 'sr'),
('kuća', 'vila', 'sr'),
('automobil', 'auto', 'sr'),
('automobil', 'kola', 'sr'),
('telefon', 'mobilni', 'sr'),
('posao', 'rad', 'sr'),
('posao', 'zaposlenje', 'sr'),

-- English synonyms
('apartment', 'flat', 'en'),
('apartment', 'condo', 'en'),
('house', 'home', 'en'),
('house', 'residence', 'en'),
('car', 'vehicle', 'en'),
('car', 'automobile', 'en'),
('phone', 'mobile', 'en'),
('phone', 'smartphone', 'en'),
('laptop', 'notebook', 'en'),
('laptop', 'computer', 'en'),
('job', 'work', 'en'),
('job', 'employment', 'en'),
('service', 'services', 'en')
ON CONFLICT DO NOTHING;

-- Create function to expand search query with synonyms
CREATE OR REPLACE FUNCTION expand_search_query(
    query_text TEXT,
    query_language VARCHAR(10) DEFAULT 'ru'
) RETURNS TEXT AS $$
DECLARE
    word TEXT;
    synonym_text TEXT;
    expanded_query TEXT := query_text;
    synonyms_array TEXT[];
BEGIN
    -- Split query into words
    FOREACH word IN ARRAY string_to_array(lower(query_text), ' ')
    LOOP
        -- Find all active synonyms for this word
        SELECT array_agg(DISTINCT synonym) INTO synonyms_array
        FROM search_synonyms
        WHERE is_active = true 
          AND language = query_language
          AND (term = word OR synonym = word);
        
        -- If synonyms found, add them to the query
        IF synonyms_array IS NOT NULL THEN
            synonym_text := array_to_string(synonyms_array, ' | ');
            expanded_query := expanded_query || ' | ' || synonym_text;
        END IF;
    END LOOP;
    
    RETURN expanded_query;
END;
$$ LANGUAGE plpgsql;

-- Create table for search query analytics (optional, for future improvements)
CREATE TABLE IF NOT EXISTS search_analytics (
    id SERIAL PRIMARY KEY,
    original_query TEXT NOT NULL,
    expanded_query TEXT,
    results_count INTEGER,
    user_id INTEGER REFERENCES users(id),
    session_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_search_analytics_created_at ON search_analytics(created_at);
CREATE INDEX IF NOT EXISTS idx_search_analytics_user_id ON search_analytics(user_id);