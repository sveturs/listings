-- Migration Rollback: 20251217030016_seed_category_attributes
-- Description: Remove all category-attribute relationships
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17

TRUNCATE TABLE category_attributes;
