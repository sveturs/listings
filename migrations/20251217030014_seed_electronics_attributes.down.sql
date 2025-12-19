-- Migration Rollback: 20251217030014_seed_electronics_attributes
-- Description: Remove electronics attributes seed data
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17

DELETE FROM attributes WHERE code IN (
    'screen_size',
    'processor',
    'ram',
    'storage_capacity',
    'operating_system',
    'connectivity',
    'battery_capacity',
    'camera_resolution',
    'refresh_rate',
    'resolution'
);
