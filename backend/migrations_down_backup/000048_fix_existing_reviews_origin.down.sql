-- Откат: очищаем entity_origin поля
UPDATE reviews
SET entity_origin_type = NULL,
    entity_origin_id = NULL;