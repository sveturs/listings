-- Migration rollback: Clear all meta_keywords
-- Description: Remove all SEO meta keywords from categories table

UPDATE categories SET meta_keywords = '{}'::jsonb;
