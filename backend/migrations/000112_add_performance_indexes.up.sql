-- Performance indexes for categories and attributes system

-- Index for parent_id lookup in marketplace_categories (for hierarchical queries)
CREATE INDEX IF NOT EXISTS idx_marketplace_categories_parent_id 
ON marketplace_categories(parent_id) 
WHERE parent_id IS NOT NULL;

-- Composite index for category_attributes table (category_id, attribute_id)
-- This optimizes queries that filter by both category and attribute
CREATE INDEX IF NOT EXISTS idx_category_attributes_category_attribute 
ON category_attributes(category_id, attribute_id);

-- Index for is_active field in marketplace_categories (for filtering active categories)
CREATE INDEX IF NOT EXISTS idx_marketplace_categories_is_active 
ON marketplace_categories(is_active) 
WHERE is_active = true;

-- Composite index for translations table (entity_type, entity_id, field_name)
-- This optimizes translation lookups
CREATE INDEX IF NOT EXISTS idx_translations_entity_type_id_field 
ON translations(entity_type, entity_id, field_name);

-- Index for language_code in translations (for filtering by language)
CREATE INDEX IF NOT EXISTS idx_translations_language_code 
ON translations(language_code);

-- Index for sorting categories by sort_order
CREATE INDEX IF NOT EXISTS idx_marketplace_categories_sort_order 
ON marketplace_categories(sort_order) 
WHERE sort_order IS NOT NULL;

-- Index for category depth queries (parent_id + sort_order for efficient tree traversal)
CREATE INDEX IF NOT EXISTS idx_marketplace_categories_parent_sort 
ON marketplace_categories(parent_id, sort_order) 
WHERE parent_id IS NOT NULL;

-- Index for attribute ordering within categories
CREATE INDEX IF NOT EXISTS idx_category_attributes_category_sort 
ON category_attributes(category_id, sort_order);

-- Index for finding attributes by type
CREATE INDEX IF NOT EXISTS idx_category_attributes_attribute_type 
ON category_attributes(attribute_type);

-- Index for required attributes lookup
CREATE INDEX IF NOT EXISTS idx_category_attributes_is_required 
ON category_attributes(is_required) 
WHERE is_required = true;

-- Index for attribute groups by category
CREATE INDEX IF NOT EXISTS idx_category_attribute_groups_category_id 
ON category_attribute_groups(category_id);

-- Index for SEO fields queries
CREATE INDEX IF NOT EXISTS idx_marketplace_categories_seo_title 
ON marketplace_categories(seo_title) 
WHERE seo_title IS NOT NULL AND seo_title != '';

-- Gin index for text search on category names and descriptions
CREATE INDEX IF NOT EXISTS idx_marketplace_categories_text_search 
ON marketplace_categories 
USING gin(to_tsvector('english', coalesce(name, '') || ' ' || coalesce(description, '')));

-- Index for efficient lookups of category attributes with translations
CREATE INDEX IF NOT EXISTS idx_category_attributes_with_lang 
ON category_attributes(category_id, attribute_type, is_active) 
WHERE is_active = true;