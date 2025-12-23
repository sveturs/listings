-- Rollback final 3 L3 keywords
UPDATE categories SET meta_keywords = NULL WHERE slug IN ('lomači-dokumenta', 'masine-sudjе', 'tesle');
