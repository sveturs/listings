-- Откат миграции - удаление добавленных переводов
DELETE FROM translations 
WHERE entity_type = 'marketplace_listing' 
AND entity_id IN (183, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268)
AND is_machine_translated = true;