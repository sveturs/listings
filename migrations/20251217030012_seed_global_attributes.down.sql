-- Migration Rollback: 20251217030012_seed_global_attributes
-- Description: Remove global attributes seed data
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17

DELETE FROM attributes WHERE code IN (
    'brand',
    'condition',
    'country_of_origin',
    'material',
    'weight',
    'dimensions',
    'warranty_months',
    'model_number',
    'year_of_manufacture',
    'energy_class'
);
