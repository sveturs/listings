-- Удаляем функцию
DROP FUNCTION IF EXISTS get_telegram_translation(VARCHAR, VARCHAR);

-- Удаляем переводы Telegram бота
DELETE FROM translations WHERE entity_type = 'telegram_bot';