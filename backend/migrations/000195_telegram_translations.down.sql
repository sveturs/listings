-- Rollback: Remove Telegram bot translations
DELETE FROM translations 
WHERE entity_type = 'telegram_bot' 
AND entity_id BETWEEN 100 AND 122;