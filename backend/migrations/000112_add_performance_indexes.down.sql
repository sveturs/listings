-- Drop performance indexes for rollback

DROP INDEX IF EXISTS idx_marketplace_categories_parent_id;
DROP INDEX IF EXISTS idx_category_attributes_category_attribute;
DROP INDEX IF EXISTS idx_marketplace_categories_is_active;
DROP INDEX IF EXISTS idx_translations_entity_type_id_field;
DROP INDEX IF EXISTS idx_translations_language_code;
DROP INDEX IF EXISTS idx_marketplace_categories_sort_order;
DROP INDEX IF EXISTS idx_marketplace_categories_parent_sort;
DROP INDEX IF EXISTS idx_category_attributes_category_sort;
DROP INDEX IF EXISTS idx_category_attributes_attribute_type;
DROP INDEX IF EXISTS idx_category_attributes_is_required;
DROP INDEX IF EXISTS idx_category_attribute_groups_category_id;
DROP INDEX IF EXISTS idx_marketplace_categories_seo_title;
DROP INDEX IF EXISTS idx_marketplace_categories_text_search;
DROP INDEX IF EXISTS idx_category_attributes_with_lang;