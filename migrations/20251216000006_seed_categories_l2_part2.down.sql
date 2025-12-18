-- Rollback: Remove L2 categories (Part 2)
-- Date: 2025-12-16
-- Purpose: Remove L2 categories for Lepota, Bebe, Sport

DELETE FROM categories WHERE level = 2 AND parent_id IN (
    SELECT id FROM categories WHERE slug IN (
        'lepota-i-zdravlje',
        'za-bebe-i-decu',
        'sport-i-turizam'
    )
);
