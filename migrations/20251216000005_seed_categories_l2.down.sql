-- Rollback: Remove L2 categories (Part 1)
-- Date: 2025-12-16
-- Purpose: Remove L2 categories for Odeća, Elektronika, Dom

DELETE FROM categories WHERE level = 2 AND parent_id IN (
    SELECT id FROM categories WHERE slug IN (
        'odeca-i-obuca',
        'elektronika',
        'dom-i-basta'
    )
);
