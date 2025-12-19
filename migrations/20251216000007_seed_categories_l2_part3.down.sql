-- Rollback: Remove L2 categories (Part 3 - Final)
-- Date: 2025-12-16
-- Purpose: Remove L2 categories for remaining 12 L1 categories

DELETE FROM categories WHERE level = 2 AND parent_id IN (
    SELECT id FROM categories WHERE slug IN (
        'automobilizam',
        'kucni-aparati',
        'nakit-i-satovi',
        'knjige-i-mediji',
        'kucni-ljubimci',
        'kancelarijski-materijal',
        'muzicki-instrumenti',
        'hrana-i-pice',
        'umetnost-i-rukotvorine',
        'industrija-i-alati',
        'usluge',
        'ostalo'
    )
);
