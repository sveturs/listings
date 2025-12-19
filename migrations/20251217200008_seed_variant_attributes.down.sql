-- Rollback seed variant attributes

DELETE FROM attribute_values WHERE attribute_id IN (3001, 3003, 3004);
DELETE FROM attributes WHERE id IN (3001, 3003, 3004);
