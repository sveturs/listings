-- Rollback telegram translations updates
-- Restore language_changed back to language_saved
UPDATE translations 
SET field_name = 'language_saved'
WHERE entity_type = 'telegram_bot' 
AND entity_id = 3 
AND field_name = 'language_changed';

-- Remove added translations
DELETE FROM translations 
WHERE entity_type = 'telegram_bot' 
AND entity_id BETWEEN 50 AND 64;