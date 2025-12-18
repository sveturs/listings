-- Migration Rollback: 20251217030013_seed_clothing_attributes
-- Description: Remove clothing attributes seed data
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17

-- Delete attribute_values first (FK constraint)
DELETE FROM attribute_values
WHERE attribute_id IN (
    SELECT id FROM attributes WHERE code IN (
        'clothing_size', 'color'
    )
);

-- Delete attributes
DELETE FROM attributes WHERE code IN (
    'clothing_size',
    'color',
    'gender',
    'fit',
    'style',
    'season',
    'neckline',
    'sleeve_length',
    'pattern',
    'closure_type'
);
